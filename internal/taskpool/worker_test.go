package taskpool

import (
	"context"
	"testing"
	"time"

	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

// waitForTaskCompletion waits until a task is completed or timeout
func waitForTaskCompletion(store *store.MemoryStore, taskID string, timeout time.Duration) bool {
	ctx := context.Background()
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		task, err := store.GetTask(ctx, taskID)
		if err == nil && task.Status == models.Completed {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

// TestNewWorker tests creating a new worker
func TestNewWorker(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(5, store)
	worker := NewWorker(1, pool)

	if worker.ID != 1 {
		t.Errorf("Expected worker ID 1, got %d", worker.ID)
	}

	if worker.TaskPool != pool {
		t.Error("TaskPool was not set correctly")
	}
}

// TestWorkerProcessTask tests worker processing a single task
func TestWorkerProcessTask(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(5, store)
	worker := NewWorker(1, pool)

	worker.Start()
	defer worker.Stop()

	task := &models.Task{
		ID:          "test-task-1",
		Title:       "Test Task",
		Description: "This is a test task",
		Duration:    1, // 1 second (worker multiplies by time.Second)
		Status:      models.Pending,
	}

	store.AddTask(task)
	pool.Tasks <- task

	// Wait for task to complete (1 second task + buffer)
	if !waitForTaskCompletion(store, task.ID, 3*time.Second) {
		t.Fatal("Task did not complete within timeout")
	}

	// Verify task status
	ctx := context.Background()
	storedTask, err := store.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve task: %v", err)
	}

	if storedTask.Status != models.Completed {
		t.Errorf("Expected task status 'completed', got '%s'", storedTask.Status)
	}
}

// TestWorkerLongRunningTask tests worker processing a long-running task
func TestWorkerLongRunningTask(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(5, store)
	worker := NewWorker(1, pool)

	worker.Start()
	defer worker.Stop()

	task := &models.Task{
		ID:          "long-task",
		Title:       "Long Task",
		Description: "Task that takes a long time",
		Duration:    1, // 1 second
		Status:      models.Pending,
	}

	store.AddTask(task)
	pool.Tasks <- task

	// Wait for task to complete (1 second task + buffer)
	if !waitForTaskCompletion(store, task.ID, 3*time.Second) {
		t.Fatal("Long task did not complete within timeout")
	}

	// Verify task status
	ctx := context.Background()
	storedTask, err := store.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve task: %v", err)
	}

	if storedTask.Status != models.Completed {
		t.Errorf("Expected task status 'completed', got '%s'", storedTask.Status)
	}
}
