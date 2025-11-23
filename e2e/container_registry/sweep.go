package container_registry

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = "test-cr-"

func init() {
	resource.AddTestSweepers("e2e_container_registry", &resource.Sweeper{
		Name: "e2e_container_registry",
		F:    sweepContainerRegistries,
	})
}

func sweepContainerRegistries(region string) error {
	cfg, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting config for region %s: %w", region, err)
	}

	client := cfg.Client()

	projectID := os.Getenv("E2E_TEST_PROJECT_ID")
	location := os.Getenv("E2E_TEST_LOCATION")

	if projectID == "" || location == "" {
		log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_LOCATION not set, skipping sweep")
		return nil
	}

	registries, err := client.GetContainerRegistryProjects(projectID, location)
	if err != nil {
		return fmt.Errorf("error listing container registries: %w", err)
	}

	log.Printf("[DEBUG] Found %d container registries in total", len(registries))

	sweptCount := 0
	for _, registry := range registries {
		if !strings.HasPrefix(registry.ProjectName, testNamePrefix) {
			log.Printf("[DEBUG] Skipping container registry %s (does not have test prefix)", registry.ProjectName)
			continue
		}

		log.Printf("[INFO] Deleting container registry: %s (ID: %d)", registry.ProjectName, registry.ID)

		registryID := fmt.Sprintf("%d", registry.ID)
		userID := "0"

		err := client.DeleteContainerRegistry(registryID, registry.ProjectName, userID, projectID, location)
		if err != nil {
			log.Printf("[ERROR] Failed to delete container registry %s: %v", registry.ProjectName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted container registry: %s", registry.ProjectName)
	}

	log.Printf("[INFO] Swept %d container registries", sweptCount)

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
