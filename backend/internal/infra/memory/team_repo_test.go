package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/internal/domain"
)

func TestRepositoryCreatesAndListsTeamsByCode(t *testing.T) {
	repository := NewRepository()
	for _, code := range []int64{20, 10} {
		team, err := domain.CreateTeam(code, "team", "organization", 3)
		if err != nil {
			t.Fatalf("CreateTeam() error = %v", err)
		}
		if err := repository.CreateTeam(context.Background(), team); err != nil {
			t.Fatalf("CreateTeam() error = %v", err)
		}
	}

	teams, err := repository.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("ListTeams() error = %v", err)
	}
	if len(teams) != 2 {
		t.Fatalf("len(teams) = %d, want 2", len(teams))
	}
	if teams[0].TeamNumber() != 10 || teams[1].TeamNumber() != 20 {
		t.Errorf("team codes = [%d, %d], want [10, 20]", teams[0].TeamNumber(), teams[1].TeamNumber())
	}
}

func TestRepositoryRejectsDuplicateTeamCode(t *testing.T) {
	repository := NewRepository()
	team, err := domain.CreateTeam(10, "team", "organization", 3)
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}
	if err := repository.CreateTeam(context.Background(), team); err != nil {
		t.Fatalf("first CreateTeam() error = %v", err)
	}

	err = repository.CreateTeam(context.Background(), team)
	var domainErr *domain.Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("second CreateTeam() error type = %T, want *domain.Error", err)
	}
	if domainErr.Type != domain.ErrAlreadyExists {
		t.Errorf("error type = %q, want %q", domainErr.Type, domain.ErrAlreadyExists)
	}
}
