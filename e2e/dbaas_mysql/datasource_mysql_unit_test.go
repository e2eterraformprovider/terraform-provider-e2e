package dbaas_mysql_test

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestDataSourceMySQLSchema_Definition(t *testing.T) {
	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	require.NotNil(t, resource)
	assert.NotNil(t, resource.Schema)
}

func TestDataSourceMySQLSchema_RequiredFields(t *testing.T) {
	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	schema := resource.Schema

	requiredFields := []string{
		tfconstants.AttrID,
	}

	for _, fieldName := range requiredFields {
		t.Run(fieldName+"_is_required", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Required, "Field %s should be required", fieldName)
		})
	}
}

func TestDataSourceMySQLSchema_ComputedFields(t *testing.T) {
	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	schema := resource.Schema

	computedFields := []string{
		tfconstants.AttrDatabaseID,
		tfconstants.AttrDatabaseName,
		tfconstants.AttrDatabaseUser,
		tfconstants.AttrStatus,
		tfconstants.AttrPublicIPAddress,
		tfconstants.AttrPrivateIPAddress,
		tfconstants.AttrIsPublicIPAttached,
		tfconstants.AttrDisk,
		tfconstants.AttrPlan,
		tfconstants.AttrDatabaseVersion,
		tfconstants.AttrParameterGroupID,
		tfconstants.AttrPowerStatus,
	}

	for _, fieldName := range computedFields {
		t.Run(fieldName+"_is_computed", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Computed, "Field %s should be computed", fieldName)
		})
	}
}

func TestDataSourceMySQLSchema_IDNotComputed(t *testing.T) {
	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	schema := resource.Schema

	fieldSchema, exists := schema[tfconstants.AttrID]
	require.True(t, exists, "Field id should exist in schema")
	assert.False(t, fieldSchema.Computed, "Field id should NOT be computed (it's required input)")
}
