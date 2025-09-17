package domain

type LogicalOperator string

const (
	LogicalOpAND    LogicalOperator = "and"
	LogicalOpOR     LogicalOperator = "or"
	LogicalOpANDNot LogicalOperator = "and_not"
)

type BooleanExpression struct {
	Condition *Condition      `json:"condition,omitempty"`
	Group     *ConditionGroup `json:"group,omitempty"`
}

type ConditionGroup struct {
	Operator LogicalOperator     `json:"operator"`
	Children []BooleanExpression `json:"children"`
}

type Condition struct {
	Attribute RuleAttribute `json:"attribute"`
	Operator  RuleOperator  `json:"operator"`
	Value     any           `json:"value"`
}
