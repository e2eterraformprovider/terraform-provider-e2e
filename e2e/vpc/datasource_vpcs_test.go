package vpc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccE2EVPCDataSource_Basic(t *testing.T) {
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCDataSourceConfig_basic(vpcName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_vpcs.test", "id"),
					resource.TestCheckResourceAttrSet("data.e2e_vpcs.test", "vpcs.#"),
				),
			},
		},
	})
}

func TestAccE2EVPCDataSource_MultipleVPCs(t *testing.T) {
	vpcName1 := fmt.Sprintf("test-vpc-1-%s", acctest.RandString(10))
	vpcName2 := fmt.Sprintf("test-vpc-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCDataSourceConfig_multiple(vpcName1, vpcName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_vpcs.test", "vpcs.#"),
				),
			},
		},
	})
}

// Configuration helpers

func testAccCheckE2EVPCDataSourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  vpc_name   = "%s"
  project_id = "%s"
  location   = "%s"
}

data "e2e_vpcs" "test" {
  project_id = "%s"
  location   = "%s"
  depends_on = [e2e_vpc.test]
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EVPCDataSourceConfig_multiple(name1, name2 string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test1" {
  vpc_name   = "%s"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_vpc" "test2" {
  vpc_name   = "%s"
  project_id = "%s"
  location   = "%s"
}

data "e2e_vpcs" "test" {
  project_id = "%s"
  location   = "%s"
  depends_on = [e2e_vpc.test1, e2e_vpc.test2]
}
`, name1, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		name2, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
