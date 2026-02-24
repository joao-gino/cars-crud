import type { Car } from "../types/car";

interface DeleteModalProps {
    car: Car;
    isLoading: boolean;
    onConfirm: () => void;
    onCancel: () => void;
}

export default function DeleteModal({
    car,
    isLoading,
    onConfirm,
    onCancel,
}: DeleteModalProps) {
    return (
        <div className="modal-overlay" onClick={onCancel}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
                <div className="modal__icon">
                    <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        width="48"
                        height="48"
                    >
                        <polyline points="3 6 5 6 21 6" />
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                        <line x1="10" y1="11" x2="10" y2="17" />
                        <line x1="14" y1="11" x2="14" y2="17" />
                    </svg>
                </div>
                <h2 className="modal__title">Delete Car</h2>
                <p className="modal__text">
                    Are you sure you want to delete{" "}
                    <strong>
                        {car.brand} {car.model} ({car.year})
                    </strong>
                    ? This action cannot be undone.
                </p>
                <div className="modal__actions">
                    <button
                        className="btn btn--ghost"
                        onClick={onCancel}
                        disabled={isLoading}
                    >
                        Cancel
                    </button>
                    <button
                        className="btn btn--danger"
                        onClick={onConfirm}
                        disabled={isLoading}
                    >
                        {isLoading ? (
                            <span className="btn-loading">
                                <span className="spinner" /> Deleting...
                            </span>
                        ) : (
                            "Delete"
                        )}
                    </button>
                </div>
            </div>
        </div>
    );
}
