package taskpool

import (
	"context"
	"fmt"
	"time"

	"github.com/shayanmkpr/task-pool/internal/application/store"
	"github.com/shayanmkpr/task-pool/internal/domain/task"
	taskEntity "github.com/shayanmkpr/task-pool/internal/domain/task"
	"github.com/shayanmkpr/task-pool/internal/logger"
)

type workerManager struct {
	workers []*Worker
	memory  *store.Store[*taskEntity.Task]
}

func NewWorkerManager(workerCount int, store *store.Store[*taskEntity.Task]) *workerManager {
	return &workerManager{
		workers: make([]*Worker, workerCount),
		memory:  store,
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

func (wm *workerManager) WaitForCompletion(ctx context.Context, log *logger.Logger, waitingTime time.Duration) {
	for {
		select {
		case <-ctx.Done():
			log.Info("Context cancelled during WaitForCompletion", "error", ctx.Err())
			return
		default:
		}

		allDone := true
		tasks, err := wm.memory.ListTasks(ctx)
		if err != nil {
			log.Error("failed to retrieve tasks", "error", err)
			return
		}
		for _, t := range tasks {
			if t.Status != task.StatusCompleted {
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
