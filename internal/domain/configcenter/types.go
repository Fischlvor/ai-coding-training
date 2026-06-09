package configcenter

import ev "ai-coding-training/internal/domain/configcenter/event"

type (
	App           = ev.App
	Environment   = ev.Environment
	ConfigGroup   = ev.ConfigGroup
	ConfigItem    = ev.ConfigItem
	ConfigVersion = ev.ConfigVersion
	ReleaseRecord = ev.ReleaseRecord
	GrayRule      = ev.GrayRule
	GrayRecord    = ev.GrayRecord
	Subscription  = ev.Subscription
	ChangeEvent   = ev.ChangeEvent

	VersionStatus    = ev.VersionStatus
	ReleaseAction    = ev.ReleaseAction
	ReleaseStatus    = ev.ReleaseStatus
	GrayRuleType     = ev.GrayRuleType
	GrayRecordStatus = ev.GrayRecordStatus
	ChangeEventType  = ev.ChangeEventType
)
