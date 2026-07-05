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

func ObtenerEquipos(c *fiber.Ctx) error {
	var (
		equipos []models.EquiposDTO
		equipo  models.EquiposDTO
		conn    = database.GetDB()
		rows    *sql.Rows
		err     error
	)

	rows, err = conn.Query(`
		SELECT
			e.equipo_id,
			e.codigo,
			e.tipo_equipo,
			e.modelo,
			e.numero_serie,
			e.accesorios,
			e.descripcion_problema,
			e.observacion,
			e.fecha_recepcion,
			e.fecha_estimada_entrega,
			e.fecha_creacion,
			e.fecha_modificacion,
			m.marca_id,
			m.nombre AS marca,
			c.cliente_id,
			c.nombres,
			c.apellidos,
			r.estado_id,
			r.nombre As estado
		FROM 
			equipos AS e
		INNER JOIN clientes AS c ON e.cliente_id = c.cliente_id
		INNER JOIN marcas m on e.marca_id = m.marca_id
		INNER JOIN estados_reparacion r on e.estado_id = r.estado_id
    ORDER BY 
			e.equipo_id DESC`)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&equipo.EquipoId,
			&equipo.Codigo,
			&equipo.TipoEquipo,
			&equipo.Modelo,
			&equipo.NumeroSerie,
			&equipo.Accesorios,
			&equipo.Descripcion,
			&equipo.Observacion,
			&equipo.FechaRecepcion,
			&equipo.FechaEstimadaEntrega,
			&equipo.FechaCreacion,
			&equipo.FechaModificacion,
			&equipo.MarcaId,
			&equipo.Marca,
			&equipo.ClienteId,
			&equipo.Nombres,
			&equipo.Apellidos,
			&equipo.EstadoId,
			&equipo.Estado)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "equipos", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		equipos = append(equipos, equipo)
	}

	if len(equipos) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(equipos)
}

func ObtenerEquipo(c *fiber.Ctx) error {
	var (
		equipo models.EquiposDTO
		conn   = database.GetDB()
		id     = c.Params("id")
		rows   *sql.Rows
		err    error
		found  = false
	)

	rows, err = conn.Query(`
		SELECT
			e.equipo_id,
			e.codigo,
			e.tipo_equipo,
			e.modelo,
			e.numero_serie,
			e.accesorios,
			e.descripcion_problema,
			e.observacion,
			e.fecha_recepcion,
			e.fecha_estimada_entrega,
			e.fecha_creacion,
			e.fecha_modificacion,
			m.marca_id,
			m.nombre AS marca,
			c.cliente_id,
			c.nombres,
			c.apellidos,
			r.estado_id,
			r.nombre As estado
		FROM 
			equipos AS e
		INNER JOIN clientes AS c ON e.cliente_id = c.cliente_id
		INNER JOIN marcas m on e.marca_id = m.marca_id
		INNER JOIN estados_reparacion r on e.estado_id = r.estado_id
		WHERE 
			e.equipo_id = $1`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&equipo.EquipoId,
			&equipo.Codigo,
			&equipo.TipoEquipo,
			&equipo.Modelo,
			&equipo.NumeroSerie,
			&equipo.Accesorios,
			&equipo.Descripcion,
			&equipo.Observacion,
			&equipo.FechaRecepcion,
			&equipo.FechaEstimadaEntrega,
			&equipo.FechaCreacion,
			&equipo.FechaModificacion,
			&equipo.MarcaId,
			&equipo.Marca,
			&equipo.ClienteId,
			&equipo.Nombres,
			&equipo.Apellidos,
			&equipo.EstadoId,
			&equipo.Estado)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "equipos", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(equipo)
}

func CrearEquipo(c *fiber.Ctx) error {
	var (
		EquipoId int
		conn     = database.GetDB()
		exist    int
		err      error
		equipo   models.EquiposDTO
		tx       *sql.Tx
	)

	if err = c.BodyParser(&equipo); err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM equipos WHERE numero_serie = $1`, strings.ToUpper(equipo.NumeroSerie)).Scan(&exist)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	if exist > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "el registro ya existe"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO equipos (
			codigo,
			tipo_equipo,
		    modelo,
			numero_serie,
		    accesorios,
			descripcion_problema,
		    observacion,
			fecha_creacion,
		    fecha_modificacion,
			fecha_recepcion,
		    fecha_estimada_entrega,
			marca_id,
			cliente_id,
		    estado_id                 
		) VALUES ($1, $2, $3,$4, $5, $6, $7, $8, $9 $10, $11, $12, $13, $14)
		RETURNING equipo_id`,
		strings.ToUpper(equipo.Codigo),
		strings.ToUpper(equipo.TipoEquipo),
		strings.ToUpper(equipo.Modelo),
		equipo.NumeroSerie,
		strings.ToUpper(equipo.Accesorios),
		strings.ToUpper(equipo.Descripcion),
		strings.ToUpper(equipo.Observacion),
		equipo.FechaRecepcion,
		equipo.FechaEstimadaEntrega,
		time.Now(),
		time.Now(),
		equipo.MarcaId,
		equipo.ClienteId,
		equipo.EstadoId).Scan(&EquipoId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "INSERT", "equipos", EquipoId, "registro creado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro creado correctamente"})
}

func ModificarEquipo(c *fiber.Ctx) error {
	var (
		EquipoId int
		conn     = database.GetDB()
		err      error
		equipo   models.Equipos
		tx       *sql.Tx
	)

	if err = c.BodyParser(&equipo); err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	err = conn.QueryRow(`SELECT equipo_id FROM equipos WHERE equipo_id = $1`, equipo.EquipoId).Scan(&EquipoId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "equipos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE equipos 
		SET codigo 							= $1,
			tipo_equipo 					= $2,
			modelo 							= $3,
			numero_serie 					= $4,
			accesorios 						= $5,
			descripcion_problema 			= $6,
			observacion 					= $7,
			fecha_recepcion 				= $8,
			fecha_estimada_entrega 			= $9,
			fecha_modificacion 				= $10,
			marca_id 						= $11,
			cliente_id 						= $12,
			estado_id 						= $13
		WHERE 
			equipo_id 				  		= $14`,
		strings.ToUpper(equipo.Codigo),
		strings.ToUpper(equipo.TipoEquipo),
		strings.ToUpper(equipo.Modelo),
		equipo.NumeroSerie,
		strings.ToUpper(equipo.Accesorios),
		strings.ToUpper(equipo.Descripcion),
		strings.ToUpper(equipo.Observacion),
		equipo.FechaRecepcion,
		equipo.FechaEstimadaEntrega,
		time.Now(),
		equipo.MarcaId,
		equipo.ClienteId,
		equipo.EstadoId,
		equipo.EquipoId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error actualizando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error actualizando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "UPDATE", "marcas", equipo.EquipoId, "registro actualizado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro actualizado correctamente"})
}

func EliminarEquipo(c *fiber.Ctx) error {
	var (
		EquipoId int
		conn     = database.GetDB()
		err      error
		tx       *sql.Tx
	)

	id, _ := strconv.Atoi(c.Params("id"))

	err = conn.QueryRow(`SELECT COUNT(*) FROM equipos WHERE equipo_id = $1`, id).Scan(&EquipoId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "equipos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})

	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM equipos WHERE equipo_id = $1`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error eliminando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error eliminando el registro"})
	}

	err = helpers.InsertLogs(tx, "DELETE", "equipos", id, "registro eliminado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro eliminado correctamente"})
}
