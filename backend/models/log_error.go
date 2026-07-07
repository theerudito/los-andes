package models

type LogError struct {
	LogErrorId   int    `json:"log_error_id"`
	Fecha        string `json:"fecha"`
	Modulo       string `json:"modulo"`
	Accion       string `json:"accion"`
	MensajeError string `json:"mensaje_error"`
}
