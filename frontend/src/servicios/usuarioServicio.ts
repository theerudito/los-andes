import api from "../helpers/fetching/axios.ts";
import type {Usuario, UsuarioLogin} from "../modelos/usuarios.ts";

export const usuarioService = {
    // Públicos
    login: async (datos: UsuarioLogin) => {
        const { data } = await api.post('/usuario/login', datos);
        return data;
    },

    resetPassword: async (datos: UsuarioLogin) => {
        const { data } = await api.put('/usuario/reset', datos);
        return data;
    },

    // Protegidos
    getUsuarios: async () => {
        const { data } = await api.get('/usuario/');
        return data;
    },
    getUsuarioById: async (id: number) => {
        const { data } = await api.get(`/usuario/${id}`);
        return data;
    },
    getUsuarioByIdentificacion: async (identificacion: string) => {
        const { data } = await api.get(`/usuario/dni/${identificacion}`);
        return data;
    },
    crearUsuario: async (payload: Usuario) => {
        const { data } = await api.post('/usuario/', payload);
        return data;
    },
    modificarUsuario: async (payload: Usuario) => {
        const { data } = await api.put('/usuario/', payload);
        return data;
    },
    eliminarUsuario: async (id: number) => {
        const { data } = await api.delete(`/usuario/${id}`);
        return data;
    },
};