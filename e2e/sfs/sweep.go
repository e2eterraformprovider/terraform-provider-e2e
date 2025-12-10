package sfs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = sweep.TestNamePrefix + "sfs-"

func init() {
	resource.AddTestSweepers("e2e_sfs", &resource.Sweeper{
		Name: "e2e_sfs",
		F:    sweepSFS,
	})
}

func sweepSFS(region string) error {
	cfg, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting config for region %s: %w", region, err)
	}

	client, err := cfg.Goe2eClientForProject("", region)
	if err != nil {
		return fmt.Errorf("error creating goe2e client for region %s: %w", region, err)
	}

	ctx := context.Background()

	sfsList, _, err := client.Sfs.ListSfss(ctx)
	if err != nil {
		return fmt.Errorf("error listing SFS: %w", err)
	}

	log.Printf("[DEBUG] Found %d SFS instances in total", len(sfsList))

	sweptCount := 0
	for _, sfs := range sfsList {
		if !strings.HasPrefix(sfs.Name, testNamePrefix) {
			log.Printf("[DEBUG] Skipping SFS %s (does not have test prefix)", sfs.Name)
			continue
		}

		log.Printf("[INFO] Deleting SFS: %s (ID: %s)", sfs.Name, sfs.ID)

		_, err := client.Sfs.DeleteSfs(ctx, sfs.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete SFS %s: %v", sfs.Name, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted SFS: %s", sfs.Name)
	}

	log.Printf("[INFO] Swept %d SFS instances", sweptCount)

	return nil
}
