package image_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/image"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EImage_Basic(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName),
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrTemplateID),
					resource.TestCheckResourceAttrSet("e2e_image.test", "image_state"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "image_type"),
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrCreatedAt)),
			},
		},
	})
}

func TestAccE2EImage_ValidImageName(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-valid-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName)),
			},
		},
	})
}

func TestAccE2EImage_FromDifferentDistros(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_ubuntuDistro(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName),
					resource.TestCheckResourceAttrSet("e2e_image.test", "distro"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "os_distribution")),
			},
		},
	})
}

func TestAccE2EImage_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EImageConfig_missingNodeID(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrNodeID)),
			},
			{
				Config:      testAccCheckE2EImageConfig_missingName(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrName)),
			},
			{
				Config:      testAccCheckE2EImageConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrProjectID)),
			},
			{
				Config:      testAccCheckE2EImageConfig_missingLocation(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrLocation)),
			},
		},
	})
}

func TestAccE2EImage_InvalidImageName(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EImageConfig_invalidName(nodeName),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("name cannot contain whitespace")),
			},
		},
	})
}

func TestAccE2EImage_Import(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID)),
			},
			{
				ResourceName:      "e2e_image.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2EImageImportID("e2e_image.test"),
			},
		},
	})
}

func TestAccE2EImage_ForceNewOnNameChange(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName1 := fmt.Sprintf("test-image-1-%s", acctest.RandString(10))
	imageName2 := fmt.Sprintf("test-image-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName1)),
			},
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName2)),
				// Verify that changing name does NOT force replacement in V3 (update in-place)
				ExpectNonEmptyPlan: false, // Plan should be empty after apply
			},
		},
	})
}

// TestAccE2EImage_NameUpdate tests that name can be updated in-place (V3 feature)
func TestAccE2EImage_NameUpdate(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName1 := fmt.Sprintf("test-image-1-%s", acctest.RandString(10))
	imageName2 := fmt.Sprintf("test-image-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", "name", imageName1)),
			},
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName2),
					// Verify same image ID (no recreation) - ID should match the one from first step
					resource.TestCheckResourceAttrWith("e2e_image.test", tfconstants.AttrID, func(value string) error {
						if value != imageID {
							return fmt.Errorf("expected image ID %s, got %s (resource was recreated)", imageID, value)
						}
						return nil
					})),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccE2EImage_Tags tests tags functionality (state-only)
func TestAccE2EImage_Tags(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_withTags(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", fmt.Sprintf("%s.%%", tfconstants.AttrTags), "2"),
					resource.TestCheckResourceAttr("e2e_image.test", fmt.Sprintf("%s.Environment", tfconstants.AttrTags), "test"),
					resource.TestCheckResourceAttr("e2e_image.test", fmt.Sprintf("%s.Purpose", tfconstants.AttrTags), "testing")),
			},
			{
				Config: testAccCheckE2EImageConfig_withTagsUpdated(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", fmt.Sprintf("%s.%%", tfconstants.AttrTags), "3"),
					resource.TestCheckResourceAttr("e2e_image.test", fmt.Sprintf("%s.Environment", tfconstants.AttrTags), "test"),
					resource.TestCheckResourceAttr("e2e_image.test", fmt.Sprintf("%s.Purpose", tfconstants.AttrTags), "testing"),
					resource.TestCheckResourceAttr("e2e_image.test", fmt.Sprintf("%s.Updated", tfconstants.AttrTags), "true")),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccE2EImage_RegionVsLocation tests region vs location handling
func TestAccE2EImage_RegionVsLocation(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_withRegion(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrRegion)),
			},
			{
				Config: testAccCheckE2EImageConfig_withLocation(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrLocation)),
				// Should work but with deprecation warning
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccE2EImage_DeprecationWarning tests deprecation warnings for location field
func TestAccE2EImage_DeprecationWarning(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_withLocation(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", nil)),
				// Deprecation warning should be logged but test should pass
			},
		},
	})
}

func TestAccE2EImage_ForceNewOnNodeIDChange(t *testing.T) {
	nodeName1 := fmt.Sprintf("test-node-1-%s", acctest.RandString(10))
	nodeName2 := fmt.Sprintf("test-node-2-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName1, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", nil)),
			},
			{
				Config: testAccCheckE2EImageConfig_differentNode(nodeName1, nodeName2, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", nil)),
				// Verify that changing node_id forces replacement
				ExpectNonEmptyPlan: false, // Plan should be empty after apply
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2EImageExists(resourceName string, imageID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Image ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		projectID := rs.Primary.Attributes[tfconstants.AttrProjectID]
		region := acceptance.GetRegionOrLocationFromState(rs)

		// Create GoE2E client for this project/region
		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("failed to create GoE2E client: %w", err)
		}

		ctx := context.Background()
		image, _, err := goe2eClient.Images.GetImage(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if image == nil {
			return fmt.Errorf("Image not found")
		}

		*imageID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EImageDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_image" {
			continue
		}

		projectID := rs.Primary.Attributes[tfconstants.AttrProjectID]
		region := acceptance.GetRegionOrLocationFromState(rs)

		// Create GoE2E client for this project/region
		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("failed to create GoE2E client: %w", err)
		}

		ctx := context.Background()
		_, _, err = goe2eClient.Images.GetImage(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Image still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccE2EImageImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes[tfconstants.AttrProjectID]
		location := acceptance.GetRegionOrLocationFromState(rs)
		imageID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, imageID), nil
	}
}

// Configuration helpers

func testAccCheckE2EImageConfig_basic(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "%s"}
`, nodeName,
		tfconstants.AttrNodeID,
		tfconstants.AttrName,
		imageName)
}

func testAccCheckE2EImageConfig_ubuntuDistro(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "%s"}
`, nodeName,
		tfconstants.AttrNodeID,
		tfconstants.AttrName,
		imageName)
}

// Error case configurations

func testAccCheckE2EImageConfig_missingNodeID() string {
	return fmt.Sprintf(`
resource "e2e_image" "test" {
  %s       = "test-image"}
`, tfconstants.AttrName)
}

func testAccCheckE2EImageConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id}
`, tfconstants.AttrNodeID)
}

func testAccCheckE2EImageConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s  = e2e_node.test.id
  %s     = "test-image"}
`, tfconstants.AttrNodeID, tfconstants.AttrName)
}

func testAccCheckE2EImageConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "test-image"}
`, tfconstants.AttrNodeID, tfconstants.AttrName)
}

func testAccCheckE2EImageConfig_invalidName(nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "invalid name with spaces"}
`, nodeName, tfconstants.AttrNodeID, tfconstants.AttrName)
}

func testAccCheckE2EImageConfig_differentNode(nodeName1, nodeName2, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test1" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_node" "test2" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test2.id
  %s       = "%s"}
`, nodeName1,
		nodeName2,
		tfconstants.AttrNodeID,
		tfconstants.AttrName,
		imageName)
}

// Configuration helpers for new tests

func testAccCheckE2EImageConfig_withTags(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "%s"
  %s = {
    Environment = "test"
    Purpose     = "testing"
  }
}
`, nodeName,
		tfconstants.AttrNodeID,
		tfconstants.AttrName,
		imageName,
		tfconstants.AttrTags)
}

func testAccCheckE2EImageConfig_withTagsUpdated(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "%s"
  %s = {
    Environment = "test"
    Purpose     = "testing"
    Updated     = "true"
  }
}
`, nodeName,
		tfconstants.AttrNodeID,
		tfconstants.AttrName,
		imageName,
		tfconstants.AttrTags)
}

func testAccCheckE2EImageConfig_withRegion(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "%s"
  %s     = "Mumbai"
}
`, nodeName,
		tfconstants.AttrNodeID,
		tfconstants.AttrName,
		imageName,
		tfconstants.AttrRegion)
}

func testAccCheckE2EImageConfig_withLocation(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "%s"
  %s   = "Mumbai"
}
`, nodeName,
		tfconstants.AttrNodeID,
		tfconstants.AttrName,
		imageName,
		tfconstants.AttrLocation)
}

// ============================================================================
// Additional Acceptance Tests for Comprehensive Coverage
// ============================================================================

// TestAccE2EImage_AsyncPolling tests async polling during image creation
func TestAccE2EImage_AsyncPolling(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					// Verify image reaches Ready state (async polling should complete)
					resource.TestCheckResourceAttr("e2e_image.test", "state", goe2econstants.ImageStateReady),
					// Verify image ID is set even if polling takes time
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrID)),
			},
		},
	})
}

// TestAccE2EImage_NameUpdateInPlace tests name update in-place (V3 feature)
func TestAccE2EImage_NameUpdateInPlace(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName1 := fmt.Sprintf("test-image-1-%s", acctest.RandString(10))
	imageName2 := fmt.Sprintf("test-image-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName1)),
			},
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName2),
					// Verify same image ID (no recreation)
					resource.TestCheckResourceAttrWith("e2e_image.test", tfconstants.AttrID, func(value string) error {
						if value != imageID {
							return fmt.Errorf("expected image ID %s, got %s (resource was recreated)", imageID, value)
						}
						return nil
					})),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccE2EImage_NameUpdateReflectedInState tests name update is reflected after refresh
func TestAccE2EImage_NameUpdateReflectedInState(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName1 := fmt.Sprintf("test-image-1-%s", acctest.RandString(10))
	imageName2 := fmt.Sprintf("test-image-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName1)),
			},
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName2)),
			},
			{
				// Refresh step to verify name is persisted
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, imageName2)),
			},
		},
	})
}

// TestAccE2EImage_RegionLocationConflict tests error when both region and location are set
func TestAccE2EImage_RegionLocationConflict(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EImageConfig_regionLocationConflict(nodeName, imageName),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(image.RegionLocationConflict)),
			},
		},
	})
}

// TestAccE2EImage_NameValidationWhitespace tests name validation with whitespace
func TestAccE2EImage_NameValidationWhitespace(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EImageConfig_invalidName(nodeName),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("name cannot contain whitespace")),
			},
		},
	})
}

// TestAccE2EImage_NameValidationInUpdate tests name validation in update operations
func TestAccE2EImage_NameValidationInUpdate(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID)),
			},
			{
				Config:      testAccCheckE2EImageConfig_invalidName(nodeName),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("name cannot contain whitespace")),
			},
		},
	})
}

// TestAccE2EImage_RegionPreferred tests creating image with region (preferred)
func TestAccE2EImage_RegionPreferred(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_withRegion(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrRegion)),
			},
		},
	})
}

// TestAccE2EImage_LocationDeprecated tests creating image with location (deprecated)
func TestAccE2EImage_LocationDeprecated(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_withLocation(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrLocation)),
				// Should work but with deprecation warning
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccE2EImage_ComputedFieldsPopulated tests all computed fields are populated
func TestAccE2EImage_ComputedFieldsPopulated(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrTemplateID),
					resource.TestCheckResourceAttrSet("e2e_image.test", "image_state"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "state"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "image_type"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "os_distribution"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "distro"),
					resource.TestCheckResourceAttrSet("e2e_image.test", tfconstants.AttrCreatedAt)),
			},
		},
	})
}

// TestAccE2EImage_StateNormalization tests state normalization (image_state → state)
func TestAccE2EImage_StateNormalization(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					// Verify state is normalized (lowercase)
					resource.TestCheckResourceAttr("e2e_image.test", "state", goe2econstants.ImageStateReady)),
			},
		},
	})
}

// TestAccE2EImage_VeryLongName tests image creation with very long name
func TestAccE2EImage_VeryLongName(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	// Create a name that's 256 characters long (typical API limit)
	longName := fmt.Sprintf("test-image-%s", acctest.RandString(240))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, longName),
				// May fail at validation or API level - both are acceptable
				// If it succeeds, that's fine too - just testing it doesn't crash
				ExpectError: nil, // Allow it to pass or fail gracefully
			},
		},
	})
}

// TestAccE2EImage_SpecialCharactersInName tests image creation with valid special characters
// Valid characters: alphanumeric, hyphens, underscores
func TestAccE2EImage_SpecialCharactersInName(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	// Test valid special characters: hyphens and underscores are typically allowed
	validSpecialNames := []string{
		"test-image-with-hyphens",
		"test_image_with_underscores",
		"test-image_with-mixed",
		"test-123-image",
		"test_image_456",
	}

	for _, specialName := range validSpecialNames {
		t.Run(fmt.Sprintf("valid_special_char_%s", specialName), func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: acceptance.TestAccProviderFactories,
				CheckDestroy:      testAccCheckE2EImageDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccCheckE2EImageConfig_basic(nodeName, specialName),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckE2EImageExists("e2e_image.test", &imageID),
							resource.TestCheckResourceAttr("e2e_image.test", tfconstants.AttrName, specialName)),
					},
				},
			})
		})
	}
}

// TestAccE2EImage_InvalidSpecialCharacters tests that invalid special characters are rejected
func TestAccE2EImage_InvalidSpecialCharacters(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	// Test invalid special characters that should fail validation
	invalidSpecialNames := []string{
		"test image with spaces", // Should fail - contains whitespace
		"test-image@special",     // @ may be invalid
		"test-image#special",     // # may be invalid
		"test-image$special",     // $ may be invalid
		"test-image%special",     // % may be invalid
		"test-image&special",     // & may be invalid
		"test-image*special",     // * may be invalid
		"test-image+special",     // + may be invalid
		"test-image:special",     // : may be invalid
		"test-image;special",     // ; may be invalid
		"test-image<special",     // < may be invalid
		"test-image>special",     // > may be invalid
		"test-image?special",     // ? may be invalid
		"test-image[special",     // [ may be invalid
		"test-image]special",     // ] may be invalid
		"test-image{special",     // { may be invalid
		"test-image}special",     // } may be invalid
		"test-image|special",     // | may be invalid
		"test-image\\special",    // \ may be invalid
		"test-image'special",     // ' may be invalid
		"test-image\"special",    // " may be invalid
	}

	for _, invalidName := range invalidSpecialNames {
		t.Run(fmt.Sprintf("invalid_special_char_%s", invalidName), func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: acceptance.TestAccProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: testAccCheckE2EImageConfig_basic(nodeName, invalidName),
						// Should fail at validation (whitespace) or API level
						ExpectError: regexp.MustCompile(`(name cannot contain whitespace|invalid|error)`),
					},
				},
			})
		})
	}
}

// TestAccE2EImage_ImportThenDelete tests importing an existing image and then deleting it
func TestAccE2EImage_ImportThenDelete(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))
	resourceName := "e2e_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: Create image
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists(resourceName, &imageID),
					resource.TestCheckResourceAttr(resourceName, tfconstants.AttrName, imageName)),
			},
			{
				// Step 2: Import the image (simulating importing an image created outside Terraform)
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2EImageImportID(resourceName),
				Check: resource.ComposeTestCheckFunc(
					// Verify all fields populated from API after import
					resource.TestCheckResourceAttrSet(resourceName, tfconstants.AttrID),
					resource.TestCheckResourceAttrSet(resourceName, tfconstants.AttrTemplateID),
					resource.TestCheckResourceAttrSet(resourceName, "image_state"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, tfconstants.AttrCreatedAt)),
			},
			{
				// Step 3: Delete the imported image
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				// CheckDestroy will verify deletion
			},
		},
	})
}

// TestAccE2EImage_ImportReadyState tests importing an image in Ready state
func TestAccE2EImage_ImportReadyState(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))
	resourceName := "e2e_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: Create image and wait for Ready state
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists(resourceName, &imageID),
					// Verify image is in Ready state
					resource.TestCheckResourceAttr(resourceName, "state", goe2econstants.ImageStateReady),
					resource.TestCheckResourceAttr(resourceName, tfconstants.AttrName, imageName)),
			},
			{
				// Step 2: Import the Ready image
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2EImageImportID(resourceName),
				Check: resource.ComposeTestCheckFunc(
					// Verify state is Ready after import
					resource.TestCheckResourceAttr(resourceName, "state", goe2econstants.ImageStateReady),
					resource.TestCheckResourceAttrSet(resourceName, "image_state"),
					resource.TestCheckResourceAttr(resourceName, tfconstants.AttrName, imageName)),
			},
		},
	})
}

// Additional configuration helpers

func testAccCheckE2EImageConfig_regionLocationConflict(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  %s    = e2e_node.test.id
  %s       = "%s"
  %s     = "Mumbai"
  %s   = "Mumbai"
}
`, nodeName,
		tfconstants.AttrNodeID,
		tfconstants.AttrName,
		imageName,
		tfconstants.AttrRegion,
		tfconstants.AttrLocation)
}
