import { useSearchParams, useNavigate } from 'react-router-dom';
import {
    Plus,
    Pencil,
    Search,
    RotateCcw,
    ArrowLeft,
    PackageCheck,
    User,
    CheckCircle2,
    XCircle,
    PrinterCheck
} from 'lucide-react';
import React, { useState } from "react";
import { useModal } from '../../store/useModal.ts';
import { ModalLista } from '../../helpers/ModalLista.ts';

// Interface 1:1 con tu API de Go
export interface EntregaDTO {
    entrega_id: number;
    fecha_entrega: string;
    trabajos_realizados: string;
    estado_final_equipo: string;
    conformidad_cliente: number; // 1 = Conforme, 0 = No Conforme
    comprobante_nro: string;
    equipo_codigo: string;
    nombres: string; // Usuario/Técnico responsable
    equipo_id?: number; // Opcional si se recibe para filtrado por URL
}

// Datos de prueba iniciales basados en tu JSON real
const entregasIniciales: EntregaDTO[] = [
    {
        entrega_id: 1,
        fecha_entrega: "14/07/2026 20:15:30",
        trabajos_realizados: "SE FORMATEO",
        estado_final_equipo: "BUEN ESTADO AUN HASTA NUEVO DIAGNOSTICO",
        conformidad_cliente: 1,
        comprobante_nro: "000001",
        equipo_codigo: "E00002",
        nombres: "SISTEMA ",
        equipo_id: 101
    },
    {
        entrega_id: 2,
        fecha_entrega: "15/07/2026 11:20:00",
        trabajos_realizados: "CAMBIO DE PANTALLA Y MANTENIMIENTO PREVENTIVO",
        estado_final_equipo: "OPERATIVO 100%",
        conformidad_cliente: 1,
        comprobante_nro: "000002",
        equipo_codigo: "E00003",
        nombres: "TECNICO SUP",
        equipo_id: 102
    }
];

export default function PaginaEntregas(): React.ReactElement {
    const { OpenModal } = useModal((state) => state);
    const navigate = useNavigate();

    // Captura estricta del equipo_id desde la URL
    const [searchParams] = useSearchParams();
    const equipoIdParam = searchParams.get('equipo_id');

    const [entregas] = useState<EntregaDTO[]>(entregasIniciales);
    const [busqueda, setBusqueda] = useState<string>('');

    // Filtrado de la tabla únicamente por equipo_id de la URL
    const entregasFiltradas = entregas.filter((e) => {
        return equipoIdParam ? e.equipo_id === Number(equipoIdParam) : true;
    });

    const handleRegresar = () => {
        navigate(-1);
    };

    const handleLimpiar = () => setBusqueda('');

    const handleBuscar = () => {
        console.log("Buscando entregas para equipo_id:", equipoIdParam, "con filtro:", busqueda);
    };

    const handleEditar = (entrega: EntregaDTO) => {
        console.log("Editar entrega ID:", entrega.entrega_id);
    };

    const handleGenerarPDFComprobante = (entrega: EntregaDTO) => {
        console.log("Generando comprobante de entrega PDF Nro:", entrega.comprobante_nro);
    };

    return (
        <div className="space-y-6 w-full">
            {/* Encabezado Superior */}
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-amber-50 text-amber-600 rounded-lg">
                            <PackageCheck className="w-6 h-6" />
                        </div>
                        <div>
                            <h1 className="text-2xl font-bold text-gray-800">
                                Gestión de Entregas {equipoIdParam && `(Equipo #${equipoIdParam})`}
                            </h1>
                            <p className="text-xs text-gray-500 mt-0.5">Comprobantes, actas de conformidad y estado final de equipos</p>
                        </div>
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

                {/* Barra de Búsqueda y Botones */}
                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">

                    {/* Input de Búsqueda */}
                    <div className="relative flex-1 min-w-[240px] max-w-md">
                        <input
                            type="text"
                            placeholder="Buscar en entregas..."
                            value={busqueda}
                            onChange={(e) => setBusqueda(e.target.value)}
                            className="w-full px-3.5 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-amber-500/20 focus:border-amber-500 transition-all"
                        />
                    </div>

                    {/* Grupo de Botones */}
                    <div className="flex flex-wrap items-center gap-2">
                        <button
                            onClick={handleLimpiar}
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
                            onClick={() => OpenModal(ModalLista.modal_entrega)}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-amber-600 hover:bg-amber-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                        >
                            <Plus className="w-3.5 h-3.5" />
                            <span>Nueva Entrega</span>
                        </button>
                    </div>

                </div>
            </div>

            {/* Tabla de Entregas */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full min-w-full text-left text-sm text-gray-600">

                        <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3.5 w-16">ID</th>
                            <th className="px-4 py-3.5 w-32">Comprobante</th>
                            <th className="px-4 py-3.5 w-32">Código Equipo</th>
                            <th className="px-4 py-3.5 min-w-[200px]">Trabajos Realizados</th>
                            <th className="px-4 py-3.5 min-w-[200px]">Estado Final</th>
                            <th className="px-4 py-3.5 w-32">Conformidad</th>
                            <th className="px-4 py-3.5 w-36">Registrado Por</th>
                            <th className="px-4 py-3.5 w-40">Fecha Entrega</th>
                            <th className="px-4 py-3.5 w-24 text-center">Acciones</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {entregasFiltradas.length > 0 ? (
                            entregasFiltradas.map((entrega) => (
                                <tr key={entrega.entrega_id} className="hover:bg-gray-50/80 transition-colors text-xs">

                                    {/* ID */}
                                    <td className="px-4 py-3.5 font-bold text-gray-900">
                                        #{entrega.entrega_id}
                                    </td>

                                    {/* Comprobante Nro */}
                                    <td className="px-4 py-3.5 font-mono font-bold text-amber-700 whitespace-nowrap">
                                        N° {entrega.comprobante_nro}
                                    </td>

                                    {/* Código Equipo */}
                                    <td className="px-4 py-3.5 font-semibold text-gray-800 whitespace-nowrap">
                                        {entrega.equipo_codigo}
                                    </td>

                                    {/* Trabajos Realizados */}
                                    <td className="px-4 py-3.5 max-w-xs truncate" title={entrega.trabajos_realizados}>
                                        <span className="font-medium text-gray-800">{entrega.trabajos_realizados || '-'}</span>
                                    </td>

                                    {/* Estado Final Equipo */}
                                    <td className="px-4 py-3.5 max-w-xs truncate" title={entrega.estado_final_equipo}>
                                        <span className="text-gray-700">{entrega.estado_final_equipo || '-'}</span>
                                    </td>

                                    {/* Conformidad Cliente (Badge) */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        {entrega.conformidad_cliente === 1 ? (
                                            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 text-[11px] font-semibold bg-green-50 text-green-700 border border-green-200 rounded-full">
                                          <CheckCircle2 className="w-3 h-3" />
                                          Conforme
                                        </span>
                                        ) : (
                                            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 text-[11px] font-semibold bg-red-50 text-red-700 border border-red-200 rounded-full">
                                          <XCircle className="w-3 h-3" />
                                          No Conforme
                                        </span>
                                        )}
                                    </td>

                                    {/* Registrado Por */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="inline-flex items-center gap-1 text-gray-700 font-medium">
                                            <User className="w-3.5 h-3.5 text-gray-400" />
                                            <span>{entrega.nombres}</span>
                                        </div>
                                    </td>

                                    {/* Fecha Entrega */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-gray-500 font-medium">
                                        {entrega.fecha_entrega}
                                    </td>

                                    {/* Acciones: PDF y Editar */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">

                                            <button
                                                className="p-1.5 text-emerald-600 bg-emerald-50 hover:bg-emerald-100 rounded-lg transition-colors border border-emerald-100 cursor-pointer"
                                                title="Imprimir Orden"
                                            >
                                                <PrinterCheck className="w-4 h-4" />
                                            </button>

                                            {/* Botón Editar */}
                                            <button
                                                onClick={() => OpenModal(ModalLista.modal_entrega)}
                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100 cursor-pointer"
                                                title="Editar entrega"
                                            >
                                                <Pencil className="w-4 h-4" />
                                            </button>

                                        </div>
                                    </td>

                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={9} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron registros de entrega para este equipo.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                {/* Pie de Tabla */}
                <div className="p-3 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de entregas: {entregasFiltradas.length}</span>
                </div>
            </div>
        </div>
    );
}