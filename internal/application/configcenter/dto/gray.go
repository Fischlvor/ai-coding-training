package configcenter

type GrayRuleDTO struct {
	ID        string
	RecordID  string
	RuleType  string
	RuleValue string
	Enabled   bool
}

type GrayRecordDTO struct {
	ID            string
	ReleaseID     string
	TargetVersion int64
	MatchedCount  int64
	Status        string
}
