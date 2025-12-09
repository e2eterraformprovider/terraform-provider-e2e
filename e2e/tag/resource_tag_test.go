package tag_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2ETag_Basic(t *testing.T) {
	var tag goe2e.Tag
	tagName := fmt.Sprintf("test-tag-%s", acctest.RandString(10))

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
	tagName := fmt.Sprintf("test-tag-%s", acctest.RandString(10))
	metadata := "Test tag for acceptance testing"

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
	tagName := fmt.Sprintf("test-tag-%s", acctest.RandString(10))

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
	tagName1 := fmt.Sprintf("test-tag-%s", acctest.RandString(10))
	tagName2 := fmt.Sprintf("test-tag-%s", acctest.RandString(10))

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
		if rs.Type != "e2e_tag" {
			continue
		}

		_, _, err := client.Tags.GetTag(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Tag still exists: %s", rs.Primary.ID)
		}

		// Verify it's a not found error
		if !containsString(err.Error(), "not found") {
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

func containsString(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
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
