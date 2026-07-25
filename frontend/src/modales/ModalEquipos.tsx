import React, {useEffect} from 'react';
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
import {useEquipos} from "../store/useEquipos.ts";
import {useMarcas} from "../store/useMarcas.ts";

export default function ModalEquipos(): React.ReactElement | null {
    const { modalName, CloseModal, OpenModal } = useModal((state) => state);
    const { form_equipo, EnviarEquipo } = useEquipos((state) => state);
    const { ObtenerMarcas, listar_marca } = useMarcas((state) => state);

    useEffect(() => {
        ObtenerMarcas();
    }, []);

    const handleChangeInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {name} = e.target;

        const value = e.target.value;

        useEquipos.setState((state) => {
            return {
                form_equipo: {
                    ...state.form_equipo,
                    [name]: value
                },
            };
        });
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

    if (modalName !== ModalLista.modal_equipo) return null;

    return (
        <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-md flex items-center justify-center z-50 p-4 transition-all duration-300">
            <div className="w-full max-w-2xl bg-white rounded-xl shadow-2xl overflow-hidden flex flex-col border border-slate-100 animate-in fade-in zoom-in-95 duration-200">

                <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-5 py-4 flex justify-between items-center shrink-0 shadow-sm">
                    <div className="flex items-center gap-2">
                        <Laptop size={18} />
                        <h2 className="font-semibold tracking-wide text-sm md:text-base">
                            Registrar Nuevo Equipo
                        </h2>
                    </div>
                    <button
                        type="button"
                        className="cursor-pointer hover:bg-white/20 transition-all rounded-full p-1.5 active:scale-95"
                        onClick={CloseModal}
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
                                    type="text"
                                    maxLength={13}
                                    placeholder="1721457494001"
                                    required
                                    className="flex-1 min-w-0 px-3 text-xs text-slate-800 outline-none font-mono font-bold tracking-wider uppercase h-full"
                                />

                                <button
                                    type="button"
                                    className="cursor-pointer bg-blue-600 hover:bg-blue-700 text-white px-3 flex items-center justify-center transition-all active:scale-95 h-full shrink-0"
                                    title="Registrar Nuevo Cliente"
                                >
                                    <Plus size={15} />
                                </button>

                                <button
                                    type="button"
                                    className="cursor-pointer bg-orange-600 text-white hover:orange-700 px-3 flex items-center justify-center transition-all active:scale-95 h-full shrink-0"
                                    title="Limpiar"
                                >
                                    <EraserIcon size={15} />
                                </button>

                                <button
                                    type="button"
                                    className="cursor-pointer bg-slate-800 hover:bg-slate-900 text-white px-3 flex items-center justify-center transition-all active:scale-95 h-full shrink-0"
                                    title="Buscar Cliente por Cédula / RUC"
                                >
                                    <Search size={15} />
                                </button>
                            </div>

                            <div className="h-10 flex items-center px-3 text-xs text-slate-700 font-bold bg-white rounded-lg border border-slate-300 shadow-sm w-full truncate">
                              Nombre Cliente
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
                                    value={form_equipo.marca_id}
                                    required
                                    className="flex-1 min-w-0 px-3 text-xs text-slate-800 outline-none font-semibold uppercase bg-white cursor-pointer h-full"
                                >
                                    <option disabled>SELECIONE UNA MARCA</option>
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
                                    className="cursor-pointer bg-blue-600 hover:bg-blue-700 text-white px-3 flex items-center justify-center transition-all active:scale-95 h-full shrink-0"
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
                                value={form_equipo.tipo_equipo}
                                onChange={handleChangeInput}
                                type="text"
                                name="tipo_equipo"
                                placeholder="Ej: LAPTOP, PC, IMPRESORA"
                                required
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
                                value={form_equipo.modelo}
                                onChange={handleChangeInput}
                                type="text"
                                name="modelo"
                                placeholder="Ej: MPS, THINKPAD E14"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <FileText size={14} className="text-blue-600" /> Número de Serie
                            </label>
                            <input
                                value={form_equipo.numero_serie}
                                onChange={handleChangeInput}
                                type="text"
                                name="numero_serie"
                                placeholder="Ej: 123456A"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-mono font-medium shadow-sm"
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                            <Package size={14} className="text-blue-600" /> Accesorios Dejados
                        </label>
                        <input
                            value={form_equipo.accesorios}
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
                                value={form_equipo.descripcion_problema}
                                onChange={handleChangeTextArea}
                                name="descripcion_problema"
                                rows={2}
                                placeholder="Ej: FORMATEO, NO ENCIENDE"
                                required
                                className="w-full p-2.5 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 uppercase placeholder-slate-400 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium resize-none shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <FileText size={14} className="text-blue-600" /> Observaciones Adicionales
                            </label>
                            <textarea
                                value={form_equipo.observacion}
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
                            <input
                                value={form_equipo.fecha_recepcion}
                                onChange={handleChangeInput}
                                type="datetime"
                                name="fecha_recepcion"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-xs font-semibold text-slate-700 mb-1 flex items-center gap-1.5">
                                <Calendar size={14} className="text-blue-600" /> Fecha Estimada Entrega
                            </label>
                            <input
                                value={form_equipo.fecha_estimada_entrega}
                                onChange={handleChangeInput}
                                type="datetime"
                                name="fecha_estimada_entrega"
                                required
                                className="w-full h-10 px-3 bg-white border border-slate-300 rounded-lg text-xs text-slate-800 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-all font-medium shadow-sm"
                            />
                        </div>
                    </div>

                    <div className="pt-3 border-t border-slate-200/80 flex items-center justify-end gap-2 shrink-0">
                        <button
                            type="button"
                            onClick={CloseModal}
                            className="cursor-pointer px-4 py-2 text-xs font-semibold text-slate-600 bg-slate-200 hover:bg-slate-300 rounded-lg transition-all active:scale-95"
                        >
                            Cancelar
                        </button>

                        <button
                            onClick={EnviarEquipo}
                            type="submit"
                            className="cursor-pointer inline-flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 rounded-lg shadow-sm transition-all active:scale-95"
                        >
                            <Save size={15} />
                            <span>Guardar Equipo</span>
                        </button>
                    </div>

                </div>

            </div>
        </div>
    );
}