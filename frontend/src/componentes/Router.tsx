import React from "react";
import { createBrowserRouter, type RouteObject } from "react-router-dom";
import MainLayout from "./Layout.tsx";
import { navItems } from "./Sidebar.tsx";
import { Login } from "../paginas/login/Login.tsx";

// Importa los componentes de las subrutas que se abren desde los botones de acción
import PaginaHistorial from "../paginas/equipos/PaginaHistorial.tsx";
import PaginaEntregas from "../paginas/equipos/PaginaEntregas.tsx"; // Ajusta la ruta a tus archivos
import PaginaPagos from "../paginas/equipos/PaginaPagos.tsx";       // Ajusta la ruta a tus archivos

// 1. Rutas adicionales que NO están explícitamente en el menú lateral Sidebar
const extraRoutes: RouteObject[] = [
    {
        path: "historial", // /historial (recibe ?equipo_id=101 por query param)
        element: <PaginaHistorial />,
    },
    {
        path: "equipos/entrega", // /equipos/entrega?equipo_id=101
        element: <PaginaEntregas />,
    },
    {
        path: "equipos/pagos", // /equipos/pagos?equipo_id=101
        element: <PaginaPagos />,
    },
];

// 2. Función para construir dinámicamente las rutas desde el Sidebar
const buildRoutes = (): RouteObject[] => {
    const routes: RouteObject[] = [];

    navItems.forEach((item) => {
        // Procesa el item principal
        if (item.path && item.element) {
            if (item.path === "/") {
                routes.push({ index: true, element: item.element });
            } else {
                routes.push({
                    path: item.path.replace(/^\//, ""), // Quita el '/' inicial para 'children'
                    element: item.element,
                });
            }
        }

        // Procesa subitems (menús desplegables del Sidebar)
        if (item.subItems) {
            item.subItems.forEach((sub) => {
                if (sub.path && sub.element) {
                    routes.push({
                        path: sub.path.replace(/^\//, ""),
                        element: sub.element,
                    });
                }
            });
        }
    });

    // Unifica las rutas del menú con las rutas adicionales operativas
    return [...routes, ...extraRoutes];
};

export const router = createBrowserRouter([
    {
        path: "/login",
        element: <Login />,
    },
    {
        path: "/",
        element: <MainLayout />,
        children: buildRoutes(),
    },
    {
        // Captura cualquier ruta inexistente
        path: "*",
        element: (
            <div className="flex h-screen items-center justify-center bg-gray-100">
                <div className="text-center">
                    <h1 className="text-4xl font-bold text-gray-800">404</h1>
                    <p className="text-gray-500 mt-2">Página no encontrada</p>
                </div>
            </div>
        ),
    },
]);