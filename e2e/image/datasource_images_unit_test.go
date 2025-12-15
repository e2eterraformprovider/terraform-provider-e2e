package image

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Schema Validation Tests
// ============================================================================

// TestDataSourceImagesSchemaDefinition tests that the datasource schema definition is correct
func TestDataSourceImagesSchemaDefinition(t *testing.T) {
	datasource := DataSourceImages()
	require.NotNil(t, datasource)
	assert.NotNil(t, datasource.Schema, "Schema should not be nil")
	assert.NotEmpty(t, datasource.Schema, "Schema should not be empty")
}

// TestDataSourceImagesSchemaRequiredVsOptionalFields tests required vs optional fields
func TestDataSourceImagesSchemaRequiredVsOptionalFields(t *testing.T) {
	datasource := DataSourceImages()
	require.NotNil(t, datasource)
	schema := datasource.Schema

	// All fields should be optional (computed) for datasource
	// region, location, project_id are optional (can use provider defaults)
	// image_list is computed
	optionalFields := []string{
		tfconstants.AttrRegion,
		tfconstants.AttrLocation,
		tfconstants.AttrProjectID,
		"image_list",
	}

	for _, fieldName := range optionalFields {
		t.Run(fieldName+"_is_optional", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.False(t, fieldSchema.Required, "Field %s should not be required", fieldName)
		})
	}
}

// TestDataSourceImagesSchemaComputedFields tests that computed fields are properly marked
func TestDataSourceImagesSchemaComputedFields(t *testing.T) {
	datasource := DataSourceImages()
	require.NotNil(t, datasource)
	schema := datasource.Schema

	// Top-level computed fields
	computedFields := []string{
		tfconstants.AttrProjectID,
		"image_list",
	}

	for _, fieldName := range computedFields {
		t.Run(fieldName+"_is_computed", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Computed, "Field %s should be computed", fieldName)
		})
	}

	// Verify region and location exist (they're optional inputs)
	optionalFields := []string{
		tfconstants.AttrRegion,
		tfconstants.AttrLocation,
	}

	for _, fieldName := range optionalFields {
		t.Run(fieldName+"_exists", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.NotNil(t, fieldSchema, "Field %s schema should not be nil", fieldName)
		})
	}

	// Verify image_list is a list type
	imageListSchema := schema["image_list"]
	require.NotNil(t, imageListSchema, "image_list field should exist")
	require.NotNil(t, imageListSchema.Elem, "image_list should have Elem")
}

// ============================================================================
// flattenSavedImages() Function Tests
// ============================================================================

// TestFlattenSavedImages_EmptyList tests with empty image list
func TestFlattenSavedImages_EmptyList(t *testing.T) {
	images := []goe2e.SavedImage{}
	result := flattenSavedImages(images)

	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Len(t, result, 0)
}

// TestFlattenSavedImages_NilList tests with nil image list
func TestFlattenSavedImages_NilList(t *testing.T) {
	var images []goe2e.SavedImage
	result := flattenSavedImages(images)

	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Len(t, result, 0)
}

// TestFlattenSavedImages_SingleImage tests with single image
func TestFlattenSavedImages_SingleImage(t *testing.T) {
	images := []goe2e.SavedImage{
		{
			TemplateID:     123,
			ImageState:     goe2econstants.ImageStatusReady,
			ImageType:      "snapshot",
			OSDistribution: "Ubuntu",
			Name:           "test-image",
			ImageID:        "img-123",
			Distro:         "ubuntu",
			SKUType:        "standard",
		},
	}

	result := flattenSavedImages(images)

	require.Len(t, result, 1)
	item := result[0].(map[string]interface{})

	assert.Equal(t, 123, item[tfconstants.AttrTemplateID])
	assert.Equal(t, "ubuntu", item["distro"])
	assert.Equal(t, "img-123", item["image_id"])
	assert.Equal(t, goe2econstants.ImageStatusReady, item["image_state"])
	assert.Equal(t, "snapshot", item["image_type"])
	assert.Equal(t, "test-image", item["name"])
	assert.Equal(t, "standard", item["sku_type"])
	assert.Equal(t, "Ubuntu", item["os_distribution"])
}

// TestFlattenSavedImages_MultipleImages tests with multiple images
func TestFlattenSavedImages_MultipleImages(t *testing.T) {
	images := []goe2e.SavedImage{
		{
			TemplateID:     123,
			ImageState:     goe2econstants.ImageStatusReady,
			ImageType:      "snapshot",
			OSDistribution: "Ubuntu",
			Name:           "image-one",
			ImageID:        "img-1",
			Distro:         "ubuntu",
			SKUType:        "standard",
		},
		{
			TemplateID:     456,
			ImageState:     goe2econstants.ImageStatusCreating,
			ImageType:      "backup",
			OSDistribution: "CentOS",
			Name:           "image-two",
			ImageID:        "img-2",
			Distro:         "centos",
			SKUType:        "premium",
		},
		{
			TemplateID:     789,
			ImageState:     goe2econstants.ImageStatusError,
			ImageType:      "snapshot",
			OSDistribution: "Debian",
			Name:           "image-three",
			ImageID:        "img-3",
			Distro:         "debian",
			SKUType:        "standard",
		},
	}

	result := flattenSavedImages(images)

	require.Len(t, result, 3)

	// Verify first image
	item1 := result[0].(map[string]interface{})
	assert.Equal(t, 123, item1[tfconstants.AttrTemplateID])
	assert.Equal(t, "img-1", item1["image_id"])
	assert.Equal(t, "image-one", item1["name"])
	assert.Equal(t, goe2econstants.ImageStatusReady, item1["image_state"])

	// Verify second image
	item2 := result[1].(map[string]interface{})
	assert.Equal(t, 456, item2[tfconstants.AttrTemplateID])
	assert.Equal(t, "img-2", item2["image_id"])
	assert.Equal(t, "image-two", item2["name"])
	assert.Equal(t, goe2econstants.ImageStatusCreating, item2["image_state"])

	// Verify third image
	item3 := result[2].(map[string]interface{})
	assert.Equal(t, 789, item3[tfconstants.AttrTemplateID])
	assert.Equal(t, "img-3", item3["image_id"])
	assert.Equal(t, "image-three", item3["name"])
	assert.Equal(t, goe2econstants.ImageStatusError, item3["image_state"])
}

// TestFlattenSavedImages_FieldMapping tests field mapping for each image
func TestFlattenSavedImages_FieldMapping(t *testing.T) {
	images := []goe2e.SavedImage{
		{
			TemplateID:     999,
			ImageState:     goe2econstants.ImageStatusReady,
			ImageType:      "custom",
			OSDistribution: "RHEL",
			Name:           "test-image-name",
			ImageID:        "test-image-id",
			Distro:         "rhel",
			SKUType:        "enterprise",
		},
	}

	result := flattenSavedImages(images)
	require.Len(t, result, 1)
	item := result[0].(map[string]interface{})

	// Verify all fields are mapped correctly
	assert.Equal(t, 999, item[tfconstants.AttrTemplateID], "template_id should be mapped")
	assert.Equal(t, "rhel", item["distro"], "distro should be mapped")
	assert.Equal(t, "test-image-id", item["image_id"], "image_id should be mapped")
	assert.Equal(t, goe2econstants.ImageStatusReady, item["image_state"], "image_state should be mapped")
	assert.Equal(t, "custom", item["image_type"], "image_type should be mapped")
	assert.Equal(t, "test-image-name", item["name"], "name should be mapped")
	assert.Equal(t, "enterprise", item["sku_type"], "sku_type should be mapped")
	assert.Equal(t, "RHEL", item["os_distribution"], "os_distribution should be mapped")
}

// TestFlattenSavedImages_DifferentStates tests with images in different states
func TestFlattenSavedImages_DifferentStates(t *testing.T) {
	tests := []struct {
		name        string
		imageState  string
		description string
	}{
		{
			name:        "Creating state",
			imageState:  goe2econstants.ImageStatusCreating,
			description: "Creating state should be preserved",
		},
		{
			name:        "Ready state",
			imageState:  goe2econstants.ImageStatusReady,
			description: "Ready state should be preserved",
		},
		{
			name:        "Error state",
			imageState:  goe2econstants.ImageStatusError,
			description: "Error state should be preserved",
		},
		{
			name:        "Deleted state",
			imageState:  goe2econstants.ImageStatusDeleted,
			description: "Deleted state should be preserved",
		},
		{
			name:        "Empty state",
			imageState:  "",
			description: "Empty state should be preserved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			images := []goe2e.SavedImage{
				{
					TemplateID:     123,
					ImageState:     tt.imageState,
					ImageType:      "snapshot",
					OSDistribution: "Ubuntu",
					Name:           "test-image",
					ImageID:        "img-123",
					Distro:         "ubuntu",
					SKUType:        "standard",
				},
			}

			result := flattenSavedImages(images)
			require.Len(t, result, 1)
			item := result[0].(map[string]interface{})

			assert.Equal(t, tt.imageState, item["image_state"], tt.description)
		})
	}
}

// TestFlattenSavedImages_EmptyZeroValues tests with images with empty/zero values
func TestFlattenSavedImages_EmptyZeroValues(t *testing.T) {
	images := []goe2e.SavedImage{
		{
			TemplateID:     0,
			ImageState:     "",
			ImageType:      "",
			OSDistribution: "",
			Name:           "",
			ImageID:        "",
			Distro:         "",
			SKUType:        "",
		},
	}

	result := flattenSavedImages(images)
	require.Len(t, result, 1)
	item := result[0].(map[string]interface{})

	// Verify empty fields are handled
	assert.Equal(t, 0, item[tfconstants.AttrTemplateID])
	assert.Equal(t, "", item["distro"])
	assert.Equal(t, "", item["image_id"])
	assert.Equal(t, "", item["image_state"])
	assert.Equal(t, "", item["image_type"])
	assert.Equal(t, "", item["name"])
	assert.Equal(t, "", item["sku_type"])
	assert.Equal(t, "", item["os_distribution"])
}

// TestFlattenSavedImages_MixedValues tests with images having mixed populated and empty fields
func TestFlattenSavedImages_MixedValues(t *testing.T) {
	images := []goe2e.SavedImage{
		{
			TemplateID:     123,
			ImageState:     goe2econstants.ImageStatusReady,
			ImageType:      "", // Empty
			OSDistribution: "Ubuntu",
			Name:           "test-image",
			ImageID:        "img-123",
			Distro:         "", // Empty
			SKUType:        "standard",
		},
	}

	result := flattenSavedImages(images)
	require.Len(t, result, 1)
	item := result[0].(map[string]interface{})

	// Verify populated fields
	assert.Equal(t, 123, item[tfconstants.AttrTemplateID])
	assert.Equal(t, goe2econstants.ImageStatusReady, item["image_state"])
	assert.Equal(t, "Ubuntu", item["os_distribution"])
	assert.Equal(t, "test-image", item["name"])
	assert.Equal(t, "img-123", item["image_id"])
	assert.Equal(t, "standard", item["sku_type"])

	// Verify empty fields
	assert.Equal(t, "", item["image_type"])
	assert.Equal(t, "", item["distro"])
}
