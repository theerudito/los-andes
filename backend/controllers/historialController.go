package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"los_andes/database"
	"los_andes/helpers"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ConsultarHistorialEquipo(c *fiber.Ctx) error {
	conn := database.GetDB()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	// Verificar primero si el equipo existe en la base de datos
	var exist int
	err = conn.QueryRow(`SELECT COUNT(*) FROM equipos WHERE equipo_id = ?`, id).Scan(&exist)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error verificando equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}
	if exist == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "el equipo no existe"})
	}

	// Consulta que une las tablas de estados y técnicos para devolver nombres legibles al frontend
	rows, err := conn.Query(`
    SELECT 
      h.historial_id,
      h.observaciones_tecnicas,
      h.fecha,
      e.nombre AS estado_nombre,
      COALESCE(t.nombres || ' ' || t.apellidos, 'SISTEMA / NO ASIGNADO') AS tecnico_nombre
    FROM historial_reparaciones h
    INNER JOIN estados_reparacion e ON h.estado_id = e.estado_id
    LEFT JOIN tecnicos t ON h.tecnico_id = t.tecnico_id
    WHERE h.equipo_id = ?
    ORDER BY h.fecha DESC`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error consultando historial "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al obtener el historial"})
	}
	defer rows.Close()

	type HistorialResponse struct {
		HistorialId           int    `json:"historial_id"`
		ObservacionesTecnicas string `json:"observaciones_tecnicas"`
		FechaCambio           string `json:"fecha_cambio"`
		EstadoNombre          string `json:"estado_nombre"`
		TecnicoNombre         string `json:"tecnico_nombre"`
	}

	var listaHistorial []HistorialResponse = []HistorialResponse{}

	for rows.Next() {
		var h HistorialResponse
		err := rows.Scan(&h.HistorialId, &h.ObservacionesTecnicas, &h.FechaCambio, &h.EstadoNombre, &h.TecnicoNombre)
		if err != nil {
			_ = helpers.InsertLogsError(conn, "historial", "error escaneando filas "+err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "error procesando los datos"})
		}
		listaHistorial = append(listaHistorial, h)
	}

	return c.Status(200).JSON(listaHistorial)
}

func ActualizarEstadoEquipo(c *fiber.Ctx) error {
	conn := database.GetDB()
	var tx *sql.Tx

	var req struct {
		EquipoId              int    `json:"equipo_id"`
		EstadoId              int    `json:"estado_id"`
		TecnicoId             *int   `json:"tecnico_id"`
		ObservacionesTecnicas string `json:"observaciones_tecnicas"`
	}

	if err := c.BodyParser(&req); err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	// 1. Obtener el estado_id ACTUAL del equipo
	var estadoActualId int
	err := conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, req.EquipoId).Scan(&estadoActualId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo no existe"})
		}
		_ = helpers.InsertLogsError(conn, "historial", "error consultando estado actual "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el estado del equipo"})
	}

	// 2. ESCUDO 1: Bloqueo si ya está en un estado final
	if estadoActualId == 6 || estadoActualId == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se puede cambiar el estado de este equipo porque ya se encuentra finalizado (Entregado/Cancelado)",
		})
	}

	// 3. ESCUDO 2: CONTROL DE FLUJO ASCENDENTE
	if req.EstadoId <= estadoActualId {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("Flujo de proceso inválido. No puedes cambiar a un estado menor o igual al actual."),
		})
	}

	// 4. ESCUDO 3: Validación para pasar a ENTREGADO
	if req.EstadoId == 6 {
		var saldo float64
		err = conn.QueryRow(`SELECT saldo FROM cuentas_reparacion WHERE equipo_id = ?`, req.EquipoId).Scan(&saldo)
		if err != nil {
			_ = helpers.InsertLogsError(conn, "historial", "error consultando saldo "+err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "error al verificar saldo financiero"})
		}
		if saldo > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("No se puede pasar a estado Entregado. La cuenta registra un saldo pendiente de: $%.2f", saldo),
			})
		}
	}

	// 5. Iniciar Transacción Atómica
	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error iniciando transacción"})
	}
	// Defer de seguridad por si el sistema llega a tirar un panic inesperado
	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	// 6. PASO A: Actualizar el estado_id en la tabla 'equipos'
	_, err = tx.Exec(`
    UPDATE equipos 
    SET estado_id = ?, 
        fecha_modificacion = ? 
    WHERE equipo_id = ?`,
		req.EstadoId,
		fechaActual,
		req.EquipoId,
	)
	if err != nil {
		_ = tx.Rollback() // <-- Liberamos la base de datos INMEDIATAMENTE
		_ = helpers.InsertLogsError(conn, "historial", "error actualizando tabla equipos "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar el estado del equipo"})
	}

	// 7. PASO B: Insertar el registro en 'historial_reparaciones'
	_, err = tx.Exec(`
    INSERT INTO historial_reparaciones (
      observaciones_tecnicas,
      fecha,                  
      tecnico_id,
      equipo_id,
      estado_id
    ) VALUES (?, ?, ?, ?, ?)`,
		strings.ToUpper(req.ObservacionesTecnicas),
		fechaActual,
		req.TecnicoId,
		req.EquipoId,
		req.EstadoId,
	)
	if err != nil {
		stringError := err.Error()
		_ = tx.Rollback() // <-- IMPORTANTE: Romper la transacción AQUÍ antes de llamar al log externo

		// Ahora que la BD está libre, el log se va a guardar sí o sí
		_ = helpers.InsertLogsError(conn, "historial", "error insertando hito de historial: "+stringError)

		return c.Status(500).JSON(fiber.Map{
			"message": "error al registrar en el historial",
			"debug":   stringError, // Te lo devuelvo en el JSON temporalmente para que veas qué restricción saltó
		})
	}

	// 8. Confirmar transacción
	err = tx.Commit()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error confirmando cambios"})
	}

	// 9. Registrar Log de Éxito
	_ = helpers.InsertLogs(conn, "UPDATE/INSERT", "historial_reparaciones", req.EquipoId, "Cambio de estado procesado correctamente")

	return c.Status(200).JSON(fiber.Map{"message": "Estado del equipo actualizado correctamente"})
}
