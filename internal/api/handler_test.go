package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

func createTestHandler() (*Handler, *store.MemoryStore, *taskpool.TaskPool) {
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(5, store)
	log := logger.NewTestLogger()
	handler := NewHandler(pool, store, log)
	return handler, store, pool
}

// TestCreateTaskSuccess tests successful task creation
func TestCreateTaskSuccess(t *testing.T) {
	handler, store, _ := createTestHandler()

	taskReq := TaskRequest{
		Title:       "Test Task",
		Description: "This is a test task",
	}
	jsonData, _ := json.Marshal(taskReq)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	taskID, ok := response["id"].(string)
	if !ok {
		t.Fatal("Response should contain 'id' field")
	}

	// Verify task was stored
	ctx := context.Background()
	storedTask, err := store.GetTask(ctx, taskID)
	if err != nil {
		t.Fatalf("Failed to retrieve stored task: %v", err)
	}

	if storedTask.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", storedTask.Title)
	}

	if storedTask.Status != models.Pending {
		t.Errorf("Expected status 'pending', got '%s'", storedTask.Status)
	}
}

// TestCreateTaskInvalidJSON tests invalid JSON input
func TestCreateTaskInvalidJSON(t *testing.T) {
	handler, _, _ := createTestHandler()

	req := httptest.NewRequest("POST", "/tasks", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateTaskMissingTitle tests task creation without title
func TestCreateTaskMissingTitle(t *testing.T) {
	handler, _, _ := createTestHandler()

	taskReq := TaskRequest{
		Title:       "",
		Description: "This is a test task",
	}
	jsonData, _ := json.Marshal(taskReq)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateTaskPoolFull tests task creation when pool is full (by design)
func TestCreateTaskPoolFull(t *testing.T) {
	store := store.NewMemoryStore()
	pool := taskpool.NewTaskPool(1, store)
	log := logger.NewTestLogger()
	handler := NewHandler(pool, store, log)

	// Fill the pool
	longTask := &models.Task{
		ID:          "long-task",
		Title:       "Long Task",
		Description: "A task that takes a long time",
		Duration:    5,
		Status:      models.Pending,
	}
	store.AddTask(longTask)
	pool.Tasks <- longTask

	// Try to add another task - should fail because pool is full (by design)
	taskReq := TaskRequest{
		Title:       "New Task",
		Description: "This task should fail to be added",
	}
	jsonData, _ := json.Marshal(taskReq)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.createTask(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}

	// Clean up
	<-pool.Tasks
}

// TestGetTaskWithIDSuccess tests successful task retrieval
func TestGetTaskWithIDSuccess(t *testing.T) {
	handler, store, _ := createTestHandler()

	task := &models.Task{
		ID:          "test-id-123",
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      models.Pending,
		Duration:    1,
	}
	store.AddTask(task)

	mux := http.NewServeMux()
	RegisterTaskRoutes(mux, handler)

	req := httptest.NewRequest("GET", "/tasks/test-id-123", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != "test-id-123" {
		t.Errorf("Expected ID 'test-id-123', got '%s'", response.ID)
	}
}

// TestGetTaskWithIDNotFound tests retrieving non-existent task
func TestGetTaskWithIDNotFound(t *testing.T) {
	handler, _, _ := createTestHandler()

	mux := http.NewServeMux()
	RegisterTaskRoutes(mux, handler)

	req := httptest.NewRequest("GET", "/tasks/non-existent-id", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestGetAllTasksSuccess tests retrieving all tasks
func TestGetAllTasksSuccess(t *testing.T) {
	handler, store, _ := createTestHandler()

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

	req := httptest.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.getAllTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []*models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(response))
	}
}

// TestRecoverPanic tests panic recovery middleware
func TestRecoverPanic(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrappedHandler := recoverPanic(panicHandler)

	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expected := `{"error": "Internal server error"}`
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("Expected response to contain '%s', got '%s'", expected, w.Body.String())
	}
}
