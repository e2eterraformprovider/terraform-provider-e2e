# E2E FaaS Implementation Proposal

## Overview
Add Function-as-a-Service (FaaS) support to the E2E Terraform provider following existing provider patterns.

## Implementation Structure

### 1. Resource: `e2e_faas_function`
**Location**: `e2e/faas/resource_faas_function.go`

**Schema**:
```hcl
resource "e2e_faas_function" "api_handler" {
  name         = "order-processor"
  namespace    = "production"
  runtime      = "python-3.11-fastapi"
  code_inline  = file("${path.module}/handler.py")

  config {
    memory_mb       = 512
    timeout_seconds = 60
    min_replicas    = 1
    max_replicas    = 5
  }

  environment = {
    DB_HOST = "postgres.example.com"
    API_KEY = var.api_key
  }

  project_id = var.project_id
  location   = "Delhi"
}

output "function_url" {
  value = e2e_faas_function.api_handler.endpoint_url
}
```

**Key Attributes**:
- `name` (Required): Function name
- `namespace` (Required): FaaS namespace
- `runtime` (Required): e.g., python-3.11-fastapi, node-18, go-1.21
- `code_inline` (Required): Function code as string
- `config` (Optional): Memory, timeout, replicas
- `environment` (Optional): Environment variables map
- `project_id` (Required): E2E project ID
- `location` (Required): Region
- Computed: `id`, `endpoint_url`, `status`, `created_at`

### 2. Data Source: `e2e_faas_function`
**Location**: `e2e/faas/datasource_faas_function.go`

**Usage**:
```hcl
data "e2e_faas_function" "existing" {
  function_id = "func-12345"
  namespace   = "production"
  project_id  = var.project_id
  location    = "Delhi"
}
```

### 3. Models
**Location**: `models/faas.go`

```go
type FaasFunction struct {
    Name        string            `json:"name"`
    Namespace   string            `json:"namespace"`
    Runtime     string            `json:"runtime"`
    Code        string            `json:"code"`
    Memory      int               `json:"memory_mb"`
    Timeout     int               `json:"timeout_seconds"`
    MinReplicas int               `json:"min_replicas"`
    MaxReplicas int               `json:"max_replicas"`
    Environment map[string]string `json:"environment"`
}

type FaasFunctionResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Namespace   string `json:"namespace"`
    EndpointURL string `json:"endpoint_url"`
    Status      string `json:"status"`
    CreatedAt   string `json:"created_at"`
}
```

### 4. Client Methods
**Location**: `client/faas.go`

```go
// API Endpoints
// POST   /faas/namespace                    - Create namespace (auto-handle)
// POST   /faas/functions                    - Create function
// GET    /faas/function/{function_id}       - Read function
// PUT    /faas/function/{function_id}       - Update function
// DELETE /faas/function/{function_id}       - Delete function
// GET    /faas/logs/{function_id}           - Get logs (future)

func (c *Client) CreateFaasFunction(fn *FaasFunction, projectID, location string) (*FaasFunctionResponse, error)
func (c *Client) GetFaasFunction(functionID, namespace, projectID, location string) (*FaasFunctionResponse, error)
func (c *Client) UpdateFaasFunction(functionID string, fn *FaasFunction, projectID, location string) error
func (c *Client) DeleteFaasFunction(functionID, namespace, projectID, location string) error
```

### 5. Provider Registration
**Location**: `e2e/provider.go`

Add to ResourcesMap:
```go
"e2e_faas_function": faas.ResourceFaasFunction(),
```

Add to DataSourcesMap:
```go
"e2e_faas_function": faas.DataSourceFaasFunction(),
```

## Implementation Order

1. **Create models** (`models/faas.go`)
2. **Create client methods** (`client/faas.go`) - Follow pattern from `client/client.go` ssh_key methods
3. **Create resource** (`e2e/faas/resource_faas_function.go`) - Follow pattern from `e2e/ssh_key/resource_ssh_key.go`
4. **Create data source** (`e2e/faas/datasource_faas_function.go`) - Follow pattern from `e2e/ssh_key/datasource_ssh_key.go`
5. **Register in provider** (`e2e/provider.go`)
6. **Test with example** (Create test HCL configuration)

## API Authentication Pattern
Follow existing pattern from client.go:
- Bearer token in Authorization header
- API key as query parameter
- Standard headers: Content-Type, User-Agent

## Notes
- Namespace creation handled automatically in CreateFaasFunction if not exists
- Function code must be provided inline (base64 encoding if needed)
- Supported runtimes should match E2E Networks FaaS offerings
- Error handling follows existing provider patterns
