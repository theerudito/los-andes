export interface CuentasReparacion {
    cuenta_id: number;
    costo_total: number;
    abono: number;
    saldo: number;
    equipo_id: number;
}

export interface CuentasReparacionDTO {
    cuenta_id: number;
    costo_total: number;
    abono: number;
    saldo: number;
    equipo_id: number;
    equipo: string;
}