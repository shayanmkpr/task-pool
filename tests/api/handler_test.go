package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shayanmkpr/task-pool/internal/api"
	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

// TestCreateTaskSuccess tests successful task creation
func TestCreateTaskSuccess(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create test request
	taskReq := api.TaskRequest{
		Title:       "Test Task",
		Description: "This is a test task",
	}
	jsonData, _ := json.Marshal(taskReq)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.CreateTask(w, req)

	// Verify
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, exists := response["id"]; !exists {
		t.Error("Response should contain 'id' field")
	}

	// Verify task was stored
	taskID := response["id"].(string)
	storedTask, err := store.GetTask(context.Background(), taskID)
	if err != nil {
		t.Fatalf("Failed to retrieve stored task: %v", err)
	}

	if storedTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", storedTask.Title)
	}
}

// TestCreateTaskInvalidMethod tests creating a task with invalid HTTP method
func TestCreateTaskInvalidMethod(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create test request with invalid method
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	// Execute
	handler.CreateTask(w, req)

	// Verify
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestCreateTaskInvalidJSON tests creating a task with invalid JSON
func TestCreateTaskInvalidJSON(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create test request with invalid JSON
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.CreateTask(w, req)

	// Verify
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateTaskMissingTitle tests creating a task without title
func TestCreateTaskMissingTitle(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create test request without title
	taskReq := api.TaskRequest{
		Title:       "", // Empty title
		Description: "This is a test task",
	}
	jsonData, _ := json.Marshal(taskReq)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.CreateTask(w, req)

	// Verify
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateTaskEmptyTitle tests creating a task with empty title
func TestCreateTaskEmptyTitle(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create test request with empty title
	taskReq := api.TaskRequest{
		Title:       "   ", // Whitespace only
		Description: "This is a test task",
	}
	jsonData, _ := json.Marshal(taskReq)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.CreateTask(w, req)

	// Verify
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetTaskWithIDSuccess tests successful retrieval of a task by ID
func TestGetTaskWithIDSuccess(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Add a task to the store
	task := &models.Task{
		ID:          "test-id-123",
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      models.Pending,
		Duration:    1,
	}
	store.AddTask(task)

	// Create test request
	req := httptest.NewRequest("GET", "/tasks/test-id-123", nil)
	req.Header.Set("Content-Type", "application/json")
	// Add the path value manually for testing
	req = req.WithContext(context.WithValue(req.Context(), "id", "test-id-123"))
	w := httptest.NewRecorder()

	// Execute
	handler.GetTaskWithID(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Task
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != "test-id-123" {
		t.Errorf("Expected ID 'test-id-123', got '%s'", response.ID)
	}
}

// TestGetTaskWithIDMissingID tests retrieving a task without providing ID
func TestGetTaskWithIDMissingID(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create test request without ID
	req := httptest.NewRequest("GET", "/tasks/", nil)
	w := httptest.NewRecorder()

	// Execute
	handler.GetTaskWithID(w, req)

	// Verify
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestGetTaskWithIDNotFound tests retrieving a non-existent task
func TestGetTaskWithIDNotFound(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create test request for non-existent task
	req := httptest.NewRequest("GET", "/tasks/non-existent-id", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "id", "non-existent-id"))
	w := httptest.NewRecorder()

	// Execute
	handler.GetTaskWithID(w, req)

	// Verify
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestGetAllTasksSuccess tests successful retrieval of all tasks
func TestGetAllTasksSuccess(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Add some tasks to the store
	task1 := &models.Task{
		ID:          "test-id-1",
		Title:       "Test Task 1",
		Description: "This is test task 1",
		Status:      models.Pending,
		Duration:    1,
	}
	task2 := &models.Task{
		ID:          "test-id-2",
		Title:       "Test Task 2",
		Description: "This is test task 2",
		Status:      models.Running,
		Duration:    2,
	}
	store.AddTask(task1)
	store.AddTask(task2)

	// Create test request
	req := httptest.NewRequest("GET", "/tasks/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.GetAllTasks(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []*models.Task
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(response))
	}
}

// TestGetAllTasksInvalidMethod tests retrieving all tasks with invalid HTTP method
func TestGetAllTasksInvalidMethod(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create test request with invalid method
	req := httptest.NewRequest("POST", "/tasks/", nil)
	w := httptest.NewRecorder()

	// Execute
	handler.GetAllTasks(w, req)

	// Verify
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestCreateTaskPoolFull tests creating a task when the pool is full
func TestCreateTaskPoolFull(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(1, store) // Pool size of 1
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Add a long-running task to fill the pool
	longTask := &models.Task{
		ID:          "long-task",
		Title:       "Long Task",
		Description: "A task that takes a long time",
		Duration:    5 * time.Second, // 5 seconds
	}
	// Add the long task directly to the pool to fill it
	pool.Tasks <- longTask

	// Create test request for new task
	taskReq := api.TaskRequest{
		Title:       "New Task",
		Description: "This task should fail to be added",
	}
	jsonData, _ := json.Marshal(taskReq)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	handler.CreateTask(w, req)

	// Verify - should return 500 Internal Server Error for full queue
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestCreateTaskWithTimeout tests creating a task with context timeout
func TestCreateTaskWithTimeout(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(1, store) // Pool size of 1
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Add a long-running task to fill the pool
	longTask := &models.Task{
		ID:          "long-task",
		Title:       "Long Task",
		Description: "A task that takes a long time",
		Duration:    5 * time.Second, // 5 seconds
	}
	// Add the long task directly to the pool to fill it
	pool.Tasks <- longTask

	// Create test request for new task with timeout context
	taskReq := api.TaskRequest{
		Title:       "New Task",
		Description: "This task should timeout",
	}
	jsonData, _ := json.Marshal(taskReq)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.CreateTask(w, req)

	// Verify - should return 408 Request Timeout
	if w.Code != http.StatusRequestTimeout {
		t.Errorf("Expected status %d, got %d", http.StatusRequestTimeout, w.Code)
	}
}

// TestGetAllTasksWithTimeout tests getting all tasks with context timeout
func TestGetAllTasksWithTimeout(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Create test request with timeout context
	req := httptest.NewRequest("GET", "/tasks/", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Execute
	handler.GetAllTasks(w, req)

	// Verify - should return 408 Request Timeout
	if w.Code != http.StatusRequestTimeout {
		t.Errorf("Expected status %d, got %d", http.StatusRequestTimeout, w.Code)
	}
}

// TestGetTaskWithIDWithTimeout tests getting a task with context timeout
func TestGetTaskWithIDWithTimeout(t *testing.T) {
	// Setup
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewLogger()
	handler := api.NewHandler(pool, store, log)

	// Add a task to the store
	task := &models.Task{
		ID:          "test-id-456",
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      models.Pending,
		Duration:    1,
	}
	store.AddTask(task)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Create test request with timeout context
	req := httptest.NewRequest("GET", "/tasks/test-id-456", nil)
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	// Add the path value manually for testing
	req = req.WithContext(context.WithValue(req.Context(), "id", "test-id-456"))
	w := httptest.NewRecorder()

	// Execute
	handler.GetTaskWithID(w, req)

	// Verify - should return 408 Request Timeout
	if w.Code != http.StatusRequestTimeout {
		t.Errorf("Expected status %d, got %d", http.StatusRequestTimeout, w.Code)
	}
}

// Helper function to access the unexported methods in the handler
// Since the methods are unexported, we need to create wrapper methods in the handler
// or use reflection to access them. For now, I'll create a helper function that
// creates a new handler with the required methods exposed for testing.
func TestPanicRecoveryInHandler(t *testing.T) {
	// This test would be for the recoverPanic middleware
	// Since we can't directly test the middleware, we'll create a mock handler
	// that panics and wrap it with the recoverPanic function
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrappedHandler := api.RecoverPanic(panicHandler)

	req := httptest.NewRequest("GET", "/tasks/", nil)
	w := httptest.NewRecorder()

	// Execute
	wrappedHandler.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Check that the response body contains the expected error message
	expected := `{"error": "Internal server error"}`
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("Expected response to contain '%s', got '%s'", expected, w.Body.String())
	}
}

// Add methods to the handler to expose the unexported methods for testing
func (h *api.Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.CreateTask(w, r)
}

func (h *api.Handler) GetTaskWithID(w http.ResponseWriter, r *http.Request) {
	h.GetTaskWithID(w, r)
}

func (h *api.Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	h.GetAllTasks(w, r)
}

// Helper function to simulate the PathValue for testing
func WithTaskID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, "id", id)
}