package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	containerRegistrySetupPath    = "container_registry/setup-container-registry/"
	containerRegistryProjectsPath = "container_registry/projects-details/"
)

// ContainerRegistryService is an interface for interacting with the Container Registry endpoints
// of the E2E Networks API.
type ContainerRegistryService interface {
	// Create a new container registry project
	CreateContainerRegistry(context.Context, *ContainerRegistryCreateRequest) (*ContainerRegistry, *Response, error)

	// List all container registry projects with pagination
	ListContainerRegistryProjects(context.Context, *ContainerRegistryListOptions) ([]ContainerRegistry, *Response, error)

	// Get a single container registry project by ID
	GetContainerRegistry(context.Context, int) (*ContainerRegistry, *Response, error)

	// Update container registry security settings
	UpdateContainerRegistry(context.Context, string, *ContainerRegistryUpdateRequest) (*Response, error)

	// Delete a container registry project
	DeleteContainerRegistry(context.Context, *ContainerRegistryDeleteRequest) (*Response, error)
}

// ContainerRegistryServiceOp handles communication with Container Registry related methods of the
// E2E Networks API.
type ContainerRegistryServiceOp struct {
	client *Client
}

var _ ContainerRegistryService = &ContainerRegistryServiceOp{}

// ContainerRegistry represents a container registry project
type ContainerRegistry struct {
	ID               int     `json:"id"`
	ProjectName      string  `json:"project_name"`
	ProjectSize      float64 `json:"project_size"`
	DomainName       string  `json:"domain_name"`
	PreventVul       bool    `json:"prevent_vul"`
	Severity         string  `json:"severity"`
	State            string  `json:"state"`
	IsPublic         bool    `json:"is_public"`
	StorageLimit     int     `json:"storage_limit"`
	Location         string  `json:"location"`
	Customer         int     `json:"customer"`
	ProjectID        int     `json:"project_id"`
	MyAccountProject int     `json:"my_account_project"`
	Deleted          bool    `json:"deleted"`
	DeletedAt        *string `json:"deleted_at"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// ContainerRegistryCreateRequest represents a request to create a container registry
type ContainerRegistryCreateRequest struct {
	ProjectName string `json:"project_name"`
	PreventVul  string `json:"prevent_vul"` // "true" or "false"
	Severity    string `json:"severity"`    // "low", "medium", "high", "critical"
}

// ContainerRegistryUpdateRequest represents a request to update a container registry
type ContainerRegistryUpdateRequest struct {
	PreventVul string `json:"prevent_vul"` // "true" or "false"
	Severity   string `json:"severity"`    // "low", "medium", "high", "critical"
}

// ContainerRegistryDeleteRequest represents a request to delete a container registry
type ContainerRegistryDeleteRequest struct {
	CRProjectID string `json:"cr_project_id"` // The container registry project ID
	ProjectName string `json:"project_name"`  // The project name
	UserID      string `json:"user_id"`       // The customer/user ID
}

// ContainerRegistryListOptions specifies optional parameters for listing container registry projects
type ContainerRegistryListOptions struct {
	Page     int `json:"page"`      // Page number for pagination (default: 1)
	PageSize int `json:"page_size"` // Number of items per page (default: 100)
}

// Response wrappers for API calls
type containerRegistryCreateRoot struct {
	Code    int                             `json:"code"`
	Message string                          `json:"message"`
	Data    containerRegistryCreateDataRoot `json:"data"`
	Errors  map[string]interface{}          `json:"errors"`
}

type containerRegistryCreateDataRoot struct {
	SetupStatus string `json:"setup_status"`
}

type containerRegistryListRoot struct {
	Code            int                    `json:"code"`
	Message         string                 `json:"message"`
	Data            []ContainerRegistry    `json:"data"`
	Errors          map[string]interface{} `json:"errors"`
	TotalPageNumber int                    `json:"total_page_number"`
	TotalCount      int                    `json:"total_count"`
}

type containerRegistryDeleteRoot struct {
	Code    int                             `json:"code"`
	Message string                          `json:"message"`
	Data    containerRegistryDeleteDataRoot `json:"data"`
	Errors  map[string]interface{}          `json:"errors"`
}

type containerRegistryDeleteDataRoot struct {
	Status string `json:"status"`
}

// CreateContainerRegistry creates a new container registry project.
func (s *ContainerRegistryServiceOp) CreateContainerRegistry(ctx context.Context, createReq *ContainerRegistryCreateRequest) (*ContainerRegistry, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.ProjectName == "" {
		return nil, nil, NewArgError("createReq.ProjectName", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, containerRegistrySetupPath, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for container registry (%s): %w", createReq.ProjectName, err)
	}

	root := new(containerRegistryCreateRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create container registry (%s): %w", createReq.ProjectName, err)
	}

	// After creation, we need to fetch the project details to get the full object
	// since the create response only returns setup_status
	listOpts := &ContainerRegistryListOptions{Page: 1, PageSize: 100}
	projects, listResp, err := s.ListContainerRegistryProjects(ctx, listOpts)
	if err != nil {
		return nil, listResp, fmt.Errorf("failed to retrieve created container registry (%s): %w", createReq.ProjectName, err)
	}

	// Find the created project by name
	for i := range projects {
		if projects[i].ProjectName == createReq.ProjectName {
			return &projects[i], resp, nil
		}
	}

	return nil, resp, fmt.Errorf("container registry (%s) created but not found in project list", createReq.ProjectName)
}

// ListContainerRegistryProjects retrieves a list of container registry projects with pagination.
func (s *ContainerRegistryServiceOp) ListContainerRegistryProjects(ctx context.Context, opts *ContainerRegistryListOptions) ([]ContainerRegistry, *Response, error) {
	// Set defaults
	if opts == nil {
		opts = &ContainerRegistryListOptions{Page: 1, PageSize: 100}
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 100
	}

	path := containerRegistryProjectsPath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing container registry projects: %w", err)
	}

	// Add pagination query parameters
	q := req.URL.Query()
	q.Set("page", fmt.Sprintf("%d", opts.Page))
	q.Set("page_size", fmt.Sprintf("%d", opts.PageSize))
	req.URL.RawQuery = q.Encode()

	root := new(containerRegistryListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list container registry projects: %w", err)
	}

	return root.Data, resp, nil
}

// GetContainerRegistry retrieves a single container registry project by ID.
func (s *ContainerRegistryServiceOp) GetContainerRegistry(ctx context.Context, id int) (*ContainerRegistry, *Response, error) {
	if id <= 0 {
		return nil, nil, NewArgError("id", "must be greater than 0")
	}

	// The API doesn't have a direct get-by-id endpoint, so we list and filter
	listOpts := &ContainerRegistryListOptions{Page: 1, PageSize: 100}
	projects, resp, err := s.ListContainerRegistryProjects(ctx, listOpts)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get container registry (ID: %d): %w", id, err)
	}

	// Find the project by ID
	for i := range projects {
		if projects[i].ID == id {
			return &projects[i], resp, nil
		}
	}

	// Return nil for not found (similar to FaaS GetFunction pattern)
	return nil, resp, nil
}

// UpdateContainerRegistry updates the security settings of a container registry project.
func (s *ContainerRegistryServiceOp) UpdateContainerRegistry(ctx context.Context, projectName string, updateReq *ContainerRegistryUpdateRequest) (*Response, error) {
	if projectName == "" {
		return nil, NewArgError("projectName", "cannot be empty")
	}
	if updateReq == nil {
		return nil, NewArgError("updateReq", "cannot be nil")
	}

	// Build the request payload
	payload := map[string]string{
		"project_name": projectName,
		"prevent_vul":  updateReq.PreventVul,
		"severity":     updateReq.Severity,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, containerRegistrySetupPath, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for updating container registry (%s): %w", projectName, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to update container registry (%s): %w", projectName, err)
	}

	return resp, nil
}

// DeleteContainerRegistry deletes a container registry project.
func (s *ContainerRegistryServiceOp) DeleteContainerRegistry(ctx context.Context, deleteReq *ContainerRegistryDeleteRequest) (*Response, error) {
	if deleteReq == nil {
		return nil, NewArgError("deleteReq", "cannot be nil")
	}
	if deleteReq.CRProjectID == "" {
		return nil, NewArgError("deleteReq.CRProjectID", "cannot be empty")
	}
	if deleteReq.ProjectName == "" {
		return nil, NewArgError("deleteReq.ProjectName", "cannot be empty")
	}
	if deleteReq.UserID == "" {
		return nil, NewArgError("deleteReq.UserID", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, containerRegistrySetupPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting container registry (ID: %s): %w", deleteReq.CRProjectID, err)
	}

	// Add delete-specific query parameters
	q := req.URL.Query()
	q.Add("cr_project_id", deleteReq.CRProjectID)
	q.Add("project_name", deleteReq.ProjectName)
	q.Add("user_id", deleteReq.UserID)
	req.URL.RawQuery = q.Encode()

	root := new(containerRegistryDeleteRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to delete container registry (ID: %s): %w", deleteReq.CRProjectID, err)
	}

	return resp, nil
}
