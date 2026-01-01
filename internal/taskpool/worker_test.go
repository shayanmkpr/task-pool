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
		Duration:    time.Duration(1), // 1 second (worker multiplies by time.Second)
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

// TestWorkerMultipleTasks tests worker processing multiple tasks sequentially
func TestWorkerMultipleTasks(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(10, store)
	worker := NewWorker(1, pool)

	worker.Start()
	defer worker.Stop()

	const numTasks = 3
	tasks := make([]*models.Task, numTasks)

	// Create and send tasks
	for i := 0; i < numTasks; i++ {
		task := &models.Task{
			ID:          "task-" + string(rune(i+'0')),
			Title:       "Task",
			Description: "Test task",
			Duration:    time.Duration(1), // 1 second
			Status:      models.Pending,
		}
		tasks[i] = task
		store.AddTask(task)
		pool.Tasks <- task
	}

	// Wait for all tasks to complete (1 second each, sequential processing = 3 seconds + buffer)
	time.Sleep(100 * time.Millisecond) // Give tasks time to start
	for _, task := range tasks {
		if !waitForTaskCompletion(store, task.ID, 5*time.Second) {
			t.Errorf("Task %s did not complete", task.ID)
		}
	}

	// Verify all tasks are completed
	ctx := context.Background()
	for _, task := range tasks {
		storedTask, err := store.GetTask(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve task %s: %v", task.ID, err)
		}
		if storedTask.Status != models.Completed {
			t.Errorf("Task %s was not completed. Status: %s", task.ID, storedTask.Status)
		}
	}
}

// TestWorkerMultipleWorkers tests multiple workers processing tasks
func TestWorkerMultipleWorkers(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(10, store)

	const numWorkers = 3
	workers := make([]*Worker, numWorkers)

	// Start multiple workers
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i+1, pool)
		worker.Start()
		workers[i] = worker
		defer worker.Stop()
	}

	const numTasks = 6
	tasks := make([]*models.Task, numTasks)

	// Create and send tasks
	for i := 0; i < numTasks; i++ {
		task := &models.Task{
			ID:          "multi-task-" + string(rune(i+'0')),
			Title:       "Multi Task",
			Description: "Test task",
			Duration:    time.Duration(1), // 1 second
			Status:      models.Pending,
		}
		tasks[i] = task
		store.AddTask(task)
		pool.Tasks <- task
	}

	// Wait for all tasks to complete (1 second each, 6 tasks with 3 workers = ~2 seconds + buffer)
	time.Sleep(100 * time.Millisecond) // Give tasks time to start
	for _, task := range tasks {
		if !waitForTaskCompletion(store, task.ID, 5*time.Second) {
			t.Errorf("Task %s did not complete", task.ID)
		}
	}

	// Verify all tasks are completed
	ctx := context.Background()
	for _, task := range tasks {
		storedTask, err := store.GetTask(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve task %s: %v", task.ID, err)
		}
		if storedTask.Status != models.Completed {
			t.Errorf("Task %s was not completed. Status: %s", task.ID, storedTask.Status)
		}
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
		Duration:    time.Duration(1), // 1 second
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
