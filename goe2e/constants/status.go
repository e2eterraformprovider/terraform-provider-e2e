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

// ============================================================================
// Autoscaling Constants
// ============================================================================

// Autoscaling scaler group status constants - API response values
// These constants represent status values returned by the E2E API for scaler groups.
const (
	AutoscalingScalerGroupStatusRunning  = "Running"
	AutoscalingScalerGroupStatusStopped  = "Stopped"
	AutoscalingScalerGroupStatusStarting = "Starting"
	AutoscalingScalerGroupStatusStopping = "Stopping"

	// Lowercase variants - observed in some normalization/legacy paths and used by V3 schema field `status`.
	AutoscalingScalerGroupStatusRunningLower  = "running"
	AutoscalingScalerGroupStatusStoppedLower  = "stopped"
	AutoscalingScalerGroupStatusStartingLower = "starting"
	AutoscalingScalerGroupStatusStoppingLower = "stopping"
)

// Autoscaling scheduled action type constants - API request/response values
// These constants represent action_type values used for scaler group scheduled actions.
const (
	AutoscalingScheduledActionTypeScaleUp     = "scale_up"
	AutoscalingScheduledActionTypeScaleDown   = "scale_down"
	AutoscalingScheduledActionTypeSetCapacity = "set_capacity"
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
	BlockStorageStatusDetached  = "Detached"
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
// Volume Attachment Constants
// ============================================================================

// Volume attachment status constants - synthetic/SDK values (not always returned by API).
const (
	VolumeAttachmentStatusAttached = "attached"
	VolumeAttachmentStatusDetached = "detached"
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

// Load balancer mode constants - API request/response values
// These constants represent the mode values used by the E2E API for load balancers
const (
	LBModeHTTP  = "HTTP"  // HTTP mode
	LBModeHTTPS = "HTTPS" // HTTPS mode
	LBModeBoth  = "Both"  // Both HTTP and HTTPS mode
)

// Load balancer port constants - API request/response values
// These constants represent standard port values used by load balancers
const (
	LBPortHTTP  = "80"  // Standard HTTP port
	LBPortHTTPS = "443" // Standard HTTPS port
)

// Load balancer balance algorithm constants - API request/response values
// These constants represent the load balancing algorithm values used by the E2E API
const (
	LBBalanceSource     = "source"     // Source-based load balancing (sticky sessions)
	LBBalanceRoundRobin = "roundrobin" // Round-robin load balancing
	LBBalanceLeastConn  = "leastconn"  // Least connections load balancing
)

// Load balancer type constants - API request/response values
// These constants represent the type values used by the E2E API for load balancers
const (
	LBTypeExternal = "External" // External load balancer (public-facing)
	LBTypeInternal = "Internal" // Internal load balancer (private)
)

// Load balancer node list type constants - API request values
// These constants represent the node list type values used by the E2E API for load balancers
const (
	LBNodeListTypeStatic  = "S" // Static node list
	LBNodeListTypeDynamic = "D" // Dynamic node list (autoscaling)
)

// Load balancer plan constants - API request/response values
// These constants represent the plan name values used by the E2E API for load balancers
const (
	LBPlanE2ELB2 = "E2E-LB-2" // E2E Load Balancer Plan 2
	LBPlanE2ELB3 = "E2E-LB-3" // E2E Load Balancer Plan 3
	LBPlanE2ELB4 = "E2E-LB-4" // E2E Load Balancer Plan 4
	LBPlanE2ELB5 = "E2E-LB-5" // E2E Load Balancer Plan 5
)

// Load balancer action constants - API action values
// These constants represent action types used in API requests for load balancer operations
const (
	LBActionRename      = "rename"       // Action type for renaming a load balancer
	LBActionUpgradePlan = "upgrade_plan" // Action type for upgrading load balancer plan
)

// Load balancer IPv6 action constants - API action values
// These constants represent IPv6 action types used in API requests for load balancer IPv6 operations
const (
	LBIPv6ActionAttach = "attach" // Action type for attaching IPv6 to load balancer
	LBIPv6ActionDetach = "detach" // Action type for detaching IPv6 from load balancer
)

// Load balancer default values - API default values
// These constants represent default values used by the E2E API for load balancer configuration
const (
	LBDefaultDomainName = "localhost" // Default domain name for health checks
	LBDefaultCheckURL   = "/"         // Default health check URL path
)

// ============================================================================
// DBaaS Constants
// ============================================================================

// DBaaS status constants - API response values
// These constants represent the various status values returned by the E2E API for DBaaS clusters
// (PostgreSQL, MySQL, MariaDB, etc.)
const (
	// DBaaSStatusCreating indicates the DBaaS cluster is still being created (API response value).
	DBaaSStatusCreating = "CREATING"

	DBaaSStatusActive     = "ACTIVE"     // DBaaS cluster is active and running
	DBaaSStatusSuspended  = "SUSPENDED"  // DBaaS cluster is suspended
	DBaaSStatusSuspending = "SUSPENDING" // DBaaS cluster is transitioning to suspended
	DBaaSStatusResuming   = "RESUMING"   // DBaaS cluster is resuming from suspended state
	DBaaSStatusStopped    = "STOPPED"    // DBaaS cluster is stopped
	DBaaSStatusRunning    = "RUNNING"    // DBaaS cluster is running
	DBaaSStatusRestarting = "RESTARTING" // DBaaS cluster is restarting
)

// ============================================================================
// DBaaS Software and Version Constants
// ============================================================================
// These constants represent values sent to / returned by the E2E API for DBaaS
// software identification and version selection.

const (
	// DBaaSSoftwarePostgreSQL is the API software name for PostgreSQL DBaaS.
	DBaaSSoftwarePostgreSQL = "PostgreSQL"
	// DBaaSSoftwareMySQL is the API software name for MySQL DBaaS.
	DBaaSSoftwareMySQL = "MySQL"
)

// PostgreSQL supported versions (API values).
const (
	PostgreSQLVersion11 = "11.0"
	PostgreSQLVersion12 = "12.0"
	PostgreSQLVersion13 = "13.0"
	PostgreSQLVersion14 = "14.0"
	PostgreSQLVersion15 = "15.0"
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

// FaaS default configuration constants - API default values
// These constants represent default values used by the E2E API for FaaS functions.
// These are API-related defaults that apply when values are not explicitly provided.
const (
	// FaaSDefaultMemoryMB is the default memory allocation in megabytes for a FaaS function.
	// This is the default value used by the API when memory_mb is not specified.
	FaaSDefaultMemoryMB = 256

	// FaaSDefaultTimeoutSeconds is the default execution timeout in seconds for a FaaS function.
	// This is the default value used by the API when timeout_seconds is not specified.
	FaaSDefaultTimeoutSeconds = 30

	// FaaSDefaultMinReplicas is the default minimum number of replicas for a FaaS function.
	// This is the default value used by the API when min_replicas is not specified.
	FaaSDefaultMinReplicas = 1

	// FaaSDefaultMaxReplicas is the default maximum number of replicas for a FaaS function.
	// This is the default value used by the API when max_replicas is not specified.
	FaaSDefaultMaxReplicas = 5
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

// Reserve IP type constants - API response values
// These constants represent the various type values returned by the E2E API for reserved IPs
const (
	ReserveIPTypeFloatingIP = "FloatingIP"
	ReserveIPTypePublicIP   = "PublicIP"
	ReserveIPTypeAddonIP    = "AddonIP"
)

// ============================================================================
// SFS (Shared File System) Constants
// ============================================================================

// SFS status constants - API response values
// These constants represent the various status values returned by the E2E API for SFS instances
const (
	// API status values (as returned by API)
	SFSStatusCreating = "Creating" // SFS is being created
	SFSStatusActive   = "Active"   // SFS is active and ready
	SFSStatusDeleting = "Deleting" // SFS is being deleted
	SFSStatusDeleted  = "Deleted"  // SFS has been deleted
	SFSStatusError    = "Error"    // SFS is in error state

	// Normalized status values (used internally for state management)
	SFSStateCreating = "creating" // Normalized: Creating
	SFSStateActive   = "active"   // Normalized: Active
	SFSStateDeleting = "deleting" // Normalized: Deleting
	SFSStateDeleted  = "deleted"  // Normalized: Deleted
	SFSStateError    = "error"    // Normalized: Error
)

// SFS desired status constants - used in wait operations
// These constants represent desired states for SFS wait/polling operations
const (
	SFSDesiredStatusDeleted = "deleted" // Desired state: deleted
	SFSDesiredStatus404     = "404"     // Desired state: not found (404)
	SFSDesiredStatusActive  = "active"  // Desired state: active
)

// ============================================================================
// Object Storage Constants
// ============================================================================

// Object storage versioning status constants - API response/request values
// These constants represent the versioning status values returned by the E2E API for object storage buckets
const (
	ObjectStorageVersioningStatusEnabled   = "Enabled"   // Versioning is enabled
	ObjectStorageVersioningStatusSuspended = "Suspended" // Versioning is suspended
)

// ============================================================================
// Image Constants
// ============================================================================

// Image status constants - API response values
// These constants represent the various status values returned by the E2E API for images
const (
	// API status values (as returned by API)
	ImageStatusCreating = "Creating" // Image is being created
	ImageStatusReady    = "Ready"    // Image is ready for use
	ImageStatusError    = "Error"    // Image is in an error state
	ImageStatusDeleted  = "Deleted"  // Image has been deleted

	// Normalized status values (used internally for state management)
	ImageStateCreating = "creating" // Normalized: Creating
	ImageStateReady    = "ready"    // Normalized: Ready
	ImageStateError    = "error"    // Normalized: Error
	ImageStateDeleted  = "deleted"  // Normalized: Deleted
)

// Image action constants - API action values
// These constants represent action types used in API requests for image operations
const (
	ImageActionSaveImages = "save_images" // Action type for saving an image from a node
	ImageActionRename     = "rename"      // Action type for renaming an image
)

// ============================================================================
// VPC Constants
// ============================================================================

// VPC status constants - API response values
// These constants represent the various status values returned by the E2E API for VPCs
const (
	VPCStatusActive = "Active" // VPC is active and ready for use
)

// ============================================================================
// Container Registry Constants
// ============================================================================

// Container Registry severity constants - API request/response values
// These constants represent the severity level values used by the E2E API
// for vulnerability scanning in container registries
const (
	ContainerRegistrySeverityLow      = "low"      // Low severity vulnerabilities
	ContainerRegistrySeverityMedium   = "medium"   // Medium severity vulnerabilities
	ContainerRegistrySeverityHigh     = "high"     // High severity vulnerabilities
	ContainerRegistrySeverityCritical = "critical" // Critical severity vulnerabilities
	ContainerRegistrySeverityNone     = "none"     // No severity threshold (disabled)
)
