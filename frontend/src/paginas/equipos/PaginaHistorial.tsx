import React, {useEffect, useState} from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import {
    Pencil,
    Search,
    RotateCcw,
    Wrench,
    FileText,
    ArrowLeft,
    Plus,
} from 'lucide-react';
import { useModal } from '../../store/useModal.ts';
import { ModalLista } from '../../helpers/ModalLista.ts';
import {useHistorial} from "../../store/useHistorial.ts";

export default function PaginaHistorial(): React.ReactElement {
    const { OpenModal } = useModal((state) => state);
    const { ObtenerHistoriales, ObtenerHistorial, DescargarPdf, listar_historial, setEquipoId } = useHistorial((state) => state);
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    const equipo_id = searchParams.get('equipo_id');

    useEffect(() => {
        if (equipo_id) {
            setEquipoId(Number(equipo_id));
        }
        ObtenerHistoriales();
    }, [equipo_id]);

    const [busqueda, setBusqueda] = useState<string>('');

    const historialFiltrado = listar_historial.filter((h) => {
        return equipo_id ? h.equipo_id === Number(equipo_id) : true;
    });

    const handleRegresar = () => {
        navigate(-1);
    };

    async function VerHistorial(historialId:number){
        await ObtenerHistorial(historialId)
        OpenModal(ModalLista.modal_historial)
    }

    const handleLimpiarBusqueda = () => setBusqueda('');

    const handleBuscar = () => {
        console.log("Ejecutando búsqueda en historial para equipo_id:", equipo_id, "con filtro:", busqueda);
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
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div className="flex items-center justify-between gap-4">
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800">
                            Historial de Reparaciones {equipo_id && `(Equipo #${equipo_id})`}
                        </h1>
                        <p className="text-xs text-gray-500 mt-0.5">Seguimiento técnico y registros de intervenciones</p>
                    </div>

                    <button
                        onClick={handleRegresar}
                        className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 border border-gray-200 rounded-lg transition-colors shadow-sm shrink-0 cursor-pointer"
                        title="Volver a la vista anterior"
                    >
                        <ArrowLeft className="w-4 h-4" />
                        <span>Regresar</span>
                    </button>
                </div>

                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">

                    <div className="relative flex-1 min-w-[240px] max-w-md">
                        <input
                            type="text"
                            placeholder="Buscar en el historial..."
                            value={busqueda}
                            onChange={(e) => setBusqueda(e.target.value)}
                            className="w-full px-3.5 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        />
                    </div>

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
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-yellow-600 hover:bg-yellow-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                        >
                            <Plus className="w-3.5 h-3.5" />
                            <span>Nuevo Historial</span>
                        </button>

                        <button
                            onClick={() => DescargarPdf(Number(equipo_id))}
                            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-semibold text-white bg-red-600 hover:bg-red-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                            title="Exportar Historial a PDF"
                        >
                            <FileText className="w-4 h-4" />
                            <span>Generar Reporte PDF</span>
                        </button>
                    </div>

                </div>
            </div>

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

                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="font-bold text-gray-800 text-sm">{item.equipo}</div>
                                        <div className="text-[11px] text-gray-400 font-mono">S/N: {item.serie}</div>
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="font-medium text-gray-800">{item.nombres_cliente} {item.apellidos_cliente}</div>
                                        <div className="text-[11px] text-gray-400">ID Cliente: #{item.cliente_id}</div>
                                    </td>

                                    <td className="px-4 py-3.5 max-w-xs truncate" title={item.observaciones_tecnicas}>
                                        <span className="text-gray-700">{item.observaciones_tecnicas || '-'}</span>
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                      <span className={`inline-block px-2.5 py-0.5 text-[11px] font-semibold rounded-full border ${getBadgeEstado(item.estado)}`}>
                                        {item.estado}
                                      </span>
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="inline-flex items-center gap-1.5 text-xs text-gray-700 font-medium">
                                            <Wrench className="w-3.5 h-3.5 text-gray-400" />
                                            <span>{item.nombres_usuario} {item.apellidos_usuario}</span>
                                        </div>
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">
                                            <button
                                                onClick={() => VerHistorial(item.historial_id)}
                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100 cursor-pointer"
                                                title="Actualizar / Editar Historial"
                                            >
                                                <Pencil className="w-4 h-4" />
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

                <div className="p-3 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de registros: {historialFiltrado.length}</span>
                </div>
            </div>
        </div>
    );
}