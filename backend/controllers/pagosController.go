package controllers

import (
	"database/sql"
	"errors"
	"los_andes/database"
	"los_andes/helpers"
	"strconv"

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
		EquipoId   int     `json:"equipo_id"`
		CostoTotal float64 `json:"costo_total"`
		Abono      float64 `json:"abono"`
	}

	if err := c.BodyParser(&req); err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	// 1. Regla de negocio: El abono no puede ser mayor al costo total
	if req.Abono > req.CostoTotal {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "El abono no puede ser superior al costo total de la reparación",
		})
	}

	// 2. CANDADO DE SEGURIDAD: Verificar si el equipo ya fue entregado formalmente
	var tieneEntrega int
	err := conn.QueryRow(`SELECT COUNT(*) FROM entregas WHERE equipo_id = ?`, req.EquipoId).Scan(&tieneEntrega)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "error verificando entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el estado del equipo"})
	}
	if tieneEntrega > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se pueden realizar movimientos de caja en un equipo con acta de entrega finalizada",
		})
	}

	// 3. Ejecutar la actualización en la base de datos
	// SQLite calcula la columna 'saldo' automáticamente tras el UPDATE
	result, err := conn.Exec(`
    UPDATE cuentas_reparacion 
    SET costo_total = ?, 
        abono = ? 
    WHERE equipo_id = ?`,
		req.CostoTotal,
		req.Abono,
		req.EquipoId,
	)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "error actualizando cuenta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al procesar el pago"})
	}

	filasAfectadas, _ := result.RowsAffected()
	if filasAfectadas == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No se encontró la cuenta asociada al equipo"})
	}

	// 4. Auditoría de caja
	_ = helpers.InsertLogs(conn, "UPDATE", "cuentas_reparacion", req.EquipoId, "Valores de cuenta/abono actualizados correctamente")

	return c.Status(200).JSON(fiber.Map{"message": "Valores de cuenta actualizados correctamente"})
}
