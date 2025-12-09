package goe2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	kubernetesPath                    = "kubernetes"
	kubernetesPlansPath               = "kubernetes/plans"
	kubernetesWorkerPlansPath         = "kubernetes/worker-plans/"
	kubernetesNodePoolsPath           = "kubernetes/node-pool-services"       // + /{clusterID}
	kubernetesClusterUpdatePath       = "kubernetes/cluster-update"           // + /{serviceID}
	kubernetesDeleteNodePoolPath      = "kubernetes/delete-node-pool-service" // + /{serviceID}
	kubernetesAddNodePoolsPath        = "kubernetes/add-node-pools"           // + /{clusterID}
	kubernetesUpdateNodePoolPath      = "kubernetes/update-node-pool"         // + /{serviceID}/
	kubernetesPersistentVolumePath    = "persistent_volume"                   // + /{clusterID}/ or /{clusterID}/{pvID}/
	kubernetesAttachSecurityGroupPath = "kubernetes/attach-security-group"    // + /{clusterID}/
	kubernetesDetachSecurityGroupPath = "kubernetes/detach-security-group"    // + /{clusterID}/
	securityGroupAttachPath           = "security_group"                      // + /{clusterID}/attach/
)

// KubernetesService is an interface for interacting with the Kubernetes endpoints
// of the E2E Networks API.
type KubernetesService interface {
	// Plan operations
	GetMasterPlans(context.Context) ([]KubernetesPlan, *Response, error)
	GetWorkerPlans(context.Context) ([]KubernetesWorkerPlan, *Response, error)

	// Cluster operations
	Create(context.Context, *KubernetesClusterCreateRequest) (*KubernetesCluster, *Response, error)
	Get(context.Context, string) (*KubernetesCluster, *Response, error)
	Delete(context.Context, string) (*Response, error)

	// Node pool operations
	GetNodePools(context.Context, string) ([]NodePoolServiceInfo, *Response, error)
	AddNodePool(context.Context, string, *NodePoolAddRequest) (*NodePoolAddResponse, *Response, error)
	UpdateNodePoolCardinality(context.Context, string, *NodePoolResizeRequest) (*Response, error)
	UpdateNodePoolDetails(context.Context, string, *NodePoolUpdateRequest) (*NodePoolUpdateResponse, *Response, error)
	DeleteNodePool(context.Context, string) (*Response, error)
	CheckNodePoolStatus(context.Context, string) ([]NodePoolServiceInfo, *Response, error)

	// Persistent volume operations
	ListPersistentVolumes(context.Context, string) ([]PersistentVolume, *Response, error)
	CreatePersistentVolume(context.Context, string, *CreatePersistentVolumeRequest) (*PersistentVolume, *Response, error)
	GetPersistentVolume(context.Context, string, string) (*PersistentVolume, *Response, error)
	DeletePersistentVolume(context.Context, string, string) (*Response, error)

	// Security group operations
	ListAttachedSecurityGroups(context.Context, string) ([]SecurityGroupAttachment, *Response, error)
	AttachSecurityGroups(context.Context, string, *AttachSecurityGroupRequest) (*Response, error)
	DetachSecurityGroups(context.Context, string, *DetachSecurityGroupRequest) (*Response, error)
}

// KubernetesServiceOp handles communication with Kubernetes related methods of the
// E2E Networks API.
type KubernetesServiceOp struct {
	client *Client
}

var _ KubernetesService = &KubernetesServiceOp{}

// KubernetesPlan represents a Kubernetes master plan
type KubernetesPlan struct {
	Plan       string                 `json:"plan"`
	K8sVersion string                 `json:"k8s_version"`
	Specs      KubernetesPlanSpecs    `json:"specs"`
	RawData    map[string]interface{} `json:"-"` // For any additional fields
}

// KubernetesPlanSpecs represents the specs of a Kubernetes plan
type KubernetesPlanSpecs struct {
	ID      string `json:"id"`
	SKUName string `json:"sku_name"`
}

// KubernetesWorkerPlan represents a Kubernetes worker plan
type KubernetesWorkerPlan struct {
	Plan    string                    `json:"plan"`
	Specs   KubernetesWorkerPlanSpecs `json:"specs"`
	RawData map[string]interface{}    `json:"-"` // For any additional fields
}

// KubernetesWorkerPlanSpecs represents the specs of a worker plan
type KubernetesWorkerPlanSpecs struct {
	ID      string `json:"id"`
	SKUName string `json:"sku_name"`
}

// KubernetesCluster represents a Kubernetes cluster
type KubernetesCluster struct {
	ID          string `json:"id"`
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	State       string `json:"state"`
	Version     string `json:"version"`
	VPCID       string `json:"vpc_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// KubernetesClusterCreateRequest represents a request to create a Kubernetes cluster
type KubernetesClusterCreateRequest struct {
	Name      string     `json:"name"`
	SlugName  string     `json:"slug_name"`
	Version   string     `json:"version"`
	VPCID     string     `json:"vpc_id"`
	SKUID     string     `json:"sku_id"`
	NodePools []NodePool `json:"node_pools"`
}

// NodePoolServiceInfo represents information about a node pool service
type NodePoolServiceInfo struct {
	ServiceID   float64 `json:"service_id"`
	ServiceName string  `json:"service_name"`
	State       string  `json:"state"`
	Cardinality int     `json:"cardinality"`
}

// NodePoolAddRequest represents a request to add node pools
type NodePoolAddRequest struct {
	NodePools []NodePool `json:"node_pools"`
}

// NodePoolAddResponse represents the response from adding node pools
type NodePoolAddResponse struct {
	Message string `json:"message"`
}

// NodePoolResizeRequest represents a request to resize a node pool
type NodePoolResizeRequest struct {
	NodePoolSize int `json:"cardinality"`
}

// NodePoolUpdateRequest represents a request to update node pool details
type NodePoolUpdateRequest struct {
	MinVms           int                `json:"min_vms"`
	Cardinality      int                `json:"cardinality"`
	MaxVms           int                `json:"max_vms"`
	PlanID           string             `json:"plan_id"`
	ElasticityPolicy []ElasticityPolicy `json:"elasticity_policies"`
	ScheduledPolicy  []SchedulePolicy   `json:"scheduled_policies"`
	PolicyType       string             `json:"policy_type,omitempty"`
	CustomParamName  string             `json:"custom_param_name,omitempty"`
	CustomParamValue string             `json:"custom_param_value,omitempty"`
}

// NodePoolUpdateResponse represents the response from updating node pool details
type NodePoolUpdateResponse struct {
	Message string `json:"message"`
}

// ElasticityPolicy represents an elasticity policy for auto-scaling
type ElasticityPolicy struct {
	Type         string `json:"type"`
	Adjust       int    `json:"adjust"`
	Parameter    string `json:"parameter"`
	Operator     string `json:"operator"`
	Value        int    `json:"value"`
	PeriodNumber int    `json:"period_number"`
	Period       int    `json:"period"`
	Cooldown     int    `json:"cooldown"`
}

// SchedulePolicy represents a scheduled scaling policy
type SchedulePolicy struct {
	Type       string `json:"type"`
	Adjust     int    `json:"adjust"`
	Recurrence string `json:"recurrence"`
}

// ElasticityDict wraps elasticity worker configuration
type ElasticityDict struct {
	Worker ElasticityWorker `json:"worker"`
}

// ElasticityWorker represents the worker configuration for elasticity
type ElasticityWorker struct {
	MinVms             int                `json:"min_vms"`
	Cardinality        int                `json:"cardinality"`
	MaxVms             int                `json:"max_vms"`
	ElasticityPolicies []ElasticityPolicy `json:"elasticity_policies"`
}

// ScheduledDict wraps scheduled worker configuration
type ScheduledDict struct {
	Worker ScheduleWorker `json:"worker"`
}

// ScheduleWorker represents the worker configuration for scheduled scaling
type ScheduleWorker struct {
	MinVms            int              `json:"min_vms"`
	Cardinality       int              `json:"cardinality"`
	MaxVms            int              `json:"max_vms"`
	ScheduledPolicies []SchedulePolicy `json:"scheduled_policies"`
}

// NodePool represents a Kubernetes node pool configuration
type NodePool struct {
	Name             string         `json:"name"`
	SlugName         string         `json:"slug_name"`
	SKUID            string         `json:"sku_id"`
	SpecsName        string         `json:"specs_name"`
	WorkerNode       int            `json:"worker_node,omitempty"`
	ElasticityDict   ElasticityDict `json:"elasticity_dict,omitempty"`
	ScheduledDict    ScheduledDict  `json:"scheduled_dict,omitempty"`
	PolicyType       string         `json:"policy_type,omitempty"` //I changed it
	CustomParamName  string         `json:"custom_param_name,omitempty"`
	CustomParamValue string         `json:"custom_param_value,omitempty"`
}

// NodePoolUpdate represents a node pool update request
type NodePoolUpdate struct {
	MinVms           int                `json:"min_vms"`
	Cardinality      int                `json:"cardinality"`
	MaxVms           int                `json:"max_vms"`
	PlanID           string             `json:"plan_id"`
	ElasticityPolicy []ElasticityPolicy `json:"elasticity_policies"`
	ScheduledPolicy  []SchedulePolicy   `json:"scheduled_policies"`
	PolicyType       string             `json:"policy_type,omitempty"`
	CustomParamName  string             `json:"custom_param_name,omitempty"`
	CustomParamValue string             `json:"custom_param_value,omitempty"`
}

// NodePoolAdd represents a request to add node pools
type NodePoolAdd struct {
	NodePools []NodePool `json:"node_pools"`
}

// NodePoolResize represents a request to resize a node pool
type NodePoolResize struct {
	NodePoolSize int `json:"cardinality"`
}

// PersistentVolume represents a Kubernetes persistent volume
type PersistentVolume struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PVSize    int    `json:"pv_size"`
	ClusterID string `json:"cluster_id"`
	Status    string `json:"status"`
	IsDynamic bool   `json:"is_dynamic"`
	CreatedAt string `json:"created_at"`
}

// CreatePersistentVolumeRequest represents a request to create a persistent volume
type CreatePersistentVolumeRequest struct {
	Name      string `json:"name"`
	PVSize    int    `json:"pv_size"`
	IsDynamic bool   `json:"is_dynamic,omitempty"`
}

// AttachSecurityGroupRequest represents a request to attach security groups to a Kubernetes cluster
type AttachSecurityGroupRequest struct {
	SecurityGroupIDs []int `json:"security_group_ids"`
}

// DetachSecurityGroupRequest represents a request to detach security groups from a Kubernetes cluster
type DetachSecurityGroupRequest struct {
	SecurityGroupIDs []int `json:"security_group_ids"`
}

// SecurityGroupAttachment represents a security group attached to a Kubernetes cluster
type SecurityGroupAttachment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Response wrappers for API calls
type kubernetesPlanRoot struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type kubernetesWorkerPlanRoot struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type kubernetesClusterRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type kubernetesClusterGetRoot struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type nodePoolsRoot struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type nodePoolAddRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type nodePoolUpdateRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type persistentVolumesRoot struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type persistentVolumeRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type securityGroupAttachmentsRoot struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

// GetMasterPlans retrieves available Kubernetes master plans.
func (s *KubernetesServiceOp) GetMasterPlans(ctx context.Context) ([]KubernetesPlan, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, kubernetesPlansPath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for Kubernetes master plans: %w", err)
	}

	root := new(kubernetesPlanRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve Kubernetes master plans: %w", err)
	}

	plans := make([]KubernetesPlan, 0, len(root.Data))
	for _, planData := range root.Data {
		plan := KubernetesPlan{
			RawData: planData,
		}

		if v, ok := planData["plan"].(string); ok {
			plan.Plan = v
		}
		if v, ok := planData["k8s_version"].(string); ok {
			plan.K8sVersion = v
		}
		if specs, ok := planData["specs"].(map[string]interface{}); ok {
			if id, ok := specs["id"].(string); ok {
				plan.Specs.ID = id
			}
			if skuName, ok := specs["sku_name"].(string); ok {
				plan.Specs.SKUName = skuName
			}
		}

		plans = append(plans, plan)
	}

	return plans, resp, nil
}

// GetWorkerPlans retrieves available Kubernetes worker plans.
func (s *KubernetesServiceOp) GetWorkerPlans(ctx context.Context) ([]KubernetesWorkerPlan, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, kubernetesWorkerPlansPath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for Kubernetes worker plans: %w", err)
	}

	root := new(kubernetesWorkerPlanRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve Kubernetes worker plans: %w", err)
	}

	plans := make([]KubernetesWorkerPlan, 0, len(root.Data))
	for _, planData := range root.Data {
		plan := KubernetesWorkerPlan{
			RawData: planData,
		}

		if v, ok := planData["plan"].(string); ok {
			plan.Plan = v
		}
		if specs, ok := planData["specs"].(map[string]interface{}); ok {
			if id, ok := specs["id"].(string); ok {
				plan.Specs.ID = id
			}
			if skuName, ok := specs["sku_name"].(string); ok {
				plan.Specs.SKUName = skuName
			}
		}

		plans = append(plans, plan)
	}

	return plans, resp, nil
}

// Create creates a new Kubernetes cluster.
func (s *KubernetesServiceOp) Create(ctx context.Context, createReq *KubernetesClusterCreateRequest) (*KubernetesCluster, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("name", "cannot be empty")
	}
	if createReq.VPCID == "" {
		return nil, nil, NewArgError("vpc_id", "cannot be empty")
	}

	// Apply the RemoveExtraFieldsFromKubernetes logic
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode create request: %w", err)
	}

	cleanedBuf, err := removeExtraFieldsFromKubernetes(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to clean request fields: %w", err)
	}

	// Create raw request from cleaned buffer
	req, err := s.client.NewRequest(ctx, http.MethodPost, kubernetesPath+"/", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for Kubernetes cluster (%s): %w", createReq.Name, err)
	}

	// Replace body with cleaned buffer
	req.Body = io.NopCloser(&cleanedBuf)
	req.ContentLength = int64(cleanedBuf.Len())

	root := new(kubernetesClusterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create Kubernetes cluster (%s): %w", createReq.Name, err)
	}

	// Parse the response
	cluster := &KubernetesCluster{}
	if data, ok := root.Data["DOCUMENT"].(map[string]interface{}); ok {
		if id, ok := data["ID"].(string); ok {
			cluster.ID = id
			cluster.ServiceID = id
		}
	}

	return cluster, resp, nil
}

// Get retrieves a Kubernetes cluster by ID.
func (s *KubernetesServiceOp) Get(ctx context.Context, clusterID string) (*KubernetesCluster, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", kubernetesPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for Kubernetes cluster (ID: %s): %w", clusterID, err)
	}

	root := new(kubernetesClusterGetRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve Kubernetes cluster (ID: %s): %w", clusterID, err)
	}

	if len(root.Data) == 0 {
		return nil, resp, fmt.Errorf("Kubernetes cluster not found (ID: %s)", clusterID)
	}

	// Parse the first cluster from the response
	clusterData := root.Data[0]
	cluster := &KubernetesCluster{}

	if serviceID, ok := clusterData["service_id"].(float64); ok {
		cluster.ServiceID = fmt.Sprintf("%.0f", serviceID)
		cluster.ID = cluster.ServiceID
	}
	if serviceName, ok := clusterData["service_name"].(string); ok {
		cluster.ServiceName = serviceName
	}
	if state, ok := clusterData["state"].(string); ok {
		cluster.State = state
	}
	if version, ok := clusterData["version"].(string); ok {
		cluster.Version = version
	}
	if createdAt, ok := clusterData["created_at"].(string); ok {
		cluster.CreatedAt = createdAt
	}

	return cluster, resp, nil
}

// Delete deletes a Kubernetes cluster.
func (s *KubernetesServiceOp) Delete(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", kubernetesPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting Kubernetes cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete Kubernetes cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// GetNodePools retrieves the node pools for a Kubernetes cluster.
func (s *KubernetesServiceOp) GetNodePools(ctx context.Context, clusterID string) ([]NodePoolServiceInfo, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", kubernetesNodePoolsPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for node pools (cluster ID: %s): %w", clusterID, err)
	}

	root := new(nodePoolsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve node pools for cluster (ID: %s): %w", clusterID, err)
	}

	nodePools := make([]NodePoolServiceInfo, 0, len(root.Data))
	for _, npData := range root.Data {
		np := NodePoolServiceInfo{}

		if serviceID, ok := npData["service_id"].(float64); ok {
			np.ServiceID = serviceID
		}
		if serviceName, ok := npData["service_name"].(string); ok {
			np.ServiceName = serviceName
		}
		if state, ok := npData["state"].(string); ok {
			np.State = state
		}
		if cardinality, ok := npData["cardinality"].(float64); ok {
			np.Cardinality = int(cardinality)
		} else if cardinality, ok := npData["cardinality"].(int); ok {
			np.Cardinality = cardinality
		}

		nodePools = append(nodePools, np)
	}

	return nodePools, resp, nil
}

// CheckNodePoolStatus checks the status of node pools for a Kubernetes cluster.
// This is an alias for GetNodePools for backward compatibility.
func (s *KubernetesServiceOp) CheckNodePoolStatus(ctx context.Context, clusterID string) ([]NodePoolServiceInfo, *Response, error) {
	return s.GetNodePools(ctx, clusterID)
}

// AddNodePool adds new node pools to a Kubernetes cluster.
func (s *KubernetesServiceOp) AddNodePool(ctx context.Context, clusterID string, addReq *NodePoolAddRequest) (*NodePoolAddResponse, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}
	if addReq == nil {
		return nil, nil, NewArgError("addReq", "cannot be nil")
	}

	// Apply the RemoveExtraFieldsFromKubernetes logic
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(addReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode add node pool request: %w", err)
	}

	cleanedBuf, err := removeExtraFieldsFromKubernetes(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to clean request fields: %w", err)
	}

	path := fmt.Sprintf("%s/%s", kubernetesAddNodePoolsPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for adding node pool (cluster ID: %s): %w", clusterID, err)
	}

	// Replace body with cleaned buffer
	req.Body = io.NopCloser(&cleanedBuf)
	req.ContentLength = int64(cleanedBuf.Len())

	root := new(nodePoolAddRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to add node pool to cluster (ID: %s): %w", clusterID, err)
	}

	result := &NodePoolAddResponse{
		Message: root.Message,
	}

	return result, resp, nil
}

// UpdateNodePoolCardinality updates the cardinality (size) of a node pool.
func (s *KubernetesServiceOp) UpdateNodePoolCardinality(ctx context.Context, nodePoolServiceID string, resizeReq *NodePoolResizeRequest) (*Response, error) {
	if nodePoolServiceID == "" {
		return nil, NewArgError("nodePoolServiceID", "cannot be empty")
	}
	if resizeReq == nil {
		return nil, NewArgError("resizeReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s", kubernetesClusterUpdatePath, nodePoolServiceID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, resizeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for updating node pool cardinality (service ID: %s): %w", nodePoolServiceID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to update node pool cardinality (service ID: %s): %w", nodePoolServiceID, err)
	}

	return resp, nil
}

// UpdateNodePoolDetails updates the details of a node pool.
func (s *KubernetesServiceOp) UpdateNodePoolDetails(ctx context.Context, nodePoolServiceID string, updateReq *NodePoolUpdateRequest) (*NodePoolUpdateResponse, *Response, error) {
	if nodePoolServiceID == "" {
		return nil, nil, NewArgError("nodePoolServiceID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", kubernetesUpdateNodePoolPath, nodePoolServiceID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for updating node pool details (service ID: %s): %w", nodePoolServiceID, err)
	}

	root := new(nodePoolUpdateRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update node pool details (service ID: %s): %w", nodePoolServiceID, err)
	}

	result := &NodePoolUpdateResponse{
		Message: root.Message,
	}

	return result, resp, nil
}

// DeleteNodePool deletes a node pool from a Kubernetes cluster.
func (s *KubernetesServiceOp) DeleteNodePool(ctx context.Context, nodePoolServiceID string) (*Response, error) {
	if nodePoolServiceID == "" {
		return nil, NewArgError("nodePoolServiceID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", kubernetesDeleteNodePoolPath, nodePoolServiceID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting node pool (service ID: %s): %w", nodePoolServiceID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete node pool (service ID: %s): %w", nodePoolServiceID, err)
	}

	return resp, nil
}

// removeExtraFieldsFromKubernetes removes extra fields from the node pools in a Kubernetes request.
// This replicates the logic from client/kubernetes_client.go RemoveExtraFieldsFromKubernetes.
func removeExtraFieldsFromKubernetes(buf *bytes.Buffer) (bytes.Buffer, error) {
	jsonData := buf.Bytes()

	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return *buf, err
	}

	nodePools, ok := data["node_pools"].([]interface{})
	if !ok {
		// If "node_pools" is not present or not a slice, return an error
		return *buf, errors.New("node_pools field is missing or invalid")
	}

	for _, nodePool := range nodePools {
		if nodePoolMap, ok := nodePool.(map[string]interface{}); ok {
			// Type assert to float64
			workerNode, workerNodePresent := nodePoolMap["worker_node"].(float64)
			if workerNodePresent {
				if workerNode == 0 {
					// If worker_node is present and its value is 0, delete the "worker_node" field
					delete(nodePoolMap, "worker_node")
				} else if workerNode >= 2 {
					// If worker_node is greater than or equal to 2, check if "elasticity_dict" is present
					// Delete optional fields if present (delete is safe on non-existent keys)
					delete(nodePoolMap, "elasticity_dict")
					delete(nodePoolMap, "scheduled_dict")
				}
			}
			policyType, policyTypePresent := nodePoolMap["policy_type"].(string)
			if !policyTypePresent || policyType == "Scheduled" {
				// If policy_type is "Scheduled", remove elasticity_dict
				delete(nodePoolMap, "elasticity_dict")
			}
			if !policyTypePresent || !contains(policyType, "Scheduled") {
				// If policy_type does not contain the keyword "Scheduled", remove scheduled_dict
				delete(nodePoolMap, "scheduled_dict")
			}
		}
	}

	newJSONData, err := json.Marshal(data)
	if err != nil {
		return *buf, err
	}

	return *bytes.NewBuffer(newJSONData), nil
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ListPersistentVolumes retrieves all persistent volumes for a Kubernetes cluster.
func (s *KubernetesServiceOp) ListPersistentVolumes(ctx context.Context, clusterID string) ([]PersistentVolume, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", kubernetesPersistentVolumePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for persistent volumes (cluster ID: %s): %w", clusterID, err)
	}

	root := new(persistentVolumesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve persistent volumes for cluster (ID: %s): %w", clusterID, err)
	}

	pvs := make([]PersistentVolume, 0, len(root.Data))
	for _, pvData := range root.Data {
		pv := PersistentVolume{}

		if id, ok := pvData["id"].(string); ok {
			pv.ID = id
		}
		if name, ok := pvData["name"].(string); ok {
			pv.Name = name
		}
		if pvSize, ok := pvData["pv_size"].(float64); ok {
			pv.PVSize = int(pvSize)
		} else if pvSize, ok := pvData["pv_size"].(int); ok {
			pv.PVSize = pvSize
		}
		if clusterID, ok := pvData["cluster_id"].(string); ok {
			pv.ClusterID = clusterID
		}
		if status, ok := pvData["status"].(string); ok {
			pv.Status = status
		}
		if isDynamic, ok := pvData["is_dynamic"].(bool); ok {
			pv.IsDynamic = isDynamic
		}
		if createdAt, ok := pvData["created_at"].(string); ok {
			pv.CreatedAt = createdAt
		}

		pvs = append(pvs, pv)
	}

	return pvs, resp, nil
}

// CreatePersistentVolume creates a new persistent volume for a Kubernetes cluster.
func (s *KubernetesServiceOp) CreatePersistentVolume(ctx context.Context, clusterID string, createReq *CreatePersistentVolumeRequest) (*PersistentVolume, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("name", "cannot be empty")
	}
	if createReq.PVSize <= 0 {
		return nil, nil, NewArgError("pv_size", "must be greater than 0")
	}

	path := fmt.Sprintf("%s/%s/", kubernetesPersistentVolumePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for persistent volume (%s) in cluster (%s): %w", createReq.Name, clusterID, err)
	}

	root := new(persistentVolumeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create persistent volume (%s) in cluster (%s): %w", createReq.Name, clusterID, err)
	}

	// Parse the response
	pv := &PersistentVolume{}
	if id, ok := root.Data["id"].(string); ok {
		pv.ID = id
	}
	if name, ok := root.Data["name"].(string); ok {
		pv.Name = name
	}
	if pvSize, ok := root.Data["pv_size"].(float64); ok {
		pv.PVSize = int(pvSize)
	} else if pvSize, ok := root.Data["pv_size"].(int); ok {
		pv.PVSize = pvSize
	}
	if clusterIDStr, ok := root.Data["cluster_id"].(string); ok {
		pv.ClusterID = clusterIDStr
	} else {
		pv.ClusterID = clusterID
	}
	if status, ok := root.Data["status"].(string); ok {
		pv.Status = status
	}
	if isDynamic, ok := root.Data["is_dynamic"].(bool); ok {
		pv.IsDynamic = isDynamic
	}
	if createdAt, ok := root.Data["created_at"].(string); ok {
		pv.CreatedAt = createdAt
	}

	return pv, resp, nil
}

// GetPersistentVolume retrieves a persistent volume by ID for a Kubernetes cluster.
func (s *KubernetesServiceOp) GetPersistentVolume(ctx context.Context, clusterID, pvID string) (*PersistentVolume, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}
	if pvID == "" {
		return nil, nil, NewArgError("pvID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/%s/", kubernetesPersistentVolumePath, clusterID, pvID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for persistent volume (ID: %s) in cluster (ID: %s): %w", pvID, clusterID, err)
	}

	root := new(persistentVolumeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil PV for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve persistent volume (ID: %s) in cluster (ID: %s): %w", pvID, clusterID, err)
	}

	// Parse the response
	pv := &PersistentVolume{}
	if id, ok := root.Data["id"].(string); ok {
		pv.ID = id
	}
	if name, ok := root.Data["name"].(string); ok {
		pv.Name = name
	}
	if pvSize, ok := root.Data["pv_size"].(float64); ok {
		pv.PVSize = int(pvSize)
	} else if pvSize, ok := root.Data["pv_size"].(int); ok {
		pv.PVSize = pvSize
	}
	if clusterIDStr, ok := root.Data["cluster_id"].(string); ok {
		pv.ClusterID = clusterIDStr
	} else {
		pv.ClusterID = clusterID
	}
	if status, ok := root.Data["status"].(string); ok {
		pv.Status = status
	}
	if isDynamic, ok := root.Data["is_dynamic"].(bool); ok {
		pv.IsDynamic = isDynamic
	}
	if createdAt, ok := root.Data["created_at"].(string); ok {
		pv.CreatedAt = createdAt
	}

	return pv, resp, nil
}

// DeletePersistentVolume deletes a persistent volume from a Kubernetes cluster.
func (s *KubernetesServiceOp) DeletePersistentVolume(ctx context.Context, clusterID, pvID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if pvID == "" {
		return nil, NewArgError("pvID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/%s/", kubernetesPersistentVolumePath, clusterID, pvID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting persistent volume (ID: %s) in cluster (ID: %s): %w", pvID, clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete persistent volume (ID: %s) in cluster (ID: %s): %w", pvID, clusterID, err)
	}

	return resp, nil
}

// ListAttachedSecurityGroups retrieves all security groups attached to a Kubernetes cluster.
func (s *KubernetesServiceOp) ListAttachedSecurityGroups(ctx context.Context, clusterID string) ([]SecurityGroupAttachment, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/attach/", securityGroupAttachPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing attached security groups (cluster ID: %s): %w", clusterID, err)
	}

	root := new(securityGroupAttachmentsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve attached security groups for cluster (ID: %s): %w", clusterID, err)
	}

	attachments := make([]SecurityGroupAttachment, 0, len(root.Data))
	for _, sgData := range root.Data {
		attachment := SecurityGroupAttachment{}

		if id, ok := sgData["id"].(float64); ok {
			attachment.ID = int(id)
		} else if id, ok := sgData["id"].(int); ok {
			attachment.ID = id
		}
		if name, ok := sgData["name"].(string); ok {
			attachment.Name = name
		}
		if description, ok := sgData["description"].(string); ok {
			attachment.Description = description
		}

		attachments = append(attachments, attachment)
	}

	return attachments, resp, nil
}

// AttachSecurityGroups attaches security groups to a Kubernetes cluster.
func (s *KubernetesServiceOp) AttachSecurityGroups(ctx context.Context, clusterID string, attachReq *AttachSecurityGroupRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if attachReq == nil {
		return nil, NewArgError("attachReq", "cannot be nil")
	}
	if len(attachReq.SecurityGroupIDs) == 0 {
		return nil, NewArgError("security_group_ids", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", kubernetesAttachSecurityGroupPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, attachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching security groups to cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach security groups to cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// DetachSecurityGroups detaches security groups from a Kubernetes cluster.
func (s *KubernetesServiceOp) DetachSecurityGroups(ctx context.Context, clusterID string, detachReq *DetachSecurityGroupRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if detachReq == nil {
		return nil, NewArgError("detachReq", "cannot be nil")
	}
	if len(detachReq.SecurityGroupIDs) == 0 {
		return nil, NewArgError("security_group_ids", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", kubernetesDetachSecurityGroupPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, detachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching security groups from cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach security groups from cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}
