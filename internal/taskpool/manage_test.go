package taskpool

import (
	"context"
	"testing"
	"time"

	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

// TestNewWorkerManager tests creating a new worker manager
func TestNewWorkerManager(t *testing.T) {
	store := store.NewMemoryStore()
	manager := NewWorkerManager(5, store)

	if len(manager.workers) != 5 {
		t.Errorf("Expected 5 workers, got %d", len(manager.workers))
	}
}

// TestInitiateWorkers tests initiating workers
func TestInitiateWorkers(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(10, store)
	manager := NewWorkerManager(3, store)

	manager.InitiateWorkers(pool)
	defer manager.ForceStopWorkers()

	// Give workers time to start
	time.Sleep(100 * time.Millisecond)

	// Verify workers were created
	for i, worker := range manager.workers {
		if worker == nil {
			t.Errorf("Worker %d was not created", i+1)
		}
		if worker.ID != i+1 {
			t.Errorf("Expected worker ID %d, got %d", i+1, worker.ID)
		}
	}
}

// TestWaitForCompletion tests waiting for all tasks to complete
func TestWaitForCompletion(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(10, store)
	manager := NewWorkerManager(2, store)
	log := logger.NewTestLogger()

	manager.InitiateWorkers(pool)
	defer manager.ForceStopWorkers()

	// Add tasks
	const numTasks = 5
	for i := 0; i < numTasks; i++ {
		task := &models.Task{
			ID:          "wait-task-" + string(rune(i+'0')),
			Title:       "Wait Task",
			Description: "Task for wait completion test",
			Duration:    time.Duration(1), // 1 second
			Status:      models.Pending,
		}
		store.AddTask(task)
		pool.Tasks <- task
	}

	// Wait for completion with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan bool)
	go func() {
		manager.WaitForCompletion(ctx, log, 100*time.Millisecond)
		done <- true
	}()

	select {
	case <-done:
		// All tasks completed - verify
		ctx := context.Background()
		tasks, err := store.ListTasks(ctx)
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		// Verify all tasks are completed
		for _, task := range tasks {
			if task.Status != models.Completed {
				t.Errorf("Task %s was not completed. Status: %s", task.ID, task.Status)
			}
		}
	case <-time.After(6 * time.Second):
		t.Error("WaitForCompletion timed out")
	}
}

// TestForceStopWorkers tests forcefully stopping all workers
func TestForceStopWorkers(t *testing.T) {
	store := store.NewMemoryStore()
	pool := NewTaskPool(10, store)
	manager := NewWorkerManager(3, store)

	manager.InitiateWorkers(pool)

	// Give workers time to start
	time.Sleep(100 * time.Millisecond)

	// Force stop all workers
	manager.ForceStopWorkers()

	// Give workers time to stop
	time.Sleep(100 * time.Millisecond)
}
