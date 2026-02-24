import { useState, useEffect, useCallback } from "react";
import type { CreateCarRequest } from "../types/car";

interface CarFormProps {
    initialData?: CreateCarRequest;
    onSubmit: (data: CreateCarRequest) => Promise<void>;
    submitLabel: string;
    isLoading: boolean;
}

interface FormErrors {
    brand?: string;
    model?: string;
    year?: string;
    color?: string;
    price?: string;
}

const CURRENT_YEAR = new Date().getFullYear();

function validate(data: CreateCarRequest): FormErrors {
    const errors: FormErrors = {};

    if (!data.brand.trim()) {
        errors.brand = "Brand is required";
    } else if (data.brand.trim().length < 2) {
        errors.brand = "Brand must be at least 2 characters";
    } else if (data.brand.trim().length > 100) {
        errors.brand = "Brand must be at most 100 characters";
    }

    if (!data.model.trim()) {
        errors.model = "Model is required";
    } else if (data.model.trim().length > 100) {
        errors.model = "Model must be at most 100 characters";
    }

    if (!data.year) {
        errors.year = "Year is required";
    } else if (data.year < 1886) {
        errors.year = "Year must be 1886 or later";
    } else if (data.year > CURRENT_YEAR + 1) {
        errors.year = `Year cannot exceed ${CURRENT_YEAR + 1}`;
    }

    if (!data.color.trim()) {
        errors.color = "Color is required";
    } else if (data.color.trim().length > 50) {
        errors.color = "Color must be at most 50 characters";
    }

    if (data.price === undefined || data.price === null || isNaN(data.price)) {
        errors.price = "Price is required";
    } else if (data.price < 0) {
        errors.price = "Price must be a positive number";
    } else if (data.price > 99999999.99) {
        errors.price = "Price is too high";
    }

    return errors;
}

export default function CarForm({
    initialData,
    onSubmit,
    submitLabel,
    isLoading,
}: CarFormProps) {
    const [formData, setFormData] = useState<CreateCarRequest>(
        initialData || {
            brand: "",
            model: "",
            year: CURRENT_YEAR,
            color: "",
            price: 0,
        }
    );
    const [errors, setErrors] = useState<FormErrors>({});
    const [touched, setTouched] = useState<Record<string, boolean>>({});

    const validateForm = useCallback(() => {
        const newErrors = validate(formData);
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    }, [formData]);

    useEffect(() => {
        if (Object.keys(touched).length > 0) {
            validateForm();
        }
    }, [formData, touched, validateForm]);

    const handleChange = (field: keyof CreateCarRequest, value: string) => {
        let parsed: string | number = value;
        if (field === "year") {
            parsed = value === "" ? 0 : parseInt(value, 10);
        } else if (field === "price") {
            parsed = value === "" ? 0 : parseFloat(value);
        }
        setFormData((prev) => ({ ...prev, [field]: parsed }));
    };

    const handleBlur = (field: string) => {
        setTouched((prev) => ({ ...prev, [field]: true }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setTouched({ brand: true, model: true, year: true, color: true, price: true });

        if (!validateForm()) return;

        await onSubmit(formData);
    };

    const isValid = Object.keys(validate(formData)).length === 0;

    return (
        <form className="car-form" onSubmit={handleSubmit} noValidate>
            <div className="form-group">
                <label className="form-label" htmlFor="brand">
                    Brand
                </label>
                <input
                    id="brand"
                    type="text"
                    className={`form-input ${touched.brand && errors.brand ? "form-input--error" : ""}`}
                    placeholder="e.g. Toyota"
                    value={formData.brand}
                    onChange={(e) => handleChange("brand", e.target.value)}
                    onBlur={() => handleBlur("brand")}
                    disabled={isLoading}
                />
                {touched.brand && errors.brand && (
                    <span className="form-error">{errors.brand}</span>
                )}
            </div>

            <div className="form-group">
                <label className="form-label" htmlFor="model">
                    Model
                </label>
                <input
                    id="model"
                    type="text"
                    className={`form-input ${touched.model && errors.model ? "form-input--error" : ""}`}
                    placeholder="e.g. Corolla"
                    value={formData.model}
                    onChange={(e) => handleChange("model", e.target.value)}
                    onBlur={() => handleBlur("model")}
                    disabled={isLoading}
                />
                {touched.model && errors.model && (
                    <span className="form-error">{errors.model}</span>
                )}
            </div>

            <div className="form-row">
                <div className="form-group">
                    <label className="form-label" htmlFor="year">
                        Year
                    </label>
                    <input
                        id="year"
                        type="number"
                        className={`form-input ${touched.year && errors.year ? "form-input--error" : ""}`}
                        placeholder={`e.g. ${CURRENT_YEAR}`}
                        value={formData.year || ""}
                        onChange={(e) => handleChange("year", e.target.value)}
                        onBlur={() => handleBlur("year")}
                        disabled={isLoading}
                        min={1886}
                        max={CURRENT_YEAR + 1}
                    />
                    {touched.year && errors.year && (
                        <span className="form-error">{errors.year}</span>
                    )}
                </div>

                <div className="form-group">
                    <label className="form-label" htmlFor="color">
                        Color
                    </label>
                    <input
                        id="color"
                        type="text"
                        className={`form-input ${touched.color && errors.color ? "form-input--error" : ""}`}
                        placeholder="e.g. White"
                        value={formData.color}
                        onChange={(e) => handleChange("color", e.target.value)}
                        onBlur={() => handleBlur("color")}
                        disabled={isLoading}
                    />
                    {touched.color && errors.color && (
                        <span className="form-error">{errors.color}</span>
                    )}
                </div>
            </div>

            <div className="form-group">
                <label className="form-label" htmlFor="price">
                    Price (USD)
                </label>
                <input
                    id="price"
                    type="number"
                    className={`form-input ${touched.price && errors.price ? "form-input--error" : ""}`}
                    placeholder="e.g. 35000.00"
                    value={formData.price || ""}
                    onChange={(e) => handleChange("price", e.target.value)}
                    onBlur={() => handleBlur("price")}
                    disabled={isLoading}
                    min={0}
                    step="0.01"
                />
                {touched.price && errors.price && (
                    <span className="form-error">{errors.price}</span>
                )}
            </div>

            <button
                type="submit"
                className="btn btn--primary btn--lg btn--full"
                disabled={isLoading || !isValid}
            >
                {isLoading ? (
                    <span className="btn-loading">
                        <span className="spinner" /> Saving...
                    </span>
                ) : (
                    submitLabel
                )}
            </button>
        </form>
    );
}
