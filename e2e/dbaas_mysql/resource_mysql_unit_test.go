package dbaas_mysql_test

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mysql"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestResourceMySQLSchema_Definition(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	require.NotNil(t, resource)
	assert.NotNil(t, resource.Schema)
}

func TestResourceMySQLSchema_RequiredFields(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	resourceSchema := resource.Schema

	requiredFields := []string{
		tfconstants.AttrVersion,
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

func TestResourceMySQLSchema_ForceNewFields(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	resourceSchema := resource.Schema

	forceNewFields := []string{
		tfconstants.AttrVersion,
		tfconstants.AttrDBaaSName,
		tfconstants.AttrDBLocation,
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

func TestResourceMySQLSchema_DatabaseBlockForceNewFields(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
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

func TestResourceMySQLSchema_DatabasePasswordNotForceNew(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
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

func TestResourceMySQLSchema_ComputedFields(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	resourceSchema := resource.Schema

	computedFields := []string{
		tfconstants.AttrDisk,
		tfconstants.AttrPublicIPAddress,
		tfconstants.AttrPrivateIPAddress,
		tfconstants.AttrPort,
	}

	for _, fieldName := range computedFields {
		t.Run(fieldName+"_is_computed", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Computed, "Field %s should be computed", fieldName)
		})
	}
}

func TestResourceMySQLSchema_SensitiveFields(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
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

func TestResourceMySQLSchema_DefaultValues(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	resourceSchema := resource.Schema

	testCases := []struct {
		fieldName       string
		expectedDefault interface{}
	}{
		{
			fieldName:       tfconstants.AttrGroup,
			expectedDefault: tfconstants.DBaaSDefaultGroupName,
		},
		{
			fieldName:       tfconstants.AttrDBLocation,
			expectedDefault: tfconstants.DBaaSDefaultDBLocation,
		},
		{
			fieldName:       tfconstants.AttrPublicIPRequired,
			expectedDefault: tfconstants.DBaaSDefaultPublicIPRequired,
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

func TestResourceMySQLSchema_DatabaseMaxItems(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	assert.Equal(t, 1, databaseSchema.MaxItems, "Database field should have MaxItems of 1")
}

func TestResourceMySQLSchema_DatabaseDbaasNumberDefault(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	fieldSchema, exists := databaseResourceSchema.Schema["dbaas_number"]
	require.True(t, exists, "Field database.dbaas_number should exist in schema")
	assert.Equal(t, tfconstants.DBaaSDefaultDBaaSNumber, fieldSchema.Default, "Field database.dbaas_number should have correct default value")
}

// ============================================================================
// Validation Function Tests
// ============================================================================

func TestResourceMySQLSchema_StatusValidation(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	resourceSchema := resource.Schema

	statusSchema, exists := resourceSchema[tfconstants.AttrStatus]
	require.True(t, exists, "Field status should exist in schema")
	require.NotNil(t, statusSchema.ValidateFunc, "Field status should have ValidateFunc")

	validStatuses := []string{
		goe2econstants.DBaaSStatusStopped,
		goe2econstants.DBaaSStatusSuspended,
		goe2econstants.DBaaSStatusRunning,
		goe2econstants.DBaaSStatusRestarting,
		tfconstants.DBaaSPowerActionStart,
		tfconstants.DBaaSPowerActionStop,
		tfconstants.DBaaSPowerActionRestart,
	}

	for _, status := range validStatuses {
		t.Run("valid_status_"+status, func(t *testing.T) {
			_, errors := statusSchema.ValidateFunc(status, "status")
			assert.Empty(t, errors, "Status %s should be valid", status)
		})
	}

	invalidStatuses := []string{
		"INVALID",
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

func TestResourceMySQL_CustomImportStateFunc(t *testing.T) {
	resource := dbaas_mysql.ResourceMySql()
	require.NotNil(t, resource.Importer, "Resource should have Importer configured")
	assert.NotNil(t, resource.Importer.State, "Importer should use State function")
}
