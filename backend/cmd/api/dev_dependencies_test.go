package main

import (
	"context"
	"net/url"
	"testing"

	"github.com/ictsc/ictsc-regalia/backend/internal/config"
)

func TestPreviewDiscord(t *testing.T) {
	d := previewDiscord{guildID: "preview", roleIDs: []string{"admin"}}
	u, err := url.Parse(d.AuthorizationURL("state-value", "challenge", "https://score.example/api/v1/admin/auth/discord/callback", true))
	if err != nil || u.Host != "score.example" || u.Query().Get("state") != "state-value" {
		t.Fatalf("invalid callback: %v", u)
	}
	a, err := d.Exchange(context.Background(), u.Query().Get("code"), "", "", true)
	if err != nil || a.GuildID != "preview" || a.Identity.Username != "preview-admin" {
		t.Fatalf("invalid admin result: %v", err)
	}
	c, err := d.Exchange(context.Background(), u.Query().Get("code"), "", "", false)
	if err != nil || c.Identity.ID == a.Identity.ID {
		t.Fatal("contestant must have a separate identity")
	}
	if _, err := d.Exchange(context.Background(), "invalid", "", "", true); err == nil {
		t.Fatal("invalid code accepted")
	}
}

func TestPreviewRequiresFakeMode(t *testing.T) {
	cfg := config.Runtime{DevFakeMode: true, DiscordAdminGuildID: "preview", DiscordAdminRoleIDs: []string{"admin"}}
	d, _, _, err := runtimeAdapters(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := d.(previewDiscord); !ok {
		t.Fatal("preview adapter was not selected")
	}
	cfg.DevFakeMode = false
	if _, _, _, err := runtimeAdapters(cfg); err == nil {
		t.Fatal("production must require real credentials")
	}
}
