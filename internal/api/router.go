package api

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		fmt.Println("=======================")
		fmt.Println("Method:", r.Method)
		fmt.Println("Path:", r.URL.Path)
		fmt.Println("Body:", string(bodyBytes))
		fmt.Println("Time:", start.Format(time.RFC3339))
		fmt.Println("=======================")

		next.ServeHTTP(w, r)
		fmt.Printf("Completed in %v\n\n", time.Since(start))
	})
}

func RegisterTaskRoutes(mux *http.ServeMux, h *Handler) {
	routes := []struct {
		method  string
		pattern string
		handler http.HandlerFunc
	}{
		{"POST", "/tasks", h.createTask},
		{"GET", "/tasks/{id}", h.getTaskWithID},
		{"GET", "/tasks", h.getAllTasks},
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
