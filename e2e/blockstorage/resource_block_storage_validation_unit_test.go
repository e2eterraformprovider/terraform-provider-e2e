package blockstorage

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestValidateBlockStorageSize_AllValidSizes(t *testing.T) {
	for _, size := range validBlockStorageSizes {
		warns, errs := validateBlockStorageSize(size, tfconstants.AttrSize)
		require.Empty(t, warns)
		require.Empty(t, errs, "size %.0f should be valid", size)
	}
}

func TestValidateBlockStorageSize_InvalidSizes_ErrorMessageIncludesValidList(t *testing.T) {
	invalidSizes := []float64{100, 300, 750, 1500, 10000}

	for _, size := range invalidSizes {
		_, errs := validateBlockStorageSize(size, tfconstants.AttrSize)
		require.NotEmpty(t, errs, "size %.0f should be invalid", size)
		require.Contains(t, errs[0].Error(), validBlockStorageSizesString)
	}
}

func TestValidateBlockStorageSize_EdgeCases(t *testing.T) {
	edgeCases := []float64{0, -1, -250, 999999}
	for _, size := range edgeCases {
		_, errs := validateBlockStorageSize(size, tfconstants.AttrSize)
		require.NotEmpty(t, errs, "size %.0f should be invalid", size)
	}
}

func TestCustomImportBlockStorage_ValidFormat(t *testing.T) {
	r := ResourceBlockStorage()

	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
	d.SetId("project-123/Mumbai/volume-456")

	out, err := customImportBlockStorage(context.Background(), d, nil)
	require.NoError(t, err)
	require.Len(t, out, 1)

	require.Equal(t, "volume-456", out[0].Id())
	require.Equal(t, "project-123", out[0].Get(tfconstants.AttrProjectID))
	require.Equal(t, "Mumbai", out[0].Get(tfconstants.AttrRegion))
}

func TestCustomImportBlockStorage_InvalidFormat(t *testing.T) {
	r := ResourceBlockStorage()

	tests := []string{
		"project-123/Mumbai",
		"project-123/Mumbai/volume-456/extra",
		"",
	}

	for _, id := range tests {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
		d.SetId(id)
		_, err := customImportBlockStorage(context.Background(), d, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid import ID format")
	}
}

func TestIsBlockStorageSizeUpgrade(t *testing.T) {
	require.True(t, isBlockStorageSizeUpgrade(250, 500))
	require.False(t, isBlockStorageSizeUpgrade(500, 250))
	require.False(t, isBlockStorageSizeUpgrade(250, 250))
}
