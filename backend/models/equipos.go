package models

import "time"

type Equipos struct {
	EquipoId             int       `json:"equipo_id"`
	Codigo               string    `json:"codigo"`
	TipoEquipo           string    `json:"tipo_equipo"`
	Modelo               string    `json:"modelo"`
	NumeroSerie          string    `json:"numero_serie"`
	Accesorios           string    `json:"accesorios"`
	Descripcion          string    `json:"descripcion_problema"`
	Observacion          string    `json:"observacion"`
	FechaRecepcion       time.Time `json:"fecha_recepcion"`
	FechaEstimadaEntrega time.Time `json:"fecha_estimada_entrega"`
	FechaCreacion        time.Time `json:"fecha_creacion"`
	FechaModificacion    time.Time `json:"fecha_modificacion"`

	MarcaId   int `json:"marca_id"`
	ClienteId int `json:"cliente_id"`
	EstadoId  int `json:"estado_id"`
}

type EquiposDTO struct {
	EquipoId             int       `json:"equipo_id"`
	Codigo               string    `json:"codigo"`
	TipoEquipo           string    `json:"tipo_equipo"`
	Modelo               string    `json:"modelo"`
	NumeroSerie          string    `json:"numero_serie"`
	Accesorios           string    `json:"accesorios"`
	Descripcion          string    `json:"descripcion_problema"`
	Observacion          string    `json:"observacion"`
	FechaRecepcion       time.Time `json:"fecha_recepcion"`
	FechaEstimadaEntrega time.Time `json:"fecha_estimada_entrega"`
	FechaCreacion        time.Time `json:"fecha_creacion"`
	FechaModificacion    time.Time `json:"fecha_modificacion"`

	MarcaId int    `json:"marca_id"`
	Marca   string `json:"marca"`

	EstadoId int    `json:"estado_id"`
	Estado   string `json:"estado"`

	ClienteId int    `json:"cliente_id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
}
