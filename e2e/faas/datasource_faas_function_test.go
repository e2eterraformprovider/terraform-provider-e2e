package faas_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceE2EFaasFunction_Basic(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceE2EFaasFunctionExists("data.e2e_faas_function.test"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "name", functionName),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "namespace", namespace),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "runtime", "python-3.11-fastapi"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "memory_mb", "256"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "timeout_seconds", "30"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "min_replicas", "1"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "max_replicas", "5"),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", "endpoint_url"),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", "status"),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccDataSourceE2EFaasFunction_WithEnvironment(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EFaasFunctionConfig_withEnvironment(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceE2EFaasFunctionExists("data.e2e_faas_function.test"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "name", functionName),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "environment.ENV", "production"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "environment.DEBUG", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceE2EFaasFunction_NotFound(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceE2EFaasFunctionConfig_nodejs(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceE2EFaasFunctionExists("data.e2e_faas_function.test"),
					resource.TestCheckResourceAttr("data.e2e_faas_function.test", "runtime", "node-18"),
					resource.TestCheckResourceAttrSet("data.e2e_faas_function.test", "endpoint_url"),
				),
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

		cfg := testAccProvider.Meta().(*config.Config)
		client := cfg.Goe2eClient()

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]
		functionID := rs.Primary.Attributes["function_id"]

		opts := &goe2e.RequestOptions{
			ProjectID: projectID,
			Location:  location,
		}

		fn, _, err := client.FaaS.GetFunction(context.Background(), functionID, opts)
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
  EOT
  project_id = "%s"
  location   = "%s"
}

data "e2e_faas_function" "test" {
  function_id = e2e_faas_function.test.id
  project_id  = "%s"
  location    = "%s"
}
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
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
  }
  project_id = "%s"
  location   = "%s"
}

data "e2e_faas_function" "test" {
  function_id = e2e_faas_function.test.id
  project_id  = "%s"
  location    = "%s"
}
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
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
  EOT
  project_id = "%s"
  location   = "%s"
}

data "e2e_faas_function" "test" {
  function_id = e2e_faas_function.test.id
  project_id  = "%s"
  location    = "%s"
}
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"),
		os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccDataSourceE2EFaasFunctionConfig_notFound() string {
	return fmt.Sprintf(`
data "e2e_faas_function" "test" {
  function_id = "non-existent-function-id-12345"
  project_id  = "%s"
  location    = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}
