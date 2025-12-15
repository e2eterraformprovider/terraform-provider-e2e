package vpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckE2EVPCExists(n string, vpc *goe2e.Vpc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no VPC ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client := cfg.Goe2eClient()

		found, _, err := client.Vpcs.GetVPC(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if found == nil {
			return fmt.Errorf("VPC not found")
		}

		*vpc = *found
		return nil
	}
}

func testAccCheckE2EVPCResourceDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client := cfg.Goe2eClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_vpc" {
			continue
		}

		_, _, err := client.Vpcs.GetVPC(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func TestAccE2EVPC_Basic(t *testing.T) {
	var vpc goe2e.Vpc
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCConfig_basic(vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVPCExists("e2e_vpc.test", &vpc),
					resource.TestCheckResourceAttr("e2e_vpc.test", tfconstants.AttrName, vpcName),
					resource.TestCheckResourceAttr("e2e_vpc.test", tfconstants.AttrIsE2EVPC, "true"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", tfconstants.AttrNetworkID),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", tfconstants.AttrCreatedAt),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", tfconstants.AttrIPv4CIDR),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", tfconstants.AttrGatewayIP),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", tfconstants.AttrPoolSize),
				),
			},
		},
	})
}

func TestAccE2EVPC_WithCustomCIDR(t *testing.T) {
	var vpc goe2e.Vpc
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))
	customCIDR := "10.0.0.0/24"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCConfig_withCustomCIDR(vpcName, customCIDR),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVPCExists("e2e_vpc.test", &vpc),
					resource.TestCheckResourceAttr("e2e_vpc.test", tfconstants.AttrName, vpcName),
					resource.TestCheckResourceAttr("e2e_vpc.test", tfconstants.AttrIPv4, customCIDR),
					resource.TestCheckResourceAttr("e2e_vpc.test", tfconstants.AttrIsE2EVPC, "false"),
				),
			},
		},
	})
}

func TestAccE2EVPC_Import(t *testing.T) {
	var vpc goe2e.Vpc
	vpcName := fmt.Sprintf("test-vpc-import-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCConfig_basic(vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVPCExists("e2e_vpc.test", &vpc),
				),
			},
			{
				ResourceName:            "e2e_vpc.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{
					// These fields may not be importable or may differ
				},
			},
		},
	})
}

// Configuration helpers

func testAccCheckE2EVPCConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  %s = "%s"
}
`, tfconstants.AttrName, name)
}

func testAccCheckE2EVPCConfig_withCustomCIDR(name, cidr string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  %s = "%s"
  %s = "%s"
  %s = false
}
`, tfconstants.AttrName, name, tfconstants.AttrIPv4, cidr, tfconstants.AttrIsE2EVPC)
}
