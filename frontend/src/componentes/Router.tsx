import { createBrowserRouter } from "react-router-dom";
import MainLayout from "./Layout.tsx";
import PaginaIndex from "../paginas/PaginaIndex.tsx";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <MainLayout />,
        children: [
            { index: true, element: <PaginaIndex /> },
        ],
    },
]);