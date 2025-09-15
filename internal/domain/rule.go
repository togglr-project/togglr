package domain

import (
	"time"
)

type RuleID string

type Rule struct {
	ID            RuleID
	ProjectID     ProjectID
	FeatureID     FeatureID
	Conditions    Conditions
	FlagVariantID FlagVariantID // which variant to assign if the condition matches
	Priority      uint8
	CreatedAt     time.Time
}

func (id RuleID) String() string {
	return string(id)
}

type RuleAttribute string

const (
	AttrUserID    RuleAttribute = "user.id"
	AttrUserEmail RuleAttribute = "user.email"
	AttrCountry   RuleAttribute = "user.country"
	AttrAppVer    RuleAttribute = "app.version"
	AttrEnv       RuleAttribute = "env"
	AttrIP        RuleAttribute = "request.ip"
)

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
