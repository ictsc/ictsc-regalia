package contestant

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"
	contestantv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1/contestantv1connect"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/connectdomain"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/contestant/session"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestProblemServiceListDeploymentsRequiresActiveSchedule(t *testing.T) {
	t.Parallel()

	cases := map[string][]*domain.ScheduleData{
		"before contest": {testSchedule("contest", time.Now().Add(time.Hour), time.Now().Add(2*time.Hour))},
		"between slots": {
			testSchedule("slot1", time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour)),
			testSchedule("slot2", time.Now().Add(time.Hour), time.Now().Add(2*time.Hour)),
		},
		"after contest": {testSchedule("contest", time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour))},
	}

	for name, schedules := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			store := newProblemServiceTestStore(schedules)
			client := newProblemServiceTestClient(t, store)

			_, err := client.ListDeployments(t.Context(), connect.NewRequest(&contestantv1.ListDeploymentsRequest{
				Code: "ABC",
			}))
			if got := connect.CodeOf(err); got != connect.CodeFailedPrecondition {
				t.Fatalf("connect.CodeOf(err) = %v, want %v (err=%v)", got, connect.CodeFailedPrecondition, err)
			}
			if store.getProblemByCodeCalls != 0 {
				t.Fatalf("GetProblemByCode called %d times, want 0", store.getProblemByCodeCalls)
			}
		})
	}
}

func TestProblemServiceListDeploymentsWithinActiveSchedule(t *testing.T) {
	t.Parallel()

	store := newProblemServiceTestStore([]*domain.ScheduleData{
		testSchedule("contest", time.Now().Add(-time.Hour), time.Now().Add(time.Hour)),
	})
	client := newProblemServiceTestClient(t, store)

	resp, err := client.ListDeployments(t.Context(), connect.NewRequest(&contestantv1.ListDeploymentsRequest{
		Code: "ABC",
	}))
	if err != nil {
		t.Fatalf("ListDeployments() error = %v", err)
	}
	if got := len(resp.Msg.GetDeployments()); got != 0 {
		t.Fatalf("len(resp.Msg.Deployments) = %d, want 0", got)
	}
	if store.getProblemByCodeCalls == 0 {
		t.Fatal("GetProblemByCode was not called")
	}
}

func TestProblemServiceDeployRequiresActiveSchedule(t *testing.T) {
	t.Parallel()

	cases := map[string][]*domain.ScheduleData{
		"before contest": {testSchedule("contest", time.Now().Add(time.Hour), time.Now().Add(2*time.Hour))},
		"between slots": {
			testSchedule("slot1", time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour)),
			testSchedule("slot2", time.Now().Add(time.Hour), time.Now().Add(2*time.Hour)),
		},
		"after contest": {testSchedule("contest", time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour))},
	}

	for name, schedules := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			store := newProblemServiceTestStore(schedules)
			client := newProblemServiceTestClient(t, store)

			_, err := client.Deploy(t.Context(), connect.NewRequest(&contestantv1.DeployRequest{
				Code: "ABC",
			}))
			if got := connect.CodeOf(err); got != connect.CodeFailedPrecondition {
				t.Fatalf("connect.CodeOf(err) = %v, want %v (err=%v)", got, connect.CodeFailedPrecondition, err)
			}
			if store.createDeploymentCalls != 0 {
				t.Fatalf("CreateDeployment called %d times, want 0", store.createDeploymentCalls)
			}
		})
	}
}

func TestProblemServiceDeployWithinActiveSchedule(t *testing.T) {
	t.Parallel()

	store := newProblemServiceTestStore([]*domain.ScheduleData{
		testSchedule("contest", time.Now().Add(-time.Hour), time.Now().Add(time.Hour)),
	})
	client := newProblemServiceTestClient(t, store)

	resp, err := client.Deploy(t.Context(), connect.NewRequest(&contestantv1.DeployRequest{
		Code: "ABC",
	}))
	if err != nil {
		t.Fatalf("Deploy() error = %v", err)
	}
	if got := resp.Msg.GetDeployment().GetRevision(); got != 1 {
		t.Fatalf("resp.Msg.Deployment.Revision = %d, want 1", got)
	}
	if store.createDeploymentCalls != 1 {
		t.Fatalf("CreateDeployment called %d times, want 1", store.createDeploymentCalls)
	}
}

func TestProblemServiceDeployChecksScheduleBeforeProblemValidation(t *testing.T) {
	t.Parallel()

	store := newProblemServiceTestStore([]*domain.ScheduleData{
		testSchedule("contest", time.Now().Add(time.Hour), time.Now().Add(2*time.Hour)),
	})
	client := newProblemServiceTestClient(t, store)

	_, err := client.Deploy(t.Context(), connect.NewRequest(&contestantv1.DeployRequest{
		Code: "!",
	}))
	if got := connect.CodeOf(err); got != connect.CodeFailedPrecondition {
		t.Fatalf("connect.CodeOf(err) = %v, want %v (err=%v)", got, connect.CodeFailedPrecondition, err)
	}
	if store.getProblemByCodeCalls != 0 {
		t.Fatalf("GetProblemByCode called %d times, want 0", store.getProblemByCodeCalls)
	}
	if store.createDeploymentCalls != 0 {
		t.Fatalf("CreateDeployment called %d times, want 0", store.createDeploymentCalls)
	}
}

type problemServiceTestStore struct {
	userID  uuid.UUID
	teamID  uuid.UUID
	problem *domain.ProblemData

	schedules []*domain.ScheduleData

	deployments []*domain.DeploymentData

	getProblemByCodeCalls int
	createDeploymentCalls int
}

func newProblemServiceTestStore(schedules []*domain.ScheduleData) *problemServiceTestStore {
	return &problemServiceTestStore{
		userID: uuid.Must(uuid.NewV4()),
		teamID: uuid.Must(uuid.NewV4()),
		problem: &domain.ProblemData{
			ID:                          uuid.Must(uuid.NewV4()),
			Code:                        "ABC",
			ProblemType:                 domain.ProblemTypeDescriptive,
			Title:                       "Test Problem",
			MaxScore:                    100,
			Category:                    "network",
			RedeployRule:                domain.RedeployRuleManual,
			SubmissionableScheduleNames: []string{"contest"},
		},
		schedules: schedules,
	}
}

func (s *problemServiceTestStore) GetTeamMemberByID(ctx context.Context, userID uuid.UUID) (*domain.TeamMemberData, error) {
	if userID != s.userID {
		return nil, domain.ErrNotFound
	}
	return &domain.TeamMemberData{
		User: &domain.UserData{
			ID:   s.userID,
			Name: "tester",
		},
		Team: &domain.TeamData{
			ID:           s.teamID,
			Code:         1,
			Name:         "team",
			Organization: "ictsc",
			MaxMembers:   3,
		},
	}, nil
}

func (*problemServiceTestStore) CountTeamMembers(context.Context, uuid.UUID) (uint, error) {
	return 1, nil
}

func (s *problemServiceTestStore) ListProblems(context.Context) ([]*domain.ProblemData, error) {
	return []*domain.ProblemData{s.problem}, nil
}

func (s *problemServiceTestStore) GetProblemByCode(_ context.Context, code string) (*domain.ProblemData, error) {
	s.getProblemByCodeCalls++
	if code != s.problem.Code {
		return nil, domain.ErrNotFound
	}
	return s.problem, nil
}

func (*problemServiceTestStore) GetDescriptiveProblem(context.Context, uuid.UUID) (*domain.DescriptiveProblemData, error) {
	return nil, domain.ErrNotFound
}

func (*problemServiceTestStore) GetTeamProblemScore(context.Context, bool, uuid.UUID, uuid.UUID) (*domain.ScoreData, error) {
	return nil, domain.ErrNotFound
}

func (*problemServiceTestStore) ListTeamProblemScoresByTeamID(context.Context, bool, uuid.UUID) ([]*domain.TeamProblemScoreData, error) {
	return nil, nil
}

func (*problemServiceTestStore) ListTeamProblemScores(context.Context, bool) ([]*domain.TeamProblemScoreData, error) {
	return nil, nil
}

func (s *problemServiceTestStore) ListDeployments(context.Context) ([]*domain.DeploymentData, error) {
	return s.deployments, nil
}

func (s *problemServiceTestStore) GetSchedule(context.Context) ([]*domain.ScheduleData, error) {
	return s.schedules, nil
}

func (s *problemServiceTestStore) RunInTx(ctx context.Context, f func(domain.DeploymentWriter) error) error {
	return f(s)
}

func (s *problemServiceTestStore) CreateDeployment(_ context.Context, input *domain.CreateDeploymentInput) error {
	s.createDeploymentCalls++
	s.deployments = append(s.deployments, &domain.DeploymentData{
		ID:          input.ID,
		TeamCode:    1,
		ProblemCode: s.problem.Code,
		Revision:    input.Revision,
		Events: []*domain.DeploymentEventData{
			{
				Status:     input.Status,
				OccurredAt: input.OccurredAt,
			},
		},
	})
	return nil
}

func (*problemServiceTestStore) UpdateDeploymentStatus(context.Context, *domain.UpdateDeploymentStatusInput) error {
	return nil
}

func newProblemServiceTestClient(t *testing.T, store *problemServiceTestStore) contestantv1connect.ProblemServiceClient {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		err := session.UserSessionStore.Write(r, w, &session.UserSession{UserID: store.userID}, &sessions.Options{
			Path:   "/",
			MaxAge: 60,
		})
		if err != nil {
			t.Fatalf("failed to write session: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	problemService := &ProblemServiceHandler{
		ListDeploymentsEffect: store,
		DeployEffect:          store,
	}
	path, handler := contestantv1connect.NewProblemServiceHandler(
		problemService,
		connect.WithInterceptors(connectdomain.NewErrorInterceptor()),
	)
	mux.Handle(path, handler)

	server := httptest.NewServer(session.NewHandler(sessions.NewCookieStore([]byte("test-secret")))(mux))
	t.Cleanup(server.Close)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	httpClient := server.Client()
	httpClient.Jar = jar

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, server.URL+"/login", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("login request error = %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("login status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}

	return contestantv1connect.NewProblemServiceClient(httpClient, server.URL)
}

func testSchedule(name string, startAt, endAt time.Time) *domain.ScheduleData {
	return &domain.ScheduleData{
		Name:    name,
		StartAt: startAt,
		EndAt:   endAt,
	}
}
