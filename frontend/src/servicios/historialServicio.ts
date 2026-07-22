import type {HistorialReparacion} from "../modelos/historial.ts";
import api from "../helpers/fetching/axios.ts";

export const historialService = {
    consultarHistorialEquipo: async (equipoId: number) => {
        const { data } = await api.get(`/historial/${equipoId}`);
        return data;
    },
    actualizarEstadoEquipo: async (payload: HistorialReparacion) => {
        const { data } = await api.put('/historial/', payload);
        return data;
    },
    reporteHistorialPdf: async (payload: { fecha_desde: string; fecha_hasta: string; equipo_id?: number | null }) => {
        const response = await api.post('/historial/reportes', payload, { responseType: 'blob' });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `reporte_historial_${payload.fecha_desde}_${payload.fecha_hasta}.pdf`);
        document.body.appendChild(link);
        link.click();
        link.remove();
    },
};