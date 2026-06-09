package configcenter

type SubscriptionDTO struct {
	ID            string
	ClientID      string
	AppID         string
	EnvironmentID string
	GroupID       string
	LastVersionNo int64
}

type ChangeEventDTO struct {
	ID            string
	AppID         string
	EnvironmentID string
	GroupID       string
	VersionNo     int64
	EventType     string
	Payload       string
}
