import React, {useEffect, useState} from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Plus,
    Pencil,
    Trash2,
    Search,
    RotateCcw,
    History,
    CreditCard,
    Tag,
    PackageCheck,
    FileText,
    Calendar,
    Filter, PrinterCheck
} from 'lucide-react';
import {useModal} from "../../store/useModal.ts";
import {ModalLista} from "../../helpers/ModalLista.ts";
import {useEquipos} from "../../store/useEquipos.ts";

export interface EquipoDTO {
    equipo_id: number;
    codigo: string;
    tipo_equipo: string;
    modelo: string;
    numero_serie: string;
    accesorios: string;
    descripcion_problema: string;
    observacion: string;
    fecha_recepcion: string;
    fecha_estimada_entrega: string;
    fecha_creacion: string;
    fecha_modificacion: string;
    marca_id: number;
    marca: string;
    estado_id: number;
    estado: string;
    cliente_id: number;
    nombres: string;
    apellidos: string;
}

export default function PaginaEquipos(): React.ReactElement {
    const OpenModal = useModal((state) => state.OpenModal);
    const {ObtenerEquipos, ObtenerEquipo, EliminarEquipo, DescargarPdf, listar_equipos} = useEquipos((state) => state);

    useEffect(() => {
        ObtenerEquipos();
    }, []);

    const navigate = useNavigate();
    const [busqueda, setBusqueda] = useState<string>('');

    // Estados para el reporte general de equipos por rango de fechas
    const [fechaDesde, setFechaDesde] = useState<string>('2026-07-01');
    const [fechaHasta, setFechaHasta] = useState<string>('2026-07-21');

    const equiposFiltrados = listar_equipos.filter((e) =>
        e.codigo.toLowerCase().includes(busqueda.toLowerCase()) ||
        e.tipo_equipo.toLowerCase().includes(busqueda.toLowerCase()) ||
        e.modelo.toLowerCase().includes(busqueda.toLowerCase()) ||
        e.numero_serie.toLowerCase().includes(busqueda.toLowerCase()) ||
        e.marca.toLowerCase().includes(busqueda.toLowerCase()) ||
        e.nombres.toLowerCase().includes(busqueda.toLowerCase()) ||
        e.apellidos.toLowerCase().includes(busqueda.toLowerCase())
    );

    const handleVerHistorial = (equipo_id: number) => {
        navigate(`/equipos/historial?equipo_id=${equipo_id}`);
    };

    const handleGestionarEntrega = (equipo_id: number) => {
        navigate(`/equipos/entrega?equipo_id=${equipo_id}`);
    };

    const handleGestionarPagos = (equipo_id: number) => {
        navigate(`/equipos/pagos?equipo_id=${equipo_id}`);
    };

    const handleGenerarPDFEquipos = () => {
        const payloadReporte = { fecha_desde: fechaDesde, fecha_hasta: fechaHasta };
        console.log("Generando Reporte PDF de Equipos en Go:", payloadReporte);
    };

    const getBadgeEstado = (estado: string) => {
        switch (estado.toLowerCase()) {
            case 'listo':
            case 'entregado':
                return 'bg-green-50 text-green-700 border-green-200';
            case 'en revisión':
            case 'en proceso':
                return 'bg-amber-50 text-amber-700 border-amber-200';
            case 'pendiente':
                return 'bg-blue-50 text-blue-700 border-blue-200';
            default:
                return 'bg-gray-50 text-gray-600 border-gray-200';
        }
    };

    return (
        <div className="space-y-6 w-full">
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div>
                    <h1 className="text-2xl font-bold text-gray-800">Listado de Equipos</h1>
                    <p className="text-xs text-gray-500 mt-0.5">Mantenimiento, recepción, cobros y entregas</p>
                </div>

                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">

                    <div className="relative flex-1 min-w-[240px] max-w-md">
                        <input
                            type="text"
                            placeholder="Buscar código, cliente, serie..."
                            value={busqueda}
                            onChange={(e) => setBusqueda(e.target.value)}
                            className="w-full px-3.5 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        />
                    </div>

                    <div className="flex flex-wrap items-center gap-2">

                        <button
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 border border-gray-200 rounded-lg transition-colors shadow-sm"
                            title="Limpiar búsqueda"
                        >
                            <RotateCcw className="w-3.5 h-3.5" />
                            <span>Limpiar</span>
                        </button>

                        <button
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-slate-800 hover:bg-slate-900 rounded-lg transition-colors shadow-sm"
                        >
                            <Search className="w-3.5 h-3.5" />
                            <span>Buscar</span>
                        </button>

                        <button
                            onClick={() => OpenModal(ModalLista.modal_marca)}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-purple-600 hover:bg-purple-700 rounded-lg transition-colors shadow-sm"
                        >
                            <Tag className="w-3.5 h-3.5" />
                            <span>Nueva Marca</span>
                        </button>

                        <button
                            onClick={() => OpenModal(ModalLista.modal_equipo)}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                        >
                            <Plus className="w-3.5 h-3.5" />
                            <span>Nuevo Equipo</span>
                        </button>

                    </div>
                </div>
            </div>

            <div className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 w-full flex flex-wrap items-center justify-between gap-4">

                <div className="flex items-center gap-2 text-gray-700 font-semibold text-xs uppercase tracking-wider">
                    <Filter className="w-4 h-4 text-blue-600" />
                    <span>Generación de Reportes PDF de Equipos</span>
                </div>

                <div className="flex flex-wrap items-center gap-3">

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

                    <button
                        onClick={handleGenerarPDFEquipos}
                        className="inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-semibold text-white bg-red-600 hover:bg-red-700 rounded-lg transition-colors shadow-sm"
                        title="Exportar Reporte de Equipos a PDF"
                    >
                        <FileText className="w-4 h-4" />
                        <span>Generar Reporte PDF</span>
                    </button>

                </div>
            </div>

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full text-left text-sm text-gray-600">

                        <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3">Código</th>
                            <th className="px-4 py-3">Equipo / Marca</th>
                            <th className="px-4 py-3">Cliente</th>
                            <th className="px-4 py-3">Problema / Obs.</th>
                            <th className="px-4 py-3">Estado</th>
                            <th className="px-4 py-3">Fechas</th>
                            <th className="px-4 py-3 text-center">Acciones</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {equiposFiltrados.length > 0 ? (
                            equiposFiltrados.map((equipo) => (
                                <tr key={equipo.equipo_id} className="hover:bg-gray-50/80 transition-colors text-xs">

                                    <td className="px-4 py-3 whitespace-nowrap font-bold text-gray-900">
                                        {equipo.codigo}
                                    </td>

                                    {/* Equipo, Modelo, Marca y Serie */}
                                    <td className="px-4 py-3 whitespace-nowrap">
                                        <div className="font-semibold text-gray-800 text-sm">{equipo.tipo_equipo} {equipo.modelo}</div>
                                        <div className="flex items-center gap-2 mt-0.5">
                                        <span className="px-1.5 py-0.5 text-[10px] font-semibold bg-gray-100 text-gray-600 rounded">
                                          {equipo.marca}
                                        </span>
                                            <span className="text-[11px] text-gray-400 font-mono">S/N: {equipo.numero_serie}</span>
                                        </div>
                                    </td>

                                    {/* Cliente */}
                                    <td className="px-4 py-3 whitespace-nowrap">
                                        <div className="font-medium text-gray-800">{equipo.nombres} {equipo.apellidos}</div>
                                        <div className="text-[11px] text-gray-400">ID: #{equipo.cliente_id}</div>
                                    </td>

                                    {/* Problema */}
                                    <td className="px-4 py-3 max-w-xs truncate" title={equipo.descripcion_problema}>
                                        <div className="text-gray-800 truncate">{equipo.descripcion_problema}</div>
                                        {equipo.accesorios && (
                                            <div className="text-[11px] text-gray-400 truncate" title={`Accesorios: ${equipo.accesorios}`}>
                                                Acc: {equipo.accesorios}
                                            </div>
                                        )}
                                    </td>

                                    {/* Estado */}
                                    <td className="px-4 py-3 whitespace-nowrap">
                                    <span className={`inline-block px-2 py-0.5 text-[11px] font-semibold rounded-full border ${getBadgeEstado(equipo.estado)}`}>
                                        {equipo.estado}
                                    </span>
                                    </td>

                                    {/* Fechas */}
                                    <td className="px-4 py-3 whitespace-nowrap text-[11px] text-gray-500">
                                        <div>Rec: {equipo.fecha_recepcion}</div>
                                        <div className="text-gray-400">Est: {equipo.fecha_estimada_entrega}</div>
                                    </td>

                                    <td className="px-4 py-3 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">
                                            <button
                                                className="p-1.5 text-aqua-600 bg-green-50 hover:bg-green-100 rounded-lg transition-colors border border-indigo-100"
                                                title="Imprimir Orden"
                                            >
                                                <PrinterCheck className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleVerHistorial(equipo.equipo_id)}
                                                className="p-1.5 text-indigo-600 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors border border-indigo-100"
                                                title="Ver Historial"
                                            >
                                                <History className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleGestionarEntrega(equipo.equipo_id)}
                                                className="p-1.5 text-amber-600 bg-amber-50 hover:bg-amber-100 rounded-lg transition-colors border border-amber-100"
                                                title="Gestionar Entrega"
                                            >
                                                <PackageCheck className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleGestionarPagos(equipo.equipo_id)}
                                                className="p-1.5 text-emerald-600 bg-emerald-50 hover:bg-emerald-100 rounded-lg transition-colors border border-emerald-100"
                                                title="Gestionar Pagos"
                                            >
                                                <CreditCard className="w-4 h-4" />
                                            </button>
                                            <button

                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100"
                                                title="Editar equipo"
                                            >
                                                <Pencil className="w-4 h-4" />
                                            </button>
                                            <button

                                                className="p-1.5 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors border border-red-100"
                                                title="Eliminar equipo"
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </button>

                                        </div>
                                    </td>

                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={7} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron equipos registrados.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                {/* Pie de Tabla */}
                <div className="p-3 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de registros: {equiposFiltrados.length}</span>
                </div>
            </div>
        </div>
    );
}