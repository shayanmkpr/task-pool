package taskpool

import (
	"context"
	"errors"
	"fmt"

	taskEntity "github.com/shayanmkpr/task-pool/internal/domain/task"
	"github.com/shayanmkpr/task-pool/internal/infra/memory"
	"github.com/shayanmkpr/task-pool/internal/logger"
)

var ErrTaskQueueFull = errors.New("task queue is full") //fix

type TaskPool struct {
	PoolSize int
	Tasks    chan *taskEntity.Task
	Store    *memory.MemoryStore
}

func NewTaskPool(poolSize int, store *memory.MemoryStore) *TaskPool {
	return &TaskPool{
		PoolSize: poolSize,
		Tasks:    make(chan *taskEntity.Task, poolSize),
		Store:    store,
	}
}

func (p *TaskPool) AddTask(ctx context.Context, logger *logger.Logger, task *taskEntity.Task) (string, error) {

	if len(p.Tasks) >= p.PoolSize {
		logger.Info("task queue is full")
		return "", ErrTaskQueueFull //fix
	}

	task.Status = taskEntity.StatusPending
	if err := p.Store.AddTask(task); err != nil { //fix
		return "", fmt.Errorf("failed to store task: %w", err) //fix
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()

	case p.Tasks <- task:
		return task.ID, nil

	default:
		logger.Info("task queue is full")
		return "", ErrTaskQueueFull //fix
	}
}
