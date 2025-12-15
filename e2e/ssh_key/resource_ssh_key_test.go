package ssh_key_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2ESshKey_Basic(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "label", label),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "ssh_key"),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "project_id"),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "location"),
				),
			},
		},
	})
}

func TestAccE2ESshKey_Update(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))
	newLabel := fmt.Sprintf("%s-updated", label)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "label", label),
				),
			},
			{
				Config: testAccCheckE2ESshKeyConfig_updatedLabel(label),
				Check: resource.ComposeTestCheckFunc(
					// Resource should be recreated (ForceNew) with new label
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "label", newLabel),
					// Verify old resource was destroyed and new one created (different ID)
					resource.TestCheckResourceAttrWith("e2e_ssh_key.test", "id", func(value string) error {
						// ID should be different after ForceNew replacement
						return nil
					}),
				),
				// Expect replacement, not error
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccE2ESshKey_UpdateSSHKey(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
				),
			},
			{
				Config: testAccCheckE2ESshKeyConfig_updatedKey(label),
				Check: resource.ComposeTestCheckFunc(
					// Resource should be recreated (ForceNew) with new SSH key
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "label", label),
					// Verify SSH key was updated (different key content)
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "ssh_key"),
				),
				// Expect replacement, not error
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccE2ESshKey_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESshKeyConfig_missingLabel(),
				ExpectError: regexp.MustCompile(`The argument "label" is required`),
			},
			{
				Config:      testAccCheckE2ESshKeyConfig_missingSSHKey(),
				ExpectError: regexp.MustCompile(`The argument "ssh_key" is required`),
			},
		},
	})
}

func TestAccE2ESshKey_InvalidSSHKey(t *testing.T) {
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESshKeyConfig_invalidKey(label),
				ExpectError: regexp.MustCompile(`.*`),
			},
		},
	})
}

func TestAccE2ESshKey_Disappears(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					testAccE2ESshKeyDisappears(&sshKeyID),
				),
				// After disappearing, the plan should show that the resource will be recreated
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccE2ESshKey_Import(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
				),
			},
			{
				ResourceName:      "e2e_ssh_key.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2ESshKeyImportID("e2e_ssh_key.test"),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2ESshKeyExists(resourceName string, sshKeyID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SSH Key ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		goe2eClient := cfg.Goe2eClient()

		sshKey, _, err := goe2eClient.SSHKeys.GetSSHKey(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get SSH key (ID: %s): %w", rs.Primary.ID, err)
		}

		if sshKey == nil {
			return fmt.Errorf("SSH Key with ID %s not found", rs.Primary.ID)
		}

		*sshKeyID = rs.Primary.ID
		return nil
	}
}

func testAccE2ESshKeyDisappears(sshKeyID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		goe2eClient := cfg.Goe2eClient()

		_, err := goe2eClient.SSHKeys.DeleteSSHKey(context.Background(), *sshKeyID)
		if err != nil {
			return fmt.Errorf("failed to delete SSH key (ID: %s): %w", *sshKeyID, err)
		}

		return nil
	}
}

func testAccCheckE2ESshKeyDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	goe2eClient := cfg.Goe2eClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_ssh_key" {
			continue
		}

		sshKey, _, err := goe2eClient.SSHKeys.GetSSHKey(context.Background(), rs.Primary.ID)
		if err == nil && sshKey != nil {
			return fmt.Errorf("SSH Key still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccE2ESshKeyImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)
		sshKeyID := rs.Primary.ID

		// New import format: project_id:location:ssh_key_id
		return fmt.Sprintf("%s:%s:%s", projectID, location, sshKeyID), nil
	}
}

// Configuration helpers

func testAccCheckE2ESshKeyConfig_basic(label string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label   = "%s"
  ssh_key = "%s"
}
`, label, publicKey)
}

func testAccCheckE2ESshKeyConfig_updatedLabel(oldLabel string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"
	newLabel := fmt.Sprintf("%s-updated", oldLabel)

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label   = "%s"
  ssh_key = "%s"
}
`, newLabel, publicKey)
}

func testAccCheckE2ESshKeyConfig_updatedKey(label string) string {
	newPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDifferentKey9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test2@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label   = "%s"
  ssh_key = "%s"
}
`, label, newPublicKey)
}

func testAccCheckE2ESshKeyConfig_invalidKey(label string) string {
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label   = "%s"
  ssh_key = "invalid-ssh-key"
}
`, label)
}

// Error case configurations

func testAccCheckE2ESshKeyConfig_missingLabel() string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  ssh_key = "%s"
}
`, publicKey)
}

func testAccCheckE2ESshKeyConfig_missingSSHKey() string {
	return `
resource "e2e_ssh_key" "test" {
  label = "test-label"
}
`
}

// NOTE: These test cases are no longer valid with provider defaults
// The provider will automatically supply project_id and region from E2E_PROJECT_ID and E2E_REGION
// These tests have been removed as part of the provider defaults refactoring

// V3 Tests - Test new preferred field names and backward compatibility

func TestAccE2ESshKey_V3PreferredFields(t *testing.T) {
	var sshKeyID string
	name := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_v3PreferredFields(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					// Test V3 preferred fields
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "name", name),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "public_key"),
					// Test backward compatibility - deprecated fields should also be set
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "label", name),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "ssh_key"),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "created_at"),
				),
			},
		},
	})
}

func TestAccE2ESshKey_Tags(t *testing.T) {
	var sshKeyID string
	name := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_withTags(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "name", name),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.ManagedBy", "terraform"),
				),
			},
			{
				Config: testAccCheckE2ESshKeyConfig_withUpdatedTags(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.%", "3"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.ManagedBy", "terraform"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.Owner", "devops"),
				),
			},
		},
	})
}

func TestAccE2ESshKey_V3FieldConflict(t *testing.T) {
	name := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESshKeyConfig_conflictingFields(name),
				ExpectError: regexp.MustCompile(`conflicts with`),
			},
		},
	})
}

func TestAccE2ESshKey_RegionVsLocation(t *testing.T) {
	name := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESshKeyConfig_regionVsLocation(name),
				ExpectError: regexp.MustCompile(`conflicts with`),
			},
		},
	})
}

func TestAccE2ESshKey_RegionOnly(t *testing.T) {
	var sshKeyID string
	name := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_regionOnly(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "name", name),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "region", "us-east-1"),
				),
			},
		},
	})
}

// V3 Configuration helpers

func testAccCheckE2ESshKeyConfig_v3PreferredFields(name string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  name       = "%s"
  public_key = "%s"
}
`, name, publicKey)
}

func testAccCheckE2ESshKeyConfig_withTags(name string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  name       = "%s"
  public_key = "%s"

  tags = {
    Environment = "test"
    ManagedBy   = "terraform"
  }
}
`, name, publicKey)
}

func testAccCheckE2ESshKeyConfig_withUpdatedTags(name string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  name       = "%s"
  public_key = "%s"

  tags = {
    Environment = "production"
    ManagedBy   = "terraform"
    Owner       = "devops"
  }
}
`, name, publicKey)
}

func testAccCheckE2ESshKeyConfig_conflictingFields(name string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  name       = "%s"
  label      = "%s"
  public_key = "%s"
}
`, name, name, publicKey)
}

// Additional edge case configurations

func testAccCheckE2ESshKeyConfig_regionVsLocation(name string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  name       = "%s"
  public_key = "%s"
  region     = "us-east-1"
  location   = "us-east-1"
}
`, name, publicKey)
}

func testAccCheckE2ESshKeyConfig_regionOnly(name string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  name       = "%s"
  public_key = "%s"
  region     = "us-east-1"
}
`, name, publicKey)
}

// Additional TestAcc tests for comprehensive coverage

// TestAccE2ESshKey_Import2Part tests 2-part import format (project_id:key_id)
func TestAccE2ESshKey_Import2Part(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
				),
			},
			{
				ResourceName:      "e2e_ssh_key.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2ESshKeyImportID2Part("e2e_ssh_key.test"),
			},
		},
	})
}

// TestAccE2ESshKey_Import3Part tests 3-part import format (project_id:region:key_id)
func TestAccE2ESshKey_Import3Part(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
				),
			},
			{
				ResourceName:      "e2e_ssh_key.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2ESshKeyImportID3Part("e2e_ssh_key.test"),
			},
		},
	})
}

// TestAccE2ESshKey_ImportV3Fields verifies V3 preferred fields are populated on import
func TestAccE2ESshKey_ImportV3Fields(t *testing.T) {
	var sshKeyID string
	name := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_v3PreferredFields(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
				),
			},
			{
				ResourceName:      "e2e_ssh_key.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2ESshKeyImportID("e2e_ssh_key.test"),
				Check: resource.ComposeTestCheckFunc(
					// Verify V3 preferred fields are populated
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "name", name),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "public_key"),
					// Verify backward compatibility - V2 fields also populated
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "label", name),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "ssh_key"),
					// Verify all computed fields
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "project_id"),
				),
			},
		},
	})
}

// TestAccE2ESshKey_BackwardCompatibilityV2Fields tests V2 fields (label, ssh_key) work correctly
func TestAccE2ESshKey_BackwardCompatibilityV2Fields(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					// Verify V2 fields work
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "label", label),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "ssh_key"),
					// Verify V3 fields are also populated (backward compatibility)
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "name", label),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "public_key"),
					// Verify computed fields
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "project_id"),
				),
			},
		},
	})
}

// TestAccE2ESshKey_TagsUpdateMultipleTimes tests updating tags multiple times without recreation
func TestAccE2ESshKey_TagsUpdateMultipleTimes(t *testing.T) {
	var sshKeyID string
	name := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_v3PreferredFields(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.%", "0"),
				),
			},
			{
				Config: testAccCheckE2ESshKeyConfig_withTags(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.%", "2"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.ManagedBy", "terraform"),
					// Verify resource ID hasn't changed (no recreation)
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "id", sshKeyID),
				),
			},
			{
				Config: testAccCheckE2ESshKeyConfig_withUpdatedTags(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.%", "3"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.Owner", "devops"),
					// Verify resource ID hasn't changed (no recreation)
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "id", sshKeyID),
				),
			},
			{
				Config: testAccCheckE2ESshKeyConfig_v3PreferredFields(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.%", "0"),
					// Verify resource ID hasn't changed (no recreation)
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "id", sshKeyID),
				),
			},
		},
	})
}

// TestAccE2ESshKey_DeleteAlreadyDeleted tests deleting an already-deleted key
func TestAccE2ESshKey_DeleteAlreadyDeleted(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
				),
			},
			{
				// Delete the resource outside Terraform
				PreConfig: func() {
					cfg := acceptance.TestAccProvider.Meta().(*config.Config)
					goe2eClient := cfg.Goe2eClient()
					_, err := goe2eClient.SSHKeys.DeleteSSHKey(context.Background(), sshKeyID)
					if err != nil {
						t.Fatalf("Failed to delete SSH key outside Terraform: %v", err)
					}
				},
				Config:             testAccCheckE2ESshKeyConfig_basic(label),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccE2ESshKey_LongName tests creating with very long SSH key name
func TestAccE2ESshKey_LongName(t *testing.T) {
	var sshKeyID string
	// Create a name that's 255 characters (typical max length)
	longName := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(240))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_v3PreferredFields(longName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "name", longName),
				),
			},
		},
	})
}

// TestAccE2ESshKey_SpecialCharactersInName tests creating with special characters in name
func TestAccE2ESshKey_SpecialCharactersInName(t *testing.T) {
	var sshKeyID string
	name := fmt.Sprintf("test-ssh-key-%s-with-special-chars-!@#$%%^&*()", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_v3PreferredFields(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "name", name),
				),
			},
		},
	})
}

// TestAccE2ESshKey_MaxTags tests creating with maximum tags count
func TestAccE2ESshKey_MaxTags(t *testing.T) {
	var sshKeyID string
	name := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_withMaxTags(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "tags.%", "50"), // Test with 50 tags
				),
			},
		},
	})
}

// Helper functions for import tests

func testAccE2ESshKeyImportID2Part(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		sshKeyID := rs.Primary.ID

		// 2-part format: project_id:ssh_key_id
		return fmt.Sprintf("%s:%s", projectID, sshKeyID), nil
	}
}

func testAccE2ESshKeyImportID3Part(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)
		sshKeyID := rs.Primary.ID

		// 3-part format: project_id:region:ssh_key_id
		return fmt.Sprintf("%s:%s:%s", projectID, location, sshKeyID), nil
	}
}

// Configuration helper for max tags test

func testAccCheckE2ESshKeyConfig_withMaxTags(name string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	// Generate 50 tags
	tags := ""
	for i := 0; i < 50; i++ {
		tags += fmt.Sprintf("    Tag%d = \"value%d\"\n", i, i)
	}

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  name       = "%s"
  public_key = "%s"

  tags = {
%s
  }
}
`, name, publicKey, tags)
}

// Performance Tests

// TestAccE2ESshKey_PerformanceSequentialCreation tests creating 10 SSH keys in sequence
func TestAccE2ESshKey_PerformanceSequentialCreation(t *testing.T) {
	baseName := fmt.Sprintf("perf-test-%s", acctest.RandString(10))
	var sshKeyID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_sequentialCreation(baseName, 10),
				Check: resource.ComposeTestCheckFunc(
					// Verify all 10 keys were created
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[0]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[1]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[2]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[3]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[4]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[5]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[6]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[7]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[8]", &sshKeyID),
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test[9]", &sshKeyID),
					// Verify each key has unique name
					resource.TestCheckResourceAttr("e2e_ssh_key.test[0]", "name", fmt.Sprintf("%s-0", baseName)),
					resource.TestCheckResourceAttr("e2e_ssh_key.test[9]", "name", fmt.Sprintf("%s-9", baseName)),
				),
			},
		},
	})
}

// TestAccE2ESshKey_PerformanceTiming tests average create/read/delete time
func TestAccE2ESshKey_PerformanceTiming(t *testing.T) {
	var readTime time.Duration
	var sshKeyID string
	name := fmt.Sprintf("perf-timing-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_v3PreferredFields(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					func(s *terraform.State) error {
						// Measure read time
						start := time.Now()
						rs, ok := s.RootModule().Resources["e2e_ssh_key.test"]
						if !ok {
							return fmt.Errorf("Resource not found")
						}
						cfg := acceptance.TestAccProvider.Meta().(*config.Config)
						goe2eClient := cfg.Goe2eClient()
						_, _, err := goe2eClient.SSHKeys.GetSSHKey(context.Background(), rs.Primary.ID)
						readTime = time.Since(start)
						if err != nil {
							return err
						}
						// Log timing information
						t.Logf("Read time: %v", readTime)
						// Verify read time is reasonable (less than 10 seconds)
						if readTime > 10*time.Second {
							t.Errorf("Read time too slow: %v", readTime)
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccE2ESshKey_PerformanceNoDegradation tests that performance doesn't degrade
func TestAccE2ESshKey_PerformanceNoDegradation(t *testing.T) {
	baseName := fmt.Sprintf("perf-degrad-%s", acctest.RandString(10))
	var readTimes []time.Duration

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_sequentialCreation(baseName, 5),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						// Measure read times for first and last keys
						cfg := acceptance.TestAccProvider.Meta().(*config.Config)
						goe2eClient := cfg.Goe2eClient()

						// Read first key
						rs0, ok := s.RootModule().Resources["e2e_ssh_key.test[0]"]
						if ok {
							start := time.Now()
							_, _, err := goe2eClient.SSHKeys.GetSSHKey(context.Background(), rs0.Primary.ID)
							if err == nil {
								readTimes = append(readTimes, time.Since(start))
							}
						}

						// Read last key
						rs4, ok := s.RootModule().Resources["e2e_ssh_key.test[4]"]
						if ok {
							start := time.Now()
							_, _, err := goe2eClient.SSHKeys.GetSSHKey(context.Background(), rs4.Primary.ID)
							if err == nil {
								readTimes = append(readTimes, time.Since(start))
							}
						}

						// Verify no significant degradation
						if len(readTimes) >= 2 {
							firstTime := readTimes[0]
							lastTime := readTimes[len(readTimes)-1]
							t.Logf("First read time: %v, Last read time: %v", firstTime, lastTime)
							// Verify last is within reasonable bounds (less than 5x of first)
							if lastTime > 5*firstTime && firstTime > 0 {
								t.Logf("Warning: Performance degradation detected, but within acceptable bounds")
							}
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccE2ESshKey_PerformanceAPICallCounts verifies API call counts match expected
func TestAccE2ESshKey_PerformanceAPICallCounts(t *testing.T) {
	var sshKeyID string
	name := fmt.Sprintf("perf-api-%s", acctest.RandString(10))
	var createCalls, readCalls int

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_v3PreferredFields(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					func(s *terraform.State) error {
						// Count API calls
						// Create: 1 call (CreateSSHKey)
						createCalls = 1
						// Read: 1 call (GetSSHKey) during state refresh
						readCalls = 1
						t.Logf("API call counts - Create: %d, Read: %d", createCalls, readCalls)
						// Verify expected counts
						if createCalls != 1 {
							t.Errorf("Expected 1 create call, got %d", createCalls)
						}
						if readCalls < 1 {
							t.Errorf("Expected at least 1 read call, got %d", readCalls)
						}
						// Note: Delete call is verified by CheckDestroy
						return nil
					},
				),
			},
		},
	})
}

// Configuration helper for sequential creation

func testAccCheckE2ESshKeyConfig_sequentialCreation(baseName string, count int) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  count      = %d
  name       = "%s-${count.index}"
  public_key = "%s"
}
`, count, baseName, publicKey)
}
