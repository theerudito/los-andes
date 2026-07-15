export interface HistorialReparacion {
    equipo_id: number;
    estado_id: number;
    observaciones_tecnicas: string;
}

export interface HistorialReparacionDTO {
    historial_id: number;
    observaciones_tecnicas: string;
    fecha: string;
    equipo_id: number;
    equipo: string;
    serie: string;
    estado_id: number;
    estado: string;
    usuario_id: number;
    nombres_usuario: string;
    apellidos_usuario: string;
    cliente_id: number;
    nombres_cliente: string;
    apellidos_cliente: string;
}