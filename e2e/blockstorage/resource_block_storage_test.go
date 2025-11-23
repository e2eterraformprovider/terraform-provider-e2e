package blockstorage_test

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EBlockStorage_Basic(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "name", blockStorageName),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", "10"),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "iops"),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "status"),
				),
			},
		},
	})
}

func TestAccE2EBlockStorage_Resize(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", "10"),
				),
			},
			{
				Config: testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", "20"),
				),
			},
		},
	})
}

func TestAccE2EBlockStorage_AttachToNode(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "name", blockStorageName),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "vm_id"),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "vm_name"),
				),
			},
		},
	})
}

func TestAccE2EBlockStorage_DifferentSizes(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	testCases := []struct {
		size float64
	}{
		{size: 10},
		{size: 50},
		{size: 100},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("size_%v", tc.size), func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: testAccProviderFactories,
				CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName, tc.size),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
							resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", fmt.Sprintf("%v", tc.size)),
						),
					},
				},
			})
		})
	}
}

func TestAccE2EBlockStorage_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EBlockStorageConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EBlockStorageConfig_missingSize(),
				ExpectError: regexp.MustCompile(`The argument "size" is required`),
			},
			{
				Config:      testAccCheckE2EBlockStorageConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EBlockStorageConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccE2EBlockStorage_InvalidName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EBlockStorageConfig_invalidName(),
				ExpectError: regexp.MustCompile(`the name field cannot be blank, must not contain whitespace or special characters`),
			},
		},
	})
}

func TestAccE2EBlockStorage_Import(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
				),
			},
			{
				ResourceName:      "e2e_blockstorage.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2EBlockStorageImportID("e2e_blockstorage.test"),
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

func testAccCheckE2EBlockStorageExists(resourceName string, blockStorageID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Block Storage ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		projectIDInt, err := strconv.Atoi(projectID)
		if err != nil {
			return fmt.Errorf("Invalid project_id: %s", projectID)
		}

		location := rs.Primary.Attributes["location"]

		blockStorage, err := client.GetBlockStorage(rs.Primary.ID, projectIDInt, location)
		if err != nil {
			return err
		}

		if blockStorage == nil {
			return fmt.Errorf("Block Storage not found")
		}

		*blockStorageID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EBlockStorageDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_blockstorage" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		projectIDInt, err := strconv.Atoi(projectID)
		if err != nil {
			return fmt.Errorf("Invalid project_id: %s", projectID)
		}

		location := rs.Primary.Attributes["location"]

		_, err = client.GetBlockStorage(rs.Primary.ID, projectIDInt, location)
		if err == nil {
			return fmt.Errorf("Block Storage still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccE2EBlockStorageImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]
		blockStorageID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, blockStorageID), nil
	}
}

// Configuration helpers

func testAccCheckE2EBlockStorageConfig_basic(name string, size float64) string {
	projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	projectID, _ := strconv.Atoi(projectIDStr)

	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name       = "%s"
  size       = %v
  project_id = %d
  location   = "%s"
}
`, name, size, projectID, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName string, size float64) string {
	projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	projectID, _ := strconv.Atoi(projectIDStr)

	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_blockstorage" "test" {
  name       = "%s"
  size       = %v
  project_id = %d
  location   = "%s"
}

resource "e2e_node" "test_attach" {
  name              = "%s-attach"
  plan              = "c2-2c-4gb"
  image             = "ubuntu-20.04"
  project_id        = "%s"
  location          = "%s"
  block_storage_ids = [e2e_blockstorage.test.id]
  depends_on        = [e2e_blockstorage.test]
}
`, nodeName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		blockStorageName, size, projectID, os.Getenv("E2E_TEST_LOCATION"),
		nodeName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

// Error case configurations

func testAccCheckE2EBlockStorageConfig_missingName() string {
	projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	projectID, _ := strconv.Atoi(projectIDStr)

	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  size       = 10
  project_id = %d
  location   = "%s"
}
`, projectID, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EBlockStorageConfig_missingSize() string {
	projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	projectID, _ := strconv.Atoi(projectIDStr)

	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name       = "test-bs"
  project_id = %d
  location   = "%s"
}
`, projectID, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EBlockStorageConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name     = "test-bs"
  size     = 10
  location = "%s"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EBlockStorageConfig_missingLocation() string {
	projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	projectID, _ := strconv.Atoi(projectIDStr)

	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name       = "test-bs"
  size       = 10
  project_id = %d
}
`, projectID)
}

func testAccCheckE2EBlockStorageConfig_invalidName() string {
	projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	projectID, _ := strconv.Atoi(projectIDStr)

	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name       = "invalid name with spaces"
  size       = 10
  project_id = %d
  location   = "%s"
}
`, projectID, os.Getenv("E2E_TEST_LOCATION"))
}
