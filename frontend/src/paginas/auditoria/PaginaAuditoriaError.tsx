import React, {useEffect, useState} from 'react';
import { Search, RotateCcw, AlertTriangle, Calendar, Filter } from 'lucide-react';
import type {LogError} from "../../modelos/logError.ts";
import {logsService} from "../../servicios/logServicio.ts";
import type {reqLog} from "../../modelos/logOk.ts";

export default function PaginaAuditoriaError(): React.ReactElement {
    const [logs, setLogs] = useState<LogError[]>([]);
    const [fechaDesde, setFechaDesde] = useState<string>("2026-07-01");
    const [fechaHasta, setFechaHasta] = useState<string>("2026-07-31");
    const [modulo, setModulo] = useState<string>('');

    async function handleBuscar () {
        const payload = {
            fecha_desde: fechaDesde,
            fecha_hasta: fechaHasta,
            modulo: modulo
        };
        const  data = await logsService.obtenerLogsError(payload)
        setLogs(data)
    }

    async function handleLimpiar () {
        setFechaDesde('');
        setFechaHasta('');
        setModulo('');
    }

    async function ObtenerLogs(){
        const obj: reqLog = {
            fecha_desde: fechaDesde,
            fecha_hasta: fechaHasta,
            modulo: modulo
        }
        const  data = await logsService.obtenerLogsError(obj)
        setLogs(data)
    }

    useEffect(() => {
        ObtenerLogs();
    }, []);

    return (
        <div className="space-y-6 w-full">
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div className="flex items-center gap-3">
                    <div className="p-2 bg-red-50 text-red-600 rounded-lg">
                        <AlertTriangle className="w-6 h-6" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800">Auditoría de Errores</h1>
                        <p className="text-xs text-gray-500 mt-0.5">Consulta los registros de excepciones y fallos del sistema</p>
                    </div>
                </div>

                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-end justify-between gap-4">
                    <div className="flex flex-wrap items-center gap-3 flex-1">
                        <div className="flex-1 min-w-[160px]">
                            <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                                <Calendar className="w-3.5 h-3.5" /> Fecha Desde
                            </label>
                            <input
                                type="date"
                                value={fechaDesde}
                                onChange={(e) => setFechaDesde(e.target.value)}
                                className="w-full px-3 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:ring-2 focus:ring-red-500/20 focus:border-red-500"
                            />
                        </div>

                        <div className="flex-1 min-w-[160px]">
                            <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                                <Calendar className="w-3.5 h-3.5" /> Fecha Hasta
                            </label>
                            <input
                                type="date"
                                value={fechaHasta}
                                onChange={(e) => setFechaHasta(e.target.value)}
                                className="w-full px-3 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:ring-2 focus:ring-red-500/20 focus:border-red-500"
                            />
                        </div>

                        <div className="flex-1 min-w-[180px]">
                            <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                                <Filter className="w-3.5 h-3.5" /> Módulo
                            </label>
                            <select
                                value={modulo}
                                onChange={(e) => setModulo(e.target.value)}
                                className="w-full px-3 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:ring-2 focus:ring-red-500/20 focus:border-red-500"
                            >
                                <option value="">Todos los módulos</option>
                                <option value="clientes">Clientes</option>
                                <option value="equipos">Equipos</option>
                                <option value="historial">Historial</option>
                                <option value="entregas">Entregas</option>
                                <option value="marcas">Marcas</option>
                                <option value="usuarios">Usuarios</option>
                                <option value="cuentas">Cuentas</option>
                                <option value="usuarios">Usuarios</option>
                                <option value="secuencial">Secuencial</option>
                            </select>
                        </div>

                    </div>

                    <div className="flex items-center gap-2">
                        <button
                            onClick={handleLimpiar}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 border border-gray-200 rounded-lg transition-colors shadow-sm"
                            title="Limpiar filtros"
                        >
                            <RotateCcw className="w-3.5 h-3.5" />
                            <span>Limpiar</span>
                        </button>
                        <button
                            onClick={handleBuscar}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-red-600 hover:bg-red-700 rounded-lg transition-colors shadow-sm"
                        >
                            <Search className="w-3.5 h-3.5" />
                            <span>Buscar</span>
                        </button>
                    </div>
                </div>
            </div>

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full text-left text-sm text-gray-600">

                        <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3.5 w-16">ID</th>
                            <th className="px-4 py-3.5 w-44">Fecha / Hora</th>
                            <th className="px-4 py-3.5 w-36">Módulo</th>
                            <th className="px-4 py-3.5 min-w-[300px]">Mensaje de Error</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {logs.length > 0 ? (
                            logs.map((log) => (
                                <tr key={log.log_error_id} className="hover:bg-red-50/20 transition-colors text-xs">
                                    <td className="px-4 py-3.5 font-bold text-gray-700">
                                        #{log.log_error_id}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap text-gray-600 font-medium">
                                        {log.fecha}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                      <span className="inline-block px-2.5 py-1 text-[11px] font-semibold bg-red-50 text-red-700 border border-red-200 rounded-md uppercase">
                        {log.modulo}
                      </span>
                                    </td>
                                    <td className="px-4 py-3.5">
                                        <div className="p-2 bg-gray-900 text-red-400 font-mono text-[11px] rounded-md border border-gray-800 break-all select-all">
                                            {log.mensaje_error}
                                        </div>
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={4} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron registros de error para los filtros seleccionados.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>
                <div className="p-3 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de errores registrados: {logs.length}</span>
                </div>
            </div>
        </div>
    );
}