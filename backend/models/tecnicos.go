package models

type Tecnicos struct {
	TecnicoId          int    `json:"tecnico_id"`
	Identificacion     string `json:"identificacion"`
	TipoIdentificacion string `json:"tipo_identificacion"`
	Nombres            string `json:"nombres"`
	Apellidos          string `json:"apellidos"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Activo             bool   `json:"activo"`
	FechaCreacion      string `json:"fecha_creacion"`
	FechaModificacion  string `json:"fecha_modificacion"`
}
