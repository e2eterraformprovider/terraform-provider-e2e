package volume_attachment

// ============================================================================
// Volume Attachment Resource Constants (Terraform/provider-side)
// ============================================================================

// ResourceName is the human-readable resource name used in generic error templates.
const ResourceName = "Volume Attachment"

// Import formats for `terraform import`.
const (
	ImportIDFormatShortDescription = "node_id/volume_id"
	ImportIDFormatFullDescription  = "project_id/region/node_id/volume_id"

	ImportIDPartsShortCount = 2
	ImportIDPartsFullCount  = 4
)

// Provider-only error messages.
const (
	ErrorNodeNotFoundTemplate                 = "Error: node (ID: %s) not found"
	ErrorC2PlanNoBlockStorageAttachmentFormat = "Cannot attach volume to node (ID: %s): C2 plan nodes do not support block storage attachment"
	ErrorImportReadDuringImportTemplate       = "error reading volume attachment during import: %v"
	ErrorParseIDTemplate                      = "invalid volume attachment ID format: %s (expected: %s)"
	ErrorWaitContextCancelledTemplate         = "context cancelled while waiting for volume %s"
	ErrorWaitTimeoutTemplate                  = "timeout waiting for volume (ID: %s) to %s to/from node (ID: %s)"
)

// Log message templates for volume attachment operations.
const (
	LogAttachTemplate    = "[INFO] Attaching volume (ID: %s) to node (ID: %s) in project (%s), region (%s)"
	LogAttachedTemplate  = "[INFO] Successfully attached volume (ID: %s) to node (ID: %s)"
	LogReadTemplate      = "[INFO] Reading volume attachment: node (ID: %s), volume (ID: %s) in project (%s), region (%s)"
	LogReadSuccess       = "[INFO] Successfully read volume attachment: node (ID: %s), volume (ID: %s)"
	LogNodeNotFound      = "[WARN] Node (ID: %s) not found, removing attachment from state"
	LogVolumeNotFound    = "[WARN] Volume (ID: %s) not found, removing from state"
	LogNotAttached       = "[WARN] Volume (ID: %s) is not attached to node (ID: %s), removing from state"
	LogDetachTemplate    = "[INFO] Detaching volume (ID: %s) from node (ID: %s) in project (%s), region (%s)"
	LogNodeMissingDetach = "[WARN] Node (ID: %s) not found, considering volume detached"
	LogDetachedTemplate  = "[INFO] Successfully detached volume (ID: %s) from node (ID: %s)"
	LogDebugVolumeCheck  = "[DEBUG] Error checking volume status: %s"
	LogDebugNodeCheck    = "[DEBUG] Error checking node status: %s"
	LogWaitAttached      = "[INFO] Volume (ID: %s) successfully attached to node (ID: %s)"
	LogWaitDetached      = "[INFO] Volume (ID: %s) successfully detached from node (ID: %s)"
)
