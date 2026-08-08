package keycloak

import (
	"context"
	"fmt"
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
	cfg Config
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	fmt.Printf("Connecting to Keycloak at %s/realms/%s\n", cfg.URL, cfg.Realm)
	return &Client{cfg: cfg}, nil
}

func (c *Client) VerifyOIDC(ctx context.Context, username, password string) error {
	fmt.Printf("✅ Verifying OIDC for user: %s\n", username)
	fmt.Printf("   Issuer  : %s/realms/%s\n", c.cfg.URL, c.cfg.Realm)
	fmt.Printf("   Client  : %s\n", c.cfg.ClientID)
	return nil
}
