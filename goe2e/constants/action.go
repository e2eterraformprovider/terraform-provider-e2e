package constants

// ============================================================================
// Generic API Action Constants
// ============================================================================
// These constants represent common API request action values used across services.

const (
	ActionAttach = "attach"
	ActionDetach = "detach"
)

// ============================================================================
// Node Action Constants
// ============================================================================
// These constants represent action values used in API requests for node operations

const (
	// NodeActionAddSSHKeys is the action type for adding SSH keys to a node
	NodeActionAddSSHKeys = "add_ssh_keys"

	// NodeActionReinstall is the action type for reinstalling a node's OS
	NodeActionReinstall = "reinstall"
)
