import { create } from "zustand";
import {toast} from "sonner";
import type {Cuenta, CuentaDTO} from "../modelos/pagos.ts";
import {pagoService} from "../servicios/pagoServicio.ts";

const initialPago = (): Cuenta => ({
    equipo_id: 0,
    costo_total: 0,
    abono: 0,
    cuenta_id : 0,
});

type Data = {
    form_pagos: Cuenta;
    listar_pagos: CuentaDTO[];
    equipo_id: number,
    isEditing: boolean;
    isLoading: boolean;
    setEquipoId: (equipo_id: number) => void;
    ObtenerPagos: () => Promise<void>;
    ObtenerPago: (id: number) => Promise<void>;
    EnviarPago: () => Promise<void>;
    DescargarPdf: (id: number) => Promise<void>;
    reset: () => void;
};

export const usePagos = create<Data>((set, get) => ({
    form_pagos: initialPago(),
    listar_pagos: [],
    equipo_id: 0,
    isEditing: false,
    isLoading: false,

    setEquipoId: (equipo_id) => set({ equipo_id }),

    ObtenerPagos: async () => {

        set({ isLoading: true });
        try {
            const data = await pagoService.consultarCuentas(get().equipo_id);
            if (Array.isArray(data)) {
                set({ listar_pagos: data, isLoading: false });
            } else {
                set({ listar_pagos: [], isLoading: false });
            }
        } catch (error: any) {
            console.error("Error al obtener lista del pagos:", error.message);
            set({ listar_pagos: [], isLoading: false });
        }
    },

    ObtenerPago: async (id?: number) => {
        const equipo_id = id || get().form_pagos.equipo_id;
        if (!equipo_id) return;

        set({ isLoading: true });
        try {
            const data = await pagoService.consultarCuenta(equipo_id);
            set({ form_pagos: data, isEditing: true, isLoading: false });
        } catch (error: any) {
            console.error(`Error al consultar pago ID ${equipo_id}:`, error);
            set({ isLoading: false });
        }
    },

    EnviarPago: async () => {
        const { form_pagos, isEditing, ObtenerPagos, reset } = get();
        set({ isLoading: true });
        try {
            const payload: Cuenta = {
                cuenta_id: Number(form_pagos.cuenta_id) || 0,
                equipo_id: Number(get().equipo_id) || 0,
                costo_total: Number(form_pagos.costo_total) || 0,
                abono: Number(form_pagos.abono) || 0
            };

            if (isEditing) {
                const data = await pagoService.actualizarCuenta(payload);
                toast.success(data.message);
            } else {
                const data = await pagoService.crearCuenta(payload);
                toast.success(data.message);
            }

            reset();

            await ObtenerPagos();

        } catch (error: any) {
            toast.error(error?.message);
            set({ isLoading: false });
        }
    },

    DescargarPdf: async (id: number) => {
        try {
            await pagoService.reporteCuentaPdf(id);
        } catch (error) {
            console.error("Error al descargar reporte en PDF:", error);
        }
    },

    reset: () =>
        set({
            form_pagos: initialPago(),
            isEditing: false,
            isLoading: false,
        }),
}));