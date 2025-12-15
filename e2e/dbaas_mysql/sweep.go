package dbaas_mysql

import (
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testMySQLDBaaSPrefix = sweep.TestNamePrefix + "mysql-"

func init() {
	resource.AddTestSweepers("e2e_dbaas_mysql", &resource.Sweeper{
		Name:         "e2e_dbaas_mysql",
		F:            sweepMySQLDBaaS,
		Dependencies: []string{
			// MySQL DBaaS instances should be swept after any dependent resources
			// Add dependencies here if needed (e.g., "e2e_node" if there are node dependencies)
		},
	})
}

func sweepMySQLDBaaS(region string) error {
	// NOTE: The E2E API client does not currently have a method to list all MySQL DBaaS instances.
	// The GetCluster method only retrieves a single MySQL DBaaS instance by ID.
	// Until a ListMySQLClusters or GetMySQLClusters method is implemented in the client,
	// automatic sweeping of test MySQL DBaaS instances is not possible.
	// Test MySQL DBaaS instances will need to be cleaned up manually or via the destroy functionality
	// in the test framework.

	log.Printf("[INFO] MySQL DBaaS sweeping not yet implemented - waiting for ListMySQLClusters API method")
	log.Printf("[INFO] Please manually clean up test MySQL DBaaS instances with prefix '%s' if needed", testMySQLDBaaSPrefix)

	// Uncomment and update the code below once the API client has a list method:
	//
	// cfg, err := sweep.SharedConfigForRegion(region)
	// if err != nil {
	// 	return fmt.Errorf("error getting config for region %s: %w", region, err)
	// }
	//
	// goe2eClient := cfg.Goe2eClient()
	// ctx := context.Background()
	//
	// projectIDStr := sweep.TestProjectID
	// if projectIDStr == "" {
	// 	log.Printf("[WARNING] E2E_TEST_PROJECT_ID not set, skipping sweep")
	// 	return nil
	// }
	//
	// // Get list of MySQL DBaaS instances
	// clusters, err := goe2eClient.DBaaSMySQL.ListClusters(ctx, projectIDStr)
	// if err != nil {
	// 	return fmt.Errorf("error listing MySQL DBaaS clusters: %w", err)
	// }
	//
	// for _, cluster := range clusters {
	// 	if strings.HasPrefix(cluster.Name, testMySQLDBaaSPrefix) {
	// 		log.Printf("[INFO] Deleting MySQL DBaaS instance: %s (ID: %d)", cluster.Name, cluster.ID)
	// 		_, err := goe2eClient.DBaaSMySQL.DeleteCluster(ctx, strconv.Itoa(cluster.ID))
	// 		if err != nil {
	// 			log.Printf("[ERROR] Error deleting MySQL DBaaS instance %s (ID: %d): %s", cluster.Name, cluster.ID, err)
	// 			continue
	// 		}
	// 	}
	// }

	return nil
}
