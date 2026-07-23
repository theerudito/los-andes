package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ConsultarCuentaEquipo(c *fiber.Ctx) error {
	var (
		conn   = database.GetDB()
		cuenta models.CuentaDTO
		rows   *sql.Rows
		err    error
		found  = false
	)

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	rows, err = conn.Query(`
    SELECT 
      c.cuenta_id,
      c.equipo_id,
      e.codigo,
      c.costo_total,
      c.abono,
      c.saldo
    FROM cuentas_reparacion c
    INNER JOIN equipos e ON c.equipo_id = e.equipo_id
    WHERE c.equipo_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "pagos", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&cuenta.CuentaId,
			&cuenta.EquipoId,
			&cuenta.EquipoCodigo,
			&cuenta.CostoTotal,
			&cuenta.Abono,
			&cuenta.Saldo)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "pagos", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true

	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(cuenta)

}

func ActualizarCuentaEquipo(c *fiber.Ctx) error {
	var (
		conn                     = database.GetDB()
		cuenta                   models.Cuenta
		claims                   *models.CustomClaims
		err                      error
		estado                   int
		costoActual, abonoActual float64
		abonoFinal, costoFinal   float64
		tx                       *sql.Tx
	)

	if err := c.BodyParser(&cuenta); err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	if claims.Rol == "TECNICO" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "solo usuario administrador o vendedor puenden realizar esta accion"})
	}

	err = conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, cuenta.EquipoId).Scan(&estado)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo especificado no existe"})
		}
		_ = helpers.InsertLogsError(conn, "equipos", "error consultando estado de equipo: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar estado del equipo"})
	}

	if estado == 6 || estado == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Operación bloqueada. No se pueden registrar pagos en equipos finalizados (Entregados/Cancelados)."})
	}

	err = conn.QueryRow(`
		SELECT costo_total, abono 
		FROM cuentas_reparacion 
		WHERE equipo_id = ?`,
		cuenta.EquipoId).Scan(&costoActual, &abonoActual)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No se encontró registro de cuenta para este equipo"})
		}
		_ = helpers.InsertLogsError(conn, "cuentas", "error consultando cuenta actual: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al leer el estado de la cuenta"})
	}

	if costoActual == 0 && abonoActual == 0 {

		if cuenta.NuevoCosto <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Debe definir un costo total inicial mayor a $0 para poder guardar la cuenta."})
		}

		if cuenta.MontoAAbonar > cuenta.NuevoCosto {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "El abono inicial no puede ser superior al costo total establecido."})
		}

		costoFinal = cuenta.NuevoCosto
		abonoFinal = cuenta.MontoAAbonar

	} else {

		if cuenta.NuevoCosto != 0 && cuenta.NuevoCosto != costoActual {

			if claims.Rol != "ADMINISTRADOR" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permiso denegado. Solo el Administrador puede modificar el costo total una vez establecido."})
			}

			costoFinal = cuenta.NuevoCosto

		} else {

			costoFinal = costoActual

		}

		abonoFinal = abonoActual + cuenta.MontoAAbonar

		if abonoFinal > costoFinal {
			saldoRestante := costoFinal - abonoActual
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("El abono ingresado supera el saldo pendiente. Saldo restante: $%.2f", saldoRestante)})
		}
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = conn.Exec(`
		UPDATE cuentas_reparacion 
		SET costo_total = ?, 
		    abono = ? 
		WHERE equipo_id = ?`,
		costoFinal,
		abonoFinal,
		cuenta.EquipoId,
	)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "error actualizando el registro: "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error actualizando el registro"})
	}

	err = helpers.InsertLogs(tx, "UPDATE", "cuentas", claims.Name, "registro actualizando correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "cuentas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro actualizando correctamente"})

}
