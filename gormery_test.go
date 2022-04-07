package gormery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombineSimpleQuery(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		t.Parallel()
		elems := make([]ConditionElement, 0)
		sql, values := CombineSimpleQuery(elems, "AND")
		assert.Empty(t, sql)
		assert.Equal(t, 0, len(values))
	})

	t.Run("Query", func(t *testing.T) {
		t.Parallel()
		elems := []ConditionElement{
			Equal("id", 1),
			Equal("name", "John"),
			Like("position", "%manager%"),
			MoreOrEqual("salary", 20000),
			Complex("OR", Equal("doc_number", "89013"), Equal("region", "ATLANTA")),
		}
		sql, values := CombineSimpleQuery(elems, "AND")
		assert.Equal(t, "id = ? AND name = ? AND position LIKE ? AND salary >= ? AND (doc_number = ? OR region = ?)", sql)
		assert.Equal(t, len(elems)+1, len(values))
	})

	t.Run("Query with complex at start", func(t *testing.T) {
		t.Parallel()
		elems := []ConditionElement{
			Complex("OR", Equal("doc_number", "89013"), Equal("region", "ATLANTA")),
			Equal("id", 1),
			Equal("name", "John"),
			Like("position", "%manager%"),
			MoreOrEqual("salary", 20000),
		}
		sql, values := CombineSimpleQuery(elems, "AND")
		assert.Equal(t, "(doc_number = ? OR region = ?) AND id = ? AND name = ? AND position LIKE ? AND salary >= ?", sql)
		assert.Equal(t, len(elems)+1, len(values))
	})
}
