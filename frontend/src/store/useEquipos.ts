import { create } from "zustand";
import type {ReqCliente} from "../modelos/clientes.ts";
import {toast} from "sonner";
import type {Equipo, EquipoDTO} from "../modelos/equipos.ts";
import {equipoService} from "../servicios/equipoServicio.ts";

const initialEquipo = (): Equipo => ({
    equipo_id: 0,
    codigo: "",
    tipo_equipo: "",
    modelo: "",
    numero_serie: "",
    accesorios: "",
    descripcion_problema: "",
    observacion: "",
    fecha_recepcion: "",
    fecha_estimada_entrega: "",
    fecha_creacion: "",
    fecha_modificacion: "",
    marca_id: 0,
    cliente_id: 0,
    estado_id: 0,
    usuario_id: 0,
});

type Data = {
    form_equipo: Equipo;
    listar_equipos: EquipoDTO[];
    isEditing: boolean;
    isLoading: boolean;
    ObtenerEquipos: () => Promise<void>;
    ObtenerEquipo: (id?: number) => Promise<void>;
    EnviarEquipo: () => Promise<void>;
    EliminarEquipo: (id: number) => Promise<void>;
    DescargarPdf: (req: ReqCliente) => Promise<void>;
    reset: () => void;
};

export const useEquipos = create<Data>((set, get) => ({
    form_equipo: initialEquipo(),
    listar_equipos: [],
    isEditing: false,
    isLoading: false,

    ObtenerEquipos: async () => {
        set({ isLoading: true });
        try {
            const data = await equipoService.getEquipos();
            if (Array.isArray(data)) {
                set({ listar_equipos: data, isLoading: false });
            } else {
                set({ listar_equipos: [], isLoading: false });
            }
        } catch (error: any) {
            console.error("Error al obtener lista de equipos:", error.message);
            set({ listar_equipos: [], isLoading: false });
        }
    },

    ObtenerEquipo: async (id?: number) => {
        const equipoId = id || get().form_equipo.equipo_id;
        if (!equipoId) return;

        set({ isLoading: true });
        try {
            const data = await equipoService.getEquipoById(equipoId);
            set({ form_equipo: data, isEditing: true, isLoading: false });
        } catch (error) {
            console.error(`Error al consultar el equipo ID ${equipoId}:`, error);
            set({ isLoading: false });
        }
    },

    EnviarEquipo: async () => {
        const { form_equipo, isEditing, ObtenerEquipos, reset } = get();
        set({ isLoading: true });

        try {
            const payload: Equipo = {
                accesorios: form_equipo.accesorios,
                cliente_id: form_equipo.cliente_id,
                codigo: form_equipo.codigo,
                descripcion_problema: form_equipo.descripcion_problema,
                equipo_id: form_equipo.equipo_id,
                estado_id: 1,
                fecha_creacion: form_equipo.fecha_creacion,
                fecha_estimada_entrega: form_equipo.fecha_estimada_entrega,
                fecha_modificacion: form_equipo.fecha_modificacion,
                fecha_recepcion: form_equipo.fecha_recepcion,
                marca_id: form_equipo.marca_id,
                modelo: form_equipo.modelo,
                numero_serie: form_equipo.numero_serie,
                observacion: form_equipo.observacion,
                tipo_equipo: form_equipo.tipo_equipo,
                usuario_id: form_equipo.usuario_id

            };

            if (isEditing) {
                const data =  await equipoService.modificarEquipo(payload);
                toast.success(data.message);
            } else {
                const data = await equipoService.crearEquipo(payload);
                toast.success(data.message);
            }

            reset();

            await ObtenerEquipos();

        } catch (error: any) {
            toast.error(error?.message);
            set({ isLoading: false });
        }
    },

    EliminarEquipo: async (id: number) => {
        set({ isLoading: true });
        try {
            const data = await equipoService.eliminarEquipo(id);
            await get().ObtenerEquipos();
            toast.success(data.message);
        } catch (error: any) {
            toast.error(error?.message);
            set({ isLoading: false });
        }
    },

    DescargarPdf: async (req: ReqCliente) => {
        try {
            await equipoService.descargarReporteEquiposPdf(req);
        } catch (error) {
            console.error("Error al descargar reporte en PDF:", error);
        }
    },

    reset: () =>
        set({
            initialEquipo: initialEquipo(),
            isEditing: false,
            isLoading: false,
        }),
}));