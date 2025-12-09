package faas_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccE2EFaasFunction_Basic(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("e2e_faas_function.test", "created_at")),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "memory_mb", "256"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "timeout_seconds", "30")),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_updated(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "memory_mb", "512"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "timeout_seconds", "60"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "min_replicas", "2"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "max_replicas", "10")),
			},
		},
	})
}

func TestAccE2EFaasFunction_WithEnvironmentVariables(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_withEnvironmentVariables(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment_variables.ENV", "production"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment_variables.DEBUG", "false")),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_withEnvironmentVariablesUpdated(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment_variables.ENV", "staging"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment_variables.DEBUG", "true"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "environment_variables.LOG_LEVEL", "debug")),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID)),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_codeUpdated(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID)),
			},
		},
	})
}

func TestAccE2EFaasFunction_DifferentRuntime(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_nodejs(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("e2e_faas_function.test", "runtime", "node-18"),
					resource.TestCheckResourceAttrSet("e2e_faas_function.test", "endpoint_url")),
			},
		},
	})
}

func TestAccE2EFaasFunction_MissingRequiredArguments(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_basic(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID)),
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
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_customScaling(functionName, namespace, 0, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "min_replicas", "0"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "max_replicas", "3")),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_customScaling(functionName, namespace, 3, 20),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "min_replicas", "3"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "max_replicas", "20")),
			},
		},
	})
}

func TestAccE2EFaasFunction_CodeWhitespace(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_codeWithWhitespace(functionName, namespace, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID)),
			},
			{
				Config:             testAccCheckE2EFaasFunctionConfig_codeWithWhitespace(functionName, namespace, true),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // Should not detect a change due to DiffSuppressFunc
			},
		},
	})
}

func TestAccE2EFaasFunction_ReplicaValidation(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EFaasFunctionConfig_invalidReplicas(functionName, namespace),
				ExpectError: regexp.MustCompile(`min_replicas \(\d+\) cannot be greater than max_replicas \(\d+\)`),
			},
		},
	})
}

// V3 Feature Tests

func TestAccE2EFaasFunction_WithTags(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_withTags(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "tags.Team", "backend")),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_withTagsUpdated(functionName, namespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "tags.Environment", "staging"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "tags.Team", "backend"),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "tags.Owner", "dev-team")),
			},
		},
	})
}

func TestAccE2EFaasFunction_WithDescription(t *testing.T) {
	var functionID string
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckE2EFaasFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckE2EFaasFunctionConfig_withDescription(functionName, namespace, "Initial description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "description", "Initial description")),
			},
			{
				Config: testAccCheckE2EFaasFunctionConfig_withDescription(functionName, namespace, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckE2EFaasFunctionExists("e2e_faas_function.test", &functionID),
					resource.TestCheckResourceAttr("e2e_faas_function.test", "description", "Updated description")),
			},
		},
	})
}

func TestAccE2EFaasFunction_CodeSourceConflict(t *testing.T) {
	functionName := fmt.Sprintf("test-func-%s", acctest.RandString(10))
	namespace := fmt.Sprintf("test-ns-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckE2EFaasFunctionConfig_bothCodeSources(functionName, namespace),
				ExpectError: regexp.MustCompile(`code_inline and code_file are mutually exclusive|conflicts with`),
			},
			{
				Config:      testAccCheckE2EFaasFunctionConfig_noCodeSource(functionName, namespace),
				ExpectError: regexp.MustCompile(`one of code_inline or code_file must be specified`),
			},
		},
	})
}

// Helper functions

func testAccPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
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

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client := cfg.Goe2eClient()

		fn, _, err := client.FaaS.GetFunction(context.Background(), rs.Primary.ID)
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
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client := cfg.Goe2eClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "e2e_faas_function" {
			continue
		}

		fn, _, err := client.FaaS.GetFunction(context.Background(), rs.Primary.ID)
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
  EOT}
`, name, namespace)
}

func testAccCheckE2EFaasFunctionConfig_updated(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name              = "%s"
  namespace         = "%s"
  runtime           = "python-3.11-fastapi"
  code_inline       = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  memory_mb         = 512
  timeout_seconds   = 60
  min_replicas      = 2
  max_replicas      = 10
}
`, name, namespace)
}

func testAccCheckE2EFaasFunctionConfig_withEnvironmentVariables(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  environment_variables = {
    ENV   = "production"
    DEBUG = "false"
  }
}
`, name, namespace)
}

func testAccCheckE2EFaasFunctionConfig_withEnvironmentVariablesUpdated(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  environment_variables = {
    ENV       = "staging"
    DEBUG     = "true"
    LOG_LEVEL = "debug"
  }
}
`, name, namespace)
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
  EOT}
`, name, namespace)
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
  EOT}
`, name, namespace)
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
  max_replicas = %d}
`, name, namespace, minReplicas, maxReplicas)
}

// Error case configurations

func testAccCheckE2EFaasFunctionConfig_missingName() string {
	return `
resource "e2e_faas_function" "test" {
  namespace   = "test-namespace"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"}
`
}

func testAccCheckE2EFaasFunctionConfig_missingNamespace() string {
	return `
resource "e2e_faas_function" "test" {
  name        = "test-function"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"}
`
}

func testAccCheckE2EFaasFunctionConfig_missingRuntime() string {
	return `
resource "e2e_faas_function" "test" {
  name        = "test-function"
  namespace   = "test-namespace"
  code_inline = "def handler(event, context): pass"}
`
}

func testAccCheckE2EFaasFunctionConfig_missingCode() string {
	return `
resource "e2e_faas_function" "test" {
  name       = "test-function"
  namespace  = "test-namespace"
  runtime    = "python-3.11-fastapi"}
`
}

func testAccCheckE2EFaasFunctionConfig_missingProjectID() string {
	return `
resource "e2e_faas_function" "test" {
  name        = "test-function"
  namespace   = "test-namespace"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"}
`
}

func testAccCheckE2EFaasFunctionConfig_missingLocation() string {
	return `
resource "e2e_faas_function" "test" {
  name        = "test-function"
  namespace   = "test-namespace"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"}
`
}

func testAccCheckE2EFaasFunctionConfig_codeWithWhitespace(name, namespace string, addWhitespace bool) string {
	code := "def handler(event, context):\n    return {\"statusCode\": 200, \"body\": \"Hello World\"}"
	if addWhitespace {
		// Add extra whitespace that should be ignored by DiffSuppressFunc
		code = "\n  def handler(event, context):\n    return {\"statusCode\": 200, \"body\": \"Hello World\"}  \n\n"
	}

	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
%s
  EOT}
`, name, namespace, code)
}

func testAccCheckE2EFaasFunctionConfig_invalidReplicas(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name         = "%s"
  namespace    = "%s"
  runtime      = "python-3.11-fastapi"
  code_inline  = "def handler(event, context): pass"
  min_replicas = 10
  max_replicas = 5
}
`, name, namespace)
}

// V3 Feature Configuration Helpers

func testAccCheckE2EFaasFunctionConfig_withTags(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  tags = {
    Environment = "production"
    Team        = "backend"
  }
}
`, name, namespace)
}

func testAccCheckE2EFaasFunctionConfig_withTagsUpdated(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  tags = {
    Environment = "staging"
    Team        = "backend"
    Owner       = "dev-team"
  }
}
`, name, namespace)
}

func testAccCheckE2EFaasFunctionConfig_withDescription(name, namespace, description string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = <<-EOT
    def handler(event, context):
        return {"statusCode": 200, "body": "Hello World"}
  EOT
  description = "%s"
}
`, name, namespace, description)
}

func testAccCheckE2EFaasFunctionConfig_bothCodeSources(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name        = "%s"
  namespace   = "%s"
  runtime     = "python-3.11-fastapi"
  code_inline = "def handler(event, context): pass"
  code_file   = "/path/to/code.zip"
}
`, name, namespace)
}

func testAccCheckE2EFaasFunctionConfig_noCodeSource(name, namespace string) string {
	return fmt.Sprintf(`
resource "e2e_faas_function" "test" {
  name      = "%s"
  namespace = "%s"
  runtime   = "python-3.11-fastapi"
}
`, name, namespace)
}
