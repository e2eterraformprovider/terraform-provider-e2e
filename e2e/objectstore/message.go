package objectstore

// ============================================================================
// Object Store Resource Name Constant
// ============================================================================

// ResourceName is the resource type name for Object Store, used in error messages
const ResourceName = "object storage bucket"

// ============================================================================
// Object Store-Specific Error Messages
// ============================================================================
// These are Object Store-specific errors that cannot be generalized due to:
// - Resource-specific import format (bucket_name or project_id:region:bucket_name)
// - Resource-specific deprecation warnings (enabling_versioning)
// - Resource-specific deletion constraints (lock policy)

// ImportIDInvalidFormat is Object Store-specific import format error
// Format: "invalid import ID format: expected 'bucket_name' or 'project_id:region:bucket_name', got '%s'"
const ImportIDInvalidFormat = "invalid import ID format: expected 'bucket_name' or 'project_id:region:bucket_name', got '%s'"

// WarnEnablingVersioningDeprecated is Object Store-specific deprecation warning
// Format: "enabling_versioning is deprecated. Use versioning_enabled instead."
const WarnEnablingVersioningDeprecated = "enabling_versioning is deprecated. Use versioning_enabled instead."

// DeleteLockPolicyEnabled is Object Store-specific error when trying to delete a bucket with lock policy enabled
// Format: "Cannot delete bucket with lock policy enabled (name: %s). Disable lock first."
const DeleteLockPolicyEnabled = "Cannot delete bucket with lock policy enabled (name: %s). Disable lock first."

// ErrorUpdatingVersioning is Object Store-specific error when updating versioning fails
// Format: "Error updating versioning (%s) for object storage bucket (name: %s) in project (%s), region (%s): %s"
const ErrorUpdatingVersioning = "Error updating versioning (%s) for object storage bucket (name: %s) in project (%s), region (%s): %s"
