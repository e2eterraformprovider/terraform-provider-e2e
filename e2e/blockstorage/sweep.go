package blockstorage

import (
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testBlockStoragePrefix = sweep.TestNamePrefix + "bs-"

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
	// cfg, err := sweep.SharedConfigForRegion(region)
	// if err != nil {
	// 	return fmt.Errorf("error getting config for region %s: %w", region, err)
	// }
	//
	// client := cfg.Client()
	// projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	// location := os.Getenv("E2E_TEST_REGION")
	//
	// if projectIDStr == "" || location == "" {
	// 	log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_REGION not set, skipping sweep")
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
