package models

type Clientes struct {
	ClienteId          int    `json:"cliente_id"`
	Identificacion     string `json:"identificacion"`
	TipoIdentificacion string `json:"tipo_identificacion"`
	Nombres            string `json:"nombres"`
	Apellidos          string `json:"apellidos"`
	Telefono           string `json:"telefono"`
	Email              string `json:"email"`
	Direccion          string `json:"direccion"`
	FechaCreacion      string `json:"fecha_creacion"`
	FechaModificacion  string `json:"fecha_modificacion"`
}
