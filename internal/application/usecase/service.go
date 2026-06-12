package usecase

import (
	"context"
	"time"

	port "ai-coding-training/internal/application/port/repository"
)

type Clock func() time.Time

type CommandSubmitter interface {
	SubmitCommand(ctx context.Context, command any) (any, error)
}

type RaftStateProvider interface {
	GetState() (term int, isLeader bool)
	GetLeader() string
}

type RaftWriteCoordinator interface {
	CommandSubmitter
	RaftStateProvider
}

type Service struct {
	Now              Clock
	Raft             CommandSubmitter
	Recovery         *RecoveryGate
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
