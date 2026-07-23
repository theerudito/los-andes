import api from "../helpers/fetching/axios.ts";
import type {reqLog} from "../modelos/logOk.ts";

export const logsService = {
    obtenerLogsError: async (req: reqLog) => {
        const { data } = await api.post('/logs-error/', req);
        return data;
    },
    obtenerLogsOk: async (req: reqLog) => {
        const { data } = await api.post('/logs-ok/', req);
        return data;
    },
};