package node_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2ENode_Basic(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", goe2econstants.NodePowerStatusOn),
					resource.TestCheckResourceAttr("e2e_node.test", "lock_node", "false"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "memory"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "disk"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "price"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "vm_id")),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName)),
			},
			{
				Config: testAccCheckE2ENodeConfig_updated(nodeNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeNameUpdated),
					resource.TestCheckResourceAttr("e2e_node.test", "label", "updated-label")),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_withSSHKeys(nodeName, sshKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_keys.#", "1")),
			},
		},
	})
}

func TestAccE2ENode_WithSecurityGroups(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_withSecurityGroups(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttrSet("e2e_node.test", "default_sg")),
			},
		},
	})
}

func TestAccE2ENode_PowerOperations(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", goe2econstants.NodePowerStatusOn)),
			},
			{
				Config: testAccCheckE2ENodeConfig_powerOff(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", goe2econstants.NodePowerStatusOff)),
			},
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", goe2econstants.NodePowerStatusOn)),
			},
		},
	})
}

func TestAccE2ENode_LockNode(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "lock_node", "false")),
			},
			{
				Config: testAccCheckE2ENodeConfig_locked(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "lock_node", "true")),
			},
		},
	})
}

func TestAccE2ENode_WithStartScript(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_withStartScript(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttrSet("e2e_node.test", "start_script")),
			},
		},
	})
}

func TestAccE2ENode_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
		ProviderFactories: acceptance.TestAccProviderFactories,
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID)),
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

// TestAccE2ENode_ForceNewImmutableFields verifies that changing immutable fields triggers recreation
func TestAccE2ENode_ForceNewImmutableFields(t *testing.T) {
	var nodeID1, nodeID2 string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_forceNew_initial(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID1),
					resource.TestCheckResourceAttr("e2e_node.test", "plan", "c2-2c-4gb"),
					resource.TestCheckResourceAttr("e2e_node.test", "image", "ubuntu-20.04")),
			},
			{
				Config: testAccCheckE2ENodeConfig_forceNew_planChange(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID2),
					resource.TestCheckResourceAttr("e2e_node.test", "plan", "c2-4c-8gb"),
					testAccCheckE2ENodeRecreated(&nodeID1, &nodeID2)),
			},
		},
	})
}

// TestAccE2ENode_ForceNewImage verifies that changing image triggers recreation
func TestAccE2ENode_ForceNewImage(t *testing.T) {
	var nodeID1, nodeID2 string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_forceNew_initial(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID1),
					resource.TestCheckResourceAttr("e2e_node.test", "image", "ubuntu-20.04")),
			},
			{
				Config: testAccCheckE2ENodeConfig_forceNew_imageChange(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID2),
					resource.TestCheckResourceAttr("e2e_node.test", "image", "ubuntu-22.04"),
					testAccCheckE2ENodeRecreated(&nodeID1, &nodeID2)),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
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

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)

		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
		if err != nil {
			return fmt.Errorf("Error creating goe2e client: %s", err)
		}

		ctx := context.Background()
		node, _, err := goe2eClient.Nodes.GetNode(ctx, rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return fmt.Errorf("Node not found")
			}
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
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_node" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
		if err != nil {
			return fmt.Errorf("Error creating goe2e client: %s", err)
		}

		_, _, err = goe2eClient.Nodes.GetNode(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Node still exists: %s", rs.Primary.ID)
		}
		if !strings.Contains(err.Error(), "not found") {
			return err
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
		location := acceptance.GetRegionOrLocationFromState(rs)
		nodeID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, nodeID), nil
	}
}

func testAccCheckE2ENodeRecreated(oldID, newID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *oldID == *newID {
			return fmt.Errorf("Expected node to be recreated, but IDs are the same: %s", *oldID)
		}
		return nil
	}
}

// Configuration helpers

func testAccCheckE2ENodeConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}
`, name)
}

func testAccCheckE2ENodeConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
  label = "updated-label"
}
`, name)
}

func testAccCheckE2ENodeConfig_withSSHKeys(name, sshKeyLabel string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  public_key = "%s"
}

resource "e2e_node" "test" {
  name     = "%s"
  plan     = "c2-2c-4gb"
  image    = "ubuntu-20.04"
  ssh_keys = [e2e_ssh_key.test.label]
}
`, sshKeyLabel, publicKey, name)
}

func testAccCheckE2ENodeConfig_withSecurityGroups(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}
`, name)
}

func testAccCheckE2ENodeConfig_powerOff(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name         = "%s"
  plan         = "c2-2c-4gb"
  image        = "ubuntu-20.04"
  power_status = "%s"
}
`, name, goe2econstants.NodePowerStatusOff)
}

func testAccCheckE2ENodeConfig_locked(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name      = "%s"
  plan      = "c2-2c-4gb"
  image     = "ubuntu-20.04"
  lock_node = true
}
`, name)
}

func testAccCheckE2ENodeConfig_withStartScript(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name         = "%s"
  plan         = "c2-2c-4gb"
  image        = "ubuntu-20.04"
  start_script = "#!/bin/bash\necho 'Hello World'"
}
`, name)
}

// Error case configurations

func testAccCheckE2ENodeConfig_missingName() string {
	return `
resource "e2e_node" "test" {
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}
`
}

func testAccCheckE2ENodeConfig_missingPlan() string {
	return `
resource "e2e_node" "test" {
  name  = "test-node"
  image = "ubuntu-20.04"
}
`
}

func testAccCheckE2ENodeConfig_missingImage() string {
	return `
resource "e2e_node" "test" {
  name = "test-node"
  plan = "c2-2c-4gb"
}
`
}

func testAccCheckE2ENodeConfig_missingProjectID() string {
	return `
resource "e2e_node" "test" {
  name     = "test-node"
  plan     = "c2-2c-4gb"
  image    = "ubuntu-20.04"}
`
}

func testAccCheckE2ENodeConfig_missingLocation() string {
	return `
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}
`
}

func testAccCheckE2ENodeConfig_invalidName() string {
	return `
resource "e2e_node" "test" {
  name       = "invalid name with spaces"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}
`
}

// ForceNew test configurations

func testAccCheckE2ENodeConfig_forceNew_initial(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}
`, name)
}

func testAccCheckE2ENodeConfig_forceNew_planChange(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-4c-8gb"
  image      = "ubuntu-20.04"}
`, name)
}

func testAccCheckE2ENodeConfig_forceNew_imageChange(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-22.04"}
`, name)
}

// V3 Feature Tests

func TestAccE2ENode_V3_SSHKeyIDs(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-v3-sshkeyids-%s", acctest.RandString(10))
	sshKeyName := fmt.Sprintf("test-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_v3_sshKeyIDs(nodeName, sshKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_key_ids.#", "1"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "ssh_key_ids.0"),
				),
			},
		},
	})
}

func TestAccE2ENode_V3_SSHKeyIDs_DataSource(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-v3-sshkey-ds-%s", acctest.RandString(10))
	sshKeyName := fmt.Sprintf("test-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_v3_sshKeyIDs_dataSource(nodeName, sshKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_key_ids.#", "1"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "ssh_key_ids.0"),
					// Verify the data source resolved correctly
					resource.TestCheckResourceAttrSet("data.e2e_ssh_key.existing", "id"),
					resource.TestCheckResourceAttr("data.e2e_ssh_key.existing", "label", sshKeyName),
				),
			},
		},
	})
}

func TestAccE2ENode_V3_RootDisk(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-v3-rootdisk-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_v3_rootDisk(nodeName, 100),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "root_disk.0.size_gb", "100"),
				),
			},
		},
	})
}

func TestAccE2ENode_V3_ReserveIPID(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-v3-reserveipid-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_v3_reserveIPID(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttrSet("e2e_node.test", "reserve_ip_id"),
				),
			},
		},
	})
}

func TestAccE2ENode_V3_NetworkInterface(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-v3-netif-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_v3_networkInterface(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "network_interface.#", "1"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "network_interface.0.vpc_id"),
					resource.TestCheckResourceAttr("e2e_node.test", "network_interface.0.assign_public_ip", "true"),
				),
			},
		},
	})
}

func TestAccE2ENode_V3_TagIDs(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-v3-tagids-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_v3_tagIDs(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccE2ENode_V3_IPv6Address(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-v3-ipv6-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_v3_ipv6(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "is_ipv6_availed", "true"),
					// IPv6 address should be set when IPv6 is enabled
					resource.TestCheckResourceAttrSet("e2e_node.test", "ipv6_address"),
				),
			},
		},
	})
}

func TestAccE2ENode_V3_DeprecationWarnings(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-v3-deprecated-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_v3_deprecatedFields(nodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
				),
				ExpectNonEmptyPlan: false, // Deprecated fields should still work
			},
		},
	})
}

// V3 Test Configuration Functions

func testAccCheckE2ENodeConfig_v3_sshKeyIDs(nodeName, sshKeyName string) string {
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCm4X3ck1X+MfL9FhvV4tGqqmJz3NZ2d7hP2gDqe1pQqE9yx0p4pWOQFLNQg4DZxBm8NtP5KzN9qdGDhPZx7Wd1JNLiPqKYp7zVnLpfN4fwDQnWwN7F0JxP4mX8c9K7T6Q+Nw4cPz4vL0xH test@example.com"
}

resource "e2e_node" "test" {
  name        = "%s"
  plan        = "C3.8GB"
  image       = "Ubuntu-20.04"
  ssh_key_ids = [e2e_ssh_key.test.id]
}
`, sshKeyName, nodeName)
}

func testAccCheckE2ENodeConfig_v3_sshKeyIDs_dataSource(nodeName, sshKeyName string) string {
	return fmt.Sprintf(`
# Create an SSH key (simulates an existing key in the account)
resource "e2e_ssh_key" "existing_key" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCm4X3ck1X+MfL9FhvV4tGqqmJz3NZ2d7hP2gDqe1pQqE9yx0p4pWOQFLNQg4DZxBm8NtP5KzN9qdGDhPZx7Wd1JNLiPqKYp7zVnLpfN4fwDQnWwN7F0JxP4mX8c9K7T6Q+Nw4cPz4vL0xH test@example.com"
}

# Look up the existing SSH key by label (name)
data "e2e_ssh_key" "existing" {
  label = e2e_ssh_key.existing_key.name
}

# Create node referencing the SSH key via data source
resource "e2e_node" "test" {
  name        = "%s"
  plan        = "C3.8GB"
  image       = "Ubuntu-20.04"
  ssh_key_ids = [data.e2e_ssh_key.existing.id]
}
`, sshKeyName, nodeName)
}

func testAccCheckE2ENodeConfig_v3_rootDisk(nodeName string, diskSize int) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "C3.8GB"
  image = "Ubuntu-20.04"

  root_disk {
    size_gb   = %d
    disk_type = "standard"
  }
}
`, nodeName, diskSize)
}

func testAccCheckE2ENodeConfig_v3_reserveIPID(nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_reserved_ip" "test" {
  name = "test-ip-%s"
}

resource "e2e_node" "test" {
  name          = "%s"
  plan          = "C3.8GB"
  image         = "Ubuntu-20.04"
  reserve_ip_id = e2e_reserved_ip.test.id
}
`, acctest.RandString(5), nodeName)
}

func testAccCheckE2ENodeConfig_v3_networkInterface(nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  name = "test-vpc-%s"
  cidr = "10.0.0.0/16"
}

resource "e2e_node" "test" {
  name  = "%s"
  plan  = "C3.8GB"
  image = "Ubuntu-20.04"

  network_interface {
    vpc_id           = e2e_vpc.test.id
    assign_public_ip = true
    enable_ipv6      = false
  }
}
`, acctest.RandString(5), nodeName)
}

func testAccCheckE2ENodeConfig_v3_tagIDs(nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name    = "%s"
  plan    = "C3.8GB"
  image   = "Ubuntu-20.04"
  tag_ids = [1, 2]  # Example tag IDs (using label API)
}
`, nodeName)
}

func testAccCheckE2ENodeConfig_v3_ipv6(nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name            = "%s"
  plan            = "C3.8GB"
  image           = "Ubuntu-20.04"
  is_ipv6_availed = true
}
`, nodeName)
}

func testAccCheckE2ENodeConfig_v3_deprecatedFields(nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "C3.8GB"
  image      = "Ubuntu-20.04"
  ssh_keys   = ["default-key"]  # Deprecated field
  reserve_ip = ""                # Deprecated field
}
`, nodeName)
}
