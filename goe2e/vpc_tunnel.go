package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	vpcTunnelsPath       = "vpc/tunnels"
	vpcTunnelPath        = "vpc/tunnels" // + /{vpc_tunnel_id}/
	vpcTunnelPausePath   = "vpc/tunnels" // + /{vpc_tunnel_id}/pause/
	vpcTunnelRestartPath = "vpc/tunnels" // + /{vpc_tunnel_id}/restart/
	vpcTunnelDeletePath  = "vpc/tunnels" // + /{vpc_tunnel_id}/delete/
)

// VPCTunnelService is an interface for interacting with VPC Tunnel endpoints
// of the E2E Networks API.
type VPCTunnelService interface {
	// List tunnels for a VPC network
	ListVPCTunnels(context.Context, string) ([]VPCTunnel, *Response, error)
	// Create a new VPC tunnel
	CreateVPCTunnel(context.Context, *CreateVPCTunnelRequest) (*VPCTunnel, *Response, error)
	// Get a VPC tunnel by ID
	GetVPCTunnel(context.Context, string) (*VPCTunnel, *Response, error)
	// Pause a VPC tunnel
	PauseVPCTunnel(context.Context, string) (*Response, error)
	// Restart a VPC tunnel
	RestartVPCTunnel(context.Context, string) (*Response, error)
	// Delete a VPC tunnel
	DeleteVPCTunnel(context.Context, string) (*Response, error)
}

// VPCTunnelServiceOp handles communication with VPC Tunnel related methods of the E2E Networks API.
type VPCTunnelServiceOp struct {
	client *Client
}

var _ VPCTunnelService = &VPCTunnelServiceOp{}

// VPCTunnel represents a VPC Tunnel (VPC Peering connection)
type VPCTunnel struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	VPCLocalNetworkID string   `json:"vpc_local_network_id"`
	VPCPeerNetworkID  string   `json:"vpc_peer_network_id"`
	IsPeerVPCExternal bool     `json:"is_peer_vpc_external"`
	Status            string   `json:"status"`
	LocalTS           []string `json:"local_traffic_selector,omitempty"`
	RemoteTS          []string `json:"remote_traffic_selector,omitempty"`
	LocalID           string   `json:"local_gateway_ip,omitempty"`
	RemoteID          string   `json:"remote_gateway_ip,omitempty"`
	CreatedAt         string   `json:"created_at"`
}

// CreateVPCTunnelRequest represents a request to create a VPC tunnel
type CreateVPCTunnelRequest struct {
	Name              string `json:"name"`
	VPCLocalNetworkID string `json:"vpc_local_network_id"`
	VPCPeerNetworkID  string `json:"vpc_peer_network_id"`
	IsPeerVPCExternal bool   `json:"is_peer_vpc_external,omitempty"`
}

// Response wrappers for API calls
type vpcTunnelRoot struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    VPCTunnel `json:"data"`
}

type vpcTunnelListRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []VPCTunnel `json:"data"`
}

// ListVPCTunnels lists all VPC tunnels for a given VPC network ID.
func (s *VPCTunnelServiceOp) ListVPCTunnels(ctx context.Context, vpcNetworkID string) ([]VPCTunnel, *Response, error) {
	if vpcNetworkID == "" {
		return nil, nil, NewArgError("vpcNetworkID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", vpcTunnelsPath, vpcNetworkID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing VPC tunnels (VPC Network ID: %s): %w", vpcNetworkID, err)
	}

	root := new(vpcTunnelListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list VPC tunnels (VPC Network ID: %s): %w", vpcNetworkID, err)
	}

	return root.Data, resp, nil
}

// CreateVPCTunnel creates a new VPC tunnel.
func (s *VPCTunnelServiceOp) CreateVPCTunnel(ctx context.Context, createReq *CreateVPCTunnelRequest) (*VPCTunnel, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("name", "cannot be empty")
	}
	if createReq.VPCLocalNetworkID == "" {
		return nil, nil, NewArgError("vpc_local_network_id", "cannot be empty")
	}
	if createReq.VPCPeerNetworkID == "" {
		return nil, nil, NewArgError("vpc_peer_network_id", "cannot be empty")
	}

	path := vpcTunnelsPath + "/"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for VPC tunnel (%s): %w", createReq.Name, err)
	}

	root := new(vpcTunnelRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create VPC tunnel (%s): %w", createReq.Name, err)
	}

	return &root.Data, resp, nil
}

// GetVPCTunnel retrieves a VPC tunnel by ID.
func (s *VPCTunnelServiceOp) GetVPCTunnel(ctx context.Context, tunnelID string) (*VPCTunnel, *Response, error) {
	if tunnelID == "" {
		return nil, nil, NewArgError("tunnelID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", vpcTunnelPath, tunnelID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for VPC tunnel (ID: %s): %w", tunnelID, err)
	}

	root := new(vpcTunnelRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve VPC tunnel (ID: %s): %w", tunnelID, err)
	}

	return &root.Data, resp, nil
}

// PauseVPCTunnel pauses a VPC tunnel.
func (s *VPCTunnelServiceOp) PauseVPCTunnel(ctx context.Context, tunnelID string) (*Response, error) {
	if tunnelID == "" {
		return nil, NewArgError("tunnelID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/pause/", vpcTunnelPausePath, tunnelID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for pausing VPC tunnel (ID: %s): %w", tunnelID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to pause VPC tunnel (ID: %s): %w", tunnelID, err)
	}

	return resp, nil
}

// RestartVPCTunnel restarts a VPC tunnel.
func (s *VPCTunnelServiceOp) RestartVPCTunnel(ctx context.Context, tunnelID string) (*Response, error) {
	if tunnelID == "" {
		return nil, NewArgError("tunnelID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/restart/", vpcTunnelRestartPath, tunnelID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for restarting VPC tunnel (ID: %s): %w", tunnelID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to restart VPC tunnel (ID: %s): %w", tunnelID, err)
	}

	return resp, nil
}

// DeleteVPCTunnel deletes a VPC tunnel.
func (s *VPCTunnelServiceOp) DeleteVPCTunnel(ctx context.Context, tunnelID string) (*Response, error) {
	if tunnelID == "" {
		return nil, NewArgError("tunnelID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/delete/", vpcTunnelDeletePath, tunnelID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting VPC tunnel (ID: %s): %w", tunnelID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete VPC tunnel (ID: %s): %w", tunnelID, err)
	}

	return resp, nil
}
