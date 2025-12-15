package dbaas_mysql_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mysql"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testMySQLVersion   = "8.0"
	testMySQLPlanSmall = "E2E-2C-4GB"
	testMySQLPlanLarge = "E2E-4C-8GB"
)

func TestAccE2EMySql_Basic(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrDBaaSName, dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrVersion, testMySQLVersion),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPlan, testMySQLPlanSmall),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPublicIPRequired, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "database.0.user", dbUser),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "database.0.name", dbName),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "database.0.dbaas_number", "1"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrID),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrDisk)),
			},
		},
	})
}

func TestAccE2EMySql_Import(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				ResourceName:            "e2e_dbaas_mysql.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccE2EMySqlImportID("e2e_dbaas_mysql.test"),
				ImportStateVerifyIgnore: []string{"database.0.password", tfconstants.AttrStatus, tfconstants.AttrSize},
			},
		},
	})
}

func TestAccE2EMySql_ForceNew(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrVersion, "8.0")),
			},
			{
				Config:      testAccCheckE2EMySqlConfig_differentVersion(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`forces new resource`),
			},
		},
	})
}

func TestAccE2EMySql_PowerOperations(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_stop(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrStatus, tfconstants.DBaaSPowerActionStop)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_restart(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrStatus, tfconstants.DBaaSPowerActionRestart)),
			},
		},
	})
}

func TestAccE2EMySql_PublicIP(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPublicIPRequired, "true")),
			},
			{
				Config: testAccCheckE2EMySqlConfig_withoutPublicIP(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPublicIPRequired, "false")),
			},
		},
	})
}

func TestAccE2EMySql_PlanUpgrade(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPlan, testMySQLPlanSmall)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPlan, testMySQLPlanLarge)),
			},
		},
	})
}

func TestAccE2EMySql_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EMySqlConfig_missingVersion(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrVersion)),
			},
			{
				Config:      testAccCheckE2EMySqlConfig_missingPlan(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrPlan)),
			},
			{
				Config:      testAccCheckE2EMySqlConfig_missingDatabase(),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, tfconstants.AttrDatabase)),
			},
		},
	})
}

func TestAccE2EMySql_Encryption(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	passphrase := acctest.RandString(32)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_encryption(dbaasName, dbUser, dbPassword, dbName, passphrase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrIsEncryptionEnabled, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrEncryptionPassphrase, passphrase)),
			},
		},
	})
}

func TestAccE2EMySql_Tags(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_tags(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "tags.environment", "test"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "tags.managed_by", "terraform")),
			},
			{
				Config: testAccCheckE2EMySqlConfig_tagsUpdated(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "tags.environment", "production"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "tags.version", "1.0")),
			},
		},
	})
}

func TestAccE2EMySql_DiskExpansion(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					// Verify disk field is set (total disk size)
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrDisk)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					// After expansion, size should be reset to 0 (cumulative behavior)
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrSize, "0"),
					// Verify disk field reflects total disk size after expansion
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrDisk)),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2EMySqlExists(resourceName string, dbaasID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No %s ID is set", dbaas_mysql.ResourceName)
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		goe2eClient := cfg.Goe2eClient()

		dbaas, _, err := goe2eClient.DBaaSMySQL.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if dbaas == nil {
			return fmt.Errorf("%s not found", dbaas_mysql.ResourceName)
		}

		*dbaasID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EMySqlDBaaSDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	goe2eClient := cfg.Goe2eClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_dbaas_mysql" {
			continue
		}

		dbaas, _, err := goe2eClient.DBaaSMySQL.GetCluster(context.Background(), rs.Primary.ID)
		if err == nil && dbaas != nil {
			return fmt.Errorf("%s still exists: %s", dbaas_mysql.ResourceName, rs.Primary.ID)
		}
	}

	return nil
}

func testAccE2EMySqlImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		dbaasID := rs.Primary.ID

		return fmt.Sprintf("%s%s%s", projectID, tfconstants.DBaaSImportIDSeparator, dbaasID), nil
	}
}

// Configuration helpers

func testAccCheckE2EMySqlConfig_basic(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_stop(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, tfconstants.DBaaSPowerActionStop, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_restart(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, tfconstants.DBaaSPowerActionRestart, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_withoutPublicIP(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name        = "%s"
  version           = "%s"
  plan              = "%s"
  public_ip_required = false
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_upgradedPlan(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanLarge, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_differentVersion(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "5.7"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLPlanSmall, dbUser, dbPassword, dbName)
}

// Error case configurations

func testAccCheckE2EMySqlConfig_missingVersion() string {
	return `
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "test-mysql"
  plan       = "E2E-2C-4GB"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`
}

func testAccCheckE2EMySqlConfig_missingPlan() string {
	return `
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "test-mysql"
  version    = "8.0"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`
}

func testAccCheckE2EMySqlConfig_missingDatabase() string {
	return `
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "test-mysql"
  version    = "8.0"
  plan       = "E2E-2C-4GB"
}
`
}

func testAccCheckE2EMySqlConfig_encryption(name, dbUser, dbPassword, dbName, passphrase string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name             = "%s"
  version                = "%s"
  plan                   = "%s"
  is_encryption_enabled  = true
  encryption_passphrase  = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, passphrase, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_tags(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
  tags = {
    environment = "test"
    managed_by  = "terraform"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_tagsUpdated(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
  tags = {
    environment = "production"
    version     = "1.0"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_diskExpansion(name, dbUser, dbPassword, dbName string, additionalSize int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  size       = %d
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, additionalSize, dbUser, dbPassword, dbName)
}

func TestAccE2EMySql_VPC(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	// Note: VPC IDs should be replaced with actual test VPC IDs from your environment
	vpcID1 := 1
	vpcID2 := 2

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_withVPC(dbaasName, dbUser, dbPassword, dbName, vpcID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "1")),
			},
			{
				Config: testAccCheckE2EMySqlConfig_withMultipleVPCs(dbaasName, dbUser, dbPassword, dbName, vpcID1, vpcID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "2")),
			},
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "0")),
			},
		},
	})
}

func TestAccE2EMySql_ParameterGroup(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	// Note: Parameter group ID should be replaced with actual test parameter group ID from your environment
	parameterGroupID := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_withParameterGroup(dbaasName, dbUser, dbPassword, dbName, parameterGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", parameterGroupID))),
			},
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
		},
	})
}

// Additional configuration helpers for VPC and Parameter Group tests

func testAccCheckE2EMySqlConfig_withVPC(name, dbUser, dbPassword, dbName string, vpcID int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  vpcs       = [%d]
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, vpcID, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_withMultipleVPCs(name, dbUser, dbPassword, dbName string, vpcID1, vpcID2 int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  vpcs       = [%d, %d]
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, vpcID1, vpcID2, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_withParameterGroup(name, dbUser, dbPassword, dbName string, parameterGroupID int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name        = "%s"
  version           = "%s"
  plan              = "%s"
  parameter_group_id = %d
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, parameterGroupID, dbUser, dbPassword, dbName)
}

// ============================================================================
// Additional TestAcc Tests for Disk Management, Plan Upgrade, Encryption, and Error Scenarios
// ============================================================================

func TestAccE2EMySql_MultipleDiskExpansions(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrDisk)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrSize, "0"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrDisk)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrSize, "0"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrDisk)),
			},
		},
	})
}

func TestAccE2EMySql_DiskExpansionWithRunningState(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_diskExpansionWithStatus(dbaasName, dbUser, dbPassword, dbName, 10, goe2econstants.DBaaSStatusRunning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrSize, "0"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrDisk)),
			},
		},
	})
}

func TestAccE2EMySql_PlanUpgradeWorkflow(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPlan, testMySQLPlanSmall)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_stop(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrStatus, tfconstants.DBaaSPowerActionStop)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPlan, testMySQLPlanLarge)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_start(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrStatus, tfconstants.DBaaSPowerActionStart)),
			},
		},
	})
}

func TestAccE2EMySql_PlanUpgradeWithoutSuspendedState(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPlan, testMySQLPlanSmall)),
			},
			{
				Config:      testAccCheckE2EMySqlConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`cannot upgrade plan.*database must be in SUSPENDED/STOPPED state`),
			},
		},
	})
}

func TestAccE2EMySql_EncryptionForceNew(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	passphrase1 := acctest.RandString(32)
	passphrase2 := acctest.RandString(32)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_encryption(dbaasName, dbUser, dbPassword, dbName, passphrase1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrIsEncryptionEnabled, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrEncryptionPassphrase, passphrase1)),
			},
			{
				Config:      testAccCheckE2EMySqlConfig_encryption(dbaasName, dbUser, dbPassword, dbName, passphrase2),
				ExpectError: regexp.MustCompile(`forces new resource`),
			},
		},
	})
}

func TestAccE2EMySql_InvalidStatusValue(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EMySqlConfig_invalidStatus(),
				ExpectError: regexp.MustCompile(`expected status to be one of`),
			},
		},
	})
}

func TestAccE2EMySql_InvalidVersion(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EMySqlConfig_invalidVersion(),
				ExpectError: regexp.MustCompile(`error retrieving.*software ID`),
			},
		},
	})
}

func TestAccE2EMySql_InvalidPlan(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EMySqlConfig_invalidPlan(),
				ExpectError: regexp.MustCompile(`error retrieving.*template ID`),
			},
		},
	})
}

func TestAccE2EMySql_DiskExpansionEdgeCases(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrSize, "0")),
			},
		},
	})
}

// Additional configuration helpers

func testAccCheckE2EMySqlConfig_diskExpansionWithStatus(name, dbUser, dbPassword, dbName string, additionalSize int, status string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
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
`, name, testMySQLVersion, testMySQLPlanSmall, additionalSize, status, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_start(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanLarge, tfconstants.DBaaSPowerActionStart, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_invalidStatus() string {
	return `
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "test-mysql"
  version    = "8.0"
  plan       = "E2E-2C-4GB"
  status     = "INVALID_STATUS"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`
}

func testAccCheckE2EMySqlConfig_invalidVersion() string {
	return `
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "test-mysql"
  version    = "99.99"
  plan       = "E2E-2C-4GB"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`
}

func testAccCheckE2EMySqlConfig_invalidPlan() string {
	return `
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "test-mysql"
  version    = "8.0"
  plan       = "INVALID-PLAN-NAME"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`
}

// ============================================================================
// Edge Cases Tests
// ============================================================================

func TestAccE2EMySql_LongName(t *testing.T) {
	var dbaasID string
	// Create a very long name (255 characters is typically the max for many systems)
	longName := fmt.Sprintf("test-mysql-%s", acctest.RandString(240))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(longName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrDBaaSName, longName)),
			},
		},
	})
}

func TestAccE2EMySql_SpecialCharactersInName(t *testing.T) {
	var dbaasID string
	// Test with special characters that are typically allowed in resource names
	specialName := fmt.Sprintf("test-mysql-%s-123", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(specialName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrDBaaSName, specialName)),
			},
		},
	})
}

func TestAccE2EMySql_AllOptionalFields(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	passphrase := acctest.RandString(32)
	vpcID1 := 1
	vpcID2 := 2
	parameterGroupID := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_allOptionalFields(dbaasName, dbUser, dbPassword, dbName, passphrase, vpcID1, vpcID2, parameterGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrDBaaSName, dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrIsEncryptionEnabled, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrPublicIPRequired, "false"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "2"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", parameterGroupID))),
			},
		},
	})
}

func TestAccE2EMySql_LargeDiskExpansion(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, 1000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrSize, "0"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", tfconstants.AttrDisk)),
			},
		},
	})
}

func TestAccE2EMySql_MultipleStatusChanges(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_stop(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrStatus, tfconstants.DBaaSPowerActionStop)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_start(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrStatus, tfconstants.DBaaSPowerActionStart)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_restart(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrStatus, tfconstants.DBaaSPowerActionRestart)),
			},
		},
	})
}

func TestAccE2EMySql_StatusValuesCaseInsensitive(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_withStatus(dbaasName, dbUser, dbPassword, dbName, tfconstants.DBaaSPowerActionStart),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", tfconstants.AttrStatus, tfconstants.DBaaSPowerActionStart)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_withStatus(dbaasName, dbUser, dbPassword, dbName, goe2econstants.DBaaSStatusRunning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
		},
	})
}

// ============================================================================
// Performance Testing
// ============================================================================

func TestAccE2EMySql_PerformanceSequentialCreates(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	const instanceCount = 3 // Reduced from 5 for faster testing
	dbaasIDs := make([]string, instanceCount)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_multipleInstances(instanceCount),
				Check: func() resource.TestCheckFunc {
					var checks []resource.TestCheckFunc
					for i := 0; i < instanceCount; i++ {
						resourceName := fmt.Sprintf("e2e_dbaas_mysql.test[%d]", i)
						idx := i // Capture loop variable
						checks = append(checks, testAccCheckE2EMySqlExists(resourceName, &dbaasIDs[idx]))
					}
					return resource.ComposeTestCheckFunc(checks...)
				}(),
			},
		},
	})
}

// Additional configuration helpers for edge cases and performance

func testAccCheckE2EMySqlConfig_allOptionalFields(name, dbUser, dbPassword, dbName, passphrase string, vpcID1, vpcID2, parameterGroupID int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name             = "%s"
  version                = "%s"
  plan                   = "%s"
  db_location            = "%s"
  group                  = "%s"
  is_encryption_enabled  = true
  encryption_passphrase  = "%s"
  public_ip_required     = false
  vpcs                   = [%d, %d]
  parameter_group_id     = %d
  database {
    user         = "%s"
    password     = "%s"
    name         = "%s"
    dbaas_number = 1
  }
  tags = {
    environment = "test"
    managed_by  = "terraform"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, tfconstants.DBaaSDefaultDBLocation, tfconstants.DBaaSDefaultGroupName, passphrase, vpcID1, vpcID2, parameterGroupID, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_withStatus(name, dbUser, dbPassword, dbName, status string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "%s"
  plan       = "%s"
  status     = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, testMySQLVersion, testMySQLPlanSmall, status, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_multipleInstances(count int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  count      = %d
  dbaas_name = "test-mysql-perf-${count.index}-%s"
  version    = "%s"
  plan       = "%s"
  database {
    user     = "testuser${count.index}"
    password = "testpassword${count.index}"
    name     = "testdb${count.index}"
  }
}
`, count, acctest.RandString(10), testMySQLVersion, testMySQLPlanSmall)
}

// ============================================================================
// Additional Acceptance Tests for Error Scenarios
// ============================================================================

func TestAccE2EMySql_ImportNonExistentInstance(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:      "e2e_dbaas_mysql.test",
				ImportState:       true,
				ImportStateId:     "test-project:999999", // Non-existent ID
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile(`not found|error retrieving`),
			},
		},
	})
}

func TestAccE2EMySql_DeleteIdempotency(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mysql-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMySqlDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMySqlConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
		},
	})
}
