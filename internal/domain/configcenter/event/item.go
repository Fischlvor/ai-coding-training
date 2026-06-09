package configcenter

type ConfigItem struct {
	ID          string
	GroupID     string
	Key         string
	Value       string
	Description string
	VersionNo   int64
	IsDeleted   bool
	CreatedAt   int64
	UpdatedAt   int64
}
