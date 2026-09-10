package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ictsc/ictsc-regalia/backend/internal/adapter/discord"
	githubadapter "github.com/ictsc/ictsc-regalia/backend/internal/adapter/github"
	"github.com/ictsc/ictsc-regalia/backend/internal/adapter/sstate"
	"github.com/ictsc/ictsc-regalia/backend/internal/config"
	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/infra/memory"
	"github.com/ictsc/ictsc-regalia/backend/internal/infra/postgres"
	redisinfra "github.com/ictsc/ictsc-regalia/backend/internal/infra/redis"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
	"github.com/ictsc/ictsc-regalia/backend/internal/session"
	"github.com/ictsc/ictsc-regalia/backend/internal/transport/httpserver"
)

func run() error {
	cfg, err := config.LoadRuntime()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store, sessions, events, cleanup, err := runtimeStores(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()
	discordClient, contentClient, deploymentClient, err := runtimeAdapters(cfg)
	if err != nil {
		return err
	}
	roles := make(map[string]struct{}, len(cfg.DiscordAdminRoleIDs))
	for _, role := range cfg.DiscordAdminRoleIDs {
		roles[role] = struct{}{}
	}
	svc := service.New(store, sessions, service.Config{
		ContestantRedirectURI: cfg.ContestantRedirectURI,
		AdminRedirectURI:      cfg.AdminRedirectURI,
		ContentRepository:     cfg.GitHubRepository,
		ContentRef:            cfg.GitHubRef,
		CallbackBaseURL:       cfg.CallbackBaseURL,
		AdminGuildID:          cfg.DiscordAdminGuildID,
		ContestantGuildID:     cfg.DiscordContestantGuildID,
		AdminRoleIDs:          roles,
	})
	svc.Discord, svc.Content, svc.Deployments, svc.Events = discordClient, contentClient, deploymentClient, events
	handler, err := httpserver.New(svc, httpserver.Options{
		SecureCookies: cfg.SecureCookies, AllowedOrigins: cfg.AllowedOrigins,
		SStateCallbackToken: cfg.SStateCallbackToken, ReadTimeout: cfg.ReadHeaderTimeout,
	})
	if err != nil {
		return fmt.Errorf("build HTTP router: %w", err)
	}
	server := &http.Server{
		Addr: cfg.Address, Handler: handler, ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout: cfg.ReadHeaderTimeout, WriteTimeout: 0, IdleTimeout: cfg.ReadHeaderTimeout * 12,
	}
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- server.ListenAndServe() }()
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	return nil
}

func runtimeStores(ctx context.Context, cfg config.Runtime) (core.Store, session.Store, service.EventBus, func(), error) {
	if cfg.DevFakeMode && cfg.DatabaseURL == "" && cfg.RedisURL == "" {
		return memory.NewCompetitionStore(), memory.NewSessionStore(), memory.NewEventBus(), func() {}, nil
	}
	if cfg.DatabaseURL == "" || cfg.RedisURL == "" {
		return nil, nil, nil, nil, fmt.Errorf("ICTSC_DATABASE_URL and ICTSC_REDIS_URL must either both be set or both be empty in development fake mode")
	}
	store, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("open PostgreSQL: %w", err)
	}
	redisStore, err := redisinfra.Open(ctx, cfg.RedisURL)
	if err != nil {
		store.Close()
		return nil, nil, nil, nil, fmt.Errorf("open Redis: %w", err)
	}
	cleanup := func() {
		_ = redisStore.Close()
		store.Close()
	}
	return store, redisStore, redisinfra.NewDeploymentBus(redisStore), cleanup, nil
}

func runtimeAdapters(cfg config.Runtime) (service.Discord, service.ContentSource, service.DeploymentGateway, error) {
	var discordClient service.Discord = disabledDiscord{}
	if cfg.DevFakeMode && cfg.DiscordClientID == "" && cfg.DiscordClientSecret == "" && cfg.DiscordAdminGuildID != "" && len(cfg.DiscordAdminRoleIDs) > 0 {
		discordClient = previewDiscord{guildID: cfg.DiscordAdminGuildID, roleIDs: cfg.DiscordAdminRoleIDs}
	}
	var contentClient service.ContentSource = disabledContent{}
	var deploymentClient service.DeploymentGateway = disabledDeployments{}
	if !cfg.DevFakeMode || cfg.DiscordClientID != "" || cfg.DiscordClientSecret != "" {
		client, err := discord.New(discord.Config{
			ClientID: cfg.DiscordClientID, ClientSecret: cfg.DiscordClientSecret,
			AdminGuildID: cfg.DiscordAdminGuildID, ContestantGuildID: cfg.DiscordContestantGuildID, AllowedAdminRoleIDs: cfg.DiscordAdminRoleIDs,
		})
		if err != nil {
			return nil, nil, nil, fmt.Errorf("configure Discord: %w", err)
		}
		discordClient = client
	}
	if !cfg.DevFakeMode || cfg.GitHubRepository != "" {
		client, err := githubadapter.New(githubadapter.Config{
			APIBaseURL: cfg.GitHubAPIBaseURL, APIToken: cfg.GitHubAPIToken,
			ManifestPath: cfg.GitHubManifestPath, Issuer: cfg.GitHubOIDCIssuer,
			Audience: cfg.GitHubOIDCAudience, RepositoryID: cfg.GitHubRepositoryID,
			Repository: cfg.GitHubRepository, Ref: cfg.GitHubRef, WorkflowRef: cfg.GitHubWorkflowRef,
			DiscoveryURL: cfg.GitHubDiscoveryURL, AllowInsecureHTTP: cfg.AllowInsecureHTTP,
		})
		if err != nil {
			return nil, nil, nil, fmt.Errorf("configure GitHub content: %w", err)
		}
		contentClient = client
	}
	if !cfg.DevFakeMode || cfg.SStateBaseURL != "" {
		client, err := sstate.New(sstate.Config{
			BaseURL: cfg.SStateBaseURL, BearerToken: cfg.SStateRequestToken, AllowInsecureHTTP: cfg.AllowInsecureHTTP,
		})
		if err != nil {
			return nil, nil, nil, fmt.Errorf("configure SState: %w", err)
		}
		deploymentClient = client
	}
	return discordClient, contentClient, deploymentClient, nil
}
