package constants

// ============================================================================
// API error indicators
// ============================================================================
// These constants represent error message substrings and codes returned by the E2E API
// when resources cannot be found, used for error handling and diagnostics

// NotFoundSubstring is the substring that appears in API error messages
// when a resource cannot be found by ID or other identifier
const NotFoundSubstring = "not found"

// NotFoundCode is the HTTP status code that indicates a resource was not found
const NotFoundCode = "404"

// ============================================================================
// SFS API Error Messages
// ============================================================================
// These constants represent error messages related to SFS API operations

// SFSNotFound is the error message format when an SFS resource cannot be found via API
// Format: "SFS %s not found" where %s is the SFS ID
const SFSNotFound = "SFS %s not found"

// SFSEnteredErrorState is the error message format when an SFS enters an error state during an operation
// Format: "SFS %s entered error state during %s operation" where %s are SFS ID and operation name
const SFSEnteredErrorState = "SFS %s entered error state during %s operation"

// SFSTimeoutWaitingForStatus is the error message format when waiting for SFS status times out
// Format: "timeout waiting for SFS %s to reach status %s" where %s are SFS ID and desired status
const SFSTimeoutWaitingForStatus = "timeout waiting for SFS %s to reach status %s"

// ClientOrServiceNil is the error message when client or service is nil (API-related validation)
const ClientOrServiceNil = "client or Sfs service is nil"

// ============================================================================
// Image API Error Messages
// ============================================================================
// These constants represent error messages related to Image API operations

// ImageTimeoutWaitingForState is the error message format when waiting for Image state times out
// Format: "timeout waiting for image %s to reach state %s" where %s are image ID and desired state
const ImageTimeoutWaitingForState = "timeout waiting for image %s to reach state %s"

// ImageEnteredErrorState is the error message format when an Image enters an error state
// Format: "image %s entered error state" where %s is the image ID
const ImageEnteredErrorState = "image %s entered error state"
