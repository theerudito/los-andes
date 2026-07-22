import api from "../helpers/fetching/axios.ts";

export interface ActualizarCuentaRequest {
    equipo_id: number;
    costo_total: number;
    abono: number;
    usuario_id: number;
}

export const pagoService = {
    consultarCuentaEquipo: async (equipoId: number) => {
        const { data } = await api.get(`/pago/${equipoId}`);
        return data;
    },
    actualizarCuentaEquipo: async (payload: ActualizarCuentaRequest) => {
        const { data } = await api.put('/pago/actualizar', payload);
        return data;
    },
    procesarEntregaEquipo: async (payload: Record<string, unknown>) => {
        const { data } = await api.post('/pago/procesar', payload);
        return data;
    },
};