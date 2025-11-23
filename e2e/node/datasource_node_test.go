package node_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccE2ENodeDataSource_Basic(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeDataSourceConfig_basic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("data.e2e_node.test", "plan", "c2-2c-4gb"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "memory"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "status"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "disk"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "price"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "public_ip_address"),
				),
			},
		},
	})
}

func TestAccE2ENodeDataSource_WithLabel(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeDataSourceConfig_withLabel(nodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("data.e2e_node.test", "label", "test-label"),
				),
			},
		},
	})
}

func TestAccE2ENodeDataSource_CheckLockStatus(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeDataSourceConfig_locked(nodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("data.e2e_node.test", "is_locked", "true"),
				),
			},
		},
	})
}

// Configuration helpers for datasource tests

func testAccCheckE2ENodeDataSourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  project_id = "%s"
  location   = "%s"
}

data "e2e_node" "test" {
  node_id    = e2e_node.test.id
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeDataSourceConfig_withLabel(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  label      = "test-label"
  project_id = "%s"
  location   = "%s"
}

data "e2e_node" "test" {
  node_id    = e2e_node.test.id
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2ENodeDataSourceConfig_locked(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  lock_node  = true
  project_id = "%s"
  location   = "%s"
}

data "e2e_node" "test" {
  node_id    = e2e_node.test.id
  project_id = "%s"
  location   = "%s"
}
`, name, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
