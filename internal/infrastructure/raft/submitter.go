package raftadapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	appcmd "ai-coding-training/internal/application/command"
	"ai-coding-training/internal/application/usecase"
)

type PendingResult struct {
	Command appcmd.ConfigCommand
	Error   error
}

type Submitter struct {
	peer    PeerFactory
	adapter CommandAdapter
	results <-chan PendingResult
	leader  string
}

func NewSubmitter(peer PeerFactory, adapter CommandAdapter, results <-chan PendingResult, leader string) *Submitter {
	return &Submitter{peer: peer, adapter: adapter, results: results, leader: leader}
}

func (s *Submitter) GetState() (int, bool) {
	if s == nil || s.peer == nil {
		return 0, false
	}
	return s.peer.GetState()
}

func (s *Submitter) GetLeader() string {
	if s == nil || s.peer == nil {
		return ""
	}
	if leader := s.peer.GetLeader(); leader != "" {
		return leader
	}
	return s.leader
}

func (s *Submitter) SubmitCommand(ctx context.Context, command any) (any, error) {
	if s == nil || s.peer == nil || s.adapter == nil {
		return nil, usecase.ErrRaftLeaderUnknown
	}
	cmd, ok := command.(appcmd.ConfigCommand)
	if !ok {
		return nil, fmt.Errorf("%w: unsupported raft command type", usecase.ErrValidation)
	}
	_, isLeader := s.peer.GetState()
	if !isLeader {
		if s.GetLeader() == "" {
			return nil, usecase.ErrRaftLeaderUnknown
		}
		return nil, usecase.ErrRaftNotLeader
	}
	envelope, err := s.adapter.ToRaftCommand(cmd)
	if err != nil {
		return nil, err
	}
	_, _, leader := s.peer.Start(envelope)
	if !leader {
		return nil, usecase.ErrRaftNotLeader
	}
	if s.results == nil {
		return nil, usecase.ErrRaftMajorityFailed
	}
	return s.waitForCommit(ctx, cmd)
}

func (s *Submitter) waitForCommit(ctx context.Context, expected appcmd.ConfigCommand) (appcmd.ConfigCommand, error) {
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return appcmd.ConfigCommand{}, usecase.ErrRaftCommitTimeout
			}
			return appcmd.ConfigCommand{}, ctx.Err()
		case result, ok := <-s.results:
			if !ok {
				return appcmd.ConfigCommand{}, usecase.ErrRaftMajorityFailed
			}
			if result.Command.ID != expected.ID {
				continue
			}
			if result.Error != nil {
				return appcmd.ConfigCommand{}, result.Error
			}
			return result.Command, nil
		case <-time.After(30 * time.Second):
			return appcmd.ConfigCommand{}, usecase.ErrRaftCommitTimeout
		}
	}
}
