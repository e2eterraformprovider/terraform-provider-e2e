package kubernetes_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceE2EKubernetes_Basic(t *testing.T) {
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EKubernetesExists("data.e2e_kubernetes.test"),
					resource.TestCheckResourceAttrPair("data.e2e_kubernetes.test", "service_id", "e2e_kubernetes.test", "id"),
					resource.TestCheckResourceAttrPair("data.e2e_kubernetes.test", "name", "e2e_kubernetes.test", "name"),
					resource.TestCheckResourceAttrPair("data.e2e_kubernetes.test", "version", "e2e_kubernetes.test", "version"),
					resource.TestCheckResourceAttrSet("data.e2e_kubernetes.test", "status"),
					resource.TestCheckResourceAttrSet("data.e2e_kubernetes.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.e2e_kubernetes.test", "master_node_id")),
			},
		},
	})
}

func TestAccDataSourceE2EKubernetes_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2EKubernetesConfig_missingServiceID(),
				ExpectError: regexp.MustCompile(`The argument "service_id" is required`),
			},
			{
				Config:      testAccDataSourceE2EKubernetesConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccDataSourceE2EKubernetesConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccDataSourceE2EKubernetes_NotFound(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2EKubernetesConfig_notFound(),
				ExpectError: regexp.MustCompile(`error finding Item with ID`),
			},
		},
	})
}

// Helper functions

func testAccCheckDataSourceE2EKubernetesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Kubernetes ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		location := acceptance.GetRegionOrLocationFromState(rs)
		projectIDStr := rs.Primary.Attributes["project_id"]
		serviceID := rs.Primary.Attributes["service_id"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, location)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		kubernetes, _, err := goe2eClient.Kubernetes.Get(context.Background(), serviceID)
		if err != nil {
			return err
		}

		if kubernetes == nil {
			return fmt.Errorf("Kubernetes cluster not found in datasource: %s", serviceID)
		}

		return nil
	}
}

// Configuration helpers

func testAccDataSourceE2EKubernetesConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name    = "%s"
  version = "%s"
  vpc_id  = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}

data "e2e_kubernetes" "test" {
  service_id = e2e_kubernetes.test.id
  project_id = e2e_kubernetes.test.project_id
  location   = e2e_kubernetes.test.location
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"),
		os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

// Error case configurations

func testAccDataSourceE2EKubernetesConfig_missingServiceID() string {
	return `
data "e2e_kubernetes" "test" {
}
`
}

func testAccDataSourceE2EKubernetesConfig_missingProjectID() string {
	return `
data "e2e_kubernetes" "test" {
  service_id = "12345"
}
`
}

func testAccDataSourceE2EKubernetesConfig_missingLocation() string {
	return `
data "e2e_kubernetes" "test" {
  service_id = "12345"
}
`
}

func testAccDataSourceE2EKubernetesConfig_notFound() string {
	return `
data "e2e_kubernetes" "test" {
  service_id = "999999999"
}
`
}
