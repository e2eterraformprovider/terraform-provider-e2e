package security_group

import (
	"fmt"
	"log"
	"os"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testSecurityGroupPrefix = "test-sg-"

func init() {
	resource.AddTestSweepers("e2e_security_group", &resource.Sweeper{
		Name: "e2e_security_group",
		F:    sweepSecurityGroups,
	})
}

func sweepSecurityGroups(region string) error {
	cfg, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting config for region %s: %w", region, err)
	}

	_ = cfg.Client()

	// Get test project ID and location from environment
	projectID := os.Getenv("E2E_TEST_PROJECT_ID")
	location := os.Getenv("E2E_TEST_LOCATION")

	if projectID == "" || location == "" {
		log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_LOCATION not set, skipping sweep")
		return nil
	}

	// Note: We would need a ListSecurityGroups method in the client
	// For now, we'll skip the sweep if the method doesn't exist
	// This is a placeholder that matches the pattern
	log.Printf("[INFO] Sweeping security groups is not fully implemented yet")
	log.Printf("[INFO] Security groups with prefix %s should be manually cleaned up if needed", testSecurityGroupPrefix)

	// Placeholder implementation:
	// response, err := client.ListSecurityGroups(projectID, location)
	// if err != nil {
	//     return fmt.Errorf("error listing security groups: %w", err)
	// }

	// sweptCount := 0
	// for _, sg := range response {
	//     if !strings.HasPrefix(sg.Name, testSecurityGroupPrefix) {
	//         continue
	//     }
	//     err := client.DeleteSecurityGroup(sg.ID, projectID, location)
	//     if err != nil {
	//         log.Printf("[ERROR] Failed to delete security group %s: %v", sg.Name, err)
	//         continue
	//     }
	//     sweptCount++
	// }
	// log.Printf("[INFO] Swept %d security groups", sweptCount)

	return nil
}

// sharedConfigForRegion returns a common config for the region
func sharedConfigForRegion(region string) (*config.Config, error) {
	apiKey := os.Getenv("SERVICE_API_KEY")
	authToken := os.Getenv("SERVICE_AUTH_TOKEN")
	apiEndpoint := os.Getenv("SERVICE_API_ENDPOINT")

	if apiKey == "" {
		return nil, fmt.Errorf("SERVICE_API_KEY must be set for acceptance tests")
	}

	if authToken == "" {
		return nil, fmt.Errorf("SERVICE_AUTH_TOKEN must be set for acceptance tests")
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
