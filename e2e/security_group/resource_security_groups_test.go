package security_group_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2ESecurityGroup_Basic(t *testing.T) {
	var sgID string
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESecurityGroupConfig_basic(sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESecurityGroupExists("e2e_security_group.test", &sgID),
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "description", "Test security group"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "default", "false"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
		},
	})
}

func TestAccE2ESecurityGroup_WithMultipleRules(t *testing.T) {
	var sgID string
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESecurityGroupConfig_multipleRules(sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESecurityGroupExists("e2e_security_group.test", &sgID),
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "3"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.rule_type", "Inbound"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.protocol_name", "Custom_TCP"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.port_range", "22"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.rule_type", "Inbound"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.protocol_name", "Custom_TCP"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.port_range", "80"),
				),
			},
		},
	})
}

func TestAccE2ESecurityGroup_Update(t *testing.T) {
	var sgID string
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESecurityGroupConfig_basic(sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESecurityGroupExists("e2e_security_group.test", &sgID),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "1"),
				),
			},
			{
				Config: testAccCheckE2ESecurityGroupConfig_updated(sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESecurityGroupExists("e2e_security_group.test", &sgID),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "2"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "description", "Updated security group"),
				),
			},
		},
	})
}

func TestAccE2ESecurityGroup_WithDifferentNetworks(t *testing.T) {
	var sgID string
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESecurityGroupConfig_differentNetworks(sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESecurityGroupExists("e2e_security_group.test", &sgID),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "3"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.network", "any"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.network", "manual"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.network_cidr", "10.0.0.0/24"),
				),
			},
		},
	})
}

func TestAccE2ESecurityGroup_MakeDefault(t *testing.T) {
	var sgID string
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESecurityGroupConfig_basic(sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESecurityGroupExists("e2e_security_group.test", &sgID),
					resource.TestCheckResourceAttr("e2e_security_group.test", "default", "false"),
				),
			},
			{
				Config: testAccCheckE2ESecurityGroupConfig_makeDefault(sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESecurityGroupExists("e2e_security_group.test", &sgID),
					resource.TestCheckResourceAttr("e2e_security_group.test", "default", "true"),
				),
			},
		},
	})
}

func TestAccE2ESecurityGroup_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESecurityGroupConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2ESecurityGroupConfig_missingRules(),
				ExpectError: regexp.MustCompile(`The argument "rules" is required`),
			},
			{
				Config:      testAccCheckE2ESecurityGroupConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2ESecurityGroupConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccE2ESecurityGroup_InvalidRuleType(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESecurityGroupConfig_invalidRuleType(sgName),
				ExpectError: regexp.MustCompile(`expected rule_type to be one of \[Inbound Outbound\]`),
			},
		},
	})
}

func TestAccE2ESecurityGroup_InvalidProtocol(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESecurityGroupConfig_invalidProtocol(sgName),
				ExpectError: regexp.MustCompile(`expected protocol_name to be one of`),
			},
		},
	})
}

func TestAccE2ESecurityGroup_InvalidNetwork(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESecurityGroupConfig_invalidNetwork(sgName),
				ExpectError: regexp.MustCompile(`expected network to be one of \[myNetwork manual any\]`),
			},
		},
	})
}

func TestAccE2ESecurityGroup_Import(t *testing.T) {
	var sgID string
	sgName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESecurityGroupConfig_basic(sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESecurityGroupExists("e2e_security_group.test", &sgID),
				),
			},
			{
				ResourceName:      "e2e_security_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Helper functions

var testAccProvider *schema.Provider

func init() {
	testAccProvider = e2e.Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SERVICE_API_KEY"); v == "" {
		t.Fatal("SERVICE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICE_AUTH_TOKEN"); v == "" {
		t.Fatal("SERVICE_AUTH_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_PROJECT_ID"); v == "" {
		t.Fatal("E2E_TEST_PROJECT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_LOCATION"); v == "" {
		t.Fatal("E2E_TEST_LOCATION must be set for acceptance tests")
	}
}

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"e2e": func() (*schema.Provider, error) {
		return e2e.Provider(), nil
	},
}

func testAccCheckE2ESecurityGroupExists(resourceName string, sgID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		name := rs.Primary.Attributes["name"]
		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		sg, err := client.GetSecurityGroup(name, projectID, location)
		if err != nil {
			return err
		}

		if sg == nil {
			return fmt.Errorf("Security Group not found")
		}

		*sgID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2ESecurityGroupDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_security_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		_, err := client.GetSecurityGroup(name, projectID, location)
		if err == nil {
			return fmt.Errorf("Security Group still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2ESecurityGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group"
  project_id  = "%s"
  location    = "%s"

  rules {
    rule_type     = "Inbound"
    protocol_name = "All"
    network       = "any"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_multipleRules(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with multiple rules"
  project_id  = "%s"
  location    = "%s"

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

  rules {
    rule_type     = "Outbound"
    protocol_name = "All"
    network       = "any"
    description   = "Allow all outbound"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Updated security group"
  project_id  = "%s"
  location    = "%s"

  rules {
    rule_type     = "Inbound"
    protocol_name = "Custom_TCP"
    port_range    = "22"
    network       = "any"
  }

  rules {
    rule_type     = "Inbound"
    protocol_name = "Custom_TCP"
    port_range    = "443"
    network       = "any"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_differentNetworks(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with different networks"
  project_id  = "%s"
  location    = "%s"

  rules {
    rule_type     = "Inbound"
    protocol_name = "All"
    network       = "any"
  }

  rules {
    rule_type     = "Inbound"
    protocol_name = "Custom_TCP"
    port_range    = "22"
    network       = "manual"
    network_cidr  = "10.0.0.0/24"
    size          = 256
  }

  rules {
    rule_type     = "Inbound"
    protocol_name = "Custom_TCP"
    port_range    = "3306"
    network       = "myNetwork"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_makeDefault(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group"
  project_id  = "%s"
  location    = "%s"
  default     = true

  rules {
    rule_type     = "Inbound"
    protocol_name = "All"
    network       = "any"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

// Error case configurations

func testAccCheckE2ESecurityGroupConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  description = "Test security group"
  project_id  = "%s"
  location    = "%s"

  rules {
    rule_type     = "Inbound"
    protocol_name = "All"
    network       = "any"
  }
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_missingRules() string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "test-sg"
  description = "Test security group"
  project_id  = "%s"
  location    = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "test-sg"
  description = "Test security group"
  location    = "%s"

  rules {
    rule_type     = "Inbound"
    protocol_name = "All"
    network       = "any"
  }
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "test-sg"
  description = "Test security group"
  project_id  = "%s"

  rules {
    rule_type     = "Inbound"
    protocol_name = "All"
    network       = "any"
  }
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccCheckE2ESecurityGroupConfig_invalidRuleType(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group"
  project_id  = "%s"
  location    = "%s"

  rules {
    rule_type     = "Invalid"
    protocol_name = "All"
    network       = "any"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_invalidProtocol(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group"
  project_id  = "%s"
  location    = "%s"

  rules {
    rule_type     = "Inbound"
    protocol_name = "InvalidProtocol"
    network       = "any"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESecurityGroupConfig_invalidNetwork(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group"
  project_id  = "%s"
  location    = "%s"

  rules {
    rule_type     = "Inbound"
    protocol_name = "All"
    network       = "invalidNetwork"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
