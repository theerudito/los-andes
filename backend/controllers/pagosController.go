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

func ConsultarCuentaEquipo(c *fiber.Ctx) error {
	conn := database.GetDB()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	type CuentaResponse struct {
		CuentaId     int     `json:"cuenta_id"`
		EquipoId     int     `json:"equipo_id"`
		EquipoCodigo string  `json:"equipo_codigo"`
		CostoTotal   float64 `json:"costo_total"`
		Abono        float64 `json:"abono"`
		Saldo        float64 `json:"saldo"`
	}

	var cuenta CuentaResponse

	// Consulta con INNER JOIN para traer también el código del equipo por estética en el frontend
	err = conn.QueryRow(`
    SELECT 
      c.cuenta_id,
      c.equipo_id,
      e.codigo,
      c.costo_total,
      c.abono,
      c.saldo
    FROM cuentas_reparacion c
    INNER JOIN equipos e ON c.equipo_id = e.equipo_id
    WHERE c.equipo_id = ?`, id).Scan(
		&cuenta.CuentaId,
		&cuenta.EquipoId,
		&cuenta.EquipoCodigo,
		&cuenta.CostoTotal,
		&cuenta.Abono,
		&cuenta.Saldo,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No se encontró una cuenta para este equipo"})
		}
		_ = helpers.InsertLogsError(conn, "cuentas", "error consultando cuenta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al obtener los datos de la cuenta"})
	}

	return c.Status(200).JSON(cuenta)
}

func ActualizarCuentaEquipo(c *fiber.Ctx) error {
	conn := database.GetDB()

	var req struct {
		EquipoId     int     `json:"equipo_id"`
		NuevoCosto   float64 `json:"costo_total"` // Costo de la reparación o diagnóstico
		MontoAAbonar float64 `json:"abono"`       // El abono que entrega en esta transacción
		UsuarioId    int     `json:"usuario_id"`  // ID del usuario logueado
	}

	if err := c.BodyParser(&req); err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	// 1. Obtener Rol del Usuario
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
		_ = helpers.InsertLogsError(conn, "cuentas", "error consultando rol: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar permisos"})
	}

	// 2. ESCUDO DE ROLES: El Técnico tiene terminantemente prohibido tocar dinero
	if rolUsuario == "TECNICO" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Permiso denegado. Su rol de Técnico no le permite registrar transacciones ni abonos.",
		})
	}

	// 3. Obtener Estado del Equipo
	var estadoId int
	err = conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, req.EquipoId).Scan(&estadoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo especificado no existe"})
		}
		_ = helpers.InsertLogsError(conn, "cuentas", "error consultando estado de equipo: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar estado del equipo"})
	}

	// 4. ESCUDO DE ESTADO FINAL: Nadie puede mover caja de equipos en estado 6 (Entregado) o 7 (Cancelado)
	if estadoId == 6 || estadoId == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Operación bloqueada. No se pueden registrar pagos en equipos finalizados (Entregados/Cancelados).",
		})
	}

	// 5. Obtener Cuenta actual en base de datos
	var costoActual, abonoActual float64
	err = conn.QueryRow(`
		SELECT costo_total, abono 
		FROM cuentas_reparacion 
		WHERE equipo_id = ?`,
		req.EquipoId,
	).Scan(&costoActual, &abonoActual)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No se encontró registro de cuenta para este equipo"})
		}
		_ = helpers.InsertLogsError(conn, "cuentas", "error consultando cuenta actual: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al leer el estado de la cuenta"})
	}

	var costoFinal float64
	var abonoFinal float64

	// Caso A: La cuenta está totalmente en cero (Costo inicial o diagnóstico inicial)
	if costoActual == 0 && abonoActual == 0 {
		if req.NuevoCosto <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Debe definir un costo total inicial mayor a $0 para poder guardar la cuenta.",
			})
		}
		if req.MontoAAbonar > req.NuevoCosto {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "El abono inicial no puede ser superior al costo total establecido.",
			})
		}
		costoFinal = req.NuevoCosto
		abonoFinal = req.MontoAAbonar

		// Caso B: Ya existe un costo fijado en la base de datos
	} else {
		// Control estricto: Solo el Administrador puede redefinir el costo total una vez guardado
		if req.NuevoCosto != 0 && req.NuevoCosto != costoActual {
			if rolUsuario != "ADMINISTRADOR" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"message": "Permiso denegado. Solo el Administrador puede modificar el costo total una vez establecido.",
				})
			}
			costoFinal = req.NuevoCosto
		} else {
			costoFinal = costoActual
		}

		// Sumar acumulativamente el nuevo abono al abono anterior
		abonoFinal = abonoActual + req.MontoAAbonar

		// Validar que el nuevo abono sumado no supere el costo total
		if abonoFinal > costoFinal {
			saldoRestante := costoFinal - abonoActual
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fmt.Sprintf("El abono ingresado supera el saldo pendiente. Saldo restante: $%.2f", saldoRestante),
			})
		}
	}

	// 6. Actualizar en Base de Datos (SQLite recalcula la columna autogenerada 'saldo')
	_, err = conn.Exec(`
		UPDATE cuentas_reparacion 
		SET costo_total = ?, 
		    abono = ? 
		WHERE equipo_id = ?`,
		costoFinal,
		abonoFinal,
		req.EquipoId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "error actualizando valores de cuenta: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al procesar la transacción"})
	}

	// 7. Auditoría de Caja
	mensajeAuditoria := fmt.Sprintf("Abono de $%.2f procesado por %s. Total Abono: $%.2f / Costo: $%.2f",
		req.MontoAAbonar, rolUsuario, abonoFinal, costoFinal)
	_ = helpers.InsertLogs(conn, "UPDATE", "cuentas_reparacion", req.EquipoId, mensajeAuditoria)

	saldoFinal := costoFinal - abonoFinal
	return c.Status(200).JSON(fiber.Map{
		"message":      "Transacción de caja registrada con éxito",
		"costo_total":  costoFinal,
		"total_abono":  abonoFinal,
		"saldo_actual": saldoFinal,
	})
}

func ProcesarEntregaEquipo(c *fiber.Ctx) error {
	conn := database.GetDB()
	var tx *sql.Tx

	var req struct {
		EquipoId           int    `json:"equipo_id"`
		UsuarioId          int    `json:"usuario_id"` // Vendedor o Administrador
		TrabajosRealizados string `json:"trabajos_realizados"`
		EstadoFinalEquipo  string `json:"estado_final_equipo"` // Ej: "REPARADO", "SIN SOLUCION"
		ConformidadCliente int    `json:"conformidad_cliente"` // 1: Conforme, 0: No conforme
		ComprobanteNro     string `json:"comprobante_nro"`     // Número de factura o recibo físico
	}

	if err := c.BodyParser(&req); err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	// 1. Obtener Rol del Usuario
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
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando rol: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar permisos"})
	}

	// 2. ESCUDO DE ROLES: Los técnicos no entregan equipos ni firman actas de conformidad
	if rolUsuario == "TECNICO" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Permiso denegado. Los Técnicos no pueden procesar entregas físicas ni firmar actas.",
		})
	}

	// 3. Obtener Estado y Saldo actual del equipo
	var estadoActualId int
	var saldo float64
	err = conn.QueryRow(`
		SELECT e.estado_id, COALESCE(cr.saldo, 0.00) 
		FROM equipos e
		LEFT JOIN cuentas_reparacion cr ON e.equipo_id = cr.equipo_id
		WHERE e.equipo_id = ?`,
		req.EquipoId,
	).Scan(&estadoActualId, &saldo)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo especificado no existe"})
		}
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando estado/saldo: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar información del equipo"})
	}

	// 4. ESCUDO DE ESTADO: Obligatorio que el equipo esté previamente "Listo para entrega (5)"
	if estadoActualId != 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Solo se pueden entregar formalmente aquellos equipos que estén en estado 'Listo para entrega' (5).",
		})
	}

	// 5. ESCUDO FINANCIERO: El saldo pendiente debe ser exactamente 0
	if saldo > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("No se puede proceder con la entrega. Registra un saldo pendiente de: $%.2f", saldo),
		})
	}

	// 6. Iniciar Transacción Atómica para garantizar consistencia
	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error iniciando transacción: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error interno del servidor"})
	}
	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	// PASO A: Guardar el acta de entrega
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
		strings.ToUpper(req.TrabajosRealizados),
		strings.ToUpper(req.EstadoFinalEquipo),
		req.ConformidadCliente,
		req.ComprobanteNro,
		req.EquipoId,
		req.UsuarioId,
	)
	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error al guardar acta de entrega: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al guardar el acta física de entrega"})
	}

	// PASO B: Actualizar el estado del equipo a 6 (Entregado)
	_, err = tx.Exec(`
		UPDATE equipos 
		SET estado_id = 6, 
		    fecha_modificacion = ? 
		WHERE equipo_id = ?`,
		fechaActual,
		req.EquipoId,
	)
	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error actualizando estado a entregado: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar el estado del equipo"})
	}

	// PASO C: Insertar hito final en 'historial_reparaciones'
	_, err = tx.Exec(`
		INSERT INTO historial_reparaciones (
			observaciones_tecnicas,
			fecha,                  
			usuario_id,
			equipo_id,
			estado_id
		) VALUES (?, ?, ?, ?, 6)`,
		"EQUIPO ENTREGADO AL CLIENTE BAJO ACTA DE CONFORMIDAD",
		fechaActual,
		req.UsuarioId,
		req.EquipoId,
	)
	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando hito final de entrega: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al registrar hito final de entrega"})
	}

	// Confirmar cambios
	err = tx.Commit()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error confirmando transacción: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al confirmar la transacción"})
	}

	// Registrar Log de Éxito
	_ = helpers.InsertLogs(conn, "INSERT/UPDATE", "entregas", req.EquipoId, "Equipo entregado formalmente por el rol: "+rolUsuario)

	return c.Status(200).JSON(fiber.Map{
		"message":   "Entrega procesada con éxito. Acta generada y equipo cerrado.",
		"equipo_id": req.EquipoId,
	})
}
