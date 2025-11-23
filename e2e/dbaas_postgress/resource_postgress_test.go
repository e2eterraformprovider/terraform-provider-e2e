package dbaas_postgress_test

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

func TestAccE2EPostgresDBaaS_Basic(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-pg-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "name", dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "version", "15"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "plan", "E2E-2C-4GB"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "public_ip_required", "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "database.0.user", dbUser),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "database.0.name", dbName),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "database.0.dbaas_number", "1"),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "is_encryption_enabled", "false"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", "id"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", "status_title"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", "num_instances"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", "project_name"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_postgress.test", "connectivity_detail"),
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
		ProviderFactories: testAccProviderFactories,
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
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "power_status", "stop"),
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
		ProviderFactories: testAccProviderFactories,
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
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "power_status", "stop"),
				),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_restart(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "power_status", "restart"),
				),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withEncryption(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "is_encryption_enabled", "true"),
				),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_withoutPublicIP(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "public_ip_required", "false"),
				),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "plan", "E2E-2C-4GB"),
				),
			},
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_postgress.test", "plan", "E2E-4C-8GB"),
				),
			},
		},
	})
}

func TestAccE2EPostgresDBaaS_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingVersion(),
				ExpectError: regexp.MustCompile(`The argument "version" is required`),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingPlan(),
				ExpectError: regexp.MustCompile(`The argument "plan" is required`),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingDatabase(),
				ExpectError: regexp.MustCompile(`The argument "database" is required`),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
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
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EPostgresDBaaSConfig_invalidPowerStatus(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`expected power_status to be one of`),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EPostgresDBaaSExists("e2e_dbaas_postgress.test", &dbaasID),
				),
			},
			{
				ResourceName:            "e2e_dbaas_postgress.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccE2EPostgresDBaaSImportID("e2e_dbaas_postgress.test"),
				ImportStateVerifyIgnore: []string{"database.0.password", "power_status", "size", "detach_public_ip"},
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

func testAccCheckE2EPostgresDBaaSExists(resourceName string, dbaasID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Postgres DBaaS ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		dbaas, err := client.GetPostgressDB(rs.Primary.ID, projectID, location)
		if err != nil {
			return err
		}

		if dbaas == nil {
			return fmt.Errorf("Postgres DBaaS not found")
		}

		*dbaasID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EPostgresDBaaSDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_dbaas_postgress" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		_, err := client.GetPostgressDB(rs.Primary.ID, projectID, location)
		if err == nil {
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
		location := rs.Primary.Attributes["location"]
		dbaasID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, dbaasID), nil
	}
}

// Configuration helpers

func testAccCheckE2EPostgresDBaaSConfig_basic(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "15"
  plan       = "E2E-2C-4GB"
  project_id = "%s"
  location   = "%s"

  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_powerOff(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name         = "%s"
  version      = "15"
  plan         = "E2E-2C-4GB"
  power_status = "stop"
  project_id   = "%s"
  location     = "%s"

  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_restart(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name         = "%s"
  version      = "15"
  plan         = "E2E-2C-4GB"
  power_status = "restart"
  project_id   = "%s"
  location     = "%s"

  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withEncryption(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name                    = "%s"
  version                 = "15"
  plan                    = "E2E-2C-4GB"
  is_encryption_enabled   = true
  project_id              = "%s"
  location                = "%s"

  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_withoutPublicIP(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name               = "%s"
  version            = "15"
  plan               = "E2E-2C-4GB"
  public_ip_required = false
  project_id         = "%s"
  location           = "%s"

  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), dbUser, dbPassword, dbName)
}

func testAccCheckE2EPostgresDBaaSConfig_upgradedPlan(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "%s"
  version    = "15"
  plan       = "E2E-4C-8GB"
  project_id = "%s"
  location   = "%s"

  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), dbUser, dbPassword, dbName)
}

// Error case configurations

func testAccCheckE2EPostgresDBaaSConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  version    = "15"
  plan       = "E2E-2C-4GB"
  project_id = "%s"
  location   = "%s"

  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EPostgresDBaaSConfig_missingVersion() string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "test-pg"
  plan       = "E2E-2C-4GB"
  project_id = "%s"
  location   = "%s"

  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EPostgresDBaaSConfig_missingPlan() string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "test-pg"
  version    = "15"
  project_id = "%s"
  location   = "%s"

  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EPostgresDBaaSConfig_missingDatabase() string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "test-pg"
  version    = "15"
  plan       = "E2E-2C-4GB"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EPostgresDBaaSConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name     = "test-pg"
  version  = "15"
  plan     = "E2E-2C-4GB"
  location = "%s"

  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EPostgresDBaaSConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name       = "test-pg"
  version    = "15"
  plan       = "E2E-2C-4GB"
  project_id = "%s"

  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
  }
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}

func testAccCheckE2EPostgresDBaaSConfig_invalidPowerStatus(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name         = "%s"
  version      = "15"
  plan         = "E2E-2C-4GB"
  power_status = "invalid"
  project_id   = "%s"
  location     = "%s"

  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), dbUser, dbPassword, dbName)
}
