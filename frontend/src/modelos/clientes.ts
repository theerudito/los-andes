export interface Cliente {
    cliente_id: number;
    identificacion: string;
    tipo_identificacion: string;
    nombres: string;
    apellidos: string;
    telefono: string;
    email: string;
    direccion: string;
    fecha_creacion: string;
    fecha_modificacion: string;
}

export interface ReqCliente {
    fecha_creacion: string;
    fecha_modificacion: string;
}