import api from "../helpers/fetching/axios.ts";
import type {Cuenta} from "../modelos/pagos.ts";

export const pagoService = {

    consultarCuentas: async (equipoId: number) => {
        const { data } = await api.get(`/pago/equipo/${equipoId}`);
        return data;
    },

    consultarCuenta: async (id: number) => {
        const { data } = await api.get(`/pago/${id}`);
        return data;
    },

    crearCuenta: async (payload: Cuenta) => {
        const { data } = await api.post('/pago/actualizar', payload);
        return data;
    },

    actualizarCuenta: async (payload: Cuenta) => {
        const { data } = await api.put('/pago/actualizar', payload);
        return data;
    },

    reporteCuentaPdf: async (id: number) => {
        const response = await api.get(`/pago/pdf/${id}`, { responseType: 'blob' });
        const blob = new Blob([response.data], { type: 'application/pdf' });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `comprobante_pago_${id}.pdf`);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
    },

};