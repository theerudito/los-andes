package models

type LogOk struct {
	LogOkId     int    `json:"log_ok_id"`
	Fecha       string `json:"fecha"`
	Modulo      string `json:"modulo"`
	Usuario     string `json:"usuario"`
	Accion      string `json:"accion"`
	Descripcion string `json:"descripcion"`
}

type LogOkDTO struct {
	FechaDesde string `json:"fecha_desde"`
	FechaHasta string `json:"fecha_hasta"`
	Modulo     string `json:"modulo"`
}
