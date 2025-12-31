package taskpool

import (
	"context"
	"fmt"

	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

type TaskPool struct {
	PoolSize int
	Tasks    chan *models.Task
	Store    *store.MemoryStore
}

func NewTaskPool(poolSize int, store *store.MemoryStore) *TaskPool {
	return &TaskPool{
		PoolSize: poolSize,
		Tasks:    make(chan *models.Task, poolSize),
		Store:    store,
	}
}

func (p *TaskPool) AddTask(ctx context.Context, logger *logger.Logger, task *models.Task) (string, error) {

	if len(p.Tasks) >= p.PoolSize {
		logger.Info("task queue is full")
		return "", fmt.Errorf("task queue is full")
	}

	task.Status = models.Pending
	p.Store.AddTask(task)

	select {
	case <-ctx.Done():
		return "", ctx.Err()

	case p.Tasks <- task:
		return task.ID, nil

	default:
		logger.Info("task queue is full")
		return "", fmt.Errorf("task queue is full")
	}
}
