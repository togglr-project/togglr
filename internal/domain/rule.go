package domain

import (
	"encoding/json"
	"time"
)

type RuleID string

type Rule struct {
	ID            RuleID
	FeatureID     FeatureID
	Condition     json.RawMessage
	FlagVariantID FlagVariantID // which variant to assign if the condition matches
	Priority      uint8
	CreatedAt     time.Time
}

func (id RuleID) String() string {
	return string(id)
}
