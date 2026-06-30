package models

import "time"

type Tecnicos struct {
	TecnicoId         int       `json:"tecnico_id"`
	Identificacion    string    `json:"identificacion"`
	Nombres           string    `json:"nombres"`
	Apellidos         string    `json:"apellidos"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Activo            bool      `json:"activo"`
	FechaCreacion     time.Time `json:"fecha_creacion"`
	FechaModificacion time.Time `json:"fecha_modificacion"`
}
