package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/internal/handler"
)

func TestTeamAPIRoutes(t *testing.T) {
	app := handler.New()

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/admin/teams",
		strings.NewReader(`{"code":12,"name":"とらこんつよいくん","organization":"トラコン大学","member_limit":3}`),
	)
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	app.Router.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d; body = %s", createResponse.Code, http.StatusCreated, createResponse.Body.String())
	}
	var created struct {
		Team struct {
			Code         int64  `json:"code"`
			Name         string `json:"name"`
			Organization string `json:"organization"`
			MemberLimit  uint   `json:"member_limit"`
		} `json:"team"`
	}
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("POST response JSON error = %v", err)
	}
	if created.Team.Code != 12 ||
		created.Team.Name != "とらこんつよいくん" ||
		created.Team.Organization != "トラコン大学" ||
		created.Team.MemberLimit != 3 {
		t.Errorf("created team = %#v, want request values", created.Team)
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/admin/teams", nil)
	listResponse := httptest.NewRecorder()
	app.Router.ServeHTTP(listResponse, listRequest)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d; body = %s", listResponse.Code, http.StatusOK, listResponse.Body.String())
	}

	var response struct {
		Teams []struct {
			Code int64 `json:"code"`
		} `json:"teams"`
	}
	if err := json.Unmarshal(listResponse.Body.Bytes(), &response); err != nil {
		t.Fatalf("response JSON error = %v", err)
	}
	if len(response.Teams) != 0 {
		t.Errorf("len(teams) = %d, want 0", len(response.Teams))
	}

	updateRequest := httptest.NewRequest(
		http.MethodPatch,
		"/admin/teams/12",
		strings.NewReader(`{"name":"新しいチーム名"}`),
	)
	updateRequest.Header.Set("Content-Type", "application/json")
	updateResponse := httptest.NewRecorder()
	app.Router.ServeHTTP(updateResponse, updateRequest)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("PATCH status = %d, want %d; body = %s", updateResponse.Code, http.StatusOK, updateResponse.Body.String())
	}

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/admin/teams/12", nil)
	deleteResponse := httptest.NewRecorder()
	app.Router.ServeHTTP(deleteResponse, deleteRequest)
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("DELETE status = %d, want %d; body = %s", deleteResponse.Code, http.StatusNoContent, deleteResponse.Body.String())
	}
}
