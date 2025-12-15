package config

import (
	"fmt"
	"log"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Config holds the provider configuration and manages all API clients.
// This follows the DigitalOcean provider pattern for clean client management.
type Config struct {
	// Provider configuration
	APIKey      string
	AuthToken   string
	APIEndpoint string

	// Provider-level defaults
	DefaultRegion    string
	DefaultProjectID string

	// Retry configuration for goe2e client
	HTTPRetryMax      int
	HTTPRetryWaitMin  int // seconds
	HTTPRetryWaitMax  int // seconds
	RequestsPerSecond float64

	// Internal clients
	goe2eClient *goe2e.Client
}

// Goe2eClient returns the new goe2e client for migrated services.
// Currently used by FaaS resources, will be used by more services as migration progresses.
func (c *Config) Goe2eClient() *goe2e.Client {
	return c.goe2eClient
}

// SetGoe2eClientForTesting sets the goe2e client for testing purposes.
// This should only be used in test files to inject mock clients.
func (c *Config) SetGoe2eClientForTesting(client *goe2e.Client) {
	c.goe2eClient = client
}

// Goe2eClientForProject creates a new goe2e client with specific projectID and region.
// This is used by resources that need to use different projectID/region than the default.
// The client is created with the same configuration options (retry, rate limiting, etc.)
// as the default client, but with the specified projectID and region.
func (c *Config) Goe2eClientForProject(projectID, region string) (*goe2e.Client, error) {
	// If a test client is set, return it instead of creating a new one
	if c.goe2eClient != nil {
		return c.goe2eClient, nil
	}

	var opts []goe2e.ClientOpt

	// Set custom base URL if provided
	if c.APIEndpoint != "" && c.APIEndpoint != "https://api.e2enetworks.com/myaccount/api/v1/" {
		opts = append(opts, goe2e.SetBaseURL(c.APIEndpoint))
	}

	// Configure retry with backoff
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

	// Create the goe2e client with the specified projectID and region
	return goe2e.NewClient(c.APIKey, c.AuthToken, projectID, region, opts...)
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

	// Create the goe2e client with placeholder region/project (will be overridden per-request)
	// Note: projectID and region are required by NewClient but are typically set at resource level
	// Using defaults here; resources should create clients with their specific values
	projectID := c.DefaultProjectID
	if projectID == "" {
		projectID = "default"
	}
	region := c.DefaultRegion
	if region == "" {
		region = "default"
	}

	goe2eClient, err := goe2e.NewClient(c.APIKey, c.AuthToken, projectID, region, opts...)
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

// GetRegionOrLocation is a helper function to handle the region/location parameter migration.
// It prefers 'region' but falls back to 'location' with a deprecation warning.
// This centralizes the migration logic to avoid code duplication across all resources.
//
// Usage in resource CRUD functions:
//
//	region, err := config.GetRegionOrLocation(d)
//	if err != nil {
//	    return diag.FromErr(err)
//	}
func GetRegionOrLocation(d *schema.ResourceData) (string, error) {
	// Prefer 'region' parameter
	if v, ok := d.GetOk("region"); ok {
		return v.(string), nil
	}

	// Fall back to 'location' with deprecation warning
	if v, ok := d.GetOk("location"); ok {
		log.Printf("[WARN] Parameter 'location' is deprecated and will be removed in v3.0.0. Please use 'region' instead")
		return v.(string), nil
	}

	// Neither parameter provided - this should be caught by schema validation
	// but we include it for safety
	return "", fmt.Errorf("either 'region' or 'location' must be specified")
}

// GetRegionOrDefault returns region from resource or provider default.
// This integrates with the existing GetRegionOrLocation() function to support
// both the region/location migration and provider-level defaults.
//
// Priority order:
// 1. Resource-level 'region' parameter
// 2. Resource-level 'location' parameter (deprecated)
// 3. Provider-level default_region
// 4. Error if none specified
//
// Usage in resource CRUD functions:
//
//	cfg := m.(*config.Config)
//	region, err := cfg.GetRegionOrDefault(d)
//	if err != nil {
//	    return diag.FromErr(err)
//	}
func (c *Config) GetRegionOrDefault(d *schema.ResourceData) (string, error) {
	// First check resource-level parameters (using existing helper logic)
	if v, ok := d.GetOk("region"); ok {
		return v.(string), nil
	}
	if v, ok := d.GetOk("location"); ok {
		log.Printf("[WARN] Parameter 'location' is deprecated and will be removed in v3.0.0. Please use 'region' instead")
		return v.(string), nil
	}

	// Fall back to provider default
	if c.DefaultRegion != "" {
		return c.DefaultRegion, nil
	}

	return "", fmt.Errorf("region must be specified (either in resource 'region'/'location' parameter or provider 'default_region' parameter/E2E_REGION environment variable)")
}

// GetProjectIDOrDefault returns project_id from resource or provider default.
//
// Priority order:
// 1. Resource-level 'project_id' parameter
// 2. Provider-level default_project_id
// 3. Error if none specified
//
// Usage in resource CRUD functions:
//
//	cfg := m.(*config.Config)
//	projectID, err := cfg.GetProjectIDOrDefault(d)
//	if err != nil {
//	    return diag.FromErr(err)
//	}
func (c *Config) GetProjectIDOrDefault(d *schema.ResourceData) (string, error) {
	if v, ok := d.GetOk("project_id"); ok {
		return v.(string), nil
	}

	if c.DefaultProjectID != "" {
		return c.DefaultProjectID, nil
	}

	return "", fmt.Errorf("project_id must be specified (either in resource 'project_id' parameter or provider 'default_project_id' parameter/E2E_PROJECT_ID environment variable)")
}
