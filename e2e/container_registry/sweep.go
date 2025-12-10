package container_registry

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = sweep.TestNamePrefix + "cr-"

func init() {
	resource.AddTestSweepers("e2e_container_registry", &resource.Sweeper{
		Name: "e2e_container_registry",
		F:    sweepContainerRegistries,
	})
}

func sweepContainerRegistries(region string) error {
	client, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	ctx := context.Background()

	// List all registries (no pagination limit, just get all)
	listOpts := &goe2e.ContainerRegistryListOptions{Page: 1, PageSize: 100}
	registries, _, err := client.ContainerRegistry.ListContainerRegistryProjects(ctx, listOpts)
	if err != nil {
		return fmt.Errorf("error listing container registries: %w", err)
	}

	log.Printf("[DEBUG] Found %d container registries in total", len(registries))

	sweptCount := 0
	for _, registry := range registries {
		if !strings.HasPrefix(registry.ProjectName, testNamePrefix) {
			log.Printf("[DEBUG] Skipping container registry %s (does not have test prefix)", registry.ProjectName)
			continue
		}

		log.Printf("[INFO] Deleting container registry: %s (ID: %d)", registry.ProjectName, registry.ID)

		registryID := fmt.Sprintf("%d", registry.ID)
		userID := strconv.Itoa(registry.Customer)

		deleteReq := &goe2e.ContainerRegistryDeleteRequest{
			CRProjectID: registryID,
			ProjectName: registry.ProjectName,
			UserID:      userID,
		}

		_, err := client.ContainerRegistry.DeleteContainerRegistry(ctx, deleteReq)
		if err != nil {
			log.Printf("[ERROR] Failed to delete container registry %s: %v", registry.ProjectName, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted container registry: %s", registry.ProjectName)
	}

	log.Printf("[INFO] Swept %d container registries", sweptCount)

	return nil
}
