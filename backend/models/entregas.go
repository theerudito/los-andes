package models

import "time"

type Entregas struct {
	EntregaId          int       `json:"entrega_id"`
	FechaEntrega       time.Time `json:"fecha_entrega"`
	TrabajosRealizados string    `json:"trabajos_realizados"`
	EstadoFinalEquipo  string    `json:"estado_final_equipo"`
	ConformidadCliente bool      `json:"conformidad_cliente"`
	ComprobanteNumero  string    `json:"comprobante_nro"`

	EquipoId  int `json:"equipo_id"`
	TecnicoId int `json:"tecnico_id"`
}

type EntregasDTO struct {
	EntregaId          int       `json:"entrega_id"`
	FechaEntrega       time.Time `json:"fecha_entrega"`
	TrabajosRealizados string    `json:"trabajos_realizados"`
	EstadoFinalEquipo  string    `json:"estado_final_equipo"`
	ConformidadCliente bool      `json:"conformidad_cliente"`
	ComprobanteNumero  string    `json:"comprobante_nro"`

	EquipoId  int    `json:"equipo_id"`
	Equipo    string `json:"equipo"`
	TecnicoId int    `json:"tecnico_id"`
	Tecnico   string `json:"tecnico"`
}
