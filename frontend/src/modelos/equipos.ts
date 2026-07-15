export interface Equipo {
    equipo_id: number;
    codigo: string;
    tipo_equipo: string;
    modelo: string;
    numero_serie: string;
    accesorios: string;
    descripcion_problema: string;
    observacion: string;
    fecha_recepcion: string;
    fecha_estimada_entrega: string;
    fecha_creacion: string;
    fecha_modificacion: string;
    marca_id: number;
    cliente_id: number;
    estado_id: number;
    usuario_id: number;
}

export interface EquipoDTO {
    equipo_id: number;
    codigo: string;
    tipo_equipo: string;
    modelo: string;
    numero_serie: string;
    accesorios: string;
    descripcion_problema: string;
    observacion: string;
    fecha_recepcion: string;
    fecha_estimada_entrega: string;
    fecha_creacion: string;
    fecha_modificacion: string;
    marca_id: number;
    marca: string;
    estado_id: number;
    estado: string;
    cliente_id: number;
    nombres: string;
    apellidos: string;
}