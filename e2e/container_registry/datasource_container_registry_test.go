package container_registry_test

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

func TestAccDataSourceE2EContainerRegistry_Basic(t *testing.T) {
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EContainerRegistryExists("data.e2e_container_registry.test"),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", "id", "e2e_container_registry.test", "id"),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", "project_name", "e2e_container_registry.test", "project_name"),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", "setup_status", "e2e_container_registry.test", "setup_status"),
					resource.TestCheckResourceAttrSet("data.e2e_container_registry.test", "severity"),
					resource.TestCheckResourceAttrSet("data.e2e_container_registry.test", "prevent_vul"),
				),
			},
		},
	})
}

func TestAccDataSourceE2EContainerRegistry_WithSeverity(t *testing.T) {
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EContainerRegistryConfig_withSeverity(projectName, "critical"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EContainerRegistryExists("data.e2e_container_registry.test"),
					resource.TestCheckResourceAttr("data.e2e_container_registry.test", "severity", "critical"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "critical"),
				),
			},
		},
	})
}

func TestAccDataSourceE2EContainerRegistry_WithPreventVul(t *testing.T) {
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EContainerRegistryConfig_withPreventVul(projectName, true, "high"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceE2EContainerRegistryExists("data.e2e_container_registry.test"),
					resource.TestCheckResourceAttr("data.e2e_container_registry.test", "prevent_vul", "true"),
					resource.TestCheckResourceAttr("data.e2e_container_registry.test", "severity", "high"),
				),
			},
		},
	})
}

func TestAccDataSourceE2EContainerRegistry_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2EContainerRegistryConfig_missingID(),
				ExpectError: regexp.MustCompile(`The argument "id" is required`),
			},
			{
				Config:      testAccDataSourceE2EContainerRegistryConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccDataSourceE2EContainerRegistryConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccDataSourceE2EContainerRegistry_NotFound(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2EContainerRegistryConfig_notFound(),
				ExpectError: regexp.MustCompile(`container registry with ID .* not found`),
			},
		},
	})
}

// Helper functions

func testAccCheckDataSourceE2EContainerRegistryExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Container Registry ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]
		id := rs.Primary.Attributes["id"]

		registries, err := client.GetContainerRegistryProjects(projectID, location)
		if err != nil {
			return err
		}

		found := false
		for _, registry := range registries {
			if fmt.Sprintf("%d", registry.ID) == id {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Container Registry not found in datasource: %s", id)
		}

		return nil
	}
}

// Configuration helpers

func testAccDataSourceE2EContainerRegistryConfig_basic(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id   = "%s"
  location     = "%s"
  project_name = "%s"
}

data "e2e_container_registry" "test" {
  id         = e2e_container_registry.test.id
  project_id = e2e_container_registry.test.project_id
  location   = e2e_container_registry.test.location
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), projectName)
}

func testAccDataSourceE2EContainerRegistryConfig_withSeverity(projectName, severity string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id   = "%s"
  location     = "%s"
  project_name = "%s"
  severity     = "%s"
}

data "e2e_container_registry" "test" {
  id         = e2e_container_registry.test.id
  project_id = e2e_container_registry.test.project_id
  location   = e2e_container_registry.test.location
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), projectName, severity)
}

func testAccDataSourceE2EContainerRegistryConfig_withPreventVul(projectName string, preventVul bool, severity string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id   = "%s"
  location     = "%s"
  project_name = "%s"
  prevent_vul  = %t
  severity     = "%s"
}

data "e2e_container_registry" "test" {
  id         = e2e_container_registry.test.id
  project_id = e2e_container_registry.test.project_id
  location   = e2e_container_registry.test.location
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), projectName, preventVul, severity)
}

// Error case configurations

func testAccDataSourceE2EContainerRegistryConfig_missingID() string {
	return fmt.Sprintf(`
data "e2e_container_registry" "test" {
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccDataSourceE2EContainerRegistryConfig_missingProjectID() string {
	return fmt.Sprintf(`
data "e2e_container_registry" "test" {
  id       = "12345"
  location = "%s"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccDataSourceE2EContainerRegistryConfig_missingLocation() string {
	return fmt.Sprintf(`
data "e2e_container_registry" "test" {
  id         = "12345"
  project_id = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccDataSourceE2EContainerRegistryConfig_notFound() string {
	return fmt.Sprintf(`
data "e2e_container_registry" "test" {
  id         = "999999999"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
