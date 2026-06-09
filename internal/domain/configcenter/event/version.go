package configcenter

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

type ConfigVersion struct {
	ID          string
	GroupID     string
	VersionNo   int64
	Title       string
	Content     string
	Status      VersionStatus
	PublishedAt int64
	CreatedAt   int64
	UpdatedAt   int64
}

type ReleaseRecord struct {
	ID            string
	GroupID       string
	VersionID     string
	TargetVersion int64
	Action        ReleaseAction
	Status        ReleaseStatus
	Remark        string
	CreatedAt     int64
	UpdatedAt     int64
}

type GrayRule struct {
	ID        string
	RecordID  string
	RuleType  GrayRuleType
	RuleValue string
	Enabled   bool
	CreatedAt int64
	UpdatedAt int64
}

type GrayRecord struct {
	ID            string
	ReleaseID     string
	TargetVersion int64
	MatchedCount  int64
	Status        GrayRecordStatus
	CreatedAt     int64
	UpdatedAt     int64
}

type Subscription struct {
	ID            string
	ClientID      string
	AppID         string
	EnvironmentID string
	GroupID       string
	LastVersionNo int64
	CreatedAt     int64
	UpdatedAt     int64
}

type ChangeEvent struct {
	ID            string
	AppID         string
	EnvironmentID string
	GroupID       string
	VersionNo     int64
	EventType     ChangeEventType
	CreatedAt     int64
	Payload       string
}
