package store

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/shayanmkpr/task-pool/internal/models"
)

// TestMemoryStoreBasic tests basic memory store functionality
func TestMemoryStoreBasic(t *testing.T) {
	store := NewMemoryStore()

	task := &models.Task{
		ID:          "basic-test-1",
		Title:       "Basic Test Task",
		Description: "Task for basic store functionality test",
		Status:      models.Pending,
		Duration:    1 * time.Second,
	}

	store.AddTask(task)

	ctx := context.Background()
	retrievedTask, err := store.GetTask(ctx, "basic-test-1")
	if err != nil {
		t.Fatalf("Failed to retrieve task: %v", err)
	}

	if retrievedTask.ID != task.ID {
		t.Errorf("Retrieved task ID doesn't match. Expected %s, got %s", task.ID, retrievedTask.ID)
	}

	if retrievedTask.Title != task.Title {
		t.Errorf("Retrieved task title doesn't match. Expected %s, got %s", task.Title, retrievedTask.Title)
	}
}

// TestMemoryStoreConcurrentSubmissions tests concurrent submissions of many tasks at once
func TestMemoryStoreConcurrentSubmissions(t *testing.T) {
	store := NewMemoryStore()

	const numTasks = 100
	var wg sync.WaitGroup
	wg.Add(numTasks)

	// Add multiple tasks concurrently
	for i := 0; i < numTasks; i++ {
		go func(id int) {
			defer wg.Done()
			task := &models.Task{
				ID:          "concurrent-task-" + string(rune(id+'0')),
				Title:       "Concurrent Task",
				Description: "Task for concurrent submission test",
				Status:      models.Pending,
				Duration:    1 * time.Second,
			}
			store.AddTask(task)
		}(i)
	}

	wg.Wait()

	// Verify all tasks were stored
	ctx := context.Background()
	tasks, err := store.ListTasks(ctx)
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != numTasks {
		t.Errorf("Expected %d tasks, got %d", numTasks, len(tasks))
	}
}

// TestMemoryStoreDuplicateIDs tests duplicate task IDs (should overwrite)
func TestMemoryStoreDuplicateIDs(t *testing.T) {
	store := NewMemoryStore()

	task1 := &models.Task{
		ID:          "duplicate-test",
		Title:       "First Task",
		Description: "First task with duplicate ID",
		Status:      models.Pending,
		Duration:    1 * time.Second,
	}

	task2 := &models.Task{
		ID:          "duplicate-test", // Same ID
		Title:       "Second Task",
		Description: "Second task with duplicate ID",
		Status:      models.Running,
		Duration:    2 * time.Second,
	}

	store.AddTask(task1)
	store.AddTask(task2) // This should overwrite the first task

	ctx := context.Background()
	retrievedTask, err := store.GetTask(ctx, "duplicate-test")
	if err != nil {
		t.Fatalf("Failed to retrieve task: %v", err)
	}

	// The second task should overwrite the first one
	if retrievedTask.Title != "Second Task" {
		t.Errorf("Task was not overwritten properly. Expected 'Second Task', got '%s'", retrievedTask.Title)
	}

	if retrievedTask.Status != models.Running {
		t.Errorf("Task status was not overwritten. Expected 'running', got '%s'", retrievedTask.Status)
	}
}

// TestMemoryStoreEmpty tests empty store
func TestMemoryStoreEmpty(t *testing.T) {
	store := NewMemoryStore()

	ctx := context.Background()
	tasks, err := store.ListTasks(ctx)
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks in empty store, got %d", len(tasks))
	}

	_, err = store.GetTask(ctx, "non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent task")
	}
}

// TestMemoryStoreUpdateTask tests updating a task
func TestMemoryStoreUpdateTask(t *testing.T) {
	store := NewMemoryStore()

	task := &models.Task{
		ID:          "update-test",
		Title:       "Original Title",
		Description: "Original Description",
		Status:      models.Pending,
		Duration:    1 * time.Second,
	}

	store.AddTask(task)

	// Update the task
	task.Title = "Updated Title"
	task.Status = models.Running
	store.UpdateTask(task)

	ctx := context.Background()
	updatedTask, err := store.GetTask(ctx, "update-test")
	if err != nil {
		t.Fatalf("Failed to retrieve updated task: %v", err)
	}

	if updatedTask.Title != "Updated Title" {
		t.Errorf("Task title was not updated. Expected 'Updated Title', got '%s'", updatedTask.Title)
	}

	if updatedTask.Status != models.Running {
		t.Errorf("Task status was not updated. Expected 'running', got '%s'", updatedTask.Status)
	}
}

// TestMemoryStoreGetTaskWithTimeout tests getting a task with context timeout
func TestMemoryStoreGetTaskWithTimeout(t *testing.T) {
	store := NewMemoryStore()

	task := &models.Task{
		ID:          "timeout-test",
		Title:       "Timeout Test Task",
		Description: "Task for timeout test",
		Status:      models.Pending,
		Duration:    1 * time.Second,
	}
	store.AddTask(task)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Nanosecond) // Ensure timeout has occurred

	_, err := store.GetTask(ctx, "timeout-test")
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}

	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Errorf("Expected context error, got %v", err)
	}
}

// TestMemoryStoreListTasksWithTimeout tests listing tasks with context timeout
func TestMemoryStoreListTasksWithTimeout(t *testing.T) {
	store := NewMemoryStore()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Nanosecond) // Ensure timeout has occurred

	_, err := store.ListTasks(ctx)
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}

	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Errorf("Expected context error, got %v", err)
	}
}

// TestMemoryStoreConcurrentReadWrite tests concurrent reads and writes
func TestMemoryStoreConcurrentReadWrite(t *testing.T) {
	store := NewMemoryStore()

	const numWriters = 10
	const numReaders = 20
	var wg sync.WaitGroup

	// Concurrent writes
	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			task := &models.Task{
				ID:          "concurrent-rw-" + string(rune(id+'0')),
				Title:       "Concurrent RW Task",
				Description: "Task for concurrent read/write test",
				Status:      models.Pending,
				Duration:    1 * time.Second,
			}
			store.AddTask(task)
		}(i)
	}

	// Concurrent reads
	wg.Add(numReaders)
	ctx := context.Background()
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			_, _ = store.ListTasks(ctx)
		}()
	}

	wg.Wait()

	// Verify all writes succeeded
	tasks, err := store.ListTasks(ctx)
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != numWriters {
		t.Errorf("Expected %d tasks, got %d", numWriters, len(tasks))
	}
}

// TestMemoryStoreLargeNumberOfTasks tests store with large number of tasks
func TestMemoryStoreLargeNumberOfTasks(t *testing.T) {
	store := NewMemoryStore()

	const numTasks = 1000
	ctx := context.Background()

	// Add many tasks
	for i := 0; i < numTasks; i++ {
		task := &models.Task{
			ID:          "large-test-" + string(rune(i+'0')),
			Title:       "Large Test Task",
			Description: "Task for large number test",
			Status:      models.Pending,
			Duration:    1 * time.Second,
		}
		store.AddTask(task)
	}

	// Verify all tasks are stored
	tasks, err := store.ListTasks(ctx)
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != numTasks {
		t.Errorf("Expected %d tasks, got %d", numTasks, len(tasks))
	}
}
