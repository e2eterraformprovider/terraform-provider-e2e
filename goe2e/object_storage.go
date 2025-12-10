package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	objectStorageBucketsPath    = "storage/buckets"           // GET (list), POST (with /{name}/)
	objectStorageBucketPath     = "storage/buckets"           // + /{name}/ for GET, DELETE
	objectStorageVersioningPath = "storage/bucket_versioning" // + /{name}/ for PUT
)

// ObjectStorageService is an interface for interacting with the Object Storage endpoints
// of the E2E Networks API.
type ObjectStorageService interface {
	// Bucket operations
	CreateBucket(context.Context, *BucketCreateRequest) (*Bucket, *Response, error)
	GetBucket(context.Context, string) (*Bucket, *Response, error)
	ListBuckets(context.Context) ([]Bucket, *Response, error)
	DeleteBucket(context.Context, string) (*Response, error)

	// Bucket versioning
	SetBucketVersioning(context.Context, string, *BucketVersioningRequest) (*BucketVersioning, *Response, error)
}

// ObjectStorageServiceOp handles communication with Object Storage related methods of the
// E2E Networks API.
type ObjectStorageServiceOp struct {
	client *Client
}

var _ ObjectStorageService = &ObjectStorageServiceOp{}

// Bucket represents an object storage bucket
type Bucket struct {
	ID                           int    `json:"id"`
	Name                         string `json:"name"`
	Status                       string `json:"status"`
	BucketSize                   string `json:"bucket_size"`
	CreatedAt                    string `json:"created_at"`
	VersioningStatus             string `json:"versioning_status"`
	LifecycleConfigurationStatus string `json:"lifecycle_configuration_status"`
}

// BucketCreateRequest represents a request to create a bucket
type BucketCreateRequest struct {
	BucketName string `json:"bucket_name"`
}

// BucketVersioningRequest represents a request to update bucket versioning
type BucketVersioningRequest struct {
	BucketName         string `json:"bucket_name"`
	NewVersioningState string `json:"new_versioning_state"`
}

// BucketVersioning represents versioning information for a bucket
type BucketVersioning struct {
	BucketName       string `json:"bucket_name"`
	VersioningStatus string `json:"versioning_status"`
}

// Response wrappers for API calls
type bucketRoot struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Bucket `json:"data"`
}

type bucketsListRoot struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []Bucket `json:"data"`
	Error   string   `json:"error"`
}

type bucketVersioningRoot struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    BucketVersioning `json:"data"`
}

// CreateBucket creates a new object storage bucket.
func (s *ObjectStorageServiceOp) CreateBucket(ctx context.Context, createReq *BucketCreateRequest) (*Bucket, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.BucketName == "" {
		return nil, nil, NewArgError("createReq.BucketName", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", objectStorageBucketsPath, createReq.BucketName)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for bucket (%s): %w", createReq.BucketName, err)
	}

	root := new(bucketRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create bucket (%s): %w", createReq.BucketName, err)
	}

	return &root.Data, resp, nil
}

// GetBucket retrieves a specific bucket by name.
func (s *ObjectStorageServiceOp) GetBucket(ctx context.Context, bucketName string) (*Bucket, *Response, error) {
	if bucketName == "" {
		return nil, nil, NewArgError("bucketName", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", objectStorageBucketPath, bucketName)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for bucket (%s): %w", bucketName, err)
	}

	root := new(bucketRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil bucket for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve bucket (%s): %w", bucketName, err)
	}

	return &root.Data, resp, nil
}

// ListBuckets retrieves all buckets in the current project and region.
func (s *ObjectStorageServiceOp) ListBuckets(ctx context.Context) ([]Bucket, *Response, error) {
	path := fmt.Sprintf("%s/", objectStorageBucketsPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing buckets: %w", err)
	}

	root := new(bucketsListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list buckets: %w", err)
	}

	return root.Data, resp, nil
}

// DeleteBucket deletes a bucket.
func (s *ObjectStorageServiceOp) DeleteBucket(ctx context.Context, bucketName string) (*Response, error) {
	if bucketName == "" {
		return nil, NewArgError("bucketName", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", objectStorageBucketPath, bucketName)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting bucket (%s): %w", bucketName, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete bucket (%s): %w", bucketName, err)
	}
	return resp, nil
}

// SetBucketVersioning updates the versioning state of a bucket.
func (s *ObjectStorageServiceOp) SetBucketVersioning(ctx context.Context, bucketName string, versioningReq *BucketVersioningRequest) (*BucketVersioning, *Response, error) {
	if bucketName == "" {
		return nil, nil, NewArgError("bucketName", "cannot be empty")
	}
	if versioningReq == nil {
		return nil, nil, NewArgError("versioningReq", "cannot be nil")
	}
	if versioningReq.NewVersioningState == "" {
		return nil, nil, NewArgError("versioningReq.NewVersioningState", "cannot be empty")
	}

	// Validate versioning state
	validStates := map[string]bool{
		"Enabled":   true,
		"Suspended": true,
	}
	if !validStates[versioningReq.NewVersioningState] {
		return nil, nil, NewArgError("versioningReq.NewVersioningState", "must be 'Enabled' or 'Suspended'")
	}

	path := fmt.Sprintf("%s/%s/", objectStorageVersioningPath, bucketName)

	// Ensure bucket name is set in request
	if versioningReq.BucketName == "" {
		versioningReq.BucketName = bucketName
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, versioningReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for bucket versioning (%s): %w", bucketName, err)
	}

	root := new(bucketVersioningRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update bucket versioning (%s): %w", bucketName, err)
	}

	return &root.Data, resp, nil
}
