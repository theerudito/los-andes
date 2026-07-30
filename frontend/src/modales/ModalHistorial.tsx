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
import { useHistorial } from "../store/useHistorial.ts";
import React from "react";
import { toast } from "sonner";

const estadosOpciones = [
    { estado_id: 1, nombre: "Recibido" },
    { estado_id: 2, nombre: "En diagnóstico" },
    { estado_id: 3, nombre: "Esperando repuestos" },
    { estado_id: 4, nombre: "En reparación" },
    { estado_id: 5, nombre: "Listo para entrega" },
    { estado_id: 6, nombre: "Entregado" },
    { estado_id: 7, nombre: "Cancelado" }
];

export default function ModalHistorial(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);
    const { form_historial, EnviarHistorial, equipo_id } = useHistorial((state) => state);

    function Enviar(e?: React.SubmitEvent) {
        if (e) e.preventDefault();

        if (!form_historial.estado_id || form_historial.estado_id <= 0) {
            toast.error("Debe seleccionar un estado válido para el equipo");
            return;
        }

        if (!form_historial.observaciones_tecnicas || form_historial.observaciones_tecnicas.trim() === "") {
            toast.error("La observación técnica es requerida");
            return;
        }

        EnviarHistorial();
    }

    const handleChangeSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const { name, value } = e.target;
        useHistorial.setState((state) => ({
            form_historial: {
                ...state.form_historial,
                [name]: Number(value)
            },
        }));
    };

    const handleChangeTextArea = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
        const { name, value } = e.target;
        useHistorial.setState((state) => ({
            form_historial: {
                ...state.form_historial,
                [name]: value,
            },
        }));
    };

    function Clear() {
        CloseModal();
        useHistorial.setState({
            form_historial: {
                equipo_id: 0,
                estado_id: 0,
                observaciones_tecnicas: "",
                historial_id: 0
            },
            isEditing: false,
        });
    }

    if (modalName !== ModalLista.modal_historial) return null;

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-[60] p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-gradient-to-r from-yellow-600 to-yellow-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <History size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Registrar Historial Técnico
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

                <form onSubmit={Enviar} className="p-5 bg-slate-50/50 flex flex-col gap-4">

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Laptop size={14} className="text-blue-600" /> ID Equipo
                            </label>
                            <input
                                value={equipo_id}
                                type="number"
                                name="equipo_id"
                                required
                                disabled
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-mono font-bold shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Tag size={14} className="text-blue-600" /> Estado del Equipo
                            </label>
                            <select
                                value={form_historial.estado_id || 2}
                                onChange={handleChangeSelect}
                                name="estado_id"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-semibold uppercase cursor-pointer shadow-sm"
                            >
                                {estadosOpciones.map((e) => (
                                    <option key={e.estado_id} value={e.estado_id}>{e.nombre}</option>
                                ))}
                            </select>
                        </div>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <FileText size={14} className="text-blue-600" /> Observaciones Técnicas
                        </label>
                        <textarea
                            value={form_historial.observaciones_tecnicas}
                            onChange={handleChangeTextArea}
                            name="observaciones_tecnicas"
                            rows={3}
                            placeholder="Ej: completado todo ok"
                            className="w-full p-2.5 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium resize-none shadow-sm"
                        />
                    </div>

                    <div className="pt-3 border-t border-slate-200/80 flex items-center justify-end gap-2 shrink-0">
                        <button
                            type="button"
                            onClick={Clear}
                            className="cursor-pointer px-4 py-2 text-xs font-semibold text-slate-600 bg-slate-200 hover:bg-slate-300 rounded-lg transition-all active:scale-95"
                        >
                            Cancelar
                        </button>

                        <button
                            type="submit"
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-yellow-600 to-yellow-600 hover:from-yellow-700 hover:to-yellow-700 rounded-lg shadow-sm transition-all active:scale-95"
                        >
                            <Save size={15} />
                            <span>Guardar Historial</span>
                        </button>
                    </div>

                </form>

            </div>
        </div>
    );
}