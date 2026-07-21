import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Footer from './Footer';
import * as React from "react";

export default function Layout(): React.ReactElement {
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