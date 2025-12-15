package objectstore

import (
	"context"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Test: resourceObjectStoreStateUpgradeV0toV1
// ============================================================================

func TestResourceObjectStoreStateUpgradeV0toV1_Basic(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrName:      "test-bucket",
		tfconstants.AttrStatus:    "Active",
		tfconstants.AttrCreatedAt: "2024-01-01T00:00:00Z",
	}

	ctx := context.Background()
	meta := &config.Config{}

	result, err := resourceObjectStoreStateUpgradeV0toV1(ctx, rawState, meta)
	require.NoError(t, err)

	// Verify all V0 fields preserved
	assert.Equal(t, "test-bucket", result[tfconstants.AttrName])
	assert.Equal(t, "Active", result[tfconstants.AttrStatus])
	assert.Equal(t, "2024-01-01T00:00:00Z", result[tfconstants.AttrCreatedAt])

	// Verify V3 fields initialized
	assert.Equal(t, false, result["versioning_enabled"])
	assert.Equal(t, false, result["encryption_enabled"])
	assert.Equal(t, false, result["lock_enabled"])
	assert.Equal(t, false, result["public_access_enabled"])
	assert.NotNil(t, result["tags"])

	// Verify computed fields initialized
	assert.Equal(t, "", result["updated_at"])
	assert.Equal(t, 0, result["created_by"])
	assert.Equal(t, "", result["bucket_size"])
	assert.Equal(t, false, result["is_cdn_attached"])
	assert.Equal(t, false, result["is_encryption_enabled"])
	assert.Equal(t, false, result["is_lock_enabled"])
	assert.Equal(t, false, result["is_public_access_enabled"])
}

func TestResourceObjectStoreStateUpgradeV0toV1_AllFields(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrName:                         "test-bucket",
		tfconstants.AttrStatus:                       "Active",
		tfconstants.AttrCreatedAt:                    "2024-01-01T00:00:00Z",
		tfconstants.AttrVersioningStatus:             goe2econstants.ObjectStorageVersioningStatusEnabled,
		tfconstants.AttrLifecycleConfigurationStatus: "Enabled",
		"enabling_versioning":                        true,
	}

	ctx := context.Background()
	meta := &config.Config{}

	result, err := resourceObjectStoreStateUpgradeV0toV1(ctx, rawState, meta)
	require.NoError(t, err)

	// Verify all V0 fields preserved
	assert.Equal(t, "test-bucket", result[tfconstants.AttrName])
	assert.Equal(t, "Active", result[tfconstants.AttrStatus])
	assert.Equal(t, goe2econstants.ObjectStorageVersioningStatusEnabled, result[tfconstants.AttrVersioningStatus])
	assert.Equal(t, "Enabled", result[tfconstants.AttrLifecycleConfigurationStatus])

	// Verify versioning_enabled set from enabling_versioning
	assert.Equal(t, true, result["versioning_enabled"])
	assert.Equal(t, true, result["enabling_versioning"]) // Preserved for backward compatibility

	// Verify V3 fields initialized
	assert.Equal(t, false, result["encryption_enabled"])
	assert.Equal(t, false, result["lock_enabled"])
	assert.Equal(t, false, result["public_access_enabled"])
	assert.NotNil(t, result["tags"])
}

func TestResourceObjectStoreStateUpgradeV0toV1_WithEnablingVersioningTrue(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrName:  "test-bucket",
		"enabling_versioning": true,
	}

	ctx := context.Background()
	meta := &config.Config{}

	result, err := resourceObjectStoreStateUpgradeV0toV1(ctx, rawState, meta)
	require.NoError(t, err)

	// Verify versioning_enabled set from enabling_versioning
	assert.Equal(t, true, result["versioning_enabled"])
	assert.Equal(t, true, result["enabling_versioning"]) // Preserved for backward compatibility
}

func TestResourceObjectStoreStateUpgradeV0toV1_WithEnablingVersioningFalse(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrName:  "test-bucket",
		"enabling_versioning": false,
	}

	ctx := context.Background()
	meta := &config.Config{}

	result, err := resourceObjectStoreStateUpgradeV0toV1(ctx, rawState, meta)
	require.NoError(t, err)

	// Verify versioning_enabled set to false
	assert.Equal(t, false, result["versioning_enabled"])
}

func TestResourceObjectStoreStateUpgradeV0toV1_PreservesDeprecated(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrName:  "test-bucket",
		"enabling_versioning": true,
	}

	ctx := context.Background()
	meta := &config.Config{}

	result, err := resourceObjectStoreStateUpgradeV0toV1(ctx, rawState, meta)
	require.NoError(t, err)

	// Verify deprecated field preserved
	assert.Equal(t, true, result["enabling_versioning"])
	// Verify new field also set
	assert.Equal(t, true, result["versioning_enabled"])

	// Verify no forced recreation (resource ID would be preserved in actual upgrade)
	// This is implicit - if upgrade succeeds, no recreation is triggered
	assert.NotNil(t, result)
}

func TestResourceObjectStoreStateUpgradeV0toV1_NoForcedRecreation(t *testing.T) {
	rawState := map[string]interface{}{
		tfconstants.AttrName:  "test-bucket",
		"enabling_versioning": false,
	}

	ctx := context.Background()
	meta := &config.Config{}

	result, err := resourceObjectStoreStateUpgradeV0toV1(ctx, rawState, meta)
	require.NoError(t, err)

	// Verify resource ID would be preserved (implicit - upgrade succeeds without error)
	// All fields are preserved or initialized, no recreation needed
	assert.Equal(t, "test-bucket", result[tfconstants.AttrName])
}
