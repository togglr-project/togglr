package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeedbackEventID uint64

type FeedbackEventType string

const (
	FeedbackEventTypeUnknown    FeedbackEventType = "unknown"
	FeedbackEventTypeEvaluation FeedbackEventType = "evaluation"
	FeedbackEventTypeSuccess    FeedbackEventType = "success"
	FeedbackEventTypeFailure    FeedbackEventType = "failure"
	FeedbackEventTypeError      FeedbackEventType = "error"
)

type FeedbackEvent struct {
	ID            FeedbackEventID
	ProjectID     ProjectID
	EnvironmentID EnvironmentID
	FeatureID     FeatureID
	FeatureKey    string
	EnvKey        string
	VariantKey    string
	AlgorithmSlug string
	EventType     FeedbackEventType
	Reward        decimal.Decimal
	Context       map[string]any
	CreatedAt     time.Time
}

type FeedbackEventDTO struct {
	ProjectID     ProjectID         `json:"projectID"`
	EnvironmentID EnvironmentID     `json:"environmentID"`
	FeatureID     FeatureID         `json:"featureID"`
	FeatureKey    string            `json:"feature_key"`
	EnvKey        string            `json:"env_key"`
	VariantKey    string            `json:"variantKey"`
	AlgorithmSlug string            `json:"algorithmSlug"`
	EventType     FeedbackEventType `json:"eventType"`
	Reward        decimal.Decimal   `json:"reward"`
	Context       map[string]any    `json:"context"`
}
