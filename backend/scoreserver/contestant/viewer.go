package contestant

import (
	"cmp"
	"context"
	"slices"
	"strings"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	adminauth "github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/infra/pg"
)

type (
	ViewerServiceHandler struct {
		contestantv1connect.UnimplementedViewerServiceHandler

		GetViewerEffect       ViewerGetEffect
		ListContestantsEffect ListContestantsEffect
		AdminEnforcer         *adminauth.Enforcer
	}
	ViewerGetEffect interface {
		domain.TeamMemberGetter
		domain.UserProfileReader
	}
	ListContestantsEffect interface {
		domain.TeamMemberProfileReader
	}
)

func newViewerServiceHandler(
	repo *pg.Repository,
	adminEnforcer *adminauth.Enforcer,
) *ViewerServiceHandler {
	return &ViewerServiceHandler{
		GetViewerEffect:       repo,
		ListContestantsEffect: repo,
		AdminEnforcer:         adminEnforcer,
	}
}

func (h *ViewerServiceHandler) GetViewer(
	ctx context.Context, _ *connect.Request[contestantv1.GetViewerRequest],
) (*connect.Response[contestantv1.GetViewerResponse], error) {
	admin, err := h.viewerAdmin(ctx)
	if err != nil {
		return nil, err
	}

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
				Name:  string(teamMember.Name()),
				Type:  contestantv1.ViewerType_VIEWER_TYPE_CONTESTANT,
				Admin: admin,
				Viewer: &contestantv1.Viewer_Contestant{
					Contestant: &contestantv1.ContestantViewer{
						Name:        string(teamMember.Name()),
						DisplayName: profile.DisplayName(),
						Impersonation: func() *contestantv1.Impersonation {
							if userSess.Impersonation == nil {
								return nil
							}
							return &contestantv1.Impersonation{
								AdminName: userSess.Impersonation.AdminName,
							}
						}(),
					},
				},
			},
		}), nil
	}

	signUpSess, err := session.SignUpSessionStore.Get(ctx)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	} else if signUpSess == nil {
		return connect.NewResponse(&contestantv1.GetViewerResponse{
			Viewer: &contestantv1.Viewer{
				Type:  contestantv1.ViewerType_VIEWER_TYPE_UNAUTHENTICATED,
				Admin: admin,
				Viewer: &contestantv1.Viewer_Unauthenticated{
					Unauthenticated: &contestantv1.UnauthenticatedViewer{},
				},
			},
		}), nil
	}

	return connect.NewResponse(&contestantv1.GetViewerResponse{
		Viewer: &contestantv1.Viewer{
			Name:  signUpSess.Discord.Username,
			Type:  contestantv1.ViewerType_VIEWER_TYPE_SIGN_UP,
			Admin: admin,
			Viewer: &contestantv1.Viewer_SignUp{
				SignUp: &contestantv1.SignUpViewer{
					Name:        signUpSess.Discord.Username,
					DisplayName: signUpSess.Discord.GlobalName,
				},
			},
		},
	}), nil
}

func (h *ViewerServiceHandler) ListContestants(
	ctx context.Context,
	_ *connect.Request[contestantv1.ListContestantsRequest],
) (*connect.Response[contestantv1.ListContestantsResponse], error) {
	ok, err := h.AdminEnforcer.Enforce(adminauth.GetViewer(ctx), "contestants", "impersonate")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, connect.NewError(connect.CodePermissionDenied, nil)
	}

	members, err := domain.ListTeamMembers(ctx, h.ListContestantsEffect)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(members, func(a, b *domain.TeamMemberProfile) int {
		if teamCmp := cmp.Compare(int64(a.Team().Code()), int64(b.Team().Code())); teamCmp != 0 {
			return teamCmp
		}
		if nameCmp := strings.Compare(string(a.Name()), string(b.Name())); nameCmp != 0 {
			return nameCmp
		}
		return strings.Compare(a.UserProfile().DisplayName(), b.UserProfile().DisplayName())
	})

	contestants := make([]*contestantv1.ContestantSummary, 0, len(members))
	for _, member := range members {
		contestants = append(contestants, &contestantv1.ContestantSummary{
			Name:        string(member.Name()),
			DisplayName: member.UserProfile().DisplayName(),
			TeamName:    member.Team().Name(),
			TeamCode:    int64(member.Team().Code()),
		})
	}

	return connect.NewResponse(&contestantv1.ListContestantsResponse{
		Contestants: contestants,
	}), nil
}

func (h *ViewerServiceHandler) viewerAdmin(ctx context.Context) (*contestantv1.ViewerAdmin, error) {
	canListContestants, err := h.AdminEnforcer.Enforce(adminauth.GetViewer(ctx), "contestants", "list")
	if err != nil {
		return nil, err
	}
	canImpersonateContestants, err := h.AdminEnforcer.Enforce(adminauth.GetViewer(ctx), "contestants", "impersonate")
	if err != nil {
		return nil, err
	}
	return &contestantv1.ViewerAdmin{
		CanListContestants:        canListContestants,
		CanImpersonateContestants: canImpersonateContestants,
	}, nil
}
