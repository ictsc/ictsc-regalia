package contestant

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type (
	ViewerServiceHandler struct {
		contestantv1connect.UnimplementedViewerServiceHandler

		GetViewerEffect ViewerGetEffect
	}
	ViewerGetEffect interface {
		domain.TeamMemberGetter
		domain.UserProfileReader
	}
)

func newViewerServiceHandler(repo *pg.Repository) *ViewerServiceHandler {
	return &ViewerServiceHandler{
		GetViewerEffect: repo,
	}
}

func (h *ViewerServiceHandler) GetViewer(
	ctx context.Context, req *connect.Request[contestantv1.GetViewerRequest],
) (*connect.Response[contestantv1.GetViewerResponse], error) {
	userSess, err := session.UserSessionStore.Get(ctx)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	} else if userSess != nil {
		teamMember, err := domain.UserID(userSess.UserID).TeamMember(ctx, h.GetViewerEffect)
		if err != nil {
			return nil, err
		}
		profile, err := teamMember.Profile(ctx, h.GetViewerEffect)
		if err != nil {
			return nil, err
		}
		return connect.NewResponse(&contestantv1.GetViewerResponse{
			Viewer: &contestantv1.Viewer{
				Name: string(teamMember.Name()),
				Type: contestantv1.ViewerType_VIEWER_TYPE_CONTESTANT,
				Viewer: &contestantv1.Viewer_Contestant{
					Contestant: &contestantv1.ContestantViewer{
						Name:        string(teamMember.Name()),
						DisplayName: profile.DisplayName(),
					},
				},
			},
		}), nil
	}

	signUpSess, err := session.SignUpSessionStore.Get(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		return nil, err
	}
	return connect.NewResponse(&contestantv1.GetViewerResponse{
		Viewer: &contestantv1.Viewer{
			Name: signUpSess.Discord.Username,
			Type: contestantv1.ViewerType_VIEWER_TYPE_SIGN_UP,
			Viewer: &contestantv1.Viewer_SignUp{
				SignUp: &contestantv1.SignUpViewer{
					Name:        signUpSess.Discord.Username,
					DisplayName: signUpSess.Discord.GlobalName,
				},
			},
		},
	}), nil
}
