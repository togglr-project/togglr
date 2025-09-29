package domain

import (
	"time"
)

type Rule struct {
	ID            RuleID            `db:"id" pk:"true"`
	ProjectID     ProjectID         `db:"project_id"`
	FeatureID     FeatureID         `db:"feature_id"`
	EnvironmentID EnvironmentID     `db:"environment_id"`
	Conditions    BooleanExpression `db:"condition" editable:"true"`
	SegmentID     *SegmentID        `db:"segment_id" editable:"true"`
	IsCustomized  bool              `db:"is_customized" editable:"true"`
	Action        RuleAction        `db:"action" editable:"true"`
	FlagVariantID *FlagVariantID    `db:"flag_variant_id" editable:"true"` // which variant to assign if the condition matches
	Priority      uint8             `db:"priority" editable:"true"`
	CreatedAt     time.Time         `db:"created_at"`
	// TODO: updatedAt
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
