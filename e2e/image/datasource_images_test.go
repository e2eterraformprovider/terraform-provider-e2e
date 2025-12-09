package image_test

import (
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccE2EImagesDataSource_Basic(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImagesDataSourceConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_images.test", "image_list.#")),
			},
		},
	})
}

func TestAccE2EImagesDataSource_VerifyImageList(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))
	imageName := fmt.Sprintf("test-image-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EImagesDataSourceConfig_basic(nodeName, imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.e2e_images.test", "image_list.0.template_id"),
					resource.TestCheckResourceAttrSet("data.e2e_images.test", "image_list.0.image_type"),
					resource.TestCheckResourceAttrSet("data.e2e_images.test", "image_list.0.image_state"),
					resource.TestCheckResourceAttrSet("data.e2e_images.test", "image_list.0.name"),
					resource.TestCheckResourceAttrSet("data.e2e_images.test", "image_list.0.image_id")),
			},
		},
	})
}

// Configuration helpers for datasource tests

func testAccCheckE2EImagesDataSourceConfig_basic(nodeName, imageName string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"}

resource "e2e_image" "test" {
  node_id    = e2e_node.test.id
  name       = "%s"}

data "e2e_images" "test" {
  depends_on = [e2e_image.test]
}
`, nodeName, imageName)
}
