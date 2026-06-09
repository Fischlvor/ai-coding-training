package configcenter

type ConfigItemDTO struct {
	ID          string
	GroupID     string
	Key         string
	Value       string
	Description string
	VersionNo   int64
}
