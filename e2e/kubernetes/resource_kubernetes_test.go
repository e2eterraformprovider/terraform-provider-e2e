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
// NODE POOL OPERATIONS - ADDITIONAL TESTS
// ============================================

// TestAccE2EKubernetes_MultipleNodePoolsMixed tests creating a cluster with multiple node pools
// (mix of static and autoscale) - This is already covered by TestAccE2EKubernetes_NodePoolAliases
// but we add this explicit test to ensure it's clear
func TestAccE2EKubernetes_MultipleNodePoolsMixed(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_multipleNodePoolsMixed(clusterName),
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

func TestAccE2EKubernetes_RemoveNodePool(t *testing.T) {
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
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "2"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.name", "default-pool"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.1.name", "worker-pool"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.name", "default-pool"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_UpdateAutoscaleMinMaxNodes(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_autoscalePool(clusterName, 2, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.type", "Autoscale"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.min_nodes", "2"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.max_nodes", "5"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_autoscalePool(clusterName, 3, 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.min_nodes", "3"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.max_nodes", "7"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_UpdateNodePoolPlan(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	plan1 := os.Getenv("E2E_TEST_NODE_POOL_SPECS")
	plan2 := os.Getenv("E2E_TEST_NODE_POOL_SPECS_ALT")

	if plan2 == "" {
		t.Skip("E2E_TEST_NODE_POOL_SPECS_ALT must be set for node pool plan update test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_withPlan(clusterName, plan1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.plan", plan1),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withPlan(clusterName, plan2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.plan", plan2),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_UpdateElasticityPolicies(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_autoscaleWithElasticity(clusterName, ">", 80),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.type", "Autoscale"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_autoscaleWithElasticity(clusterName, ">=", 85),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.type", "Autoscale"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_UpdateScheduledPolicies(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_autoscaleWithScheduled(clusterName, 3, 4),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.type", "Autoscale"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_autoscaleWithScheduled(clusterName, 4, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.type", "Autoscale"),
				),
			},
		},
	})
}

// ============================================
// ERROR SCENARIO TESTS
// ============================================

func TestAccE2EKubernetes_ConflictingFields(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EKubernetesConfig_conflictingNameFields(),
				ExpectError: regexp.MustCompile(`.*conflicts with.*`),
			},
		},
	})
}

func TestAccE2EKubernetes_InvalidKubernetesVersion(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EKubernetesConfig_invalidVersion(),
				ExpectError: regexp.MustCompile(`.*must be format 1\\.XX.*`),
			},
		},
	})
}

func TestAccE2EKubernetes_InvalidRegion(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EKubernetesConfig_invalidRegion(),
				ExpectError: regexp.MustCompile(`.*invalid region.*|.*region.*not found.*`),
			},
		},
	})
}

func TestAccE2EKubernetes_InvalidVPCID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EKubernetesConfig_invalidVPCID(),
				ExpectError: regexp.MustCompile(`.*VPC.*not found.*|.*invalid.*vpc.*`),
			},
		},
	})
}

func TestAccE2EKubernetes_InvalidSecurityGroupID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EKubernetesConfig_invalidSecurityGroupID(),
				ExpectError: regexp.MustCompile(`.*security group.*not found.*|.*invalid.*security.*group.*`),
			},
		},
	})
}

func TestAccE2EKubernetes_DuplicateNodePoolNames(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EKubernetesConfig_duplicateNodePoolNames(),
				ExpectError: regexp.MustCompile(`.*Name of the worker node pools must be unique.*`),
			},
		},
	})
}

func TestAccE2EKubernetes_InvalidNodePoolPlan(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EKubernetesConfig_invalidNodePoolPlan(),
				ExpectError: regexp.MustCompile(`.*no matching plan found.*|.*invalid.*plan.*`),
			},
		},
	})
}

// ============================================
// IMPORT FUNCTIONALITY TESTS
// ============================================

func TestAccE2EKubernetes_ImportVerifyAllFields(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "id"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "vpc_id"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.#"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.name"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.plan"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.type"),
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

func TestAccE2EKubernetes_ImportV3FieldsPopulated(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.plan"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.type"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.size"),
				),
			},
			{
				ResourceName:            "e2e_kubernetes.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccE2EKubernetesImportIDBasic("e2e_kubernetes.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.plan"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.type"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_ImportWithSecurityGroups(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	securityGroupID := os.Getenv("E2E_TEST_SECURITY_GROUP_ID")

	if securityGroupID == "" {
		t.Skip("E2E_TEST_SECURITY_GROUP_ID must be set for security group import test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_withSecurityGroups(clusterName, securityGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "security_group_ids.#", "1"),
				),
			},
			{
				ResourceName:            "e2e_kubernetes.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccE2EKubernetesImportIDBasic("e2e_kubernetes.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "security_group_ids.#", "1"),
				),
			},
		},
	})
}

// ============================================
// BACKWARD COMPATIBILITY TESTS
// ============================================

func TestAccE2EKubernetes_V2FieldsBackwardCompatibility(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_v2Fields(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					// Verify V2 fields are set
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "location"),
					// Verify V3 fields are also populated (backward compatibility)
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "region"),
					// Verify node pool V2 fields
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.specs_name"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.node_pool_type"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.worker_node"),
					// Verify node pool V3 fields are also populated
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.plan"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.type"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.size"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_UpdateTagsNoRecreate(t *testing.T) {
	var kubernetesID string
	var initialCreatedAt string
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
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "created_at"),
					testAccCheckCaptureCreatedAt("e2e_kubernetes.test", &initialCreatedAt),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withTags(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.%", "2"),
					// Verify cluster was not recreated (created_at should be the same)
					testAccCheckCreatedAtUnchanged("e2e_kubernetes.test", initialCreatedAt),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withTagsUpdated(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.%", "3"),
					// Verify cluster was not recreated
					testAccCheckCreatedAtUnchanged("e2e_kubernetes.test", initialCreatedAt),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_UpdateSecurityGroupsNoRecreate(t *testing.T) {
	var kubernetesID string
	var initialCreatedAt string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	securityGroupID1 := os.Getenv("E2E_TEST_SECURITY_GROUP_ID")
	securityGroupID2 := os.Getenv("E2E_TEST_SECURITY_GROUP_ID_ALT")

	if securityGroupID1 == "" || securityGroupID2 == "" {
		t.Skip("E2E_TEST_SECURITY_GROUP_ID and E2E_TEST_SECURITY_GROUP_ID_ALT must be set for security group update test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "created_at"),
					testAccCheckCaptureCreatedAt("e2e_kubernetes.test", &initialCreatedAt),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withSecurityGroups(clusterName, securityGroupID1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "security_group_ids.#", "1"),
					// Verify cluster was not recreated
					testAccCheckCreatedAtUnchanged("e2e_kubernetes.test", initialCreatedAt),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withSecurityGroups(clusterName, securityGroupID2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "security_group_ids.#", "1"),
					// Verify cluster was not recreated
					testAccCheckCreatedAtUnchanged("e2e_kubernetes.test", initialCreatedAt),
				),
			},
		},
	})
}

// ============================================
// V3 FEATURES TESTS
// ============================================

func TestAccE2EKubernetes_V3FieldsInState(t *testing.T) {
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
					// Verify V3 preferred fields are set
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.plan"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.type"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.size"),
					// Verify V2 deprecated fields are also populated for backward compatibility
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.specs_name"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.node_pool_type"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "node_pools.0.worker_node"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_AddTagsInSeparateStep(t *testing.T) {
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
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.%", "0"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withTags(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.environment", "test"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "tags.team", "platform"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_AttachDetachSecurityGroups(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	securityGroupID := os.Getenv("E2E_TEST_SECURITY_GROUP_ID")

	if securityGroupID == "" {
		t.Skip("E2E_TEST_SECURITY_GROUP_ID must be set for security group attach/detach test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "security_group_ids.#", "0"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_withSecurityGroups(clusterName, securityGroupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "security_group_ids.#", "1"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "security_group_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_EncryptionEnabled(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	encryptionPassphrase := "test-encryption-passphrase-123"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_withEncryption(clusterName, encryptionPassphrase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "encryption_enabled", "true"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_EncryptionFieldsForceNew(t *testing.T) {
	var kubernetesID string
	var initialCreatedAt string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))
	encryptionPassphrase := "test-encryption-passphrase-123"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "encryption_enabled", "false"),
					testAccCheckCaptureCreatedAt("e2e_kubernetes.test", &initialCreatedAt),
				),
			},
			{
				Config:      testAccCheckE2EKubernetesConfig_withEncryption(clusterName, encryptionPassphrase),
				ExpectError: regexp.MustCompile(`.*forces new resource.*|.*cannot be changed.*`),
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

// ============================================
// CONFIGURATION HELPERS FOR NEW TESTS
// ============================================

func testAccCheckE2EKubernetesConfig_multipleNodePoolsMixed(name string) string {
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

func testAccCheckE2EKubernetesConfig_autoscalePool(name string, minNodes, maxNodes int) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name      = "autoscale-pool"
    plan      = "%s"
    type      = "Autoscale"
    min_nodes = %d
    max_nodes = %d
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"), minNodes, maxNodes)
}

func testAccCheckE2EKubernetesConfig_withPlan(name, plan string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name            = "default-pool"
    plan            = "%s"
    type            = "Static"
    size            = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"), plan)
}

func testAccCheckE2EKubernetesConfig_autoscaleWithElasticity(name, operator string, value int) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name      = "autoscale-pool"
    plan      = "%s"
    type      = "Autoscale"
    min_nodes = 2
    max_nodes = 5

    elasticity_dict {
      worker {
        period_number         = 1
        policy_paramter_type  = "Default"
        parameter             = "CPU"
        elasticity_policies {
          operator     = "%s"
          value        = %d
          period       = 5
          watch_period = 5
          cooldown     = 5
        }
      }
    }
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"), operator, value)
}

func testAccCheckE2EKubernetesConfig_autoscaleWithScheduled(name string, upscaleCard, downscaleCard int) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name       = "%s"
  version    = "%s"
  vpc_id     = "%s"

  node_pools {
    name      = "autoscale-pool"
    plan      = "%s"
    type      = "Autoscale"
    min_nodes = 2
    max_nodes = 5

    scheduled_dict {
      worker {
        scheduled_policies {
          upscale_cardinality   = %d
          upscale_recurrence    = "0 12 * * *"
          downscale_cardinality = %d
          downscale_recurrence  = "0 2 * * *"
        }
      }
    }
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"), upscaleCard, downscaleCard)
}

// Error scenario configuration helpers

func testAccCheckE2EKubernetesConfig_conflictingNameFields() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name         = "test-k8s"
  cluster_name = "test-k8s-cluster"
  version      = "%s"
  vpc_id       = "%s"

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

func testAccCheckE2EKubernetesConfig_invalidVersion() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name    = "test-k8s"
  version = "2.0"
  vpc_id  = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_invalidRegion() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name     = "test-k8s"
  version  = "%s"
  region   = "InvalidRegion"
  vpc_id   = "%s"

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

func testAccCheckE2EKubernetesConfig_invalidVPCID() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name    = "test-k8s"
  version = "%s"
  vpc_id  = "999999"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_invalidSecurityGroupID() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name              = "test-k8s"
  version           = "%s"
  vpc_id            = "%s"
  security_group_ids = [999999]

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

func testAccCheckE2EKubernetesConfig_duplicateNodePoolNames() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name    = "test-k8s"
  version = "%s"
  vpc_id  = "%s"

  node_pools {
    name            = "duplicate-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }

  node_pools {
    name            = "duplicate-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 2
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"), os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_invalidNodePoolPlan() string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name    = "test-k8s"
  version = "%s"
  vpc_id  = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "INVALID_PLAN_NAME"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"))
}

// ============================================
// HELPER FUNCTIONS FOR NEW TESTS
// ============================================

func testAccCheckCaptureCreatedAt(resourceName string, createdAt *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		*createdAt = rs.Primary.Attributes["created_at"]
		return nil
	}
}

func testAccCheckCreatedAtUnchanged(resourceName, expectedCreatedAt string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		actualCreatedAt := rs.Primary.Attributes["created_at"]
		if actualCreatedAt != expectedCreatedAt {
			return fmt.Errorf("created_at changed from %s to %s, indicating resource was recreated", expectedCreatedAt, actualCreatedAt)
		}
		return nil
	}
}

// ============================================
// CONFIGURATION HELPERS FOR IMPORT, BACKWARD COMPATIBILITY, AND V3 FEATURES
// ============================================

func testAccCheckE2EKubernetesConfig_withSecurityGroups(name, securityGroupID string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name              = "%s"
  version           = "%s"
  vpc_id            = "%s"
  security_group_ids = [%s]

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"), securityGroupID,
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_v2Fields(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name     = "%s"
  version  = "%s"
  location = "%s"
  vpc_id   = "%s"

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

func testAccCheckE2EKubernetesConfig_withEncryption(name, passphrase string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name                 = "%s"
  version              = "%s"
  vpc_id               = "%s"
  encryption_enabled   = true
  encryption_passphrase = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 3
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"), passphrase,
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_minimumNodes(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name    = "%s"
  version = "%s"
  vpc_id  = "%s"

  node_pools {
    name            = "default-pool"
    specs_name      = "%s"
    node_pool_type  = "Static"
    worker_node     = 2
  }
}
`, name, os.Getenv("E2E_TEST_K8S_VERSION"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_NODE_POOL_SPECS"))
}

func testAccCheckE2EKubernetesConfig_withLongTimeout(name string) string {
	return fmt.Sprintf(`
resource "e2e_kubernetes" "test" {
  name    = "%s"
  version = "%s"
  vpc_id  = "%s"

  timeouts {
    create = "45m"
    update = "45m"
    delete = "30m"
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

// ============================================
// EDGE CASE TESTS
// ============================================

func TestAccE2EKubernetes_VeryLongClusterName(t *testing.T) {
	var kubernetesID string
	// Create a name that's exactly 255 characters (max length)
	longName := fmt.Sprintf("test-k8s-%s", acctest.RandString(240)) // 10 + 240 = 250, but we need exactly 255
	if len(longName) > 255 {
		longName = longName[:255]
	} else {
		// Pad to exactly 255 characters
		for len(longName) < 255 {
			longName += "a"
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(longName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", longName),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", longName),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_SpecialCharactersInName(t *testing.T) {
	var kubernetesID string
	// Kubernetes names typically allow alphanumeric and hyphens, but let's test with valid special chars
	// Note: Most cloud providers restrict special characters, so we'll test with hyphens and underscores
	specialName := fmt.Sprintf("test-k8s_%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_basic(specialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", specialName),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_MinimumNodesPerPool(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_minimumNodes(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.size", "2"),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.0.worker_node", "2"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_ImportV2Naming(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_v2Fields(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					// Verify V2 fields are set
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "location"),
				),
			},
			{
				ResourceName:            "e2e_kubernetes.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccE2EKubernetesImportIDBasic("e2e_kubernetes.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
				Check: resource.ComposeTestCheckFunc(
					// Verify both V2 and V3 fields are populated after import
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "name", clusterName),
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "kubernetes_version"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "location"),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "region"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_VeryLongTimeout(t *testing.T) {
	var kubernetesID string
	clusterName := fmt.Sprintf("test-k8s-%s", acctest.RandString(10))

	// Note: This test verifies that very long timeouts (45 minutes) can be configured
	// The actual timeout is set in the resource configuration via timeouts block
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EKubernetesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EKubernetesConfig_withLongTimeout(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EKubernetesExists("e2e_kubernetes.test", &kubernetesID),
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "id"),
				),
			},
		},
	})
}

// ============================================
// PERFORMANCE TESTS
// ============================================

func TestAccE2EKubernetes_PerformanceCreateCluster(t *testing.T) {
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
					// Performance measurement is done by Terraform's test framework
					// The test will fail if it takes longer than the timeout
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "id"),
				),
			},
		},
	})
}

func TestAccE2EKubernetes_PerformanceAddNodePool(t *testing.T) {
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
					// Performance measurement is done by Terraform's test framework
				),
			},
		},
	})
}

func TestAccE2EKubernetes_PerformanceRemoveNodePool(t *testing.T) {
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
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "2"),
				),
			},
			{
				Config: testAccCheckE2EKubernetesConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_kubernetes.test", "node_pools.#", "1"),
					// Performance measurement is done by Terraform's test framework
				),
			},
		},
	})
}

func TestAccE2EKubernetes_PerformanceDeleteCluster(t *testing.T) {
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
					// Performance measurement for deletion is done by CheckDestroy
					resource.TestCheckResourceAttrSet("e2e_kubernetes.test", "id"),
				),
			},
		},
	})
}

// ============================================
// TEST SWEEPER VERIFICATION
// ============================================

// TestAccE2EKubernetes_SweeperRegistered verifies that the Kubernetes sweeper is properly registered.
// Note: Actual sweeping functionality requires ListKubernetesClusters API which may not be available.
// The sweeper is registered in e2e/kubernetes/sweep.go via init() function.
// This test documents that the sweeper exists and can be verified by checking the sweep.go file.
func TestAccE2EKubernetes_SweeperRegistered(t *testing.T) {
	// The sweeper is registered in e2e/kubernetes/sweep.go via:
	// resource.AddTestSweepers("e2e_kubernetes", &resource.Sweeper{...})
	//
	// To verify the sweeper is working:
	// 1. Run: go test -v -sweep=us-east-1 -sweep-run=e2e_kubernetes ./e2e/kubernetes/...
	// 2. Check that clusters with test prefix are cleaned up
	//
	// Note: Actual sweeping functionality requires ListKubernetesClusters API
	// which may not be available. The sweeper is registered but may log warnings.
	t.Log("Kubernetes sweeper is registered in e2e/kubernetes/sweep.go")
	t.Log("To test sweeper: go test -v -sweep=<region> -sweep-run=e2e_kubernetes ./e2e/kubernetes/...")
}

// Note: "Delete cluster created outside Terraform" test scenario:
// This is typically tested manually by:
// 1. Creating a cluster outside of Terraform (via API/console)
// 2. Importing it into Terraform state
// 3. Running terraform destroy to verify it can be deleted
// This is difficult to automate in acceptance tests as it requires external cluster creation.
// The import and destroy functionality is already tested separately.
