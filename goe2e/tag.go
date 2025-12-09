package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

const (
	tagBasePath     = "label"
	tagDetailPath   = "label" // + /{id}/
	tagMappingPath  = "label/mapping"
	tagResourcePath = "label/mapping" // + /{resource_type}/{resource_id}/
)

// TagService is an interface for interacting with the Tag/Label endpoints
// of the E2E Networks API.
type TagService interface {
	// Tag CRUD operations
	CreateTag(context.Context, *TagCreateRequest) (*Tag, *Response, error)
	ListTags(context.Context) ([]Tag, *Response, error)
	GetTag(context.Context, string) (*Tag, *Response, error)
	UpdateTag(context.Context, string, *TagUpdateRequest) (*Tag, *Response, error)
	DeleteTag(context.Context, string) (*Response, error)

	// Tag mapping operations (attach/detach to resources)
	AttachTags(context.Context, string, string, []int) (*Response, error)
	DetachTags(context.Context, string, string, []int) (*Response, error)
	GetResourceTags(context.Context, string, string) ([]TagMapping, *Response, error)
}

// TagServiceOp handles communication with Tag related methods of the
// E2E Networks API.
type TagServiceOp struct {
	client *Client
}

var _ TagService = &TagServiceOp{}

// Tag represents an E2E tag (label)
type Tag struct {
	LabelID   int    `json:"label_id"`
	LabelName string `json:"label_name"`
	Metadata  string `json:"metadata,omitempty"`
}

// TagCreateRequest represents a request to create a tag
type TagCreateRequest struct {
	LabelName string `json:"label_name"`
	Metadata  string `json:"metadata,omitempty"`
}

// TagUpdateRequest represents a request to update a tag
type TagUpdateRequest struct {
	LabelName *string `json:"label_name,omitempty"`
	Metadata  *string `json:"metadata,omitempty"`
}

// TagMapping represents a tag attachment to a resource
type TagMapping struct {
	ResourceID   int    `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	LabelMapping []Tag  `json:"label_mapping"`
}

// Response wrappers for API calls
type tagCreateRoot struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    tagCreateData `json:"data"`
	Errors  interface{}   `json:"errors"`
}

type tagCreateData struct {
	LabelID int `json:"label_id"`
}

type tagListRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []Tag       `json:"data"`
	Errors  interface{} `json:"errors"`
}

type tagDeleteRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

type tagMappingRoot struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []TagMapping `json:"data"`
	Errors  interface{}  `json:"errors"`
}

type tagAttachRequest struct {
	Attach []int `json:"attach,omitempty"`
	Detach []int `json:"detach,omitempty"`
}

// CreateTag creates a new tag.
func (s *TagServiceOp) CreateTag(ctx context.Context, createReq *TagCreateRequest) (*Tag, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.LabelName == "" {
		return nil, nil, NewArgError("label_name", "cannot be empty")
	}
	path := tagBasePath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for tag (%s): %w", createReq.LabelName, err)
	}

	root := new(tagCreateRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create tag (%s): %w", createReq.LabelName, err)
	}

	tag := &Tag{
		LabelID:   root.Data.LabelID,
		LabelName: createReq.LabelName,
		Metadata:  createReq.Metadata,
	}

	return tag, resp, nil
}

// ListTags retrieves all tags for a project and location.
func (s *TagServiceOp) ListTags(ctx context.Context) ([]Tag, *Response, error) {
	path := tagBasePath + "/"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing tags: %w", err)
	}

	root := new(tagListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list tags: %w", err)
	}

	return root.Data, resp, nil
}

// GetTag retrieves a specific tag by ID.
func (s *TagServiceOp) GetTag(ctx context.Context, tagID string) (*Tag, *Response, error) {
	if tagID == "" {
		return nil, nil, NewArgError("tagID", "cannot be empty")
	}
	// Get all tags and find the specific one
	tags, resp, err := s.ListTags(ctx)
	if err != nil {
		return nil, resp, err
	}

	tagIDInt, err := strconv.Atoi(tagID)
	if err != nil {
		return nil, resp, NewArgError("tagID", "must be a valid integer")
	}

	for _, tag := range tags {
		if tag.LabelID == tagIDInt {
			return &tag, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("tag with ID %s not found", tagID)
}

// UpdateTag updates an existing tag.
// Note: The E2E API may not support tag updates. This method is provided for completeness.
func (s *TagServiceOp) UpdateTag(ctx context.Context, tagID string, updateReq *TagUpdateRequest) (*Tag, *Response, error) {
	if tagID == "" {
		return nil, nil, NewArgError("tagID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}
	// Note: E2E API may not support tag updates
	// This would require a PUT/PATCH endpoint like: PUT /label/{id}/
	return nil, nil, fmt.Errorf("tag updates are not supported by the E2E API")
}

// DeleteTag deletes a tag by ID.
func (s *TagServiceOp) DeleteTag(ctx context.Context, tagID string) (*Response, error) {
	if tagID == "" {
		return nil, NewArgError("tagID", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/", tagDetailPath, tagID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting tag (%s): %w", tagID, err)
	}

	root := new(tagDeleteRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, fmt.Errorf("failed to delete tag (%s): %w", tagID, err)
	}

	return resp, nil
}

// AttachTags attaches tags to a resource.
func (s *TagServiceOp) AttachTags(ctx context.Context, resourceType, resourceID string, tagIDs []int) (*Response, error) {
	if resourceType == "" {
		return nil, NewArgError("resourceType", "cannot be empty")
	}
	if resourceID == "" {
		return nil, NewArgError("resourceID", "cannot be empty")
	}
	if len(tagIDs) == 0 {
		return nil, NewArgError("tagIDs", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/%s/", tagResourcePath, resourceType, resourceID)
	attachReq := tagAttachRequest{
		Attach: tagIDs,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, attachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching tags to %s/%s: %w", resourceType, resourceID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach tags to %s/%s: %w", resourceType, resourceID, err)
	}

	return resp, nil
}

// DetachTags detaches tags from a resource.
func (s *TagServiceOp) DetachTags(ctx context.Context, resourceType, resourceID string, tagIDs []int) (*Response, error) {
	if resourceType == "" {
		return nil, NewArgError("resourceType", "cannot be empty")
	}
	if resourceID == "" {
		return nil, NewArgError("resourceID", "cannot be empty")
	}
	if len(tagIDs) == 0 {
		return nil, NewArgError("tagIDs", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/%s/", tagResourcePath, resourceType, resourceID)
	detachReq := tagAttachRequest{
		Detach: tagIDs,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, detachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching tags from %s/%s: %w", resourceType, resourceID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach tags from %s/%s: %w", resourceType, resourceID, err)
	}

	return resp, nil
}

// GetResourceTags retrieves all tags attached to resources of a specific type.
func (s *TagServiceOp) GetResourceTags(ctx context.Context, resourceType, resourceID string) ([]TagMapping, *Response, error) {
	if resourceType == "" {
		return nil, nil, NewArgError("resourceType", "cannot be empty")
	}
	path := fmt.Sprintf("%s/%s/", tagMappingPath, resourceType)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for getting tags for %s: %w", resourceType, err)
	}

	root := new(tagMappingRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get tags for %s: %w", resourceType, err)
	}

	// Filter by resourceID if provided
	if resourceID != "" {
		resourceIDInt, _ := strconv.Atoi(resourceID)
		var filtered []TagMapping
		for _, mapping := range root.Data {
			if mapping.ResourceID == resourceIDInt {
				filtered = append(filtered, mapping)
			}
		}
		return filtered, resp, nil
	}

	return root.Data, resp, nil
}
