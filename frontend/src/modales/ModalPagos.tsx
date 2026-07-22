import React, { useState } from 'react';
import {
    X,
    Save,
    CreditCard,
    DollarSign,
    User,
    Laptop
} from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts';

export default function ModalPagos(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);

    if (modalName !== ModalLista.modal_pago) return null;

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-[60] p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-gradient-to-r from-emerald-600 to-teal-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <CreditCard size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Registrar / Actualizar Pago
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
                                <Laptop size={14} className="text-emerald-600" /> ID Equipo
                            </label>
                            <input
                                type="number"
                                name="equipo_id"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all font-mono font-bold shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <User size={14} className="text-emerald-600" /> ID Usuario (Cobrador)
                            </label>
                            <input
                                type="number"
                                name="usuario_id"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all font-mono font-bold shadow-sm"
                            />
                        </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <DollarSign size={14} className="text-emerald-600" /> Costo Total ($)
                            </label>
                            <input
                                type="number"
                                step="0.01"
                                min="0"
                                name="costo_total"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all font-bold shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <DollarSign size={14} className="text-emerald-600" /> Abono ($)
                            </label>
                            <input
                                type="number"
                                step="0.01"
                                min="0"
                                name="abono"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all font-bold shadow-sm text-emerald-600"
                            />
                        </div>
                    </div>

                    <div className="p-3 bg-white rounded-lg border border-slate-200 flex justify-between items-center shadow-sm">
                        <span className="text-xs font-medium text-slate-600">Saldo Pendiente Estimado:</span>
                        <span className={`text-sm font-bold font-mono ${modalName <= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
                            ${modalName.toFixed(2)}
                        </span>
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
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 rounded-lg shadow-sm transition-all active:scale-95"
                        >
                            <Save size={15} />
                            <span>Guardar Pago</span>
                        </button>
                    </div>

                </div>

            </div>
        </div>
    );
}