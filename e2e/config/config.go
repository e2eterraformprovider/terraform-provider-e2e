package config

import (
	"log"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

// Config holds the provider configuration and manages all API clients.
// This follows the DigitalOcean provider pattern for clean client management.
type Config struct {
	// Provider configuration
	APIKey      string
	AuthToken   string
	APIEndpoint string

	// Retry configuration for goe2e client
	HTTPRetryMax      int
	HTTPRetryWaitMin  int // seconds
	HTTPRetryWaitMax  int // seconds
	RequestsPerSecond float64

	// Internal clients
	legacyClient *client.Client
	goe2eClient  *goe2e.Client
}

// Client returns the legacy client for existing resources.
// This is used by all non-FaaS resources during the migration period.
func (c *Config) Client() *client.Client {
	return c.legacyClient
}

// Goe2eClient returns the new goe2e client for migrated services.
// Currently used by FaaS resources, will be used by more services as migration progresses.
func (c *Config) Goe2eClient() *goe2e.Client {
	return c.goe2eClient
}

// NewConfig creates and initializes a Config with both clients.
// It sets up the legacy client and the goe2e client with retry configuration.
func NewConfig(apiKey, authToken, apiEndpoint string) (*Config, error) {
	cfg := &Config{
		APIKey:      apiKey,
		AuthToken:   authToken,
		APIEndpoint: apiEndpoint,

		// Default retry configuration
		HTTPRetryMax:     4,
		HTTPRetryWaitMin: 1,
		HTTPRetryWaitMax: 30,
	}

	// Initialize legacy client
	cfg.legacyClient = client.NewClient(apiKey, authToken, apiEndpoint)

	// Initialize goe2e client with retry configuration
	if err := cfg.initGoe2eClient(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// initGoe2eClient initializes the goe2e client with retry and rate limiting.
// This follows DigitalOcean's pattern for configuring the client with functional options.
func (c *Config) initGoe2eClient() error {
	var opts []goe2e.ClientOpt

	// Set custom base URL if provided
	if c.APIEndpoint != "" && c.APIEndpoint != "https://api.e2enetworks.com/myaccount/api/v1/" {
		opts = append(opts, goe2e.SetBaseURL(c.APIEndpoint))
	}

	// Configure retry with backoff (similar to DigitalOcean's pattern)
	if c.HTTPRetryMax > 0 {
		retryConfig := goe2e.RetryConfig{
			RetryMax:     c.HTTPRetryMax,
			RetryWaitMin: goe2e.PtrTo(time.Duration(c.HTTPRetryWaitMin) * time.Second),
			RetryWaitMax: goe2e.PtrTo(time.Duration(c.HTTPRetryWaitMax) * time.Second),
		}
		opts = append(opts, goe2e.WithRetryAndBackoffs(retryConfig))
	}

	// Set user agent
	userAgent := "terraform-provider-e2e/1.0"
	opts = append(opts, goe2e.SetUserAgent(userAgent))

	// Configure rate limiting if specified
	if c.RequestsPerSecond > 0.0 {
		opts = append(opts, goe2e.SetStaticRateLimit(c.RequestsPerSecond))
	}

	// Create the goe2e client
	goe2eClient, err := goe2e.NewClient(c.APIKey, c.AuthToken, opts...)
	if err != nil {
		return err
	}

	// TODO: Add logging transport similar to DO's pattern
	// This would allow logging all API requests for debugging
	// clientTransport := logging.NewTransport("E2E", goe2eClient.HTTPClient.Transport)
	// goe2eClient.HTTPClient.Transport = clientTransport

	c.goe2eClient = goe2eClient

	// Log successful initialization
	log.Printf("[INFO] E2E provider clients initialized successfully")

	return nil
}
