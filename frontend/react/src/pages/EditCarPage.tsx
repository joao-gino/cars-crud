import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import toast from "react-hot-toast";
import CarForm from "../components/CarForm";
import { getCar, updateCar } from "../api/cars";
import type { CreateCarRequest } from "../types/car";

export default function EditCarPage() {
    const { id } = useParams<{ id: string }>();
    const [initialData, setInitialData] = useState<CreateCarRequest | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [isFetching, setIsFetching] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        if (!id) return;

        const fetchCar = async () => {
            try {
                const car = await getCar(id);
                setInitialData({
                    brand: car.brand,
                    model: car.model,
                    year: car.year,
                    color: car.color,
                    price: car.price,
                });
            } catch {
                toast.error("Car not found");
                navigate("/cars");
            } finally {
                setIsFetching(false);
            }
        };

        fetchCar();
    }, [id, navigate]);

    const handleSubmit = async (data: CreateCarRequest) => {
        if (!id) return;

        setIsLoading(true);
        try {
            const car = await updateCar(id, data);
            toast.success(`${car.brand} ${car.model} updated successfully!`);
            navigate("/cars");
        } catch {
            toast.error("Failed to update car. Please try again.");
        } finally {
            setIsLoading(false);
        }
    };

    if (isFetching) {
        return (
            <div className="form-page">
                <div className="page-header">
                    <div>
                        <div className="skeleton skeleton--title" />
                        <div className="skeleton skeleton--text skeleton--short" />
                    </div>
                </div>
                <div className="form-card">
                    <div className="skeleton skeleton--input" />
                    <div className="skeleton skeleton--input" />
                    <div className="skeleton skeleton--input" />
                </div>
            </div>
        );
    }

    if (!initialData) return null;

    return (
        <div className="form-page">
            <div className="page-header">
                <div>
                    <h2 className="page-title">Edit Car</h2>
                    <p className="page-subtitle">
                        Update the details of {initialData.brand} {initialData.model}
                    </p>
                </div>
                <button className="btn btn--ghost" onClick={() => navigate("/cars")}>
                    ← Back to Cars
                </button>
            </div>

            <div className="form-card">
                <CarForm
                    initialData={initialData}
                    onSubmit={handleSubmit}
                    submitLabel="Update Car"
                    isLoading={isLoading}
                />
            </div>
        </div>
    );
}
