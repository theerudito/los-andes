import React, {useEffect, useState} from 'react';
import {Plus, Pencil, Trash2} from 'lucide-react';
import {useModal} from "../../store/useModal.ts";
import {ModalLista} from "../../helpers/ModalLista.ts";
import {useUsuarios} from "../../store/useUsuarios.ts";
import {toast} from "sonner";

export default function PaginaUsuarios(): React.ReactElement {
    const {OpenModal} = useModal((state) => state);
    const {ObtenerUsuarios, ObtenerUsuario, EliminarUsuario, listar_usuario} = useUsuarios((state) => state);

    const [busqueda, setBusqueda] = useState<string>('');

    const usuariosFiltrados = listar_usuario.filter((u) =>
        u.nombres.toLowerCase().includes(busqueda.toLowerCase()) ||
        u.apellidos.toLowerCase().includes(busqueda.toLowerCase()) ||
        u.identificacion.includes(busqueda) ||
        u.email.toLowerCase().includes(busqueda.toLowerCase())
    );

    function VerUsuario(id: number) {
        OpenModal(ModalLista.modal_usuario)
        ObtenerUsuario(id)
    }

    useEffect(() => {
        ObtenerUsuarios();
        toast.dismiss();
    }, []);


    function Eliminar(id: number) {

        const toastId = `delete-confirm-${id}`;

        toast.custom(
            (t) => (
                <div
                    className="flex items-center justify-between gap-3 w-auto max-w-[calc(100vw-2rem)] sm:max-w-xs bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 px-3 py-2 rounded-xl shadow-lg text-xs select-none transition-all">
                    <span
                        className="text-zinc-700 dark:text-zinc-300 font-medium truncate shrink">¿Eliminar elemento?</span>

                    <div className="flex items-center gap-1.5 shrink-0">
                        <button
                            onClick={() => toast.dismiss(toastId)}
                            className="px-2.5 py-1 text-zinc-500 hover:text-zinc-800 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg font-medium transition-colors cursor-pointer"
                        >
                            No
                        </button>

                        <button
                            onClick={() => {
                                toast.dismiss(toastId);
                                EliminarUsuario(id)
                            }}
                            className="px-2.5 py-1 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-lg shadow-sm shadow-red-500/20 active:scale-95 transition-all cursor-pointer"
                        >
                            Sí
                        </button>
                    </div>
                </div>
            ),
            {
                id: toastId,
                duration: Infinity
            }
        );
    }

    return (
        <div className="space-y-6 w-full">
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 w-full">
                <div>
                    <h1 className="text-2xl font-bold text-gray-800">Gestión de Usuarios</h1>
                    <p className="text-xs text-gray-500 mt-0.5">Control de cuentas y accesos al sistema</p>
                </div>
                <div className="mt-5 pt-4 border-t border-gray-100 flex flex-wrap items-center justify-between gap-3">
                    <div className="relative flex-1 min-w-[240px] max-w-md">
                        <input
                            type="text"
                            placeholder="Buscar por nombre, cédula o email..."
                            value={busqueda}
                            onChange={(e) => setBusqueda(e.target.value)}
                            className="w-full px-3.5 py-2 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all"
                        />
                    </div>
                    <div className="flex flex-wrap items-center gap-2">
                        <button
                            onClick={() => OpenModal(ModalLista.modal_usuario)}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                        >
                            <Plus className="w-3.5 h-3.5"/>
                            <span>Nuevo Usuario</span>
                        </button>
                    </div>
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
                            <th className="px-4 py-3.5 min-w-[180px]">Email</th>
                            <th className="px-4 py-3.5 w-24 text-center">Acciones</th>
                        </tr>
                        </thead>

                        <tbody className="divide-y divide-gray-100">
                        {usuariosFiltrados.length > 0 ? (
                            usuariosFiltrados.map((usuario) => (
                                <tr key={usuario.usuario_id} className="hover:bg-gray-50/80 transition-colors">
                                    <td className="px-4 py-3.5 font-medium text-gray-900">
                                        #{usuario.usuario_id}
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap font-medium text-gray-800">
                                        ({usuario.tipo_identificacion}) {usuario.identificacion}
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap font-medium text-gray-800">
                                        {usuario.nombres} {usuario.apellidos}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap text-gray-600">
                                        {usuario.email}
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">
                                            <button
                                                onClick={() => VerUsuario(usuario.usuario_id)}
                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100"
                                                title="Editar usuario"
                                            >
                                                <Pencil className="w-4 h-4"/>
                                            </button>
                                            <button
                                                onClick={() => Eliminar(usuario.usuario_id)}
                                                className="p-1.5 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors border border-red-100"
                                                title="Eliminar usuario"
                                            >
                                                <Trash2 className="w-4 h-4"/>
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan={8} className="px-4 py-8 text-center text-gray-400">
                                    No se encontraron usuarios registrados.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                <div
                    className="p-4 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de registros: {usuariosFiltrados.length}</span>
                </div>
            </div>
        </div>
    );
}