package query

type GetItemQuery struct {
	ID string
}

type ListItemsQuery struct {
	GroupID        string
	IncludeDeleted bool
}
