import { create } from "zustand";
import type {Cliente, RPT_Clientes} from "../modelos/clientes.ts";
import {clienteService} from "../servicios/clienteServicio.ts";
import {toast} from "sonner";

const initialCliente = (): Cliente => ({
    cliente_id: 0,
    identificacion: "",
    tipo_identificacion: "",
    nombres: "",
    apellidos: "",
    telefono: "",
    email: "",
    direccion: "",
    fecha_creacion: "",
    fecha_modificacion: "",
});

type Data = {
    form_cliente: Cliente;
    listar_clientes: Cliente[];
    isEditing: boolean;
    clienteId : number;
    isLoading: boolean;
    ObtenerClientes: () => Promise<void>;
    ObtenerCliente: (id?: number) => Promise<void>;
    ObtenerClientePorIdentifiacion: (identificacion: string) => Promise<boolean>;
    EnviarCliente: () => Promise<Cliente | null>;
    EliminarCliente: (id: number) => Promise<void>;
    DescargarPdf: (req: RPT_Clientes) => Promise<void>;
    reset: () => void;
};

export const useClientes = create<Data>((set, get) => ({
    form_cliente: initialCliente(),
    listar_clientes: [],
    isEditing: false,
    isLoading: false,
    clienteId: 0,

    ObtenerClientes: async () => {
        set({ isLoading: true });
        try {
            const data = await clienteService.getClientes();
            if (Array.isArray(data)) {
                set({ listar_clientes: data, isLoading: false });
            } else {
                set({ listar_clientes: [], isLoading: false });
            }
        } catch (error: any) {
            console.error("Error al obtener lista de clientes:", error.message);
            set({ listar_clientes: [], isLoading: false });
        }
    },

    ObtenerClientePorIdentifiacion: async (identificacion: string) => {
        try {
            const data = await clienteService.getClienteByIdentificacion(identificacion);

            if (data && data.cliente_id) {
                const cliente = {
                    cliente_id: data.cliente_id,
                    identificacion: data.identificacion,
                    tipo_identificacion: data.tipo_identificacion,
                    nombres: data.nombres,
                    apellidos: data.apellidos,
                    telefono: data.telefono,
                    email: data.email,
                    direccion: data.direccion,
                    fecha_creacion: data.fecha_creacion,
                    fecha_modificacion: data.fecha_modificacion,
                };

                set({ form_cliente: cliente, clienteId: data.cliente_id });
                return true;
            }

            set({ clienteId: 0 });
            return false;
        } catch (error) {
            set({ clienteId: 0 });
            return false;
        }
    },

    ObtenerCliente: async (id?: number) => {
        const clienteId = id || get().form_cliente.cliente_id;
        if (!clienteId) return;

        set({ isLoading: true });
        try {
            const data = await clienteService.getClienteById(clienteId);
            set({ form_cliente: data, isEditing: true, isLoading: false });
        } catch (error) {
            console.error(`Error al consultar cliente ID ${clienteId}:`, error);
            set({ isLoading: false });
        }
    },

    EnviarCliente: async () => {
        const { form_cliente, isEditing, ObtenerClientes, reset } = get();
        set({ isLoading: true });

        try {
            const payload: Cliente = {
                apellidos: form_cliente.apellidos,
                cliente_id: form_cliente.cliente_id,
                direccion: form_cliente.direccion,
                email: form_cliente.email,
                fecha_creacion: form_cliente.fecha_creacion,
                fecha_modificacion: form_cliente.fecha_modificacion,
                identificacion: form_cliente.identificacion,
                nombres: form_cliente.nombres,
                telefono: form_cliente.telefono,
                tipo_identificacion: form_cliente.tipo_identificacion
            };

            let clienteProcesado: Cliente;

            if (isEditing) {
                const res = await clienteService.modificarCliente(payload);
                toast.success(res.message || "Cliente actualizado correctamente");

                clienteProcesado = {
                    ...payload,
                    cliente_id: res.cliente_id || payload.cliente_id
                };
            } else {
                const res = await clienteService.crearCliente(payload);
                toast.success(res.message || "Cliente creado correctamente");

                clienteProcesado = {
                    ...payload,
                    cliente_id: res.cliente_id
                };
            }

            set({ clienteId: clienteProcesado.cliente_id, isLoading: false });

            reset();
            await ObtenerClientes();

            return clienteProcesado;

        } catch (error: any) {
            toast.error(error?.response?.data?.message || error?.message || "Error al procesar cliente");
            set({ isLoading: false });
            return null;
        }
    },

    EliminarCliente: async (id: number) => {
        set({ isLoading: true });
        try {
            const data = await clienteService.eliminarCliente(id);
            await get().ObtenerClientes();
            toast.success(data.message);
        } catch (error: any) {
            toast.error(error?.message);
            set({ isLoading: false });
        }
    },

    DescargarPdf: async (req: RPT_Clientes) => {
        try {
            await clienteService.reporteClientePdf(req);
        } catch (error: any) {
            let mensajeError = "Error al procesar el reporte";

            if (error?.response?.data instanceof Blob) {
                try {
                    const textoError = await error.response.data.text();
                    const jsonError = JSON.parse(textoError);
                    mensajeError = jsonError.error || mensajeError;
                } catch {
                    mensajeError = "Error al generar el archivo PDF";
                }
            } else if (error?.response?.data?.error) {
                mensajeError = error.response.data.error;
            }

            toast.error(mensajeError);
        }
    },

    reset: () =>
        set({
            form_cliente: initialCliente(),
            isEditing: false,
            isLoading: false,
        }),
}));