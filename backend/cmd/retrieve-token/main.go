package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"slices"

	"github.com/cockroachdb/errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// nolint:gochecknonglobals
var (
	flagIssuer      = flag.String("issuer", "", "Issuer")
	flagClientID    = flag.String("client-id", "", "Client ID")
	flagRedirectURI = flag.String("redirect-uri", "http://localhost:8080/callback", "Redirect URI")

	additionalScopes []string
)

func main() {
	flag.Func("scope", "Additional Scopes", func(s string) error {
		additionalScopes = append(additionalScopes, s)
		return nil
	})
	flag.Parse()
	os.Exit(start())
}

type IssuerClaims struct {
	GrantTypesSupported               []string `json:"grant_types_supported,omitempty"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported,omitempty"`
	ScopesSupported                   []string `json:"scopes_supported,omitempty"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported,omitempty"`
}

// nolint:funlen,cyclop
func start() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	issuer := *flagIssuer
	if issuer == "" {
		log.Println("issuer is required")
		return 1
	}

	clientID := *flagClientID
	if clientID == "" {
		log.Println("client-id is required")
		return 1
	}

	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		log.Println("CLIENT_SECRET is required")
		return 1
	}

	httpClient := http.DefaultClient
	ctx = oidc.ClientContext(ctx, httpClient)

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Println("failed to create provider:", err.Error())
		return 1
	}
	var issuerClaims IssuerClaims
	if err := provider.Claims(&issuerClaims); err != nil {
		log.Println("failed to get claims:", err.Error())
		return 1
	}

	scopes := []string{oidc.ScopeOpenID}
	scopes = append(scopes, additionalScopes...)
	if len(issuerClaims.ScopesSupported) > 0 {
		for _, scope := range scopes {
			if slices.Index(issuerClaims.ScopesSupported, scope) == -1 {
				log.Println("unsupported scope:", scope)
				return 1
			}
		}
	}

	verifier := provider.VerifierContext(ctx, &oidc.Config{ClientID: clientID})

	endpoint := provider.Endpoint()
	if len(issuerClaims.TokenEndpointAuthMethodsSupported) > 0 {
		if slices.Contains(issuerClaims.TokenEndpointAuthMethodsSupported, "client_secret_post") {
			endpoint.AuthStyle = oauth2.AuthStyleInParams
		} else if slices.Contains(issuerClaims.TokenEndpointAuthMethodsSupported, "client_secret_basic") {
			endpoint.AuthStyle = oauth2.AuthStyleInHeader
		}
	}
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     endpoint,
		Scopes:       scopes,
	}

	// auth code
	redirectURI, err := url.Parse(*flagRedirectURI)
	if err != nil {
		log.Println("failed to parse redirect uri:", err.Error())
		return 1
	}
	oauthConfig.RedirectURL = redirectURI.String()

	oauth2State := oauth2.GenerateVerifier()
	var authCodeURLOpts []oauth2.AuthCodeOption
	var exchangeOpts []oauth2.AuthCodeOption
	if slices.Contains(issuerClaims.CodeChallengeMethodsSupported, "S256") {
		pkceVerifier := oauth2.GenerateVerifier()
		authCodeURLOpts = append(authCodeURLOpts, oauth2.S256ChallengeOption(pkceVerifier))
		exchangeOpts = append(exchangeOpts, oauth2.VerifierOption(pkceVerifier))
	}
	authCodeURL := oauthConfig.AuthCodeURL(oauth2State, authCodeURLOpts...)
	log.Println("auth code url:", authCodeURL)

	authCode, authState, err := waitForCallback(ctx, redirectURI)
	if err != nil {
		log.Println("failed to wait for callback:", err.Error())
		return 1
	}
	if authState != oauth2State {
		log.Println("invalid state")
		return 1
	}

	token, err := oauthConfig.Exchange(ctx, authCode, exchangeOpts...)
	if err != nil {
		log.Println("failed to exchange token:", err.Error())
		return 1
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("no id token")
		return 1
	}

	if _, err := verifier.Verify(ctx, rawIDToken); err != nil {
		log.Println("failed to verify id token:", err.Error())
		return 1
	}

	//nolint:forbidigo
	fmt.Println(rawIDToken)

	return 0
}

func waitForCallback(ctx context.Context, redirectURI *url.URL) (string, string, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	redirectPath := redirectURI.Path
	if redirectPath == "" {
		redirectPath = "/"
	}

	var authCode, authState string
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+redirectPath, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "<html><script>close()</script></html>")

		if errMsg := req.FormValue("error"); errMsg != "" {
			cancel(errors.Newf("%s: %s", errMsg, req.FormValue("error_description")))
		}
		authCode = req.FormValue("code")
		if authCode == "" {
			cancel(errors.New("no auth code"))
		}
		authState = req.FormValue("state")
		if authState == "" {
			cancel(errors.New("no auth state"))
		}

		cancel(nil)
	})

	log.Println("starting server at", redirectURI.Host)
	//nolint:gosec
	server := &http.Server{Addr: redirectURI.Host, Handler: mux}
	go func() {
		<-ctx.Done()
		_ = server.Close()
	}()
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return "", "", errors.Wrap(err, "failed to start server")
	}
	if err := context.Cause(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return "", "", (err)
	}
	return authCode, authState, nil
}
