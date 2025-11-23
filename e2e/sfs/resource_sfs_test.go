package sfs_test

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

func TestAccE2ESFS_Basic(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_basic(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "name", sfsName),
					resource.TestCheckResourceAttr("e2e_sfs.test", "disk_size", "100"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "disk_iops", "1000"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "is_encryption_enabled", "false"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "status"),
				),
			},
		},
	})
}

func TestAccE2ESFS_WithEncryption(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_withEncryption(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "is_encryption_enabled", "true"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "encryption_passphrase", "test-passphrase-123"),
				),
			},
		},
	})
}

func TestAccE2ESFS_NameValidation(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESFSConfig_invalidName(),
				ExpectError: regexp.MustCompile(`name cannot contain whitespace`),
			},
		},
	})
}

func TestAccE2ESFS_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESFSConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2ESFSConfig_missingPlan(),
				ExpectError: regexp.MustCompile(`The argument "plan" is required`),
			},
			{
				Config:      testAccCheckE2ESFSConfig_missingVpcID(),
				ExpectError: regexp.MustCompile(`The argument "vpc_id" is required`),
			},
			{
				Config:      testAccCheckE2ESFSConfig_missingDiskSize(),
				ExpectError: regexp.MustCompile(`The argument "disk_size" is required`),
			},
			{
				Config:      testAccCheckE2ESFSConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2ESFSConfig_missingDiskIOPS(),
				ExpectError: regexp.MustCompile(`The argument "disk_iops" is required`),
			},
			{
				Config:      testAccCheckE2ESFSConfig_missingRegion(),
				ExpectError: regexp.MustCompile(`The argument "region" is required`),
			},
		},
	})
}

func TestAccE2ESFS_Import(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_basic(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
				),
			},
			{
				ResourceName:            "e2e_sfs.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption_passphrase"},
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
	if v := os.Getenv("E2E_TEST_REGION"); v == "" {
		t.Fatal("E2E_TEST_REGION must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_VPC_ID"); v == "" {
		t.Fatal("E2E_TEST_VPC_ID must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_SFS_PLAN"); v == "" {
		t.Fatal("E2E_TEST_SFS_PLAN must be set for acceptance tests")
	}
}

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"e2e": func() (*schema.Provider, error) {
		return e2e.Provider(), nil
	},
}

func testAccCheckE2ESFSExists(resourceName string, sfsID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SFS ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		region := rs.Primary.Attributes["region"]

		sfs, err := client.GetSfs(rs.Primary.ID, projectID, region)
		if err != nil {
			return err
		}

		if sfs == nil {
			return fmt.Errorf("SFS not found")
		}

		*sfsID = rs.Primary.ID

		return nil
	}
}

func testAccCheckE2ESFSDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_sfs" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := rs.Primary.Attributes["region"]

		_, err := client.GetSfs(rs.Primary.ID, projectID, region)
		if err == nil {
			return fmt.Errorf("SFS still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2ESFSConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "%s"
  plan                    = "%s"
  vpc_id                  = "%s"
  disk_size               = 100
  project_id              = "%s"
  disk_iops               = 1000
  region                  = "%s"
  is_encryption_enabled   = false
}
`, name, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2ESFSConfig_withEncryption(name string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "%s"
  plan                    = "%s"
  vpc_id                  = "%s"
  disk_size               = 100
  project_id              = "%s"
  disk_iops               = 1000
  region                  = "%s"
  is_encryption_enabled   = true
  encryption_passphrase   = "test-passphrase-123"
}
`, name, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2ESFSConfig_invalidName() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "test sfs with spaces"
  plan                    = "%s"
  vpc_id                  = "%s"
  disk_size               = 100
  project_id              = "%s"
  disk_iops               = 1000
  region                  = "%s"
  is_encryption_enabled   = false
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

// Missing required argument configurations

func testAccCheckE2ESFSConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  plan       = "%s"
  vpc_id     = "%s"
  disk_size  = 100
  project_id = "%s"
  disk_iops  = 1000
  region     = "%s"
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2ESFSConfig_missingPlan() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  vpc_id     = "%s"
  disk_size  = 100
  project_id = "%s"
  disk_iops  = 1000
  region     = "%s"
}
`, os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2ESFSConfig_missingVpcID() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  plan       = "%s"
  disk_size  = 100
  project_id = "%s"
  disk_iops  = 1000
  region     = "%s"
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2ESFSConfig_missingDiskSize() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  plan       = "%s"
  vpc_id     = "%s"
  project_id = "%s"
  disk_iops  = 1000
  region     = "%s"
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2ESFSConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name      = "test-sfs"
  plan      = "%s"
  vpc_id    = "%s"
  disk_size = 100
  disk_iops = 1000
  region    = "%s"
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2ESFSConfig_missingDiskIOPS() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  plan       = "%s"
  vpc_id     = "%s"
  disk_size  = 100
  project_id = "%s"
  region     = "%s"
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"))
}

func testAccCheckE2ESFSConfig_missingRegion() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  plan       = "%s"
  vpc_id     = "%s"
  disk_size  = 100
  project_id = "%s"
  disk_iops  = 1000
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"), os.Getenv("E2E_TEST_PROJECT_ID"))
}
