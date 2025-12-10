package reserve_ip_test

import (
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceE2EReserveIps_List(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceE2EReserveIpsConfig_list(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.#")),
			},
		},
	})
}

func TestAccDataSourceE2EReserveIps_WithReserveIP(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceE2EReserveIpsConfig_withReserveIP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.#"),
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.0.reserve_id"),
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.0.ip_address"),
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.0.status"),
					resource.TestCheckResourceAttrSet("data.e2e_reserve_ips.test", "reserve_ips_list.0.bought_at")),
			},
		},
	})
}

func TestAccDataSourceE2EReserveIps_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
	return `
data "e2e_reserve_ips" "test" {}
`
}

func testAccCheckDataSourceE2EReserveIpsConfig_withReserveIP() string {
	return `
resource "e2e_reserve_ip" "test" {}

data "e2e_reserve_ips" "test" {  depends_on = [e2e_reserve_ip.test]
}
`
}

func testAccCheckDataSourceE2EReserveIpsConfig_missingRegion() string {
	return `
data "e2e_reserve_ips" "test" {}
`
}

func testAccCheckDataSourceE2EReserveIpsConfig_missingProjectID() string {
	return `
data "e2e_reserve_ips" "test" {}
`
}
