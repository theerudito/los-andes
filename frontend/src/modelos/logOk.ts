export interface LogOk {
    log_ok_id: number;
    fecha: string;
    modulo: string;
    usuario: string;
    accion: string;
    descripcion: string;
}

export interface LogOkDTO {
    fecha_desde: string;
    fecha_hasta: string;
    modulo: string;
}