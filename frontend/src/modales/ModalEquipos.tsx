import React, { useEffect } from 'react';
import {
    X,
    Save,
    Laptop,
    Search,
    Calendar,
    Tag,
    User,
    FileText,
    Wrench,
    Package,
    Plus,
    EraserIcon
} from 'lucide-react';
import { useModal } from '../store/useModal.ts';
import { ModalLista } from '../helpers/ModalLista.ts';
import { useEquipos } from "../store/useEquipos.ts";
import { useMarcas } from "../store/useMarcas.ts";
import { useClientes } from "../store/useClientes.ts";
import { toast } from "sonner";
import ReactDatePicker from "react-datepicker";

export default function ModalEquipos(): React.ReactElement | null {
    const { modalName, CloseModal, OpenModal } = useModal((state) => state);
    const { form_equipo, EnviarEquipo } = useEquipos((state) => state);
    const { ObtenerMarcas, listar_marca } = useMarcas((state) => state);
    const { form_cliente, ObtenerClientePorIdentifiacion, clienteId } = useClientes((state) => state);

    const formatDateToDbStr = (date: Date | null): string => {
        if (!date) return '';
        const pad = (n: number) => n.toString().padStart(2, '0');
        const year = date.getFullYear();
        const month = pad(date.getMonth() + 1);
        const day = pad(date.getDate());
        const hours = pad(date.getHours());
        const minutes = pad(date.getMinutes());
        const seconds = pad(date.getSeconds());
        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    };

    const parseStringToDate = (dateStr?: string): Date => {
        if (!dateStr) return new Date();
        const normalized = dateStr.replace(' ', 'T');
        const date = new Date(normalized);
        return isNaN(date.getTime()) ? new Date() : date;
    };

    const handleDateChange = (campo: 'fecha_recepcion' | 'fecha_estimada_entrega', date: Date | null) => {
        const formattedDate = formatDateToDbStr(date || new Date());
        useEquipos.setState((state) => ({
            form_equipo: {
                ...state.form_equipo,
                [campo]: formattedDate,
            },
        }));
    };

    const handleChangeInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;

        useEquipos.setState((state) => ({
            form_equipo: {
                ...state.form_equipo,
                [name]: value,
                cliente_id: clienteId
            },
        }));

        useClientes.setState((state) => ({
            form_cliente: {
                ...state.form_cliente,
                [name]: value
            },
        }));
    };

    const handleChangeTextArea = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
        const { name, value } = e.target;
        useEquipos.setState((state) => ({
            form_equipo: {
                ...state.form_equipo,
                [name]: value,
            },
        }));
    };

    const handleChangeSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const { name, value } = e.target;
        useEquipos.setState((state) => ({
            form_equipo: {
                ...state.form_equipo,
                [name]: Number(value)
            },
        }));
    };

    function LimpiarIdentificacion() {
        useClientes.setState({
            form_cliente: {
                cliente_id: 0,
                identificacion: "",
                tipo_identificacion: "",
                nombres: "",
                apellidos: "",
                telefono: "",
                email: "",
                direccion: "",
                fecha_creacion: "",
                fecha_modificacion: "",
            },
            clienteId: 0,
            isEditing: false,
        });

        useEquipos.setState((state) => ({
            form_equipo: {
                ...state.form_equipo,
                cliente_id: 0,
                identificacion: "",
                nombres: "",
                apellidos: "",
            }
        }));
    }

    function Clear() {
        CloseModal();
        useEquipos.setState({
            form_equipo: {
                equipo_id: 0,
                codigo: "",
                tipo_equipo: "",
                modelo: "",
                numero_serie: "",
                accesorios: "",
                descripcion_problema: "",
                observacion: "",
                fecha_recepcion: "",
                fecha_estimada_entrega: "",
                fecha_creacion: "",
                fecha_modificacion: "",
                marca_id: 0,
                cliente_id: 0,
                estado_id: 0,
                usuario_id: 0,
            },
            isEditing: false,
        });
        LimpiarIdentificacion();
    }

    function AbrirModalCliente() {
        OpenModal(ModalLista.modal_cliente);
    }

    async function ObtenerCliente() {
        if (!form_equipo.identificacion?.trim()) {
            toast.error("Ingrese un número de cédula o RUC");
            return;
        }

        await ObtenerClientePorIdentifiacion(form_equipo.identificacion);

        const clienteActual = useClientes.getState().form_cliente;
        const idClienteActual = useClientes.getState().clienteId;

        if (idClienteActual > 0 && clienteActual) {
            useEquipos.setState((state) => ({
                form_equipo: {
                    ...state.form_equipo,
                    cliente_id: clienteActual.cliente_id,
                    identificacion: clienteActual.identificacion,
                    nombres: clienteActual.nombres,
                    apellidos: clienteActual.apellidos,
                }
            }));

            toast.success("Cliente cargado correctamente");
        } else {
            toast.error("No se encontró ningún cliente registrado con esa identificación");
        }
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const currentClienteId = useClientes.getState().clienteId || form_equipo.cliente_id;

        if (!currentClienteId || currentClienteId <= 0) {
            toast.error("Debe buscar y seleccionar un cliente válido antes de guardar");
            return;
        }

        if (!form_equipo.marca_id || form_equipo.marca_id <= 0) {
            toast.error("Debe seleccionar una Marca para el equipo");
            return;
        }

        if (!form_equipo.tipo_equipo?.trim()) {
            toast.error("Ingrese el tipo de equipo (Ej: Laptop, PC)");
            return;
        }

        if (!form_equipo.modelo?.trim()) {
            toast.error("Ingrese el modelo del equipo");
            return;
        }

        if (!form_equipo.numero_serie?.trim()) {
            toast.error("Ingrese el número de serie del equipo");
            return;
        }

        if (!form_equipo.descripcion_problema?.trim()) {
            toast.error("Ingrese la descripción del problema");
            return;
        }

        if (!form_equipo.fecha_recepcion?.trim()) {
            toast.error("Seleccione la fecha de recepción");
            return;
        }

        if (!form_equipo.fecha_estimada_entrega?.trim()) {
            toast.error("Seleccione la fecha estimada de entrega");
            return;
        }

        try {
            await EnviarEquipo();
            Clear();
        } catch (error) {
            console.error("Error al procesar el guardado del equipo:", error);
        }
    };

    useEffect(() => {
        if (modalName === ModalLista.modal_equipo) {
            ObtenerMarcas();

            const hoyStr = formatDateToDbStr(new Date());
            useEquipos.setState((state) => ({
                form_equipo: {
                    ...state.form_equipo,
                    fecha_recepcion: state.form_equipo.fecha_recepcion || hoyStr,
                    fecha_estimada_entrega: state.form_equipo.fecha_estimada_entrega || hoyStr,
                }
            }));
        }
    }, [modalName]);

    if (modalName !== ModalLista.modal_equipo) return null;

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-50 p-4 transition-all duration-300">
            <form onSubmit={handleSubmit} className="w-full max-w-2xl bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <Laptop size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Registrar Nuevo Equipo
                        </h2>
                    </div>
                    <button
                        type="button"
                        className="cursor-pointer hover:bg-white/20 transition-all rounded-full p-1.5 active:scale-95 focus:outline-none"
                        onClick={Clear}
                    >
                        <X size={18} />
                    </button>
                </div>

                <div className="p-5 bg-slate-50/50 flex flex-col gap-4 max-h-[80vh] overflow-y-auto">

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <User size={14} className="text-blue-600" /> Identificación Cliente (Cédula / RUC)
                        </label>
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 items-center">

                            <div className="flex items-stretch shadow-sm rounded-lg overflow-hidden border border-slate-300 focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500 transition-all bg-white h-10 divide-x divide-slate-200 w-full">
                                <input
                                    value={form_equipo.identificacion || ''}
                                    onChange={handleChangeInput}
                                    onKeyDown={(e) => {
                                        if (e.key === 'Enter') {
                                            e.preventDefault();
                                            ObtenerCliente();
                                        }
                                    }}
                                    type="text"
                                    name="identificacion"
                                    maxLength={13}
                                    placeholder="1721457494001"
                                    className="flex-1 min-w-0 px-3 text-xs text-slate-800 outline-none font-mono font-bold tracking-wider uppercase h-full focus:outline-none"
                                />

                                <button
                                    onClick={AbrirModalCliente}
                                    type="button"
                                    className="cursor-pointer bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white px-3 flex items-center justify-center transition-all active:scale-95 h-full shrink-0 focus:outline-none"
                                    title="Registrar Nuevo Cliente"
                                >
                                    <Plus size={15} />
                                </button>

                                <button
                                    onClick={LimpiarIdentificacion}
                                    type="button"
                                    className="cursor-pointer bg-orange-600 hover:bg-orange-700 active:bg-orange-800 text-white px-3 flex items-center justify-center transition-all active:scale-95 h-full shrink-0 focus:outline-none"
                                    title="Limpiar datos de cliente"
                                >
                                    <EraserIcon size={15} />
                                </button>

                                <button
                                    onClick={ObtenerCliente}
                                    type="button"
                                    className="cursor-pointer bg-slate-800 hover:bg-slate-900 active:bg-slate-950 text-white px-3 flex items-center justify-center transition-all active:scale-95 h-full shrink-0 focus:outline-none"
                                    title="Buscar Cliente por Cédula / RUC"
                                >
                                    <Search size={15} />
                                </button>
                            </div>

                            <div className="h-10 flex items-center px-3 text-xs text-slate-700 font-bold bg-white rounded-lg border border-slate-300 shadow-sm w-full truncate">
                                {form_equipo.nombres || form_equipo.apellidos || form_cliente.nombres || form_cliente.apellidos
                                    ? `${form_equipo.nombres || form_cliente.nombres || ''} ${form_equipo.apellidos || form_cliente.apellidos || ''}`.trim()
                                    : ""}
                            </div>
                        </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Tag size={14} className="text-blue-600" /> Marca
                            </label>
                            <div className="flex items-stretch shadow-sm rounded-lg overflow-hidden border border-slate-300 focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500 transition-all bg-white h-10 divide-x divide-slate-200">
                                <select
                                    onChange={handleChangeSelect}
                                    name="marca_id"
                                    value={form_equipo.marca_id || 0}
                                    className="flex-1 min-w-0 px-3 text-xs text-slate-800 outline-none font-semibold uppercase bg-white cursor-pointer h-full focus:outline-none"
                                >
                                    <option value={0} disabled>SELECCIONE UNA MARCA</option>
                                    {
                                        listar_marca.map((item) => (
                                            <option value={item.marca_id} key={item.marca_id}>
                                                {item.nombre}
                                            </option>
                                        ))
                                    }
                                </select>

                                <button
                                    onClick={() => OpenModal(ModalLista.modal_marca)}
                                    type="button"
                                    className="cursor-pointer bg-blue-600 hover:bg-blue-700 active:bg-blue-800 text-white px-3 flex items-center justify-center transition-all active:scale-95 h-full shrink-0 focus:outline-none"
                                    title="Agregar Nueva Marca"
                                >
                                    <Plus size={15} />
                                </button>
                            </div>
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Laptop size={14} className="text-blue-600" /> Tipo Equipo
                            </label>
                            <input
                                value={form_equipo.tipo_equipo || ''}
                                onChange={handleChangeInput}
                                type="text"
                                name="tipo_equipo"
                                placeholder="Ej: LAPTOP, PC, IMPRESORA"
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Package size={14} className="text-blue-600" /> Modelo
                            </label>
                            <input
                                value={form_equipo.modelo || ''}
                                onChange={handleChangeInput}
                                type="text"
                                name="modelo"
                                placeholder="Ej: MPS, THINKPAD E14"
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <FileText size={14} className="text-blue-600" /> Número de Serie
                            </label>
                            <input
                                value={form_equipo.numero_serie || ''}
                                onChange={handleChangeInput}
                                type="text"
                                name="numero_serie"
                                placeholder="Ej: 123456A"
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-mono font-medium shadow-sm"
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <Package size={14} className="text-blue-600" /> Accesorios Dejados
                        </label>
                        <input
                            value={form_equipo.accesorios || ''}
                            onChange={handleChangeInput}
                            type="text"
                            name="accesorios"
                            placeholder="Ej: CARGADOR, MOUSE, MOCHILA"
                            className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                        />
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Wrench size={14} className="text-blue-600" /> Descripción del Problema
                            </label>
                            <textarea
                                value={form_equipo.descripcion_problema || ''}
                                onChange={handleChangeTextArea}
                                name="descripcion_problema"
                                rows={2}
                                placeholder="Ej: FORMATEO, NO ENCIENDE"
                                className="w-full p-2.5 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium resize-none shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <FileText size={14} className="text-blue-600" /> Observaciones Adicionales
                            </label>
                            <textarea
                                value={form_equipo.observacion || ''}
                                onChange={handleChangeTextArea}
                                name="observacion"
                                rows={2}
                                placeholder="Ej: RAYONES EN LA TAPA"
                                className="w-full p-2.5 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium resize-none shadow-sm"
                            />
                        </div>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Calendar size={14} className="text-blue-600" /> Fecha Recepción
                            </label>
                            <ReactDatePicker
                                selected={parseStringToDate(form_equipo.fecha_recepcion)}
                                onChange={(date: Date | null) => handleDateChange('fecha_recepcion', date)}
                                showTimeSelect
                                timeIntervals={15}
                                timeCaption="Hora"
                                dateFormat="dd/MM/yyyy HH:mm"
                                locale="es"
                                wrapperClassName="w-full"
                                popperClassName="!z-[9999]"
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 font-medium outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all shadow-sm cursor-pointer"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Calendar size={14} className="text-blue-600" /> Fecha Estimada Entrega
                            </label>
                            <ReactDatePicker
                                selected={parseStringToDate(form_equipo.fecha_estimada_entrega)}
                                onChange={(date: Date | null) => handleDateChange('fecha_estimada_entrega', date)}
                                showTimeSelect
                                timeIntervals={15}
                                timeCaption="Hora"
                                dateFormat="dd/MM/yyyy HH:mm"
                                locale="es"
                                wrapperClassName="w-full"
                                popperClassName="!z-[9999]"
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 font-medium outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all shadow-sm cursor-pointer"
                            />
                        </div>
                    </div>

                    <div className="pt-3 border-t border-slate-200/80 flex items-center justify-end gap-2 shrink-0">
                        <button
                            type="button"
                            onClick={Clear}
                            className="cursor-pointer px-4 py-2 text-xs font-semibold text-slate-600 bg-slate-200 hover:bg-slate-300 active:bg-slate-400 rounded-lg transition-all active:scale-95 focus:outline-none focus:ring-0"
                        >
                            Cancelar
                        </button>

                        <button
                            type="submit"
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 active:from-blue-800 active:to-indigo-800 rounded-lg shadow-sm transition-all active:scale-95 focus:outline-none focus:ring-0"
                        >
                            <Save size={15} />
                            <span>Guardar Equipo</span>
                        </button>
                    </div>

                </div>

            </form>
        </div>
    );
}