package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/shayanmkpr/task-pool/internal/domain/task"
)

type MemoryStore struct {
	mu    sync.RWMutex          // for reading memory safe
	tasks map[string]*task.Task // assigining ids to tasks
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[string]*task.Task),
	}
}

func (s *MemoryStore) AddTask(task *task.Task) error { //fix
	if task == nil { //fix
		return errors.New("task cannot be nil") //fix
	}
	if task.ID == "" { //fix
		return errors.New("task ID cannot be empty") //fix
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
	return nil
}

func (s *MemoryStore) GetTask(ctx context.Context, id string) (*task.Task, error) {
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

func (s *MemoryStore) ListTasks(ctx context.Context) ([]*task.Task, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]*task.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *MemoryStore) UpdateTask(task *task.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}
