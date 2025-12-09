package constants

// ============================================================================
// Node Constants
// ============================================================================

// Node status constants - API response values
// These constants represent the various status values returned by the E2E API for nodes
const (
	NodeStatusRunning      = "Running"
	NodeStatusReinstalling = "Reinstalling"
	NodeStatusCreating     = "Creating"
	NodeStatusFailed       = "Failed"
	NodeStatusPoweredOff   = "Powered off"
	NodeStatusSaving       = "Saving"
)

// Node power transition status constants - API response values
// These constants represent transient status values returned by the E2E API during node power operations
const (
	NodeStatusPowering = "Powering" // Node is powering on
	NodeStatusStopping = "Stopping" // Node is stopping/powering off
)

// Node power status constants - API action/response values
const (
	NodePowerStatusOn  = "power_on"
	NodePowerStatusOff = "power_off"
)

// Node action constants - API action values
const (
	NodeActionReboot   = "reboot"
	NodeActionLockVM   = "lock_vm"
	NodeActionUnlockVM = "unlock_vm"
)

// Node LCM state constants - API response values
// These constants represent lifecycle management states returned by the E2E API for nodes.
// LCM (Life Cycle Management) states indicate transient operational states during node operations
// such as volume attachment/detachment, disk resizing, etc.
const (
	NodeLCMStateHotplugPrologPoweroff = "HOTPLUG_PROLOG_POWEROFF"
	NodeLCMStateHotplugEpilogPoweroff = "HOTPLUG_EPILOG_POWEROFF"
	NodeLCMStateHotplug               = "Hotplug"
	NodeLCMStateDiskResize            = "DISK_RESIZE"
	NodeLCMStateDiskResizePoweroff    = "DISK_RESIZE_POWEROFF"
)

// ============================================================================
// Block Storage Constants
// ============================================================================

// Block storage status constants - API response values
const (
	BlockStorageStatusAttached  = "Attached"
	BlockStorageStatusSaving    = "Saving"
	BlockStorageStatusCreating  = "Creating"
	BlockStorageStatusAvailable = "Available"
	BlockStorageStatusError     = "ERROR"
)

// Block storage action constants - API action values
const (
	BlockStorageActionAttach = "attach"
	BlockStorageActionDetach = "detach"
)

// ============================================================================
// Load Balancer Constants
// ============================================================================

// Load balancer status constants - API response values
// These constants represent the various status values returned by the E2E API for load balancers
const (
	// API status values (as returned by API)
	LBStatusCreating           = "Creating"                   // Load balancer is being created
	LBStatusDeploying          = "Deploying"                  // Load balancer is being deployed
	LBStatusRunning            = "Running"                    // Load balancer is running normally
	LBStatusRunningAPI         = "RUNNING"                    // API returns RUNNING (uppercase) in lb_status.status
	LBStatusPoweredOff         = "Powered off"                // Load balancer is powered off
	LBStatusPoweredOffAPI      = "STOP"                       // API returns STOP (uppercase) in lb_status.status
	LBStatusUpgrading          = "Upgrading"                  // Load balancer plan is being upgraded
	LBStatusUpgradingAPI       = "UPDATING"                   // API returns UPDATING (uppercase) in lb_status.status
	LBStatusError              = "Error"                      // Load balancer is in error state
	LBStatusBackendUnavailable = "Backend Status Unavailable" // Backend status cannot be determined
	LBStatusBackendFailure     = "Backend Connection Failure" // Backend connection failed

	// Normalized status values (used internally for state management)
	LBStateCreating  = "creating"  // Normalized: Creating, Deploying
	LBStateRunning   = "running"   // Normalized: Running
	LBStateStopped   = "stopped"   // Normalized: Powered off, STOP
	LBStateUpgrading = "upgrading" // Normalized: Upgrading, UPDATING
	LBStateError     = "error"     // Normalized: Error, Failed
)

// ============================================================================
// DBaaS Constants
// ============================================================================

// DBaaS status constants - API response values
// These constants represent the various status values returned by the E2E API for DBaaS clusters
// (PostgreSQL, MySQL, MariaDB, etc.)
const (
	DBaaSStatusActive     = "ACTIVE"     // DBaaS cluster is active and running
	DBaaSStatusSuspended  = "SUSPENDED"  // DBaaS cluster is suspended
	DBaaSStatusSuspending = "SUSPENDING" // DBaaS cluster is transitioning to suspended
	DBaaSStatusResuming   = "RESUMING"   // DBaaS cluster is resuming from suspended state
	DBaaSStatusStopped    = "STOPPED"    // DBaaS cluster is stopped
	DBaaSStatusRunning    = "RUNNING"    // DBaaS cluster is running
	DBaaSStatusRestarting = "RESTARTING" // DBaaS cluster is restarting
)

// ============================================================================
// FaaS Constants
// ============================================================================

// FaaS function status constants - API response values
// These constants represent the various status values returned by the E2E API for FaaS functions
const (
	FaaSStatusBuilding = "Building" // FaaS function is being built
	FaaSStatusUpdating = "Updating" // FaaS function is being updated
	FaaSStatusPending  = "Pending"  // FaaS function is pending
	FaaSStatusReady    = "Ready"    // FaaS function is ready
	FaaSStatusError    = "Error"    // FaaS function is in error state
	FaaSStatusFailed   = "Failed"   // FaaS function has failed
)

// ============================================================================
// Reserve IP Constants
// ============================================================================

// Reserve IP status constants - API response values
// These constants represent the various status values returned by the E2E API for reserved IPs
const (
	ReserveIPStatusAvailable = "Available"
	ReserveIPStatusAttached  = "Attached"
)
