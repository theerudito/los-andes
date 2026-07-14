package controllers

import (
	"database/sql"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"

	"github.com/gofiber/fiber/v2"
)

func ObtenerLogsError(c *fiber.Ctx) error {

	var (
		logs []models.LogError
		conn = database.GetDB()
		rows *sql.Rows
		err  error
		lg   models.LogOkDTO
	)

	if err = c.BodyParser(&lg); err != nil {
		_ = helpers.InsertLogsError(conn, "log_error", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	query := `
		SELECT
			log_error_id,
			fecha,
			modulo,
			mensaje_error
		FROM 
			log_error 
		WHERE DATE(fecha) BETWEEN ? AND ?`

	args := []interface{}{lg.FechaDesde, lg.FechaHasta}

	if lg.Modulo != "" {
		query += " AND modulo = ?"
		args = append(args, lg.Modulo)
	}

	query += " " + "ORDER" + " BY fecha DESC"

	rows, err = conn.Query(query, args...)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "log_error", "Error al ejecutar la consulta: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {

		var log models.LogError

		err = rows.Scan(
			&log.LogErrorId,
			&log.Fecha,
			&log.Modulo,
			&log.MensajeError)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "log_error", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		logs = append(logs, log)
	}

	if len(logs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(logs)

}

func ObtenerLogsOk(c *fiber.Ctx) error {
	var (
		logs []models.LogOk
		conn = database.GetDB()
		rows *sql.Rows
		err  error
		lg   models.LogOkDTO
	)

	if err = c.BodyParser(&lg); err != nil {
		_ = helpers.InsertLogsError(conn, "logs_ok", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	query := `
		SELECT
			log_ok_id, 
			fecha, 
			modulo, 
			usuario, 
			accion, 
			descripcion	
		FROM 
			log_ok
		WHERE DATE(fecha) BETWEEN ? AND ?`

	args := []interface{}{lg.FechaDesde, lg.FechaHasta}

	if lg.Modulo != "" {
		query += " AND modulo = ?"
		args = append(args, lg.Modulo)
	}

	query += " " + "ORDER" + " BY fecha DESC"

	rows, err = conn.Query(query, args...)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "log_ok", "Error al ejecutar la consulta: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {

		var log models.LogOk

		err = rows.Scan(
			&log.LogOkId,
			&log.Fecha,
			&log.Modulo,
			&log.Usuario,
			&log.Accion,
			&log.Descripcion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "logs_ok", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		logs = append(logs, log)
	}

	if len(logs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(logs)
}
