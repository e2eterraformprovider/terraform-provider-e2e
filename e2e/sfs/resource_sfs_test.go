package sfs_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2ESFS_Basic(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "status")),
			},
		},
	})
}

func TestAccE2ESFS_WithEncryption(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_withEncryption(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "is_encryption_enabled", "true"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "encryption_passphrase", "test-passphrase-123")),
			},
		},
	})
}

func TestAccE2ESFS_NameValidation(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
		ProviderFactories: acceptance.TestAccProviderFactories,
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_basic(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID)),
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

// V3 Field Tests

func TestAccE2ESFS_V3Fields(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_v3Fields(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "name", sfsName),
					resource.TestCheckResourceAttr("e2e_sfs.test", "size_gb", "100"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "iops", "1000"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "encryption_enabled", "false"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "state"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "mount_endpoint")),
			},
		},
	})
}

func TestAccE2ESFS_V3FieldsWithEncryption(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_v3FieldsEncrypted(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "encryption_enabled", "true"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "encryption_passphrase", "test-passphrase-v3")),
			},
		},
	})
}

func TestAccE2ESFS_DeprecatedDiskSize(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_basic(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					// Verify both V2 and V3 field names work
					resource.TestCheckResourceAttr("e2e_sfs.test", "disk_size", "100"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "size_gb")),
			},
		},
	})
}

func TestAccE2ESFS_DeprecatedEncryptionFlag(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_withEncryption(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					// Verify both V2 and V3 flag names work
					resource.TestCheckResourceAttr("e2e_sfs.test", "is_encryption_enabled", "true"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "encryption_enabled")),
			},
		},
	})
}

func TestAccE2ESFS_Tags(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_withTags(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.Purpose", "testing")),
			},
		},
	})
}

func TestAccE2ESFS_MountEndpoint(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESFSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESFSConfig_v3Fields(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "mount_endpoint"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "private_endpoint")),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
	if v := os.Getenv("E2E_TEST_VPC_ID"); v == "" {
		t.Fatal("E2E_TEST_VPC_ID must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_SFS_PLAN"); v == "" {
		t.Fatal("E2E_TEST_SFS_PLAN must be set for acceptance tests")
	}
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

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)

		client, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		sfs, _, err := client.Sfs.GetSfs(context.Background(), rs.Primary.ID)
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
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_sfs" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)

		client, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		_, _, err = client.Sfs.GetSfs(context.Background(), rs.Primary.ID)
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
  disk_size               = 100  disk_iops               = 1000  is_encryption_enabled   = false
}
`, name, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_withEncryption(name string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "%s"
  plan                    = "%s"
  vpc_id                  = "%s"
  disk_size               = 100  disk_iops               = 1000  is_encryption_enabled   = true
  encryption_passphrase   = "test-passphrase-123"
}
`, name, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_invalidName() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "test sfs with spaces"
  plan                    = "%s"
  vpc_id                  = "%s"
  disk_size               = 100  disk_iops               = 1000  is_encryption_enabled   = false
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

// Missing required argument configurations

func testAccCheckE2ESFSConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  plan       = "%s"
  vpc_id     = "%s"
  disk_size  = 100  disk_iops  = 1000}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_missingPlan() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  vpc_id     = "%s"
  disk_size  = 100  disk_iops  = 1000}
`, os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_missingVpcID() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  plan       = "%s"
  disk_size  = 100  disk_iops  = 1000}
`, os.Getenv("E2E_TEST_SFS_PLAN"))
}

func testAccCheckE2ESFSConfig_missingDiskSize() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  plan       = "%s"
  vpc_id     = "%s"  disk_iops  = 1000}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name      = "test-sfs"
  plan      = "%s"
  vpc_id    = "%s"
  disk_size = 100
  disk_iops = 1000}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_missingDiskIOPS() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  plan       = "%s"
  vpc_id     = "%s"
  disk_size  = 100}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_missingRegion() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name       = "test-sfs"
  plan       = "%s"
  vpc_id     = "%s"
  disk_size  = 100  disk_iops  = 1000
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

// V3 Field Configuration Helpers

func testAccCheckE2ESFSConfig_v3Fields(name string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                = "%s"
  plan                = "%s"
  vpc_id              = "%s"
  size_gb             = 100
  iops                = 1000
  encryption_enabled  = false
}
`, name, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_v3FieldsEncrypted(name string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                = "%s"
  plan                = "%s"
  vpc_id              = "%s"
  size_gb             = 100
  iops                = 1000
  encryption_enabled  = true
  encryption_passphrase = "test-passphrase-v3"
}
`, name, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_withTags(name string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                = "%s"
  plan                = "%s"
  vpc_id              = "%s"
  size_gb             = 100
  iops                = 1000
  encryption_enabled  = false
  tags = {
    Environment = "test"
    Purpose     = "testing"
  }
}
`, name, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}
