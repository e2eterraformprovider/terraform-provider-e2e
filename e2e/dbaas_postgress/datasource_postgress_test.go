package dbaas_postgress_test

import (
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testPostgresDataSourcePlanSmall = "E2E-2C-4GB"
	testPostgresDataSourceVersion   = "15"
)

func TestAccE2EPostgresDBaaSDataSource_Basic(t *testing.T) {
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
				Config: testAccCheckE2EPostgresDBaaSDataSourceConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrID),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrDatabaseID),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrDatabaseName),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrDatabaseUser),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrPlan),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrDatabaseVersion),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrDisk),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", tfconstants.AttrDatabaseUser, dbUser),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", tfconstants.AttrDatabaseName, dbName),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", tfconstants.AttrPlan, testPostgresDataSourcePlanSmall),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", tfconstants.AttrDatabaseVersion, testPostgresDataSourceVersion)),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EPostgresDBaaSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EPostgresDBaaSDataSourceConfig_withPublicIP(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrPublicIPAddress),
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_postgress.test", tfconstants.AttrPrivateIPAddress),
					resource.TestCheckResourceAttr("data.e2e_dbaas_postgress.test", tfconstants.AttrIsPublicIPAttached, "true")),
			},
		},
	})
}

// Configuration helpers

func testAccCheckE2EPostgresDBaaSDataSourceConfig_basic(name, dbUser, dbPassword, dbName string) string {
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

data "e2e_dbaas_postgress" "test" {
  %s = e2e_dbaas_postgress.test.%s
}
`, name, testPostgresDataSourceVersion, testPostgresDataSourcePlanSmall,
		dbUser, dbPassword, dbName, tfconstants.AttrID, tfconstants.AttrID)
}

func testAccCheckE2EPostgresDBaaSDataSourceConfig_withPublicIP(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_postgress" "test" {
  name               = "%s"
  version            = "%s"
  plan               = "%s"
  public_ip_required = true
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
  }
}

data "e2e_dbaas_postgress" "test" {
  %s = e2e_dbaas_postgress.test.%s
}
`, name, testPostgresDataSourceVersion, testPostgresDataSourcePlanSmall,
		dbUser, dbPassword, dbName, tfconstants.AttrID, tfconstants.AttrID)
}
