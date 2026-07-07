package models

type LogOk struct {
	LogOkId     int    `json:"log_ok_id"`
	Fecha       string `json:"fecha"`
	Modulo      string `json:"modulo"`
	Accion      string `json:"accion"`
	Descripcion string `json:"descripcion"`
}
