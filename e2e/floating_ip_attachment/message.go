package floating_ip_attachment

// ============================================================================
// Floating IP Attachment Import Format Constants
// ============================================================================

// FloatingIPAttachmentImportFormat describes the expected import ID format for floating IP attachments.
// Format: "project_id/region/ip_address"
const FloatingIPAttachmentImportFormat = "project_id/region/ip_address"

// ============================================================================
// Floating IP Attachment-Specific Error Messages
// ============================================================================
// These are floating IP attachment-specific errors that cannot be generalized due to:
// - Resource-specific validation rules
// - Resource-specific error context

// ErrNodeIDsCannotBeEmpty is the error message when node_ids is empty during creation.
const ErrNodeIDsCannotBeEmpty = "node_ids cannot be empty"

// ErrNodeIDsCannotBeEmptyWithContext is the error message when node_ids becomes empty during update.
// This provides additional context that a floating IP attachment must have at least one node attached.
const ErrNodeIDsCannotBeEmptyWithContext = "node_ids cannot be empty. A floating IP attachment must have at least one node attached"

// ErrAttachingFloatingIP is the error message format when attaching a floating IP fails.
// Format: "Error attaching floating IP (%s) to nodes in project (%s), region (%s): %s"
const ErrAttachingFloatingIP = "Error attaching floating IP (%s) to nodes in project (%s), region (%s): %s"

// ErrDetachingFloatingIP is the error message format when detaching a floating IP fails.
// Format: "Error detaching floating IP (%s) from nodes in project (%s), region (%s): %s"
const ErrDetachingFloatingIP = "Error detaching floating IP (%s) from nodes in project (%s), region (%s): %s"

// ErrRetrievingReservedIPs is the error message format when retrieving reserved IPs fails.
// Format: "Error retrieving reserved IPs in project (%s), region (%s): %s"
const ErrRetrievingReservedIPs = "Error retrieving reserved IPs in project (%s), region (%s): %s"

// ErrReadingFloatingIPAttachmentDuringImport is the error message when reading fails during import.
// Format: "error reading floating IP attachment during import: %v"
const ErrReadingFloatingIPAttachmentDuringImport = "error reading floating IP attachment during import: %v"
