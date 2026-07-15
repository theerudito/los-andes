export interface CuentaDTO {
    cuenta_id: number;
    equipo_id: number;
    equipo_codigo: string;
    costo_total: number;
    abono: number;
    saldo: number;
}

export interface Cuenta {
    equipo_id: number;
    costo_total: number;
    abono: number;
}