import { jwtDecode } from "jwt-decode";

interface TokenPayload {
    name: string;
    user_id: number;
}

export const ObtenerToken = (): TokenPayload | null => {
    try {
        const token = localStorage.getItem("token");

        if (!token || token === "null" || token === "undefined") {
            return null;
        }

        const decoded = jwtDecode<TokenPayload>(token);

        return {
            name: decoded.name,
            user_id: decoded.user_id,
        };
    } catch {
        return null;
    }
};