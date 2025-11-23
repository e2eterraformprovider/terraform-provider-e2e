package vpc

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testVPCPrefix = "test-vpc-"

func init() {
	resource.AddTestSweepers("e2e_vpc", &resource.Sweeper{
		Name:         "e2e_vpc",
		F:            sweepVPCs,
		Dependencies: []string{
			// Add dependencies if needed (e.g., nodes, load balancers)
		},
	})
}

func sweepVPCs(region string) error {
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

	// Get list of VPCs
	response, err := client.GetVpcs(location, projectID)
	if err != nil {
		return fmt.Errorf("error listing VPCs: %w", err)
	}

	if response == nil {
		log.Printf("[WARNING] No VPCs found")
		return nil
	}

	// Parse the response
	data, ok := response["data"].([]interface{})
	if !ok {
		log.Printf("[WARNING] No VPC data found")
		return nil
	}

	log.Printf("[DEBUG] Found %d VPCs in total", len(data))

	sweptCount := 0
	for _, item := range data {
		vpc := item.(map[string]interface{})
		vpcName := vpc["name"].(string)

		if !strings.HasPrefix(vpcName, testVPCPrefix) {
			log.Printf("[DEBUG] Skipping VPC %s (does not have test prefix)", vpcName)
			continue
		}

		vpcIDFloat := vpc["network_id"].(float64)
		vpcID := fmt.Sprintf("%.0f", vpcIDFloat)
		log.Printf("[INFO] Deleting VPC: %s (ID: %s)", vpcName, vpcID)

		_, err := client.DeleteVpc(vpcID, projectID, location)
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPC %s: %v", vpcName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted VPC: %s", vpcName)
	}

	log.Printf("[INFO] Swept %d VPCs", sweptCount)
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
