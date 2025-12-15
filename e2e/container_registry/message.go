package container_registry

// ============================================================================
// Container Registry Resource Name Constant
// ============================================================================

// ResourceName is the resource type name for Container Registry, used in error messages
const ResourceName = "Container Registry"

// ============================================================================
// Container Registry Deprecation Messages
// ============================================================================

// DeprecationMessageSetupStatus is the deprecation message for the setup_status field
const DeprecationMessageSetupStatus = "Use 'status' instead. This parameter will be removed in version 3.0.0"

// DeprecationMessageSetupStatusAlternative is the alternative message suggesting the new field
const DeprecationMessageSetupStatusAlternative = "DEPRECATED: Use 'status' instead"

// DeprecationVersionV3 is the version when deprecated fields will be removed
const DeprecationVersionV3 = "3.0.0"

// ============================================================================
// Container Registry-Specific Error Messages
// ============================================================================
// NOTE: Container Registry operations are not region-scoped in this provider implementation
// (the resource uses the global goe2e client without a region/project context in the error).
// As a result, the generic templates in `e2e/constants/message.go` that include project/region
// are not always a good fit here. Keep container-registry-specific CRUD errors here for clarity,
// and use `e2e/constants/message.go` templates only when you actually have those dimensions.

// ErrorInvalidID is the error message when the container registry ID is invalid
const ErrorInvalidID = "invalid container registry ID: %w"

// ErrorCreateResponseEmpty is the error message when the create response is empty
const ErrorCreateResponseEmpty = "container registry created but response was empty"

// ErrorNotFound is the error message format when a container registry is not found
// Format: "Container Registry project with ID %s not found; removing from state"
const ErrorNotFound = "Container Registry project with ID %s not found; removing from state"

// ErrorSetField is the error message format when setting a field fails
// Format: "failed to set %s: %w"
const ErrorSetField = "failed to set %s: %w"

// ErrorReadRegistry is the error message format when reading a container registry fails
// Format: "failed to read container registry (ID: %s): %w"
const ErrorReadRegistry = "failed to read container registry (ID: %s): %w"

// ErrorDeleteRegistry is the error message format when deleting a container registry fails
// Format: "failed to delete Container Registry: %w"
const ErrorDeleteRegistry = "failed to delete Container Registry: %w"

// ErrorUpdateRegistry is the error message format when updating a container registry fails
// Format: "failed to update container registry: %w"
const ErrorUpdateRegistry = "failed to update container registry: %w"

// ErrorCreateRegistry is the error message format when creating a container registry fails
// Format: "failed to create Container Registry: %w"
const ErrorCreateRegistry = "failed to create Container Registry: %w"

// ============================================================================
// Container Registry Log Messages
// ============================================================================

// LogDeleteWarning is the warning log message when fetching registry details for deletion fails
const LogDeleteWarning = "[WARN] Failed to fetch container registry details for deletion, using default customer ID: %v"

// LogRegistryNotFound is the info log message when a container registry is not found
const LogRegistryNotFound = "[INFO] Container Registry project with ID %s not found; removing from state"

// ============================================================================
// Container Registry Test Constants
// ============================================================================

// TestResourceNamePrefix is the prefix used for test container registry names
const TestResourceNamePrefix = "test-cr-"

// TestResourceType is the Terraform resource type for container registry
const TestResourceType = "e2e_container_registry"
