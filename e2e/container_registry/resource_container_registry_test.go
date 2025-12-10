package container_registry_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EContainerRegistry_Basic(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "project_name", projectName),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "prevent_vul", "false"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "low"),
					resource.TestCheckResourceAttrSet("e2e_container_registry.test", "setup_status")),
			},
		},
	})
}

func TestAccE2EContainerRegistry_Update(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "prevent_vul", "false"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "low")),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_updated(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "prevent_vul", "true"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "high")),
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
				ProviderFactories: acceptance.TestAccProviderFactories,
				CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, severity),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", severity)),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName, true, "critical"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "prevent_vul", "true"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "critical")),
			},
		},
	})
}

func TestAccE2EContainerRegistry_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID)),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, "low"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "low")),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, "critical"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "severity", "critical")),
			},
		},
	})
}

func TestAccE2EContainerRegistry_Tags(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withTags(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "tags.ManagedBy", "terraform")),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withTagsUpdated(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "tags.ManagedBy", "terraform"),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "tags.Owner", "devops")),
			},
		},
	})
}

func TestAccE2EContainerRegistry_SeverityValidation(t *testing.T) {
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EContainerRegistryConfig_invalidSeverity(projectName),
				ExpectError: regexp.MustCompile(`expected severity to be one of`),
			},
		},
	})
}

func TestAccE2EContainerRegistry_ForceNew(t *testing.T) {
	var registryID1, registryID2 string
	projectName1 := fmt.Sprintf("test-cr-%s", acctest.RandString(10))
	projectName2 := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID1),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "project_name", projectName1)),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID2),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "project_name", projectName2),
					testAccCheckE2EContainerRegistryRecreated(&registryID1, &registryID2)),
			},
		},
	})
}

func TestAccE2EContainerRegistry_DeprecatedSetupStatus(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttrSet("e2e_container_registry.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_container_registry.test", "setup_status"),
					// Verify both status and setup_status have the same value (backward compatibility)
					resource.TestCheckResourceAttrPair("e2e_container_registry.test", "status", "e2e_container_registry.test", "setup_status")),
			},
		},
	})
}

func TestAccE2EContainerRegistry_V2Compatibility(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("test-cr-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				// Test that location (deprecated) still works for backward compatibility
				Config: testAccCheckE2EContainerRegistryConfig_withLocation(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "project_name", projectName)),
			},
			{
				// Test that region (new) works
				Config: testAccCheckE2EContainerRegistryConfig_withRegion(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists("e2e_container_registry.test", &registryID),
					resource.TestCheckResourceAttr("e2e_container_registry.test", "project_name", projectName)),
			},
			{
				// Test that using both location and region causes conflict
				Config:      testAccCheckE2EContainerRegistryConfig_locationAndRegion(projectName),
				ExpectError: regexp.MustCompile(`.*conflicts with.*`),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
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

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client := cfg.Goe2eClient()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("invalid registry ID: %w", err)
		}

		ctx := context.Background()
		registry, _, err := client.ContainerRegistry.GetContainerRegistry(ctx, id)
		if err != nil {
			return err
		}

		if registry == nil {
			return fmt.Errorf("Container Registry not found: %s", rs.Primary.ID)
		}

		*registryID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EContainerRegistryDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client := cfg.Goe2eClient()

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_container_registry" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			// If we can't parse the ID, skip (it was likely already destroyed)
			continue
		}

		registry, _, err := client.ContainerRegistry.GetContainerRegistry(ctx, id)
		if err != nil {
			// If we get an error, the registry is likely already deleted
			continue
		}

		if registry != nil {
			return fmt.Errorf("Container Registry still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckE2EContainerRegistryRecreated(oldID, newID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *oldID == *newID {
			return fmt.Errorf("Expected container registry to be recreated, but IDs are the same: %s", *oldID)
		}
		return nil
	}
}

// Configuration helpers

func testAccCheckE2EContainerRegistryConfig_basic(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {  project_name = "%s"
  prevent_vul  = false
  severity     = "low"
}
`, projectName)
}

func testAccCheckE2EContainerRegistryConfig_updated(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {  project_name = "%s"
  prevent_vul  = true
  severity     = "high"
}
`, projectName)
}

func testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, severity string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {  project_name = "%s"
  severity     = "%s"
}
`, projectName, severity)
}

func testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName string, preventVul bool, severity string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {  project_name = "%s"
  prevent_vul  = %t
  severity     = "%s"
}
`, projectName, preventVul, severity)
}

func testAccCheckE2EContainerRegistryConfig_withTags(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"

  tags = {
    Environment = "test"
    ManagedBy   = "terraform"
  }
}
`, projectName)
}

func testAccCheckE2EContainerRegistryConfig_withTagsUpdated(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"

  tags = {
    Environment = "production"
    ManagedBy   = "terraform"
    Owner       = "devops"
  }
}
`, projectName)
}

// Error case configurations

func testAccCheckE2EContainerRegistryConfig_missingProjectID() string {
	return `
resource "e2e_container_registry" "test" {  project_name = "test-project"
}
`
}

func testAccCheckE2EContainerRegistryConfig_missingLocation() string {
	return `
resource "e2e_container_registry" "test" {  project_name = "test-project"
}
`
}

func testAccCheckE2EContainerRegistryConfig_missingProjectName() string {
	return `
resource "e2e_container_registry" "test" {}
`
}

func testAccCheckE2EContainerRegistryConfig_invalidSeverity(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"
  severity     = "invalid-severity"
}
`, projectName)
}

func testAccCheckE2EContainerRegistryConfig_withLocation(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"
  location     = "us-east-1"
}
`, projectName)
}

func testAccCheckE2EContainerRegistryConfig_withRegion(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"
  region       = "us-east-1"
}
`, projectName)
}

func testAccCheckE2EContainerRegistryConfig_locationAndRegion(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"
  location     = "us-east-1"
  region       = "us-west-1"
}
`, projectName)
}
