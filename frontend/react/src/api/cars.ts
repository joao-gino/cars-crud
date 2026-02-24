import apiClient from "./client";
import type {
    Car,
    CreateCarRequest,
    UpdateCarRequest,
    SuccessResponse,
    PaginatedResponse,
} from "../types/car";

export async function listCars(
    offset: number = 0,
    limit: number = 10
): Promise<PaginatedResponse<Car>> {
    const response = await apiClient.get<PaginatedResponse<Car>>(
        `/api/v1/cars?offset=${offset}&limit=${limit}`
    );
    return response.data;
}

export async function getCar(id: string): Promise<Car> {
    const response = await apiClient.get<SuccessResponse<Car>>(
        `/api/v1/cars/${id}`
    );
    return response.data.data;
}

export async function createCar(data: CreateCarRequest): Promise<Car> {
    const response = await apiClient.post<SuccessResponse<Car>>(
        "/api/v1/cars",
        data
    );
    return response.data.data;
}

export async function updateCar(
    id: string,
    data: UpdateCarRequest
): Promise<Car> {
    const response = await apiClient.put<SuccessResponse<Car>>(
        `/api/v1/cars/${id}`,
        data
    );
    return response.data.data;
}

export async function deleteCar(id: string): Promise<void> {
    await apiClient.delete(`/api/v1/cars/${id}`);
}
