package configcenter

type ConfigCommandType string

const (
	ConfigCommandTypeCreateApp       ConfigCommandType = "create_app"
	ConfigCommandTypeUpdateApp       ConfigCommandType = "update_app"
	ConfigCommandTypeDeleteApp       ConfigCommandType = "delete_app"
	ConfigCommandTypeCreateGroup     ConfigCommandType = "create_group"
	ConfigCommandTypeUpdateGroup     ConfigCommandType = "update_group"
	ConfigCommandTypeDeleteGroup     ConfigCommandType = "delete_group"
	ConfigCommandTypeCreateItem      ConfigCommandType = "create_item"
	ConfigCommandTypeUpdateItem      ConfigCommandType = "update_item"
	ConfigCommandTypeDeleteItem      ConfigCommandType = "delete_item"
	ConfigCommandTypeSaveDraft       ConfigCommandType = "save_draft"
	ConfigCommandTypePublishVersion  ConfigCommandType = "publish_version"
	ConfigCommandTypeRollbackVersion ConfigCommandType = "rollback_version"
)

type ConfigCommand struct {
	ID        string
	Type      ConfigCommandType
	Aggregate string
	Payload   map[string]any
}
