package loadbalancer

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testLoadBalancerPrefix = "test-lb-"

func init() {
	resource.AddTestSweepers("e2e_loadbalancer", &resource.Sweeper{
		Name:         "e2e_loadbalancer",
		F:            sweepLoadBalancers,
		Dependencies: []string{
			// Add dependencies if needed
		},
	})
}

func sweepLoadBalancers(region string) error {
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

	// Get list of load balancers
	response, err := client.GetLoadBalancers(location, projectID)
	if err != nil {
		return fmt.Errorf("error listing load balancers: %w", err)
	}

	if response == nil {
		log.Printf("[WARNING] No load balancers found")
		return nil
	}

	// Parse the response
	data, ok := response["data"].([]interface{})
	if !ok {
		log.Printf("[WARNING] No load balancer data found")
		return nil
	}

	log.Printf("[DEBUG] Found %d load balancers in total", len(data))

	sweptCount := 0
	for _, item := range data {
		lb := item.(map[string]interface{})
		lbName := lb["name"].(string)

		if !strings.HasPrefix(lbName, testLoadBalancerPrefix) {
			log.Printf("[DEBUG] Skipping load balancer %s (does not have test prefix)", lbName)
			continue
		}

		lbIDFloat := lb["id"].(float64)
		lbID := fmt.Sprintf("%.0f", lbIDFloat)
		log.Printf("[INFO] Deleting load balancer: %s (ID: %s)", lbName, lbID)

		err := client.DeleteLoadBalancer(lbID, location, projectID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete load balancer %s: %v", lbName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted load balancer: %s", lbName)
	}

	log.Printf("[INFO] Swept %d load balancers", sweptCount)
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
