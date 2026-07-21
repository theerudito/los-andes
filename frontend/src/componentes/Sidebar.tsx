import { NavLink } from 'react-router-dom';
import {
    Menu,
    X,
    Home,
    Users,
    Laptop,
    Settings,
    LayoutDashboard,
    LogOut,
    ChevronDown, Bolt, Cpu, Timeline, LaptopMinimalCheck
} from 'lucide-react';

import PaginaIndex from '../paginas/PaginaIndex';
import PaginaClientes from "../paginas/clientes/PaginaClientes.tsx";
import PaginaEquipos from "../paginas/equipos/PaginaEquipos.tsx";
import {useState} from "react";
import * as React from "react";
import PaginaHistorial from "../paginas/equipos/PaginaHistorial.tsx";
import PaginaAuditoriaError from "../paginas/auditoria/PaginaAuditoriaError.tsx";
import PaginaAuditoriaOk from "../paginas/auditoria/PaginaAuditoriaOk.tsx";
import PaginaUsuarios from "../paginas/configuracion/PaginaUsuarios.tsx";

export interface SubNavItem {
    label: string;
    path: string;
    element?: React.ReactNode;
}

export interface NavItem {
    label: string;
    path?: string;
    icon: React.ElementType;
    element?: React.ReactNode;
    subItems?: SubNavItem[];
}

export const navItems: NavItem[] = [
    {
        label: 'Inicio',
        path: '/',
        icon: Home,
        element: <PaginaIndex />
    },
    {
        label: 'Clientes',
        icon: Users,
        subItems: [
            {
                label: 'Listado Clientes',
                path: '/clientes/listado',
                element: <PaginaClientes />
            }
        ]
    },
    {
        label: 'Equipos',
        icon: Laptop,
        subItems: [
            {
                label: 'Listado Equipos',
                path: '/equipos/listado',
                element: <PaginaEquipos />
            }
        ]
    },
    {
        label: 'Auditoria',
        icon: Cpu,
        subItems: [
            {
                label: 'Logs Error',
                path: '/auditoria/logs-errors',
                element: <PaginaAuditoriaError />
            },
            {
                label: 'Logs Ok',
                path: '/auditoria/logs-ok',
                element: <PaginaAuditoriaOk />
            },
        ]
    },
    {
        label: 'Configuración',
        icon: Bolt,
        subItems: [
            {
                label: 'Listado Usuarios',
                path: '/usuarios/listado',
                element: <PaginaUsuarios />
            }
        ]
    },
];

export default function Sidebar(): React.ReactElement {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [openSubmenu, setOpenSubmenu] = useState<string | null>(null);

    const toggleSidebar = (): void => setIsOpen(!isOpen);
    const closeSidebar = (): void => setIsOpen(false);

    const toggleSubmenu = (label: string): void => {
        setOpenSubmenu(openSubmenu === label ? null : label);
    };

    const handleLogout = (): void => {
        console.log("Cerrando sesión...");
    };

    const linkStyle = ({ isActive }: { isActive: boolean }): string =>
        `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors text-sm font-medium ${
            isActive
                ? 'bg-blue-600 text-white'
                : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'
        }`;

    const subLinkStyle = ({ isActive }: { isActive: boolean }): string =>
        `block px-4 py-2 rounded-lg text-sm transition-colors font-medium ${
            isActive
                ? 'bg-blue-50 text-blue-600 font-semibold'
                : 'text-gray-500 hover:bg-gray-100 hover:text-gray-900'
        }`;

    return (
        <>
            <div className="md:hidden fixed top-4 left-4 z-50">
                <button
                    onClick={toggleSidebar}
                    className="p-2 rounded-lg bg-white shadow-md text-gray-700 hover:bg-gray-50 focus:outline-none"
                    aria-label="Toggle menu"
                >
                    {isOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
                </button>
            </div>

            {isOpen && (
                <div
                    onClick={closeSidebar}
                    className="fixed inset-0 bg-black/50 z-40 md:hidden transition-opacity"
                />
            )}

            <aside
                className={`fixed top-0 left-0 z-40 h-screen w-64 bg-white border-r border-gray-200 flex flex-col justify-between transition-transform duration-300 ease-in-out md:translate-x-0 ${
                    isOpen ? 'translate-x-0' : '-translate-x-full'
                }`}
            >
                <div className="p-5 flex-1 flex flex-col min-h-0">
                    <div className="flex items-center gap-2 mb-6 px-2 shrink-0">
                        <LayoutDashboard className="w-7 h-7 text-blue-600" />
                        <span className="font-bold text-xl text-gray-800">Panel Admin</span>
                    </div>

                    <nav className="space-y-1 overflow-y-auto flex-1 pr-1">
                        {navItems.map((item) => {
                            const IconComponent = item.icon;
                            const hasSubmenu = Boolean(item.subItems && item.subItems.length > 0);
                            const isSubmenuOpen = openSubmenu === item.label;

                            if (hasSubmenu) {
                                return (
                                    <div key={item.label} className="space-y-1">
                                        <button
                                            onClick={() => toggleSubmenu(item.label)}
                                            className="w-full flex items-center justify-between px-4 py-3 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-100 hover:text-gray-900 transition-colors"
                                        >
                                            <div className="flex items-center gap-3">
                                                <IconComponent className="w-5 h-5 shrink-0" />
                                                <span>{item.label}</span>
                                            </div>
                                            <ChevronDown
                                                className={`w-4 h-4 transition-transform duration-200 ${
                                                    isSubmenuOpen ? 'rotate-180' : ''
                                                }`}
                                            />
                                        </button>

                                        {isSubmenuOpen && (
                                            <div className="pl-9 space-y-1">
                                                {item.subItems?.map((subItem) => (
                                                    <NavLink
                                                        key={subItem.path}
                                                        to={subItem.path}
                                                        className={subLinkStyle}
                                                        onClick={closeSidebar}
                                                    >
                                                        {subItem.label}
                                                    </NavLink>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                );
                            }

                            return (
                                <NavLink
                                    key={item.path}
                                    to={item.path!}
                                    className={linkStyle}
                                    onClick={closeSidebar}
                                >
                                    <IconComponent className="w-5 h-5 shrink-0" />
                                    <span>{item.label}</span>
                                </NavLink>
                            );
                        })}
                    </nav>
                </div>

                <div className="p-4 border-t border-gray-100 shrink-0">
                    <div className="flex items-center justify-between gap-3">
                        <div className="flex items-center gap-3 min-w-0">
                            <div className="w-9 h-9 rounded-full bg-blue-100 text-blue-600 font-bold flex items-center justify-center shrink-0">
                                U
                            </div>
                            <div className="truncate">
                                <p className="text-sm font-medium text-gray-700 truncate">Jorge Loor</p>
                            </div>
                        </div>

                        <button
                            onClick={handleLogout}
                            className="p-2 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors shrink-0"
                            title="Cerrar sesión"
                            aria-label="Cerrar sesión"
                        >
                            <LogOut className="w-5 h-5" />
                        </button>
                    </div>
                </div>
            </aside>
        </>
    );
}