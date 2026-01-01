package taskpool

import (
	"testing"
	"time"

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
