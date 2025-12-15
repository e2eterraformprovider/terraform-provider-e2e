package image

// ============================================================================
// Image Resource Name Constant
// ============================================================================

// ResourceName is the resource type name for Image, used in error messages
const ResourceName = "Image"

// ============================================================================
// Image-Specific Error Messages
// ============================================================================
// These are Image-specific errors that cannot be generalized due to:
// - Resource-specific import format (image_id or project_id/region/image_id)
// - Resource-specific validation (name cannot contain whitespace)
// - Resource-specific deletion constraints (cloning operations)

// ImportIDInvalidFormat is Image-specific import format error
// Format: "invalid import ID format, expected 'image_id' or 'project_id/region/image_id', got: %s"
const ImportIDInvalidFormat = "invalid import ID format, expected 'image_id' or 'project_id/region/image_id', got: %s"

// DeleteCloningOps is Image-specific error when trying to delete an image with ongoing cloning operations
// Format: "cannot delete image with ongoing cloning operations (cloning_ops: %s)"
const DeleteCloningOps = "cannot delete image with ongoing cloning operations (cloning_ops: %s)"

// WarnImageRunningVMs is Image-specific warning when deleting an image with running VMs
// Format: "Image has %s running VMs, deletion will proceed"
const WarnImageRunningVMs = "Image has %s running VMs, deletion will proceed"

// WarnLocationDeprecated is Image-specific deprecation warning for location field
// Format: "Parameter 'location' is deprecated and will be removed in v4.0. Please use 'region' instead"
const WarnLocationDeprecated = "Parameter 'location' is deprecated and will be removed in v4.0. Please use 'region' instead"

// RegionLocationConflict is Image-specific error when both region and location are set
// Format: "cannot set both 'region' and 'location' parameters"
const RegionLocationConflict = "cannot set both 'region' and 'location' parameters"

// ErrorCheckingImageState is Image-specific error when checking image state fails
// Format: "error checking image state: %w" where %w is the underlying error
const ErrorCheckingImageState = "error checking image state: %w"
