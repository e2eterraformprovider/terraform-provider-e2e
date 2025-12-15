package reserve_ip

// ============================================================================
// Reserve IP Resource Name Constant
// ============================================================================

// ResourceName is the resource type name for Reserve IP, used in error messages
const ResourceName = "Reserve IP"

// ============================================================================
// Reserve IP-Specific Error Messages
// ============================================================================
// These are Reserve IP-specific errors that cannot be generalized due to:
// - Resource-specific import format (project_id/region/ip_address)
// - Resource-specific warning messages

// WarnReserveIPAttached is Reserve IP-specific warning when deleting an attached IP
// Format: "Reserved IP (%s) is currently attached. The API will handle detachment automatically."
const WarnReserveIPAttached = "Reserved IP (%s) is currently attached. The API will handle detachment automatically."
