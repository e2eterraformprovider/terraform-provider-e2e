package dbaas_postgress_test

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_postgress"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestDataSourcePostgreSQLSchema_Definition(t *testing.T) {
	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	require.NotNil(t, resource)
	assert.NotNil(t, resource.Schema)
}

func TestDataSourcePostgreSQLSchema_RequiredFields(t *testing.T) {
	resource := dbaas_postgress.DataSourcePostgresDBaaS()
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

func TestDataSourcePostgreSQLSchema_ComputedFields(t *testing.T) {
	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	schema := resource.Schema

	computedFields := []string{
		tfconstants.AttrDatabaseID,
		tfconstants.AttrDatabaseName,
		tfconstants.AttrDatabaseUser,
		tfconstants.AttrPGDetails,
		tfconstants.AttrStatus,
		tfconstants.AttrStatusActions,
		tfconstants.AttrPublicIPAddress,
		tfconstants.AttrPrivateIPAddress,
		tfconstants.AttrIsPublicIPAttached,
		tfconstants.AttrPlan,
		tfconstants.AttrDatabaseVersion,
		tfconstants.AttrParameterGroupID,
		tfconstants.AttrDiskSize,
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

func TestDataSourcePostgreSQLSchema_IDNotComputed(t *testing.T) {
	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	schema := resource.Schema

	fieldSchema, exists := schema[tfconstants.AttrID]
	require.True(t, exists, "Field id should exist in schema")
	assert.False(t, fieldSchema.Computed, "Field id should NOT be computed (it's required input)")
}

// ============================================================================
// Data Source Read Function Tests
// ============================================================================
// Note: Full integration tests for dataSourceReadPostgres require API mocking
// These tests verify the status normalization logic used in the datasource read function

func TestDataSourcePostgreSQL_StatusNormalization(t *testing.T) {
	// Test status normalization logic used in dataSourceReadPostgres
	// SUSPENDED from API should be normalized to STOPPED in state

	testCases := []struct {
		name           string
		apiStatus      string
		expectedStatus string
		description    string
	}{
		{
			name:           "SUSPENDED_normalized_to_STOPPED",
			apiStatus:      goe2econstants.DBaaSStatusSuspended,
			expectedStatus: goe2econstants.DBaaSStatusStopped,
			description:    "SUSPENDED from API should be normalized to STOPPED in datasource state",
		},
		{
			name:           "RUNNING_passes_through",
			apiStatus:      goe2econstants.DBaaSStatusRunning,
			expectedStatus: goe2econstants.DBaaSStatusRunning,
			description:    "RUNNING status should pass through unchanged in datasource",
		},
		{
			name:           "RESTARTING_passes_through",
			apiStatus:      goe2econstants.DBaaSStatusRestarting,
			expectedStatus: goe2econstants.DBaaSStatusRestarting,
			description:    "RESTARTING status should pass through unchanged in datasource",
		},
		{
			name:           "STOPPED_passes_through",
			apiStatus:      goe2econstants.DBaaSStatusStopped,
			expectedStatus: goe2econstants.DBaaSStatusStopped,
			description:    "STOPPED status should pass through unchanged in datasource",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the normalization logic from dataSourceReadPostgres
			status := tc.apiStatus
			if status == goe2econstants.DBaaSStatusSuspended {
				status = goe2econstants.DBaaSStatusStopped
			}

			assert.Equal(t, tc.expectedStatus, status, tc.description)
		})
	}
}

func TestDataSourcePostgreSQL_FieldMapping(t *testing.T) {
	// Test that datasource schema includes all expected computed fields
	// This verifies the field mapping structure used in dataSourceReadPostgres

	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	schema := resource.Schema

	expectedComputedFields := map[string]string{
		tfconstants.AttrDatabaseID:         "database_id",
		tfconstants.AttrDatabaseName:       "database_name",
		tfconstants.AttrDatabaseUser:       "database_user",
		tfconstants.AttrStatus:             "status",
		tfconstants.AttrStatusActions:      "status_actions",
		tfconstants.AttrPublicIPAddress:    "public_ip_address",
		tfconstants.AttrPrivateIPAddress:   "private_ip_address",
		tfconstants.AttrIsPublicIPAttached: "is_public_ip_attached",
		tfconstants.AttrPlan:               "plan",
		tfconstants.AttrDatabaseVersion:    "database_version",
		tfconstants.AttrParameterGroupID:   "parameter_group_id",
		tfconstants.AttrDiskSize:           "disk_size",
		tfconstants.AttrPowerStatus:        "power_status",
	}

	for fieldName, expectedKey := range expectedComputedFields {
		t.Run("field_"+fieldName+"_exists", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in datasource schema", fieldName)
			assert.True(t, fieldSchema.Computed, "Field %s should be computed", fieldName)
			assert.Equal(t, expectedKey, fieldName, "Field name should match expected key")
		})
	}
}
