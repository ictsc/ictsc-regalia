package domain

import "fmt"

const (
	MinTeamNumber int64 = 1  // チーム番号の最小値
	MaxTeamNumber int64 = 99 // チーム番号の最大値
)

// チームのドメインモデルを表す構造体
type Team struct {
	teamNumber       int64  // チーム番号
	teamName         string // チーム名
	teamOrganization string // チームの所属組織
	maxTeamMembers   uint   // チームの最大メンバー数
}

// チームを作成する関数
func CreateTeam(teamNumber int64, teamName string, teamOrganization string, maxTeamMembers uint) (*Team, error) {
	if teamNumber < MinTeamNumber || teamNumber > MaxTeamNumber {
		return nil, fmt.Errorf("teamNumber must be between %d and %d", MinTeamNumber, MaxTeamNumber)
	}

	if teamName == "" {
		return nil, fmt.Errorf("teamName cannot be empty")
	}

	if teamOrganization == "" {
		return nil, fmt.Errorf("teamOrganization cannot be empty")
	}

	if maxTeamMembers == 0 {
		return nil, fmt.Errorf("maxTeamMembers must be greater than 0")
	}

	return &Team{
		teamNumber:       teamNumber,
		teamName:         teamName,
		teamOrganization: teamOrganization,
		maxTeamMembers:   maxTeamMembers,
	}, nil
}

func (t *Team) TeamNumber() int64 {
	return t.teamNumber
}

func (t *Team) TeamName() string {
	return t.teamName
}

func (t *Team) TeamOrganization() string {
	return t.teamOrganization
}

func (t *Team) MaxTeamMembers() uint {
	return t.maxTeamMembers
}
