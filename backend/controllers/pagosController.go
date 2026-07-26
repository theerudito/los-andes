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

	"github.com/gofiber/fiber/v2"
	"github.com/phpdave11/gofpdf"
)

func ConsultarCuentaEquipo(c *fiber.Ctx) error {
	var (
		conn    = database.GetDB()
		cuentas []models.CuentasDTO
		rows    *sql.Rows
		err     error
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
		var cuenta models.CuentasDTO
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

		cuentas = append(cuentas, cuenta)

	}

	if len(cuentas) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(cuentas)

}

func ConsultarCuenta(c *fiber.Ctx) error {
	var (
		conn   = database.GetDB()
		cuenta models.CuentasDTO
		rows   *sql.Rows
		err    error
	)

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de cuenta inválido"})
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
    WHERE c.cuenta_id = ?`, id)

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
	}

	/*
		if len(cuentas) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
		}
	*/

	return c.JSON(cuenta)

}

func ActualizarCuentaEquipo(c *fiber.Ctx) error {
	var (
		conn                     = database.GetDB()
		cuenta                   models.Cuentas
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

	/*
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
		} */

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

func ComprobantePago(c *fiber.Ctx) error {
	var (
		conn = database.GetDB()
		id   = c.Params("id")
		data struct {
			CuentaID              int
			CostoTotal            float64
			Abono                 float64
			Saldo                 float64
			EquipoCodigo          string
			TipoEquipo            string
			Modelo                string
			NumeroSerie           string
			Marca                 string
			ClienteIdentificacion string
			ClienteNombres        string
			ClienteApellidos      string
			ClienteTelefono       string
			ClienteEmail          string
		}
	)

	cuentaID, err := strconv.Atoi(id)
	if err != nil || cuentaID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID de cuenta inválido"})
	}

	query := `
		SELECT 
			c.cuenta_id,
			COALESCE(c.costo_total, 0.00),
			COALESCE(c.abono, 0.00),
			COALESCE(c.saldo, 0.00),

			COALESCE(e.codigo, ''),
			COALESCE(e.tipo_equipo, ''),
			COALESCE(e.modelo, ''),
			COALESCE(e.numero_serie, ''),
			COALESCE(m.nombre, ''),

			COALESCE(cl.identificacion, ''),
			COALESCE(cl.nombres, ''),
			COALESCE(cl.apellidos, ''),
			COALESCE(cl.telefono, ''),
			COALESCE(cl.email, '')
		FROM cuentas_reparacion c
		INNER JOIN equipos e ON c.equipo_id = e.equipo_id
		INNER JOIN clientes cl ON e.cliente_id = cl.cliente_id
		INNER JOIN marcas m ON e.marca_id = m.marca_id
		WHERE c.cuenta_id = ?;`

	err = database.GetDB().QueryRow(query, cuentaID).Scan(
		&data.CuentaID,
		&data.CostoTotal,
		&data.Abono,
		&data.Saldo,
		&data.EquipoCodigo,
		&data.TipoEquipo,
		&data.Modelo,
		&data.NumeroSerie,
		&data.Marca,
		&data.ClienteIdentificacion,
		&data.ClienteNombres,
		&data.ClienteApellidos,
		&data.ClienteTelefono,
		&data.ClienteEmail,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Cuenta de pago no encontrada"})
	} else if err != nil {
		_ = helpers.InsertLogsError(conn, "pagos", "Error al consultar comprobante de pago: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error al consultar los datos del pago"})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(110, 7, "COMPROBANTE DE PAGO", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(70, 7, fmt.Sprintf("RECIBO N°: %06d", data.CuentaID), "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 5, "Sistema de Gestión de Mantenimiento de Computadoras")
	pdf.Ln(8)

	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(180, 6, " 1. DATOS DEL CLIENTE", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	nombreCliente := fmt.Sprintf("%s %s", data.ClienteNombres, data.ClienteApellidos)
	pdf.CellFormat(30, 6, "Cliente:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, helpers.Limitar(nombreCliente, 60), "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Identificación:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.ClienteIdentificacion, "", 0, "L", false, 0, "")
	pdf.CellFormat(25, 6, "Teléfono:", "", 0, "L", false, 0, "")
	pdf.CellFormat(65, 6, data.ClienteTelefono, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Email:", "L,B", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.ClienteEmail, "R,B", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, " 2. DATOS DEL EQUIPO EN SERVICIO", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(30, 6, "Código Equipo:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.EquipoCodigo, "", 0, "L", false, 0, "")
	pdf.CellFormat(25, 6, "Tipo / Marca:", "", 0, "L", false, 0, "")
	pdf.CellFormat(65, 6, fmt.Sprintf("%s - %s", data.TipoEquipo, data.Marca), "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Modelo:", "L,B", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.Modelo, "B", 0, "L", false, 0, "")
	pdf.CellFormat(25, 6, "N° Serie:", "B", 0, "L", false, 0, "")
	pdf.CellFormat(65, 6, data.NumeroSerie, "R,B", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, " 3. RESUMEN DE PAGOS Y SALDO", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 8, fmt.Sprintf("Costo Total: $%.2f", data.CostoTotal), "L,B", 0, "C", false, 0, "")
	pdf.CellFormat(60, 8, fmt.Sprintf("Total Abondado: $%.2f", data.Abono), "B", 0, "C", false, 0, "")
	pdf.CellFormat(60, 8, fmt.Sprintf("Saldo Restante: $%.2f", data.Saldo), "R,B", 1, "C", false, 0, "")

	if data.Saldo <= 0 {
		pdf.Ln(4)
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(0, 128, 0)
		pdf.CellFormat(180, 6, "*** SERVICIO PAGADO EN SU TOTALIDAD ***", "", 1, "C", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.Ln(25)

	yFirmas := pdf.GetY()
	pdf.Line(25, yFirmas, 85, yFirmas)
	pdf.Line(110, yFirmas, 170, yFirmas)

	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(25, yFirmas+2)
	pdf.CellFormat(60, 5, "Firma del Cliente", "", 0, "C", false, 0, "")

	pdf.SetXY(110, yFirmas+2)
	pdf.CellFormat(60, 5, "Firma / Sello de Caja", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "pagos", "Error al procesar PDF de comprobante: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error al generar el comprobante de pago"})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="comprobante_pago_%d.pdf"`, data.CuentaID))
	return c.Send(buf.Bytes())
}
