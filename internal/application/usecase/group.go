package usecase

import (
	"context"
	"fmt"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
)

func (s Service) CreateGroup(ctx context.Context, in dto.ConfigGroupDTO) (dto.ConfigGroupDTO, error) {
	if err := validateName(in.Name); err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	if err := validateID(in.AppID, "app_id"); err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	if err := validateID(in.EnvironmentID, "environment_id"); err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	exists, err := s.Groups.ExistsByName(ctx, in.AppID, in.EnvironmentID, in.Name, "")
	if err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	if exists {
		return dto.ConfigGroupDTO{}, fmt.Errorf("%w: group name already exists", ErrConflict)
	}
	created, err := s.Groups.Create(ctx, entity.ConfigGroup{ID: in.ID, AppID: in.AppID, EnvironmentID: in.EnvironmentID, Name: in.Name, Description: in.Description, CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()})
	if err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	return dto.ConfigGroupDTO{ID: created.ID, AppID: created.AppID, EnvironmentID: created.EnvironmentID, Name: created.Name, Description: created.Description}, nil
}

func (s Service) UpdateGroup(ctx context.Context, in dto.ConfigGroupDTO) (dto.ConfigGroupDTO, error) {
	if err := validateID(in.ID, "id"); err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	if err := validateName(in.Name); err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	if err := validateID(in.AppID, "app_id"); err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	if err := validateID(in.EnvironmentID, "environment_id"); err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	if ok, err := s.Groups.ExistsByName(ctx, in.AppID, in.EnvironmentID, in.Name, in.ID); err != nil {
		return dto.ConfigGroupDTO{}, err
	} else if ok {
		return dto.ConfigGroupDTO{}, fmt.Errorf("%w: group name already exists", ErrConflict)
	}
	updated, err := s.Groups.Update(ctx, entity.ConfigGroup{ID: in.ID, AppID: in.AppID, EnvironmentID: in.EnvironmentID, Name: in.Name, Description: in.Description, UpdatedAt: s.now().Unix()})
	if err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	return dto.ConfigGroupDTO{ID: updated.ID, AppID: updated.AppID, EnvironmentID: updated.EnvironmentID, Name: updated.Name, Description: updated.Description}, nil
}

func (s Service) DeleteGroup(ctx context.Context, id string) error {
	return s.Groups.Delete(ctx, id)
}

func (s Service) GetGroup(ctx context.Context, id string) (dto.ConfigGroupDTO, error) {
	v, err := s.Groups.Get(ctx, id)
	if err != nil {
		return dto.ConfigGroupDTO{}, err
	}
	return dto.ConfigGroupDTO{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, Name: v.Name, Description: v.Description}, nil
}
func (s Service) ListGroups(ctx context.Context, appID, environmentID string) ([]dto.ConfigGroupDTO, error) {
	items, err := s.Groups.ListByScope(ctx, appID, environmentID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ConfigGroupDTO, 0, len(items))
	for _, v := range items {
		out = append(out, dto.ConfigGroupDTO{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, Name: v.Name, Description: v.Description})
	}
	return out, nil
}
