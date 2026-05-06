package boolee

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

import (
	"log"
	"log/slog"
)

var (
	logger = slog.New(slog.NewJSONHandler(log.Writer(), &slog.HandlerOptions{Level: slog.LevelDebug}))
)

func WithBoolee(items ...string) func(*gorm.DB) *gorm.DB {
	if len(items) < 3 {
		panic("at least 3 items are required: columns, operation, comparator")
	}

	comparator := items[len(items)-1]
	operation := items[len(items)-2]
	columns := items[:len(items)-2]

	return func(db *gorm.DB) *gorm.DB {
		for _, column := range columns {
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

			db.Where(query.String())
			logger.Debug("boolee query dispatched", "query", query.String())
		}

		return db
	}
}
