package postgres

import (
	"context"
	"fmt"

	taskEntity "github.com/shayanmkpr/task-pool/internal/domain/task"
	"github.com/shayanmkpr/task-pool/internal/infra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresAdapter[T taskEntity.Task] struct {
	db *gorm.DB
}

func NewPostgresAdapter[T taskEntity.Task](dsn string) (infra.Persistence[T], error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	var entity T
	if err := db.AutoMigrate(&entity); err != nil {
		return nil, err
	}

	return &postgresAdapter[T]{
		db: db,
	}, nil
}

func (p *postgresAdapter[T]) Save(
	ctx context.Context,
	entity *T,
) error {
	return p.db.WithContext(ctx).Create(entity).Error
}

func (p *postgresAdapter[T]) Update(
	ctx context.Context,
	entity *T,
) error {
	return p.db.WithContext(ctx).Save(entity).Error
}

func (p *postgresAdapter[T]) Delete(
	ctx context.Context,
	id infra.ID,
) error {

	var entity T

	return p.db.WithContext(ctx).
		Delete(&entity, id).
		Error
}
func (p *postgresAdapter[T]) Get(
	ctx context.Context,
	id infra.ID,
) (*T, error) {

	var entity T

	err := p.db.WithContext(ctx).
		First(&entity, id).
		Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (p *postgresAdapter[T]) List(
	ctx context.Context,
	filter infra.QueryFilter,
) ([]T, error) {

	var results []T

	db := p.db.WithContext(ctx)
	db = BuildFilter(db, filter)

	err := db.Find(&results).Error
	return results, err
}

func (p *postgresAdapter[T]) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func BuildFilter(db *gorm.DB, filter infra.QueryFilter) *gorm.DB {

	for key, value := range filter {

		switch {
		case key == "":
			continue

		default:
			db = db.Where(fmt.Sprintf("%s ?", key), value)
		}
	}

	return db
}
