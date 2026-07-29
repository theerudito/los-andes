package models

type DashboardResponse struct {
	Equipos  int `json:"equipos"`
	Entregas int `json:"entregas"`
	Clientes int `json:"clientes"`
	Errores  int `json:"errors"`
}
