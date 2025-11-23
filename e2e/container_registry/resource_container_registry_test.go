package container_registry_test

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

func TestAccE2EContainerRegistry_Basic(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "project_name", projectName),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "prevent_vul", "false"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "low"),
					resource.TestCheckResourceAttrSet("e2e_container_registry.test", "setup_status"),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_Update(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "prevent_vul", "false"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "low"),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_updated(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "prevent_vul", "true"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "high"),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_WithSeverityLevels(t *testing.T) {
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	severityLevels := []string{"low", "medium", "high", "critical"}

	for _, severity := range severityLevels {
		t.Run(fmt.Sprintf("severity_%s", severity), func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: testAccProviderFactories,
				CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, severity),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", severity),
						),
					},
				},
			})
		})
	}
}

func TestAccE2EContainerRegistry_WithPreventVulnerability(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName, true, "critical"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "prevent_vul", "true"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "critical"),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EContainerRegistryConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EContainerRegistryConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
			{
				Config:      testAccCheckE2EContainerRegistryConfig_missingProjectName(),
				ExpectError: regexp.MustCompile(`The argument "project_name" is required`),
			},
		},
	})
}

func TestAccE2EContainerRegistry_Import(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
				),
			},
			{
				ResourceName:      "e2e_container_registry.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccE2EContainerRegistry_UpdateSeverityOnly(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, "low"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "low"),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, "critical"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "critical"),
				),
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

func testAccCheckE2EContainerRegistryExists(resourceName string, registryID *string) resource.TestCheckFunc {
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

		registries, err := client.GetContainerRegistryProjects(projectID, location)
		if err != nil {
			return err
		}

		found := false
		for _, registry := range registries {
			if fmt.Sprintf("%d", registry.ID) == rs.Primary.ID {
				found = true
				*registryID = rs.Primary.ID
				break
			}
		}

		if !found {
			return fmt.Errorf("Container Registry not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckE2EContainerRegistryDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_container_registry" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		registries, err := client.GetContainerRegistryProjects(projectID, location)
		if err != nil {
			continue
		}

		for _, registry := range registries {
			if fmt.Sprintf("%d", registry.ID) == rs.Primary.ID {
				return fmt.Errorf("Container Registry still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2EContainerRegistryConfig_basic(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id   = "%s"
  location     = "%s"
  project_name = "%s"
  prevent_vul  = false
  severity     = "low"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), projectName)
}

func testAccCheckE2EContainerRegistryConfig_updated(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id   = "%s"
  location     = "%s"
  project_name = "%s"
  prevent_vul  = true
  severity     = "high"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), projectName)
}

func testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, severity string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id   = "%s"
  location     = "%s"
  project_name = "%s"
  severity     = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), projectName, severity)
}

func testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName string, preventVul bool, severity string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id   = "%s"
  location     = "%s"
  project_name = "%s"
  prevent_vul  = %t
  severity     = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), projectName, preventVul, severity)
}

// Error case configurations

func testAccCheckE2EContainerRegistryConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  location     = "%s"
  project_name = "test-project"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EContainerRegistryConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id   = "%s"
  project_name = "test-project"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccCheckE2EContainerRegistryConfig_missingProjectName() string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
