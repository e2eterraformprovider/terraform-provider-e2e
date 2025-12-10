package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	firewallCreatePath = "fortigate/create"
	firewallListPath   = "fortigate/list"
	firewallDetailPath = "fortigate" // + /{id}/ for detail operations
)

// FirewallService is an interface for interacting with the Firewall (Fortigate) endpoints
// of the E2E Networks API.
type FirewallService interface {
	// Firewall CRUD operations
	CreateFirewall(context.Context, *FirewallCreateRequest) (*Firewall, *Response, error)
	GetFirewalls(context.Context) ([]Firewall, *Response, error)
	GetFirewall(context.Context, string) (*Firewall, *Response, error)
	UpdateFirewall(context.Context, string, *FirewallUpdateRequest) (*Firewall, *Response, error)
	DeleteFirewall(context.Context, string) (*Response, error)
}

// FirewallServiceOp handles communication with Firewall related methods of the
// E2E Networks API.
type FirewallServiceOp struct {
	client *Client
}

// Compile-time check to ensure FirewallServiceOp implements FirewallService
var _ FirewallService = &FirewallServiceOp{}

// Firewall represents a Fortigate firewall appliance
type Firewall struct {
	ID               string `json:"id"`
	VMID             int    `json:"vm_id"`
	Name             string `json:"name"`
	Label            string `json:"label"`
	Plan             string `json:"plan"`
	Status           string `json:"status"`
	PublicIPAddress  string `json:"public_ip_address"`
	PrivateIPAddress string `json:"private_ip_address"`
	IPv6Address      string `json:"ipv6_address,omitempty"`
	Memory           string `json:"memory"`
	Disk             string `json:"disk"`
	VCPUs            string `json:"vcpus,omitempty"`
	CreatedAt        string `json:"created_at"`
	IsActive         bool   `json:"is_active"`
	IsLocked         bool   `json:"is_locked"`
	IsFortigateVM    bool   `json:"is_fortigate_vm"`
	VPCID            string `json:"vpc_id,omitempty"`
	CNID             string `json:"cn_id,omitempty"`
}

// FirewallCreateRequest represents a request to create a firewall appliance
type FirewallCreateRequest struct {
	Name                 string   `json:"name"`
	Label                string   `json:"label,omitempty"`
	Plan                 string   `json:"plan"`
	Image                string   `json:"image"`
	VPCID                string   `json:"vpc_id"`
	CNID                 string   `json:"cn_id,omitempty"`
	SSHKeys              []string `json:"ssh_keys,omitempty"`
	StartScripts         []string `json:"start_scripts,omitempty"`
	Backups              bool     `json:"backups,omitempty"`
	EnableBitninja       bool     `json:"enable_bitninja,omitempty"`
	DisablePassword      bool     `json:"disable_password,omitempty"`
	IsSavedImage         bool     `json:"is_saved_image,omitempty"`
	SavedImageTemplateID int      `json:"saved_image_template_id,omitempty"`
	ReserveIP            string   `json:"reserve_ip,omitempty"`
	IsIPv6Availed        bool     `json:"is_ipv6_availed,omitempty"`
	DefaultPublicIP      bool     `json:"default_public_ip,omitempty"`
}

// FirewallUpdateRequest represents a request to update a firewall appliance
type FirewallUpdateRequest struct {
	Name  *string `json:"name,omitempty"`
	Label *string `json:"label,omitempty"`
}

// Response wrappers for API calls
type firewallRoot struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    Firewall `json:"data"`
}

type firewallListRoot struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    []Firewall `json:"data"`
}

// CreateFirewall creates a new firewall appliance.
func (s *FirewallServiceOp) CreateFirewall(ctx context.Context, createReq *FirewallCreateRequest) (*Firewall, *Response, error) {
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
	if createReq.VPCID == "" {
		return nil, nil, NewArgError("vpc_id", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, firewallCreatePath+"/", createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for firewall (%s): %w", createReq.Name, err)
	}

	root := new(firewallRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create firewall (%s): %w", createReq.Name, err)
	}

	return &root.Data, resp, nil
}

// GetFirewalls retrieves all firewall appliances.
func (s *FirewallServiceOp) GetFirewalls(ctx context.Context) ([]Firewall, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, firewallListPath+"/", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing firewalls: %w", err)
	}

	root := new(firewallListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list firewalls: %w", err)
	}

	return root.Data, resp, nil
}

// GetFirewall retrieves a firewall by ID.
func (s *FirewallServiceOp) GetFirewall(ctx context.Context, firewallID string) (*Firewall, *Response, error) {
	if firewallID == "" {
		return nil, nil, NewArgError("firewallID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", firewallDetailPath, firewallID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for firewall (ID: %s): %w", firewallID, err)
	}

	root := new(firewallRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil firewall for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve firewall (ID: %s): %w", firewallID, err)
	}

	return &root.Data, resp, nil
}

// UpdateFirewall updates a firewall appliance.
func (s *FirewallServiceOp) UpdateFirewall(ctx context.Context, firewallID string, updateReq *FirewallUpdateRequest) (*Firewall, *Response, error) {
	if firewallID == "" {
		return nil, nil, NewArgError("firewallID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", firewallDetailPath, firewallID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for updating firewall (ID: %s): %w", firewallID, err)
	}

	root := new(firewallRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update firewall (ID: %s): %w", firewallID, err)
	}

	return &root.Data, resp, nil
}

// DeleteFirewall deletes a firewall appliance.
func (s *FirewallServiceOp) DeleteFirewall(ctx context.Context, firewallID string) (*Response, error) {
	if firewallID == "" {
		return nil, NewArgError("firewallID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", firewallDetailPath, firewallID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting firewall (ID: %s): %w", firewallID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete firewall (ID: %s): %w", firewallID, err)
	}
	return resp, nil
}
