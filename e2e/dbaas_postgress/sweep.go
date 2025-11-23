package dbaas_postgress

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testPostgresDBaaSPrefix = "test-pg-"

func init() {
	resource.AddTestSweepers("e2e_dbaas_postgress", &resource.Sweeper{
		Name:         "e2e_dbaas_postgress",
		F:            sweepPostgresDBaaS,
		Dependencies: []string{
			// Add dependencies if needed
		},
	})
}

func sweepPostgresDBaaS(region string) error {
	cfg, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting config for region %s: %w", region, err)
	}

	client := cfg.Client()

	// Get test project ID and location from environment
	projectID := os.Getenv("E2E_TEST_PROJECT_ID")
	location := os.Getenv("E2E_TEST_LOCATION")

	if projectID == "" || location == "" {
		log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_LOCATION not set, skipping sweep")
		return nil
	}

	// Get list of DBaaS instances
	response, err := client.GetPostgresDBList(location, projectID)
	if err != nil {
		return fmt.Errorf("error listing Postgres DBaaS instances: %w", err)
	}

	if response == nil {
		log.Printf("[WARNING] No Postgres DBaaS instances found")
		return nil
	}

	// Parse the response
	data, ok := response["data"].([]interface{})
	if !ok {
		log.Printf("[WARNING] No Postgres DBaaS data found")
		return nil
	}

	log.Printf("[DEBUG] Found %d Postgres DBaaS instances in total", len(data))

	sweptCount := 0
	for _, item := range data {
		dbaas := item.(map[string]interface{})
		dbaasName := dbaas["name"].(string)

		if !strings.HasPrefix(dbaasName, testPostgresDBaaSPrefix) {
			log.Printf("[DEBUG] Skipping Postgres DBaaS %s (does not have test prefix)", dbaasName)
			continue
		}

		dbaasIDFloat := dbaas["id"].(float64)
		dbaasID := fmt.Sprintf("%.0f", dbaasIDFloat)
		log.Printf("[INFO] Deleting Postgres DBaaS: %s (ID: %s)", dbaasName, dbaasID)

		err := client.DeletePostgressDB(dbaasID, projectID, location)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Postgres DBaaS %s: %v", dbaasName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted Postgres DBaaS: %s", dbaasName)
	}

	log.Printf("[INFO] Swept %d Postgres DBaaS instances", sweptCount)
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
