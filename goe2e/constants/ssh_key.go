package constants

// ============================================================================
// SSH Key API Constants
// ============================================================================

// SSH Key API endpoint paths - used in goe2e SDK client implementation
const (
	// SSHKeysPath is the API endpoint for SSH key operations (list, create, get)
	SSHKeysPath = "ssh_keys"

	// DeleteSSHKeyPath is the API endpoint for deleting an SSH key
	DeleteSSHKeyPath = "delete_ssh_key"
)

// SSH Key API error indicators
// These constants represent error message substrings returned by the E2E API
// when SSH key operations fail, used for error handling and diagnostics
const (
	// SSHKeyNotFoundSubstring is the substring that appears in API error messages
	// when an SSH key cannot be found by ID or label
	SSHKeyNotFoundSubstring = "not found"
)
