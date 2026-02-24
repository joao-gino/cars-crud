import {
    createContext,
    useContext,
    useState,
    useEffect,
    useCallback,
    type ReactNode,
} from "react";
import { setAuthInterceptors } from "../api/client";

interface AuthContextType {
    token: string | null;
    isAuthenticated: boolean;
    login: (token: string) => void;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [token, setToken] = useState<string | null>(null);

    const logout = useCallback(() => {
        setToken(null);
    }, []);

    const login = useCallback((newToken: string) => {
        setToken(newToken);
    }, []);

    useEffect(() => {
        setAuthInterceptors(() => token, logout);
    }, [token, logout]);

    return (
        <AuthContext.Provider
            value={{
                token,
                isAuthenticated: !!token,
                login,
                logout,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth(): AuthContextType {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
