# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the **E2E Networks Terraform Provider**, a production-ready Infrastructure-as-Code provider for managing E2E Cloud resources. The provider is built using the Terraform Plugin SDK v2 and is undergoing a migration from a legacy client architecture to a modern Go SDK called `goe2e`.

### Technology Stack

- **Go 1.24+**
- **Terraform Plugin SDK v2** (v2.27.0)
- **hashicorp/go-retryablehttp** for API client with automatic retry
- **Terraform 0.13+** for running the provider

### Key Repository Links

- Terraform Registry: https://registry.terraform.io/providers/e2eterraformprovider/e2e/latest
- E2E API Docs: https://docs.e2enetworks.com/api/myaccount/
- E2E Console: https://myaccount.e2enetworks.com/

## Essential Commands

### Building and Installing

```bash
# Build the provider (auto-tidies dependencies)
make build

# Install provider to Terraform plugin directory
# Note: Currently hardcoded to darwin_arm64
make install
```

### Testing

```bash
# Unit tests (no cloud access required)
make test

# Acceptance tests (creates real resources, requires API credentials)
# Set these environment variables first:
export E2E_API_KEY="your-api-key"
export E2E_AUTH_TOKEN="your-auth-token"
export E2E_TEST_PROJECT_ID="your-project-id"
export E2E_TEST_REGION="your-region"

make testacc

# Run specific acceptance test
make testacc TESTARGS="-run TestAccE2ENode_Basic"
```

### Code Quality

```bash
# Format code
make fmt

# Run linters (installs golangci-lint if missing)
make lint

# Run go vet
make vet
```

### Upgrading goe2e SDK

```bash
# Uncomment the go get line in Makefile first, then:
make upgrade_goe2e
```

### Other Tasks

```bash
# Update CHANGELOG.md for a new version
make changelog VERSION=2.2.8
```

## Architecture

### High-Level Structure

```
terraform-provider-e2e/
├── e2e/                  # Terraform provider implementation
│   ├── provider.go       # Provider registration, schema, and configuration
│   ├── config/           # Config management and client initialization
│   ├── <resource>/       # One directory per resource/data source
│   │   ├── resource_*.go
│   │   ├── datasource_*.go
│   │   ├── helpers.go    # Resource-specific helpers
│   │   ├── *_test.go     # Acceptance tests
│   │   └── sweep.go      # Test cleanup
│   ├── util/             # Shared utilities (wait, set operations, schema helpers)
│   ├── constants/        # Shared constants for attribute names
│   └── sweep/            # Shared sweeper configuration
│
├── goe2e/                # Modern Go SDK for E2E Networks API
│   ├── goe2e.go          # Core client implementation
│   ├── <service>.go      # Service implementations (faas.go, node.go, etc.)
│   ├── constants/        # API constants
│   └── *_test.go         # SDK unit tests
│
├── docs/                 # Terraform provider documentation
└── Makefile             # Build, test, and development commands
```

### Provider Configuration Pattern

The provider follows the **DigitalOcean provider pattern** for client management:

1. **Config struct** (`e2e/config/config.go`): Holds provider configuration and manages API clients
2. **Client support**:
   - `goe2eClient`: New SDK client goe2e
3. **Provider-level defaults**: Supports `default_region` and `default_project_id`
4. **Client factory method**: `Goe2eClientForProject()` creates resource-specific clients

### goe2e SDK Design Philosophy

The `goe2e` SDK follows **DigitalOcean's godo** design patterns:

1. **Required parameters in client**: `projectID` and `region` are stored in the client, eliminating repetitive parameter passing
2. **Functional options pattern**: Optional configuration via `goe2e.WithWorkspace()`, `goe2e.SetUserAgent()`, etc.
3. **Clean service methods**: Methods focus on business logic, not infrastructure concerns
4. **Single validation point**: Infrastructure requirements validated at client creation
5. **Standard cloud terminology**: Uses `region` externally while mapping to API's `location` internally

**Client creation example:**

```go
client, err := goe2e.NewClient(apiKey, authToken, projectID, region,
    goe2e.WithRetryAndBackoffs(retryConfig),
    goe2e.SetUserAgent("terraform-provider-e2e/1.0"),
)
```

### Resource Implementation Pattern

Each resource follows this structure:

1. **Schema definition**: Using constants from `e2e/constants/` for attribute names
2. **CRUD operations**: Create, Read, Update, Delete using `diag.Diagnostics` return type
3. **Client retrieval**: Get config and create resource-specific client
4. **Region/ProjectID handling**: Use `cfg.GetRegionOrDefault(d)` and `cfg.GetProjectIDOrDefault(d)`
5. **State management**: Map API responses to Terraform state
6. **Error handling**: Return `diag.FromErr()` for errors

**Typical resource function signature:**

```go
func resourceNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    cfg := m.(*config.Config)
    region, err := cfg.GetRegionOrDefault(d)
    projectID, err := cfg.GetProjectIDOrDefault(d)
    client, err := cfg.Goe2eClientForProject(projectID, region)
    // ... implementation
}
```

### goe2e

When working on resources:

- Use `cfg.Goe2eClientForProject()` for new SDK
- Check `goe2e/README.md` for SDK examples
- Follow existing patterns in resources (e.g., `e2e/faas/`, `e2e/node/`)

### Testing Architecture

**Unit Tests** (`*_unit_test.go`):

- Mock-based testing without API calls
- Located alongside production code
- Run with `make test`

**Acceptance Tests** (`*_test.go` in resource packages):

- Create real cloud resources
- Use `acceptance.TestAccProviderFactories` for provider setup
- Follow naming: `TestAccE2E<Resource>_<Scenario>`
- Use `resource.ParallelTest()` for concurrent execution
- Require environment variables: `E2E_API_KEY`, `E2E_AUTH_TOKEN`, `E2E_TEST_PROJECT_ID`, `E2E_TEST_REGION`

**Sweepers** (`sweep.go` files):

- Clean up test resources
- Registered per resource type
- Use `sweep.SharedGoe2eClientForTests()` helper
- Prefix test resources with `test-` (see `sweep.TestNamePrefix`)

**State Upgrade Tests** (`*_state_upgrade_test.go`):

- Test Terraform state migrations between provider versions
- Located in resource directories

### Important Constants

**Shared attribute names** (`e2e/constants/`):

- `AttrRegion`, `AttrLocation`: Region/location parameters (migration in progress)
- `AttrProjectID`: Project identifier
- `AttrName`, `AttrImage`, `AttrPlan`: Common resource attributes

**API constants** (`goe2e/constants/`):

- Node states, plan types, image formats, etc.
- Shared across SDK and provider

### Region/Location Migration

The provider is migrating from `location` to `region` parameter:

- **Prefer `region`** in all new code
- **Support `location`** for backward compatibility with deprecation warning
- **Helper function**: `config.GetRegionOrLocation(d)` handles the migration
- **Provider defaults**: `default_region` supported at provider level

### Utilities and Helpers

**`e2e/util/`**:

- `wait.go`: Resource state waiting functions (e.g., wait for node creation)
- `set.go`: Set operations for Terraform sets
- `schema.go`: Schema helpers
- `errors.go`: Error formatting utilities

**`e2e/config/schema_helpers.go`**:

- `RegionSchema()`, `LocationSchema()`: Region/location schema definitions
- `ProjectIDSchemaResource()`: Project ID schema for resources

**`goe2e` pointer helpers**:

- `goe2e.String()`, `goe2e.Int()`, `goe2e.Bool()`: Create pointers for optional fields
- `goe2e.PtrTo[T]()`: Generic pointer helper

## Development Workflow Guidance

### Adding a New Resource

1. Create directory: `e2e/<resource_name>/`
2. Create files:
   - `resource_<name>.go`: Resource implementation
   - `datasource_<name>.go`: Data source (if applicable)
   - `resource_<name>_test.go`: Acceptance tests
   - `helpers.go`: Resource-specific helpers
   - `sweep.go`: Test cleanup
3. Implement in `goe2e/` if API methods don't exist:
   - Add service interface and implementation
   - Write unit tests in `goe2e/<service>_test.go`
4. Register in `e2e/provider.go`:
   - Add to `ResourcesMap` or `DataSourcesMap`
5. Generate docs: Follow Terraform provider docs conventions

### Working with goe2e SDK

- **Read `goe2e/README.md`** for comprehensive usage examples
- **Service structure**: Each service (FaaS, Nodes, etc.) has its own file
- **Context support**: All methods accept `context.Context` as first parameter
- **Response pattern**: Methods return `(data, *Response, error)`
- **Testing**: Use table-driven tests with mock HTTP responses

### Common Patterns to Follow

1. **Use constants** for attribute names (`e2e/constants/`)
2. **Use `diag.Diagnostics`** for error returns in resource functions
3. **Handle region/project defaults** via `cfg.GetRegionOrDefault()` and `cfg.GetProjectIDOrDefault()`
4. **Name test resources** with `test-` prefix for sweepers
5. **Use `util` package** for common operations (waiting, set operations)
6. **Validate inputs** early in resource CRUD functions
7. **Use resource-specific clients** via `cfg.Goe2eClientForProject()`

### Environment Variables for Testing

Required for acceptance tests:

- `E2E_API_KEY`: API key
- `E2E_AUTH_TOKEN`: Authentication token
- `E2E_TEST_PROJECT_ID`: Project ID for tests
- `E2E_TEST_REGION`: Region for tests (e.g., "Mumbai", "Delhi", "Chennai")

Optional:

- `E2E_API_ENDPOINT`: Custom API endpoint (defaults to production)
- `E2E_REGION`: Provider-level default region
- `E2E_PROJECT_ID`: Provider-level default project

### Debugging Tips

- **Enable TF_LOG**: `export TF_LOG=DEBUG` for verbose Terraform logs
- **Check API responses**: The `goe2e.Response` includes full HTTP response
- **Use sweepers**: Clean up stuck resources with `go test -sweep=<region>`
- **Test incrementally**: Use `-run` flag to run specific tests

## Current Branch Context

**Main branch for PRs**: `develop`

This branch contains:

- Deletion of legacy `client/` and `models/` directories
- Migration to `goe2e` SDK for most resources
- New test infrastructure (`e2e/acceptance/`, `e2e/sweep/`)
- New utility helpers and constants
- Documentation updates for all resources
- Unit tests for critical resources (autoscaling, SFS, Kubernetes, etc.)

Major changes include removal of old client abstraction in favor of direct `goe2e` usage throughout the provider.
