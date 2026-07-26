import type {Historial} from "../modelos/historial.ts";
import api from "../helpers/fetching/axios.ts";

export const historialService = {

    consultarHistoriales: async (equipoId: number) => {
        const { data } = await api.get(`/historial/equipo/${equipoId}`);
        return data;
    },

    consultarHistorial: async (id: number) => {
        const { data } = await api.get(`/historial/${id}`);
        return data;
    },

    crearHistorial: async (payload: Historial) => {
        const { data } = await api.post('/historial/', payload);
        return data;
    },

    actualizarHistorial: async (payload: Historial) => {
        const { data } = await api.put('/historial/', payload);
        return data;
    },

    reporteHistorialPdf: async (id: number) => {
        const response = await api.get(`/historial/pdf/${id}`, { responseType: 'blob' });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `reporte_historial.pdf`);
        document.body.appendChild(link);
        link.click();
        link.remove();
    },

};