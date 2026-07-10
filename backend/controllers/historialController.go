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

	var exist int
	err = conn.QueryRow(`SELECT COUNT(*) FROM equipos WHERE equipo_id = ?`, id).Scan(&exist)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error verificando equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}
	if exist == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "el equipo no existe"})
	}

	rows, err := conn.Query(`
    SELECT 
      h.historial_id,
      h.observaciones_tecnicas,
      h.fecha,
      e.nombre AS estado_nombre,
      COALESCE(u.nombres || ' ' || u.apellidos, '') AS nombres
    FROM historial_reparaciones h
    INNER JOIN estados_reparacion e ON h.estado_id = e.estado_id
    LEFT JOIN usuarios u ON h.usuario_id = u.usuario_id
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
		TecnicoNombre         string `json:"nombres"`
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
		UsuarioId             int    `json:"usuario_id"` // ID del usuario logueado en la sesión
		ObservacionesTecnicas string `json:"observaciones_tecnicas"`
	}

	if err := c.BodyParser(&req); err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	// 1. Obtener el Rol del Usuario
	var rolUsuario string
	err := conn.QueryRow(`
		SELECT r.nombre 
		FROM usuarios u
		INNER JOIN roles r ON u.rol_id = r.rol_id
		WHERE u.usuario_id = ? AND u.activo = 1`,
		req.UsuarioId,
	).Scan(&rolUsuario)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Usuario no autorizado o inexistente"})
		}
		_ = helpers.InsertLogsError(conn, "historial", "error consultando rol "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar permisos del usuario"})
	}

	// 2. Obtener el estado_id ACTUAL del equipo
	var estadoActualId int
	err = conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, req.EquipoId).Scan(&estadoActualId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo no existe"})
		}
		_ = helpers.InsertLogsError(conn, "historial", "error consultando estado actual "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el estado del equipo"})
	}

	// ==========================================
	//  RESTRICCIONES DE CONTROL DE ACCESO (ROLES)
	// ==========================================
	switch rolUsuario {
	case "TECNICO":
		// El técnico solo opera en estados de taller (1 al 5)
		if req.EstadoId >= 6 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Permiso denegado. Como Técnico no estás autorizado para Entregar (6) o Cancelar (7) equipos.",
			})
		}

	case "VENDEDOR":
		// El vendedor solo realiza transacciones de salida (6 o 7)
		if req.EstadoId < 6 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Permiso denegado. El perfil de Vendedor solo puede procesar Entregas (6) o Cancelaciones (7).",
			})
		}

	case "ADMINISTRADOR":
		// Acceso total a cualquier estado

	default:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Rol de usuario no reconocido"})
	}

	// ==========================================
	//  REGLAS DE NEGOCIO (ESTADOS Y SALDOS)
	// ==========================================

	// Evitar modificar equipos que ya fueron entregados o cancelados definitivamente
	if estadoActualId == 6 || estadoActualId == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se puede modificar el estado del equipo porque ya se encuentra en un estado final (Entregado/Cancelado).",
		})
	}

	// Evitar retrocesos al estado 'Recibido' (1)
	if estadoActualId > 1 && req.EstadoId == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Operación inválida. Un equipo en revisión no puede regresar al estado 'Recibido'.",
		})
	}

	// BLOQUEO: Para pasar a "Entregado (6)" se DEBE usar el otro endpoint específico
	if req.EstadoId == 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Para entregar el equipo debe usar el proceso de facturación/entregas para generar el acta obligatoria.",
		})
	}

	// Validación financiera obligatoria para pasar a CANCELADO (7) (Cobro de diagnóstico)
	if req.EstadoId == 7 {
		var saldo float64
		err = conn.QueryRow(`SELECT saldo FROM cuentas_reparacion WHERE equipo_id = ?`, req.EquipoId).Scan(&saldo)
		if err != nil {
			_ = helpers.InsertLogsError(conn, "historial", "error consultando saldo para cancelar "+err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "error al verificar saldo de cuenta"})
		}
		if saldo > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("No se puede cancelar el equipo. Registra un cobro de diagnóstico pendiente: $%.2f", saldo),
			})
		}
	}

	// Evitar el reenvío redundante del mismo estado actual
	if req.EstadoId == estadoActualId {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "El equipo ya se encuentra en el estado solicitado.",
		})
	}

	// 3. Ejecución en Transacción Atómica
	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error de base de datos"})
	}
	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	// Actualizar el estado principal del equipo
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
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "historial", "error actualizando equipos "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar estado del equipo"})
	}

	// Insertar el hito de cambio en 'historial_reparaciones'
	_, err = tx.Exec(`
		INSERT INTO historial_reparaciones (
			observaciones_tecnicas,
			fecha,                  
			usuario_id,
			equipo_id,
			estado_id
		) VALUES (?, ?, ?, ?, ?)`,
		strings.ToUpper(req.ObservacionesTecnicas),
		fechaActual,
		req.UsuarioId,
		req.EquipoId,
		req.EstadoId,
	)
	if err != nil {
		stringError := err.Error()
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "historial", "error insertando historial: "+stringError)
		return c.Status(500).JSON(fiber.Map{
			"message": "error al guardar el historial técnico",
			"debug":   stringError,
		})
	}

	// Confirmar cambios
	err = tx.Commit()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al procesar los cambios"})
	}

	// Registrar Log de Éxito
	_ = helpers.InsertLogs(conn, "UPDATE/INSERT", "historial_reparaciones", req.EquipoId, "Cambio de estado procesado por el rol: "+rolUsuario)

	return c.Status(200).JSON(fiber.Map{"message": "Estado del equipo actualizado correctamente"})
}
