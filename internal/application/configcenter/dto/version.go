package configcenter

type ConfigVersionDTO struct {
	ID          string
	GroupID     string
	VersionNo   int64
	Title       string
	Content     string
	Status      string
	PublishedAt int64
}

type ReleaseRecordDTO struct {
	ID            string
	GroupID       string
	VersionID     string
	TargetVersion int64
	Action        string
	Status        string
	Remark        string
}
