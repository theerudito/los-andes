package models

type ReqReportesClientes struct {
	Formato     int    `json:"formato"`
	Fecha_Desde string `json:"fecha_desde"`
	Fecha_Hasta string `json:"fecha_hasta"`
}
