import axios from 'axios';
import { url_base } from "./url.ts";

const api = axios.create({
    baseURL: url_base,
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token && token !== 'null' && token !== 'undefined' && config.headers) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

api.interceptors.response.use(
    (response) => {
        if (response.data && response.data.error) {
            if (response.data.error === 'Token inválido' || response.data.error === 'Token expirado') {
                console.error('Sesión vencida o token inválido');
                localStorage.removeItem('token');
            }
            return Promise.reject(new Error(response.data.error));
        }
        return response;
    },
    (error) => {
        const mensajeError =
            error.response?.data?.message ||
            error.response?.data?.error ||
            error.message ||
            'Error en la petición';

        error.message = mensajeError;

        if (error.response?.status === 401 || error.response?.status === 403) {
            console.error('Sesión expirada o no autorizada');
            localStorage.removeItem('token');
        }

        return Promise.reject(error);
    }
);

export default api;