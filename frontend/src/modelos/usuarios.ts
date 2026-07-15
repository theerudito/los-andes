export interface Usuario {
    usuario_id: number;
    identificacion: string;
    tipo_identificacion: string;
    nombres: string;
    apellidos: string;
    email: string;
    password?: string;
    activo: boolean;
    fecha_creacion: string;
    fecha_modificacion: string;
    rol_id: number;
}

export interface UsuarioLogin {
    identificacion: string;
    password: string;
}

export interface UsuarioJWT {
    usuario_id: number;
    nombres: string;
    rol: string;
}