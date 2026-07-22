import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Pencil,
    Trash2,
    Search,
    RotateCcw,
    Wrench,
    FileText,
    Calendar,
    Filter,
    ArrowLeft,
    Plus,
    Tag
} from 'lucide-react';
import { useModal } from '../../store/useModal.ts';
import { ModalLista } from '../../helpers/ModalLista.ts';

export interface HistorialReparacionesDTO {
    historial_id: number;
    observaciones_tecnicas: string;
    fecha: string;
    equipo_id: number;
    equipo: string;
    serie: string;
    estado_id: number;
    estado: string;
    usuario_id: number;
    nombres_usuario: string;
    apellidos_usuario: string;
    cliente_id: number;
    nombres_cliente: string;
    apellidos_cliente: string;
}

const historialInicial: HistorialReparacionesDTO[] = [
    {
        historial_id: 1,
        observaciones_tecnicas: "Se realizó cambio de pasta térmica y limpieza general de componentes internos.",
        fecha: "2026-07-21 09:15",
        equipo_id: 101,
        equipo: "Laptop ThinkPad E14",
        serie: "LNV-98765432",
        estado_id: 2,
        estado: "Reparado",
        usuario_id: 2,
        nombres_usuario: "Soporte",
        apellidos_usuario: "Técnico",
        cliente_id: 1,
        nombres_cliente: "Juan Carlos",
        apellidos_cliente: "Pérez Gómez"
    },
    {
        historial_id: 2,
        observaciones_tecnicas: "Diagnóstico inicial: Falla detectada en módulo de memoria RAM de 8GB.",
        fecha: "2026-07-20 14:30",
        equipo_id: 102,
        equipo: "PC Escritorio Custom Gamer",
        serie: "SN-PC-2026-99",
        estado_id: 1,
        estado: "En Revisión",
        usuario_id: 1,
        nombres_usuario: "Admin",
        apellidos_usuario: "Sistema",
        cliente_id: 2,
        nombres_cliente: "Empresa Tech",
        apellidos_cliente: "S.A."
    }
];

export default function PaginaHistorial(): React.ReactElement {
    const { OpenModal } = useModal((state) => state);
    const navigate = useNavigate();

    // Captura del parámetro directamente de la ruta /equipos/:equipo_id/historial
    const { equipo_id } = useParams<{ equipo_id: string }>();

    const [historial, setHistorial] = useState<HistorialReparacionesDTO[]>(historialInicial);
    const [busqueda, setBusqueda] = useState<string>('');

    // Estados para el reporte PDF por rango de fechas
    const [fechaDesde, setFechaDesde] = useState<string>('2026-07-01');
    const [fechaHasta, setFechaHasta] = useState<string>('2026-07-21');

    // Filtrado de la tabla únicamente por equipo_id proveniente de la URL
    const historialFiltrado = historial.filter((h) => {
        return equipo_id ? h.equipo_id === Number(equipo_id) : true;
    });

    const handleRegresar = () => {
        navigate(-1);
    };

    const handleLimpiarBusqueda = () => setBusqueda('');

    const handleBuscar = () => {
        console.log("Ejecutando búsqueda en historial para equipo_id:", equipo_id, "con filtro:", busqueda);
    };

    const handleActualizarHistorial = (item: HistorialReparacionesDTO) => {
        console.log("Actualizar observación/historial ID:", item.historial_id);
    };

    const handleEliminar = (id: number) => {
        if (confirm('¿Estás seguro de que deseas eliminar este registro del historial?')) {
            setHistorial(historial.filter((h) => h.historial_id !== id));
        }
    };

    const handleGenerarPDFHistorial = () => {
        const payloadReporte = {
            fecha_desde: fechaDesde,
            fecha_hasta: fechaHasta,
            equipo_id: equipo_id ? Number(equipo_id) : null
        };
        console.log("Generando Reporte PDF del Historial en Go:", payloadReporte);
    };

    const getBadgeEstado = (estado: string) => {
        switch (estado.toLowerCase()) {
            case 'reparado':
            case 'entregado':
            case 'completado':
                return 'bg-green-50 text-green-700 border-green-200';
            case 'en revisión':
            case 'en proceso':
            case 'diagnóstico':
                return 'bg-amber-50 text-amber-700 border-amber-200';
            case 'pendiente':
                return 'bg-blue-50 text-blue-700 border-blue-200';
            default:
                return 'bg-gray-50 text-gray-600 border-gray-200';
        }
    };

    return (
        <div className="space-y-6 w-full">
            {/* Encabezado Superior */}
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div className="flex items-center justify-between gap-4">
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800">
                            Historial de Reparaciones {equipo_id && `(Equipo #${equipo_id})`}
                        </h1>
                        <p className="text-xs text-gray-500 mt-0.5">Seguimiento técnico y registros de intervenciones</p>
                    </div>

                    {/* Botón Regresar */}
                    <button
                        onClick={handleRegresar}
                        className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 border border-gray-200 rounded-lg transition-colors shadow-sm shrink-0 cursor-pointer"
                        title="Volver a la vista anterior"
                    >
                        <ArrowLeft className="w-4 h-4" />
                        <span>Regresar</span>
                    </button>
                </div>

                {/* Barra de Búsqueda y Botones abajo */}
                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">

                    {/* Input de Búsqueda */}
                    <div className="relative flex-1 min-w-[240px] max-w-md">
                        <input
                            type="text"
                            placeholder="Buscar en el historial..."
                            value={busqueda}
                            onChange={(e) => setBusqueda(e.target.value)}
                            className="w-full px-3.5 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        />
                    </div>

                    {/* Grupo de Botones */}
                    <div className="flex flex-wrap items-center gap-2">
                        <button
                            onClick={handleLimpiarBusqueda}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 border border-gray-200 rounded-lg transition-colors shadow-sm cursor-pointer"
                            title="Limpiar búsqueda"
                        >
                            <RotateCcw className="w-3.5 h-3.5" />
                            <span>Limpiar</span>
                        </button>

                        <button
                            onClick={handleBuscar}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-slate-800 hover:bg-slate-900 rounded-lg transition-colors shadow-sm cursor-pointer"
                        >
                            <Search className="w-3.5 h-3.5" />
                            <span>Buscar</span>
                        </button>


                        <button
                            onClick={() => OpenModal(ModalLista.modal_historial)}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                        >
                            <Plus className="w-3.5 h-3.5" />
                            <span>Nuevo Historial</span>
                        </button>
                    </div>

                </div>
            </div>

            {/* BLOQUE DE REPORTES PDF */}
            <div className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 w-full flex flex-wrap items-center justify-between gap-4">

                <div className="flex items-center gap-2 text-gray-700 font-semibold text-xs uppercase tracking-wider">
                    <Filter className="w-4 h-4 text-blue-600" />
                    <span>Generación de Reportes PDF de Historial</span>
                </div>

                <div className="flex flex-wrap items-center gap-3">

                    {/* Fecha Desde */}
                    <div className="flex items-center gap-2">
                        <label className="text-xs font-medium text-gray-600 flex items-center gap-1">
                            <Calendar className="w-3.5 h-3.5 text-gray-400" /> Desde:
                        </label>
                        <input
                            type="date"
                            value={fechaDesde}
                            onChange={(e) => setFechaDesde(e.target.value)}
                            className="px-2.5 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:border-blue-500"
                        />
                    </div>

                    {/* Fecha Hasta */}
                    <div className="flex items-center gap-2">
                        <label className="text-xs font-medium text-gray-600 flex items-center gap-1">
                            <Calendar className="w-3.5 h-3.5 text-gray-400" /> Hasta:
                        </label>
                        <input
                            type="date"
                            value={fechaHasta}
                            onChange={(e) => setFechaHasta(e.target.value)}
                            className="px-2.5 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:border-blue-500"
                        />
                    </div>

                    <div className="w-[1px] h-6 bg-gray-200 hidden sm:block" />

                    {/* Botón Descarga PDF */}
                    <button
                        onClick={handleGenerarPDFHistorial}
                        className="inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-semibold text-white bg-red-600 hover:bg-red-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                        title="Exportar Historial a PDF"
                    >
                        <FileText className="w-4 h-4" />
                        <span>Generar Reporte PDF</span>
                    </button>

                </div>
            </div>

            {/* Tabla de Historial */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full min-w-full text-left text-sm text-gray-600">

                        <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3.5 w-16">ID</th>
                            <th className="px-4 py-3.5 min-w-[160px]">Equipo / Serie</th>
                            <th className="px-4 py-3.5 min-w-[160px]">Cliente</th>
                            <th className="px-4 py-3.5 min-w-[220px]">Observaciones Técnicas</th>
                            <th className="px-4 py-3.5 w-32">Estado</th>
                            <th className="px-4 py-3.5 min-w-[140px]">Técnico</th>
                            <th className="px-4 py-3.5 w-36">Fecha</th>
                            <th className="px-4 py-3.5 w-24 text-center">Acciones</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {historialFiltrado.length > 0 ? (
                            historialFiltrado.map((item) => (
                                <tr key={item.historial_id} className="hover:bg-gray-50/80 transition-colors text-xs">
                                    <td className="px-4 py-3.5 font-bold text-gray-900">
                                        #{item.historial_id}
                                    </td>

                                    {/* Equipo y Serie */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="font-bold text-gray-800 text-sm">{item.equipo}</div>
                                        <div className="text-[11px] text-gray-400 font-mono">S/N: {item.serie}</div>
                                    </td>

                                    {/* Cliente */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="font-medium text-gray-800">{item.nombres_cliente} {item.apellidos_cliente}</div>
                                        <div className="text-[11px] text-gray-400">ID Cliente: #{item.cliente_id}</div>
                                    </td>

                                    {/* Observaciones Técnicas */}
                                    <td className="px-4 py-3.5 max-w-xs truncate" title={item.observaciones_tecnicas}>
                                        <span className="text-gray-700">{item.observaciones_tecnicas || '-'}</span>
                                    </td>

                                    {/* Estado */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                      <span className={`inline-block px-2.5 py-0.5 text-[11px] font-semibold rounded-full border ${getBadgeEstado(item.estado)}`}>
                                        {item.estado}
                                      </span>
                                    </td>

                                    {/* Técnico Responsable */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="inline-flex items-center gap-1.5 text-xs text-gray-700 font-medium">
                                            <Wrench className="w-3.5 h-3.5 text-gray-400" />
                                            <span>{item.nombres_usuario} {item.apellidos_usuario}</span>
                                        </div>
                                    </td>

                                    {/* Fecha de registro */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-xs text-gray-400">
                                        {item.fecha}
                                    </td>

                                    {/* Acciones */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">
                                            <button
                                                onClick={() => OpenModal(ModalLista.modal_historial)}
                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100 cursor-pointer"
                                                title="Actualizar / Editar Historial"
                                            >
                                                <Pencil className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleEliminar(item.historial_id)}
                                                className="p-1.5 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors border border-red-100 cursor-pointer"
                                                title="Eliminar registro"
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={8} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron registros de historial técnico para este equipo.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                {/* Pie de Tabla */}
                <div className="p-3 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de registros: {historialFiltrado.length}</span>
                </div>
            </div>
        </div>
    );
}