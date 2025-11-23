package kubernetes_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EKubernetes_Basic(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "created_at"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EKubernetesConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EKubernetesConfig_missingVersion(),
				ExpectError: regexp.MustCompile(`The argument "version" is required`),
			},
			{
				Config:      testAccCheckE2EKubernetesConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EKubernetesConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
			{
				Config:      testAccCheckE2EKubernetesConfig_missingVpcID(),
				ExpectError: regexp.MustCompile(`The argument "vpc_id" is required`),
			},
			{
				Config:      testAccCheckE2EKubernetesConfig_missingNodePools(),
				ExpectError: regexp.MustCompile(`The argument "node_pools" is required`),
			},
		},
	})
}

func TestAccE2EKubernetes_Import(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
				),
			},
			{
				ResourceName:      "e2e_kubernetes.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Helper functions

var testAccProvider *schema.Provider

func init() {
	testAccProvider = e2e.Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SERVICE_API_KEY"); v == "" {
		t.Fatal("SERVICE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICE_AUTH_TOKEN"); v == "" {
		t.Fatal("SERVICE_AUTH_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_PROJECT_ID"); v == "" {
		t.Fatal("E2E_TEST_PROJECT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_LOCATION"); v == "" {
		t.Fatal("E2E_TEST_LOCATION must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_VPC_ID"); v == "" {
		t.Fatal("E2E_TEST_VPC_ID must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_K8S_VERSION"); v == "" {
		t.Fatal("E2E_TEST_K8S_VERSION must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_NODE_POOL_SPECS"); v == "" {
		t.Fatal("E2E_TEST_NODE_POOL_SPECS must be set for acceptance tests")
	}
}

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"e2e": func() (*schema.Provider, error) {
		return e2e.Provider(), nil
	},
}

func testAccCheckE2EKubernetesExists(resourceName string, kubernetesID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Kubernetes ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		location := rs.Primary.Attributes["location"]
		projectIDStr := rs.Primary.Attributes["project_id"]
		projectID := 0
		fmt.Sscanf(projectIDStr, "%d", &projectID)

		kubernetes, err := client.GetKubernetesServiceInfo(rs.Primary.ID, location, projectID)
		if err != nil {
			return err
		}

		if kubernetes == nil {
			return fmt.Errorf("Kubernetes cluster not found")
		}

		*kubernetesID = rs.Primary.ID

		return nil
	}
}

func testAccCheckE2EKubernetesDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_kubernetes" {
			continue
		}

		location := rs.Primary.Attributes["location"]
		projectIDStr := rs.Primary.Attributes["project_id"]
		projectID := 0
		fmt.Sscanf(projectIDStr, "%d", &projectID)

		_, err := client.GetKubernetesServiceInfo(rs.Primary.ID, location, projectID)
		if err == nil {
			return fmt.Errorf("Kubernetes cluster still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2EKubernetesConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  project_id = %s
  location   = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_PROJECT_ID"),
		os.Getenv("E2E_TEST_LOCATION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

// Missing required argument configurations

func testAccCheckE2EKubernetesConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  version    = "%s"
  project_id = %s
  location   = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_PROJECT_ID"),
		os.Getenv("E2E_TEST_LOCATION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingVersion() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  project_id = %s
  location   = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  version    = "%s"
  location   = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  version    = "%s"
  project_id = %s
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_PROJECT_ID"),
		os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingVpcID() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  version    = "%s"
  project_id = %s
  location   = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_PROJECT_ID"),
		os.Getenv("E2E_TEST_LOCATION"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingNodePools() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  version    = "%s"
  project_id = %s
  location   = "%s"
  vpc_id     = "%s"
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_PROJECT_ID"),
		os.Getenv("E2E_TEST_LOCATION"), os.Getenv("E2E_TEST_VPC_ID"))
}
