package usecase

import (
	"context"
	"fmt"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
)

func (s Service) CreateApp(ctx context.Context, in dto.AppDTO) (dto.AppDTO, error) {
	if err := validateName(in.Name); err != nil {
		return dto.AppDTO{}, err
	}
	exists, err := s.Apps.ExistsByName(ctx, in.Name, "")
	if err != nil {
		return dto.AppDTO{}, err
	}
	if exists {
		return dto.AppDTO{}, fmt.Errorf("%w: app name already exists", ErrConflict)
	}
	created, err := s.Apps.Create(ctx, entity.App{ID: in.ID, Name: in.Name, Description: in.Description, CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()})
	if err != nil {
		return dto.AppDTO{}, err
	}
	return dto.AppDTO{ID: created.ID, Name: created.Name, Description: created.Description}, nil
}

func (s Service) UpdateApp(ctx context.Context, in dto.AppDTO) (dto.AppDTO, error) {
	if err := validateID(in.ID, "id"); err != nil {
		return dto.AppDTO{}, err
	}
	if err := validateName(in.Name); err != nil {
		return dto.AppDTO{}, err
	}
	if ok, err := s.Apps.ExistsByName(ctx, in.Name, in.ID); err != nil {
		return dto.AppDTO{}, err
	} else if ok {
		return dto.AppDTO{}, fmt.Errorf("%w: app name already exists", ErrConflict)
	}
	updated, err := s.Apps.Update(ctx, entity.App{ID: in.ID, Name: in.Name, Description: in.Description, UpdatedAt: s.now().Unix()})
	if err != nil {
		return dto.AppDTO{}, err
	}
	return dto.AppDTO{ID: updated.ID, Name: updated.Name, Description: updated.Description}, nil
}

func (s Service) DeleteApp(ctx context.Context, id string) error {
	return s.Apps.Delete(ctx, id)
}

func (s Service) GetApp(ctx context.Context, id string) (dto.AppDTO, error) {
	v, err := s.Apps.Get(ctx, id)
	if err != nil {
		return dto.AppDTO{}, err
	}
	return dto.AppDTO{ID: v.ID, Name: v.Name, Description: v.Description}, nil
}
func (s Service) ListApps(ctx context.Context) ([]dto.AppDTO, error) {
	items, err := s.Apps.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.AppDTO, 0, len(items))
	for _, v := range items {
		out = append(out, dto.AppDTO{ID: v.ID, Name: v.Name, Description: v.Description})
	}
	return out, nil
}
