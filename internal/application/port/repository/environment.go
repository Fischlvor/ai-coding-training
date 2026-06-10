package repository

import (
	"context"

	entity "ai-coding-training/internal/domain/entity"
)

type EnvironmentRepository interface {
	Create(ctx context.Context, environment entity.Environment) (entity.Environment, error)
	Update(ctx context.Context, environment entity.Environment) (entity.Environment, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (entity.Environment, error)
	List(ctx context.Context) ([]entity.Environment, error)
	ExistsByName(ctx context.Context, name string, excludeID string) (bool, error)
}
