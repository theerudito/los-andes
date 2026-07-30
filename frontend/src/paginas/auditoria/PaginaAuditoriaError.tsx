import React, { useEffect, useState } from 'react';
import { Search, AlertTriangle, Calendar, Filter } from 'lucide-react';
import type { LogError } from "../../modelos/logError.ts";
import { logsService } from "../../servicios/logServicio.ts";
import type { reqLog } from "../../modelos/logOk.ts";
import { formatearFecha } from "../../helpers/formatearFecha.ts";
import ReactDatePicker from "react-datepicker";

export default function PaginaAuditoriaError(): React.ReactElement {
    const [logs, setLogs] = useState<LogError[]>([]);
    const [fechaDesde, setFechaDesde] = useState<Date | null>(new Date());
    const [fechaHasta, setFechaHasta] = useState<Date | null>(new Date());
    const [modulo, setModulo] = useState<string>('');

    async function handleBuscar() {
        const payload = {
            fecha_desde: formatearFecha(fechaDesde),
            fecha_hasta: formatearFecha(fechaHasta),
            modulo: modulo
        };
        const data = await logsService.obtenerLogsError(payload);
        setLogs(data);
    }

    async function ObtenerLogs() {
        const obj: reqLog = {
            fecha_desde: formatearFecha(fechaDesde),
            fecha_hasta: formatearFecha(fechaHasta),
            modulo: modulo
        };
        
        const data = await logsService.obtenerLogsError(obj);

        setLogs(data);
    }

    useEffect(() => {
        ObtenerLogs();
    }, []);

    const handleRangoChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const opcion = e.target.value;
        const hoy = new Date();

        switch (opcion) {
            case "hoy":
                setFechaDesde(hoy);
                setFechaHasta(hoy);
                break;

            case "mes_actual": {
                const primerDia = new Date(hoy.getFullYear(), hoy.getMonth(), 1);
                const ultimoDia = new Date(hoy.getFullYear(), hoy.getMonth() + 1, 0);
                setFechaDesde(primerDia);
                setFechaHasta(ultimoDia);
                break;
            }

            case "mes_anterior": {
                const primerDia = new Date(hoy.getFullYear(), hoy.getMonth() - 1, 1);
                const ultimoDia = new Date(hoy.getFullYear(), hoy.getMonth(), 0);
                setFechaDesde(primerDia);
                setFechaHasta(ultimoDia);
                break;
            }

            case "anio_actual": {
                const primerDia = new Date(hoy.getFullYear(), 0, 1);
                const ultimoDia = new Date(hoy.getFullYear(), 11, 31);
                setFechaDesde(primerDia);
                setFechaHasta(ultimoDia);
                break;
            }

            case "anio_anterior": {
                const primerDia = new Date(hoy.getFullYear() - 1, 0, 1);
                const ultimoDia = new Date(hoy.getFullYear() - 1, 11, 31);
                setFechaDesde(primerDia);
                setFechaHasta(ultimoDia);
                break;
            }

            default:
                break;
        }
    };

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

                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-end justify-end gap-3 w-full">
                    <div>
                        <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                            <Filter className="w-3.5 h-3.5 text-gray-400" /> Período
                        </label>
                        <select
                            defaultValue=""
                            onChange={handleRangoChange}
                            className="w-36 px-2.5 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 font-medium focus:outline-none focus:ring-0 focus:border-red-500 cursor-pointer"
                        >
                            <option value="hoy">Hoy</option>
                            <option value="mes_actual">Mes actual</option>
                            <option value="mes_anterior">Mes anterior</option>
                            <option value="anio_actual">Año actual</option>
                            <option value="anio_anterior">Año anterior</option>
                        </select>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                            <Calendar className="w-3.5 h-3.5 text-gray-400" /> Fecha Desde
                        </label>
                        <ReactDatePicker
                            selected={fechaDesde}
                            onChange={(date: Date | null) => setFechaDesde(date)}
                            dateFormat="dd/MM/yyyy"
                            locale="es"
                            popperClassName="!z-[9999]"
                            className="w-28 px-2.5 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 font-mono font-medium focus:outline-none focus:ring-0 focus:border-red-500 cursor-pointer"
                        />
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                            <Calendar className="w-3.5 h-3.5 text-gray-400" /> Fecha Hasta
                        </label>
                        <ReactDatePicker
                            selected={fechaHasta}
                            onChange={(date: Date | null) => setFechaHasta(date)}
                            dateFormat="dd/MM/yyyy"
                            locale="es"
                            popperClassName="!z-[9999]"
                            className="w-28 px-2.5 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 font-mono font-medium focus:outline-none focus:ring-0 focus:border-red-500 cursor-pointer"
                        />
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                            <Filter className="w-3.5 h-3.5 text-gray-400" /> Módulo
                        </label>
                        <select
                            value={modulo}
                            onChange={(e) => setModulo(e.target.value)}
                            className="w-44 px-3 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:ring-0 focus:border-red-500 cursor-pointer"
                        >
                            <option value="">Todos los módulos</option>
                            <option value="clientes">Clientes</option>
                            <option value="equipos">Equipos</option>
                            <option value="historial">Historial</option>
                            <option value="entregas">Entregas</option>
                            <option value="marcas">Marcas</option>
                            <option value="usuarios">Usuarios</option>
                            <option value="cuentas">Cuentas</option>
                            <option value="secuencial">Secuencial</option>
                        </select>
                    </div>

                    <div>
                        <button
                            type="button"
                            onClick={handleBuscar}
                            className="inline-flex items-center gap-1.5 px-4 py-1.5 text-xs font-semibold text-white bg-red-600 hover:bg-red-700 active:bg-red-800 focus:outline-none focus:ring-0 focus:bg-red-600 focus-visible:outline-none focus-visible:ring-0 focus-visible:bg-red-600 rounded-lg transition-colors shadow-sm cursor-pointer"
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