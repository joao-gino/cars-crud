import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "react-hot-toast";
import { AuthProvider } from "./hooks/useAuth";
import ProtectedRoute from "./components/ProtectedRoute";
import Layout from "./components/Layout";
import LoginPage from "./pages/LoginPage";
import CarsPage from "./pages/CarsPage";
import CreateCarPage from "./pages/CreateCarPage";
import EditCarPage from "./pages/EditCarPage";
import "./App.css";

export default function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <Routes>
                    <Route path="/login" element={<LoginPage />} />
                    <Route
                        element={
                            <ProtectedRoute>
                                <Layout />
                            </ProtectedRoute>
                        }
                    >
                        <Route path="/cars" element={<CarsPage />} />
                        <Route path="/cars/new" element={<CreateCarPage />} />
                        <Route path="/cars/:id/edit" element={<EditCarPage />} />
                    </Route>
                    <Route path="*" element={<Navigate to="/cars" replace />} />
                </Routes>
            </BrowserRouter>
            <Toaster
                position="top-right"
                toastOptions={{
                    duration: 4000,
                    style: {
                        background: "var(--color-surface)",
                        color: "var(--color-text)",
                        border: "1px solid var(--color-border)",
                        borderRadius: "12px",
                        fontSize: "0.875rem",
                    },
                    success: {
                        iconTheme: {
                            primary: "var(--color-success)",
                            secondary: "var(--color-surface)",
                        },
                    },
                    error: {
                        iconTheme: {
                            primary: "var(--color-danger)",
                            secondary: "var(--color-surface)",
                        },
                    },
                }}
            />
        </AuthProvider>
    );
}
