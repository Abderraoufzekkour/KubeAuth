package keycloak

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
	AdminUser    string
	AdminPass    string
}

type Client struct {
	cfg        Config
	httpClient *http.Client
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenClaims struct {
	Username string   `json:"preferred_username"`
	Email    string   `json:"email"`
	Groups   []string `json:"groups"`
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *Client) VerifyOIDC(ctx context.Context, username, password string) error {
	fmt.Printf("Connecting to Keycloak at %s/realms/%s\n", c.cfg.URL, c.cfg.Realm)

	token, err := c.fetchToken(ctx, username, password)
	if err != nil {
		return fmt.Errorf("failed to fetch token: %w", err)
	}

	fmt.Println("Token fetched successfully.")

	claims, err := c.decodeToken(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to decode token: %w", err)
	}

	fmt.Println("-------------------------------------------")
	fmt.Printf("Username : %s\n", claims.Username)
	fmt.Printf("Email    : %s\n", claims.Email)
	fmt.Printf("Groups   : %v\n", claims.Groups)
	fmt.Println("-------------------------------------------")
	fmt.Println("Suggested kube-apiserver flags:")
	fmt.Printf("  --oidc-issuer-url=%s/realms/%s\n", c.cfg.URL, c.cfg.Realm)
	fmt.Printf("  --oidc-client-id=%s\n", c.cfg.ClientID)
	fmt.Printf("  --oidc-username-claim=preferred_username\n")
	fmt.Printf("  --oidc-groups-claim=groups\n")
	fmt.Printf("  --oidc-username-prefix=oidc:\n")
	fmt.Printf("  --oidc-groups-prefix=oidc:\n")

	return nil
}

func (c *Client) fetchToken(ctx context.Context, username, password string) (*TokenResponse, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		c.cfg.URL, c.cfg.Realm)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", c.cfg.ClientID)
	data.Set("username", username)
	data.Set("password", password)
	if c.cfg.ClientSecret != "" {
		data.Set("client_secret", c.cfg.ClientSecret)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not reach Keycloak: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak returned %d: %s", resp.StatusCode, string(body))
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *Client) decodeToken(accessToken string) (*TokenClaims, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	payload := parts[1]
	payload = strings.ReplaceAll(payload, "-", "+")
	payload = strings.ReplaceAll(payload, "_", "/")
	padding := 4 - len(payload)%4
	if padding != 4 {
		payload += strings.Repeat("=", padding)
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	var claims TokenClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &claims, nil
}
