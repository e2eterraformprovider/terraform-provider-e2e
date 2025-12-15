package volume_attachment

import (
	"strings"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

func TestParseVolumeAttachmentID(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		wantNodeID   string
		wantVolumeID string
		wantError    bool
		errorMsg     string
	}{
		{
			name:         "valid ID - node_id/volume_id",
			id:           "node-123/volume-456",
			wantNodeID:   "node-123",
			wantVolumeID: "volume-456",
			wantError:    false,
		},
		{
			name:         "valid ID - numeric IDs",
			id:           "12345/67890",
			wantNodeID:   "12345",
			wantVolumeID: "67890",
			wantError:    false,
		},
		{
			name:         "valid ID - UUID-like IDs",
			id:           "550e8400-e29b-41d4-a716-446655440000/550e8400-e29b-41d4-a716-446655440001",
			wantNodeID:   "550e8400-e29b-41d4-a716-446655440000",
			wantVolumeID: "550e8400-e29b-41d4-a716-446655440001",
			wantError:    false,
		},
		{
			name:      "invalid ID - single part",
			id:        "node-123",
			wantError: true,
			errorMsg:  ErrorParseIDTemplate,
		},
		{
			name:      "invalid ID - empty string",
			id:        "",
			wantError: true,
			errorMsg:  ErrorParseIDTemplate,
		},
		{
			name:      "invalid ID - three parts",
			id:        "project/node/volume",
			wantError: true,
			errorMsg:  ErrorParseIDTemplate,
		},
		{
			name:      "invalid ID - four parts",
			id:        "project/region/node/volume",
			wantError: true,
			errorMsg:  ErrorParseIDTemplate,
		},
		{
			name:      "invalid ID - wrong delimiter",
			id:        "node-123:volume-456",
			wantError: true,
			errorMsg:  ErrorParseIDTemplate,
		},
		{
			name:         "valid ID - empty node ID (edge case)",
			id:           "/volume-456",
			wantNodeID:   "",
			wantVolumeID: "volume-456",
			wantError:    false,
		},
		{
			name:         "valid ID - empty volume ID (edge case)",
			id:           "node-123/",
			wantNodeID:   "node-123",
			wantVolumeID: "",
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeID, volumeID, err := parseVolumeAttachmentID(tt.id)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), strings.Split(tt.errorMsg, "%")[0]) {
					t.Errorf("error message mismatch. got: %v, want to contain: %s", err.Error(), tt.errorMsg)
				}
				// Verify error message includes expected format
				if !strings.Contains(err.Error(), ImportIDFormatShortDescription) {
					t.Errorf("error message should include expected format: %s", ImportIDFormatShortDescription)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if nodeID != tt.wantNodeID {
					t.Errorf("nodeID mismatch. got: %s, want: %s", nodeID, tt.wantNodeID)
				}
				if volumeID != tt.wantVolumeID {
					t.Errorf("volumeID mismatch. got: %s, want: %s", volumeID, tt.wantVolumeID)
				}
			}
		})
	}
}

// ============================================================================
// Unit Test Coverage Notes
// ============================================================================
//
// This file contains unit tests for helper functions and constants used by the
// Volume Attachment resource. The CRUD operations (Create, Read, Delete) are
// not unit tested here because:
//
// 1. Create/Delete operations use async/polling via waitForVolumeAttachment()
//    which requires time-based testing and complex mocking of multiple services
//    (NodeService, BlockStorageService, VolumeAttachmentService)
//
// 2. Read operations require complex state verification across multiple API
//    calls (GetBlockStorage, GetNode) with VMDetail parsing
//
// 3. The async/polling behavior (waitForVolumeAttachment) is a private function
//    that uses time.NewTicker and requires deterministic time control for proper
//    unit testing, which is better suited for integration/acceptance tests
//
// Comprehensive CRUD testing is provided via acceptance tests:
// - TestAccE2EVolumeAttachment_Basic (in resource_volume_attachment_test.go)
// - TestAccE2EVolumeAttachment_Import (in resource_volume_attachment_test.go)
//
// Unit tests focus on:
// - Helper functions (parseVolumeAttachmentID)
// - Error message constants validation
// - Log message constants validation
// - Import format constants validation
//
// This approach follows the project's testing strategy where unit tests focus
// on testable, isolated logic while acceptance tests cover end-to-end behavior
// including async operations, polling, and real API interactions.

// Test error message constants contain expected format strings
func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		contains string
	}{
		{
			name:     "ErrorParseIDTemplate",
			constant: ErrorParseIDTemplate,
			contains: "invalid volume attachment ID format",
		},
		{
			name:     "ErrorImportReadDuringImportTemplate",
			constant: ErrorImportReadDuringImportTemplate,
			contains: "error reading volume attachment during import",
		},
		{
			name:     "ErrorWaitContextCancelledTemplate",
			constant: ErrorWaitContextCancelledTemplate,
			contains: "context cancelled",
		},
		{
			name:     "ErrorWaitTimeoutTemplate",
			constant: ErrorWaitTimeoutTemplate,
			contains: "timeout waiting",
		},
		{
			name:     "ErrorNodeNotFoundTemplate",
			constant: ErrorNodeNotFoundTemplate,
			contains: "node",
		},
		{
			name:     "ErrorC2PlanNoBlockStorageAttachmentFormat",
			constant: ErrorC2PlanNoBlockStorageAttachmentFormat,
			contains: "C2 plan nodes do not support block storage attachment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.constant, tt.contains) {
				t.Errorf("constant %s should contain %q, got: %q", tt.name, tt.contains, tt.constant)
			}
		})
	}
}

// Test log message constants contain expected format strings
func TestLogMessages(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		contains string
	}{
		{
			name:     "LogAttachTemplate",
			constant: LogAttachTemplate,
			contains: "Attaching volume",
		},
		{
			name:     "LogAttachedTemplate",
			constant: LogAttachedTemplate,
			contains: "Successfully attached",
		},
		{
			name:     "LogReadTemplate",
			constant: LogReadTemplate,
			contains: "Reading volume attachment",
		},
		{
			name:     "LogReadSuccess",
			constant: LogReadSuccess,
			contains: "Successfully read volume attachment",
		},
		{
			name:     "LogNodeNotFound",
			constant: LogNodeNotFound,
			contains: "Node",
		},
		{
			name:     "LogVolumeNotFound",
			constant: LogVolumeNotFound,
			contains: "Volume",
		},
		{
			name:     "LogNotAttached",
			constant: LogNotAttached,
			contains: "not attached",
		},
		{
			name:     "LogDetachTemplate",
			constant: LogDetachTemplate,
			contains: "Detaching volume",
		},
		{
			name:     "LogDetachedTemplate",
			constant: LogDetachedTemplate,
			contains: "Successfully detached",
		},
		{
			name:     "LogWaitAttached",
			constant: LogWaitAttached,
			contains: "successfully attached",
		},
		{
			name:     "LogWaitDetached",
			constant: LogWaitDetached,
			contains: "successfully detached",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.constant, tt.contains) {
				t.Errorf("constant %s should contain %q, got: %q", tt.name, tt.contains, tt.constant)
			}
		})
	}
}

// Test constants are used correctly
func TestConstantsUsage(t *testing.T) {
	// Verify constants are accessible and have expected values
	tests := []struct {
		name     string
		constant interface{}
		expected interface{}
	}{
		{
			name:     "ImportIDFormatShortDescription",
			constant: ImportIDFormatShortDescription,
			expected: "node_id/volume_id",
		},
		{
			name:     "ImportIDFormatFullDescription",
			constant: ImportIDFormatFullDescription,
			expected: "project_id/region/node_id/volume_id",
		},
		{
			name:     "ImportIDPartsShortCount",
			constant: ImportIDPartsShortCount,
			expected: 2,
		},
		{
			name:     "ImportIDPartsFullCount",
			constant: ImportIDPartsFullCount,
			expected: 4,
		},
		{
			name:     "ResourceName",
			constant: ResourceName,
			expected: "Volume Attachment",
		},
		{
			name:     "VolumeAttachmentImportDelimiter",
			constant: tfconstants.VolumeAttachmentImportDelimiter,
			expected: "/",
		},
		{
			name:     "VolumeAttachmentVMDetailKeyVMID",
			constant: tfconstants.VolumeAttachmentVMDetailKeyVMID,
			expected: "vm_id",
		},
		{
			name:     "VolumeAttachmentVMDetailKeyVMName",
			constant: tfconstants.VolumeAttachmentVMDetailKeyVMName,
			expected: "vm_name",
		},
		{
			name:     "VolumeAttachmentVMIDNullValue",
			constant: tfconstants.VolumeAttachmentVMIDNullValue,
			expected: "null",
		},
		{
			name:     "BlockStorageActionAttach",
			constant: goe2econstants.BlockStorageActionAttach,
			expected: "attach",
		},
		{
			name:     "BlockStorageActionDetach",
			constant: goe2econstants.BlockStorageActionDetach,
			expected: "detach",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("constant %s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}
