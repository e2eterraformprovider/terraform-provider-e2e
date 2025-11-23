package kubernetes

import (
	"fmt"
	"log"
	"os"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testNamePrefix = "test-k8s-"

func init() {
	resource.AddTestSweepers("e2e_kubernetes", &resource.Sweeper{
		Name: "e2e_kubernetes",
		F:    sweepKubernetes,
	})
}

func sweepKubernetes(region string) error {
	cfg, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting config for region %s: %w", region, err)
	}

	_ = cfg.Client()

	projectID := os.Getenv("E2E_TEST_PROJECT_ID")
	location := os.Getenv("E2E_TEST_LOCATION")

	if projectID == "" || location == "" {
		log.Printf("[WARNING] E2E_TEST_PROJECT_ID or E2E_TEST_LOCATION not set, skipping sweep")
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

func sharedConfigForRegion(region string) (*config.Config, error) {
	apiKey := os.Getenv("SERVICE_API_KEY")
	authToken := os.Getenv("SERVICE_AUTH_TOKEN")
	apiEndpoint := os.Getenv("SERVICE_API_ENDPOINT")

	if apiKey == "" {
		return nil, fmt.Errorf("SERVICE_API_KEY must be set for acceptance tests")
	}

	if authToken == "" {
		return nil, fmt.Errorf("SERVICE_AUTH_TOKEN must be set for acceptance tests")
	}

	if apiEndpoint == "" {
		apiEndpoint = "https://api.e2enetworks.com/myaccount/api/v1/"
	}

	cfg, err := config.NewConfig(apiKey, authToken, apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error creating config: %w", err)
	}

	return cfg, nil
}
