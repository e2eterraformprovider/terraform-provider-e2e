package loadbalancer

// ============================================================================
// Load Balancer-Specific Constants
// ============================================================================
// This file contains constants specific to the load balancer resource that are
// not part of the API contract or Terraform provider-wide configuration.
// These constants are used internally within the loadbalancer package.
//
// For API-related constants (status values, modes, etc.), see: goe2e/constants/status.go
// For Terraform-wide configuration (timeouts, poll intervals, etc.), see: e2e/constants/internal.go
// For Terraform attribute names, see: e2e/constants/attrs.go

// ============================================================================
// Production Error Message Constants
// ============================================================================
// These constants define error messages used in production code for consistent
// error handling and reporting across the load balancer resource.

const (
	// Import errors
	ErrImportFormat = "import format error: import ID must be in format 'project_id/region/load_balancer_id'"

	// Configuration validation errors
	ErrConflictRegionLocation = "cannot set both 'region' and 'location' parameters"
	ErrMissingNameOrLbName    = "either 'name' or 'lb_name' must be provided"

	// Client creation errors
	ErrCreatingGoe2eClient = "Error creating goe2e client: %s"

	// Load balancer creation errors
	ErrCreatingLoadBalancer = "Error creating load balancer (name: %s) in project (%s), region (%s): %s"
	ErrInsufficientCredits  = "Cannot create load balancer (name: %s) in project (%s), region (%s): insufficient credits. Please add credits to your account"
	ErrWaitingForReady      = "Error waiting for load balancer to become ready: %s"
	ErrSettingTags          = "error setting tags: %w"

	// Load balancer read errors
	ErrRetrievingLoadBalancer = "Error retrieving load balancer (ID: %s) in project (%s), region (%s): %s"
	ErrSettingStatus          = "Error setting load balancer status for ID (%s) in project (%s), region (%s): %s"

	// Load balancer update errors
	ErrInvalidPowerStatus        = "Cannot change power status for load balancer (ID: %s): load balancer is in %s state in project (%s), region (%s). Power status can only be changed when load balancer is in %s or %s state"
	ErrUpdatingPowerStatus       = "Error updating power status for load balancer (ID: %s) in project (%s), region (%s): %s"
	ErrWaitingForPowerAction     = "Error waiting for load balancer power action to complete: %s"
	ErrPlanDowngrade             = "Cannot downgrade plan for load balancer (ID: %s) from %s to %s in project (%s), region (%s): plan downgrades are not supported. Please specify a plan equal to or higher than the current plan"
	ErrUpgradingPlan             = "Error upgrading plan for load balancer (ID: %s) from %s to %s in project (%s), region (%s): %s"
	ErrWaitingForPlanUpgrade     = "Error waiting for load balancer plan upgrade to complete: %s"
	ErrMustBePoweredOn           = "Cannot update load balancer (ID: %s): load balancer is in %s state in project (%s), region (%s). Load balancer must be powered on to update configuration"
	ErrRenamingLoadBalancer      = "Error renaming load balancer (ID: %s) in project (%s), region (%s): %s"
	ErrUpdatingIPv6              = "Error updating IPv6 configuration for load balancer (ID: %s) in project (%s), region (%s): %s"
	ErrUpdatingBackendConfig     = "Error updating backend configuration for load balancer (ID: %s) in project (%s), region (%s): %s"
	ErrInsufficientCreditsUpdate = "Cannot update load balancer (ID: %s) in project (%s), region (%s): insufficient credits. Please add credits to your account"

	// Load balancer delete errors
	ErrInvalidDeleteStatus  = "Cannot delete load balancer (ID: %s): load balancer is in %s state in project (%s), region (%s). Load balancer can only be deleted when not in %s, %s, or %s state"
	ErrDeletingLoadBalancer = "Error deleting load balancer (ID: %s) in project (%s), region (%s): %s"

	// Helper function errors
	ErrInvalidPlanFormat    = "invalid plan name format: %s"
	ErrMissingServerID      = "either 'node_id' or 'id' must be provided for server"
	ErrNodeNotRunning       = "Node with id %s is not in running state"
	ErrVPCNotActive         = "Can not attach vpc currently, vpc is in %s state"
	ErrEmptyLoadBalancerID  = "load balancer ID cannot be empty"
	ErrCreatingRequest      = "failed to create request for Load Balancer (ID: %s): %w"
	ErrLoadBalancerNotFound = "load balancer with ID %s not found"
	ErrRetrievingLBDetails  = "failed to retrieve Load Balancer (ID: %s): %w"
	ErrUnmarshalingResponse = "failed to unmarshal Load Balancer response (ID: %s): %w"

	// Waiting/polling errors
	ErrContextCancelled = "context cancelled while waiting for load balancer %s to reach status %s"
	ErrWaitTimeout      = "timeout waiting for load balancer %s to reach status %s after %d minutes"
)

// ============================================================================
// ACL Condition Constants
// ============================================================================
// ACL condition constants - API values for ACL condition types
// These constants represent the condition types used in load balancer ACL rules
const (
	ACLConditionPathBeg = "path_beg" // Path begins with
	ACLConditionPathEnd = "path_end" // Path ends with
	ACLConditionPathSub = "path_sub" // Path contains substring
	ACLConditionPathDir = "path_dir" // Path directory match
	ACLConditionPathReg = "path_reg" // Path regex match
	ACLConditionHdrBeg  = "hdr_beg"  // Header begins with
	ACLConditionHdrEnd  = "hdr_end"  // Header ends with
	ACLConditionHdrSub  = "hdr_sub"  // Header contains substring
	ACLConditionHdrReg  = "hdr_reg"  // Header regex match
)

// ============================================================================
// Test Data Constants
// ============================================================================
// Test data constants - used in unit tests
// These constants represent common test values used across loadbalancer tests
const (
	TestBackendName      = "test-backend"
	TestServerIP         = "10.0.0.1"
	TestServerPort       = "8080"
	TestBalanceAlgorithm = "roundrobin"
	TestDomainName       = "example.com"
	TestCheckURL         = "/health"
	TestACLName          = "test-acl"
	TestACLPath          = "/api"
	TestBucketName       = "test-bucket"
	TestAccessKey        = "test-access-key"
	TestSecretKey        = "test-secret-key"
)

// ============================================================================
// Validation Arrays
// ============================================================================
// Array constants for validation and testing
var (
	// ValidBalanceAlgorithms lists all valid load balancing algorithms
	// Referenced from goe2e/constants for consistency
	ValidBalanceAlgorithms = []string{
		"source",
		"roundrobin",
		"leastconn",
	}

	// ValidLBModes lists all valid load balancer modes
	// Referenced from goe2e/constants for consistency
	ValidLBModes = []string{
		"HTTP",
		"HTTPS",
		"Both",
	}
)
