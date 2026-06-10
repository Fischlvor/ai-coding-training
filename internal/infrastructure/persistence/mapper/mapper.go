package mapper

import (
	"ai-coding-training/internal/application/dto"
	entity "ai-coding-training/internal/domain/entity"
	vo "ai-coding-training/internal/domain/valueobject"
)

type Mapper interface {
	AppToDTO(entity.App) dto.AppDTO
	AppFromDTO(dto.AppDTO) entity.App
	EnvironmentToDTO(entity.Environment) dto.EnvironmentDTO
	EnvironmentFromDTO(dto.EnvironmentDTO) entity.Environment
	GroupToDTO(entity.ConfigGroup) dto.ConfigGroupDTO
	GroupFromDTO(dto.ConfigGroupDTO) entity.ConfigGroup
	ItemToDTO(entity.ConfigItem) dto.ConfigItemDTO
	ItemFromDTO(dto.ConfigItemDTO) entity.ConfigItem
	VersionToDTO(entity.ConfigVersion) dto.ConfigVersionDTO
	VersionFromDTO(dto.ConfigVersionDTO) entity.ConfigVersion
	ReleaseToDTO(entity.ReleaseRecord) dto.ReleaseRecordDTO
	ReleaseFromDTO(dto.ReleaseRecordDTO) entity.ReleaseRecord
	GrayRuleToDTO(entity.GrayRule) dto.GrayRuleDTO
	GrayRuleFromDTO(dto.GrayRuleDTO) entity.GrayRule
	GrayRecordToDTO(entity.GrayRecord) dto.GrayRecordDTO
	GrayRecordFromDTO(dto.GrayRecordDTO) entity.GrayRecord
	SubscriptionToDTO(entity.Subscription) dto.SubscriptionDTO
	SubscriptionFromDTO(dto.SubscriptionDTO) entity.Subscription
	ChangeEventToDTO(entity.ChangeEvent) dto.ChangeEventDTO
	ChangeEventFromDTO(dto.ChangeEventDTO) entity.ChangeEvent
}

type IdentityMapper struct{}

func (IdentityMapper) AppToDTO(v entity.App) dto.AppDTO {
	return dto.AppDTO{ID: v.ID, Name: v.Name, Description: v.Description}
}
func (IdentityMapper) AppFromDTO(v dto.AppDTO) entity.App {
	return entity.App{ID: v.ID, Name: v.Name, Description: v.Description}
}
func (IdentityMapper) EnvironmentToDTO(v entity.Environment) dto.EnvironmentDTO {
	return dto.EnvironmentDTO{ID: v.ID, Name: v.Name, Description: v.Description}
}
func (IdentityMapper) EnvironmentFromDTO(v dto.EnvironmentDTO) entity.Environment {
	return entity.Environment{ID: v.ID, Name: v.Name, Description: v.Description}
}
func (IdentityMapper) GroupToDTO(v entity.ConfigGroup) dto.ConfigGroupDTO {
	return dto.ConfigGroupDTO{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, Name: v.Name, Description: v.Description}
}
func (IdentityMapper) GroupFromDTO(v dto.ConfigGroupDTO) entity.ConfigGroup {
	return entity.ConfigGroup{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, Name: v.Name, Description: v.Description}
}
func (IdentityMapper) ItemToDTO(v entity.ConfigItem) dto.ConfigItemDTO {
	return dto.ConfigItemDTO{ID: v.ID, GroupID: v.GroupID, Key: v.Key, Value: v.Value, Description: v.Description, VersionNo: v.VersionNo}
}
func (IdentityMapper) ItemFromDTO(v dto.ConfigItemDTO) entity.ConfigItem {
	return entity.ConfigItem{ID: v.ID, GroupID: v.GroupID, Key: v.Key, Value: v.Value, Description: v.Description, VersionNo: v.VersionNo}
}
func (IdentityMapper) VersionToDTO(v entity.ConfigVersion) dto.ConfigVersionDTO {
	return dto.ConfigVersionDTO{ID: v.ID, GroupID: v.GroupID, VersionNo: v.VersionNo, Title: v.Title, Content: v.Content, Status: string(v.Status), PublishedAt: v.PublishedAt}
}
func (IdentityMapper) VersionFromDTO(v dto.ConfigVersionDTO) entity.ConfigVersion {
	return entity.ConfigVersion{ID: v.ID, GroupID: v.GroupID, VersionNo: v.VersionNo, Title: v.Title, Content: v.Content, Status: vo.VersionStatus(v.Status), PublishedAt: v.PublishedAt}
}
func (IdentityMapper) ReleaseToDTO(v entity.ReleaseRecord) dto.ReleaseRecordDTO {
	return dto.ReleaseRecordDTO{ID: v.ID, GroupID: v.GroupID, VersionID: v.VersionID, TargetVersion: v.TargetVersion, Action: string(v.Action), Status: string(v.Status), Remark: v.Remark}
}
func (IdentityMapper) ReleaseFromDTO(v dto.ReleaseRecordDTO) entity.ReleaseRecord {
	return entity.ReleaseRecord{ID: v.ID, GroupID: v.GroupID, VersionID: v.VersionID, TargetVersion: v.TargetVersion, Action: vo.ReleaseAction(v.Action), Status: vo.ReleaseStatus(v.Status), Remark: v.Remark}
}
func (IdentityMapper) GrayRuleToDTO(v entity.GrayRule) dto.GrayRuleDTO {
	return dto.GrayRuleDTO{ID: v.ID, RecordID: v.RecordID, RuleType: string(v.RuleType), RuleValue: v.RuleValue, Enabled: v.Enabled}
}
func (IdentityMapper) GrayRuleFromDTO(v dto.GrayRuleDTO) entity.GrayRule {
	return entity.GrayRule{ID: v.ID, RecordID: v.RecordID, RuleType: vo.GrayRuleType(v.RuleType), RuleValue: v.RuleValue, Enabled: v.Enabled}
}
func (IdentityMapper) GrayRecordToDTO(v entity.GrayRecord) dto.GrayRecordDTO {
	return dto.GrayRecordDTO{ID: v.ID, ReleaseID: v.ReleaseID, TargetVersion: v.TargetVersion, MatchedCount: v.MatchedCount, Status: string(v.Status)}
}
func (IdentityMapper) GrayRecordFromDTO(v dto.GrayRecordDTO) entity.GrayRecord {
	return entity.GrayRecord{ID: v.ID, ReleaseID: v.ReleaseID, TargetVersion: v.TargetVersion, MatchedCount: v.MatchedCount, Status: vo.GrayRecordStatus(v.Status)}
}
func (IdentityMapper) SubscriptionToDTO(v entity.Subscription) dto.SubscriptionDTO {
	return dto.SubscriptionDTO{ID: v.ID, ClientID: v.ClientID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, GroupID: v.GroupID, LastVersionNo: v.LastVersionNo}
}
func (IdentityMapper) SubscriptionFromDTO(v dto.SubscriptionDTO) entity.Subscription {
	return entity.Subscription{ID: v.ID, ClientID: v.ClientID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, GroupID: v.GroupID, LastVersionNo: v.LastVersionNo}
}
func (IdentityMapper) ChangeEventToDTO(v entity.ChangeEvent) dto.ChangeEventDTO {
	return dto.ChangeEventDTO{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, GroupID: v.GroupID, VersionNo: v.VersionNo, EventType: string(v.EventType), Payload: v.Payload}
}
func (IdentityMapper) ChangeEventFromDTO(v dto.ChangeEventDTO) entity.ChangeEvent {
	return entity.ChangeEvent{ID: v.ID, AppID: v.AppID, EnvironmentID: v.EnvironmentID, GroupID: v.GroupID, VersionNo: v.VersionNo, EventType: vo.ChangeEventType(v.EventType), Payload: v.Payload}
}
