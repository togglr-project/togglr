package domain

import (
	"time"
)

type FeedbackEventID uint64

type FeedbackEvent struct {
	ID          FeedbackEventID
	FeatureID   FeatureID
	AlgorithmID AlgorithmID
	VariantKey  string
	EventType   string
	Reward      float64
	Context     map[string]any
	CreatedAt   time.Time
}

type FeedbackEventDTO struct {
	FeatureID   FeatureID      `json:"featureID"`
	AlgorithmID *AlgorithmID   `json:"algorithmID"`
	VariantKey  string         `json:"variantKey"`
	EventType   string         `json:"eventType"`
	Reward      float64        `json:"reward"`
	Context     map[string]any `json:"context"`
}
