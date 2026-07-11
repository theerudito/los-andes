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

func ConsultarEntregaPorEquipo(c *fiber.Ctx) error {

	var (
		conn    = database.GetDB()
		rows    *sql.Rows
		err     error
		entrega models.EntregaDTO
		id      int
		found   = false
	)

	id, err = strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	rows, err = conn.Query(`
    SELECT 
      en.entrega_id,
      en.fecha_entrega,
      en.trabajos_realizados,
      en.estado_final_equipo,
      en.conformidad_cliente,
      en.comprobante_nro,
      eq.codigo,
      (u.nombres || ' ' || u.apellidos) AS nombres
    FROM entregas en
    INNER JOIN equipos eq ON en.equipo_id = eq.equipo_id
    INNER JOIN usuarios u ON en.usuario_id = u.usuario_id
    WHERE en.equipo_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "Error al ejecutar la consulta "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&entrega.EntregaId,
			&entrega.FechaEntrega,
			&entrega.TrabajosRealizados,
			&entrega.EstadoFinalEquipo,
			&entrega.ConformidadCliente,
			&entrega.ComprobanteNro,
			&entrega.EquipoCodigo,
			&entrega.Usuario)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "usuarios", "Error al leer los registros "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true

	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(entrega)

}

func RegistrarEntrega(c *fiber.Ctx) error {

	var (
		entrega models.Entrega
		conn    = database.GetDB()
		err     error
		tx      *sql.Tx
		claims  *models.CustomClaims
		existe  int
		estado  int
		saldo   float64
	)

	if err := c.BodyParser(&entrega); err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	if claims.Rol == "TECNICO" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "solo usuario administrador o venderor puenden realizar esta accion"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM entregas WHERE equipo_id = ?`, entrega.EquipoId).Scan(&existe)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error verificando duplicidad de entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el historial de entregas"})
	}

	if existe > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ya existe una entrega registrada"})
	}

	err = conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, entrega.EquipoId).Scan(&existe)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo no existe"})
		}
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando estado del equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el equipo"})
	}

	if estado != 5 {
		if estado == 6 || estado == 7 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Este equipo ya fue procesado y cerrado anteriormente"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se puede registrar la entrega. El equipo debe estar en estado 'Listo para entrega' primero.",
		})
	}

	err = conn.QueryRow(`SELECT saldo FROM cuentas_reparacion WHERE equipo_id = ?`, entrega.EquipoId).Scan(&saldo)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando saldo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el saldo de la cuenta"})
	}

	if saldo > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("No se puede entregar el equipo. Registra un saldo pendiente de: $%.2f", saldo),
		})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error iniciando transacción"})
	}

	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	comprobanteNro, err := helpers.ObtenerCodigo(conn, "O")

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error obteniendo secuencial "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error generando número de comprobante"})
	}

	_, err = tx.Exec(`
		INSERT INTO entregas (
			fecha_entrega,
			trabajos_realizados,
			estado_final_equipo,
			conformidad_cliente,
			comprobante_nro,
			equipo_id,
			usuario_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		fechaActual,
		strings.ToUpper(entrega.TrabajosRealizados),
		strings.ToUpper(entrega.EstadoFinalEquipo),
		entrega.ConformidadCliente,
		comprobanteNro,
		entrega.EquipoId,
		claims.UserId,
	)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al registrar el acta de entrega"})
	}

	_, err = tx.Exec(`
		UPDATE equipos 
		SET estado_id = 6, 
		    fecha_modificacion = ? 
		WHERE equipo_id = ?`,
		fechaActual,
		entrega.EquipoId,
	)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error actualizando estado del equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar el estado final del equipo"})
	}

	_, err = tx.Exec(`
		INSERT INTO historial_reparaciones (
			observaciones_tecnicas,
			fecha, 
			usuario_id, 
			equipo_id,
			estado_id
		) VALUES (?, ?, ?, ?, 6)`,
		strings.ToUpper(entrega.Observaciones),
		fechaActual,
		claims.UserId,
		entrega.EquipoId,
	)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando historial de cierre "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al registrar el hito de cierre"})
	}

	err = helpers.InsertLogs(tx, "INSERT", "entregas", claims.Name, "registro creado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al confirmar la entrega"})
	}

	_ = helpers.ActualizarCodigo(conn, "O")

	return c.Status(201).JSON(fiber.Map{"message": "registro creado correctamente"})

}

func ProcesarEntregaEquipo(c *fiber.Ctx) error {
	var (
		conn    = database.GetDB()
		tx      *sql.Tx
		entrega models.EntregaEquipo
		claims  *models.CustomClaims
		err     error
		estado  int
		saldo   float64
	)

	if err := c.BodyParser(&entrega); err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	if claims.Rol == "TECNICO" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "solo usuario administrador o vendedor puenden realizar esta accion"})
	}

	err = conn.QueryRow(`
		SELECT e.estado_id, COALESCE(cr.saldo, 0.00) 
		FROM equipos e
		LEFT JOIN cuentas_reparacion cr ON e.equipo_id = cr.equipo_id
		WHERE e.equipo_id = ?`,
		entrega.EquipoId,
	).Scan(&estado, &saldo)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo especificado no existe"})
		}
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando estado/saldo: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar información del equipo"})
	}

	if estado != 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Solo se pueden entregar formalmente aquellos equipos que estén en estado 'Listo para entrega' (5)."})
	}

	if saldo > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("No se puede proceder con la entrega. Registra un saldo pendiente de: $%.2f", saldo)})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error iniciando transacción: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error interno del servidor"})
	}

	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	_, err = tx.Exec(`
		INSERT INTO entregas (
			fecha_entrega,
			trabajos_realizados,
			estado_final_equipo,
			conformidad_cliente,
			comprobante_nro,
			equipo_id,
			usuario_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		fechaActual,
		strings.ToUpper(entrega.TrabajosRealizados),
		strings.ToUpper(entrega.EstadoFinalEquipo),
		entrega.ConformidadCliente,
		entrega.ComprobanteNro,
		entrega.EquipoId,
		claims.UserId,
	)
	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error al guardar acta de entrega: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al guardar el acta física de entrega"})
	}

	_, err = tx.Exec(`
		UPDATE equipos 
		SET estado_id = 6, 
		    fecha_modificacion = ? 
		WHERE equipo_id = ?`,
		fechaActual,
		entrega.EquipoId,
	)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error actualizando estado a entregado: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar el estado del equipo"})
	}

	_, err = tx.Exec(`
		INSERT INTO historial_reparaciones (
			observaciones_tecnicas,
			fecha,                  
			usuario_id,
			equipo_id,
			estado_id
		) VALUES (?, ?, ?, ?, 6)`,
		strings.ToUpper(entrega.Observaciones),
		fechaActual,
		claims.UserId,
		entrega.EquipoId,
	)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error actualizando el registro: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error actualizando el registro"})
	}

	err = helpers.InsertLogs(tx, "UPDATE", "entregas", claims.Name, "registro actualizando correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro actualizando correctamente"})

}
