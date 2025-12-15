package sfs

// ============================================================================
// SFS Resource Name Constant
// ============================================================================

// ResourceName is the resource type name for SFS, used in error messages
const ResourceName = "SFS"

// ============================================================================
// SFS-Specific Error Messages
// ============================================================================
// These are SFS-specific errors that cannot be generalized due to:
// - V2/V3 field migration specifics (size_gb/disk_size, iops/disk_iops)
// - Resource-specific state transitions (Creating state)
// - Resource-specific import format

// CreateSizeRequired is SFS-specific because it mentions both size_gb and disk_size (V2/V3 migration)
// Format: "Error creating SFS (name: %s): size_gb or disk_size must be specified"
const CreateSizeRequired = "Error creating SFS (name: %s): size_gb or disk_size must be specified"

// CreateIOPSRequired is SFS-specific because it mentions both iops and disk_iops (V2/V3 migration)
// Format: "Error creating SFS (name: %s): iops or disk_iops must be specified"
const CreateIOPSRequired = "Error creating SFS (name: %s): iops or disk_iops must be specified"

// DeleteCreatingState is SFS-specific because it mentions "Creating" state specifically
// Format: "Cannot delete SFS (ID: %s): SFS is in Creating state in project (%s), region (%s). Please wait for SFS creation to complete"
const DeleteCreatingState = "Cannot delete SFS (ID: %s): SFS is in Creating state in project (%s), region (%s). Please wait for SFS creation to complete"

// ImportIDInvalidFormat is SFS-specific because it mentions SFS-specific import format
// Format: "invalid import ID format: %s. Expected either <sfs_id> or <project_id>/<region>/<sfs_id>"
const ImportIDInvalidFormat = "invalid import ID format: %s. Expected either <sfs_id> or <project_id>/<region>/<sfs_id>"
