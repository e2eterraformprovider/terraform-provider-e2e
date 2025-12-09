package image

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

// normalizeImageState normalizes API image_state to Terraform state
// Maps API states (Creating, Ready, Error, etc.) to normalized states (creating, ready, error, deleted)
func normalizeImageState(imageState string) string {
	// Normalize common states
	switch strings.ToLower(imageState) {
	case "creating":
		return "creating"
	case "ready":
		return "ready"
	case "error":
		return "error"
	case "deleted":
		return "deleted"
	default:
		return strings.ToLower(imageState)
	}
}

// waitForImageState polls the image until it reaches the desired state
// Used for async operations like image creation (Creating -> Ready)
func waitForImageState(ctx context.Context, client *goe2e.Client, imageID, desiredState string, timeout time.Duration) error {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeoutTimer.C:
			return fmt.Errorf("timeout waiting for image %s to reach state %s", imageID, desiredState)
		case <-ticker.C:
			image, _, err := client.Images.GetImage(ctx, imageID)
			if err != nil {
				// If desired state is "deleted", 404 is expected
				if strings.Contains(err.Error(), "not found") && desiredState == "deleted" {
					return nil
				}
				return fmt.Errorf("error checking image state: %w", err)
			}

			currentState := normalizeImageState(image.ImageState)
			if currentState == desiredState {
				return nil
			}

			// If image enters error state, return error
			if currentState == "error" {
				return fmt.Errorf("image %s entered error state", imageID)
			}
		}
	}
}

// flattenImageResponse converts a goe2e.SavedImage to Terraform state map
func flattenImageResponse(image *goe2e.SavedImage) map[string]interface{} {
	result := make(map[string]interface{})

	result["template_id"] = image.TemplateID
	result["image_state"] = image.ImageState
	result["state"] = normalizeImageState(image.ImageState)
	result["image_type"] = image.ImageType
	result["os_distribution"] = image.OSDistribution
	result["name"] = image.Name
	result["image_id"] = image.ImageID
	result["distro"] = image.Distro
	result["sku_type"] = image.SKUType
	result["image_size"] = image.ImageSize
	result["cloning_ops"] = image.CloningOps
	result["running_vms"] = image.RunningVMs
	result["is_windows"] = image.IsWindows
	result["creation_time"] = image.CreationTime
	result["vm_info"] = image.VMInfo

	return result
}
