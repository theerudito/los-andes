package controllers

import (
	"database/sql"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"

	"github.com/gofiber/fiber/v2"
)

func ObtenerEstados(c *fiber.Ctx) error {
	var (
		estados []models.Estados_Reparacion
		estado  models.Estados_Reparacion
		conn    = database.GetDB()
		rows    *sql.Rows
		err     error
	)

	rows, err = conn.Query(`
		SELECT
			e.estado_id,
			e.nombre
		FROM 
			estados_reparacion AS e
    ORDER BY 
			e.estado_id DESC`)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "estados_reparacion", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&estado.EstadoId, &estado.Nombre)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "estados_reparacion", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		estados = append(estados, estado)

	}

	if len(estados) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(estados)

}

func ObtenerEstado(c *fiber.Ctx) error {

	var (
		estado models.Estados_Reparacion
		conn   = database.GetDB()
		id     = c.Params("id")
		rows   *sql.Rows
		err    error
		found  = false
	)

	rows, err = conn.Query(`
		SELECT
			e.estado_id,
			e.nombre
		FROM 
			estados_reparacion AS e
		WHERE 
			e.estado_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "estados_reparacion", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&estado.EstadoId, &estado.Nombre)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "estados_reparacion", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true

	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(estado)

}
