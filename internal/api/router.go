package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func RegisterTaskRoutes(mux *http.ServeMux, h *Handler) {
	routes := []struct {
		method  string
		pattern string
		handler http.HandlerFunc
	}{
		{"POST", "/tasks", h.createTask},
		{"GET", "/tasks/{id}", h.getTaskWithID},
		{"GET", "/tasks/", h.getAllTasks},
	}

	fmt.Println("\nRegistered routes:")
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.method, route.pattern)
		mux.HandleFunc(pattern, recoverPanic(route.handler))
		fmt.Printf("  %-6s %s\n", route.method, route.pattern)
	}
	fmt.Println()
}

func recoverPanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered in HTTP handler",
					"error", err,
					"url", r.URL.String(),
					"method", r.Method,
					"stack", string(debug.Stack()),
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error"}`))
			}
		}()
		next(w, r)
	}
}
