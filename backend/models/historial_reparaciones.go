package models

type HistorialReparaciones struct {
	EquipoId              int    `json:"equipo_id"`
	EstadoId              int    `json:"estado_id"`
	ObservacionesTecnicas string `json:"observaciones_tecnicas"`
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
