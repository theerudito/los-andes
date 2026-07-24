import React, {useEffect, useState} from 'react';
import {
    Plus,
    Pencil,
    Trash2,
    Search,
    RotateCcw,
    FileText,
    Calendar,
    Filter
} from 'lucide-react';
import {useModal} from "../../store/useModal.ts";
import {ModalLista} from "../../helpers/ModalLista.ts";
import {useClientes} from "../../store/useClientes.ts";
import type {ReqCliente} from "../../modelos/clientes.ts";

export default function PaginaClientes(): React.ReactElement {
    const {OpenModal} = useModal((state) => state);
    const {ObtenerClientes, listar_clientes, EliminarCliente, ObtenerCliente, DescargarPdf} = useClientes((state) => state);

    const [busqueda, setBusqueda] = useState<string>('');
    const [fechaDesde, setFechaDesde] = useState<string>('2026-07-01');
    const [fechaHasta, setFechaHasta] = useState<string>('2026-07-16');

    useEffect(() => {
        ObtenerClientes();
    }, []);

    function VerCliente (id: number) {
        OpenModal(ModalLista.modal_cliente)
        ObtenerCliente(id)
    }

    function VerReporte() {
        const obj: ReqCliente = {
            fecha_desde: fechaDesde,
            fecha_hasta: fechaHasta
        };
        DescargarPdf(obj);
    }

    return (
        <div className="space-y-6 w-full">

            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div>
                    <h1 className="text-2xl font-bold text-gray-800">Listado de Clientes</h1>
                    <p className="text-xs text-gray-500 mt-0.5">Gestión de clientes y generación de reportes</p>
                </div>

                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">

                    <div className="relative flex-1 min-w-[240px] max-w-md">
                        <input
                            type="text"
                            placeholder="Buscar por nombre, identificación o correo..."
                            value={busqueda}
                            onChange={(e) => setBusqueda(e.target.value)}
                            className="w-full px-3.5 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        />
                    </div>

                    <div className="flex flex-wrap items-center gap-2">
                        <button
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 border border-gray-200 rounded-lg transition-colors shadow-sm cursor-pointer"
                            title="Limpiar búsqueda"
                        >
                            <RotateCcw className="w-3.5 h-3.5"/>
                            <span>Limpiar</span>
                        </button>

                        <button
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-slate-800 hover:bg-slate-900 rounded-lg transition-colors shadow-sm cursor-pointer"
                        >
                            <Search className="w-3.5 h-3.5"/>
                            <span>Buscar</span>
                        </button>

                        <button
                            onClick={() => OpenModal(ModalLista.modal_cliente)}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                        >
                            <Plus className="w-3.5 h-3.5"/>
                            <span>Nuevo Cliente</span>
                        </button>
                    </div>
                </div>
            </div>

            <div
                className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 w-full flex flex-wrap items-center justify-between gap-4">

                <div className="flex items-center gap-2 text-gray-700 font-semibold text-xs uppercase tracking-wider">
                    <Filter className="w-4 h-4 text-blue-600"/>
                    <span>Generación de Reportes PDF</span>
                </div>

                <div className="flex flex-wrap items-center gap-3">

                    <div className="flex items-center gap-2">
                        <label className="text-xs font-medium text-gray-600 flex items-center gap-1">
                            <Calendar className="w-3.5 h-3.5 text-gray-400"/> Desde:
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
                            <Calendar className="w-3.5 h-3.5 text-gray-400"/> Hasta:
                        </label>
                        <input
                            type="date"
                            value={fechaHasta}
                            onChange={(e) => setFechaHasta(e.target.value)}
                            className="px-2.5 py-1.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 focus:outline-none focus:border-blue-500"
                        />
                    </div>

                    <div className="w-[1px] h-6 bg-gray-200 hidden sm:block"/>
                    <button
                        onClick={VerReporte}
                        className="inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-semibold text-white bg-red-600 hover:bg-red-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                        title="Exportar a PDF"
                    >
                        <FileText className="w-4 h-4"/>
                        <span>Generar Reporte PDF</span>
                    </button>

                </div>
            </div>

            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full min-w-full text-left text-sm text-gray-600">

                        <thead
                            className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3.5 w-16">ID</th>
                            <th className="px-4 py-3.5 w-44">Identificación</th>
                            <th className="px-4 py-3.5 min-w-[180px]">Nombres / Apellidos</th>
                            <th className="px-4 py-3.5 min-w-[180px]">Contacto</th>
                            <th className="px-4 py-3.5 min-w-[200px]">Dirección</th>
                            <th className="px-4 py-3.5 w-24 text-center">Acciones</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {listar_clientes.length > 0 ? (
                            listar_clientes.map((cliente) => (
                                <tr key={cliente.cliente_id} className="hover:bg-gray-50/80 transition-colors">
                                    <td className="px-4 py-3.5 font-medium text-gray-900">
                                        #{cliente.cliente_id}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="font-medium text-gray-800">{cliente.identificacion}</div>
                                        <span
                                            className="inline-block px-2 py-0.5 text-[10px] font-semibold bg-gray-100 text-gray-600 rounded">
                                            {cliente.tipo_identificacion}
                                        </span>
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap font-medium text-gray-800">
                                        {cliente.nombres} {cliente.apellidos}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div>{cliente.email}</div>
                                        <div className="text-xs text-gray-400">{cliente.telefono}</div>
                                    </td>
                                    <td className="px-4 py-3.5 truncate">
                                        {cliente.direccion || '-'}
                                    </td>


                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">
                                            <button
                                                onClick={() => VerCliente(cliente.cliente_id)}
                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100 cursor-pointer"
                                                title="Editar cliente"
                                            >
                                                <Pencil className="w-4 h-4"/>
                                            </button>
                                            <button
                                                onClick={() => EliminarCliente(cliente.cliente_id)}
                                                className="p-1.5 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors border border-red-100 cursor-pointer"
                                                title="Eliminar cliente"
                                            >
                                                <Trash2 className="w-4 h-4"/>
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr >
                                <td colSpan={7} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron clientes registrados.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                <div
                    className="p-4 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de registros: {listar_clientes.length}</span>
                </div>
            </div>
        </div>
    );
}