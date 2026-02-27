package store

import (
	"context"

	taskEntity "github.com/shayanmkpr/task-pool/internal/domain/task"
	"github.com/shayanmkpr/task-pool/internal/infra"
)

type Store[T infra.Entity] struct {
	adapter infra.Persistence[T]
}

func NewStore[T infra.Entity](adapter infra.Persistence[T]) *Store[T] {
	return &Store[T]{
		adapter: adapter,
	}
}

func (s *Store[T]) AddTask(ctx context.Context, task *taskEntity.Task) error {
	err := s.adapter.Save(ctx, task)
	return err
}

func (s *Store[T]) GetTask(ctx context.Context, id string) (*taskEntity.Task, error) {
	task, err := s.adapter.Get(ctx, id)
	return task, err
}

func (s *Store[T]) ListTasks(ctx context.Context) ([]T, error) {
	tasks, err := s.adapter.List(ctx, infra.QueryFilter{})
	return tasks, err
}

func (s *Store[T]) UpdateTask(ctx context.Context, task *taskEntity.Task) error {
	return s.adapter.Update(ctx, task)
}
