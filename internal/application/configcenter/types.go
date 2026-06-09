package configcenter

import (
	cmd "ai-coding-training/internal/application/configcenter/command"
	dto "ai-coding-training/internal/application/configcenter/dto"
)

type (
	CommandType = cmd.CommandType
	Command     = cmd.Command

	AppDTO          = dto.AppDTO
	EnvironmentDTO   = dto.EnvironmentDTO
	ConfigGroupDTO   = dto.ConfigGroupDTO
	ConfigItemDTO    = dto.ConfigItemDTO
	ConfigVersionDTO = dto.ConfigVersionDTO
	ReleaseRecordDTO = dto.ReleaseRecordDTO
	GrayRuleDTO      = dto.GrayRuleDTO
	GrayRecordDTO    = dto.GrayRecordDTO
	SubscriptionDTO  = dto.SubscriptionDTO
	ChangeEventDTO   = dto.ChangeEventDTO
)