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

// FaaS validation limits - Terraform-specific validation constraints
// These constants represent validation limits enforced by the Terraform provider.
// They are Terraform-specific business rules, not API contract values.
const (
	// FaaSMinMemoryMB is the minimum memory allocation in megabytes allowed for a FaaS function.
	// This is enforced by Terraform validation, not by the API.
	FaaSMinMemoryMB = 128

	// FaaSMinTimeoutSeconds is the minimum execution timeout in seconds allowed for a FaaS function.
	// This is enforced by Terraform validation, not by the API.
	FaaSMinTimeoutSeconds = 1

	// FaaSMaxTimeoutSeconds is the maximum execution timeout in seconds allowed for a FaaS function.
	// This is enforced by Terraform validation, not by the API.
	FaaSMaxTimeoutSeconds = 900
)

// FaaS operation timeout constants - Terraform-specific timeout configuration
// These constants control how long the provider waits for FaaS operations to complete.
// All durations are typed as time.Duration to avoid unit mixing at callsites.
const (
	// FaaSCreateTimeout is the maximum time to wait for a FaaS function creation to complete.
	// This includes the time for the function to transition from Building/Pending to Ready state.
	FaaSCreateTimeout = 10 * time.Minute

	// FaaSUpdateTimeout is the maximum time to wait for a FaaS function update to complete.
	// This includes the time for the function to transition from Updating to Ready state.
	FaaSUpdateTimeout = 10 * time.Minute

	// FaaSDeleteTimeout is the maximum time to wait for a FaaS function deletion to complete.
	FaaSDeleteTimeout = 2 * time.Minute

	// FaaSPollInterval is the interval between polling attempts when waiting for FaaS function state changes.
	// This controls how frequently the provider checks if a FaaS function has transitioned states.
	FaaSPollInterval = 10 * time.Second
)

// Load Balancer operation timeout constants - Terraform-specific timeout configuration
// These constants control how long the provider waits for Load Balancer operations to complete.
// All durations are typed as time.Duration to avoid unit mixing at callsites.
const (
	// LBCreateTimeout is the maximum time to wait for a load balancer creation to complete.
	// This includes the time for the load balancer to transition from Creating to Running state.
	LBCreateTimeout = 15 * time.Minute

	// LBDeleteTimeout is the maximum time to wait for a load balancer deletion to complete.
	// This includes polling until the load balancer returns 404 (not found).
	LBDeleteTimeout = 10 * time.Minute

	// LBPowerActionTimeout is the maximum time to wait for load balancer power actions to complete.
	// This includes power on/off operations and the transition to the target power state.
	LBPowerActionTimeout = 5 * time.Minute

	// LBPlanUpgradeTimeout is the maximum time to wait for a load balancer plan upgrade to complete.
	// This includes the time for the load balancer to transition back to Running state after upgrade.
	LBPlanUpgradeTimeout = 10 * time.Minute

	// LBPollInterval is the interval between polling attempts when waiting for load balancer state changes.
	// This controls how frequently the provider checks if a load balancer has transitioned states.
	LBPollInterval = 10 * time.Second
)

// Load Balancer validation constants - Terraform-specific validation constraints
// These constants represent validation constraints enforced by the Terraform provider for load balancers.
var (
	// LBTCPDisallowedPorts is the list of ports that are not allowed for TCP backends.
	// These ports are reserved for internal use or conflict with HTTP/HTTPS frontends.
	LBTCPDisallowedPorts = []string{"8080", "10050", "9101", "80", "443"}
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

	// AutoscalingScalerGroupProvisionStatusAllowed are the allowed values for the legacy V2
	// `provision_status` field on autoscaling scaler groups.
	// These values are part of the API contract, but the *list* is a Terraform validation construct.
	AutoscalingScalerGroupProvisionStatusAllowed = []string{
		goe2econstants.AutoscalingScalerGroupStatusRunning,
		goe2econstants.AutoscalingScalerGroupStatusStopped,
	}

	// AutoscalingScalerGroupStatusAllowed are the allowed values for the V3 `status` field on
	// autoscaling scaler groups (lowercase variants).
	AutoscalingScalerGroupStatusAllowed = []string{
		goe2econstants.AutoscalingScalerGroupStatusRunningLower,
		goe2econstants.AutoscalingScalerGroupStatusStoppedLower,
	}

	// AutoscalingScalerGroupRunningStates are states treated as "running" for Terraform-side
	// state requirement checks (handles both legacy and V3 values).
	AutoscalingScalerGroupRunningStates = []string{
		goe2econstants.AutoscalingScalerGroupStatusRunning,
		goe2econstants.AutoscalingScalerGroupStatusRunningLower,
	}

	// AutoscalingScalerGroupStoppedStates are states treated as "stopped" for Terraform-side
	// state requirement checks (handles both legacy and V3 values).
	AutoscalingScalerGroupStoppedStates = []string{
		goe2econstants.AutoscalingScalerGroupStatusStopped,
		goe2econstants.AutoscalingScalerGroupStatusStoppedLower,
	}
)

// SFS field migration keys - Terraform-specific field names used in state migration
// These are used internally by the Terraform provider for state upgrade logic
const (
	// FieldMigrationKeyDiskSize is the V2 field name for disk size (migrated to size_gb in V3)
	FieldMigrationKeyDiskSize = "disk_size"
	// FieldMigrationKeySizeGB is the V3 field name for disk size
	FieldMigrationKeySizeGB = "size_gb"
	// FieldMigrationKeyDiskIOPS is the V2 field name for disk IOPS (migrated to iops in V3)
	FieldMigrationKeyDiskIOPS = "disk_iops"
	// FieldMigrationKeyIOPS is the V3 field name for IOPS
	FieldMigrationKeyIOPS = "iops"
	// FieldMigrationKeyIsEncryptionEnabled is the V2 field name for encryption flag (migrated to encryption_enabled in V3)
	FieldMigrationKeyIsEncryptionEnabled = "is_encryption_enabled"
	// FieldMigrationKeyEncryptionEnabled is the V3 field name for encryption flag
	FieldMigrationKeyEncryptionEnabled = "encryption_enabled"
)

// ============================================================================
// DBaaS Terraform Defaults and Migration Keys
// ============================================================================
// These constants are Terraform-provider specific (schema defaults, state upgrade keys)
// and are not part of the goe2e API contract.

const (
	// DBaaSImportIDFormatDescription describes the expected import ID format for DBaaS resources.
	DBaaSImportIDFormatDescription = "project_id:dbaas_id"

	// DBaaSImportIDSeparator is the separator used in DBaaS import IDs.
	DBaaSImportIDSeparator = ":"
)

// DBaaS Terraform schema defaults.
const (
	DBaaSDefaultGroupName            = "Default"
	DBaaSDefaultDBaaSNumber          = 1
	DBaaSDefaultPublicIPRequired     = true
	DBaaSDefaultPublicIPEnabled      = true
	DBaaSDefaultIsEncryptionEnabled  = false
	DBaaSDefaultEncryptionPassphrase = ""
	DBaaSDefaultParameterGroupID     = 0
	DBaaSDefaultDBLocation           = "Delhi"
)

// DBaaS power action constants used in legacy state and some schema validation paths.
// These are Terraform/provider values (not API statuses).
const (
	DBaaSPowerActionStart   = "start"
	DBaaSPowerActionStop    = "stop"
	DBaaSPowerActionRestart = "restart"
)

// DBaaS state upgrade keys (legacy schema field names).
const (
	FieldMigrationKeyVPCList        = "vpc_list"
	FieldMigrationKeyDetachPublicIP = "detach_public_ip"
	FieldMigrationKeyPowerStatus    = "power_status"
)

// Container Registry Terraform defaults - Terraform-specific default values
// These constants represent default values used by the Terraform provider schema,
// not API defaults. They are enforced at the Terraform layer.
const (
	// ContainerRegistryDefaultSeverity is the default vulnerability severity threshold
	// used in the Terraform schema when severity is not specified.
	ContainerRegistryDefaultSeverity = "low"

	// ContainerRegistryDefaultPreventVul is the default value for prevent_vulnerabilities
	// used in the Terraform schema when prevent_vul is not specified.
	ContainerRegistryDefaultPreventVul = false

	// ContainerRegistryDefaultCustomerID is the default customer/user ID used when
	// the customer ID cannot be determined from the API response.
	ContainerRegistryDefaultCustomerID = "0"
)

// ============================================================================
// Kubernetes Terraform Defaults and Configuration
// ============================================================================
// These constants are Terraform-provider specific (timeouts, import formats, validation)
// and are not part of the goe2e API contract.

// Kubernetes operation timeout constants - Terraform-specific timeout configuration
// These constants control how long the provider waits for Kubernetes operations to complete.
// All durations are typed as time.Duration to avoid unit mixing at callsites.
const (
	// KubernetesCreateTimeout is the maximum time to wait for a Kubernetes cluster creation to complete.
	// This includes the time for the cluster to transition from Creating/Provisioning to Running state.
	KubernetesCreateTimeout = 30 * time.Minute

	// KubernetesUpdateTimeout is the maximum time to wait for a Kubernetes cluster update to complete.
	// This includes the time for the cluster to transition from Updating to Running state.
	KubernetesUpdateTimeout = 30 * time.Minute

	// KubernetesDeleteTimeout is the maximum time to wait for a Kubernetes cluster deletion to complete.
	// This includes polling until the cluster returns 404 (not found).
	KubernetesDeleteTimeout = 20 * time.Minute

	// KubernetesStatusCheckDelay is the interval between polling attempts when waiting for cluster state changes.
	// This controls how frequently the provider checks if a cluster has transitioned states.
	KubernetesStatusCheckDelay = 30 * time.Second

	// KubernetesStatusCheckMinTimeout is the minimum time between polling attempts.
	// This prevents excessive API calls during rapid state transitions.
	KubernetesStatusCheckMinTimeout = 10 * time.Second
)

// Kubernetes import format constants - Terraform-specific import configuration
const (
	// KubernetesImportDelimiter is the delimiter used in Kubernetes import IDs.
	// Format: "project_id:region:cluster_id" or "cluster_id"
	KubernetesImportDelimiter = ":"
)

// Kubernetes validation constants - Terraform-specific validation constraints
const (
	// KubernetesVersionRegex is the regex pattern for validating Kubernetes version format.
	// Format: 1.XX (e.g., 1.20, 1.21, 1.22)
	KubernetesVersionRegex = `^1\.\d{2}$`
)

// ============================================================================
// Volume Attachment Terraform Configuration
// ============================================================================
// These constants are Terraform-provider specific (timeouts, import formats, API response field names)
// and are not part of the goe2e API contract.

// Volume attachment operation timeout constants - Terraform-specific timeout configuration
// These constants control how long the provider waits for volume attachment operations to complete.
// All durations are typed as time.Duration to avoid unit mixing at callsites.
const (
	// VolumeAttachmentPollInterval is the interval between polling attempts when waiting for volume attachment state changes.
	// This controls how frequently the provider checks if a volume attachment has transitioned states.
	VolumeAttachmentPollInterval = 5 * time.Second

	// VolumeAttachmentWaitTimeout is the maximum time to wait for a volume attachment/detachment operation to complete.
	VolumeAttachmentWaitTimeout = 3 * time.Minute
)

// Volume attachment import format constants - Terraform-specific import configuration
const (
	// VolumeAttachmentImportDelimiter is the delimiter used in volume attachment import IDs.
	// Format: "node_id/volume_id" or "project_id/region/node_id/volume_id"
	VolumeAttachmentImportDelimiter = "/"
)

// Volume attachment API response field names - API response field keys
// These constants represent field names used in API responses (e.g., VMDetail map keys).
// These are API contract values but are Terraform-specific in that they're used for parsing responses.
const (
	// VolumeAttachmentVMDetailKeyVMID is the key for VM ID in volume VMDetail map.
	VolumeAttachmentVMDetailKeyVMID = "vm_id"

	// VolumeAttachmentVMDetailKeyVMName is the key for VM name in volume VMDetail map.
	VolumeAttachmentVMDetailKeyVMName = "vm_name"
)

// Volume attachment state values - Terraform-specific state checking values
// These constants represent values used for state checking in wait operations.
const (
	// VolumeAttachmentVMIDNullValue is the string value representing a null/empty VM ID in API responses.
	// This is used when checking if a volume is detached (vm_id is empty or "null").
	VolumeAttachmentVMIDNullValue = "null"
)
