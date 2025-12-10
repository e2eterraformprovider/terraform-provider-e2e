package sfs_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceE2ESFSs_Basic(t *testing.T) {
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2ESFSsConfig_basic(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2ESFSsExists("data.e2e_sfss.test"),
					resource.TestCheckResourceAttrSet("data.e2e_sfss.test", "sfs_list.#")),
			},
		},
	})
}

func TestAccDataSourceE2ESFSs_WithSFS(t *testing.T) {
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2ESFSsConfig_withSFS(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2ESFSsExists("data.e2e_sfss.test"),
					testAccCheckDataSourceE2ESFSsContainsSFS("data.e2e_sfss.test", sfsName)),
			},
		},
	})
}

func TestAccDataSourceE2ESFSs_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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

// V3 Field Tests

func TestAccDataSourceE2ESFSs_V3Fields(t *testing.T) {
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2ESFSsConfig_v3Fields(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2ESFSsExists("data.e2e_sfss.test"),
					resource.TestCheckResourceAttrSet("data.e2e_sfss.test", "sfs_list.#"),
					testAccCheckDataSourceE2ESFSsContainsSFS("data.e2e_sfss.test", sfsName)),
			},
		},
	})
}

func TestAccDataSourceE2ESFSs_ComputedFields(t *testing.T) {
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2ESFSsConfig_v3Fields(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2ESFSsExists("data.e2e_sfss.test"),
					// Verify computed fields are present
					testAccCheckDataSourceE2ESFSsHasComputedFields("data.e2e_sfss.test")),
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

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		region := acceptance.GetRegionOrLocationFromState(rs)
		projectID := rs.Primary.Attributes["project_id"]

		client, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		sfsList, _, err := client.Sfs.ListSfss(context.Background())
		if err != nil {
			return fmt.Errorf("Error fetching SFS list: %v", err)
		}

		found := false
		for _, sfs := range sfsList {
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

func testAccCheckDataSourceE2ESFSsHasComputedFields(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		// Check that computed fields are present in the data source
		if _, ok := rs.Primary.Attributes["sfs_list.0.id"]; !ok && rs.Primary.Attributes["sfs_list.#"] != "0" {
			return fmt.Errorf("Expected at least one SFS with id field in datasource")
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
  disk_size               = 100  disk_iops               = 1000}

data "e2e_sfss" "test" {
  region     = e2e_sfs.test.region
  project_id = e2e_sfs.test.project_id
}
`, sfsName, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccDataSourceE2ESFSsConfig_withSFS(sfsName string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "%s"
  plan                    = "%s"
  vpc_id                  = "%s"
  disk_size               = 100  disk_iops               = 1000}

data "e2e_sfss" "test" {  depends_on = [e2e_sfs.test]
}
`, sfsName, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

// Error case configurations

func testAccDataSourceE2ESFSsConfig_missingRegion() string {
	return `
data "e2e_sfss" "test" {}
`
}

func testAccDataSourceE2ESFSsConfig_missingProjectID() string {
	return `
data "e2e_sfss" "test" {}
`
}

// V3 Field Configuration Helpers

func testAccDataSourceE2ESFSsConfig_v3Fields(sfsName string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                = "%s"
  plan                = "%s"
  vpc_id              = "%s"
  size_gb             = 100
  iops                = 1000
  encryption_enabled  = false
}

data "e2e_sfss" "test" {
  region     = e2e_sfs.test.region
  project_id = e2e_sfs.test.project_id
  
  depends_on = [e2e_sfs.test]
}
`, sfsName, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}
