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
        //const token = localStorage.getItem('token');
        const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJuYW1lIjoiU0lTVEVNQSIsInJvbCI6IlNJU1RFTUEiLCJhdWQiOlsiIl0sImV4cCI6MTc4NDgwNjQyNSwiaWF0IjoxNzg0NzcwNDI1fQ.ysaSelX-xW1nJ8qjD03tWFton788L85FbV-GQy9yoAw';
        if (token && config.headers) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

api.interceptors.response.use(
    (response) => {
        if (response.data && response.data.error) {
            if (response.data.error === 'Token inválido') {
                console.error('Sesión vencida o token inválido');
                localStorage.removeItem('token');
            }
            return Promise.reject(new Error(response.data.error));
        }
        return response;
    },
    (error) => {
        const mensajeError = error.response?.data?.error || 'Error en la petición';

        if (error.response?.status === 401 || error.response?.status === 403) {
            console.error('Sesión expirada o no autorizada');
            localStorage.removeItem('token');
        }

        return Promise.reject(new Error(mensajeError));
    }
);

export default api;