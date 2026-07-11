package models

type Usuarios struct {
	UsuarioId          int    `json:"usuario_id"`
	Identificacion     string `json:"identificacion"`
	TipoIdentificacion string `json:"tipo_identificacion"`
	Nombres            string `json:"nombres"`
	Apellidos          string `json:"apellidos"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Activo             bool   `json:"activo"`
	FechaCreacion      string `json:"fecha_creacion"`
	FechaModificacion  string `json:"fecha_modificacion"`
	RolId              int    `json:"rol_id"`
}

type UsuarioLogin struct {
	Identificacion string `json:"identificacion"`
	Password       string `json:"password"`
}

type UsuarioJWT struct {
	UsuarioId int    `json:"usuario_id"`
	Nombres   string `json:"nombres"`
	Rol       string `json:"rol"`
}
