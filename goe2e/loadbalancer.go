package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	loadBalancersPath       = "appliances/load-balancers"
	loadBalancerPath        = "appliances"                // + /{id}/
	loadBalancerActionsPath = "appliances/load-balancers" // + /{id}/actions/
	loadBalancerIPv6Path    = "appliances/load-balancers" // + /{id}/ipv6/
)

// LoadBalancerService is an interface for interacting with the Load Balancer
// endpoints of the E2E Networks API.
type LoadBalancerService interface {
	// Load Balancer CRUD operations
	CreateLoadBalancer(context.Context, *LoadBalancerCreateRequest) (*LoadBalancer, *Response, error)
	GetLoadBalancer(context.Context, string) (*LoadBalancer, *Response, error)
	UpdateLoadBalancer(context.Context, string, *LoadBalancerUpdateRequest) (*LoadBalancer, *Response, error)
	DeleteLoadBalancer(context.Context, string) (*Response, error)

	// Load Balancer actions
	UpdateLoadBalancerAction(context.Context, string, *LoadBalancerActionRequest) (*Response, error)
	UpdateIPv6(context.Context, string, *IPv6ActionRequest) (*Response, error)
}

// LoadBalancerServiceOp handles communication with Load Balancer related methods
// of the E2E Networks API.
type LoadBalancerServiceOp struct {
	client *Client
}

var _ LoadBalancerService = &LoadBalancerServiceOp{}

// LoadBalancer represents an E2E Cloud Load Balancer instance
type LoadBalancer struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	PlanName           string                 `json:"plan_name"`
	LBType             string                 `json:"lb_type"`
	LBMode             string                 `json:"lb_mode"`
	LBPort             string                 `json:"lb_port"`
	NodeListType       string                 `json:"node_list_type"`
	CheckBoxEnable     string                 `json:"checkbox_enable"`
	LBReserveIP        string                 `json:"lb_reserve_ip,omitempty"`
	SSLCertificateID   string                 `json:"ssl_certificate_id,omitempty"`
	SSLContext         map[string]interface{} `json:"ssl_context,omitempty"`
	EnableBitNinja     bool                   `json:"enable_bitninja"`
	Backends           []LBBackend            `json:"backends"`
	ACLList            []LBACLList            `json:"acl_list"`
	ACLMap             []LBACLMap             `json:"acl_map"`
	VPCList            []LBVPCDetail          `json:"vpc_list"`
	EnableEOSLogger    LBEOSDetail            `json:"enable_eos_logger,omitempty"`
	TCPBackend         []LBTCPBackend         `json:"tcp_backend"`
	IsIPv6Attached     bool                   `json:"is_ipv6_attached"`
	DefaultBackend     string                 `json:"default_backend"`
	Status             string                 `json:"status"`
	PowerStatus        string                 `json:"power_status,omitempty"`
	PublicIPAddress    string                 `json:"public_ip"`
	PrivateIPAddress   string                 `json:"private_ip"`
	HostTargetIPv6     string                 `json:"host_target_ipv6,omitempty"`
	RAM                string                 `json:"ram"`
	Disk               string                 `json:"disk"`
	VCPU               float64                `json:"vcpu"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
	IsCreditSufficient bool                   `json:"is_credit_sufficient"`
	Location           string                 `json:"location"`
}

// LBBackend represents a backend configuration for a Load Balancer
type LBBackend struct {
	Name           string     `json:"name"`
	Balance        string     `json:"balance"`
	CheckboxEnable bool       `json:"checkbox_enable"`
	DomainName     string     `json:"domain_name"`
	CheckURL       string     `json:"check_url"`
	Servers        []LBServer `json:"servers"`
	HTTPCheck      bool       `json:"http_check"`
	ScalerID       string     `json:"scaler_id,omitempty"`
	ScalerPort     string     `json:"scaler_port,omitempty"`
}

// LBServer represents a backend server configuration
type LBServer struct {
	BackendName string `json:"backend_name"`
	BackendIP   string `json:"backend_ip"`
	BackendPort string `json:"backend_port"`
}

// LBACLList represents an ACL rule configuration
type LBACLList struct {
	ACLName         string `json:"acl_name"`
	ACLCondition    string `json:"acl_condition"`
	ACLMatchingPath string `json:"acl_matching_path"`
}

// LBACLMap represents an ACL routing rule
type LBACLMap struct {
	ACLName           string `json:"acl_name"`
	ACLConditionState bool   `json:"acl_condition_state"`
	ACLBackend        string `json:"acl_backend"`
}

// LBVPCDetail represents a VPC attached to a Load Balancer
type LBVPCDetail struct {
	VPCName   string  `json:"vpc_name,omitempty"`
	IPv4CIDR  string  `json:"ipv4_cidr,omitempty"`
	NetworkID float64 `json:"network_id,omitempty"`
}

// LBEOSDetail represents EOS logging configuration
type LBEOSDetail struct {
	ApplianceID int    `json:"appliance_id"`
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
	Bucket      string `json:"bucket"`
}

// LBTCPBackend represents a TCP backend configuration
type LBTCPBackend struct {
	BackendName string     `json:"backend_name"`
	Port        string     `json:"port"`
	Balance     string     `json:"balance"`
	Servers     []LBServer `json:"servers"`
}

// LoadBalancerCreateRequest represents a request to create a Load Balancer
type LoadBalancerCreateRequest struct {
	PlanName         string                 `json:"plan_name"`
	LBName           string                 `json:"lb_name"`
	LBType           string                 `json:"lb_type,omitempty"`
	LBMode           string                 `json:"lb_mode"`
	LBPort           string                 `json:"lb_port"`
	NodeListType     string                 `json:"node_list_type"`
	CheckBoxEnable   string                 `json:"checkbox_enable"`
	LBReserveIP      string                 `json:"lb_reserve_ip"`
	SSLCertificateID string                 `json:"ssl_certificate_id"`
	SSLContext       map[string]interface{} `json:"ssl_context"`
	EnableBitNinja   bool                   `json:"enable_bitninja"`
	Backends         []LBBackend            `json:"backends"`
	ACLList          []LBACLList            `json:"acl_list"`
	ACLMap           []LBACLMap             `json:"acl_map"`
	VPCList          []LBVPCDetail          `json:"vpc_list"`
	EnableEOSLogger  LBEOSDetail            `json:"enable_eos_logger,omitempty"`
	TCPBackend       []LBTCPBackend         `json:"tcp_backend"`
	IsIPv6Attached   bool                   `json:"is_ipv6_attached"`
	DefaultBackend   string                 `json:"default_backend"`
	Location         string                 `json:"location"`
}

// LoadBalancerUpdateRequest represents a request to update a Load Balancer
type LoadBalancerUpdateRequest struct {
	PlanName         string                 `json:"plan_name,omitempty"`
	LBName           string                 `json:"lb_name,omitempty"`
	CheckBoxEnable   string                 `json:"checkbox_enable,omitempty"`
	SSLCertificateID string                 `json:"ssl_certificate_id,omitempty"`
	SSLContext       map[string]interface{} `json:"ssl_context,omitempty"`
	Backends         []LBBackend            `json:"backends,omitempty"`
	ACLList          []LBACLList            `json:"acl_list,omitempty"`
	ACLMap           []LBACLMap             `json:"acl_map,omitempty"`
	VPCList          []LBVPCDetail          `json:"vpc_list,omitempty"`
	EnableEOSLogger  LBEOSDetail            `json:"enable_eos_logger,omitempty"`
	TCPBackend       []LBTCPBackend         `json:"tcp_backend,omitempty"`
	DefaultBackend   string                 `json:"default_backend,omitempty"`
}

// LoadBalancerActionRequest represents a request for Load Balancer actions
// (power, rename, upgrade plan)
type LoadBalancerActionRequest struct {
	Type     string `json:"type"`
	Name     string `json:"name,omitempty"`
	PlanName string `json:"plan_name,omitempty"`
}

// IPv6ActionRequest represents a request for IPv6 operations
type IPv6ActionRequest struct {
	Action     string `json:"action"`
	DetachIPv6 string `json:"detach_ipv6,omitempty"`
}

// Response wrappers for API calls
type loadBalancerRoot struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    LoadBalancer `json:"data"`
}

// CreateLoadBalancer creates a new Load Balancer.
func (s *LoadBalancerServiceOp) CreateLoadBalancer(ctx context.Context, createReq *LoadBalancerCreateRequest) (*LoadBalancer, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.LBName == "" {
		return nil, nil, NewArgError("LBName", "cannot be empty")
	}
	if createReq.PlanName == "" {
		return nil, nil, NewArgError("PlanName", "cannot be empty")
	}
	if createReq.LBMode == "" {
		return nil, nil, NewArgError("LBMode", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, loadBalancersPath, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for Load Balancer (%s): %w", createReq.LBName, err)
	}

	root := new(loadBalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create Load Balancer (%s): %w", createReq.LBName, err)
	}

	return &root.Data, resp, nil
}

// GetLoadBalancer retrieves a Load Balancer by ID.
func (s *LoadBalancerServiceOp) GetLoadBalancer(ctx context.Context, lbID string) (*LoadBalancer, *Response, error) {
	if lbID == "" {
		return nil, nil, NewArgError("lbID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", loadBalancerPath, lbID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for Load Balancer (ID: %s): %w", lbID, err)
	}

	root := new(loadBalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil load balancer for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve Load Balancer (ID: %s): %w", lbID, err)
	}

	return &root.Data, resp, nil
}

// UpdateLoadBalancer updates a Load Balancer's backend configuration.
func (s *LoadBalancerServiceOp) UpdateLoadBalancer(ctx context.Context, lbID string, updateReq *LoadBalancerUpdateRequest) (*LoadBalancer, *Response, error) {
	if lbID == "" {
		return nil, nil, NewArgError("lbID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", loadBalancerPath, lbID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for updating Load Balancer (ID: %s): %w", lbID, err)
	}

	root := new(loadBalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update Load Balancer (ID: %s): %w", lbID, err)
	}

	return &root.Data, resp, nil
}

// DeleteLoadBalancer deletes a Load Balancer.
func (s *LoadBalancerServiceOp) DeleteLoadBalancer(ctx context.Context, lbID string) (*Response, error) {
	if lbID == "" {
		return nil, NewArgError("lbID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", loadBalancerPath, lbID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting Load Balancer (ID: %s): %w", lbID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete Load Balancer (ID: %s): %w", lbID, err)
	}
	return resp, nil
}

// UpdateLoadBalancerAction performs an action on a Load Balancer (power, rename, upgrade).
func (s *LoadBalancerServiceOp) UpdateLoadBalancerAction(ctx context.Context, lbID string, actionReq *LoadBalancerActionRequest) (*Response, error) {
	if lbID == "" {
		return nil, NewArgError("lbID", "cannot be empty")
	}
	if actionReq == nil {
		return nil, NewArgError("actionReq", "cannot be nil")
	}
	if actionReq.Type == "" {
		return nil, NewArgError("Type", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/actions/", loadBalancerActionsPath, lbID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, actionReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Load Balancer action (ID: %s, Type: %s): %w", lbID, actionReq.Type, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to execute action on Load Balancer (ID: %s, Type: %s): %w", lbID, actionReq.Type, err)
	}
	return resp, nil
}

// UpdateIPv6 performs IPv6 attachment/detachment on a Load Balancer.
func (s *LoadBalancerServiceOp) UpdateIPv6(ctx context.Context, lbID string, ipv6Req *IPv6ActionRequest) (*Response, error) {
	if lbID == "" {
		return nil, NewArgError("lbID", "cannot be empty")
	}
	if ipv6Req == nil {
		return nil, NewArgError("ipv6Req", "cannot be nil")
	}
	if ipv6Req.Action == "" {
		return nil, NewArgError("Action", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/ipv6/", loadBalancerIPv6Path, lbID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, ipv6Req)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for IPv6 action on Load Balancer (ID: %s, Action: %s): %w", lbID, ipv6Req.Action, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to execute IPv6 action on Load Balancer (ID: %s, Action: %s): %w", lbID, ipv6Req.Action, err)
	}
	return resp, nil
}
