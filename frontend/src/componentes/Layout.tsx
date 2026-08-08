import { Outlet, Navigate } from 'react-router-dom';
import Sidebar from './Sidebar';
import Footer from './Footer';
import React from "react";
import { useUsuarios } from "../store/useUsuarios.ts";

export default function Layout(): React.ReactElement {
    const { isLogin } = useUsuarios();
    const token = localStorage.getItem('token');

    if (!isLogin || !token) {
        return <Navigate to="/login" replace />;
    }

    return (
        <div className="h-screen w-screen flex overflow-hidden bg-gray-50">
            <Sidebar />
            <div className="flex-1 flex flex-col md:ml-64 h-full min-w-0">
                <main className="flex-1 overflow-y-auto p-6 md:p-8 pt-16 md:pt-8">
                    <Outlet />
                </main>
                <Footer />
            </div>
        </div>
    );
}