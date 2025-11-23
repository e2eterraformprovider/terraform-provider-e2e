package faas_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EFaasFunction_Basic(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "name", functionName),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "namespace", namespace),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "runtime", "python-3.11-fastapi"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "memory_mb", "256"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "timeout_seconds", "30"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "min_replicas", "1"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "max_replicas", "5"),
					resource.TestCheckResourceAttrSet("e2e_faas_function.test", "endpoint_url"),
					resource.TestCheckResourceAttrSet("e2e_faas_function.test", "status"),
					resource.TestCheckResourceAttrSet("e2e_faas_function.test", "created_at"),
				),
			},
		},
	})
}

func TestAccE2EFaasFunction_Update(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "memory_mb", "256"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "timeout_seconds", "30"),
				),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_updated(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "memory_mb", "512"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "timeout_seconds", "60"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "min_replicas", "2"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "max_replicas", "10"),
				),
			},
		},
	})
}

func TestAccE2EFaasFunction_WithEnvironment(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_withEnvironment(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment.ENV", "production"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment.DEBUG", "false"),
				),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_withEnvironmentUpdated(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment.ENV", "staging"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment.DEBUG", "true"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment.LOG_LEVEL", "debug"),
				),
			},
		},
	})
}

func TestAccE2EFaasFunction_CodeUpdate(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
				),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_codeUpdated(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
				),
			},
		},
	})
}

func TestAccE2EFaasFunction_DifferentRuntime(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_nodejs(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_faas_function.test", "runtime", "node-18"),
					resource.TestCheckResourceAttrSet("e2e_faas_function.test", "endpoint_url"),
				),
			},
		},
	})
}

func TestAccE2EFaasFunction_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EFaasFunctionConfig_missingName(),
				ExpectError: regexp.MustCompile(`The argument "name" is required`),
			},
			{
				Config:      testAccCheckE2EFaasFunctionConfig_missingNamespace(),
				ExpectError: regexp.MustCompile(`The argument "namespace" is required`),
			},
			{
				Config:      testAccCheckE2EFaasFunctionConfig_missingRuntime(),
				ExpectError: regexp.MustCompile(`The argument "runtime" is required`),
			},
			{
				Config:      testAccCheckE2EFaasFunctionConfig_missingCode(),
				ExpectError: regexp.MustCompile(`The argument "code_inline" is required`),
			},
			{
				Config:      testAccCheckE2EFaasFunctionConfig_missingProjectID(),
				ExpectError: regexp.MustCompile(`The argument "project_id" is required`),
			},
			{
				Config:      testAccCheckE2EFaasFunctionConfig_missingLocation(),
				ExpectError: regexp.MustCompile(`The argument "location" is required`),
			},
		},
	})
}

func TestAccE2EFaasFunction_Import(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
				),
			},
			{
				ResourceName:      "e2e_faas_function.test",
				ImportState:       true,
				ImportStateVerify: true,
				// code_inline is not returned by the API, so it won't match
				ImportStateVerifyIgnore: []string{"code_inline"},
			},
		},
	})
}

func TestAccE2EFaasFunction_ScalingParameters(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_customScaling(functionName, namespace, 0, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "min_replicas", "0"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "max_replicas", "3"),
				),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_customScaling(functionName, namespace, 3, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "min_replicas", "3"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "max_replicas", "20"),
				),
			},
		},
	})
}

// Helper functions

var testAccProvider *schema.Provider

func init() {
	testAccProvider = e2e.Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SERVICE_API_KEY"); v == "" {
		t.Fatal("SERVICE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("SERVICE_AUTH_TOKEN"); v == "" {
		t.Fatal("SERVICE_AUTH_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_PROJECT_ID"); v == "" {
		t.Fatal("E2E_TEST_PROJECT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("E2E_TEST_LOCATION"); v == "" {
		t.Fatal("E2E_TEST_LOCATION must be set for acceptance tests")
	}
}

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"e2e": func() (*schema.Provider, error) {
		return e2e.Provider(), nil
	},
}

func testAccCheckE2EFaasFunctionExists(resourceName string, functionID *string) resource.TestCheckFunc {
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

		opts := &goe2e.RequestOptions{
			ProjectID: projectID,
			Location:  location,
		}

		fn, _, err := client.FaaS.GetFunction(context.Background(), rs.Primary.ID, opts)
		if err != nil {
			return err
		}

		if fn == nil {
			return fmt.Errorf("FaaS function not found")
		}

		*functionID = fn.ID

		return nil
	}
}

func testAccCheckE2EFaasFunctionDestroy(s *terraform.State) error {
	cfg := testAccProvider.Meta().(*config.Config)
	client := cfg.Goe2eClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_faas_function" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		location := rs.Primary.Attributes["location"]

		opts := &goe2e.RequestOptions{
			ProjectID: projectID,
			Location:  location,
		}

		fn, _, err := client.FaaS.GetFunction(context.Background(), rs.Primary.ID, opts)
		if err == nil && fn != nil {
			return fmt.Errorf("FaaS function still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

// Configuration helpers

func testAccCheckE2EFaasFunctionConfig_basic(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name           = "%s"
  namespace      = "%s"
  runtime        = "python-3.11-fastapi"
  code_inline    = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  project_id     = "%s"
  location       = "%s"
}
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_updated(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name            = "%s"
  namespace       = "%s"
  runtime         = "python-3.11-fastapi"
  code_inline     = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  memory_mb       = 512
  timeout_seconds = 60
  min_replicas    = 2
  max_replicas    = 10
  project_id      = "%s"
  location        = "%s"
}
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_withEnvironment(name, namespace string) string {
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
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_withEnvironmentUpdated(name, namespace string) string {
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
    ENV       = "staging"
    DEBUG     = "true"
    LOG_LEVEL = "debug"
  }
  project_id = "%s"
  location   = "%s"
}
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_codeUpdated(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World Updated"}
  EOT
  project_id = "%s"
  location   = "%s"
}
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_nodejs(name, namespace string) string {
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
`, name, namespace, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_customScaling(name, namespace string, minReplicas, maxReplicas int) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name         = "%s"
  namespace    = "%s"
  runtime      = "python-3.11-fastapi"
  code_inline  = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  min_replicas = %d
  max_replicas = %d
  project_id   = "%s"
  location     = "%s"
}
`, name, namespace, minReplicas, maxReplicas, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

// Error case configurations

func testAccCheckE2EFaasFunctionConfig_missingName() string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  namespace   = "test-namespace"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"
  project_id  = "%s"
  location    = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_missingNamespace() string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "test-function"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"
  project_id  = "%s"
  location    = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_missingRuntime() string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "test-function"
  namespace   = "test-namespace"
  code_inline = "def handler(event, context): pass"
  project_id  = "%s"
  location    = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_missingCode() string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name       = "test-function"
  namespace  = "test-namespace"
  runtime    = "python-3.11-fastapi"
  project_id = "%s"
  location   = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"), os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_missingProjectID() string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "test-function"
  namespace   = "test-namespace"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"
  location    = "%s"
}
`, os.Getenv("E2E_TEST_LOCATION"))
}

func testAccCheckE2EFaasFunctionConfig_missingLocation() string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "test-function"
  namespace   = "test-namespace"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"
  project_id  = "%s"
}
`, os.Getenv("E2E_TEST_PROJECT_ID"))
}
