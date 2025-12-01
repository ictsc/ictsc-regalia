package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	adminv1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
	"github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1/adminv1connect"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	var option Option
	flag.Func("endpoint", "Admin API Endpoint", func(s string) error {
		u, err := url.Parse(s)
		if err != nil {
			return errors.WithStack(err)
		}
		option.Endpoint = u
		return nil
	})
	now := time.Now()
	flag.TextVar(
		&option.InvitationExpires,
		"invitation-expires",
		time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()),
		"Invitation Expires",
	)
	flag.Parse()
	if option.Endpoint == nil {
		log.Println("endpoint is required")
		os.Exit(1)
	}
	option.Token = os.Getenv("ICTSCORE_TOKEN")

	os.Exit(start(option))
}

type Option struct {
	Endpoint          *url.URL
	Token             string
	InvitationExpires time.Time
}

type teamEntry struct {
	ID           uint32 `json:"id"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Count        uint32 `json:"count"`

	InvitationCode string `json:"invitation_code,omitempty"`
}

func (e *teamEntry) Is(team *adminv1.Team) bool {
	return e.Name == team.GetName() && e.Organization == team.GetOrganization() && e.Count == team.GetMemberLimit()
}

func start(option Option) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var entries []*teamEntry
	if err := json.NewDecoder(os.Stdin).Decode(&entries); err != nil {
		log.Printf("Failed to decode input: %v", err)
		return 1
	}

	httpClient := http.DefaultClient
	if option.Token != "" {
		httpClient = oauth2.NewClient(
			context.WithValue(ctx, oauth2.HTTPClient, httpClient),
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: option.Token}),
		)
	}
	teamClient := adminv1connect.NewTeamServiceClient(
		httpClient, option.Endpoint.String(), connect.WithGRPC(),
	)
	invitationClient := adminv1connect.NewInvitationServiceClient(
		httpClient, option.Endpoint.String(), connect.WithGRPC(),
	)

	invitaionCodeResp, err := invitationClient.ListInvitationCodes(ctx, connect.NewRequest(&adminv1.ListInvitationCodesRequest{}))
	if err != nil {
		log.Printf("Failed to get invitation codes: %v", err)
		return 1
	}
	invitationCodes := invitaionCodeResp.Msg.GetInvitationCodes()

	for _, entry := range entries {
		//nolint:nestif // ログのロジックを外に出すのはそれはそれで問題
		if teamResp, err := teamClient.GetTeam(
			ctx, connect.NewRequest(&adminv1.GetTeamRequest{Code: int64(entry.ID)}),
		); err != nil {
			if connectErr := (*connect.Error)(nil); errors.As(err, &connectErr) && connectErr.Code() == connect.CodeNotFound {
				if err := createTeam(ctx, teamClient, entry); err != nil {
					log.Printf("[team:%d] failed to create team: %v", entry.ID, err)
					continue
				}
				log.Printf("[team:%d] created", entry.ID)
			} else {
				log.Printf("[team:%d] Failed to get team: %v", entry.ID, err)
				continue
			}
		} else if team := teamResp.Msg.GetTeam(); !entry.Is(team) {
			if err := updateTeam(ctx, teamClient, team, entry); err != nil {
				log.Printf("[team:%d] failed to update team: %v", entry.ID, err)
				continue
			} else {
				log.Printf("[team:%d] updated", entry.ID)
			}
		} else {
			log.Printf("[team:%d] unchanged", entry.ID)
		}

		idx := slices.IndexFunc(invitationCodes, func(ic *adminv1.InvitationCode) bool {
			return ic.GetTeamCode() == int64(entry.ID)
		})
		if idx >= 0 {
			entry.InvitationCode = invitationCodes[idx].GetCode()
		} else {
			code, err := createInvitationCode(ctx, invitationClient, entry, option.InvitationExpires)
			if err != nil {
				log.Printf("[team:%d] failed to create invitation code: %v", entry.ID, err)
				continue
			}
			log.Printf("[team:%d] invitation code created", entry.ID)
			entry.InvitationCode = code
		}
	}

	if err := json.NewEncoder(os.Stdout).Encode(entries); err != nil {
		log.Printf("Failed to encode output: %v", err)
		return 1
	}

	return 0
}

func createTeam(ctx context.Context, client adminv1connect.TeamServiceClient, entry *teamEntry) error {
	if _, err := client.CreateTeam(ctx, connect.NewRequest(&adminv1.CreateTeamRequest{
		Team: &adminv1.Team{
			Code:         int64(entry.ID),
			Name:         entry.Name,
			Organization: entry.Organization,
			MemberLimit:  entry.Count,
		},
	})); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func updateTeam(ctx context.Context, client adminv1connect.TeamServiceClient, team *adminv1.Team, entry *teamEntry) error {
	if _, err := client.UpdateTeam(ctx, connect.NewRequest(&adminv1.UpdateTeamRequest{
		Team: &adminv1.Team{
			Code:         team.GetCode(),
			Name:         entry.Name,
			Organization: entry.Organization,
			MemberLimit:  entry.Count,
		}})); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func createInvitationCode(ctx context.Context, client adminv1connect.InvitationServiceClient, entry *teamEntry, expiresAt time.Time) (string, error) {
	invitationCode := &adminv1.InvitationCode{
		TeamCode:  int64(entry.ID),
		ExpiresAt: timestamppb.New(expiresAt),
	}
	if entry.InvitationCode != "" {
		invitationCode.Code = entry.InvitationCode
	}

	resp, err := client.CreateInvitationCode(ctx, connect.NewRequest(&adminv1.CreateInvitationCodeRequest{
		InvitationCode: invitationCode,
	}))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return resp.Msg.GetInvitationCode().GetCode(), nil
}
