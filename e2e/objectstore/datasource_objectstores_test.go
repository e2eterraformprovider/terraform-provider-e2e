package objectstore_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceE2EObjectStores_Basic(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EObjectStoresConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EObjectStoresExists("data.e2e_objectstores.test"),
					resource.TestCheckResourceAttrSet("data.e2e_objectstores.test", "bucket_list.#"),
				),
			},
		},
	})
}

func TestAccDataSourceE2EObjectStores_WithBucket(t *testing.T) {
	bucketName := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EObjectStoresConfig_withBucket(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EObjectStoresExists("data.e2e_objectstores.test"),
					testAccCheckDataSourceE2EObjectStoresContainsBucket("data.e2e_objectstores.test", bucketName),
				),
			},
		},
	})
}

func TestAccDataSourceE2EObjectStores_MultipleBuckets(t *testing.T) {
	bucketName1 := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))
	bucketName2 := fmt.Sprintf("test-bucket-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EObjectStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EObjectStoresConfig_multipleBuckets(bucketName1, bucketName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EObjectStoresExists("data.e2e_objectstores.test"),
					testAccCheckDataSourceE2EObjectStoresContainsBucket("data.e2e_objectstores.test", bucketName1),
					testAccCheckDataSourceE2EObjectStoresContainsBucket("data.e2e_objectstores.test", bucketName2),
				),
			},
		},
	})
}

func TestAccDataSourceE2EObjectStores_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2EObjectStoresConfig_missingRegion(),
				ExpectError: regexp.MustCompile(`The argument "region" is required`),
			},
			{
				Config:      testAccDataSourceE2EObjectStoresConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
		},
	})
}

func TestAccDataSourceE2EObjectStores_EmptyList(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EObjectStoresConfig_empty(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EObjectStoresExists("data.e2e_objectstores.test"),
					resource.TestCheckResourceAttrSet("data.e2e_objectstores.test", "bucket_list.#"),
				),
			},
		},
	})
}

// Helper functions

func testAccCheckDataSourceE2EObjectStoresExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Object Stores datasource ID is set")
		}

		return nil
	}
}

func testAccCheckDataSourceE2EObjectStoresContainsBucket(resourceName, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		region := rs.Primary.Attributes["region"]
		projectID := rs.Primary.Attributes["project_id"]

		response, err := client.GetBuckets(region, projectID)
		if err != nil {
			return fmt.Errorf("Error fetching object stores: %v", err)
		}

		found := false
		for _, bucket := range response.Data {
			if bucket.Name == bucketName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Bucket %s not found in datasource", bucketName)
		}

		return nil
	}
}

// Configuration helpers

func testAccDataSourceE2EObjectStoresConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name       = "%s"
  project_id = %s
  region     = "%s"
}

data "e2e_objectstores" "test" {
  region     = e2e_objectstore.test.region
  project_id = e2e_objectstore.test.project_id
}
`, bucketName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccDataSourceE2EObjectStoresConfig_withBucket(bucketName string) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test" {
  name       = "%s"
  project_id = %s
  region     = "%s"
}

data "e2e_objectstores" "test" {
  region     = "%s"
  project_id = %s

  depends_on = [e2e_objectstore.test]
}
`, bucketName, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"),
		os.Getenv("E2E_TEST_REGION"), os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccDataSourceE2EObjectStoresConfig_multipleBuckets(bucketName1, bucketName2 string) string {
	return fmt.Sprintf(`
resource "e2e_objectstore" "test1" {
  name       = "%s"
  project_id = %s
  region     = "%s"
}

resource "e2e_objectstore" "test2" {
  name       = "%s"
  project_id = %s
  region     = "%s"
}

data "e2e_objectstores" "test" {
  region     = "%s"
  project_id = %s

  depends_on = [e2e_objectstore.test1, e2e_objectstore.test2]
}
`, bucketName1, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"),
		bucketName2, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"),
		os.Getenv("E2E_TEST_REGION"), os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccDataSourceE2EObjectStoresConfig_empty() string {
	return fmt.Sprintf(`
data "e2e_objectstores" "test" {
  region     = "%s"
  project_id = %s
}
`, os.Getenv("E2E_TEST_REGION"), os.Getenv("E2E_TEST_PROJECT_ID"))
}

// Error case configurations

func testAccDataSourceE2EObjectStoresConfig_missingRegion() string {
	return fmt.Sprintf(`
data "e2e_objectstores" "test" {
  project_id = %s
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccDataSourceE2EObjectStoresConfig_missingProjectID() string {
	return fmt.Sprintf(`
data "e2e_objectstores" "test" {
  region = "%s"
}
`, os.Getenv("E2E_TEST_REGION"))
}
