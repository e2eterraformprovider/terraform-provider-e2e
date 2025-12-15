package blockstorage_test

import (
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccE2EBlockStorageDataSource_Basic(t *testing.T) {
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageDataSourceConfig_basic(blockStorageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_blockstorage.test", tfconstants.AttrName, blockStorageName),
					resource.TestCheckResourceAttrSet("data.e2e_blockstorage.test", tfconstants.AttrSize),
					resource.TestCheckResourceAttrSet("data.e2e_blockstorage.test", tfconstants.AttrIOPS),
					resource.TestCheckResourceAttrSet("data.e2e_blockstorage.test", tfconstants.AttrStatus)),
			},
		},
	})
}

func TestAccE2EBlockStorageDataSource_VerifyAttributes(t *testing.T) {
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageDataSourceConfig_basic(blockStorageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_blockstorage.test", tfconstants.AttrName, blockStorageName),
					resource.TestCheckResourceAttr("data.e2e_blockstorage.test", tfconstants.AttrSize, "250")),
			},
		},
	})
}

// Configuration helpers for datasource tests

func testAccCheckE2EBlockStorageDataSourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name = "%s"
  size = 250
}

data "e2e_blockstorage" "test" {
  block_id = e2e_blockstorage.test.id
}
`, name)
}
