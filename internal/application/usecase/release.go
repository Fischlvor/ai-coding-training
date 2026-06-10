package usecase

import (
	"context"
	"fmt"
	"strings"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
	vo "ai-coding-training/internal/domain/valueobject"
)

func (s Service) PublishVersion(ctx context.Context, versionID string, targetVersion int64, remark string) (dto.ReleaseRecordDTO, error) {
	if err := validateID(versionID, "version_id"); err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	if targetVersion <= 0 {
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: target_version is required", ErrValidation)
	}
	if s.VersionRepo == nil {
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: version repository is not configured", ErrValidation)
	}
	if s.Releases == nil {
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: release repository is not configured", ErrValidation)
	}
	version, err := s.VersionRepo.GetByVersionNo(ctx, "", targetVersion)
	if err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	if strings.TrimSpace(version.ID) == "" {
		return dto.ReleaseRecordDTO{}, ErrNotFound
	}
	if existing, err := s.Releases.GetByID(ctx, versionID); err == nil {
		if existing.Action == vo.ReleaseActionPublish && existing.TargetVersion == targetVersion {
			return dto.ReleaseRecordDTO{ID: existing.ID, GroupID: existing.GroupID, VersionID: existing.VersionID, TargetVersion: existing.TargetVersion, Action: string(existing.Action), Status: string(existing.Status), Remark: existing.Remark}, nil
		}
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: release id already exists", ErrConflict)
	} else if err != ErrNotFound {
		return dto.ReleaseRecordDTO{}, err
	}
	record := entity.ReleaseRecord{ID: versionID, GroupID: version.GroupID, VersionID: version.ID, TargetVersion: targetVersion, Action: vo.ReleaseActionPublish, Status: vo.ReleaseStatusPending, Remark: remark, CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()}
	created, err := s.Releases.Create(ctx, record)
	if err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	published := version
	published.Status = vo.VersionStatusPublished
	published.PublishedAt = s.now().Unix()
	published.UpdatedAt = s.now().Unix()
	if _, err := s.VersionRepo.SaveDraft(ctx, published); err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	updated, err := s.Releases.UpdateStatus(ctx, created.ID, vo.ReleaseStatusSuccess, remark)
	if err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	return dto.ReleaseRecordDTO{ID: updated.ID, GroupID: updated.GroupID, VersionID: updated.VersionID, TargetVersion: updated.TargetVersion, Action: string(updated.Action), Status: string(updated.Status), Remark: updated.Remark}, nil
}

func (s Service) RollbackVersion(ctx context.Context, versionID string, targetVersion int64, remark string) (dto.ReleaseRecordDTO, error) {
	if err := validateID(versionID, "version_id"); err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	if targetVersion <= 0 {
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: target_version is required", ErrValidation)
	}
	if s.VersionRepo == nil {
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: version repository is not configured", ErrValidation)
	}
	if s.Releases == nil {
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: release repository is not configured", ErrValidation)
	}
	version, err := s.VersionRepo.GetByVersionNo(ctx, "", targetVersion)
	if err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	if strings.TrimSpace(version.ID) == "" {
		return dto.ReleaseRecordDTO{}, ErrNotFound
	}
	if existing, err := s.Releases.GetByID(ctx, versionID); err == nil {
		if existing.Action == vo.ReleaseActionRollback && existing.TargetVersion == targetVersion {
			return dto.ReleaseRecordDTO{ID: existing.ID, GroupID: existing.GroupID, VersionID: existing.VersionID, TargetVersion: existing.TargetVersion, Action: string(existing.Action), Status: string(existing.Status), Remark: existing.Remark}, nil
		}
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: release id already exists", ErrConflict)
	} else if err != ErrNotFound {
		return dto.ReleaseRecordDTO{}, err
	}
	record := entity.ReleaseRecord{ID: versionID, GroupID: version.GroupID, VersionID: version.ID, TargetVersion: targetVersion, Action: vo.ReleaseActionRollback, Status: vo.ReleaseStatusPending, Remark: remark, CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()}
	created, err := s.Releases.Create(ctx, record)
	if err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	updated, err := s.Releases.UpdateStatus(ctx, created.ID, vo.ReleaseStatusSuccess, remark)
	if err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	return dto.ReleaseRecordDTO{ID: updated.ID, GroupID: updated.GroupID, VersionID: updated.VersionID, TargetVersion: updated.TargetVersion, Action: string(updated.Action), Status: string(updated.Status), Remark: updated.Remark}, nil
}

func (s Service) GetReleaseRecord(ctx context.Context, id string) (dto.ReleaseRecordDTO, error) {
	if err := validateID(id, "release_id"); err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	if s.Releases == nil {
		return dto.ReleaseRecordDTO{}, fmt.Errorf("%w: release repository is not configured", ErrValidation)
	}
	v, err := s.Releases.GetByID(ctx, id)
	if err != nil {
		return dto.ReleaseRecordDTO{}, err
	}
	return dto.ReleaseRecordDTO{ID: v.ID, GroupID: v.GroupID, VersionID: v.VersionID, TargetVersion: v.TargetVersion, Action: string(v.Action), Status: string(v.Status), Remark: v.Remark}, nil
}

func (s Service) ListReleaseRecords(ctx context.Context, groupID string) ([]dto.ReleaseRecordDTO, error) {
	if err := validateID(groupID, "group_id"); err != nil {
		return nil, err
	}
	if s.Releases == nil {
		return nil, fmt.Errorf("%w: release repository is not configured", ErrValidation)
	}
	items, err := s.Releases.ListByGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ReleaseRecordDTO, 0, len(items))
	for _, v := range items {
		out = append(out, dto.ReleaseRecordDTO{ID: v.ID, GroupID: v.GroupID, VersionID: v.VersionID, TargetVersion: v.TargetVersion, Action: string(v.Action), Status: string(v.Status), Remark: v.Remark})
	}
	return out, nil
}
