import React, { useState } from 'react';
import { Plus, Pencil, Trash2, Search, RotateCcw, Shield } from 'lucide-react';
import {useModal} from "../../store/useModal.ts";
import {ModalLista} from "../../helpers/ModalLista.ts";

export interface Usuario {
    usuario_id: number;
    identificacion: string;
    tipo_identificacion: string;
    nombres: string;
    apellidos: string;
    email: string;
    password?: string;
    activo: boolean;
    fecha_creacion: string;
    fecha_modificacion: string;
    rol_id: number;
}

// Datos de prueba iniciales
const usuariosIniciales: Usuario[] = [
    {
        usuario_id: 1,
        identificacion: "1712345678",
        tipo_identificacion: "CEDULA",
        nombres: "Admin",
        apellidos: "Sistema",
        email: "admin@empresa.com",
        activo: true,
        fecha_creacion: "2026-01-01 08:00",
        fecha_modificacion: "2026-01-01 08:00",
        rol_id: 1 // Administrador
    },
    {
        usuario_id: 2,
        identificacion: "1798765432",
        tipo_identificacion: "CEDULA",
        nombres: "Técnico",
        apellidos: "Soporte",
        email: "tecnico@empresa.com",
        activo: false,
        fecha_creacion: "2026-02-15 10:30",
        fecha_modificacion: "2026-03-01 12:00",
        rol_id: 2 // Técnico
    }
];

export default function PaginaUsuarios(): React.ReactElement {
    const { OpenModal } = useModal((state) => state);


    const [usuarios, setUsuarios] = useState<Usuario[]>(usuariosIniciales);
    const [busqueda, setBusqueda] = useState<string>('');

    // Filtrado de usuarios
    const usuariosFiltrados = usuarios.filter((u) =>
        u.nombres.toLowerCase().includes(busqueda.toLowerCase()) ||
        u.apellidos.toLowerCase().includes(busqueda.toLowerCase()) ||
        u.identificacion.includes(busqueda) ||
        u.email.toLowerCase().includes(busqueda.toLowerCase())
    );

    const handleLimpiar = () => setBusqueda('');

    const handleBuscar = () => {
        console.log("Buscando usuarios:", busqueda);
    };

    const handleNuevo = () => {
        console.log("Registrar nuevo usuario");
    };

    const handleEditar = (usuario: Usuario) => {
        console.log("Editar usuario:", usuario);
    };

    const handleEliminar = (id: number) => {
        if (confirm('¿Estás seguro de que deseas eliminar este usuario?')) {
            setUsuarios(usuarios.filter((u) => u.usuario_id !== id));
        }
    };

    // Helper para mostrar etiqueta del rol
    const getNombreRol = (rolId: number) => {
        switch (rolId) {
            case 1:
                return 'Administrador';
            case 2:
                return 'Técnico';
            default:
                return 'Operador';
        }
    };

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
                            onClick={() => OpenModal(ModalLista.modal_usuario)}
                            className="inline-flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                        >
                            <Plus className="w-3.5 h-3.5" />
                            <span>Nuevo Usuario</span>
                        </button>

                    </div>
                </div>
            </div>

            {/* Contenedor de la Tabla */}
            <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden w-full flex flex-col">
                <div className="overflow-x-auto max-h-[600px] overflow-y-auto w-full">
                    <table className="w-full min-w-full text-left text-sm text-gray-600">

                        {/* Cabecera Fija */}
                        <thead className="sticky top-0 bg-gray-50 border-b border-gray-200 text-xs uppercase text-gray-500 font-semibold z-10">
                        <tr>
                            <th className="px-4 py-3.5 w-16">ID</th>
                            <th className="px-4 py-3.5 w-44">Identificación</th>
                            <th className="px-4 py-3.5 min-w-[180px]">Nombres / Apellidos</th>
                            <th className="px-4 py-3.5 min-w-[180px]">Email</th>
                            <th className="px-4 py-3.5 w-32">Rol</th>
                            <th className="px-4 py-3.5 w-28">Estado</th>
                            <th className="px-4 py-3.5 w-36">Creación</th>
                            <th className="px-4 py-3.5 w-24 text-center">Acciones</th>
                        </tr>
                        </thead>

                        {/* Cuerpo de la Tabla */}
                        <tbody className="divide-y divide-gray-100">
                        {usuariosFiltrados.length > 0 ? (
                            usuariosFiltrados.map((usuario) => (
                                <tr key={usuario.usuario_id} className="hover:bg-gray-50/80 transition-colors">
                                    <td className="px-4 py-3.5 font-medium text-gray-900">
                                        #{usuario.usuario_id}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        <div className="font-medium text-gray-800">{usuario.identificacion}</div>
                                        <span className="inline-block px-2 py-0.5 text-[10px] font-semibold bg-gray-100 text-gray-600 rounded">
                        {usuario.tipo_identificacion}
                      </span>
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap font-medium text-gray-800">
                                        {usuario.nombres} {usuario.apellidos}
                                    </td>
                                    <td className="px-4 py-3.5 whitespace-nowrap text-gray-600">
                                        {usuario.email}
                                    </td>

                                    {/* Rol */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                      <span className="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium bg-slate-100 text-slate-700 rounded-md">
                        <Shield className="w-3 h-3 text-slate-500" />
                          {getNombreRol(usuario.rol_id)}
                      </span>
                                    </td>

                                    {/* Estado Booleano (Activo/Inactivo) */}
                                    <td className="px-4 py-3.5 whitespace-nowrap">
                                        {usuario.activo ? (
                                            <span className="inline-block px-2.5 py-0.5 text-xs font-semibold bg-green-50 text-green-700 border border-green-200 rounded-full">
                          Activo
                        </span>
                                        ) : (
                                            <span className="inline-block px-2.5 py-0.5 text-xs font-semibold bg-red-50 text-red-600 border border-red-200 rounded-full">
                          Inactivo
                        </span>
                                        )}
                                    </td>

                                    <td className="px-4 py-3.5 whitespace-nowrap text-xs text-gray-400">
                                        {usuario.fecha_creacion}
                                    </td>

                                    {/* Acciones */}
                                    <td className="px-4 py-3.5 whitespace-nowrap text-center">
                                        <div className="flex items-center justify-center gap-1.5">
                                            <button
                                                onClick={() => handleEditar(usuario)}
                                                className="p-1.5 text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors border border-blue-100"
                                                title="Editar usuario"
                                            >
                                                <Pencil className="w-4 h-4" />
                                            </button>
                                            <button
                                                onClick={() => handleEliminar(usuario.usuario_id)}
                                                className="p-1.5 text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors border border-red-100"
                                                title="Eliminar usuario"
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
                                    No se encontraron usuarios registrados.
                                </td>
                            </tr>
                        )}
                        </tbody>
                    </table>
                </div>

                {/* Pie de Tabla */}
                <div className="p-4 border-t border-gray-100 text-xs text-gray-500 flex justify-between items-center bg-gray-50/30">
                    <span>Total de registros: {usuariosFiltrados.length}</span>
                </div>
            </div>
        </div>
    );
}