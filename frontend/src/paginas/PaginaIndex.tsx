import React, { useEffect, useState } from 'react';
import {
    Laptop,
    Users,
    CheckCircle2,
    AlertTriangle,
    ShieldCheck,
    TrendingUp,
    Loader2
} from 'lucide-react';
import api from "../helpers/fetching/axios.ts";

interface DashboardData {
    equipos: number;
    entregas: number;
    clientes: number;
    errors: number;
}

export default function PaginaIndex(): React.ReactElement {
    const [data, setData] = useState<DashboardData | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {

       async function  Obtener () {
            try {
                setLoading(true);
                const response = await api.get('/dashboard/');
                setData(response.data);
            } catch (err: any) {
                console.error("Error al cargar el dashboard:", err);
            } finally {
                setLoading(false);
            }
        }
        Obtener();
    }, []);

    const estadisticas = [
        {
            titulo: "Equipos",
            valor: data ? data.equipos : "0",
            cambio: "Registrados en sistema",
            icono: Laptop,
            colorIcono: "text-blue-600",
            bgIcono: "bg-blue-50",
            border: "border-blue-100",
        },
        {
            titulo: "Entregas",
            valor: data ? data.entregas : "0",
            cambio: "Registrados en sistema",
            icono: CheckCircle2,
            colorIcono: "text-emerald-600",
            bgIcono: "bg-emerald-50",
            border: "border-emerald-100",
        },
        {
            titulo: "Clientes",
            valor: data ? data.clientes: "0",
            cambio: "Registrados en sistema",
            icono: Users,
            colorIcono: "text-indigo-600",
            bgIcono: "bg-indigo-50",
            border: "border-indigo-100",
        },
        {
            titulo: "Errores",
            valor: data ? data.errors : "0",
            cambio: "Registrados en sistema",
            icono: AlertTriangle,
            colorIcono: "text-amber-600",
            bgIcono: "bg-amber-50",
            border: "border-amber-100",
        },
    ];

    return (
        <div className="space-y-6 w-full">
            <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-slate-900 via-indigo-950 to-blue-900 p-6 sm:p-8 text-white shadow-lg">
                <div className="relative z-10 space-y-2 max-w-xl">
                    <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/10 text-xs font-medium text-blue-200 backdrop-blur-md border border-white/10">
                        <ShieldCheck className="w-3.5 h-3.5 text-blue-400" />
                        <span>Sistema de Gestión Técnica v2.4</span>
                    </div>
                    <h1 className="text-2xl sm:text-3xl font-extrabold tracking-tight">
                        ¡Bienvenido al Panel Control!
                    </h1>
                    <p className="text-slate-300 text-xs sm:text-sm leading-relaxed">
                        Monitorea el ingreso de equipos, gestiona ordenes de servicio y revisa la auditoría del sistema en tiempo real.
                    </p>
                </div>

                <div className="absolute -right-10 -bottom-10 w-60 h-60 bg-blue-500/10 rounded-full blur-3xl pointer-events-none" />
                <div className="absolute right-40 -top-10 w-40 h-40 bg-indigo-500/10 rounded-full blur-2xl pointer-events-none" />
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                {estadisticas.map((item, index) => {
                    const Icono = item.icono;
                    return (
                        <div
                            key={index}
                            className={`bg-white p-5 rounded-xl shadow-sm border ${item.border} hover:shadow-md transition-shadow flex flex-col justify-between`}
                        >
                            <div className="flex items-center justify-between">
                                <span className="text-xs font-semibold text-gray-500">{item.titulo}</span>
                                <div className={`p-2.5 rounded-xl ${item.bgIcono} ${item.colorIcono}`}>
                                    <Icono />
                                </div>
                            </div>
                            <div className="mt-4">
                                <h3 className="text-2xl font-bold text-gray-800 tracking-tight flex items-center gap-2">
                                    {loading ? (
                                        <Loader2 className="w-6 h-6 animate-spin text-gray-400" />
                                    ) : (
                                        item.valor
                                    )}
                                </h3>
                                <p className="text-[11px] font-medium text-gray-400 mt-1 flex items-center gap-1">
                                    <TrendingUp className="w-3 h-3 text-emerald-500" />
                                    {item.cambio}
                                </p>
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}