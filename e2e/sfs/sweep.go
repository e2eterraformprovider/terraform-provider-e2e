package sfs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = "test-sfs-"

func init() {
	resource.AddTestSweepers("e2e_sfs", &resource.Sweeper{
		Name: "e2e_sfs",
		F:    sweepSFS,
	})
}

func sweepSFS(region string) error {
	cfg, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting config for region %s: %w", region, err)
	}

	client := cfg.Client()

	projectID := os.Getenv("E2E_TEST_PROJECT_ID")
	testRegion := os.Getenv("E2E_TEST_REGION")

	if projectID == "" || testRegion == "" {
		log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_REGION not set, skipping sweep")
		return nil
	}

	response, err := client.GetSfss(testRegion, projectID)
	if err != nil {
		return fmt.Errorf("error listing SFS: %w", err)
	}

	log.Printf("[DEBUG] Found %d SFS instances in total", len(response.Data))

	sweptCount := 0
	for _, sfs := range response.Data {
		if !strings.HasPrefix(sfs.Name, testNamePrefix) {
			log.Printf("[DEBUG] Skipping SFS %s (does not have test prefix)", sfs.Name)
			continue
		}

		log.Printf("[INFO] Deleting SFS: %s (ID: %d)", sfs.Name, sfs.ID)

		sfsID := fmt.Sprintf("%d", sfs.ID)
		err := client.DeleteSFs(sfsID, projectID, testRegion)
		if err != nil {
			log.Printf("[ERROR] Failed to delete SFS %s: %v", sfs.Name, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted SFS: %s", sfs.Name)
	}

	log.Printf("[INFO] Swept %d SFS instances", sweptCount)

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
