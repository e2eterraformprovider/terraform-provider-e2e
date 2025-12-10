package objectstore

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = sweep.TestNamePrefix + "bucket-"

func init() {
	resource.AddTestSweepers("e2e_objectstore", &resource.Sweeper{
		Name: "e2e_objectstore",
		F:    sweepObjectStores,
	})
}

func sweepObjectStores(region string) error {
	goe2eClient, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	ctx := context.Background()
	buckets, _, err := goe2eClient.ObjectStorage.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("error listing object store buckets: %w", err)
	}

	log.Printf("[DEBUG] Found %d object store buckets in total", len(buckets))

	sweptCount := 0
	for _, bucket := range buckets {
		if !strings.HasPrefix(bucket.Name, testNamePrefix) {
			log.Printf("[DEBUG] Skipping object store bucket %s (does not have test prefix)", bucket.Name)
			continue
		}

		log.Printf("[INFO] Deleting object store bucket: %s", bucket.Name)

		_, err := goe2eClient.ObjectStorage.DeleteBucket(ctx, bucket.Name)
		if err != nil {
			log.Printf("[ERROR] Failed to delete object store bucket %s: %v", bucket.Name, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted object store bucket: %s", bucket.Name)
	}

	log.Printf("[INFO] Swept %d object store buckets", sweptCount)

	return nil
}
