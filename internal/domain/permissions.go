package domain

// PermKey is a canonical permission identifier used across layers.
type PermKey string

const (
	// Project-level
	PermProjectView   PermKey = "project.view"
	PermProjectManage PermKey = "project.manage"
	PermProjectCreate PermKey = "project.create"

	// Feature-level
	PermFeatureView   PermKey = "feature.view"
	PermFeatureToggle PermKey = "feature.toggle"
	PermFeatureManage PermKey = "feature.manage"

	// Segments & Scheduling
	PermSegmentManage  PermKey = "segment.manage"
	PermScheduleManage PermKey = "schedule.manage"

	// Audit & Membership
	PermAuditView        PermKey = "audit.view"
	PermMembershipManage PermKey = "membership.manage"
)
