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

	// 1. Validar el cuerpo de la solicitud
	if err = c.BodyParser(&equipo); err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	// 2. Verificar duplicados por número de serie
	err = conn.QueryRow(`SELECT COUNT(*) FROM equipos WHERE numero_serie = ?`, strings.ToUpper(equipo.NumeroSerie)).Scan(&exist)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	if exist > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "el registro ya existe"})
	}

	// 3. Iniciar la transacción
	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	// En caso de cualquier error intermedio, se deshacen todos los cambios de la transacción
	defer tx.Rollback()

	// 4. Obtener el secuencial para el código del equipo
	codigo, err = helpers.ObtenerCodigo(conn, "E")

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error obteniendo el codigo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error obteniendo el codigo"})
	}

	// 5. Insertar en la tabla 'equipos' y obtener el ID generado
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

	// 6. Insertar el hito inicial en 'historial_reparaciones'
	_, err = tx.Exec(`
    INSERT INTO historial_reparaciones (
      observaciones_tecnicas,
      fecha,
      tecnico_id,
      equipo_id,
      estado_id
    ) VALUES (?, ?, ?, ?, ?)`,
		"INGRESO INICIAL: "+strings.ToUpper(equipo.Descripcion),
		helpers.FechaActual(),
		equipo.TecnicoId,
		EquipoId,
		1,
	)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando historial inicial "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error registrando historial de ingreso"})
	}

	// 7. Inicializar la billetera/cuenta en 'cuentas_reparacion'
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

	// 8. Confirmar la transacción en la Base de Datos si todo salió bien
	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	// 9. Auditoría externa a la transacción (Log OK)
	err = helpers.InsertLogs(conn, "INSERT", "equipos", EquipoId, "registro creado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	// 10. Actualizar el contador secuencial para el siguiente código de equipo
	err = helpers.ActualizarCodigo(conn, "E")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error actualizando el codigo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error actualizando el codigo"})
	}

	// Respuesta exitosa al Frontend
	return c.Status(201).JSON(fiber.Map{
		"message":   "registro creado correctamente",
		"equipo_id": EquipoId,
		"codigo":    codigo,
	})
}

func ModificarEquipo(c *fiber.Ctx) error {
	var (
		EquipoId int
		conn     = database.GetDB()
		err      error
		equipo   models.Equipos
		tx       *sql.Tx
	)

	// 1. Validar el cuerpo de la solicitud
	if err = c.BodyParser(&equipo); err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	// 2. Validar si el equipo existe en la base de datos
	err = conn.QueryRow(`SELECT equipo_id FROM equipos WHERE equipo_id = ?`, equipo.EquipoId).Scan(&EquipoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "El registro no existe"})
		}
		_ = helpers.InsertLogsError(conn, "equipos", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	// 3. CANDADO DE SEGURIDAD: Validar si el equipo ya fue entregado formalmente
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

	// 5. Actualizar datos informativos del equipo (No altera el estado_id)
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

	// 6. Confirmar cambios
	err = tx.Commit()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	// 7. Auditoría de acción exitosa
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

	// 1. Obtener y validar el ID de los parámetros de la URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID inválido"})
	}

	// 2. Validar si el equipo existe y obtener su código y estado actual
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

	// =========================================================================
	// ESCUDO 1: CANDADO POR ESTADO (Regla de Negocio)
	// Bloquea el borrado si el equipo ya está en su fase final
	// =========================================================================
	if estadoId == 6 || estadoId == 7 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No se puede eliminar el equipo porque se encuentra en un estado finalizado (Entregado/Cancelado)",
		})
	}

	// =========================================================================
	// ESCUDO 2: CANDADO DE CAJA ANTI-FRAUDE / ANTI-ERROR
	// Bloquea el borrado si el cliente ya dejó dinero abonado (así esté en estado 1)
	// =========================================================================
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

	// =========================================================================
	// 3. Iniciar Transacción Atómica (Solo si pasó los dos candados de arriba)
	// =========================================================================
	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}
	defer tx.Rollback()

	// 4. Romper el RESTRICT de la base de datos eliminando primero la entrega si existiera
	_, err = tx.Exec(`DELETE FROM entregas WHERE equipo_id = ?`, id)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error eliminando entrega asociada "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al procesar la eliminación de entregas"})
	}

	// 5. Eliminar el equipo principal (ON DELETE CASCADE limpia historial y cuenta con abono en 0.00)
	_, err = tx.Exec(`DELETE FROM equipos WHERE equipo_id = ?`, id)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error deleting el equipo "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error al eliminar el registro central"})
	}

	// 6. Confirmar la transacción en bloque
	err = tx.Commit()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	// 7. Auditoría de eliminación completa en log_ok
	descripcionLog := "Eliminación total del equipo " + codigoEquipo + " (Sin abonos pendientes y en fase operativa)"
	err = helpers.InsertLogs(conn, "DELETE", "equipos", id, descripcionLog)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "equipos", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Equipo y todo su historial eliminados correctamente"})
}
