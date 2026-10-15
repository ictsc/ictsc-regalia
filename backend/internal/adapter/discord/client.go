package discord

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const (
	defaultAuthorizationEndpoint = "https://discord.com/oauth2/authorize"
	defaultTokenEndpoint         = "https://discord.com/api/v10/oauth2/token"
	defaultAPIBaseURL            = "https://discord.com/api/v10"
	defaultTimeout               = 10 * time.Second
	maxResponseBytes             = 1 << 20
)

var (
	ErrAdminRoleRequired = errors.New("discord admin role is required")
	ErrGuildMembership   = errors.New("discord guild membership is required")
)

var _ service.Discord = (*Client)(nil)

type Config struct {
	ClientID              string
	ClientSecret          string
	AdminGuildID          string
	ContestantGuildID     string
	AllowedAdminRoleIDs   []string
	AuthorizationEndpoint string
	TokenEndpoint         string
	APIBaseURL            string
	HTTPClient            *http.Client
}

type Client struct {
	clientID              string
	clientSecret          string
	adminGuildID          string
	contestantGuildID     string
	allowedAdminRoleIDs   map[string]struct{}
	authorizationEndpoint string
	tokenEndpoint         string
	apiBaseURL            string
	httpClient            *http.Client
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type GuildMember struct {
	GuildID string
	RoleIDs []string
}

type HTTPError struct {
	Operation  string
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("discord %s returned HTTP %d", e.Operation, e.StatusCode)
}

func New(config Config) (*Client, error) {
	if strings.TrimSpace(config.ClientID) == "" {
		return nil, errors.New("discord client ID is required")
	}
	if strings.TrimSpace(config.ClientSecret) == "" {
		return nil, errors.New("discord client secret is required")
	}
	if config.AuthorizationEndpoint == "" {
		config.AuthorizationEndpoint = defaultAuthorizationEndpoint
	}
	if config.TokenEndpoint == "" {
		config.TokenEndpoint = defaultTokenEndpoint
	}
	if config.APIBaseURL == "" {
		config.APIBaseURL = defaultAPIBaseURL
	}
	for name, rawURL := range map[string]string{
		"authorization endpoint": config.AuthorizationEndpoint,
		"token endpoint":         config.TokenEndpoint,
		"API base URL":           config.APIBaseURL,
	} {
		parsed, err := url.Parse(rawURL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return nil, fmt.Errorf("discord %s must be an absolute URL", name)
		}
	}
	if config.AdminGuildID != "" && !isSnowflake(config.AdminGuildID) {
		return nil, errors.New("discord admin guild ID must be a numeric snowflake")
	}
	if config.ContestantGuildID != "" && !isSnowflake(config.ContestantGuildID) {
		return nil, errors.New("discord contestant guild ID must be a numeric snowflake")
	}
	allowedRoles := make(map[string]struct{}, len(config.AllowedAdminRoleIDs))
	for _, roleID := range config.AllowedAdminRoleIDs {
		if !isSnowflake(roleID) {
			return nil, errors.New("discord admin role IDs must be numeric snowflakes")
		}
		allowedRoles[roleID] = struct{}{}
	}
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{
		clientID:              config.ClientID,
		clientSecret:          config.ClientSecret,
		adminGuildID:          config.AdminGuildID,
		contestantGuildID:     config.ContestantGuildID,
		allowedAdminRoleIDs:   allowedRoles,
		authorizationEndpoint: strings.TrimRight(config.AuthorizationEndpoint, "/"),
		tokenEndpoint:         config.TokenEndpoint,
		apiBaseURL:            strings.TrimRight(config.APIBaseURL, "/"),
		httpClient:            httpClient,
	}, nil
}

func (c *Client) AuthorizationURL(state, codeChallenge, redirectURI string, admin bool) string {
	values := url.Values{
		"client_id":             {c.clientID},
		"redirect_uri":          {redirectURI},
		"response_type":         {"code"},
		"scope":                 {strings.Join(scopes(admin || c.contestantGuildID != ""), " ")},
		"state":                 {state},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
	}
	return c.authorizationEndpoint + "?" + values.Encode()
}

func (c *Client) Exchange(ctx context.Context, code, codeVerifier, redirectURI string, admin bool) (service.DiscordResult, error) {
	if code == "" || codeVerifier == "" || redirectURI == "" {
		return service.DiscordResult{}, errors.New("discord authorization code, PKCE verifier, and redirect URI are required")
	}
	token, err := c.exchangeToken(ctx, code, codeVerifier, redirectURI)
	if err != nil {
		return service.DiscordResult{}, err
	}
	identity, err := c.User(ctx, token.AccessToken)
	if err != nil {
		return service.DiscordResult{}, err
	}
	result := service.DiscordResult{Identity: identity}
	guildID := c.contestantGuildID
	if admin {
		guildID = c.adminGuildID
	}
	if !admin && guildID == "" {
		return result, nil
	}
	if guildID == "" {
		return service.DiscordResult{}, errors.New("discord admin guild ID is not configured")
	}
	member, err := c.GuildMember(ctx, token.AccessToken, guildID)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			return service.DiscordResult{}, &guildMembershipError{cause: err}
		}
		return service.DiscordResult{}, err
	}
	result.GuildID = member.GuildID
	result.RoleIDs = member.RoleIDs
	return result, nil
}

func (c *Client) User(ctx context.Context, accessToken string) (core.DiscordIdentity, error) {
	var response struct {
		ID         string  `json:"id"`
		Username   string  `json:"username"`
		GlobalName *string `json:"global_name"`
	}
	if err := c.getJSON(ctx, c.apiBaseURL+"/users/@me", accessToken, "get current user", &response); err != nil {
		return core.DiscordIdentity{}, err
	}
	if !isSnowflake(response.ID) || strings.TrimSpace(response.Username) == "" {
		return core.DiscordIdentity{}, errors.New("discord current user response is invalid")
	}
	displayName := response.Username
	if response.GlobalName != nil && strings.TrimSpace(*response.GlobalName) != "" {
		displayName = *response.GlobalName
	}
	return core.DiscordIdentity{ID: response.ID, Username: response.Username, DisplayName: displayName}, nil
}

func (c *Client) GuildMember(ctx context.Context, accessToken, guildID string) (GuildMember, error) {
	if !isSnowflake(guildID) {
		return GuildMember{}, errors.New("discord guild ID must be a numeric snowflake")
	}
	var response struct {
		Roles []string `json:"roles"`
	}
	endpoint := c.apiBaseURL + "/users/@me/guilds/" + url.PathEscape(guildID) + "/member"
	if err := c.getJSON(ctx, endpoint, accessToken, "get current guild member", &response); err != nil {
		return GuildMember{}, err
	}
	for _, roleID := range response.Roles {
		if !isSnowflake(roleID) {
			return GuildMember{}, errors.New("discord guild member response contains an invalid role ID")
		}
	}
	sort.Strings(response.Roles)
	response.Roles = compactStrings(response.Roles)
	return GuildMember{GuildID: guildID, RoleIDs: response.Roles}, nil
}

func (c *Client) HasAllowedAdminRole(roleIDs []string) bool {
	for _, roleID := range roleIDs {
		for allowedRoleID := range c.allowedAdminRoleIDs {
			if subtle.ConstantTimeCompare([]byte(roleID), []byte(allowedRoleID)) == 1 {
				return true
			}
		}
	}
	return false
}

func (c *Client) RequireAllowedAdminRole(roleIDs []string) error {
	if !c.HasAllowedAdminRole(roleIDs) {
		return ErrAdminRoleRequired
	}
	return nil
}

func GeneratePKCE() (verifier, challenge string, err error) {
	bytes := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, bytes); err != nil {
		return "", "", fmt.Errorf("generate PKCE verifier: %w", err)
	}
	verifier = base64.RawURLEncoding.EncodeToString(bytes)
	return verifier, S256Challenge(verifier), nil
}

func S256Challenge(verifier string) string {
	digest := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(digest[:])
}

func (c *Client) exchangeToken(ctx context.Context, code, codeVerifier, redirectURI string) (Token, error) {
	form := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"code_verifier": {codeVerifier},
		"redirect_uri":  {redirectURI},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return Token{}, fmt.Errorf("create discord token request: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return Token{}, fmt.Errorf("exchange discord authorization code: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Token{}, readHTTPError("exchange authorization code", response)
	}
	var token Token
	if err := decodeJSON(response.Body, &token); err != nil {
		return Token{}, fmt.Errorf("decode discord token response: %w", err)
	}
	if token.AccessToken == "" || !strings.EqualFold(token.TokenType, "Bearer") {
		return Token{}, errors.New("discord token response is invalid")
	}
	return token, nil
}

func (c *Client) getJSON(ctx context.Context, endpoint, accessToken, operation string, target any) error {
	if strings.TrimSpace(accessToken) == "" {
		return errors.New("discord access token is required")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create discord %s request: %w", operation, err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("discord %s: %w", operation, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return readHTTPError(operation, response)
	}
	if err := decodeJSON(response.Body, target); err != nil {
		return fmt.Errorf("decode discord %s response: %w", operation, err)
	}
	return nil
}

func readHTTPError(operation string, response *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
	return &HTTPError{Operation: operation, StatusCode: response.StatusCode, Body: strings.TrimSpace(string(body))}
}

func decodeJSON(reader io.Reader, target any) error {
	decoder := json.NewDecoder(io.LimitReader(reader, maxResponseBytes))
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("unexpected trailing JSON value")
		}
		return err
	}
	return nil
}

func scopes(admin bool) []string {
	if admin {
		return []string{"identify", "guilds.members.read"}
	}
	return []string{"identify"}
}

func isSnowflake(value string) bool {
	if value == "" {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func compactStrings(values []string) []string {
	if len(values) < 2 {
		return values
	}
	result := values[:1]
	for _, value := range values[1:] {
		if value != result[len(result)-1] {
			result = append(result, value)
		}
	}
	return result
}
