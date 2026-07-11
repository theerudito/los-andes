package models

type EntregaDTO struct {
	EntregaId          int    `json:"entrega_id"`
	FechaEntrega       string `json:"fecha_entrega"`
	TrabajosRealizados string `json:"trabajos_realizados"`
	EstadoFinalEquipo  string `json:"estado_final_equipo"`
	ConformidadCliente int    `json:"conformidad_cliente"`
	ComprobanteNro     string `json:"comprobante_nro"`
	EquipoCodigo       string `json:"equipo_codigo"`
	Usuario            string `json:"nombres"`
}

type Entrega struct {
	EquipoId           int    `json:"equipo_id"`
	TrabajosRealizados string `json:"trabajos_realizados"`
	Observaciones      string `json:"observaciones"`
	EstadoFinalEquipo  string `json:"estado_final_equipo"`
	ConformidadCliente int    `json:"conformidad_cliente"`
}

type EntregaEquipo struct {
	EquipoId           int    `json:"equipo_id"`
	UsuarioId          int    `json:"usuario_id"`
	TrabajosRealizados string `json:"trabajos_realizados"`
	EstadoFinalEquipo  string `json:"estado_final_equipo"`
	ConformidadCliente int    `json:"conformidad_cliente"`
	ComprobanteNro     string `json:"comprobante_nro"`
	Observaciones      string `json:"observaciones"`
}
