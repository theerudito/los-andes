import api from "../helpers/fetching/axios.ts";
import type {Cliente, ReqCliente} from "../modelos/clientes.ts";

export const clienteService = {

    getClientes: async () => {
        const { data } = await api.get('/cliente/');
        return data;
    },
    getClienteById: async (id: number) => {
        const { data } = await api.get(`/cliente/${id}`);
        return data;
    },
    getClienteByIdentificacion: async (identificacion: string) => {
        const { data } = await api.get(`/cliente/dni/${identificacion}`);
        return data;
    },
    crearCliente: async (payload: Cliente) => {
        const { data } = await api.post('/cliente/', payload);
        return data;
    },
    modificarCliente: async (payload: Cliente) => {
        const { data } = await api.put('/cliente/', payload);
        return data;
    },
    eliminarCliente: async (id: number) => {
        const { data } = await api.delete(`/cliente/${id}`);
        return data;
    },
    reporteClientePdf: async (req: ReqCliente) => {
        const response = await api.post('/cliente/reportes', req, { responseType: 'blob' });
        descargarBlob(response.data, 'reporte_clientes.pdf');
    },
};

const descargarBlob = (data: BlobPart, filename: string) => {
    const url = window.URL.createObjectURL(new Blob([data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filename);
    document.body.appendChild(link);
    link.click();
    link.remove();
};