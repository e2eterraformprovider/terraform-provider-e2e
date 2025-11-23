package image_test

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

func TestAccE2EImage_Basic(t *testing.T) {
	var imageID string
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("e2e_image.test", "creation_time"),
				),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", "name", imageName),
				),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_ubuntuDistro(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
					resource.TestCheckResourceAttr("e2e_image.test", "name", imageName),
					resource.TestCheckResourceAttrSet("e2e_image.test", "distro"),
					resource.TestCheckResourceAttrSet("e2e_image.test", "os_distribution"),
				),
			},
		},
	})
}

func TestAccE2EImage_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
		ProviderFactories: testAccProviderFactories,
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImageConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EImageExists("e2e_image.test", &imageID),
				),
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

func testAccCheckE2EImageExists(resourceName string, imageID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Image ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]

		image, err := client.GetImage(rs.Primary.ID, projectID)
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
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_image" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]

		_, err := client.GetImage(rs.Primary.ID, projectID)
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
		location := rs.Primary.Attributes["location"]
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
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  name       = "%s"
  project_id = "%s"
  location   = "%s"
}
`, nodeName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		imageName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EImageConfig_ubuntuDistro(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  name       = "%s"
  project_id = "%s"
  location   = "%s"
}
`, nodeName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		imageName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

// Error case configurations

func testAccCheckE2EImageConfig_missingNodeID() string {
	return fmt.Sprintf(`
resource "e2e_image" "test" {
  name       = "test-image"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EImageConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EImageConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_image" "test" {
  node_id  = e2e_node.test.id
  name     = "test-image"
  location = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EImageConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "test-node"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  name       = "test-image"
  project_id = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccCheckE2EImageConfig_invalidName(nodeName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  name       = "invalid name with spaces"
  project_id = "%s"
  location   = "%s"
}
`, nodeName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
