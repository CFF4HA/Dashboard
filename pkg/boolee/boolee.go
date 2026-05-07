package boolee

import (
	"fmt"
	"log"
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

var (
	logger = slog.New(slog.NewJSONHandler(log.Writer(), &slog.HandlerOptions{Level: slog.LevelDebug}))
)

func buildQuery(column, operation, comparator string) string {
	var query strings.Builder
	var element_buffer strings.Builder

	isInString := false
	for _, symbol := range operation {
		if isInString && symbol == '\'' {
			isInString = false
			continue
		} else if isInString && symbol != '\'' {
			element_buffer.WriteRune(symbol)
			continue
		}

		switch symbol {
		case ' ', '\t', '\n':
			continue
		case '\'':
			isInString = true
		case '&':
			if element_buffer.Len() != 0 {
				query.WriteString(fmt.Sprintf("%s %s '%s'", column, comparator, element_buffer.String()))
				element_buffer.Reset()
			}
			query.WriteString(" AND ")
		case '|':
			if element_buffer.Len() != 0 {
				query.WriteString(fmt.Sprintf("%s %s '%s'", column, comparator, element_buffer.String()))
				element_buffer.Reset()
			}
			query.WriteString(" OR ")
		case '!':
			if element_buffer.Len() != 0 {
				query.WriteString(fmt.Sprintf("%s %s '%s'", column, comparator, element_buffer.String()))
				element_buffer.Reset()
			}
			query.WriteString(" NOT ")
		case '(':
			query.WriteString("( ")
			element_buffer.Reset()
		case ')':
			if element_buffer.Len() != 0 {
				query.WriteString(fmt.Sprintf("%s %s '%s'", column, comparator, element_buffer.String()))
				element_buffer.Reset()
			}
			query.WriteString(") ")
		}
	}

	if element_buffer.Len() > 0 {
		query.WriteString(fmt.Sprintf("%s %s '%s'", column, comparator, element_buffer.String()))
	}

	return query.String()
}

// WithBoolee applies the boolean operation to each column independently, ANDing
// all column conditions together.
func WithBoolee(items ...string) func(*gorm.DB) *gorm.DB {
	if len(items) < 3 {
		panic("at least 3 items are required: columns, operation, comparator")
	}

	comparator := items[len(items)-1]
	operation := items[len(items)-2]
	columns := items[:len(items)-2]

	return func(db *gorm.DB) *gorm.DB {
		for _, column := range columns {
			q := buildQuery(column, operation, comparator)
			db = db.Where(q)
			logger.Debug("boolee query dispatched", "query", q)
		}
		return db
	}
}

// WithBooleeAny applies the boolean operation across all columns with OR semantics —
// a row matches if ANY column satisfies the expression.
func WithBooleeAny(items ...string) func(*gorm.DB) *gorm.DB {
	if len(items) < 3 {
		panic("at least 3 items are required: columns, operation, comparator")
	}

	comparator := items[len(items)-1]
	operation := items[len(items)-2]
	columns := items[:len(items)-2]

	return func(db *gorm.DB) *gorm.DB {
		parts := make([]string, 0, len(columns))
		for _, column := range columns {
			if q := buildQuery(column, operation, comparator); q != "" {
				parts = append(parts, "("+q+")")
			}
		}
		if len(parts) == 0 {
			return db
		}
		combined := strings.Join(parts, " OR ")
		logger.Debug("boolee any-query dispatched", "query", combined)
		return db.Where(combined)
	}
}
