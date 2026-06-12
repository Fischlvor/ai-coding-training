package usecase

import (
	"context"
	"fmt"
	"strings"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
	vo "ai-coding-training/internal/domain/valueobject"
)

func (s Service) RegisterSubscription(ctx context.Context, in dto.SubscriptionDTO) (dto.SubscriptionDTO, error) {
	if err := validateID(in.ClientID, "client_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if err := validateID(in.AppID, "app_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if err := validateID(in.EnvironmentID, "environment_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if err := validateID(in.GroupID, "group_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if s.SubscriptionRepo == nil {
		return dto.SubscriptionDTO{}, fmt.Errorf("%w: subscription repository is not configured", ErrValidation)
	}
	sub := entity.Subscription{ID: in.ID, ClientID: in.ClientID, AppID: in.AppID, EnvironmentID: in.EnvironmentID, GroupID: in.GroupID, LastVersionNo: in.LastVersionNo, CreatedAt: s.now().Unix(), UpdatedAt: s.now().Unix()}
	created, err := s.SubscriptionRepo.Create(ctx, sub)
	if err != nil {
		return dto.SubscriptionDTO{}, err
	}
	return subscriptionToDTO(created), nil
}

func (s Service) UpsertSubscription(ctx context.Context, in dto.SubscriptionDTO) (dto.SubscriptionDTO, error) {
	if err := validateID(in.ClientID, "client_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if err := validateID(in.AppID, "app_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if err := validateID(in.EnvironmentID, "environment_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if err := validateID(in.GroupID, "group_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if s.SubscriptionRepo == nil {
		return dto.SubscriptionDTO{}, fmt.Errorf("%w: subscription repository is not configured", ErrValidation)
	}
	existing, err := s.SubscriptionRepo.GetByClientAndGroup(ctx, in.ClientID, in.AppID, in.EnvironmentID, in.GroupID)
	if err == nil && strings.TrimSpace(existing.ID) != "" {
		existing.LastVersionNo = in.LastVersionNo
		existing.UpdatedAt = s.now().Unix()
		updated, err := s.SubscriptionRepo.Update(ctx, existing)
		if err != nil {
			return dto.SubscriptionDTO{}, err
		}
		return subscriptionToDTO(updated), nil
	}
	if err != nil && err != ErrNotFound {
		return dto.SubscriptionDTO{}, err
	}
	return s.RegisterSubscription(ctx, in)
}

func (s Service) RecordChangeEvent(ctx context.Context, in dto.ChangeEventDTO) (dto.ChangeEventDTO, error) {
	if err := validateID(in.AppID, "app_id"); err != nil {
		return dto.ChangeEventDTO{}, err
	}
	if err := validateID(in.EnvironmentID, "environment_id"); err != nil {
		return dto.ChangeEventDTO{}, err
	}
	if err := validateID(in.GroupID, "group_id"); err != nil {
		return dto.ChangeEventDTO{}, err
	}
	if in.VersionNo <= 0 {
		return dto.ChangeEventDTO{}, fmt.Errorf("%w: version_no is required", ErrValidation)
	}
	if strings.TrimSpace(in.EventType) == "" {
		return dto.ChangeEventDTO{}, fmt.Errorf("%w: event_type is required", ErrValidation)
	}
	if s.ChangeEvents == nil {
		return dto.ChangeEventDTO{}, fmt.Errorf("%w: change event repository is not configured", ErrValidation)
	}
	if strings.TrimSpace(in.ID) != "" {
		existing, err := s.ChangeEvents.GetByID(ctx, in.ID)
		if err == nil && existing.ID != "" {
			return changeEventToDTO(existing), nil
		}
		if err != nil && err != ErrNotFound {
			return dto.ChangeEventDTO{}, err
		}
	}
	event := entity.ChangeEvent{ID: in.ID, AppID: in.AppID, EnvironmentID: in.EnvironmentID, GroupID: in.GroupID, VersionNo: in.VersionNo, EventType: vo.ChangeEventType(in.EventType), Payload: in.Payload, CreatedAt: s.now().Unix()}
	created, err := s.ChangeEvents.Create(ctx, event)
	if err != nil {
		return dto.ChangeEventDTO{}, err
	}
	return changeEventToDTO(created), nil
}

func (s Service) ListChangeEvents(ctx context.Context, appID, environmentID, groupID string, sinceVersion int64) ([]dto.ChangeEventDTO, error) {
	if s.ChangeEvents == nil {
		return nil, fmt.Errorf("%w: change event repository is not configured", ErrValidation)
	}
	items, err := s.ChangeEvents.ListByGroup(ctx, appID, environmentID, groupID, sinceVersion)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ChangeEventDTO, 0, len(items))
	for _, item := range items {
		out = append(out, changeEventToDTO(item))
	}
	return out, nil
}

func (s Service) ListPendingEvents(ctx context.Context, subscriptionID string) ([]dto.ChangeEventDTO, error) {
	if err := validateID(subscriptionID, "subscription_id"); err != nil {
		return nil, err
	}
	if s.ChangeEvents == nil {
		return nil, fmt.Errorf("%w: change event repository is not configured", ErrValidation)
	}
	items, err := s.ChangeEvents.ListPendingBySubscription(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	return uniqueChangeEvents(items), nil
}

func (s Service) CompensateSubscription(ctx context.Context, clientID, subscriptionID string, lastVersionNo int64) ([]dto.ChangeEventDTO, error) {
	if err := validateID(subscriptionID, "subscription_id"); err != nil {
		return nil, err
	}
	if lastVersionNo < 0 {
		return nil, fmt.Errorf("%w: last_version_no must be non-negative", ErrValidation)
	}
	if s.SubscriptionRepo == nil {
		return nil, fmt.Errorf("%w: subscription repository is not configured", ErrValidation)
	}
	if s.ChangeEvents == nil {
		return nil, fmt.Errorf("%w: change event repository is not configured", ErrValidation)
	}
	sub, err := s.SubscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(clientID) != "" && sub.ClientID != clientID {
		return nil, fmt.Errorf("%w: subscription client mismatch", ErrValidation)
	}
	if lastVersionNo > sub.LastVersionNo {
		sub.LastVersionNo = lastVersionNo
		sub.UpdatedAt = s.now().Unix()
		if _, err := s.SubscriptionRepo.Update(ctx, sub); err != nil {
			return nil, err
		}
	}
	items, err := s.ChangeEvents.ListByGroup(ctx, sub.AppID, sub.EnvironmentID, sub.GroupID, sub.LastVersionNo)
	if err != nil {
		return nil, err
	}
	return uniqueChangeEvents(items), nil
}

func (s Service) AckSubscription(ctx context.Context, subscriptionID string, lastVersionNo int64) (dto.SubscriptionDTO, error) {
	if err := validateID(subscriptionID, "subscription_id"); err != nil {
		return dto.SubscriptionDTO{}, err
	}
	if s.SubscriptionRepo == nil {
		return dto.SubscriptionDTO{}, fmt.Errorf("%w: subscription repository is not configured", ErrValidation)
	}
	sub, err := s.SubscriptionRepo.GetByID(ctx, subscriptionID)
	if err != nil {
		return dto.SubscriptionDTO{}, err
	}
	sub.LastVersionNo = lastVersionNo
	sub.UpdatedAt = s.now().Unix()
	updated, err := s.SubscriptionRepo.Update(ctx, sub)
	if err != nil {
		return dto.SubscriptionDTO{}, err
	}
	return subscriptionToDTO(updated), nil
}

func subscriptionToDTO(sub entity.Subscription) dto.SubscriptionDTO {
	return dto.SubscriptionDTO{ID: sub.ID, ClientID: sub.ClientID, AppID: sub.AppID, EnvironmentID: sub.EnvironmentID, GroupID: sub.GroupID, LastVersionNo: sub.LastVersionNo}
}

func changeEventToDTO(ev entity.ChangeEvent) dto.ChangeEventDTO {
	return dto.ChangeEventDTO{ID: ev.ID, AppID: ev.AppID, EnvironmentID: ev.EnvironmentID, GroupID: ev.GroupID, VersionNo: ev.VersionNo, EventType: string(ev.EventType), Payload: ev.Payload}
}

func uniqueChangeEvents(items []entity.ChangeEvent) []dto.ChangeEventDTO {
	out := make([]dto.ChangeEventDTO, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		key := item.ID
		if strings.TrimSpace(key) == "" {
			key = fmt.Sprintf("%s:%s:%s:%d:%s", item.AppID, item.EnvironmentID, item.GroupID, item.VersionNo, item.EventType)
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, changeEventToDTO(item))
	}
	return out
}
