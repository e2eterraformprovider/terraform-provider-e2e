package blockstorage_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccE2EBlockStorageDataSource_Basic(t *testing.T) {
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageDataSourceConfig_basic(blockStorageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_blockstorage.test", "name", blockStorageName),
					resource.TestCheckResourceAttrSet("data.e2e_blockstorage.test", "size"),
					resource.TestCheckResourceAttrSet("data.e2e_blockstorage.test", "iops"),
					resource.TestCheckResourceAttrSet("data.e2e_blockstorage.test", "status"),
				),
			},
		},
	})
}

func TestAccE2EBlockStorageDataSource_VerifyAttributes(t *testing.T) {
	blockStorageName := fmt.Sprintf("test-bs-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageDataSourceConfig_basic(blockStorageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_blockstorage.test", "name", blockStorageName),
					resource.TestCheckResourceAttr("data.e2e_blockstorage.test", "size", "10"),
				),
			},
		},
	})
}

// Configuration helpers for datasource tests

func testAccCheckE2EBlockStorageDataSourceConfig_basic(name string) string {
	projectIDStr := os.Getenv("E2E_TEST_PROJECT_ID")
	projectID, _ := strconv.Atoi(projectIDStr)

	return fmt.Sprintf(`
resource "e2e_blockstorage" "test" {
  name       = "%s"
  size       = 10
  project_id = %d
  location   = "%s"
}

data "e2e_blockstorage" "test" {
  block_id   = e2e_blockstorage.test.id
  project_id = %d
  location   = "%s"
}
`, name, projectID, os.Getenv("E2E_TEST_LOCATION"),
		projectID, os.Getenv("E2E_TEST_LOCATION"))
}
