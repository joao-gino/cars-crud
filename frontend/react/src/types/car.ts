export interface Car {
    id: string;
    brand: string;
    model: string;
    year: number;
    color: string;
    price: number;
    created_at: string;
    updated_at: string;
}

export interface CreateCarRequest {
    brand: string;
    model: string;
    year: number;
    color: string;
    price: number;
}

export interface UpdateCarRequest {
    brand?: string;
    model?: string;
    year?: number;
    color?: string;
    price?: number;
}

export interface SuccessResponse<T> {
    data: T;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    offset: number;
    limit: number;
}

export interface ErrorResponse {
    error: string;
}

export interface ValidateRequest {
    api_key: string;
}

export interface ValidateResponse {
    token: string;
}
