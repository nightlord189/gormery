// Package gormery предоставляет дополнительные утилиты для построения SQL-запросов,
// работая как обертка над GORM, упрощая фильтрацию и условия в запросах.
package gormery

import (
	"fmt"
	"strings"
	"time"
)

// ConditionElement представляет элемент условия SQL запроса
type ConditionElement struct {
	Field    string      // Имя поля или полная часть SQL выражения
	Operator string      // Оператор сравнения (=, >, <, LIKE и т.д.)
	Value    interface{} // Значение для сравнения
}

// Equal создает условие равенства field = value
func Equal(field string, value interface{}) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "=",
		Value:    value,
	}
}

// NotEqual создает условие неравенства field != value
func NotEqual(field string, value interface{}) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "!=",
		Value:    value,
	}
}

// More создает условие field > value
func More(field string, value interface{}) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: ">",
		Value:    value,
	}
}

// Less создает условие field < value
func Less(field string, value interface{}) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "<",
		Value:    value,
	}
}

// MoreOrEqual создает условие field >= value
func MoreOrEqual(field string, value interface{}) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: ">=",
		Value:    value,
	}
}

// LessOrEqual создает условие field <= value
func LessOrEqual(field string, value interface{}) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "<=",
		Value:    value,
	}
}

// Like создает условие field LIKE value
func Like(field string, value string) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "LIKE",
		Value:    value,
	}
}

// In создает условие field IN (values)
func In(field string, values interface{}) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "IN",
		Value:    values,
	}
}

// IsNull создает условие field IS NULL
func IsNull(field string) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "IS NULL",
		Value:    nil,
	}
}

// IsNotNull создает условие field IS NOT NULL
func IsNotNull(field string) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "IS NOT NULL",
		Value:    nil,
	}
}

// Between создает условие field BETWEEN value1 AND value2
func Between(field string, value1, value2 interface{}) ConditionElement {
	return ConditionElement{
		Field:    field,
		Operator: "BETWEEN",
		Value:    []interface{}{value1, value2},
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
		switch cond.Operator {
		case "IS NULL", "IS NOT NULL":
			// Для операторов, которые не требуют значения
			queryParts = append(queryParts, fmt.Sprintf("%s %s", cond.Field, cond.Operator))
		case "IN":
			// Для оператора IN, который требует особого форматирования
			switch values := cond.Value.(type) {
			case []string:
				placeholders := make([]string, len(values))
				for i, v := range values {
					placeholders[i] = "?"
					params = append(params, v)
				}
				queryParts = append(queryParts, fmt.Sprintf("%s %s (%s)", cond.Field, cond.Operator, strings.Join(placeholders, ", ")))
			case []int:
				placeholders := make([]string, len(values))
				for i, v := range values {
					placeholders[i] = "?"
					params = append(params, v)
				}
				queryParts = append(queryParts, fmt.Sprintf("%s %s (%s)", cond.Field, cond.Operator, strings.Join(placeholders, ", ")))
			default:
				// Для других типов, если требуется поддержка
				queryParts = append(queryParts, fmt.Sprintf("%s %s (?)", cond.Field, cond.Operator))
				params = append(params, cond.Value)
			}
		case "BETWEEN":
			// Для оператора BETWEEN, который требует двух значений
			values, ok := cond.Value.([]interface{})
			if ok && len(values) == 2 {
				queryParts = append(queryParts, fmt.Sprintf("%s %s ? AND ?", cond.Field, cond.Operator))
				params = append(params, values[0], values[1])
			}
		default:
			// Для стандартных операторов (=, >, <, !=, LIKE, ...)
			queryParts = append(queryParts, fmt.Sprintf("%s %s ?", cond.Field, cond.Operator))
			params = append(params, cond.Value)
		}
	}

	return strings.Join(queryParts, fmt.Sprintf(" %s ", logicalOperator)), params
}

// FormatTimeValue форматирует time.Time для базы данных
func FormatTimeValue(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// BuildOrderBy создает часть ORDER BY запроса
func BuildOrderBy(fields []string, directions []string) string {
	if len(fields) == 0 {
		return ""
	}

	orderClauses := make([]string, len(fields))
	for i, field := range fields {
		direction := "ASC" // По умолчанию сортировка по возрастанию
		if i < len(directions) && (directions[i] == "DESC" || directions[i] == "desc") {
			direction = "DESC"
		}
		orderClauses[i] = fmt.Sprintf("%s %s", field, direction)
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(orderClauses, ", "))
}

// BuildLimit создает часть LIMIT запроса
func BuildLimit(limit int) string {
	if limit <= 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d", limit)
}

// BuildOffset создает часть OFFSET запроса
func BuildOffset(offset int) string {
	if offset <= 0 {
		return ""
	}
	return fmt.Sprintf("OFFSET %d", offset)
}

// BuildPagination создает части LIMIT и OFFSET для пагинации
func BuildPagination(page, pageSize int) (string, string) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	limit := BuildLimit(pageSize)
	offset := BuildOffset((page - 1) * pageSize)

	return limit, offset
}
