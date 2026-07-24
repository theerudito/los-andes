import { create } from "zustand";
import type {Usuario, UsuarioLogin} from "../modelos/usuarios.ts";
import {usuarioService} from "../servicios/usuarioServicio.ts";
import { toast } from "sonner";
import {ObtenerToken} from "../helpers/jwtDedoce.ts";

const initialUsuario = (): Usuario => ({
    usuario_id: 0,
    identificacion: "",
    tipo_identificacion: "",
    nombres: "",
    apellidos: "",
    email: "",
    password: "",
    activo: false,
    fecha_creacion: "",
    fecha_modificacion: "",
    rol_id: 0,
});

type Data = {
    form_usuario: Usuario;
    listar_usuario: Usuario[];
    isEditing: boolean;
    isLoading: boolean;
    isLogin: boolean;
    ObtenerUsuarios: () => Promise<void>;
    ObtenerUsuario: (id: number) => Promise<void>;
    EnviarUsuario: () => Promise<void>;
    EliminarUsuario: (id: number) => Promise<void>;
    LoginUsuario: (usuario: UsuarioLogin) => Promise<void>;
    ResetUsuario: (usuario: UsuarioLogin) => Promise<void>;
    Logout: () => void;
    reset: () => void;
};

export const useUsuarios = create<Data>((set, get) => ({
    form_usuario: initialUsuario(),
    listar_usuario: [],
    isEditing: false,
    isLoading: false,
    isLogin: !!localStorage.getItem("token"),

    ObtenerUsuarios: async () => {
        set({ isLoading: true });
        try {
            const data = await usuarioService.getUsuarios();
            if (Array.isArray(data)) {
                set({ listar_usuario: data, isLoading: false });
            } else {
                set({ listar_usuario: [], isLoading: false });
            }
        } catch (error: any) {
            console.error("Error al obtener lista de usuarios:", error.message);
            set({ listar_usuario: [], isLoading: false });
        }
    },

    ObtenerUsuario: async (id?: number) => {
        const usuario_id = id || get().form_usuario.usuario_id;
        if (!usuario_id) return;

        set({ isLoading: true });
        try {
            const data = await usuarioService.getUsuarioById(usuario_id);
            set({ form_usuario: data, isEditing: true, isLoading: false });
        } catch (error) {
            console.error(`Error al consultar el usuario ID ${usuario_id}:`, error);
            set({ isLoading: false });
        }
    },

    EnviarUsuario: async () => {
        const { form_usuario, isEditing, ObtenerUsuarios, reset } = get();
        set({ isLoading: true });

        try {
            const payload: Usuario = {
                activo: form_usuario.activo,
                apellidos: form_usuario.apellidos,
                email: form_usuario.email,
                fecha_creacion: form_usuario.fecha_creacion,
                fecha_modificacion: form_usuario.fecha_modificacion,
                identificacion: form_usuario.identificacion,
                nombres: form_usuario.nombres,
                password: form_usuario.password,
                rol_id: form_usuario.rol_id,
                tipo_identificacion: form_usuario.tipo_identificacion,
                usuario_id: form_usuario.usuario_id
            };

            if (isEditing) {
                const data = await usuarioService.modificarUsuario(payload);
                toast.success(data.message);
            } else {
                const data = await usuarioService.crearUsuario(payload);
                toast.success(data.message);
            }

            reset();

            await ObtenerUsuarios();

        } catch (error: any) {
            toast.error(error?.message);
            set({ isLoading: false });
        }
    },

    EliminarUsuario: async (id: number) => {

        const usuario = ObtenerToken()

        if (usuario.user_id === id) {
            toast.info("No es posible borrar a un usuario ya logeado");
            return;
        }

        set({ isLoading: true });
        try {
            const data = await usuarioService.eliminarUsuario(id);
            await get().ObtenerUsuarios();
            toast.success(data.message);
        } catch (error: any) {
            toast.error(error?.message);
        } finally {
            set({ isLoading: false });
        }
    },

    LoginUsuario: async (usuario: UsuarioLogin) => {
        try {
            const data =  await usuarioService.login(usuario);
            localStorage.setItem("token", data.message);
            set({ isLogin: true });
            toast.success("Logueo Correctamente");
        } catch (error: any) {
            toast.error(error?.message);
        }
    },

    ResetUsuario: async (usuario: UsuarioLogin) => {
        try {
            await usuarioService.resetPassword(usuario);
            toast.success("Contraseña Reseteada Correctamente");
        } catch (error: any) {
            toast.error(error?.message);
        }
    },

    Logout: () => {
        localStorage.removeItem("token");
        set({ isLogin: false });
    },

    reset: () =>
        set({
            form_usuario: initialUsuario(),
            isEditing: false,
            isLoading: false,
        }),
}));