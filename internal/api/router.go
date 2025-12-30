package api

import "net/http"

func RegisterTaskRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /tasks", h.createTask)
	mux.HandleFunc("GET /tasks/{id}", h.getTaskWithID)
	mux.HandleFunc("GET /tasks/", h.getAllTasks)
}
