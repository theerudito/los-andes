import * as React from "react";

export default function Footer(): React.ReactElement {
    return (
        <footer className="bg-white border-t border-gray-200 py-4 px-6 text-center text-sm text-gray-500 shrink-0">
            <p>Hecho por Between Bytes Software. © {new Date().getFullYear()}</p>
        </footer>
    );
}