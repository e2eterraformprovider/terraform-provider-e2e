package node_test

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

func TestAccE2ENode_Basic(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "plan", "c2-2c-4gb"),
					resource.TestCheckResourceAttr("e2e_node.test", "image", "ubuntu-20.04"),
					resource.TestCheckResourceAttr("e2e_node.test", "label", "default"),
					resource.TestCheckResourceAttr("e2e_node.test", "backup", "false"),
					resource.TestCheckResourceAttr("e2e_node.test", "default_public_ip", "false"),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", "power_on"),
					resource.TestCheckResourceAttr("e2e_node.test", "lock_node", "false"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "memory"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "disk"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "price"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "vm_id"),
				),
			},
		},
	})
}

func TestAccE2ENode_Update(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	nodeNameUpdated := fmt.Sprintf("test-node-updated-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
				),
			},
			{
				Config: testAccCheckE2ENodeConfig_updated(nodeNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeNameUpdated),
					resource.TestCheckResourceAttr("e2e_node.test", "label", "updated-label"),
				),
			},
		},
	})
}

func TestAccE2ENode_WithSSHKeys(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	sshKeyLabel := fmt.Sprintf("test-ssh-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_withSSHKeys(nodeName, sshKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_keys.#", "1"),
				),
			},
		},
	})
}

func TestAccE2ENode_WithSecurityGroups(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_withSecurityGroups(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttrSet("e2e_node.test", "default_sg"),
				),
			},
		},
	})
}

func TestAccE2ENode_PowerOperations(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", "power_on"),
				),
			},
			{
				Config: testAccCheckE2ENodeConfig_powerOff(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", "power_off"),
				),
			},
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", "power_on"),
				),
			},
		},
	})
}

func TestAccE2ENode_LockNode(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "lock_node", "false"),
				),
			},
			{
				Config: testAccCheckE2ENodeConfig_locked(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "lock_node", "true"),
				),
			},
		},
	})
}

func TestAccE2ENode_WithStartScript(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_withStartScript(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttrSet("e2e_node.test", "start_script"),
				),
			},
		},
	})
}

func TestAccE2ENode_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ENodeConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2ENodeConfig_missingPlan(),
				ExpectError: regexp.MustCompile(`The argument "plan" is required`),
			},
			{
				Config:      testAccCheckE2ENodeConfig_missingImage(),
				ExpectError: regexp.MustCompile(`The argument "image" is required`),
			},
			{
				Config:      testAccCheckE2ENodeConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2ENodeConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccE2ENode_InvalidName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ENodeConfig_invalidName(),
				ExpectError: regexp.MustCompile(`the name field cannot be blank, must not contain whitespace or special characters`),
			},
		},
	})
}

func TestAccE2ENode_Import(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
				),
			},
			{
				ResourceName:            "e2e_node.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccE2ENodeImportID("e2e_node.test"),
				ImportStateVerifyIgnore: []string{"start_script", "reboot_node", "reinstall_node"},
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

func testAccCheckE2ENodeExists(resourceName string, nodeID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		node, err := client.GetNode(rs.Primary.ID, projectID, location)
		if err != nil {
			return err
		}

		if node == nil {
			return fmt.Errorf("Node not found")
		}

		*nodeID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2ENodeDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_node" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		_, err := client.GetNode(rs.Primary.ID, projectID, location)
		if err == nil {
			return fmt.Errorf("Node still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccE2ENodeImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]
		nodeID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, nodeID), nil
	}
}

// Configuration helpers

func testAccCheckE2ENodeConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  label      = "updated-label"
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_withSSHKeys(name, sshKeyLabel string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  public_key = "%s"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
  ssh_keys   = [e2e_ssh_key.test.label]
}
`, sshKeyLabel, publicKey, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_withSecurityGroups(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_powerOff(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name         = "%s"
  plan         = "c2-2c-4gb"
  image        = "ubuntu-20.04"
  power_status = "power_off"
  project_id   = "%s"
  location     = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_locked(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  lock_node  = true
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_withStartScript(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name         = "%s"
  plan         = "c2-2c-4gb"
  image        = "ubuntu-20.04"
  start_script = "#!/bin/bash\necho 'Hello World'"
  project_id   = "%s"
  location     = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

// Error case configurations

func testAccCheckE2ENodeConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_missingPlan() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_missingImage() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name     = "test-node"
  plan     = "c2-2c-4gb"
  image    = "ubuntu-20.04"
  location = "%s"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccCheckE2ENodeConfig_invalidName() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "invalid name with spaces"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
