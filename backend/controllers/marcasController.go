package controllers

import (
	"database/sql"
	"errors"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ObtenerMarcas(c *fiber.Ctx) error {

	var (
		marcas []models.Marcas

		conn = database.GetDB()
		rows *sql.Rows
		err  error
	)

	rows, err = conn.Query(`
		SELECT
			m.marca_id,
			m.nombre,
			COALESCE(strftime('%d/%m/%Y', m.fecha_creacion), '') AS fecha_creacion,
      COALESCE(strftime('%d/%m/%Y', m.fecha_modificacion), '') AS fecha_modificacion
		FROM 
			marcas AS m
    ORDER BY 
			m.marca_id DESC`)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {

		var marca models.Marcas

		err = rows.Scan(
			&marca.MarcaId,
			&marca.Nombre,
			&marca.FechaCreacion,
			&marca.FechaModificacion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "marcas", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		marcas = append(marcas, marca)
	}

	if len(marcas) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(marcas)

}

func ObtenerMarca(c *fiber.Ctx) error {

	var (
		marca models.Marcas
		conn  = database.GetDB()
		id    = c.Params("id")
		rows  *sql.Rows
		err   error
		found = false
	)

	rows, err = conn.Query(`
		SELECT
			m.marca_id,
			m.nombre,
			COALESCE(strftime('%d/%m/%Y', m.fecha_creacion), '') AS fecha_creacion,
      COALESCE(strftime('%d/%m/%Y', m.fecha_modificacion), '') AS fecha_modificacion
		FROM 
			marcas AS m
		WHERE 
			m.marca_id = $1`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&marca.MarcaId,
			&marca.Nombre,
			&marca.FechaCreacion,
			&marca.FechaModificacion,
		)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "marcas", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(marca)

}

func CrearMarca(c *fiber.Ctx) error {

	var (
		MarcaId int
		conn    = database.GetDB()
		exist   int
		err     error
		marca   models.Marcas
		tx      *sql.Tx
	)

	if err = c.BodyParser(&marca); err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM marcas WHERE nombre = $1`, strings.ToUpper(marca.Nombre)).Scan(&exist)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	if exist > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "el registro ya existe"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO marcas (
			nombre,
			fecha_creacion,
			fecha_modificacion
		) VALUES ($1, $2, $3)
		RETURNING marca_id`,
		strings.ToUpper(marca.Nombre),
		time.Now(),
		time.Now()).Scan(&MarcaId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error insertando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "INSERT", "marcas", MarcaId, "registro creado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro creado correctamente"})

}

func ModificarMarca(c *fiber.Ctx) error {

	var (
		MarcaId int
		conn    = database.GetDB()
		err     error
		marca   models.Marcas
		tx      *sql.Tx
	)

	if err = c.BodyParser(&marca); err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	err = conn.QueryRow(`SELECT marca_id ROM marcas WHERE marca_id = $1`, marca.MarcaId).Scan(&MarcaId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "marcas", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE marcas 
		SET nombre 							= $1,
			fecha_modificacion 		= $2
		WHERE 
			marca_id 				  		= $3`,
		strings.ToUpper(marca.Nombre),
		time.Now(),
		marca.MarcaId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error actualizando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error actualizando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "UPDATE", "marcas", marca.MarcaId, "registro actualizado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro actualizado correctamente"})

}

func EliminarMarca(c *fiber.Ctx) error {

	var (
		MarcaId int
		conn    = database.GetDB()
		err     error
		tx      *sql.Tx
	)

	id, _ := strconv.Atoi(c.Params("id"))

	err = conn.QueryRow(`SELECT COUNT(*) FROM marcas WHERE marca_id = $1`, id).Scan(&MarcaId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "marcas", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})

	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM marcas WHERE marca_id = $1`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error eliminando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error eliminando el registro"})
	}

	err = helpers.InsertLogs(tx, "DELETE", "marcas", id, "registro eliminado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "marcas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro eliminado correctamente"})

}
