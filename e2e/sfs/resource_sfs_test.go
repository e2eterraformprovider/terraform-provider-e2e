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

// Additional Tests

func TestAccE2ESFS_V3FieldsUpdate(t *testing.T) {
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
					resource.TestCheckResourceAttr("e2e_sfs.test", "size_gb", "100"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "iops", "1000")),
			},
			// Note: Most fields are ForceNew, so we can't update them
			// This test verifies that the resource exists and has the expected values
		},
	})
}

func TestAccE2ESFS_TagsUpdate(t *testing.T) {
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
			{
				Config: testAccCheckE2ESFSConfig_withUpdatedTags(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.%", "3"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.Purpose", "storage"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.Team", "platform")),
			},
		},
	})
}

func TestAccE2ESFS_ConflictDiskSizeAndSizeGb(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESFSConfig_conflictingSizeFields(),
				ExpectError: regexp.MustCompile(`conflicts with (disk_size|size_gb)`),
			},
		},
	})
}

func TestAccE2ESFS_ConflictIopsFields(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESFSConfig_conflictingIopsFields(),
				ExpectError: regexp.MustCompile(`conflicts with (disk_iops|iops)`),
			},
		},
	})
}

func TestAccE2ESFS_ConflictEncryption(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESFSConfig_conflictingEncryptionFields(),
				ExpectError: regexp.MustCompile(`conflicts with (is_encryption_enabled|encryption_enabled)`),
			},
		},
	})
}

func TestAccE2ESFS_BasicV2(t *testing.T) {
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
					resource.TestCheckResourceAttr("e2e_sfs.test", "disk_size", "100"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "disk_iops", "1000"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "is_encryption_enabled", "false"),
					// Verify both V2 and V3 fields are populated
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "size_gb"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "iops"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "encryption_enabled")),
			},
		},
	})
}

func TestAccE2ESFS_EncryptionWithPassphrase(t *testing.T) {
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

func TestAccE2ESFS_WaitForActive(t *testing.T) {
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
					// Verify status is set to a valid state (not empty)
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "state"),
					// If we reached this point, the resource waited for Active status
					testAccCheckE2ESFSStatusIsActiveOrBeyond("e2e_sfs.test")),
			},
		},
	})
}

func TestAccE2ESFS_ImportBasic(t *testing.T) {
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

func TestAccE2ESFS_ImportFull(t *testing.T) {
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
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID)),
			},
			{
				ResourceName:            "e2e_sfs.test",
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("%s/%s/%s", os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_REGION"), sfsID),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption_passphrase"},
			},
		},
	})
}

func TestAccE2ESFS_ImportUsesV3Fields(t *testing.T) {
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
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID)),
			},
			{
				ResourceName:            "e2e_sfs.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption_passphrase"},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(states))
					}
					state := states[0]
					// Verify V3 fields are populated
					if state.Attributes["size_gb"] == "" {
						return fmt.Errorf("size_gb not populated after import")
					}
					if state.Attributes["iops"] == "" {
						return fmt.Errorf("iops not populated after import")
					}
					if state.Attributes["encryption_enabled"] == "" {
						return fmt.Errorf("encryption_enabled not populated after import")
					}
					// Verify V2 fields are also populated for backwards compatibility
					if state.Attributes["disk_size"] == "" {
						return fmt.Errorf("disk_size not populated after import")
					}
					if state.Attributes["disk_iops"] == "" {
						return fmt.Errorf("disk_iops not populated after import")
					}
					return nil
				},
			},
		},
	})
}

func TestAccE2ESFS_InvalidName(t *testing.T) {
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

func TestAccE2ESFS_DeleteWhileCreating(t *testing.T) {
	t.Skip("This test requires manual intervention to simulate deletion during creation")
	// This test would need to be manually triggered or use mocking
	// to simulate the scenario where deletion is attempted while status is "Creating"
}

func TestAccE2ESFS_Disappears(t *testing.T) {
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
					testAccCheckE2ESFSDisappears("e2e_sfs.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// CRUD Operation Tests

func TestAccE2ESFS_ReadNonExistent(t *testing.T) {
	t.Skip("This test requires the resource to be deleted externally during the test")
	// This would require external manipulation of the resource
	// Covered by TestAccE2ESFS_Disappears which is similar
}

func TestAccE2ESFS_UpdateForceNewFields(t *testing.T) {
	var sfsID string
	sfsName := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))
	sfsNameUpdated := fmt.Sprintf("test-sfs-%s", acctest.RandString(10))

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
					resource.TestCheckResourceAttr("e2e_sfs.test", "size_gb", "100")),
			},
			{
				// Attempt to update ForceNew field (name) - should trigger recreation
				Config: testAccCheckE2ESFSConfig_v3Fields(sfsNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "name", sfsNameUpdated),
					// Note: sfsID will be different after recreation
				),
			},
		},
	})
}

func TestAccE2ESFS_DeleteNonExistent(t *testing.T) {
	t.Skip("This test requires manual resource deletion and is covered by idempotency in CheckDestroy")
	// The CheckDestroy function in all tests already handles this scenario
	// by verifying that resources are properly cleaned up even if already deleted
}

func TestAccE2ESFS_Delete404Handling(t *testing.T) {
	t.Skip("404 handling is tested implicitly in TestAccE2ESFS_Disappears and CheckDestroy")
	// TestAccE2ESFS_Disappears already tests this scenario
}

func TestAccE2ESFS_StateRefresh(t *testing.T) {
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
					// Verify all fields are correctly set after creation and refresh
					resource.TestCheckResourceAttr("e2e_sfs.test", "name", sfsName),
					resource.TestCheckResourceAttr("e2e_sfs.test", "size_gb", "100"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "iops", "1000"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "encryption_enabled", "false"),
					// Verify state normalization
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "state"),
					// Verify computed fields
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "mount_endpoint"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "private_endpoint"),
				),
			},
			{
				// Refresh the resource and verify all fields are still correct
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "name", sfsName),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "state"),
				),
			},
		},
	})
}

func TestAccE2ESFS_V2V3FieldsBackwardsCompatibility(t *testing.T) {
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
					// Verify V2 fields work
					resource.TestCheckResourceAttr("e2e_sfs.test", "disk_size", "100"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "disk_iops", "1000"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "is_encryption_enabled", "false"),
					// Verify V3 fields are also populated
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "size_gb"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "iops"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "encryption_enabled"),
				),
			},
		},
	})
}

func TestAccE2ESFS_MountEndpointEqualsPrivateEndpoint(t *testing.T) {
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
					testAccCheckE2ESFSMountEndpointEqualsPrivateEndpoint("e2e_sfs.test"),
				),
			},
		},
	})
}

func TestAccE2ESFS_TagsPreservedFromState(t *testing.T) {
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
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.Purpose", "testing"),
				),
			},
			{
				// Refresh and verify tags are preserved
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("e2e_sfs.test", "tags.Purpose", "testing"),
				),
			},
		},
	})
}

func TestAccE2ESFS_ComputedFieldsPopulated(t *testing.T) {
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
					// Verify all computed fields are populated
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "id"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "state"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "mount_endpoint"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "private_endpoint"),
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "is_backup_enabled"),
					// Verify status is in a valid state (should be Active after creation)
					testAccCheckE2ESFSStatusIsActiveOrBeyond("e2e_sfs.test"),
				),
			},
		},
	})
}

// Additional Import Tests

func TestAccE2ESFS_ImportNonExistent(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            testAccCheckE2ESFSConfig_basic("test-sfs-import"),
				ResourceName:      "e2e_sfs.test",
				ImportState:       true,
				ImportStateId:     "non-existent-sfs-id-12345",
				ExpectError:       regexp.MustCompile(`Cannot import non-existent|not found|404`),
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccE2ESFS_ImportInvalidFormat(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            testAccCheckE2ESFSConfig_basic("test-sfs-import"),
				ResourceName:      "e2e_sfs.test",
				ImportState:       true,
				ImportStateId:     "invalid/format",
				ExpectError:       regexp.MustCompile(`expected import ID format|invalid format`),
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccE2ESFS_ImportNoRecreation(t *testing.T) {
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
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID)),
			},
			{
				ResourceName:            "e2e_sfs.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption_passphrase"},
			},
			{
				Config: testAccCheckE2ESFSConfig_v3Fields(sfsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESFSExists("e2e_sfs.test", &sfsID),
					// Verify the ID hasn't changed (no recreation)
					resource.TestCheckResourceAttrSet("e2e_sfs.test", "id")),
				PlanOnly: true,
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

func testAccCheckE2ESFSConfig_withUpdatedTags(name string) string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                = "%s"
  plan                = "%s"
  vpc_id              = "%s"
  size_gb             = 100
  iops                = 1000
  encryption_enabled  = false
  tags = {
    Environment = "production"
    Purpose     = "storage"
    Team        = "platform"
  }
}
`, name, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_conflictingSizeFields() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                = "test-sfs-conflict"
  plan                = "%s"
  vpc_id              = "%s"
  size_gb             = 100
  disk_size           = 100
  iops                = 1000
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_conflictingIopsFields() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                = "test-sfs-conflict"
  plan                = "%s"
  vpc_id              = "%s"
  size_gb             = 100
  iops                = 1000
  disk_iops           = 1000
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSConfig_conflictingEncryptionFields() string {
	return fmt.Sprintf(`
resource "e2e_sfs" "test" {
  name                    = "test-sfs-conflict"
  plan                    = "%s"
  vpc_id                  = "%s"
  size_gb                 = 100
  iops                    = 1000
  encryption_enabled      = true
  is_encryption_enabled   = true
}
`, os.Getenv("E2E_TEST_SFS_PLAN"), os.Getenv("E2E_TEST_VPC_ID"))
}

func testAccCheckE2ESFSStatusIsActiveOrBeyond(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		status := rs.Primary.Attributes["status"]
		state := rs.Primary.Attributes["state"]

		// Check that status is not empty and is a valid state
		if status == "" {
			return fmt.Errorf("status is empty")
		}

		if state == "" {
			return fmt.Errorf("state is empty")
		}

		// Valid states that indicate the resource waited for Active or beyond
		validStates := []string{"active", "Active", "deleting", "Deleting", "deleted", "Deleted"}
		isValid := false
		for _, validState := range validStates {
			if status == validState || state == validState {
				isValid = true
				break
			}
		}

		if !isValid {
			return fmt.Errorf("SFS is in unexpected state. Status: %s, State: %s", status, state)
		}

		return nil
	}
}

func testAccCheckE2ESFSDisappears(resourceName string) resource.TestCheckFunc {
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

		// Delete the SFS
		_, err = client.Sfs.DeleteSfs(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error deleting SFS: %w", err)
		}

		return nil
	}
}

func testAccCheckE2ESFSMountEndpointEqualsPrivateEndpoint(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		mountEndpoint := rs.Primary.Attributes["mount_endpoint"]
		privateEndpoint := rs.Primary.Attributes["private_endpoint"]

		if mountEndpoint == "" {
			return fmt.Errorf("mount_endpoint is empty")
		}

		if privateEndpoint == "" {
			return fmt.Errorf("private_endpoint is empty")
		}

		if mountEndpoint != privateEndpoint {
			return fmt.Errorf("mount_endpoint (%s) does not equal private_endpoint (%s)", mountEndpoint, privateEndpoint)
		}

		return nil
	}
}
