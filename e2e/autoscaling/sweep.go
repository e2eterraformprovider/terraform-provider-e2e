package autoscaling

import (
	"fmt"
	"log"
	"os"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = "test-sg-"

func init() {
	resource.AddTestSweepers("e2e_scaler_group", &resource.Sweeper{
		Name: "e2e_scaler_group",
		F:    sweepScalerGroups,
	})
}

func sweepScalerGroups(region string) error {
	cfg, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting config for region %s: %w", region, err)
	}

	_ = cfg.Client()

	projectID := os.Getenv("E2E_TEST_PROJECT_ID")
	location := os.Getenv("E2E_TEST_LOCATION")

	if projectID == "" || location == "" {
		log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_LOCATION not set, skipping sweep")
		return nil
	}

	// Note: This assumes there's a method to list all scaler groups
	// If such a method doesn't exist, we can't sweep them automatically
	log.Printf("[INFO] Scaler group sweeping would require list API method")
	log.Printf("[INFO] Please manually delete scaler groups with prefix '%s' if needed", testNamePrefix)

	// TODO: Implement when ListScalerGroups API is available
	// groups, err := client.ListScalerGroups(projectID, location)
	// if err != nil {
	// 	return fmt.Errorf("error listing scaler groups: %w", err)
	// }
	//
	// log.Printf("[DEBUG] Found %d scaler groups in total", len(groups))
	//
	// sweptCount := 0
	// for _, group := range groups {
	// 	if !strings.HasPrefix(group.Name, testNamePrefix) {
	// 		log.Printf("[DEBUG] Skipping scaler group %s (does not have test prefix)", group.Name)
	// 		continue
	// 	}
	//
	// 	log.Printf("[INFO] Deleting scaler group: %s (ID: %s)", group.Name, group.ID)
	//
	// 	err := client.DeleteScalerGroup(group.ID, projectID, location)
	// 	if err != nil {
	// 		log.Printf("[ERROR] Failed to delete scaler group %s: %v", group.Name, err)
	// 		continue
	// 	}
	//
	// 	sweptCount++
	// 	log.Printf("[INFO] Successfully deleted scaler group: %s", group.Name)
	// }
	//
	// log.Printf("[INFO] Swept %d scaler groups", sweptCount)

	return nil
}

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
