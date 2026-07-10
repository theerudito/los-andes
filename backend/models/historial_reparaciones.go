package models

type HistorialReparaciones struct {
	HistorialId           int    `json:"historial_id"`
	ObservacionesTecnicas string `json:"observaciones_tecnicas"`
	FechaCambio           string `json:"fecha_cambio"`

	EquipoId  int `json:"equipo_id"`
	UsuarioId int `json:"usuario_id"`
	EstadoId  int `json:"estado_id"`
}

type HistorialReparacionesDTO struct {
	HistorialId           int    `json:"historial_id"`
	ObservacionesTecnicas string `json:"observaciones_tecnicas"`
	FechaCambio           string `json:"fecha_cambio"`

	EquipoId int    `json:"equipo_id"`
	Equipo   string `json:"equipo"`

	EstadoId int    `json:"estado_id"`
	Estado   string `json:"estado"`

	UsuarioId int `json:"usuario_id"`
	Nombres   int `json:"nombres"`
	Apellidos int `json:"apellidos"`
}
