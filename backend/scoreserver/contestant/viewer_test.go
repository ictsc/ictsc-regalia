package contestant

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid/v5"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	adminauth "github.com/ictsc/ictsc-regalia/backend/scoreserver/admin/auth"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/config"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestViewerServiceListContestantsAllowsListPermission(t *testing.T) {
	t.Parallel()

	enforcer, err := adminauth.NewEnforcer(config.AdminAuthz{
		Policy: "p, role:readonly, contestants, list",
	})
	if err != nil {
		t.Fatalf("NewEnforcer() error = %v", err)
	}

	client := newViewerServiceTestClient(t, viewerServiceTestStore{
		members: []*domain.TeamMemberProfileData{
			newViewerServiceTestMember("alice", "Alice", 10, "Team A"),
		},
	}, enforcer, &adminauth.Viewer{
		Name:   "readonly",
		Groups: []string{"role:readonly"},
	})

	resp, err := client.ListContestants(t.Context(), connect.NewRequest(&contestantv1.ListContestantsRequest{}))
	if err != nil {
		t.Fatalf("ListContestants() error = %v", err)
	}

	if got := len(resp.Msg.GetContestants()); got != 1 {
		t.Fatalf("len(resp.Msg.Contestants) = %d, want 1", got)
	}
	if got := resp.Msg.GetContestants()[0].GetName(); got != "alice" {
		t.Fatalf("resp.Msg.Contestants[0].Name = %q, want %q", got, "alice")
	}
}

func TestViewerServiceListContestantsRejectsImpersonateOnlyPermission(t *testing.T) {
	t.Parallel()

	enforcer, err := adminauth.NewEnforcer(config.AdminAuthz{
		Policy: "p, role:impersonator, contestants, impersonate",
	})
	if err != nil {
		t.Fatalf("NewEnforcer() error = %v", err)
	}

	client := newViewerServiceTestClient(t, viewerServiceTestStore{}, enforcer, &adminauth.Viewer{
		Name:   "impersonator",
		Groups: []string{"role:impersonator"},
	})

	_, err = client.ListContestants(t.Context(), connect.NewRequest(&contestantv1.ListContestantsRequest{}))
	if got := connect.CodeOf(err); got != connect.CodePermissionDenied {
		t.Fatalf("connect.CodeOf(err) = %v, want %v (err=%v)", got, connect.CodePermissionDenied, err)
	}
}

type viewerServiceTestStore struct {
	members []*domain.TeamMemberProfileData
}

func (s viewerServiceTestStore) ListTeamMembers(context.Context) ([]*domain.TeamMemberProfileData, error) {
	return s.members, nil
}

func (viewerServiceTestStore) ListTeamMembersByTeamID(context.Context, uuid.UUID) ([]*domain.TeamMemberProfileData, error) {
	return nil, nil
}

type viewerServiceTestAuthenticator struct {
	viewer *adminauth.Viewer
}

func (a viewerServiceTestAuthenticator) HandleRequest(*http.Request) (*adminauth.Viewer, error) {
	return a.viewer, nil
}

func newViewerServiceTestClient(
	t *testing.T,
	store viewerServiceTestStore,
	enforcer *adminauth.Enforcer,
	viewer *adminauth.Viewer,
) contestantv1connect.ViewerServiceClient {
	t.Helper()

	path, handler := contestantv1connect.NewViewerServiceHandler(&ViewerServiceHandler{
		ListContestantsEffect: store,
		AdminEnforcer:         enforcer,
	})

	mux := http.NewServeMux()
	mux.Handle(path, adminauth.WithAuthn(handler, viewerServiceTestAuthenticator{viewer: viewer}))

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return contestantv1connect.NewViewerServiceClient(server.Client(), server.URL)
}

func newViewerServiceTestMember(name, displayName string, teamCode int64, teamName string) *domain.TeamMemberProfileData {
	return &domain.TeamMemberProfileData{
		User: &domain.UserData{
			Name: name,
		},
		Team: &domain.TeamData{
			Code:         teamCode,
			Name:         teamName,
			Organization: "ictsc",
			MaxMembers:   3,
		},
		Profile: &domain.ProfileData{
			DisplayName: displayName,
		},
	}
}
