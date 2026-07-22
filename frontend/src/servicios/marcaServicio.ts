import api from "../helpers/fetching/axios.ts";
import type {Marca} from "../modelos/marcas.ts";

export const marcaService = {
    getMarcas: async () => {
        const { data } = await api.get('/marca/');
        return data;
    },
    getMarcaById: async (id: number) => {
        const { data } = await api.get(`/marca/${id}`);
        return data;
    },
    crearMarca: async (payload: Marca) => {
        const { data } = await api.post('/marca/', payload);
        return data;
    },
    modificarMarca: async (payload: Marca) => {
        const { data } = await api.put('/marca/', payload);
        return data;
    },
    eliminarMarca: async (id: number) => {
        const { data } = await api.delete(`/marca/${id}`);
        return data;
    },
};

export const estadoService = {
    getEstados: async () => {
        const { data } = await api.get('/estado/');
        return data;
    },
    getEstadoById: async (id: number) => {
        const { data } = await api.get(`/estado/${id}`);
        return data;
    },
};