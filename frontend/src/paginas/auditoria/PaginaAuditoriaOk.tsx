import React, { useState } from 'react';
import { Search, RotateCcw, CheckCircle2, Calendar, Filter, User } from 'lucide-react';

export interface LogOk {
    log_ok_id: number;
    fecha: string;
    modulo: string;
    usuario: string;
    accion: string;
    descripcion: string;
}

// Datos de prueba iniciales basados en la respuesta de tu API
const logsIniciales: LogOk[] = [
    {
        log_ok_id: 18,
        fecha: "14/07/2026 20:31:58",
        modulo: "entregas",
        usuario: "SISTEMA",
        accion: "INSERT",
        descripcion: "registro creado correctamente"
    },
    {
        log_ok_id: 17,
        fecha: "14/07/2026 20:22:28",
        modulo: "historial_reparaciones",
        usuario: "SISTEMA",
        accion: "UPDATE",
        descripcion: "registro actualizado correctamente"
    },
    {
        log_ok_id: 16,
        fecha: "14/07/2026 20:21:05",
        modulo: "cuentas",
        usuario: "SISTEMA",
        accion: "UPDATE",
        descripcion: "registro actualizando correctamente"
    },
    {
        log_ok_id: 15,
        fecha: "14/07/2026 20:19:29",
        modulo: "historial_reparaciones",
        usuario: "SISTEMA",
        accion: "UPDATE",
        descripcion: "registro actualizado correctamente"
    },
    {
        log_ok_id: 14,
        fecha: "14/07/2026 19:44:57",
        modulo: "equipos",
        usuario: "SISTEMA",
        accion: "INSERT",
        descripcion: "registro creado correctamente"
    },
    {
        log_ok_id: 13,
        fecha: "14/07/2026 19:44:32",
        modulo: "equipos",
        usuario: "SISTEMA",
        accion: "DELETE",
        descripcion: "registro eliminado correctamente"
    },
    {
        log_ok_id: 12,
        fecha: "14/07/2026 19:42:01",
        modulo: "equipos",
        usuario: "SISTEMA",
        accion: "UPDATE",
        descripcion: "registro actualizado correctamente"
    }
];

export default function PaginaAuditoriaOk(): React.ReactElement {
    const [logs, setLogs] = useState<LogOk[]>(logsIniciales);

    // Estados de los filtros para la petición al backend
    const [fechaDesde, setFechaDesde] = useState<string>('2026-07-14');
    const [fechaHasta, setFechaHasta] = useState<string>('2026-07-14');
    const [modulo, setModulo] = useState<string>('');

    const handleBuscar = () => {
        const payload = {
            fecha_desde: fechaDesde,
            fecha_hasta: fechaHasta,
            modulo: modulo
        };
        console.log("Consultando auditoría de operaciones exitosas:", payload);
        // Aquí ejecutas la llamada a la API enviando `payload`
    };

    const handleLimpiar = () => {
        setFechaDesde('');
        setFechaHasta('');
        setModulo('');
    };

    // Helper para asignar color a la insignia según la acción (INSERT, UPDATE, DELETE)
    const getBadgeAccion = (accion: string) => {
        switch (accion.toUpperCase()) {
            case 'INSERT':
                return 'bg-emerald-50 text-emerald-700 border-emerald-200';
            case 'UPDATE':
                return 'bg-blue-50 text-blue-700 border-blue-200';
            case 'DELETE':
                return 'bg-red-50 text-red-700 border-red-200';
            default:
                return 'bg-gray-50 text-gray-700 border-gray-200';
        }
    };

    return (
        <div className="space-y-6 w-full">
            {/* Encabezado Superior con Formulario de Filtros */}
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div className="flex items-center gap-3">
                    <div className="p-2 bg-emerald-50 text-emerald-600 rounded-lg">
                        <CheckCircle2 className="w-6 h-6" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800">Auditoría de Operaciones (OK)</h1>
                        <p className="text-xs text-gray-500 mt-0.5">Historial de registros creados, modificados o eliminados con éxito</p>
                    </div>
                </div>

                {/* Sección de Filtros por Fecha y Módulo */}
                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-end justify-between gap-4">
                    <div className="flex flex-wrap items-center gap-3 flex-1">

                        {/* Input Fecha Desde */}
                        <div className="flex-1 min-w-[160px]">
                            <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                                <Calendar className="w-3.5 h-3.5" /> Fecha Desde
                            </label>
                            <input
                                type="date"
                                value={fechaDesde}
                                onChange={(e) => setFechaDesde(e.target.value)}
                                className="w-full px-3 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all"
                            />
                        </div>

                        {/* Input Fecha Hasta */}
                        <div className="flex-1 min-w-[160px]">
                            <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                                <Calendar className="w-3.5 h-3.5" /> Fecha Hasta
                            </label>
                            <input
                                type="date"
                                value={fechaHasta}
                                onChange={(e) => setFechaHasta(e.target.value)}
                                className="w-full px-3 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all"
                            />
                        </div>

                        {/* Select Módulo */}
                        <div className="flex-1 min-w-[180px]">
                            <label className="block text-xs font-semibold text-gray-600 mb-1 flex items-center gap-1">
                                <Filter className="w-3.5 h-3.5" /> Módulo
                            </label>
                            <select
                                value={modulo}
                                onChange={(e) => setModulo(e.target.value)}
                                className="w-full px-3 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all"
                            >
                                <option value="">Todos los módulos</option>
                                <option value="clientes">Clientes</option>
                                <option value="equipos">Equipos</option>
                                <option value="historial_reparaciones">Historial Reparaciones</option>
                                <option value="entregas">Entregas</option>
                                <option value="cuentas">Cuentas / Pagos</option>
                                <option value="usuarios">Usuarios</option>
                                <option value="marcas">Marcas</option>
                            </select>
                        </div>

                    </div>

                    {/* Botones de Acción */}
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
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg transition-colors shadow-sm"
                        >
                            <Search className="w-3.5 h-3.5" />
                            <span>Buscar</span>
                        </button>
                    </div>
                </div>
            </div>

            {/* Tabla de Registros de Auditoría OK */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full text-left text-sm text-gray-600">

                        <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3.5 w-16">ID</th>
                            <th className="px-4 py-3.5 w-44">Fecha / Hora</th>
                            <th className="px-4 py-3.5 w-32">Acción</th>
                            <th className="px-4 py-3.5 w-44">Módulo</th>
                            <th className="px-4 py-3.5 w-36">Usuario</th>
                            <th className="px-4 py-3.5 min-w-[250px]">Descripción</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {logs.length > 0 ? (
                            logs.map((log) => (
                                <tr key={log.log_ok_id} className="hover:bg-gray-50/80 transition-colors text-xs">

                                    {/* ID */}
                                    <td className="px-4 py-3.5 font-bold text-gray-700">
                                        #{log.log_ok_id}
                                    </td>

                                    {/* Fecha */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-gray-600 font-medium">
                                        {log.fecha}
                                    </td>

                                    {/* Acción (Badge) */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                      <span className={`inline-block px-2.5 py-0.5 text-[11px] font-bold border rounded-md uppercase ${getBadgeAccion(log.accion)}`}>
                        {log.accion}
                      </span>
                                    </td>

                                    {/* Módulo */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                      <span className="inline-block px-2 py-0.5 text-[11px] font-semibold bg-gray-100 text-gray-700 rounded-md">
                        {log.modulo}
                      </span>
                                    </td>

                                    {/* Usuario */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="inline-flex items-center gap-1.5 font-medium text-gray-700">
                                            <User className="w-3.5 h-3.5 text-gray-400" />
                                            <span>{log.usuario}</span>
                                        </div>
                                    </td>

                                    {/* Descripción */}
                                    <td className="px-4 py-3.5 text-gray-700">
                                        {log.descripcion}
                                    </td>

                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={6} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron registros de auditoría para los filtros seleccionados.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                {/* Pie de Tabla */}
                <div className="p-3 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de registros exitosos: {logs.length}</span>
                </div>
            </div>
        </div>
    );
}