import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Plus,
    Pencil,
    Trash2,
    History,
    CreditCard,
    Tag,
    PackageCheck,
    FileText,
    Calendar,
    Filter,
    PrinterCheck,
    EraserIcon,
    Search,
    RotateCcw
} from 'lucide-react';
import { useModal } from "../../store/useModal.ts";
import { ModalLista } from "../../helpers/ModalLista.ts";
import { useEquipos } from "../../store/useEquipos.ts";
import type { RPT_Equipos } from "../../modelos/equipos.ts";
import ModalMarcas from "../../modales/ModalMarcas.tsx";
import ReactDatePicker from "react-datepicker";
import { formatearFecha } from "../../helpers/formatearFecha.ts";
import { useClientes } from "../../store/useClientes.ts";
import {toast} from "sonner";

const estadosOpciones = [
    { estado_id: 0, nombre: "Todos" },
    { estado_id: 1, nombre: "Recibido" },
    { estado_id: 2, nombre: "En diagnóstico" },
    { estado_id: 3, nombre: "Esperando repuestos" },
    { estado_id: 4, nombre: "En reparación" },
    { estado_id: 5, nombre: "Listo para entrega" },
    { estado_id: 6, nombre: "Entregado" },
    { estado_id: 7, nombre: "Cancelado" }
];

export default function PaginaEquipos(): React.ReactElement {
    const OpenModal = useModal((state) => state.OpenModal);
    const { ObtenerClientePorIdentifiacion, form_cliente } = useClientes((state) => state);
    const { ObtenerEquipos, ObtenerEquipo, EliminarEquipo, DescargarPdf, DescargarOrdenPdf, listar_equipos } = useEquipos((state) => state);

    const [identificacion, setIdentificacion] = useState<string>("");
    const [estado, setEstado] = useState<number>(0);
    const [busqueda, setBusqueda] = useState<string>('');
    const [periodoSelect, setPeriodoSelect] = useState<string>("hoy");
    const [fechaDesde, setFechaDesde] = useState<Date | null>(new Date());
    const [fechaHasta, setFechaHasta] = useState<Date | null>(new Date());

    const navigate = useNavigate();

    useEffect(() => {
        ObtenerEquipos();
        toast.dismiss();
    }, []);

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

    function VerEquipo(id: number) {
        OpenModal(ModalLista.modal_equipo);
        ObtenerEquipo(id);
    }

    function VerReporte() {
        const obj: RPT_Equipos = {
            fecha_desde: formatearFecha(fechaDesde),
            fecha_hasta: formatearFecha(fechaHasta),
            cliente_id: identificacion !== "" ? form_cliente.cliente_id : 0,
            estado: estado
        };
        DescargarPdf(obj);
    }

    const handleRangoChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const opcion = e.target.value;
        setPeriodoSelect(opcion);
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

    async function ObtenerCliente() {
        if (!identificacion.trim()) {
            toast.error("Ingrese un número de cédula o RUC");
            return;
        }

        await ObtenerClientePorIdentifiacion(identificacion);

        if (form_cliente.identificacion !== ""){
            toast.success("Cliente cargado correctamente");
        } else {
            toast.error("No se encontró ningún cliente registrado con esa identificación");
        }
    }

    function ResetField() {
        setIdentificacion("");
        form_cliente.identificacion = "";
        form_cliente.nombres = "";
        form_cliente.apellidos = "";
    }

    function Reset() {
        setIdentificacion("");
        setEstado(0);
        setPeriodoSelect("hoy");
        setFechaDesde(new Date());
        setFechaHasta(new Date());
        form_cliente.identificacion = "";
        form_cliente.nombres = "";
        form_cliente.apellidos = "";
    }

    function Eliminar (id:number){
        toast.custom(
            (t) => (
                <div className="flex items-center justify-between gap-3 w-auto max-w-[calc(100vw-2rem)] sm:max-w-xs bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 px-3 py-2 rounded-xl shadow-lg text-xs select-none transition-all">
                    <span className="text-zinc-700 dark:text-zinc-300 font-medium truncate shrink">¿Eliminar elemento?</span>

                    <div className="flex items-center gap-1.5 shrink-0">
                        <button
                            onClick={() => toast.dismiss(t)}
                            className="px-2.5 py-1 text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg font-medium transition-colors cursor-pointer"
                        >
                            No
                        </button>

                        <button
                            onClick={() => {
                                toast.dismiss(t);
                                EliminarEquipo(id)
                            }}
                            className="px-2.5 py-1 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-lg shadow-sm shadow-red-500/20 active:scale-95 transition-all cursor-pointer"
                        >
                            Sí
                        </button>
                    </div>
                </div>
            ),
            { duration: Infinity }
        );
    }

    return (
        <>
            <ModalMarcas />
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
                                onClick={() => OpenModal(ModalLista.modal_marca)}
                                className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-purple-600 hover:bg-purple-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                            >
                                <Tag className="w-3.5 h-3.5" />
                                <span>Nueva Marca</span>
                            </button>

                            <button
                                onClick={() => OpenModal(ModalLista.modal_equipo)}
                                className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm cursor-pointer"
                            >
                                <Plus className="w-3.5 h-3.5" />
                                <span>Nuevo Equipo</span>
                            </button>
                        </div>
                    </div>
                </div>

                <div className="bg-white p-5 rounded-xl shadow-sm border border-gray-100 w-full space-y-4">
                    <div className="flex items-center gap-2 text-gray-800 font-bold text-xs uppercase tracking-wider border-b border-gray-100 pb-3">
                        <Filter className="w-4 h-4 text-blue-600" />
                        <span>Generación de Reportes PDF de Equipos</span>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-12 gap-4 items-end">
                        <div className="xl:col-span-3 flex flex-col gap-1.5">
                            <label className="text-xs font-semibold text-gray-600 flex items-center gap-1 truncate">
                                <Filter className="w-3.5 h-3.5 text-gray-400 shrink-0" />
                                {!form_cliente.identificacion ? "Cliente" : `Cliente: ${form_cliente.nombres} ${form_cliente.apellidos}`}
                            </label>
                            <div className="flex items-center h-9">
                                <input
                                    value={identificacion}
                                    onChange={(e) => setIdentificacion(e.target.value)}
                                    name="identificacion"
                                    type="number"
                                    placeholder="Identificación / Cédula"
                                    className="w-full h-full px-3 bg-gray-50 border border-r-0 border-gray-200 rounded-l-lg text-xs text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:bg-white transition-all"
                                />

                                <button
                                    onClick={ResetField}
                                    type="button"
                                    className="h-full bg-orange-500 hover:bg-orange-600 active:bg-orange-700 text-white px-3 flex items-center justify-center transition-all shrink-0 cursor-pointer"
                                    title="Limpiar datos de cliente"
                                >
                                    <EraserIcon size={14} />
                                </button>

                                <button
                                    onClick={ObtenerCliente}
                                    type="button"
                                    className="h-full bg-slate-800 hover:bg-slate-900 active:bg-slate-950 text-white px-3 rounded-r-lg flex items-center justify-center transition-all shrink-0 cursor-pointer"
                                    title="Buscar Cliente por Cédula / RUC"
                                >
                                    <Search size={14} />
                                </button>
                            </div>
                        </div>

                        <div className="xl:col-span-2 flex flex-col gap-1.5">
                            <label className="text-xs font-semibold text-gray-600 flex items-center gap-1">
                                <Filter className="w-3.5 h-3.5 text-gray-400" /> Estado
                            </label>
                            <select
                                name="estado_id"
                                value={estado}
                                required
                                onChange={(e) => setEstado(Number(e.target.value))}
                                className="w-full h-9 px-3 bg-gray-50 border border-gray-200 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium uppercase cursor-pointer"
                            >
                                {estadosOpciones.map((e) => (
                                    <option key={e.estado_id} value={e.estado_id}>
                                        {e.nombre}
                                    </option>
                                ))}
                            </select>
                        </div>

                        <div className="xl:col-span-2 flex flex-col gap-1.5">
                            <label className="text-xs font-semibold text-gray-600 flex items-center gap-1">
                                <Filter className="w-3.5 h-3.5 text-gray-400" /> Período
                            </label>
                            <select
                                value={periodoSelect}
                                onChange={handleRangoChange}
                                className="w-full h-9 px-2.5 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 font-medium focus:outline-none focus:border-blue-500 cursor-pointer"
                            >
                                <option value="hoy">Hoy</option>
                                <option value="mes_actual">Mes actual</option>
                                <option value="mes_anterior">Mes anterior</option>
                                <option value="anio_actual">Año actual</option>
                                <option value="anio_anterior">Año anterior</option>
                            </select>
                        </div>

                        <div className="xl:col-span-3 flex items-center gap-2">
                            <div className="flex flex-col gap-1.5 w-1/2">
                                <label className="text-xs font-semibold text-gray-600 flex items-center gap-1">
                                    <Calendar className="w-3.5 h-3.5 text-gray-400" /> Desde
                                </label>
                                <ReactDatePicker
                                    selected={fechaDesde}
                                    onChange={(date: Date | null) => setFechaDesde(date)}
                                    dateFormat="dd/MM/yyyy"
                                    locale="es"
                                    popperClassName="!z-[9999]"
                                    className="w-full h-9 px-2 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 font-mono font-medium text-center focus:outline-none focus:border-blue-500 cursor-pointer"
                                />
                            </div>

                            <div className="flex flex-col gap-1.5 w-1/2">
                                <label className="text-xs font-semibold text-gray-600 flex items-center gap-1">
                                    <Calendar className="w-3.5 h-3.5 text-gray-400" /> Hasta
                                </label>
                                <ReactDatePicker
                                    selected={fechaHasta}
                                    onChange={(date: Date | null) => setFechaHasta(date)}
                                    dateFormat="dd/MM/yyyy"
                                    locale="es"
                                    popperClassName="!z-[9999]"
                                    className="w-full h-9 px-2 bg-gray-50 border border-gray-200 rounded-lg text-xs text-gray-800 font-mono font-medium text-center focus:outline-none focus:border-blue-500 cursor-pointer"
                                />
                            </div>
                        </div>

                        <div className="xl:col-span-2 flex items-center gap-2">
                            <button
                                onClick={Reset}
                                type="button"
                                className="w-1/2 h-9 inline-flex items-center justify-center gap-1.5 px-2 text-xs font-semibold text-gray-700 bg-gray-100 hover:bg-gray-200 active:scale-95 rounded-lg transition-all shadow-sm cursor-pointer border border-gray-200"
                                title="Limpiar todos los filtros"
                            >
                                <RotateCcw className="w-3.5 h-3.5" />
                                <span>Reset</span>
                            </button>

                            <button
                                onClick={VerReporte}
                                type="button"
                                className="w-1/2 h-9 inline-flex items-center justify-center gap-1.5 px-2 text-xs font-semibold text-white bg-red-600 hover:bg-red-700 active:scale-95 rounded-lg transition-all shadow-sm cursor-pointer"
                                title="Exportar a PDF"
                            >
                                <FileText className="w-3.5 h-3.5" />
                                <span>Generar PDF</span>
                            </button>
                        </div>
                    </div>
                </div>

                <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                    <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                        <table className="w-full text-left text-sm text-gray-600">
                            <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                            <tr>
                                <th className="px-4 py-3">Id</th>
                                <th className="px-4 py-3">Cliente</th>
                                <th className="px-4 py-3">Orden</th>
                                <th className="px-4 py-3">Equipo / Marca</th>
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
                                            #{equipo.equipo_id}
                                        </td>

                                        <td className="px-4 py-3 whitespace-nowrap">
                                            <div className="font-medium text-gray-800">{equipo.nombres} {equipo.apellidos}</div>
                                            <div className="text-[11px] text-gray-400">ID: #{equipo.cliente_id}</div>
                                        </td>

                                        <td className="px-4 py-3 whitespace-nowrap font-bold text-gray-900">
                                            N° {equipo.codigo}
                                        </td>

                                        <td className="px-4 py-3 whitespace-nowrap">
                                            <div className="font-semibold text-gray-800 text-sm">{equipo.tipo_equipo} {equipo.modelo}</div>
                                            <div className="flex items-center gap-2 mt-0.5">
                                                    <span className="px-1.5 py-0.5 text-[10px] font-semibold bg-gray-100 text-gray-600 rounded">
                                                        {equipo.marca}
                                                    </span>
                                                <span className="text-[11px] text-gray-400 font-mono">S/N: {equipo.numero_serie}</span>
                                            </div>
                                        </td>

                                        <td className="px-4 py-3 max-w-xs truncate" title={equipo.descripcion_problema}>
                                            <div className="text-gray-800 truncate">{equipo.descripcion_problema}</div>
                                            {equipo.accesorios && (
                                                <div className="text-[11px] text-gray-400 truncate" title={`Accesorios: ${equipo.accesorios}`}>
                                                    Acc: {equipo.accesorios}
                                                </div>
                                            )}
                                        </td>

                                        <td className="px-4 py-3 whitespace-nowrap">
                                                <span className={`inline-block px-2 py-0.5 text-[11px] font-semibold rounded-full border ${getBadgeEstado(equipo.estado)}`}>
                                                    {equipo.estado}
                                                </span>
                                        </td>

                                        <td className="px-4 py-3 whitespace-nowrap text-[11px] text-gray-500">
                                            <div>Rec: {equipo.fecha_recepcion}</div>
                                            <div className="text-gray-400">Est: {equipo.fecha_estimada_entrega}</div>
                                        </td>

                                        <td className="px-4 py-3 whitespace-nowrap text-center">
                                            <div className="flex items-center justify-center gap-1.5">
                                                <button
                                                    onClick={() => DescargarOrdenPdf(equipo.equipo_id)}
                                                    className="p-1.5 text-emerald-600 bg-green-50 hover:bg-green-100 rounded-lg transition-colors border border-green-100"
                                                    title="Imprimir Orden"
                                                >
                                                    <PrinterCheck className="w-4 h-4" />
                                                </button>
                                                <button
                                                    onClick={() => handleVerHistorial(equipo.equipo_id)}
                                                    className="p-1.5 text-yellow-600 bg-yellow-50 hover:bg-yellow-100 rounded-lg transition-colors border border-yellow-100"
                                                    title="Ver Historial"
                                                >
                                                    <History className="w-4 h-4" />
                                                </button>
                                                <button
                                                    onClick={() => handleGestionarPagos(equipo.equipo_id)}
                                                    className="p-1.5 text-emerald-600 bg-emerald-50 hover:bg-emerald-100 rounded-lg transition-colors border border-emerald-100"
                                                    title="Gestionar Pagos"
                                                >
                                                    <CreditCard className="w-4 h-4" />
                                                </button>
                                                <button
                                                    onClick={() => handleGestionarEntrega(equipo.equipo_id)}
                                                    className="p-1.5 text-amber-600 bg-amber-50 hover:bg-amber-100 rounded-lg transition-colors border border-amber-100"
                                                    title="Gestionar Entrega"
                                                >
                                                    <PackageCheck className="w-4 h-4" />
                                                </button>
                                                <button
                                                    onClick={() => VerEquipo(equipo.equipo_id)}
                                                    className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100"
                                                    title="Editar equipo"
                                                >
                                                    <Pencil className="w-4 h-4" />
                                                </button>
                                                <button
                                                    onClick={() => Eliminar(equipo.equipo_id)}
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
                                    <td colSpan={8} className="px-4 py-8 text-center text-gray-400">
                                        No se encontraron equipos registrados.
                                    </td>
                                </tr>
                            )}
                            </tbody>
                        </table>
                    </div>

                    <div className="p-3 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                        <span>Total de registros: {equiposFiltrados.length}</span>
                    </div>
                </div>
            </div>
        </>
    );
}