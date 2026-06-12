package usecase

import (
	"context"
	"fmt"

	appcmd "ai-coding-training/internal/application/command"
	"ai-coding-training/internal/application/dto"
)

func (s Service) SubmitConfigCommand(ctx context.Context, cmd appcmd.ConfigCommand) (appcmd.ConfigCommand, error) {
	if s.Raft == nil {
		return appcmd.ConfigCommand{}, fmt.Errorf("%w: raft submitter is not configured", ErrValidation)
	}
	if s.Recovery.IsRecovering() {
		return appcmd.ConfigCommand{}, ErrRaftRecovering
	}
	if state, ok := s.Raft.(RaftStateProvider); ok {
		_, isLeader := state.GetState()
		if !isLeader {
			if state.GetLeader() == "" {
				return appcmd.ConfigCommand{}, ErrRaftLeaderUnknown
			}
			return appcmd.ConfigCommand{}, ErrRaftNotLeader
		}
	}
	result, err := s.Raft.SubmitCommand(ctx, cmd)
	if err != nil {
		return appcmd.ConfigCommand{}, normalizeRaftWriteError(err)
	}
	submitted, ok := result.(appcmd.ConfigCommand)
	if !ok {
		return appcmd.ConfigCommand{}, fmt.Errorf("%w: unexpected raft command result", ErrValidation)
	}
	return submitted, nil
}

func normalizeRaftWriteError(err error) error {
	if err == nil {
		return nil
	}
	if err == context.DeadlineExceeded || err == context.Canceled {
		return fmt.Errorf("%w: %v", ErrRaftCommitTimeout, err)
	}
	return err
}

func (s Service) CreateAppWithRaft(ctx context.Context, in dto.AppDTO) (dto.AppDTO, error) {
	cmd := appcmd.ConfigCommand{
		ID:        in.ID,
		Type:      appcmd.ConfigCommandTypeCreateApp,
		Aggregate: "app",
		Payload: map[string]any{
			"name":        in.Name,
			"description": in.Description,
		},
	}
	if _, err := s.SubmitConfigCommand(ctx, cmd); err != nil {
		return dto.AppDTO{}, err
	}
	return s.CreateApp(ctx, in)
}
