package infra

import (
	"context"

	taskEntity "github.com/shayanmkpr/task-pool/internal/domain/task"
)

type QueryFilter map[string]any

type ID string

// list the types included in the Persistence.
type Entity interface {
	taskEntity.Task
}

type Persistence[T Entity] interface {
	Save(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id ID) error
	Get(ctx context.Context, id ID) (*T, error)
	List(ctx context.Context, filter QueryFilter) ([]T, error)
	Close() error
}
