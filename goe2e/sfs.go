package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	sfsPath                 = "efs"
	sfsCreatePath           = "efs/create"
	sfsDeletePath           = "efs/delete"
	sfsBackupActivatePath   = "efs" // + /{sfs_id}/backup/activate/
	sfsBackupDeactivatePath = "efs" // + /{sfs_id}/backup/deactivate/
)

// SfsService is an interface for interacting with SFS (Shared File System) endpoints
// of the E2E Networks API.
type SfsService interface {
	// SFS CRUD operations
	CreateSfs(context.Context, *SfsCreateRequest) (*Sfs, *Response, error)
	GetSfs(context.Context, string) (*Sfs, *Response, error)
	DeleteSfs(context.Context, string) (*Response, error)
	ListSfss(context.Context) ([]Sfs, *Response, error)

	// Backup operations
	ActivateSFSBackup(context.Context, string, *ActivateSFSBackupRequest) (*Response, error)
	DeactivateSFSBackup(context.Context, string) (*Response, error)
}

// SfsServiceOp handles communication with SFS related methods of the E2E Networks API.
type SfsServiceOp struct {
	client *Client
}

var _ SfsService = &SfsServiceOp{}

// Sfs represents an E2E Shared File System instance
type Sfs struct {
	ID                  string `json:"efs_id,omitempty"`
	Name                string `json:"name"`
	Status              string `json:"status"`
	VPCID               string `json:"vpc_id"`
	DiskSize            string `json:"efs_disk_size"`
	DiskIOPS            int    `json:"disk_iops"`
	PlanName            string `json:"plan_name"`
	IsEncryptionEnabled bool   `json:"isEncryptionEnabled"`
	PrivateIPAddress    string `json:"private_endpoint,omitempty"`
	IsBackupEnabled     bool   `json:"is_backup_enabled,omitempty"`
}

// SfsCreateRequest represents a request to create an SFS
type SfsCreateRequest struct {
	Name                 string `json:"efs_name"`
	Plan                 string `json:"efs_plan_name"`
	VPCID                string `json:"vpc_id"`
	DiskSize             int    `json:"efs_disk_size"`
	DiskIOPS             int    `json:"efs_disk_iops"`
	IsEncryptionEnabled  bool   `json:"isEncryptionEnabled"`
	EncryptionPassphrase string `json:"encryption_passphrase,omitempty"`
}

// SFSBackupConfig represents the backup configuration for an SFS
type SFSBackupConfig struct {
	Enabled        bool `json:"enabled"`
	BackupNow      bool `json:"backup_now,omitempty"`
	StartingMinute int  `json:"starting_minute,omitempty"`
}

// ActivateSFSBackupRequest represents a request to activate SFS backup
type ActivateSFSBackupRequest struct {
	*SFSBackupConfig
}

// SfsCreateResponse represents the API response for creating an SFS
type sfsCreateRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// SfsGetResponse represents the API response for getting an SFS
type sfsGetRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// SfsListResponse represents the API response for listing SFS instances
type sfsListRoot struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Sfs  `json:"data"`
	Error   string `json:"error,omitempty"`
}

// CreateSfs creates a new SFS instance.
func (s *SfsServiceOp) CreateSfs(ctx context.Context, createReq *SfsCreateRequest) (*Sfs, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("name", "cannot be empty")
	}
	if createReq.Plan == "" {
		return nil, nil, NewArgError("plan", "cannot be empty")
	}
	if createReq.VPCID == "" {
		return nil, nil, NewArgError("vpc_id", "cannot be empty")
	}
	if createReq.DiskSize <= 0 {
		return nil, nil, NewArgError("disk_size", "must be greater than 0")
	}
	if createReq.DiskIOPS <= 0 {
		return nil, nil, NewArgError("disk_iops", "must be greater than 0")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, sfsCreatePath+"/", createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for SFS (%s): %w", createReq.Name, err)
	}

	root := new(sfsCreateRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create SFS (%s): %w", createReq.Name, err)
	}

	// Extract SFS ID from response data
	var sfsID string
	if id, ok := root.Data["efs_id"]; ok {
		switch v := id.(type) {
		case float64:
			sfsID = fmt.Sprintf("%.0f", v)
		case int:
			sfsID = fmt.Sprintf("%d", v)
		case string:
			sfsID = v
		}
	}

	// Build Sfs object from response
	sfs := &Sfs{
		ID:   sfsID,
		Name: createReq.Name,
	}

	return sfs, resp, nil
}

// GetSfs retrieves an SFS instance by ID.
func (s *SfsServiceOp) GetSfs(ctx context.Context, sfsID string) (*Sfs, *Response, error) {
	if sfsID == "" {
		return nil, nil, NewArgError("sfsID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", sfsPath, sfsID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for SFS (ID: %s): %w", sfsID, err)
	}

	root := new(sfsGetRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to retrieve SFS (ID: %s): %w", sfsID, err)
	}

	// Convert map[string]interface{} to Sfs struct
	sfs := s.convertMapToSfs(root.Data)
	sfs.ID = sfsID

	return sfs, resp, nil
}

// DeleteSfs deletes an SFS instance by ID.
func (s *SfsServiceOp) DeleteSfs(ctx context.Context, sfsID string) (*Response, error) {
	if sfsID == "" {
		return nil, NewArgError("sfsID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", sfsDeletePath, sfsID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting SFS (ID: %s): %w", sfsID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete SFS (ID: %s): %w", sfsID, err)
	}

	return resp, nil
}

// ListSfss retrieves all SFS instances for a project and location.
func (s *SfsServiceOp) ListSfss(ctx context.Context) ([]Sfs, *Response, error) {
	path := sfsPath + "/"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing SFS instances: %w", err)
	}

	root := new(sfsListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list SFS instances: %w", err)
	}

	return root.Data, resp, nil
}

// convertMapToSfs converts a map[string]interface{} to an Sfs struct.
// This helper function handles the dynamic nature of the API response.
func (s *SfsServiceOp) convertMapToSfs(data map[string]interface{}) *Sfs {
	sfs := &Sfs{}

	if name, ok := data["name"].(string); ok {
		sfs.Name = name
	}
	if status, ok := data["status"].(string); ok {
		sfs.Status = status
	}
	if vpcID, ok := data["vpc_id"].(string); ok {
		sfs.VPCID = vpcID
	}
	if diskSize, ok := data["efs_disk_size"].(string); ok {
		sfs.DiskSize = diskSize
	}
	if planName, ok := data["plan_name"].(string); ok {
		sfs.PlanName = planName
	}
	if privateIP, ok := data["private_endpoint"].(string); ok {
		sfs.PrivateIPAddress = privateIP
	}
	if isBackup, ok := data["is_backup_enabled"].(bool); ok {
		sfs.IsBackupEnabled = isBackup
	}
	if isEncryption, ok := data["isEncryptionEnabled"].(bool); ok {
		sfs.IsEncryptionEnabled = isEncryption
	}

	// Handle disk_iops which might be int or float64
	if diskIOPS, ok := data["disk_iops"]; ok {
		switch v := diskIOPS.(type) {
		case int:
			sfs.DiskIOPS = v
		case float64:
			sfs.DiskIOPS = int(v)
		}
	}

	return sfs
}

// ActivateSFSBackup activates backup for an SFS instance.
func (s *SfsServiceOp) ActivateSFSBackup(ctx context.Context, sfsID string, activateReq *ActivateSFSBackupRequest) (*Response, error) {
	if sfsID == "" {
		return nil, NewArgError("sfsID", "cannot be empty")
	}
	if activateReq == nil {
		return nil, NewArgError("activateReq", "cannot be nil")
	}
	if activateReq.SFSBackupConfig == nil {
		return nil, NewArgError("SFSBackupConfig", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/backup/activate/", sfsBackupActivatePath, sfsID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, activateReq.SFSBackupConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for activating backup on SFS (ID: %s): %w", sfsID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to activate backup on SFS (ID: %s): %w", sfsID, err)
	}

	return resp, nil
}

// DeactivateSFSBackup deactivates backup for an SFS instance.
func (s *SfsServiceOp) DeactivateSFSBackup(ctx context.Context, sfsID string) (*Response, error) {
	if sfsID == "" {
		return nil, NewArgError("sfsID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/backup/deactivate/", sfsBackupDeactivatePath, sfsID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deactivating backup on SFS (ID: %s): %w", sfsID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to deactivate backup on SFS (ID: %s): %w", sfsID, err)
	}

	return resp, nil
}
