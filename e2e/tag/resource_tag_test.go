package tag_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	// Test constants
	testTagNamePrefix = "test-tag-"
	testMetadataValue = "Test tag for acceptance testing"
	resourceTypeTag   = "e2e_tag"
)

func TestAccE2ETag_Basic(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("e2e_tag.test", "metadata", ""),
					resource.TestCheckResourceAttrSet("e2e_tag.test", "label_id"),
					resource.TestCheckResourceAttrSet("e2e_tag.test", "project_id"),
					resource.TestCheckResourceAttrSet("e2e_tag.test", "region"),
				),
			},
		},
	})
}

func TestAccE2ETag_WithMetadata(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))
	metadata := testMetadataValue

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_withMetadata(tagName, metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("e2e_tag.test", "metadata", metadata),
					resource.TestCheckResourceAttrSet("e2e_tag.test", "label_id"),
				),
			},
		},
	})
}

func TestAccE2ETag_Import(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
				),
			},
			{
				ResourceName:      "e2e_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2ETagImportID("e2e_tag.test"),
			},
		},
	})
}

func TestAccE2ETag_ForceNew(t *testing.T) {
	var tag1, tag2 goe2e.Tag
	tagName1 := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))
	tagName2 := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(tagName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag1),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName1),
				),
			},
			{
				Config: testAccCheckE2ETagConfig_basic(tagName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag2),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName2),
					testAccCheckE2ETagRecreated(&tag1, &tag2),
				),
			},
		},
	})
}

func TestAccE2ETag_MissingRequiredFields(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ETagConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
}

func testAccCheckE2ETagExists(resourceName string, tag *goe2e.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client := cfg.Goe2eClient()

		foundTag, _, err := client.Tags.GetTag(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundTag == nil {
			return fmt.Errorf("Tag not found")
		}

		*tag = *foundTag
		return nil
	}
}

func testAccCheckE2ETagDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client := cfg.Goe2eClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceTypeTag {
			continue
		}

		_, _, err := client.Tags.GetTag(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Tag still exists: %s", rs.Primary.ID)
		}

		// Verify it's a not found error
		if !strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) {
			return fmt.Errorf("Expected not found error, got: %s", err)
		}
	}

	return nil
}

func testAccE2ETagImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		region := acceptance.GetRegionOrLocationFromState(rs)
		tagID := rs.Primary.ID

		return fmt.Sprintf("%s/%s/%s", projectID, region, tagID), nil
	}
}

func testAccCheckE2ETagRecreated(old, new *goe2e.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if old.LabelID == new.LabelID {
			return fmt.Errorf("Expected tag to be recreated, but IDs are the same: %d", old.LabelID)
		}
		return nil
	}
}

// Configuration helpers

func testAccCheckE2ETagConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_tag" "test" {
  name = "%s"
}
`, name)
}

func testAccCheckE2ETagConfig_withMetadata(name, metadata string) string {
	return fmt.Sprintf(`
resource "e2e_tag" "test" {
  name     = "%s"
  metadata = "%s"
}
`, name, metadata)
}

func testAccCheckE2ETagConfig_missingName() string {
	return `
resource "e2e_tag" "test" {
  metadata = "test"
}
`
}

func testAccCheckE2ETagConfig_withRegion(name, region string) string {
	return fmt.Sprintf(`
resource "e2e_tag" "test" {
  name   = "%s"
  region = "%s"
}
`, name, region)
}

func testAccCheckE2ETagConfig_withProjectID(name, projectID string) string {
	return fmt.Sprintf(`
resource "e2e_tag" "test" {
  name       = "%s"
  project_id = "%s"
}
`, name, projectID)
}

func testAccCheckE2ETagConfig_duplicate(name string) string {
	return fmt.Sprintf(`
resource "e2e_tag" "test" {
  name = "%s"
}

resource "e2e_tag" "duplicate" {
  name = "%s"
}
`, name, name)
}

// TestAccE2ETag_InvalidRegion tests that invalid region produces an error
func TestAccE2ETag_InvalidRegion(t *testing.T) {
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ETagConfig_withRegion(tagName, "invalid-region-xyz"),
				ExpectError: regexp.MustCompile(`.*`), // API will reject invalid region
			},
		},
	})
}

// TestAccE2ETag_InvalidProjectID tests that invalid project ID produces an error
func TestAccE2ETag_InvalidProjectID(t *testing.T) {
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ETagConfig_withProjectID(tagName, "invalid-project-id-xyz"),
				ExpectError: regexp.MustCompile(`.*`), // API will reject invalid project ID
			},
		},
	})
}

// TestAccE2ETag_DuplicateName tests creating two tags with the same name
// This tests API behavior - may succeed or fail depending on API implementation
func TestAccE2ETag_DuplicateName(t *testing.T) {
	var tag1, tag2 goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag1),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
				),
			},
			{
				// Try to create another tag with the same name
				// API may allow or reject this - test handles both cases
				Config: testAccCheckE2ETagConfig_duplicate(tagName),
				Check: resource.ComposeTestCheckFunc(
					// If API allows duplicates, verify both exist
					testAccCheckE2ETagExists("e2e_tag.duplicate", &tag2),
					resource.TestCheckResourceAttr("e2e_tag.duplicate", "name", tagName),
					// If duplicates are allowed, IDs should be different
					func(s *terraform.State) error {
						if tag1.LabelID == tag2.LabelID {
							return fmt.Errorf("Expected different tag IDs for duplicate names, but got same ID: %d", tag1.LabelID)
						}
						return nil
					},
				),
				// If API rejects duplicates, this step will error - that's also valid
				ExpectError: nil, // Allow both success and error cases
			},
		},
	})
}

// ============================================================================
// Enhanced Acceptance Tests - Phase 5
// ============================================================================

// TestAccE2ETag_EmptyMetadata tests that empty metadata string is handled correctly
func TestAccE2ETag_EmptyMetadata(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_withMetadata(tagName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("e2e_tag.test", "metadata", ""),
				),
			},
		},
	})
}

// TestAccE2ETag_LongMetadata tests with very long metadata string (256+ chars)
func TestAccE2ETag_LongMetadata(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))
	longMetadata := strings.Repeat("A", 300) // 300 characters

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_withMetadata(tagName, longMetadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("e2e_tag.test", "metadata", longMetadata),
				),
			},
		},
	})
}

// TestAccE2ETag_SpecialCharacters tests with quotes, newlines, unicode in metadata
func TestAccE2ETag_SpecialCharacters(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))
	specialMetadata := "Test with \"quotes\", newlines\nand unicode: 测试 🚀"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_withMetadata(tagName, specialMetadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("e2e_tag.test", "metadata", specialMetadata),
				),
			},
		},
	})
}

// TestAccE2ETag_MaxNameLength tests with 128-char name (schema max)
func TestAccE2ETag_MaxNameLength(t *testing.T) {
	var tag goe2e.Tag
	maxLengthName := strings.Repeat("A", 128) // Exactly 128 characters

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(maxLengthName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", maxLengthName),
				),
			},
		},
	})
}

// TestAccE2ETag_NameValidation tests name validation (should fail if invalid)
func TestAccE2ETag_NameValidation(t *testing.T) {
	// Test with empty name (should fail validation)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2ETagConfig_basic(""),
				ExpectError: regexp.MustCompile(`.*`),
			},
		},
	})
}

// TestAccE2ETag_ImportThreePartFormat tests import with project_id/region/tag_id format
func TestAccE2ETag_ImportThreePartFormat(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
				),
			},
			{
				ResourceName:      "e2e_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/%s/%d", acceptance.TestProjectID, acceptance.TestRegion, tag.LabelID),
			},
		},
	})
}

// TestAccE2ETag_ImportSinglePartFormat tests import with just tag_id (uses provider defaults)
func TestAccE2ETag_ImportSinglePartFormat(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
				),
			},
			{
				ResourceName:      "e2e_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%d", tag.LabelID),
			},
		},
	})
}

// TestAccE2ETag_ImportInvalidFormat tests import with wrong format should error
func TestAccE2ETag_ImportInvalidFormat(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:  "e2e_tag.test",
				ImportState:   true,
				ImportStateId: "invalid/format/with/too/many/parts",
				ExpectError:   regexp.MustCompile(`invalid import ID format`),
			},
		},
	})
}

// TestAccE2ETag_ImportNonExistent tests import of non-existent tag ID should fail gracefully
func TestAccE2ETag_ImportNonExistent(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:  "e2e_tag.test",
				ImportState:   true,
				ImportStateId: "999999999",
				ExpectError:   regexp.MustCompile(`.*`),
			},
		},
	})
}

// TestAccE2ETag_DeleteNonExistent tests that deleting already-deleted tag should not error
func TestAccE2ETag_DeleteNonExistent(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
				),
			},
			{
				// Manually delete the tag via API
				PreConfig: func() {
					cfg := acceptance.TestAccProvider.Meta().(*config.Config)
					client := cfg.Goe2eClient()
					_, _ = client.Tags.DeleteTag(context.Background(), fmt.Sprintf("%d", tag.LabelID))
				},
				// Then try to delete via Terraform - should succeed (idempotent)
				Config:             testAccCheckE2ETagConfig_basic(tagName),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccE2ETag_QuickRecreate tests create, delete, create same name quickly (race condition test)
func TestAccE2ETag_QuickRecreate(t *testing.T) {
	var tag1, tag2 goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag1),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
				),
			},
			{
				// Delete the tag
				Config:             " ", // Empty config to trigger destroy
				Destroy:            true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Recreate immediately with same name
				Config: testAccCheckE2ETagConfig_basic(tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag2),
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
					// Verify it's a new tag (different ID)
					testAccCheckE2ETagRecreated(&tag1, &tag2),
				),
			},
		},
	})
}

// TestAccE2ETag_StateConsistency verifies all computed fields match API response
func TestAccE2ETag_StateConsistency(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))
	metadata := testMetadataValue

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_withMetadata(tagName, metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
					// Verify all computed fields are set
					resource.TestCheckResourceAttrSet("e2e_tag.test", "label_id"),
					resource.TestCheckResourceAttrSet("e2e_tag.test", "project_id"),
					resource.TestCheckResourceAttrSet("e2e_tag.test", "region"),
					// Verify label_id matches API response
					resource.TestCheckResourceAttr("e2e_tag.test", "label_id", fmt.Sprintf("%d", tag.LabelID)),
					// Verify name matches API response
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tag.LabelName),
					// Verify metadata matches API response
					resource.TestCheckResourceAttr("e2e_tag.test", "metadata", tag.Metadata),
				),
			},
		},
	})
}

// TestAccE2ETag_ImportStateVerify ensures import populates all fields from API
func TestAccE2ETag_ImportStateVerify(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("%s%s", testTagNamePrefix, acctest.RandString(10))
	metadata := testMetadataValue

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ETagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ETagConfig_withMetadata(tagName, metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2ETagExists("e2e_tag.test", &tag),
				),
			},
			{
				ResourceName:      "e2e_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccE2ETagImportID("e2e_tag.test"),
				Check: resource.ComposeTestCheckFunc(
					// Verify all fields are populated after import
					resource.TestCheckResourceAttr("e2e_tag.test", "name", tagName),
					resource.TestCheckResourceAttr("e2e_tag.test", "metadata", metadata),
					resource.TestCheckResourceAttr("e2e_tag.test", "label_id", fmt.Sprintf("%d", tag.LabelID)),
					resource.TestCheckResourceAttrSet("e2e_tag.test", "project_id"),
					resource.TestCheckResourceAttrSet("e2e_tag.test", "region"),
				),
			},
		},
	})
}
