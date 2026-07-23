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
			COALESCE(e.codigo, '') AS codigo,
			COALESCE(e.tipo_equipo, '') AS tipo_equipo,
			COALESCE(e.modelo, '') AS modelo,
			COALESCE(e.numero_serie, '') AS numero_serie,
			COALESCE(e.accesorios, '') AS accesorios,
			COALESCE(e.descripcion_problema, '') AS descripcion_problema,
			COALESCE(e.observacion, '') AS observacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_recepcion), '') AS fecha_recepcion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_estimada_entrega), '') AS fecha_estimada_entrega,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_creacion), '') AS fecha_creacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_modificacion), '') AS fecha_modificacion,
			COALESCE(m.marca_id, 0) AS marca_id,
			COALESCE(m.nombre, '') AS marca,
			COALESCE(c.cliente_id, 0) AS cliente_id,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(r.estado_id, 0) AS estado_id,
			COALESCE(r.nombre, '') AS r_nombre
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
			COALESCE(e.codigo, '') AS codigo,
			COALESCE(e.tipo_equipo, '') AS tipo_equipo,
			COALESCE(e.modelo, '') AS modelo,
			COALESCE(e.numero_serie, '') AS numero_serie,
			COALESCE(e.accesorios, '') AS accesorios,
			COALESCE(e.descripcion_problema, '') AS descripcion_problema,
			COALESCE(e.observacion, '') AS observacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_recepcion), '') AS fecha_recepcion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_estimada_entrega), '') AS fecha_estimada_entrega,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_creacion), '') AS fecha_creacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_modificacion), '') AS fecha_modificacion,
			COALESCE(m.marca_id, 0) AS marca_id,
			COALESCE(m.nombre, '') AS marca,
			COALESCE(c.cliente_id, 0) AS cliente_id,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(r.estado_id, 0) AS estado_id,
			COALESCE(r.nombre, '') AS r_nombre
		FROM 
			equipos AS e
		INNER JOIN clientes AS c ON e.cliente_id = c.cliente_id
		INNER JOIN marcas m on e.marca_id = m.marca_id
		INNER JOIN estados_reparacion r on e.estado_id = r.estado_id
		WHERE 
			e.equipo_id = ?`, id)

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

func ObtenerEquipoPorTipo(c *fiber.Ctx) error {
	var (
		conn = database.GetDB()
		err  error
	)

	tipoBusqueda := c.Params("tipo") // 1 nombres - apellido / 2 identificacion / 3 codigo / 4 serie
	valorBusqueda := c.Params("valor")

	if tipoBusqueda == "" || valorBusqueda == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Debe proporcionar el 'tipo' y el 'valor' de búsqueda",
		})
	}

	baseQuery := `
		SELECT
			e.equipo_id,
			COALESCE(e.codigo, '') AS codigo,
			COALESCE(e.tipo_equipo, '') AS tipo_equipo,
			COALESCE(e.modelo, '') AS modelo,
			COALESCE(e.numero_serie, '') AS numero_serie,
			COALESCE(e.accesorios, '') AS accesorios,
			COALESCE(e.descripcion_problema, '') AS descripcion_problema,
			COALESCE(e.observacion, '') AS observacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_recepcion), '') AS fecha_recepcion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_estimada_entrega), '') AS fecha_estimada_entrega,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_creacion), '') AS fecha_creacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_modificacion), '') AS fecha_modificacion,
			COALESCE(m.marca_id, 0) AS marca_id,
			COALESCE(m.nombre, '') AS marca,
			COALESCE(c.cliente_id, 0) AS cliente_id,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(r.estado_id, 0) AS estado_id,
			COALESCE(r.nombre, '') AS r_nombre
		FROM 
			equipos AS e
		INNER JOIN clientes AS c ON e.cliente_id = c.cliente_id
		INNER JOIN marcas m ON e.marca_id = m.marca_id
		INNER JOIN estados_reparacion r ON e.estado_id = r.estado_id`

	var finalQuery string
	var args []interface{}

	// 2. Cláusulas WHERE dinámicas
	switch tipoBusqueda {
	case "1": // Nombre del cliente
		finalQuery = baseQuery + ` WHERE (c.nombres || ' ' || c.apellidos) LIKE ?`
		args = append(args, "%"+valorBusqueda+"%")

	case "2": // Identificación del cliente
		finalQuery = baseQuery + ` WHERE c.identificacion = ?`
		args = append(args, valorBusqueda)

	case "3": // Número de serie
		finalQuery = baseQuery + ` WHERE e.numero_serie = ?`
		args = append(args, valorBusqueda)

	case "4": // Código del equipo
		finalQuery = baseQuery + ` WHERE e.codigo = ?`
		args = append(args, valorBusqueda)

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tipo de búsqueda inválido. Use: 1 (Nombre), 2 (Identificación), 3 (Serie), 4 (Código)",
		})
	}

	rows, err := conn.Query(finalQuery, args...)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Error al ejecutar la consulta: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	listaEquipos := []models.EquiposDTO{}

	for rows.Next() {
		var equipo models.EquiposDTO
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
			&equipo.Estado,
		)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "equipos", "Error al leer los registros: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		listaEquipos = append(listaEquipos, equipo)
	}

	if len(listaEquipos) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(listaEquipos)
}

func CrearEquipo(c *fiber.Ctx) error {
	var (
		EquipoId int
		conn     = database.GetDB()
		exist    int
		err      error
		equipo   models.Equipos
		tx       *sql.Tx
		codigo   string
		claims   *models.CustomClaims
	)

	if err = c.BodyParser(&equipo); err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM equipos WHERE numero_serie = ?`, strings.ToUpper(equipo.NumeroSerie)).Scan(&exist)

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

	codigo, err = helpers.ObtenerCodigo(conn, "O")

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error obteniendo el codigo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error obteniendo el codigo"})
	}

	err = tx.QueryRow(`
    INSERT INTO equipos (
      codigo,
      tipo_equipo,
      modelo,
      numero_serie,
      accesorios,
      descripcion_problema,
      observacion,
      fecha_recepcion,
      fecha_estimada_entrega,
      fecha_creacion,
      fecha_modificacion,
      marca_id,
      cliente_id,
      estado_id                 
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    RETURNING equipo_id`,
		codigo,
		strings.ToUpper(equipo.TipoEquipo),
		strings.ToUpper(equipo.Modelo),
		strings.ToUpper(equipo.NumeroSerie),
		strings.ToUpper(equipo.Accesorios),
		strings.ToUpper(equipo.Descripcion),
		strings.ToUpper(equipo.Observacion),
		equipo.FechaRecepcion,
		equipo.FechaEstimadaEntrega,
		helpers.FechaActual(),
		helpers.FechaActual(),
		equipo.MarcaId,
		equipo.ClienteId,
		1).Scan(&EquipoId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando el registro"})
	}

	_, err = tx.Exec(`
    INSERT INTO historial_reparaciones (
      observaciones_tecnicas,
      fecha,
      usuario_id,
      equipo_id,
      estado_id
    ) VALUES (?, ?, ?, ?, ?)`,
		"INGRESO INICIAL: "+strings.ToUpper(equipo.Descripcion),
		helpers.FechaActual(),
		claims.UserId,
		EquipoId,
		1,
	)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando historial inicial "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error registrando historial de ingreso"})
	}

	_, err = tx.Exec(`
    INSERT INTO cuentas_reparacion (
      costo_total,
      abono,
      equipo_id
    ) VALUES (?, ?, ?)`,
		0.00,
		0.00,
		EquipoId,
	)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error inicializando cuenta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error inicializando cuenta de cobro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "INSERT", "equipos", claims.Name, "registro creado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = helpers.ActualizarCodigo(conn, "E")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error actualizando el codigo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error actualizando el codigo"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro creado correctamente"})

}

func ModificarEquipo(c *fiber.Ctx) error {
	var (
		EquipoId, tieneEntrega int
		conn                   = database.GetDB()
		err                    error
		equipo                 models.Equipos
		tx                     *sql.Tx
		claims                 *models.CustomClaims
	)

	if err = c.BodyParser(&equipo); err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	err = conn.QueryRow(`SELECT equipo_id FROM equipos WHERE equipo_id = ?`, equipo.EquipoId).Scan(&EquipoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El registro no existe"})
		}
		_ = helpers.InsertLogsError(conn, "equipos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM entregas WHERE equipo_id = ?`, equipo.EquipoId).Scan(&tieneEntrega)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error verificando entrega "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar el estado de entrega"})
	}

	if tieneEntrega > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se puede modificar este equipo porque ya cuenta con un acta de entrega formal (Caso Cerrado)",
		})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
    UPDATE equipos 
    SET tipo_equipo           = ?,
      modelo                  = ?,
      numero_serie            = ?,
      accesorios              = ?,
      descripcion_problema    = ?,
      observacion             = ?,
      fecha_recepcion         = ?,
      fecha_estimada_entrega  = ?,
      fecha_modificacion      = ?,
      marca_id                = ?,
      cliente_id              = ?
    WHERE 
      equipo_id               = ?`,
		strings.ToUpper(equipo.TipoEquipo),
		strings.ToUpper(equipo.Modelo),
		strings.ToUpper(equipo.NumeroSerie),
		strings.ToUpper(equipo.Accesorios),
		strings.ToUpper(equipo.Descripcion),
		strings.ToUpper(equipo.Observacion),
		equipo.FechaRecepcion,
		equipo.FechaEstimadaEntrega,
		helpers.FechaActual(),
		equipo.MarcaId,
		equipo.ClienteId,
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

	err = helpers.InsertLogs(conn, "UPDATE", "equipos", claims.Name, "registro actualizado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro actualizado correctamente"})
}

func EliminarEquipo(c *fiber.Ctx) error {
	var (
		exist        int
		estadoId     int
		abono        float64
		conn         = database.GetDB()
		err          error
		tx           *sql.Tx
		codigoEquipo string
		claims       *models.CustomClaims
	)

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID inválido"})
	}

	err = conn.QueryRow(`
    SELECT 
			COUNT(*), 
			COALESCE(codigo, ''), 
			estado_id 
    FROM 
			equipos 
    WHERE equipo_id = ?`, id).Scan(&exist, &codigoEquipo, &estadoId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	if exist == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "el registro no existe"})
	}

	if estadoId == 6 || estadoId == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se puede eliminar el equipo porque se encuentra en un estado finalizado (Entregado/Cancelado)",
		})
	}

	err = conn.QueryRow(`SELECT abono FROM cuentas_reparacion WHERE equipo_id = ?`, id).Scan(&abono)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error consultando abonos "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al verificar la cuenta financiera"})
	}

	if abono > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("No se puede eliminar el equipo porque registra un abono de $%.2f en caja. Primero debe devolver o reversar el dinero.", abono),
		})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM entregas WHERE equipo_id = ?`, id)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error eliminando entrega asociada "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al procesar la eliminación de entregas"})
	}

	_, err = tx.Exec(`DELETE FROM equipos WHERE equipo_id = ?`, id)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error deleting el equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al eliminar el registro central"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "DELETE", "equipos", claims.Name, "registro eliminado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro eliminados correctamente"})

}

func ReporteEquipos(c *fiber.Ctx) error {
	var (
		req     models.ReqReportesEquipos
		equipos []models.EquiposDTO
		equipo  models.EquiposDTO
		conn    = database.GetDB()
		rows    *sql.Rows
		err     error
	)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error al procesar los parámetros de entrada"})
	}

	query := `
		SELECT
			e.equipo_id,
			COALESCE(e.codigo, '') AS codigo,
			COALESCE(e.tipo_equipo, '') AS tipo_equipo,
			COALESCE(e.modelo, '') AS modelo,
			COALESCE(e.numero_serie, '') AS numero_serie,
			COALESCE(e.accesorios, '') AS accesorios,
			COALESCE(e.descripcion_problema, '') AS descripcion_problema,
			COALESCE(e.observacion, '') AS observacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_recepcion), '') AS fecha_recepcion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_estimada_entrega), '') AS fecha_estimada_entrega,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_creacion), '') AS fecha_creacion,
			COALESCE(strftime('%d/%m/%Y %H:%M:%S', e.fecha_modificacion), '') AS fecha_modificacion,
			COALESCE(m.marca_id, 0) AS marca_id,
			COALESCE(m.nombre, '') AS marca,
			COALESCE(c.cliente_id, 0) AS cliente_id,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(r.estado_id, 0) AS estado_id,
			COALESCE(r.nombre, '') AS r_nombre
		FROM 
			equipos AS e
		INNER JOIN clientes AS c ON e.cliente_id = c.cliente_id
		INNER JOIN marcas m ON e.marca_id = m.marca_id
		INNER JOIN estados_reparacion r ON e.estado_id = r.estado_id
		WHERE 
			DATE (e.fecha_recepcion) BETWEEN ? AND ?`

	args := []any{req.Fecha_Desde, req.Fecha_Hasta}

	/*if lg.Modulo != "" {
		query += " AND modulo = ?"
		args = append(args, lg.Modulo)
	}*/

	query += "ORDER BY e.equipo_id DESC"

	rows, err = conn.Query(query, args...)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "log_error", "Error al ejecutar la consulta: "+err.Error())
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
			_ = helpers.InsertLogsError(conn, "equipos", "Error al leer los registros: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		equipos = append(equipos, equipo)
	}

	if len(equipos) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros para el rango de fechas asignado"})
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 6, "REPORTE GENERAL DE EQUIPOS EN SOPORTE")
	pdf.Ln(6)

	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 5, "Sistema de Gestion de Mantenimiento de Computadoras")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 9)
	pdf.Cell(0, 5, fmt.Sprintf("Filtro desde: %s hasta %s", req.Fecha_Desde, req.Fecha_Hasta))
	pdf.Ln(10)

	pdf.Line(10, pdf.GetY(), 287, pdf.GetY())
	pdf.Ln(2)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(25, 6, "Codigo", "B", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "Tipo Equipo", "B", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, "Marca", "B", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Modelo", "B", 0, "L", false, 0, "")
	pdf.CellFormat(65, 6, "Cliente", "B", 0, "L", false, 0, "")
	pdf.CellFormat(32, 6, "F. Recepcion", "B", 0, "C", false, 0, "")
	pdf.CellFormat(60, 6, "Estado", "B", 0, "L", false, 0, "")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 9)
	for _, eq := range equipos {
		clienteCompleto := fmt.Sprintf("%s %s", eq.Nombres, eq.Apellidos)

		pdf.CellFormat(25, 6, eq.Codigo, "", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, helpers.Limitar(eq.TipoEquipo, 15), "", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, helpers.Limitar(eq.Marca, 15), "", 0, "L", false, 0, "")
		pdf.CellFormat(35, 6, helpers.Limitar(eq.Modelo, 18), "", 0, "L", false, 0, "")
		pdf.CellFormat(65, 6, helpers.Limitar(clienteCompleto, 34), "", 0, "L", false, 0, "")
		pdf.CellFormat(32, 6, eq.FechaRecepcion, "", 0, "C", false, 0, "")
		pdf.CellFormat(60, 6, helpers.Limitar(eq.Estado, 20), "", 0, "L", false, 0, "")
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Error al procesar salida PDF: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar el archivo PDF"})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", `attachment; filename="reporte_equipos.pdf"`)
	return c.Send(buf.Bytes())

}

func OrdenIngreso(c *fiber.Ctx) error {
	var (
		data models.ReqOrdenIngresoData
		conn = database.GetDB()
		id   = c.Params("id")
	)

	equipoID, err := strconv.Atoi(id)
	if err != nil || equipoID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de equipo no valido"})
	}

	query := `
		SELECT 
			COALESCE(e.codigo, '') AS codigo,
			COALESCE(e.tipo_equipo, '') AS tipo_equipo,
			COALESCE(e.modelo, '') AS modelo,
			COALESCE(e.numero_serie, '') AS numero_serie,
			COALESCE(e.accesorios, '') AS accesorios,
			COALESCE(e.descripcion_problema, '') AS descripcion_problema,
			COALESCE(e.observacion, '') AS observacion,
			COALESCE(strftime('%d/%m/%Y', e.fecha_recepcion), '') AS fecha_recepcion,
			COALESCE(strftime('%d/%m/%Y', e.fecha_estimada_entrega), '') AS fecha_estimada_entrega,
			COALESCE(m.nombre, '') AS marca,
			COALESCE(est.nombre, '') AS estado,
			COALESCE(c.identificacion, '') AS identificacion,
			COALESCE(c.nombres, '') AS nombres,
			COALESCE(c.apellidos, '') AS apellidos,
			COALESCE(c.telefono, '') AS telefono,
			COALESCE(c.email, '') AS email,
			COALESCE(c.direccion, '') AS direccion
		FROM equipos e
		INNER JOIN clientes c ON e.cliente_id = c.cliente_id
		INNER JOIN marcas m ON e.marca_id = m.marca_id
		INNER JOIN estados_reparacion est ON e.estado_id = est.estado_id
		WHERE e.equipo_id = ?;`

	err = conn.QueryRow(query, equipoID).Scan(
		&data.Codigo,
		&data.TipoEquipo,
		&data.Modelo,
		&data.NumeroSerie,
		&data.Accesorios,
		&data.DescripcionProblema,
		&data.Observacion,
		&data.FechaRecepcion,
		&data.FechaEstimadaEntrega,
		&data.Marca,
		&data.Estado,
		&data.ClienteIdentificacion,
		&data.ClienteNombres,
		&data.ClienteApellidos,
		&data.ClienteTelefono,
		&data.ClienteEmail,
		&data.ClienteDireccion,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Equipo no encontrado"})
	} else if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Error al consultar orden de ingreso: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al consultar los datos del equipo"})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(110, 7, "ORDEN DE INGRESO / RECEPCION", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(70, 7, fmt.Sprintf("ORDEN: %s", data.Codigo), "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 5, "Sistema de Gestion de Mantenimiento de Computadoras")
	pdf.Ln(8)

	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(180, 6, " 1. INFORMACION DEL CLIENTE", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	nombreCliente := fmt.Sprintf("%s %s", data.ClienteNombres, data.ClienteApellidos)
	pdf.CellFormat(30, 6, "Cliente:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, helpers.Limitar(nombreCliente, 60), "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Identificacion:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.ClienteIdentificacion, "", 0, "L", false, 0, "")
	pdf.CellFormat(25, 6, "Telefono:", "", 0, "L", false, 0, "")
	pdf.CellFormat(65, 6, data.ClienteTelefono.String, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Email:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.ClienteEmail, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Direccion:", "L,B", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, helpers.Limitar(data.ClienteDireccion.String, 65), "R,B", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, " 2. DATOS Y ESTADO DEL EQUIPO", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "", 9)

	pdf.CellFormat(30, 6, "Tipo Equipo:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(55, 6, data.TipoEquipo, "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Marca:", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.Marca, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Modelo:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(55, 6, data.Modelo.String, "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Numero de Serie:", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.NumeroSerie.String, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Fecha Recepcion:", "L", 0, "L", false, 0, "")
	pdf.CellFormat(55, 6, data.FechaRecepcion.String, "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Fecha Estimada:", "", 0, "L", false, 0, "")
	pdf.CellFormat(60, 6, data.FechaEstimadaEntrega.String, "R", 1, "L", false, 0, "")

	pdf.CellFormat(30, 6, "Estado Inicial:", "L,B", 0, "L", false, 0, "")
	pdf.CellFormat(150, 6, data.Estado, "R,B", 1, "L", false, 0, "")

	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(180, 6, " 3. ACCESORIOS Y PROBLEMA REPORTADO", "1", 1, "L", true, 0, "")

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(180, 5, "Accesorios Entregados:", "L,R", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	txtAccesorios := data.Accesorios.String
	if txtAccesorios == "" {
		txtAccesorios = "Ninguno"
	}
	pdf.MultiCell(180, 5, txtAccesorios, "L,R,B", "L", false)

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(180, 5, "Descripcion del Problema:", "L,R", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(180, 5, data.DescripcionProblema, "L,R,B", "L", false)

	if data.Observacion.Valid && data.Observacion.String != "" {
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(180, 5, "Observaciones Adicionales:", "L,R", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 9)
		pdf.MultiCell(180, 5, data.Observacion.String, "L,R,B", "L", false)
	}

	pdf.Ln(12)

	pdf.SetFont("Arial", "I", 8)
	pdf.MultiCell(180, 4, "Nota: El taller no se hace responsable por la informacion contenida en los discos duros ni por fallas ocultas no descritas en este documento. Pasados los 30 dias de notificada la reparacion, los equipos no retirados podran generar costo de bodegaje.", "", "C", false)

	pdf.Ln(25)

	yFirmas := pdf.GetY()
	pdf.Line(25, yFirmas, 85, yFirmas)
	pdf.Line(110, yFirmas, 170, yFirmas)

	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(25, yFirmas+2)
	pdf.CellFormat(60, 5, "Firma del Cliente", "", 0, "C", false, 0, "")

	pdf.SetXY(110, yFirmas+2)
	pdf.CellFormat(60, 5, "Firma del Usuario / Tecnico", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Error al procesar PDF de Orden de Ingreso: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar el comprobante PDF"})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="orden_ingreso_%s.pdf"`, data.Codigo))
	return c.Send(buf.Bytes())
}
