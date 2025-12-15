package constants

import "fmt"

// ============================================================================
// Generic Resource Operation Error Templates
// ============================================================================
// These templates can be reused across all resources by substituting the resource name
// This reduces duplication and ensures consistent error message formatting

// ResourceOperationErrorTemplate is a generic template for resource operation errors
// Format: "Error {operation} {resource} ({identifier}) in project (%s), region (%s): %s"
// Usage: fmt.Errorf(ResourceOperationErrorTemplate, "creating", "SFS", name, projectID, region, err)
const ResourceOperationErrorTemplate = "Error %s %s (%s) in project (%s), region (%s): %s"

// ResourceOperationByIDErrorTemplate is a generic template for resource operations using ID
// Format: "Error {operation} {resource} (ID: %s) in project (%s), region (%s): %s"
// Usage: fmt.Errorf(ResourceOperationByIDErrorTemplate, "retrieving", "SFS", id, projectID, region, err)
const ResourceOperationByIDErrorTemplate = "Error %s %s (ID: %s) in project (%s), region (%s): %s"

// ResourceOperationWaitErrorTemplate is a generic template for waiting operations
// Format: "Error waiting for {resource} (ID: %s) to become {state} in project (%s), region (%s): %s"
// Usage: fmt.Errorf(ResourceOperationWaitErrorTemplate, "SFS", id, "active", projectID, region, err)
const ResourceOperationWaitErrorTemplate = "Error waiting for %s (ID: %s) to become %s in project (%s), region (%s): %s"

// ResourceDataSourceListErrorTemplate is a generic template for listing resources in data sources
// Format: "Error listing {resource} instances in project (%s), region (%s): %s"
// Usage: fmt.Errorf(ResourceDataSourceListErrorTemplate, "SFS", projectID, region, err)
const ResourceDataSourceListErrorTemplate = "Error listing %s instances in project (%s), region (%s): %s"

// ResourceDeleteStateErrorTemplate is a generic template for delete operations when resource is in wrong state
// Format: "Cannot delete {resource} (ID: %s): {resource} is in {state} state in project (%s), region (%s). {message}"
// Usage: fmt.Errorf(ResourceDeleteStateErrorTemplate, "SFS", id, "Creating", projectID, region, "Please wait for SFS creation to complete")
const ResourceDeleteStateErrorTemplate = "Cannot delete %s (ID: %s): %s is in %s state in project (%s), region (%s). %s"

// ResourceCreateRequiredFieldTemplate is a generic template for required field errors during creation
// Format: "Error creating {resource} (name: %s): {field_description} must be specified"
// Usage: fmt.Errorf(ResourceCreateRequiredFieldTemplate, "SFS", name, "size_gb or disk_size")
const ResourceCreateRequiredFieldTemplate = "Error creating %s (name: %s): %s must be specified"

// ResourceCreateInvalidResponseTemplate is a generic template for invalid API responses
// Format: "Error creating {resource} (name: %s) in project (%s), region (%s): unable to retrieve valid '{field}' from API response"
// Usage: fmt.Errorf(ResourceCreateInvalidResponseTemplate, "SFS", name, projectID, region, "efs_id")
const ResourceCreateInvalidResponseTemplate = "Error creating %s (name: %s) in project (%s), region (%s): unable to retrieve valid '%s' from API response"

// ============================================================================
// Generic State Setting Error Helper
// ============================================================================

// ErrorSettingStateFormat returns a formatted error message for state setting failures
// This eliminates the need for individual constants for each field
// Usage: fmt.Errorf(ErrorSettingStateFormat("name"), err)
func ErrorSettingStateFormat(field string) string {
	return fmt.Sprintf("error setting %s: %%w", field)
}

// ============================================================================
// Generic Validation Error Messages
// ============================================================================

// NameExpectedString is a generic error message when name validation receives a non-string type
// Can be reused across all resources that validate names
const NameExpectedString = "expected name to be string"

// NameCannotContainWhitespaceTemplate is a generic template for name whitespace validation
// Format: "name cannot contain whitespace. Got %s"
// Usage: fmt.Errorf(NameCannotContainWhitespaceTemplate, value)
const NameCannotContainWhitespaceTemplate = "name cannot contain whitespace. Got %s"

// DatabaseConfigurationRequired is a generic error message for missing nested database config blocks.
const DatabaseConfigurationRequired = "database configuration is required"

// ============================================================================
// Generic Import ID Error Templates
// ============================================================================

// ImportIDInvalidFormatTemplate is a generic template for invalid import ID format
// Format: "invalid import ID format: %s. Expected: {format_description}"
// Usage: fmt.Errorf(ImportIDInvalidFormatTemplate, id, "either <resource_id> or <project_id>/<region>/<resource_id>")
const ImportIDInvalidFormatTemplate = "invalid import ID format: %s. Expected: %s"

// ============================================================================
// Operation Constants (Used in Generic Templates)
// ============================================================================
// These constants represent common operation names used in error messages

// OperationCreating represents the "creating" operation
const OperationCreating = "creating"

// OperationRetrieving represents the "retrieving" operation
const OperationRetrieving = "retrieving"

// OperationDeleting represents the "deleting" operation
const OperationDeleting = "deleting"

// OperationUpdating represents the "updating" operation
const OperationUpdating = "updating"

// OperationListing represents the "listing" operation
const OperationListing = "listing"

// ============================================================================
// Client Creation Errors (Generic)
// ============================================================================

// ErrorCreatingGoe2eClient is the error message format when goe2e client creation fails
// This is generic and reusable across all resources
// Format: "Error creating goe2e client: %s" where %s is the error
const ErrorCreatingGoe2eClient = "Error creating goe2e client: %s"

// ============================================================================
// SSH Key Import and Error Handling Constants
// ============================================================================
// These are Terraform provider-specific constants for SSH key resource import
// and error checking, not part of the API contract.

// SSHKeyImportIDFormatDescription describes the expected format for importing SSH keys
// into Terraform. Users must provide either two or three colon-separated parts.
const SSHKeyImportIDFormatDescription = "project_id:ssh_key_id or project_id:region:ssh_key_id"

// SSHKeyImportIDRegionRequired is the error message shown when region is not specified
// in the import ID and no default region is configured in the provider
const SSHKeyImportIDRegionRequired = "region must be specified in import ID or provider default_region must be set. Use format: project_id:region:ssh_key_id"

// SSHKeyImportIDInvalidFormat is the error message shown when the import ID format is invalid
const SSHKeyImportIDInvalidFormat = "invalid import ID format. Expected: project_id:ssh_key_id or project_id:region:ssh_key_id"
