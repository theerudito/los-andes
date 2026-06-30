package models

import "time"

type Marcas struct {
	MarcaId           int       `json:"marca_id"`
	Nombre            string    `json:"nombre"`
	FechaCreacion     time.Time `json:"fecha_creacion"`
	FechaModificacion time.Time `json:"fecha_modificacion"`
}
