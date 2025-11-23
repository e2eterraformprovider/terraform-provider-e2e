package ssh_key_test

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

func TestAccE2ESshKey_Basic(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
					resource.TestCheckResourceAttr("e2e_ssh_key.test", "label", label),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "ssh_key"),
					resource.TestCheckResourceAttrSet("e2e_ssh_key.test", "timestamp"),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
				Config:      testAccCheckE2ESshKeyConfig_updatedLabel(label),
				ExpectError: regexp.MustCompile(`label cannot be updated once you add the ssh key`),
			},
		},
	})
}

func TestAccE2ESshKey_UpdateSSHKey(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ESshKeyExists("e2e_ssh_key.test", &sshKeyID),
				),
			},
			{
				Config:      testAccCheckE2ESshKeyConfig_updatedKey(label),
				ExpectError: regexp.MustCompile(`ssh_key cannot be updated once you add the ssh key`),
			},
		},
	})
}

func TestAccE2ESshKey_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESshKeyConfig_missingLabel(),
				ExpectError: regexp.MustCompile(`The argument "label" is required`),
			},
			{
				Config:      testAccCheckE2ESshKeyConfig_missingSSHKey(),
				ExpectError: regexp.MustCompile(`The argument "ssh_key" is required`),
			},
			{
				Config:      testAccCheckE2ESshKeyConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2ESshKeyConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccE2ESshKey_InvalidSSHKey(t *testing.T) {
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ESshKeyConfig_invalidKey(label),
				ExpectError: regexp.MustCompile(`.*`),
			},
		},
	})
}

func TestAccE2ESshKey_Import(t *testing.T) {
	var sshKeyID string
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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

func testAccCheckE2ESshKeyExists(resourceName string, sshKeyID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SSH Key ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		sshKey, err := client.GetSshKeyByPk(rs.Primary.ID, projectID, location)
		if err != nil {
			return err
		}

		if sshKey == nil {
			return fmt.Errorf("SSH Key not found")
		}

		*sshKeyID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2ESshKeyDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_ssh_key" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		_, err := client.GetSshKeyByPk(rs.Primary.ID, projectID, location)
		if err == nil {
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
		location := rs.Primary.Attributes["location"]
		sshKeyID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, sshKeyID), nil
	}
}

// Configuration helpers

func testAccCheckE2ESshKeyConfig_basic(label string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  ssh_key    = "%s"
  project_id = "%s"
  location   = "%s"
}
`, label, publicKey, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESshKeyConfig_updatedLabel(oldLabel string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"
	newLabel := fmt.Sprintf("%s-updated", oldLabel)

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  ssh_key    = "%s"
  project_id = "%s"
  location   = "%s"
}
`, newLabel, publicKey, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESshKeyConfig_updatedKey(label string) string {
	newPublicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDifferentKey9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test2@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  ssh_key    = "%s"
  project_id = "%s"
  location   = "%s"
}
`, label, newPublicKey, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESshKeyConfig_invalidKey(label string) string {
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  ssh_key    = "invalid-ssh-key"
  project_id = "%s"
  location   = "%s"
}
`, label, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

// Error case configurations

func testAccCheckE2ESshKeyConfig_missingLabel() string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  ssh_key    = "%s"
  project_id = "%s"
  location   = "%s"
}
`, publicKey, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESshKeyConfig_missingSSHKey() string {
	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "test-label"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESshKeyConfig_missingProjectID() string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label    = "test-label"
  ssh_key  = "%s"
  location = "%s"
}
`, publicKey, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ESshKeyConfig_missingLocation() string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "test-label"
  ssh_key    = "%s"
  project_id = "%s"
}
`, publicKey, os.Getenv("E2E_TEST_PROJECT_ID"))
}
