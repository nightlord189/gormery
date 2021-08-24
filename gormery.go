package gormery

import "fmt"

type ConditionElement struct {
	Operator string
	Field    string
	Value    interface{}
}

func CombineSimpleQuery(elements []ConditionElement, relation string) (string, []interface{}) {
	if len(elements) == 0 {
		return "", nil
	}
	sql := ""
	values := make([]interface{}, len(elements))
	for index, query := range elements {
		if sql == "" {
			sql += fmt.Sprintf("%s %s ?", query.Field, query.Operator)
		} else {
			sql += fmt.Sprintf(" %s %s %s ?", relation, query.Field, query.Operator)
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
