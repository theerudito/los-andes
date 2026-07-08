package controllers

import (
	"database/sql"
	"errors"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ObtenerTecnicos(c *fiber.Ctx) error {

	var (
		tecnicos []models.Tecnicos
		tecnico  models.Tecnicos
		conn     = database.GetDB()
		rows     *sql.Rows
		err      error
	)

	rows, err = conn.Query(`
		SELECT
			t.tecnico_id,
			t.identificacion,
			t.tipo_identificacion,
			t.nombres,
			t.apellidos,
			t.email,
			t.activo,
			t.fecha_creacion,
      t.fecha_modificacion
		FROM 
			tecnicos AS t
    ORDER BY 
			t.tecnico_id DESC`)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&tecnico.TecnicoId,
			&tecnico.Identificacion,
			&tecnico.TipoIdentificacion,
			&tecnico.Nombres,
			&tecnico.Apellidos,
			&tecnico.Email,
			&tecnico.Activo,
			&tecnico.FechaCreacion,
			&tecnico.FechaModificacion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "tecnicos", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		tecnicos = append(tecnicos, tecnico)
	}

	if len(tecnicos) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(tecnicos)

}

func ObtenerTecnico(c *fiber.Ctx) error {

	var (
		tecnico models.Tecnicos
		conn    = database.GetDB()
		id      = c.Params("id")
		rows    *sql.Rows
		err     error
		found   = false
	)

	rows, err = conn.Query(`
		SELECT
			t.tecnico_id,
			t.identificacion,
			t.tipo_identificacion,
			t.nombres,
			t.apellidos,
			t.email,
			t.activo,
			t.fecha_creacion,
      t.fecha_modificacion
		FROM 
			tecnicos AS t
		WHERE 
			t.tecnico_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&tecnico.TecnicoId,
			&tecnico.Identificacion,
			&tecnico.TipoIdentificacion,
			&tecnico.Nombres,
			&tecnico.Apellidos,
			&tecnico.Email,
			&tecnico.Activo,
			&tecnico.FechaCreacion,
			&tecnico.FechaModificacion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "tecnicos", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true

	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(tecnico)

}

func ObtenerTecnicoPorIdentificacion(c *fiber.Ctx) error {

	var (
		tecnicos []models.Tecnicos
		tecnico  models.Tecnicos
		conn     = database.GetDB()
		valor    = c.Params("identificacion")
		rows     *sql.Rows
		err      error
	)

	rows, err = conn.Query(`
		SELECT
			t.tecnico_id,
			t.identificacion,
			t.tipo_identificacion,
			t.nombres,
			t.apellidos,
			t.email,
			t.activo,
			t.fecha_creacion,
      t.fecha_modificacion
		FROM 
			tecnicos AS t
		WHERE 
			t.identificacion = ?`, valor)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&tecnico.TecnicoId,
			&tecnico.Identificacion,
			&tecnico.TipoIdentificacion,
			&tecnico.Nombres,
			&tecnico.Apellidos,
			&tecnico.Email,
			&tecnico.Activo,
			&tecnico.FechaCreacion,
			&tecnico.FechaModificacion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "tecnicos", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		tecnicos = append(tecnicos, tecnico)

	}

	if len(tecnicos) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(tecnicos)

}

func CrearTecnico(c *fiber.Ctx) error {

	var (
		TecnicoId int
		conn      = database.GetDB()
		exist     int
		err       error
		tecnico   models.Tecnicos
		tx        *sql.Tx
	)

	if err = c.BodyParser(&tecnico); err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM tecnicos WHERE identificacion = ?`, strings.ToUpper(tecnico.Identificacion)).Scan(&exist)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	if exist > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "el registro ya existe"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	passHash, err := helpers.EncriptarDato(tecnico.Password)

	print(passHash)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error incriptando el dato "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error incriptando el dato"})
	}

	err = tx.QueryRow(`
		INSERT INTO tecnicos (
			identificacion,
			tipo_identificacion,
			nombres,
			apellidos,
			email,
			password,
			activo,
			fecha_creacion,
			fecha_modificacion
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING tecnico_id`,
		tecnico.Identificacion,
		helpers.TipoIdentificacion(tecnico.Identificacion),
		strings.ToUpper(tecnico.Nombres),
		strings.ToUpper(tecnico.Apellidos),
		tecnico.Email,
		passHash,
		true,
		helpers.FechaActual(),
		helpers.FechaActual()).Scan(&TecnicoId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error insertando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "INSERT", "tecnicos", TecnicoId, "registro creado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro creado correctamente"})

}

func ModificarTecnico(c *fiber.Ctx) error {

	var (
		TecnicoId int
		conn      = database.GetDB()
		err       error
		tecnico   models.Tecnicos
		tx        *sql.Tx
	)

	if err = c.BodyParser(&tecnico); err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	err = conn.QueryRow(`SELECT tecnico_id FROM tecnicos WHERE tecnico_id = ?`, tecnico.TecnicoId).Scan(&TecnicoId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "tecnicos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE tecnicos 
		SET identificacion 				= ?,
			  tipo_identificacion 	= ?,
				nombres 							= ?,
				apellidos 						= ?,
				email 								= ?,
				activo 								= ?,
				fecha_modificacion 		= ?
		WHERE 
			tecnico_id 				  		= ?`,
		tecnico.Identificacion,
		helpers.TipoIdentificacion(tecnico.Identificacion),
		strings.ToUpper(tecnico.Nombres),
		strings.ToUpper(tecnico.Apellidos),
		tecnico.Email,
		true,
		helpers.FechaActual(),
		tecnico.TecnicoId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error actualizando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error actualizando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "UPDATE", "tecnicos", tecnico.TecnicoId, "registro actualizado correctamente")

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro actualizado correctamente"})

}

func EliminarTecnico(c *fiber.Ctx) error {

	var (
		TecnicoId int
		conn      = database.GetDB()
		err       error
		tx        *sql.Tx
	)

	id, _ := strconv.Atoi(c.Params("id"))

	err = conn.QueryRow(`SELECT COUNT(*) FROM tecnicos WHERE tecnico_id = ?`, id).Scan(&TecnicoId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "tecnicos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})

	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM tecnicos WHERE tecnico_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error eliminando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error eliminando el registro"})
	}

	err = helpers.InsertLogs(tx, "DELETE", "tecnicos", id, "registro eliminado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "tecnicos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro eliminado correctamente"})

}
