package gormery

import "fmt"

type ConditionElement struct {
	Operator         string
	Field            string
	Value            interface{}
	Children         []ConditionElement
	ChildrenRelation string
}

func CombineSimpleQuery(elements []ConditionElement, relation string) (string, []interface{}) {
	if len(elements) == 0 {
		return "", nil
	}
	sql := ""
	values := make([]interface{}, 0, len(elements))
	for _, query := range elements {
		if sql == "" {
			if len(query.Children) == 0 {
				sql += fmt.Sprintf("%s %s ?", query.Field, query.Operator)
				values = append(values, query.Value)
			} else {
				querySql, queryValues := CombineSimpleQuery(query.Children, query.ChildrenRelation)
				sql += fmt.Sprintf("(%s)", querySql)
				values = append(values, queryValues...)
			}
		} else {
			if len(query.Children) == 0 {
				sql += fmt.Sprintf(" %s %s %s ?", relation, query.Field, query.Operator)
				values = append(values, query.Value)
			} else {
				querySql, queryValues := CombineSimpleQuery(query.Children, query.ChildrenRelation)
				sql += fmt.Sprintf(" %s (%s)", relation, querySql)
				values = append(values, queryValues...)
			}
		}
	}
	return sql, values
}

func Complex(relation string, elems ...ConditionElement) ConditionElement {
	return ConditionElement{
		Children:         elems,
		ChildrenRelation: relation,
	}
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
