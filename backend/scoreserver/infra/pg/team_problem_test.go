package pg_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/pkg/pgtest"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

func TestListProblemsScoreByTeamID(t *testing.T) {
	t.Parallel()

	// DB をセットアップし、リポジトリを生成
	db := pgtest.SetupDB(t)
	repo := pg.NewRepository(db)
	teamID := uuid.FromStringOrNil("a1de8fe6-26c8-42d7-b494-dea48e409091")

	actual, err := repo.ListProblemsScoreByTeamID(t.Context(), teamID)
	if err != nil {
		t.Fatalf("ListProblemsScoreByTeamID failed: %v", err)
	}

	// 期待される結果
	expected := []*domain.TeamProblemScoreData{
		{
			// チーム "トラブルシューターズ" が問題Aを解いた結果
			ProblemID:   uuid.FromStringOrNil("16643c32-c686-44ba-996b-2fbe43b54513"),
			MarkedScore: 80,
			Penalty:     0,
			TotalScore:  80,
			// maxScore は出力に含まれていないため無視する
		},
		{
			// チーム "トラブルシューターズ" が問題Bを解いた結果
			ProblemID:   uuid.FromStringOrNil("24f6aef0-5dcd-4032-825b-d1b19174a6f2"),
			MarkedScore: 80,
			Penalty:     0,
			TotalScore:  80,
		},
	}

	// maxScore フィールドは比較対象から除外する
	opts := cmpopts.IgnoreFields(domain.Score{}, "max")

	if diff := cmp.Diff(expected, actual, opts); diff != "" {
		t.Errorf("ListProblemsScoreByTeamID mismatch (-want +got):\n%s", diff)
	}
}
