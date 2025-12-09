package loadbalancer_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2ELoadBalancer_Basic(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_name", lbName),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "plan", "E2E-LB-2"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_mode", "HTTP"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_type", "External"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "node_list_type", "S"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "power_status", "power_on"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "enable_bitninja", "false"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "is_ipv6_attached", "false"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "public_ip"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "private_ip"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "ram"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "disk"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "vcpu"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "status")),
			},
		},
	})
}

func TestAccE2ELoadBalancer_Update(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	lbNameUpdated := fmt.Sprintf("test-lb-updated-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_name", lbName)),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_updated(lbNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_name", lbNameUpdated)),
			},
		},
	})
}

func TestAccE2ELoadBalancer_HTTPS(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_https(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_mode", "HTTPS")),
			},
		},
	})
}

func TestAccE2ELoadBalancer_WithBackends(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_withBackends(lbName, nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "backends.#", "1"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "backends.0.name", "backend-1"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "backends.0.balance", "roundrobin")),
			},
		},
	})
}

func TestAccE2ELoadBalancer_PowerOperations(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "power_status", "power_on")),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_powerOff(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "power_status", "power_off")),
			},
		},
	})
}

func TestAccE2ELoadBalancer_PlanUpgrade(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "plan", "E2E-LB-2")),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_upgradedPlan(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "plan", "E2E-LB-3")),
			},
		},
	})
}

func TestAccE2ELoadBalancer_IPv6(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_withIPv6(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "is_ipv6_attached", "true")),
			},
		},
	})
}

func TestAccE2ELoadBalancer_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ELoadBalancerConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "lb_name" is required`),
			},
			{
				Config:      testAccCheckE2ELoadBalancerConfig_missingPlan(),
				ExpectError: regexp.MustCompile(`The argument "plan" is required`),
			},
			{
				Config:      testAccCheckE2ELoadBalancerConfig_missingMode(),
				ExpectError: regexp.MustCompile(`The argument "lb_mode" is required`),
			},
			{
				Config:      testAccCheckE2ELoadBalancerConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2ELoadBalancerConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccE2ELoadBalancer_InvalidPlan(t *testing.T) {
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ELoadBalancerConfig_invalidPlan(lbName),
				ExpectError: regexp.MustCompile(`expected plan to be one of`),
			},
		},
	})
}

func TestAccE2ELoadBalancer_InvalidMode(t *testing.T) {
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ELoadBalancerConfig_invalidMode(lbName),
				ExpectError: regexp.MustCompile(`expected lb_mode to be one of`),
			},
		},
	})
}

func TestAccE2ELoadBalancer_Import(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID)),
			},
			{
				ResourceName:            "e2e_loadbalancer.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccE2ELoadBalancerImportID("e2e_loadbalancer.test"),
				ImportStateVerifyIgnore: []string{"ssl_context", "enable_eos_logger"},
			},
		},
	})
}

func TestAccE2ELoadBalancer_ForceNewLbMode(t *testing.T) {
	var lbID1, lbID2 string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID1),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_mode", "HTTP")),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_https(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_mode", "HTTPS"),
					// Verify that resource was recreated (different ID)
					testAccCheckE2ELoadBalancerRecreated(&lbID1, &lbID2)),
			},
		},
	})
}

func TestAccE2ELoadBalancer_ForceNewLbType(t *testing.T) {
	var lbID1, lbID2 string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID1),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_type", "External")),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_internal(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "lb_type", "Internal"),
					// Verify that resource was recreated (different ID)
					testAccCheckE2ELoadBalancerRecreated(&lbID1, &lbID2)),
			},
		},
	})
}

func TestAccE2ELoadBalancer_ForceNewNodeListType(t *testing.T) {
	var lbID1, lbID2 string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID1),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "node_list_type", "S")),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_dynamicNodes(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "node_list_type", "D"),
					// Verify that resource was recreated (different ID)
					testAccCheckE2ELoadBalancerRecreated(&lbID1, &lbID2)),
			},
		},
	})
}

func TestAccE2ELoadBalancer_PlanDowngrade(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_upgradedPlan(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "plan", "E2E-LB-3")),
			},
			{
				Config:      testAccCheckE2ELoadBalancerConfig_basic(lbName),
				ExpectError: regexp.MustCompile(`Cannot downgrade plan`),
			},
		},
	})
}

func TestAccE2ELoadBalancer_ForceNewPlan(t *testing.T) {
	var lbID1, lbID2 string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_basic(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID1),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "plan", "E2E-LB-2")),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_planForceNew(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "plan", "E2E-LB-4"),
					// Verify that resource was recreated (different ID) due to ForceNew
					testAccCheckE2ELoadBalancerRecreated(&lbID1, &lbID2)),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2ELoadBalancerExists(resourceName string, lbID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LoadBalancer ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)

		// Use goe2e client
		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("Error creating goe2e client: %s", err)
		}

		lb, _, err := goe2eClient.LoadBalancer.GetLoadBalancer(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if lb == nil {
			return fmt.Errorf("LoadBalancer not found")
		}

		*lbID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2ELoadBalancerDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_loadbalancer" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)

		// Use goe2e client
		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("Error creating goe2e client: %s", err)
		}

		lb, _, err := goe2eClient.LoadBalancer.GetLoadBalancer(context.Background(), rs.Primary.ID)
		if err == nil && lb != nil {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccE2ELoadBalancerImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		lbID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, region, lbID), nil
	}
}

func testAccCheckE2ELoadBalancerRecreated(oldID, newID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *oldID == *newID {
			return fmt.Errorf("Expected load balancer to be recreated, but IDs are the same: %s", *oldID)
		}
		return nil
	}
}

// Configuration helpers

func testAccCheckE2ELoadBalancerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "E2E-LB-2"
  lb_mode     = "HTTP"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "E2E-LB-2"
  lb_mode     = "HTTP"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_https(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "E2E-LB-2"
  lb_mode     = "HTTPS"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_withBackends(lbName, nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "E2E-LB-2"
  lb_mode     = "HTTP"
  backends {
    name    = "backend-1"
    balance = "roundrobin"
    servers {
      id   = e2e_node.backend.id
      port = "80"
    }
  }
}
`, nodeName,
		lbName)
}

func testAccCheckE2ELoadBalancerConfig_powerOff(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name      = "%s"
  plan    = "E2E-LB-2"
  lb_mode      = "HTTP"
  power_status = "power_off"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_upgradedPlan(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "E2E-LB-3"
  lb_mode     = "HTTP"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_withIPv6(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name          = "%s"
  plan        = "E2E-LB-2"
  lb_mode          = "HTTP"
  is_ipv6_attached = true}
`, name)
}

// Error case configurations

func testAccCheckE2ELoadBalancerConfig_missingName() string {
	return `
resource "e2e_loadbalancer" "test" {
  plan  = "E2E-LB-2"
  lb_mode    = "HTTP"}
`
}

func testAccCheckE2ELoadBalancerConfig_missingPlan() string {
	return `
resource "e2e_loadbalancer" "test" {
  lb_name    = "test-lb"
  lb_mode    = "HTTP"}
`
}

func testAccCheckE2ELoadBalancerConfig_missingMode() string {
	return `
resource "e2e_loadbalancer" "test" {
  lb_name    = "test-lb"
  plan  = "E2E-LB-2"}
`
}

func testAccCheckE2ELoadBalancerConfig_missingProjectID() string {
	return `
resource "e2e_loadbalancer" "test" {
  lb_name   = "test-lb"
  plan = "E2E-LB-2"
  lb_mode   = "HTTP"}
`
}

func testAccCheckE2ELoadBalancerConfig_missingLocation() string {
	return `
resource "e2e_loadbalancer" "test" {
  lb_name    = "test-lb"
  plan  = "E2E-LB-2"
  lb_mode    = "HTTP"}
`
}

func testAccCheckE2ELoadBalancerConfig_invalidPlan(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name    = "%s"
  plan  = "INVALID-PLAN"
  lb_mode    = "HTTP"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_invalidMode(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name    = "%s"
  plan  = "E2E-LB-2"
  lb_mode    = "INVALID"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_internal(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "E2E-LB-2"
  lb_mode     = "HTTP"
  lb_type     = "Internal"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_dynamicNodes(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name        = "%s"
  plan      = "E2E-LB-2"
  lb_mode        = "HTTP"
  node_list_type = "D"}
`, name)
}

func testAccCheckE2ELoadBalancerConfig_planForceNew(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "E2E-LB-4"
  lb_mode     = "HTTP"}
`, name)
}
