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

export interface Entrega {
    equipo_id: number;
    trabajos_realizados: string;
    observaciones: string;
    estado_final_equipo: string;
    conformidad_cliente: number;
}

export interface EntregaEquipo {
    equipo_id: number;
    usuario_id: number;
    trabajos_realizados: string;
    estado_final_equipo: string;
    conformidad_cliente: number;
    comprobante_nro: string;
    observaciones: string;
}