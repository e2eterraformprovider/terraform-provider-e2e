package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	scalerGroupsPath              = "scaler/scalegroups"
	scalerGroupPath               = "scaler/scalegroups"                 // + /{id}/
	scalerGroupUpdatePath         = "scaler/scalegroups/update"          // + /{id}/
	scalerGroupStartPath          = "scaler/scalegroups"                 // + /{id}/start/
	scalerGroupStopPath           = "scaler/scalegroups"                 // + /{id}/stop/
	scalerGroupVPCActionPath      = "scaler/scalegroups"                 // + /{id}/vpc/action/
	scalerGroupSecurityGroupsPath = "scaler/scalegroups/security_groups" // + /{id}/
	scalerGroupPublicIPPath       = "scaler/scalegroups"                 // + /{id}/public_ip/action/
)

// AutoscalingService is an interface for interacting with the Scaler Group endpoints
// of the E2E Networks API.
type AutoscalingService interface {
	// Scaler Group CRUD operations
	CreateScalerGroup(context.Context, *ScalerGroupCreateRequest) (*ScalerGroup, *Response, error)
	GetScalerGroup(context.Context, string) (*ScalerGroup, *Response, error)
	UpdateScalerGroup(context.Context, string, *ScalerGroupUpdateRequest) (*Response, error)
	DeleteScalerGroup(context.Context, string) (*Response, error)
	ListScalerGroups(context.Context) ([]ScalerGroup, *Response, error)

	// Scaler Group status operations
	StartScalerGroup(context.Context, string) (*Response, error)
	StopScalerGroup(context.Context, string) (*Response, error)
	UpdateScalerGroupStatus(context.Context, string, string) (*Response, error)

	// Desired node count operations
	UpdateDesiredNodeCount(context.Context, string, int) (*Response, error)

	// VPC operations
	AttachVPCToScalerGroup(context.Context, string, *VPCAttachRequest) (*Response, error)
	DetachVPCFromScalerGroup(context.Context, string, string) (*Response, error)
	GetAttachedVPCsForScalerGroup(context.Context, string) ([]VPCPartial, *Response, error)

	// Security Group operations
	AttachSecurityGroupToScalerGroup(context.Context, string, int) (*Response, error)
	DetachSecurityGroupFromScalerGroup(context.Context, string, int) (*Response, error)

	// Public IP operations
	AttachPublicIPToScalerGroup(context.Context, string) (*PublicIPActionResponse, *Response, error)
	DetachPublicIPFromScalerGroup(context.Context, string) (*PublicIPActionResponse, *Response, error)
	GetPublicIPStatus(context.Context, string) (*PublicIPStatus, *Response, error)
}

// AutoscalingServiceOp handles communication with Scaler Group related methods of the
// E2E Networks API.
type AutoscalingServiceOp struct {
	client *Client
}

var _ AutoscalingService = &AutoscalingServiceOp{}

// ScalerGroup represents a scaler group (autoscaling group) in E2E Cloud
type ScalerGroup struct {
	ID                      string            `json:"id"`
	Name                    string            `json:"name"`
	VMImageName             string            `json:"vm_image_name"`
	ProvisionStatus         string            `json:"provision_status"`
	Running                 int               `json:"running"`
	Desired                 int               `json:"desired"`
	Tags                    string            `json:"tags"`
	MinNodes                int               `json:"min_nodes"`
	MaxNodes                int               `json:"max_nodes"`
	CustomerID              int               `json:"customer_id"`
	PlanID                  int               `json:"plan_id"`
	PlanName                string            `json:"plan_name"`
	ImageID                 int               `json:"image_id"`
	VPCNames                string            `json:"vpc_names"`
	Nodes                   []ScalerGroupNode `json:"nodes"`
	PolicyType              string            `json:"policy_type"`
	Policy                  string            `json:"policy"`
	UpscalePolicyValue      int               `json:"upscale_policy_value"`
	DownscalePolicyValue    int               `json:"downscale_policy_value"`
	WaitForPeriod           int               `json:"wait_for_period"`
	WaitPeriod              int               `json:"wait_period"`
	Cooldown                int               `json:"cooldown"`
	PolicyMeasure           string            `json:"policy_measure"`
	PolicyUpscaleOperator   string            `json:"policy_upscale_operator"`
	PolicyDownscaleOperator string            `json:"policy_downscale_operator"`
	ParameterEvaluatedValue float64           `json:"parameter_evaluated_value"`
	ScheduledPolicyOp       string            `json:"scheduled_policy_op"`
	UpscaleRecurrence       string            `json:"upscale_recurrence"`
	UpscaleAdjust           int               `json:"upscale_adjust"`
	DownscaleRecurrence     string            `json:"downscale_recurrence"`
	DownscaleAdjust         int               `json:"downscale_adjust"`
}

// ScalerGroupNode represents a node in a scaler group
type ScalerGroupNode struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	IP       []string `json:"ip"`
	PublicIP string   `json:"public_ip"`
	Status   string   `json:"status"`
	RealCPU  string   `json:"real_cpu"`
}

// ScalerGroupCreateRequest represents a request to create a scaler group
type ScalerGroupCreateRequest struct {
	Name                 string            `json:"name"`
	PlanID               string            `json:"plan_id"`
	PlanName             string            `json:"plan_name"`
	SlugName             string            `json:"slug_name"`
	SKUID                string            `json:"sku_id"`
	VMImageID            string            `json:"vm_image_id"`
	VMImageName          string            `json:"vm_image_name"`
	VMTemplateID         int               `json:"vm_template_id"`
	MyAccountSGID        int               `json:"my_account_sg_id"`
	IsEncryptionEnabled  bool              `json:"isEncryptionEnabled"`
	EncryptionPassphrase string            `json:"encryption_passphrase,omitempty"`
	IsPublicIPRequired   bool              `json:"is_public_ip_required"`
	MinNodes             string            `json:"min_nodes"`
	MaxNodes             string            `json:"max_nodes"`
	Desired              string            `json:"desired"`
	PolicyType           string            `json:"policy_type,omitempty"`
	Policy               []ElasticPolicy   `json:"policy,omitempty"`
	ScheduledPolicy      []ScheduledPolicy `json:"scheduled_policy,omitempty"`
	VPC                  []VPCDetail       `json:"vpc,omitempty"`
}

// ElasticPolicy represents an elastic scaling policy
type ElasticPolicy struct {
	Type          string `json:"type"`
	Adjust        int    `json:"adjust"`
	Parameter     string `json:"parameter"`
	Operator      string `json:"operator"`
	Value         string `json:"value"`
	PeriodNumber  string `json:"period_number"`
	PeriodSeconds string `json:"period"`
	Cooldown      string `json:"cooldown"`
}

// ScheduledPolicy represents a scheduled scaling policy
type ScheduledPolicy struct {
	Type       string `json:"type"`
	Adjust     string `json:"adjust"`
	Recurrence string `json:"recurrence"`
}

// VPCDetail represents VPC details
type VPCDetail struct {
	Name      string         `json:"name,omitempty"`
	NetworkID int            `json:"network_id"`
	IPv4CIDR  string         `json:"ipv4_cidr,omitempty"`
	State     string         `json:"state,omitempty"`
	Subnets   []SubnetDetail `json:"subnets,omitempty"`
}

// VPCPartial represents partial VPC information
type VPCPartial struct {
	Name      string `json:"name"`
	NetworkID int    `json:"network_id"`
	IPv4CIDR  string `json:"ipv4_cidr"`
}

// ScalerGroupUpdateRequest represents a request to update a scaler group
type ScalerGroupUpdateRequest struct {
	Name            string            `json:"name"`
	PlanID          string            `json:"plan_id"`
	MinNodes        int               `json:"min_nodes"`
	MaxNodes        int               `json:"max_nodes"`
	PolicyType      string            `json:"policy_type,omitempty"`
	Policy          []ElasticPolicy   `json:"policy"`
	ScheduledPolicy []ScheduledPolicy `json:"scheduled_policy"`
}

// VPCAttachRequest represents a request to attach VPCs
type VPCAttachRequest struct {
	VPC []VPCDetail `json:"vpc"`
}

// PublicIPStatus represents the public IP status
type PublicIPStatus struct {
	IsPublicIPRequired bool `json:"is_public_ip_required"`
}

// PublicIPActionResponse represents the response from public IP actions
type PublicIPActionResponse struct {
	Code    int                    `json:"code"`
	Data    string                 `json:"data"`
	Errors  map[string]interface{} `json:"errors"`
	Message string                 `json:"message"`
}

// Response wrappers for API calls
type scalerGroupRoot struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
	Data    ScalerGroup       `json:"data"`
}

type scalerGroupCreateRoot struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Errors  map[string]string        `json:"errors,omitempty"`
	Data    scalerGroupCreateDetails `json:"data"`
}

type scalerGroupCreateDetails struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	VMImageName     string `json:"vm_image_name"`
	ProvisionStatus string `json:"provision_status"`
	Running         int    `json:"running"`
	Desired         int    `json:"desired"`
	Tags            string `json:"tags"`
	MinNodes        int    `json:"min_nodes"`
	MaxNodes        int    `json:"max_nodes"`
	CustomerID      int    `json:"customer_id"`
	PlanID          int    `json:"plan_id"`
	ImageID         int    `json:"image_id"`
	VPCNames        string `json:"vpc_names"`
}

type scalerGroupListRoot struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []ScalerGroup `json:"data"`
}

type vpcPartialListRoot struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []VPCPartial `json:"data"`
}

type publicIPStatusRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
	Data    PublicIPStatus         `json:"data"`
}

// CreateScalerGroup creates a new scaler group.
func (s *AutoscalingServiceOp) CreateScalerGroup(ctx context.Context, createReq *ScalerGroupCreateRequest) (*ScalerGroup, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("name", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, scalerGroupsPath, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for scaler group (%s): %w", createReq.Name, err)
	}

	root := new(scalerGroupCreateRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create scaler group (%s): %w", createReq.Name, err)
	}

	// Convert create response to ScalerGroup
	scalerGroup := &ScalerGroup{
		ID:              root.Data.ID,
		Name:            root.Data.Name,
		VMImageName:     root.Data.VMImageName,
		ProvisionStatus: root.Data.ProvisionStatus,
		Running:         root.Data.Running,
		Desired:         root.Data.Desired,
		MinNodes:        root.Data.MinNodes,
		MaxNodes:        root.Data.MaxNodes,
		PlanID:          root.Data.PlanID,
		ImageID:         root.Data.ImageID,
	}

	return scalerGroup, resp, nil
}

// GetScalerGroup retrieves a scaler group by ID.
func (s *AutoscalingServiceOp) GetScalerGroup(ctx context.Context, asgID string) (*ScalerGroup, *Response, error) {
	if asgID == "" {
		return nil, nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", scalerGroupPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for scaler group (ID: %s): %w", asgID, err)
	}

	root := new(scalerGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve scaler group (ID: %s): %w", asgID, err)
	}

	return &root.Data, resp, nil
}

// UpdateScalerGroup updates a scaler group configuration.
func (s *AutoscalingServiceOp) UpdateScalerGroup(ctx context.Context, asgID string, updateReq *ScalerGroupUpdateRequest) (*Response, error) {
	if asgID == "" {
		return nil, NewArgError("asgID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, NewArgError("updateReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", scalerGroupUpdatePath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for updating scaler group (ID: %s): %w", asgID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to update scaler group (ID: %s): %w", asgID, err)
	}

	return resp, nil
}

// DeleteScalerGroup deletes a scaler group.
func (s *AutoscalingServiceOp) DeleteScalerGroup(ctx context.Context, asgID string) (*Response, error) {
	if asgID == "" {
		return nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", scalerGroupPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting scaler group (ID: %s): %w", asgID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete scaler group (ID: %s): %w", asgID, err)
	}
	return resp, nil
}

// ListScalerGroups lists all scaler groups.
func (s *AutoscalingServiceOp) ListScalerGroups(ctx context.Context) ([]ScalerGroup, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, scalerGroupsPath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing scaler groups: %w", err)
	}

	root := new(scalerGroupListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list scaler groups: %w", err)
	}

	return root.Data, resp, nil
}

// StartScalerGroup starts a scaler group.
func (s *AutoscalingServiceOp) StartScalerGroup(ctx context.Context, asgID string) (*Response, error) {
	return s.UpdateScalerGroupStatus(ctx, asgID, "Running")
}

// StopScalerGroup stops a scaler group.
func (s *AutoscalingServiceOp) StopScalerGroup(ctx context.Context, asgID string) (*Response, error) {
	return s.UpdateScalerGroupStatus(ctx, asgID, "Stopped")
}

// UpdateScalerGroupStatus updates the status of a scaler group (start/stop).
func (s *AutoscalingServiceOp) UpdateScalerGroupStatus(ctx context.Context, asgID string, status string) (*Response, error) {
	if asgID == "" {
		return nil, NewArgError("asgID", "cannot be empty")
	}

	var path string
	switch status {
	case "Stopped":
		path = fmt.Sprintf("%s/%s/stop/", scalerGroupStopPath, asgID)
	case "Running":
		path = fmt.Sprintf("%s/%s/start/", scalerGroupStartPath, asgID)
	default:
		return nil, fmt.Errorf("unsupported status value: %s (must be 'Running' or 'Stopped')", status)
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, struct{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to create request for updating scaler group status (ID: %s, status: %s): %w", asgID, status, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to update scaler group status (ID: %s, status: %s): %w", asgID, status, err)
	}

	return resp, nil
}

// UpdateDesiredNodeCount updates the desired node count of a scaler group.
func (s *AutoscalingServiceOp) UpdateDesiredNodeCount(ctx context.Context, asgID string, desired int) (*Response, error) {
	if asgID == "" {
		return nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", scalerGroupPath, asgID)

	payload := map[string]int{"cardinality": desired}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for updating desired node count (ID: %s): %w", asgID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to update desired node count (ID: %s): %w", asgID, err)
	}

	return resp, nil
}

// AttachVPCToScalerGroup attaches VPCs to a scaler group.
func (s *AutoscalingServiceOp) AttachVPCToScalerGroup(ctx context.Context, asgID string, attachReq *VPCAttachRequest) (*Response, error) {
	if asgID == "" {
		return nil, NewArgError("asgID", "cannot be empty")
	}
	if attachReq == nil {
		return nil, NewArgError("attachReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/vpc/action/", scalerGroupVPCActionPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, attachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching VPC to scaler group (ID: %s): %w", asgID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach VPC to scaler group (ID: %s): %w", asgID, err)
	}

	return resp, nil
}

// DetachVPCFromScalerGroup detaches a VPC from a scaler group.
func (s *AutoscalingServiceOp) DetachVPCFromScalerGroup(ctx context.Context, asgID string, vpcID string) (*Response, error) {
	if asgID == "" {
		return nil, NewArgError("asgID", "cannot be empty")
	}
	if vpcID == "" {
		return nil, NewArgError("vpcID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/vpc/action/", scalerGroupVPCActionPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching VPC from scaler group (ID: %s): %w", asgID, err)
	}

	// Add vpc_id as query parameter
	req.URL.RawQuery = url.Values{"vpc_id": []string{vpcID}}.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach VPC from scaler group (ID: %s): %w", asgID, err)
	}

	return resp, nil
}

// GetAttachedVPCsForScalerGroup retrieves the list of VPCs attached to a scaler group.
func (s *AutoscalingServiceOp) GetAttachedVPCsForScalerGroup(ctx context.Context, asgID string) ([]VPCPartial, *Response, error) {
	if asgID == "" {
		return nil, nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/vpc/action/", scalerGroupVPCActionPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting attached VPCs (ID: %s): %w", asgID, err)
	}

	root := new(vpcPartialListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get attached VPCs (ID: %s): %w", asgID, err)
	}

	return root.Data, resp, nil
}

// AttachSecurityGroupToScalerGroup attaches a security group to a scaler group.
func (s *AutoscalingServiceOp) AttachSecurityGroupToScalerGroup(ctx context.Context, asgID string, sgID int) (*Response, error) {
	if asgID == "" {
		return nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", scalerGroupSecurityGroupsPath, asgID)

	payload := map[string]int{"security_group_id": sgID}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching security group to scaler group (ID: %s): %w", asgID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach security group to scaler group (ID: %s): %w", asgID, err)
	}

	return resp, nil
}

// DetachSecurityGroupFromScalerGroup detaches a security group from a scaler group.
func (s *AutoscalingServiceOp) DetachSecurityGroupFromScalerGroup(ctx context.Context, asgID string, sgID int) (*Response, error) {
	if asgID == "" {
		return nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", scalerGroupSecurityGroupsPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching security group from scaler group (ID: %s): %w", asgID, err)
	}

	// Add security_group_id as query parameter
	req.URL.RawQuery = url.Values{"security_group_id": []string{strconv.Itoa(sgID)}}.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach security group from scaler group (ID: %s): %w", asgID, err)
	}

	return resp, nil
}

// AttachPublicIPToScalerGroup attaches a public IP to a scaler group.
func (s *AutoscalingServiceOp) AttachPublicIPToScalerGroup(ctx context.Context, asgID string) (*PublicIPActionResponse, *Response, error) {
	if asgID == "" {
		return nil, nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/public_ip/action/", scalerGroupPublicIPPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for attaching public IP to scaler group (ID: %s): %w", asgID, err)
	}

	root := new(PublicIPActionResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to attach public IP to scaler group (ID: %s): %w", asgID, err)
	}

	return root, resp, nil
}

// DetachPublicIPFromScalerGroup detaches a public IP from a scaler group.
func (s *AutoscalingServiceOp) DetachPublicIPFromScalerGroup(ctx context.Context, asgID string) (*PublicIPActionResponse, *Response, error) {
	if asgID == "" {
		return nil, nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/public_ip/action/", scalerGroupPublicIPPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for detaching public IP from scaler group (ID: %s): %w", asgID, err)
	}

	root := new(PublicIPActionResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to detach public IP from scaler group (ID: %s): %w", asgID, err)
	}

	return root, resp, nil
}

// GetPublicIPStatus retrieves the public IP status of a scaler group.
func (s *AutoscalingServiceOp) GetPublicIPStatus(ctx context.Context, asgID string) (*PublicIPStatus, *Response, error) {
	if asgID == "" {
		return nil, nil, NewArgError("asgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/public_ip/action/", scalerGroupPublicIPPath, asgID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting public IP status (ID: %s): %w", asgID, err)
	}

	root := new(publicIPStatusRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get public IP status (ID: %s): %w", asgID, err)
	}

	return &root.Data, resp, nil
}
