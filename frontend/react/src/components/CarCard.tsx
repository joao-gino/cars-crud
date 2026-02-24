import { useNavigate } from "react-router-dom";
import type { Car } from "../types/car";

interface CarCardProps {
    car: Car;
    onDelete: (car: Car) => void;
}

const COLOR_MAP: Record<string, string> = {
    white: "#f0f0f0",
    black: "#1a1a2e",
    red: "#e74c3c",
    blue: "#3498db",
    green: "#2ecc71",
    silver: "#bdc3c7",
    gray: "#95a5a6",
    grey: "#95a5a6",
    yellow: "#f1c40f",
    orange: "#e67e22",
    brown: "#8b4513",
    purple: "#9b59b6",
    pink: "#e91e90",
    gold: "#ffd700",
    beige: "#f5f5dc",
    navy: "#001f3f",
    maroon: "#800000",
};

function getColorHex(color: string): string {
    return COLOR_MAP[color.toLowerCase()] || "#6c5ce7";
}

function formatPrice(price: number): string {
    return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: "USD",
    }).format(price);
}

export default function CarCard({ car, onDelete }: CarCardProps) {
    const navigate = useNavigate();

    return (
        <div className="car-card" onClick={() => navigate(`/cars/${car.id}/edit`)}>
            <div
                className="car-card__color-strip"
                style={{ background: getColorHex(car.color) }}
            />
            <div className="car-card__body">
                <div className="car-card__header">
                    <h3 className="car-card__title">
                        {car.brand} {car.model}
                    </h3>
                    <span className="car-card__year">{car.year}</span>
                </div>

                <div className="car-card__details">
                    <div className="car-card__detail">
                        <span className="car-card__label">Color</span>
                        <div className="car-card__color-info">
                            <span
                                className="car-card__color-dot"
                                style={{ background: getColorHex(car.color) }}
                            />
                            <span className="car-card__value">
                                {car.color.charAt(0).toUpperCase() + car.color.slice(1)}
                            </span>
                        </div>
                    </div>
                    <div className="car-card__detail">
                        <span className="car-card__label">Price</span>
                        <span className="car-card__price">{formatPrice(car.price)}</span>
                    </div>
                </div>

                <div className="car-card__actions">
                    <button
                        className="btn btn--ghost btn--sm"
                        onClick={(e) => {
                            e.stopPropagation();
                            navigate(`/cars/${car.id}/edit`);
                        }}
                    >
                        Edit
                    </button>
                    <button
                        className="btn btn--danger btn--sm"
                        onClick={(e) => {
                            e.stopPropagation();
                            onDelete(car);
                        }}
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    );
}
