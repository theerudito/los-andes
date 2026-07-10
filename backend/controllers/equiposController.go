package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"los_andes/database"
	"los_andes/helpers"
	"los_andes/models"
	"strconv"
	"strings"

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
			e.codigo,
			e.tipo_equipo,
			COALESCE(e.modelo, ''),
			COALESCE(e.numero_serie, ''),
			COALESCE(e.accesorios, ''),
			e.descripcion_problema,
			COALESCE(e.observacion, ''),
			COALESCE(e.fecha_recepcion, ''),
			COALESCE(e.fecha_estimada_entrega, ''),
			e.fecha_creacion,
			e.fecha_modificacion,
			m.marca_id,
			m.nombre AS marca,
			c.cliente_id,
			c.nombres,
			c.apellidos,
			r.estado_id,
			r.nombre AS estado
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
	)

	if err = c.BodyParser(&equipo); err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
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

	codigo, err = helpers.ObtenerCodigo(conn, "E")

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
		equipo.UsuarioId,
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

	err = helpers.InsertLogs(conn, "INSERT", "equipos", EquipoId, "registro creado correctamente")
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

	err = conn.QueryRow(`SELECT equipo_id FROM equipos WHERE equipo_id = ?`, equipo.EquipoId).Scan(&EquipoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El registro no existe"})
		}
		_ = helpers.InsertLogsError(conn, "equipos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	var tieneEntrega int
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

	// 4. Iniciar Transacción
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

	err = helpers.InsertLogs(conn, "UPDATE", "equipos", equipo.EquipoId, "registro actualizado correctamente")
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
	)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID inválido"})
	}

	err = conn.QueryRow(`
    SELECT COUNT(*), COALESCE(codigo, ''), estado_id 
    FROM equipos 
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

	err = helpers.InsertLogs(conn, "DELETE", "equipos", id, "Eliminación total del equipo "+codigoEquipo+" (Sin abonos pendientes y en fase operativa)")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro eliminados correctamente"})
}
