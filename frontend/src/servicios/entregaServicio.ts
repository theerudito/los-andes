import api from "../helpers/fetching/axios.ts";
import type {Entrega} from "../modelos/entregas.ts";

export const entregaService = {

    consultarEntregas: async (equipoId: number) => {
        const { data } = await api.get(`/entrega/equipo/${equipoId}`);
        return data;
    },

    consultarEntrega: async (id: number) => {
        const { data } = await api.get(`/entrega/${id}`);
        return data;
    },

    crearEntrega: async (payload: Entrega) => {
        const { data } = await api.post('/entrega/', payload);
        return data;
    },

    modificarEntrega: async (payload: Entrega) => {
        const { data } = await api.put('/entrega/', payload);
        return data;
    },

    descargarOrdenEntregaPdf: async (id: number) => {
        const response = await api.get(`/entrega/pdf/${id}`, { responseType: 'blob' });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `orden_entrega.pdf`);
        document.body.appendChild(link);
        link.click();
        link.remove();
    },

};