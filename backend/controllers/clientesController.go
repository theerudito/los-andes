package controllers

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/phpdave11/gofpdf"
)

func ObtenerClientes(c *fiber.Ctx) error {

	var (
		clientes []models.Clientes
		conn     = database.GetDB()
		rows     *sql.Rows
		err      error
	)

	rows, err = conn.Query(`
		SELECT
			c.cliente_id,
			COALESCE(c.identificacion, '') AS identificacion,
			COALESCE(c.tipo_identificacion, '') AS tipo_identificacion,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(c.telefono, '') AS telefono,
			COALESCE(c.email, '') AS email,
			COALESCE(c.direccion, '') AS direccion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', c.fecha_creacion), '') AS fecha_creacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', c.fecha_modificacion), '') AS fecha_modificacion
		FROM 
			clientes AS c
    ORDER BY 
			c.cliente_id DESC`)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "Error al ejecutar la consulta")
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
			_ = helpers.InsertLogsError(conn, "clientes", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		clientes = append(clientes, cliente)
	}

	if len(clientes) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(clientes)

}

func ObtenerCliente(c *fiber.Ctx) error {
	var (
		cliente models.Clientes
		conn    = database.GetDB()
		id      = c.Params("id")
		rows    *sql.Rows
		err     error
		found   = false
	)

	rows, err = conn.Query(`
		SELECT
			c.cliente_id,
			COALESCE(c.identificacion, '') AS identificacion,
			COALESCE(c.tipo_identificacion, '') AS tipo_identificacion,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(c.telefono, '') AS telefono,
			COALESCE(c.email, '') AS email,
			COALESCE(c.direccion, '') AS direccion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', c.fecha_creacion), '') AS fecha_creacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', c.fecha_modificacion), '') AS fecha_modificacion
		FROM 
			clientes AS c
		WHERE 
			c.cliente_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
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
			&cliente.FechaModificacion,
		)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "clientes", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true
	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(cliente)

}

func ObtenerClientePorIdentificacion(c *fiber.Ctx) error {

	var (
		clientes []models.Clientes
		cliente  models.Clientes
		conn     = database.GetDB()
		valor    = c.Params("identificacion")
		rows     *sql.Rows
		err      error
	)

	rows, err = conn.Query(`
		SELECT
			c.cliente_id,
			COALESCE(c.identificacion, '') AS identificacion,
			COALESCE(c.tipo_identificacion, '') AS tipo_identificacion,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(c.telefono, '') AS telefono,
			COALESCE(c.email, '') AS email,
			COALESCE(c.direccion, '') AS direccion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', c.fecha_creacion), '') AS fecha_creacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', c.fecha_modificacion), '') AS fecha_modificacion
		FROM 
			clientes AS c
		WHERE 
			c.identificacion = ?`, valor)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "Error al ejecutar la consulta")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
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
			_ = helpers.InsertLogsError(conn, "clientes", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		clientes = append(clientes, cliente)
	}

	if len(clientes) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(clientes)

}

func CrearCliente(c *fiber.Ctx) error {
	var (
		ClienteId int
		conn      = database.GetDB()
		exist     int
		err       error
		cliente   models.Clientes
		tx        *sql.Tx
		claims    *models.CustomClaims
	)

	if err = c.BodyParser(&cliente); err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM clientes WHERE identificacion = ?`, strings.ToUpper(cliente.Identificacion)).Scan(&exist)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	if exist > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "el registro ya existe"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO clientes (
			identificacion,
			tipo_identificacion,
			nombres,
			apellidos,
			telefono,
			email,
			direccion,
			fecha_creacion,
			fecha_modificacion
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING cliente_id`,
		cliente.Identificacion,
		helpers.TipoIdentificacion(cliente.Identificacion),
		strings.ToUpper(cliente.Nombres),
		strings.ToUpper(cliente.Apellidos),
		cliente.Telefono,
		cliente.Email,
		strings.ToUpper(cliente.Direccion),
		helpers.FechaActual(),
		helpers.FechaActual()).Scan(&ClienteId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error insertando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "INSERT", "clientes", claims.Name, "registro creado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro creado correctamente"})
}

func ModificarCliente(c *fiber.Ctx) error {

	var (
		ClienteId int
		conn      = database.GetDB()
		err       error
		cliente   models.Clientes
		tx        *sql.Tx
		claims    *models.CustomClaims
	)

	if err = c.BodyParser(&cliente); err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	err = conn.QueryRow(`SELECT cliente_id FROM clientes WHERE cliente_id = ?`, cliente.ClienteId).Scan(&ClienteId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			_ = helpers.InsertLogsError(conn, "cliente", err.Error())
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "clientes", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE clientes 
		SET identificacion 			= ?,
			tipo_identificacion 	= ?,
			nombres 							= ?,
			apellidos 						= ?,
			telefono 							= ?,
			email 								= ?,
			direccion 						= ?,
			fecha_modificacion 		= ?
		WHERE cliente_id 				= ?`,
		cliente.Identificacion,
		helpers.TipoIdentificacion(cliente.Identificacion),
		strings.ToUpper(cliente.Nombres),
		strings.ToUpper(cliente.Apellidos),
		cliente.Telefono,
		cliente.Email,
		strings.ToUpper(cliente.Direccion),
		helpers.FechaActual(),
		cliente.ClienteId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error actualizando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error actualizando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "UPDATE", "clientes", claims.Name, "registro actualizado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro actualizado correctamente"})

}

func EliminarCliente(c *fiber.Ctx) error {
	var (
		ClienteId int
		conn      = database.GetDB()
		err       error
		tx        *sql.Tx
		claims    *models.CustomClaims
	)

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	id, _ := strconv.Atoi(c.Params("id"))

	err = conn.QueryRow(`SELECT cliente_id FROM clientes WHERE cliente_id = ?`, id).Scan(&ClienteId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "clientes", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM clientes WHERE cliente_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error eliminando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error eliminando el registro"})
	}

	err = helpers.InsertLogs(tx, "DELETE", "clientes", claims.Name, "registro eliminado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro eliminado correctamente"})

}

func ReporteCliente(c *fiber.Ctx) error {
	var (
		clientes []models.Clientes
		req      models.ReqReportesClientes
		conn     = database.GetDB()
		rows     *sql.Rows
		err      error
	)

	if err := c.BodyParser(&req); err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "El contenido del json es incorrecto: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	rows, err = conn.Query(`
		SELECT
			c.cliente_id,
			COALESCE(c.identificacion, '') AS identificacion,
			COALESCE(c.tipo_identificacion, '') AS tipo_identificacion,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(c.telefono, '') AS telefono,
			COALESCE(c.email, '') AS email,
			COALESCE(c.direccion, '') AS direccion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', c.fecha_creacion), '') AS fecha_creacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', c.fecha_modificacion), '') AS fecha_modificacion
		FROM 
			clientes AS c
		WHERE 
		     DATE (c.fecha_creacion) BETWEEN ? AND ?
		ORDER BY 
			c.nombres ASC   
			`, req.Fecha_Desde, req.Fecha_Hasta)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "Error al ejecutar la consulta: "+err.Error())
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
			&cliente.FechaModificacion,
		)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "clientes", "Error al leer los registros: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		clientes = append(clientes, cliente)

	}

	if len(clientes) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 6, "REPORTE GENERAL DE CLIENTES")
	pdf.Ln(6)

	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 5, "Sistema de Gestion de Mantenimiento de Computadoras")
	pdf.Ln(10)

	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(1)

	pdf.SetFont("Arial", "B", 10)

	pdf.CellFormat(65, 6, "Nombre Completo", "B", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Identificacion", "B", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "Telefono", "B", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Email", "B", 0, "L", false, 0, "")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 9)

	for _, cli := range clientes {
		nombreCompleto := fmt.Sprintf("%s %s", cli.Nombres, cli.Apellidos)
		identificacion := fmt.Sprintf("%s", cli.Identificacion)

		pdf.CellFormat(65, 6, helpers.Limitar(nombreCompleto, 32), "", 0, "L", false, 0, "")
		pdf.CellFormat(35, 6, identificacion, "", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, cli.Telefono, "", 0, "L", false, 0, "")
		pdf.CellFormat(60, 6, helpers.Limitar(cli.Email, 30), "", 0, "L", false, 0, "")
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes", "Error al procesar salida PDF: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar el archivo PDF"})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", `attachment; filename="reporte_clientes.pdf"`)
	return c.Send(buf.Bytes())
}
