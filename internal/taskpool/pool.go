package taskpool

import (
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

func (p *TaskPool) AddTask(task *models.Task) error {
	task.Status = models.Pending
	p.Store.AddTask(task)
	select {
	case p.Tasks <- task:
		return nil
	default:
		return fmt.Errorf("task queue is full")
	}
}
