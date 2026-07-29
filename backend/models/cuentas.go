package models

type Cuentas struct {
	CuentaId     int     `json:"cuenta_id"`
	CostoTotal   float64 `json:"costo_total"`
	Abono        float64 `json:"abono"`
	Saldo        float64 `json:"saldo"`
	EquipoId     int     `json:"equipo_id"`
	NuevoCosto   float64 `json:"nuevo_costo,omitempty"`
	MontoAAbonar float64 `json:"monto_a_abonar,omitempty"`
}

type CuentasDTO struct {
	CuentaId   int     `json:"cuenta_id"`
	CostoTotal float64 `json:"costo_total"`
	Abono      float64 `json:"abono"`
	Saldo      float64 `json:"saldo"`
	Codigo     string  `json:"codigo"`
	EquipoId   int     `json:"equipo_id"`
	Equipo     string  `json:"equipo"`
	Estado     string  `json:"estado"`

	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
}
