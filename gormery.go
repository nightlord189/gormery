package gormery

import "fmt"

type Relation int

const (
	And Relation = iota + 1
	Or
)

type Operator int

const (
	Eq Operator = iota + 1
	NotEq
	Like
	NotLike
	More
	Less
	MoreOrEq
	LessOrEq
)

type ConditionElement struct {
	Oper  Operator
	Field string
	Value interface{}
}

func CombineSimpleQuery(elements []ConditionElement, relation Relation) (string, []interface{}) {
	if len(elements) == 0 {
		return "", nil
	}
	sql := ""
	values := make([]interface{}, len(elements))
	conjunction := ""
	switch relation {
	case And:
		conjunction = "AND"
	case Or:
		conjunction = "OR"
	}
	for index, query := range elements {
		if sql == "" {
			sql += fmt.Sprintf("%s %s ?", query.Field, stringifyOperator(query.Oper))
		} else {
			sql += fmt.Sprintf(" %s %s %s ?", conjunction, query.Field, stringifyOperator(query.Oper))
		}
		values[index] = query.Value
	}
	return sql, values
}

func stringifyOperator(operator Operator) string {
	switch operator {
	case Eq:
		return "="
	case NotEq:
		return "<>"
	case Like:
		return "LIKE"
	case NotLike:
		return "NOT LIKE"
	case More:
		return ">"
	case Less:
		return "<"
	case MoreOrEq:
		return ">="
	case LessOrEq:
		return "<="
	}
	return "OPERATOR_ERROR"
}

func stringifyConjunction(relation Relation) string {
	switch relation {
	case And:
		return "AND"
	case Or:
		return "OR"
	}
	return "CONJUNCTION_ERROR"
}
