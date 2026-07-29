import { create } from "zustand";
import {toast} from "sonner";
import type {Entrega, EntregaDTO} from "../modelos/entregas.ts";
import {entregaService} from "../servicios/entregaServicio.ts";

const initialEntrega = (): Entrega => ({
    equipo_id: 0,
    entrega_id: 0,
    trabajos_realizados: "",
    estado_final_equipo: "",
    conformidad_cliente: 1,
});

type Data = {
    form_entrega: Entrega;
    listar_entrega: EntregaDTO[];
    equipo_id: number,
    isEditing: boolean;
    isLoading: boolean;
    setEquipoId: (equipo_id: number) => void;
    ObtenerEntregas: () => Promise<void>;
    ObtenerEntrega: (id: number) => Promise<void>;
    EnviarEntrega: () => Promise<void>;
    DescargarPdf: (id: number) => Promise<void>;
    reset: () => void;
};

export const useEntregas = create<Data>((set, get) => ({
    form_entrega: initialEntrega(),
    listar_entrega: [],
    equipo_id: 0,
    isEditing: false,
    isLoading: false,

    setEquipoId: (equipo_id) => set({ equipo_id }),

    ObtenerEntregas: async () => {

        set({ isLoading: true });
        try {
            const data = await entregaService.consultarEntregas(get().equipo_id);
            if (Array.isArray(data)) {
                set({ listar_entrega: data, isLoading: false });
            } else {
                set({ listar_entrega: [], isLoading: false });
            }
        } catch (error: any) {
            console.error("Error al obtener lista de la entrega:", error.message);
            set({ listar_entrega: [], isLoading: false });
        }
    },

    ObtenerEntrega: async (id?: number) => {
        const entregaId = id || get().form_entrega.entrega_id;
        if (!entregaId) return;

        set({ isLoading: true });
        try {
            const data = await entregaService.consultarEntrega(entregaId);
            set({ form_entrega: data, isEditing: true, isLoading: false });
        } catch (error) {
            console.error(`Error al consultar entrega ID ${entregaId}:`, error);
            set({ isLoading: false });
        }
    },

    EnviarEntrega: async () => {
        const { form_entrega, isEditing, ObtenerEntregas, reset } = get();
        set({ isLoading: true });
        try {
            const payload: Entrega = {
                conformidad_cliente: form_entrega.conformidad_cliente,
                entrega_id: form_entrega.entrega_id,
                equipo_id: get().equipo_id,
                estado_final_equipo: form_entrega.estado_final_equipo,
                trabajos_realizados: form_entrega.trabajos_realizados
            };

            if (isEditing) {
                const data = await entregaService.modificarEntrega(payload);
                toast.success(data.message);
            } else {
                const data = await entregaService.crearEntrega(payload);
                toast.success(data.message);
            }

            reset();

            await ObtenerEntregas();

        } catch (error: any) {
            toast.error(error?.message);
            set({ isLoading: false });
        }
    },

    DescargarPdf: async (id: number) => {
        try {
            await entregaService.descargarOrdenEntregaPdf(id);
        } catch (error) {
            console.error("Error al descargar reporte en PDF:", error);
        }
    },

    reset: () =>
        set({
            form_entrega: initialEntrega(),
            isEditing: false,
            isLoading: false,
        }),
}));