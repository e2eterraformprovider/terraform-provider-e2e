package kubernetes

import (
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = sweep.TestNamePrefix + "k8s-"

func init() {
	resource.AddTestSweepers("e2e_kubernetes", &resource.Sweeper{
		Name: "e2e_kubernetes",
		F:    sweepKubernetes,
	})
}

func sweepKubernetes(region string) error {
	_, err := sweep.SharedGoe2eClientForTests()
	if err != nil {
		log.Printf("[WARNING] %v - skipping sweep", err)
		return nil
	}

	// Note: This assumes there's a method to list all Kubernetes clusters
	// If such a method doesn't exist, we can't sweep them automatically
	log.Printf("[INFO] Kubernetes cluster sweeping would require list API method")
	log.Printf("[INFO] Please manually delete Kubernetes clusters with prefix '%s' if needed", testNamePrefix)

	// TODO: Implement when ListKubernetesClusters API is available
	// clusters, err := client.ListKubernetesClusters(projectID, location)
	// if err != nil {
	// 	return fmt.Errorf("error listing Kubernetes clusters: %w", err)
	// }
	//
	// log.Printf("[DEBUG] Found %d Kubernetes clusters in total", len(clusters))
	//
	// sweptCount := 0
	// for _, cluster := range clusters {
	// 	if !strings.HasPrefix(cluster.Name, testNamePrefix) {
	// 		log.Printf("[DEBUG] Skipping Kubernetes cluster %s (does not have test prefix)", cluster.Name)
	// 		continue
	// 	}
	//
	// 	log.Printf("[INFO] Deleting Kubernetes cluster: %s (ID: %s)", cluster.Name, cluster.ID)
	//
	// 	err := client.DeleteKubernetesService(cluster.ID, location, projectID)
	// 	if err != nil {
	// 		log.Printf("[ERROR] Failed to delete Kubernetes cluster %s: %v", cluster.Name, err)
	// 		continue
	// 	}
	//
	// 	sweptCount++
	// 	log.Printf("[INFO] Successfully deleted Kubernetes cluster: %s", cluster.Name)
	// }
	//
	// log.Printf("[INFO] Swept %d Kubernetes clusters", sweptCount)

	return nil
}
