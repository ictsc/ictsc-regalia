package domain

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
)

func FixTeam1(tb testing.TB, data *TeamData) *Team {
	tb.Helper()

	//nolint:mnd
	teamData := &TeamData{
		ID:           uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091"),
		Code:         1,
		Name:         "トラブルシューターズ",
		Organization: "ICTSC Association",
		MaxMembers:   6,
	}
	if data != nil {
		teamData.merge(data)
	}

	team, err := teamData.parse()
	if err != nil {
		tb.Fatal(err)
	}
	return team
}

func FixInvitationCode1(tb testing.TB, data *InvitationCodeData) *InvitationCode {
	tb.Helper()

	var teamData *TeamData
	if data != nil {
		teamData = data.Team
	}

	codeData := &InvitationCodeData{
		ID:        uuid.FromStringOrNil("ad3f83d3-65be-4884-8a03-adb11a8127ef"),
		Team:      FixTeam1(tb, teamData).Data(),
		Code:      "LHNZXGSF7L59WCG9",
		CreatedAt: must(time.Parse(time.RFC3339, "2025-02-02T08:10:00Z")),
		ExpiresAt: must(time.Parse(time.RFC3339, "2038-04-02T15:00:00Z")),
	}
	if data != nil {
		if !data.ID.IsNil() {
			codeData.ID = data.ID
		}
		if data.Code != "" {
			codeData.Code = data.Code
		}
		if data.CreatedAt != (time.Time{}) {
			codeData.CreatedAt = data.CreatedAt
		}
		if data.ExpiresAt != (time.Time{}) {
			codeData.ExpiresAt = data.ExpiresAt
		}
	}

	code, err := codeData.parse()
	if err != nil {
		tb.Fatal(err)
	}
	return code
}

func FixUser1(tb testing.TB, data *UserData) *User {
	tb.Helper()

	userData := &UserData{
		ID:   uuid.FromStringOrNil("3a4ca027-5e02-4ade-8e2d-eddb39adc235"),
		Name: "alice",
	}
	if data != nil {
		if !data.ID.IsNil() {
			userData.ID = data.ID
		}
		if data.Name != "" {
			userData.Name = data.Name
		}
	}

	user, err := userData.parse()
	if err != nil {
		tb.Fatal(err)
	}
	return user
}

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}
