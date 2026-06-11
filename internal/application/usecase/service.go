package usecase

import (
	"time"

	port "ai-coding-training/internal/application/port/repository"
)

type Clock func() time.Time

type Service struct {
	Now              Clock
	Versions         port.VersionAllocator
	VersionRepo      port.VersionRepository
	Releases         port.ReleaseRecordRepository
	GrayRules        port.GrayRuleRepository
	GrayRecords      port.GrayRecordRepository
	SubscriptionRepo port.SubscriptionRepository
	ChangeEvents     port.ChangeEventRepository
	Apps             port.AppRepository
	Environments     port.EnvironmentRepository
	Groups           port.ConfigGroupRepository
	Items            port.ConfigItemRepository
}
