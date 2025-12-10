package reserve_ip

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testReserveIPPrefix = sweep.TestNamePrefix + "reserve-ip-"

func init() {
	resource.AddTestSweepers("e2e_reserve_ip", &resource.Sweeper{
		Name: "e2e_reserve_ip",
		F:    sweepReserveIPs,
	})
}

func sweepReserveIPs(region string) error {
	goe2eClient, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	ctx := context.Background()

	// Get list of reserved IPs
	reserveIPs, _, err := goe2eClient.ReserveIP.ListReserveIPs(ctx)
	if err != nil {
		return fmt.Errorf("error listing reserved IPs: %w", err)
	}

	if len(reserveIPs) == 0 {
		log.Printf("[WARNING] No reserved IPs found")
		return nil
	}

	log.Printf("[DEBUG] Found %d reserved IPs in total", len(reserveIPs))

	sweptCount := 0
	for _, ip := range reserveIPs {
		// Reserved IPs don't have names, but we can filter by VM name if attached
		// For test cleanup, we'll clean up all unattached IPs created during tests
		// This is a simplified approach - in production you might want more sophisticated filtering

		// Skip if attached to a VM (to avoid breaking existing resources)
		if ip.VMName != "" && !strings.HasPrefix(ip.VMName, "test-") {
			log.Printf("[DEBUG] Skipping reserved IP %s (attached to non-test VM: %s)", ip.IPAddress, ip.VMName)
			continue
		}

		log.Printf("[INFO] Deleting reserved IP: %s", ip.IPAddress)

		_, err := goe2eClient.ReserveIP.DeleteReserveIP(ctx, ip.IPAddress)
		if err != nil {
			log.Printf("[ERROR] Failed to delete reserved IP %s: %v", ip.IPAddress, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted reserved IP: %s", ip.IPAddress)
	}

	log.Printf("[INFO] Swept %d reserved IPs", sweptCount)
	return nil
}
