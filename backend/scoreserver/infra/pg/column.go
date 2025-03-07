package pg

import "strings"

type columns []string

func (c columns) String(alias string) string {
	if alias == "" {
		return strings.Join(c, ", ")
	}
	var builder strings.Builder
	for i, col := range c {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteRune('"')
		builder.WriteString(alias)
		builder.WriteString("\".\"")
		builder.WriteString(col)
		builder.WriteRune('"')
	}
	return builder.String()
}

func (c columns) As(alias string) string {
	var builder strings.Builder
	for i, col := range c {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteRune('"')
		builder.WriteString(alias)
		builder.WriteString("\".\"")
		builder.WriteString(col)
		builder.WriteString("\" AS \"")
		builder.WriteString(alias)
		builder.WriteRune('.')
		builder.WriteString(col)
		builder.WriteRune('"')
	}
	return builder.String()
}
