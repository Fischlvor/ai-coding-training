package usecase

import (
	"context"
	"fmt"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
)

func (s Service) CreateItem(ctx context.Context, in dto.ConfigItemDTO) (dto.ConfigItemDTO, error) {
	if err := validateName(in.Key); err != nil {
		return dto.ConfigItemDTO{}, err
	}
	if err := validateID(in.GroupID, "group_id"); err != nil {
		return dto.ConfigItemDTO{}, err
	}
	exists, err := s.Items.ExistsByKey(ctx, in.GroupID, in.Key, "")
	if err != nil {
		return dto.ConfigItemDTO{}, err
	}
	if exists {
		return dto.ConfigItemDTO{}, fmt.Errorf("%w: config key already exists", ErrConflict)
	}
	var versionNo int64
	if s.Versions != nil {
		versionNo, err = s.Versions.NextVersion(ctx, in.GroupID)
		if err != nil {
			return dto.ConfigItemDTO{}, err
		}
	}
	created, err := s.Items.Create(ctx, entity.ConfigItem{ID: in.ID, GroupID: in.GroupID, Key: in.Key, Value: in.Value, Description: in.Description, VersionNo: versionNo, CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()})
	if err != nil {
		return dto.ConfigItemDTO{}, err
	}
	return dto.ConfigItemDTO{ID: created.ID, GroupID: created.GroupID, Key: created.Key, Value: created.Value, Description: created.Description, VersionNo: created.VersionNo}, nil
}

func (s Service) UpdateItem(ctx context.Context, in dto.ConfigItemDTO) (dto.ConfigItemDTO, error) {
	if err := validateID(in.ID, "id"); err != nil {
		return dto.ConfigItemDTO{}, err
	}
	if err := validateName(in.Key); err != nil {
		return dto.ConfigItemDTO{}, err
	}
	if err := validateID(in.GroupID, "group_id"); err != nil {
		return dto.ConfigItemDTO{}, err
	}
	if ok, err := s.Items.ExistsByKey(ctx, in.GroupID, in.Key, in.ID); err != nil {
		return dto.ConfigItemDTO{}, err
	} else if ok {
		return dto.ConfigItemDTO{}, fmt.Errorf("%w: config key already exists", ErrConflict)
	}
	updated, err := s.Items.Update(ctx, entity.ConfigItem{ID: in.ID, GroupID: in.GroupID, Key: in.Key, Value: in.Value, Description: in.Description, VersionNo: in.VersionNo, UpdatedAt: s.now().Unix()})
	if err != nil {
		return dto.ConfigItemDTO{}, err
	}
	return dto.ConfigItemDTO{ID: updated.ID, GroupID: updated.GroupID, Key: updated.Key, Value: updated.Value, Description: updated.Description, VersionNo: updated.VersionNo}, nil
}

func (s Service) DeleteItem(ctx context.Context, id string) error { return s.Items.Delete(ctx, id) }
func (s Service) GetItem(ctx context.Context, id string) (dto.ConfigItemDTO, error) {
	v, err := s.Items.Get(ctx, id)
	if err != nil {
		return dto.ConfigItemDTO{}, err
	}
	return dto.ConfigItemDTO{ID: v.ID, GroupID: v.GroupID, Key: v.Key, Value: v.Value, Description: v.Description, VersionNo: v.VersionNo}, nil
}
func (s Service) ListItems(ctx context.Context, groupID string, includeDeleted bool) ([]dto.ConfigItemDTO, error) {
	items, err := s.Items.ListByGroup(ctx, groupID, includeDeleted)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ConfigItemDTO, 0, len(items))
	for _, v := range items {
		out = append(out, dto.ConfigItemDTO{ID: v.ID, GroupID: v.GroupID, Key: v.Key, Value: v.Value, Description: v.Description, VersionNo: v.VersionNo})
	}
	return out, nil
}
