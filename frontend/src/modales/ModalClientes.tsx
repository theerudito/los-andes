import React, { useState } from 'react';
import { X, Save, User, CreditCard, Phone, Mail, MapPin } from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts'; // O la ruta correspondiente de tu helper

export interface FormClienteBody {
    identificacion: string;
    nombres: string;
    apellidos: string;
    telefono: string;
    email: string;
    direccion: string;
}

const formInicial: FormClienteBody = {
    identificacion: '',
    nombres: '',
    apellidos: '',
    telefono: '',
    email: '',
    direccion: ''
};

export default function ModalCliente(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);
    const [formData, setFormData] = useState<FormClienteBody>(formInicial);

    if (modalName !== ModalLista.modal_cliente) return null;

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value
        }));
    };

    const handleReset = () => {
        setFormData(formInicial);
    };

    const handleGuardar = (e: React.FormEvent) => {
        e.preventDefault();

        // Body JSON exacto listo para enviar al backend en Go
        const payload: FormClienteBody = {
            identificacion: formData.identificacion.trim(),
            nombres: formData.nombres.trim(),
            apellidos: formData.apellidos.trim(),
            telefono: formData.telefono.trim(),
            email: formData.email.trim(),
            direccion: formData.direccion.trim()
        };

        console.log("Payload enviado a Backend Go (REQ):", payload);

        // Aquí ejecutas la función para guardar (API / Zustand Store)

        handleReset();
        CloseModal();
    };

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-50 p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                {/* Header del Modal */}
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

                <form onSubmit={handleGuardar} className="p-5 bg-slate-50/50 flex flex-col gap-4">

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <CreditCard size={14} className="text-blue-600" /> Identificación / Cédula / RUC
                        </label>
                        <input
                            type="text"
                            name="identificacion"
                            value={formData.identificacion}
                            onChange={handleChange}
                            placeholder="Ej: 1721457494001"
                            required
                            className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium"
                        />
                    </div>

                    {/* Nombres y Apellidos en 2 Columnas */}
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <User size={14} className="text-blue-600" /> Nombres
                            </label>
                            <input
                                type="text"
                                name="nombres"
                                value={formData.nombres}
                                onChange={handleChange}
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
                                type="text"
                                name="apellidos"
                                value={formData.apellidos}
                                onChange={handleChange}
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
                                type="text"
                                name="telefono"
                                value={formData.telefono}
                                onChange={handleChange}
                                placeholder="Ej: 0960806054"
                                className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Mail size={14} className="text-blue-600" /> Correo Electrónico
                            </label>
                            <input
                                type="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
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
                            name="direccion"
                            rows={2}
                            value={formData.direccion}
                            onChange={handleChange}
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
                            type="submit"
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 rounded-lg shadow-sm transition-all active:scale-95"
                        >
                            <Save size={15} />
                            <span>Guardar</span>
                        </button>
                    </div>

                </form>

            </div>
        </div>
    );
}