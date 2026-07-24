import React from "react";
import { createBrowserRouter, type RouteObject } from "react-router-dom";
import MainLayout from "./Layout.tsx";
import { navItems } from "./Sidebar.tsx";
import { Login } from "../paginas/login/Login.tsx";
import PaginaHistorial from "../paginas/equipos/PaginaHistorial.tsx";
import PaginaEntregas from "../paginas/equipos/PaginaEntregas.tsx";
import PaginaPagos from "../paginas/equipos/PaginaPagos.tsx";

const extraRoutes: RouteObject[] = [
    {
        path: "equipos/historial",
        element: <PaginaHistorial />,
    },
    {
        path: "equipos/entrega",
        element: <PaginaEntregas />,
    },
    {
        path: "equipos/pagos",
        element: <PaginaPagos />,
    },
];

const buildRoutes = (): RouteObject[] => {
    const routes: RouteObject[] = [];

    navItems.forEach((item) => {
        if (item.path && item.element) {
            if (item.path === "/") {
                routes.push({ index: true, element: item.element });
            } else {
                routes.push({
                    path: item.path.replace(/^\//, ""),
                    element: item.element,
                });
            }
        }

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