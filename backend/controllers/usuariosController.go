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

func ObtenerUsuarios(c *fiber.Ctx) error {

	var (
		usuarios []models.Usuarios
		usuario  models.Usuarios
		conn     = database.GetDB()
		rows     *sql.Rows
		err      error
	)

	rows, err = conn.Query(`
		SELECT
			u.usuario_id,
			u.identificacion,
			u.tipo_identificacion,
			u.nombres,
			u.apellidos,
			u.email,
			u.activo,
			u.rol_id,
			u.fecha_creacion,
      u.fecha_modificacion
		FROM 
			usuarios AS u
    ORDER BY 
			u.usuario_id DESC`)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "Error al ejecutar la consulta "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&usuario.UsuarioId,
			&usuario.Identificacion,
			&usuario.TipoIdentificacion,
			&usuario.Nombres,
			&usuario.Apellidos,
			&usuario.Email,
			&usuario.Activo,
			&usuario.RolId,
			&usuario.FechaCreacion,
			&usuario.FechaModificacion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "usuarios", "Error al leer los registros "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		usuarios = append(usuarios, usuario)
	}

	if len(usuarios) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(usuarios)

}

func ObtenerUsuario(c *fiber.Ctx) error {

	var (
		usuario models.Usuarios
		conn    = database.GetDB()
		id      = c.Params("id")
		rows    *sql.Rows
		err     error
		found   = false
	)

	rows, err = conn.Query(`
		SELECT
			u.usuario_id,
			u.identificacion,
			u.tipo_identificacion,
			u.nombres,
			u.apellidos,
			u.email,
			u.activo,
			u.rol_id,
			u.fecha_creacion,
      u.fecha_modificacion
		FROM 
			usuarios AS u
		WHERE 
			u.usuario_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "Error al ejecutar la consulta "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&usuario.UsuarioId,
			&usuario.Identificacion,
			&usuario.TipoIdentificacion,
			&usuario.Nombres,
			&usuario.Apellidos,
			&usuario.Email,
			&usuario.Activo,
			&usuario.RolId,
			&usuario.FechaCreacion,
			&usuario.FechaModificacion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "usuarios", "Error al leer los registros "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		found = true

	}

	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(usuario)

}

func ObtenerUsuarioPorIdentificacion(c *fiber.Ctx) error {

	var (
		usuarios []models.Usuarios
		usuario  models.Usuarios
		conn     = database.GetDB()
		valor    = c.Params("identificacion")
		rows     *sql.Rows
		err      error
	)

	rows, err = conn.Query(`
		SELECT
			u.usuario_id,
			u.identificacion,
			u.tipo_identificacion,
			u.nombres,
			u.apellidos,
			u.email,
			u.activo,
			u.rol_id,
			u.fecha_creacion,
      u.fecha_modificacion
		FROM 
			usuarios AS u
		WHERE 
			u.identificacion = ?`, valor)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "Error al ejecutar la consulta "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al ejecutar la consulta"})
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&usuario.UsuarioId,
			&usuario.Identificacion,
			&usuario.TipoIdentificacion,
			&usuario.Nombres,
			&usuario.Apellidos,
			&usuario.Email,
			&usuario.Activo,
			&usuario.RolId,
			&usuario.FechaCreacion,
			&usuario.FechaModificacion)

		if err != nil {
			_ = helpers.InsertLogsError(conn, "usuarios", "Error al leer los registros")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer los registros"})
		}

		usuarios = append(usuarios, usuario)

	}

	if len(usuarios) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No se encontraron registros"})
	}

	return c.JSON(usuarios)

}

func CrearUsuario(c *fiber.Ctx) error {

	var (
		conn      = database.GetDB()
		exist     int
		err       error
		usuario   models.Usuarios
		tx        *sql.Tx
		UsuarioId int
		passHash  string
		claims    *models.CustomClaims
	)

	if err = c.BodyParser(&usuario); err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	if claims.Rol == "TECNICO" || claims.Rol == "VENDEDOR" {
		return c.Status(409).JSON(fiber.Map{"messaje": "solo usuarios administrador puenden realizar esta accion"})
	}

	err = conn.QueryRow(`SELECT COUNT(*) FROM usuarios WHERE identificacion = ?`, strings.ToUpper(usuario.Identificacion)).Scan(&exist)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	if exist > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "el registro ya existe"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	passHash, err = helpers.EncriptarDato(usuario.Password)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error incriptando el dato "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error incriptando el dato"})
	}

	err = tx.QueryRow(`
		INSERT INTO usuarios (
			identificacion,
			tipo_identificacion,
			nombres,
			apellidos,
			email,
			password,
			activo,
			rol_id,
			fecha_creacion,
			fecha_modificacion
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING usuario_id`,
		usuario.Identificacion,
		helpers.TipoIdentificacion(usuario.Identificacion),
		strings.ToUpper(usuario.Nombres),
		strings.ToUpper(usuario.Apellidos),
		usuario.Email,
		passHash,
		true,
		usuario.RolId,
		helpers.FechaActual(),
		helpers.FechaActual()).Scan(&UsuarioId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error insertando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "INSERT", "usuarios", claims.Name, "registro creado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro creado correctamente"})

}

func ModificarUsuario(c *fiber.Ctx) error {

	var (
		UsuarioId int
		conn      = database.GetDB()
		err       error
		usuario   models.Usuarios
		tx        *sql.Tx
		claims    *models.CustomClaims
	)

	if err = c.BodyParser(&usuario); err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	if claims.Rol == "TECNICO" || claims.Rol == "VENDEDOR" {
		return c.Status(409).JSON(fiber.Map{"messaje": "solo usuarios administrador puenden realizar esta accion"})
	}

	err = conn.QueryRow(`SELECT usuario_id FROM usuarios WHERE usuario_id = ?`, usuario.UsuarioId).Scan(&UsuarioId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "usuarios", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})
	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE usuarios 
		SET identificacion 				= ?,
			  tipo_identificacion 	= ?,
				nombres 							= ?,
				apellidos 						= ?,
				email 								= ?,
				activo 								= ?,
				rol_id								= ?,
				fecha_modificacion 		= ?
		WHERE 
			usuario_id 				  		= ?`,
		usuario.Identificacion,
		helpers.TipoIdentificacion(usuario.Identificacion),
		strings.ToUpper(usuario.Nombres),
		strings.ToUpper(usuario.Apellidos),
		usuario.Email,
		true,
		usuario.RolId,
		helpers.FechaActual(),
		usuario.UsuarioId)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error actualizando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error actualizando el registro"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "UPDATE", "usuarios", claims.Name, "registro actualizado correctamente")

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registro actualizado correctamente"})

}

func EliminarUsuario(c *fiber.Ctx) error {

	var (
		UsuarioId int
		conn      = database.GetDB()
		err       error
		tx        *sql.Tx
		claims    *models.CustomClaims
	)

	if claims.Rol == "TECNICO" || claims.Rol == "VENDEDOR" {
		return c.Status(409).JSON(fiber.Map{"messaje": "solo usuario administrador puenden realizar esta accion"})
	}

	id, _ := strconv.Atoi(c.Params("id"))

	err = conn.QueryRow(`SELECT COUNT(*) FROM usuarios WHERE usuario_id = ?`, id).Scan(&UsuarioId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
		}

		_ = helpers.InsertLogsError(conn, "usuarios", "error ejecutando la consulta "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error ejecutando la consulta"})

	}

	tx, err = conn.Begin()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error iniciando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error iniciando transacción"})
	}

	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM usuarios WHERE usuario_id = ?`, id)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error eliminando el registro "+err.Error())
		return c.Status(500).JSON(fiber.Map{"message": "error eliminando el registro"})
	}

	err = helpers.InsertLogs(tx, "DELETE", "usuarios", claims.Name, "registro eliminado correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error insertando la auditoria "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error insertando la auditoria"})
	}

	err = tx.Commit()

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error confirmando transacción "+err.Error())
		return c.Status(500).JSON(fiber.Map{"messaje": "error confirmando transacción"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "registro eliminado correctamente"})

}

func LoginUsuario(c *fiber.Ctx) error {
	var (
		conn                                       = database.GetDB()
		err                                        error
		login                                      models.UsuarioLogin
		passwordDB, nombres, apellidos, token, rol string
		usuarioID                                  int
	)

	if err = c.BodyParser(&login); err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cuerpo de solicitud inválido"})
	}

	err = conn.QueryRow(`
		SELECT
			COALESCE(u.usuario_id, 0) AS usuario_id,
			COALESCE(u.password, '') AS password,
			COALESCE(u.nombres, '') AS nombres,
			COALESCE(u.apellidos, '') AS apellidos,
			r.nombre AS rol
		FROM usuarios AS u
		INNER JOIN roles AS r ON u.rol_id = r.rol_id
		WHERE identificacion = ?`,
		login.Identificacion).Scan(
		&usuarioID,
		&passwordDB,
		&nombres,
		&apellidos,
		&rol)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "usuario no existe"})
	}

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	passwordDB, err = helpers.DesencriptarDato(passwordDB)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error desencriptando la contraseña"})
	}

	if passwordDB != login.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "contraseña incorrecta"})
	}

	token, err = helpers.GenerateToken(models.UsuarioJWT{UsuarioId: usuarioID, Nombres: helpers.ObtenerPalabra(nombres, apellidos), Rol: rol})

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error generando token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": token})

}

func ResetUsuario(c *fiber.Ctx) error {

	type ResetDatos struct {
		Identificacion string `json:"identificacion"`
		Password       string `json:"password"`
	}

	var (
		usuarioID int
		conn      = database.GetDB()
		err       error
		reset     ResetDatos
		tx        *sql.Tx
		passHash  string
		claims    *models.CustomClaims
	)

	if err = c.BodyParser(&reset); err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "Cuerpo de solicitud inválido")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cuerpo de solicitud inválido"})
	}

	claims, err = helpers.ReadClaims(c)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", "error al leer los clains "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error al leer los clains"})
	}

	if claims.Rol == "TECNICO" || claims.Rol == "VENDEDOR" {
		return c.Status(409).JSON(fiber.Map{"messaje": "solo usuario administrador puenden realizar esta accion"})
	}

	err = conn.QueryRow(`SELECT usuario_id FROM usuarios WHERE identificacion = ?`, reset.Identificacion).Scan(&usuarioID)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "registro no existe"})
	}

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	tx, err = conn.Begin()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error iniciando transacción"})
	}

	defer tx.Rollback()

	passHash, err = helpers.EncriptarDato(reset.Password)
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error encriptando la contraseña"})
	}

	_, err = tx.Exec(`
		UPDATE usuarios
		SET
			password 					 = ?,
			fecha_modificacion = ?
		WHERE
			usuario_id = ?
	`,
		passHash,
		helpers.FechaActual(),
		usuarioID,
	)

	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	err = tx.Commit()
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error confirmando transacción"})
	}

	err = helpers.InsertLogs(conn, "UPDATE", "usuarios", claims.Name, "contraseña actualizada correctamente")
	if err != nil {
		_ = helpers.InsertLogsError(conn, "usuarios", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "error insertando la auditoría"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "contraseña actualizada correctamente"})

}
