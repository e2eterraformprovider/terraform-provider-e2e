package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

// API paths
const (
	reserveIPPath = "reserve_ip"
)

// Reserve IP type constants - API response values
const (
	ReserveIPTypeFloating = "FloatingIP"
	ReserveIPTypePublic   = "PublicIP"
	ReserveIPTypeAddon    = "AddonIP"
)

// ReserveIPService is an interface for interacting with the Reserve IP endpoints
// of the E2E Networks API.
type ReserveIPService interface {
	ListReserveIPs(context.Context) ([]ReserveIP, *Response, error)
	GetReserveIP(context.Context, string) (*ReserveIP, *Response, error)
	CreateReserveIP(context.Context) (*ReserveIP, *Response, error)
	DeleteReserveIP(context.Context, string) (*Response, error)
	AttachFloatingIP(context.Context, *FloatingIPAttachmentRequest) (*Response, error)
	DetachFloatingIP(context.Context, *FloatingIPDetachmentRequest) (*Response, error)
}

// ReserveIPServiceOp handles communication with Reserve IP related methods of the
// E2E Networks API.
type ReserveIPServiceOp struct {
	client *Client
}

// Compile-time check to ensure ReserveIPServiceOp implements ReserveIPService
var _ ReserveIPService = &ReserveIPServiceOp{}

// ReserveIP represents a reserved IP address
type ReserveIP struct {
	ReserveID               string                   `json:"reserve_id"`
	ApplianceType           string                   `json:"appliance_type"`
	IPAddress               string                   `json:"ip_address"`
	ReservedType            string                   `json:"reserved_type"`
	Status                  string                   `json:"status"`
	VMID                    int                      `json:"vm_id,omitempty"`
	VMName                  string                   `json:"vm_name,omitempty"`
	BoughtAt                string                   `json:"bought_at"`
	URN                     string                   `json:"urn,omitempty"`
	ProjectName             string                   `json:"project_name,omitempty"`
	AttachedNodes           []AttachedNode           `json:"attached_nodes,omitempty"`
	FloatingIPAttachedNodes []FloatingIPAttachedNode `json:"floating_ip_attached_nodes,omitempty"`
}

// AttachedNode represents a node attached to a floating IP
type AttachedNode struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// FloatingIPAttachedNode represents a node attached to a floating IP with full details
// This is used when retrieving floating IP details with attached nodes information
type FloatingIPAttachedNode struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	VMID                int    `json:"vm_id,omitempty"`
	IPAddressPublic     string `json:"ip_address_public,omitempty"`
	IPAddressPrivate    string `json:"ip_address_private,omitempty"`
	StatusName          string `json:"status_name,omitempty"`
	SecurityGroupStatus string `json:"security_group_status,omitempty"`
}

// FloatingIPAttachmentRequest represents a request to attach a floating IP to nodes
type FloatingIPAttachmentRequest struct {
	IPAddress string   `json:"ip_address"`
	NodeIDs   []string `json:"node_ids"`
}

// FloatingIPDetachmentRequest represents a request to detach a floating IP from nodes
type FloatingIPDetachmentRequest struct {
	IPAddress string   `json:"ip_address"`
	NodeIDs   []string `json:"node_ids"`
}

// Response wrappers for API calls
type reserveIPRoot struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    ReserveIP `json:"data"`
}

type reserveIPListRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []ReserveIP `json:"data"`
}

// ListReserveIPs retrieves all reserved IPs for the current project/region
func (s *ReserveIPServiceOp) ListReserveIPs(ctx context.Context) ([]ReserveIP, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, reserveIPPath+"/", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing reserve IPs: %w", err)
	}

	root := new(reserveIPListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list reserve IPs: %w", err)
	}

	return root.Data, resp, nil
}

// GetReserveIP retrieves a specific reserved IP by IP address
func (s *ReserveIPServiceOp) GetReserveIP(ctx context.Context, ipAddress string) (*ReserveIP, *Response, error) {
	if ipAddress == "" {
		return nil, nil, NewArgError("ipAddress", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", reserveIPPath, ipAddress)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for reserve IP (IP: %s): %w", ipAddress, err)
	}

	root := new(reserveIPRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve reserve IP (IP: %s): %w", ipAddress, err)
	}

	return &root.Data, resp, nil
}

// CreateReserveIP creates a new reserved IP
// The API endpoint is POST /reserve_ip/ with no request body
func (s *ReserveIPServiceOp) CreateReserveIP(ctx context.Context) (*ReserveIP, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, reserveIPPath+"/", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for creating reserve IP: %w", err)
	}

	root := new(reserveIPRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return ErrorResponse directly without wrapping
		if errResp, ok := err.(*ErrorResponse); ok {
			return nil, resp, errResp
		}
		return nil, resp, fmt.Errorf("failed to create reserve IP: %w", err)
	}

	return &root.Data, resp, nil
}

// DeleteReserveIP deletes a reserved IP
func (s *ReserveIPServiceOp) DeleteReserveIP(ctx context.Context, ipAddress string) (*Response, error) {
	if ipAddress == "" {
		return nil, NewArgError("ipAddress", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", reserveIPPath, ipAddress)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting reserve IP (IP: %s): %w", ipAddress, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete reserve IP (IP: %s): %w", ipAddress, err)
	}
	return resp, nil
}

// AttachFloatingIP attaches a floating IP to one or more nodes
// The API endpoint is POST /reserve_ip/{ip_address}/attach/
func (s *ReserveIPServiceOp) AttachFloatingIP(ctx context.Context, req *FloatingIPAttachmentRequest) (*Response, error) {
	if req == nil {
		return nil, NewArgError("req", "cannot be nil")
	}
	if req.IPAddress == "" {
		return nil, NewArgError("req.IPAddress", "cannot be empty")
	}
	if len(req.NodeIDs) == 0 {
		return nil, NewArgError("req.NodeIDs", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/attach/", reserveIPPath, req.IPAddress)

	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching floating IP (IP: %s): %w", req.IPAddress, err)
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach floating IP (IP: %s): %w", req.IPAddress, err)
	}
	return resp, nil
}

// DetachFloatingIP detaches a floating IP from one or more nodes
// The API endpoint is POST /reserve_ip/{ip_address}/detach/
func (s *ReserveIPServiceOp) DetachFloatingIP(ctx context.Context, req *FloatingIPDetachmentRequest) (*Response, error) {
	if req == nil {
		return nil, NewArgError("req", "cannot be nil")
	}
	if req.IPAddress == "" {
		return nil, NewArgError("req.IPAddress", "cannot be empty")
	}
	if len(req.NodeIDs) == 0 {
		return nil, NewArgError("req.NodeIDs", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/detach/", reserveIPPath, req.IPAddress)

	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching floating IP (IP: %s): %w", req.IPAddress, err)
	}

	resp, err := s.client.Do(ctx, httpReq, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach floating IP (IP: %s): %w", req.IPAddress, err)
	}
	return resp, nil
}
