package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/httpresponse"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// Liveness indicates if the application server process is running.
// (Used by Kubernetes/Docker to restart dead containers)
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	httpresponse.Success(
		w,
		http.StatusOK,
		"Service is alive",
		map[string]string{"status": "up"},
	)
}

// Readiness indicates if the application is healthy enough to handle live traffic.
// Checks external dependencies like database connectivity with a strict timeout.
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	// Enforce 2-second timeout so health checks fail fast rather than hanging connections
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		httpresponse.Error(
			w,
			http.StatusServiceUnavailable,
			"DATABASE_UNAVAILABLE",
			"Database is unreachable",
		)
		return
	}

	httpresponse.Success(
		w,
		http.StatusOK,
		"Service is ready",
		map[string]string{
			"status":   "up",
			"database": "ok",
		},
	)
}
