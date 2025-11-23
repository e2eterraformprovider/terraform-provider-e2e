package node

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNodePrefix = "test-node-"

func init() {
	resource.AddTestSweepers("e2e_node", &resource.Sweeper{
		Name:         "e2e_node",
		F:            sweepNodes,
		Dependencies: []string{
			// Add dependencies if needed (e.g., if nodes depend on other resources)
		},
	})
}

func sweepNodes(region string) error {
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

	// Get list of nodes
	response, err := client.GetNodes(projectID, location)
	if err != nil {
		return fmt.Errorf("error listing nodes: %w", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		log.Printf("[WARNING] No nodes found or invalid response format")
		return nil
	}

	log.Printf("[DEBUG] Found %d nodes in total", len(data))

	sweptCount := 0
	for _, item := range data {
		node, ok := item.(map[string]interface{})
		if !ok {
			log.Printf("[WARNING] Invalid node format, skipping")
			continue
		}

		nodeName, ok := node["name"].(string)
		if !ok {
			log.Printf("[WARNING] Node name not found, skipping")
			continue
		}

		if !strings.HasPrefix(nodeName, testNodePrefix) {
			log.Printf("[DEBUG] Skipping node %s (does not have test prefix)", nodeName)
			continue
		}

		nodeID, ok := node["id"].(float64)
		if !ok {
			log.Printf("[WARNING] Invalid node ID format for %s, skipping", nodeName)
			continue
		}

		nodeIDStr := fmt.Sprintf("%.0f", nodeID)
		log.Printf("[INFO] Deleting node: %s (ID: %s)", nodeName, nodeIDStr)

		err := client.DeleteNode(nodeIDStr, projectID, location)
		if err != nil {
			log.Printf("[ERROR] Failed to delete node %s: %v", nodeName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted node: %s", nodeName)
	}

	log.Printf("[INFO] Swept %d nodes", sweptCount)
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
