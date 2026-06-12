package usecase

import (
	"context"
	"sync/atomic"
)

// RecoveryGate prevents new writes from being acknowledged while a node is
// replaying committed Raft log entries after restart. The gate is intentionally
// small so it can be owned by the application layer without depending on Raft
// internals.
type RecoveryGate struct {
	recovering atomic.Bool
}

func NewRecoveryGate(recovering bool) *RecoveryGate {
	g := &RecoveryGate{}
	g.recovering.Store(recovering)
	return g
}

func (g *RecoveryGate) IsRecovering() bool {
	return g != nil && g.recovering.Load()
}

func (g *RecoveryGate) Begin() {
	if g != nil {
		g.recovering.Store(true)
	}
}

func (g *RecoveryGate) Complete() {
	if g != nil {
		g.recovering.Store(false)
	}
}

// RecoveryController is implemented by infrastructure components that know how
// to replay durable committed log entries into the application state machine.
type RecoveryController interface {
	Recover(ctx context.Context) error
	IsRecovering() bool
}
