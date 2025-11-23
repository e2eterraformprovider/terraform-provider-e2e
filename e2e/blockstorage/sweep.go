package blockstorage

import (
	"fmt"
	"log"
	"os"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testBlockStoragePrefix = "test-bs-"

func init() {
	resource.AddTestSweepers("e2e_blockstorage", &resource.Sweeper{
		Name: "e2e_blockstorage",
		F:    sweepBlockStorages,
		Dependencies: []string{
			// Block storages should be swept after nodes if there are any dependencies
			"e2e_node",
		},
	})
}

func sweepBlockStorages(region string) error {
	// NOTE: The E2E API client does not currently have a method to list all block storages.
	// The GetBlockStorage method only retrieves a single block storage by ID.
	// Until a ListBlockStorages or GetBlockStorages method is implemented in the client,
	// automatic sweeping of test block storages is not possible.
	// Test block storages will need to be cleaned up manually or via the destroy functionality
	// in the test framework.

	log.Printf("[INFO] Block storage sweeping not yet implemented - waiting for ListBlockStorages API method")
	log.Printf("[INFO] Please manually clean up test block storages with prefix '%s' if needed", testBlockStoragePrefix)

	// Uncomment and update the code below once the API client has a list method:
	//
	// cfg, err := sharedConfigForRegion(region)
	// if err != nil {
	// 	return fmt.Errorf("error getting config for region %s: %w", region, err)
	// }
	//
	// client := cfg.Client()
	// projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	// location := os.Getenv("E2E_TEST_LOCATION")
	//
	// if projectIDStr == "" || location == "" {
	// 	log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_LOCATION not set, skipping sweep")
	// 	return nil
	// }
	//
	// projectID, err := strconv.Atoi(projectIDStr)
	// if err != nil {
	// 	return fmt.Errorf("error converting project ID to int: %w", err)
	// }
	//
	// // Get list of block storages
	// response, err := client.ListBlockStorages(projectID, location)
	// if err != nil {
	// 	return fmt.Errorf("error listing block storages: %w", err)
	// }
	//
	// ... rest of sweep logic

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
