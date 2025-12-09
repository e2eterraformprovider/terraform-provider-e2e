package autoscaling_test

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

func TestAccDataSourceE2EScalerGroup_Basic(t *testing.T) {
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EScalerGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EScalerGroupExists("data.e2e_scaler_group.test"),
					resource.TestCheckResourceAttrPair("data.e2e_scaler_group.test", "id", "e2e_scaler_group.test", "id"),
					resource.TestCheckResourceAttrPair("data.e2e_scaler_group.test", "name", "e2e_scaler_group.test", "name"),
					resource.TestCheckResourceAttrPair("data.e2e_scaler_group.test", "min_nodes", "e2e_scaler_group.test", "min_nodes"),
					resource.TestCheckResourceAttrPair("data.e2e_scaler_group.test", "max_nodes", "e2e_scaler_group.test", "max_nodes"),
					resource.TestCheckResourceAttrPair("data.e2e_scaler_group.test", "desired", "e2e_scaler_group.test", "desired"),
					resource.TestCheckResourceAttrSet("data.e2e_scaler_group.test", "plan"),
					resource.TestCheckResourceAttrSet("data.e2e_scaler_group.test", "vm_image_name"),
					resource.TestCheckResourceAttrSet("data.e2e_scaler_group.test", "provision_status")),
			},
		},
	})
}

func TestAccDataSourceE2EScalerGroup_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2EScalerGroupConfig_missingID(),
				ExpectError: regexp.MustCompile(`The argument "id" is required`),
			},
			{
				Config:      testAccDataSourceE2EScalerGroupConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccDataSourceE2EScalerGroupConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccDataSourceE2EScalerGroup_NotFound(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2EScalerGroupConfig_notFound(),
				ExpectError: regexp.MustCompile(`failed to read scaler group`),
			},
		},
	})
}

// Helper functions

func testAccCheckDataSourceE2EScalerGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Scaler Group ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)

		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)
		id := rs.Primary.Attributes["id"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
		if err != nil {
			return fmt.Errorf("failed to create GoE2E client: %w", err)
		}

		group, _, err := goe2eClient.Autoscaling.GetScalerGroup(context.Background(), id)
		if err != nil {
			return err
		}

		if group == nil {
			return fmt.Errorf("Scaler Group not found in datasource: %s", id)
		}

		return nil
	}
}

// Configuration helpers

func testAccDataSourceE2EScalerGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "%s"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}

data "e2e_scaler_group" "test" {
  id         = e2e_scaler_group.test.id
  project_id = e2e_scaler_group.test.project_id
  location   = e2e_scaler_group.test.location
}
`, name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

// Error case configurations

func testAccDataSourceE2EScalerGroupConfig_missingID() string {
	return `
data "e2e_scaler_group" "test" {}
`
}

func testAccDataSourceE2EScalerGroupConfig_missingProjectID() string {
	return `
data "e2e_scaler_group" "test" {
  id       = "12345"}
`
}

func testAccDataSourceE2EScalerGroupConfig_missingLocation() string {
	return `
data "e2e_scaler_group" "test" {
  id         = "12345"}
`
}

func testAccDataSourceE2EScalerGroupConfig_notFound() string {
	return `
data "e2e_scaler_group" "test" {
  id         = "999999999"}
`
}
