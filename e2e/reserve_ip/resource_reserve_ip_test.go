package reserve_ip_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EReserveIP_Basic(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "ip_address"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "bought_at"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "reserve_id"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "reserved_type"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "project_name"),
					resource.TestCheckResourceAttr("e2e_reserve_ip.test", "project_id", os.Getenv("E2E_TEST_PROJECT_ID")),
					resource.TestCheckResourceAttr("e2e_reserve_ip.test", "location", os.Getenv("E2E_TEST_LOCATION")),
				),
			},
		},
	})
}

func TestAccE2EReserveIP_StatusCheck(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "status"),
				),
			},
		},
	})
}

func TestAccE2EReserveIP_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EReserveIPConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EReserveIPConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccE2EReserveIP_Import(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
				),
			},
			{
				ResourceName:      "e2e_reserve_ip.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2EReserveIPImportID("e2e_reserve_ip.test"),
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

func testAccCheckE2EReserveIPExists(resourceName string, reserveIPID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Reserve IP ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]
		ipAddress := rs.Primary.Attributes["ip_address"]

		res, err := client.GetReservedIps(projectID, location)
		if err != nil {
			return err
		}

		found := false
		for _, item := range res.Data {
			if item.IPAddress == ipAddress {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Reserve IP not found")
		}

		*reserveIPID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EReserveIPDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_reserve_ip" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]
		ipAddress := rs.Primary.Attributes["ip_address"]

		res, err := client.GetReservedIps(projectID, location)
		if err != nil {
			// If we get an error, assume it's because the resource doesn't exist
			continue
		}

		for _, item := range res.Data {
			if item.IPAddress == ipAddress {
				return fmt.Errorf("Reserve IP still exists: %s", ipAddress)
			}
		}
	}

	return nil
}

func testAccE2EReserveIPImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]
		reserveIPID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, reserveIPID), nil
	}
}

// Configuration helpers

func testAccCheckE2EReserveIPConfig_basic() string {
	return fmt.Sprintf(`
resource "e2e_reserve_ip" "test" {
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

// Error case configurations

func testAccCheckE2EReserveIPConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_reserve_ip" "test" {
  location = "%s"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EReserveIPConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_reserve_ip" "test" {
  project_id = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}
