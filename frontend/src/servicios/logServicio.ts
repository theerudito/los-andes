import api from "../helpers/fetching/axios.ts";
import type {LogError} from "../modelos/logError.ts";

export const logsService = {
    obtenerLogsError: async (req: LogError) => {
        const { data } = await api.post('/logs-error/', req || {});
        return data;
    },
    obtenerLogsOk: async (req: LogError) => {
        const { data } = await api.post('/logs-ok/', req || {});
        return data;
    },
};