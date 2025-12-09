package dbaas_mysql_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "dbaas_name", dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "version", "8.0"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "plan", "E2E-2C-4GB"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "public_ip_required", "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "database.0.user", dbUser),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "database.0.name", dbName),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "database.0.dbaas_number", "1"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", "id"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mysql.test", "disk")),
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
				ImportStateVerifyIgnore: []string{"database.0.password", "status", "size"},
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
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "version", "8.0")),
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
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "status", "stop")),
			},
			{
				Config: testAccCheckE2EMySqlConfig_restart(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "status", "restart")),
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
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "public_ip_required", "true")),
			},
			{
				Config: testAccCheckE2EMySqlConfig_withoutPublicIP(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "public_ip_required", "false")),
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
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "plan", "E2E-2C-4GB")),
			},
			{
				Config: testAccCheckE2EMySqlConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "plan", "E2E-4C-8GB")),
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
				ExpectError: regexp.MustCompile(`The argument "version" is required`),
			},
			{
				Config:      testAccCheckE2EMySqlConfig_missingPlan(),
				ExpectError: regexp.MustCompile(`The argument "plan" is required`),
			},
			{
				Config:      testAccCheckE2EMySqlConfig_missingDatabase(),
				ExpectError: regexp.MustCompile(`The argument "database" is required`),
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
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "is_encryption_enabled", "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "encryption_passphrase", passphrase)),
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
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID)),
			},
			{
				Config: testAccCheckE2EMySqlConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMySqlExists("e2e_dbaas_mysql.test", &dbaasID),
					// After expansion, size should be reset to 0 (cumulative behavior)
					resource.TestCheckResourceAttr("e2e_dbaas_mysql.test", "size", "0")),
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
			return fmt.Errorf("No MySQL DBaaS ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		goe2eClient := cfg.Goe2eClient()

		dbaas, _, err := goe2eClient.DBaaSMySQL.GetCluster(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if dbaas == nil {
			return fmt.Errorf("MySQL DBaaS not found")
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
			return fmt.Errorf("MySQL DBaaS still exists: %s", rs.Primary.ID)
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

		return fmt.Sprintf("%s:%s", projectID, dbaasID), nil
	}
}

// Configuration helpers

func testAccCheckE2EMySqlConfig_basic(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "8.0"
  plan       = "E2E-2C-4GB"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_stop(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "8.0"
  plan       = "E2E-2C-4GB"
  status     = "stop"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_restart(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "8.0"
  plan       = "E2E-2C-4GB"
  status     = "restart"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_withoutPublicIP(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name        = "%s"
  version           = "8.0"
  plan              = "E2E-2C-4GB"
  public_ip_required = false
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_upgradedPlan(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "8.0"
  plan       = "E2E-4C-8GB"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_differentVersion(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "5.7"
  plan       = "E2E-2C-4GB"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, dbUser, dbPassword, dbName)
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
  version                = "8.0"
  plan                   = "E2E-2C-4GB"
  is_encryption_enabled  = true
  encryption_passphrase  = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, passphrase, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_tags(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "8.0"
  plan       = "E2E-2C-4GB"
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
`, name, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_tagsUpdated(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "8.0"
  plan       = "E2E-2C-4GB"
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
`, name, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMySqlConfig_diskExpansion(name, dbUser, dbPassword, dbName string, additionalSize int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mysql" "test" {
  dbaas_name = "%s"
  version    = "8.0"
  plan       = "E2E-2C-4GB"
  size       = %d
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, additionalSize, dbUser, dbPassword, dbName)
}
