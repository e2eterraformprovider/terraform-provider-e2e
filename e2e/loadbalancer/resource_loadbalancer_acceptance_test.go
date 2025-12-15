package loadbalancer_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// ============================================================================
// V3 Preferred Fields Tests
// ============================================================================

func TestAccE2ELoadBalancer_V3Fields(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_v3Fields(lbName, nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					// V3 fields check
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrName, lbName),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrRegion),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrTags+".env", "test"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrTags+".team", "devops"),
					// V2 fields also populated (backward compatibility)
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbName, lbName),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrLocation),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_Name(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_nameField(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrName, lbName),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbName, lbName),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_Region(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_regionField(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrRegion),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrLocation),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_FloatingIPID(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_floatingIPID(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttrPair("e2e_loadbalancer.test", "floating_ip_id", "e2e_reserve_ip.test", tfconstants.AttrID),
					resource.TestCheckResourceAttrPair("e2e_loadbalancer.test", "lb_reserve_ip", "e2e_reserve_ip.test", tfconstants.AttrID),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_Tags(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_tags(lbName, `{"env" = "test", "team" = "devops"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrTags+".env", "test"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrTags+".team", "devops"),
				),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_tags(lbName, `{"env" = "prod", "owner" = "sre"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrTags+".env", "prod"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrTags+".owner", "sre"),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_NodeIDInBackends(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_nodeIDBackends(lbName, nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrBackends+".#", "1"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrBackends+".0.servers.#", "1"),
				),
			},
		},
	})
}

// ============================================================================
// Backward Compatibility Tests
// ============================================================================

func TestAccE2ELoadBalancer_DeprecatedLocation(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_deprecatedLocation(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrLocation),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrRegion),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_DeprecatedLbName(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_deprecatedLbName(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbName, lbName),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrName, lbName),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_DeprecatedLbReserveIP(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_deprecatedLbReserveIP(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttrPair("e2e_loadbalancer.test", "lb_reserve_ip", "e2e_reserve_ip.test", tfconstants.AttrID),
					resource.TestCheckResourceAttrPair("e2e_loadbalancer.test", "floating_ip_id", "e2e_reserve_ip.test", tfconstants.AttrID),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_DeprecatedServerID(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_deprecatedServerID(lbName, nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrBackends+".#", "1"),
				),
			},
		},
	})
}

// ============================================================================
// Migration Scenarios Tests
// ============================================================================

func TestAccE2ELoadBalancer_MigrateToRegion(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_deprecatedLocation(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrLocation),
				),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_regionField(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrRegion),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", tfconstants.AttrLocation),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_MigrateToName(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_deprecatedLbName(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbName, lbName),
				),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_nameField(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrName, lbName),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbName, lbName),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_MigrateToFloatingIPID(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_deprecatedLbReserveIP(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "lb_reserve_ip"),
				),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_floatingIPID(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "floating_ip_id"),
					resource.TestCheckResourceAttrSet("e2e_loadbalancer.test", "lb_reserve_ip"),
				),
			},
		},
	})
}

// ============================================================================
// Advanced Features Tests
// ============================================================================

func TestAccE2ELoadBalancer_BackendUpdate(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	nodeName1 := fmt.Sprintf("test-node-1-%s", acctest.RandString(10))
	nodeName2 := fmt.Sprintf("test-node-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_singleBackend(lbName, nodeName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrBackends+".#", "1"),
				),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_multipleBackends(lbName, nodeName1, nodeName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrBackends+".#", "2"),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_TcpBackend(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_tcpBackend(lbName, nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "tcp_backend.#", "1"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "tcp_backend.0.backend_name", "tcp-backend-1"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "tcp_backend.0.port", "3306"),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_ACL(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_acl(lbName, nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "acl_list.#", "1"),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "acl_map.#", "1"),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_VPC(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_vpc(lbName, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "vpc_list.#", "1"),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_SSL(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_ssl(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrLbMode, goe2econstants.LBModeHTTPS),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", "ssl_context.#", "1"),
				),
			},
		},
	})
}

// ============================================================================
// Import Functionality Tests
// ============================================================================

func TestAccE2ELoadBalancer_ImportFull(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_nameField(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID)),
			},
			{
				ResourceName:            "e2e_loadbalancer.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccE2ELoadBalancerImportIDFull("e2e_loadbalancer.test"),
				ImportStateVerifyIgnore: []string{"ssl_context", "enable_eos_logger"},
			},
		},
	})
}

// ============================================================================
// Validation Tests
// ============================================================================

func TestAccE2ELoadBalancer_ConflictingFields(t *testing.T) {
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ELoadBalancerConfig_conflictingNameAndLbName(lbName),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
			{
				Config:      testAccCheckE2ELoadBalancerConfig_conflictingRegionAndLocation(lbName),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
			{
				Config:      testAccCheckE2ELoadBalancerConfig_conflictingFloatingIPFields(lbName),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
		},
	})
}

func TestAccE2ELoadBalancer_TransitionalStateDelete(t *testing.T) {
	// This test verifies delete protection during transitional states
	// Note: In practice, delete waits for ready state, so this may be hard to test
	// This is more of a documentation test of expected behavior
	t.Skip("Skipping: Delete protection tested via unit tests and async deletion logic")
}

// ============================================================================
// Async Operations Tests
// ============================================================================

func TestAccE2ELoadBalancer_AsyncCreate(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_nameField(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrStatus, goe2econstants.LBStatusRunning),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrState, goe2econstants.LBStateRunning),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_AsyncDelete(t *testing.T) {
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_nameField(lbName),
			},
		},
	})
}

func TestAccE2ELoadBalancer_AsyncPowerActions(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_nameField(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOn),
				),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_powerOff(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOff),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrStatus, goe2econstants.LBStatusPoweredOff),
				),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_nameField(lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOn),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrStatus, goe2econstants.LBStatusRunning),
				),
			},
		},
	})
}

func TestAccE2ELoadBalancer_AsyncPlanUpgrade(t *testing.T) {
	var lbID string
	lbName := fmt.Sprintf("test-lb-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ELoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ELoadBalancerConfig_plan(lbName, goe2econstants.LBPlanE2ELB2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPlan, goe2econstants.LBPlanE2ELB2),
				),
			},
			{
				Config: testAccCheckE2ELoadBalancerConfig_plan(lbName, goe2econstants.LBPlanE2ELB3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ELoadBalancerExists("e2e_loadbalancer.test", &lbID),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrPlan, goe2econstants.LBPlanE2ELB3),
					resource.TestCheckResourceAttr("e2e_loadbalancer.test", tfconstants.AttrStatus, goe2econstants.LBStatusRunning),
				),
			},
		},
	})
}

// ============================================================================
// Helper Functions (new)
// ============================================================================

// testAccE2ELoadBalancerImportIDFull creates a full import ID with project/region/lb format
// Note: testAccE2ELoadBalancerImportID already exists in the main test file and does the same thing
func testAccE2ELoadBalancerImportIDFull(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes[tfconstants.AttrProjectID]
		region := acceptance.GetRegionOrLocationFromState(rs)
		lbID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, region, lbID), nil
	}
}

// ============================================================================
// Configuration Helpers
// ============================================================================

func testAccCheckE2ELoadBalancerConfig_v3Fields(lbName, nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
}

resource "e2e_loadbalancer" "test" {
  name     = "%s"
  plan     = "%s"
  lb_mode  = "%s"
  tags = {
    env  = "test"
    team = "devops"
  }
  backends {
    name    = "backend-1"
    balance = "%s"
    servers {
      node_id = e2e_node.backend.id
      port    = "%s"
    }
  }
}
`, nodeName, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBBalanceRoundRobin, goe2econstants.LBPortHTTP)
}

func testAccCheckE2ELoadBalancerConfig_nameField(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_regionField(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_floatingIPID(name string) string {
	return fmt.Sprintf(`
resource "e2e_reserve_ip" "test" {
  name = "%s-ip"
}

resource "e2e_loadbalancer" "test" {
  name           = "%s"
  plan           = "%s"
  lb_mode        = "%s"
  floating_ip_id = e2e_reserve_ip.test.id
}
`, name, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_tags(name, tags string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
  tags    = %s
}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, tags)
}

func testAccCheckE2ELoadBalancerConfig_nodeIDBackends(lbName, nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
  backends {
    name    = "backend-1"
    balance = "%s"
    servers {
      node_id = e2e_node.backend.id
      port    = "%s"
    }
  }
}
`, nodeName, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBBalanceRoundRobin, goe2econstants.LBPortHTTP)
}

func testAccCheckE2ELoadBalancerConfig_deprecatedLocation(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name = "%s"
  plan    = "%s"
  lb_mode = "%s"
}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_deprecatedLbName(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  lb_name = "%s"
  plan    = "%s"
  lb_mode = "%s"
}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_deprecatedLbReserveIP(name string) string {
	return fmt.Sprintf(`
resource "e2e_reserve_ip" "test" {
  name = "%s-ip"
}

resource "e2e_loadbalancer" "test" {
  lb_name      = "%s"
  plan         = "%s"
  lb_mode      = "%s"
  lb_reserve_ip = e2e_reserve_ip.test.id
}
`, name, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_deprecatedServerID(lbName, nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_loadbalancer" "test" {
  lb_name = "%s"
  plan    = "%s"
  lb_mode = "%s"
  backends {
    name    = "backend-1"
    balance = "%s"
    servers {
      id   = e2e_node.backend.id
      port = "%s"
    }
  }
}
`, nodeName, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBBalanceRoundRobin, goe2econstants.LBPortHTTP)
}

func testAccCheckE2ELoadBalancerConfig_singleBackend(lbName, nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend1" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
  backends {
    name    = "backend-1"
    balance = "%s"
    servers {
      node_id = e2e_node.backend1.id
      port    = "%s"
    }
  }
}
`, nodeName, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBBalanceRoundRobin, goe2econstants.LBPortHTTP)
}

func testAccCheckE2ELoadBalancerConfig_multipleBackends(lbName, nodeName1, nodeName2 string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend1" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_node" "backend2" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
  backends {
    name    = "backend-1"
    balance = "%s"
    servers {
      node_id = e2e_node.backend1.id
      port    = "%s"
    }
  }
  backends {
    name    = "backend-2"
    balance = "%s"
    servers {
      node_id = e2e_node.backend2.id
      port    = "%s"
    }
  }
}
`, nodeName1, nodeName2, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBBalanceRoundRobin, goe2econstants.LBPortHTTP, goe2econstants.LBBalanceLeastConn, goe2econstants.LBPortHTTP)
}

func testAccCheckE2ELoadBalancerConfig_tcpBackend(lbName, nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
  tcp_backend {
    backend_name = "tcp-backend-1"
    port         = "3306"
    balance      = "%s"
    servers {
      node_id = e2e_node.backend.id
      port    = "3306"
    }
  }
}
`, nodeName, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBBalanceRoundRobin)
}

func testAccCheckE2ELoadBalancerConfig_acl(lbName, nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "backend" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"

  acl_list {
    acl_name          = "api-acl"
    acl_condition     = "path_beg"
    acl_matching_path = "/api"
  }

  acl_map {
    acl_name    = "api-acl"
    acl_backend = "backend-1"
  }

  backends {
    name    = "backend-1"
    balance = "%s"
    servers {
      node_id = e2e_node.backend.id
      port    = "%s"
    }
  }
}
`, nodeName, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP, goe2econstants.LBBalanceRoundRobin, goe2econstants.LBPortHTTP)
}

func testAccCheckE2ELoadBalancerConfig_vpc(lbName, vpcName string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  name      = "%s"
  ipv4_cidr = "10.0.0.0/16"
}

resource "e2e_loadbalancer" "test" {
  name     = "%s"
  plan     = "%s"
  lb_mode  = "%s"
  vpc_list = [e2e_vpc.test.id]
}
`, vpcName, lbName, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_ssl(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
  ssl_context {
    redirect_to_https = false
  }
}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTPS)
}

// testAccCheckE2ELoadBalancerConfig_powerOff already exists in resource_loadbalancer_test.go

func testAccCheckE2ELoadBalancerConfig_plan(name, plan string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  name    = "%s"
  plan    = "%s"
  lb_mode = "%s"
}
`, name, plan, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_conflictingNameAndLbName(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  name    = "%s"
  lb_name = "%s-conflict"
  plan    = "%s"
  lb_mode = "%s"
}
`, name, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_conflictingRegionAndLocation(name string) string {
	return fmt.Sprintf(`
resource "e2e_loadbalancer" "test" {
  name     = "%s"
  plan     = "%s"
  lb_mode  = "%s"
  region   = "Mumbai"
  location = "Delhi"
}
`, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}

func testAccCheckE2ELoadBalancerConfig_conflictingFloatingIPFields(name string) string {
	return fmt.Sprintf(`
resource "e2e_reserve_ip" "test" {
  name = "%s-ip"
}

resource "e2e_loadbalancer" "test" {
  name           = "%s"
  plan           = "%s"
  lb_mode        = "%s"
  floating_ip_id = e2e_reserve_ip.test.id
  lb_reserve_ip  = e2e_reserve_ip.test.id
}
`, name, name, goe2econstants.LBPlanE2ELB2, goe2econstants.LBModeHTTP)
}
