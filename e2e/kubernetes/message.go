package kubernetes

// ============================================================================
// Kubernetes-Specific Constants
// ============================================================================
// This file contains constants specific to the kubernetes resource that are
// not part of the API contract or Terraform provider-wide configuration.
// These constants are used internally within the kubernetes package.
//
// For API-related constants (status values, node pool types, etc.), see: goe2e/constants/kubernetes.go
// For Terraform-wide configuration (timeouts, poll intervals, etc.), see: e2e/constants/internal.go
// For Terraform attribute names, see: e2e/constants/attrs.go

// ============================================================================
// Production Error Message Constants
// ============================================================================
// These constants define error messages used in production code for consistent
// error handling and reporting across the kubernetes resource.

const (
	// Node pool validation errors
	ErrNodePoolDuplicateNames = "Name of the worker node pools must be unique!"

	// Node pool type errors
	ErrNodePoolTypeRequired = "node pool type (type or node_pool_type) is required"
	ErrNodePoolTypeInvalid  = "invalid node_pool_type: %s"

	// Node pool configuration errors
	ErrNodePoolPlanNotFound         = "no matching plan found for plan: %s"
	ErrNodePoolStaticSizeRequired   = "size (or worker_node) is required for Static node pools"
	ErrNodePoolAutoscaleMinRequired = "in case of Autoscale node type, the 'min_nodes' (or 'min_vms') field is required"
	ErrNodePoolAutoscaleMaxRequired = "in case of Autoscale node type, the 'max_nodes' (or 'max_vms') field is required"

	// Cluster operation errors
	ErrClusterVersionRequired = "kubernetes_version or version is required"
	ErrClusterNilResponse     = "Cluster creation returned nil response"

	// Import errors
	ErrImportInvalidFormat   = "invalid import ID format: expected 'cluster_id' or 'project_id:region:cluster_id', got '%s'"
	ErrImportClusterIDEmpty  = "cluster_id cannot be empty"
	ErrImportClusterNotFound = "Kubernetes cluster (ID: %s) not found"

	// State setting errors
	ErrSettingSlugName  = "error setting slug_name: %w"
	ErrSettingTags      = "error setting tags: %w"
	ErrSettingNodePools = "error setting node_pools: %s"

	// Node pool update errors
	ErrNodePoolSizeTooSmall     = "Cannot update node pool '%s' in Kubernetes cluster (ID: %s): node_pool_size must be at least 2 (current value: %d)"
	ErrNodePoolTypeImmutable    = "Cannot update node pool type for node pool '%s' in Kubernetes cluster (ID: %s): this field is immutable after node pool creation"
	ErrNodePoolDeleteNotRunning = "Cannot delete node pool '%s' from Kubernetes cluster (ID: %s): node pool must be in Running state before deletion"
	ErrNodePoolDeleteLastPool   = "Cannot delete node pool from Kubernetes cluster (ID: %s): at least one node pool must be present in a Kubernetes cluster"
	ErrNodePoolNotFound         = "Cannot delete node pool '%s' from Kubernetes cluster (ID: %s): node pool does not exist in project (%s), region (%s)"

	// Elasticity and scheduled policy errors
	ErrElasticityInvalidFormat   = "Invalid format for Elast"
	ErrScheduledInvalidFormat    = "Invalid format for Scheduled Dictionary"
	ErrUpscaleCardinalityRange   = "upscale cardinality must be between min nodes and max nodes"
	ErrDownscaleCardinalityRange = "downscale cardinality must be between min nodes and max nodes"
)
