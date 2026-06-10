package repository

import (
	"context"

	entity "ai-coding-training/internal/domain/entity"
)

type ConfigItemRepository interface {
	Create(ctx context.Context, item entity.ConfigItem) (entity.ConfigItem, error)
	Update(ctx context.Context, item entity.ConfigItem) (entity.ConfigItem, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (entity.ConfigItem, error)
	ListByGroup(ctx context.Context, groupID string, includeDeleted bool) ([]entity.ConfigItem, error)
	ExistsByKey(ctx context.Context, groupID, key string, excludeID string) (bool, error)
}
