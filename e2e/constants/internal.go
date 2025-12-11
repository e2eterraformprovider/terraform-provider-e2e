package constants

import (
	"time"

	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

// Provider-specific constants (not part of API contract)
//
// These constants are used internally by the Terraform provider and are NOT part of the API SDK.
// They represent Terraform-specific configuration, business logic, or implementation details
// that are not exposed through the goe2e API client.
//
// When to add constants here:
//   - Terraform-specific configuration values (e.g., plan prefixes, resource naming patterns)
//   - Provider-specific timeout/retry configuration
//   - Internal state management values not returned by the API
//   - Terraform schema attribute names or internal identifiers
//
// When NOT to add constants here:
//   - API response values (status codes, state strings returned by API) → goe2e/constants/
//   - API request parameters → goe2e/constants/
//   - Values that represent the API contract → goe2e/constants/
//
// If a constant represents a value returned by or sent to the API, it belongs in goe2e/constants/
// to ensure consistency across all API consumers, not just the Terraform provider.

const (
	// PREFIX_C2_NODE is the prefix used to identify C2 node plans in Terraform resource logic.
	// This is used for Terraform-specific plan validation and resource behavior.
	PREFIX_C2_NODE = "C2"

	// ============================================================================
	// SSH Key Import and Error Handling Constants
	// ============================================================================
	// These are Terraform provider-specific constants for SSH key resource import
	// and error checking, not part of the API contract.

	// SSHKeyImportIDFormatDescription describes the expected format for importing SSH keys
	// into Terraform. Users must provide either two or three colon-separated parts.
	SSHKeyImportIDFormatDescription = "project_id:ssh_key_id or project_id:region:ssh_key_id"

	// SSHKeyImportIDRegionRequired is the error message shown when region is not specified
	// in the import ID and no default region is configured in the provider
	SSHKeyImportIDRegionRequired = "region must be specified in import ID or provider default_region must be set. Use format: project_id:region:ssh_key_id"

	// SSHKeyImportIDInvalidFormat is the error message shown when the import ID format is invalid
	SSHKeyImportIDInvalidFormat = "invalid import ID format. Expected: project_id:ssh_key_id or project_id:region:ssh_key_id"

	// SSHKeyNotFoundCheckSubstring is used to detect "not found" errors in API responses.
	// This substring appears in error messages when an SSH key resource cannot be found.
	SSHKeyNotFoundCheckSubstring = "not found"
)

// State change wait configuration constants
// These constants control polling behavior when waiting for resources to transition between states.
// All durations are typed as time.Duration to avoid unit mixing at callsites.
const (
	// StateChangePollInterval is the interval between polling attempts when waiting for state changes.
	// This controls how frequently the provider checks if a resource has transitioned states.
	// This is used as the 'delay' parameter in WaitForState.
	StateChangePollInterval = 3 * time.Second

	// StateChangeRetryBackoff is the minimum time between polling attempts.
	// This prevents excessive API calls during rapid state transitions.
	// This is used as the 'minTimeout' parameter in WaitForState.
	StateChangeRetryBackoff = 3 * time.Second

	// StateChangeDefaultDelay is the default delay between polling attempts when not specified.
	// This is used as the fallback default in WaitForState when delay parameter is 0.
	StateChangeDefaultDelay = 10 * time.Second

	// StateChangeTimeoutShort is used for operations expected to complete quickly (e.g., DBaaS suspend).
	StateChangeTimeoutShort = 5 * time.Minute

	// StateChangeTimeoutDefault is used for standard resource state transitions (e.g., node power state).
	StateChangeTimeoutDefault = 10 * time.Minute

	// DefaultNotFoundChecks is the number of "not found" responses to allow before failing.
	// This is used as the fallback default in WaitForState when notFoundChecks parameter is 0.
	DefaultNotFoundChecks = 60

	// NodeLCMReadyState is a synthetic state used internally by the Terraform provider
	// to indicate that a node LCM operation has completed (i.e., the LCM state is not
	// in any of the excluded pending states). This is NOT an API value - it's created
	// by the wait logic when the LCM state exits the pending states.
	NodeLCMReadyState = "ready"

	// WaitStateDeleted is a synthetic state used internally by Terraform wait logic
	// to indicate that a resource was not found (404 response). This is NOT an API value
	// but is used internally for state machine logic in wait operations.
	WaitStateDeleted = "deleted"

	// WaitStateUnknown is a synthetic state used internally by Terraform wait logic
	// to indicate that a resource status cannot be determined (e.g., empty status field).
	// This is NOT an API value but is used internally for state machine logic in wait operations.
	WaitStateUnknown = "unknown"
)

// Wait state groups
// These represent common state combinations used in wait operations.
// They are grouped by domain/operation type for clarity and reusability.
// Note: These reference API constants from goe2e/constants/status.go to ensure consistency.
var (
	// DBaaSSuspendPendingStates are the states that indicate a DBaaS instance is transitioning to suspended.
	// These reference API status constants from goe2e/constants/status.go.
	DBaaSSuspendPendingStates = []string{
		goe2econstants.DBaaSStatusActive,
		goe2econstants.DBaaSStatusSuspending,
		goe2econstants.DBaaSStatusResuming,
	}

	// DBaaSSuspendTargetStates are the target states when suspending a DBaaS instance.
	// These reference API status constants from goe2e/constants/status.go.
	DBaaSSuspendTargetStates = []string{
		goe2econstants.DBaaSStatusSuspended,
	}

	// NodeLCMReadyTargetStates are the target states when waiting for node LCM operations to complete.
	// This uses a synthetic state (NodeLCMReadyState) indicating the LCM state is not in an excluded pending state.
	NodeLCMReadyTargetStates = []string{NodeLCMReadyState}

	// FaaSPendingStates are the states that indicate a FaaS function is being prepared or updated.
	// These reference API status constants from goe2e/constants/status.go.
	FaaSPendingStates = []string{
		goe2econstants.FaaSStatusBuilding,
		goe2econstants.FaaSStatusUpdating,
		goe2econstants.FaaSStatusPending,
	}

	// NodePowerPendingStates are the states that indicate a node is transitioning power state.
	// These reference API status constants from goe2e/constants/status.go.
	NodePowerPendingStates = []string{
		goe2econstants.NodeStatusPowering,
		goe2econstants.NodeStatusStopping,
	}
)
