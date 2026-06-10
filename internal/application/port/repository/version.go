package repository

import (
	"context"

	entity "ai-coding-training/internal/domain/entity"
)

type VersionAllocator interface {
	NextVersion(ctx context.Context, groupID string) (int64, error)
}

type VersionRepository interface {
	SaveDraft(ctx context.Context, version entity.ConfigVersion) (entity.ConfigVersion, error)
	ListByGroup(ctx context.Context, groupID string) ([]entity.ConfigVersion, error)
	GetByVersionNo(ctx context.Context, groupID string, versionNo int64) (entity.ConfigVersion, error)
	GetLatestDraft(ctx context.Context, groupID string) (entity.ConfigVersion, error)
	ExistsByVersionNo(ctx context.Context, groupID string, versionNo int64, excludeID string) (bool, error)
}
