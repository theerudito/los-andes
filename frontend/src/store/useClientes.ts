import { create } from "zustand";
import type {Cliente, ReqCliente} from "../modelos/clientes.ts";
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
    isLoading: boolean;
    ObtenerClientes: () => Promise<void>;
    ObtenerCliente: (id?: number) => Promise<void>;
    EnviarCliente: () => Promise<void>;
    EliminarCliente: (id: number) => Promise<void>;
    DescargarPdf: (req: ReqCliente) => Promise<void>;
    reset: () => void;
};

export const useClientes = create<Data>((set, get) => ({
    form_cliente: initialCliente(),
    listar_clientes: [],
    isEditing: false,
    isLoading: false,

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

            if (isEditing) {
                const data =  await clienteService.modificarCliente(payload);
                toast.success(data.message);
            } else {
                const data = await clienteService.crearCliente(payload);
                toast.success(data.message);
            }

            reset();

            await ObtenerClientes();

        } catch (error: any) {
            toast.error(error?.message);
            set({ isLoading: false });
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

    DescargarPdf: async (req: ReqCliente) => {
        try {
            await clienteService.reporteClientePdf(req);
        } catch (error) {
            console.error("Error al descargar reporte en PDF:", error);
        }
    },

    reset: () =>
        set({
            form_cliente: initialCliente(),
            isEditing: false,
            isLoading: false,
        }),
}));