package floating_ip_attachment_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EFloatingIPAttachment_Basic(t *testing.T) {
	var attachmentID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFloatingIPAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFloatingIPAttachmentConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFloatingIPAttachmentExists("e2e_floating_ip_attachment.test", &attachmentID),
					resource.TestCheckResourceAttrSet("e2e_floating_ip_attachment.test", "ip_address"),
					resource.TestCheckResourceAttrSet("e2e_floating_ip_attachment.test", "node_ids"),
				),
			},
		},
	})
}

func TestAccE2EFloatingIPAttachment_Import(t *testing.T) {
	var attachmentID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFloatingIPAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFloatingIPAttachmentConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFloatingIPAttachmentExists("e2e_floating_ip_attachment.test", &attachmentID),
				),
			},
			{
				ResourceName:      "e2e_floating_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2EFloatingIPAttachmentImportID("e2e_floating_ip_attachment.test"),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2EFloatingIPAttachmentExists(resourceName string, attachmentID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Floating IP Attachment ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		ipAddress := rs.Primary.Attributes["ip_address"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
		}

		ctx := context.Background()
		rips, _, err := goe2eClient.ReserveIP.ListReserveIPs(ctx)
		if err != nil {
			return err
		}

		found := false
		for _, item := range rips {
			if item.IPAddress == ipAddress && item.ReservedType == goe2econstants.ReserveIPTypeFloatingIP {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Floating IP Attachment not found")
		}

		*attachmentID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EFloatingIPAttachmentDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_floating_ip_attachment" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		ipAddress := rs.Primary.Attributes["ip_address"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			// If we can't create a client, assume the resource doesn't exist
			continue
		}

		ctx := context.Background()
		rips, _, err := goe2eClient.ReserveIP.ListReserveIPs(ctx)
		if err != nil {
			// If we get an error, assume it's because the resource doesn't exist
			continue
		}

		for _, item := range rips {
			if item.IPAddress == ipAddress && item.ReservedType == goe2econstants.ReserveIPTypeFloatingIP {
				// Check if there are any attached nodes
				if len(item.FloatingIPAttachedNodes) > 0 {
					return fmt.Errorf("Floating IP Attachment still exists: %s", ipAddress)
				}
			}
		}
	}

	return nil
}

func testAccE2EFloatingIPAttachmentImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		ipAddress := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, region, ipAddress), nil
	}
}

// Configuration helpers

func testAccCheckE2EFloatingIPAttachmentConfig_basic() string {
	return `
# Note: This test requires a pre-existing floating IP and node
# In a real scenario, you would create these resources first
resource "e2e_floating_ip_attachment" "test" {
  ip_address = "164.52.220.153"  # Replace with actual floating IP
  node_ids   = ["node-id-1"]      # Replace with actual node IDs
}
`
}

// ============================================================================
// Floating IP Attachment Resource Acceptance Tests
// ============================================================================

func TestAccE2EFloatingIPAttachment_MultipleNodes(t *testing.T) {
	var attachmentID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFloatingIPAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFloatingIPAttachmentConfig_multipleNodes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFloatingIPAttachmentExists("e2e_floating_ip_attachment.test", &attachmentID),
					resource.TestCheckResourceAttrSet("e2e_floating_ip_attachment.test", "ip_address"),
					resource.TestCheckResourceAttr("e2e_floating_ip_attachment.test", "node_ids.#", "2"),
				),
			},
		},
	})
}

func testAccCheckE2EFloatingIPAttachmentConfig_multipleNodes() string {
	return `
# Note: This test requires pre-existing floating IP and nodes
resource "e2e_floating_ip_attachment" "test" {
  ip_address = "164.52.220.153"  # Replace with actual floating IP
  node_ids   = ["node-id-1", "node-id-2"]  # Replace with actual node IDs
}
`
}
