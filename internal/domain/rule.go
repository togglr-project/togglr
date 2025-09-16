package domain

import (
	"time"
)

type Rule struct {
	ID            RuleID
	ProjectID     ProjectID
	FeatureID     FeatureID
	Conditions    Conditions
	Action        RuleAction
	FlagVariantID *FlagVariantID // which variant to assign if the condition matches
	Priority      uint8
	CreatedAt     time.Time
}

type RuleID string

func (id RuleID) String() string {
	return string(id)
}

type RuleAction string

const (
	RuleActionAssign  RuleAction = "assign"
	RuleActionInclude RuleAction = "include"
	RuleActionExclude RuleAction = "exclude"
)

func (action RuleAction) String() string {
	return string(action)
}

type RuleAttribute string

const (
	RuleAttributeUserID RuleAttribute = "user.id"
)

func (attr RuleAttribute) String() string {
	return string(attr)
}

type RuleOperator string

const (
	OpEq         RuleOperator = "eq"         // equals
	OpNotEq      RuleOperator = "neq"        // not equals
	OpIn         RuleOperator = "in"         // in list
	OpNotIn      RuleOperator = "not_in"     // not in list
	OpGt         RuleOperator = "gt"         // greater than
	OpGte        RuleOperator = "gte"        // greater or equal
	OpLt         RuleOperator = "lt"         // less than
	OpLte        RuleOperator = "lte"        // less or equal
	OpRegex      RuleOperator = "regex"      // regex match
	OpPercentage RuleOperator = "percentage" // percentage rollout
)

type Condition struct {
	Attribute RuleAttribute `json:"attribute"`
	Operator  RuleOperator  `json:"operator"`
	Value     any           `json:"value"`
}

type Conditions []Condition
