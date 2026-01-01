package taskpool

import (
	"context"
	"testing"
	"time"

	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

// TestNewTaskPool tests creating a new task pool
func TestNewTaskPool(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(5, store)

	if pool.PoolSize != 5 {
		t.Errorf("Expected pool size 5, got %d", pool.PoolSize)
	}

	if cap(pool.Tasks) != 5 {
		t.Errorf("Expected task channel capacity 5, got %d", cap(pool.Tasks))
	}
}

// TestAddTaskSuccess tests successfully adding a task
func TestAddTaskSuccess(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(5, store)
	log := logger.NewTestLogger()

	task := &models.Task{
		ID:          "test-task-1",
		Title:       "Test Task",
		Description: "This is a test task",
		Duration:    1 * time.Second,
	}

	ctx := context.Background()
	taskID, err := pool.AddTask(ctx, log, task)
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	if taskID != "test-task-1" {
		t.Errorf("Expected task ID 'test-task-1', got '%s'", taskID)
	}

	// Verify task was stored
	storedTask, err := store.GetTask(ctx, taskID)
	if err != nil {
		t.Fatalf("Failed to retrieve stored task: %v", err)
	}

	if storedTask.Status != models.Pending {
		t.Errorf("Expected task status 'pending', got '%s'", storedTask.Status)
	}

	// Verify task was added to channel
	select {
	case taskFromChannel := <-pool.Tasks:
		if taskFromChannel.ID != taskID {
			t.Errorf("Task in channel doesn't match. Expected %s, got %s", taskID, taskFromChannel.ID)
		}
	case <-time.After(1 * time.Second):
		t.Error("Task was not added to channel")
	}
}

// TestAddTaskPoolFull tests adding task when pool is full (by design)
func TestAddTaskPoolFull(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(1, store)
	log := logger.NewTestLogger()

	// Fill the pool
	longTask := &models.Task{
		ID:          "long-task",
		Title:       "Long Task",
		Description: "A task that takes a long time",
		Duration:    5 * time.Second,
	}
	store.AddTask(longTask)
	pool.Tasks <- longTask

	// Try to add another task - should fail (by design)
	newTask := &models.Task{
		ID:          "new-task",
		Title:       "New Task",
		Description: "This task should fail",
		Duration:    1 * time.Second,
	}

	ctx := context.Background()
	_, err := pool.AddTask(ctx, log, newTask)
	if err == nil {
		t.Error("Expected error when pool is full")
	}

	if err.Error() != "task queue is full" {
		t.Errorf("Expected 'task queue is full' error, got '%v'", err)
	}

	// Clean up
	<-pool.Tasks
}
