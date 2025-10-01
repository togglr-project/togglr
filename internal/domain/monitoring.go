package domain

import "time"

// ErrorReport represents a reported SDK error for a feature in an environment.
type ErrorReport struct {
	EventID       string
	ProjectID     ProjectID
	FeatureID     FeatureID
	EnvironmentID EnvironmentID
	ErrorType     string
	ErrorMessage  string
	Context       map[string]any
	CreatedAt     time.Time
}

// FeatureHealth is an aggregate health snapshot for a feature/environment.
type FeatureHealth struct {
	FeatureID     FeatureID
	EnvironmentID EnvironmentID
	Enabled       bool
	Status        string  // healthy | degraded | disabled
	ErrorRate     float64 // percentage 0..1
	LastErrorAt   time.Time
}
