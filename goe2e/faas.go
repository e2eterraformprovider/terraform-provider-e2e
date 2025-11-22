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
	CreateNamespace(context.Context, string, *RequestOptions) (*FaasNamespace, *Response, error)
	DeleteNamespace(context.Context, string, *RequestOptions) (*Response, error)

	// Function operations
	CreateFunction(context.Context, *FaasFunctionCreateRequest, *RequestOptions) (*FaasFunction, *Response, error)
	GetFunction(context.Context, string, *RequestOptions) (*FaasFunction, *Response, error)
	UpdateFunction(context.Context, string, *FaasFunctionUpdateRequest, *RequestOptions) (*FaasFunction, *Response, error)
	DeleteFunction(context.Context, string, *RequestOptions) (*Response, error)

	// Logs
	GetLogs(context.Context, string, *RequestOptions) (*FaasLogs, *Response, error)
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
func (s *FaasServiceOp) CreateNamespace(ctx context.Context, namespace string, opts *RequestOptions) (*FaasNamespace, *Response, error) {
	if namespace == "" {
		return nil, nil, NewArgError("namespace", "cannot be empty")
	}
	if opts == nil {
		return nil, nil, NewArgError("opts", "cannot be nil")
	}

	namespaceReq := map[string]string{"name": namespace}

	req, err := s.client.NewRequest(ctx, http.MethodPost, faasNamespacePath, namespaceReq)
	if err != nil {
		return nil, nil, err
	}

	// Add E2E required query parameters
	q := req.URL.Query()
	q.Add("apikey", s.client.apiKey)
	q.Add("project_id", opts.ProjectID)
	q.Add("location", opts.Location)
	req.URL.RawQuery = q.Encode()

	root := new(faasNamespaceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.Data, resp, nil
}

// DeleteNamespace deletes a FaaS namespace.
func (s *FaasServiceOp) DeleteNamespace(ctx context.Context, namespace string, opts *RequestOptions) (*Response, error) {
	if namespace == "" {
		return nil, NewArgError("namespace", "cannot be empty")
	}
	if opts == nil {
		return nil, NewArgError("opts", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, faasNamespacePath, nil)
	if err != nil {
		return nil, err
	}

	// Add E2E query parameters
	q := req.URL.Query()
	q.Add("apikey", s.client.apiKey)
	q.Add("project_id", opts.ProjectID)
	q.Add("location", opts.Location)
	q.Add("namespace", namespace)
	req.URL.RawQuery = q.Encode()

	return s.client.Do(ctx, req, nil)
}

// CreateFunction creates a new FaaS function.
func (s *FaasServiceOp) CreateFunction(ctx context.Context, createReq *FaasFunctionCreateRequest, opts *RequestOptions) (*FaasFunction, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if opts == nil {
		return nil, nil, NewArgError("opts", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, faasFunctionsPath, createReq)
	if err != nil {
		return nil, nil, err
	}

	// Add E2E query parameters
	q := req.URL.Query()
	q.Add("apikey", s.client.apiKey)
	q.Add("project_id", opts.ProjectID)
	q.Add("location", opts.Location)
	req.URL.RawQuery = q.Encode()

	root := new(faasFunctionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.Data, resp, nil
}

// GetFunction retrieves a FaaS function by ID.
func (s *FaasServiceOp) GetFunction(ctx context.Context, functionID string, opts *RequestOptions) (*FaasFunction, *Response, error) {
	if functionID == "" {
		return nil, nil, NewArgError("functionID", "cannot be empty")
	}
	if opts == nil {
		return nil, nil, NewArgError("opts", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", faasFunctionPath, functionID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add E2E query parameters
	q := req.URL.Query()
	q.Add("apikey", s.client.apiKey)
	q.Add("project_id", opts.ProjectID)
	q.Add("location", opts.Location)
	req.URL.RawQuery = q.Encode()

	root := new(faasFunctionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil function for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, err
	}

	return &root.Data, resp, nil
}

// UpdateFunction updates a FaaS function.
func (s *FaasServiceOp) UpdateFunction(ctx context.Context, functionID string, updateReq *FaasFunctionUpdateRequest, opts *RequestOptions) (*FaasFunction, *Response, error) {
	if functionID == "" {
		return nil, nil, NewArgError("functionID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}
	if opts == nil {
		return nil, nil, NewArgError("opts", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", faasFunctionPath, functionID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return nil, nil, err
	}

	// Add E2E query parameters
	q := req.URL.Query()
	q.Add("apikey", s.client.apiKey)
	q.Add("project_id", opts.ProjectID)
	q.Add("location", opts.Location)
	req.URL.RawQuery = q.Encode()

	root := new(faasFunctionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.Data, resp, nil
}

// DeleteFunction deletes a FaaS function.
func (s *FaasServiceOp) DeleteFunction(ctx context.Context, functionID string, opts *RequestOptions) (*Response, error) {
	if functionID == "" {
		return nil, NewArgError("functionID", "cannot be empty")
	}
	if opts == nil {
		return nil, NewArgError("opts", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", faasFunctionPath, functionID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	// Add E2E query parameters
	q := req.URL.Query()
	q.Add("apikey", s.client.apiKey)
	q.Add("project_id", opts.ProjectID)
	q.Add("location", opts.Location)
	req.URL.RawQuery = q.Encode()

	return s.client.Do(ctx, req, nil)
}

// GetLogs retrieves logs for a FaaS function.
func (s *FaasServiceOp) GetLogs(ctx context.Context, functionID string, opts *RequestOptions) (*FaasLogs, *Response, error) {
	if functionID == "" {
		return nil, nil, NewArgError("functionID", "cannot be empty")
	}
	if opts == nil {
		return nil, nil, NewArgError("opts", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", faasLogsPath, functionID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add E2E query parameters
	q := req.URL.Query()
	q.Add("apikey", s.client.apiKey)
	q.Add("project_id", opts.ProjectID)
	q.Add("location", opts.Location)
	req.URL.RawQuery = q.Encode()

	root := new(faasLogsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	logs := &FaasLogs{Logs: root.Data}
	return logs, resp, nil
}
