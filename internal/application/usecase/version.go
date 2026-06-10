package usecase

import (
	"context"
	"fmt"
	"strings"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
	vo "ai-coding-training/internal/domain/valueobject"
)

func (s Service) SaveDraftVersion(ctx context.Context, in dto.ConfigVersionDTO) (dto.ConfigVersionDTO, error) {
	if err := validateID(in.GroupID, "group_id"); err != nil {
		return dto.ConfigVersionDTO{}, err
	}
	if strings.TrimSpace(in.Title) == "" {
		return dto.ConfigVersionDTO{}, fmt.Errorf("%w: title is required", ErrValidation)
	}
	if strings.TrimSpace(in.Content) == "" {
		return dto.ConfigVersionDTO{}, fmt.Errorf("%w: content is required", ErrValidation)
	}
	versionNo := in.VersionNo
	if versionNo == 0 && s.Versions != nil {
		generated, err := s.Versions.NextVersion(ctx, in.GroupID)
		if err != nil {
			return dto.ConfigVersionDTO{}, err
		}
		versionNo = generated
	}
	version := entity.ConfigVersion{
		ID:          in.ID,
		GroupID:     in.GroupID,
		VersionNo:   versionNo,
		Title:       in.Title,
		Content:     in.Content,
		Status:      vo.VersionStatusDraft,
		PublishedAt: in.PublishedAt,
		CreatedAt:   s.now().Unix(),
		UpdatedAt:   s.now().Unix(),
	}
	if s.VersionRepo == nil {
		return dto.ConfigVersionDTO{}, fmt.Errorf("%w: version repository is not configured", ErrValidation)
	}
	created, err := s.VersionRepo.SaveDraft(ctx, version)
	if err != nil {
		return dto.ConfigVersionDTO{}, err
	}
	return dto.ConfigVersionDTO{ID: created.ID, GroupID: created.GroupID, VersionNo: created.VersionNo, Title: created.Title, Content: created.Content, Status: string(created.Status), PublishedAt: created.PublishedAt}, nil
}

func (s Service) ListVersions(ctx context.Context, groupID string) ([]dto.ConfigVersionDTO, error) {
	if err := validateID(groupID, "group_id"); err != nil {
		return nil, err
	}
	if s.VersionRepo == nil {
		return nil, fmt.Errorf("%w: version repository is not configured", ErrValidation)
	}
	items, err := s.VersionRepo.ListByGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ConfigVersionDTO, 0, len(items))
	for _, v := range items {
		out = append(out, dto.ConfigVersionDTO{ID: v.ID, GroupID: v.GroupID, VersionNo: v.VersionNo, Title: v.Title, Content: v.Content, Status: string(v.Status), PublishedAt: v.PublishedAt})
	}
	return out, nil
}
