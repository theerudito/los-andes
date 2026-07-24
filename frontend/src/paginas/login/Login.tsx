import { User, Lock, Eye, EyeOff, KeyRound, ArrowLeft } from 'lucide-react';
import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useUsuarios } from "../../store/useUsuarios.ts";
import type { UsuarioLogin } from "../../modelos/usuarios.ts";

const inicial: UsuarioLogin = {
    identificacion: "",
    password: ""
};

export function Login(): React.ReactElement {
    const navigate = useNavigate();

    const { LoginUsuario, ResetUsuario, isLogin } = useUsuarios((state) => state);

    const [isResetView, setIsResetView] = useState<boolean>(false);
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const [datos, setDatos] = useState<UsuarioLogin>(inicial);

    useEffect(() => {
        if (isLogin) {
            navigate('/', { replace: true });
        }
    }, [isLogin, navigate]);

    const toggleView = (): void => {
        setIsResetView(!isResetView);
        setShowPassword(false);
        setDatos(inicial);
    };

    const handleChangeInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setDatos((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const handleLoginSubmit = async () => {
        await LoginUsuario(datos);
    };

    const handleResetSubmit = async () => {
        await ResetUsuario(datos);
    };

    return (
        <div className="min-h-screen w-screen flex items-center justify-center bg-gray-100 p-4">
            <div className="w-full max-w-md bg-white rounded-2xl shadow-xl border border-gray-100 p-8">

                {!isResetView ? (
                    <div className="space-y-6">
                        <div className="text-center">
                            <div className="w-14 h-14 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center mx-auto mb-3">
                                <Lock className="w-7 h-7" />
                            </div>
                            <h2 className="text-2xl font-bold text-gray-800">Iniciar Sesión</h2>
                            <p className="text-sm text-gray-500 mt-1">
                                Ingresa a tu cuenta para acceder al sistema
                            </p>
                        </div>

                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Identificación
                                </label>
                                <div className="relative">
                                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                                        <User className="w-5 h-5" />
                                    </div>
                                    <input
                                        type="text"
                                        name="identificacion"
                                        required
                                        value={datos.identificacion}
                                        onChange={handleChangeInput}
                                        placeholder="Identificacion o Email"
                                        className="w-full pl-10 pr-4 py-2.5 bg-gray-50 border border-gray-300 rounded-lg text-sm text-gray-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Contraseña
                                </label>
                                <div className="relative">
                                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                                        <Lock className="w-5 h-5" />
                                    </div>
                                    <input
                                        type={showPassword ? 'text' : 'password'}
                                        name="password"
                                        required
                                        value={datos.password}
                                        onChange={handleChangeInput}
                                        placeholder="••••••••"
                                        className="w-full pl-10 pr-10 py-2.5 bg-gray-50 border border-gray-300 rounded-lg text-sm text-gray-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                                    />
                                    <button
                                        type="button"
                                        onClick={() => setShowPassword(!showPassword)}
                                        className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 focus:outline-none"
                                        tabIndex={-1}
                                    >
                                        {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                                    </button>
                                </div>
                            </div>

                            <div className="flex justify-end">
                                <button
                                    type="button"
                                    onClick={toggleView}
                                    className="text-sm font-medium text-blue-600 hover:text-blue-700 hover:underline focus:outline-none cursor-pointer"
                                >
                                    ¿Olvidaste tu contraseña?
                                </button>
                            </div>
                        </div>

                        <button
                            type="button"
                            onClick={handleLoginSubmit}
                            className="w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg shadow-md hover:shadow-lg transition-all focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer"
                        >
                            Login
                        </button>
                    </div>
                ) : (
                    <div className="space-y-6">
                        <div className="text-center">
                            <div className="w-14 h-14 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center mx-auto mb-3">
                                <KeyRound className="w-7 h-7" />
                            </div>
                            <h2 className="text-2xl font-bold text-gray-800">Restablecer Contraseña</h2>
                            <p className="text-sm text-gray-500 mt-1">
                                Ingresa tus datos para cambiar tu clave de acceso
                            </p>
                        </div>

                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Identificación
                                </label>
                                <div className="relative">
                                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                                        <User className="w-5 h-5" />
                                    </div>
                                    <input
                                        type="text"
                                        name="identificacion"
                                        required
                                        value={datos.identificacion}
                                        onChange={handleChangeInput}
                                        placeholder="Identificacion o Email"
                                        className="w-full pl-10 pr-4 py-2.5 bg-gray-50 border border-gray-300 rounded-lg text-sm text-gray-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Nueva Contraseña
                                </label>
                                <div className="relative">
                                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                                        <Lock className="w-5 h-5" />
                                    </div>
                                    <input
                                        type={showPassword ? 'text' : 'password'}
                                        name="password"
                                        required
                                        value={datos.password}
                                        onChange={handleChangeInput}
                                        placeholder="••••••••"
                                        className="w-full pl-10 pr-10 py-2.5 bg-gray-50 border border-gray-300 rounded-lg text-sm text-gray-900 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                                    />
                                    <button
                                        type="button"
                                        onClick={() => setShowPassword(!showPassword)}
                                        className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 focus:outline-none"
                                        tabIndex={-1}
                                    >
                                        {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                                    </button>
                                </div>
                            </div>
                        </div>

                        <button
                            type="button"
                            onClick={handleResetSubmit}
                            className="w-full py-3 px-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg shadow-md hover:shadow-lg transition-all focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer"
                        >
                            Resetear Contraseña
                        </button>

                        <div className="text-center pt-2">
                            <button
                                type="button"
                                onClick={toggleView}
                                className="inline-flex items-center gap-2 text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors focus:outline-none cursor-pointer"
                            >
                                <ArrowLeft className="w-4 h-4" />
                                <span>Volver al inicio de sesión</span>
                            </button>
                        </div>
                    </div>
                )}

            </div>
        </div>
    );
}