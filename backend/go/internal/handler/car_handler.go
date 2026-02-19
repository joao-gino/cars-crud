package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gino/cars-crud/internal/domain"
	"github.com/gino/cars-crud/internal/usecase"
)

type CarHandler struct {
	usecase *usecase.CarUsecase
}

func NewCarHandler(uc *usecase.CarUsecase) *CarHandler {
	return &CarHandler{usecase: uc}
}

func (h *CarHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/cars", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.GetAll)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// Create godoc
// @Summary      Create a car
// @Description  Create a new car entry
// @Tags         cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        car  body      domain.CreateCarRequest  true  "Car data"
// @Success      201  {object}  SuccessResponse{data=domain.Car}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/cars [post]
func (h *CarHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Brand == "" || req.Model == "" || req.Year == 0 {
		respondError(w, http.StatusBadRequest, "brand, model, and year are required")
		return
	}

	car, err := h.usecase.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create car")
		return
	}

	respondJSON(w, http.StatusCreated, SuccessResponse{Data: car})
}

// GetAll godoc
// @Summary      List all cars
// @Description  Get a paginated list of all cars
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        offset  query     int  false  "Offset"  default(0)
// @Param        limit   query     int  false  "Limit"   default(10)
// @Success      200     {object}  PaginatedResponse{data=[]domain.Car}
// @Failure      500     {object}  ErrorResponse
// @Router       /api/v1/cars [get]
func (h *CarHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	cars, total, err := h.usecase.GetAll(r.Context(), offset, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list cars")
		return
	}

	respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:   cars,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	})
}

// GetByID godoc
// @Summary      Get a car
// @Description  Get a single car by its ID
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car ID (UUID)"
// @Success      200  {object}  SuccessResponse{data=domain.Car}
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/v1/cars/{id} [get]
func (h *CarHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid car id")
		return
	}

	car, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			respondError(w, http.StatusNotFound, "car not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get car")
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Data: car})
}

// Update godoc
// @Summary      Update a car
// @Description  Update an existing car by ID
// @Tags         cars
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "Car ID (UUID)"
// @Param        car  body      domain.UpdateCarRequest  true  "Car data to update"
// @Success      200  {object}  SuccessResponse{data=domain.Car}
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/cars/{id} [put]
func (h *CarHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid car id")
		return
	}

	var req domain.UpdateCarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	car, err := h.usecase.Update(r.Context(), id, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			respondError(w, http.StatusNotFound, "car not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update car")
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Data: car})
}

// Delete godoc
// @Summary      Delete a car
// @Description  Soft-delete a car by ID
// @Tags         cars
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Car ID (UUID)"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/cars/{id} [delete]
func (h *CarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid car id")
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete car")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
