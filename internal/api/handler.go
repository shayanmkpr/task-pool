package api

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

const (
	maxRequestBodySize = 1 << 20 // 1MB //fix
	minTaskDuration    = 1       // seconds //fix
	maxTaskDuration    = 5       // seconds //fix
	maxTitleLength     = 200     // characters //fix
	maxDescLength      = 1000    // characters //fix
)

type Handler struct {
	pool   *taskpool.TaskPool
	store  *store.MemoryStore
	logger *logger.Logger
}

func NewHandler(pool *taskpool.TaskPool, store *store.MemoryStore, logger *logger.Logger) *Handler {
	return &Handler{
		pool:   pool,
		store:  store,
		logger: logger,
	}
}

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("createTask handler called", "method", r.Method, "url", r.URL.String())

	if r.Method != http.MethodPost {
		h.logger.Warn("method not allowed", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize) //fix

	var req TaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		h.logger.Warn("title is required")
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if len(title) > maxTitleLength { //fix
		h.logger.Warn("title too long")                        //fix
		http.Error(w, "title too long", http.StatusBadRequest) //fix
		return                                                 //fix
	}
	if len(req.Description) > maxDescLength { //fix
		h.logger.Warn("description too long")                        //fix
		http.Error(w, "description too long", http.StatusBadRequest) //fix
		return                                                       //fix
	}

	ctx := r.Context()

	// Generate a Unique ID
	newUUID := uuid.New().String()

	h.logger.Info("adding task to pool", "task_id", newUUID, "title", req.Title)

	taskID, err := h.pool.AddTask(ctx, h.logger, &models.Task{
		ID:          newUUID,
		Title:       title, //fix
		Description: req.Description,
		Duration:    rand.Intn(maxTaskDuration-minTaskDuration+1) + minTaskDuration, //fix
	})
	if err != nil {
		h.logger.Error("failed to add task to pool", "error", err)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request cancelled", http.StatusRequestTimeout)
			return
		}
		// check if the error is becuase the task queue is full
		if errors.Is(err, taskpool.ErrTaskQueueFull) { //fix
			http.Error(w, "task queue is full", http.StatusTooManyRequests) //fix
			return                                                          //fix
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("task added successfully", "task_id", taskID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"id": taskID,
	}); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		return
	}
	h.logger.Info("response sent successfully", "task_id", taskID)
}

func (h *Handler) getTaskWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.logger.Info("getTaskWithID handler called", "task_id", id, "method", r.Method, "url", r.URL.String())

	if id == "" {
		h.logger.Warn("task id is required")
		http.Error(w, "task id is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context() // graceful shutdown

	h.logger.Info("retrieving task from store", "task_id", id)

	task, err := h.store.GetTask(ctx, id)
	if err != nil {
		h.logger.Error("failed to retrieve task", "error", err, "task_id", id)
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	h.logger.Info("task retrieved successfully", "task_id", id)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.logger.Error("failed to encode response", "error", err, "task_id", id)
		return
	}
	h.logger.Info("response sent successfully", "task_id", id)
}

func (h *Handler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("getAllTasks handler called", "method", r.Method, "url", r.URL.String())

	if r.Method != http.MethodGet {
		h.logger.Warn("method not allowed", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	h.logger.Info("retrieving all tasks from store")

	tasks, err := h.store.ListTasks(ctx)
	if err != nil {
		h.logger.Error("failed to retrieve tasks", "error", err)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request cancelled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("tasks retrieved successfully", "count", len(tasks))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		return
	}
	h.logger.Info("response sent successfully", "count", len(tasks))
}
