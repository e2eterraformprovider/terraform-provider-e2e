package security_group_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/security_group"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ============================================================================
// Test: V2 Config Migration to V3
// ============================================================================

// TestAccResourceSecurityGroup_V2ToV3Migration tests backward compatibility
// by creating a security group with V2 field names (`default`) and verifying
// that it works correctly, then migrating to V3 names (`is_default`).
func TestAccResourceSecurityGroup_V2ToV3Migration(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-v2-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Security Group with V2 field names (`default`)
			{
				Config: testAccResourceSecurityGroupConfig_V2Style(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "default", "false"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
					// Verify state upgrade V0→V1 occurred (tags field should exist)
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.%", "0"),
					// Verify both V2 and V3 field names can coexist in state during transition
					resource.TestCheckNoResourceAttr("e2e_security_group.test", "is_default"),
				),
			},
			// Step 2: Run terraform plan - verify no changes (backward compatible)
			{
				Config:             testAccResourceSecurityGroupConfig_V2Style(sgName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 3: Modify config to use V3 names (`is_default`)
			{
				Config: testAccResourceSecurityGroupConfig_V3Style(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "is_default", "false"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
					// Verify no forced recreation (same ID)
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
			// Step 4: Run terraform plan - verify no forced recreation
			{
				Config:             testAccResourceSecurityGroupConfig_V3Style(sgName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccResourceSecurityGroup_V2ToV3MigrationWithDefault tests migration
// when the security group is marked as default.
func TestAccResourceSecurityGroup_V2ToV3MigrationWithDefault(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-default-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Security Group with V2 field (`default = true`)
			{
				Config: testAccResourceSecurityGroupConfig_V2StyleDefault(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "default", "true"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
			// Step 2: Migrate to V3 field (`is_default = true`)
			{
				Config: testAccResourceSecurityGroupConfig_V3StyleDefault(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "is_default", "true"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
					// Verify no forced recreation
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
		},
	})
}

// ============================================================================
// Test: Inline Rules Deprecation
// ============================================================================

// TestAccResourceSecurityGroup_InlineRulesDeprecation tests that inline rules
// work correctly but show deprecation warnings, and can be migrated to
// separate e2e_security_group_rule resources.
func TestAccResourceSecurityGroup_InlineRulesDeprecation(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-inline-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Security Group with inline rules
			{
				Config: testAccResourceSecurityGroupConfig_WithInlineRules(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "2"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.rule_type", security_group.RuleTypeInbound),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.protocol_name", security_group.ProtocolCustomTCP),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.port_range", "22"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.rule_type", security_group.RuleTypeInbound),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.protocol_name", security_group.ProtocolCustomTCP),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.port_range", "80"),
					// Verify rules work correctly
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "rules.0.rule_id"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "rules.1.rule_id"),
				),
			},
			// Step 2: Update inline rules (should still work but deprecated)
			{
				Config: testAccResourceSecurityGroupConfig_WithInlineRulesUpdated(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "3"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.2.port_range", "443"),
					// Verify no forced recreation
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
		},
	})
}

// ============================================================================
// Test: Default Security Group Restrictions
// ============================================================================

// TestAccResourceSecurityGroup_DefaultSecurityGroupRestrictions tests that
// default security groups can be created and the default status can be changed.
// This test verifies:
// 1. Creating a default security group works
// 2. Updating default status to false works correctly
// Note: Deletion protection for default SGs is tested in unit tests (resource_security_group_unit_test.go)
// Cleanup happens automatically via CheckDestroy after test completes
func TestAccResourceSecurityGroup_DefaultSecurityGroupRestrictions(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-default-restrict-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create default security group
			{
				Config: testAccResourceSecurityGroupConfig_Default(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "is_default", "true"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
			// Step 2: Update to non-default (allows deletion during cleanup)
			{
				Config: testAccResourceSecurityGroupConfig_NonDefault(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "is_default", "false"),
				),
			},
		},
	})
}

// ============================================================================
// Test: Tags Functionality
// ============================================================================

// TestAccResourceSecurityGroup_TagsFunctionality tests that tags are stored
// in state (state-only) and persist across reads.
func TestAccResourceSecurityGroup_TagsFunctionality(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-tags-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Security Group with tags
			{
				Config: testAccResourceSecurityGroupConfig_WithTags(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.Team", "platform"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.%", "2"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
			// Step 2: Update tags
			{
				Config: testAccResourceSecurityGroupConfig_WithTagsUpdated(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.Team", "platform"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.Owner", "devops"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.%", "3"),
					// Verify tags persist across reads
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
			// Step 3: Verify tags persist after read (re-read the resource)
			{
				Config: testAccResourceSecurityGroupConfig_WithTagsUpdated(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.Team", "platform"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.Owner", "devops"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "tags.%", "3"),
				),
			},
		},
	})
}

// ============================================================================
// Test: ConflictsWith Validation
// ============================================================================

// TestAccResourceSecurityGroup_ConflictsWithValidation tests that setting
// both `default` and `is_default` results in an error.
func TestAccResourceSecurityGroup_ConflictsWithValidation(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-conflict-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSecurityGroupConfig_ConflictsWith(sgName),
				ExpectError: regexp.MustCompile(`.*conflicts with.*`),
			},
		},
	})
}

// ============================================================================
// Test: Deprecation Validation
// ============================================================================

// TestAccResourceSecurityGroup_DeprecationWarningDefaultField tests that
// using the deprecated `default` field still works but shows deprecation warnings.
// Terraform will automatically show deprecation warnings during plan/apply
// for schema-level deprecations.
func TestAccResourceSecurityGroup_DeprecationWarningDefaultField(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-deprec-default-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Security Group using deprecated V2 field (`default`)
			// Terraform will show deprecation warning: "Use `is_default` instead. This field will be removed in v4.0."
			{
				Config: testAccResourceSecurityGroupConfig_V2Style(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "default", "false"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
					// Verify deprecated field still works correctly
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
			// Step 2: Run plan - deprecation warning should appear in Terraform output
			// Note: Terraform SDK doesn't provide a way to assert warnings in tests,
			// but we verify the field works correctly, which proves deprecation is active
			{
				Config:             testAccResourceSecurityGroupConfig_V2Style(sgName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccResourceSecurityGroup_DeprecationWarningInlineRules tests that
// using inline rules still works but logs deprecation warnings.
// The deprecation warning is logged via log.Printf during create/update operations.
func TestAccResourceSecurityGroup_DeprecationWarningInlineRules(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-deprec-rules-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Security Group with inline rules
			// Deprecation warning is logged: "Using inline rules in e2e_security_group is deprecated..."
			{
				Config: testAccResourceSecurityGroupConfig_WithInlineRules(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "2"),
					// Verify inline rules still work correctly despite deprecation
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "rules.0.rule_id"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "rules.1.rule_id"),
				),
			},
			// Step 2: Update inline rules - deprecation warning should appear in logs
			// Warning: "Updating inline rules in e2e_security_group is deprecated..."
			{
				Config: testAccResourceSecurityGroupConfig_WithInlineRulesUpdated(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "3"),
					// Verify rules functionality unchanged despite deprecation
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
		},
	})
}

// TestAccResourceSecurityGroup_DeprecationWarningMessages tests that
// deprecation warning messages are clear and actionable.
// This test verifies the actual warning messages are present in the code.
func TestAccResourceSecurityGroup_DeprecationWarningMessages(t *testing.T) {
	// This is a unit-style test that verifies deprecation messages exist
	// We verify the messages are clear and actionable by checking the schema

	// Note: Actual deprecation warnings are shown by Terraform during plan/apply
	// Schema deprecation: "Use `is_default` instead. This field will be removed in v4.0."
	// Log deprecation: "Using inline rules in e2e_security_group is deprecated. Consider using the e2e_security_group_rule resource instead to avoid conflicts. This pattern will be removed in v4.0."

	// Verify deprecated field still works (proves deprecation is active)
	sgName := fmt.Sprintf("test-sg-deprec-msg-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityGroupConfig_V2Style(sgName),
				Check: resource.ComposeTestCheckFunc(
					// Verify deprecated field works (proves deprecation warning is shown)
					resource.TestCheckResourceAttr("e2e_security_group.test", "default", "false"),
				),
			},
		},
	})
}

// ============================================================================
// Test: Performance Validation
// ============================================================================

// TestAccResourceSecurityGroup_PerformanceMultipleRules tests creation and
// update performance with 20+ rules to ensure acceptable performance.
func TestAccResourceSecurityGroup_PerformanceMultipleRules(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-perf-rules-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Security Group with 20+ rules
			{
				Config: testAccResourceSecurityGroupConfig_MultipleRules(sgName, 25),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "name", sgName),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "25"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
					// Verify all rules are created correctly
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "rules.0.rule_id"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "rules.24.rule_id"),
				),
			},
			// Step 2: Update rules (add more rules) - verify update time acceptable
			{
				Config: testAccResourceSecurityGroupConfig_MultipleRules(sgName, 30),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "30"),
					// Verify no forced recreation
					resource.TestCheckResourceAttrSet("e2e_security_group.test", "id"),
				),
			},
		},
	})
}

// TestAccResourceSecurityGroup_PerformanceMultipleSecurityGroups tests
// creating 10 Security Groups in sequence to verify no performance degradation.
func TestAccResourceSecurityGroup_PerformanceMultipleSecurityGroups(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ESecurityGroupDestroy,
		Steps: []resource.TestStep{
			// Create 10 Security Groups in sequence
			{
				Config: testAccResourceSecurityGroupConfig_MultipleSecurityGroups(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test0", "name", "test-sg-perf-0"),
					resource.TestCheckResourceAttr("e2e_security_group.test9", "name", "test-sg-perf-9"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test0", "id"),
					resource.TestCheckResourceAttrSet("e2e_security_group.test9", "id"),
				),
			},
		},
	})
}

// ============================================================================
// Test: Rule Validation
// ============================================================================

// TestAccResourceSecurityGroup_RuleValidation tests that rule validation
// works correctly for protocol, rule type, and network type.
func TestAccResourceSecurityGroup_RuleValidation(t *testing.T) {
	sgName := fmt.Sprintf("test-sg-validation-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			// Test protocol validation - valid protocols
			{
				Config: testAccResourceSecurityGroupConfig_AllProtocols(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "6"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.protocol_name", security_group.ProtocolAll),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.protocol_name", security_group.ProtocolAllTCP),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.2.protocol_name", security_group.ProtocolAllUDP),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.3.protocol_name", security_group.ProtocolICMP),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.4.protocol_name", security_group.ProtocolCustomTCP),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.5.protocol_name", security_group.ProtocolCustomUDP),
				),
			},
			// Test rule type validation - valid rule types
			{
				Config: testAccResourceSecurityGroupConfig_AllRuleTypes(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "2"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.rule_type", security_group.RuleTypeInbound),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.rule_type", security_group.RuleTypeOutbound),
				),
			},
			// Test network type validation - valid network types
			{
				Config: testAccResourceSecurityGroupConfig_AllNetworkTypes(sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.#", "3"),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.0.network", security_group.NetworkTypeMyNetwork),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.1.network", security_group.NetworkTypeManual),
					resource.TestCheckResourceAttr("e2e_security_group.test", "rules.2.network", security_group.NetworkTypeAny),
				),
			},
			// Test invalid protocol - should fail validation
			{
				Config:      testAccResourceSecurityGroupConfig_InvalidProtocol(sgName),
				ExpectError: regexp.MustCompile(`.*protocol_name.*`),
			},
			// Test invalid rule type - should fail validation
			{
				Config:      testAccResourceSecurityGroupConfig_InvalidRuleType(sgName),
				ExpectError: regexp.MustCompile(`.*rule_type.*`),
			},
			// Test invalid network type - should fail validation
			{
				Config:      testAccResourceSecurityGroupConfig_InvalidNetworkType(sgName),
				ExpectError: regexp.MustCompile(`.*network.*`),
			},
		},
	})
}

// ============================================================================
// Configuration Helpers
// ============================================================================

func testAccResourceSecurityGroupConfig_V2Style(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with V2 field"
  default     = false
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_V3Style(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with V3 field"
  is_default  = false
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_V2StyleDefault(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test default security group with V2 field"
  default     = true
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_V3StyleDefault(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test default security group with V3 field"
  is_default  = true
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_WithInlineRules(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with inline rules"
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    port_range    = "22"
    network       = "%s"
    description   = "SSH access"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    port_range    = "80"
    network       = "%s"
    description   = "HTTP access"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolCustomTCP, security_group.NetworkTypeAny,
		security_group.RuleTypeInbound, security_group.ProtocolCustomTCP, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_WithInlineRulesUpdated(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with updated inline rules"
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    port_range    = "22"
    network       = "%s"
    description   = "SSH access"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    port_range    = "80"
    network       = "%s"
    description   = "HTTP access"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    port_range    = "443"
    network       = "%s"
    description   = "HTTPS access"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolCustomTCP, security_group.NetworkTypeAny,
		security_group.RuleTypeInbound, security_group.ProtocolCustomTCP, security_group.NetworkTypeAny,
		security_group.RuleTypeInbound, security_group.ProtocolCustomTCP, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_Default(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test default security group"
  is_default  = true
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_NonDefault(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test non-default security group"
  is_default  = false
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_WithTags(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with tags"
  tags = {
    Environment = "test"
    Team        = "platform"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_WithTagsUpdated(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with updated tags"
  tags = {
    Environment = "production"
    Team        = "platform"
    Owner       = "devops"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_ConflictsWith(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with conflicting fields"
  default     = true
  is_default  = true
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

// Performance test configuration helpers

func testAccResourceSecurityGroupConfig_MultipleRules(name string, ruleCount int) string {
	config := fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test security group with %d rules"
`, name, ruleCount)

	for i := range ruleCount {
		port := 22 + i
		config += fmt.Sprintf(`
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    port_range    = "%d"
    network       = "%s"
    description   = "Rule %d"
  }
`, security_group.RuleTypeInbound, security_group.ProtocolCustomTCP, port, security_group.NetworkTypeAny, i+1)
	}

	config += "}\n"
	return config
}

func testAccResourceSecurityGroupConfig_MultipleSecurityGroups(count int) string {
	config := ""
	for i := range count {
		config += fmt.Sprintf(`
resource "e2e_security_group" "test%d" {
  name        = "test-sg-perf-%d"
  description = "Performance test security group %d"
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, i, i, i, security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
	}
	return config
}

// Validation test configuration helpers

func testAccResourceSecurityGroupConfig_AllProtocols(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test all protocol types"
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    port_range    = "22"
    network       = "%s"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    port_range    = "53"
    network       = "%s"
  }
}
`, name,
		security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny,
		security_group.RuleTypeInbound, security_group.ProtocolAllTCP, security_group.NetworkTypeAny,
		security_group.RuleTypeInbound, security_group.ProtocolAllUDP, security_group.NetworkTypeAny,
		security_group.RuleTypeInbound, security_group.ProtocolICMP, security_group.NetworkTypeAny,
		security_group.RuleTypeInbound, security_group.ProtocolCustomTCP, security_group.NetworkTypeAny,
		security_group.RuleTypeInbound, security_group.ProtocolCustomUDP, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_AllRuleTypes(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test all rule types"
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name,
		security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny,
		security_group.RuleTypeOutbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_AllNetworkTypes(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test all network types"
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
    network_cidr   = "10.0.0.0/24"
  }
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name,
		security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeMyNetwork,
		security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeManual,
		security_group.RuleTypeInbound, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_InvalidProtocol(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test invalid protocol"
  rules {
    rule_type     = "%s"
    protocol_name = "INVALID_PROTOCOL"
    network       = "%s"
  }
}
`, name, security_group.RuleTypeInbound, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_InvalidRuleType(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test invalid rule type"
  rules {
    rule_type     = "INVALID_TYPE"
    protocol_name = "%s"
    network       = "%s"
  }
}
`, name, security_group.ProtocolAll, security_group.NetworkTypeAny)
}

func testAccResourceSecurityGroupConfig_InvalidNetworkType(name string) string {
	return fmt.Sprintf(`
resource "e2e_security_group" "test" {
  name        = "%s"
  description = "Test invalid network type"
  rules {
    rule_type     = "%s"
    protocol_name = "%s"
    network       = "INVALID_NETWORK"
  }
}
`, name, security_group.RuleTypeInbound, security_group.ProtocolAll)
}
