import api from "../helpers/fetching/axios.ts";
import type {Equipo} from "../modelos/equipos.ts";

export const equipoService = {
    getEquipos: async () => {
        const { data } = await api.get('/equipo/');
        return data;
    },
    getEquipoById: async (id: number) => {
        const { data } = await api.get(`/equipo/${id}`);
        return data;
    },
    getEquipoPorTipo: async (tipo: string, valor: string) => {
        const { data } = await api.get(`/equipo/${tipo}/${valor}`);
        return data;
    },
    crearEquipo: async (payload: Equipo) => {
        const { data } = await api.post('/equipo/', payload);
        return data;
    },
    modificarEquipo: async (payload: Equipo) => {
        const { data } = await api.put('/equipo/', payload);
        return data;
    },
    eliminarEquipo: async (id: number) => {
        const { data } = await api.delete(`/equipo/${id}`);
        return data;
    },
    descargarOrdenIngresoPdf: async (id: number) => {
        const response = await api.get(`/equipo/orden-ingreso/${id}`, { responseType: 'blob' });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `orden_ingreso_${id}.pdf`);
        document.body.appendChild(link);
        link.click();
        link.remove();
    },
    descargarReporteEquiposPdf: async (payload: Record<string, unknown>) => {
        const response = await api.post('/equipo/reportes', payload, { responseType: 'blob' });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'reporte_equipos.pdf');
        document.body.appendChild(link);
        link.click();
        link.remove();
    },
};