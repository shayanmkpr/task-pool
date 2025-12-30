package api

import (
	"net/http"

	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

type Handler struct {
	pool *taskpool.TaskPool
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	// parse request
	// add task to pool
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getTaskWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// fetch task
}

func (h *Handler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// fetch task
}
