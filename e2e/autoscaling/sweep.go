package autoscaling

import (
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = sweep.TestNamePrefix + "sg-"

func init() {
	resource.AddTestSweepers("e2e_scaler_group", &resource.Sweeper{
		Name: "e2e_scaler_group",
		F:    sweepScalerGroups,
	})
}

func sweepScalerGroups(region string) error {
	_, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	// Note: This assumes there's a method to list all scaler groups
	// If such a method doesn't exist, we can't sweep them automatically
	log.Printf("[INFO] Scaler group sweeping would require list API method")
	log.Printf("[INFO] Please manually delete scaler groups with prefix '%s' if needed", testNamePrefix)

	// TODO: Implement when ListScalerGroups API is available
	// groups, err := client.ListScalerGroups(projectID, location)
	// if err != nil {
	// 	return fmt.Errorf("error listing scaler groups: %w", err)
	// }
	//
	// log.Printf("[DEBUG] Found %d scaler groups in total", len(groups))
	//
	// sweptCount := 0
	// for _, group := range groups {
	// 	if !strings.HasPrefix(group.Name, testNamePrefix) {
	// 		log.Printf("[DEBUG] Skipping scaler group %s (does not have test prefix)", group.Name)
	// 		continue
	// 	}
	//
	// 	log.Printf("[INFO] Deleting scaler group: %s (ID: %s)", group.Name, group.ID)
	//
	// 	err := client.DeleteScalerGroup(group.ID, projectID, location)
	// 	if err != nil {
	// 		log.Printf("[ERROR] Failed to delete scaler group %s: %v", group.Name, err)
	// 		continue
	// 	}
	//
	// 	sweptCount++
	// 	log.Printf("[INFO] Successfully deleted scaler group: %s", group.Name)
	// }
	//
	// log.Printf("[INFO] Swept %d scaler groups", sweptCount)

	return nil
}
