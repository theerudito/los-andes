import React, { useState } from 'react';
import {
    X,
    Save,
    User,
    CreditCard,
    Mail,
    Lock,
    ShieldCheck,
    UserCheck
} from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts';

// Estructura exacto del JSON request enviado a Go
export interface FormUsuarioBody {
    identificacion: string;
    nombres: string;
    apellidos: string;
    email: string;
    password: string;
    rol_id: number;
}

// Opciones de prueba para el select de Roles
const rolesOpciones = [
    { rol_id: 1, nombre: "ADMINISTRADOR" },
    { rol_id: 2, nombre: "TÉCNICO" },
    { rol_id: 3, nombre: "RECEPCIONISTA" },
    { rol_id: 4, nombre: "CLIENTE" },
];

const formInicial: FormUsuarioBody = {
    identificacion: '',
    nombres: '',
    apellidos: '',
    email: '',
    password: '',
    rol_id: 3, // Valor por defecto según tu JSON
};

export default function ModalUsuario(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);
    const [formData, setFormData] = useState<FormUsuarioBody>(formInicial);

    if (modalName !== ModalLista.modal_usuario) return null;

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
    ) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: name === 'rol_id' ? Number(value) : value,
        }));
    };

    const handleGuardar = (e: React.FormEvent) => {
        e.preventDefault();

        // Payload formateado 1:1 para el Backend Go
        const payload: FormUsuarioBody = {
            identificacion: formData.identificacion.trim(),
            nombres: formData.nombres.toUpperCase().trim(),
            apellidos: formData.apellidos.toUpperCase().trim(),
            email: formData.email.trim(),
            password: formData.password,
            rol_id: Number(formData.rol_id),
        };

        console.log("Payload enviado a Backend Go (REQ):", payload);

        setFormData(formInicial);
        CloseModal();
    };

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-[60] p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                {/* Header del Modal */}
                <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <UserCheck size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Registrar Nuevo Usuario
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

                {/* Formulario */}
                <form onSubmit={handleGuardar} className="p-5 bg-slate-50/50 flex flex-col gap-4">

                    {/* Identificación / Cédula */}
                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <CreditCard size={14} className="text-blue-600" /> Identificación (Cédula / RUC)
                        </label>
                        <input
                            type="text"
                            name="identificacion"
                            maxLength={13}
                            value={formData.identificacion}
                            onChange={handleChange}
                            placeholder="Ej: 1721457495"
                            required
                            className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-mono font-bold tracking-wider shadow-sm"
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
                                placeholder="Ej: JORGE"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
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
                                placeholder="Ej: LOOR"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>
                    </div>

                    {/* Correo Electrónico y Contraseña en 2 Columnas */}
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Mail size={14} className="text-blue-600" /> Correo Electrónico
                            </label>
                            <input
                                type="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                                placeholder="Ej: erudito.tv@gmail.com"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Lock size={14} className="text-blue-600" /> Contraseña
                            </label>
                            <input
                                type="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                placeholder="••••••••"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>
                    </div>

                    {/* Rol (Select) */}
                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <ShieldCheck size={14} className="text-blue-600" /> Rol de Usuario
                        </label>
                        <select
                            name="rol_id"
                            value={formData.rol_id}
                            onChange={handleChange}
                            required
                            className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-semibold uppercase cursor-pointer shadow-sm"
                        >
                            {rolesOpciones.map((r) => (
                                <option key={r.rol_id} value={r.rol_id}>
                                    {r.nombre}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Footer con Botones Cancelar y Guardar */}
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
                            <span>Guardar Usuario</span>
                        </button>
                    </div>

                </form>

            </div>
        </div>
    );
}