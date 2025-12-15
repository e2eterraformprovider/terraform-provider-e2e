package image

import (
	"context"
	"fmt"
	"strings"
	"time"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

// File-local timeout and polling constants for image operations
// These are Terraform-specific timing configurations, not part of the API contract
const (
	imageCreateTimeout   = 30 * time.Minute
	imageDeleteTimeout   = 5 * time.Minute
	imagePollingInterval = 5 * time.Second
)

// normalizeImageState normalizes API image_state to Terraform state
// Maps API states (Creating, Ready, Error, etc.) to normalized states (creating, ready, error, deleted)
func normalizeImageState(imageState string) string {
	// Normalize common states
	switch strings.ToLower(imageState) {
	case strings.ToLower(goe2econstants.ImageStatusCreating):
		return goe2econstants.ImageStateCreating
	case strings.ToLower(goe2econstants.ImageStatusReady):
		return goe2econstants.ImageStateReady
	case strings.ToLower(goe2econstants.ImageStatusError):
		return goe2econstants.ImageStateError
	case strings.ToLower(goe2econstants.ImageStatusDeleted):
		return goe2econstants.ImageStateDeleted
	default:
		return strings.ToLower(imageState)
	}
}

// waitForImageState polls the image until it reaches the desired state
// Used for async operations like image creation (Creating -> Ready)
func waitForImageState(ctx context.Context, client *goe2e.Client, imageID, desiredState string, timeout time.Duration) error {
	// Helper function to check the image state
	checkState := func() (bool, error) {
		image, _, err := client.Images.GetImage(ctx, imageID)
		if err != nil {
			// If desired state is "deleted", 404 is expected
			if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) && desiredState == goe2econstants.ImageStateDeleted {
				return true, nil
			}
			return false, fmt.Errorf(ErrorCheckingImageState, err)
		}

		currentState := normalizeImageState(image.ImageState)
		if currentState == desiredState {
			return true, nil
		}

		// If image enters error state, return error
		if currentState == goe2econstants.ImageStateError {
			return false, fmt.Errorf(goe2econstants.ImageEnteredErrorState, imageID)
		}

		return false, nil
	}

	// Check immediately first (before waiting for ticker)
	if done, err := checkState(); done || err != nil {
		return err
	}

	ticker := time.NewTicker(imagePollingInterval)
	defer ticker.Stop()
	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeoutTimer.C:
			return fmt.Errorf(goe2econstants.ImageTimeoutWaitingForState, imageID, desiredState)
		case <-ticker.C:
			if done, err := checkState(); done || err != nil {
				return err
			}
		}
	}
}

// flattenImageResponse converts a goe2e.SavedImage to Terraform state map
func flattenImageResponse(image *goe2e.SavedImage) map[string]interface{} {
	result := make(map[string]interface{})

	result[tfconstants.AttrTemplateID] = image.TemplateID
	result["image_state"] = image.ImageState
	result["state"] = normalizeImageState(image.ImageState)
	result["image_type"] = image.ImageType
	result["os_distribution"] = image.OSDistribution
	result[tfconstants.AttrName] = image.Name
	result["image_id"] = image.ImageID
	result["distro"] = image.Distro
	result["sku_type"] = image.SKUType
	result["image_size"] = image.ImageSize
	result["cloning_ops"] = image.CloningOps
	result["running_vms"] = image.RunningVMs
	result["is_windows"] = image.IsWindows
	result[tfconstants.AttrCreatedAt] = image.CreationTime
	result["vm_info"] = image.VMInfo

	return result
}
