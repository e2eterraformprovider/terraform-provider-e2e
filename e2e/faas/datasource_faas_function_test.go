package faas_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceE2EFaasFunction_Basic(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceE2EFaasFunctionExists("data.e2e_faas_function.test"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrName, functionName),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrNamespace, namespace),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrRuntime, "python-3.11-fastapi"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrMemoryMB, fmt.Sprintf("%d", goe2econstants.FaaSDefaultMemoryMB)),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrTimeoutSeconds, fmt.Sprintf("%d", goe2econstants.FaaSDefaultTimeoutSeconds)),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrMinReplicas, fmt.Sprintf("%d", goe2econstants.FaaSDefaultMinReplicas)),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrMaxReplicas, fmt.Sprintf("%d", goe2econstants.FaaSDefaultMaxReplicas)),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", tfconstants.AttrEndpointURL),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", tfconstants.AttrStatus),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", tfconstants.AttrCreatedAt),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", tfconstants.AttrUpdatedAt)),
			},
		},
	})
}

func TestAccDataSourceE2EFaasFunction_WithEnvironment(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EFaasFunctionConfig_withEnvironment(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceE2EFaasFunctionExists("data.e2e_faas_function.test"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrName, functionName),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "environment.ENV", "production"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "environment.DEBUG", "false")),
			},
		},
	})
}

func TestAccDataSourceE2EFaasFunction_NotFound(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceE2EFaasFunctionConfig_notFound(),
				ExpectError: regexp.MustCompile(`FaaS function with ID .* not found`),
			},
		},
	})
}

func TestAccDataSourceE2EFaasFunction_DifferentRuntime(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EFaasFunctionConfig_nodejs(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceE2EFaasFunctionExists("data.e2e_faas_function.test"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", tfconstants.AttrRuntime, "node-18"),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", tfconstants.AttrEndpointURL)),
			},
		},
	})
}

// Helper functions

func testAccDataSourceE2EFaasFunctionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No FaaS function ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client := cfg.Goe2eClient()

		functionID := rs.Primary.Attributes[tfconstants.AttrFunctionID]

		fn, _, err := client.FaaS.GetFunction(context.Background(), functionID)
		if err != nil {
			return err
		}

		if fn == nil {
			return fmt.Errorf("FaaS function not found")
		}

		return nil
	}
}

// Configuration helpers

func testAccDataSourceE2EFaasFunctionConfig_basic(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT}

data "e2e_faas_function" "test" {
  function_id = e2e_faas_function.test.id}
`, name, namespace)
}

func testAccDataSourceE2EFaasFunctionConfig_withEnvironment(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  environment = {
    ENV   = "production"
    DEBUG = "false"
  }}

data "e2e_faas_function" "test" {
  function_id = e2e_faas_function.test.id}
`, name, namespace)
}

func testAccDataSourceE2EFaasFunctionConfig_nodejs(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "node-18"
  code_inline = <<-EOT
    module.exports = async (event, context) => {
      return {statusCode: 200, body: "Hello from Node.js"};
    };
  EOT}

data "e2e_faas_function" "test" {
  function_id = e2e_faas_function.test.id}
`, name, namespace)
}

func testAccDataSourceE2EFaasFunctionConfig_notFound() string {
	return `
data "e2e_faas_function" "test" {
  function_id = "non-existent-function-id-12345"}
`
}
