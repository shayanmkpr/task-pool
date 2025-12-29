package taskpool

import (
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

func (p *TaskPool) AddTask(task *models.Task) {
	task.Status = models.Pending
	p.Store.AddTask(task)
	p.Tasks <- task // adding to a buffered channel that works as a limitted queue.
}
