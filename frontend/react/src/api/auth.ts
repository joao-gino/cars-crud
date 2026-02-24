import apiClient from "./client";
import type { SuccessResponse, ValidateResponse } from "../types/car";

export async function validateApiKey(apiKey: string): Promise<string> {
    const response = await apiClient.post<SuccessResponse<ValidateResponse>>(
        "/auth/validate",
        { api_key: apiKey }
    );
    return response.data.data.token;
}
