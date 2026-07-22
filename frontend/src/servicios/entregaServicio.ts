import api from "../helpers/fetching/axios.ts";

export interface RegistrarEntregaRequest {
    equipo_id: number;
    trabajos_realizados: string;
    estado_final_equipo: string;
    conformidad_cliente: number;
    observaciones: string;
}

export const entregaService = {
    consultarEntregaPorEquipo: async (equipoId: number) => {
        const { data } = await api.get(`/entrega/${equipoId}`);
        return data;
    },
    registrarEntrega: async (payload: RegistrarEntregaRequest) => {
        const { data } = await api.post('/entrega/', payload);
        return data;
    },
    descargarOrdenEntregaPdf: async (entregaId: number) => {
        const response = await api.get(`/entrega/orden-entrega/${entregaId}`, { responseType: 'blob' });
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `orden_entrega_${entregaId}.pdf`);
        document.body.appendChild(link);
        link.click();
        link.remove();
    },
};