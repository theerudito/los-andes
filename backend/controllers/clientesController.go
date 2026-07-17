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
	"github.com/xuri/excelize/v2"
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
		WHERE 
			c.cliente_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "movie", "Error al ejecutar la consulta")
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
			_ = helpers.InsertLogsError(conn, "movie", "Error al leer los registros")
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
		WHERE 
			c.identificacion = ?`, valor)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "movie", "Error al ejecutar la consulta")
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
	}

	rows, err = conn.Query(`
		SELECT
			c.cliente_id,
			c.identificacion,
			c.tipo_identificacion,
			c.nombres,
			c.apellidos,
			COALESCE(c.telefono, ''),
			c.email,
			COALESCE(c.direccion, ''),
			c.fecha_creacion,
			c.fecha_modificacion
		FROM 
			clientes AS c
		WHERE 
		     DATE (c.fecha_creacion) BETWEEN ? AND ?
		ORDER BY 
			c.nombres ASC`, req.Fecha_Desde, req.Fecha_Hasta)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "clientes_reporte", "Error al ejecutar la consulta: "+err.Error())
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
			_ = helpers.InsertLogsError(conn, "clientes_reporte", "Error al leer los registros: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		clientes = append(clientes, cliente)

	}

	if len(clientes) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	switch req.Formato {
	// PDF con Gofpdf
	case 1:

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
			_ = helpers.InsertLogsError(conn, "clientes_reporte", "Error al procesar salida PDF: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar el archivo PDF"})
		}

		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", `attachment; filename="reporte_clientes.pdf"`)
		return c.Send(buf.Bytes())

	// EXCEL
	case 2:

		f := excelize.NewFile()

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		sheetName := "Clientes"

		index, err := f.NewSheet(sheetName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear hoja de Excel"})
		}

		f.SetActiveSheet(index)

		_ = f.DeleteSheet("Sheet1")

		styleHeader, err := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"2F4F4F"}, Pattern: 1},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error de estilos"})
		}

		styleData, err := f.NewStyle(&excelize.Style{
			Border: []excelize.Border{
				{Type: "left", Color: "FFFFFF", Style: 1},
				{Type: "top", Color: "D3D3D3", Style: 1},
				{Type: "right", Color: "D3D3D3", Style: 1},
				{Type: "bottom", Color: "D3D3D3", Style: 1},
			},
		})

		headers := []string{"ID Cliente", "Identificación", "Tipo Identificacion", "Nombres", "Apellidos", "Teléfono", "Email", "Dirección", "Fecha Creación"}
		for colIdx, text := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
			_ = f.SetCellValue(sheetName, cell, text)
			_ = f.SetCellStyle(sheetName, cell, cell, styleHeader)
		}
		_ = f.SetRowHeight(sheetName, 1, 25)

		rowIdx := 2
		for _, cli := range clientes {
			_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIdx), cli.ClienteId)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIdx), cli.Identificacion)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIdx), cli.TipoIdentificacion)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIdx), cli.Nombres)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIdx), cli.Apellidos)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowIdx), cli.Telefono)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowIdx), cli.Email)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowIdx), cli.Direccion)
			_ = f.SetCellValue(sheetName, fmt.Sprintf("I%d", rowIdx), cli.FechaCreacion)

			rangeCells := fmt.Sprintf("A%d:I%d", rowIdx, rowIdx)
			_ = f.SetCellStyle(sheetName, rangeCells, rangeCells, styleData)
			_ = f.SetRowHeight(sheetName, rowIdx, 20)

			rowIdx++
		}

		cols, _ := f.GetCols(sheetName)
		for colIdx, colCells := range cols {
			maxLen := 0
			for _, cellVal := range colCells {
				if len(cellVal) > maxLen {
					maxLen = len(cellVal)
				}
			}
			colName, _ := excelize.ColumnNumberToName(colIdx + 1)
			_ = f.SetColWidth(sheetName, colName, colName, float64(maxLen+3))
		}

		buffer, err := f.WriteToBuffer()
		if err != nil {
			_ = helpers.InsertLogsError(conn, "clientes_reporte", "Error al escribir Excel en Buffer: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar archivo Excel"})
		}

		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", `attachment; filename="reporte_clientes.xlsx"`)

		return c.Send(buffer.Bytes())

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Formato de reporte no válido"})
	}
}
