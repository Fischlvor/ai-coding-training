package usecase

import (
	"context"
	"fmt"

	appcmd "ai-coding-training/internal/application/command"
	"ai-coding-training/internal/application/dto"
)

// ApplyCommittedCommand applies a command after Raft has committed it. This is
// the single state-machine entry for configuration metadata, version and release
// updates; callers must only invoke it for committed commands.
func (s Service) ApplyCommittedCommand(ctx context.Context, cmd appcmd.ConfigCommand) (any, error) {
	switch cmd.Type {
	case appcmd.ConfigCommandTypeCreateApp:
		return s.CreateApp(ctx, dto.AppDTO{ID: cmd.ID, Name: stringPayload(cmd, "name"), Description: stringPayload(cmd, "description")})
	case appcmd.ConfigCommandTypeUpdateApp:
		return s.UpdateApp(ctx, dto.AppDTO{ID: cmd.ID, Name: stringPayload(cmd, "name"), Description: stringPayload(cmd, "description")})
	case appcmd.ConfigCommandTypeDeleteApp:
		return nil, s.DeleteApp(ctx, cmd.ID)
	case appcmd.ConfigCommandTypeCreateGroup:
		return s.CreateGroup(ctx, dto.ConfigGroupDTO{ID: cmd.ID, AppID: stringPayload(cmd, "app_id"), EnvironmentID: stringPayload(cmd, "environment_id"), Name: stringPayload(cmd, "name"), Description: stringPayload(cmd, "description")})
	case appcmd.ConfigCommandTypeUpdateGroup:
		return s.UpdateGroup(ctx, dto.ConfigGroupDTO{ID: cmd.ID, AppID: stringPayload(cmd, "app_id"), EnvironmentID: stringPayload(cmd, "environment_id"), Name: stringPayload(cmd, "name"), Description: stringPayload(cmd, "description")})
	case appcmd.ConfigCommandTypeDeleteGroup:
		return nil, s.DeleteGroup(ctx, cmd.ID)
	case appcmd.ConfigCommandTypeCreateItem:
		return s.CreateItem(ctx, dto.ConfigItemDTO{ID: cmd.ID, GroupID: stringPayload(cmd, "group_id"), Key: stringPayload(cmd, "key"), Value: stringPayload(cmd, "value"), Description: stringPayload(cmd, "description")})
	case appcmd.ConfigCommandTypeUpdateItem:
		return s.UpdateItem(ctx, dto.ConfigItemDTO{ID: cmd.ID, GroupID: stringPayload(cmd, "group_id"), Key: stringPayload(cmd, "key"), Value: stringPayload(cmd, "value"), Description: stringPayload(cmd, "description"), VersionNo: int64Payload(cmd, "version_no")})
	case appcmd.ConfigCommandTypeDeleteItem:
		return nil, s.DeleteItem(ctx, cmd.ID)
	case appcmd.ConfigCommandTypeSaveDraft:
		return s.SaveDraftVersion(ctx, dto.ConfigVersionDTO{ID: cmd.ID, GroupID: stringPayload(cmd, "group_id"), VersionNo: int64Payload(cmd, "version_no"), Title: stringPayload(cmd, "title"), Content: stringPayload(cmd, "content")})
	case appcmd.ConfigCommandTypePublishVersion:
		return s.PublishVersion(ctx, cmd.ID, int64Payload(cmd, "target_version"), stringPayload(cmd, "remark"))
	case appcmd.ConfigCommandTypeRollbackVersion:
		return s.RollbackVersion(ctx, cmd.ID, int64Payload(cmd, "target_version"), stringPayload(cmd, "remark"))
	default:
		return nil, fmt.Errorf("%w: unsupported config command type %q", ErrValidation, cmd.Type)
	}
}

func stringPayload(cmd appcmd.ConfigCommand, key string) string {
	if cmd.Payload == nil {
		return ""
	}
	value, ok := cmd.Payload[key]
	if !ok || value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return fmt.Sprint(value)
}

func int64Payload(cmd appcmd.ConfigCommand, key string) int64 {
	if cmd.Payload == nil {
		return 0
	}
	switch value := cmd.Payload[key].(type) {
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	default:
		return 0
	}
}
