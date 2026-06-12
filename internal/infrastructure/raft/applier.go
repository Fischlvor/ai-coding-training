package raftadapter

import (
	"context"
	"fmt"
	"sync"

	appcmd "ai-coding-training/internal/application/command"
	"ai-coding-training/internal/application/dto"
)

type ApplyStateMachine interface {
	ApplyCommittedCommand(context.Context, appcmd.ConfigCommand) (any, error)
	RecordChangeEvent(context.Context, dto.ChangeEventDTO) (dto.ChangeEventDTO, error)
}

type ApplyDispatcher struct {
	adapter CommandAdapter
	machine ApplyStateMachine
	results chan<- PendingResult

	mu      sync.Mutex
	applied map[string]struct{}
}

func NewApplyDispatcher(adapter CommandAdapter, machine ApplyStateMachine, results chan<- PendingResult) *ApplyDispatcher {
	return &ApplyDispatcher{adapter: adapter, machine: machine, results: results, applied: make(map[string]struct{})}
}

func (d *ApplyDispatcher) HandleApply(msg ApplyMsg) {
	if d == nil || !msg.CommandValid || d.adapter == nil || d.machine == nil {
		return
	}
	cmd, ok, err := d.adapter.FromApplyMsg(msg)
	if err != nil || !ok {
		d.notify(PendingResult{Command: cmd, Error: err})
		return
	}
	if d.isDuplicate(cmd, msg.CommandIndex) {
		d.notify(PendingResult{Command: cmd})
		return
	}
	if _, err := d.machine.ApplyCommittedCommand(context.Background(), cmd); err != nil {
		d.notify(PendingResult{Command: cmd, Error: err})
		return
	}
	if event, ok, err := d.adapter.ApplyMsgToEvent(msg); err != nil {
		d.notify(PendingResult{Command: cmd, Error: err})
		return
	} else if ok {
		if _, err := d.machine.RecordChangeEvent(context.Background(), event); err != nil {
			d.notify(PendingResult{Command: cmd, Error: err})
			return
		}
	}
	d.notify(PendingResult{Command: cmd})
}

func (d *ApplyDispatcher) isDuplicate(cmd appcmd.ConfigCommand, index int) bool {
	key := cmd.ID
	if key == "" {
		key = fmt.Sprintf("%s:%s:%d", cmd.Type, cmd.Aggregate, index)
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.applied[key]; ok {
		return true
	}
	d.applied[key] = struct{}{}
	return false
}

func (d *ApplyDispatcher) notify(result PendingResult) {
	if d == nil || d.results == nil {
		return
	}
	select {
	case d.results <- result:
	default:
	}
}
