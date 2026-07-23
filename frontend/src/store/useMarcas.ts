import { create } from "zustand";
import type {Marca} from "../modelos/marcas.ts";
import {marcaService} from "../servicios/marcaServicio.ts";

const initialMarca = (): Marca => ({
    marca_id: 0,
    nombre: "",
    fecha_creacion: "",
    fecha_modificacion: ""
});

type Data = {
    form_marca: Marca;
    listar_marca: Marca[];
    isEditing: boolean;
    isLoading: boolean;
    ObtenerMarcas: () => Promise<void>;
    ObtenerMarca: (id: number) => Promise<void>;
    EnviarMarca: () => Promise<void>;
    EliminarMarca: (id: number) => Promise<void>;
    reset: () => void;
};

export const useMarcas = create<Data>((set, get) => ({
    form_marca: initialMarca(),
    listar_marca: [],
    isEditing: false,
    isLoading: false,

    ObtenerMarcas: async () => {
        set({ isLoading: true });
        try {
            const data = await marcaService.getMarcas();
            if (Array.isArray(data)) {
                set({ listar_marca: data, isLoading: false });
            } else {
                set({ listar_marca: [], isLoading: false });
            }
        } catch (error: any) {
            console.error("Error al obtener lista de marcas:", error.message);
            set({ listar_marca: [], isLoading: false });
        }
    },

    ObtenerMarca: async (id?: number) => {
        const marca_id = id || get().listar_marca.marca_id;
        if (!marca_id) return;

        set({ isLoading: true });
        try {
            const data = await marcaService.getMarcaById(marca_id);
            set({ form_marca: data, isEditing: true, isLoading: false });
        } catch (error) {
            console.error(`Error al consultar marca ID ${marca_id}:`, error);
            set({ isLoading: false });
        }
    },

    EnviarMarca: async () => {
        const { form_marca, isEditing, ObtenerMarcas, reset } = get();
        set({ isLoading: true });
        try {
            const payload: Marca = {
                fecha_creacion: form_marca.fecha_creacion,
                fecha_modificacion: form_marca.fecha_modificacion,
                marca_id: form_marca.marca_id,
                nombre: form_marca.nombre
            };

            if (isEditing) {
                await marcaService.modificarMarca(payload);
            } else {
                await marcaService.crearMarca(payload);
            }

            reset();

            await ObtenerMarcas();

        } catch (error) {
            console.error(isEditing == true ? "Error al modificar la marca:" : "Error al crear la marca", error);
            set({ isLoading: false });
        }
    },

    EliminarMarca: async (id: number) => {
        set({ isLoading: true });
        try {
            await marcaService.eliminarMarca(id);
            await get().ObtenerMarcas();
        } catch (error) {
            console.error(`Error al eliminar marca ID ${id}:`, error);
            set({ isLoading: false });
        }
    },

    reset: () =>
        set({
            form_marca: initialMarca(),
            isEditing: false,
            isLoading: false,
        }),
}));