import React, { useState } from 'react';
import {
    X,
    Save,
    Tag,
    FileText,
    Laptop,
    History
} from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts';

const estadosOpciones = [
    { estado_id: 1, nombre: "EN REVISIÓN" },
    { estado_id: 2, nombre: "DIAGNÓSTICO" },
    { estado_id: 3, nombre: "EN PROCESO" },
    { estado_id: 4, nombre: "REPARADO" },
    { estado_id: 5, nombre: "COMPLETADO" },
];

export default function ModalHistorial(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);

    if (modalName !== ModalLista.modal_historial) return null;

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-[60] p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <History size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Registrar Historial Técnico
                        </h2>
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

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Laptop size={14} className="text-blue-600" /> ID Equipo
                            </label>
                            <input
                                type="number"
                                name="equipo_id"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-mono font-bold shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Tag size={14} className="text-blue-600" /> Estado del Equipo
                            </label>
                            <select
                                name="estado_id"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-semibold uppercase cursor-pointer shadow-sm"
                            >
                                {estadosOpciones.map((e) => (
                                    <option key={e.estado_id} value={e.estado_id}>
                                        {e.nombre}
                                    </option>
                                ))}
                            </select>
                        </div>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <FileText size={14} className="text-blue-600" /> Observaciones Técnicas
                        </label>
                        <textarea
                            name="observaciones_tecnicas"
                            rows={3}
                            placeholder="Ej: completado todo ok"
                            required
                            className="w-full p-2.5 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium resize-none shadow-sm"
                        />
                    </div>

                    <div className="pt-3 border-t border-slate-200/80 flex items-center justify-end gap-2 shrink-0">
                        <button
                            type="button"
                            onClick={CloseModal}
                            className="cursor-pointer px-4 py-2 text-xs font-semibold text-slate-600 bg-slate-200 hover:bg-slate-300 rounded-lg transition-all active:scale-95"
                        >
                            Cancelar
                        </button>

                        <button
                            type="submit"
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 rounded-lg shadow-sm transition-all active:scale-95"
                        >
                            <Save size={15} />
                            <span>Guardar Historial</span>
                        </button>
                    </div>

                </div>

            </div>
        </div>
    );
}