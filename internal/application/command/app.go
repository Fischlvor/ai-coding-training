package configcenter

type CommandType string

const (
	CommandTypeUpsertApp CommandType = "upsert_app"
)

type Command struct {
	ID        string
	Type      CommandType
	Aggregate string
	VersionNo int64
	Payload   []byte
	CreatedAt int64
}
