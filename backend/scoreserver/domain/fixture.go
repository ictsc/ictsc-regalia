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

func FixTeamMember1(tb testing.TB, data *TeamMemberData) *TeamMember {
	tb.Helper()

	var (
		teamData *TeamData
		userData *UserData
	)
	if data != nil {
		teamData = data.Team
		userData = data.User
	}
	return &TeamMember{
		user: FixUser1(tb, userData),
		team: FixTeam1(tb, teamData),
	}
}

func FixDescriptiveProblem1(tb testing.TB, data *DescriptiveProblemData) *DescriptiveProblem {
	tb.Helper()

	if data != nil {
		tb.Fatal("additional data is not supported")
	}

	//nolint:mnd
	problemData := &DescriptiveProblemData{
		Problem: &ProblemData{
			ID:           uuid.FromStringOrNil("24f6aef0-5dcd-4032-825b-d1b19174a6f2"),
			Code:         "ZZB",
			ProblemType:  ProblemTypeDescriptive,
			Title:        "Problem 1",
			MaxScore:     100,
			Category:     "Network",
			RedeployRule: RedeployRuleUnredeployable,
		},
		Content: &ProblemContentData{
			Body:        "This is a problem.",
			Explanation: "This is an explanation.",
		},
	}
	problem, err := problemData.parse()
	if err != nil {
		tb.Fatal(err)
	}
	return problem
}

func FixDescriptiveProblemManual(tb testing.TB) *DescriptiveProblem {
	tb.Helper()

	//nolint:mnd
	problemData := &DescriptiveProblemData{
		Problem: &ProblemData{
			ID:           uuid.FromStringOrNil("35f7bf01-6ede-5043-926c-e2c20c285b03"),
			Code:         "MAN",
			ProblemType:  ProblemTypeDescriptive,
			Title:        "Manual Redeploy Problem",
			MaxScore:     100,
			Category:     "Network",
			RedeployRule: RedeployRuleManual,
		},
		Content: &ProblemContentData{
			Body:        "This is a manual redeploy problem.",
			Explanation: "This is a manual redeploy explanation.",
		},
	}
	problem, err := problemData.parse()
	if err != nil {
		tb.Fatal(err)
	}
	return problem
}

func FixNotice1(tb testing.TB, data *NoticeData) *Notice {
	tb.Helper()

	if data != nil {
		tb.Fatal("additional data is not supported")
	}

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour) //nolint:mnd //

	noticeData := &NoticeData{
		Slug:          "test-notice",
		Title:         "テストお知らせ", //nolint:gosmopolitan // fixture
		Markdown:      "これはサンプルです。\n",
		EffectiveFrom: yesterday,
	}

	if data != nil {
		if data.Slug != "" {
			noticeData.Slug = data.Slug
		}
		if data.Title != "" {
			noticeData.Title = data.Title
		}
		if data.Markdown != "" {
			noticeData.Markdown = data.Markdown
		}
		if !data.EffectiveFrom.IsZero() {
			noticeData.EffectiveFrom = data.EffectiveFrom
		}
	}

	notice, err := noticeData.parse()
	if err != nil {
		tb.Fatal(err)
	}

	return notice
}

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}
