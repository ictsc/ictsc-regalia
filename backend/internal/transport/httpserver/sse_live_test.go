package httpserver

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
)

func TestDeploymentSSEDeliversLiveDatabaseState(t *testing.T) {
	fixture := newContractFixture(t)
	contestant := fixture.seedContestant(t)
	fixture.seedActiveContent(t)
	queuedID := "11111111-1111-4111-8111-111111111111"
	deployment, err := fixture.store.CreateDeployment(context.Background(), core.Deployment{
		RequestID: "22222222-2222-4222-8222-222222222222", TeamCode: contestant.TeamCode, ProblemCode: "P1",
		ContentCommit: testContentCommit, RequestedAt: testNow, LatestStatus: core.DeploymentQueued,
		Events: []core.DeploymentEvent{{EventID: queuedID, OccurredAt: testNow, Status: core.DeploymentQueued}},
	})
	if err != nil {
		t.Fatal(err)
	}
	cookie := fixture.createSession(t, session.Data{Kind: session.KindContestant, ContestantName: contestant.Name}, service.ContestantTTL)
	server := httptest.NewServer(fixture.handler)
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/api/v1/contestant/problems/P1/deployments/stream", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.AddCookie(cookie)
	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	reader := bufio.NewReader(response.Body)
	for range 4 {
		if _, err := reader.ReadString('\n'); err != nil {
			t.Fatalf("read snapshot: %v", err)
		}
	}

	liveID := "33333333-3333-4333-8333-333333333333"
	updated, _, err := fixture.store.AppendDeploymentEvent(context.Background(), deployment.RequestID, core.DeploymentEvent{
		EventID: liveID, OccurredAt: testNow.Add(time.Second), Status: core.DeploymentDeploying,
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(updated)
	if err := fixture.service.Events.Publish(context.Background(), payload); err != nil {
		t.Fatal(err)
	}
	event, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("read live event: %v", err)
	}
	id, _ := reader.ReadString('\n')
	data, _ := reader.ReadString('\n')
	if event != "event: deployment\n" || id != "id: "+liveID+"\n" {
		t.Fatalf("live SSE header = %q %q", event, id)
	}
	if !strings.Contains(data, `"status":"DEPLOYING"`) {
		t.Fatalf("live SSE data = %q", data)
	}
}
