package usecase

import (
	"context"
	"fmt"
	"strings"

	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
	vo "ai-coding-training/internal/domain/valueobject"
)

type SDKClient struct {
	Service
}

type SDKFetchInput struct {
	AppID         string
	EnvironmentID string
	GroupID       string
	ClientID      string
}

type SDKSubscribeInput struct {
	ClientID      string
	AppID         string
	EnvironmentID string
	GroupID       string
	LastVersionNo int64
}

func (s Service) NewSDKClient() SDKClient { return SDKClient{Service: s} }

func (c SDKClient) FetchConfig(ctx context.Context, in SDKFetchInput) (dto.ConfigSnapshotDTO, error) {
	if err := validateID(in.AppID, "app_id"); err != nil {
		return dto.ConfigSnapshotDTO{}, err
	}
	if err := validateID(in.EnvironmentID, "environment_id"); err != nil {
		return dto.ConfigSnapshotDTO{}, err
	}
	if err := validateID(in.GroupID, "group_id"); err != nil {
		return dto.ConfigSnapshotDTO{}, err
	}
	if c.VersionRepo == nil {
		return dto.ConfigSnapshotDTO{}, fmt.Errorf("%w: version repository is not configured", ErrValidation)
	}
	versions, err := c.VersionRepo.ListByGroup(ctx, in.GroupID)
	if err != nil {
		return dto.ConfigSnapshotDTO{}, err
	}
	var latest entity.ConfigVersion
	for _, version := range versions {
		if version.VersionNo >= latest.VersionNo {
			latest = version
		}
	}
	matched, rule := c.matchGray(in.ClientID, latest)
	return dto.ConfigSnapshotDTO{
		AppID:           in.AppID,
		EnvironmentID:   in.EnvironmentID,
		GroupID:         in.GroupID,
		VersionNo:       latest.VersionNo,
		Title:           latest.Title,
		Content:         latest.Content,
		GrayMatched:     matched,
		GrayMatchedRule: rule,
		UpdatedAt:       latest.UpdatedAt,
	}, nil
}

func (c SDKClient) Subscribe(ctx context.Context, in SDKSubscribeInput) (dto.SubscriptionDTO, []dto.ChangeEventDTO, error) {
	if err := validateID(in.ClientID, "client_id"); err != nil {
		return dto.SubscriptionDTO{}, nil, err
	}
	sub, err := c.RegisterSubscription(ctx, dto.SubscriptionDTO{ClientID: in.ClientID, AppID: in.AppID, EnvironmentID: in.EnvironmentID, GroupID: in.GroupID, LastVersionNo: in.LastVersionNo})
	if err != nil {
		return dto.SubscriptionDTO{}, nil, err
	}
	events, err := c.ListPendingEvents(ctx, sub.ID)
	if err != nil {
		return dto.SubscriptionDTO{}, nil, err
	}
	return sub, events, nil
}

func (c SDKClient) RefreshCache(ctx context.Context, clientID, subscriptionID string, lastVersionNo int64) ([]dto.ChangeEventDTO, error) {
	return c.CompensateSubscription(ctx, clientID, subscriptionID, lastVersionNo)
}

func (c SDKClient) AckEvent(ctx context.Context, subscriptionID string, lastVersionNo int64) (dto.SubscriptionDTO, error) {
	return c.AckSubscription(ctx, subscriptionID, lastVersionNo)
}

func (c SDKClient) matchGray(clientID string, version entity.ConfigVersion) (bool, string) {
	if strings.TrimSpace(clientID) == "" {
		return false, ""
	}
	if version.ID == "" {
		return false, ""
	}
	if c.GrayRecords == nil {
		return false, ""
	}
	return version.Status == vo.VersionStatusPublished, "default"
}
