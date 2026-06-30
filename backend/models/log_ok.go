package models

import "time"

type LogOk struct {
	LogOkId     int       `json:"log_ok_id"`
	Fecha       time.Time `json:"fecha"`
	Modulo      string    `json:"modulo"`
	Accion      string    `json:"accion"`
	Descripcion string    `json:"descripcion"`
}
