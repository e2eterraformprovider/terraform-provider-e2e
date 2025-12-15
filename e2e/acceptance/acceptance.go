package acceptance

import (
	"context"
	"os"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const TestNamePrefix = "test-"

var (
	TestAccProvider          *schema.Provider
	TestAccProviders         map[string]*schema.Provider
	TestAccProviderFactories map[string]func() (*schema.Provider, error)

	// Test environment variables - loaded once and reused across all tests
	TestProjectID string
	TestRegion    string
)

func init() {
	TestAccProvider = e2e.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"e2e": TestAccProvider,
	}
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"e2e": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}

func TestAccPreCheck(t *testing.T) {
	// Skip acceptance tests unless TF_ACC is explicitly set to "1"
	// This prevents accidental execution when TF_ACC=0 or other values are set
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Skipping acceptance test. Set TF_ACC=1 to run acceptance tests.")
	}

	if v := os.Getenv("E2E_API_KEY"); v == "" {
		t.Fatal("E2E_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_AUTH_TOKEN"); v == "" {
		t.Fatal("E2E_AUTH_TOKEN must be set for acceptance tests")
	}

	// Load test environment variables once
	// These are used for test-specific values (e.g., in assertions)
	TestProjectID = os.Getenv("E2E_TEST_PROJECT_ID")
	if TestProjectID == "" {
		t.Fatal("E2E_TEST_PROJECT_ID must be set for acceptance tests")
	}

	TestRegion = os.Getenv("E2E_TEST_REGION")
	if TestRegion == "" {
		t.Fatal("E2E_TEST_REGION must be set for acceptance tests")
	}

	// Also ensure provider-level defaults are set
	// These are used by the provider to set default project_id/region on resources
	if v := os.Getenv("E2E_PROJECT_ID"); v == "" {
		// Set it from test vars if not already set
		_ = os.Setenv("E2E_PROJECT_ID", TestProjectID)
	}
	if v := os.Getenv("E2E_REGION"); v == "" {
		// Set it from test vars if not already set
		_ = os.Setenv("E2E_REGION", TestRegion)
	}

	err := TestAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}

// TestAccPreCheckSyntaxOnly is a PreCheck function for syntax/argument validation tests
// that don't require API calls. These tests use ExpectError to validate Terraform
// syntax/argument checking and don't need credentials.
//
// Usage: Use this PreCheck for tests that only validate syntax/arguments and don't
// make actual API calls (e.g., tests with ExpectError for missing required arguments).
func TestAccPreCheckSyntaxOnly(t *testing.T) {
	// Syntax validation tests can run without TF_ACC since they don't make API calls
	// They only validate Terraform configuration syntax/arguments
	// No credentials needed - Terraform validates before any API calls
}

// GetRegionOrLocationFromState is a helper function to handle the region/location parameter migration in tests.
// It prefers 'region' but falls back to 'location' for backwards compatibility.
// This mirrors the logic in config.GetRegionOrLocation() but works with terraform.ResourceState.
//
// Usage in test functions:
//
//	region := acceptance.GetRegionOrLocationFromState(rs)
func GetRegionOrLocationFromState(rs *terraform.ResourceState) string {
	// Prefer 'region' parameter
	if region := rs.Primary.Attributes["region"]; region != "" {
		return region
	}

	// Fall back to 'location' for backwards compatibility
	return rs.Primary.Attributes["location"]
}
