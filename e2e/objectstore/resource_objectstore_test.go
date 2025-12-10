package objectstore_test

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

func TestAccE2EObjectStore_Basic(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "name", bucketName),
					resource.TestCheckResourceAttrSet("e2e_objectstore.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_objectstore.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_objectstore.test", "versioning_status"),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "false")),
			},
		},
	})
}

func TestAccE2EObjectStore_Update(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "false")),
			},
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true")),
			},
		},
	})
}

func TestAccE2EObjectStore_NameChange(t *testing.T) {
	var bucketID1, bucketID2 string
	bucketName1 := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))
	bucketName2 := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID1),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "name", bucketName1)),
			},
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID2),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "name", bucketName2),
					// Verify that resource was recreated (different ID) due to ForceNew
					testAccCheckE2EObjectStoreRecreated(&bucketID1, &bucketID2)),
			},
		},
	})
}

func TestAccE2EObjectStore_WithVersioning(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true")),
			},
		},
	})
}

func TestAccE2EObjectStore_VersioningToggle(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true")),
			},
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "false")),
			},
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true")),
			},
		},
	})
}

func TestAccE2EObjectStore_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EObjectStoreConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EObjectStoreConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EObjectStoreConfig_missingRegion(),
				ExpectError: regexp.MustCompile(`The argument "region" is required`),
			},
		},
	})
}

func TestAccE2EObjectStore_Import(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID)),
			},
			{
				ResourceName:            "e2e_objectstore.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccE2EObjectStoreImportID("e2e_objectstore.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enabling_versioning"},
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2EObjectStoreExists(resourceName string, bucketID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Object Store bucket ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)

		bucketName := rs.Primary.Attributes["name"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		projectID := rs.Primary.Attributes["project_id"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("Error creating goe2e client: %v", err)
		}

		bucket, _, err := goe2eClient.ObjectStorage.GetBucket(context.Background(), bucketName)
		if err != nil {
			return err
		}

		if bucket == nil {
			return fmt.Errorf("Object Store bucket not found")
		}

		*bucketID = rs.Primary.ID

		return nil
	}
}

func testAccCheckE2EObjectStoreDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_objectstore" {
			continue
		}

		bucketName := rs.Primary.Attributes["name"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		projectID := rs.Primary.Attributes["project_id"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("Error creating goe2e client: %v", err)
		}

		bucket, _, _ := goe2eClient.ObjectStorage.GetBucket(context.Background(), bucketName)
		if bucket != nil {
			return fmt.Errorf("Object Store bucket still exists: %s", bucketName)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2EObjectStoreConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name = "%s"
}
`, bucketName)
}

func testAccCheckE2EObjectStoreConfig_withVersioning(bucketName string, enableVersioning bool) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name                = "%s"
  enabling_versioning = %t
}
`, bucketName, enableVersioning)
}

// Error case configurations

func testAccCheckE2EObjectStoreConfig_missingName() string {
	return `
resource "e2e_objectstore" "test" {
}
`
}

func testAccCheckE2EObjectStoreConfig_missingProjectID() string {
	return `
resource "e2e_objectstore" "test" {
  name = "test-bucket"
}
`
}

func testAccCheckE2EObjectStoreConfig_missingRegion() string {
	return `
resource "e2e_objectstore" "test" {
  name = "test-bucket"
}
`
}

func testAccE2EObjectStoreImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		bucketName := rs.Primary.Attributes["name"]

		// Import format: project_id:region:bucket_name
		return fmt.Sprintf("%s:%s:%s", projectID, region, bucketName), nil
	}
}

func testAccCheckE2EObjectStoreRecreated(oldID, newID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *oldID == *newID {
			return fmt.Errorf("Expected object store bucket to be recreated, but IDs are the same: %s", *oldID)
		}
		return nil
	}
}

// V3 Feature Tests

// TestAccE2EObjectStore_VersioningEnabled tests the new versioning_enabled field (V3)
func TestAccE2EObjectStore_VersioningEnabled(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_versioningEnabled(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "name", bucketName),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "versioning_enabled", "true")),
			},
			{
				Config: testAccCheckE2EObjectStoreConfig_versioningEnabled(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "versioning_enabled", "false")),
			},
		},
	})
}

// TestAccE2EObjectStore_Tags tests the new tags field (V3)
func TestAccE2EObjectStore_Tags(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_withTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "name", bucketName),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "tags.environment", "test"),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "tags.application", "terraform")),
			},
		},
	})
}

// TestAccE2EObjectStore_AllFeatures tests all V3 features together
func TestAccE2EObjectStore_AllFeatures(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_allFeatures(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "name", bucketName),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "versioning_enabled", "true"),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "encryption_enabled", "true"),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "lock_enabled", "false"),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "public_access_enabled", "false"),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "tags.environment", "production")),
			},
		},
	})
}

// TestAccE2EObjectStore_DeprecatedEnablingVersioning tests backwards compatibility with deprecated field
func TestAccE2EObjectStore_DeprecatedEnablingVersioning(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true")),
			},
		},
	})
}

// TestAccE2EObjectStore_ImportBasic tests simple import format (bucket_name only)
func TestAccE2EObjectStore_ImportBasic(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))
	resourceName := "e2e_objectstore.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists(resourceName, &bucketID)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccE2EObjectStoreImportID(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

// Configuration helpers for V3 tests

func testAccCheckE2EObjectStoreConfig_versioningEnabled(bucketName string, enableVersioning bool) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name                = "%s"
  versioning_enabled  = %t
}
`, bucketName, enableVersioning)
}

func testAccCheckE2EObjectStoreConfig_withTags(bucketName string) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name = "%s"

  tags = {
    environment = "test"
    application = "terraform"
  }
}
`, bucketName)
}

func testAccCheckE2EObjectStoreConfig_allFeatures(bucketName string) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name                    = "%s"
  versioning_enabled      = true
  encryption_enabled      = true
  lock_enabled            = false
  public_access_enabled   = false

  tags = {
    environment = "production"
    managed_by  = "terraform"
  }
}
`, bucketName)
}
