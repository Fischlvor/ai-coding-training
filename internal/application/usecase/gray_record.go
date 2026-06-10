package usecase

import (
	"context"
	"fmt"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
	vo "ai-coding-training/internal/domain/valueobject"
)

func (s Service) MatchGrayTarget(ctx context.Context, releaseID string, matchedCount int64) (dto.GrayRecordDTO, error) {
	if err := validateID(releaseID, "release_id"); err != nil {
		return dto.GrayRecordDTO{}, err
	}
	if s.GrayRecords == nil {
		return dto.GrayRecordDTO{}, fmt.Errorf("%w: gray record repository is not configured", ErrValidation)
	}
	record := entity.GrayRecord{ID: releaseID, ReleaseID: releaseID, TargetVersion: matchedCount, MatchedCount: matchedCount, Status: vo.GrayRecordStatus(vo.GrayRecordStatusPending), CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()}
	created, err := s.GrayRecords.Create(ctx, record)
	if err != nil {
		return dto.GrayRecordDTO{}, err
	}
	updated, err := s.GrayRecords.UpdateStatus(ctx, created.ID, vo.GrayRecordStatus(vo.GrayRecordStatusApplied), matchedCount)
	if err != nil {
		return dto.GrayRecordDTO{}, err
	}
	return dto.GrayRecordDTO{ID: updated.ID, ReleaseID: updated.ReleaseID, TargetVersion: updated.TargetVersion, MatchedCount: updated.MatchedCount, Status: string(updated.Status)}, nil
}

func (s Service) ListGrayRecords(ctx context.Context, releaseID string) ([]dto.GrayRecordDTO, error) {
	if err := validateID(releaseID, "release_id"); err != nil {
		return nil, err
	}
	if s.GrayRecords == nil {
		return nil, fmt.Errorf("%w: gray record repository is not configured", ErrValidation)
	}
	items, err := s.GrayRecords.ListByRelease(ctx, releaseID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.GrayRecordDTO, 0, len(items))
	for _, v := range items {
		out = append(out, dto.GrayRecordDTO{ID: v.ID, ReleaseID: v.ReleaseID, TargetVersion: v.TargetVersion, MatchedCount: v.MatchedCount, Status: string(v.Status)})
	}
	return out, nil
}
