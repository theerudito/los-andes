import * as React from "react";

export default function Footer(): React.ReactElement {
    return (
        <footer className="bg-white border-t border-gray-200 py-4 px-6 text-center text-sm text-gray-500 shrink-0">
            <p>© {new Date().getFullYear()} Mi Proyecto. Todos los derechos reservados.</p>
        </footer>
    );
}