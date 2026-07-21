import React, { useState } from 'react';
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

export interface Cliente {
    cliente_id: number;
    identificacion: string;
    tipo_identificacion: string;
    nombres: string;
    apellidos: string;
    telefono: string;
    email: string;
    direccion: string;
    fecha_creacion: string;
    fecha_modificacion: string;
}

const clientesIniciales: Cliente[] = [
    {
        cliente_id: 1,
        identificacion: "1712345678",
        tipo_identificacion: "CEDULA",
        nombres: "Juan Carlos",
        apellidos: "Pérez Gómez",
        telefono: "0991234567",
        email: "juan.perez@example.com",
        direccion: "Av. Central 123 y Calle 4",
        fecha_creacion: "2026-01-15 10:30",
        fecha_modificacion: "2026-02-01 14:20"
    },
    {
        cliente_id: 2,
        identificacion: "1798765432001",
        tipo_identificacion: "RUC",
        nombres: "Empresa Tech",
        apellidos: "S.A.",
        telefono: "022345678",
        email: "contacto@tech.com",
        direccion: "Panamericana Norte Km 5",
        fecha_creacion: "2026-02-10 09:00",
        fecha_modificacion: "2026-02-10 09:00"
    }
];

export default function PaginaClientes(): React.ReactElement {
    const [clientes, setClientes] = useState<Cliente[]>(clientesIniciales);
    const [busqueda, setBusqueda] = useState<string>('');

    // Estados para el reporte por rango de fechas
    const [fechaDesde, setFechaDesde] = useState<string>('2026-07-01');
    const [fechaHasta, setFechaHasta] = useState<string>('2026-07-16');

    const clientesFiltrados = clientes.filter((c) =>
        c.nombres.toLowerCase().includes(busqueda.toLowerCase()) ||
        c.apellidos.toLowerCase().includes(busqueda.toLowerCase()) ||
        c.identificacion.includes(busqueda) ||
        c.email.toLowerCase().includes(busqueda.toLowerCase())
    );

    const handleLimpiar = () => setBusqueda('');

    const handleBuscar = () => {
        console.log("Buscando cliente:", busqueda);
    };

    const handleNuevo = () => {
        console.log("Nuevo cliente");
    };

    const handleEditar = (cliente: Cliente) => {
        console.log("Editar cliente:", cliente);
    };

    const handleEliminar = (id: number) => {
        if (confirm('¿Estás seguro de que deseas eliminar este cliente?')) {
            setClientes(clientes.filter((c) => c.cliente_id !== id));
        }
    };

    // Manejador para la generación del reporte PDF
    const handleGenerarPDF = () => {
        const payloadReporte = { fecha_desde: fechaDesde, fecha_hasta: fechaHasta };
        console.log("Descargando Reporte PDF de Clientes en Go:", payloadReporte);
        // Aquí ejecutas la petición a tu endpoint de Go para generar el PDF
    };

    return (
        <div className="space-y-6 w-full">

            {/* Encabezado Superior: Título + Barra de Búsqueda y Botones de Acción */}
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div>
                    <h1 className="text-2xl font-bold text-gray-800">Listado de Clientes</h1>
                    <p className="text-xs text-gray-500 mt-0.5">Gestión de clientes y generación de reportes</p>
                </div>

                {/* Barra de Búsqueda y Acciones Rápidas */}
                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">

                    {/* Input de Búsqueda general */}
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
                            onClick={handleLimpiar}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 border border-gray-200 rounded-lg transition-colors shadow-sm"
                            title="Limpiar búsqueda"
                        >
                            <RotateCcw className="w-3.5 h-3.5" />
                            <span>Limpiar</span>
                        </button>

                        <button
                            onClick={handleBuscar}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-slate-800 hover:bg-slate-900 rounded-lg transition-colors shadow-sm"
                        >
                            <Search className="w-3.5 h-3.5" />
                            <span>Buscar</span>
                        </button>

                        <button
                            onClick={handleNuevo}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                        >
                            <Plus className="w-3.5 h-3.5" />
                            <span>Nuevo Cliente</span>
                        </button>
                    </div>
                </div>
            </div>

            {/* BLOQUE DE REPORTES PDF (Filtro por Rango de Fechas) */}
            <div className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 w-full flex flex-wrap items-center justify-between gap-4">

                <div className="flex items-center gap-2 text-gray-700 font-semibold text-xs uppercase tracking-wider">
                    <Filter className="w-4 h-4 text-blue-600" />
                    <span>Generación de Reportes PDF</span>
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

                    {/* Separador */}
                    <div className="w-[1px] h-6 bg-gray-200 hidden sm:block" />

                    {/* Botón Único de Descarga PDF */}
                    <button
                        onClick={handleGenerarPDF}
                        className="inline-flex items-center gap-1.5 px-3.5 py-1.5 text-xs font-semibold text-white bg-red-600 hover:bg-red-700 rounded-lg transition-colors shadow-sm"
                        title="Exportar a PDF"
                    >
                        <FileText className="w-4 h-4" />
                        <span>Generar Reporte PDF</span>
                    </button>

                </div>
            </div>

            {/* Contenedor de la Tabla de Clientes */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full min-w-full text-left text-sm text-gray-600">

                        <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3.5 w-16">ID</th>
                            <th className="px-4 py-3.5 w-44">Identificación</th>
                            <th className="px-4 py-3.5 min-w-[180px]">Nombres / Apellidos</th>
                            <th className="px-4 py-3.5 min-w-[180px]">Contacto</th>
                            <th className="px-4 py-3.5 min-w-[200px]">Dirección</th>
                            <th className="px-4 py-3.5 w-36">Creación</th>
                            <th className="px-4 py-3.5 w-24 text-center">Acciones</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {clientesFiltrados.length > 0 ? (
                            clientesFiltrados.map((cliente) => (
                                <tr key={cliente.cliente_id} className="hover:bg-gray-50/80 transition-colors">
                                    <td className="px-4 py-3.5 font-medium text-gray-900">
                                        #{cliente.cliente_id}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="font-medium text-gray-800">{cliente.identificacion}</div>
                                        <span className="inline-block px-2 py-0.5 text-[10px] font-semibold bg-gray-100 text-gray-600 rounded">
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
                                    <td className="px-4 py-3.5 truncate" title={cliente.direccion}>
                                        {cliente.direccion || '-'}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap text-xs text-gray-400">
                                        {cliente.fecha_creacion}
                                    </td>

                                    {/* Acciones */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">
                                            <button
                                                onClick={() => handleEditar(cliente)}
                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100"
                                                title="Editar cliente"
                                            >
                                                <Pencil className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleEliminar(cliente.cliente_id)}
                                                className="p-1.5 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors border border-red-100"
                                                title="Eliminar cliente"
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
                                    No se encontraron clientes registrados.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                {/* Pie de Tabla */}
                <div className="p-4 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de registros: {clientesFiltrados.length}</span>
                </div>
            </div>
        </div>
    );
}