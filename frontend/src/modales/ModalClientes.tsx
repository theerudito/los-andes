import React from "react";
import { X, Save, User, CreditCard, Phone, Mail, MapPin } from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts';
import {useClientes} from "../store/useClientes.ts";

export default function ModalCliente(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);

    const { form_cliente, EnviarCliente } = useClientes((state) => state);

    const handleChangeInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name } = e.target;

        const value = e.target.value;

        useClientes.setState((state) => {
            return {
                form_cliente: {
                    ...state.form_cliente,
                    [name]: value
                },
            };
        });
    };

    const handleChangeTextArea = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
        const { name, value } = e.target;
        useClientes.setState((state) => ({
            form_cliente: {
                ...state.form_cliente,
                [name]: value,
            },
        }));
    };

    if (modalName !== ModalLista.modal_cliente) return null;

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-50 p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <User size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Registrar Nuevo Cliente
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

                <div  className="p-5 bg-slate-50/50 flex flex-col gap-4">

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <CreditCard size={14} className="text-blue-600" /> Identificación / Cédula / RUC
                        </label>
                        <input
                            value={form_cliente.identificacion}
                            onChange={handleChangeInput}
                            type="text"
                            name="identificacion"
                            placeholder="Ej: 1721457494001"
                            required
                            className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium"
                        />
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <User size={14} className="text-blue-600" /> Nombres
                            </label>
                            <input
                                value={form_cliente.nombres}
                                onChange={handleChangeInput}
                                type="text"
                                name="nombres"
                                placeholder="Ej: Juan Carlos"
                                required
                                className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <User size={14} className="text-blue-600" /> Apellidos
                            </label>
                            <input
                                value={form_cliente.apellidos}
                                onChange={handleChangeInput}
                                type="text"
                                name="apellidos"
                                placeholder="Ej: Pérez Gómez"
                                required
                                className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium"
                            />
                        </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Phone size={14} className="text-blue-600" /> Teléfono
                            </label>
                            <input
                                value={form_cliente.telefono}
                                onChange={handleChangeInput}
                                type="text"
                                name="telefono"
                                placeholder="Ej: 0960806054"
                                className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Mail size={14} className="text-blue-600" /> Correo Electrónico
                            </label>
                            <input
                                value={form_cliente.email}
                                onChange={handleChangeInput}
                                type="email"
                                name="email"
                                placeholder="Ej: cliente@correo.com"
                                className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium"
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <MapPin size={14} className="text-blue-600" /> Dirección
                        </label>
                        <textarea
                            value={form_cliente.direccion}
                            onChange={handleChangeTextArea}
                            name="direccion"
                            rows={2}
                            placeholder="Ej: Libertad del Toachi"
                            className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium resize-none"
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
                            onClick={EnviarCliente}
                            type="submit"
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 rounded-lg shadow-sm transition-all active:scale-95"
                        >
                            <Save size={15} />
                            <span>Guardar</span>
                        </button>
                    </div>

                </div>

            </div>
        </div>
    );
}