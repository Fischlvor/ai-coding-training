package usecase

import (
	"context"
	"fmt"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
)

func (s Service) CreateEnvironment(ctx context.Context, in dto.EnvironmentDTO) (dto.EnvironmentDTO, error) {
	if err := validateName(in.Name); err != nil {
		return dto.EnvironmentDTO{}, err
	}
	exists, err := s.Environments.ExistsByName(ctx, in.Name, "")
	if err != nil {
		return dto.EnvironmentDTO{}, err
	}
	if exists {
		return dto.EnvironmentDTO{}, fmt.Errorf("%w: environment name already exists", ErrConflict)
	}
	created, err := s.Environments.Create(ctx, entity.Environment{ID: in.ID, Name: in.Name, Description: in.Description, CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()})
	if err != nil {
		return dto.EnvironmentDTO{}, err
	}
	return dto.EnvironmentDTO{ID: created.ID, Name: created.Name, Description: created.Description}, nil
}

func (s Service) UpdateEnvironment(ctx context.Context, in dto.EnvironmentDTO) (dto.EnvironmentDTO, error) {
	if err := validateID(in.ID, "id"); err != nil {
		return dto.EnvironmentDTO{}, err
	}
	if err := validateName(in.Name); err != nil {
		return dto.EnvironmentDTO{}, err
	}
	if ok, err := s.Environments.ExistsByName(ctx, in.Name, in.ID); err != nil {
		return dto.EnvironmentDTO{}, err
	} else if ok {
		return dto.EnvironmentDTO{}, fmt.Errorf("%w: environment name already exists", ErrConflict)
	}
	updated, err := s.Environments.Update(ctx, entity.Environment{ID: in.ID, Name: in.Name, Description: in.Description, UpdatedAt: s.now().Unix()})
	if err != nil {
		return dto.EnvironmentDTO{}, err
	}
	return dto.EnvironmentDTO{ID: updated.ID, Name: updated.Name, Description: updated.Description}, nil
}

func (s Service) DeleteEnvironment(ctx context.Context, id string) error {
	return s.Environments.Delete(ctx, id)
}
func (s Service) GetEnvironment(ctx context.Context, id string) (dto.EnvironmentDTO, error) {
	v, err := s.Environments.Get(ctx, id)
	if err != nil {
		return dto.EnvironmentDTO{}, err
	}
	return dto.EnvironmentDTO{ID: v.ID, Name: v.Name, Description: v.Description}, nil
}
func (s Service) ListEnvironments(ctx context.Context) ([]dto.EnvironmentDTO, error) {
	items, err := s.Environments.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.EnvironmentDTO, 0, len(items))
	for _, v := range items {
		out = append(out, dto.EnvironmentDTO{ID: v.ID, Name: v.Name, Description: v.Description})
	}
	return out, nil
}
