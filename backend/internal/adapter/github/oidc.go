package github

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/service"
)

const (
	defaultIssuer        = "https://token.actions.githubusercontent.com"
	defaultAPIBaseURL    = "https://api.github.com"
	defaultManifestPath  = "content/manifest.yaml"
	defaultHTTPTimeout   = 15 * time.Second
	defaultJWKSCacheTTL  = 5 * time.Minute
	defaultClockSkew     = 30 * time.Second
	maxOIDCDocumentBytes = 1 << 20
	maxJWTBytes          = 32 << 10
)

var ErrInvalidActionsToken = errors.New("invalid GitHub Actions OIDC token")

var _ service.ContentSource = (*Client)(nil)

type Config struct {
	APIBaseURL        string
	APIToken          string
	ManifestPath      string
	Issuer            string
	Audience          string
	RepositoryID      string
	Repository        string
	Ref               string
	WorkflowRef       string
	DiscoveryURL      string
	HTTPClient        *http.Client
	Now               func() time.Time
	JWKSCacheTTL      time.Duration
	ClockSkew         time.Duration
	AllowInsecureHTTP bool
}

type Client struct {
	apiBaseURL   string
	apiToken     string
	manifestPath string
	contentRoot  string

	issuer       string
	audience     string
	repositoryID string
	repository   string
	ref          string
	workflowRef  string
	discoveryURL string
	httpClient   *http.Client
	now          func() time.Time
	cacheTTL     time.Duration
	clockSkew    time.Duration

	keysMu       sync.Mutex
	keys         map[string]*rsa.PublicKey
	keysExpireAt time.Time
}

type verificationClaims struct {
	Issuer         string        `json:"iss"`
	Audience       audienceClaim `json:"aud"`
	Subject        string        `json:"sub"`
	ExpiresAt      int64         `json:"exp"`
	NotBefore      *int64        `json:"nbf"`
	IssuedAt       int64         `json:"iat"`
	RepositoryID   string        `json:"repository_id"`
	Repository     string        `json:"repository"`
	Ref            string        `json:"ref"`
	JobWorkflowRef string        `json:"job_workflow_ref"`
}

type audienceClaim []string

func (a *audienceClaim) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*a = []string{single}
		return nil
	}
	var multiple []string
	if err := json.Unmarshal(data, &multiple); err != nil {
		return errors.New("aud must be a string or an array of strings")
	}
	*a = multiple
	return nil
}

type tokenHeader struct {
	Algorithm string `json:"alg"`
	KeyID     string `json:"kid"`
	Type      string `json:"typ"`
}

type discoveryDocument struct {
	Issuer  string `json:"issuer"`
	JWKSURI string `json:"jwks_uri"`
}

type jwksDocument struct {
	Keys []jsonWebKey `json:"keys"`
}

type jsonWebKey struct {
	KeyType   string `json:"kty"`
	Use       string `json:"use"`
	KeyID     string `json:"kid"`
	Algorithm string `json:"alg"`
	Modulus   string `json:"n"`
	Exponent  string `json:"e"`
}

func New(config Config) (*Client, error) {
	if config.APIBaseURL == "" {
		config.APIBaseURL = defaultAPIBaseURL
	}
	if config.ManifestPath == "" {
		config.ManifestPath = defaultManifestPath
	}
	if config.Issuer == "" {
		config.Issuer = defaultIssuer
	}
	if config.DiscoveryURL == "" {
		config.DiscoveryURL = strings.TrimRight(config.Issuer, "/") + "/.well-known/openid-configuration"
	}
	if config.Audience == "" || config.RepositoryID == "" || config.Repository == "" || config.Ref == "" || config.WorkflowRef == "" {
		return nil, errors.New("GitHub OIDC audience, repository ID, repository, ref, and workflow ref are required")
	}
	if _, _, err := splitRepository(config.Repository); err != nil {
		return nil, fmt.Errorf("configured GitHub repository: %w", err)
	}
	if !safeRepositoryPath(config.ManifestPath) {
		return nil, errors.New("GitHub manifest path must be a clean repository-relative path")
	}
	apiURL, err := validateEndpoint(config.APIBaseURL, config.AllowInsecureHTTP)
	if err != nil {
		return nil, fmt.Errorf("GitHub API base URL: %w", err)
	}
	discoveryURL, err := validateEndpoint(config.DiscoveryURL, config.AllowInsecureHTTP)
	if err != nil {
		return nil, fmt.Errorf("GitHub OIDC discovery URL: %w", err)
	}
	issuerURL, err := validateEndpoint(config.Issuer, config.AllowInsecureHTTP)
	if err != nil {
		return nil, fmt.Errorf("GitHub OIDC issuer: %w", err)
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: defaultHTTPTimeout}
	}
	if config.Now == nil {
		config.Now = func() time.Time { return time.Now().UTC() }
	}
	if config.JWKSCacheTTL <= 0 {
		config.JWKSCacheTTL = defaultJWKSCacheTTL
	}
	if config.ClockSkew <= 0 {
		config.ClockSkew = defaultClockSkew
	}
	return &Client{
		apiBaseURL:   strings.TrimRight(apiURL.String(), "/"),
		apiToken:     config.APIToken,
		manifestPath: config.ManifestPath,
		contentRoot:  repositoryDir(config.ManifestPath),
		issuer:       strings.TrimRight(issuerURL.String(), "/"),
		audience:     config.Audience,
		repositoryID: config.RepositoryID,
		repository:   config.Repository,
		ref:          config.Ref,
		workflowRef:  config.WorkflowRef,
		discoveryURL: discoveryURL.String(),
		httpClient:   config.HTTPClient,
		now:          config.Now,
		cacheTTL:     config.JWKSCacheTTL,
		clockSkew:    config.ClockSkew,
	}, nil
}

func (c *Client) VerifyActionsToken(ctx context.Context, token string) (service.ActionsClaims, error) {
	header, claims, signingInput, signature, err := parseJWT(token)
	if err != nil {
		return service.ActionsClaims{}, invalidToken(err)
	}
	if header.Algorithm != "RS256" || header.KeyID == "" || (header.Type != "" && header.Type != "JWT") {
		return service.ActionsClaims{}, invalidToken(errors.New("unsupported JWT header"))
	}
	key, err := c.signingKey(ctx, header.KeyID)
	if err != nil {
		var upstream interface{ MachineTokenUpstream() bool }
		if errors.As(err, &upstream) && upstream.MachineTokenUpstream() {
			return service.ActionsClaims{}, err
		}
		return service.ActionsClaims{}, invalidToken(err)
	}
	digest := sha256.Sum256([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], signature); err != nil {
		return service.ActionsClaims{}, invalidToken(errors.New("JWT signature verification failed"))
	}
	if err := c.validateClaims(claims); err != nil {
		return service.ActionsClaims{}, invalidToken(err)
	}
	return service.ActionsClaims{
		Issuer:       claims.Issuer,
		Audience:     c.audience,
		RepositoryID: claims.RepositoryID,
		Repository:   claims.Repository,
		Ref:          claims.Ref,
		WorkflowRef:  claims.JobWorkflowRef,
		Subject:      claims.Subject,
	}, nil
}

func (c *Client) validateClaims(claims verificationClaims) error {
	if claims.Issuer != c.issuer {
		return errors.New("iss claim does not match")
	}
	if !containsString(claims.Audience, c.audience) {
		return errors.New("aud claim does not match")
	}
	if claims.RepositoryID != c.repositoryID || claims.Repository != c.repository {
		return errors.New("repository claims do not match")
	}
	if claims.Ref != c.ref {
		return errors.New("ref claim does not match")
	}
	if claims.JobWorkflowRef != c.workflowRef {
		return errors.New("job_workflow_ref claim does not match")
	}
	if claims.Subject == "" {
		return errors.New("sub claim is required")
	}
	now := c.now()
	if claims.ExpiresAt == 0 || !now.Before(time.Unix(claims.ExpiresAt, 0).Add(c.clockSkew)) {
		return errors.New("token is expired")
	}
	if claims.NotBefore != nil && now.Add(c.clockSkew).Before(time.Unix(*claims.NotBefore, 0)) {
		return errors.New("token is not valid yet")
	}
	if claims.IssuedAt == 0 || now.Add(c.clockSkew).Before(time.Unix(claims.IssuedAt, 0)) {
		return errors.New("iat claim is invalid")
	}
	if claims.ExpiresAt <= claims.IssuedAt {
		return errors.New("exp claim must be later than iat")
	}
	return nil
}

func (c *Client) signingKey(ctx context.Context, keyID string) (*rsa.PublicKey, error) {
	c.keysMu.Lock()
	defer c.keysMu.Unlock()
	if c.keys != nil && c.now().Before(c.keysExpireAt) {
		if key := c.keys[keyID]; key != nil {
			return key, nil
		}
	}
	keys, err := c.fetchKeys(ctx)
	if err != nil {
		return nil, &oidcUpstreamError{cause: err}
	}
	c.keys = keys
	c.keysExpireAt = c.now().Add(c.cacheTTL)
	key := keys[keyID]
	if key == nil {
		return nil, errors.New("JWT signing key was not found")
	}
	return key, nil
}

func (c *Client) fetchKeys(ctx context.Context) (map[string]*rsa.PublicKey, error) {
	var discovery discoveryDocument
	if err := c.getExternalJSON(ctx, c.discoveryURL, &discovery); err != nil {
		return nil, fmt.Errorf("fetch OIDC discovery document: %w", err)
	}
	if strings.TrimRight(discovery.Issuer, "/") != c.issuer {
		return nil, errors.New("OIDC discovery issuer does not match")
	}
	jwksURL, err := url.Parse(discovery.JWKSURI)
	if err != nil || jwksURL.Scheme == "" || jwksURL.Host == "" {
		return nil, errors.New("OIDC discovery jwks_uri is invalid")
	}
	discoveryURL, _ := url.Parse(c.discoveryURL)
	if jwksURL.Scheme != discoveryURL.Scheme || jwksURL.Host != discoveryURL.Host {
		return nil, errors.New("OIDC jwks_uri must use the discovery origin")
	}
	var document jwksDocument
	if err := c.getExternalJSON(ctx, jwksURL.String(), &document); err != nil {
		return nil, fmt.Errorf("fetch OIDC JWK set: %w", err)
	}
	keys := make(map[string]*rsa.PublicKey)
	for _, jwk := range document.Keys {
		if jwk.KeyType != "RSA" || jwk.KeyID == "" || (jwk.Use != "" && jwk.Use != "sig") || (jwk.Algorithm != "" && jwk.Algorithm != "RS256") {
			continue
		}
		key, err := rsaKey(jwk)
		if err != nil {
			continue
		}
		keys[jwk.KeyID] = key
	}
	if len(keys) == 0 {
		return nil, errors.New("OIDC JWK set contains no usable RS256 key")
	}
	return keys, nil
}

func (c *Client) getExternalJSON(ctx context.Context, endpoint string, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", response.StatusCode)
	}
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxOIDCDocumentBytes))
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return requireJSONEOF(decoder)
}

func parseJWT(token string) (tokenHeader, verificationClaims, string, []byte, error) {
	if token == "" || len(token) > maxJWTBytes {
		return tokenHeader{}, verificationClaims{}, "", nil, errors.New("JWT length is invalid")
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return tokenHeader{}, verificationClaims{}, "", nil, errors.New("JWT compact serialization is invalid")
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return tokenHeader{}, verificationClaims{}, "", nil, errors.New("JWT header encoding is invalid")
	}
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return tokenHeader{}, verificationClaims{}, "", nil, errors.New("JWT claims encoding is invalid")
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return tokenHeader{}, verificationClaims{}, "", nil, errors.New("JWT signature encoding is invalid")
	}
	var header tokenHeader
	if err := decodeJSONBytes(headerBytes, &header); err != nil {
		return tokenHeader{}, verificationClaims{}, "", nil, fmt.Errorf("JWT header: %w", err)
	}
	var claims verificationClaims
	if err := decodeJSONBytes(claimsBytes, &claims); err != nil {
		return tokenHeader{}, verificationClaims{}, "", nil, fmt.Errorf("JWT claims: %w", err)
	}
	return header, claims, parts[0] + "." + parts[1], signature, nil
}

func rsaKey(jwk jsonWebKey) (*rsa.PublicKey, error) {
	modulus, err := base64.RawURLEncoding.DecodeString(jwk.Modulus)
	if err != nil || len(modulus) == 0 {
		return nil, errors.New("invalid RSA modulus")
	}
	exponentBytes, err := base64.RawURLEncoding.DecodeString(jwk.Exponent)
	if err != nil || len(exponentBytes) == 0 || len(exponentBytes) > 4 {
		return nil, errors.New("invalid RSA exponent")
	}
	exponent := 0
	for _, value := range exponentBytes {
		exponent = exponent<<8 | int(value)
	}
	if exponent < 3 || exponent%2 == 0 {
		return nil, errors.New("invalid RSA exponent")
	}
	key := &rsa.PublicKey{N: new(big.Int).SetBytes(modulus), E: exponent}
	if key.N.BitLen() < 2048 {
		return nil, errors.New("RSA modulus must be at least 2048 bits")
	}
	return key, nil
}

func decodeJSONBytes(data []byte, target any) error {
	if err := rejectDuplicateKeys(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return requireJSONEOF(decoder)
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := walkJSONValue(decoder); err != nil {
		return err
	}
	return requireJSONEOF(decoder)
}

func walkJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		seen := make(map[string]struct{})
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return errors.New("JSON object key is not a string")
			}
			if _, exists := seen[key]; exists {
				return fmt.Errorf("duplicate JSON object key %q", key)
			}
			seen[key] = struct{}{}
			if err := walkJSONValue(decoder); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim('}') {
			return errors.New("JSON object is not closed")
		}
	case '[':
		for decoder.More() {
			if err := walkJSONValue(decoder); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim(']') {
			return errors.New("JSON array is not closed")
		}
	default:
		return errors.New("unexpected JSON delimiter")
	}
	return nil
}

func requireJSONEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("unexpected trailing JSON value")
		}
		return err
	}
	return nil
}

func invalidToken(err error) error {
	return fmt.Errorf("%w: %v", ErrInvalidActionsToken, err)
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func validateEndpoint(value string, allowInsecure bool) (*url.URL, error) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil {
		return nil, errors.New("must be an absolute URL without user information")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("may not contain a query or fragment")
	}
	if parsed.Scheme != "https" && !(allowInsecure && parsed.Scheme == "http") {
		return nil, errors.New("must use HTTPS")
	}
	return parsed, nil
}
