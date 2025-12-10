package security_group

import (
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testSecurityGroupPrefix = sweep.TestNamePrefix + "sg-"

func init() {
	resource.AddTestSweepers("e2e_security_group", &resource.Sweeper{
		Name: "e2e_security_group",
		F:    sweepSecurityGroups,
	})
}

func sweepSecurityGroups(region string) error {
	_, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	// Note: We would need a ListSecurityGroups method in the client
	// For now, we'll skip the sweep if the method doesn't exist
	// This is a placeholder that matches the pattern
	log.Printf("[INFO] Sweeping security groups is not fully implemented yet")
	log.Printf("[INFO] Security groups with prefix %s should be manually cleaned up if needed", testSecurityGroupPrefix)

	// Placeholder implementation:
	// response, err := client.ListSecurityGroups(projectID, location)
	// if err != nil {
	//     return fmt.Errorf("error listing security groups: %w", err)
	// }

	// sweptCount := 0
	// for _, sg := range response {
	//     if !strings.HasPrefix(sg.Name, testSecurityGroupPrefix) {
	//         continue
	//     }
	//     err := client.DeleteSecurityGroup(sg.ID, projectID, location)
	//     if err != nil {
	//         log.Printf("[ERROR] Failed to delete security group %s: %v", sg.Name, err)
	//         continue
	//     }
	//     sweptCount++
	// }
	// log.Printf("[INFO] Swept %d security groups", sweptCount)

	return nil
}
