package raftadapter

import (
	"context"
	"sync"

	appcmd "ai-coding-training/internal/application/command"
	"ai-coding-training/internal/application/usecase"
)

// RecoveryManager coordinates node restart recovery. Raft itself restores its
// persisted term, vote, log and snapshot from Persister during construction;
// this manager owns the application-level recovery window and replays committed
// ApplyMsg entries that Raft emits after restart.
type RecoveryManager struct {
	gate    *usecase.RecoveryGate
	adapter CommandAdapter
	machine interface {
		ApplyCommittedCommand(context.Context, appcmd.ConfigCommand) (any, error)
	}
	applyCh <-chan ApplyMsg
	done    chan struct{}
	once    sync.Once
}

func NewRecoveryManager(gate *usecase.RecoveryGate, adapter CommandAdapter, machine interface {
	ApplyCommittedCommand(context.Context, appcmd.ConfigCommand) (any, error)
}, applyCh <-chan ApplyMsg) *RecoveryManager {
	if gate == nil {
		gate = usecase.NewRecoveryGate(true)
	}
	return &RecoveryManager{gate: gate, adapter: adapter, machine: machine, applyCh: applyCh, done: make(chan struct{})}
}

func (m *RecoveryManager) IsRecovering() bool {
	return m != nil && m.gate.IsRecovering()
}

func (m *RecoveryManager) Recover(ctx context.Context) error {
	if m == nil || m.gate == nil {
		return nil
	}
	m.gate.Begin()
	defer m.gate.Complete()
	defer m.once.Do(func() { close(m.done) })

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-m.applyCh:
			if !ok {
				return nil
			}
			if !msg.CommandValid {
				continue
			}
			cmd, decoded, err := m.adapter.FromApplyMsg(msg)
			if err != nil {
				return err
			}
			if !decoded {
				continue
			}
			if _, err := m.machine.ApplyCommittedCommand(ctx, cmd); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (m *RecoveryManager) Done() <-chan struct{} {
	if m == nil {
		done := make(chan struct{})
		close(done)
		return done
	}
	return m.done
}
