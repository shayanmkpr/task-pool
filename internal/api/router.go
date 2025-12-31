package api

import (
	"log/slog"
	"net/http"
)

func RegisterTaskRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /tasks", recoverPanic(h.createTask))
	mux.HandleFunc("GET /tasks/{id}", recoverPanic(h.getTaskWithID))
	mux.HandleFunc("GET /tasks/", recoverPanic(h.getAllTasks))
}

func recoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered in HTTP handler", "error", err, "url", r.URL.String(), "method", r.Method)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}
