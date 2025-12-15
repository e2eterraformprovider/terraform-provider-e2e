package autoscaling_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/security_group"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EScalerGroup_Basic(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "name", groupName),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_nodes", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_nodes", "5"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "is_encryption_enabled", "false"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "is_public_ip_required", "true"),
					resource.TestCheckResourceAttrSet("e2e_scaler_group.test", "plan_id"),
					resource.TestCheckResourceAttrSet("e2e_scaler_group.test", "vm_image_id")),
			},
		},
	})
}

func TestAccE2EScalerGroup_Update(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_nodes", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_nodes", "5"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired", "2")),
			},
			{
				Config: testAccCheckE2EScalerGroupConfig_updated(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_nodes", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_nodes", "10"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired", "3")),
			},
		},
	})
}

func TestAccE2EScalerGroup_WithEncryption(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_withEncryption(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "is_encryption_enabled", "true"),
					// Do not assert the raw secret value in tests; it's a Sensitive field.
					resource.TestCheckResourceAttrSet("e2e_scaler_group.test", "encryption_passphrase")),
			},
		},
	})
}

func TestAccE2EScalerGroup_ValidationErrors(t *testing.T) {
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EScalerGroupConfig_minGreaterThanDesired(groupName),
				ExpectError: regexp.MustCompile(`min_nodes .* cannot be greater than desired`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_desiredGreaterThanMax(groupName),
				ExpectError: regexp.MustCompile(`desired .* cannot be greater than max_nodes`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_minNodesTooLow(groupName),
				ExpectError: regexp.MustCompile(`must be at least 1`),
			},
		},
	})
}

func TestAccE2EScalerGroup_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingPlanName(),
				ExpectError: regexp.MustCompile(`The argument "plan" is required`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingVMImageName(),
				ExpectError: regexp.MustCompile(`The argument "vm_image_name" is required`),
			},
		},
	})
}

func TestAccE2EScalerGroup_Import(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID)),
			},
			{
				ResourceName:            "e2e_scaler_group.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption_passphrase"},
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
	if v := os.Getenv("E2E_TEST_PLAN_NAME"); v == "" {
		t.Fatal("E2E_TEST_PLAN_NAME must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_VM_IMAGE_NAME"); v == "" {
		t.Fatal("E2E_TEST_VM_IMAGE_NAME must be set for acceptance tests")
	}
}

func testAccCheckE2EScalerGroupExists(resourceName string, scalerGroupID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Scaler Group ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)

		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
		if err != nil {
			return fmt.Errorf("failed to create GoE2E client: %w", err)
		}

		group, _, err := goe2eClient.Autoscaling.GetScalerGroup(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if group == nil {
			return fmt.Errorf("Scaler Group not found")
		}

		*scalerGroupID = rs.Primary.ID

		return nil
	}
}

func testAccCheckE2EScalerGroupDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_scaler_group" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
		if err != nil {
			return fmt.Errorf("failed to create GoE2E client: %w", err)
		}

		group, _, err := goe2eClient.Autoscaling.GetScalerGroup(context.Background(), rs.Primary.ID)
		if err == nil && group != nil {
			return fmt.Errorf("Scaler Group still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2EScalerGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "%s"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  is_public_ip_required  = true
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "%s"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  is_public_ip_required  = true
  min_nodes              = 2
  max_nodes              = 10
  desired                = 3
}
`, name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_withEncryption(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "%s"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = true
  encryption_passphrase  = "test-passphrase-123"
  is_public_ip_required  = true
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

// Validation error configurations

func testAccCheckE2EScalerGroupConfig_minGreaterThanDesired(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "%s"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 5
  max_nodes              = 10
  desired                = 3
}
`, name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_desiredGreaterThanMax(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "%s"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 10
}
`, name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_minNodesTooLow(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "%s"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 0
  max_nodes              = 5
  desired                = 2
}
`, name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

// Missing required argument configurations

func testAccCheckE2EScalerGroupConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "test-sg"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "test-sg"
  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  plan              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_missingPlanName() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "test-sg"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_missingVMImageName() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {  name                   = "test-sg"
  plan              = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PLAN_NAME"))
}

// V3 Field Names Tests

func TestAccE2EScalerGroup_V3Fields(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-v3-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_v3Fields(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "name", groupName),
					// V3 field names
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "image", os.Getenv("E2E_TEST_VM_IMAGE_NAME")),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_size", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_size", "5"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired_capacity", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusStoppedLower),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "enable_encryption", "false"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "assign_public_ip", "true"),
					// Verify V2 fields are also populated (backwards compatibility)
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_nodes", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_nodes", "5"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "running_node_count", "0"),
				),
			},
		},
	})
}

// Structured Blocks Tests

func TestAccE2EScalerGroup_ScalingPolicy(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-policy-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_scalingPolicy(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "name", groupName),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "scaling_policy.#", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "scaling_policy.0.type", "scale_up"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "scaling_policy.0.metric", "cpu_utilization"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "scaling_policy.1.type", "scale_down"),
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_ScheduledAction(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-scheduled-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_scheduledAction(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "name", groupName),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "scheduled_action.#", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "scheduled_action.0.name", "morning-scale-up"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "scheduled_action.0.action_type", "scale_up"),
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_VPCConfig(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-vpc-%s", acctest.RandString(10))
	vpcName := os.Getenv("E2E_TEST_VPC_NAME")

	if vpcName == "" {
		t.Skip("E2E_TEST_VPC_NAME not set, skipping VPC config test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_vpcConfig(groupName, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "name", groupName),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "vpc_config.#", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "vpc_config.0.name", vpcName),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["e2e_scaler_group.test"]
						if !ok {
							return fmt.Errorf("resource not found: e2e_scaler_group.test")
						}
						projectID := rs.Primary.Attributes["project_id"]
						location := acceptance.GetRegionOrLocationFromState(rs)
						return waitForAttachedVPCNames(projectID, location, rs.Primary.ID, []string{vpcName}, 0)
					},
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_NetworkConfig(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-net-%s", acctest.RandString(10))
	vpcName := os.Getenv("E2E_TEST_VPC_NAME")

	if vpcName == "" {
		t.Skip("E2E_TEST_VPC_NAME not set, skipping network_config test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_networkConfig(groupName, vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "name", groupName),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "network_config.#", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "network_config.0.assign_public_ip", "true"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "network_config.0.vpc_names.#", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "network_config.0.vpc_names.0", vpcName),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["e2e_scaler_group.test"]
						if !ok {
							return fmt.Errorf("resource not found: e2e_scaler_group.test")
						}
						projectID := rs.Primary.Attributes["project_id"]
						location := acceptance.GetRegionOrLocationFromState(rs)
						// network_config assigns public IP = true by default in the test config.
						if err := waitForAttachedVPCNames(projectID, location, rs.Primary.ID, []string{vpcName}, 0); err != nil {
							return err
						}
						return waitForPublicIPRequired(projectID, location, rs.Primary.ID, true, 0)
					},
				),
			},
		},
	})
}

// Status Change Tests

func TestAccE2EScalerGroup_StatusChange(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-status-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_statusStopped(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusStoppedLower),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "provision_status", goe2econstants.AutoscalingScalerGroupStatusStopped),
				),
			},
			{
				Config: testAccCheckE2EScalerGroupConfig_statusRunning(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusRunningLower),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "provision_status", goe2econstants.AutoscalingScalerGroupStatusRunning),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["e2e_scaler_group.test"]
						if !ok {
							return fmt.Errorf("resource not found: e2e_scaler_group.test")
						}
						projectID := rs.Primary.Attributes["project_id"]
						location := acceptance.GetRegionOrLocationFromState(rs)
						return waitForScalerGroupStatusNormalized(projectID, location, rs.Primary.ID, goe2econstants.AutoscalingScalerGroupStatusRunning, 0)
					},
				),
			},
			{
				Config: testAccCheckE2EScalerGroupConfig_statusStopped(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusStoppedLower),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["e2e_scaler_group.test"]
						if !ok {
							return fmt.Errorf("resource not found: e2e_scaler_group.test")
						}
						projectID := rs.Primary.Attributes["project_id"]
						location := acceptance.GetRegionOrLocationFromState(rs)
						return waitForScalerGroupStatusNormalized(projectID, location, rs.Primary.ID, goe2econstants.AutoscalingScalerGroupStatusStopped, 0)
					},
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_EdgeCases_ConflictsAndMissingPairs(t *testing.T) {
	groupName := fmt.Sprintf("test-sg-edge-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			// Conflicts (schema ConflictsWith)
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictImage(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictMin(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictMax(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictDesired(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictStatus(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictVPC(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictEncryption(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictPublicIP(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictPolicy(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_conflictScheduledPolicy(groupName),
				ExpectError: regexp.MustCompile(`ConflictsWith|conflicts with`),
			},

			// Missing required "pair" (CustomizeDiff)
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingImagePair(groupName),
				ExpectError: regexp.MustCompile(`either 'vm_image_name' or 'image' must be specified`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingSizePairs(groupName),
				ExpectError: regexp.MustCompile(`either 'min_nodes' or 'min_size' must be specified|either 'max_nodes' or 'max_size' must be specified|either 'desired' or 'desired_capacity' must be specified`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingMaxPair(groupName),
				ExpectError: regexp.MustCompile(`either 'max_nodes' or 'max_size' must be specified`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_missingDesiredPair(groupName),
				ExpectError: regexp.MustCompile(`either 'desired' or 'desired_capacity' must be specified`),
			},
		},
	})
}

func TestAccE2EScalerGroup_EdgeCases_ScheduledActionValidation(t *testing.T) {
	groupName := fmt.Sprintf("test-sg-sched-validate-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EScalerGroupConfig_scheduledAction_missingTargetCapacity(groupName),
				ExpectError: regexp.MustCompile(`target_capacity must be set`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_scheduledAction_missingAdjustment(groupName),
				ExpectError: regexp.MustCompile(`adjustment must be set`),
			},
		},
	})
}

func TestAccE2EScalerGroup_EdgeCases_ScalingPolicyValidation(t *testing.T) {
	groupName := fmt.Sprintf("test-sg-policy-validate-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EScalerGroupConfig_scalingPolicy_invalidMetric(groupName),
				ExpectError: regexp.MustCompile(`metric`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_scalingPolicy_invalidOperator(groupName),
				ExpectError: regexp.MustCompile(`operator`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_scalingPolicy_invalidType(groupName),
				ExpectError: regexp.MustCompile(`type`),
			},
		},
	})
}

func TestAccE2EScalerGroup_ValidationErrors_V3(t *testing.T) {
	groupName := fmt.Sprintf("test-sg-validate-v3-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EScalerGroupConfig_v3MinGreaterThanDesired(groupName),
				ExpectError: regexp.MustCompile(`min_size/min_nodes .* cannot be greater than desired_capacity/desired`),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_v3DesiredGreaterThanMax(groupName),
				ExpectError: regexp.MustCompile(`desired_capacity/desired .* cannot be greater than max_size/max_nodes`),
			},
		},
	})
}

func TestAccE2EScalerGroup_Bounds_MinEqualsDesiredEqualsMaxAllowed_V3(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-bounds-eq-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_v3MinEqDesiredEqMax(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_size", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_size", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired_capacity", "2"),
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_Perf_DesiredOnlyUpdateNoDrift(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-perf-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_desiredOnly(groupName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired_capacity", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_size", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_size", "5"),
				),
			},
			{
				Config: testAccCheckE2EScalerGroupConfig_desiredOnly(groupName, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired_capacity", "3"),
					// Assert no drift on unrelated fields for the "desired-only" fast path.
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_size", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_size", "5"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "image", os.Getenv("E2E_TEST_VM_IMAGE_NAME")),
				),
			},
		},
	})
}

// State Requirement Validation Tests

func TestAccE2EScalerGroup_SecurityGroupUpdateRequiresRunning(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basicStopped(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusStoppedLower),
				),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_securityGroupUpdateWhileStopped(groupName),
				ExpectError: regexp.MustCompile(`security group updates require scaler group to be in 'Running' state`),
			},
		},
	})
}

func TestAccE2EScalerGroup_VPCUpdateRequiresStopped(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-vpc-update-%s", acctest.RandString(10))
	vpcName := os.Getenv("E2E_TEST_VPC_NAME")

	if vpcName == "" {
		t.Skip("E2E_TEST_VPC_NAME not set, skipping VPC update test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basicRunning(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusRunningLower),
				),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_vpcUpdateWhileRunning(groupName, vpcName),
				ExpectError: regexp.MustCompile(`VPC updates require scaler group to be in 'Stopped' state`),
			},
		},
	})
}

func TestAccE2EScalerGroup_PublicIPUpdateRequiresStoppedAndVPC(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-ip-%s", acctest.RandString(10))
	vpcName := os.Getenv("E2E_TEST_VPC_NAME")

	if vpcName == "" {
		t.Skip("E2E_TEST_VPC_NAME not set, skipping public IP update test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basicRunning(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusRunningLower),
				),
			},
			{
				Config:      testAccCheckE2EScalerGroupConfig_publicIPUpdateWhileRunning(groupName),
				ExpectError: regexp.MustCompile(`public IP updates require scaler group to be in 'Stopped' state`),
			},
		},
	})
}

func TestAccE2EScalerGroup_PublicIPUpdateRequiresVPCWhenStopped(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-ip-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basicStopped(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusStoppedLower),
				),
			},
			{
				// Even though the group is stopped, public IP updates require an attached VPC.
				Config:      testAccCheckE2EScalerGroupConfig_publicIPUpdateWhileStoppedNoVPC(groupName),
				ExpectError: regexp.MustCompile(`public IP updates require at least one VPC to be attached`),
			},
		},
	})
}

// TestAccE2EScalerGroup_DeprecationWarningMessages is a unit-style acceptance test:
// it verifies deprecated V2 fields remain functional (and therefore deprecation signaling is active).
//
// Note: This repo's terraform-plugin-sdk does not expose ResourceDiff.AddWarning, so the provider
// emits deprecation warnings via log output during CustomizeDiff.
func TestAccE2EScalerGroup_DeprecationWarningMessages(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-deprec-msg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basic(groupName), // V2 config (deprecated fields)
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "vm_image_name", os.Getenv("E2E_TEST_VM_IMAGE_NAME")),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_nodes", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_nodes", "5"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired", "2"),
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_SecurityGroupUpdateSucceedsWhenRunning(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-sg-ok-%s", acctest.RandString(10))
	sgName := fmt.Sprintf("test-sg-linked-%s", acctest.RandString(8))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_runningWithSecurityGroup(groupName, sgName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "status", goe2econstants.AutoscalingScalerGroupStatusRunningLower),
				),
			},
			{
				Config: testAccCheckE2EScalerGroupConfig_runningWithSecurityGroup(groupName, sgName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					func(s *terraform.State) error {
						sgRes, ok := s.RootModule().Resources["e2e_security_group.sg"]
						if !ok {
							return fmt.Errorf("resource not found: e2e_security_group.sg")
						}
						sgID, err := strconv.Atoi(sgRes.Primary.ID)
						if err != nil {
							return fmt.Errorf("invalid sg id: %q", sgRes.Primary.ID)
						}
						return waitForSecurityGroupIDsInState(s, "e2e_scaler_group.test", []int{sgID}, 0)
					},
				),
			},
		},
	})
}

// Test Configuration Helpers

func testAccCheckE2EScalerGroupConfig_v3Fields(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
  enable_encryption = false
  assign_public_ip = true
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_scalingPolicy(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  scaling_policy {
    type               = "scale_up"
    adjustment         = 2
    metric             = "cpu_utilization"
    operator           = ">"
    threshold          = "80"
    evaluation_periods = 3
    period_seconds     = 60
    cooldown_seconds   = 300
  }

  scaling_policy {
    type               = "scale_down"
    adjustment         = 1
    metric             = "cpu_utilization"
    operator           = "<"
    threshold          = "20"
    evaluation_periods = 3
    period_seconds     = 60
    cooldown_seconds   = 300
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_scheduledAction(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  scheduled_action {
    name        = "morning-scale-up"
    action_type = "scale_up"
    adjustment  = 2
    recurrence  = "0 9 * * *"
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_vpcConfig(name, vpcName string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  vpc_config {
    name = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower, vpcName)
}

func testAccCheckE2EScalerGroupConfig_networkConfig(name, vpcName string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  network_config {
    assign_public_ip = true
    vpc_names        = ["%s"]
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower, vpcName)
}

func testAccCheckE2EScalerGroupConfig_statusStopped(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_statusRunning(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusRunningLower)
}

func testAccCheckE2EScalerGroupConfig_basicStopped(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_basicRunning(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusRunningLower)
}

func testAccCheckE2EScalerGroupConfig_securityGroupUpdateWhileStopped(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
  security_group_ids = [999]  # Attempting to update while stopped should fail
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_vpcUpdateWhileRunning(name, vpcName string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  vpc_config {
    name = "%s"
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusRunningLower, vpcName)
}

func testAccCheckE2EScalerGroupConfig_publicIPUpdateWhileRunning(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
  assign_public_ip = false  # Attempting to update while running should fail
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusRunningLower)
}

func testAccCheckE2EScalerGroupConfig_publicIPUpdateWhileStoppedNoVPC(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
  assign_public_ip = false  # Attempting to update without VPC should fail
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_runningWithSecurityGroup(groupName, sgName string, attach bool) string {
	attachLine := ""
	if attach {
		// SG id is a string; autoscaling schema expects int list.
		attachLine = "  security_group_ids = [tonumber(e2e_security_group.sg.id)]\n"
	}
	return fmt.Sprintf(`
resource "e2e_security_group" "sg" {
  name = "%s"
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}

resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
%s
}
`, sgName,
		security_group.RuleTypeInbound,
		security_group.ProtocolAll,
		security_group.NetworkTypeAny,
		groupName,
		os.Getenv("E2E_TEST_PLAN_NAME"),
		os.Getenv("E2E_TEST_VM_IMAGE_NAME"),
		goe2econstants.AutoscalingScalerGroupStatusRunningLower,
		attachLine,
	)
}

func testAccCheckE2EScalerGroupConfig_conflictImage(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  vm_image_name    = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictMin(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_nodes        = 1
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictMax(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_nodes        = 5
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictDesired(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired          = 2
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictStatus(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  provision_status = "%s"
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStopped, goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictVPC(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  vpc = [{ name = "dummy" }]
  vpc_config { name = "dummy" }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictEncryption(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name                = "%s"
  plan                = "%s"
  image               = "%s"
  min_size            = 1
  max_size            = 5
  desired_capacity    = 2
  status              = "%s"
  is_encryption_enabled = true
  enable_encryption     = true
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictPublicIP(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name                = "%s"
  plan                = "%s"
  image               = "%s"
  min_size            = 1
  max_size            = 5
  desired_capacity    = 2
  status              = "%s"
  is_public_ip_required = true
  assign_public_ip      = true
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictPolicy(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  policy = [{
    type           = "upscale"
    adjust         = 1
    parameter      = "cpu"
    operator       = ">"
    value          = "80"
    period_number  = "1"
    period_seconds = "60"
    cooldown       = "60"
  }]

  scaling_policy {
    type               = "scale_up"
    adjustment         = 1
    metric             = "cpu_utilization"
    operator           = ">"
    threshold          = "80"
    evaluation_periods = 1
    period_seconds     = 60
    cooldown_seconds   = 60
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_conflictScheduledPolicy(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  scheduled_policy = [{
    type       = "upscale"
    adjust     = "1"
    recurrence = "0 9 * * *"
  }]

  scheduled_action {
    name        = "morning-scale-up"
    action_type = "scale_up"
    adjustment  = 1
    recurrence  = "0 9 * * *"
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_missingImagePair(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_missingSizePairs(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name   = "%s"
  plan   = "%s"
  image  = "%s"
  status = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_missingMaxPair(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  desired_capacity = 1
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_missingDesiredPair(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name     = "%s"
  plan     = "%s"
  image    = "%s"
  min_size = 1
  max_size = 2
  status   = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_scheduledAction_missingTargetCapacity(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  scheduled_action {
    name        = "set-capacity"
    action_type = "set_capacity"
    recurrence  = "0 9 * * *"
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_scheduledAction_missingAdjustment(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  scheduled_action {
    name        = "scale-up"
    action_type = "scale_up"
    recurrence  = "0 9 * * *"
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_scalingPolicy_invalidMetric(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  scaling_policy {
    type               = "scale_up"
    adjustment         = 1
    metric             = "disk_utilization"
    operator           = ">"
    threshold          = "80"
    evaluation_periods = 1
    period_seconds     = 60
    cooldown_seconds   = 60
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_scalingPolicy_invalidOperator(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  scaling_policy {
    type               = "scale_up"
    adjustment         = 1
    metric             = "cpu_utilization"
    operator           = "!="
    threshold          = "80"
    evaluation_periods = 1
    period_seconds     = 60
    cooldown_seconds   = 60
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_scalingPolicy_invalidType(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 2
  status           = "%s"

  scaling_policy {
    type               = "scale_out"
    adjustment         = 1
    metric             = "cpu_utilization"
    operator           = ">"
    threshold          = "80"
    evaluation_periods = 1
    period_seconds     = 60
    cooldown_seconds   = 60
  }
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_v3MinGreaterThanDesired(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 5
  max_size         = 10
  desired_capacity = 3
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_v3DesiredGreaterThanMax(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = 10
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_v3MinEqDesiredEqMax(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 2
  max_size         = 2
  desired_capacity = 2
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}

func testAccCheckE2EScalerGroupConfig_desiredOnly(name string, desired int) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  name             = "%s"
  plan             = "%s"
  image            = "%s"
  min_size         = 1
  max_size         = 5
  desired_capacity = %d
  status           = "%s"
}
`, name, os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"), desired, goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}
