package image

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testImagePrefix = "test-image-"

func init() {
	resource.AddTestSweepers("e2e_image", &resource.Sweeper{
		Name:         "e2e_image",
		F:            sweepImages,
		Dependencies: []string{
			// Images should be swept before nodes if there are any dependencies
		},
	})
}

func sweepImages(region string) error {
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

	// Get list of saved images
	response, err := client.GetSavedImages(location, projectID)
	if err != nil {
		return fmt.Errorf("error listing images: %w", err)
	}

	if response == nil || len(response.Data) == 0 {
		log.Printf("[WARNING] No images found")
		return nil
	}

	log.Printf("[DEBUG] Found %d images in total", len(response.Data))

	sweptCount := 0
	for _, image := range response.Data {
		imageName := image.Name
		if !strings.HasPrefix(imageName, testImagePrefix) {
			log.Printf("[DEBUG] Skipping image %s (does not have test prefix)", imageName)
			continue
		}

		imageID := image.Image_id
		log.Printf("[INFO] Deleting image: %s (ID: %s)", imageName, imageID)

		err := client.DeleteImage(imageID, projectID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete image %s: %v", imageName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted image: %s", imageName)
	}

	log.Printf("[INFO] Swept %d images", sweptCount)
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
