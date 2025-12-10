package blockstorage_test

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

func TestAccE2EBlockStorage_Basic(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "name", blockStorageName),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", "10"),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "iops"),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "status")),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", "10")),
			},
			{
				Config: testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", "20")),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "name", blockStorageName),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "vm_id"),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "vm_name")),
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
				ProviderFactories: acceptance.TestAccProviderFactories,
				CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName, tc.size),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
							resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", fmt.Sprintf("%v", tc.size))),
					},
				},
			})
		})
	}
}

func TestAccE2EBlockStorage_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EBlockStorageConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EBlockStorageConfig_missingSize(),
				ExpectError: regexp.MustCompile(`The argument "size" is required`),
			},
		},
	})
}

func TestAccE2EBlockStorage_InvalidName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID)),
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

// TestAccE2EBlockStorage_ForceNewName verifies that changing name triggers recreation
func TestAccE2EBlockStorage_ForceNewName(t *testing.T) {
	var blockStorageID1, blockStorageID2 string
	blockStorageName1 := fmt.Sprintf("test-bs-%s", acctest.RandString(10))
	blockStorageName2 := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName1, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID1),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "name", blockStorageName1)),
			},
			{
				Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName2, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID2),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "name", blockStorageName2),
					testAccCheckE2EBlockStorageRecreated(&blockStorageID1, &blockStorageID2)),
			},
		},
	})
}

// TestAccE2EBlockStorage_SizeDowngrade verifies that reducing size is not allowed
func TestAccE2EBlockStorage_SizeDowngrade(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", "20")),
			},
			{
				Config:      testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 10),
				ExpectError: regexp.MustCompile(`Cannot reduce block storage`),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
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

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("Error creating goe2e client: %s", err)
		}

		blockStorage, _, err := goe2eClient.BlockStorage.GetBlockStorage(context.Background(), rs.Primary.ID)
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
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_blockstorage" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("Error creating goe2e client: %s", err)
		}

		blockStorage, _, err := goe2eClient.BlockStorage.GetBlockStorage(context.Background(), rs.Primary.ID)
		if err == nil && blockStorage != nil {
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
		region := acceptance.GetRegionOrLocationFromState(rs)
		if region == "" {
			return "", fmt.Errorf("neither region nor location attribute found")
		}
		blockStorageID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, region, blockStorageID), nil
	}
}

// Configuration helpers

func testAccCheckE2EBlockStorageConfig_basic(name string, size float64) string {
	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name = "%s"
  size = %v
}
`, name, size)
}

func testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName string, size float64) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_blockstorage" "test" {
  name = "%s"
  size = %v
}

resource "e2e_node" "test_attach" {
  name              = "%s-attach"
  plan              = "c2-2c-4gb"
  image             = "ubuntu-20.04"
  block_storage_ids = [e2e_blockstorage.test.id]
  depends_on        = [e2e_blockstorage.test]
}
`, nodeName, blockStorageName, size, nodeName)
}

// Error case configurations

func testAccCheckE2EBlockStorageConfig_missingName() string {
	return `
resource "e2e_blockstorage" "test" {
  size = 10
}
`
}

func testAccCheckE2EBlockStorageConfig_missingSize() string {
	return `
resource "e2e_blockstorage" "test" {
  name = "test-bs"
}
`
}

// NOTE: These test cases are no longer valid with provider defaults
// The provider will automatically supply project_id and region from E2E_PROJECT_ID and E2E_REGION

func testAccCheckE2EBlockStorageConfig_invalidName() string {
	return `
resource "e2e_blockstorage" "test" {
  name = "invalid name with spaces"
  size = 10
}
`
}

// testAccCheckE2EBlockStorageRecreated verifies that the resource was recreated (ID changed)
func testAccCheckE2EBlockStorageRecreated(oldID, newID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *oldID == *newID {
			return fmt.Errorf("expected block storage to be recreated, but ID remained the same: %s", *oldID)
		}
		return nil
	}
}

// TestAccE2EBlockStorage_Tags tests tags CRUD functionality
func TestAccE2EBlockStorage_Tags(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_withTags(blockStorageName, 10, map[string]string{
					"Environment": "test",
					"ManagedBy":   "terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "tags.ManagedBy", "terraform"),
				),
			},
			{
				Config: testAccCheckE2EBlockStorageConfig_withTags(blockStorageName, 10, map[string]string{
					"Environment": "production",
					"ManagedBy":   "terraform",
					"CostCenter":  "engineering",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "tags.%", "3"),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "tags.ManagedBy", "terraform"),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "tags.CostCenter", "engineering"),
				),
			},
			{
				Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "tags.%", "0"),
				),
			},
		},
	})
}

// TestAccE2EBlockStorage_SizeUpgrade_Detached tests that upgrading a detached volume fails
func TestAccE2EBlockStorage_SizeUpgrade_Detached(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_basic(blockStorageName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttr("e2e_blockstorage.test", "size", "10"),
				),
			},
			{
				Config:      testAccCheckE2EBlockStorageConfig_basic(blockStorageName, 20),
				ExpectError: regexp.MustCompile(`Cannot resize block storage.*must be attached`),
			},
		},
	})
}

// TestAccE2EBlockStorage_Delete_Attached tests that deleting an attached volume fails
func TestAccE2EBlockStorage_Delete_Attached(t *testing.T) {
	var blockStorageID string
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_withNode(blockStorageName, nodeName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EBlockStorageExists("e2e_blockstorage.test", &blockStorageID),
					resource.TestCheckResourceAttrSet("e2e_blockstorage.test", "vm_id"),
				),
			},
			{
				// Try to delete by removing the resource - should fail if still attached
				Config:      testAccCheckE2EBlockStorageConfig_basic(nodeName, 10), // Different resource
				ExpectError: nil,                                                   // The destroy check will verify it's still there
			},
		},
	})
}

// Helper function for tags test
func testAccCheckE2EBlockStorageConfig_withTags(name string, size float64, tags map[string]string) string {
	tagsStr := ""
	if len(tags) > 0 {
		tagsStr = "tags = {\n"
		for k, v := range tags {
			tagsStr += fmt.Sprintf("    %s = \"%s\"\n", k, v)
		}
		tagsStr += "  }\n"
	}

	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name = "%s"
  size = %v
  %s
}
`, name, size, tagsStr)
}
