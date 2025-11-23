package dbaas_postgress_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccE2EPostgresDBaaSDataSource_Basic(t *testing.T) {
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
				Config: testAccCheckE2EPostgresDBaaSDataSourceConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "id"),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "database_id"),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "database_name"),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "database_user"),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "status"),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "plan"),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "database_version"),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "disk"),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", "database_user", dbUser),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", "database_name", dbName),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", "plan", "E2E-2C-4GB"),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", "database_version", "15"),
				),
			},
		},
	})
}

func TestAccE2EPostgresDBaaSDataSource_WithPublicIP(t *testing.T) {
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
				Config: testAccCheckE2EPostgresDBaaSDataSourceConfig_withPublicIP(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "public_ip"),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", "private_ip"),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", "is_public_ip_attached", "true"),
				),
			},
		},
	})
}

// Configuration helpers

func testAccCheckE2EPostgresDBaaSDataSourceConfig_basic(name, dbUser, dbPassword, dbName string) string {
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

data "e2e_dbaas_postgress" "test" {
  id         = e2e_dbaas_postgress.test.id
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		dbUser, dbPassword, dbName,
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EPostgresDBaaSDataSourceConfig_withPublicIP(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name               = "%s"
  version            = "15"
  plan               = "E2E-2C-4GB"
  public_ip_required = true
  project_id         = "%s"
  location           = "%s"

  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}

data "e2e_dbaas_postgress" "test" {
  id         = e2e_dbaas_postgress.test.id
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		dbUser, dbPassword, dbName,
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
