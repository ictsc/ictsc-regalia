package domain

import (
	"testing"

	"github.com/gofrs/uuid/v5"
)

func FixTeam1(tb testing.TB) *Team {
	tb.Helper()

	return &team{
		teamID:       teamID(uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091")),
		code:         teamCode(1),
		name:         "トラブルシューターズ",
		organization: "ICTSC Association",
	}
}
