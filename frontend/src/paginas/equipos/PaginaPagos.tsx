import React, { useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import {
    Pencil,
    Search,
    RotateCcw,
    FileText,
    ArrowLeft,
    CreditCard,
    Laptop,
    X,
    CheckCircle2,
    Clock
} from 'lucide-react';

// Interface 1:1 con tu API de Go
export interface CuentaDTO {
    cuenta_id: number;
    equipo_id: number;
    equipo_codigo: string;
    costo_total: number;
    abono: number;
    saldo: number;
}

// Datos de prueba iniciales basados en tu respuesta JSON real
const cuentasIniciales: CuentaDTO[] = [
    {
        cuenta_id: 2,
        equipo_id: 2,
        equipo_codigo: "E00002",
        costo_total: 25,
        abono: 25,
        saldo: 0
    },
    {
        cuenta_id: 3,
        equipo_id: 3,
        equipo_codigo: "E00003",
        costo_total: 60,
        abono: 20,
        saldo: 40
    }
];

export default function PaginaPagos(): React.ReactElement {
    const navigate = useNavigate();
    const [searchParams, setSearchParams] = useSearchParams();
    const equipoIdParam = searchParams.get('equipo_id');

    const [cuentas] = useState<CuentaDTO[]>(cuentasIniciales);
    const [busqueda, setBusqueda] = useState<string>('');

    // Filtrado de cuentas por equipo_id (URL) + búsqueda general
    const cuentasFiltradas = cuentas.filter((c) => {
        const coincideEquipoId = equipoIdParam ? c.equipo_id === Number(equipoIdParam) : true;
        const coincideTexto =
            c.equipo_codigo.toLowerCase().includes(busqueda.toLowerCase()) ||
            c.cuenta_id.toString().includes(busqueda);

        return coincideEquipoId && coincideTexto;
    });

    // Totales acumulados
    const totalCosto = cuentasFiltradas.reduce((acc, c) => acc + c.costo_total, 0);
    const totalAbonos = cuentasFiltradas.reduce((acc, c) => acc + c.abono, 0);
    const totalSaldo = cuentasFiltradas.reduce((acc, c) => acc + c.saldo, 0);

    const handleRegresar = () => {
        navigate(-1);
    };

    const handleLimpiar = () => setBusqueda('');

    const handleQuitarFiltroEquipo = () => {
        searchParams.delete('equipo_id');
        setSearchParams(searchParams);
    };

    const handleBuscar = () => {
        console.log("Buscando cuentas/pagos:", { busqueda, equipoIdParam });
    };

    const handleEditar = (cuenta: CuentaDTO) => {
        console.log("Registrar/Editar cobro o abono para cuenta ID:", cuenta.cuenta_id);
        // Aquí abres tu modal de edición para registrar el abono o ajustar costo total
    };

    const handleGenerarPDFRecibo = (cuenta: CuentaDTO) => {
        console.log("Generando recibo PDF de pago para la cuenta ID:", cuenta.cuenta_id);
    };

    return (
        <div className="space-y-6 w-full">
            {/* Encabezado Superior */}
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-emerald-50 text-emerald-600 rounded-lg">
                            <CreditCard className="w-6 h-6" />
                        </div>
                        <div>
                            <h1 className="text-2xl font-bold text-gray-800">Gestión de Cuentas y Pagos</h1>
                            <p className="text-xs text-gray-500 mt-0.5">Control de costos totales, abonos y saldos pendientes</p>
                        </div>
                    </div>

                    {/* Botón Regresar */}
                    <button
                        onClick={handleRegresar}
                        className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 border border-gray-200 rounded-lg transition-colors shadow-sm shrink-0"
                        title="Volver a la vista anterior"
                    >
                        <ArrowLeft className="w-4 h-4" />
                        <span>Regresar</span>
                    </button>
                </div>

                {/* Indicador de Filtro por Equipo ID */}
                {equipoIdParam && (
                    <div className="mt-3 inline-flex items-center gap-2 px-3 py-1.5 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-lg text-xs font-semibold">
                        <Laptop className="w-4 h-4" />
                        <span>Filtrado por Equipo ID: #{equipoIdParam}</span>
                        <button
                            onClick={handleQuitarFiltroEquipo}
                            className="p-0.5 hover:bg-emerald-100 rounded-full transition-colors ml-1"
                            title="Mostrar todos los registros"
                        >
                            <X className="w-3.5 h-3.5" />
                        </button>
                    </div>
                )}

                {/* Barra de Búsqueda y Botones */}
                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">

                    {/* Input de Búsqueda */}
                    <div className="relative flex-1 min-w-[240px] max-w-md">
                        <input
                            type="text"
                            placeholder="Buscar por código de equipo..."
                            value={busqueda}
                            onChange={(e) => setBusqueda(e.target.value)}
                            className="w-full px-3.5 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500/20 focus:border-emerald-500 transition-all"
                        />
                    </div>

                    {/* Grupo de Botones */}
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
                    </div>

                </div>
            </div>

            {/* Tabla de Cuentas */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full min-w-full text-left text-sm text-gray-600">

                        <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3.5 w-16">ID</th>
                            <th className="px-4 py-3.5 w-36">Código Equipo</th>
                            <th className="px-4 py-3.5 w-32 text-right">Costo Total</th>
                            <th className="px-4 py-3.5 w-32 text-right">Abono</th>
                            <th className="px-4 py-3.5 w-32 text-right">Saldo</th>
                            <th className="px-4 py-3.5 w-36 text-center">Estado Pago</th>
                            <th className="px-4 py-3.5 w-24 text-center">Acciones</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {cuentasFiltradas.length > 0 ? (
                            cuentasFiltradas.map((cuenta) => (
                                <tr key={cuenta.cuenta_id} className="hover:bg-gray-50/80 transition-colors text-xs">

                                    {/* Cuenta ID */}
                                    <td className="px-4 py-3.5 font-bold text-gray-900">
                                        #{cuenta.cuenta_id}
                                    </td>

                                    {/* Código Equipo */}
                                    <td className="px-4 py-3.5 font-semibold text-gray-800 whitespace-nowrap">
                                        {cuenta.equipo_codigo}
                                    </td>

                                    {/* Costo Total */}
                                    <td className="px-4 py-3.5 font-semibold text-gray-800 text-right whitespace-nowrap">
                                        ${cuenta.costo_total.toFixed(2)}
                                    </td>

                                    {/* Abono */}
                                    <td className="px-4 py-3.5 font-bold text-emerald-600 text-right whitespace-nowrap">
                                        ${cuenta.abono.toFixed(2)}
                                    </td>

                                    {/* Saldo */}
                                    <td className="px-4 py-3.5 font-bold text-red-600 text-right whitespace-nowrap">
                                        ${cuenta.saldo.toFixed(2)}
                                    </td>

                                    {/* Estado del Pago Badge */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        {cuenta.saldo === 0 ? (
                                            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 text-[11px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 rounded-full">
                          <CheckCircle2 className="w-3 h-3" />
                          Pagado
                        </span>
                                        ) : (
                                            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 text-[11px] font-semibold bg-amber-50 text-amber-700 border border-amber-200 rounded-full">
                          <Clock className="w-3 h-3" />
                          Pendiente
                        </span>
                                        )}
                                    </td>

                                    {/* Acciones: PDF y Editar (para gestionar el pago/abono) */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">

                                            {/* Botón Recibo PDF */}
                                            <button
                                                onClick={() => handleGenerarPDFRecibo(cuenta)}
                                                className="p-1.5 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors border border-red-100"
                                                title="Descargar Comprobante / Recibo PDF"
                                            >
                                                <FileText className="w-4 h-4" />
                                            </button>

                                            {/* Botón Editar (Gestión de Cobro/Abono) */}
                                            <button
                                                onClick={() => handleEditar(cuenta)}
                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100"
                                                title="Gestionar Abono / Editar Cuenta"
                                            >
                                                <Pencil className="w-4 h-4" />
                                            </button>

                                        </div>
                                    </td>

                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={7} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron registros de cuentas o pagos.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                {/* Pie de Tabla con Resumen Monetario */}
                <div className="p-3 border-t border-gray-100 text-xs text-gray-600 flex flex-wrap justify-between items-center bg-gray-50/50 gap-2">
                    <span>Total registros: <strong className="text-gray-800">{cuentasFiltradas.length}</strong></span>

                    <div className="flex items-center gap-4 text-xs font-bold">
                        <div className="text-gray-700">Total Costo: <span className="text-gray-900">${totalCosto.toFixed(2)}</span></div>
                        <div className="text-emerald-700">Total Abonado: <span>${totalAbonos.toFixed(2)}</span></div>
                        <div className="text-red-600">Total Pendiente: <span>${totalSaldo.toFixed(2)}</span></div>
                    </div>
                </div>
            </div>
        </div>
    );
}