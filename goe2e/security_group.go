package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	securityGroupPath  = "security_group"
	securityGroupsPath = "security_group/"
)

// SecurityGroupService is an interface for interacting with security group endpoints
// of the E2E Networks API.
type SecurityGroupService interface {
	// CRUD operations
	CreateSecurityGroup(context.Context, *SecurityGroupCreateRequest) (*SecurityGroup, *Response, error)
	GetSecurityGroup(context.Context, string) (*SecurityGroup, *Response, error)
	UpdateSecurityGroup(context.Context, string, *SecurityGroupUpdateRequest) (*SecurityGroup, *Response, error)
	DeleteSecurityGroup(context.Context, string) (*Response, error)
	GetSecurityGroupList(context.Context) ([]*SecurityGroup, *Response, error)

	// Actions
	MakeDefaultSecurityGroup(context.Context, string) (*Response, error)
	AttachSecurityGroup(context.Context, int, *SecurityGroupAttachRequest) (*Response, error)
	DetachSecurityGroup(context.Context, int, *SecurityGroupAttachRequest) (*Response, error)
}

// SecurityGroupServiceOp handles communication with security group related methods
// of the E2E Networks API.
type SecurityGroupServiceOp struct {
	client *Client
}

var _ SecurityGroupService = &SecurityGroupServiceOp{}

// SecurityGroup represents a security group
type SecurityGroup struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsDefault    bool   `json:"is_default"`
	IsAllTraffic bool   `json:"is_all_traffic_rule,omitempty"`
	Rules        []Rule `json:"rules"`
	ProjectID    string `json:"project_id,omitempty"`
	Location     string `json:"location,omitempty"`
}

// Rule represents a security group rule
type Rule struct {
	ID           int     `json:"id,omitempty"`
	RuleType     string  `json:"rule_type"`
	ProtocolName string  `json:"protocol_name"`
	PortRange    string  `json:"port_range"`
	Network      string  `json:"network"`
	NetworkCIDR  string  `json:"network_cidr,omitempty"`
	NetworkSize  *int    `json:"network_size,omitempty"`
	Description  string  `json:"description"`
	VPCID        *string `json:"vpc_id,omitempty"`
}

// SecurityGroupCreateRequest represents a request to create a security group
type SecurityGroupCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default,omitempty"`
	Rules       []Rule `json:"rules"`
}

// SecurityGroupUpdateRequest represents a request to update a security group
type SecurityGroupUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Rules       []Rule `json:"rules"`
}

// SecurityGroupAttachRequest represents a request to attach/detach a security group
type SecurityGroupAttachRequest struct {
	SecurityGroupIDs []int `json:"security_group_id"`
}

// Response wrappers for API calls
type securityGroupRoot struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    SecurityGroup `json:"data"`
	Errors  interface{}   `json:"errors"`
}

type securityGroupListRoot struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    []SecurityGroup `json:"data"`
	Errors  interface{}     `json:"errors"`
}

// CreateSecurityGroup creates a new security group
func (s *SecurityGroupServiceOp) CreateSecurityGroup(ctx context.Context, createReq *SecurityGroupCreateRequest) (*SecurityGroup, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("name", "cannot be empty")
	}

	path := securityGroupsPath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for security group (%s): %w", createReq.Name, err)
	}

	root := new(securityGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create security group (%s): %w", createReq.Name, err)
	}

	return &root.Data, resp, nil
}

// GetSecurityGroup retrieves a security group by ID
func (s *SecurityGroupServiceOp) GetSecurityGroup(ctx context.Context, sgID string) (*SecurityGroup, *Response, error) {
	if sgID == "" {
		return nil, nil, NewArgError("sgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", securityGroupPath, sgID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for security group (ID: %s): %w", sgID, err)
	}

	root := new(securityGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve security group (ID: %s): %w", sgID, err)
	}

	return &root.Data, resp, nil
}

// UpdateSecurityGroup updates a security group
func (s *SecurityGroupServiceOp) UpdateSecurityGroup(ctx context.Context, sgID string, updateReq *SecurityGroupUpdateRequest) (*SecurityGroup, *Response, error) {
	if sgID == "" {
		return nil, nil, NewArgError("sgID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", securityGroupPath, sgID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for updating security group (ID: %s): %w", sgID, err)
	}

	root := new(securityGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update security group (ID: %s): %w", sgID, err)
	}

	return &root.Data, resp, nil
}

// DeleteSecurityGroup deletes a security group
func (s *SecurityGroupServiceOp) DeleteSecurityGroup(ctx context.Context, sgID string) (*Response, error) {
	if sgID == "" {
		return nil, NewArgError("sgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", securityGroupPath, sgID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting security group (ID: %s): %w", sgID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete security group (ID: %s): %w", sgID, err)
	}
	return resp, nil
}

// GetSecurityGroupList retrieves all security groups
func (s *SecurityGroupServiceOp) GetSecurityGroupList(ctx context.Context) ([]*SecurityGroup, *Response, error) {
	path := securityGroupsPath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for security group list: %w", err)
	}

	root := new(securityGroupListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve security group list: %w", err)
	}

	// Convert to slice of pointers
	result := make([]*SecurityGroup, len(root.Data))
	for i := range root.Data {
		result[i] = &root.Data[i]
	}

	return result, resp, nil
}

// MakeDefaultSecurityGroup marks a security group as default
func (s *SecurityGroupServiceOp) MakeDefaultSecurityGroup(ctx context.Context, sgID string) (*Response, error) {
	if sgID == "" {
		return nil, NewArgError("sgID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/mark-default/", securityGroupPath, sgID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for making security group default (ID: %s): %w", sgID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to make security group default (ID: %s): %w", sgID, err)
	}
	return resp, nil
}

// AttachSecurityGroup attaches a security group to a node
func (s *SecurityGroupServiceOp) AttachSecurityGroup(ctx context.Context, vmID int, attachReq *SecurityGroupAttachRequest) (*Response, error) {
	if vmID <= 0 {
		return nil, NewArgError("vmID", "must be greater than 0")
	}
	if attachReq == nil {
		return nil, NewArgError("attachReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d/attach/", securityGroupPath, vmID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, attachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching security group to node (VM ID: %d): %w", vmID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach security group to node (VM ID: %d): %w", vmID, err)
	}
	return resp, nil
}

// DetachSecurityGroup detaches a security group from a node
func (s *SecurityGroupServiceOp) DetachSecurityGroup(ctx context.Context, vmID int, detachReq *SecurityGroupAttachRequest) (*Response, error) {
	if vmID <= 0 {
		return nil, NewArgError("vmID", "must be greater than 0")
	}
	if detachReq == nil {
		return nil, NewArgError("detachReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d/detach/", securityGroupPath, vmID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, detachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching security group from node (VM ID: %d): %w", vmID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach security group from node (VM ID: %d): %w", vmID, err)
	}
	return resp, nil
}
