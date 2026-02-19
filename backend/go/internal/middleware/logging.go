package middleware

import (
	"net/http"
	"time"

	"github.com/gino/cars-crud/internal/domain"
	"github.com/gino/cars-crud/internal/queue"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func RequestLogger(producer *queue.LogProducer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rec, r)

			logEntry := domain.RequestLog{
				Method:     r.Method,
				Path:       r.URL.Path,
				StatusCode: rec.statusCode,
				Duration:   time.Since(start).Milliseconds(),
				IP:         r.RemoteAddr,
				UserAgent:  r.UserAgent(),
				Timestamp:  start,
			}

			_ = producer.Publish(r.Context(), logEntry)
		})
	}
}
