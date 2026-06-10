package repository

import (
	vo "ai-coding-training/internal/domain/valueobject"
	"context"

	entity "ai-coding-training/internal/domain/entity"
)

type GrayRuleRepository interface {
	Create(ctx context.Context, rule entity.GrayRule) (entity.GrayRule, error)
	Update(ctx context.Context, rule entity.GrayRule) (entity.GrayRule, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (entity.GrayRule, error)
	ListByRecord(ctx context.Context, recordID string) ([]entity.GrayRule, error)
}

type GrayRecordRepository interface {
	Create(ctx context.Context, record entity.GrayRecord) (entity.GrayRecord, error)
	UpdateStatus(ctx context.Context, id string, status vo.GrayRecordStatus, matchedCount int64) (entity.GrayRecord, error)
	Get(ctx context.Context, id string) (entity.GrayRecord, error)
	ListByRelease(ctx context.Context, releaseID string) ([]entity.GrayRecord, error)
}
