package valueobject

type VersionStatus string

type ReleaseAction string

type ReleaseStatus string

type GrayRuleType string

type GrayRecordStatus string

type ChangeEventType string

const (
	VersionStatusDraft     VersionStatus = "draft"
	VersionStatusPublished VersionStatus = "published"
	VersionStatusArchived  VersionStatus = "archived"
)

const (
	ReleaseActionPublish  ReleaseAction = "publish"
	ReleaseActionRollback ReleaseAction = "rollback"
)

const (
	ReleaseStatusPending ReleaseStatus = "pending"
	ReleaseStatusSuccess ReleaseStatus = "success"
	ReleaseStatusFailed  ReleaseStatus = "failed"
)

const (
	GrayRuleTypePercent   GrayRuleType = "percent"
	GrayRuleTypeTag       GrayRuleType = "tag"
	GrayRuleTypeWhitelist GrayRuleType = "whitelist"
)

const (
	GrayRecordStatusPending GrayRecordStatus = "pending"
	GrayRecordStatusApplied GrayRecordStatus = "applied"
	GrayRecordStatusFailed  GrayRecordStatus = "failed"
)

const (
	ChangeEventTypePublished ChangeEventType = "published"
	ChangeEventTypeRollback  ChangeEventType = "rollback"
	ChangeEventTypeGrayPush  ChangeEventType = "gray_push"
	ChangeEventTypeSync      ChangeEventType = "sync"
)
