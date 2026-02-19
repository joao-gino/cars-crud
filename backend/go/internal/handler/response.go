package handler

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type PaginatedResponse struct {
	Data   interface{} `json:"data"`
	Total  int64       `json:"total"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
