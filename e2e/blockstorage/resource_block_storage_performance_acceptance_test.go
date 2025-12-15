package blockstorage_test

import (
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Performance tests are still TF_ACC gated via testAccPreCheck.
// They are intentionally lightweight (assert correctness + basic multi-resource behavior),
// and do not enforce strict timing thresholds to avoid flakiness.

func TestAccE2EBlockStorage_Perf_CreateMultiple(t *testing.T) {
	prefix := fmt.Sprintf("test-bs-perf-%s", acctest.RandString(8))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_multiple(prefix, 10, 250),
				Check: func() resource.TestCheckFunc {
					checks := make([]resource.TestCheckFunc, 0, 10)
					for i := 0; i < 10; i++ {
						addr := fmt.Sprintf("e2e_blockstorage.perf[%d]", i)
						checks = append(checks,
							resource.TestCheckResourceAttrSet(addr, tfconstants.AttrID),
							resource.TestCheckResourceAttr(addr, tfconstants.AttrSize, "250"),
						)
					}
					return resource.ComposeTestCheckFunc(checks...)
				}(),
			},
		},
	})
}

func TestAccE2EBlockStorage_Perf_LargeUpgrade(t *testing.T) {
	nodeName := fmt.Sprintf("test-bs-perf-node-%s", acctest.RandString(8))
	volumeName := fmt.Sprintf("test-bs-perf-vol-%s", acctest.RandString(8))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EBlockStorageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EBlockStorageConfig_withVolumeAttachment(nodeName, volumeName, 16000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_blockstorage.test", tfconstants.AttrSize, "16000"),
					resource.TestCheckResourceAttrSet("e2e_volume_attachment.test", tfconstants.AttrID),
				),
			},
			{
				Config: testAccCheckE2EBlockStorageConfig_withVolumeAttachment(nodeName, volumeName, 24000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_blockstorage.test", tfconstants.AttrSize, "24000"),
				),
			},
		},
	})
}

func testAccCheckE2EBlockStorageConfig_multiple(prefix string, count int, size int) string {
	return fmt.Sprintf(`
resource "e2e_blockstorage" "perf" {
  count = %d
  name  = format("%s-%%02d", count.index)
  size  = %d
}
`, count, prefix, size)
}

func testAccCheckE2EBlockStorageConfig_withVolumeAttachment(nodeName, volumeName string, size int) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "c2-2c-4gb"
  image = "ubuntu-20.04"
}

resource "e2e_blockstorage" "test" {
  name = "%s"
  size = %d
}

resource "e2e_volume_attachment" "test" {
  node_id   = e2e_node.test.id
  volume_id = e2e_blockstorage.test.id
}
`, nodeName, volumeName, size)
}
