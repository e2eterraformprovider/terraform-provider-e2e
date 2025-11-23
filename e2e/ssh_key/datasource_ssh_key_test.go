package ssh_key_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceE2ESshKey_Basic(t *testing.T) {
	label := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceE2ESshKeyConfig_basic(label),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_ssh_key.test", "label", label),
					resource.TestCheckResourceAttrSet("data.e2e_ssh_key.test", "ssh_key"),
					resource.TestCheckResourceAttrSet("data.e2e_ssh_key.test", "timestamp"),
					resource.TestCheckResourceAttrSet("data.e2e_ssh_key.test", "project_name"),
					resource.TestCheckResourceAttrPair(
						"data.e2e_ssh_key.test", "id",
						"e2e_ssh_key.test", "id",
					),
				),
			},
		},
	})
}

func TestAccDataSourceE2ESshKey_NonExistent(t *testing.T) {
	label := fmt.Sprintf("non-existent-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataSourceE2ESshKeyConfig_nonExistent(label),
				ExpectError: regexp.MustCompile(`error finding ssh key with label`),
			},
		},
	})
}

func TestAccDataSourceE2ESshKeys_List(t *testing.T) {
	label1 := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))
	label2 := fmt.Sprintf("test-ssh-key-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceE2ESshKeysConfig_list(label1, label2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_ssh_keys.test", "ssh_key_list.#"),
					resource.TestCheckResourceAttrSet("data.e2e_ssh_keys.test", "ssh_key_list.0.pk"),
					resource.TestCheckResourceAttrSet("data.e2e_ssh_keys.test", "ssh_key_list.0.label"),
					resource.TestCheckResourceAttrSet("data.e2e_ssh_keys.test", "ssh_key_list.0.ssh_key"),
					resource.TestCheckResourceAttrSet("data.e2e_ssh_keys.test", "ssh_key_list.0.timestamp"),
				),
			},
		},
	})
}

func TestAccDataSourceE2ESshKey_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataSourceE2ESshKeyConfig_missingLabel(),
				ExpectError: regexp.MustCompile(`The argument "label" is required`),
			},
			{
				Config:      testAccCheckDataSourceE2ESshKeyConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckDataSourceE2ESshKeyConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

// Configuration helpers for data source tests

func testAccCheckDataSourceE2ESshKeyConfig_basic(label string) string {
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test" {
  label      = "%s"
  ssh_key    = "%s"
  project_id = "%s"
  location   = "%s"
}

data "e2e_ssh_key" "test" {
  label      = e2e_ssh_key.test.label
  project_id = "%s"
  location   = "%s"
}
`, label, publicKey, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckDataSourceE2ESshKeyConfig_nonExistent(label string) string {
	return fmt.Sprintf(`
data "e2e_ssh_key" "test" {
  label      = "%s"
  project_id = "%s"
  location   = "%s"
}
`, label, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckDataSourceE2ESshKeysConfig_list(label1, label2 string) string {
	publicKey1 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcbi1cXHVf9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test1@example.com"
	publicKey2 := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDifferentKey9aJxdQTPwmBce/dL7eLEb7NWUoGZ9ZJ7YKhCZ1DD+UM8wR48eGFw24q5V4L3T6T2EbAm1y9ByC1E0Cn/Vn8T4X3d3KZW7VDkP8xKdZO2Y7ePJJJyNLpU8VabCxI3PbL3EkT2aTCKtU/yLlGqYLfzQEO/T9E2eSqQMqIqvjEQnWDCQQfFxHvVZxQP5s2qJaF9P3cH4VbS4v5pJ0NJrS8Iv8OZaCP4LkqPFXq4T3qZ8MJT0XbY2J9KMQ5wY8TyT3X8pMJ3cVnU9fT8XqJ3V2nT8X5T2kTqT3X8V2nT8X5T2kTqT3X8V2nT8X5T2k test2@example.com"

	return fmt.Sprintf(`
resource "e2e_ssh_key" "test1" {
  label      = "%s"
  ssh_key    = "%s"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_ssh_key" "test2" {
  label      = "%s"
  ssh_key    = "%s"
  project_id = "%s"
  location   = "%s"
}

data "e2e_ssh_keys" "test" {
  project_id = "%s"
  location   = "%s"
  depends_on = [e2e_ssh_key.test1, e2e_ssh_key.test2]
}
`, label1, publicKey1, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		label2, publicKey2, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckDataSourceE2ESshKeyConfig_missingLabel() string {
	return fmt.Sprintf(`
data "e2e_ssh_key" "test" {
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckDataSourceE2ESshKeyConfig_missingProjectID() string {
	return fmt.Sprintf(`
data "e2e_ssh_key" "test" {
  label    = "test-label"
  location = "%s"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckDataSourceE2ESshKeyConfig_missingLocation() string {
	return fmt.Sprintf(`
data "e2e_ssh_key" "test" {
  label      = "test-label"
  project_id = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}
