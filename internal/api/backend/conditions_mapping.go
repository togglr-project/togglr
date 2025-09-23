package apibackend

import (
	"encoding/json"

	"github.com/go-faster/jx"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// exprFromAPI converts generated RuleConditionExpression to domain BooleanExpression.
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

// exprToAPI converts domain BooleanExpression to generated RuleConditionExpression.
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
		b, _ := json.Marshal(in.Value)
		raw = b
	}
	return generatedapi.RuleCondition{
		Attribute: generatedapi.RuleAttribute(in.Attribute),
		Operator:  generatedapi.RuleOperator(in.Operator),
		Value:     raw,
	}
}
