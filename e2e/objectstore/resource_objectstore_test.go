package objectstore_test

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

func TestAccE2EObjectStore_Basic(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "name", bucketName),
					resource.TestCheckResourceAttrSet("e2e_objectstore.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_objectstore.test", "created_on"),
					resource.TestCheckResourceAttrSet("e2e_objectstore.test", "versioning_status"),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "false"),
				),
			},
		},
	})
}

func TestAccE2EObjectStore_Update(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "false"),
				),
			},
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true"),
				),
			},
		},
	})
}

func TestAccE2EObjectStore_NameChange(t *testing.T) {
	bucketName1 := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))
	bucketName2 := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_objectstore.test", "name", bucketName1),
				),
			},
			{
				Config:      testAccCheckE2EObjectStoreConfig_basic(bucketName2),
				ExpectError: regexp.MustCompile(`cannot change the bucket name`),
			},
		},
	})
}

func TestAccE2EObjectStore_WithVersioning(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true"),
				),
			},
		},
	})
}

func TestAccE2EObjectStore_VersioningToggle(t *testing.T) {
	var bucketID string
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true"),
				),
			},
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "false"),
				),
			},
			{
				Config: testAccCheckE2EObjectStoreConfig_withVersioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
					resource.TestCheckResourceAttr("e2e_objectstore.test", "enabling_versioning", "true"),
				),
			},
		},
	})
}

func TestAccE2EObjectStore_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EObjectStoreConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EObjectStoreExists("e2e_objectstore.test", &bucketID),
				),
			},
			{
				ResourceName:            "e2e_objectstore.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enabling_versioning"},
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
	if v := os.Getenv("E2E_TEST_REGION"); v == "" {
		t.Fatal("E2E_TEST_REGION must be set for acceptance tests")
	}
}

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"e2e": func() (*schema.Provider, error) {
		return e2e.Provider(), nil
	},
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

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		bucketName := rs.Primary.Attributes["name"]
		region := rs.Primary.Attributes["region"]
		projectID := rs.Primary.Attributes["project_id"]

		bucket, err := client.GetBucket(bucketName, region, projectID)
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
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_objectstore" {
			continue
		}

		bucketName := rs.Primary.Attributes["name"]
		region := rs.Primary.Attributes["region"]
		projectID := rs.Primary.Attributes["project_id"]

		_, err := client.GetBucket(bucketName, region, projectID)
		if err == nil {
			return fmt.Errorf("Object Store bucket still exists: %s", bucketName)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2EObjectStoreConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name       = "%s"
  project_id = %s
  region     = "%s"
}
`, bucketName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2EObjectStoreConfig_withVersioning(bucketName string, enableVersioning bool) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name                = "%s"
  project_id          = %s
  region              = "%s"
  enabling_versioning = %t
}
`, bucketName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"), enableVersioning)
}

// Error case configurations

func testAccCheckE2EObjectStoreConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  project_id = %s
  region     = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2EObjectStoreConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name   = "test-bucket"
  region = "%s"
}
`, os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2EObjectStoreConfig_missingRegion() string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name       = "test-bucket"
  project_id = %s
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}
