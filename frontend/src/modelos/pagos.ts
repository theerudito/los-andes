export interface Cuenta {
    cuenta_id: number;
    equipo_id: number;
    costo_total: number;
    abono: number;
}

export interface CuentaDTO {
    cuenta_id: number;
    equipo_id: number;
    codigo: string;
    costo_total: number;
    abono: number;
    saldo: number;
    estado: string;
    nombres: string;
    apellidos: string;
}


