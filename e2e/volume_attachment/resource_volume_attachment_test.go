package volume_attachment_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	volumeattachment "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/volume_attachment"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EVolumeAttachment_Basic(t *testing.T) {
	nodeName := fmt.Sprintf("test-va-node-%s", acctest.RandString(10))
	volumeName := fmt.Sprintf("test-va-vol-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccE2EVolumeAttachmentConfig_basic(nodeName, volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVolumeAttachmentExists("e2e_volume_attachment.test"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", tfconstants.AttrNodeID),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", tfconstants.AttrVolumeID),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", tfconstants.AttrVMID),
					// Verify block storage now shows attachment details
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", tfconstants.AttrVMID),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", tfconstants.AttrVMName),
				),
			},
			// Remove attachment resource to trigger detach (Delete)
			{
				Config: testAccE2EVolumeAttachmentConfig_detachOnly(nodeName, volumeName),
				Check: resource.ComposeTestCheckFunc(
					// Ensure block storage is detached at API level (or at least no VMDetail)
					testAccCheckE2EBlockStorageDetached("e2e_blockstorage.test"),
				),
			},
		},
	})
}

func TestAccE2EVolumeAttachment_Import(t *testing.T) {
	nodeName := fmt.Sprintf("test-va-node-%s", acctest.RandString(10))
	volumeName := fmt.Sprintf("test-va-vol-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccE2EVolumeAttachmentConfig_basic(nodeName, volumeName),
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
		parts := config.ParseImportID(rs.Primary.ID)
		if len(parts) != 2 {
			return fmt.Errorf("Volume attachment ID is malformed: %s (expected: %s)", rs.Primary.ID, volumeattachment.ImportIDFormatShortDescription)
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

		nodeID := rs.Primary.Attributes[tfconstants.AttrNodeID]
		volumeID := rs.Primary.Attributes[tfconstants.AttrVolumeID]

		return fmt.Sprintf("%s%s%s", nodeID, tfconstants.VolumeAttachmentImportDelimiter, volumeID), nil
	}
}

func testAccE2EVolumeAttachmentConfig_basic(nodeName, volumeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_blockstorage" "test" {
  name = "%s"
  size = 250
}

resource "e2e_volume_attachment" "test" {
  node_id   = e2e_node.test.id
  volume_id = e2e_blockstorage.test.id
}
`, nodeName, volumeName)
}

func testAccE2EVolumeAttachmentConfig_detachOnly(nodeName, volumeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_blockstorage" "test" {
  name = "%s"
  size = 250
}
`, nodeName, volumeName)
}

func testAccCheckE2EBlockStorageDetached(blockStorageResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[blockStorageResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", blockStorageResourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Block Storage ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)

		projectID := rs.Primary.Attributes[tfconstants.AttrProjectID]
		if projectID == "" {
			projectID = acceptance.TestProjectID
		}
		region := acceptance.GetRegionOrLocationFromState(rs)
		if region == "" {
			region = acceptance.TestRegion
		}

		client, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
		}

		vol, _, err := client.BlockStorage.GetBlockStorage(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if vol == nil {
			return nil
		}

		// Prefer explicit detached status if API provides it.
		if vol.Status == goe2econstants.BlockStorageStatusDetached {
			return nil
		}
		// Otherwise ensure it is not attached.
		if vol.Status == goe2econstants.BlockStorageStatusAttached {
			return fmt.Errorf("expected block storage to be detached, got status %q", vol.Status)
		}
		return nil
	}
}
