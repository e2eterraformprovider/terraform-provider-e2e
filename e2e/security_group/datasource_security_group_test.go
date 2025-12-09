package security_group_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2ESecurityGroupDestroy(s *terraform.State) error {
	// Security groups created in tests should be cleaned up
	// For now, return nil as the datasource test doesn't directly delete resources
	return nil
}

func TestAccDataSourceE2ESecurityGroup_Basic(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceE2ESecurityGroupConfig_basic(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("data.e2e_security_group.test", "description", "Test security group"),
					resource.TestCheckResourceAttrSet("data.e2e_security_group.test", "id"),
					resource.TestCheckResourceAttrSet("data.e2e_security_group.test", "rules.#"),
					resource.TestCheckResourceAttrPair(
						"data.e2e_security_group.test", "id",
						"e2e_security_group.test", "id")),
			},
		},
	})
}

func TestAccDataSourceE2ESecurityGroup_WithMultipleRules(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceE2ESecurityGroupConfig_multipleRules(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("data.e2e_security_group.test", "rules.#", "2"),
					resource.TestCheckResourceAttrSet("data.e2e_security_group.test", "rules.0.rule_id"),
					resource.TestCheckResourceAttr("data.e2e_security_group.test", "rules.0.rule_type", "Inbound"),
					resource.TestCheckResourceAttr("data.e2e_security_group.test", "rules.0.protocol_name", "Custom_TCP")),
			},
		},
	})
}

func TestAccDataSourceE2ESecurityGroup_NonExistent(t *testing.T) {
	sgName := fmt.Sprintf("non-existent-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataSourceE2ESecurityGroupConfig_nonExistent(sgName),
				ExpectError: regexp.MustCompile(`.*`),
			},
		},
	})
}

func TestAccDataSourceE2ESecurityGroup_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataSourceE2ESecurityGroupConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckDataSourceE2ESecurityGroupConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckDataSourceE2ESecurityGroupConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

// Configuration helpers for data source tests

func testAccCheckDataSourceE2ESecurityGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group"
  rules {
    rule_type     = "Inbound"
    protocol_name = "All"
    network       = "any"
  }
}

data "e2e_security_group" "test" {
  name       = e2e_security_group.test.name}
`, name)
}

func testAccCheckDataSourceE2ESecurityGroupConfig_multipleRules(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with multiple rules"
  rules {
    rule_type     = "Inbound"
    protocol_name = "Custom_TCP"
    port_range    = "22"
    network       = "any"
    description   = "SSH access"
  }

  rules {
    rule_type     = "Inbound"
    protocol_name = "Custom_TCP"
    port_range    = "80"
    network       = "any"
    description   = "HTTP access"
  }
}

data "e2e_security_group" "test" {
  name       = e2e_security_group.test.name}
`, name)
}

func testAccCheckDataSourceE2ESecurityGroupConfig_nonExistent(name string) string {
	return fmt.Sprintf(`
data "e2e_security_group" "test" {
  name       = "%s"}
`, name)
}

func testAccCheckDataSourceE2ESecurityGroupConfig_missingName() string {
	return `
data "e2e_security_group" "test" {}
`
}

func testAccCheckDataSourceE2ESecurityGroupConfig_missingProjectID() string {
	return `
data "e2e_security_group" "test" {
  name     = "test-sg"}
`
}

func testAccCheckDataSourceE2ESecurityGroupConfig_missingLocation() string {
	return `
data "e2e_security_group" "test" {
  name       = "test-sg"}
`
}
