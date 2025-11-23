package blockstorage

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
	cfg, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting config for region %s: %w", region, err)
	}

	client := cfg.Client()

	// Get test project ID and location from environment
	projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	location := os.Getenv("E2E_TEST_LOCATION")

	if projectIDStr == "" || location == "" {
		log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_LOCATION not set, skipping sweep")
		return nil
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		return fmt.Errorf("error converting project ID to int: %w", err)
	}

	// Get list of block storages
	response, err := client.GetBlockStorages(projectID, location)
	if err != nil {
		return fmt.Errorf("error listing block storages: %w", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		log.Printf("[WARNING] No block storages found or invalid response format")
		return nil
	}

	log.Printf("[DEBUG] Found %d block storages in total", len(data))

	sweptCount := 0
	for _, item := range data {
		blockStorage, ok := item.(map[string]interface{})
		if !ok {
			log.Printf("[WARNING] Invalid block storage format, skipping")
			continue
		}

		blockStorageName, ok := blockStorage["name"].(string)
		if !ok {
			log.Printf("[WARNING] Block storage name not found, skipping")
			continue
		}

		if !strings.HasPrefix(blockStorageName, testBlockStoragePrefix) {
			log.Printf("[DEBUG] Skipping block storage %s (does not have test prefix)", blockStorageName)
			continue
		}

		// Check if block storage is attached, skip if it is
		status, ok := blockStorage["status"].(string)
		if ok && status == "Attached" {
			log.Printf("[INFO] Skipping attached block storage: %s", blockStorageName)
			continue
		}

		blockStorageID, ok := blockStorage["id"].(float64)
		if !ok {
			log.Printf("[WARNING] Invalid block storage ID format for %s, skipping", blockStorageName)
			continue
		}

		blockStorageIDStr := fmt.Sprintf("%.0f", blockStorageID)
		log.Printf("[INFO] Deleting block storage: %s (ID: %s)", blockStorageName, blockStorageIDStr)

		err := client.DeleteBlockStorage(blockStorageIDStr, projectID, location)
		if err != nil {
			log.Printf("[ERROR] Failed to delete block storage %s: %v", blockStorageName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted block storage: %s", blockStorageName)
	}

	log.Printf("[INFO] Swept %d block storages", sweptCount)
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
