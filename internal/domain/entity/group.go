package entity

type ConfigGroup struct {
	ID            string
	AppID         string
	EnvironmentID string
	Name          string
	Description   string
	CreatedAt     int64
	UpdatedAt     int64
}
