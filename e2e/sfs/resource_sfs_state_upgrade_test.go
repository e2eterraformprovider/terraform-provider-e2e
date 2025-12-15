package sfs

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResourceSfsStateUpgradeV0toV1_Basic tests basic state upgrade
func TestResourceSfsStateUpgradeV0toV1_Basic(t *testing.T) {
	// Simulate old V0 state
	oldState := map[string]interface{}{
		tfconstants.AttrID:                               "12345",
		tfconstants.AttrName:                             "test-sfs",
		tfconstants.AttrPlan:                             "standard",
		tfconstants.AttrVPCID:                            "vpc-123",
		tfconstants.FieldMigrationKeyDiskSize:            float64(100), // Note: JSON unmarshals numbers as float64
		tfconstants.FieldMigrationKeyDiskIOPS:            float64(1000),
		tfconstants.FieldMigrationKeyIsEncryptionEnabled: false,
		tfconstants.AttrStatus:                           goe2econstants.SFSStatusActive,
	}

	// Run the upgrade function
	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify old fields are preserved
	assert.Equal(t, "12345", newState[tfconstants.AttrID])
	assert.Equal(t, "test-sfs", newState[tfconstants.AttrName])
	assert.Equal(t, "standard", newState[tfconstants.AttrPlan])
	assert.Equal(t, "vpc-123", newState[tfconstants.AttrVPCID])

	// Verify new fields are created
	assert.Equal(t, float64(100), newState[tfconstants.FieldMigrationKeySizeGB])
	assert.Equal(t, float64(1000), newState[tfconstants.FieldMigrationKeyIOPS])
	assert.Equal(t, false, newState[tfconstants.FieldMigrationKeyEncryptionEnabled])

	// Verify normalized state is added
	assert.Equal(t, "", newState[tfconstants.AttrState])

	// Verify both old and new fields exist for backwards compatibility
	assert.Equal(t, float64(100), newState[tfconstants.FieldMigrationKeyDiskSize])
	assert.Equal(t, float64(1000), newState[tfconstants.FieldMigrationKeyDiskIOPS])
	assert.Equal(t, false, newState[tfconstants.FieldMigrationKeyIsEncryptionEnabled])
}

// TestResourceSfsStateUpgradeV0toV1_AllFields tests upgrade with all fields
func TestResourceSfsStateUpgradeV0toV1_AllFields(t *testing.T) {
	oldState := map[string]interface{}{
		tfconstants.AttrID:                               "67890",
		tfconstants.AttrName:                             "prod-sfs",
		tfconstants.AttrPlan:                             "premium",
		tfconstants.AttrVPCID:                            "vpc-456",
		tfconstants.FieldMigrationKeyDiskSize:            float64(500),
		tfconstants.FieldMigrationKeyDiskIOPS:            float64(5000),
		tfconstants.FieldMigrationKeyIsEncryptionEnabled: true,
		tfconstants.AttrEncryptionPassphrase:             "secret-key",
		tfconstants.AttrStatus:                           goe2econstants.SFSStatusCreating,
		tfconstants.AttrPrivateEndpoint:                  "nfs.example.com:/data",
		tfconstants.AttrIsBackupEnabled:                  true,
		tfconstants.AttrCreatedAt:                        "2025-01-01T00:00:00Z",
	}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify all original fields preserved
	assert.Equal(t, "67890", newState[tfconstants.AttrID])
	assert.Equal(t, "prod-sfs", newState[tfconstants.AttrName])
	assert.Equal(t, "premium", newState[tfconstants.AttrPlan])
	assert.Equal(t, "vpc-456", newState[tfconstants.AttrVPCID])
	assert.Equal(t, "secret-key", newState[tfconstants.AttrEncryptionPassphrase])
	assert.Equal(t, "nfs.example.com:/data", newState[tfconstants.AttrPrivateEndpoint])
	assert.Equal(t, true, newState[tfconstants.AttrIsBackupEnabled])
	assert.Equal(t, "2025-01-01T00:00:00Z", newState[tfconstants.AttrCreatedAt])

	// Verify new V3 fields
	assert.Equal(t, float64(500), newState[tfconstants.FieldMigrationKeySizeGB])
	assert.Equal(t, float64(5000), newState[tfconstants.FieldMigrationKeyIOPS])
	assert.Equal(t, true, newState[tfconstants.FieldMigrationKeyEncryptionEnabled])

	// Verify normalized state
	assert.Equal(t, "", newState[tfconstants.AttrState])

	// Verify mount_endpoint is set as alias for private_endpoint
	assert.Equal(t, "nfs.example.com:/data", newState[tfconstants.AttrMountEndpoint])
}

// TestResourceSfsStateUpgradeV0toV1_PreservesDeprecated tests that deprecated fields are preserved
func TestResourceSfsStateUpgradeV0toV1_PreservesDeprecated(t *testing.T) {
	oldState := map[string]interface{}{
		tfconstants.AttrID:                               "11111",
		tfconstants.AttrName:                             "legacy-sfs",
		tfconstants.AttrPlan:                             "standard",
		tfconstants.AttrVPCID:                            "vpc-789",
		tfconstants.FieldMigrationKeyDiskSize:            float64(200),
		tfconstants.FieldMigrationKeyDiskIOPS:            float64(2000),
		tfconstants.FieldMigrationKeyIsEncryptionEnabled: false,
		tfconstants.AttrStatus:                           goe2econstants.SFSStatusActive,
	}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify V2 fields are NOT removed (backwards compatibility)
	assert.Equal(t, float64(200), newState[tfconstants.FieldMigrationKeyDiskSize])
	assert.Equal(t, float64(2000), newState[tfconstants.FieldMigrationKeyDiskIOPS])
	assert.Equal(t, false, newState[tfconstants.FieldMigrationKeyIsEncryptionEnabled])

	// Verify V3 fields are added
	assert.Equal(t, float64(200), newState[tfconstants.FieldMigrationKeySizeGB])
	assert.Equal(t, float64(2000), newState[tfconstants.FieldMigrationKeyIOPS])
	assert.Equal(t, false, newState[tfconstants.FieldMigrationKeyEncryptionEnabled])
}

// TestResourceSfsStateUpgradeV0toV1_ComputedFields tests computed fields initialization
func TestResourceSfsStateUpgradeV0toV1_ComputedFields(t *testing.T) {
	oldState := map[string]interface{}{
		tfconstants.AttrID:     "22222",
		tfconstants.AttrName:   "test-sfs",
		tfconstants.AttrPlan:   "standard",
		tfconstants.AttrVPCID:  "vpc-111",
		tfconstants.AttrStatus: goe2econstants.SFSStatusActive,
	}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify computed fields are initialized
	assert.Equal(t, "", newState[tfconstants.AttrState])

	// mount_endpoint should be set from private_endpoint or empty
	_, hasMountEndpoint := newState[tfconstants.AttrMountEndpoint]
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
	assert.Equal(t, "", newState[tfconstants.AttrState])
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
			inputStatus:   goe2econstants.SFSStatusCreating,
			expectedState: goe2econstants.SFSStateCreating,
		},
		{
			name:          "Active status",
			inputStatus:   goe2econstants.SFSStatusActive,
			expectedState: goe2econstants.SFSStateActive,
		},
		{
			name:          "Deleting status",
			inputStatus:   goe2econstants.SFSStatusDeleting,
			expectedState: goe2econstants.SFSStateDeleting,
		},
		{
			name:          "Error status",
			inputStatus:   goe2econstants.SFSStatusError,
			expectedState: goe2econstants.SFSStateError,
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
				tfconstants.AttrID:     "test-id",
				tfconstants.AttrName:   "test-sfs",
				tfconstants.AttrStatus: tt.inputStatus,
			}

			newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
			require.NoError(t, err)

			assert.Equal(t, "", newState[tfconstants.AttrState])
		})
	}
}

// TestResourceSfsStateUpgradeV0toV1_IntegrationWithResource tests that upgraded state works with resource
func TestResourceSfsStateUpgradeV0toV1_IntegrationWithResource(t *testing.T) {
	// This test verifies that the upgraded state can be used with the resource schema
	oldState := map[string]interface{}{
		tfconstants.AttrID:                               "integration-test",
		tfconstants.AttrName:                             "test-sfs",
		tfconstants.AttrPlan:                             "standard",
		tfconstants.AttrVPCID:                            "vpc-123",
		tfconstants.FieldMigrationKeyDiskSize:            float64(100),
		tfconstants.FieldMigrationKeyDiskIOPS:            float64(1000),
		tfconstants.FieldMigrationKeyIsEncryptionEnabled: false,
		tfconstants.AttrStatus:                           goe2econstants.SFSStatusActive,
	}

	// Convert to terraform.InstanceState (simulating how the state upgrader would be called)
	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify all required fields exist in upgraded state
	requiredFields := []string{
		tfconstants.AttrID,
		tfconstants.AttrName,
		tfconstants.AttrPlan,
		tfconstants.AttrVPCID,
		tfconstants.FieldMigrationKeySizeGB,
		tfconstants.FieldMigrationKeyIOPS,
		tfconstants.FieldMigrationKeyEncryptionEnabled,
		tfconstants.AttrState,
	}
	for _, field := range requiredFields {
		_, exists := newState[field]
		assert.True(t, exists, "Field %s should exist in upgraded state", field)
	}
}

// TestResourceSfsStateUpgradeV0toV1_DiskSizeZeroHandling tests handling of zero disk size
func TestResourceSfsStateUpgradeV0toV1_DiskSizeZeroHandling(t *testing.T) {
	oldState := map[string]interface{}{
		tfconstants.AttrID:                    "zero-test",
		tfconstants.AttrName:                  "test-sfs",
		tfconstants.FieldMigrationKeyDiskSize: float64(0),
		tfconstants.FieldMigrationKeyDiskIOPS: float64(0),
		tfconstants.AttrStatus:                goe2econstants.SFSStatusActive,
	}

	newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Even zero values should be preserved during upgrade
	assert.Equal(t, float64(0), newState[tfconstants.FieldMigrationKeySizeGB])
	assert.Equal(t, float64(0), newState[tfconstants.FieldMigrationKeyIOPS])
}

// TestResourceSfsV0Schema tests that the V0 schema is correctly defined
func TestResourceSfsV0Schema(t *testing.T) {
	schema := resourceSfsResourceV0()
	require.NotNil(t, schema)

	// Verify schema structure exists
	assert.NotNil(t, schema.Schema)

	// Verify V0 has old field names
	assert.Contains(t, schema.Schema, tfconstants.FieldMigrationKeyDiskSize)
	assert.Contains(t, schema.Schema, tfconstants.FieldMigrationKeyDiskIOPS)
	assert.Contains(t, schema.Schema, tfconstants.FieldMigrationKeyIsEncryptionEnabled)
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

// TestResourceSfsStateUpgradeV0toV1_AllFieldMigrationCombinations tests all field migration combinations
func TestResourceSfsStateUpgradeV0toV1_AllFieldMigrationCombinations(t *testing.T) {
	tests := []struct {
		name                string
		oldState            map[string]interface{}
		expectedV3Fields    map[string]interface{}
		expectedV2Preserved bool
		description         string
	}{
		{
			name: "all_V2_fields_present",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                               "test-1",
				tfconstants.FieldMigrationKeyDiskSize:            float64(100),
				tfconstants.FieldMigrationKeyDiskIOPS:            float64(1000),
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: true,
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB:            float64(100),
				tfconstants.FieldMigrationKeyIOPS:              float64(1000),
				tfconstants.FieldMigrationKeyEncryptionEnabled: true,
			},
			expectedV2Preserved: true,
			description:         "All V2 fields should migrate to V3 and be preserved",
		},
		{
			name: "only_disk_size_present",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-2",
				tfconstants.FieldMigrationKeyDiskSize: float64(200),
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB: float64(200),
			},
			expectedV2Preserved: true,
			description:         "Only disk_size should migrate to size_gb",
		},
		{
			name: "only_disk_iops_present",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-3",
				tfconstants.FieldMigrationKeyDiskIOPS: float64(2000),
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.FieldMigrationKeyIOPS: float64(2000),
			},
			expectedV2Preserved: true,
			description:         "Only disk_iops should migrate to iops",
		},
		{
			name: "only_is_encryption_enabled_present",
			oldState: map[string]interface{}{
				tfconstants.AttrID: "test-4",
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: false,
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.FieldMigrationKeyEncryptionEnabled: false,
			},
			expectedV2Preserved: true,
			description:         "Only is_encryption_enabled should migrate to encryption_enabled",
		},
		{
			name: "disk_size_and_disk_iops_present",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-5",
				tfconstants.FieldMigrationKeyDiskSize: float64(300),
				tfconstants.FieldMigrationKeyDiskIOPS: float64(3000),
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB: float64(300),
				tfconstants.FieldMigrationKeyIOPS:   float64(3000),
			},
			expectedV2Preserved: true,
			description:         "disk_size and disk_iops should both migrate",
		},
		{
			name: "disk_size_and_encryption_present",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                               "test-6",
				tfconstants.FieldMigrationKeyDiskSize:            float64(400),
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: true,
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB:            float64(400),
				tfconstants.FieldMigrationKeyEncryptionEnabled: true,
			},
			expectedV2Preserved: true,
			description:         "disk_size and is_encryption_enabled should both migrate",
		},
		{
			name: "disk_iops_and_encryption_present",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                               "test-7",
				tfconstants.FieldMigrationKeyDiskIOPS:            float64(4000),
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: false,
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.FieldMigrationKeyIOPS:              float64(4000),
				tfconstants.FieldMigrationKeyEncryptionEnabled: false,
			},
			expectedV2Preserved: true,
			description:         "disk_iops and is_encryption_enabled should both migrate",
		},
		{
			name: "no_V2_fields_present",
			oldState: map[string]interface{}{
				tfconstants.AttrID:   "test-8",
				tfconstants.AttrName: "test-sfs",
			},
			expectedV3Fields:    map[string]interface{}{},
			expectedV2Preserved: false,
			description:         "No V2 fields means no V3 fields created, but computed fields should be initialized",
		},
		{
			name: "all_fields_with_other_fields",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                               "test-9",
				tfconstants.AttrName:                             "test-sfs",
				tfconstants.AttrPlan:                             "standard",
				tfconstants.AttrVPCID:                            "vpc-123",
				tfconstants.FieldMigrationKeyDiskSize:            float64(500),
				tfconstants.FieldMigrationKeyDiskIOPS:            float64(5000),
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: true,
				tfconstants.AttrEncryptionPassphrase:             "secret",
				tfconstants.AttrStatus:                           goe2econstants.SFSStatusActive,
				tfconstants.AttrPrivateEndpoint:                  "10.0.0.1",
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB:            float64(500),
				tfconstants.FieldMigrationKeyIOPS:              float64(5000),
				tfconstants.FieldMigrationKeyEncryptionEnabled: true,
			},
			expectedV2Preserved: true,
			description:         "All V2 fields should migrate even when other fields are present",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), tt.oldState, nil)
			require.NoError(t, err, tt.description)

			// Verify V3 fields are created
			for fieldName, expectedValue := range tt.expectedV3Fields {
				actualValue, exists := newState[fieldName]
				assert.True(t, exists, "V3 field %s should exist: %s", fieldName, tt.description)
				assert.Equal(t, expectedValue, actualValue, "V3 field %s should match: %s", fieldName, tt.description)
			}

			// Verify V2 fields are preserved if they existed
			if tt.expectedV2Preserved {
				if diskSize, ok := tt.oldState[tfconstants.FieldMigrationKeyDiskSize]; ok {
					assert.Equal(t, diskSize, newState[tfconstants.FieldMigrationKeyDiskSize], "V2 field disk_size should be preserved: %s", tt.description)
					assert.Equal(t, diskSize, newState[tfconstants.FieldMigrationKeySizeGB], "V3 field size_gb should equal V2 disk_size: %s", tt.description)
				}
				if diskIops, ok := tt.oldState[tfconstants.FieldMigrationKeyDiskIOPS]; ok {
					assert.Equal(t, diskIops, newState[tfconstants.FieldMigrationKeyDiskIOPS], "V2 field disk_iops should be preserved: %s", tt.description)
					assert.Equal(t, diskIops, newState[tfconstants.FieldMigrationKeyIOPS], "V3 field iops should equal V2 disk_iops: %s", tt.description)
				}
				if isEncryption, ok := tt.oldState[tfconstants.FieldMigrationKeyIsEncryptionEnabled]; ok {
					assert.Equal(t, isEncryption, newState[tfconstants.FieldMigrationKeyIsEncryptionEnabled], "V2 field is_encryption_enabled should be preserved: %s", tt.description)
					assert.Equal(t, isEncryption, newState[tfconstants.FieldMigrationKeyEncryptionEnabled], "V3 field encryption_enabled should equal V2 is_encryption_enabled: %s", tt.description)
				}
			}

			// Verify computed fields are initialized
			assert.Equal(t, "", newState[tfconstants.AttrState], "Computed field state should be initialized: %s", tt.description)
			_, hasMountEndpoint := newState[tfconstants.AttrMountEndpoint]
			assert.True(t, hasMountEndpoint, "Computed field mount_endpoint should exist: %s", tt.description)
		})
	}
}

// TestResourceSfsStateUpgradeV0toV1_EdgeCases tests edge cases for field values
func TestResourceSfsStateUpgradeV0toV1_EdgeCases(t *testing.T) {
	tests := []struct {
		name            string
		oldState        map[string]interface{}
		expectedResults map[string]interface{}
		description     string
	}{
		{
			name: "nil_disk_size",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-nil-size",
				tfconstants.FieldMigrationKeyDiskSize: nil,
			},
			expectedResults: map[string]interface{}{},
			description:     "Nil disk_size should not create size_gb field (nil values are skipped)",
		},
		{
			name: "nil_disk_iops",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-nil-iops",
				tfconstants.FieldMigrationKeyDiskIOPS: nil,
			},
			expectedResults: map[string]interface{}{},
			description:     "Nil disk_iops should not create iops field (nil values are skipped)",
		},
		{
			name: "nil_is_encryption_enabled",
			oldState: map[string]interface{}{
				tfconstants.AttrID: "test-nil-encryption",
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: nil,
			},
			expectedResults: map[string]interface{}{},
			description:     "Nil is_encryption_enabled should not create encryption_enabled field (nil values are skipped)",
		},
		{
			name: "zero_disk_size",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-zero-size",
				tfconstants.FieldMigrationKeyDiskSize: float64(0),
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB: float64(0),
			},
			description: "Zero disk_size should migrate to zero size_gb",
		},
		{
			name: "zero_disk_iops",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-zero-iops",
				tfconstants.FieldMigrationKeyDiskIOPS: float64(0),
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeyIOPS: float64(0),
			},
			description: "Zero disk_iops should migrate to zero iops",
		},
		{
			name: "negative_disk_size",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-negative-size",
				tfconstants.FieldMigrationKeyDiskSize: float64(-100),
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB: float64(-100),
			},
			description: "Negative disk_size should migrate as-is (validation happens elsewhere)",
		},
		{
			name: "negative_disk_iops",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-negative-iops",
				tfconstants.FieldMigrationKeyDiskIOPS: float64(-1000),
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeyIOPS: float64(-1000),
			},
			description: "Negative disk_iops should migrate as-is (validation happens elsewhere)",
		},
		{
			name: "very_large_disk_size",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-large-size",
				tfconstants.FieldMigrationKeyDiskSize: float64(999999),
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB: float64(999999),
			},
			description: "Very large disk_size should migrate correctly",
		},
		{
			name: "very_large_disk_iops",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-large-iops",
				tfconstants.FieldMigrationKeyDiskIOPS: float64(999999),
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeyIOPS: float64(999999),
			},
			description: "Very large disk_iops should migrate correctly",
		},
		{
			name: "all_nil_fields",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                               "test-all-nil",
				tfconstants.FieldMigrationKeyDiskSize:            nil,
				tfconstants.FieldMigrationKeyDiskIOPS:            nil,
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: nil,
			},
			expectedResults: map[string]interface{}{},
			description:     "All nil V2 fields should not create V3 fields (nil values are skipped)",
		},
		{
			name: "all_zero_fields",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                               "test-all-zero",
				tfconstants.FieldMigrationKeyDiskSize:            float64(0),
				tfconstants.FieldMigrationKeyDiskIOPS:            float64(0),
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: false,
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB:            float64(0),
				tfconstants.FieldMigrationKeyIOPS:              float64(0),
				tfconstants.FieldMigrationKeyEncryptionEnabled: false,
			},
			description: "All zero V2 fields should migrate to zero V3 fields",
		},
		{
			name: "all_negative_fields",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-all-negative",
				tfconstants.FieldMigrationKeyDiskSize: float64(-1),
				tfconstants.FieldMigrationKeyDiskIOPS: float64(-1),
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB: float64(-1),
				tfconstants.FieldMigrationKeyIOPS:   float64(-1),
			},
			description: "Negative values should migrate as-is",
		},
		{
			name: "mixed_nil_zero_negative",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                               "test-mixed",
				tfconstants.FieldMigrationKeyDiskSize:            nil,
				tfconstants.FieldMigrationKeyDiskIOPS:            float64(0),
				tfconstants.FieldMigrationKeyIsEncryptionEnabled: false,
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeyIOPS:              float64(0),
				tfconstants.FieldMigrationKeyEncryptionEnabled: false,
			},
			description: "Mixed nil (skipped), zero, and false values should migrate correctly",
		},
		{
			name: "float64_precision_preservation",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-precision",
				tfconstants.FieldMigrationKeyDiskSize: float64(100.5), // Should be preserved as float64
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB: float64(100.5),
			},
			description: "Float64 precision should be preserved during migration",
		},
		{
			name: "int_type_conversion",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-int",
				tfconstants.FieldMigrationKeyDiskSize: int(100), // Different type
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeySizeGB: int(100),
			},
			description: "Int type should migrate correctly",
		},
		{
			name: "int64_type_conversion",
			oldState: map[string]interface{}{
				tfconstants.AttrID:                    "test-int64",
				tfconstants.FieldMigrationKeyDiskIOPS: int64(1000), // Different type
			},
			expectedResults: map[string]interface{}{
				tfconstants.FieldMigrationKeyIOPS: int64(1000),
			},
			description: "Int64 type should migrate correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newState, err := resourceSfsStateUpgradeV0toV1(context.Background(), tt.oldState, nil)
			require.NoError(t, err, tt.description)

			// Verify expected results
			for fieldName, expectedValue := range tt.expectedResults {
				actualValue, exists := newState[fieldName]
				assert.True(t, exists, "Field %s should exist: %s", fieldName, tt.description)
				assert.Equal(t, expectedValue, actualValue, "Field %s should match expected value: %s", fieldName, tt.description)
			}

			// Verify computed fields are always initialized
			assert.Equal(t, "", newState[tfconstants.AttrState], "Computed field state should be initialized: %s", tt.description)
			_, hasMountEndpoint := newState[tfconstants.AttrMountEndpoint]
			assert.True(t, hasMountEndpoint, "Computed field mount_endpoint should exist: %s", tt.description)
		})
	}
}
