import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

export default function Layout() {
    const { logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate("/login");
    };

    return (
        <div className="app-layout">
            <aside className="sidebar">
                <div className="sidebar-brand">
                    <span className="brand-icon">🚗</span>
                    <h1 className="brand-title">Cars CRUD</h1>
                </div>

                <nav className="sidebar-nav">
                    <NavLink
                        to="/cars"
                        end
                        className={({ isActive }) =>
                            `nav-link ${isActive ? "nav-link--active" : ""}`
                        }
                    >
                        <svg
                            className="nav-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            strokeWidth="2"
                        >
                            <rect x="3" y="3" width="7" height="7" rx="1" />
                            <rect x="14" y="3" width="7" height="7" rx="1" />
                            <rect x="3" y="14" width="7" height="7" rx="1" />
                            <rect x="14" y="14" width="7" height="7" rx="1" />
                        </svg>
                        All Cars
                    </NavLink>
                    <NavLink
                        to="/cars/new"
                        className={({ isActive }) =>
                            `nav-link ${isActive ? "nav-link--active" : ""}`
                        }
                    >
                        <svg
                            className="nav-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            strokeWidth="2"
                        >
                            <circle cx="12" cy="12" r="10" />
                            <line x1="12" y1="8" x2="12" y2="16" />
                            <line x1="8" y1="12" x2="16" y2="12" />
                        </svg>
                        New Car
                    </NavLink>
                </nav>

                <button className="logout-btn" onClick={handleLogout}>
                    <svg
                        className="nav-icon"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                    >
                        <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                        <polyline points="16 17 21 12 16 7" />
                        <line x1="21" y1="12" x2="9" y2="12" />
                    </svg>
                    Logout
                </button>
            </aside>

            <main className="main-content">
                <Outlet />
            </main>
        </div>
    );
}
