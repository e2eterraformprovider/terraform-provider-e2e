package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	backupActivatePath   = "cdpbackup/activate"   // + /{node_id}/
	backupDeactivatePath = "cdpbackup/deactivate" // + /{node_id}/
	backupStatusPath     = "cdpbackup"            // + /{node_id}/
	backupAgentPath      = "cdpbackup"            // + /{node_id}/cdp-agent/
	backupListPath       = "cdpbackup"
)

// BackupService is an interface for interacting with the backup endpoints
// of the E2E Networks API.
type BackupService interface {
	// Activate backup for a node
	ActivateNodeBackup(context.Context, string, *BackupConfig) (*BackupStatus, *Response, error)
	// Deactivate backup for a node
	DeactivateNodeBackup(context.Context, string) (*Response, error)
	// Get backup status for a node
	GetNodeBackupStatus(context.Context, string) (*BackupStatus, *Response, error)
	// Get backup agent status for a node
	GetBackupAgentStatus(context.Context, string) (*BackupAgentStatus, *Response, error)
	// List all nodes with backup status
	ListNodeBackupStatus(context.Context) ([]BackupStatus, *Response, error)
}

// BackupServiceOp handles communication with backup related methods of the
// E2E Networks API.
type BackupServiceOp struct {
	client *Client
}

var _ BackupService = &BackupServiceOp{}

// BackupConfig represents the backup configuration for a node
type BackupConfig struct {
	PlanID               int    `json:"plan_id"`
	Type                 string `json:"type"` // HOURLY, DAILY, WEEKLY, MONTHLY
	ExcludeFileFolder    string `json:"exclude_file_folder,omitempty"`
	BackupNow            bool   `json:"backup_now,omitempty"`
	CompressionType      string `json:"compression_type,omitempty"`  // ZLib, GZip, None
	CompressionLevel     string `json:"compression_level,omitempty"` // Low, Medium, High
	IsEncryptionRequired bool   `json:"is_encryption_required,omitempty"`
	EncryptionPassphrase string `json:"encryption_passphrase,omitempty"`
	HoursOfDay           string `json:"hours_of_day,omitempty"` // Comma-separated hours
	StartingMinute       int    `json:"starting_minute,omitempty"`
	DBEnabled            bool   `json:"db_enabled,omitempty"`
	DBUsername           string `json:"db_username,omitempty"`
	DBPassword           string `json:"db_password,omitempty"`
}

// BackupStatus represents the current backup status of a node
type BackupStatus struct {
	Status            string `json:"status"`
	Detail            string `json:"detail,omitempty"`
	LastRecoveryPoint string `json:"last_recovery_point,omitempty"`
}

// BackupAgentStatus represents the backup agent status
type BackupAgentStatus struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

// ActivateBackupRequest represents a request to activate backup
type ActivateBackupRequest struct {
	*BackupConfig
}

// Response wrappers for API calls
type backupStatusRoot struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    BackupStatus `json:"data"`
}

type backupStatusListRoot struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []BackupStatus `json:"data"`
}

type backupAgentStatusRoot struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    BackupAgentStatus `json:"data"`
}

// ActivateNodeBackup activates backup for a node.
func (s *BackupServiceOp) ActivateNodeBackup(ctx context.Context, nodeID string, config *BackupConfig) (*BackupStatus, *Response, error) {
	if nodeID == "" {
		return nil, nil, NewArgError("nodeID", "cannot be empty")
	}
	if config == nil {
		return nil, nil, NewArgError("config", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/", backupActivatePath, nodeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for activating backup on node (%s): %w", nodeID, err)
	}

	root := new(backupStatusRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to activate backup on node (%s): %w", nodeID, err)
	}

	return &root.Data, resp, nil
}

// DeactivateNodeBackup deactivates backup for a node.
func (s *BackupServiceOp) DeactivateNodeBackup(ctx context.Context, nodeID string) (*Response, error) {
	if nodeID == "" {
		return nil, NewArgError("nodeID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", backupDeactivatePath, nodeID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deactivating backup on node (%s): %w", nodeID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to deactivate backup on node (%s): %w", nodeID, err)
	}

	return resp, nil
}

// GetNodeBackupStatus retrieves the backup status of a node.
func (s *BackupServiceOp) GetNodeBackupStatus(ctx context.Context, nodeID string) (*BackupStatus, *Response, error) {
	if nodeID == "" {
		return nil, nil, NewArgError("nodeID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", backupStatusPath, nodeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for backup status on node (%s): %w", nodeID, err)
	}

	root := new(backupStatusRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil status for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve backup status on node (%s): %w", nodeID, err)
	}

	return &root.Data, resp, nil
}

// GetBackupAgentStatus checks the backup agent status for a node.
func (s *BackupServiceOp) GetBackupAgentStatus(ctx context.Context, nodeID string) (*BackupAgentStatus, *Response, error) {
	if nodeID == "" {
		return nil, nil, NewArgError("nodeID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/cdp-agent/", backupAgentPath, nodeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for backup agent status on node (%s): %w", nodeID, err)
	}

	root := new(backupAgentStatusRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil status for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve backup agent status on node (%s): %w", nodeID, err)
	}

	return &root.Data, resp, nil
}

// ListNodeBackupStatus lists all nodes with backup status.
func (s *BackupServiceOp) ListNodeBackupStatus(ctx context.Context) ([]BackupStatus, *Response, error) {
	path := fmt.Sprintf("%s/", backupListPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for listing backup status: %w", err)
	}

	root := new(backupStatusListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list backup status: %w", err)
	}

	return root.Data, resp, nil
}

// ParseHoursOfDay converts a slice of hour strings to comma-separated format.
func ParseHoursOfDay(hours []string) string {
	return strings.Join(hours, ",")
}

// ParseHoursOfDayFromString converts comma-separated hours back to slice.
func ParseHoursOfDayFromString(hours string) []string {
	if hours == "" {
		return []string{}
	}
	parts := strings.Split(hours, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
