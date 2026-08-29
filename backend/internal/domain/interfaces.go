package domain

import "context"

// TeamListerはチーム一覧の取得に必要な処理を表す
type TeamLister interface {
	ListTeams(ctx context.Context) ([]*Team, error)
}

// TeamCreatorはチームの作成に必要な処理を表す
type TeamCreator interface {
	CreateTeam(ctx context.Context, team *Team) error
}

// TeamRepositoryはチームAPIが利用するRepositoryを表す
type TeamRepository interface {
	TeamLister
	TeamCreator
}
