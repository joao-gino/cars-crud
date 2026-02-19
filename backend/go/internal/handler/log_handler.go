package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/gino/cars-crud/internal/repository"
)

type LogHandler struct {
	repo repository.LogRepository
}

func NewLogHandler(repo repository.LogRepository) *LogHandler {
	return &LogHandler{repo: repo}
}

func (h *LogHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/logs", func(r chi.Router) {
		r.Get("/", h.GetAll)
	})
}

// GetAll godoc
// @Summary      List request logs
// @Description  Get a paginated list of all request logs from MongoDB
// @Tags         logs
// @Produce      json
// @Security     BearerAuth
// @Param        offset  query     int  false  "Offset"  default(0)
// @Param        limit   query     int  false  "Limit"   default(20)
// @Success      200     {object}  PaginatedResponse{data=[]domain.RequestLog}
// @Failure      401     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /api/v1/logs [get]
func (h *LogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	logs, total, err := h.repo.GetAll(r.Context(), offset, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list logs")
		return
	}

	respondJSON(w, http.StatusOK, PaginatedResponse{
		Data:   logs,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	})
}
