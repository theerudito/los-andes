package models

import "database/sql"

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

type ReqReportesEntregas struct {
	Fecha_Desde string `json:"fecha_desde"`
	Fecha_Hasta string `json:"fecha_hasta"`
}

type ReqEntregaData struct {
	// Datos de la Entrega
	ComprobanteNro     sql.NullString
	FechaEntrega       sql.NullString
	TrabajosRealizados sql.NullString
	EstadoFinal        sql.NullString
	ConformidadCliente int

	// Datos del Equipo
	Codigo      sql.NullString
	TipoEquipo  sql.NullString
	Modelo      sql.NullString
	NumeroSerie sql.NullString
	Marca       sql.NullString

	// Datos del Cliente
	ClienteIdentificacion sql.NullString
	ClienteNombres        sql.NullString
	ClienteApellidos      sql.NullString
	ClienteTelefono       sql.NullString
	ClienteEmail          sql.NullString
	ClienteDireccion      sql.NullString

	CostoTotal float64
	Abono      float64
	Saldo      float64
}
