package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type FeedbackEventID uint64

type FeedbackEvent struct {
	ID            FeedbackEventID
	FeatureID     FeatureID
	AlgorithmSlug string
	VariantKey    string
	EventType     string
	Reward        decimal.Decimal
	Context       map[string]any
	CreatedAt     time.Time
}

type FeedbackEventDTO struct {
	FeatureID     FeatureID       `json:"featureID"`
	AlgorithmSlug *string         `json:"algorithmSlug"`
	VariantKey    string          `json:"variantKey"`
	EventType     string          `json:"eventType"`
	Reward        decimal.Decimal `json:"reward"`
	Context       map[string]any  `json:"context"`
}
