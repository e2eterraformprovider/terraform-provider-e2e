package faas

// ============================================================================
// FaaS Resource Name Constant
// ============================================================================

// ResourceName is the human-readable resource name used in logs and errors.
const ResourceName = "FaaS Function"

// ============================================================================
// FaaS Validation / Error / Warning Messages
// ============================================================================
// Keep these stable because unit/acceptance tests and user experience depend on them.

const (
	// Schema / config validation
	ErrMinReplicasGreaterThanMaxFmt = "min_replicas (%d) cannot be greater than max_replicas (%d)"
	ErrCodeInlineOrFileRequired     = "one of code_inline or code_file must be specified"
	ErrCodeInlineAndFileExclusive   = "code_inline and code_file are mutually exclusive"

	// Read-after-write invariants
	ErrFunctionNotFoundAfterCreate = "function not found after creation"
	ErrFunctionNotFoundAfterUpdate = "function not found after update"

	// Data source / read
	ErrFunctionNotFoundByIDFmt = "FaaS function with ID %s not found"

	// Warnings / logs
	LogNamespaceCreateWarningFmt  = "[WARN] Namespace creation returned error (may already exist): %v"
	LogCodeFileNotImplemented     = "[WARN] code_file support for actual file upload is not yet implemented. Treating as inline code."
	LogFunctionNotFoundWarningFmt = "[WARN] FaaS function with ID %s not found"
)
