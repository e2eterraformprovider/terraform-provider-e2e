package ssh_key

import tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"

// ============================================================================
// SSH Key Resource Name Constant
// ============================================================================

// ResourceName is the human-readable resource name used in errors/logs.
const ResourceName = "SSH Key"

// ============================================================================
// SSH Key Import / Validation Messages
// ============================================================================
//
// Historically these lived in `e2e/constants/message.go`. We alias them here so the ssh_key
// package can follow the same `e2e/<resource>/message.go` convention without breaking
// existing references to tfconstants.

const (
	ImportIDFormatDescription = tfconstants.SSHKeyImportIDFormatDescription
	ImportIDRegionRequired    = tfconstants.SSHKeyImportIDRegionRequired
	ImportIDInvalidFormat     = tfconstants.SSHKeyImportIDInvalidFormat
)
