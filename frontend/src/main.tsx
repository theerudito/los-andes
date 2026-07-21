import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { ToasUI } from "./componentes/ToasUI.tsx";

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <ToasUI />
        <App />
    </StrictMode>,
)