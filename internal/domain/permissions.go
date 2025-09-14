package domain

// PermKey is a canonical permission identifier used across layers.
type PermKey string

const (
	PermProjectView      PermKey = "project.view"
	PermProjectManage    PermKey = "project.manage"
	PermProjectCreate    PermKey = "project.create"
	PermFeatureView      PermKey = "feature.view"
	PermFeatureToggle    PermKey = "feature.toggle"
	PermFeatureManage    PermKey = "feature.manage"
	PermRuleManage       PermKey = "rule.manage"
	PermAuditView        PermKey = "audit.view"
	PermMembershipManage PermKey = "membership.manage"
)
