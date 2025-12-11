package constants

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
