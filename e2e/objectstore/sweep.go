package objectstore

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = "test-bucket-"

func init() {
	resource.AddTestSweepers("e2e_objectstore", &resource.Sweeper{
		Name: "e2e_objectstore",
		F:    sweepObjectStores,
	})
}

func sweepObjectStores(region string) error {
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

	response, err := client.GetBuckets(testRegion, projectID)
	if err != nil {
		return fmt.Errorf("error listing object store buckets: %w", err)
	}

	log.Printf("[DEBUG] Found %d object store buckets in total", len(response.Data))

	sweptCount := 0
	for _, bucket := range response.Data {
		if !strings.HasPrefix(bucket.Name, testNamePrefix) {
			log.Printf("[DEBUG] Skipping object store bucket %s (does not have test prefix)", bucket.Name)
			continue
		}

		log.Printf("[INFO] Deleting object store bucket: %s", bucket.Name)

		err := client.DeleteBucket(bucket.Name, testRegion, projectID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete object store bucket %s: %v", bucket.Name, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted object store bucket: %s", bucket.Name)
	}

	log.Printf("[INFO] Swept %d object store buckets", sweptCount)

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
