package repository

import (
	"context"

	entity "ai-coding-training/internal/domain/entity"
)

type AppRepository interface {
	Create(ctx context.Context, app entity.App) (entity.App, error)
	Update(ctx context.Context, app entity.App) (entity.App, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (entity.App, error)
	List(ctx context.Context) ([]entity.App, error)
	ExistsByName(ctx context.Context, name string, excludeID string) (bool, error)
}
