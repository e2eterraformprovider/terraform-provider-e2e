package goe2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	vpcBasePath = "vpc"
	vpcListPath = "vpc/list"
)

// VpcService is an interface for interacting with VPC endpoints
// of the E2E Networks API.
type VpcService interface {
	// VPC CRUD operations
	CreateVPC(context.Context, *VpcCreateRequest) (*Vpc, *Response, error)
	GetVPC(context.Context, string) (*Vpc, *Response, error)
	ListVPCs(context.Context) ([]Vpc, *Response, error)
	DeleteVPC(context.Context, string) (*Response, error)

	// VPC utility operations
	GetVPCByName(context.Context, string) (*VpcWithSubnets, *Response, error)
}

// VpcServiceOp handles communication with VPC related methods of the E2E Networks API.
type VpcServiceOp struct {
	client *Client
}

var _ VpcService = &VpcServiceOp{}

// Vpc represents a Virtual Private Cloud
type Vpc struct {
	ID        float64 `json:"network_id"`
	Name      string  `json:"name"`
	State     string  `json:"state"`
	CreatedAt string  `json:"created_at"`
	IPv4CIDR  string  `json:"ipv4_cidr"`
	GatewayIP string  `json:"gateway_ip"`
	PoolSize  float64 `json:"pool_size"`
	IsActive  bool    `json:"is_active"`
}

// SubnetDetail represents a subnet within a VPC
type SubnetDetail struct {
	ID         int    `json:"id"`
	SubnetName string `json:"subnet_name"`
	CIDR       string `json:"cidr"`
	UsedIPs    int    `json:"usedIPs"`
	TotalIPs   int    `json:"totalIPs"`
}

// VpcWithSubnets represents a VPC with its subnets (from vpc/list endpoint)
type VpcWithSubnets struct {
	Name      string         `json:"name,omitempty"`
	NetworkID int            `json:"network_id"`
	IPv4CIDR  string         `json:"ipv4_cidr,omitempty"`
	State     string         `json:"state,omitempty"`
	Subnets   []SubnetDetail `json:"subnets,omitempty"`
}

// VpcCreateRequest represents a request to create a VPC
type VpcCreateRequest struct {
	VpcName  string `json:"vpc_name"`
	IPv4     string `json:"ipv4"`
	IsE2EVpc bool   `json:"is_e2e_vpc"`
}

// VpcListResponse represents the API response for listing VPCs
type VpcListResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []Vpc       `json:"data"`
	Error   interface{} `json:"error"`
}

// VpcResponse represents the API response for single VPC operations
type vpcRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

// CreateVPC creates a new VPC
func (s *VpcServiceOp) CreateVPC(ctx context.Context, createReq *VpcCreateRequest) (*Vpc, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.VpcName == "" {
		return nil, nil, NewArgError("vpc_name", "cannot be empty")
	}

	path := vpcBasePath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for creating VPC (%s): %w", createReq.VpcName, err)
	}

	root := new(vpcRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create VPC (%s): %w", createReq.VpcName, err)
	}

	vpc := &Vpc{
		Name:     createReq.VpcName,
		IsActive: true,
	}

	return vpc, resp, nil
}

// GetVPC retrieves a VPC by its ID
func (s *VpcServiceOp) GetVPC(ctx context.Context, vpcID string) (*Vpc, *Response, error) {
	if vpcID == "" {
		return nil, nil, NewArgError("vpcID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", vpcBasePath, vpcID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting VPC (%s): %w", vpcID, err)
	}

	root := new(vpcRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get VPC (%s): %w", vpcID, err)
	}

	vpc := &Vpc{}
	if root.Data != nil {
		body, _ := json.Marshal(root.Data)
		_ = json.Unmarshal(body, vpc)
	}
	return vpc, resp, nil
}

// ListVPCs retrieves all VPCs in a project
func (s *VpcServiceOp) ListVPCs(ctx context.Context) ([]Vpc, *Response, error) {
	path := vpcListPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing VPCs: %w", err)
	}

	root := new(VpcListResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list VPCs: %w", err)
	}

	vpcs := []Vpc{}
	if root.Data != nil {
		vpcs = root.Data
	}

	return vpcs, resp, nil
}

// DeleteVPC deletes a VPC
func (s *VpcServiceOp) DeleteVPC(ctx context.Context, vpcID string) (*Response, error) {
	if vpcID == "" {
		return nil, NewArgError("vpcID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", vpcBasePath, vpcID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting VPC (%s): %w", vpcID, err)
	}

	root := new(vpcRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to delete VPC (%s): %w", vpcID, err)
	}

	return resp, nil
}

// vpcListWithSubnetsRoot represents the response from vpc/list endpoint with subnets
type vpcListWithSubnetsRoot struct {
	Data []VpcWithSubnets `json:"data"`
}

// GetVPCByName retrieves a VPC by name with subnets included.
// This uses the vpc/list endpoint and filters by name.
func (s *VpcServiceOp) GetVPCByName(ctx context.Context, name string) (*VpcWithSubnets, *Response, error) {
	if name == "" {
		return nil, nil, NewArgError("name", "cannot be empty")
	}

	path := vpcListPath + "/?page_no=1&per_page=100"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing VPCs: %w", err)
	}

	root := new(vpcListWithSubnetsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list VPCs: %w", err)
	}

	for _, vpc := range root.Data {
		if vpc.Name == name {
			return &vpc, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("no VPC found with name %q", name)
}
