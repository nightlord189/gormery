package gormery

import "fmt"

type Relation int

const (
	andRelation Relation = iota + 1
	orRelation
)

type ConditionElement struct {
	Operator string
	Field    string
	Value    interface{}
}

func CombineSimpleQuery(elements []ConditionElement, relation Relation) (string, []interface{}) {
	if len(elements) == 0 {
		return "", nil
	}
	sql := ""
	values := make([]interface{}, len(elements))
	conjunction := stringifyConjunction(relation)
	for index, query := range elements {
		if sql == "" {
			sql += fmt.Sprintf("%s %s ?", query.Field, query.Operator)
		} else {
			sql += fmt.Sprintf(" %s %s %s ?", conjunction, query.Field, query.Operator)
		}
		values[index] = query.Value
	}
	return sql, values
}

func Equal(field string, value interface{}) ConditionElement {
	return NewConditionElement("=", field, value)
}

func NotEqual(field string, value interface{}) ConditionElement {
	return NewConditionElement("<>", field, value)
}

func Like(field string, value interface{}) ConditionElement {
	return NewConditionElement("LIKE", field, value)
}

func NotLike(field string, value interface{}) ConditionElement {
	return NewConditionElement("NOT LIKE", field, value)
}

func More(field string, value interface{}) ConditionElement {
	return NewConditionElement(">", field, value)
}

func MoreOrEqual(field string, value interface{}) ConditionElement {
	return NewConditionElement(">=", field, value)
}

func Less(field string, value interface{}) ConditionElement {
	return NewConditionElement("<", field, value)
}

func LessOrEqual(field string, value interface{}) ConditionElement {
	return NewConditionElement("<=", field, value)
}

func NewConditionElement(operator, field string, value interface{}) ConditionElement {
	return ConditionElement{
		Operator: operator,
		Field:    field,
		Value:    value,
	}
}

func stringifyConjunction(relation Relation) string {
	switch relation {
	case andRelation:
		return "AND"
	case orRelation:
		return "OR"
	}
	return "CONJUNCTION_ERROR"
}
