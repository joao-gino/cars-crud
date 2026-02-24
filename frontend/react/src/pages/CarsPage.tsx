import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import { listCars, deleteCar } from "../api/cars";
import CarCard from "../components/CarCard";
import Pagination from "../components/Pagination";
import DeleteModal from "../components/DeleteModal";
import EmptyState from "../components/EmptyState";
import type { Car } from "../types/car";

const LIMIT = 9;

export default function CarsPage() {
    const [cars, setCars] = useState<Car[]>([]);
    const [total, setTotal] = useState(0);
    const [offset, setOffset] = useState(0);
    const [isLoading, setIsLoading] = useState(true);
    const [carToDelete, setCarToDelete] = useState<Car | null>(null);
    const [isDeleting, setIsDeleting] = useState(false);
    const navigate = useNavigate();

    const fetchCars = useCallback(async () => {
        setIsLoading(true);
        try {
            const response = await listCars(offset, LIMIT);
            setCars(response.data || []);
            setTotal(response.total);
        } catch {
            toast.error("Failed to load cars");
        } finally {
            setIsLoading(false);
        }
    }, [offset]);

    useEffect(() => {
        fetchCars();
    }, [fetchCars]);

    const handleDelete = async () => {
        if (!carToDelete) return;

        setIsDeleting(true);
        try {
            await deleteCar(carToDelete.id);
            toast.success(
                `${carToDelete.brand} ${carToDelete.model} deleted successfully`
            );
            setCarToDelete(null);
            fetchCars();
        } catch {
            toast.error("Failed to delete car");
        } finally {
            setIsDeleting(false);
        }
    };

    return (
        <div className="cars-page">
            <div className="page-header">
                <div>
                    <h2 className="page-title">Cars</h2>
                    <p className="page-subtitle">
                        {total} {total === 1 ? "car" : "cars"} in your collection
                    </p>
                </div>
                <button
                    className="btn btn--primary"
                    onClick={() => navigate("/cars/new")}
                >
                    <svg
                        className="btn-icon"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        width="18"
                        height="18"
                    >
                        <circle cx="12" cy="12" r="10" />
                        <line x1="12" y1="8" x2="12" y2="16" />
                        <line x1="8" y1="12" x2="16" y2="12" />
                    </svg>
                    Add Car
                </button>
            </div>

            {isLoading ? (
                <div className="cars-grid">
                    {Array.from({ length: 6 }).map((_, i) => (
                        <div key={i} className="car-card car-card--skeleton">
                            <div className="skeleton skeleton--strip" />
                            <div className="car-card__body">
                                <div className="skeleton skeleton--title" />
                                <div className="skeleton skeleton--text" />
                                <div className="skeleton skeleton--text skeleton--short" />
                            </div>
                        </div>
                    ))}
                </div>
            ) : cars.length === 0 ? (
                <EmptyState
                    title="No cars yet"
                    description="Add your first car to get started with your collection."
                    actionLabel="Add Your First Car"
                    onAction={() => navigate("/cars/new")}
                />
            ) : (
                <>
                    <div className="cars-grid">
                        {cars.map((car) => (
                            <CarCard key={car.id} car={car} onDelete={setCarToDelete} />
                        ))}
                    </div>

                    <Pagination
                        offset={offset}
                        limit={LIMIT}
                        total={total}
                        onPageChange={setOffset}
                    />
                </>
            )}

            {carToDelete && (
                <DeleteModal
                    car={carToDelete}
                    isLoading={isDeleting}
                    onConfirm={handleDelete}
                    onCancel={() => setCarToDelete(null)}
                />
            )}
        </div>
    );
}
