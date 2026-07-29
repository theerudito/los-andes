import React, { useEffect } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import {
    Pencil,
    ArrowLeft,
    CreditCard,
    CheckCircle2,
    Clock,
    PrinterCheck
} from 'lucide-react';
import { useModal } from '../../store/useModal.ts';
import { ModalLista } from '../../helpers/ModalLista.ts';
import { usePagos } from "../../store/usePagos.ts";

export default function PaginaPagos(): React.ReactElement {
    const { OpenModal } = useModal((state) => state);
    const { ObtenerPagos, ObtenerPago, DescargarPdf, listar_pagos, setEquipoId } = usePagos((state) => state);
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    const equipo_id = searchParams.get('equipo_id');

    useEffect(() => {
        if (equipo_id) {
            setEquipoId(Number(equipo_id));
        }
        ObtenerPagos();
    }, [equipo_id]);

    const totalCosto = listar_pagos.reduce((acc, c) => acc + c.costo_total, 0);
    const totalAbonos = listar_pagos.reduce((acc, c) => acc + c.abono, 0);
    const totalSaldo = listar_pagos.reduce((acc, c) => acc + c.saldo, 0);

    const handleRegresar = () => {
        navigate(-1);
    };

    async function VerPago(pagoId: number) {
        await ObtenerPago(pagoId);
        OpenModal(ModalLista.modal_pago);
    }

    return (
        <div className="space-y-6 w-full">
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-emerald-50 text-emerald-600 rounded-lg">
                            <CreditCard className="w-6 h-6" />
                        </div>
                        <div>
                            <h1 className="text-2xl font-bold text-gray-800">
                                Gestión de Cuentas y Pagos {equipo_id && `(Equipo #${equipo_id})`}
                            </h1>
                            <p className="text-xs text-gray-500 mt-0.5">Control de costos totales, abonos y saldos pendientes</p>
                        </div>
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
            </div>

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
                        {listar_pagos.length > 0 ? (
                            listar_pagos.map((item) => (
                                <tr key={item.cuenta_id} className="hover:bg-gray-50/80 transition-colors text-xs">

                                    <td className="px-4 py-3.5 font-bold text-gray-900">
                                        #{item.cuenta_id}
                                    </td>

                                    {/* 🛠️ FIX: Reemplazado item.equipo e item.codigo por item.equipo_codigo */}
                                    <td className="px-4 py-3.5 font-semibold text-gray-800 whitespace-nowrap">
                                        ORDEN: {item.equipo_codigo}
                                    </td>

                                    <td className="px-4 py-3.5 font-semibold text-gray-800 text-right whitespace-nowrap">
                                        ${item.costo_total.toFixed(2)}
                                    </td>

                                    <td className="px-4 py-3.5 font-bold text-emerald-600 text-right whitespace-nowrap">
                                        ${item.abono.toFixed(2)}
                                    </td>

                                    <td className="px-4 py-3.5 font-bold text-red-600 text-right whitespace-nowrap">
                                        ${item.saldo.toFixed(2)}
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        {item.estado === "Pendiente" ? (
                                            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 text-[11px] font-semibold bg-amber-50 text-amber-700 border border-amber-200 rounded-full">
                                                    <Clock className="w-3 h-3" />
                                                    Pendiente
                                                </span>
                                        ) : (
                                            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 text-[11px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 rounded-full">
                                                    <CheckCircle2 className="w-3 h-3" />
                                                    Pagado
                                                </span>
                                        )}
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">

                                            {(item.costo_total > 0 || item.abono > 0) && (
                                                <button
                                                    onClick={() => DescargarPdf(item.cuenta_id)}
                                                    className="p-1.5 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors border border-red-100 cursor-pointer"
                                                    title="Descargar Comprobante / Recibo PDF"
                                                >
                                                    <PrinterCheck className="w-4 h-4" />
                                                </button>
                                            )}

                                            {item.estado === "Pendiente" && (
                                                <button
                                                    onClick={() => VerPago(item.cuenta_id)}
                                                    className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100 cursor-pointer"
                                                    title="Gestionar Abono / Editar Cuenta"
                                                >
                                                    <Pencil className="w-4 h-4" />
                                                </button>
                                            )}

                                        </div>
                                    </td>

                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={7} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron registros de cuentas o pagos para este equipo.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                <div className="p-3 border-t border-gray-100 text-xs text-gray-600 flex flex-wrap justify-between items-center bg-gray-50/50 gap-2">
                    <span>Total registros: <strong className="text-gray-800">{listar_pagos.length}</strong></span>

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