package repository

import (
	"context"

	entity "ai-coding-training/internal/domain/entity"
)

type ConfigGroupRepository interface {
	Create(ctx context.Context, group entity.ConfigGroup) (entity.ConfigGroup, error)
	Update(ctx context.Context, group entity.ConfigGroup) (entity.ConfigGroup, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (entity.ConfigGroup, error)
	ListByScope(ctx context.Context, appID, environmentID string) ([]entity.ConfigGroup, error)
	ExistsByName(ctx context.Context, appID, environmentID, name string, excludeID string) (bool, error)
}
