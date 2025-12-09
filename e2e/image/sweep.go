package image

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testImagePrefix = sweep.TestNamePrefix + "image-"

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
	goe2eClient, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	ctx := context.Background()

	// Get list of saved images
	images, _, err := goe2eClient.Images.GetSavedImages(ctx)
	if err != nil {
		return fmt.Errorf("error listing images: %w", err)
	}

	if len(images) == 0 {
		log.Printf("[WARNING] No images found")
		return nil
	}

	log.Printf("[DEBUG] Found %d images in total", len(images))

	sweptCount := 0
	for _, image := range images {
		imageName := image.Name
		if !strings.HasPrefix(imageName, testImagePrefix) {
			log.Printf("[DEBUG] Skipping image %s (does not have test prefix)", imageName)
			continue
		}

		imageID := image.ImageID
		log.Printf("[INFO] Deleting image: %s (ID: %s)", imageName, imageID)

		_, _, err := goe2eClient.Images.DeleteImage(ctx, imageID)
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
