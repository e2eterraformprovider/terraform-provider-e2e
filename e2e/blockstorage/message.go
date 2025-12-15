package blockstorage

// ============================================================================
// Block Storage Resource Constants (Terraform/provider-side)
// ============================================================================

// ResourceName is the human-readable resource name used in generic error templates.
const ResourceName = "Block Storage"

// Import ID format description for `terraform import`.
const (
	ImportIDFormatDescription = "project_id/region/block_storage_id"
	ImportIDPartsCount        = 3
)

// Block storage-specific error templates (non-generic or provider-only semantics).
const (
	// ErrorValidateSizeTemplate is used when plan lookup/validation fails.
	// Format: "Error validating block storage size for project (%s), region (%s): %s"
	ErrorValidateSizeTemplate = "Error validating block storage size for project (%s), region (%s): %s"

	// ErrorCreateNilResponseTemplate indicates the create call returned a nil object.
	// Format: "Error creating Block Storage (%s): received nil response"
	ErrorCreateNilResponseTemplate = "Error creating " + ResourceName + " (%s): received nil response"

	// ErrorResizeVMIDMissingTemplate indicates we could not find a VM ID in vm_detail.
	// Format: "Cannot resize block storage (ID: %s): VM ID not found in block storage details"
	ErrorResizeVMIDMissingTemplate = "Cannot resize block storage (ID: %s): VM ID not found in block storage details"

	// ErrorResizeReduceNotAllowedTemplate indicates a requested size decrease.
	// Format: "Cannot reduce block storage (ID: %s) size: only upgrades (increases) are allowed"
	ErrorResizeReduceNotAllowedTemplate = "Cannot reduce block storage (ID: %s) size: only upgrades (increases) are allowed"

	// ErrorResizeRequiresAttachmentTemplate indicates block storage must be attached to resize.
	// Format: "Cannot resize block storage (ID: %s): block storage must be attached to a node"
	ErrorResizeRequiresAttachmentTemplate = "Cannot resize block storage (ID: %s): block storage must be attached to a node"

	// ErrorResizeConcurrentDiskOperationTemplate indicates a concurrent disk resize on the same VM.
	// Format: "Cannot resize block storage (ID: %s, name: %s): currently resizing another disk on same virtual machine (VM ID: %.0f). Please wait"
	ErrorResizeConcurrentDiskOperationTemplate = "Cannot resize block storage (ID: %s, name: %s): currently resizing another disk on same virtual machine (VM ID: %.0f). Please wait"

	// ErrorUpdateInErrorStateTemplate indicates updates are blocked when in ERROR.
	// Format: "Cannot update block storage (ID: %s): block storage is in ERROR state in project (%s), region (%s)"
	ErrorUpdateInErrorStateTemplate = "Cannot update block storage (ID: %s): block storage is in ERROR state in project (%s), region (%s)"

	// ErrorDeleteInStateTemplate indicates delete is blocked in certain transient states.
	// Format: "Cannot delete block storage (ID: %s): block storage is in %s state in project (%s), region (%s)"
	ErrorDeleteInStateTemplate = "Cannot delete block storage (ID: %s): block storage is in %s state in project (%s), region (%s)"

	// ErrorDeleteWhileAttachedTemplate indicates delete is blocked when attached.
	// Format: "Cannot delete block storage (ID: %s): block storage is attached to node (VM ID: %s). Detach it first"
	ErrorDeleteWhileAttachedTemplate = "Cannot delete block storage (ID: %s): block storage is attached to node (VM ID: %s). Detach it first"

	// LogNotFoundRemoveFromStateTemplate logs when a resource is not found and is removed from state.
	// Format: "[WARN] Block storage (ID: %s) not found, removing from state"
	LogNotFoundRemoveFromStateTemplate = "[WARN] Block storage (ID: %s) not found, removing from state"

	// LogImportTemplate logs import context.
	// Format: "[INFO] Importing block storage: id=%s, project=%s, region=%s"
	LogImportTemplate = "[INFO] Importing block storage: id=%s, project=%s, region=%s"
)
