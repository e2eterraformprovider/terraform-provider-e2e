package volume_attachment

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testVolumeAttachmentPrefix = sweep.TestNamePrefix + "va-"

func init() {
	resource.AddTestSweepers("e2e_volume_attachment", &resource.Sweeper{
		Name:         "e2e_volume_attachment",
		F:            sweepVolumeAttachments,
		Dependencies: []string{
			// Volume attachments should be swept before nodes and block storages
			// to avoid dependency issues during cleanup
		},
	})
}

func sweepVolumeAttachments(region string) error {
	client, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	ctx := context.Background()

	// Get list of nodes using goe2e client
	nodes, _, err := client.Nodes.ListNodes(ctx)
	if err != nil {
		return fmt.Errorf("error listing nodes: %w", err)
	}

	if len(nodes) == 0 {
		log.Printf("[INFO] No nodes found, skipping volume attachment sweep")
		return nil
	}

	log.Printf("[DEBUG] Found %d nodes to check for volume attachments", len(nodes))

	sweptCount := 0
	for _, node := range nodes {
		nodeName := node.Name
		// Only process nodes with test prefix
		if !strings.HasPrefix(nodeName, sweep.TestNamePrefix) {
			log.Printf("[DEBUG] Skipping node %s (does not have test prefix)", nodeName)
			continue
		}

		nodeID := node.ID
		log.Printf("[DEBUG] Checking node %s (ID: %s) for volume attachments", nodeName, nodeID)

		// Get node details to check for attached volumes
		nodeDetails, _, err := client.Nodes.GetNode(ctx, nodeID)
		if err != nil {
			log.Printf("[WARNING] Failed to get details for node %s: %v", nodeName, err)
			continue
		}

		if nodeDetails == nil {
			log.Printf("[DEBUG] Node %s not found, skipping", nodeName)
			continue
		}

		// Check if node has any attached volumes by checking block storages
		// Since we can't list all block storages, we'll need to iterate through
		// block storages that might be attached. However, the API doesn't provide
		// a direct way to list attachments.
		//
		// Alternative approach: List all block storages with test prefix and check
		// if they're attached to test nodes. But block storage list is also not available.
		//
		// For now, we'll document that volume attachments are cleaned up when
		// nodes are deleted (as part of node sweep), but we can't directly sweep
		// attachments without a list API.

		// Note: Volume attachments are typically cleaned up automatically when
		// nodes are deleted. However, if we need to explicitly detach volumes,
		// we would need a way to list them. Since GetAttachments is not implemented,
		// we'll log a warning and rely on node/block storage sweeps.

		log.Printf("[INFO] Volume attachment sweeping requires GetAttachments API or block storage list API")
		log.Printf("[INFO] Volume attachments will be cleaned up when test nodes are deleted")
		log.Printf("[INFO] If you need to explicitly detach volumes, use the block storage sweep after node sweep")
	}

	log.Printf("[INFO] Volume attachment sweep completed (swept %d attachments)", sweptCount)
	return nil
}

// Note: Full volume attachment sweeping requires either:
// 1. A ListBlockStorages API method to list all block storages, or
// 2. A GetAttachments API method to list attachments for a node
//
// Currently, neither API method is available. Volume attachments are automatically
// cleaned up when test nodes are deleted (as part of node sweep).
//
// When the API becomes available, the sweep logic should:
// 1. List all block storages with test prefix
// 2. For each block storage, check if it's attached to a test node
// 3. If attached to a test node, detach it using VolumeAttachment.DetachVolume
