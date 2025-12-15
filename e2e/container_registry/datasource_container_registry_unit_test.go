package container_registry

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestDataSourceContainerRegistrySchema validates the schema structure
func TestDataSourceContainerRegistrySchema(t *testing.T) {
	dataSource := DataSourceContainerRegistry()
	assert.NotNil(t, dataSource)
	assert.NotNil(t, dataSource.Schema)
	assert.NotNil(t, dataSource.ReadContext)
}

// TestDataSourceContainerRegistrySchema_RequiredFields verifies ID field is required
func TestDataSourceContainerRegistrySchema_RequiredFields(t *testing.T) {
	dataSource := DataSourceContainerRegistry()

	// Test ID field
	idField, ok := dataSource.Schema[tfconstants.AttrID]
	assert.True(t, ok, "id field should exist")
	assert.NotNil(t, idField)
	assert.True(t, idField.Required, "id should be required")
	assert.Equal(t, schema.TypeString, idField.Type)
}

// TestDataSourceContainerRegistrySchema_ComputedFields verifies all output fields are computed
func TestDataSourceContainerRegistrySchema_ComputedFields(t *testing.T) {
	dataSource := DataSourceContainerRegistry()

	computedFields := []string{
		tfconstants.AttrProjectName,
		tfconstants.AttrStatus,
		tfconstants.AttrSetupStatus,
		tfconstants.AttrSeverity,
		tfconstants.AttrPreventVulnerabilities,
	}

	for _, fieldName := range computedFields {
		t.Run("computed_"+fieldName, func(t *testing.T) {
			field, ok := dataSource.Schema[fieldName]
			assert.True(t, ok, "field %s should exist", fieldName)
			assert.NotNil(t, field)
			assert.True(t, field.Computed, "field %s should be computed", fieldName)
			assert.False(t, field.Required, "field %s should not be required", fieldName)
			assert.False(t, field.Optional, "field %s should not be optional", fieldName)
		})
	}

	// ProjectID is special - it uses a helper that makes it optional+computed for data sources
	projectIDField, ok := dataSource.Schema[tfconstants.AttrProjectID]
	assert.True(t, ok, "project_id field should exist")
	assert.NotNil(t, projectIDField)
	assert.True(t, projectIDField.Computed, "project_id should be computed")
}

// TestDataSourceContainerRegistrySchema_DeprecatedFields verifies setup_status has deprecation
func TestDataSourceContainerRegistrySchema_DeprecatedFields(t *testing.T) {
	dataSource := DataSourceContainerRegistry()

	setupStatusField, ok := dataSource.Schema[tfconstants.AttrSetupStatus]
	assert.True(t, ok, "setup_status field should exist")
	assert.NotNil(t, setupStatusField)
	assert.Equal(t, DeprecationMessageSetupStatus, setupStatusField.Deprecated)
	assert.Equal(t, DeprecationMessageSetupStatusAlternative, setupStatusField.Description)
	assert.True(t, setupStatusField.Computed, "setup_status should be computed")
}

// TestDataSourceContainerRegistrySchema_FieldTypes validates field types
func TestDataSourceContainerRegistrySchema_FieldTypes(t *testing.T) {
	dataSource := DataSourceContainerRegistry()

	testCases := []struct {
		fieldName    string
		expectedType schema.ValueType
	}{
		{tfconstants.AttrID, schema.TypeString},
		{tfconstants.AttrProjectName, schema.TypeString},
		{tfconstants.AttrStatus, schema.TypeString},
		{tfconstants.AttrSetupStatus, schema.TypeString},
		{tfconstants.AttrSeverity, schema.TypeString},
		{tfconstants.AttrPreventVulnerabilities, schema.TypeBool},
		{tfconstants.AttrRegion, schema.TypeString},
		{tfconstants.AttrLocation, schema.TypeString},
	}

	for _, tc := range testCases {
		t.Run("type_"+tc.fieldName, func(t *testing.T) {
			field, ok := dataSource.Schema[tc.fieldName]
			assert.True(t, ok, "field %s should exist", tc.fieldName)
			assert.NotNil(t, field)
			assert.Equal(t, tc.expectedType, field.Type, "field %s should have correct type", tc.fieldName)
		})
	}
}

// TestDataSourceContainerRegistrySchema_RegionLocationFields validates region/location schema
func TestDataSourceContainerRegistrySchema_RegionLocationFields(t *testing.T) {
	dataSource := DataSourceContainerRegistry()

	// Region field (new)
	regionField, ok := dataSource.Schema[tfconstants.AttrRegion]
	assert.True(t, ok, "region field should exist")
	assert.NotNil(t, regionField)
	assert.Equal(t, schema.TypeString, regionField.Type)

	// Location field (deprecated)
	locationField, ok := dataSource.Schema[tfconstants.AttrLocation]
	assert.True(t, ok, "location field should exist for backward compatibility")
	assert.NotNil(t, locationField)
	assert.Equal(t, schema.TypeString, locationField.Type)
}

// TestDataSourceContainerRegistrySchema_Descriptions validates all fields have descriptions
func TestDataSourceContainerRegistrySchema_Descriptions(t *testing.T) {
	dataSource := DataSourceContainerRegistry()

	fieldsToCheck := []string{
		tfconstants.AttrID,
		tfconstants.AttrProjectName,
		tfconstants.AttrStatus,
		tfconstants.AttrSetupStatus,
		tfconstants.AttrSeverity,
		tfconstants.AttrPreventVulnerabilities,
	}

	for _, fieldName := range fieldsToCheck {
		t.Run("description_"+fieldName, func(t *testing.T) {
			field, ok := dataSource.Schema[fieldName]
			assert.True(t, ok, "field %s should exist", fieldName)
			assert.NotNil(t, field)
			assert.NotEmpty(t, field.Description, "field %s should have a description", fieldName)
		})
	}
}

// TestDataSourceContainerRegistrySchema_NoOptionalFields verifies no fields are optional
func TestDataSourceContainerRegistrySchema_NoOptionalFields(t *testing.T) {
	dataSource := DataSourceContainerRegistry()

	// In a data source, fields should be either required (input) or computed (output)
	// No fields should be both optional and not computed (except region/location which are handled by helpers)
	for fieldName, field := range dataSource.Schema {
		// Skip region/location as they use special helpers
		if fieldName == tfconstants.AttrRegion || fieldName == tfconstants.AttrLocation {
			continue
		}

		t.Run("not_optional_"+fieldName, func(t *testing.T) {
			if field.Optional {
				assert.True(t, field.Computed, "optional field %s should also be computed", fieldName)
			}
		})
	}
}
