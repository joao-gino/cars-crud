import { useState } from "react";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import { useAuth } from "../hooks/useAuth";
import { validateApiKey } from "../api/auth";

export default function LoginPage() {
    const [apiKey, setApiKey] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!apiKey.trim()) {
            toast.error("Please enter your API key");
            return;
        }

        setIsLoading(true);
        try {
            const token = await validateApiKey(apiKey.trim());
            login(token);
            toast.success("Authenticated successfully!");
            navigate("/cars");
        } catch {
            toast.error("Invalid API key. Please try again.");
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="login-page">
            <div className="login-card">
                <div className="login-card__header">
                    <span className="login-card__icon">🚗</span>
                    <h1 className="login-card__title">Cars CRUD</h1>
                    <p className="login-card__subtitle">
                        Enter your API key to access the dashboard
                    </p>
                </div>

                <form onSubmit={handleSubmit} className="login-form">
                    <div className="form-group">
                        <label className="form-label" htmlFor="api-key">
                            API Key
                        </label>
                        <input
                            id="api-key"
                            type="password"
                            className="form-input"
                            placeholder="Enter your API key..."
                            value={apiKey}
                            onChange={(e) => setApiKey(e.target.value)}
                            disabled={isLoading}
                            autoFocus
                        />
                    </div>

                    <button
                        type="submit"
                        className="btn btn--primary btn--lg btn--full"
                        disabled={isLoading || !apiKey.trim()}
                    >
                        {isLoading ? (
                            <span className="btn-loading">
                                <span className="spinner" /> Authenticating...
                            </span>
                        ) : (
                            "Sign In"
                        )}
                    </button>
                </form>

                <p className="login-card__hint">
                    The API key is configured in the backend <code>.env</code> file
                </p>
            </div>
        </div>
    );
}
