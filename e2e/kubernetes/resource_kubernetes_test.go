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

func TestAccE2EKubernetes_Basic(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1")),
			},
		},
	})
}

func TestAccE2EKubernetes_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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

func TestAccE2EKubernetes_ImportBasic(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "status"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
				),
			},
			{
				ResourceName:            "e2e_kubernetes.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccE2EKubernetesImportIDBasic("e2e_kubernetes.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func TestAccE2EKubernetes_ImportFull(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "status"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
				),
			},
			{
				ResourceName:            "e2e_kubernetes.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccE2EKubernetesImportIDFull("e2e_kubernetes.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func TestAccE2EKubernetes_ImportMultipleNodePools(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_multipleNodePools(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "status"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "2"),
				),
			},
			{
				ResourceName:            "e2e_kubernetes.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccE2EKubernetesImportIDBasic("e2e_kubernetes.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func TestAccE2EKubernetes_NameChange(t *testing.T) {
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	clusterNameUpdated := fmt.Sprintf("test-k8s-updated-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", clusterName)),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", clusterNameUpdated)),
			},
		},
	})
}

func TestAccE2EKubernetes_VersionChange(t *testing.T) {
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	version1 := os.Getenv("E2E_TEST_K8S_VERSION")
	version2 := os.Getenv("E2E_TEST_K8S_VERSION_ALT")

	if version2 == "" {
		t.Skip("E2E_TEST_K8S_VERSION_ALT must be set for version change test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_withVersion(clusterName, version1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "version", version1)),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withVersion(clusterName, version2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "version", version2)),
			},
		},
	})
}

func TestAccE2EKubernetes_VPCChange(t *testing.T) {
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	vpc1 := os.Getenv("E2E_TEST_VPC_ID")
	vpc2 := os.Getenv("E2E_TEST_VPC_ID_ALT")

	if vpc2 == "" {
		t.Skip("E2E_TEST_VPC_ID_ALT must be set for VPC change test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_withVPC(clusterName, vpc1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "vpc_id", vpc1)),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withVPC(clusterName, vpc2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "vpc_id", vpc2)),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
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

func testAccCheckE2EKubernetesExists(resourceName string, kubernetesID *string) resource.TestCheckFunc {
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

		goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, location)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		kubernetes, _, err := goe2eClient.Kubernetes.Get(context.Background(), rs.Primary.ID)
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
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_kubernetes" {
			continue
		}

		location := acceptance.GetRegionOrLocationFromState(rs)
		projectIDStr := rs.Primary.Attributes["project_id"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, location)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		_, _, err = goe2eClient.Kubernetes.Get(context.Background(), rs.Primary.ID)
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
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

// Missing required argument configurations

func testAccCheckE2EKubernetesConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingVersion() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`,
		os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  version    = "%s"  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"),
		os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"),
		os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingVpcID() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  version    = "%s"
  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_missingNodePools() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "test-k8s"
  version    = "%s"
  vpc_id     = "%s"
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2EKubernetesConfig_withVersion(name, version string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, version, os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_withVPC(name, vpcID string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
   vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), vpcID,
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_multipleNodePools(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }

  node_pools {
    name            = "worker-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 2
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

// Import ID helper functions

func testAccE2EKubernetesImportIDBasic(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		// Simple format: just cluster_id (uses provider defaults)
		return rs.Primary.ID, nil
	}
}

func testAccE2EKubernetesImportIDFull(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		clusterID := rs.Primary.ID

		// Full format: project_id:region:cluster_id
		return fmt.Sprintf("%s:%s:%s", projectID, region, clusterID), nil
	}
}

// ============================================
// V3 FEATURES ACCEPTANCE TESTS
// ============================================

func TestAccE2EKubernetes_PreferredFields(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_preferredFields(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					// Verify deprecated fields are not set when using preferred fields
					resource.TestCheckNoResourceAttr("e2e_kubernetes.test", "name"),
					resource.TestCheckNoResourceAttr("e2e_kubernetes.test", "version"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_Tags(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_withTags(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.environment", "test"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.team", "platform"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withTagsUpdated(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.%", "3"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.environment", "test"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.team", "platform"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.updated", "true"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_NodePoolAliases(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_nodePoolAliases(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "2"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.name", "static-pool"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.type", "Static"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.size", "3"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.1.name", "autoscale-pool"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.1.type", "Autoscale"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.1.min_nodes", "2"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.1.max_nodes", "5"),
				),
			},
		},
	})
}

// ============================================
// DEPRECATION COMPATIBILITY TESTS
// ============================================

func TestAccE2EKubernetes_DeprecatedLocation(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_deprecatedLocation(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "region"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_DeprecatedName(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_deprecatedName(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					// Verify that deprecated 'name' field is still read as 'cluster_name'
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", clusterName),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_DeprecatedNodePoolFields(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_deprecatedNodePoolFields(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.name", "deprecated-pool"),
				),
			},
		},
	})
}

// ============================================
// NODE POOL OPERATIONS TESTS
// ============================================

func TestAccE2EKubernetes_AddNodePool(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_addNodePool(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "2"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.name", "default-pool"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.1.name", "new-pool"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_ResizeStaticPool(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.size", "3"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_resizeStaticPool(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.node_pool_size", "5"),
				),
			},
		},
	})
}

// ============================================
// CONFIGURATION HELPERS FOR NEW TESTS
// ============================================

func testAccCheckE2EKubernetesConfig_preferredFields(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  cluster_name       = "%s"
  kubernetes_version = "%s"
  vpc_id             = "%s"

  node_pools {
    name   = "default-pool"
    plan   = "%s"
    type   = "Static"
    size   = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_withTags(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  tags = {
    environment = "test"
    team        = "platform"
  }

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_withTagsUpdated(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  tags = {
    environment = "test"
    team        = "platform"
    updated     = "true"
  }

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_nodePoolAliases(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name = "static-pool"
    plan = "%s"
    type = "Static"
    size = 3
  }

  node_pools {
    name      = "autoscale-pool"
    plan      = "%s"
    type      = "Autoscale"
    min_nodes = 2
    max_nodes = 5
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_deprecatedLocation(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
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
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_REGION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_deprecatedName(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_deprecatedNodePoolFields(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "deprecated-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_addNodePool(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }

  node_pools {
    name            = "new-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 2
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_resizeStaticPool(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
    node_pool_size  = 5
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}
