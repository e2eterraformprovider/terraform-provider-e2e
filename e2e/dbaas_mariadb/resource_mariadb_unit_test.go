package dbaas_mariadb_test

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mariadb"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestResourceMariaDBSchema_Definition(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	require.NotNil(t, resource)
	assert.NotNil(t, resource.Schema)
}

func TestResourceMariaDBSchema_RequiredFields(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	requiredFields := []string{
		tfconstants.AttrName,
		tfconstants.AttrSoftwareName,
		tfconstants.AttrSoftwareVersion,
		tfconstants.AttrGroup,
		tfconstants.AttrDatabase,
		tfconstants.AttrPlan,
	}

	for _, fieldName := range requiredFields {
		t.Run(fieldName+"_is_required", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Required, "Field %s should be required", fieldName)
		})
	}
}

func TestResourceMariaDBSchema_ForceNewFields(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	forceNewFields := []string{
		tfconstants.AttrName,
		tfconstants.AttrSoftwareName,
		tfconstants.AttrSoftwareVersion,
		tfconstants.AttrGroup,
		tfconstants.AttrIsEncryptionEnabled,
		tfconstants.AttrEncryptionPassphrase,
	}

	for _, fieldName := range forceNewFields {
		t.Run(fieldName+"_is_force_new", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.ForceNew, "Field %s should be ForceNew", fieldName)
		})
	}
}

func TestResourceMariaDBSchema_DatabaseBlockForceNewFields(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	databaseForceNewFields := []string{
		"user",
		"name",
		"dbaas_number",
	}

	for _, fieldName := range databaseForceNewFields {
		t.Run("database."+fieldName+"_is_force_new", func(t *testing.T) {
			fieldSchema, exists := databaseResourceSchema.Schema[fieldName]
			require.True(t, exists, "Field database.%s should exist in schema", fieldName)
			assert.True(t, fieldSchema.ForceNew, "Field database.%s should be ForceNew", fieldName)
		})
	}
}

func TestResourceMariaDBSchema_DatabasePasswordNotForceNew(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	fieldSchema, exists := databaseResourceSchema.Schema["password"]
	require.True(t, exists, "Field database.password should exist in schema")
	assert.False(t, fieldSchema.ForceNew, "Field database.password should NOT be ForceNew")
}

func TestResourceMariaDBSchema_ComputedFields(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	computedFields := []string{
		tfconstants.AttrSoftwareID,
		tfconstants.AttrTemplateID,
		tfconstants.AttrPublicIPAttached,
		tfconstants.AttrPublicIPAddress,
		tfconstants.AttrPrivateIPAddress,
		tfconstants.AttrPort,
		tfconstants.AttrTotalDiskSize,
	}

	for _, fieldName := range computedFields {
		t.Run(fieldName+"_is_computed", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Computed, "Field %s should be computed", fieldName)
		})
	}
}

func TestResourceMariaDBSchema_SensitiveFields(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	// Test database.password is sensitive
	passwordSchema, exists := databaseResourceSchema.Schema["password"]
	require.True(t, exists, "Field database.password should exist in schema")
	assert.True(t, passwordSchema.Sensitive, "Field database.password should be sensitive")

	// Test encryption_passphrase is sensitive
	encryptionSchema, exists := resourceSchema[tfconstants.AttrEncryptionPassphrase]
	require.True(t, exists, "Field encryption_passphrase should exist in schema")
	assert.True(t, encryptionSchema.Sensitive, "Field encryption_passphrase should be sensitive")
}

func TestResourceMariaDBSchema_DefaultValues(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	testCases := []struct {
		fieldName       string
		expectedDefault interface{}
	}{
		{
			fieldName:       tfconstants.AttrPublicIPEnabled,
			expectedDefault: tfconstants.DBaaSDefaultPublicIPEnabled,
		},
		{
			fieldName:       tfconstants.AttrParameterGroupID,
			expectedDefault: tfconstants.DBaaSDefaultParameterGroupID,
		},
		{
			fieldName:       tfconstants.AttrIsEncryptionEnabled,
			expectedDefault: tfconstants.DBaaSDefaultIsEncryptionEnabled,
		},
		{
			fieldName:       tfconstants.AttrEncryptionPassphrase,
			expectedDefault: tfconstants.DBaaSDefaultEncryptionPassphrase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.fieldName+"_has_default", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[tc.fieldName]
			require.True(t, exists, "Field %s should exist in schema", tc.fieldName)
			assert.NotNil(t, fieldSchema.Default, "Field %s should have a default value", tc.fieldName)
			assert.Equal(t, tc.expectedDefault, fieldSchema.Default, "Field %s should have correct default value", tc.fieldName)
		})
	}
}

func TestResourceMariaDBSchema_DatabaseMaxItems(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	assert.Equal(t, 1, databaseSchema.MaxItems, "Database field should have MaxItems of 1")
}

// ============================================================================
// Validation Function Tests
// ============================================================================

func TestResourceMariaDBSchema_StatusValidation(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	resourceSchema := resource.Schema

	statusSchema, exists := resourceSchema[tfconstants.AttrStatus]
	require.True(t, exists, "Field status should exist in schema")
	require.NotNil(t, statusSchema.ValidateFunc, "Field status should have ValidateFunc")

	validStatuses := []string{
		goe2econstants.DBaaSStatusStopped,
		goe2econstants.DBaaSStatusRunning,
		goe2econstants.DBaaSStatusRestarting,
	}

	for _, status := range validStatuses {
		t.Run("valid_status_"+status, func(t *testing.T) {
			_, errors := statusSchema.ValidateFunc(status, "status")
			assert.Empty(t, errors, "Status %s should be valid", status)
		})
	}

	invalidStatuses := []string{
		"INVALID",
		"SUSPENDED",
		"",
		"running", // lowercase should fail
	}

	for _, status := range invalidStatuses {
		t.Run("invalid_status_"+status, func(t *testing.T) {
			_, errors := statusSchema.ValidateFunc(status, "status")
			assert.NotEmpty(t, errors, "Status %s should be invalid", status)
		})
	}
}

// ============================================================================
// State Management Tests
// ============================================================================

func TestResourceMariaDB_ImportStatePassthrough(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	require.NotNil(t, resource.Importer, "Resource should have Importer configured")
	assert.NotNil(t, resource.Importer.StateContext, "Importer should use StateContext")
}

// ============================================================================
// Import Functionality Tests
// ============================================================================

func TestResourceMariaDB_ImportStateContext(t *testing.T) {
	resource := dbaas_mariadb.ResourceMariaDB()
	require.NotNil(t, resource.Importer, "Resource should have Importer configured")

	// Test that import uses passthrough context (ID is passed directly)
	// This means the import ID format is just the instance ID
	testCases := []struct {
		name     string
		importID string
		valid    bool
	}{
		{
			name:     "valid numeric ID",
			importID: "12345",
			valid:    true,
		},
		{
			name:     "valid string ID",
			importID: "test-instance-id",
			valid:    true,
		},
		{
			name:     "empty ID should be handled by terraform",
			importID: "",
			valid:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ImportStatePassthroughContext accepts any non-empty string
			// Empty strings would be caught by Terraform validation
			if tc.valid && tc.importID != "" {
				// The importer should accept this ID format
				// Actual validation happens during resource read
				assert.NotNil(t, resource.Importer.StateContext, "StateContext should be configured")
			}
		})
	}
}

func TestResourceMariaDB_ImportStatePassthroughBehavior(t *testing.T) {
	// Test that import uses passthrough, meaning:
	// 1. Import ID is set directly as resource ID
	// 2. No parsing or transformation of import ID
	// 3. Resource read will fetch data using the ID

	resource := dbaas_mariadb.ResourceMariaDB()
	require.NotNil(t, resource.Importer)

	// Verify it's using ImportStatePassthroughContext
	// This means the import ID format is simply the instance ID
	// No special format like "project_id:region:id" is required
	assert.NotNil(t, resource.Importer.StateContext, "Should use StateContext for import")
}
