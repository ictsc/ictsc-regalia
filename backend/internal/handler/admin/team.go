package admin

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ictsc/ictsc-regalia/backend/internal/domain"
)

// Teamはチームのレスポンスを表す
type Team struct {
	Code         int64  `json:"code" doc:"チーム番号"`
	Name         string `json:"name" doc:"チーム名"`
	Organization string `json:"organization" doc:"所属組織"`
	MemberLimit  uint   `json:"member_limit" doc:"チームの最大メンバー数"`
}

// ListTeamsOutputはチーム一覧APIのレスポンスを表す
type ListTeamsOutput struct {
	Body struct {
		Teams []Team `json:"teams"`
	}
}

// CreateTeamInputはチーム作成APIのリクエストを表す
type CreateTeamInput struct {
	Body struct {
		Code         int64  `json:"code" minimum:"1" maximum:"99" doc:"チーム番号"`
		Name         string `json:"name" minLength:"1" doc:"チーム名"`
		Organization string `json:"organization" minLength:"1" doc:"所属組織"`
		MemberLimit  uint   `json:"member_limit" minimum:"1" doc:"チームの最大メンバー数"`
	}
}

// CreateTeamOutputはチーム作成APIのレスポンスを表す
type CreateTeamOutput struct {
	Body struct {
		Team Team `json:"team"`
	}
}

// UpdateTeamInputはチーム情報変更APIのリクエストを表す
type UpdateTeamInput struct {
	Code int64 `path:"code" minimum:"1" maximum:"99" doc:"チーム番号"`
	Body struct {
		Name         *string `json:"name,omitempty" minLength:"1" doc:"チーム名"`
		Organization *string `json:"organization,omitempty" minLength:"1" doc:"所属組織"`
		MemberLimit  *uint   `json:"member_limit,omitempty" minimum:"1" doc:"チームの最大メンバー数"`
	}
}

// UpdateTeamOutputはチーム情報変更APIのレスポンスを表す
type UpdateTeamOutput struct {
	Body struct {
		Team Team `json:"team"`
	}
}

// DeleteTeamInputはチーム削除APIのリクエストを表す
type DeleteTeamInput struct {
	Code int64 `path:"code" minimum:"1" maximum:"99" doc:"チーム番号"`
}

// DeleteTeamOutputはチーム削除APIのレスポンスを表す
type DeleteTeamOutput struct{}

// RegisterTeamRoutesは運営向けチームAPIを登録する
func RegisterTeamRoutes(api huma.API) {
	huma.Get(api, "/admin/teams", listTeams)
	huma.Post(api, "/admin/teams", createTeam, func(operation *huma.Operation) {
		operation.DefaultStatus = http.StatusCreated
		operation.Errors = []int{http.StatusBadRequest}
		operation.Summary = "Create admin team"
	})
	huma.Patch(api, "/admin/teams/{code}", updateTeam, func(operation *huma.Operation) {
		operation.Errors = []int{http.StatusBadRequest, http.StatusNotFound}
		operation.Summary = "Update admin team"
	})
	huma.Delete(api, "/admin/teams/{code}", deleteTeam, func(operation *huma.Operation) {
		operation.DefaultStatus = http.StatusNoContent
		operation.Errors = []int{http.StatusNotFound}
		operation.Summary = "Delete admin team"
	})
}

// listTeamsは空のチーム一覧を返すスタブ
func listTeams(
	ctx context.Context,
	input *struct{},
) (*ListTeamsOutput, error) {
	return &ListTeamsOutput{
		Body: struct {
			Teams []Team `json:"teams"`
		}{
			Teams: []Team{},
		},
	}, nil
}

// createTeamは入力からチームのドメインモデルを作成する
func createTeam(
	ctx context.Context,
	input *CreateTeamInput,
) (*CreateTeamOutput, error) {
	team, err := domain.CreateTeam(
		input.Body.Code,
		input.Body.Name,
		input.Body.Organization,
		input.Body.MemberLimit,
	)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	return &CreateTeamOutput{
		Body: struct {
			Team Team `json:"team"`
		}{
			Team: Team{
				Code:         team.TeamNumber(),
				Name:         team.TeamName(),
				Organization: team.TeamOrganization(),
				MemberLimit:  team.MaxTeamMembers(),
			},
		},
	}, nil
}

// updateTeamは空のレスポンスを返すスタブ
func updateTeam(
	ctx context.Context,
	input *UpdateTeamInput,
) (*UpdateTeamOutput, error) {
	return &UpdateTeamOutput{}, nil
}

// deleteTeamは空のレスポンスを返すスタブ
func deleteTeam(
	ctx context.Context,
	input *DeleteTeamInput,
) (*DeleteTeamOutput, error) {
	return &DeleteTeamOutput{}, nil
}
