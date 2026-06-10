package query

type GetGroupQuery struct {
	ID string
}

type ListGroupsQuery struct {
	AppID         string
	EnvironmentID string
}
