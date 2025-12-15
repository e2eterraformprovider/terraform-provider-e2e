package container_registry_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/container_registry"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EContainerRegistry_Basic(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrProjectName, projectName),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "false"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityLow),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrSetupStatus)),
			},
		},
	})
}

func TestAccE2EContainerRegistry_Update(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "false"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, "low")),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_updated(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "true"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, "high")),
			},
		},
	})
}

func TestAccE2EContainerRegistry_WithSeverityLevels(t *testing.T) {
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

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
							resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, severity)),
					},
				},
			})
		})
	}
}

func TestAccE2EContainerRegistry_WithPreventVulnerability(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName, true, "critical"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "true"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, "critical")),
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
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID)),
			},
			{
				ResourceName:      container_registry.TestResourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccE2EContainerRegistry_UpdateSeverityOnly(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, "low"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, "low")),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, "critical"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, "critical")),
			},
		},
	})
}

func TestAccE2EContainerRegistry_Tags(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withTags(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrTags+".Environment", "test"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrTags+".ManagedBy", "terraform")),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withTagsUpdated(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrTags+".Environment", "production"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrTags+".ManagedBy", "terraform"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrTags+".Owner", "devops")),
			},
		},
	})
}

func TestAccE2EContainerRegistry_SeverityValidation(t *testing.T) {
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

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
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID1),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrProjectName, projectName1)),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID2),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrProjectName, projectName2),
					testAccCheckE2EContainerRegistryRecreated(&registryID1, &registryID2)),
			},
		},
	})
}

func TestAccE2EContainerRegistry_DeprecatedSetupStatus(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrSetupStatus),
					// Verify both status and setup_status have the same value (backward compatibility)
					resource.TestCheckResourceAttrPair(container_registry.TestResourceType+".test", tfconstants.AttrStatus, container_registry.TestResourceType+".test", tfconstants.AttrSetupStatus)),
			},
		},
	})
}

func TestAccE2EContainerRegistry_V2Compatibility(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				// Test that location (deprecated) still works for backward compatibility
				Config: testAccCheckE2EContainerRegistryConfig_withLocation(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrProjectName, projectName)),
			},
			{
				// Test that region (new) works
				Config: testAccCheckE2EContainerRegistryConfig_withRegion(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrProjectName, projectName)),
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

// ============================================================================
// DATA SOURCE ACCEPTANCE TESTS
// ============================================================================

func TestAccE2EContainerRegistryDataSource_Basic(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryDataSourceConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", tfconstants.AttrID, container_registry.TestResourceType+".test", tfconstants.AttrID),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", tfconstants.AttrProjectName, container_registry.TestResourceType+".test", tfconstants.AttrProjectName),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", tfconstants.AttrStatus, container_registry.TestResourceType+".test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", tfconstants.AttrSetupStatus, container_registry.TestResourceType+".test", tfconstants.AttrSetupStatus),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", tfconstants.AttrSeverity, container_registry.TestResourceType+".test", tfconstants.AttrSeverity),
					resource.TestCheckResourceAttrPair("data.e2e_container_registry.test", tfconstants.AttrPreventVulnerabilities, container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistryDataSource_SecurityFields(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryDataSourceConfig_security(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr("data.e2e_container_registry.test", tfconstants.AttrPreventVulnerabilities, "true"),
					resource.TestCheckResourceAttr("data.e2e_container_registry.test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityCritical),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistryDataSource_NonExistent(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EContainerRegistryDataSourceConfig_nonExistent(),
				ExpectError: regexp.MustCompile(`container registry.*not found|error.*reading`),
			},
		},
	})
}

func TestAccE2EContainerRegistryDataSource_InvalidID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EContainerRegistryDataSourceConfig_invalidID(),
				ExpectError: regexp.MustCompile(`invalid.*ID|error.*parsing`),
			},
		},
	})
}

// Data source configuration helpers

func testAccCheckE2EContainerRegistryDataSourceConfig_basic(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"
  prevent_vul  = false
  severity     = "low"
}

data "e2e_container_registry" "test" {
  id = e2e_container_registry.test.id
}
`, projectName)
}

func testAccCheckE2EContainerRegistryDataSourceConfig_security(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"
  prevent_vul  = true
  severity     = "critical"
}

data "e2e_container_registry" "test" {
  id = e2e_container_registry.test.id
}
`, projectName)
}

func testAccCheckE2EContainerRegistryDataSourceConfig_nonExistent() string {
	return `
data "e2e_container_registry" "test" {
  id = "999999999"
}
`
}

func testAccCheckE2EContainerRegistryDataSourceConfig_invalidID() string {
	return `
data "e2e_container_registry" "test" {
  id = "invalid-id-format"
}
`
}

// ============================================================================
// ADDITIONAL EDGE CASE ACCEPTANCE TESTS
// ============================================================================

func TestAccE2EContainerRegistry_ComputedFields(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_basic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					// Verify all computed fields are set
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrDomainName),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrProjectSize),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrStorageLimit),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrIsPublic),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrCreatedAt),
					resource.TestCheckResourceAttrSet(container_registry.TestResourceType+".test", tfconstants.AttrUpdatedAt),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_SeverityNone(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, goe2econstants.ContainerRegistrySeverityNone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityNone),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_AllSeverityTransitions(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, goe2econstants.ContainerRegistrySeverityLow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityLow),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, goe2econstants.ContainerRegistrySeverityMedium),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityMedium),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, goe2econstants.ContainerRegistrySeverityHigh),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityHigh),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, goe2econstants.ContainerRegistrySeverityCritical),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityCritical),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withSeverity(projectName, goe2econstants.ContainerRegistrySeverityNone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityNone),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_RemoveTags(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withTags(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrTags+".Environment", "test"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrTags+".ManagedBy", "terraform"),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withEmptyTags(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrTags+".%", "0"),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_PreventVulToggle(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName, false, goe2econstants.ContainerRegistrySeverityLow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "false"),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName, true, goe2econstants.ContainerRegistrySeverityLow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "true"),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName, false, goe2econstants.ContainerRegistrySeverityLow),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "false"),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_MultipleRegistries(t *testing.T) {
	var registryID1, registryID2 string
	projectName1 := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))
	projectName2 := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_multiple(projectName1, projectName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test1", &registryID1),
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test2", &registryID2),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test1", tfconstants.AttrProjectName, projectName1),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test2", tfconstants.AttrProjectName, projectName2),
				),
			},
		},
	})
}

func TestAccE2EContainerRegistry_ImportWithTags(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withTags(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
				),
			},
			{
				ResourceName:      container_registry.TestResourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				// Tags are state-only, so they won't be preserved on import
				ImportStateVerifyIgnore: []string{tfconstants.AttrTags},
			},
		},
	})
}

func TestAccE2EContainerRegistry_UpdatePreventVulOnly(t *testing.T) {
	var registryID string
	projectName := fmt.Sprintf("%s%s", container_registry.TestResourceNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName, false, goe2econstants.ContainerRegistrySeverityMedium),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "false"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityMedium),
				),
			},
			{
				Config: testAccCheckE2EContainerRegistryConfig_withPreventVul(projectName, true, goe2econstants.ContainerRegistrySeverityMedium),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EContainerRegistryExists(container_registry.TestResourceType+".test", &registryID),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrPreventVulnerabilities, "true"),
					resource.TestCheckResourceAttr(container_registry.TestResourceType+".test", tfconstants.AttrSeverity, goe2econstants.ContainerRegistrySeverityMedium),
				),
			},
		},
	})
}

// Additional configuration helpers

func testAccCheckE2EContainerRegistryConfig_withEmptyTags(projectName string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test" {
  project_name = "%s"

  tags = {}
}
`, projectName)
}

func testAccCheckE2EContainerRegistryConfig_multiple(projectName1, projectName2 string) string {
	return fmt.Sprintf(`
resource "e2e_container_registry" "test1" {
  project_name = "%s"
  prevent_vul  = false
  severity     = "low"
}

resource "e2e_container_registry" "test2" {
  project_name = "%s"
  prevent_vul  = true
  severity     = "high"
}
`, projectName1, projectName2)
}
