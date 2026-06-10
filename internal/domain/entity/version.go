package entity

import "ai-coding-training/internal/domain/valueobject"

type ConfigVersion struct {
	ID          string
	GroupID     string
	VersionNo   int64
	Title       string
	Content     string
	Status      valueobject.VersionStatus
	PublishedAt int64
	CreatedAt   int64
	UpdatedAt   int64
}

type ReleaseRecord struct {
	ID            string
	GroupID       string
	VersionID     string
	TargetVersion int64
	Action        valueobject.ReleaseAction
	Status        valueobject.ReleaseStatus
	Remark        string
	CreatedAt     int64
	UpdatedAt     int64
}

type GrayRule struct {
	ID        string
	RecordID  string
	RuleType  valueobject.GrayRuleType
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
	Status        valueobject.GrayRecordStatus
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
	EventType     valueobject.ChangeEventType
	CreatedAt     int64
	Payload       string
}
