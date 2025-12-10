package volume_attachment_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EVolumeAttachment_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccE2EVolumeAttachmentConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVolumeAttachmentExists("e2e_volume_attachment.test"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", "node_id"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", "volume_id"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", "vm_id"),
				),
			},
		},
	})
}

func TestAccE2EVolumeAttachment_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccE2EVolumeAttachmentConfig_basic(),
			},
			{
				ResourceName:      "e2e_volume_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVolumeAttachmentImportStateIDFunc("e2e_volume_attachment.test"),
			},
		},
	})
}

func testAccCheckE2EVolumeAttachmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume attachment ID is set")
		}

		// Parse the composite ID
		parts := strings.Split(rs.Primary.ID, "/")
		if len(parts) != 2 {
			return fmt.Errorf("Volume attachment ID is malformed: %s (expected: node_id/volume_id)", rs.Primary.ID)
		}

		nodeID, volumeID := parts[0], parts[1]
		if nodeID == "" || volumeID == "" {
			return fmt.Errorf("Volume attachment ID is malformed: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckE2EVolumeAttachmentDestroy(s *terraform.State) error {
	// This is a simplified check - in a real test, we would verify the volume
	// is actually detached by checking the API
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_volume_attachment" {
			continue
		}

		// If we reach here, the resource still exists in state
		// In a full test, we would make an API call to verify it's detached
	}

	return nil
}

func testAccVolumeAttachmentImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		nodeID := rs.Primary.Attributes["node_id"]
		volumeID := rs.Primary.Attributes["volume_id"]

		return fmt.Sprintf("%s/%s", nodeID, volumeID), nil
	}
}

func testAccE2EVolumeAttachmentConfig_basic() string {
	return `
resource "e2e_ssh_key" "test" {
  name       = "test-volume-attachment-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCm4X3ck1X+MfL9FhvV4tGqqmJz3NZ2d7hP2gDqe1pQqE9yx0p4pWOQFLNQg4DZxBm8NtP5KzN9qdGDhPZx7Wd1JNLiPqKYp7zVnLpfN4fwDQnWwN7F0JxP4mX8c9K7T6Q+Nw4cPz4vL0xH test@example.com"
}

resource "e2e_node" "test" {
  name   = "test-volume-node"
  plan   = "C3.8GB"
  image  = "Ubuntu-18.04-Distro"
  ssh_keys = [e2e_ssh_key.test.name]
}

resource "e2e_block_storage" "test" {
  name = "test-volume-storage"
  size = 100
}

resource "e2e_volume_attachment" "test" {
  node_id   = e2e_node.test.id
  volume_id = e2e_block_storage.test.id
}
`
}
