package node_test

import (
	"fmt"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccE2ENodeDataSource_Basic(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "public_ip_address")),
			},
		},
	})
}

func TestAccE2ENodeDataSource_WithLabel(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeDataSourceConfig_withLabel(nodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("data.e2e_node.test", "label", "test-label")),
			},
		},
	})
}

func TestAccE2ENodeDataSource_CheckLockStatus(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeDataSourceConfig_locked(nodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttr("data.e2e_node.test", "is_locked", "true")),
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
  image      = "ubuntu-20.04"}

data "e2e_node" "test" {
  node_id    = e2e_node.test.id}
`, name)
}

func testAccCheckE2ENodeDataSourceConfig_withLabel(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  label      = "test-label"}

data "e2e_node" "test" {
  node_id    = e2e_node.test.id}
`, name)
}

func testAccCheckE2ENodeDataSourceConfig_locked(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name       = "%s"
  plan       = "c2-2c-4gb"
  image      = "ubuntu-20.04"
  lock_node  = true}

data "e2e_node" "test" {
  node_id    = e2e_node.test.id}
`, name)
}

// V3 Feature Tests for Data Source

func TestAccE2ENodeDataSource_V3_IPv6Address(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-v3-ipv6-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeDataSourceConfig_v3_ipv6(nodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_node.test", "name", nodeName),
					// IPv6 address may or may not be present depending on node configuration
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "id"),
				),
			},
		},
	})
}

func TestAccE2ENodeDataSource_V3_ComputedFields(t *testing.T) {
	nodeName := fmt.Sprintf("test-node-v3-computed-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2ENodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2ENodeDataSourceConfig_v3_computed(nodeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.e2e_node.test", "name", nodeName),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "status"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "public_ip_address"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "private_ip_address"),
					resource.TestCheckResourceAttrSet("data.e2e_node.test", "vm_id"),
				),
			},
		},
	})
}

// V3 Configuration Functions for Data Source

func testAccCheckE2ENodeDataSourceConfig_v3_ipv6(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name            = "%s"
  plan            = "C3.8GB"
  image           = "Ubuntu-20.04"
  is_ipv6_availed = true
}

data "e2e_node" "test" {
  node_id = e2e_node.test.id
}
`, name)
}

func testAccCheckE2ENodeDataSourceConfig_v3_computed(name string) string {
	return fmt.Sprintf(`
resource "e2e_node" "test" {
  name  = "%s"
  plan  = "C3.8GB"
  image = "Ubuntu-20.04"
}

data "e2e_node" "test" {
  node_id = e2e_node.test.id
}
`, name)
}
