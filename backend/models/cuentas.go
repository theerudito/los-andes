package models

type CuentaDTO struct {
	CuentaId     int     `json:"cuenta_id"`
	EquipoId     int     `json:"equipo_id"`
	EquipoCodigo string  `json:"equipo_codigo"`
	CostoTotal   float64 `json:"costo_total"`
	Abono        float64 `json:"abono"`
	Saldo        float64 `json:"saldo"`
}

type Cuenta struct {
	EquipoId     int     `json:"equipo_id"`
	NuevoCosto   float64 `json:"costo_total"`
	MontoAAbonar float64 `json:"abono"`
}
