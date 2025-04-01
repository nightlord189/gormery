// Package gormery предоставляет дополнительные утилиты для построения SQL-запросов,
// работая как обертка над GORM, упрощая фильтрацию и условия в запросах.
package gormery

import (
	"fmt"
	"strings"
)

// ConditionElement представляет интерфейс для элементов условия SQL запроса
type ConditionElement interface {
	ToSQL() (string, []interface{})
}

// SimpleCondition представляет простое условие SQL запроса
type SimpleCondition struct {
	Field    string      // Имя поля или полная часть SQL выражения
	Operator string      // Оператор сравнения (=, >, <, LIKE и т.д.)
	Value    interface{} // Значение для сравнения
}

// ToSQL конвертирует простое условие в строку SQL и параметры
func (c SimpleCondition) ToSQL() (string, []interface{}) {
	switch c.Operator {
	case "IS NULL", "IS NOT NULL":
		// Для операторов, которые не требуют значения
		return fmt.Sprintf("%s %s", c.Field, c.Operator), nil
	case "IN":
		// Для оператора IN, который требует особого форматирования
		switch values := c.Value.(type) {
		case []string:
			placeholders := make([]string, len(values))
			params := make([]interface{}, len(values))
			for i, v := range values {
				placeholders[i] = "?"
				params[i] = v
			}
			return fmt.Sprintf("%s %s (%s)", c.Field, c.Operator, strings.Join(placeholders, ", ")), params
		case []int:
			placeholders := make([]string, len(values))
			params := make([]interface{}, len(values))
			for i, v := range values {
				placeholders[i] = "?"
				params[i] = v
			}
			return fmt.Sprintf("%s %s (%s)", c.Field, c.Operator, strings.Join(placeholders, ", ")), params
		default:
			// Для других типов, если требуется поддержка
			return fmt.Sprintf("%s %s (?)", c.Field, c.Operator), []interface{}{c.Value}
		}
	case "BETWEEN":
		// Для оператора BETWEEN, который требует двух значений
		values, ok := c.Value.([]interface{})
		if ok && len(values) == 2 {
			return fmt.Sprintf("%s %s ? AND ?", c.Field, c.Operator), values
		}
		return "", nil
	default:
		// Для стандартных операторов (=, >, <, !=, LIKE, ...)
		return fmt.Sprintf("%s %s ?", c.Field, c.Operator), []interface{}{c.Value}
	}
}

// ComplexCondition представляет сложное условие с вложенными элементами
type ComplexCondition struct {
	LogicalOperator string             // Логический оператор (AND, OR)
	Conditions      []ConditionElement // Список условий
}

// ToSQL конвертирует сложное условие в строку SQL и параметры
func (c ComplexCondition) ToSQL() (string, []interface{}) {
	if len(c.Conditions) == 0 {
		return "", nil
	}

	var queryParts []string
	var params []interface{}

	for _, cond := range c.Conditions {
		sql, elems := cond.ToSQL()
		if sql != "" {
			queryParts = append(queryParts, sql)
			if elems != nil {
				params = append(params, elems...)
			}
		}
	}

	if len(queryParts) == 0 {
		return "", nil
	}

	return fmt.Sprintf("(%s)", strings.Join(queryParts, fmt.Sprintf(" %s ", c.LogicalOperator))), params
}

// Equal создает условие равенства field = value
func Equal(field string, value interface{}) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "=",
		Value:    value,
	}
}

// NotEqual создает условие неравенства field != value
func NotEqual(field string, value interface{}) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "!=",
		Value:    value,
	}
}

// More создает условие field > value
func More(field string, value interface{}) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: ">",
		Value:    value,
	}
}

// Less создает условие field < value
func Less(field string, value interface{}) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "<",
		Value:    value,
	}
}

// MoreOrEqual создает условие field >= value
func MoreOrEqual(field string, value interface{}) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: ">=",
		Value:    value,
	}
}

// LessOrEqual создает условие field <= value
func LessOrEqual(field string, value interface{}) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "<=",
		Value:    value,
	}
}

// Like создает условие field LIKE value
func Like(field string, value string) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "LIKE",
		Value:    value,
	}
}

// In создает условие field IN (values)
func In(field string, values interface{}) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "IN",
		Value:    values,
	}
}

// IsNull создает условие field IS NULL
func IsNull(field string) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "IS NULL",
		Value:    nil,
	}
}

// IsNotNull создает условие field IS NOT NULL
func IsNotNull(field string) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "IS NOT NULL",
		Value:    nil,
	}
}

// Between создает условие field BETWEEN value1 AND value2
func Between(field string, value1, value2 interface{}) ConditionElement {
	return SimpleCondition{
		Field:    field,
		Operator: "BETWEEN",
		Value:    []interface{}{value1, value2},
	}
}

// Complex создает сложное условие с несколькими ConditionElement, объединенными через логический оператор
func Complex(logicalOperator string, conditions ...ConditionElement) ConditionElement {
	return ComplexCondition{
		LogicalOperator: logicalOperator,
		Conditions:      conditions,
	}
}

// CombineSimpleQuery объединяет условия в SQL запрос с выбранным логическим оператором
// Возвращает SQL строку и срез параметров для подстановки
func CombineSimpleQuery(conditions []ConditionElement, logicalOperator string) (string, []interface{}) {
	if len(conditions) == 0 {
		return "", nil
	}

	var queryParts []string
	var params []interface{}

	for _, cond := range conditions {
		sql, elems := cond.ToSQL()
		if sql != "" {
			queryParts = append(queryParts, sql)
			if elems != nil {
				params = append(params, elems...)
			}
		}
	}

	if len(queryParts) == 0 {
		return "", nil
	}

	return strings.Join(queryParts, fmt.Sprintf(" %s ", logicalOperator)), params
}
