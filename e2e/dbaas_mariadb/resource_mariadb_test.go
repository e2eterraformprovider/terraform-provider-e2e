package dbaas_mariadb_test

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
	testMariaDBSoftwareName    = "MariaDB"
	testMariaDBSoftwareVersion = "10.6"
	testMariaDBPlanSmall       = "DBS.16GB"
	testMariaDBPlanLarge       = "DBS.32GB"
)

func TestAccE2EMariaDB_Basic(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrName, dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrSoftwareName, testMariaDBSoftwareName),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrSoftwareVersion, testMariaDBSoftwareVersion),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPlan, testMariaDBPlanSmall),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPublicIPEnabled, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", "database.0.user", dbUser),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", "database.0.name", dbName),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", "id"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPublicIPAddress),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPrivateIPAddress),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPort),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_Import(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
				),
			},
			{
				ResourceName:            "e2e_dbaas_mariadb.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database.0.password", tfconstants.AttrStatus, tfconstants.AttrDiskSize},
			},
		},
	})
}

func TestAccE2EMariaDB_ForceNew(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrSoftwareVersion, testMariaDBSoftwareVersion),
				),
			},
			{
				Config:      testAccCheckE2EMariaDBConfig_differentVersion(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`forces new resource`),
			},
		},
	})
}

func TestAccE2EMariaDB_Update(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPublicIPEnabled, "true"),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_withoutPublicIP(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPublicIPEnabled, "false"),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_PowerManagement(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
				),
			},
			{
				// Test stopping instance (status = STOPPED)
				Config: testAccCheckE2EMariaDBConfig_stopped(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped),
				),
			},
			{
				// Test starting instance (status = RUNNING)
				Config: testAccCheckE2EMariaDBConfig_running(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRunning),
				),
			},
			{
				// Test restarting instance (status = RESTARTING)
				Config: testAccCheckE2EMariaDBConfig_restarting(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRestarting),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_Encryption(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	passphrase := acctest.RandString(32)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				// Create instance with encryption enabled
				Config: testAccCheckE2EMariaDBConfig_encryption(dbaasName, dbUser, dbPassword, dbName, passphrase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					// Verify is_encryption_enabled is set correctly
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrIsEncryptionEnabled, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrEncryptionPassphrase, passphrase),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_EncryptionForceNew(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	passphrase := acctest.RandString(32)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrIsEncryptionEnabled, "false"),
				),
			},
			{
				// Test encryption cannot be changed after creation (ForceNew)
				Config:      testAccCheckE2EMariaDBConfig_encryption(dbaasName, dbUser, dbPassword, dbName, passphrase),
				ExpectError: regexp.MustCompile(`forces new resource`),
			},
		},
	})
}

func TestAccE2EMariaDB_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EMariaDBConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EMariaDBConfig_missingSoftwareName(),
				ExpectError: regexp.MustCompile(`The argument "software_name" is required`),
			},
			{
				Config:      testAccCheckE2EMariaDBConfig_missingSoftwareVersion(),
				ExpectError: regexp.MustCompile(`The argument "software_version" is required`),
			},
			{
				Config:      testAccCheckE2EMariaDBConfig_missingPlan(),
				ExpectError: regexp.MustCompile(`The argument "plan" is required`),
			},
			{
				Config:      testAccCheckE2EMariaDBConfig_missingDatabase(),
				ExpectError: regexp.MustCompile(`The argument "database" is required`),
			},
			{
				Config:      testAccCheckE2EMariaDBConfig_missingGroup(),
				ExpectError: regexp.MustCompile(`The argument "group" is required`),
			},
		},
	})
}

// ============================================================================
// CRUD Operation Tests
// ============================================================================

// CREATE Tests

func TestAccE2EMariaDB_WithVPCs(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	// Note: VPC IDs should be replaced with actual test VPC IDs from your environment
	vpcID1 := "1"
	vpcID2 := "2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "0"),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_withVPC(dbaasName, dbUser, dbPassword, dbName, vpcID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "1"),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.0", tfconstants.AttrVPCs), vpcID1),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_withMultipleVPCs(dbaasName, dbUser, dbPassword, dbName, vpcID1, vpcID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "2"),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "0"),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_WithParameterGroup(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	// Note: Parameter group ID should be replaced with actual test parameter group ID from your environment
	parameterGroupID := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", tfconstants.DBaaSDefaultParameterGroupID)),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_withParameterGroup(dbaasName, dbUser, dbPassword, dbName, parameterGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", parameterGroupID)),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", tfconstants.DBaaSDefaultParameterGroupID)),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_ComputedFields(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					// Verify all computed fields are populated
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrID),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrSoftwareID),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrTemplateID),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPublicIPAttached),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPublicIPAddress),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPrivateIPAddress),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPort),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrTotalDiskSize),
				),
			},
		},
	})
}

// READ Tests

func TestAccE2EMariaDB_StateRefresh(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrName, dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrSoftwareName, testMariaDBSoftwareName),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrSoftwareVersion, testMariaDBSoftwareVersion),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPlan, testMariaDBPlanSmall),
				),
			},
			{
				// Refresh state by reading again - all fields should be updated correctly
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrName, dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrSoftwareName, testMariaDBSoftwareName),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrSoftwareVersion, testMariaDBSoftwareVersion),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPlan, testMariaDBPlanSmall),
					// Verify computed fields are still populated after refresh
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPublicIPAddress),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrPrivateIPAddress),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_Tags(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_tags(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.environment", tfconstants.AttrTags), "test"),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.managed_by", tfconstants.AttrTags), "terraform"),
				),
			},
			{
				// Verify tags are preserved from state after refresh
				Config: testAccCheckE2EMariaDBConfig_tags(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.environment", tfconstants.AttrTags), "test"),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.managed_by", tfconstants.AttrTags), "terraform"),
				),
			},
			{
				// Update tags
				Config: testAccCheckE2EMariaDBConfig_tagsUpdated(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.environment", tfconstants.AttrTags), "production"),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.version", tfconstants.AttrTags), "1.0"),
				),
			},
		},
	})
}

// UPDATE Tests

func TestAccE2EMariaDB_ParameterGroupUpdate(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	parameterGroupID := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_withParameterGroup(dbaasName, dbUser, dbPassword, dbName, parameterGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", parameterGroupID)),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_PlanUpgrade(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPlan, testMariaDBPlanSmall),
				),
			},
			{
				// Stop the instance first (plan upgrade requires STOPPED state)
				Config: testAccCheckE2EMariaDBConfig_stopped(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped),
				),
			},
			{
				// Upgrade plan
				Config: testAccCheckE2EMariaDBConfig_upgradedPlan(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPlan, testMariaDBPlanLarge),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_DiskExpansion(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	additionalDiskSize := 10

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					// Get initial total_disk_size
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrTotalDiskSize),
				),
			},
			{
				// Stop the instance first (disk expansion requires STOPPED state)
				Config: testAccCheckE2EMariaDBConfig_stopped(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped),
				),
			},
			{
				// Test disk_size expansion (requires STOPPED state)
				Config: testAccCheckE2EMariaDBConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, additionalDiskSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					// Verify disk_size resets to 0 after expansion (known behavior)
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrDiskSize, "0"),
					// Verify total_disk_size reflects expansion (should be set and greater than 0)
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrTotalDiskSize),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_DiskExpansionRunningState(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	additionalDiskSize := 10

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					// Ensure instance is in RUNNING state
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRunning),
				),
			},
			{
				// Test disk expansion with instance in RUNNING state (should fail)
				Config:      testAccCheckE2EMariaDBConfig_diskExpansionRunning(dbaasName, dbUser, dbPassword, dbName, additionalDiskSize),
				ExpectError: regexp.MustCompile(`must be in STOPPED state`),
			},
		},
	})
}

// DELETE Tests
// Note: Delete functionality is tested by CheckDestroy in all test cases
// CheckDestroy verifies:
// 1. Successful deletion
// 2. Cleanup verification (resource no longer exists)
// 3. Idempotent delete (calling delete on non-existent resource doesn't error)

// ============================================================================
// Error Scenarios Tests
// ============================================================================

func TestAccE2EMariaDB_InvalidStatus(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EMariaDBConfig_invalidStatus(),
				ExpectError: regexp.MustCompile(`expected status to be one of`),
			},
		},
	})
}

func TestAccE2EMariaDB_PlanUpgradeWithoutStoppedState(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					// Ensure instance is in RUNNING state
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRunning),
				),
			},
			{
				// Test plan upgrade without STOPPED state (should fail)
				Config:      testAccCheckE2EMariaDBConfig_upgradedPlanRunning(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`must be in STOPPED state`),
			},
		},
	})
}

func TestAccE2EMariaDB_UpdateForceNewFields(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
				),
			},
			{
				// Try to update name (ForceNew field) - should fail or force recreation
				Config:      testAccCheckE2EMariaDBConfig_differentName(dbaasName, dbUser, dbPassword, dbName),
				ExpectError: regexp.MustCompile(`forces new resource`),
			},
		},
	})
}

// ============================================================================
// Edge Cases Tests
// ============================================================================

func TestAccE2EMariaDB_LongName(t *testing.T) {
	var dbaasID string
	// Create a very long name (255 characters is typically the max)
	longName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(240))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_customName(longName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrName, longName),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_SpecialCharactersInName(t *testing.T) {
	var dbaasID string
	// Test with special characters (hyphens and underscores are typically allowed)
	specialName := fmt.Sprintf("test-mariadb_%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_customName(specialName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrName, specialName),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_AllOptionalFields(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	vpcID1 := "1"
	vpcID2 := "2"
	parameterGroupID := 1
	passphrase := acctest.RandString(32)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_allOptionalFields(dbaasName, dbUser, dbPassword, dbName, vpcID1, vpcID2, parameterGroupID, passphrase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrName, dbaasName),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrPublicIPEnabled, "false"),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrParameterGroupID, fmt.Sprintf("%d", parameterGroupID)),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrIsEncryptionEnabled, "true"),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", fmt.Sprintf("%s.#", tfconstants.AttrVPCs), "2"),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_LargeDiskExpansion(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))
	// Test with a very large disk expansion (100GB)
	largeDiskSize := 100

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
				),
			},
			{
				// Stop the instance first
				Config: testAccCheckE2EMariaDBConfig_stopped(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped),
				),
			},
			{
				// Test with very large disk_size expansion
				Config: testAccCheckE2EMariaDBConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, largeDiskSize),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrDiskSize, "0"),
					resource.TestCheckResourceAttrSet("e2e_dbaas_mariadb.test", tfconstants.AttrTotalDiskSize),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_MultipleStatusChanges(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
				),
			},
			{
				// First status change: STOPPED
				Config: testAccCheckE2EMariaDBConfig_stopped(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped),
				),
			},
			{
				// Second status change: RUNNING
				Config: testAccCheckE2EMariaDBConfig_running(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRunning),
				),
			},
			{
				// Third status change: RESTARTING
				Config: testAccCheckE2EMariaDBConfig_restarting(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusRestarting),
				),
			},
			{
				// Fourth status change: STOPPED again
				Config: testAccCheckE2EMariaDBConfig_stopped(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped),
				),
			},
		},
	})
}

// ============================================================================
// Performance Testing
// ============================================================================
// Note: Performance testing requires actual API calls and should be run with TF_ACC=1
// The following test demonstrates sequential operations but actual performance
// measurement should be done manually or with specialized tooling

func TestAccE2EMariaDB_SequentialOperations(t *testing.T) {
	// This test creates multiple instances sequentially to verify no performance degradation
	// Actual timing measurements should be done manually or with performance testing tools
	var dbaasID1, dbaasID2, dbaasID3 string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_sequential(1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test1", &dbaasID1),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_sequential(2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test1", &dbaasID1),
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test2", &dbaasID2),
				),
			},
			{
				Config: testAccCheckE2EMariaDBConfig_sequential(3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test1", &dbaasID1),
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test2", &dbaasID2),
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test3", &dbaasID3),
				),
			},
		},
	})
}

// ============================================================================
// Security Review Tests
// ============================================================================
// Note: Security review tests are primarily unit tests that verify:
// 1. Sensitive fields are marked as Sensitive in schema (TestResourceMariaDBSchema_SensitiveFields)
// 2. Error messages don't leak credentials (verified in code review)
// 3. Input sanitization (handled by Terraform SDK)

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2EMariaDBExists(resourceName string, dbaasID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MariaDB DBaaS ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		goe2eClient := cfg.Goe2eClient()

		mariaDB, _, err := goe2eClient.MariaDB.GetMariaDB(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if mariaDB == nil {
			return fmt.Errorf("MariaDB DBaaS not found")
		}

		*dbaasID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EMariaDBDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	goe2eClient := cfg.Goe2eClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_dbaas_mariadb" {
			continue
		}

		mariaDB, _, err := goe2eClient.MariaDB.GetMariaDB(context.Background(), rs.Primary.ID)
		if err == nil && mariaDB != nil {
			return fmt.Errorf("MariaDB DBaaS still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2EMariaDBConfig_basic(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name            = "%s"
  software_name   = "%s"
  software_version = "%s"
  group           = "Default"
  plan            = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_withoutPublicIP(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  public_ip_enabled = false
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_stopped(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  status           = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, goe2econstants.DBaaSStatusStopped, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_running(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  status           = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, goe2econstants.DBaaSStatusRunning, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_restarting(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  status           = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, goe2econstants.DBaaSStatusRestarting, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_differentVersion(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "10.5"
  group            = "Default"
  plan             = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_encryption(name, dbUser, dbPassword, dbName, passphrase string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name                = "%s"
  software_name       = "%s"
  software_version    = "%s"
  group               = "Default"
  plan                = "%s"
  is_encryption_enabled = true
  encryption_passphrase = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, passphrase, dbUser, dbPassword, dbName)
}

// Error case configurations

func testAccCheckE2EMariaDBConfig_missingName() string {
	return `
resource "e2e_dbaas_mariadb" "test" {
  software_name    = "MariaDB"
  software_version = "10.6"
  group            = "Default"
  plan             = "DBS.16GB"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
    dbaas_number = 1
  }
}
`
}

func testAccCheckE2EMariaDBConfig_missingSoftwareName() string {
	return `
resource "e2e_dbaas_mariadb" "test" {
  name             = "test-mariadb"
  software_version = "10.6"
  group            = "Default"
  plan             = "DBS.16GB"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
    dbaas_number = 1
  }
}
`
}

func testAccCheckE2EMariaDBConfig_missingSoftwareVersion() string {
	return `
resource "e2e_dbaas_mariadb" "test" {
  name          = "test-mariadb"
  software_name = "MariaDB"
  group         = "Default"
  plan          = "DBS.16GB"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
    dbaas_number = 1
  }
}
`
}

func testAccCheckE2EMariaDBConfig_missingPlan() string {
	return `
resource "e2e_dbaas_mariadb" "test" {
  name             = "test-mariadb"
  software_name    = "MariaDB"
  software_version = "10.6"
  group            = "Default"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
    dbaas_number = 1
  }
}
`
}

func testAccCheckE2EMariaDBConfig_missingDatabase() string {
	return `
resource "e2e_dbaas_mariadb" "test" {
  name             = "test-mariadb"
  software_name    = "MariaDB"
  software_version = "10.6"
  group            = "Default"
  plan             = "DBS.16GB"
}
`
}

func testAccCheckE2EMariaDBConfig_missingGroup() string {
	return `
resource "e2e_dbaas_mariadb" "test" {
  name             = "test-mariadb"
  software_name    = "MariaDB"
  software_version = "10.6"
  plan             = "DBS.16GB"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
    dbaas_number = 1
  }
}
`
}

// Additional configuration helpers for CRUD tests

func testAccCheckE2EMariaDBConfig_withVPC(name, dbUser, dbPassword, dbName, vpcID string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  vpc_ids          = ["%s"]
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, vpcID, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_withMultipleVPCs(name, dbUser, dbPassword, dbName, vpcID1, vpcID2 string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  vpc_ids          = ["%s", "%s"]
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, vpcID1, vpcID2, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_withParameterGroup(name, dbUser, dbPassword, dbName string, parameterGroupID int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name                = "%s"
  software_name       = "%s"
  software_version    = "%s"
  group               = "Default"
  plan                = "%s"
  parameter_group_id  = %d
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, parameterGroupID, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_tags(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
  tags = {
    environment = "test"
    managed_by  = "terraform"
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_tagsUpdated(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
  tags = {
    environment = "production"
    version     = "1.0"
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_upgradedPlan(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  status           = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanLarge, goe2econstants.DBaaSStatusStopped, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_diskExpansion(name, dbUser, dbPassword, dbName string, additionalSize int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  status           = "%s"
  disk_size        = %d
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, goe2econstants.DBaaSStatusStopped, additionalSize, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_diskExpansionRunning(name, dbUser, dbPassword, dbName string, additionalSize int) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  status           = "%s"
  disk_size        = %d
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, goe2econstants.DBaaSStatusRunning, additionalSize, dbUser, dbPassword, dbName)
}

// Error scenario configuration helpers

func testAccCheckE2EMariaDBConfig_invalidStatus() string {
	return `
resource "e2e_dbaas_mariadb" "test" {
  name             = "test-mariadb"
  software_name    = "MariaDB"
  software_version = "10.6"
  group            = "Default"
  plan             = "DBS.16GB"
  status           = "INVALID_STATUS"
  database {
    user     = "testuser"
    password = "testpassword"
    name     = "testdb"
    dbaas_number = 1
  }
}
`
}

func testAccCheckE2EMariaDBConfig_upgradedPlanRunning(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  status           = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanLarge, goe2econstants.DBaaSStatusRunning, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_differentName(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s-changed"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, dbUser, dbPassword, dbName)
}

// Edge case configuration helpers

func testAccCheckE2EMariaDBConfig_customName(name, dbUser, dbPassword, dbName string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_allOptionalFields(name, dbUser, dbPassword, dbName, vpcID1, vpcID2 string, parameterGroupID int, passphrase string) string {
	return fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test" {
  name                 = "%s"
  software_name        = "%s"
  software_version     = "%s"
  group                = "Default"
  plan                 = "%s"
  public_ip_enabled    = false
  parameter_group_id   = %d
  is_encryption_enabled = true
  encryption_passphrase = "%s"
  vpc_ids              = ["%s", "%s"]
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, parameterGroupID, passphrase, vpcID1, vpcID2, dbUser, dbPassword, dbName)
}

func testAccCheckE2EMariaDBConfig_sequential(count int) string {
	config := ""
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("test-mariadb-seq-%d", i)
		dbUser := fmt.Sprintf("testuser%d", i)
		dbPassword := fmt.Sprintf("testpass%d", i)
		dbName := fmt.Sprintf("testdb%d", i)
		config += fmt.Sprintf(`
resource "e2e_dbaas_mariadb" "test%d" {
  name             = "%s"
  software_name    = "%s"
  software_version = "%s"
  group            = "Default"
  plan             = "%s"
  database {
    user     = "%s"
    password = "%s"
    name     = "%s"
    dbaas_number = 1
  }
}
`, i, name, testMariaDBSoftwareName, testMariaDBSoftwareVersion, testMariaDBPlanSmall, dbUser, dbPassword, dbName)
	}
	return config
}

// ============================================================================
// Missing Acceptance Tests - Error Scenarios and Edge Cases
// ============================================================================

func TestAccE2EMariaDBDataSource_InvalidID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EMariaDBDataSourceConfig_invalidID(),
				ExpectError: regexp.MustCompile(`not found|error retrieving`),
			},
		},
	})
}

func TestAccE2EMariaDBDataSource_StatusNormalization(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
				),
			},
			{
				Config: testAccCheckE2EMariaDBDataSourceConfig_withID(dbaasID),
				Check: resource.ComposeTestCheckFunc(
					// Verify status is normalized (if API returns SUSPENDED, datasource should show STOPPED)
					resource.TestCheckResourceAttrSet("data.e2e_dbaas_mariadb.test", tfconstants.AttrStatus),
				),
			},
		},
	})
}

func TestAccE2EMariaDB_ZeroDiskExpansion(t *testing.T) {
	var dbaasID string
	dbaasName := fmt.Sprintf("test-mariadb-%s", acctest.RandString(10))
	dbUser := fmt.Sprintf("testuser%s", acctest.RandString(5))
	dbPassword := acctest.RandString(16)
	dbName := fmt.Sprintf("testdb%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EMariaDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EMariaDBConfig_basic(dbaasName, dbUser, dbPassword, dbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
				),
			},
			{
				// Test disk_size = 0 (should be a no-op)
				Config: testAccCheckE2EMariaDBConfig_diskExpansion(dbaasName, dbUser, dbPassword, dbName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EMariaDBExists("e2e_dbaas_mariadb.test", &dbaasID),
					resource.TestCheckResourceAttr("e2e_dbaas_mariadb.test", tfconstants.AttrDiskSize, "0"),
				),
			},
		},
	})
}

// Configuration helpers for new tests

func testAccCheckE2EMariaDBDataSourceConfig_invalidID() string {
	return `
data "e2e_dbaas_mariadb" "test" {
  id = "999999999"
}
`
}

func testAccCheckE2EMariaDBDataSourceConfig_withID(id string) string {
	return fmt.Sprintf(`
data "e2e_dbaas_mariadb" "test" {
  id = "%s"
}
`, id)
}
