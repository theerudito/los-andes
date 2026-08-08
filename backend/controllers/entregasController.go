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

func ConsultarEntregaPorEquipo(c *fiber.Ctx) error {

	var (
		conn     = database.GetDB()
		rows     *sql.Rows
		err      error
		entregas []models.EntregaDTO
		id       int
	)

	id, err = strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	rows, err = conn.Query(`
		SELECT 
		  en.entrega_id,
		  eq.equipo_id,
		  COALESCE(strftime('%d/%m/%Y %H:%M:%S', en.fecha_entrega), '') AS fecha_entrega,
		  en.trabajos_realizados,
		  en.estado_final_equipo,
		  en.conformidad_cliente,
		  en.comprobante_nro,
		  eq.codigo,
		  (u.nombres || ' ' || u.apellidos) AS usuario,
		  cli.nombres,
		  cli.apellidos
		FROM entregas en
		INNER JOIN equipos eq ON en.equipo_id = eq.equipo_id
		INNER JOIN clientes cli ON eq.cliente_id = cli.cliente_id
		INNER JOIN usuarios u ON en.usuario_id = u.usuario_id
		WHERE en.equipo_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Error al ejecutar la consulta "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		var entrega models.EntregaDTO
		err = rows.Scan(
			&entrega.EntregaId,
			&entrega.EquipoId,
			&entrega.FechaEntrega,
			&entrega.TrabajosRealizados,
			&entrega.EstadoFinalEquipo,
			&entrega.ConformidadCliente,
			&entrega.ComprobanteNro,
			&entrega.EquipoCodigo,
			&entrega.Usuario,
			&entrega.Nombres,
			&entrega.Apellidos)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "entregas", "Error al leer los registros "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		entregas = append(entregas, entrega)

	}

	if len(entregas) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(entregas)

}

func ConsultarEntrega(c *fiber.Ctx) error {

	var (
		conn    = database.GetDB()
		rows    *sql.Rows
		err     error
		entrega models.EntregaDTO
		id      int
		found   bool
	)

	id, err = strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de equipo inválido"})
	}

	rows, err = conn.Query(`
		SELECT 
		  en.entrega_id,
		  eq.equipo_id,
		  COALESCE(strftime('%d/%m/%Y %H:%M:%S', en.fecha_entrega), '') AS fecha_entrega,
		  en.trabajos_realizados,
		  en.estado_final_equipo,
		  en.conformidad_cliente,
		  en.comprobante_nro,
		  eq.codigo,
		  (u.nombres || ' ' || u.apellidos) AS nombres
		FROM entregas en
		INNER JOIN equipos eq ON en.equipo_id = eq.equipo_id
		INNER JOIN usuarios u ON en.usuario_id = u.usuario_id
		WHERE en.entrega_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Error al ejecutar la consulta "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&entrega.EntregaId,
			&entrega.EquipoId,
			&entrega.FechaEntrega,
			&entrega.TrabajosRealizados,
			&entrega.EstadoFinalEquipo,
			&entrega.ConformidadCliente,
			&entrega.ComprobanteNro,
			&entrega.EquipoCodigo,
			&entrega.Usuario)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "entregas", "Error al leer los registros "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true

	}

	if found == false {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(entrega)

}

func RegistrarEntrega(c *fiber.Ctx) error {

	var (
		entrega models.Entrega
		conn    = database.GetDB()
		err     error
		tx      *sql.Tx
		claims  *models.CustomClaims
		existe  int
		estado  int
		saldo   float64
	)

	if err := c.BodyParser(&entrega); err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error al leer los claims "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los claims"})
	}

	if claims.Rol == "TECNICO" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "solo usuario administrador o vendedor pueden realizar esta acción"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM entregas WHERE equipo_id = ?`, entrega.EquipoId).Scan(&existe)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error verificando duplicidad de entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el historial de entregas"})
	}

	if existe > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ya existe una entrega registrada"})
	}

	err = conn.QueryRow(`SELECT estado_id FROM equipos WHERE equipo_id = ?`, entrega.EquipoId).Scan(&estado)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El equipo no existe"})
		}
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando estado del equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el equipo"})
	}

	if estado != 5 {
		if estado == 6 || estado == 7 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Este equipo ya fue procesado y cerrado anteriormente"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se puede registrar la entrega. El equipo debe estar en estado 'Listo para entrega' primero.",
		})
	}

	err = conn.QueryRow(`SELECT saldo FROM cuentas_reparacion WHERE equipo_id = ?`, entrega.EquipoId).Scan(&saldo)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando saldo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el saldo de la cuenta"})
	}

	if saldo > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("No se puede entregar el equipo. Registra un saldo pendiente de: $%.2f", saldo),
		})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error iniciando transacción"})
	}

	defer tx.Rollback()

	fechaActual := helpers.FechaActual()

	var comprobanteNro string

	for {
		comprobanteNro, err = helpers.ObtenerCodigo(tx, "E")
		if err != nil {
			_ = tx.Rollback()
			_ = helpers.InsertLogsError(conn, "entregas", "error obteniendo secuencial "+err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "error generando número de comprobante"})
		}

		var count int
		err = tx.QueryRow(`SELECT COUNT(*) FROM entregas WHERE comprobante_nro = ?`, comprobanteNro).Scan(&count)
		if err != nil {
			_ = tx.Rollback()
			_ = helpers.InsertLogsError(conn, "entregas", "error verificando correlativo "+err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "error verificando el secuencial generado"})
		}

		if count == 0 {
			break
		}

		err = helpers.ActualizarCodigo(tx, "E")
		if err != nil {
			_ = tx.Rollback()
			_ = helpers.InsertLogsError(conn, "entregas", "error incrementando secuencial "+err.Error())
			return c.Status(500).JSON(fiber.Map{"message": "error al ajustar el secuencial"})
		}
	}

	_, err = tx.Exec(`
		INSERT INTO entregas (
			fecha_entrega,
			trabajos_realizados,
			estado_final_equipo,
			conformidad_cliente,
			comprobante_nro,
			equipo_id,
			usuario_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		fechaActual,
		strings.ToUpper(entrega.TrabajosRealizados),
		strings.ToUpper(entrega.EstadoFinalEquipo),
		entrega.ConformidadCliente,
		comprobanteNro,
		entrega.EquipoId,
		claims.UserId,
	)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al registrar el acta de entrega: " + err.Error()})
	}

	_, err = tx.Exec(`
		UPDATE equipos 
		SET estado_id = 6, 
		    fecha_modificacion = ? 
		WHERE equipo_id = ?`,
		fechaActual,
		entrega.EquipoId,
	)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error actualizando estado del equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar el estado final del equipo"})
	}

	_, err = tx.Exec(`
		INSERT INTO historial_reparaciones (
			observaciones_tecnicas,
			fecha, 
			usuario_id, 
			equipo_id,
			estado_id
		) VALUES (?, ?, ?, ?, 6)`,
		strings.ToUpper(entrega.Observaciones),
		fechaActual,
		claims.UserId,
		entrega.EquipoId,
	)

	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando historial de cierre "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al registrar el hito de cierre"})
	}

	err = helpers.ActualizarCodigo(tx, "E")
	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error actualizando secuencial 'E' "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar el secuencial"})
	}

	err = helpers.InsertLogs(tx, "INSERT", "entregas", claims.Name, "registro creado correctamente")
	if err != nil {
		_ = tx.Rollback()
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al confirmar la entrega"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro creado correctamente"})
}

func ModificarEntrega(c *fiber.Ctx) error {
	var (
		entrega models.Entrega
		conn    = database.GetDB()
		err     error
		tx      *sql.Tx
		claims  *models.CustomClaims
		existe  int
	)

	if err := c.BodyParser(&entrega); err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error al leer los claims "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los claims"})
	}

	if claims.Rol == "TECNICO" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Solo el usuario Administrador o Vendedor pueden realizar esta acción"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM entregas WHERE entrega_id = ?`, entrega.EntregaId).Scan(&existe)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error consultando acta de entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar la entrega"})
	}

	if existe == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El acta de entrega especificada no existe"})
	}

	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE entregas 
		SET trabajos_realizados = ?,
		    estado_final_equipo = ?,
		    conformidad_cliente = ?
		WHERE entrega_id = ?`,
		strings.ToUpper(entrega.TrabajosRealizados),
		strings.ToUpper(entrega.EstadoFinalEquipo),
		entrega.ConformidadCliente,
		entrega.EntregaId,
	)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error actualizando entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al actualizar el acta de entrega: " + err.Error()})
	}

	err = helpers.InsertLogs(tx, "UPDATE", "entregas", claims.Name, "registro actualizado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error insertando la auditoria"})
	}

	err = tx.Commit()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al confirmar la transacción"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro actualizado correctamente"})
}

func OrdenEntrega(c *fiber.Ctx) error {
	var (
		data models.ReqEntregaData
		conn = database.GetDB()
		id   = c.Params("id")
	)

	entregaID, err := strconv.Atoi(id)

	if err != nil || entregaID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de entrega no válido"})
	}

	query := `
		SELECT 
			COALESCE(en.comprobante_nro, 'S/N') AS comprobante_nro,
			COALESCE(strftime('%d/%m/%Y', en.fecha_entrega), '') AS fecha_entrega,
			COALESCE(en.trabajos_realizados, '') AS trabajos_realizados,
			COALESCE(en.estado_final_equipo, '') AS estado_final_equipo,
			COALESCE(en.conformidad_cliente, 1) AS conformidad_cliente,

			COALESCE(e.codigo, '') AS codigo,
			COALESCE(e.tipo_equipo, '') AS tipo_equipo,
			COALESCE(e.modelo, '') AS modelo,
			COALESCE(e.numero_serie, '') AS numero_serie,
			COALESCE(m.nombre, '') AS marca,

			COALESCE(c.identificacion, '') AS identificacion,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(c.telefono, '') AS telefono,
			COALESCE(c.email, '') AS email,
			COALESCE(c.direccion, '') AS direccion,

			COALESCE(cr.costo_total, 0.00) AS costo_total,
			COALESCE(cr.abono, 0.00) AS abono,
			COALESCE(cr.saldo, 0.00) AS saldo
		FROM entregas en
		INNER JOIN equipos e ON en.equipo_id = e.equipo_id
		INNER JOIN clientes c ON e.cliente_id = c.cliente_id
		INNER JOIN marcas m ON e.marca_id = m.marca_id
		LEFT JOIN cuentas_reparacion cr ON e.equipo_id = cr.equipo_id
		WHERE en.entrega_id = ?;`

	err = conn.QueryRow(query, entregaID).Scan(
		&data.ComprobanteNro,
		&data.FechaEntrega,
		&data.TrabajosRealizados,
		&data.EstadoFinal,
		&data.ConformidadCliente,
		&data.Codigo,
		&data.TipoEquipo,
		&data.Modelo,
		&data.NumeroSerie,
		&data.Marca,
		&data.ClienteIdentificacion,
		&data.ClienteNombres,
		&data.ClienteApellidos,
		&data.ClienteTelefono,
		&data.ClienteEmail,
		&data.ClienteDireccion,
		&data.CostoTotal,
		&data.Abono,
		&data.Saldo,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Informacion de entrega no encontrada"})
	} else if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Error al consultar acta de entrega: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al consultar los datos de entrega"})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(110, 7, "ACTA DE ENTREGA DE EQUIPO", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(70, 7, fmt.Sprintf("ENTREGA: %s", data.ComprobanteNro.String), "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 5, "Sistema de Gestion de Mantenimiento de Computadoras")
	pdf.Ln(8)

	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(180, 6, " 1. INFORMACION DEL CLIENTE", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	nombreCliente := fmt.Sprintf("%s %s", data.ClienteNombres.String, data.ClienteApellidos.String)
	pdf.CellFormat(30, 6, "Cliente:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, helpers.Limitar(nombreCliente, 60), "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Identificacion:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.ClienteIdentificacion.String, "", 0, "L", false, 0, "")
	pdf.CellFormat(25, 6, "Telefono:", "", 0, "L", false, 0, "")
	pdf.CellFormat(65, 6, data.ClienteTelefono.String, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Email:", "L,B", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.ClienteEmail.String, "R,B", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, " 2. DETALLES DEL EQUIPO Y SALIDA", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)

	pdf.CellFormat(30, 6, "Codigo Equipo:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(55, 6, data.Codigo.String, "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Fecha de Entrega:", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.FechaEntrega.String, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Tipo / Marca:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(55, 6, fmt.Sprintf("%s - %s", data.TipoEquipo.String, data.Marca.String), "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Numero de Serie:", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.NumeroSerie.String, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Estado Final:", "L,B", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.EstadoFinal.String, "R,B", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, " 3. TRABAJOS Y SERVICIOS REALIZADOS", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(180, 5, data.TrabajosRealizados.String, "L,R,B", "L", false)

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, " 4. ESTADO DE CUENTA Y PAGOS", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(60, 6, fmt.Sprintf("Costo Total: $%.2f", data.CostoTotal), "L,B", 0, "C", false, 0, "")
	pdf.CellFormat(60, 6, fmt.Sprintf("Abonos / Pagos: $%.2f", data.Abono), "B", 0, "C", false, 0, "")
	pdf.CellFormat(60, 6, fmt.Sprintf("Saldo Pendiente: $%.2f", data.Saldo), "R,B", 1, "C", false, 0, "")

	if data.Saldo <= 0 {
		pdf.Ln(3)
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(0, 128, 0)
		pdf.CellFormat(180, 5, "*** CUENTA CANCELADA EN SU TOTALIDAD ***", "", 1, "C", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.Ln(10)

	pdf.SetFont("Arial", "I", 8)
	pdf.MultiCell(180, 4, "Declaracion de Conformidad: Al firmar este documento, el cliente declara recibir el equipo en mencion a entera satisfaccion, habiendo verificado el correcto funcionamiento de los trabajos realizados.", "", "C", false)

	pdf.Ln(25)

	yFirmas := pdf.GetY()
	pdf.Line(25, yFirmas, 85, yFirmas)
	pdf.Line(110, yFirmas, 170, yFirmas)

	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(25, yFirmas+2)
	pdf.CellFormat(60, 5, "Firma de Conformidad Cliente", "", 0, "C", false, 0, "")

	pdf.SetXY(110, yFirmas+2)
	pdf.CellFormat(60, 5, "Firma Usuario", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "entregas", "Error al procesar PDF de Entrega: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar el comprobante de entrega"})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="acta_entrega_%s.pdf"`, data.Codigo.String))
	return c.Send(buf.Bytes())
}
