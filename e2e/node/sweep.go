package node

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNodePrefix = sweep.TestNamePrefix + "node-"

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
		log.Printf("[WARNING] No nodes found")
		return nil
	}

	log.Printf("[DEBUG] Found %d nodes in total", len(nodes))

	sweptCount := 0
	for _, node := range nodes {
		nodeName := node.Name
		if !strings.HasPrefix(nodeName, testNodePrefix) {
			log.Printf("[DEBUG] Skipping node %s (does not have test prefix)", nodeName)
			continue
		}

		nodeIDStr := node.ID
		log.Printf("[INFO] Deleting node: %s (ID: %s)", nodeName, nodeIDStr)

		ctx := context.Background()
		_, err = client.Nodes.DeleteNode(ctx, nodeIDStr)
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
