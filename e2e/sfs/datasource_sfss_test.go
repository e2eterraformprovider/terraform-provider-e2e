package sfs_test

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

func TestAccDataSourceE2ESFSs_Basic(t *testing.T) {
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2ESFSsConfig_basic(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2ESFSsExists("data.e2e_sfss.test"),
					resource.TestCheckResourceAttrSet("data.e2e_sfss.test", "sfs_list.#"),
				),
			},
		},
	})
}

func TestAccDataSourceE2ESFSs_WithSFS(t *testing.T) {
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2ESFSsConfig_withSFS(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2ESFSsExists("data.e2e_sfss.test"),
					testAccCheckDataSourceE2ESFSsContainsSFS("data.e2e_sfss.test", sfsName),
				),
			},
		},
	})
}

func TestAccDataSourceE2ESFSs_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2ESFSsConfig_missingRegion(),
				ExpectError: regexp.MustCompile(`The argument "region" is required`),
			},
			{
				Config:      testAccDataSourceE2ESFSsConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
		},
	})
}

// Helper functions

func testAccCheckDataSourceE2ESFSsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SFS datasource ID is set")
		}

		return nil
	}
}

func testAccCheckDataSourceE2ESFSsContainsSFS(resourceName, sfsName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		region := rs.Primary.Attributes["region"]
		projectID := rs.Primary.Attributes["project_id"]

		response, err := client.GetSfss(region, projectID)
		if err != nil {
			return fmt.Errorf("Error fetching SFS list: %v", err)
		}

		found := false
		for _, sfs := range response.Data {
			if sfs.Name == sfsName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("SFS %s not found in datasource", sfsName)
		}

		return nil
	}
}

// Configuration helpers

func testAccDataSourceE2ESFSsConfig_basic(sfsName string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "%s"
  plan                    = "%s"
  vpc_id                  = "%s"
  disk_size               = 100
  project_id              = "%s"
  disk_iops               = 1000
  region                  = "%s"
}

data "e2e_sfss" "test" {
  region     = e2e_sfs.test.region
  project_id = e2e_sfs.test.project_id
}
`, sfsName, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccDataSourceE2ESFSsConfig_withSFS(sfsName string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "%s"
  plan                    = "%s"
  vpc_id                  = "%s"
  disk_size               = 100
  project_id              = "%s"
  disk_iops               = 1000
  region                  = "%s"
}

data "e2e_sfss" "test" {
  region     = "%s"
  project_id = "%s"

  depends_on = [e2e_sfs.test]
}
`, sfsName, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"),
		os.Getenv("E2E_TEST_REGION"), os.Getenv("E2E_TEST_PROJECT_ID"))
}

// Error case configurations

func testAccDataSourceE2ESFSsConfig_missingRegion() string {
	return fmt.Sprintf(`
data "e2e_sfss" "test" {
  project_id = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccDataSourceE2ESFSsConfig_missingProjectID() string {
	return fmt.Sprintf(`
data "e2e_sfss" "test" {
  region = "%s"
}
`, os.Getenv("E2E_TEST_REGION"))
}
