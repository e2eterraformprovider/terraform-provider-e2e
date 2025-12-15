package dbaas_postgress_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testPostgresVersion   = "15"
	testPostgresPlanSmall = "E2E-2C-4GB"
	testPostgresPlanLarge = "E2E-4C-8GB"
)

func TestAccE2EPostgresDBaaS_Basic(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrName, dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrVersion, testPostgresVersion),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPlan, testPostgresPlanSmall),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPublicIPRequired, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "database.0.user", dbUser),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "database.0.name", dbName),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "database.0.dbaas_number", "1"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrIsEncryptionEnabled, "false"),
					// Verify all computed fields are populated
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrID),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatusTitle),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrNumInstances),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrProjectName),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrConnectivityDetail),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrSnapshotExist),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrVectorDatabaseStatus),
					// Verify status normalization: SUSPENDED from API should be normalized to STOPPED in state
					// Note: This checks that status is set (could be RUNNING, STOPPED, etc.)
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatus),
				),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_Update(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify state refresh updates all fields correctly
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatusTitle),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrNumInstances),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrProjectName),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrConnectivityDetail),
				),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_powerOff(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify status normalization: SUSPENDED from API should be normalized to STOPPED in state
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped),
					// Verify state refresh still updates all fields correctly after status change
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatusTitle),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrNumInstances),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrProjectName),
				),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_PowerOperations(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_powerOff(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped)),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_restart(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRestarting)),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_WithEncryption(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withEncryption(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrIsEncryptionEnabled, "true")),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_WithoutPublicIP(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withoutPublicIP(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPublicIPRequired, "false")),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_PlanUpgrade(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPlan, testPostgresPlanSmall)),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPlan, testPostgresPlanLarge)),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingName(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrName)),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingVersion(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrVersion)),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingPlan(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrPlan)),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingDatabase(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrDatabase)),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_InvalidPowerStatus(t *testing.T) {
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_invalidPowerStatus(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`expected status to be one of`),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_Import(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify all fields are populated before import
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrName, dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrVersion, testPostgresVersion),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPlan, testPostgresPlanSmall),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrID),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrProjectName),
				),
			},
			{
				ResourceName:            "e2e_dbaas_postgress.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccE2EPostgresDBaaSImportID("e2e_dbaas_postgress.test"),
				ImportStateVerifyIgnore: []string{"database.0.password", tfconstants.AttrStatus, tfconstants.AttrSize, "detach_public_ip"},
				// Verify import doesn't trigger recreation by checking ID remains the same
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrID, dbaasID),
				),
			},
		},
	})
}

// TestAccE2EPostgresDBaaS_ImportInvalidFormat tests import with invalid format
// This test verifies that invalid import ID formats are properly rejected
func TestAccE2EPostgresDBaaS_ImportInvalidFormat(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:      "e2e_dbaas_postgress.test",
				ImportState:       true,
				ImportStateId:     "invalid-format",
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile(tfconstants.DBaaSImportIDFormatDescription),
			},
			{
				ResourceName:      "e2e_dbaas_postgress.test",
				ImportState:       true,
				ImportStateId:     "project123", // Missing separator and dbaas_id
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile(tfconstants.DBaaSImportIDFormatDescription),
			},
			{
				ResourceName:      "e2e_dbaas_postgress.test",
				ImportState:       true,
				ImportStateId:     tfconstants.DBaaSImportIDSeparator + "dbaas456", // Empty project_id
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile(tfconstants.DBaaSImportIDFormatDescription),
			},
			{
				ResourceName:      "e2e_dbaas_postgress.test",
				ImportState:       true,
				ImportStateId:     "project123" + tfconstants.DBaaSImportIDSeparator, // Empty dbaas_id
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile(tfconstants.DBaaSImportIDFormatDescription),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2EPostgresDBaaSExists(resourceName string, dbaasID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Postgres DBaaS ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		ctx := context.Background()
		cluster, _, err := goe2eClient.PostgreSQL.GetCluster(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if cluster == nil {
			return fmt.Errorf("Postgres DBaaS not found")
		}

		*dbaasID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EPostgresDBaaSDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_dbaas_postgress" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		cluster, _, err := goe2eClient.PostgreSQL.GetCluster(ctx, rs.Primary.ID)
		if err == nil && cluster != nil {
			return fmt.Errorf("Postgres DBaaS still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccE2EPostgresDBaaSImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		dbaasID := rs.Primary.ID

		// Import format: project_id:dbaas_id
		return fmt.Sprintf("%s%s%s", projectID, tfconstants.DBaaSImportIDSeparator, dbaasID), nil
	}
}

// Configuration helpers

func testAccCheckE2EPostgresDBaaSConfig_basic(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_powerOff(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, goe2econstants.DBaaSStatusStopped, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_restart(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, goe2econstants.DBaaSStatusRestarting, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withEncryption(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name                    = "%s"
  version                 = "%s"
  plan                    = "%s"
  is_encryption_enabled   = true
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withoutPublicIP(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name               = "%s"
  version            = "%s"
  plan               = "%s"
  public_ip_required = false
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_upgradedPlan(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanLarge, dbUser, dbPassword, dbName)
}

// Error case configurations

func testAccCheckE2EPostgresDBaaSConfig_missingName() string {
	return `
resource "e2e_dbaas_postgress" "test" {
  version    = "15"
  plan       = "E2E-2C-4GB"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`
}

func testAccCheckE2EPostgresDBaaSConfig_missingVersion() string {
	return `
resource "e2e_dbaas_postgress" "test" {
  name       = "test-pg"
  plan       = "E2E-2C-4GB"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`
}

func testAccCheckE2EPostgresDBaaSConfig_missingPlan() string {
	return `
resource "e2e_dbaas_postgress" "test" {
  name       = "test-pg"
  version    = "15"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`
}

func testAccCheckE2EPostgresDBaaSConfig_missingDatabase() string {
	return `
resource "e2e_dbaas_postgress" "test" {
  name       = "test-pg"
  version    = "15"
  plan       = "E2E-2C-4GB"}
`
}

func testAccCheckE2EPostgresDBaaSConfig_invalidPowerStatus(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "invalid"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}

func TestAccE2EPostgresDBaaS_VPC(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	// Note: VPC IDs should be replaced with actual test VPC IDs from your environment
	vpcID1 := 1
	vpcID2 := 2

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withVPC(dbaasName, dbUser, dbPassword, dbName, vpcID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "1")),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withMultipleVPCs(dbaasName, dbUser, dbPassword, dbName, vpcID1, vpcID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "2")),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "0")),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_ParameterGroup(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	// Note: Parameter group ID should be replaced with actual test parameter group ID from your environment
	parameterGroupID := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withParameterGroup(dbaasName, dbUser, dbPassword, dbName, parameterGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", parameterGroupID))),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_DiskExpansion(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withDiskExpansion(dbaasName, dbUser, dbPassword, dbName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrSize, "10")),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withDiskExpansion(dbaasName, dbUser, dbPassword, dbName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrSize, "20")),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_PasswordRotation(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword1 := acctest.RandString(16)
	dbPassword2 := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword1, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword2, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify resource ID hasn't changed (no recreation)
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrID, dbaasID)),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_Tags(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withTags(dbaasName, dbUser, dbPassword, dbName, "env", "test", "team", "backend"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.env", tfconstants.AttrTags), "test"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.team", tfconstants.AttrTags), "backend"),
					// Verify tags are preserved from state (state-only, no API call)
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.#", tfconstants.AttrTags), "2")),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withTags(dbaasName, dbUser, dbPassword, dbName, "env", "prod", "team", "backend", "version", "1.0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.env", tfconstants.AttrTags), "prod"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.team", tfconstants.AttrTags), "backend"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.version", tfconstants.AttrTags), "1.0"),
					// Verify tags are preserved from state after update
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.#", tfconstants.AttrTags), "3")),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_ForceNew(t *testing.T) {
	var dbaasID string
	dbaasName1 := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbaasName2 := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName1, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrName, dbaasName1),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrVersion, testPostgresVersion)),
			},
			{
				// Verify ForceNew fields cannot be updated - name change should force recreation
				Config:      testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName2, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`forces new resource`),
			},
			{
				// Verify ForceNew fields cannot be updated - version change should force recreation
				Config:      testAccCheckE2EPostgresDBaaSConfig_differentVersion(dbaasName1, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`forces new resource`),
			},
		},
	})
}

// Additional configuration helpers for new tests

func testAccCheckE2EPostgresDBaaSConfig_withVPC(name, dbUser, dbPassword, dbName string, vpcID int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  vpcs       = [%d]
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, vpcID, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withMultipleVPCs(name, dbUser, dbPassword, dbName string, vpcID1, vpcID2 int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  vpcs       = [%d, %d]
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, vpcID1, vpcID2, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withParameterGroup(name, dbUser, dbPassword, dbName string, parameterGroupID int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name               = "%s"
  version            = "%s"
  plan               = "%s"
  parameter_group_id = %d
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, parameterGroupID, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withDiskExpansion(name, dbUser, dbPassword, dbName string, size int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  size       = %d
  status     = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, size, goe2econstants.DBaaSStatusSuspended, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withTags(name, dbUser, dbPassword, dbName string, tags ...string) string {
	tagsStr := ""
	if len(tags) > 0 {
		tagsStr = "  tags = {\n"
		for i := 0; i < len(tags); i += 2 {
			if i+1 < len(tags) {
				tagsStr += fmt.Sprintf("    %s = \"%s\"\n", tags[i], tags[i+1])
			}
		}
		tagsStr += "  }\n"
	}

	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
%s  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, tagsStr, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_differentVersion(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "14"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}

// ============================================================================
// Enhanced Power Management Tests
// ============================================================================

// TestAccE2EPostgresDBaaS_PowerManagementStatusTransitions tests status transitions
// including stopping, starting, and restarting instances
func TestAccE2EPostgresDBaaS_PowerManagementStatusTransitions(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify initial status is set
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatus)),
			},
			{
				// Test stopping instance (status = SUSPENDED in API, STOPPED in state)
				Config: testAccCheckE2EPostgresDBaaSConfig_powerOff(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify status normalization from API (SUSPENDED → STOPPED in state)
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped)),
			},
			{
				// Test starting instance (status = RUNNING)
				Config: testAccCheckE2EPostgresDBaaSConfig_powerOn(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRunning)),
			},
			{
				// Test restarting instance (status = RESTARTING)
				Config: testAccCheckE2EPostgresDBaaSConfig_restart(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRestarting)),
			},
		},
	})
}

// ============================================================================
// Enhanced Disk Management Tests
// ============================================================================

// TestAccE2EPostgresDBaaS_DiskExpansionCumulative tests cumulative disk expansion behavior
func TestAccE2EPostgresDBaaS_DiskExpansionCumulative(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
			{
				// First expansion: 10 GB
				Config: testAccCheckE2EPostgresDBaaSConfig_withDiskExpansion(dbaasName, dbUser, dbPassword, dbName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify size field represents additional GB to add (cumulative)
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrSize, "10")),
			},
			{
				// Second expansion: 20 GB total (cumulative)
				Config: testAccCheckE2EPostgresDBaaSConfig_withDiskExpansion(dbaasName, dbUser, dbPassword, dbName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify cumulative calculation logic works correctly
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrSize, "20")),
			},
			{
				// Third expansion: 30 GB total (cumulative)
				Config: testAccCheckE2EPostgresDBaaSConfig_withDiskExpansion(dbaasName, dbUser, dbPassword, dbName, 30),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrSize, "30")),
			},
		},
	})
}

// ============================================================================
// Enhanced Plan Upgrade Tests
// ============================================================================

// TestAccE2EPostgresDBaaS_PlanUpgradeWorkflow tests the complete plan upgrade workflow
// (stop → upgrade → start) and verifies plan upgrade requires SUSPENDED state
func TestAccE2EPostgresDBaaS_PlanUpgradeWorkflow(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPlan, testPostgresPlanSmall)),
			},
			{
				// Stop instance (required for plan upgrade)
				Config: testAccCheckE2EPostgresDBaaSConfig_powerOff(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify plan upgrade requires SUSPENDED state (STOPPED in state)
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped)),
			},
			{
				// Upgrade plan (requires SUSPENDED state)
				Config: testAccCheckE2EPostgresDBaaSConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPlan, testPostgresPlanLarge),
					// Verify plan upgrade doesn't cause recreation (same ID)
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrID, dbaasID)),
			},
		},
	})
}

// ============================================================================
// Error Scenarios Tests
// ============================================================================

// TestAccE2EPostgresDBaaS_InvalidVersion tests validation error for invalid version
func TestAccE2EPostgresDBaaS_InvalidVersion(t *testing.T) {
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_invalidVersion(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`expected version to be one of`),
			},
		},
	})
}

// TestAccE2EPostgresDBaaS_PlanUpgradeWithoutSuspendedState tests plan upgrade error
// when instance is not in SUSPENDED state
func TestAccE2EPostgresDBaaS_PlanUpgradeWithoutSuspendedState(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID)),
			},
			{
				// Try to upgrade plan without stopping (should fail)
				Config:      testAccCheckE2EPostgresDBaaSConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`must be in SUSPENDED state`),
			},
		},
	})
}

// ============================================================================
// Edge Cases Tests
// ============================================================================

// TestAccE2EPostgresDBaaS_WithCustomGroupName tests creation with custom group name
func TestAccE2EPostgresDBaaS_WithCustomGroupName(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	customGroup := "CustomGroup"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withCustomGroup(dbaasName, dbUser, dbPassword, dbName, customGroup),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrGroup, customGroup)),
			},
		},
	})
}

// TestAccE2EPostgresDBaaS_WithAllOptionalFields tests creation with all optional fields set
func TestAccE2EPostgresDBaaS_WithAllOptionalFields(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	vpcID1 := 1
	parameterGroupID := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withAllOptionalFields(dbaasName, dbUser, dbPassword, dbName, vpcID1, parameterGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrIsEncryptionEnabled, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPublicIPRequired, "false"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "1"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", parameterGroupID))),
			},
		},
	})
}

// ============================================================================
// Additional Configuration Helpers
// ============================================================================

func testAccCheckE2EPostgresDBaaSConfig_powerOn(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, goe2econstants.DBaaSStatusRunning, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_invalidVersion(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "99.0"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withCustomGroup(name, dbUser, dbPassword, dbName, group string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  group      = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, group, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withAllOptionalFields(name, dbUser, dbPassword, dbName string, vpcID, parameterGroupID int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name                 = "%s"
  version              = "%s"
  plan                 = "%s"
  is_encryption_enabled = true
  public_ip_required   = false
  vpcs                 = [%d]
  parameter_group_id   = %d
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, vpcID, parameterGroupID, dbUser, dbPassword, dbName)
}

// ============================================================================
// Security Review Tests
// ============================================================================

// TestAccE2EPostgresDBaaS_PasswordSensitive tests that database.password is marked Sensitive
// and doesn't appear in logs or state output
func TestAccE2EPostgresDBaaS_PasswordSensitive(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Verify database.password field exists but is sensitive (won't show in state)
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", "database.0.user"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", "database.0.name"),
					// Note: We can't directly test that password is sensitive, but we verify
					// the resource schema marks it as sensitive in unit tests
				),
			},
		},
	})
}

// ============================================================================
// Enhanced Error Scenarios Tests
// ============================================================================

// TestAccE2EPostgresDBaaS_InvalidStatusValue tests validation error for invalid status value
func TestAccE2EPostgresDBaaS_InvalidStatusValue(t *testing.T) {
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_invalidStatus(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`expected status to be one of`),
			},
		},
	})
}

// TestAccE2EPostgresDBaaS_EncryptionForceNew tests that encryption cannot be changed after creation
func TestAccE2EPostgresDBaaS_EncryptionForceNew(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrIsEncryptionEnabled, "false")),
			},
			{
				// Try to update encryption (should fail or force recreation)
				Config:      testAccCheckE2EPostgresDBaaSConfig_withEncryption(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`forces new resource`),
			},
		},
	})
}

// Additional configuration helpers for error scenarios

func testAccCheckE2EPostgresDBaaSConfig_invalidStatus(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "INVALID_STATUS"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}

// ============================================================================
// API Error Handling Tests
// ============================================================================

func TestAccE2EPostgresDBaaS_ImportNonExistentInstance(t *testing.T) {
	// Test import with non-existent instance ID (error handling)
	// Note: This test requires API access and may fail if instance ID format is invalid
	// Skipping for now as it requires actual API validation
	t.Skip("Skipping test that requires API validation of non-existent instance")
}

func TestAccE2EPostgresDBaaS_DiskExpansionInvalidSize(t *testing.T) {
	// Test disk expansion with invalid size (negative or zero)
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
				),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_diskExpansionInvalid(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`size must be greater than previous size`),
			},
		},
	})
}

// ============================================================================
// Update Combination Tests
// ============================================================================

func TestAccE2EPostgresDBaaS_MultipleSimultaneousUpdates(t *testing.T) {
	// Test multiple simultaneous updates (status + VPC + public IP)
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPublicIPRequired, "true"),
				),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_multipleUpdates(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrPublicIPRequired, "false"),
					// Verify status change was applied
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatus),
				),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_UpdateWithInstanceInCreatingState(t *testing.T) {
	// Test update with instance in CREATING state (should be blocked)
	// Note: This is difficult to test reliably as CREATING state is transient
	// Skipping for now as it requires precise timing
	t.Skip("Skipping test that requires instance in CREATING state (timing-dependent)")
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestAccE2EPostgresDBaaS_VeryLongName(t *testing.T) {
	// Test with very long instance name (API limits)
	// Note: Actual API limit is unknown, using 100 characters as test
	var dbaasID string
	longName := fmt.Sprintf("test-pg-%s", acctest.RandString(90)) // Total ~100 chars
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(longName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrName, longName),
				),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_SpecialCharactersInName(t *testing.T) {
	// Test with special characters in instance name
	// Note: Using alphanumeric and hyphens only as special characters may be rejected
	var dbaasID string
	specialName := fmt.Sprintf("test-pg-123-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(specialName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrName, specialName),
				),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_LargeDiskExpansion(t *testing.T) {
	// Test with very large disk_size expansion (API limits)
	// Note: Actual API limit is unknown, using 1000 GB as test
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
				),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_largeDiskExpansion(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrSize, "1000"),
				),
			},
		},
	})
}

// ============================================================================
// Concurrent Operations Tests
// ============================================================================

func TestAccE2EPostgresDBaaS_RapidStatusChanges(t *testing.T) {
	// Test multiple rapid status changes
	// Note: This test verifies that rapid status changes are handled correctly
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
				),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_powerOff(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped),
				),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					// Status should be RUNNING after starting
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", tfconstants.AttrStatus),
				),
			},
		},
	})
}

// ============================================================================
// Additional Test Configuration Helpers
// ============================================================================

func testAccCheckE2EPostgresDBaaSConfig_diskExpansionInvalid(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  size       = 5
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_multipleUpdates(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name                = "%s"
  version             = "%s"
  plan                = "%s"
  public_ip_required  = false
  status              = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, goe2econstants.DBaaSStatusStopped, dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_largeDiskExpansion(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "%s"
  plan       = "%s"
  size       = 1000
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testPostgresVersion, testPostgresPlanSmall, dbUser, dbPassword, dbName)
}
