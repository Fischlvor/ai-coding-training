package usecase

import "errors"

var (
	ErrValidation         = errors.New("validation failed")
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrRaftNotLeader      = errors.New("raft node is not leader")
	ErrRaftLeaderUnknown  = errors.New("raft leader is unavailable")
	ErrRaftMajorityFailed = errors.New("raft majority is not formed")
	ErrRaftCommitTimeout  = errors.New("raft command commit timeout")
	ErrRaftRecovering     = errors.New("raft node is recovering")
)
