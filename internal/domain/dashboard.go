package domain

import "time"

// Dashboard domain models reflecting DB views.

type HealthStatus string

const (
	HealthGreen  HealthStatus = "green"
	HealthYellow HealthStatus = "yellow"
	HealthRed    HealthStatus = "red"
)

type ProjectHealth struct {
	ProjectID                  string
	ProjectName                string
	EnvironmentID              string
	EnvironmentKey             string
	TotalFeatures              uint
	EnabledFeatures            uint
	DisabledFeatures           uint
	AutoDisableManagedFeatures uint
	UncategorizedFeatures      uint
	GuardedFeatures            uint
	PendingFeatures            uint
	PendingGuardedFeatures     uint
	HealthStatus               HealthStatus
}

type CategoryHealth struct {
	ProjectID                  string
	ProjectName                string
	EnvironmentID              string
	EnvironmentKey             string
	CategoryID                 string
	CategoryName               string
	CategorySlug               string
	TotalFeatures              uint
	EnabledFeatures            uint
	DisabledFeatures           uint
	PendingFeatures            uint
	GuardedFeatures            uint
	AutoDisableManagedFeatures uint
	PendingGuardedFeatures     uint
	HealthStatus               HealthStatus
}

type RecentChange struct {
	Entity   string
	EntityID string
	Action   string
}

type RecentActivity struct {
	ProjectID      string
	EnvironmentID  string
	EnvironmentKey string
	ProjectName    string
	RequestID      string
	Actor          string
	CreatedAt      time.Time
	Status         string
	Changes        []RecentChange
}

type RiskyFeature struct {
	ProjectID      string
	ProjectName    string
	EnvironmentID  string
	EnvironmentKey string
	FeatureID      string
	FeatureName    string
	Enabled        bool
	HasPending     bool
	RiskyTags      string
}

type PendingSummary struct {
	ProjectID             string
	ProjectName           string
	EnvironmentID         string
	EnvironmentKey        string
	TotalPending          uint
	PendingFeatureChanges uint
	PendingGuardedChanges uint
	OldestRequestAt       *time.Time
}

// DashboardOverview aggregates all parts. FeatureActivity parts are optional for now.
type DashboardOverview struct {
	Projects       []ProjectHealth
	Categories     []CategoryHealth
	RecentActivity []RecentActivity
	RiskyFeatures  []RiskyFeature
	PendingSummary []PendingSummary
	Upcoming       []FeatureUpcoming // not from view yet
	Recent         []FeatureRecent   // not from view yet
}

type FeatureUpcoming struct {
	FeatureID   string
	FeatureName string
	NextState   string
	At          time.Time
}

type FeatureRecent struct {
	FeatureID   string
	FeatureName string
	Action      string
	Actor       string
	At          time.Time
}
