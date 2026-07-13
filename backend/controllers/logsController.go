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
		clientes []models.Clientes
		conn     = database.GetDB()
		rows     *sql.Rows
		err      error
	)

	rows, err = conn.Query(`
		SELECT
			c.cliente_id,
			c.identificacion,
			c.tipo_identificacion,
			c.nombres,
			c.apellidos,
			c.telefono,
			c.email,
			c.direccion,
			c.fecha_creacion,
			c.fecha_modificacion
		FROM 
			clientes AS c
    ORDER BY 
			c.cliente_id DESC`)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "movie", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {

		var cliente models.Clientes

		err = rows.Scan(
			&cliente.ClienteId,
			&cliente.Identificacion,
			&cliente.TipoIdentificacion,
			&cliente.Nombres,
			&cliente.Apellidos,
			&cliente.Telefono,
			&cliente.Email,
			&cliente.Direccion,
			&cliente.FechaCreacion,
			&cliente.FechaModificacion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "movie", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		clientes = append(clientes, cliente)
	}

	if len(clientes) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(clientes)

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

	rows, err = conn.Query(`
		SELECT
			log_ok_id, 
			fecha, 
			modulo, 
			usuario, 
			accion, 
			descripcion	
		FROM 
			log_ok
		WHERE DATE(fecha) BETWEEN ? AND ?
    	ORDER BY 
			fecha DESC`, lg.FechaDesde, lg.FechaDesde)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "logs_ok", "Error al ejecutar la consulta")
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
