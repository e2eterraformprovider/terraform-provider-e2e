package vpc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testVPCPrefix = sweep.TestNamePrefix + "vpc-"

func init() {
	resource.AddTestSweepers("e2e_vpc", &resource.Sweeper{
		Name:         "e2e_vpc",
		F:            sweepVPCs,
		Dependencies: []string{
			// VPCs may be dependencies for other resources (nodes, etc.)
			// Add dependencies if needed
		},
	})
}

func sweepVPCs(region string) error {
	client, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	ctx := context.Background()
	// Get list of VPCs using goe2e client
	vpcs, _, err := client.Vpcs.ListVPCs(ctx)
	if err != nil {
		return fmt.Errorf("error listing VPCs: %w", err)
	}

	if len(vpcs) == 0 {
		log.Printf("[WARNING] No VPCs found")
		return nil
	}

	log.Printf("[DEBUG] Found %d VPCs in total", len(vpcs))

	sweptCount := 0
	for _, vpc := range vpcs {
		vpcName := vpc.Name
		if !strings.HasPrefix(vpcName, testVPCPrefix) {
			log.Printf("[DEBUG] Skipping VPC %s (does not have test prefix)", vpcName)
			continue
		}

		// Convert network_id (float64) to string for deletion
		vpcIDStr := strconv.FormatFloat(vpc.ID, 'f', -1, 64)
		log.Printf("[INFO] Deleting VPC: %s (ID: %s)", vpcName, vpcIDStr)

		ctx := context.Background()
		_, err = client.Vpcs.DeleteVPC(ctx, vpcIDStr)
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPC %s: %v", vpcName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted VPC: %s", vpcName)
	}

	log.Printf("[INFO] Swept %d VPCs", sweptCount)
	return nil
}
