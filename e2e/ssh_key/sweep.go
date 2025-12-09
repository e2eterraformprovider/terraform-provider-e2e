package ssh_key

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testSSHKeyPrefix = sweep.TestNamePrefix + "ssh-key-"

func init() {
	resource.AddTestSweepers("e2e_ssh_key", &resource.Sweeper{
		Name: "e2e_ssh_key",
		F:    sweepSSHKeys,
	})
}

func sweepSSHKeys(region string) error {
	goe2eClient, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	ctx := context.Background()

	// Get list of SSH keys
	sshKeys, _, err := goe2eClient.SSHKeys.ListSSHKeys(ctx)
	if err != nil {
		return fmt.Errorf("error listing SSH keys: %w", err)
	}

	if len(sshKeys) == 0 {
		log.Printf("[WARNING] No SSH keys found")
		return nil
	}

	log.Printf("[DEBUG] Found %d SSH keys in total", len(sshKeys))

	sweptCount := 0
	for _, sshKey := range sshKeys {
		label := sshKey.Label
		if !strings.HasPrefix(label, testSSHKeyPrefix) {
			log.Printf("[DEBUG] Skipping SSH key %s (does not have test prefix)", label)
			continue
		}

		sshKeyID := fmt.Sprintf("%d", sshKey.PK)
		log.Printf("[INFO] Deleting SSH key: %s (ID: %s)", label, sshKeyID)

		_, err := goe2eClient.SSHKeys.DeleteSSHKey(ctx, sshKeyID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete SSH key %s: %v", label, err)
			continue
		}

		sweptCount++
		log.Printf("[INFO] Successfully deleted SSH key: %s", label)
	}

	log.Printf("[INFO] Swept %d SSH keys", sweptCount)
	return nil
}
