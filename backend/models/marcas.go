package models

type Marcas struct {
	MarcaId           int    `json:"marca_id"`
	Nombre            string `json:"nombre"`
	FechaCreacion     string `json:"fecha_creacion"`
	FechaModificacion string `json:"fecha_modificacion"`
}
