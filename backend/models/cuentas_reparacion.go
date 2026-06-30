package models

type CuentasReparacion struct {
	CuentaId   int     `json:"cuenta_id"`
	CostoTotal float64 `json:"costo_total"`
	Abono      float64 `json:"abono"`
	Saldo      float64 `json:"saldo"`

	EquipoId int `json:"equipo_id"`
}

type CuentasReparacionDTO struct {
	CuentaId   int     `json:"cuenta_id"`
	CostoTotal float64 `json:"costo_total"`
	Abono      float64 `json:"abono"`
	Saldo      float64 `json:"saldo"`

	EquipoId int    `json:"equipo_id"`
	Equipo   string `json:"equipo"`
}
