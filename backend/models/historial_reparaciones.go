package models

type HistorialReparaciones struct {
	EquipoId              int    `json:"equipo_id"`
	EstadoId              int    `json:"estado_id"`
	ObservacionesTecnicas string `json:"observaciones_tecnicas"`
}

type HistorialReparacionesDTO struct {
	HistorialId           int    `json:"historial_id"`
	ObservacionesTecnicas string `json:"observaciones_tecnicas"`
	Fecha                 string `json:"fecha"`

	EquipoId int    `json:"equipo_id"`
	Equipo   string `json:"equipo"`
	Serie    string `json:"serie"`

	EstadoId int    `json:"estado_id"`
	Estado   string `json:"estado"`

	UsuarioId         int    `json:"usuario_id"`
	Nombres_Usuario   string `json:"nombres_usuario"`
	Apellidos_Usuario string `json:"apellidos_usuario"`

	ClienteId         int    `json:"cliente_id"`
	Nombres_Cliente   string `json:"nombres_cliente"`
	Apellidos_Cliente string `json:"apellidos_cliente"`
}
