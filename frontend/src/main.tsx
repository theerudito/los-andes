import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {ToasUI} from "./componentes/ToasUI.tsx";
import ModalCliente from "./modales/ModalClientes.tsx";
import ModalEquipos from "./modales/ModalEquipos.tsx";
import ModalMarcas from "./modales/ModalMarcas.tsx";
import ModalUsuario from "./modales/ModalUsuario.tsx";
import ModalEntregas from "./modales/ModalEntregas.tsx";
import ModalPagos from "./modales/ModalPagos.tsx";
import ModalHistorial from "./modales/ModalHistorial.tsx";

createRoot(document.getElementById('root')!).render(
    <StrictMode>
            <ModalMarcas/>
            <ModalCliente/>
            <ModalEquipos/>
            <ModalUsuario/>
            <ModalEntregas/>
            <ModalPagos/>
            <ModalHistorial/>
            <ModalUsuario/>
            <ToasUI/>
            <App/>
    </StrictMode>,
)