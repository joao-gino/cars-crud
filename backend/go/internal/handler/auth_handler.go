package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gino/cars-crud/internal/domain"
)

type AuthHandler struct {
	apiKey    string
	jwtSecret string
}

func NewAuthHandler(apiKey, jwtSecret string) *AuthHandler {
	return &AuthHandler{apiKey: apiKey, jwtSecret: jwtSecret}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/validate", h.Validate)
	})
}

// Validate godoc
// @Summary      Validate API Key
// @Description  Validates an API key and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      domain.ValidateRequest  true  "API Key"
// @Success      200   {object}  SuccessResponse{data=domain.ValidateResponse}
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /auth/validate [post]
func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var req domain.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.APIKey == "" {
		respondError(w, http.StatusBadRequest, "api_key is required")
		return
	}

	if req.APIKey != h.apiKey {
		respondError(w, http.StatusUnauthorized, "invalid api key")
		return
	}

	claims := jwt.MapClaims{
		"iss": "cars-crud-api",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Data: domain.ValidateResponse{Token: signed},
	})
}
