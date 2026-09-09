package core

const DefaultTeamColor = "#A6E35F"

func TeamColor(color string) string {
	if color == "" {
		return DefaultTeamColor
	}
	return color
}
