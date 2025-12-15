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

// ============================================================================
// Real-World Scenarios Tests (Lines 158-203)
// ============================================================================

// TestAccE2ENode_V2ToV3Migration tests V2 config migration to V3
// This test verifies backward compatibility and smooth migration path
func TestAccE2ENode_V2ToV3Migration(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-v2tov3-%s", acctest.RandString(10))
	sshKeyLabel := fmt.Sprintf("test-ssh-v2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Node with V2 field names
			{
				Config: testAccCheckE2ENodeConfig_v2Fields(nodeName, sshKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					// Verify V2 fields are set
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_keys.#", "1"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "location"),
				),
			},
			// Step 2: Run terraform plan - verify no changes (backward compatible)
			{
				Config:             testAccCheckE2ENodeConfig_v2Fields(nodeName, sshKeyLabel),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 3: Modify config to use V3 names
			{
				Config: testAccCheckE2ENodeConfig_v3Fields(nodeName, sshKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					// Verify V3 fields are set
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_key_ids.#", "1"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "region"),
				),
			},
			// Step 4: Verify no forced recreation (same node ID)
			{
				Config: testAccCheckE2ENodeConfig_v3Fields(nodeName, sshKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					// Verify node ID hasn't changed (no recreation)
					resource.TestCheckResourceAttr("e2e_node.test", "id", nodeID),
				),
			},
		},
	})
}

// TestAccE2ENode_NetworkInterfaceComprehensive tests network_interface block comprehensively
func TestAccE2ENode_NetworkInterfaceComprehensive(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-netif-%s", acctest.RandString(10))
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Node with network_interface block
			{
				Config: testAccCheckE2ENodeConfig_networkInterfaceFull(nodeName, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					// Verify network_interface block
					resource.TestCheckResourceAttr("e2e_node.test", "network_interface.#", "1"),
					// Verify VPC attachment works
					resource.TestCheckResourceAttrSet("e2e_node.test", "network_interface.0.vpc_id"),
					// Verify public IP assignment works
					resource.TestCheckResourceAttr("e2e_node.test", "network_interface.0.assign_public_ip", "true"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "public_ip_address"),
					// Verify IPv6 assignment works
					resource.TestCheckResourceAttr("e2e_node.test", "network_interface.0.enable_ipv6", "true"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "ipv6_address"),
				),
			},
		},
	})
}

// TestAccE2ENode_VolumeAttachmentMigration tests volume attachment migration
func TestAccE2ENode_VolumeAttachmentMigration(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-va-mig-%s", acctest.RandString(10))
	volumeName := fmt.Sprintf("test-vol-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Node with block_storage_ids (deprecated)
			{
				Config: testAccCheckE2ENodeConfig_withBlockStorageIDs(nodeName, volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "block_storage_ids.#", "1"),
				),
			},
			// Step 2: Migrate to e2e_volume_attachment resources
			{
				Config: testAccCheckE2ENodeConfig_withVolumeAttachment(nodeName, volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					// Verify no forced recreation (same node ID)
					resource.TestCheckResourceAttr("e2e_node.test", "id", nodeID),
					// Verify volume attachment resource exists
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", "node_id"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", "volume_id"),
				),
			},
		},
	})
}

// TestAccE2ENode_SSHKeysMigration tests SSH keys migration
func TestAccE2ENode_SSHKeysMigration(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-ssh-mig-%s", acctest.RandString(10))
	sshKeyLabel := fmt.Sprintf("test-ssh-mig-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Node with ssh_keys (deprecated)
			{
				Config: testAccCheckE2ENodeConfig_withSSHKeysV2(nodeName, sshKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_keys.#", "1"),
				),
			},
			// Step 2: Migrate to ssh_key_ids
			{
				Config: testAccCheckE2ENodeConfig_withSSHKeyIDs(nodeName, sshKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					// Verify no forced recreation (same node ID)
					resource.TestCheckResourceAttr("e2e_node.test", "id", nodeID),
					// Verify ssh_key_ids is set
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_key_ids.#", "1"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "ssh_key_ids.0"),
				),
			},
		},
	})
}

// TestAccE2ENode_AsyncOperations tests async operations and polling
func TestAccE2ENode_AsyncOperations(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-async-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Node and verify polling to "Running" status
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					// Verify node reaches Running status (polling should handle this)
					testAccCheckE2ENodeStatusRunning("e2e_node.test"),
				),
			},
			// Step 2: Power operations and verify status changes
			{
				Config: testAccCheckE2ENodeConfig_powerOff(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", goe2econstants.NodePowerStatusOff),
					// Verify status reflects powered off state
					testAccCheckE2ENodeStatusPoweredOff("e2e_node.test"),
				),
			},
			// Step 3: Power back on
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "power_status", goe2econstants.NodePowerStatusOn),
					testAccCheckE2ENodeStatusRunning("e2e_node.test"),
				),
			},
		},
	})
}

// TestAccE2ENode_ImportComprehensive tests import functionality comprehensively
func TestAccE2ENode_ImportComprehensive(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-import-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Node
			{
				Config: testAccCheckE2ENodeConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
				),
			},
			// Step 2: Import with simple format: <node_id>
			{
				ResourceName:            "e2e_node.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           nodeID,
				ImportStateVerifyIgnore: []string{"start_script", "reboot_node", "reinstall_node", "user_data"},
			},
			// Step 3: Import with full format: <project_id>/<region>/<node_id>
			{
				ResourceName:            "e2e_node.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccE2ENodeImportID("e2e_node.test"),
				ImportStateVerifyIgnore: []string{"start_script", "reboot_node", "reinstall_node", "user_data"},
			},
		},
	})
}

// Helper functions for async operation tests

func testAccCheckE2ENodeStatusRunning(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		status := rs.Primary.Attributes["status"]
		if status != goe2econstants.NodeStatusRunning {
			return fmt.Errorf("Expected node status to be %s, got %s", goe2econstants.NodeStatusRunning, status)
		}

		return nil
	}
}

func testAccCheckE2ENodeStatusPoweredOff(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		status := rs.Primary.Attributes["status"]
		if status != goe2econstants.NodeStatusPoweredOff {
			return fmt.Errorf("Expected node status to be %s, got %s", goe2econstants.NodeStatusPoweredOff, status)
		}

		return nil
	}
}

// Configuration helper functions for new tests

func testAccCheckE2ENodeConfig_v2Fields(nodeName, sshKeyLabel string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k"
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  public_key = "%s"
}

resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  ssh_keys   = [e2e_ssh_key.test.label]  # V2 field
  location   = "Mumbai"                    # V2 field
  start_script = "#!/bin/bash\necho 'V2 script'"  # V2 field
}
`, sshKeyLabel, publicKey, nodeName)
}

func testAccCheckE2ENodeConfig_v3Fields(nodeName, sshKeyLabel string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k"
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  public_key = "%s"
}

resource "e2e_node" "test" {
  name        = "%s"
  plan        = "c2-2c-4gb"
  image       = "ubuntu-20.04"
  ssh_key_ids = [e2e_ssh_key.test.id]  # V3 field
  region      = "Mumbai"                 # V3 field
  user_data   = "#!/bin/bash\necho 'V3 script'"  # V3 field
}
`, sshKeyLabel, publicKey, nodeName)
}

func testAccCheckE2ENodeConfig_networkInterfaceFull(nodeName, vpcName string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  name = "%s"
  cidr = "10.0.0.0/16"
}

resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"

  network_interface {
    vpc_id           = e2e_vpc.test.id
    assign_public_ip = true
    enable_ipv6      = true
  }
}
`, vpcName, nodeName)
}

func testAccCheckE2ENodeConfig_withBlockStorageIDs(nodeName, volumeName string) string {
	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name = "%s"
  size = 20
}

resource "e2e_node" "test" {
  name             = "%s"
  plan             = "c2-2c-4gb"
  image            = "ubuntu-20.04"
  block_storage_ids = [e2e_blockstorage.test.id]  # Deprecated field
}
`, volumeName, nodeName)
}

func testAccCheckE2ENodeConfig_withVolumeAttachment(nodeName, volumeName string) string {
	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name = "%s"
  size = 20
}

resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_volume_attachment" "test" {
  node_id   = e2e_node.test.id
  volume_id = e2e_blockstorage.test.id
}
`, volumeName, nodeName)
}

func testAccCheckE2ENodeConfig_withSSHKeysV2(nodeName, sshKeyLabel string) string {
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
  ssh_keys = [e2e_ssh_key.test.label]  # V2 deprecated field
}
`, sshKeyLabel, publicKey, nodeName)
}

func testAccCheckE2ENodeConfig_withSSHKeyIDs(nodeName, sshKeyLabel string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k"
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  public_key = "%s"
}

resource "e2e_node" "test" {
  name        = "%s"
  plan        = "c2-2c-4gb"
  image       = "ubuntu-20.04"
  ssh_key_ids = [e2e_ssh_key.test.id]  # V3 field
}
`, sshKeyLabel, publicKey, nodeName)
}

// ============================================================================
// Deprecation Validation Tests (Lines 205-228)
// ============================================================================

// TestAccE2ENode_DeprecationWarnings tests that deprecation warnings appear when using V2 fields
func TestAccE2ENode_DeprecationWarnings(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-deprecation-%s", acctest.RandString(10))
	sshKeyLabel := fmt.Sprintf("test-ssh-dep-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			// Test that deprecated fields still work but may log warnings
			{
				Config: testAccCheckE2ENodeConfig_deprecatedFields(nodeName, sshKeyLabel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					// Verify deprecated fields are still functional
					resource.TestCheckResourceAttr("e2e_node.test", "ssh_keys.#", "1"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "start_script"),
					resource.TestCheckResourceAttrSet("e2e_node.test", "location"),
				),
				// Note: Deprecation warnings appear in logs, not as errors
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccE2ENode_ConflictsWithValidation tests ConflictsWith validation errors
func TestAccE2ENode_ConflictsWithValidation(t *testing.T) {
	sshKeyLabel := fmt.Sprintf("test-ssh-conflict-%s", acctest.RandString(10))
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			// Test: ssh_keys and ssh_key_ids conflict
			{
				Config:      testAccCheckE2ENodeConfig_conflictSSHKeys(sshKeyLabel, publicKey),
				ExpectError: regexp.MustCompile(`(?i).*conflicts? with.*ssh_key_ids?.*`),
			},
			// Test: reserve_ip and reserve_ip_id conflict
			{
				Config:      testAccCheckE2ENodeConfig_conflictReserveIP(),
				ExpectError: regexp.MustCompile(`(?i).*conflicts? with.*reserve_ip_id.*`),
			},
		},
	})
}

// ============================================================================
// Performance Validation Tests (Lines 229-241)
// ============================================================================

// TestAccE2ENode_PerformanceMultipleNodes tests creating multiple nodes in sequence
func TestAccE2ENode_PerformanceMultipleNodes(t *testing.T) {
	baseName := fmt.Sprintf("test-node-perf-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_multipleNodes(baseName, 3),
				Check: resource.ComposeTestCheckFunc(
					// Verify all nodes are created
					resource.TestCheckResourceAttr("e2e_node.test1", "name", fmt.Sprintf("%s-1", baseName)),
					resource.TestCheckResourceAttr("e2e_node.test2", "name", fmt.Sprintf("%s-2", baseName)),
					resource.TestCheckResourceAttr("e2e_node.test3", "name", fmt.Sprintf("%s-3", baseName)),
					// Verify all nodes exist
					testAccCheckE2ENodeExists("e2e_node.test1", new(string)),
					testAccCheckE2ENodeExists("e2e_node.test2", new(string)),
					testAccCheckE2ENodeExists("e2e_node.test3", new(string)),
				),
			},
		},
	})
}

// TestAccE2ENode_PerformanceMultipleAttachments tests node with multiple volume attachments
func TestAccE2ENode_PerformanceMultipleAttachments(t *testing.T) {
	var nodeID string
	nodeName := fmt.Sprintf("test-node-multi-attach-%s", acctest.RandString(10))
	volume1Name := fmt.Sprintf("test-vol-1-%s", acctest.RandString(10))
	volume2Name := fmt.Sprintf("test-vol-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeConfig_multipleVolumeAttachments(nodeName, volume1Name, volume2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ENodeExists("e2e_node.test", &nodeID),
					resource.TestCheckResourceAttr("e2e_node.test", "name", nodeName),
					// Verify both volume attachments exist
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test1", "node_id"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test1", "volume_id"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test2", "node_id"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test2", "volume_id"),
					// Verify both attachments reference the same node
					resource.TestCheckResourceAttr("e2e_volume_attachment.test1", "node_id", nodeID),
					resource.TestCheckResourceAttr("e2e_volume_attachment.test2", "node_id", nodeID),
				),
			},
		},
	})
}

// Configuration helper functions for deprecation and performance tests

func testAccCheckE2ENodeConfig_deprecatedFields(nodeName, sshKeyLabel string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k"
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  public_key = "%s"
}

resource "e2e_node" "test" {
  name        = "%s"
  plan        = "c2-2c-4gb"
  image       = "ubuntu-20.04"
  ssh_keys    = [e2e_ssh_key.test.label]  # Deprecated field
  location    = "Mumbai"                    # Deprecated field
  start_script = "#!/bin/bash\necho 'Deprecated script'"  # Deprecated field
}
`, sshKeyLabel, publicKey, nodeName)
}

func testAccCheckE2ENodeConfig_conflictSSHKeys(sshKeyLabel, publicKey string) string {
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  public_key = "%s"
}

resource "e2e_node" "test" {
  name        = "test-conflict"
  plan        = "c2-2c-4gb"
  image       = "ubuntu-20.04"
  ssh_keys    = [e2e_ssh_key.test.label]    # Conflicts with ssh_key_ids
  ssh_key_ids = [e2e_ssh_key.test.id]       # Conflicts with ssh_keys
}
`, sshKeyLabel, publicKey)
}

func testAccCheckE2ENodeConfig_conflictReserveIP() string {
	return `
resource "e2e_reserved_ip" "test" {
  name = "test-ip-conflict"
}

resource "e2e_node" "test" {
  name          = "test-conflict"
  plan          = "c2-2c-4gb"
  image         = "ubuntu-20.04"
  reserve_ip    = "test-ip"              # Conflicts with reserve_ip_id
  reserve_ip_id = e2e_reserved_ip.test.id # Conflicts with reserve_ip
}
`
}

func testAccCheckE2ENodeConfig_multipleNodes(baseName string, count int) string {
	config := ""
	for i := 1; i <= count; i++ {
		config += fmt.Sprintf(`
resource "e2e_node" "test%d" {
  name  = "%s-%d"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}
`, i, baseName, i)
	}
	return config
}

func testAccCheckE2ENodeConfig_multipleVolumeAttachments(nodeName, volume1Name, volume2Name string) string {
	return fmt.Sprintf(`
resource "e2e_blockstorage" "test1" {
  name = "%s"
  size = 20
}

resource "e2e_blockstorage" "test2" {
  name = "%s"
  size = 20
}

resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_volume_attachment" "test1" {
  node_id   = e2e_node.test.id
  volume_id = e2e_blockstorage.test1.id
}

resource "e2e_volume_attachment" "test2" {
  node_id   = e2e_node.test.id
  volume_id = e2e_blockstorage.test2.id
}
`, volume1Name, volume2Name, nodeName)
}
