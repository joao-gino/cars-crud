import axios from "axios";

const apiClient = axios.create({
    baseURL: "",
    headers: {
        "Content-Type": "application/json",
    },
});

let getTokenFn: (() => string | null) | null = null;
let onUnauthorizedFn: (() => void) | null = null;

export function setAuthInterceptors(
    getToken: () => string | null,
    onUnauthorized: () => void
) {
    getTokenFn = getToken;
    onUnauthorizedFn = onUnauthorized;
}

apiClient.interceptors.request.use((config) => {
    if (getTokenFn) {
        const token = getTokenFn();
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
    }
    return config;
});

apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401 && onUnauthorizedFn) {
            onUnauthorizedFn();
        }
        return Promise.reject(error);
    }
);

export default apiClient;
