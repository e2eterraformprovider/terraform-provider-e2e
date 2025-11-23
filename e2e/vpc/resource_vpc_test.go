package vpc_test

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

func TestAccE2EVPC_Basic(t *testing.T) {
	var vpcID string
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCConfig_basic(vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVPCExists("e2e_vpc.test", &vpcID),
					resource.TestCheckResourceAttr("e2e_vpc.test", "vpc_name", vpcName),
					resource.TestCheckResourceAttr("e2e_vpc.test", "is_e2e_vpc", "true"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", "network_id"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", "state"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", "ipv4_cidr"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", "gateway_ip"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", "pool_size"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test", "is_active"),
				),
			},
		},
	})
}

func TestAccE2EVPC_CustomCIDR(t *testing.T) {
	var vpcID string
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCConfig_customCIDR(vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVPCExists("e2e_vpc.test", &vpcID),
					resource.TestCheckResourceAttr("e2e_vpc.test", "vpc_name", vpcName),
					resource.TestCheckResourceAttr("e2e_vpc.test", "is_e2e_vpc", "false"),
					resource.TestCheckResourceAttr("e2e_vpc.test", "ipv4", "10.0.0.0/24"),
				),
			},
		},
	})
}

func TestAccE2EVPC_Update(t *testing.T) {
	var vpcID string
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))
	vpcNameUpdated := fmt.Sprintf("test-vpc-updated-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCConfig_basic(vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVPCExists("e2e_vpc.test", &vpcID),
					resource.TestCheckResourceAttr("e2e_vpc.test", "vpc_name", vpcName),
				),
			},
			{
				Config:      testAccCheckE2EVPCConfig_basic(vpcNameUpdated),
				ExpectError: regexp.MustCompile(`vpc_name cannot be updated`),
			},
		},
	})
}

func TestAccE2EVPC_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EVPCConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "vpc_name" is required`),
			},
			{
				Config:      testAccCheckE2EVPCConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EVPCConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccE2EVPC_Import(t *testing.T) {
	var vpcID string
	vpcName := fmt.Sprintf("test-vpc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCConfig_basic(vpcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EVPCExists("e2e_vpc.test", &vpcID),
				),
			},
			{
				ResourceName:      "e2e_vpc.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2EVPCImportID("e2e_vpc.test"),
			},
		},
	})
}

func TestAccE2EVPC_MultipleVPCs(t *testing.T) {
	vpcName1 := fmt.Sprintf("test-vpc-1-%s", acctest.RandString(10))
	vpcName2 := fmt.Sprintf("test-vpc-2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EVPCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EVPCConfig_multiple(vpcName1, vpcName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_vpc.test1", "vpc_name", vpcName1),
					resource.TestCheckResourceAttr("e2e_vpc.test2", "vpc_name", vpcName2),
					resource.TestCheckResourceAttrSet("e2e_vpc.test1", "network_id"),
					resource.TestCheckResourceAttrSet("e2e_vpc.test2", "network_id"),
				),
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

func testAccCheckE2EVPCExists(resourceName string, vpcID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Client()

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		vpc, err := client.GetVpc(rs.Primary.ID, projectID, location)
		if err != nil {
			return err
		}

		if vpc == nil {
			return fmt.Errorf("VPC not found")
		}

		*vpcID = rs.Primary.ID
		return nil
	}
}

func testAccCheckE2EVPCDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_vpc" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		_, err := client.GetVpc(rs.Primary.ID, projectID, location)
		if err == nil {
			return fmt.Errorf("VPC still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccE2EVPCImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]
		vpcID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, vpcID), nil
	}
}

// Configuration helpers

func testAccCheckE2EVPCConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  vpc_name   = "%s"
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EVPCConfig_customCIDR(name string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  vpc_name   = "%s"
  is_e2e_vpc = false
  ipv4       = "10.0.0.0/24"
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EVPCConfig_multiple(name1, name2 string) string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test1" {
  vpc_name   = "%s"
  project_id = "%s"
  location   = "%s"
}

resource "e2e_vpc" "test2" {
  vpc_name   = "%s"
  project_id = "%s"
  location   = "%s"
}
`, name1, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		name2, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

// Error case configurations

func testAccCheckE2EVPCConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EVPCConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  vpc_name = "test-vpc"
  location = "%s"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EVPCConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_vpc" "test" {
  vpc_name   = "test-vpc"
  project_id = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}
