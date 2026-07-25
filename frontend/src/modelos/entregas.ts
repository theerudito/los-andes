export interface Entrega {
    equipo_id: number;
    entrega_id: number;
    trabajos_realizados: string;
    observaciones: string;
    estado_final_equipo: string;
    conformidad_cliente: number;
}

export interface EntregaDTO {
    entrega_id: number;
    fecha_entrega: string;
    trabajos_realizados: string;
    estado_final_equipo: string;
    conformidad_cliente: number;
    comprobante_nro: string;
    equipo_codigo: string;
    nombres: string;
}


