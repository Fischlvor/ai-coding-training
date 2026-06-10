package usecase

import (
	vo "ai-coding-training/internal/domain/valueobject"
	"context"
	"fmt"
	"strings"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
)

func (s Service) CreateGrayRule(ctx context.Context, in dto.GrayRuleDTO) (dto.GrayRuleDTO, error) {
	if err := validateID(in.RecordID, "record_id"); err != nil {
		return dto.GrayRuleDTO{}, err
	}
	if strings.TrimSpace(in.RuleType) == "" {
		return dto.GrayRuleDTO{}, fmt.Errorf("%w: rule_type is required", ErrValidation)
	}
	if strings.TrimSpace(in.RuleValue) == "" {
		return dto.GrayRuleDTO{}, fmt.Errorf("%w: rule_value is required", ErrValidation)
	}
	if s.GrayRules == nil {
		return dto.GrayRuleDTO{}, fmt.Errorf("%w: gray rule repository is not configured", ErrValidation)
	}
	rule := entity.GrayRule{ID: in.ID, RecordID: in.RecordID, RuleType: vo.GrayRuleType(in.RuleType), RuleValue: in.RuleValue, Enabled: in.Enabled, CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()}
	created, err := s.GrayRules.Create(ctx, rule)
	if err != nil {
		return dto.GrayRuleDTO{}, err
	}
	return dto.GrayRuleDTO{ID: created.ID, RecordID: created.RecordID, RuleType: string(created.RuleType), RuleValue: created.RuleValue, Enabled: created.Enabled}, nil
}

func (s Service) UpdateGrayRule(ctx context.Context, in dto.GrayRuleDTO) (dto.GrayRuleDTO, error) {
	if err := validateID(in.ID, "id"); err != nil {
		return dto.GrayRuleDTO{}, err
	}
	if s.GrayRules == nil {
		return dto.GrayRuleDTO{}, fmt.Errorf("%w: gray rule repository is not configured", ErrValidation)
	}
	rule := entity.GrayRule{ID: in.ID, RecordID: in.RecordID, RuleType: vo.GrayRuleType(in.RuleType), RuleValue: in.RuleValue, Enabled: in.Enabled, UpdatedAt: s.now().Unix()}
	updated, err := s.GrayRules.Update(ctx, rule)
	if err != nil {
		return dto.GrayRuleDTO{}, err
	}
	return dto.GrayRuleDTO{ID: updated.ID, RecordID: updated.RecordID, RuleType: string(updated.RuleType), RuleValue: updated.RuleValue, Enabled: updated.Enabled}, nil
}

func (s Service) DeleteGrayRule(ctx context.Context, id string) error {
	if s.GrayRules == nil {
		return fmt.Errorf("%w: gray rule repository is not configured", ErrValidation)
	}
	return s.GrayRules.Delete(ctx, id)
}
