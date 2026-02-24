import { useState } from "react";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import CarForm from "../components/CarForm";
import { createCar } from "../api/cars";
import type { CreateCarRequest } from "../types/car";

export default function CreateCarPage() {
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (data: CreateCarRequest) => {
        setIsLoading(true);
        try {
            const car = await createCar(data);
            toast.success(`${car.brand} ${car.model} created successfully!`);
            navigate("/cars");
        } catch {
            toast.error("Failed to create car. Please try again.");
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="form-page">
            <div className="page-header">
                <div>
                    <h2 className="page-title">New Car</h2>
                    <p className="page-subtitle">Add a new car to your collection</p>
                </div>
                <button className="btn btn--ghost" onClick={() => navigate("/cars")}>
                    ← Back to Cars
                </button>
            </div>

            <div className="form-card">
                <CarForm
                    onSubmit={handleSubmit}
                    submitLabel="Create Car"
                    isLoading={isLoading}
                />
            </div>
        </div>
    );
}
