package repository

import (
	vo "ai-coding-training/internal/domain/valueobject"
	"context"

	entity "ai-coding-training/internal/domain/entity"
)

type ReleaseRecordRepository interface {
	Create(ctx context.Context, record entity.ReleaseRecord) (entity.ReleaseRecord, error)
	UpdateStatus(ctx context.Context, id string, status vo.ReleaseStatus, remark string) (entity.ReleaseRecord, error)
	GetByID(ctx context.Context, id string) (entity.ReleaseRecord, error)
	ListByGroup(ctx context.Context, groupID string) ([]entity.ReleaseRecord, error)
}
