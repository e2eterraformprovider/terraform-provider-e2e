package autoscaling_test

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

func TestAccE2EScalerGroup_Basic(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("e2e_scaler_group.test", "vm_image_id"),
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_Update(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_nodes", "1"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_nodes", "5"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired", "2"),
				),
			},
			{
				Config: testAccCheckE2EScalerGroupConfig_updated(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "min_nodes", "2"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "max_nodes", "10"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "desired", "3"),
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_WithEncryption(t *testing.T) {
	var scalerGroupID string
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_withEncryption(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "is_encryption_enabled", "true"),
					resource.TestCheckResourceAttr("e2e_scaler_group.test", "encryption_passphrase", "test-passphrase-123"),
				),
			},
		},
	})
}

func TestAccE2EScalerGroup_ValidationErrors(t *testing.T) {
	groupName := fmt.Sprintf("test-sg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
		ProviderFactories: testAccProviderFactories,
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
				ExpectError: regexp.MustCompile(`The argument "plan_name" is required`),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EScalerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EScalerGroupConfig_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EScalerGroupExists("e2e_scaler_group.test", &scalerGroupID),
				),
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
	if v := os.Getenv("E2E_TEST_PLAN_NAME"); v == "" {
		t.Fatal("E2E_TEST_PLAN_NAME must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_VM_IMAGE_NAME"); v == "" {
		t.Fatal("E2E_TEST_VM_IMAGE_NAME must be set for acceptance tests")
	}
}

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"e2e": func() (*schema.Provider, error) {
		return e2e.Provider(), nil
	},
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

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		group, err := client.GetScalerGroup(rs.Primary.ID, projectID, location)
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
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_scaler_group" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		_, err := client.GetScalerGroup(rs.Primary.ID, projectID, location)
		if err == nil {
			return fmt.Errorf("Scaler Group still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2EScalerGroupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  name                   = "%s"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  is_public_ip_required  = true
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  name                   = "%s"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  is_public_ip_required  = true
  min_nodes              = 2
  max_nodes              = 10
  desired                = 3
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_withEncryption(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  name                   = "%s"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = true
  encryption_passphrase  = "test-passphrase-123"
  is_public_ip_required  = true
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

// Validation error configurations

func testAccCheckE2EScalerGroupConfig_minGreaterThanDesired(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  name                   = "%s"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 5
  max_nodes              = 10
  desired                = 3
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_desiredGreaterThanMax(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  name                   = "%s"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 10
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_minNodesTooLow(name string) string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  name                   = "%s"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 0
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), name,
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

// Missing required argument configurations

func testAccCheckE2EScalerGroupConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  location               = "%s"
  name                   = "test-sg"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_LOCATION"), os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  name                   = "test-sg"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  plan_name              = "%s"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PLAN_NAME"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_missingPlanName() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  name                   = "test-sg"
  vm_image_name          = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), os.Getenv("E2E_TEST_VM_IMAGE_NAME"))
}

func testAccCheckE2EScalerGroupConfig_missingVMImageName() string {
	return fmt.Sprintf(`
resource "e2e_scaler_group" "test" {
  project_id             = "%s"
  location               = "%s"
  name                   = "test-sg"
  plan_name              = "%s"
  is_encryption_enabled  = false
  min_nodes              = 1
  max_nodes              = 5
  desired                = 2
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"), os.Getenv("E2E_TEST_PLAN_NAME"))
}
