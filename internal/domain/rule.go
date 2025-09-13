package domain

import (
	"encoding/json"
	"time"
)

type RuleID string

type Rule struct {
	ID             RuleID
	FeatureID      FeatureID
	Condition      json.RawMessage
	Variant        string
	RolloutPercent uint8
	Priority       uint8
	CreatedAt      time.Time
}

func (id RuleID) String() string {
	return string(id)
}
