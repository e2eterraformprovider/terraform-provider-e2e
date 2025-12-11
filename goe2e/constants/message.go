package constants

// ============================================================================
// SSH Key API error indicators
// ============================================================================
// These constants represent error message substrings returned by the E2E API
// when SSH key operations fail, used for error handling and diagnostics

// SSHKeyNotFoundSubstring is the substring that appears in API error messages
// when an SSH key cannot be found by ID or label
const SSHKeyNotFoundSubstring = "not found"
