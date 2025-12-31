package taskpool

import (
	"context"
	"fmt"

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

func (p *TaskPool) AddTask(ctx context.Context, task *models.Task) (string, error) {
	task.Status = models.Pending
	p.Store.AddTask(task)

	select {
	case <-ctx.Done():
		return "", ctx.Err()

	case p.Tasks <- task:
		return task.ID, nil

	default:
		return "", fmt.Errorf("task queue is full")
	}
}
