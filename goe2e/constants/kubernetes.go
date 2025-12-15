package constants

// ============================================================================
// Kubernetes Constants
// ============================================================================

// Kubernetes cluster status constants - API response values
// These constants represent the various status values returned by the E2E API for Kubernetes clusters
const (
	KubernetesClusterStatusRunning      = "Running"      // Cluster is running normally
	KubernetesClusterStatusCreating     = "Creating"     // Cluster is being created
	KubernetesClusterStatusProvisioning = "Provisioning" // Cluster is being provisioned
	KubernetesClusterStatusUpdating     = "Updating"     // Cluster is being updated
	KubernetesClusterStatusDeleting     = "Deleting"     // Cluster is being deleted
	KubernetesClusterStatusDeleted      = "Deleted"      // Cluster has been deleted
	KubernetesClusterStatusFailed       = "Failed"       // Cluster is in failed state
	KubernetesClusterStatusError        = "Error"        // Cluster is in error state
)

// Kubernetes node pool type constants - API request/response values
// These constants represent the node pool type values used by the E2E API
const (
	KubernetesNodePoolTypeStatic    = "Static"    // Static node pool (fixed size)
	KubernetesNodePoolTypeAutoscale = "Autoscale" // Autoscale node pool (dynamic size)
)

// Kubernetes policy type constants - API request/response values
// These constants represent policy type values used in elasticity and scheduled policies
const (
	KubernetesPolicyTypeChange      = "CHANGE"      // Elasticity policy type
	KubernetesPolicyTypeCardinality = "CARDINALITY" // Scheduled policy type
	KubernetesPolicyTypeDefault     = "Default"     // Default policy parameter type
	KubernetesPolicyTypeCustom      = "Custom"      // Custom policy parameter type
	KubernetesPolicyTypeScheduled   = "Scheduled"   // Scheduled policy type
)

// Kubernetes policy parameter constants - API request/response values
// These constants represent parameter values used in elasticity policies
const (
	KubernetesPolicyParameterCPU    = "CPU"    // CPU parameter
	KubernetesPolicyParameterMemory = "Memory" // Memory parameter
)
