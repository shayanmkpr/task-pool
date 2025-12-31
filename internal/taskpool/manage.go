package taskpool

import (
	"fmt"
	"time"

	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

type workerManager struct {
	workers []*Worker
	store   *store.MemoryStore
}

func NewWorkerManager(workerCount int, store *store.MemoryStore) *workerManager {
	return &workerManager{
		workers: make([]*Worker, workerCount),
		store:   store,
	}
}

func (wm *workerManager) InitiateWorkers(pool *TaskPool) {
	for i := range wm.workers {
		w := NewWorker(i+1, pool)
		w.Start()
		wm.workers[i] = w
	}
}

func (wm *workerManager) MonitorWorkers(log *logger.Logger) {
	for _, worker := range wm.workers {
		w := worker
		go func() {
			for assigned := range w.Assigned {
				if assigned != nil {
					fmt.Printf("worker %d assigned to %v \n", w.ID, assigned.ID)
					log.Info(fmt.Sprintf("worker %d assigned to %v \n", w.ID, assigned.ID))
				} else {
					fmt.Printf("worker %d is free \n", w.ID)
					log.Info(fmt.Sprintf("worker %d is free\n", w.ID))
				}
			}
		}()
	}
}

func (wm *workerManager) WaitForCompletion(log *logger.Logger, waitingTime time.Duration) {
	for {
		allDone := true
		tasks := wm.store.ListTasks()
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

func (wm *workerManager) ForceStopWorkers() {
	for _, w := range wm.workers {
		w.Stop()
	}
}
