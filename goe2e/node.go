package goe2e

import (
	"context"
	"fmt"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

const (
	nodesPath             = "nodes"
	nodeSecurityGroupPath = "nodes/security_group"
	nodeVPCPath           = "nodes/vpc"
	nodeLCMStatePath      = "nodes/lcm_state"
	nodePlanUpgradePath   = "nodes/plan_upgrade"
	nodeActionsPath       = "nodes/actions"
)

// Node action constants for API requests
// These are aliases/re-exports from goe2e/constants/status.go for convenience
const (
	NodeActionPowerOn  = constants.NodePowerStatusOn  // Reuse from goe2e/constants/status.go
	NodeActionPowerOff = constants.NodePowerStatusOff // Reuse from goe2e/constants/status.go
	NodeActionReboot   = constants.NodeActionReboot   // Reuse from goe2e/constants/status.go
	NodeActionLockVM   = constants.NodeActionLockVM   // Reuse from goe2e/constants/status.go
	NodeActionUnlockVM = constants.NodeActionUnlockVM // Reuse from goe2e/constants/status.go
)

// NodeService is an interface for interacting with node endpoints
// of the E2E Networks API.
type NodeService interface {
	// Node CRUD operations
	CreateNode(context.Context, *NodeCreateRequest) (*Node, *Response, error)
	GetNode(context.Context, string) (*Node, *Response, error)
	ListNodes(context.Context) ([]Node, *Response, error)
	UpdateNode(context.Context, string, *NodeUpdateRequest) (*Node, *Response, error)
	DeleteNode(context.Context, string) (*Response, error)

	// Node actions
	PowerOn(context.Context, string) (*Response, error)
	PowerOff(context.Context, string) (*Response, error)
	Reboot(context.Context, string) (*Response, error)
	Reinstall(context.Context, string, *NodeReinstallRequest) (*Response, error)
	SaveImage(context.Context, string, *NodeSaveImageRequest) (*NodeSaveImageResult, *Response, error)
	LockNode(context.Context, string) (*Response, error)
	UnlockNode(context.Context, string) (*Response, error)

	// Security group operations
	AttachSecurityGroup(context.Context, string, *SecurityGroupRequest) (*Response, error)
	DetachSecurityGroup(context.Context, string, *SecurityGroupRequest) (*Response, error)
	GetSecurityGroupList(context.Context) ([]SecurityGroupInfo, *Response, error)

	// VPC operations
	AttachVPC(context.Context, string, *NodeVPCAttachRequest) (*Response, error)
	DetachVPC(context.Context, string) (*Response, error)

	// Node state operations
	GetLCMState(context.Context, string) (*NodeLCMState, *Response, error)
	UpgradePlan(context.Context, string, *NodePlanUpgradeRequest) (*Response, error)

	// SSH key operations
	UpdateSSH(context.Context, string, *SSHUpdateRequest) (*Response, error)
}

// NodeServiceOp handles communication with node related methods of the E2E Networks API.
type NodeServiceOp struct {
	client *Client
}

var _ NodeService = &NodeServiceOp{}

// Node represents an E2E Cloud node instance
type Node struct {
	ID                    string `json:"id"`
	VMID                  int    `json:"vm_id"`
	Name                  string `json:"name"`
	Label                 string `json:"label"`
	Plan                  string `json:"plan"`
	Status                string `json:"status"`
	PublicIPAddress       string `json:"public_ip_address"`
	PrivateIPAddress      string `json:"private_ip_address"`
	IPv6Address           string `json:"ipv6_address,omitempty"`
	Memory                string `json:"memory"`
	Disk                  string `json:"disk"`
	VCPUs                 string `json:"vcpus,omitempty"`
	CreatedAt             string `json:"created_at"`
	IsActive              bool   `json:"is_active"`
	IsLocked              bool   `json:"is_locked"`
	DefaultSGID           int    `json:"default_sg,omitempty"`
	Price                 string `json:"price,omitempty"`
	BitNinjaLicenseActive bool   `json:"is_bitninja_license_active,omitempty"`
}

// NodeCreateRequest represents a request to create a node
type NodeCreateRequest struct {
	Name                 string   `json:"name"`
	Label                string   `json:"label,omitempty"`
	Plan                 string   `json:"plan"`
	Image                string   `json:"image"`
	Backup               bool     `json:"backup,omitempty"`
	DefaultPublicIP      bool     `json:"default_public_ip,omitempty"`
	DisablePassword      bool     `json:"disable_password,omitempty"`
	EnableBitNinja       bool     `json:"enable_bitninja,omitempty"`
	EnableIPv6           bool     `json:"is_ipv6_availed,omitempty"`
	IsSavedImage         bool     `json:"is_saved_image,omitempty"`
	SavedImageTemplateID int      `json:"saved_image_template_id,omitempty"`
	ReserveIP            string   `json:"reserve_ip,omitempty"`
	VPCID                string   `json:"vpc_id,omitempty"`
	SecurityGroupID      int      `json:"security_group_id,omitempty"`
	SSHKeys              []string `json:"ssh_keys,omitempty"`
	StartScript          string   `json:"start_script,omitempty"`
	Disk                 int      `json:"disk,omitempty"`
	ImageID              int      `json:"image_id,omitempty"`
}

// NodeUpdateRequest represents a request to update a node
type NodeUpdateRequest struct {
	Name  string `json:"name,omitempty"`
	Label string `json:"label,omitempty"`
}

// NodeReinstallRequest represents a request to reinstall a node
type NodeReinstallRequest struct {
	Image string `json:"image,omitempty"`
}

// NodeSaveImageRequest represents a request to save/create an image from a node
type NodeSaveImageRequest struct {
	ActionType string `json:"action_type"`
	Name       string `json:"name"`
}

// NodeSaveImageResult represents the result of saving an image from a node
type NodeSaveImageResult struct {
	ImageID string `json:"id"`
}

// SecurityGroupRequest represents a request to attach/detach a security group
type SecurityGroupRequest struct {
	SecurityGroupList []int `json:"security_group_list,omitempty"`
}

// SecurityGroupInfo represents a simplified security group info (used by node service)
type SecurityGroupInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// NodeVPCAttachRequest represents a request to attach a VPC to a node
type NodeVPCAttachRequest struct {
	VPCID string `json:"vpc_id"`
}

// NodeLCMState represents the lifecycle state of a node
type NodeLCMState struct {
	LCMState string `json:"lcm_state"`
	State    string `json:"state,omitempty"`
}

// NodePlanUpgradeRequest represents a request to upgrade a node's plan
type NodePlanUpgradeRequest struct {
	Plan  string `json:"plan"`
	Image string `json:"image,omitempty"`
}

// SSHUpdateRequest represents a request to update SSH keys for a node
type SSHUpdateRequest struct {
	Action  string                   `json:"type"`     // e.g., "add_ssh_keys"
	SSHKeys []map[string]interface{} `json:"ssh_keys"` // Format: [{"label": "ssh-key-1", "ssh_key": "key"}]
}

// NodeActionRequest represents a generic action request for node operations
type NodeActionRequest struct {
	Action string `json:"action"`
}

// NodeReinstallActionRequest represents a reinstall action request for node operations
type NodeReinstallActionRequest struct {
	Action string `json:"action"`
	Image  string `json:"image"`
}

// Response wrapper for API calls
type nodeRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

// NodeListResponse represents the API response for listing nodes
type NodeListResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []Node      `json:"data"`
	Errors  interface{} `json:"errors"`
}

// CreateNode creates a new node
func (s *NodeServiceOp) CreateNode(ctx context.Context, createReq *NodeCreateRequest) (*Node, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("name", "cannot be empty")
	}
	if createReq.Plan == "" {
		return nil, nil, NewArgError("plan", "cannot be empty")
	}
	if createReq.Image == "" {
		return nil, nil, NewArgError("image", "cannot be empty")
	}
	path := nodesPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for creating node (%s): %w", createReq.Name, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create node (%s): %w", createReq.Name, err)
	}

	node := &Node{
		Name: createReq.Name,
	}

	return node, resp, nil
}

// GetNode retrieves a node by ID
func (s *NodeServiceOp) GetNode(ctx context.Context, nodeID string) (*Node, *Response, error) {
	if nodeID == "" {
		return nil, nil, NewArgError("nodeID", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil node for 404 (not found)
		if IsNotFoundResponse(resp) {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to get node (%s): %w", nodeID, err)
	}

	node := &Node{}
	return node, resp, nil
}

// ListNodes retrieves a list of all nodes
func (s *NodeServiceOp) ListNodes(ctx context.Context) ([]Node, *Response, error) {
	path := nodesPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing nodes: %w", err)
	}

	root := new(NodeListResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list nodes: %w", err)
	}

	nodes := []Node{}
	if root.Data != nil {
		nodes = root.Data
	}

	return nodes, resp, nil
}

// UpdateNode updates a node's attributes
func (s *NodeServiceOp) UpdateNode(ctx context.Context, nodeID string, updateReq *NodeUpdateRequest) (*Node, *Response, error) {
	if nodeID == "" {
		return nil, nil, NewArgError("nodeID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}
	path := fmt.Sprintf("%s/%s/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for updating node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update node (%s): %w", nodeID, err)
	}

	node := &Node{}
	return node, resp, nil
}

// DeleteNode deletes a node
func (s *NodeServiceOp) DeleteNode(ctx context.Context, nodeID string) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to delete node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// PowerOn powers on a node
func (s *NodeServiceOp) PowerOn(ctx context.Context, nodeID string) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	actionReq := &NodeActionRequest{Action: NodeActionPowerOn}
	path := fmt.Sprintf("%s/%s/actions/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, actionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for powering on node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to power on node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// PowerOff powers off a node
func (s *NodeServiceOp) PowerOff(ctx context.Context, nodeID string) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	actionReq := &NodeActionRequest{Action: NodeActionPowerOff}
	path := fmt.Sprintf("%s/%s/actions/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, actionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for powering off node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to power off node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// Reboot reboots a node
func (s *NodeServiceOp) Reboot(ctx context.Context, nodeID string) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	actionReq := &NodeActionRequest{Action: NodeActionReboot}
	path := fmt.Sprintf("%s/%s/actions/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, actionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for rebooting node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to reboot node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// Reinstall reinstalls a node's OS
func (s *NodeServiceOp) Reinstall(ctx context.Context, nodeID string, reinstallReq *NodeReinstallRequest) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	if reinstallReq == nil {
		return nil, NewArgError("reinstallReq", "cannot be nil")
	}
	if reinstallReq.Image == "" {
		return nil, NewArgError("image", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/actions/", nodesPath, nodeID)
	actionReq := &NodeReinstallActionRequest{Action: constants.NodeActionReinstall, Image: reinstallReq.Image}
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, actionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for reinstalling node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to reinstall node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// SaveImage saves an image from a node (creates a new image from the node's current state)
func (s *NodeServiceOp) SaveImage(ctx context.Context, nodeID string, saveReq *NodeSaveImageRequest) (*NodeSaveImageResult, *Response, error) {
	if nodeID == "" {
		return nil, nil, NewArgError("nodeID", "cannot be empty")
	}
	if saveReq == nil {
		return nil, nil, NewArgError("saveReq", "cannot be nil")
	}
	if saveReq.ActionType == "" {
		return nil, nil, NewArgError("ActionType", "cannot be empty")
	}
	if saveReq.Name == "" {
		return nil, nil, NewArgError("Name", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, saveReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for saving image from node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to save image from node (%s): %w", nodeID, err)
	}

	// Extract the image ID from the response data
	result := &NodeSaveImageResult{}
	if data, ok := root.Data.(map[string]interface{}); ok {
		if id, ok := data["id"].(string); ok {
			result.ImageID = id
		} else if id, ok := data["id"].(float64); ok {
			// Handle numeric ID (convert to string)
			result.ImageID = fmt.Sprintf("%.0f", id)
		}
	}

	return result, resp, nil
}

// LockNode locks a node
func (s *NodeServiceOp) LockNode(ctx context.Context, nodeID string) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	actionReq := &NodeActionRequest{Action: NodeActionLockVM}
	path := fmt.Sprintf("%s/%s/actions/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, actionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for locking node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to lock node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// UnlockNode unlocks a node
func (s *NodeServiceOp) UnlockNode(ctx context.Context, nodeID string) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	actionReq := &NodeActionRequest{Action: NodeActionUnlockVM}
	path := fmt.Sprintf("%s/%s/actions/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, actionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for unlocking node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to unlock node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// AttachSecurityGroup attaches a security group to a node
func (s *NodeServiceOp) AttachSecurityGroup(ctx context.Context, nodeID string, sgReq *SecurityGroupRequest) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	if sgReq == nil {
		return nil, NewArgError("sgReq", "cannot be nil")
	}
	path := fmt.Sprintf("%s/%s/security_group/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, sgReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching security group to node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to attach security group to node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// DetachSecurityGroup detaches a security group from a node
func (s *NodeServiceOp) DetachSecurityGroup(ctx context.Context, nodeID string, sgReq *SecurityGroupRequest) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	if sgReq == nil {
		return nil, NewArgError("sgReq", "cannot be nil")
	}
	path := fmt.Sprintf("%s/%s/security_group/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, sgReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching security group from node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to detach security group from node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// GetSecurityGroupList retrieves the list of security groups
func (s *NodeServiceOp) GetSecurityGroupList(ctx context.Context) ([]SecurityGroupInfo, *Response, error) {
	path := "security_group/"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting security group list: %w", err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get security group list: %w", err)
	}

	sgs := []SecurityGroupInfo{}
	return sgs, resp, nil
}

// AttachVPC attaches a VPC to a node
func (s *NodeServiceOp) AttachVPC(ctx context.Context, nodeID string, vpcReq *NodeVPCAttachRequest) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	if vpcReq == nil {
		return nil, NewArgError("vpcReq", "cannot be nil")
	}
	path := fmt.Sprintf("%s/%s/vpc/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, vpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching VPC to node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to attach VPC to node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// DetachVPC detaches a VPC from a node
func (s *NodeServiceOp) DetachVPC(ctx context.Context, nodeID string) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/vpc/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching VPC from node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to detach VPC from node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// GetLCMState retrieves the lifecycle state of a node
func (s *NodeServiceOp) GetLCMState(ctx context.Context, nodeID string) (*NodeLCMState, *Response, error) {
	if nodeID == "" {
		return nil, nil, NewArgError("nodeID", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/lcm_state/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting LCM state of node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get LCM state of node (%s): %w", nodeID, err)
	}

	lcmState := &NodeLCMState{}
	return lcmState, resp, nil
}

// UpgradePlan upgrades a node's plan
func (s *NodeServiceOp) UpgradePlan(ctx context.Context, nodeID string, upgradeReq *NodePlanUpgradeRequest) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	if upgradeReq == nil {
		return nil, NewArgError("upgradeReq", "cannot be nil")
	}
	if upgradeReq.Plan == "" {
		return nil, NewArgError("plan", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/plan_upgrade/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, upgradeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for upgrading plan of node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to upgrade plan of node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// UpdateSSH updates SSH keys for a node
func (s *NodeServiceOp) UpdateSSH(ctx context.Context, nodeID string, sshReq *SSHUpdateRequest) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}
	if sshReq == nil {
		return nil, NewArgError("sshReq", "cannot be nil")
	}
	if sshReq.Action == "" {
		return nil, NewArgError("action", "cannot be empty")
	}
	if len(sshReq.SSHKeys) == 0 {
		return nil, NewArgError("ssh_keys", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/actions/", nodesPath, nodeID)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, sshReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for updating SSH keys for node (%s): %w", nodeID, err)
	}

	root := new(nodeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to update SSH keys for node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// GenerateSSHKeyMap converts a slice of SSH key strings/interfaces into the format
// required by the API: [{"label": "ssh-key-1", "ssh_key": "key"}, ...]
func GenerateSSHKeyMap(keys []interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	for i, key := range keys {
		sshKeyMap := make(map[string]interface{})
		label := fmt.Sprintf("ssh-key-%d", i+1)
		sshKeyMap["label"] = label
		sshKeyMap["ssh_key"] = key
		result = append(result, sshKeyMap)
	}

	return result
}
