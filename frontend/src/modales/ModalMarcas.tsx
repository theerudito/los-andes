import {
    Pencil,
    Trash2,
    Save,
    X,
    EraserIcon,
    Tag
} from "lucide-react";

import { useModal } from "../store/useModal.ts";
import { ModalLista } from "../helpers/ModalLista.ts";
import {useMarcas} from "../store/useMarcas.ts";
import React, {useEffect} from "react";

export default function ModalMarcas() {
    const { modalName, CloseModal } = useModal((state) => state);
    const { ObtenerMarcas, ObtenerMarca, EnviarMarca, EliminarMarca, form_marca, listar_marca} = useMarcas((state) => state);

    useEffect(() => {
        ObtenerMarcas();
    }, []);

    if (modalName !== ModalLista.modal_marca) return null;

    const handleChangeInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name } = e.target;
        const value = e.target.value;
        useMarcas.setState((state) => {
            return {
                form_marca: {
                    ...state.form_marca,
                    [name]: value.toUpperCase()
                },
            };
        });
    };

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-[60] p-4 transition-all duration-300">
            <div className="w-full max-w-xl bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-blue-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <Tag size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">Gestión de Marcas</h2>
                    </div>
                    <button
                        type="button"
                        className="cursor-pointer hover:bg-white/20 transition-all rounded-full p-1.5 active:scale-95"
                        onClick={CloseModal}
                    >
                        <X size={18} />
                    </button>
                </div>

                <div className="p-5 bg-slate-50/50 flex flex-col gap-4">
                    <div className="flex shadow-sm rounded-lg overflow-hidden border border-slate-300 focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500 transition-all bg-white h-10">
                        <input
                            type="text"
                            name="nombre"
                            value={form_marca.nombre}
                            onChange={handleChangeInput}
                            placeholder="Ingrese la nueva marca"
                            className="flex-1 px-3 py-2 text-xs outline-none text-slate-700 placeholder-slate-400 uppercase font-medium"
                        />

                        <div className="flex shrink-0 border-l border-slate-200">
                            <button
                                type="button"
                                className="cursor-pointer bg-orange-50 text-orange-600 hover:bg-orange-100 px-3.5 flex items-center justify-center transition-all border-r border-orange-200/60 active:scale-95"
                                title="Limpiar"
                            >
                                <EraserIcon size={16} />
                            </button>

                            <button
                                onClick={EnviarMarca}
                                type="button"
                                className="cursor-pointer bg-blue-600 hover:bg-blue-700 text-white px-4 flex items-center justify-center transition-all gap-1 text-xs font-semibold active:scale-95"
                                title="Guardar"
                            >
                                <Save size={15} />
                            </button>
                        </div>
                    </div>
                </div>

                <div className="px-5 pb-5 bg-slate-50/50 flex-1">
                    <div className="border border-slate-200 rounded-xl overflow-hidden bg-white shadow-sm">
                        <div className="max-h-60 overflow-y-auto scrollbar-thin scrollbar-thumb-slate-300">
                            <table className="w-full text-xs border-collapse table-fixed">
                                <thead className="bg-slate-100 text-slate-700 border-b border-slate-200 sticky top-0 z-10">
                                <tr>
                                    <th className="w-[15%] py-2.5 px-4 font-bold text-left uppercase tracking-wider text-[10px]">
                                        ID
                                    </th>
                                    <th className="w-[60%] py-2.5 px-4 font-bold text-left uppercase tracking-wider text-[10px]">
                                        Nombre
                                    </th>
                                    <th className="w-[25%] py-2.5 px-4 font-bold text-center uppercase tracking-wider text-[10px]">
                                        Acciones
                                    </th>
                                </tr>
                                </thead>

                                <tbody className="divide-y divide-slate-100">
                                {listar_marca && listar_marca.length > 0 ? (
                                    listar_marca.map((item) => (
                                        <tr key={item.marca_id} className="hover:bg-slate-50 transition-colors">
                                            <td className="py-2.5 px-4 font-semibold text-slate-900 truncate">
                                                #{item.marca_id}
                                            </td>
                                            <td className="py-2.5 px-4 text-slate-700 font-medium truncate" title={`${item.nombre} `}>
                                                {item.nombre}
                                            </td>
                                            <td className="py-2.5 px-4 text-center">
                                                <div className="flex items-center justify-center gap-1.5">
                                                    <button
                                                        onClick={() => ObtenerMarca(item.marca_id)}
                                                        type="button"
                                                        className="p-1 text-blue-600 hover:bg-blue-50 rounded transition-colors"
                                                        title="Editar"
                                                    >
                                                        <Pencil className="w-3.5 h-3.5" />
                                                    </button>
                                                    <button
                                                        onClick={() => EliminarMarca(item.marca_id)}
                                                        type="button"
                                                        className="p-1 text-rose-600 hover:bg-rose-50 rounded transition-colors"
                                                        title="Eliminar"
                                                    >
                                                        <Trash2 className="w-3.5 h-3.5" />
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={3} className="py-6 text-center text-slate-400 font-medium">
                                            No hay registros disponibles
                                        </td>
                                    </tr>
                                )}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    );
}