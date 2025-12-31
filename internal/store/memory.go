package store

import (
	"context"
	"errors"
	"sync"

	"github.com/shayanmkpr/task-pool/internal/models"
)

type MemoryStore struct {
	mu    sync.RWMutex            // for reading memory safe
	tasks map[string]*models.Task // assigining ids to tasks
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[string]*models.Task),
	}
}

func (s *MemoryStore) AddTask(task *models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *MemoryStore) GetTask(ctx context.Context, id string) (*models.Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]

	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (s *MemoryStore) ListTasks(ctx context.Context) ([]*models.Task, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *MemoryStore) UpdateTask(task *models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}
