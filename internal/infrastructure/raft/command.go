package raftadapter

import appcmd "ai-coding-training/internal/application/command"

// CommandAdapter converts application-layer commands into opaque Raft payloads
// and restores apply results back to the application layer.
type CommandAdapter interface {
	ToRaftCommand(appcmd.Command) interface{}
	FromApplyMsg(ApplyMsg) (appcmd.Command, bool)
}

type OpaqueCommandAdapter struct{}

func (OpaqueCommandAdapter) ToRaftCommand(cmd appcmd.Command) interface{} {
	return cmd
}

func (OpaqueCommandAdapter) FromApplyMsg(msg ApplyMsg) (appcmd.Command, bool) {
	cmd, ok := msg.Command.(appcmd.Command)
	return cmd, ok
}
