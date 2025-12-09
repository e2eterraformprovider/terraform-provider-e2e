package sweep

import (
	"fmt"
	"os"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

const TestNamePrefix = "test-"

// TestProjectID and TestRegion are loaded from environment variables for use in sweepers
var (
	TestProjectID string
	TestRegion    string
)

func init() {
	// Load test environment variables once during init
	TestProjectID = os.Getenv("E2E_TEST_PROJECT_ID")
	TestRegion = os.Getenv("E2E_TEST_REGION")
}

func SharedConfigForRegion(region string) (*config.Config, error) {
	apiKey := os.Getenv("E2E_API_KEY")
	authToken := os.Getenv("E2E_AUTH_TOKEN")
	apiEndpoint := os.Getenv("E2E_API_ENDPOINT")

	if apiKey == "" {
		return nil, fmt.Errorf("E2E_API_KEY must be set for acceptance tests")
	}

	if authToken == "" {
		return nil, fmt.Errorf("E2E_AUTH_TOKEN must be set for acceptance tests")
	}

	if apiEndpoint == "" {
		apiEndpoint = "https://api.e2enetworks.com/myaccount/api/v1/"
	}

	cfg, err := config.NewConfig(apiKey, authToken, apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating config: %w", err)
	}

	return cfg, nil
}

// SharedGoe2eClientForTests returns a configured goe2e client for sweeper tests.
// It handles all setup: validating required environment variables (E2E_TEST_PROJECT_ID, E2E_TEST_REGION)
// and creating a client configured for the test project and region.
//
// Returns:
//   - goe2e.Client: Ready-to-use client configured for test project/region
//   - error: If required env vars are not set or client creation fails
//
// Usage in sweepers:
//
//	client, err := sweep.SharedGoe2eClientForTests()
//	if err != nil {
//		log.Printf("[WARNING] %v - skipping sweep", err)
//		return nil
//	}
//	// Use client directly
func SharedGoe2eClientForTests() (*goe2e.Client, error) {
	// Validate required environment variables
	if TestProjectID == "" {
		return nil, fmt.Errorf("E2E_TEST_PROJECT_ID not set")
	}
	if TestRegion == "" {
		return nil, fmt.Errorf("E2E_TEST_REGION not set")
	}

	// Get base config
	cfg, err := SharedConfigForRegion(TestRegion)
	if err != nil {
		return nil, fmt.Errorf("error creating base config: %w", err)
	}

	// Create client for specific project and region
	client, err := cfg.Goe2eClientForProject(TestProjectID, TestRegion)
	if err != nil {
		return nil, fmt.Errorf("error creating goe2e client: %w", err)
	}

	return client, nil
}
