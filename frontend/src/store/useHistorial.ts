import { create } from "zustand";
import {toast} from "sonner";
import type {Historial, HistorialDTO} from "../modelos/historial.ts";
import {historialService} from "../servicios/historialServicio.ts";

const initialHistorial = (): Historial => ({
    equipo_id: 0,
    estado_id: 0,
    observaciones_tecnicas: "",
    historial_id: 0
});

type Data = {
    form_historial: Historial;
    listar_historial: HistorialDTO[];
    equipo_id: number,
    isEditing: boolean;
    isLoading: boolean;
    setEquipoId: (equipo_id: number) => void;
    ObtenerHistoriales: () => Promise<void>;
    ObtenerHistorial: (id: number) => Promise<void>;
    EnviarHistorial: () => Promise<void>;
    DescargarPdf: (id: number) => Promise<void>;
    reset: () => void;
};

export const useHistorial = create<Data>((set, get) => ({
    form_historial: initialHistorial(),
    listar_historial: [],
    equipo_id: 0,
    isEditing: false,
    isLoading: false,

    setEquipoId: (equipo_id) => set({ equipo_id }),

    ObtenerHistoriales: async () => {

        set({ isLoading: true });
        try {
            const data = await historialService.consultarHistoriales(get().equipo_id);
            if (Array.isArray(data)) {
                set({ listar_historial: data, isLoading: false });
            } else {
                set({ listar_historial: [], isLoading: false });
            }
        } catch (error: any) {
            console.error("Error al obtener lista del historial:", error.message);
            set({ listar_historial: [], isLoading: false });
        }
    },

    ObtenerHistorial: async (id?: number) => {
        const historial_id = id || get().form_historial.historial_id;
        if (!historial_id) return;

        set({ isLoading: true });
        try {
            const data = await historialService.consultarHistorial(historial_id);
            set({ form_historial: data, isEditing: true, isLoading: false });
        } catch (error: any) {
            console.error(`Error al consultar historial ID ${historial_id}:`, error);
            set({ isLoading: false });
        }
    },

    EnviarHistorial: async () => {
        const { form_historial, isEditing, ObtenerHistoriales, reset } = get();
        set({ isLoading: true });
        try {
            const payload: Historial = {
                historial_id: form_historial.historial_id,
                equipo_id: get().equipo_id,
                estado_id: form_historial.estado_id,
                observaciones_tecnicas: form_historial.observaciones_tecnicas
            };

            if (isEditing) {
                const data = await historialService.actualizarHistorial(payload);
                toast.success(data.message);
            } else {
                const data = await historialService.crearHistorial(payload);
                toast.success(data.message);
            }

            reset();

            await ObtenerHistoriales();

        } catch (error: any) {
            toast.error(error?.message);
            set({ isLoading: false });
        }
    },

    DescargarPdf: async (id: number) => {
        try {
            await historialService.reporteHistorialPdf(id);
        } catch (error) {
            console.error("Error al descargar reporte en PDF:", error);
        }
    },

    reset: () =>
        set({
            form_historial: initialHistorial(),
            isEditing: false,
            isLoading: false,
        }),
}));