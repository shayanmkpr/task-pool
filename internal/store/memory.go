package store

import (
	"errors"
	"log"
	"sync"
	"time"

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

func (s *MemoryStore) GetTask(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (s *MemoryStore) ListTasks() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *MemoryStore) UpdateTask(task *models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
}

func (s *MemoryStore) WaitForCompletion(waitingTime time.Duration) {
	for {
		allDone := true
		tasks := s.ListTasks()
		for _, t := range tasks {
			if t.Status != models.Completed {
				allDone = false
				break
			}
		}
		if allDone {
			log.Println("All tasks completed")
			return
		}
		time.Sleep(waitingTime)
	}
}
