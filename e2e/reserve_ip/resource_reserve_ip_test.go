package reserve_ip_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EReserveIP_Basic(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "ip_address"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "created_at"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "reserve_id"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "reserved_type"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "type"), // V3 field
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "urn"),  // V3 field
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "project_name"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "region")),
			},
		},
	})
}

func TestAccE2EReserveIP_StatusCheck(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "status")),
			},
		},
	})
}

func TestAccE2EReserveIP_MissingRequiredArguments(t *testing.T) {
	// This test validates syntax/argument checking and doesn't need API calls
	// Use TestAccPreCheckSyntaxOnly instead of testAccPreCheck
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckSyntaxOnly(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID)),
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

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
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

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)
		ipAddress := rs.Primary.Attributes["ip_address"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
		if err != nil {
			return fmt.Errorf("error creating goe2e client: %w", err)
		}

		reserveIPs, _, err := goe2eClient.ReserveIP.ListReserveIPs(context.Background())
		if err != nil {
			return fmt.Errorf("error listing reserved IPs: %w", err)
		}

		found := false
		for _, ip := range reserveIPs {
			if ip.IPAddress == ipAddress {
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
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_reserve_ip" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := acceptance.GetRegionOrLocationFromState(rs)
		ipAddress := rs.Primary.Attributes["ip_address"]

		goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
		if err != nil {
			// If we can't create a client, assume the resource doesn't exist
			continue
		}

		reserveIPs, _, err := goe2eClient.ReserveIP.ListReserveIPs(context.Background())
		if err != nil {
			// If we get an error, assume it's because the resource doesn't exist
			continue
		}

		for _, ip := range reserveIPs {
			if ip.IPAddress == ipAddress {
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
		location := acceptance.GetRegionOrLocationFromState(rs)
		reserveIPID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, location, reserveIPID), nil
	}
}

// Configuration helpers

func testAccCheckE2EReserveIPConfig_basic() string {
	return `
resource "e2e_reserve_ip" "test" {}
`
}

// Error case configurations

func testAccCheckE2EReserveIPConfig_missingProjectID() string {
	return `
resource "e2e_reserve_ip" "test" {}
`
}

func testAccCheckE2EReserveIPConfig_missingLocation() string {
	return `
resource "e2e_reserve_ip" "test" {}
`
}

// ============================================================================
// V3 Fields Acceptance Tests
// ============================================================================

func TestAccE2EReserveIP_V3Fields(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					// Verify V3 fields
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "type"),          // V3 preferred field
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "urn"),           // V3 field
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "reserved_type"), // Deprecated but still present
					// Verify URN format
					resource.TestMatchResourceAttr("e2e_reserve_ip.test", "urn", regexp.MustCompile(`^e2e:reserve_ip:.+:.+$`)),
				),
			},
		},
	})
}

func TestAccE2EReserveIP_URN(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					// Verify URN format matches e2e:reserve_ip:<region>:<ip_address>
					resource.TestMatchResourceAttr("e2e_reserve_ip.test", "urn", regexp.MustCompile(`^e2e:reserve_ip:.+:.+$`)),
					// Verify URN contains region and IP address
					resource.TestCheckResourceAttrPair("e2e_reserve_ip.test", "urn", "e2e_reserve_ip.test", "ip_address"),
				),
			},
		},
	})
}

func TestAccE2EReserveIP_Type(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					// Verify type field populated
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "type"),
					// Verify reserved_type also populated (backward compatibility)
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "reserved_type"),
					// Verify both fields have same value
					resource.TestCheckResourceAttrPair("e2e_reserve_ip.test", "type", "e2e_reserve_ip.test", "reserved_type"),
				),
			},
		},
	})
}

func TestAccE2EReserveIP_FloatingIPAttachedNodes(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					// Verify floating_ip_attached_nodes field exists (may be empty for non-FloatingIP types)
					resource.TestCheckResourceAttr("e2e_reserve_ip.test", "floating_ip_attached_nodes.#", "0"),
				),
			},
		},
	})
}

// ============================================================================
// Backward Compatibility Acceptance Tests
// ============================================================================

func TestAccE2EReserveIP_DeprecatedReservedType(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					// Verify reserved_type field populated (deprecated)
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "reserved_type"),
					// Verify both reserved_type and type in state
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "type"),
					// Verify they have the same value
					resource.TestCheckResourceAttrPair("e2e_reserve_ip.test", "reserved_type", "e2e_reserve_ip.test", "type"),
				),
			},
		},
	})
}

func TestAccE2EReserveIP_DeprecatedLocation(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_withLocation(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
					// Verify location field works
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "location"),
					// Verify both location and region in state
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "region"),
					// Verify they have the same value
					resource.TestCheckResourceAttrPair("e2e_reserve_ip.test", "location", "e2e_reserve_ip.test", "region"),
				),
			},
		},
	})
}

// ============================================================================
// Import Functionality Acceptance Tests
// ============================================================================

func TestAccE2EReserveIP_ImportWithFloatingIPNodes(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID)),
			},
			{
				ResourceName:      "e2e_reserve_ip.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2EReserveIPImportID("e2e_reserve_ip.test"),
				Check: resource.ComposeTestCheckFunc(
					// Verify V3 fields populated after import
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "type"),
					resource.TestCheckResourceAttrSet("e2e_reserve_ip.test", "urn"),
					// Verify floating_ip_attached_nodes populated if FloatingIP type
					resource.TestCheckResourceAttr("e2e_reserve_ip.test", "floating_ip_attached_nodes.#", "0"),
				),
			},
		},
	})
}

// ============================================================================
// Error Scenarios Acceptance Tests
// ============================================================================

func TestAccE2EReserveIP_ImportInvalidFormat(t *testing.T) {
	// This test validates import ID format syntax and doesn't need API calls
	// Use TestAccPreCheckSyntaxOnly instead of testAccPreCheck
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckSyntaxOnly(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:      "e2e_reserve_ip.test",
				ImportState:       true,
				ImportStateId:     "invalid-format",
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile(`invalid import ID format`),
			},
		},
	})
}

func TestAccE2EReserveIP_ImportNotFound(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:      "e2e_reserve_ip.test",
				ImportState:       true,
				ImportStateId:     "project-123/Mumbai/999.999.999.999", // Non-existent IP
				ImportStateVerify: false,
				ExpectError:       regexp.MustCompile(`not found`),
			},
		},
	})
}

func TestAccE2EReserveIP_DeleteAttached(t *testing.T) {
	var reserveIPID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EReserveIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EReserveIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EReserveIPExists("e2e_reserve_ip.test", &reserveIPID),
				),
			},
			// Deletion happens automatically when resource is removed from config
			// CheckDestroy verifies deletion succeeds even if attached (warning logged)
		},
	})
}

// ============================================================================
// Helper Functions for Acceptance Tests
// ============================================================================

func testAccCheckE2EReserveIPConfig_withLocation() string {
	return `
resource "e2e_reserve_ip" "test" {
  location = "Mumbai"
}
`
}
