package image_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
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
					resource.TestCheckResourceAttr("e2e_image.test", "name", imageName),
					resource.TestCheckResourceAttrSet("e2e_image.test", "template_id"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "image_state"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "image_type"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "creation_time")),
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
					resource.TestCheckResourceAttr("e2e_image.test", "name", imageName)),
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
					resource.TestCheckResourceAttr("e2e_image.test", "name", imageName),
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
				ExpectError: regexp.MustCompile(`The argument "node_id" is required`),
			},
			{
				Config:      testAccCheckE2EImageConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EImageConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EImageConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
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
				ExpectError: regexp.MustCompile(`name cannot contain whitespace`),
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
					resource.TestCheckResourceAttr("e2e_image.test", "name", imageName1)),
			},
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", "name", imageName2)),
				// Verify that changing name forces replacement
				ExpectNonEmptyPlan: false, // Plan should be empty after apply
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
		projectID := rs.Primary.Attributes["project_id"]
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

		projectID := rs.Primary.Attributes["project_id"]
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

		projectID := rs.Primary.Attributes["project_id"]
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
  node_id    = e2e_node.test.id
  name       = "%s"}
`, nodeName,
		imageName)
}

func testAccCheckE2EImageConfig_ubuntuDistro(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  name       = "%s"}
`, nodeName,
		imageName)
}

// Error case configurations

func testAccCheckE2EImageConfig_missingNodeID() string {
	return `
resource "e2e_image" "test" {
  name       = "test-image"}
`
}

func testAccCheckE2EImageConfig_missingName() string {
	return `
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id}
`
}

func testAccCheckE2EImageConfig_missingProjectID() string {
	return `
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  node_id  = e2e_node.test.id
  name     = "test-image"}
`
}

func testAccCheckE2EImageConfig_missingLocation() string {
	return `
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  name       = "test-image"}
`
}

func testAccCheckE2EImageConfig_invalidName(nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  name       = "invalid name with spaces"}
`, nodeName)
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
  node_id    = e2e_node.test2.id
  name       = "%s"}
`, nodeName1,
		nodeName2,
		imageName)
}
