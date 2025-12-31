package application

import (
	"github.com/shayanmkpr/task-pool/internal/models"
	"github.com/shayanmkpr/task-pool/internal/taskpool"
)

type TaskService struct {
	pool *taskpool.TaskPool
}

func NewTaskService(pool *taskpool.TaskPool) *TaskService {
	return &TaskService{
		pool: pool,
	}
}

func (s *TaskService) AddTask(task *models.Task) (string, error) {
	return s.pool.AddTask(task)
}
