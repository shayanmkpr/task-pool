package taskpool

import (
	"time"

	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

type workerManager struct {
	Workers []*Worker
	Store   *store.MemoryStore
}

func NewWorkerManager(workerCount int, store *store.MemoryStore) *workerManager {
	return &workerManager{
		Workers: make([]*Worker, workerCount),
		Store:   store,
	}
}

func (wm *workerManager) InitiateWorkers(pool *TaskPool) {
	for i := range wm.Workers {
		w := NewWorker(i+1, pool)
		w.Start()
		wm.Workers[i] = w
	}
}

func (wm *workerManager) WaitForCompletion(log *logger.Logger, waitingTime time.Duration) {
	for {
		allDone := true
		tasks := wm.Store.ListTasks()
		for _, t := range tasks {
			if t.Status != models.Completed {
				allDone = false
				break
			}
		}
		if allDone {
			return
		}
		time.Sleep(waitingTime)
	}
}
