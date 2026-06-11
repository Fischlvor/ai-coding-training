package repository

import (
	"context"

	entity "ai-coding-training/internal/domain/entity"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription entity.Subscription) (entity.Subscription, error)
	Update(ctx context.Context, subscription entity.Subscription) (entity.Subscription, error)
	GetByID(ctx context.Context, id string) (entity.Subscription, error)
	GetByClientAndGroup(ctx context.Context, clientID, appID, environmentID, groupID string) (entity.Subscription, error)
	ListByGroup(ctx context.Context, appID, environmentID, groupID string) ([]entity.Subscription, error)
	Delete(ctx context.Context, id string) error
}

type ChangeEventRepository interface {
	Create(ctx context.Context, event entity.ChangeEvent) (entity.ChangeEvent, error)
	GetByID(ctx context.Context, id string) (entity.ChangeEvent, error)
	ListByGroup(ctx context.Context, appID, environmentID, groupID string, sinceVersion int64) ([]entity.ChangeEvent, error)
	ListPendingBySubscription(ctx context.Context, subscriptionID string) ([]entity.ChangeEvent, error)
}
