package ssh_key

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testSSHKeyPrefix = "test-ssh-key-"

func init() {
	resource.AddTestSweepers("e2e_ssh_key", &resource.Sweeper{
		Name: "e2e_ssh_key",
		F:    sweepSSHKeys,
	})
}

func sweepSSHKeys(region string) error {
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

	// Get list of SSH keys
	response, err := client.GetSshKeys(location, projectID)
	if err != nil {
		return fmt.Errorf("error listing SSH keys: %w", err)
	}

	if response == nil || len(response.Data) == 0 {
		log.Printf("[WARNING] No SSH keys found")
		return nil
	}

	log.Printf("[DEBUG] Found %d SSH keys in total", len(response.Data))

	sweptCount := 0
	for _, sshKey := range response.Data {
		label := sshKey.Label
		if !strings.HasPrefix(label, testSSHKeyPrefix) {
			log.Printf("[DEBUG] Skipping SSH key %s (does not have test prefix)", label)
			continue
		}

		sshKeyID := fmt.Sprintf("%d", sshKey.Pk)
		log.Printf("[INFO] Deleting SSH key: %s (ID: %s)", label, sshKeyID)

		err := client.DeleteSshKey(sshKeyID, projectID, location)
		if err != nil {
			log.Printf("[ERROR] Failed to delete SSH key %s: %v", label, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted SSH key: %s", label)
	}

	log.Printf("[INFO] Swept %d SSH keys", sweptCount)
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
