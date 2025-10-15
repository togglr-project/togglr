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
	FeatureID     FeatureID
	EnvironmentID EnvironmentID
	VariantKey    string
	EventType     FeedbackEventType
	Reward        decimal.Decimal
	Context       map[string]any
	CreatedAt     time.Time
}

type FeedbackEventDTO struct {
	FeatureID     FeatureID         `json:"featureID"`
	EnvironmentID EnvironmentID     `json:"environmentID"`
	VariantKey    string            `json:"variantKey"`
	EventType     FeedbackEventType `json:"eventType"`
	Reward        decimal.Decimal   `json:"reward"`
	Context       map[string]any    `json:"context"`
}
