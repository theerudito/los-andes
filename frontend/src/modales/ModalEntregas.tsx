import React, { useState } from 'react';
import {
    X,
    Save,
    PackageCheck,
    Laptop,
    FileText,
    CheckCircle2,
    MessageSquare,
    Wrench
} from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts';

export default function ModalEntregas(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);

    if (modalName !== ModalLista.modal_entrega) return null;

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-[60] p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-gradient-to-r from-amber-600 to-orange-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <PackageCheck size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Registrar Entrega de Equipo
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

                <div className="p-5 bg-slate-50/50 flex flex-col gap-4 max-h-[80vh] overflow-y-auto">

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Laptop size={14} className="text-amber-600" /> ID Equipo
                            </label>
                            <input
                                type="number"
                                name="equipo_id"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-amber-500 focus:ring-1 focus:ring-amber-500 transition-all font-mono font-bold shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <CheckCircle2 size={14} className="text-amber-600" /> Conformidad Cliente
                            </label>
                            <select
                                name="conformidad_cliente"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-amber-500 focus:ring-1 focus:ring-amber-500 transition-all font-semibold uppercase cursor-pointer shadow-sm"
                            >
                                <option value={1}>CONFORME</option>
                                <option value={0}>NO CONFORME</option>
                            </select>
                        </div>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <Wrench size={14} className="text-amber-600" /> Trabajos Realizados
                        </label>
                        <textarea
                            name="trabajos_realizados"
                            rows={2}
                            placeholder="Ej: SE FORMATEO Y MANTENIMIENTO GENERAL"
                            required
                            className="w-full p-2.5 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-amber-500 focus:ring-1 focus:ring-amber-500 transition-all font-medium resize-none shadow-sm"
                        />
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <FileText size={14} className="text-amber-600" /> Estado Final del Equipo
                        </label>
                        <textarea
                            name="estado_final_equipo"
                            rows={2}
                            placeholder="Ej: BUEN ESTADO AUN HASTA NUEVO DIAGNOSTICO"
                            required
                            className="w-full p-2.5 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-amber-500 focus:ring-1 focus:ring-amber-500 transition-all font-medium resize-none shadow-sm"
                        />
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <MessageSquare size={14} className="text-amber-600" /> Observaciones Adicionales
                        </label>
                        <textarea
                            name="observaciones"
                            rows={2}
                            placeholder="Ej: LISTO Y RECIBIDO POR EL CLIENTE"
                            className="w-full p-2.5 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-amber-500 focus:ring-1 focus:ring-amber-500 transition-all font-medium resize-none shadow-sm"
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
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-amber-600 to-orange-600 hover:from-amber-700 hover:to-orange-700 rounded-lg shadow-sm transition-all active:scale-95"
                        >
                            <Save size={15} />
                            <span>Guardar Entrega</span>
                        </button>
                    </div>

                </div>

            </div>
        </div>
    );
}