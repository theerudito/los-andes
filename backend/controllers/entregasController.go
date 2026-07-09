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

func ConsultarEntregaPorEquipo(c *fiber.Ctx) error {
	conn := database.GetDB()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	type EntregaResponse struct {
		EntregaId          int    `json:"entrega_id"`
		FechaEntrega       string `json:"fecha_entrega"`
		TrabajosRealizados string `json:"trabajos_realizados"`
		EstadoFinalEquipo  string `json:"estado_final_equipo"`
		ConformidadCliente int    `json:"conformidad_cliente"`
		ComprobanteNro     string `json:"comprobante_nro"`
		EquipoCodigo       string `json:"equipo_codigo"`
		TecnicoNombre      string `json:"tecnico_nombre"`
	}

	var ent EntregaResponse

	err = conn.QueryRow(`
    SELECT 
      en.entrega_id,
      en.fecha_entrega,
      en.trabajos_realizados,
      en.estado_final_equipo,
      en.conformidad_cliente,
      en.comprobante_nro,
      eq.codigo,
      (t.nombres || ' ' || t.apellidos) AS tecnico_nombre
    FROM entregas en
    INNER JOIN equipos eq ON en.equipo_id = eq.equipo_id
    INNER JOIN tecnicos t ON en.tecnico_id = t.tecnico_id
    WHERE en.equipo_id = ?`, id).Scan(
		&ent.EntregaId,
		&ent.FechaEntrega,
		&ent.TrabajosRealizados,
		&ent.EstadoFinalEquipo,
		&ent.ConformidadCliente,
		&ent.ComprobanteNro,
		&ent.EquipoCodigo,
		&ent.TecnicoNombre,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Este equipo aún no cuenta con un registro de entrega"})
		}
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al consultar los datos de la entrega"})
	}

	return c.Status(200).JSON(ent)
}

func RegistrarEntrega(c *fiber.Ctx) error {
	conn := database.GetDB()
	var tx *sql.Tx

	var req struct {
		EquipoId           int    `json:"equipo_id"`
		TecnicoId          int    `json:"tecnico_id"`
		TrabajosRealizados string `json:"trabajos_realizados"`
		EstadoFinalEquipo  string `json:"estado_final_equipo"`
		ConformidadCliente int    `json:"conformidad_cliente"`
	}

	if err := c.BodyParser(&req); err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	// 1. Obtener el estado actual del equipo antes de liquidar
	var estadoActualId int
	err := conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, req.EquipoId).Scan(&estadoActualId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo no existe"})
		}
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando estado del equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el equipo"})
	}

	// =========================================================================
	// NUEVO CANDADO DE FLUJO: Solo se permite entregar si está "Listo para entrega" (5)
	// =========================================================================
	if estadoActualId != 5 {
		// Si ya es 6 o 7, avisamos que ya fue cerrado
		if estadoActualId == 6 || estadoActualId == 7 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Este equipo ya fue procesado y cerrado anteriormente"})
		}
		// Si es del 1 al 4, bloqueamos porque no ha terminado el taller
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se puede registrar la entrega. El equipo debe estar en estado 'Listo para entrega' primero.",
		})
	}
	// =========================================================================

	// 2. REGLA DE NEGOCIO FINANCIERA: Verificar que el saldo de la cuenta sea 0.00
	var saldo float64
	err = conn.QueryRow(`SELECT saldo FROM cuentas_reparacion WHERE equipo_id = ?`, req.EquipoId).Scan(&saldo)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando saldo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el saldo de la cuenta"})
	}
	if saldo > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("No se puede entregar el equipo. Registra un saldo pendiente de: $%.2f", saldo),
		})
	}

	// 3. Iniciar Transacción Atómica
	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error iniciando transacción"})
	}
	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	// 4. Obtener secuencial para el número de comprobante de entrega
	comprobanteNro, err := helpers.ObtenerCodigo(conn, "O")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error obteniendo secuencial "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error generando número de comprobante"})
	}

	// 5. PASO A: Insertar en la tabla 'entregas'
	_, err = tx.Exec(`
    INSERT INTO entregas (
      fecha_entrega,
      trabajos_realizados,
      estado_final_equipo,
      conformidad_cliente,
      comprobante_nro,
      equipo_id,
      tecnico_id
    ) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		fechaActual,
		strings.ToUpper(req.TrabajosRealizados),
		strings.ToUpper(req.EstadoFinalEquipo),
		req.ConformidadCliente,
		comprobanteNro,
		req.EquipoId,
		req.TecnicoId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al registrar el acta de entrega"})
	}

	// 6. PASO B: Actualizar el estado global del equipo a 'Entregado' (ID 6)
	_, err = tx.Exec(`
    UPDATE equipos 
    SET estado_id = 6, 
        fecha_modificacion = ? 
    WHERE equipo_id = ?`,
		fechaActual,
		req.EquipoId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error actualizando estado del equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar el estado final del equipo"})
	}

	// 7. PASO C: Registrar el hito de cierre en el historial_reparaciones
	_, err = tx.Exec(`
    INSERT INTO historial_reparaciones (
      observaciones_tecnicas,
      fecha_cambio,
      tecnico_id,
      equipo_id,
      estado_id
    ) VALUES (?, ?, ?, ?, 6)`, // 6 = Entregado
		"EQUIPO ENTREGADO FORMALMENTE AL CLIENTE. COMPROBANTE NRO: "+comprobanteNro,
		fechaActual,
		req.TecnicoId,
		req.EquipoId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando historial de cierre "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al registrar el hito de cierre"})
	}

	// 8. Confirmar todos los cambios en la base de datos
	err = tx.Commit()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al confirmar la entrega"})
	}

	// 9. Incrementar el secuencial del comprobante de entrega
	_ = helpers.ActualizarCodigo(conn, "O")

	// 10. Log de Auditoría OK
	_ = helpers.InsertLogs(conn, "INSERT", "entregas", req.EquipoId, "Entrega y acta de salida generada correctamente")

	return c.Status(201).JSON(fiber.Map{
		"message":         "Entrega procesada y registrada con éxito",
		"comprobante_nro": comprobanteNro,
	})
}
