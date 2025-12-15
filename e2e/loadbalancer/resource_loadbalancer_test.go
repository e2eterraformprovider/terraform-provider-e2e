package loadbalancer_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbName, lbName),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPlan, goe2econstants.LBPlanE2ELB2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbMode, goe2econstants.LBModeHTTP),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbType, goe2econstants.LBTypeExternal),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "node_list_type", goe2econstants.LBNodeListTypeStatic),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOn),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "enable_bitninja", "false"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "is_ipv6_attached", "false"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "public_ip"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "private_ip"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrRAM),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrDisk),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrVCPU),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrStatus)),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbName, lbName)),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_updated(lbNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbName, lbNameUpdated)),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbMode, goe2econstants.LBModeHTTPS)),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrBackends+".#", "1"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrBackends+".0.name", "backend-1"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrBackends+".0.balance", goe2econstants.LBBalanceRoundRobin)),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOn)),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_powerOff(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOff)),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPlan, goe2econstants.LBPlanE2ELB2)),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_upgradedPlan(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPlan, goe2econstants.LBPlanE2ELB3)),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbMode, goe2econstants.LBModeHTTP)),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_https(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbMode, goe2econstants.LBModeHTTPS),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbType, goe2econstants.LBTypeExternal)),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_internal(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbType, goe2econstants.LBTypeInternal),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "node_list_type", goe2econstants.LBNodeListTypeStatic)),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_dynamicNodes(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "node_list_type", goe2econstants.LBNodeListTypeDynamic),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPlan, goe2econstants.LBPlanE2ELB3)),
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
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPlan, goe2econstants.LBPlanE2ELB2)),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_planForceNew(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID2),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPlan, goe2econstants.LBPlanE2ELB4),
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
  plan   = "%s"
  lb_mode     = "%s"}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "%s"
  lb_mode     = "%s"}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_https(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "%s"
  lb_mode     = "%s"}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTPS)
}

func testAccCheckE2ELoadBalancerConfig_withBackends(lbName, nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "%s"
  lb_mode     = "%s"
  backends {
    name    = "backend-1"
    balance = "%s"
    servers {
      id   = e2e_node.backend.id
      port = "%s"
    }
  }
}
`, nodeName, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP,
		goe2econstants.LBBalanceRoundRobin, goe2econstants.LBPortHTTP)
}

func testAccCheckE2ELoadBalancerConfig_powerOff(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name      = "%s"
  plan    = "%s"
  lb_mode      = "%s"
  power_status = "%s"}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.NodePowerStatusOff)
}

func testAccCheckE2ELoadBalancerConfig_upgradedPlan(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "%s"
  lb_mode     = "%s"}
`, name, goe2econstants.LBPlanE2ELB3, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_withIPv6(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name          = "%s"
  plan        = "%s"
  lb_mode          = "%s"
  is_ipv6_attached = true}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
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
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name    = "test-lb"
  lb_mode    = "%s"}
`, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_missingMode() string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name    = "test-lb"
  plan  = "%s"}
`, goe2econstants.LBPlanE2ELB2)
}

func testAccCheckE2ELoadBalancerConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name   = "test-lb"
  plan = "%s"
  lb_mode   = "%s"}
`, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name    = "test-lb"
  plan  = "%s"
  lb_mode    = "%s"}
`, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_invalidPlan(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name    = "%s"
  plan  = "INVALID-PLAN"
  lb_mode    = "%s"}
`, name, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_invalidMode(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name    = "%s"
  plan  = "%s"
  lb_mode    = "INVALID"}
`, name, goe2econstants.LBPlanE2ELB2)
}

func testAccCheckE2ELoadBalancerConfig_internal(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "%s"
  lb_mode     = "%s"
  lb_type     = "%s"}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBTypeInternal)
}

func testAccCheckE2ELoadBalancerConfig_dynamicNodes(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name        = "%s"
  plan      = "%s"
  lb_mode        = "%s"
  node_list_type = "%s"}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBNodeListTypeDynamic)
}

func testAccCheckE2ELoadBalancerConfig_planForceNew(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name     = "%s"
  plan   = "%s"
  lb_mode     = "%s"}
`, name, goe2econstants.LBPlanE2ELB4, goe2econstants.LBModeHTTP)
}
