package mapper

import (
	app "ai-coding-training/internal/application/configcenter"
	domain "ai-coding-training/internal/domain/configcenter"
)

type Mapper interface {
	AppToDTO(domain.App) app.AppDTO
	AppFromDTO(app.AppDTO) domain.App
	EnvironmentToDTO(domain.Environment) app.EnvironmentDTO
	EnvironmentFromDTO(app.EnvironmentDTO) domain.Environment
	GroupToDTO(domain.ConfigGroup) app.ConfigGroupDTO
	GroupFromDTO(app.ConfigGroupDTO) domain.ConfigGroup
	ItemToDTO(domain.ConfigItem) app.ConfigItemDTO
	ItemFromDTO(app.ConfigItemDTO) domain.ConfigItem
	VersionToDTO(domain.ConfigVersion) app.ConfigVersionDTO
	VersionFromDTO(app.ConfigVersionDTO) domain.ConfigVersion
	ReleaseToDTO(domain.ReleaseRecord) app.ReleaseRecordDTO
	ReleaseFromDTO(app.ReleaseRecordDTO) domain.ReleaseRecord
	GrayRuleToDTO(domain.GrayRule) app.GrayRuleDTO
	GrayRuleFromDTO(app.GrayRuleDTO) domain.GrayRule
	GrayRecordToDTO(domain.GrayRecord) app.GrayRecordDTO
	GrayRecordFromDTO(app.GrayRecordDTO) domain.GrayRecord
	SubscriptionToDTO(domain.Subscription) app.SubscriptionDTO
	SubscriptionFromDTO(app.SubscriptionDTO) domain.Subscription
	ChangeEventToDTO(domain.ChangeEvent) app.ChangeEventDTO
	ChangeEventFromDTO(app.ChangeEventDTO) domain.ChangeEvent
}

type IdentityMapper struct{}

func (IdentityMapper) AppToDTO(v domain.App) app.AppDTO { return app.AppDTO{ID: v.ID, Name: v.Name, Description: v.Description} }
func (IdentityMapper) AppFromDTO(v app.AppDTO) domain.App { return domain.App{ID: v.ID, Name: v.Name, Description: v.Description} }
func (IdentityMapper) EnvironmentToDTO(v domain.Environment) app.EnvironmentDTO { return app.EnvironmentDTO{ID: v.ID, Name: v.Name, Description: v.Description} }
func (IdentityMapper) EnvironmentFromDTO(v app.EnvironmentDTO) domain.Environment { return domain.Environment{ID: v.ID, Name: v.Name, Description: v.Description} }
func (IdentityMapper) GroupToDTO(v domain.ConfigGroup) app.ConfigGroupDTO { return app.ConfigGroupDTO{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, Name: v.Name, Description: v.Description} }
func (IdentityMapper) GroupFromDTO(v app.ConfigGroupDTO) domain.ConfigGroup { return domain.ConfigGroup{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, Name: v.Name, Description: v.Description} }
func (IdentityMapper) ItemToDTO(v domain.ConfigItem) app.ConfigItemDTO { return app.ConfigItemDTO{ID: v.ID, GroupID: v.GroupID, Key: v.Key, Value: v.Value, Description: v.Description, VersionNo: v.VersionNo} }
func (IdentityMapper) ItemFromDTO(v app.ConfigItemDTO) domain.ConfigItem { return domain.ConfigItem{ID: v.ID, GroupID: v.GroupID, Key: v.Key, Value: v.Value, Description: v.Description, VersionNo: v.VersionNo} }
func (IdentityMapper) VersionToDTO(v domain.ConfigVersion) app.ConfigVersionDTO { return app.ConfigVersionDTO{ID: v.ID, GroupID: v.GroupID, VersionNo: v.VersionNo, Title: v.Title, Content: v.Content, Status: string(v.Status), PublishedAt: v.PublishedAt} }
func (IdentityMapper) VersionFromDTO(v app.ConfigVersionDTO) domain.ConfigVersion { return domain.ConfigVersion{ID: v.ID, GroupID: v.GroupID, VersionNo: v.VersionNo, Title: v.Title, Content: v.Content, Status: domain.VersionStatus(v.Status), PublishedAt: v.PublishedAt} }
func (IdentityMapper) ReleaseToDTO(v domain.ReleaseRecord) app.ReleaseRecordDTO { return app.ReleaseRecordDTO{ID: v.ID, GroupID: v.GroupID, VersionID: v.VersionID, TargetVersion: v.TargetVersion, Action: string(v.Action), Status: string(v.Status), Remark: v.Remark} }
func (IdentityMapper) ReleaseFromDTO(v app.ReleaseRecordDTO) domain.ReleaseRecord { return domain.ReleaseRecord{ID: v.ID, GroupID: v.GroupID, VersionID: v.VersionID, TargetVersion: v.TargetVersion, Action: domain.ReleaseAction(v.Action), Status: domain.ReleaseStatus(v.Status), Remark: v.Remark} }
func (IdentityMapper) GrayRuleToDTO(v domain.GrayRule) app.GrayRuleDTO { return app.GrayRuleDTO{ID: v.ID, RecordID: v.RecordID, RuleType: string(v.RuleType), RuleValue: v.RuleValue, Enabled: v.Enabled} }
func (IdentityMapper) GrayRuleFromDTO(v app.GrayRuleDTO) domain.GrayRule { return domain.GrayRule{ID: v.ID, RecordID: v.RecordID, RuleType: domain.GrayRuleType(v.RuleType), RuleValue: v.RuleValue, Enabled: v.Enabled} }
func (IdentityMapper) GrayRecordToDTO(v domain.GrayRecord) app.GrayRecordDTO { return app.GrayRecordDTO{ID: v.ID, ReleaseID: v.ReleaseID, TargetVersion: v.TargetVersion, MatchedCount: v.MatchedCount, Status: string(v.Status)} }
func (IdentityMapper) GrayRecordFromDTO(v app.GrayRecordDTO) domain.GrayRecord { return domain.GrayRecord{ID: v.ID, ReleaseID: v.ReleaseID, TargetVersion: v.TargetVersion, MatchedCount: v.MatchedCount, Status: domain.GrayRecordStatus(v.Status)} }
func (IdentityMapper) SubscriptionToDTO(v domain.Subscription) app.SubscriptionDTO { return app.SubscriptionDTO{ID: v.ID, ClientID: v.ClientID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, GroupID: v.GroupID, LastVersionNo: v.LastVersionNo} }
func (IdentityMapper) SubscriptionFromDTO(v app.SubscriptionDTO) domain.Subscription { return domain.Subscription{ID: v.ID, ClientID: v.ClientID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, GroupID: v.GroupID, LastVersionNo: v.LastVersionNo} }
func (IdentityMapper) ChangeEventToDTO(v domain.ChangeEvent) app.ChangeEventDTO { return app.ChangeEventDTO{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, GroupID: v.GroupID, VersionNo: v.VersionNo, EventType: string(v.EventType), Payload: v.Payload} }
func (IdentityMapper) ChangeEventFromDTO(v app.ChangeEventDTO) domain.ChangeEvent { return domain.ChangeEvent{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, GroupID: v.GroupID, VersionNo: v.VersionNo, EventType: domain.ChangeEventType(v.EventType), Payload: v.Payload} }
