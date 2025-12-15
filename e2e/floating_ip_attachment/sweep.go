package floating_ip_attachment

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNodePrefix = sweep.TestNamePrefix + "node-"

func init() {
	resource.AddTestSweepers("e2e_floating_ip_attachment", &resource.Sweeper{
		Name: "e2e_floating_ip_attachment",
		F:    sweepFloatingIPAttachments,
		Dependencies: []string{
			// Floating IP attachments depend on nodes, so sweep after nodes
			"e2e_node",
		},
	})
}

// sweepFloatingIPAttachments sweeps floating IP attachments that are attached to test nodes.
// This is a safe sweeper that only detaches floating IPs from nodes with test prefixes.
func sweepFloatingIPAttachments(region string) error {
	client, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	ctx := context.Background()

	// Get list of reserved IPs using goe2e client
	rips, _, err := client.ReserveIP.ListReserveIPs(ctx)
	if err != nil {
		return fmt.Errorf("error listing reserved IPs: %w", err)
	}

	if len(rips) == 0 {
		log.Printf("[DEBUG] No reserved IPs found")
		return nil
	}

	log.Printf("[DEBUG] Found %d reserved IPs in total", len(rips))

	sweptCount := 0
	for _, rip := range rips {
		// Only process FloatingIP types
		if rip.ReservedType != goe2econstants.ReserveIPTypeFloatingIP {
			log.Printf("[DEBUG] Skipping reserved IP %s (type: %s, not FloatingIP)", rip.IPAddress, rip.ReservedType)
			continue
		}

		// Check if this floating IP has any attached nodes
		if len(rip.FloatingIPAttachedNodes) == 0 {
			log.Printf("[DEBUG] Skipping floating IP %s (no attached nodes)", rip.IPAddress)
			continue
		}

		// Collect node IDs for detachment
		// Safety: We check node names to ensure we only detach from test nodes
		testNodeIDs := make([]string, 0)

		for _, node := range rip.FloatingIPAttachedNodes {
			// Check if node name has test prefix (safety check)
			if strings.HasPrefix(node.Name, testNodePrefix) {
				nodeIDStr := fmt.Sprintf("%d", node.ID)
				testNodeIDs = append(testNodeIDs, nodeIDStr)
			}
		}

		// Safety check: Only detach if we found at least one test node
		// This ensures we don't accidentally detach from production nodes
		if len(testNodeIDs) == 0 {
			log.Printf("[DEBUG] Skipping floating IP %s (no test nodes attached)", rip.IPAddress)
			continue
		}

		// Additional safety: Only proceed if we're in the test region/project
		// This is handled by SharedGoe2eClientForTests which uses E2E_TEST_PROJECT_ID and E2E_TEST_REGION

		ipAddress := rip.IPAddress
		log.Printf("[INFO] Detaching floating IP %s from %d node(s)", ipAddress, len(testNodeIDs))

		// Use the DetachFloatingIP method with proper request type
		detachReq := &goe2e.FloatingIPDetachmentRequest{
			IPAddress: ipAddress,
			NodeIDs:   testNodeIDs,
		}

		_, err = client.ReserveIP.DetachFloatingIP(ctx, detachReq)

		if err != nil {
			log.Printf("[ERROR] Failed to detach floating IP %s: %v", ipAddress, err)
			// Continue with other attachments even if one fails
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully detached floating IP %s from %d node(s)", ipAddress, len(testNodeIDs))
	}

	log.Printf("[INFO] Swept %d floating IP attachment(s)", sweptCount)
	return nil
}
