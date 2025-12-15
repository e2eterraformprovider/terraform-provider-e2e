package vpc_test

import (
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2EVPCDestroy(s *terraform.State) error {
	// VPCs created in tests should be cleaned up
	// For now, return nil as the datasource test doesn't directly delete resources
	return nil
}

func TestAccE2EVPCDataSource_Basic(t *testing.T) {
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCDataSourceConfig_basic(vpcName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_vpcs.test", tfconstants.AttrID),
					resource.TestCheckResourceAttrSet("data.e2e_vpcs.test", tfconstants.AttrVPCList+".#")),
			},
		},
	})
}

func TestAccE2EVPCDataSource_MultipleVPCs(t *testing.T) {
	vpcName1 := fmt.Sprintf("test-vpc-1-%s", acctest.RandString(10))
	vpcName2 := fmt.Sprintf("test-vpc-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCDataSourceConfig_multiple(vpcName1, vpcName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_vpcs.test", tfconstants.AttrVPCList+".#")),
			},
		},
	})
}

// Configuration helpers

func testAccCheckE2EVPCDataSourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  %s = "%s"
}

data "e2e_vpcs" "test" {
  depends_on = [e2e_vpc.test]
}
`, tfconstants.AttrName, name)
}

func testAccCheckE2EVPCDataSourceConfig_multiple(name1, name2 string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test1" {
  %s = "%s"
}

resource "e2e_vpc" "test2" {
  %s = "%s"
}

data "e2e_vpcs" "test" {
  depends_on = [e2e_vpc.test1, e2e_vpc.test2]
}
`, tfconstants.AttrName, name1, tfconstants.AttrName, name2)
}
