package domain

import (
	"testing"
)

// 新しいチームを作成する関数のテスト
func TestCreateTeam(t *testing.T) {
	team, err := CreateTeam(
		12,
		"とらこんつよいくん",
		"トラコン大学",
		3,
	)
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}

	if got := team.TeamNumber(); got != 12 {
		t.Errorf("TeamNumber() = %d, want 12", got)
	}

	if got := team.TeamName(); got != "とらこんつよいくん" {
		t.Errorf("TeamName() = %q, want %q", got, "とらこんつよいくん")
	}

	if got := team.TeamOrganization(); got != "トラコン大学" {
		t.Errorf("TeamOrganization() = %q, want %q", got, "トラコン大学")
	}

	if got := team.MaxTeamMembers(); got != 3 {
		t.Errorf("MaxTeamMembers() = %d, want 3", got)
	}
}

// 新しいチームを作成する関数がバリデーションエラーを返すことを確認するテスト
func TestCreateTeamAcceptsBoundaryCodes(t *testing.T) {
	tests := []struct {
		name string
		code int64
	}{
		{
			name: "minimum code",
			code: MinTeamNumber,
		},
		{
			name: "maximum code",
			code: MaxTeamNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, err := CreateTeam(tt.code, "とらこんつよいくん", "トラコン大学", 3)
			if err != nil {
				t.Fatalf("CreateTeam() error = %v", err)
			}

			if got := team.TeamNumber(); got != tt.code {
				t.Errorf("TeamNumber() = %d, want %d", got, tt.code)
			}
		})
	}
}

// 新しいチームを作成する関数が入力項目ごとにエラーを返すことを確認するテスト
func TestCreateTeamReturnsValidationError(t *testing.T) {
	tests := []struct {
		name             string
		teamNumber       int64
		teamName         string
		teamOrganization string
		maxTeamMembers   uint
		want             string
	}{
		{
			name:             "invalid team number",
			teamNumber:       0,
			teamName:         "team",
			teamOrganization: "organization",
			maxTeamMembers:   3,
			want:             "teamNumber must be between 1 and 99",
		},
		{
			name:             "empty team name",
			teamNumber:       1,
			teamOrganization: "organization",
			maxTeamMembers:   3,
			want:             "teamName cannot be empty",
		},
		{
			name:           "empty team organization",
			teamNumber:     1,
			teamName:       "team",
			maxTeamMembers: 3,
			want:           "teamOrganization cannot be empty",
		},
		{
			name:             "invalid member limit",
			teamNumber:       1,
			teamName:         "team",
			teamOrganization: "organization",
			want:             "maxTeamMembers must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team, err := CreateTeam(tt.teamNumber, tt.teamName, tt.teamOrganization, tt.maxTeamMembers)
			if team != nil {
				t.Fatal("CreateTeam() returned a team for invalid input")
			}
			if err == nil {
				t.Fatal("CreateTeam() error = nil, want validation error")
			}
			if err.Error() != tt.want {
				t.Errorf("CreateTeam() error = %q, want %q", err, tt.want)
			}
		})
	}
}
