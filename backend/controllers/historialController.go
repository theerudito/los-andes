package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ConsultarHistorialEquipo(c *fiber.Ctx) error {
	var (
		conn        = database.GetDB()
		rows        *sql.Rows
		err         error
		exist       int
		historiales []models.HistorialReparacionesDTO
		historial   models.HistorialReparacionesDTO
	)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM equipos WHERE equipo_id = ?`, id).Scan(&exist)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error verificando equipo "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	if exist == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "el equipo no existe"})
	}

	rows, err = conn.Query(`
		SELECT
			h.historial_id,
			h.observaciones_tecnicas,
			h.fecha,
			h.equipo_id,
			COALESCE(eq.nombre, ''),
			h.estado_id,
			e.nombre,
			h.usuario_id,
			COALESCE(u.nombres, '') AS nombres,
			COALESCE(u.apellidos, '') AS apellidos 
		FROM historial_reparaciones h
		INNER JOIN equipos eq ON h.equipo_id = eq.equipo_id
		INNER JOIN estados_reparacion e ON h.estado_id = e.estado_id
		LEFT JOIN usuarios u ON h.usuario_id = u.usuario_id
		WHERE h.equipo_id = ?
		ORDER BY h.fecha DESC`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error consultando historial "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al obtener el historial"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&historial.HistorialId,
			&historial.ObservacionesTecnicas,
			&historial.FechaCambio,
			&historial.EquipoId,
			&historial.Equipo,
			&historial.EstadoId,
			&historial.Estado,
			&historial.UsuarioId,
			&historial.Nombres,
			&historial.Apellidos,
		)
		if err != nil {
			_ = helpers.InsertLogsError(conn, "historial", "error escaneando filas "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error procesando los datos"})
		}
		historiales = append(historiales, historial)
	}

	if err = rows.Err(); err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error recorriendo filas "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error procesando los datos"})
	}

	return c.Status(200).JSON(historiales)
}

func ActualizarEstadoEquipo(c *fiber.Ctx) error {
	var (
		conn      = database.GetDB()
		err       error
		estado    int
		claims    *models.CustomClaims
		saldo     float64
		tx        *sql.Tx
		historial models.HistorialReparaciones
	)

	if err := c.BodyParser(&historial); err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	err = conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, historial.EquipoId).Scan(&estado)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo no existe"})
		}
		_ = helpers.InsertLogsError(conn, "historial", "error consultando estado actual "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el estado del equipo"})
	}

	switch claims.Rol {
	case "TECNICO":
		if historial.EstadoId >= 6 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. Como Técnico no estás autorizado para Entregar o Cancelar equipos."})
		}
	case "VENDEDOR":
		if historial.EstadoId < 6 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. El perfil de Vendedor solo puede procesar Entregas o Cancelaciones."})
		}
	default:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Rol de usuario no reconocido"})
	}

	if estado == 6 || estado == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "No se puede modificar el estado del equipo porque ya se encuentra en un estado final (Entregado/Cancelado)."})
	}

	if estado > 1 && historial.EstadoId == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Operación inválida. Un equipo en revisión no puede regresar al estado 'Recibido'."})
	}

	if historial.EstadoId == 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Para entregar el equipo debe usar el proceso de facturación/entregas para generar el acta obligatoria."})
	}

	if historial.EstadoId == 7 {
		var err = conn.QueryRow(`SELECT saldo FROM cuentas_reparacion WHERE equipo_id = ?`, historial.EquipoId).Scan(&saldo)
		if err != nil {
			_ = helpers.InsertLogsError(conn, "historial", "error consultando saldo para cancelar "+err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "error al verificar saldo de cuenta"})
		}

		if saldo > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("No se puede cancelar el equipo. Registra un cobro de diagnóstico pendiente: $%.2f", saldo)})
		}
	}

	if historial.EstadoId == estado {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "El equipo ya se encuentra en el estado solicitado.",
		})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error de base de datos"})
	}

	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	_, err = tx.Exec(`
		UPDATE equipos 
		SET estado_id = ?, 
		    fecha_modificacion = ? 
		WHERE equipo_id = ?`,
		historial.EstadoId,
		fechaActual,
		historial.EquipoId,
	)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "historial", "error actualizando equipos "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar estado del equipo"})
	}

	_, err = tx.Exec(`
		INSERT INTO historial_reparaciones (
			observaciones_tecnicas,
			fecha,                  
			usuario_id,
			equipo_id,
			estado_id
		) VALUES (?, ?, ?, ?, ?)`,
		strings.ToUpper(historial.ObservacionesTecnicas),
		fechaActual,
		claims.UserId,
		historial.EquipoId,
		historial.EstadoId)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "historial", "error insertando historial: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al guardar el historial técnico"})
	}

	err = helpers.InsertLogs(conn, "UPDATE", "historial_reparaciones", claims.Name, "registro actualizado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al procesar los cambios"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro actualizado correctamente"})

}
