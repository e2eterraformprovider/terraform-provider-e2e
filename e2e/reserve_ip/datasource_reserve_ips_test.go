package reserve_ip_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceE2EReserveIps_List(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceE2EReserveIpsConfig_list(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.#"),
				),
			},
		},
	})
}

func TestAccDataSourceE2EReserveIps_WithReserveIP(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceE2EReserveIpsConfig_withReserveIP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.#"),
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.0.reserve_id"),
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.0.status"),
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.0.bought_at"),
				),
			},
		},
	})
}

func TestAccDataSourceE2EReserveIps_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataSourceE2EReserveIpsConfig_missingRegion(),
				ExpectError: regexp.MustCompile(`The argument "region" is required`),
			},
			{
				Config:      testAccCheckDataSourceE2EReserveIpsConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
		},
	})
}

// Configuration helpers for data source tests

func testAccCheckDataSourceE2EReserveIpsConfig_list() string {
	return fmt.Sprintf(`
data "e2e_reserve_ips" "test" {
  project_id = "%s"
  region     = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckDataSourceE2EReserveIpsConfig_withReserveIP() string {
	return fmt.Sprintf(`
resource "e2e_reserve_ip" "test" {
  project_id = "%s"
  location   = "%s"
}

data "e2e_reserve_ips" "test" {
  project_id = "%s"
  region     = "%s"
  depends_on = [e2e_reserve_ip.test]
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckDataSourceE2EReserveIpsConfig_missingRegion() string {
	return fmt.Sprintf(`
data "e2e_reserve_ips" "test" {
  project_id = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccCheckDataSourceE2EReserveIpsConfig_missingProjectID() string {
	return fmt.Sprintf(`
data "e2e_reserve_ips" "test" {
  region = "%s"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}
