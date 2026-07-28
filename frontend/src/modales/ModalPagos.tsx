import {
    X,
    Save,
    CreditCard,
    DollarSign,
    Laptop,
    Receipt,
    Info,
    CheckCircle2,
    AlertCircle
} from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts';
import { usePagos } from "../store/usePagos.ts";

export default function ModalPagos(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);
    const { form_pagos, EnviarPago, equipo_id } = usePagos((state) => state);

    const costoTotal = Number(form_pagos.costo_total || 0);
    const abonoTotal = Number(form_pagos.abono || 0);
    const saldoPendiente = Math.max(0, costoTotal - abonoTotal);
    const estaPagado = costoTotal > 0 && saldoPendiente === 0;

    const handleChangeInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;

        usePagos.setState((state) => ({
            form_pagos: {
                ...state.form_pagos,
                [name]: value
            }
        }));
    };

    function Clear() {
        CloseModal();
        usePagos.setState({
            form_pagos: {
                equipo_id: 0,
                costo_total: 0,
                abono: 0,
                cuenta_id: 0,
            },
            isEditing: false,
        });
    }

    if (modalName !== ModalLista.modal_pago) return null;

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-[60] p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-gradient-to-r from-emerald-600 to-teal-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <CreditCard size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Gestionar Estado de Cuenta
                        </h2>
                    </div>
                    <button
                        type="button"
                        className="cursor-pointer hover:bg-white/20 transition-all rounded-full p-1.5 active:scale-95"
                        onClick={Clear}
                    >
                        <X size={18} />
                    </button>
                </div>

                <div className="p-5 bg-slate-50/50 flex flex-col gap-4">

                    <div className="flex items-start gap-2.5 p-3 bg-blue-50/80 border border-blue-100 rounded-lg text-blue-800 text-xs">
                        <Info size={16} className="text-blue-600 shrink-0 mt-0.5" />
                        <p className="leading-relaxed">
                            Ingresa o actualiza el <strong className="font-semibold">Costo Total</strong> del trabajo y el <strong className="font-semibold">Abono Acumulado</strong> recibido hasta la fecha.
                        </p>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <Laptop size={14} className="text-emerald-600" /> Equipo Asociado
                        </label>
                        <input
                            value={`Equipo #${equipo_id}`}
                            type="text"
                            disabled
                            className="w-full h-10 px-3 bg-slate-100 border border-slate-300 rounded-lg text-xs text-slate-800 font-mono font-bold shadow-sm cursor-not-allowed"
                        />
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
                                value={form_pagos.costo_total}
                                onChange={handleChangeInput}
                                required
                                placeholder="0.00"
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all font-bold shadow-sm"
                            />
                            <span className="text-[10px] text-slate-400 mt-1 block">Presupuesto total del servicio</span>
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Receipt size={14} className="text-emerald-600" /> Abono Acumulado ($)
                            </label>
                            <input
                                type="number"
                                step="0.01"
                                min="0"
                                name="abono"
                                value={form_pagos.abono}
                                onChange={handleChangeInput}
                                required
                                placeholder="0.00"
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 transition-all font-bold shadow-sm text-emerald-600"
                            />
                            <span className="text-[10px] text-slate-400 mt-1 block">Total entregado por el cliente</span>
                        </div>
                    </div>

                    <div className="p-3.5 bg-white rounded-xl border border-slate-200/80 shadow-sm flex flex-col gap-2">
                        <div className="text-[11px] font-semibold text-slate-500 uppercase tracking-wider">
                            Resumen de la Cuenta
                        </div>
                        <div className="grid grid-cols-3 gap-2 text-center pt-1 border-t border-slate-100">
                            <div className="flex flex-col">
                                <span className="text-[10px] text-slate-400 font-medium">Costo</span>
                                <span className="text-xs font-bold text-slate-700 font-mono">${costoTotal.toFixed(2)}</span>
                            </div>
                            <div className="flex flex-col">
                                <span className="text-[10px] text-slate-400 font-medium">Abonado</span>
                                <span className="text-xs font-bold text-emerald-600 font-mono">${abonoTotal.toFixed(2)}</span>
                            </div>
                            <div className="flex flex-col">
                                <span className="text-[10px] text-slate-400 font-medium">Pendiente</span>
                                <span className={`text-xs font-bold font-mono ${saldoPendiente <= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
                                    ${saldoPendiente.toFixed(2)}
                                </span>
                            </div>
                        </div>

                        <div className="mt-1 pt-2 border-t border-slate-100 flex items-center justify-between">
                            <span className="text-xs text-slate-500">Estado de Cuenta:</span>
                            {estaPagado ? (
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 text-[11px] font-bold text-emerald-700 bg-emerald-50 border border-emerald-200 rounded-full">
                                    <CheckCircle2 size={12} /> PAGADO
                                </span>
                            ) : (
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 text-[11px] font-bold text-amber-700 bg-amber-50 border border-amber-200 rounded-full">
                                    <AlertCircle size={12} /> PENDIENTE
                                </span>
                            )}
                        </div>
                    </div>

                    <div className="pt-2 border-t border-slate-200/80 flex items-center justify-end gap-2 shrink-0">
                        <button
                            type="button"
                            onClick={Clear}
                            className="cursor-pointer px-4 py-2 text-xs font-semibold text-slate-600 bg-slate-200 hover:bg-slate-300 rounded-lg transition-all active:scale-95"
                        >
                            Cancelar
                        </button>

                        <button
                            onClick={EnviarPago}
                            type="button"
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 rounded-lg shadow-sm transition-all active:scale-95"
                        >
                            <Save size={15} />
                            <span>Guardar Cambios</span>
                        </button>
                    </div>

                </div>

            </div>
        </div>
    );
}