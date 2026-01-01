package taskpool

import (
	"context"
	"testing"

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
		Duration:    5,
	}
	store.AddTask(longTask)
	pool.Tasks <- longTask

	// Try to add another task - should fail (by design)
	newTask := &models.Task{
		ID:          "new-task",
		Title:       "New Task",
		Description: "This task should fail",
		Duration:    1,
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
