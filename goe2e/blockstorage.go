package goe2e

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	blockStoragePath      = "block_storage"
	blockStoragePlansPath = "block_storage/plans"

	errBlockStorageCreateExtractID = "failed to extract block storage ID from create response"
)

// BlockStorageService is an interface for interacting with the Block Storage endpoints
// of the E2E Networks API.
type BlockStorageService interface {
	// Block Storage lifecycle operations
	CreateBlockStorage(context.Context, *BlockStorageCreateRequest) (*BlockStorage, *Response, error)
	GetBlockStorage(context.Context, string) (*BlockStorage, *Response, error)
	DeleteBlockStorage(context.Context, string) (*Response, error)

	// Block Storage upgrade operations
	UpgradeBlockStorage(context.Context, string, *BlockStorageUpgradeRequest) (*Response, error)

	// Block Storage attachment operations
	AttachBlockStorage(context.Context, string, *BlockStorageAttachRequest) (*Response, error)
	DetachBlockStorage(context.Context, string, *BlockStorageAttachRequest) (*Response, error)

	// Block Storage plan operations
	GetBlockStoragePlans(context.Context) ([]BlockStoragePlan, *Response, error)
}

// BlockStorageServiceOp handles communication with Block Storage related methods of the
// E2E Networks API.
type BlockStorageServiceOp struct {
	client *Client
}

var _ BlockStorageService = &BlockStorageServiceOp{}

// BlockStorageCreateRequest represents a request to create block storage.
type BlockStorageCreateRequest struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
	IOPS string  `json:"iops"`
}

// BlockStorageUpgradeRequest represents a request to upgrade block storage.
type BlockStorageUpgradeRequest struct {
	Name string  `json:"name"`
	Size float64 `json:"block_storage_size"`
	VMID float64 `json:"vm_id"`
}

// BlockStorageAttachRequest represents a request to attach/detach block storage.
type BlockStorageAttachRequest struct {
	VMID int `json:"vm_id"`
}

// ResponseTemplate represents the template information in block storage responses
type ResponseTemplate struct {
	DevPrefix    string `json:"DEV_PREFIX"`
	Driver       string `json:"DRIVER"`
	TotalIOPSSec string `json:"TOTAL_IOPS_SEC"`
}

// BlockStorage represents block storage volume information.
type BlockStorage struct {
	BlockID   int                    `json:"block_id"`
	Name      string                 `json:"name"`
	Size      int                    `json:"size"`
	Status    string                 `json:"status"`
	Template  ResponseTemplate       `json:"template"`
	VMDetail  map[string]interface{} `json:"vm_detail"`
	CreatedOn string                 `json:"created_on"`
}

// BlockStoragePlan represents a block storage plan/pricing tier.
type BlockStoragePlan struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	IOPS  string  `json:"iops"`
}

// Response wrappers for API calls
type blockStorageRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type blockStorageDetailRoot struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    BlockStorage `json:"data"`
}

type blockStoragePlansRoot struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    []BlockStoragePlan `json:"data"`
}

type blockStorageOperationRoot struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CreateBlockStorage creates a new block storage volume.
func (s *BlockStorageServiceOp) CreateBlockStorage(ctx context.Context, createReq *BlockStorageCreateRequest) (*BlockStorage, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("createReq.Name", "cannot be empty")
	}
	if createReq.Size <= 0 {
		return nil, nil, NewArgError("createReq.Size", "must be greater than 0")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, blockStoragePath+"/", createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for block storage (%s): %w", createReq.Name, err)
	}

	root := new(blockStorageRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create block storage (%s): %w", createReq.Name, err)
	}

	// The create response returns the storage data in a generic map,
	// we need to fetch the full block storage details
	if idVal, ok := root.Data["id"].(float64); ok {
		blockStorageID := fmt.Sprintf("%.0f", idVal)
		return s.GetBlockStorage(ctx, blockStorageID)
	}

	return nil, resp, errors.New(errBlockStorageCreateExtractID)
}

// GetBlockStorage retrieves a block storage volume by ID.
func (s *BlockStorageServiceOp) GetBlockStorage(ctx context.Context, blockStorageID string) (*BlockStorage, *Response, error) {
	if blockStorageID == "" {
		return nil, nil, NewArgError("blockStorageID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", blockStoragePath, blockStorageID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for block storage (ID: %s): %w", blockStorageID, err)
	}

	root := new(blockStorageDetailRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil storage for 404 (not found)
		if IsNotFoundResponse(resp) {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve block storage (ID: %s): %w", blockStorageID, err)
	}

	return &root.Data, resp, nil
}

// DeleteBlockStorage deletes a block storage volume.
func (s *BlockStorageServiceOp) DeleteBlockStorage(ctx context.Context, blockStorageID string) (*Response, error) {
	if blockStorageID == "" {
		return nil, NewArgError("blockStorageID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", blockStoragePath, blockStorageID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting block storage (ID: %s): %w", blockStorageID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete block storage (ID: %s): %w", blockStorageID, err)
	}
	return resp, nil
}

// UpgradeBlockStorage upgrades a block storage volume (plan/size).
func (s *BlockStorageServiceOp) UpgradeBlockStorage(ctx context.Context, blockStorageID string, upgradeReq *BlockStorageUpgradeRequest) (*Response, error) {
	if blockStorageID == "" {
		return nil, NewArgError("blockStorageID", "cannot be empty")
	}
	if upgradeReq == nil {
		return nil, NewArgError("upgradeReq", "cannot be nil")
	}
	if upgradeReq.Size <= 0 {
		return nil, NewArgError("upgradeReq.Size", "must be greater than 0")
	}

	path := fmt.Sprintf("%s/%s/vm/upgrade/", blockStoragePath, blockStorageID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, upgradeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for upgrading block storage (ID: %s): %w", blockStorageID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to upgrade block storage (ID: %s): %w", blockStorageID, err)
	}
	return resp, nil
}

// AttachBlockStorage attaches a block storage volume to a node (VM).
func (s *BlockStorageServiceOp) AttachBlockStorage(ctx context.Context, blockStorageID string, attachReq *BlockStorageAttachRequest) (*Response, error) {
	if blockStorageID == "" {
		return nil, NewArgError("blockStorageID", "cannot be empty")
	}
	if attachReq == nil {
		return nil, NewArgError("attachReq", "cannot be nil")
	}
	if attachReq.VMID <= 0 {
		return nil, NewArgError("attachReq.VMID", "must be greater than 0")
	}

	path := fmt.Sprintf("%s/%s/vm/attach/", blockStoragePath, blockStorageID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, attachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching block storage (ID: %s) to VM (%d): %w", blockStorageID, attachReq.VMID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach block storage (ID: %s) to VM (%d): %w", blockStorageID, attachReq.VMID, err)
	}
	return resp, nil
}

// DetachBlockStorage detaches a block storage volume from a node (VM).
func (s *BlockStorageServiceOp) DetachBlockStorage(ctx context.Context, blockStorageID string, detachReq *BlockStorageAttachRequest) (*Response, error) {
	if blockStorageID == "" {
		return nil, NewArgError("blockStorageID", "cannot be empty")
	}
	if detachReq == nil {
		return nil, NewArgError("detachReq", "cannot be nil")
	}
	if detachReq.VMID <= 0 {
		return nil, NewArgError("detachReq.VMID", "must be greater than 0")
	}

	path := fmt.Sprintf("%s/%s/vm/detach/", blockStoragePath, blockStorageID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, detachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching block storage (ID: %s) from VM (%d): %w", blockStorageID, detachReq.VMID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach block storage (ID: %s) from VM (%d): %w", blockStorageID, detachReq.VMID, err)
	}
	return resp, nil
}

// GetBlockStoragePlans retrieves available block storage plans.
func (s *BlockStorageServiceOp) GetBlockStoragePlans(ctx context.Context) ([]BlockStoragePlan, *Response, error) {
	path := blockStoragePlansPath + "/"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for block storage plans: %w", err)
	}

	root := new(blockStoragePlansRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve block storage plans: %w", err)
	}

	return root.Data, resp, nil
}
