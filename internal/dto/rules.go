package dto

import (
	"encoding/json"

	"github.com/go-faster/jx"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainRuleToAPI converts domain Rule to generated API Rule.
func DomainRuleToAPI(rule domain.Rule) (generatedapi.Rule, error) {
	expr, err := exprToAPI(rule.Conditions)
	if err != nil {
		return generatedapi.Rule{}, err
	}

	var segmentID generatedapi.OptString
	if rule.SegmentID != nil {
		segmentID = generatedapi.NewOptString(rule.SegmentID.String())
	}

	return generatedapi.Rule{
		ID:            rule.ID.String(),
		FeatureID:     rule.FeatureID.String(),
		Conditions:    expr,
		SegmentID:     segmentID,
		IsCustomized:  rule.IsCustomized,
		Action:        generatedapi.RuleAction(rule.Action),
		FlagVariantID: flagVariantRef2OptString(rule.FlagVariantID),
		Priority:      int(rule.Priority),
		CreatedAt:     rule.CreatedAt,
	}, nil
}

// DomainRulesToAPI converts slice of domain Rules to slice of generated API Rules.
func DomainRulesToAPI(rules []domain.Rule) ([]generatedapi.Rule, error) {
	resp := make([]generatedapi.Rule, 0, len(rules))

	for _, rule := range rules {
		apiRule, err := DomainRuleToAPI(rule)
		if err != nil {
			return nil, err
		}

		resp = append(resp, apiRule)
	}

	return resp, nil
}

// APIRuleToDomain converts generated API Rule to domain Rule.
func APIRuleToDomain(rule generatedapi.Rule) (domain.Rule, error) {
	var conditions domain.BooleanExpression

	var err error

	if rule.Conditions.Condition.IsSet() {
		conditions, err = exprFromAPI(rule.Conditions)
		if err != nil {
			return domain.Rule{}, err
		}
	}

	var segmentID *domain.SegmentID

	if rule.SegmentID.IsSet() {
		sid := domain.SegmentID(rule.SegmentID.Value)
		segmentID = &sid
	}

	return domain.Rule{
		ID:            domain.RuleID(rule.ID),
		FeatureID:     domain.FeatureID(rule.FeatureID),
		Conditions:    conditions,
		SegmentID:     segmentID,
		IsCustomized:  rule.IsCustomized,
		Action:        domain.RuleAction(rule.Action),
		FlagVariantID: optString2FlagVariantIDRef(rule.FlagVariantID),
		Priority:      uint8(rule.Priority),
		CreatedAt:     rule.CreatedAt,
	}, nil
}

// Helper functions for condition expressions.
func exprFromAPI(in generatedapi.RuleConditionExpression) (domain.BooleanExpression, error) {
	if in.Condition.IsSet() {
		c, _ := in.Condition.Get()

		dc, err := conditionFromAPI(c)
		if err != nil {
			return domain.BooleanExpression{}, err
		}

		return domain.BooleanExpression{Condition: &dc}, nil
	}

	if in.Group.IsSet() {
		g, _ := in.Group.Get()

		dg, err := groupFromAPI(g)
		if err != nil {
			return domain.BooleanExpression{}, err
		}

		return domain.BooleanExpression{Group: &dg}, nil
	}
	// empty expression
	return domain.BooleanExpression{}, nil
}

func groupFromAPI(in generatedapi.RuleConditionGroup) (domain.ConditionGroup, error) {
	children := make([]domain.BooleanExpression, 0, len(in.Children))

	for _, ch := range in.Children {
		e, err := exprFromAPI(ch)
		if err != nil {
			return domain.ConditionGroup{}, err
		}

		children = append(children, e)
	}

	return domain.ConditionGroup{
		Operator: domain.LogicalOperator(in.Operator),
		Children: children,
	}, nil
}

func conditionFromAPI(in generatedapi.RuleCondition) (domain.Condition, error) {
	var val any
	if len(in.Value) > 0 {
		if err := json.Unmarshal(in.Value, &val); err != nil {
			return domain.Condition{}, err
		}
	}

	return domain.Condition{
		Attribute: domain.RuleAttribute(in.Attribute),
		Operator:  domain.RuleOperator(in.Operator),
		Value:     val,
	}, nil
}

func exprToAPI(in domain.BooleanExpression) (generatedapi.RuleConditionExpression, error) {
	if in.Condition != nil {
		c := conditionToAPI(*in.Condition)

		return generatedapi.RuleConditionExpression{Condition: generatedapi.NewOptRuleCondition(c)}, nil
	}

	if in.Group != nil {
		g, err := groupToAPI(*in.Group)
		if err != nil {
			return generatedapi.RuleConditionExpression{}, err
		}

		return generatedapi.RuleConditionExpression{Group: generatedapi.NewOptRuleConditionGroup(g)}, nil
	}

	return generatedapi.RuleConditionExpression{}, nil
}

func groupToAPI(in domain.ConditionGroup) (generatedapi.RuleConditionGroup, error) {
	children := make([]generatedapi.RuleConditionExpression, 0, len(in.Children))

	for _, ch := range in.Children {
		e, err := exprToAPI(ch)
		if err != nil {
			return generatedapi.RuleConditionGroup{}, err
		}

		children = append(children, e)
	}

	return generatedapi.RuleConditionGroup{
		Operator: generatedapi.LogicalOperator(in.Operator),
		Children: children,
	}, nil
}

func conditionToAPI(in domain.Condition) generatedapi.RuleCondition {
	var raw jx.Raw

	if in.Value != nil {
		b, err := json.Marshal(in.Value)
		if err != nil {
			raw = jx.Raw("{}")
		} else {
			raw = b
		}
	}

	return generatedapi.RuleCondition{
		Attribute: generatedapi.RuleAttribute(in.Attribute),
		Operator:  generatedapi.RuleOperator(in.Operator),
		Value:     raw,
	}
}
