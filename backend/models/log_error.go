package models

import "time"

type LogError struct {
	LogErrorId   int       `json:"log_error_id"`
	Fecha        time.Time `json:"fecha"`
	Modulo       string    `json:"modulo"`
	Accion       string    `json:"accion"`
	MensajeError string    `json:"mensaje_error"`
}
