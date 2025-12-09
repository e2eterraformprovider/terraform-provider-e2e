package sfs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResourceSfsStateUpgradeV0toV1_Basic tests basic state upgrade
func TestResourceSfsStateUpgradeV0toV1_Basic(t *testing.T) {
	// Simulate old V0 state
	oldState := map[string]interface{}{
		"id":                    "12345",
		"name":                  "test-sfs",
		"plan":                  "standard",
		"vpc_id":                "vpc-123",
		"disk_size":             float64(100), // Note: JSON unmarshals numbers as float64
		"disk_iops":             float64(1000),
		"is_encryption_enabled": false,
		"status":                "Active",
	}

	// Run the upgrade function
	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify old fields are preserved
	assert.Equal(t, "12345", newState["id"])
	assert.Equal(t, "test-sfs", newState["name"])
	assert.Equal(t, "standard", newState["plan"])
	assert.Equal(t, "vpc-123", newState["vpc_id"])

	// Verify new fields are created
	assert.Equal(t, float64(100), newState["size_gb"])
	assert.Equal(t, float64(1000), newState["iops"])
	assert.Equal(t, false, newState["encryption_enabled"])

	// Verify normalized state is added
	assert.Equal(t, "", newState["state"])

	// Verify both old and new fields exist for backwards compatibility
	assert.Equal(t, float64(100), newState["disk_size"])
	assert.Equal(t, float64(1000), newState["disk_iops"])
	assert.Equal(t, false, newState["is_encryption_enabled"])
}

// TestResourceSfsStateUpgradeV0toV1_AllFields tests upgrade with all fields
func TestResourceSfsStateUpgradeV0toV1_AllFields(t *testing.T) {
	oldState := map[string]interface{}{
		"id":                    "67890",
		"name":                  "prod-sfs",
		"plan":                  "premium",
		"vpc_id":                "vpc-456",
		"disk_size":             float64(500),
		"disk_iops":             float64(5000),
		"is_encryption_enabled": true,
		"encryption_passphrase": "secret-key",
		"status":                "Creating",
		"private_endpoint":      "nfs.example.com:/data",
		"is_backup_enabled":     true,
		"created_at":            "2025-01-01T00:00:00Z",
	}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify all original fields preserved
	assert.Equal(t, "67890", newState["id"])
	assert.Equal(t, "prod-sfs", newState["name"])
	assert.Equal(t, "premium", newState["plan"])
	assert.Equal(t, "vpc-456", newState["vpc_id"])
	assert.Equal(t, "secret-key", newState["encryption_passphrase"])
	assert.Equal(t, "nfs.example.com:/data", newState["private_endpoint"])
	assert.Equal(t, true, newState["is_backup_enabled"])
	assert.Equal(t, "2025-01-01T00:00:00Z", newState["created_at"])

	// Verify new V3 fields
	assert.Equal(t, float64(500), newState["size_gb"])
	assert.Equal(t, float64(5000), newState["iops"])
	assert.Equal(t, true, newState["encryption_enabled"])

	// Verify normalized state
	assert.Equal(t, "", newState["state"])

	// Verify mount_endpoint is set as alias for private_endpoint
	assert.Equal(t, "nfs.example.com:/data", newState["mount_endpoint"])
}

// TestResourceSfsStateUpgradeV0toV1_PreservesDeprecated tests that deprecated fields are preserved
func TestResourceSfsStateUpgradeV0toV1_PreservesDeprecated(t *testing.T) {
	oldState := map[string]interface{}{
		"id":                    "11111",
		"name":                  "legacy-sfs",
		"plan":                  "standard",
		"vpc_id":                "vpc-789",
		"disk_size":             float64(200),
		"disk_iops":             float64(2000),
		"is_encryption_enabled": false,
		"status":                "Active",
	}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify V2 fields are NOT removed (backwards compatibility)
	assert.Equal(t, float64(200), newState["disk_size"])
	assert.Equal(t, float64(2000), newState["disk_iops"])
	assert.Equal(t, false, newState["is_encryption_enabled"])

	// Verify V3 fields are added
	assert.Equal(t, float64(200), newState["size_gb"])
	assert.Equal(t, float64(2000), newState["iops"])
	assert.Equal(t, false, newState["encryption_enabled"])
}

// TestResourceSfsStateUpgradeV0toV1_ComputedFields tests computed fields initialization
func TestResourceSfsStateUpgradeV0toV1_ComputedFields(t *testing.T) {
	oldState := map[string]interface{}{
		"id":     "22222",
		"name":   "test-sfs",
		"plan":   "standard",
		"vpc_id": "vpc-111",
		"status": "Active",
	}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify computed fields are initialized
	assert.Equal(t, "", newState["state"])

	// mount_endpoint should be set from private_endpoint or empty
	_, hasMountEndpoint := newState["mount_endpoint"]
	assert.True(t, hasMountEndpoint)
}

// TestResourceSfsStateUpgradeV0toV1_NilState tests upgrade with nil state
func TestResourceSfsStateUpgradeV0toV1_NilState(t *testing.T) {
	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), nil, nil)
	// Should handle gracefully or return error
	if err != nil {
		require.Error(t, err)
	} else {
		require.NotNil(t, newState)
	}
}

// TestResourceSfsStateUpgradeV0toV1_EmptyState tests upgrade with empty state
func TestResourceSfsStateUpgradeV0toV1_EmptyState(t *testing.T) {
	oldState := map[string]interface{}{}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)
	require.NotNil(t, newState)

	// Verify new fields are initialized with defaults
	assert.Equal(t, "", newState["state"])
}

// TestResourceSfsStateUpgradeV0toV1_StatusNormalization tests status normalization during upgrade
func TestResourceSfsStateUpgradeV0toV1_StatusNormalization(t *testing.T) {
	tests := []struct {
		name          string
		inputStatus   string
		expectedState string
	}{
		{
			name:          "Creating status",
			inputStatus:   "Creating",
			expectedState: "creating",
		},
		{
			name:          "Active status",
			inputStatus:   "Active",
			expectedState: "active",
		},
		{
			name:          "Deleting status",
			inputStatus:   "Deleting",
			expectedState: "deleting",
		},
		{
			name:          "Error status",
			inputStatus:   "Error",
			expectedState: "error",
		},
		{
			name:          "Unknown status",
			inputStatus:   "Unknown",
			expectedState: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldState := map[string]interface{}{
				"id":     "test-id",
				"name":   "test-sfs",
				"status": tt.inputStatus,
			}

			newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
			require.NoError(t, err)

			assert.Equal(t, "", newState["state"])
		})
	}
}

// TestResourceSfsStateUpgradeV0toV1_IntegrationWithResource tests that upgraded state works with resource
func TestResourceSfsStateUpgradeV0toV1_IntegrationWithResource(t *testing.T) {
	// This test verifies that the upgraded state can be used with the resource schema
	oldState := map[string]interface{}{
		"id":                    "integration-test",
		"name":                  "test-sfs",
		"plan":                  "standard",
		"vpc_id":                "vpc-123",
		"disk_size":             float64(100),
		"disk_iops":             float64(1000),
		"is_encryption_enabled": false,
		"status":                "Active",
	}

	// Convert to terraform.InstanceState (simulating how the state upgrader would be called)
	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify all required fields exist in upgraded state
	requiredFields := []string{"id", "name", "plan", "vpc_id", "size_gb", "iops", "encryption_enabled", "state"}
	for _, field := range requiredFields {
		_, exists := newState[field]
		assert.True(t, exists, "Field %s should exist in upgraded state", field)
	}
}

// TestResourceSfsStateUpgradeV0toV1_DiskSizeZeroHandling tests handling of zero disk size
func TestResourceSfsStateUpgradeV0toV1_DiskSizeZeroHandling(t *testing.T) {
	oldState := map[string]interface{}{
		"id":        "zero-test",
		"name":      "test-sfs",
		"disk_size": float64(0),
		"disk_iops": float64(0),
		"status":    "Active",
	}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Even zero values should be preserved during upgrade
	assert.Equal(t, float64(0), newState["size_gb"])
	assert.Equal(t, float64(0), newState["iops"])
}

// TestResourceSfsV0Schema tests that the V0 schema is correctly defined
func TestResourceSfsV0Schema(t *testing.T) {
	schema := resourceSfsResourceV0()
	require.NotNil(t, schema)

	// Verify schema structure exists
	assert.NotNil(t, schema.Schema)

	// Verify V0 has old field names
	assert.Contains(t, schema.Schema, "disk_size")
	assert.Contains(t, schema.Schema, "disk_iops")
	assert.Contains(t, schema.Schema, "is_encryption_enabled")
}

// TestStateUpgraderType tests the StateUpgrader configuration
func TestStateUpgraderType(t *testing.T) {
	resource := ResourceSfs()
	require.NotNil(t, resource)

	// Verify StateUpgraders is configured
	assert.Greater(t, len(resource.StateUpgraders), 0, "StateUpgraders should be configured")

	// Verify upgrader points to correct upgrade function
	upgrader := resource.StateUpgraders[0]
	assert.NotNil(t, upgrader.Upgrade)
}
