package image

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// State Upgrade V0 → V1 Tests
// ============================================================================

func TestResourceImageStateUpgradeV0toV1_V0SchemaDefinition(t *testing.T) {
	v0Resource := resourceImageResourceV0()
	require.NotNil(t, v0Resource)
	assert.NotNil(t, v0Resource.Schema)

	// Verify V0 schema has expected fields
	schema := v0Resource.Schema
	assert.Contains(t, schema, tfconstants.AttrRegion)
	assert.Contains(t, schema, tfconstants.AttrLocation)
	assert.Contains(t, schema, tfconstants.AttrProjectID)
	assert.Contains(t, schema, tfconstants.AttrNodeID)
	assert.Contains(t, schema, tfconstants.AttrName)
	assert.Contains(t, schema, tfconstants.AttrTemplateID)
	assert.Contains(t, schema, "image_state")
	assert.Contains(t, schema, "image_type")
	assert.Contains(t, schema, "os_distribution")
	assert.Contains(t, schema, "distro")
	assert.Contains(t, schema, tfconstants.AttrCreatedAt)
}

func TestResourceImageStateUpgradeV0toV1_AddsComputedFields(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrNodeID: "node-123",
		tfconstants.AttrName:   "test-image",
	}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)
	require.NotNil(t, upgraded)

	// Verify new computed fields are added
	assert.Contains(t, upgraded, "state")
	assert.Contains(t, upgraded, "image_size")
	assert.Contains(t, upgraded, "cloning_ops")
	assert.Contains(t, upgraded, "running_vms")
	assert.Contains(t, upgraded, "is_windows")
	assert.Contains(t, upgraded, "sku_type")
	assert.Contains(t, upgraded, "vm_info")
	assert.Contains(t, upgraded, "tags")
}

func TestResourceImageStateUpgradeV0toV1_StateFieldWithImageState(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrNodeID: "node-123",
		"image_state":          goe2econstants.ImageStatusReady,
	}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	// State should be normalized from image_state
	assert.Equal(t, goe2econstants.ImageStateReady, upgraded["state"])
}

func TestResourceImageStateUpgradeV0toV1_StateFieldWithoutImageState(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrNodeID: "node-123",
	}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	// State should be empty string if image_state not present
	assert.Equal(t, DefaultEmptyString, upgraded["state"])
}

func TestResourceImageStateUpgradeV0toV1_DefaultValues(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrNodeID: "node-123",
	}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	// Verify default values
	assert.Equal(t, DefaultEmptyString, upgraded["image_size"])
	assert.Equal(t, DefaultEmptyString, upgraded["cloning_ops"])
	assert.Equal(t, DefaultEmptyString, upgraded["running_vms"])
	assert.Equal(t, DefaultIsWindows, upgraded["is_windows"])
	assert.Equal(t, DefaultEmptyString, upgraded["sku_type"])
	assert.Equal(t, DefaultVMInfo, upgraded["vm_info"])
	assert.Equal(t, DefaultTags, upgraded["tags"])
}

func TestResourceImageStateUpgradeV0toV1_PreservesExistingFields(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrNodeID:    "node-123",
		tfconstants.AttrName:      "test-image",
		"image_state":             goe2econstants.ImageStatusReady,
		"image_type":              "snapshot",
		"os_distribution":         "Ubuntu",
		"distro":                  "ubuntu",
		tfconstants.AttrCreatedAt: "2024-01-01T00:00:00Z",
	}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	// Verify existing fields are preserved
	assert.Equal(t, "node-123", upgraded[tfconstants.AttrNodeID])
	assert.Equal(t, "test-image", upgraded[tfconstants.AttrName])
	assert.Equal(t, goe2econstants.ImageStatusReady, upgraded["image_state"])
	assert.Equal(t, "snapshot", upgraded["image_type"])
	assert.Equal(t, "Ubuntu", upgraded["os_distribution"])
	assert.Equal(t, "ubuntu", upgraded["distro"])
	assert.Equal(t, "2024-01-01T00:00:00Z", upgraded[tfconstants.AttrCreatedAt])
}

func TestResourceImageStateUpgradeV0toV1_PreservesExistingComputedFields(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrNodeID: "node-123",
		"state":                "ready",
		"image_size":           "100 GB",
		"cloning_ops":          "1",
		"running_vms":          "2",
		"is_windows":           true,
		"sku_type":             "standard",
		"vm_info":              []interface{}{map[string]interface{}{"vm_id": 123}},
		"tags":                 map[string]interface{}{"key": "value"},
	}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	// Verify existing computed fields are preserved (not overwritten)
	assert.Equal(t, "ready", upgraded["state"])
	assert.Equal(t, "100 GB", upgraded["image_size"])
	assert.Equal(t, "1", upgraded["cloning_ops"])
	assert.Equal(t, "2", upgraded["running_vms"])
	assert.Equal(t, true, upgraded["is_windows"])
	assert.Equal(t, "standard", upgraded["sku_type"])
	assert.Equal(t, []interface{}{map[string]interface{}{"vm_id": 123}}, upgraded["vm_info"])
	assert.Equal(t, map[string]interface{}{"key": "value"}, upgraded["tags"])
}

func TestResourceImageStateUpgradeV0toV1_EmptyState(t *testing.T) {
	rawState := map[string]interface{}{}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	// Should add all computed fields with defaults
	assert.Contains(t, upgraded, "state")
	assert.Contains(t, upgraded, "image_size")
	assert.Contains(t, upgraded, "cloning_ops")
	assert.Contains(t, upgraded, "running_vms")
	assert.Contains(t, upgraded, "is_windows")
	assert.Contains(t, upgraded, "sku_type")
	assert.Contains(t, upgraded, "vm_info")
	assert.Contains(t, upgraded, "tags")
}

func TestResourceImageStateUpgradeV0toV1_AllFieldsPresent(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrRegion:     "Mumbai",
		tfconstants.AttrLocation:   "Mumbai",
		tfconstants.AttrProjectID:  "proj-123",
		tfconstants.AttrNodeID:     "node-123",
		tfconstants.AttrName:       "test-image",
		tfconstants.AttrTemplateID: 456,
		"image_state":              goe2econstants.ImageStatusReady,
		"image_type":               "snapshot",
		"os_distribution":          "Ubuntu",
		"distro":                   "ubuntu",
		tfconstants.AttrCreatedAt:  "2024-01-01T00:00:00Z",
		"state":                    "ready",
		"image_size":               "100 GB",
		"cloning_ops":              "0",
		"running_vms":              "1",
		"is_windows":               false,
		"sku_type":                 "standard",
		"vm_info":                  []interface{}{},
		"tags":                     map[string]interface{}{},
	}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	// Verify all fields are present
	assert.Equal(t, "Mumbai", upgraded[tfconstants.AttrRegion])
	assert.Equal(t, "Mumbai", upgraded[tfconstants.AttrLocation])
	assert.Equal(t, "proj-123", upgraded[tfconstants.AttrProjectID])
	assert.Equal(t, "node-123", upgraded[tfconstants.AttrNodeID])
	assert.Equal(t, "test-image", upgraded[tfconstants.AttrName])
	assert.Equal(t, 456, upgraded[tfconstants.AttrTemplateID])
	assert.Equal(t, goe2econstants.ImageStatusReady, upgraded["image_state"])
	assert.Equal(t, "snapshot", upgraded["image_type"])
	assert.Equal(t, "Ubuntu", upgraded["os_distribution"])
	assert.Equal(t, "ubuntu", upgraded["distro"])
	assert.Equal(t, "2024-01-01T00:00:00Z", upgraded[tfconstants.AttrCreatedAt])
	assert.Equal(t, "ready", upgraded["state"])
	assert.Equal(t, "100 GB", upgraded["image_size"])
	assert.Equal(t, "0", upgraded["cloning_ops"])
	assert.Equal(t, "1", upgraded["running_vms"])
	assert.Equal(t, false, upgraded["is_windows"])
	assert.Equal(t, "standard", upgraded["sku_type"])
	assert.Equal(t, []interface{}{}, upgraded["vm_info"])
	assert.Equal(t, map[string]interface{}{}, upgraded["tags"])
}

func TestResourceImageStateUpgradeV0toV1_StateNormalization(t *testing.T) {
	tests := []struct {
		name          string
		imageState    string
		expectedState string
	}{
		{
			name:          "Creating normalized",
			imageState:    goe2econstants.ImageStatusCreating,
			expectedState: goe2econstants.ImageStateCreating,
		},
		{
			name:          "Ready normalized",
			imageState:    goe2econstants.ImageStatusReady,
			expectedState: goe2econstants.ImageStateReady,
		},
		{
			name:          "Error normalized",
			imageState:    goe2econstants.ImageStatusError,
			expectedState: goe2econstants.ImageStateError,
		},
		{
			name:          "Deleted normalized",
			imageState:    goe2econstants.ImageStatusDeleted,
			expectedState: goe2econstants.ImageStateDeleted,
		},
		{
			name:          "Lowercase pass through",
			imageState:    "ready",
			expectedState: "ready",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawState := map[string]interface{}{
				tfconstants.AttrNodeID: "node-123",
				"image_state":          tt.imageState,
			}

			upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedState, upgraded["state"])
		})
	}
}

func TestResourceImageStateUpgradeV0toV1_NoErrors(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrNodeID: "node-123",
		tfconstants.AttrName:   "test-image",
	}

	_, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	assert.NoError(t, err, "State upgrade should not return errors for valid V0 state")
}

func TestResourceImageStateUpgradeV0toV1_MinimalValidState(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrNodeID: "node-123",
		tfconstants.AttrName:   "test-image",
	}

	upgraded, err := resourceImageStateUpgradeV0toV1(context.Background(), rawState, nil)
	require.NoError(t, err)

	// Verify all required fields are present
	assert.Contains(t, upgraded, tfconstants.AttrNodeID)
	assert.Contains(t, upgraded, tfconstants.AttrName)

	// Verify computed fields are added
	assert.Contains(t, upgraded, "state")
	assert.Contains(t, upgraded, "image_size")
	assert.Contains(t, upgraded, "cloning_ops")
	assert.Contains(t, upgraded, "running_vms")
	assert.Contains(t, upgraded, "is_windows")
	assert.Contains(t, upgraded, "sku_type")
	assert.Contains(t, upgraded, "vm_info")
	assert.Contains(t, upgraded, "tags")
}
