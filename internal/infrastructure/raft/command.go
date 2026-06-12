package raftadapter

import (
	"encoding/json"
	"fmt"

	appcmd "ai-coding-training/internal/application/command"
	"ai-coding-training/internal/application/dto"
)

// RaftCommandEnvelope is the opaque command shape stored in the Raft log.
// It deliberately keeps the business command serialized so the Raft layer does
// not depend on or inspect configuration-center business semantics.
type RaftCommandEnvelope struct {
	Kind    string
	Payload []byte
}

// CommandAdapter converts application-layer commands into opaque Raft payloads
// and restores apply results back to the application layer.
type CommandAdapter interface {
	ToRaftCommand(appcmd.ConfigCommand) (RaftCommandEnvelope, error)
	FromApplyMsg(ApplyMsg) (appcmd.ConfigCommand, bool, error)
	ApplyMsgToEvent(ApplyMsg) (dto.ChangeEventDTO, bool, error)
}

type OpaqueCommandAdapter struct{}

func (OpaqueCommandAdapter) ToRaftCommand(cmd appcmd.ConfigCommand) (RaftCommandEnvelope, error) {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return RaftCommandEnvelope{}, fmt.Errorf("marshal raft command: %w", err)
	}
	return RaftCommandEnvelope{Kind: "config_command", Payload: payload}, nil
}

func (a OpaqueCommandAdapter) FromApplyMsg(msg ApplyMsg) (appcmd.ConfigCommand, bool, error) {
	if !msg.CommandValid {
		return appcmd.ConfigCommand{}, false, nil
	}
	switch command := msg.Command.(type) {
	case RaftCommandEnvelope:
		return decodeConfigCommandEnvelope(command)
	case appcmd.ConfigCommand:
		return command, true, nil
	default:
		return appcmd.ConfigCommand{}, false, nil
	}
}

func (a OpaqueCommandAdapter) ApplyMsgToEvent(msg ApplyMsg) (dto.ChangeEventDTO, bool, error) {
	cmd, ok, err := a.FromApplyMsg(msg)
	if err != nil || !ok {
		return dto.ChangeEventDTO{}, ok, err
	}
	return ConfigCommandToEvent(cmd, msg.CommandIndex), true, nil
}

func decodeConfigCommandEnvelope(envelope RaftCommandEnvelope) (appcmd.ConfigCommand, bool, error) {
	if envelope.Kind != "config_command" {
		return appcmd.ConfigCommand{}, false, nil
	}
	var cmd appcmd.ConfigCommand
	if err := json.Unmarshal(envelope.Payload, &cmd); err != nil {
		return appcmd.ConfigCommand{}, true, fmt.Errorf("unmarshal raft command: %w", err)
	}
	return cmd, true, nil
}

func ConfigCommandToEvent(cmd appcmd.ConfigCommand, index int) dto.ChangeEventDTO {
	payload, _ := json.Marshal(cmd.Payload)
	return dto.ChangeEventDTO{
		ID:        eventID(cmd.ID, index),
		AppID:     stringPayload(cmd.Payload, "app_id"),
		GroupID:   stringPayload(cmd.Payload, "group_id"),
		VersionNo: int64Payload(cmd.Payload, "version_no"),
		EventType: string(cmd.Type),
		Payload:   string(payload),
	}
}

func eventID(commandID string, index int) string {
	if commandID != "" {
		return commandID
	}
	return fmt.Sprintf("raft-%d", index)
}

func stringPayload(payload map[string]any, key string) string {
	if payload == nil {
		return ""
	}
	value, _ := payload[key].(string)
	return value
}

func int64Payload(payload map[string]any, key string) int64 {
	if payload == nil {
		return 0
	}
	switch value := payload[key].(type) {
	case int64:
		return value
	case int:
		return int64(value)
	case float64:
		return int64(value)
	case json.Number:
		v, _ := value.Int64()
		return v
	default:
		return 0
	}
}
