package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	volumeAttachPath = "nodes/volume/attach"
	volumeDetachPath = "nodes/volume/detach"
)

// VolumeAttachmentService is an interface for interacting with volume attachment endpoints
// of the E2E Networks API.
type VolumeAttachmentService interface {
	// Volume attachment operations
	AttachVolume(context.Context, *VolumeAttachmentRequest) (*VolumeAttachment, *Response, error)
	DetachVolume(context.Context, *VolumeDetachmentRequest) (*Response, error)
	GetAttachments(context.Context, string) ([]VolumeAttachment, *Response, error)
}

// VolumeAttachmentServiceOp handles communication with volume attachment related methods
// of the E2E Networks API.
type VolumeAttachmentServiceOp struct {
	client *Client
}

var _ VolumeAttachmentService = &VolumeAttachmentServiceOp{}

// VolumeAttachment represents a volume attached to a node
type VolumeAttachment struct {
	NodeID     string `json:"node_id"`
	VolumeID   string `json:"volume_id"`
	DeviceName string `json:"device_name,omitempty"`
	Status     string `json:"status,omitempty"`
}

// VolumeAttachmentRequest represents a request to attach a volume to a node
type VolumeAttachmentRequest struct {
	NodeID   string `json:"node_id"`
	VolumeID string `json:"volume_id"`
}

// VolumeDetachmentRequest represents a request to detach a volume from a node
type VolumeDetachmentRequest struct {
	NodeID   string `json:"node_id"`
	VolumeID string `json:"volume_id"`
}

// Response wrappers for API calls
type volumeAttachmentRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

// AttachVolume attaches a volume to a node.
func (s *VolumeAttachmentServiceOp) AttachVolume(ctx context.Context, attachReq *VolumeAttachmentRequest) (*VolumeAttachment, *Response, error) {
	if attachReq == nil {
		return nil, nil, NewArgError("attachReq", "cannot be nil")
	}
	if attachReq.NodeID == "" {
		return nil, nil, NewArgError("node_id", "cannot be empty")
	}
	if attachReq.VolumeID == "" {
		return nil, nil, NewArgError("volume_id", "cannot be empty")
	}
	path := volumeAttachPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, attachReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for attaching volume (%s) to node (%s): %w", attachReq.VolumeID, attachReq.NodeID, err)
	}

	root := new(volumeAttachmentRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to attach volume (%s) to node (%s): %w", attachReq.VolumeID, attachReq.NodeID, err)
	}

	attachment := &VolumeAttachment{
		NodeID:   attachReq.NodeID,
		VolumeID: attachReq.VolumeID,
		Status:   "attached",
	}

	return attachment, resp, nil
}

// DetachVolume detaches a volume from a node.
func (s *VolumeAttachmentServiceOp) DetachVolume(ctx context.Context, detachReq *VolumeDetachmentRequest) (*Response, error) {
	if detachReq == nil {
		return nil, NewArgError("detachReq", "cannot be nil")
	}
	if detachReq.NodeID == "" {
		return nil, NewArgError("node_id", "cannot be empty")
	}
	if detachReq.VolumeID == "" {
		return nil, NewArgError("volume_id", "cannot be empty")
	}
	path := volumeDetachPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, detachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching volume (%s) from node (%s): %w", detachReq.VolumeID, detachReq.NodeID, err)
	}

	root := new(volumeAttachmentRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to detach volume (%s) from node (%s): %w", detachReq.VolumeID, detachReq.NodeID, err)
	}

	return resp, nil
}

// GetAttachments gets all volume attachments for a node.
// Note: This is a helper function that may require querying node details.
func (s *VolumeAttachmentServiceOp) GetAttachments(ctx context.Context, nodeID string) ([]VolumeAttachment, *Response, error) {
	if nodeID == "" {
		return nil, nil, NewArgError("nodeID", "cannot be empty")
	}
	// Note: E2E API may not have a dedicated endpoint to list attachments
	// This would require getting node details and parsing attached volumes
	// For now, return a not implemented error
	return nil, nil, fmt.Errorf("listing volume attachments is not yet implemented - use node details API")
}
