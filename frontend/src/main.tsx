import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import 'react-datepicker/dist/react-datepicker.css';
import App from './App.tsx'
import {ToasUI} from "./componentes/ToasUI.tsx";
import ModalCliente from "./modales/ModalClientes.tsx";
import ModalEquipos from "./modales/ModalEquipos.tsx";
import ModalUsuario from "./modales/ModalUsuario.tsx";
import ModalEntregas from "./modales/ModalEntregas.tsx";
import ModalPagos from "./modales/ModalPagos.tsx";
import ModalHistorial from "./modales/ModalHistorial.tsx";
import { es } from "date-fns/locale";
import {registerLocale} from "react-datepicker";
registerLocale("es", es);

createRoot(document.getElementById('root')!).render(
    <StrictMode>
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