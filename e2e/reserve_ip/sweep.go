package reserve_ip

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testReserveIPPrefix = "test-reserve-ip-"

func init() {
	resource.AddTestSweepers("e2e_reserve_ip", &resource.Sweeper{
		Name: "e2e_reserve_ip",
		F:    sweepReserveIPs,
	})
}

func sweepReserveIPs(region string) error {
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

	// Get list of reserved IPs
	response, err := client.GetReservedIps(projectID, location)
	if err != nil {
		return fmt.Errorf("error listing reserved IPs: %w", err)
	}

	if response == nil || len(response.Data) == 0 {
		log.Printf("[WARNING] No reserved IPs found")
		return nil
	}

	log.Printf("[DEBUG] Found %d reserved IPs in total", len(response.Data))

	sweptCount := 0
	for _, reserveIP := range response.Data {
		// Reserved IPs don't have names, but we can filter by VM name if attached
		// For test cleanup, we'll clean up all unattached IPs created during tests
		// This is a simplified approach - in production you might want more sophisticated filtering

		// Skip if attached to a VM (to avoid breaking existing resources)
		if reserveIP.VMName != "" && !strings.HasPrefix(reserveIP.VMName, "test-") {
			log.Printf("[DEBUG] Skipping reserved IP %s (attached to non-test VM: %s)", reserveIP.IPAddress, reserveIP.VMName)
			continue
		}

		log.Printf("[INFO] Deleting reserved IP: %s", reserveIP.IPAddress)

		err := client.DeleteReserveIP(reserveIP.IPAddress, projectID, location)
		if err != nil {
			log.Printf("[ERROR] Failed to delete reserved IP %s: %v", reserveIP.IPAddress, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted reserved IP: %s", reserveIP.IPAddress)
	}

	log.Printf("[INFO] Swept %d reserved IPs", sweptCount)
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
