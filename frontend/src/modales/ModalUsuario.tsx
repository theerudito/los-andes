import React, { useState } from 'react';
import {
    X,
    Save,
    User,
    CreditCard,
    Mail,
    Lock,
    ShieldCheck,
    UserCheck,
    Eye,
    EyeOff
} from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts';


const rolesOpciones = [
    { rol_id: 1, nombre: "ADMINISTRADOR" },
    { rol_id: 2, nombre: "TÉCNICO" },
    { rol_id: 3, nombre: "RECEPCIONISTA" },
    { rol_id: 4, nombre: "CLIENTE" },
];

export default function ModalUsuario(): React.ReactElement | null {
    const { modalName, CloseModal } = useModal((state) => state);
    const [verPassword, setVerPassword] = useState<boolean>(false);

    if (modalName !== ModalLista.modal_usuario) return null;
    
    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-[60] p-4 transition-all duration-300">
            <div className="w-full max-w-lg bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

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

                <div className="p-5 bg-slate-50/50 flex flex-col gap-4">

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <CreditCard size={14} className="text-blue-600" /> Identificación (Cédula / RUC)
                        </label>
                        <input
                            type="text"
                            name="identificacion"
                            maxLength={13}
                            placeholder="Ej: 1721457495"
                            required
                            className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-mono font-bold tracking-wider shadow-sm"
                        />
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <User size={14} className="text-blue-600" /> Nombres
                            </label>
                            <input
                                type="text"
                                name="nombres"
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
                                placeholder="Ej: LOOR"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Mail size={14} className="text-blue-600" /> Correo Electrónico
                            </label>
                            <input
                                type="email"
                                name="email"
                                placeholder="Ej: erudito.tv@gmail.com"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Lock size={14} className="text-blue-600" /> Contraseña
                            </label>
                            <div className="relative flex items-center">
                                <input
                                    type={verPassword ? "text" : "password"}
                                    name="password"
                                    placeholder="••••••••"
                                    required
                                    className="w-full h-10 pl-3 pr-10 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                                />
                                <button
                                    type="button"
                                    onClick={() => setVerPassword(!verPassword)}
                                    className="absolute right-3 text-slate-400 hover:text-slate-600 cursor-pointer transition-colors p-1"
                                    title={verPassword ? "Ocultar contraseña" : "Mostrar contraseña"}
                                >
                                    {verPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                                </button>
                            </div>
                        </div>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <ShieldCheck size={14} className="text-blue-600" /> Rol de Usuario
                        </label>
                        <select
                            name="rol_id"
                            required
                            className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-semibold uppercase cursor-pointer shadow-sm"
                        >
                                <option >
                                    Rol
                                </option>
                        </select>
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
                            <span>Guardar Usuario</span>
                        </button>
                    </div>

                </div>

            </div>
        </div>
    );
}