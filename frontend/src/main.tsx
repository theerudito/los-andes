import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {ToasUI} from "./componentes/ToasUI.tsx";

createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <ToasUI/>
      <p className="bg-blue-600 text-white font-bold p-4 rounded-lg shadow-lg text-center">
          ¡Si ves este recuadro azul centrado y con texto blanco, Tailwind CSS está funcionando correctamente!
      </p>
      <App/>
  </StrictMode>,
)
