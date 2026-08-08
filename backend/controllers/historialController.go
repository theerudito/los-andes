package controllers

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/phpdave11/gofpdf"
)

func ConsultarHistorialEquipo(c *fiber.Ctx) error {
	var (
		conn        = database.GetDB()
		rows        *sql.Rows
		err         error
		historiales []models.HistorialDTO
		historial   models.HistorialDTO
	)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	rows, err = conn.Query(`
		SELECT
			h.historial_id,
			h.observaciones_tecnicas,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', h.fecha), '') AS fecha,
			e.equipo_id,
			e.tipo_equipo,
			e.numero_serie,
			e.codigo,
			r.nombre AS estado,
			r.estado_id,
			u.usuario_id,
			u.nombres AS usuario_nombres,
			u.apellidos AS usuario_apellidos,
			c.cliente_id,
			c.nombres AS cliente_nombres,
			c.apellidos AS cliente_apellidos
		FROM historial_reparaciones h
			INNER JOIN equipos e ON h.equipo_id = e.equipo_id
			INNER JOIN clientes c ON e.cliente_id = c.cliente_id
			INNER JOIN estados_reparacion r ON h.estado_id = r.estado_id
			INNER JOIN usuarios u ON h.usuario_id = u.usuario_id
		WHERE e.equipo_id = ?
		ORDER BY h.historial_id DESC`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error al obtener el historial "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al obtener el historial"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&historial.HistorialId,
			&historial.ObservacionesTecnicas,
			&historial.Fecha,
			&historial.EquipoId,
			&historial.Equipo,
			&historial.Serie,
			&historial.Codigo,
			&historial.Estado,
			&historial.EstadoId,
			&historial.UsuarioId,
			&historial.Nombres_Usuario,
			&historial.Apellidos_Usuario,
			&historial.ClienteId,
			&historial.Nombres_Cliente,
			&historial.Apellidos_Cliente)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "historial", "error escaneando filas "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error procesando los datos"})
		}

		historiales = append(historiales, historial)

	}

	if len(historiales) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.Status(200).JSON(historiales)
}

func ConsultarHistorial(c *fiber.Ctx) error {
	var (
		conn      = database.GetDB()
		rows      *sql.Rows
		err       error
		exist     int
		historial models.Historial
	)

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID del historial inválido"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM historial_reparaciones WHERE historial_id = ?`, id).Scan(&exist)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error verificando equipo "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	rows, err = conn.Query(`
		SELECT
			h.historial_id,
			h.observaciones_tecnicas,
			e.equipo_id,
			r.estado_id
		FROM historial_reparaciones h
			INNER JOIN equipos e ON h.equipo_id = e.equipo_id
			INNER JOIN clientes c ON e.cliente_id = c.cliente_id
			INNER JOIN estados_reparacion r ON h.estado_id = r.estado_id
		WHERE h.historial_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error al obtener el historial "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al obtener el historial"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&historial.HistorialId,
			&historial.ObservacionesTecnicas,
			&historial.EquipoId,
			&historial.EstadoId)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "clientes", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

	}

	if exist == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(historial)
}

func CrearEstadoEquipo(c *fiber.Ctx) error {
	var (
		conn      = database.GetDB()
		err       error
		estado    int
		claims    *models.CustomClaims
		saldo     float64
		tx        *sql.Tx
		historial models.Historial
		existe    int
	)

	if err := c.BodyParser(&historial); err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	if historial.EstadoId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Debe seleccionar un estado válido para el equipo."})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error al leer los claims "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los claims"})
	}

	err = conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, historial.EquipoId).Scan(&estado)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo no existe"})
		}
		_ = helpers.InsertLogsError(conn, "historial", "error consultando estado actual "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al verificar el estado del equipo"})
	}

	switch claims.Rol {
	case "ADMINISTRADOR":
		if historial.EstadoId < 2 || historial.EstadoId > 7 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. Rango de estados no válido."})
		}
	case "TECNICO":
		if historial.EstadoId < 2 || historial.EstadoId > 5 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. Como Técnico solo puedes gestionar estados entre 'En diagnóstico' (2) y 'Listo para entregar' (5)."})
		}
	case "VENDEDOR":
		if historial.EstadoId != 7 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. El perfil de Vendedor solo está autorizado para registrar Cancelaciones (7)."})
		}
	default:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. Tu rol no tiene autorización para registrar estados en el historial."})
	}

	if estado == 6 || estado == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "No se puede modificar el estado del equipo porque ya se encuentra en un estado final (Entregado/Cancelado)."})
	}

	if historial.EstadoId != 7 {
		if estado == 1 && historial.EstadoId != 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Secuencia inválida. El equipo está en 'Recibido' (1) y debe pasar obligatoriamente a 'En diagnóstico' (2)."})
		}
		if estado > 1 && historial.EstadoId != estado+1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Secuencia inválida. El estado actual es %d y debe avanzar al estado %d.", estado, estado+1)})
		}
	}

	if historial.EstadoId == 7 {
		err = conn.QueryRow(`SELECT COALESCE(saldo, 0) FROM cuentas_reparacion WHERE equipo_id = ?`, historial.EquipoId).Scan(&saldo)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			_ = helpers.InsertLogsError(conn, "cuentas", "error consultando saldo para cancelar "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al verificar saldo de cuenta"})
		}

		if saldo > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("No se puede cancelar el equipo. Registra un cobro de diagnóstico pendiente: $%.2f", saldo)})
		}
	}

	err = conn.QueryRow(`
		SELECT COUNT(1) 
		FROM historial_reparaciones 
		WHERE equipo_id = ? AND estado_id = ?`,
		historial.EquipoId, historial.EstadoId,
	).Scan(&existe)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error verificando duplicidad de estado "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error de verificación de historial"})
	}

	if existe > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Este estado ya fue registrado anteriormente para este equipo."})
	}

	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error iniciando transacción "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error de base de datos"})
	}
	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	_, err = tx.Exec(`
		UPDATE equipos 
		SET estado_id = ?, fecha_modificacion = ? 
		WHERE equipo_id = ?`,
		historial.EstadoId, fechaActual, historial.EquipoId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error actualizando equipos "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al actualizar estado del equipo"})
	}

	_, err = tx.Exec(`
		INSERT INTO historial_reparaciones (
			observaciones_tecnicas, fecha, usuario_id, equipo_id, estado_id
		) VALUES (?, ?, ?, ?, ?)`,
		strings.ToUpper(historial.ObservacionesTecnicas), fechaActual, claims.UserId, historial.EquipoId, historial.EstadoId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error insertando historial: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al guardar el historial técnico"})
	}

	if err = tx.Commit(); err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error confirmando transacción "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al procesar los cambios"})
	}

	_ = helpers.InsertLogs(conn, "INSERT", "historial", claims.Name, "nuevo estado registrado correctamente")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "registro creado correctamente"})
}

func ActualizarEstadoEquipo(c *fiber.Ctx) error {
	var (
		conn            = database.GetDB()
		err             error
		estadoActual    int
		claims          *models.CustomClaims
		tx              *sql.Tx
		historial       models.Historial
		existe          int
		maxEstado       int
		estadoAnterior  sql.NullInt64
		estadoSiguiente sql.NullInt64
	)

	if err := c.BodyParser(&historial); err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error al leer los claims "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los claims"})
	}

	if historial.HistorialId <= 0 || historial.EstadoId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Parámetros inválidos."})
	}

	err = conn.QueryRow(`
		SELECT equipo_id 
		FROM historial_reparaciones 
		WHERE historial_id = ?`, historial.HistorialId).Scan(&historial.EquipoId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El registro de historial no existe"})
		}
		_ = helpers.InsertLogsError(conn, "historial", "error consultando el registro de historial "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al verificar el historial"})
	}

	err = conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, historial.EquipoId).Scan(&estadoActual)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error consultando estado actual del equipo "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al verificar el estado del equipo"})
	}

	switch claims.Rol {
	case "ADMINISTRADOR":
		if historial.EstadoId < 2 || historial.EstadoId > 7 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. Rango de estados no válido."})
		}
	case "TECNICO":
		if historial.EstadoId < 2 || historial.EstadoId > 5 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. Como Técnico solo puedes gestionar estados entre 'En diagnóstico' (2) y 'Listo para entregar' (5)."})
		}
	case "VENDEDOR":
		if historial.EstadoId != 7 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. El perfil de Vendedor solo está autorizado para modificar estados hacia Cancelado (7)."})
		}
	default:
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. Tu rol no tiene autorización para modificar el historial del equipo."})
	}

	if estadoActual == 6 || estadoActual == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "No se puede modificar el historial porque el equipo ya se encuentra en un estado final (Entregado/Cancelado)."})
	}

	if estadoActual == 1 || historial.EstadoId == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Operación inválida. No se permite modificar el estado inicial (Recibido)."})
	}

	_ = conn.QueryRow(`
		SELECT estado_id 
		FROM historial_reparaciones 
		WHERE equipo_id = ? AND historial_id < ? 
		ORDER BY historial_id DESC LIMIT 1`,
		historial.EquipoId, historial.HistorialId,
	).Scan(&estadoAnterior)

	_ = conn.QueryRow(`
		SELECT estado_id 
		FROM historial_reparaciones 
		WHERE equipo_id = ? AND historial_id > ? 
		ORDER BY historial_id ASC LIMIT 1`,
		historial.EquipoId, historial.HistorialId,
	).Scan(&estadoSiguiente)

	if historial.EstadoId != 7 {
		if estadoAnterior.Valid {
			if estadoAnterior.Int64 == 1 && historial.EstadoId != 2 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Secuencia inválida. El primer paso después de 'Recibido' (1) debe ser 'En diagnóstico' (2)."})
			}
			if estadoAnterior.Int64 > 1 && historial.EstadoId != int(estadoAnterior.Int64)+1 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": fmt.Sprintf("Secuencia inválida. El nuevo estado debe ser el correlativo %d.", estadoAnterior.Int64+1),
				})
			}
		} else if historial.EstadoId != 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Secuencia inválida. El primer estado debe ser 'En diagnóstico' (2)."})
		}
	}

	if estadoSiguiente.Valid {
		if historial.EstadoId >= int(estadoSiguiente.Int64) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("Secuencia inválida. El estado debe ser menor al siguiente registro registrado (%d).", estadoSiguiente.Int64),
			})
		}
	}

	err = conn.QueryRow(`
		SELECT COUNT(1) 
		FROM historial_reparaciones 
		WHERE equipo_id = ? AND estado_id = ? AND historial_id != ?`,
		historial.EquipoId, historial.EstadoId, historial.HistorialId,
	).Scan(&existe)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error verificando duplicidad de estado "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error de verificación de historial"})
	}

	if existe > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Este estado ya existe en otro registro del historial de este equipo. Elija un estado diferente."})
	}

	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error iniciando transacción "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error de base de datos"})
	}
	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	_, err = tx.Exec(`
		UPDATE historial_reparaciones 
		SET observaciones_tecnicas = ?, estado_id = ?, fecha = ?, usuario_id = ? 
		WHERE historial_id = ?`,
		strings.ToUpper(historial.ObservacionesTecnicas), historial.EstadoId, fechaActual, claims.UserId, historial.HistorialId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error actualizando registro de historial "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al actualizar el historial técnico"})
	}

	err = tx.QueryRow(`
		SELECT COALESCE(MAX(estado_id), 1) 
		FROM historial_reparaciones 
		WHERE equipo_id = ?`, historial.EquipoId).Scan(&maxEstado)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error consultando max estado "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error recalculando el estado del equipo"})
	}

	_, err = tx.Exec(`
		UPDATE equipos 
		SET estado_id = ?, fecha_modificacion = ? 
		WHERE equipo_id = ?`,
		maxEstado, fechaActual, historial.EquipoId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error actualizando estado del equipo "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al actualizar estado del equipo"})
	}

	if err = tx.Commit(); err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error confirmando transacción "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al procesar los cambios"})
	}

	_ = helpers.InsertLogs(conn, "UPDATE", "historial", claims.Name, "registro de historial actualizado correctamente")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "registro actualizado correctamente"})
}

func ReporteHistorial(c *fiber.Ctx) error {
	var (
		conn        = database.GetDB()
		historiales []models.HistorialDTO
		rows        *sql.Rows
		err         error
	)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de equipo inválido"})
	}

	rows, err = conn.Query(`
		SELECT
			h.historial_id,
			COALESCE(h.observaciones_tecnicas, '') AS observaciones_tecnicas,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', h.fecha), '') AS fecha,
			e.equipo_id,
			e.tipo_equipo,
			COALESCE(e.numero_serie, '') AS numero_serie,
			r.nombre AS estado,
			r.estado_id,
			COALESCE(u.usuario_id, 0) AS usuario_id,
			COALESCE(u.nombres, 'N/A') AS usuario_nombres,
			COALESCE(u.apellidos, '') AS usuario_apellidos,
			c.cliente_id,
			c.nombres AS cliente_nombres,
			c.apellidos AS cliente_apellidos
		FROM historial_reparaciones h
			INNER JOIN equipos e ON h.equipo_id = e.equipo_id
			INNER JOIN clientes c ON e.cliente_id = c.cliente_id
			INNER JOIN estados_reparacion r ON h.estado_id = r.estado_id
			LEFT JOIN usuarios u ON h.usuario_id = u.usuario_id
		WHERE e.equipo_id = ?
		ORDER BY h.historial_id ASC`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "error al obtener el historial "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error al obtener el historial"})
	}
	defer rows.Close()

	for rows.Next() {
		var h models.HistorialDTO
		err = rows.Scan(
			&h.HistorialId,
			&h.ObservacionesTecnicas,
			&h.Fecha,
			&h.EquipoId,
			&h.Equipo,
			&h.Serie,
			&h.Estado,
			&h.EstadoId,
			&h.UsuarioId,
			&h.Nombres_Usuario,
			&h.Apellidos_Usuario,
			&h.ClienteId,
			&h.Nombres_Cliente,
			&h.Apellidos_Cliente,
		)
		if err != nil {
			_ = helpers.InsertLogsError(conn, "historial", "Error al leer los registros: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}
		historiales = append(historiales, h)
	}

	if len(historiales) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros de historial para este equipo"})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	primerRegistro := historiales[0]
	ultimoRegistro := historiales[len(historiales)-1]

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(110, 7, tr("HISTORIAL TÉCNICO DE EQUIPO"), "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(70, 7, tr(fmt.Sprintf("EQUIPO ID: %d", primerRegistro.EquipoId)), "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 5, tr("Sistema de Gestión de Mantenimiento de Computadoras"))
	pdf.Ln(8)

	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(180, 6, tr(" 1. INFORMACIÓN DEL CLIENTE"), "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	nombreCliente := fmt.Sprintf("%s %s", primerRegistro.Nombres_Cliente, primerRegistro.Apellidos_Cliente)
	pdf.CellFormat(30, 6, tr("Cliente:"), "L", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, tr(helpers.Limitar(nombreCliente, 60)), "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, tr("ID Cliente:"), "L,B", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, tr(fmt.Sprintf("%d", primerRegistro.ClienteId)), "R,B", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, tr(" 2. DETALLES DEL EQUIPO"), "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(30, 6, tr("Tipo Equipo:"), "L", 0, "L", false, 0, "")
	pdf.CellFormat(55, 6, tr(primerRegistro.Equipo), "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, tr("Número de Serie:"), "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, tr(primerRegistro.Serie), "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, tr("Estado Actual:"), "L,B", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, tr(ultimoRegistro.Estado), "R,B", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, tr(" 3. REGISTRO CONTINUO DE EVENTOS Y MANTENIMIENTO"), "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(35, 6, tr("Fecha / Hora"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 6, tr("Estado"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(70, 6, tr("Observaciones Técnicas"), "1", 0, "L", false, 0, "")
	pdf.CellFormat(40, 6, tr("Técnico"), "1", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 8.5)

	for _, item := range historiales {
		tecnico := fmt.Sprintf("%s %s", item.Nombres_Usuario, item.Apellidos_Usuario)
		if strings.TrimSpace(tecnico) == "" || item.UsuarioId == 0 {
			tecnico = "Sistema"
		}

		obs := strings.TrimSpace(item.ObservacionesTecnicas)
		if obs == "" {
			obs = "Sin observaciones."
		}

		if pdf.GetY() > 260 {
			pdf.AddPage()
			pdf.SetFont("Arial", "B", 9)
			pdf.CellFormat(35, 6, tr("Fecha / Hora"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(35, 6, tr("Estado"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(70, 6, tr("Observaciones Técnicas"), "1", 0, "L", false, 0, "")
			pdf.CellFormat(40, 6, tr("Técnico"), "1", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 8.5)
		}

		yInicial := pdf.GetY()

		pdf.CellFormat(35, 6, tr(item.Fecha), "L,B", 0, "C", false, 0, "")
		pdf.CellFormat(35, 6, tr(helpers.Limitar(item.Estado, 18)), "L,B", 0, "C", false, 0, "")

		currX := pdf.GetX()
		pdf.SetXY(currX, yInicial)
		pdf.MultiCell(70, 6, tr(obs), "L,B", "L", false)

		yFinal := pdf.GetY()
		if yFinal == yInicial {
			yFinal = yInicial + 6
		}

		pdf.SetXY(currX+70, yInicial)
		pdf.CellFormat(40, yFinal-yInicial, tr(helpers.Limitar(tecnico, 22)), "L,R,B", 1, "L", false, 0, "")

		pdf.SetY(yFinal)
	}

	pdf.Ln(10)

	pdf.SetFont("Arial", "I", 8)
	pdf.MultiCell(180, 4, tr("Nota de Trazabilidad: Este documento contiene la traza oficial de las revisiones, diagnósticos y mantenimientos ejecutados sobre el equipo en mención."), "", "C", false)

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "historial", "Error al procesar PDF de Historial: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar el reporte de historial"})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="historial_equipo_%d.pdf"`, primerRegistro.EquipoId))
	return c.Send(buf.Bytes())
}
