package taskpool

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/shayanmkpr/task-pool/internal/adapter"
	"github.com/shayanmkpr/task-pool/internal/logger"
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/store"
)

var ErrTaskQueueFull = errors.New("task queue is full") //fix

type TaskPool struct {
	PoolSize int
	Tasks    chan *models.Task
	Store    *store.MemoryStore
	CBack    adapter.CallBack
}

func NewTaskPool(poolSize int, store *store.MemoryStore) *TaskPool {
	return &TaskPool{
		PoolSize: poolSize,
		Tasks:    make(chan *models.Task, poolSize),
		Store:    store,
		CBack: adapter.CallBack{
			EndPoint: "http://127.0.0.1:8000",
			Client:   &http.Client{},
		},
	}
}

func (p *TaskPool) AddTask(ctx context.Context, logger *logger.Logger, task *models.Task) (string, error) {

	if len(p.Tasks) >= p.PoolSize {
		logger.Info("task queue is full")
		return "", ErrTaskQueueFull //fix
	}

	task.Status = models.Pending
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
