package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	faasNamespacePath = "faas/namespace"
	faasFunctionsPath = "faas/functions"
	faasFunctionPath  = "faas/function" // + /{id}/
	faasLogsPath      = "faas/logs"     // + /{id}/
)

// FaasService is an interface for interacting with the FaaS endpoints
// of the E2E Networks API.
type FaasService interface {
	// Namespace operations
	CreateNamespace(context.Context, string) (*FaasNamespace, *Response, error)
	DeleteNamespace(context.Context, string) (*Response, error)

	// Function operations
	CreateFunction(context.Context, *FaasFunctionCreateRequest) (*FaasFunction, *Response, error)
	GetFunction(context.Context, string) (*FaasFunction, *Response, error)
	UpdateFunction(context.Context, string, *FaasFunctionUpdateRequest) (*FaasFunction, *Response, error)
	DeleteFunction(context.Context, string) (*Response, error)

	// Logs
	GetLogs(context.Context, string) (*FaasLogs, *Response, error)
}

// FaasServiceOp handles communication with FaaS related methods of the
// E2E Networks API.
type FaasServiceOp struct {
	client *Client
}

var _ FaasService = &FaasServiceOp{}

// FaasNamespace represents a FaaS namespace
type FaasNamespace struct {
	Name string `json:"name"`
}

// FaasFunction represents a FaaS function
type FaasFunction struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Runtime     string            `json:"runtime"`
	MemoryMB    int               `json:"memory_mb"`
	Timeout     int               `json:"timeout_seconds"`
	MinReplicas int               `json:"min_replicas"`
	MaxReplicas int               `json:"max_replicas"`
	Environment map[string]string `json:"environment,omitempty"`
	EndpointURL string            `json:"endpoint_url"`
	Status      string            `json:"status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// FaasFunctionCreateRequest represents a request to create a function
type FaasFunctionCreateRequest struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Runtime     string            `json:"runtime"`
	Code        string            `json:"code_inline"`
	MemoryMB    int               `json:"memory_mb"`
	Timeout     int               `json:"timeout_seconds"`
	MinReplicas int               `json:"min_replicas"`
	MaxReplicas int               `json:"max_replicas"`
	Environment map[string]string `json:"environment,omitempty"`
}

// FaasFunctionUpdateRequest represents a request to update a function
type FaasFunctionUpdateRequest struct {
	Code        *string           `json:"code_inline,omitempty"`
	MemoryMB    *int              `json:"memory_mb,omitempty"`
	Timeout     *int              `json:"timeout_seconds,omitempty"`
	MinReplicas *int              `json:"min_replicas,omitempty"`
	MaxReplicas *int              `json:"max_replicas,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// FaasLogs represents function logs
type FaasLogs struct {
	Logs []FaasLogEntry `json:"data"`
}

// FaasLogEntry represents a single log entry
type FaasLogEntry struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp,omitempty"`
}

// Response wrappers for API calls
type faasNamespaceRoot struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    FaasNamespace `json:"data"`
}

type faasFunctionRoot struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    FaasFunction `json:"data"`
}

type faasLogsRoot struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []FaasLogEntry `json:"data"`
}

// CreateNamespace creates a new FaaS namespace.
func (s *FaasServiceOp) CreateNamespace(ctx context.Context, namespace string) (*FaasNamespace, *Response, error) {
	if namespace == "" {
		return nil, nil, NewArgError("namespace", "cannot be empty")
	}

	namespaceReq := map[string]string{"name": namespace}

	req, err := s.client.NewRequest(ctx, http.MethodPost, faasNamespacePath, namespaceReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for FaaS namespace (%s): %w", namespace, err)
	}

	root := new(faasNamespaceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create FaaS namespace (%s): %w", namespace, err)
	}

	return &root.Data, resp, nil
}

// DeleteNamespace deletes a FaaS namespace.
func (s *FaasServiceOp) DeleteNamespace(ctx context.Context, namespace string) (*Response, error) {
	if namespace == "" {
		return nil, NewArgError("namespace", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, faasNamespacePath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting FaaS namespace (%s): %w", namespace, err)
	}

	// Add additional query parameter
	q := req.URL.Query()
	q.Add("namespace", namespace)
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete FaaS namespace (%s): %w", namespace, err)
	}
	return resp, nil
}

// CreateFunction creates a new FaaS function.
func (s *FaasServiceOp) CreateFunction(ctx context.Context, createReq *FaasFunctionCreateRequest) (*FaasFunction, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, faasFunctionsPath, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for FaaS function (%s) in namespace (%s): %w", createReq.Name, createReq.Namespace, err)
	}

	root := new(faasFunctionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create FaaS function (%s) in namespace (%s): %w", createReq.Name, createReq.Namespace, err)
	}

	return &root.Data, resp, nil
}

// GetFunction retrieves a FaaS function by ID.
func (s *FaasServiceOp) GetFunction(ctx context.Context, functionID string) (*FaasFunction, *Response, error) {
	if functionID == "" {
		return nil, nil, NewArgError("functionID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", faasFunctionPath, functionID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for FaaS function (ID: %s): %w", functionID, err)
	}

	root := new(faasFunctionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil function for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve FaaS function (ID: %s): %w", functionID, err)
	}

	return &root.Data, resp, nil
}

// UpdateFunction updates a FaaS function.
func (s *FaasServiceOp) UpdateFunction(ctx context.Context, functionID string, updateReq *FaasFunctionUpdateRequest) (*FaasFunction, *Response, error) {
	if functionID == "" {
		return nil, nil, NewArgError("functionID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", faasFunctionPath, functionID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for updating FaaS function (ID: %s): %w", functionID, err)
	}

	root := new(faasFunctionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update FaaS function (ID: %s): %w", functionID, err)
	}

	return &root.Data, resp, nil
}

// DeleteFunction deletes a FaaS function.
func (s *FaasServiceOp) DeleteFunction(ctx context.Context, functionID string) (*Response, error) {
	if functionID == "" {
		return nil, NewArgError("functionID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", faasFunctionPath, functionID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting FaaS function (ID: %s): %w", functionID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete FaaS function (ID: %s): %w", functionID, err)
	}
	return resp, nil
}

// GetLogs retrieves logs for a FaaS function.
func (s *FaasServiceOp) GetLogs(ctx context.Context, functionID string) (*FaasLogs, *Response, error) {
	if functionID == "" {
		return nil, nil, NewArgError("functionID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", faasLogsPath, functionID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for FaaS function (ID: %s) logs: %w", functionID, err)
	}

	root := new(faasLogsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve logs for FaaS function (ID: %s): %w", functionID, err)
	}

	logs := &FaasLogs{Logs: root.Data}
	return logs, resp, nil
}
