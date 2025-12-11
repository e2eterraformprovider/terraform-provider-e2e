package goe2e

import (
	"context"
	"fmt"
	"net/http"

	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

// SSHKeyService is an interface for interacting with SSH key endpoints
// of the E2E Networks API.
type SSHKeyService interface {
	// SSH key CRUD operations
	CreateSSHKey(context.Context, *SSHKeyCreateRequest) (*SSHKey, *Response, error)
	GetSSHKey(context.Context, string) (*SSHKey, *Response, error)
	GetSSHKeyByLabel(context.Context, string) (*SSHKey, *Response, error)
	ListSSHKeys(context.Context) ([]SSHKey, *Response, error)
	DeleteSSHKey(context.Context, string) (*Response, error)
}

// SSHKeyServiceOp handles communication with SSH key related methods of the E2E Networks API.
type SSHKeyServiceOp struct {
	client *Client
}

var _ SSHKeyService = &SSHKeyServiceOp{}

// SSHKey represents an SSH public key
type SSHKey struct {
	PK        int    `json:"pk"`
	Label     string `json:"label"`
	SSHKey    string `json:"ssh_key"`
	Timestamp string `json:"timestamp,omitempty"`
}

// SSHKeyCreateRequest represents a request to create an SSH key
type SSHKeyCreateRequest struct {
	Label  string `json:"label"`
	SSHKey string `json:"ssh_key"`
}

// SSHKeyListResponse represents the API response for listing SSH keys
type SSHKeyListResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []SSHKey    `json:"data"`
	Error   interface{} `json:"error"`
}

// SSHKeyResponse represents the API response for single SSH key operations
type sshKeyRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

// CreateSSHKey creates a new SSH key
func (s *SSHKeyServiceOp) CreateSSHKey(ctx context.Context, createReq *SSHKeyCreateRequest) (*SSHKey, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Label == "" {
		return nil, nil, NewArgError("label", "cannot be empty")
	}
	if createReq.SSHKey == "" {
		return nil, nil, NewArgError("ssh_key", "cannot be empty")
	}

	path := goe2econstants.SSHKeysPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for creating SSH key (%s): %w", createReq.Label, err)
	}

	root := new(sshKeyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create SSH key (%s): %w", createReq.Label, err)
	}

	sshKey := &SSHKey{
		Label:  createReq.Label,
		SSHKey: createReq.SSHKey,
	}

	return sshKey, resp, nil
}

// GetSSHKey retrieves an SSH key by its primary key (ID)
func (s *SSHKeyServiceOp) GetSSHKey(ctx context.Context, pk string) (*SSHKey, *Response, error) {
	if pk == "" {
		return nil, nil, NewArgError("pk", "cannot be empty")
	}

	path := goe2econstants.SSHKeysPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting SSH key (%s): %w", pk, err)
	}

	// Add specific query parameter for pk
	q := req.URL.Query()
	q.Add("pk", pk)
	req.URL.RawQuery = q.Encode()

	root := new(SSHKeyListResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get SSH key (%s): %w", pk, err)
	}

	// Find the SSH key with matching PK
	if len(root.Data) > 0 {
		for _, key := range root.Data {
			if key.PK == 0 || fmt.Sprintf("%d", key.PK) == pk {
				return &key, resp, nil
			}
		}
	}

	return nil, resp, fmt.Errorf("SSH key with ID %s %s", pk, goe2econstants.SSHKeyNotFoundSubstring)
}

// GetSSHKeyByLabel retrieves an SSH key by its label
func (s *SSHKeyServiceOp) GetSSHKeyByLabel(ctx context.Context, label string) (*SSHKey, *Response, error) {
	if label == "" {
		return nil, nil, NewArgError("label", "cannot be empty")
	}

	path := goe2econstants.SSHKeysPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting SSH key by label (%s): %w", label, err)
	}

	// Add specific query parameter for label
	q := req.URL.Query()
	q.Add("label", label)
	req.URL.RawQuery = q.Encode()

	root := new(SSHKeyListResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get SSH key by label (%s): %w", label, err)
	}

	// Find the SSH key with matching label
	if len(root.Data) > 0 {
		for _, key := range root.Data {
			if key.Label == label {
				return &key, resp, nil
			}
		}
	}

	return nil, resp, fmt.Errorf("SSH key with label %s %s", label, goe2econstants.SSHKeyNotFoundSubstring)
}

// ListSSHKeys retrieves all SSH keys for a project and location
func (s *SSHKeyServiceOp) ListSSHKeys(ctx context.Context) ([]SSHKey, *Response, error) {
	path := goe2econstants.SSHKeysPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing SSH keys: %w", err)
	}

	root := new(SSHKeyListResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list SSH keys: %w", err)
	}

	return root.Data, resp, nil
}

// DeleteSSHKey deletes an SSH key by its primary key (ID)
func (s *SSHKeyServiceOp) DeleteSSHKey(ctx context.Context, pk string) (*Response, error) {
	if pk == "" {
		return nil, NewArgError("pk", "cannot be empty")
	}

	path := goe2econstants.DeleteSSHKeyPath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting SSH key (%s): %w", pk, err)
	}

	// Add specific query parameter for pk
	q := req.URL.Query()
	q.Add("pk", pk)
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete SSH key (%s): %w", pk, err)
	}

	return resp, nil
}
