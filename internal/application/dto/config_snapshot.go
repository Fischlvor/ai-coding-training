package dto

type ConfigSnapshotDTO struct {
	AppID           string
	EnvironmentID   string
	GroupID         string
	VersionNo       int64
	Title           string
	Content         string
	GrayMatched     bool
	GrayMatchedRule string
	UpdatedAt       int64
}
