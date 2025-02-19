package pg

import "strings"

type columns []string

func (c columns) String() string {
	return strings.Join(c, ", ")
}

func (c columns) As(alias string) string {
	var builder strings.Builder
	for i, col := range c {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(col)
		builder.WriteString(" AS ")
		builder.WriteRune('"')
		builder.WriteString(alias)
		builder.WriteRune('.')
		builder.WriteString(col)
		builder.WriteRune('"')
	}
	return builder.String()
}
