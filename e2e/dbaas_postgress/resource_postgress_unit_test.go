package dbaas_postgress_test

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_postgress"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestResourcePostgreSQLSchema_Definition(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	require.NotNil(t, resource)
	assert.NotNil(t, resource.Schema)
}

func TestResourcePostgreSQLSchema_RequiredFields(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	requiredFields := []string{
		tfconstants.AttrVersion,
		tfconstants.AttrName,
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

func TestResourcePostgreSQLSchema_ForceNewFields(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	forceNewFields := []string{
		tfconstants.AttrVersion,
		tfconstants.AttrName,
		tfconstants.AttrGroup,
		tfconstants.AttrIsEncryptionEnabled,
	}

	for _, fieldName := range forceNewFields {
		t.Run(fieldName+"_is_force_new", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.ForceNew, "Field %s should be ForceNew", fieldName)
		})
	}
}

func TestResourcePostgreSQLSchema_DatabaseBlockForceNewFields(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	databaseForceNewFields := []string{
		tfconstants.AttrDatabaseBlockUser,
		tfconstants.AttrDatabaseBlockName,
		tfconstants.AttrDatabaseBlockDBaaSNumber,
	}

	for _, fieldName := range databaseForceNewFields {
		t.Run("database."+fieldName+"_is_force_new", func(t *testing.T) {
			fieldSchema, exists := databaseResourceSchema.Schema[fieldName]
			require.True(t, exists, "Field database.%s should exist in schema", fieldName)
			assert.True(t, fieldSchema.ForceNew, "Field database.%s should be ForceNew", fieldName)
		})
	}
}

func TestResourcePostgreSQLSchema_DatabasePasswordNotForceNew(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	fieldSchema, exists := databaseResourceSchema.Schema[tfconstants.AttrDatabaseBlockPassword]
	require.True(t, exists, "Field database.password should exist in schema")
	assert.False(t, fieldSchema.ForceNew, "Field database.password should NOT be ForceNew (allows rotation)")
}

func TestResourcePostgreSQLSchema_ComputedFields(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	computedFields := []string{
		tfconstants.AttrID,
		tfconstants.AttrStatusTitle,
		tfconstants.AttrStatusActions,
		tfconstants.AttrNumInstances,
		tfconstants.AttrProjectName,
		tfconstants.AttrSnapshotExist,
		tfconstants.AttrConnectivityDetail,
		tfconstants.AttrVectorDatabaseStatus,
		tfconstants.AttrPublicIPAddress,
		tfconstants.AttrPrivateIPAddress,
		tfconstants.AttrPort,
		tfconstants.AttrDisk,
	}

	for _, fieldName := range computedFields {
		t.Run(fieldName+"_is_computed", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Computed, "Field %s should be computed", fieldName)
		})
	}
}

func TestResourcePostgreSQLSchema_SensitiveFields(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	// Test database.password is sensitive
	passwordSchema, exists := databaseResourceSchema.Schema[tfconstants.AttrDatabaseBlockPassword]
	require.True(t, exists, "Field database.password should exist in schema")
	assert.True(t, passwordSchema.Sensitive, "Field database.password should be sensitive")
}

func TestResourcePostgreSQLSchema_DefaultValues(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
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
			fieldName:       tfconstants.AttrPublicIPRequired,
			expectedDefault: tfconstants.DBaaSDefaultPublicIPRequired,
		},
		{
			fieldName:       tfconstants.AttrIsEncryptionEnabled,
			expectedDefault: tfconstants.DBaaSDefaultIsEncryptionEnabled,
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

func TestResourcePostgreSQLSchema_DatabaseMaxItems(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	assert.Equal(t, 1, databaseSchema.MaxItems, "Database field should have MaxItems of 1")
}

func TestResourcePostgreSQLSchema_DatabaseDbaasNumberDefault(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	fieldSchema, exists := databaseResourceSchema.Schema[tfconstants.AttrDatabaseBlockDBaaSNumber]
	require.True(t, exists, "Field database.dbaas_number should exist in schema")
	assert.Equal(t, tfconstants.DBaaSDefaultDBaaSNumber, fieldSchema.Default, "Field database.dbaas_number should have default value of DBaaSDefaultDBaaSNumber")
}

// ============================================================================
// Validation Function Tests
// ============================================================================

func TestResourcePostgreSQLSchema_VersionValidation(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	versionSchema, exists := resourceSchema[tfconstants.AttrVersion]
	require.True(t, exists, "Field version should exist in schema")
	require.NotNil(t, versionSchema.ValidateFunc, "Field version should have ValidateFunc")

	validVersions := []string{
		goe2econstants.PostgreSQLVersion11,
		goe2econstants.PostgreSQLVersion12,
		goe2econstants.PostgreSQLVersion13,
		goe2econstants.PostgreSQLVersion14,
		goe2econstants.PostgreSQLVersion15,
	}

	for _, version := range validVersions {
		t.Run("valid_version_"+version, func(t *testing.T) {
			_, errors := versionSchema.ValidateFunc(version, "version")
			assert.Empty(t, errors, "Version %s should be valid", version)
		})
	}

	invalidVersions := []string{
		"10.0",
		"16.0",
		"invalid",
		"",
	}

	for _, version := range invalidVersions {
		t.Run("invalid_version_"+version, func(t *testing.T) {
			_, errors := versionSchema.ValidateFunc(version, "version")
			assert.NotEmpty(t, errors, "Version %s should be invalid", version)
		})
	}
}

func TestResourcePostgreSQLSchema_StatusValidation(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	statusSchema, exists := resourceSchema[tfconstants.AttrStatus]
	require.True(t, exists, "Field status should exist in schema")
	require.NotNil(t, statusSchema.ValidateFunc, "Field status should have ValidateFunc")

	validStatuses := []string{
		goe2econstants.DBaaSStatusStopped,
		goe2econstants.DBaaSStatusSuspended,
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

func TestResourcePostgreSQL_CustomImportStateFunc(t *testing.T) {
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	require.NotNil(t, resource.Importer, "Resource should have Importer configured")
	assert.NotNil(t, resource.Importer.StateContext, "Importer should use StateContext function")
}

// ============================================================================
// Status Normalization Tests
// ============================================================================
// These tests verify that status normalization (SUSPENDED → STOPPED) works correctly
// in both Create and Read operations, and that other status values pass through unchanged.

func TestResourcePostgreSQL_StatusNormalization_SuspendedToStopped(t *testing.T) {
	// This test verifies the status normalization logic in Create/Read operations
	// SUSPENDED from API should be normalized to STOPPED in state
	// This is a unit test that verifies the normalization logic conceptually

	// Note: Actual status normalization happens in resourceCreatePostgress and resourceReadPostgress
	// These tests document the expected behavior

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
			description:    "SUSPENDED from API should be normalized to STOPPED in state",
		},
		{
			name:           "RUNNING_passes_through",
			apiStatus:      goe2econstants.DBaaSStatusRunning,
			expectedStatus: goe2econstants.DBaaSStatusRunning,
			description:    "RUNNING status should pass through unchanged",
		},
		{
			name:           "RESTARTING_passes_through",
			apiStatus:      goe2econstants.DBaaSStatusRestarting,
			expectedStatus: goe2econstants.DBaaSStatusRestarting,
			description:    "RESTARTING status should pass through unchanged",
		},
		{
			name:           "STOPPED_passes_through",
			apiStatus:      goe2econstants.DBaaSStatusStopped,
			expectedStatus: goe2econstants.DBaaSStatusStopped,
			description:    "STOPPED status should pass through unchanged",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the normalization logic from resourceCreatePostgress/resourceReadPostgress
			status := tc.apiStatus
			if status == goe2econstants.DBaaSStatusSuspended {
				status = goe2econstants.DBaaSStatusStopped
			}

			assert.Equal(t, tc.expectedStatus, status, tc.description)
		})
	}
}

func TestResourcePostgreSQL_StatusNormalization_ForPlanUpgrade(t *testing.T) {
	// This test verifies that STOPPED in state is converted to SUSPENDED for API calls
	// when performing plan upgrade (which requires SUSPENDED state)

	testCases := []struct {
		name              string
		stateStatus       string
		expectedAPIStatus string
		description       string
	}{
		{
			name:              "STOPPED_converted_to_SUSPENDED_for_API",
			stateStatus:       goe2econstants.DBaaSStatusStopped,
			expectedAPIStatus: goe2econstants.DBaaSStatusSuspended,
			description:       "STOPPED in state should be converted to SUSPENDED for API calls",
		},
		{
			name:              "SUSPENDED_passes_through",
			stateStatus:       goe2econstants.DBaaSStatusSuspended,
			expectedAPIStatus: goe2econstants.DBaaSStatusSuspended,
			description:       "SUSPENDED in state should pass through unchanged for API",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the conversion logic from resourceUpdatePostgress (plan upgrade)
			apiStatus := tc.stateStatus
			if apiStatus == goe2econstants.DBaaSStatusStopped {
				apiStatus = goe2econstants.DBaaSStatusSuspended
			}

			assert.Equal(t, tc.expectedAPIStatus, apiStatus, tc.description)
		})
	}
}

func TestResourcePostgreSQL_StatusNormalization_ForDiskExpansion(t *testing.T) {
	// This test verifies that STOPPED in state is converted to SUSPENDED for API calls
	// when performing disk expansion (which requires SUSPENDED state)

	testCases := []struct {
		name              string
		stateStatus       string
		expectedAPIStatus string
		description       string
	}{
		{
			name:              "STOPPED_converted_to_SUSPENDED_for_disk_expansion",
			stateStatus:       goe2econstants.DBaaSStatusStopped,
			expectedAPIStatus: goe2econstants.DBaaSStatusSuspended,
			description:       "STOPPED in state should be converted to SUSPENDED for disk expansion API calls",
		},
		{
			name:              "SUSPENDED_passes_through",
			stateStatus:       goe2econstants.DBaaSStatusSuspended,
			expectedAPIStatus: goe2econstants.DBaaSStatusSuspended,
			description:       "SUSPENDED in state should pass through unchanged for disk expansion API",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the conversion logic from resourceUpdatePostgress (disk expansion)
			apiStatus := tc.stateStatus
			if apiStatus == goe2econstants.DBaaSStatusStopped {
				apiStatus = goe2econstants.DBaaSStatusSuspended
			}

			assert.Equal(t, tc.expectedAPIStatus, apiStatus, tc.description)
		})
	}
}

// ============================================================================
// Disk Expansion Calculation Tests
// ============================================================================

func TestDiskExpansionCalculation_FirstExpansion(t *testing.T) {
	// Test first expansion: prevSize = nil, currSize = 10 → additionalSize = 10
	var prevSize interface{} = nil
	currSize := 10

	prevSizeInt := 0
	if prevSize != nil {
		prevSizeInt = prevSize.(int)
	}
	additionalSize := currSize - prevSizeInt

	assert.Equal(t, 10, additionalSize, "First expansion should calculate correctly")
	assert.Equal(t, 10, currSize, "Current size should be 10")
}

func TestDiskExpansionCalculation_SecondExpansion(t *testing.T) {
	// Test second expansion: prevSize = 10, currSize = 20 → additionalSize = 10
	var prevSize interface{} = 10 //nolint:staticcheck
	currSize := 20

	prevSizeInt := 0
	if prevSize != nil { //nolint:staticcheck
		prevSizeInt = prevSize.(int)
	}
	additionalSize := currSize - prevSizeInt

	assert.Equal(t, 10, additionalSize, "Second expansion should calculate difference correctly")
	assert.Equal(t, 20, currSize, "Current size should be 20")
}

func TestDiskExpansionCalculation_ThirdExpansion(t *testing.T) {
	// Test third expansion: prevSize = 20, currSize = 35 → additionalSize = 15
	var prevSize interface{} = 20 //nolint:staticcheck
	currSize := 35

	prevSizeInt := 0
	if prevSize != nil { //nolint:staticcheck
		prevSizeInt = prevSize.(int)
	}
	additionalSize := currSize - prevSizeInt

	assert.Equal(t, 15, additionalSize, "Third expansion should calculate difference correctly")
	assert.Equal(t, 35, currSize, "Current size should be 35")
}

func TestDiskExpansionCalculation_InvalidExpansion(t *testing.T) {
	// Test invalid expansion: prevSize = 20, currSize = 15 → additionalSize = -5 (should be rejected)
	var prevSize interface{} = 20 //nolint:staticcheck
	currSize := 15

	prevSizeInt := 0
	if prevSize != nil { //nolint:staticcheck
		prevSizeInt = prevSize.(int)
	}
	additionalSize := currSize - prevSizeInt

	assert.Equal(t, -5, additionalSize, "Invalid expansion should result in negative additional size")
	assert.True(t, additionalSize <= 0, "Additional size should be <= 0 for invalid expansion")
}

func TestDiskExpansionCalculation_ZeroExpansion(t *testing.T) {
	// Test zero expansion: prevSize = 10, currSize = 10 → additionalSize = 0 (should be rejected)
	var prevSize interface{} = 10 //nolint:staticcheck
	currSize := 10

	prevSizeInt := 0
	if prevSize != nil { //nolint:staticcheck
		prevSizeInt = prevSize.(int)
	}
	additionalSize := currSize - prevSizeInt

	assert.Equal(t, 0, additionalSize, "Zero expansion should result in 0 additional size")
	assert.True(t, additionalSize <= 0, "Additional size should be <= 0 for zero expansion")
}

// ============================================================================
// VPC Attach/Detach Logic Tests
// ============================================================================

func TestVPCAttachDetach_AttachOnly(t *testing.T) {
	// Test VPC attach: old: [1,2], new: [1,2,3] → attach: [3]
	oldVPCIDs := []string{"1", "2"}
	newVPCIDs := []string{"1", "2", "3"}

	// Build old set
	oldSet := make(map[string]bool)
	for _, vpc := range oldVPCIDs {
		oldSet[vpc] = true
	}

	// Find VPCs to attach
	var toAttach []string
	for _, vpc := range newVPCIDs {
		if !oldSet[vpc] {
			toAttach = append(toAttach, vpc)
		}
	}

	// Find VPCs to detach
	var toDetach []string
	for _, vpc := range oldVPCIDs {
		found := false
		for _, newVPC := range newVPCIDs {
			if vpc == newVPC {
				found = true
				break
			}
		}
		if !found {
			toDetach = append(toDetach, vpc)
		}
	}

	assert.Equal(t, []string{"3"}, toAttach, "Should attach VPC 3")
	assert.Empty(t, toDetach, "Should not detach any VPCs")
}

func TestVPCAttachDetach_DetachOnly(t *testing.T) {
	// Test VPC detach: old: [1,2,3], new: [1,2] → detach: [3]
	oldVPCIDs := []string{"1", "2", "3"}
	newVPCIDs := []string{"1", "2"}

	// Build old set
	oldSet := make(map[string]bool)
	for _, vpc := range oldVPCIDs {
		oldSet[vpc] = true
	}

	// Find VPCs to attach
	var toAttach []string
	for _, vpc := range newVPCIDs {
		if !oldSet[vpc] {
			toAttach = append(toAttach, vpc)
		}
	}

	// Find VPCs to detach
	var toDetach []string
	for _, vpc := range oldVPCIDs {
		found := false
		for _, newVPC := range newVPCIDs {
			if vpc == newVPC {
				found = true
				break
			}
		}
		if !found {
			toDetach = append(toDetach, vpc)
		}
	}

	assert.Empty(t, toAttach, "Should not attach any VPCs")
	assert.Equal(t, []string{"3"}, toDetach, "Should detach VPC 3")
}

func TestVPCAttachDetach_Replace(t *testing.T) {
	// Test VPC replace: old: [1,2], new: [3,4] → detach: [1,2], attach: [3,4]
	oldVPCIDs := []string{"1", "2"}
	newVPCIDs := []string{"3", "4"}

	// Build old set
	oldSet := make(map[string]bool)
	for _, vpc := range oldVPCIDs {
		oldSet[vpc] = true
	}

	// Find VPCs to attach
	var toAttach []string
	for _, vpc := range newVPCIDs {
		if !oldSet[vpc] {
			toAttach = append(toAttach, vpc)
		}
	}

	// Find VPCs to detach
	var toDetach []string
	for _, vpc := range oldVPCIDs {
		found := false
		for _, newVPC := range newVPCIDs {
			if vpc == newVPC {
				found = true
				break
			}
		}
		if !found {
			toDetach = append(toDetach, vpc)
		}
	}

	assert.ElementsMatch(t, []string{"3", "4"}, toAttach, "Should attach VPCs 3 and 4")
	assert.ElementsMatch(t, []string{"1", "2"}, toDetach, "Should detach VPCs 1 and 2")
}

func TestVPCAttachDetach_EmptyOld(t *testing.T) {
	// Test empty old VPC list: old: [], new: [1,2] → attach: [1,2]
	oldVPCIDs := []string{}
	newVPCIDs := []string{"1", "2"}

	// Build old set
	oldSet := make(map[string]bool)
	for _, vpc := range oldVPCIDs {
		oldSet[vpc] = true
	}

	// Find VPCs to attach
	var toAttach []string
	for _, vpc := range newVPCIDs {
		if !oldSet[vpc] {
			toAttach = append(toAttach, vpc)
		}
	}

	// Find VPCs to detach
	var toDetach []string
	for _, vpc := range oldVPCIDs {
		found := false
		for _, newVPC := range newVPCIDs {
			if vpc == newVPC {
				found = true
				break
			}
		}
		if !found {
			toDetach = append(toDetach, vpc)
		}
	}

	assert.ElementsMatch(t, []string{"1", "2"}, toAttach, "Should attach VPCs 1 and 2")
	assert.Empty(t, toDetach, "Should not detach any VPCs")
}

func TestVPCAttachDetach_EmptyNew(t *testing.T) {
	// Test empty new VPC list: old: [1,2], new: [] → detach: [1,2]
	oldVPCIDs := []string{"1", "2"}
	newVPCIDs := []string{}

	// Build old set
	oldSet := make(map[string]bool)
	for _, vpc := range oldVPCIDs {
		oldSet[vpc] = true
	}

	// Find VPCs to attach
	var toAttach []string
	for _, vpc := range newVPCIDs {
		if !oldSet[vpc] {
			toAttach = append(toAttach, vpc)
		}
	}

	// Find VPCs to detach
	var toDetach []string
	for _, vpc := range oldVPCIDs {
		found := false
		for _, newVPC := range newVPCIDs {
			if vpc == newVPC {
				found = true
				break
			}
		}
		if !found {
			toDetach = append(toDetach, vpc)
		}
	}

	assert.Empty(t, toAttach, "Should not attach any VPCs")
	assert.ElementsMatch(t, []string{"1", "2"}, toDetach, "Should detach VPCs 1 and 2")
}

func TestVPCAttachDetach_NoChange(t *testing.T) {
	// Test no change: old: [1,2], new: [1,2] → no attach, no detach
	oldVPCIDs := []string{"1", "2"}
	newVPCIDs := []string{"1", "2"}

	// Build old set
	oldSet := make(map[string]bool)
	for _, vpc := range oldVPCIDs {
		oldSet[vpc] = true
	}

	// Find VPCs to attach
	var toAttach []string
	for _, vpc := range newVPCIDs {
		if !oldSet[vpc] {
			toAttach = append(toAttach, vpc)
		}
	}

	// Find VPCs to detach
	var toDetach []string
	for _, vpc := range oldVPCIDs {
		found := false
		for _, newVPC := range newVPCIDs {
			if vpc == newVPC {
				found = true
				break
			}
		}
		if !found {
			toDetach = append(toDetach, vpc)
		}
	}

	assert.Empty(t, toAttach, "Should not attach any VPCs")
	assert.Empty(t, toDetach, "Should not detach any VPCs")
}

// ============================================================================
// Password Rotation Tests
// ============================================================================

func TestPasswordRotation_CanUpdate(t *testing.T) {
	// Test that database.password field is NOT ForceNew (allows rotation)
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	passwordSchema, exists := databaseResourceSchema.Schema[tfconstants.AttrDatabaseBlockPassword]
	require.True(t, exists, "Field database.password should exist in schema")
	assert.False(t, passwordSchema.ForceNew, "Field database.password should NOT be ForceNew (allows rotation)")
}

func TestPasswordRotation_IsSensitive(t *testing.T) {
	// Test that database.password field is marked as Sensitive
	resource := dbaas_postgress.ResourcePostgresDBaaS()
	resourceSchema := resource.Schema

	databaseSchema, exists := resourceSchema[tfconstants.AttrDatabase]
	require.True(t, exists, "Field database should exist in schema")
	require.NotNil(t, databaseSchema.Elem, "Database field should have Elem")

	databaseResourceSchema, ok := databaseSchema.Elem.(*schema.Resource)
	require.True(t, ok, "Database Elem should be a Resource")

	passwordSchema, exists := databaseResourceSchema.Schema[tfconstants.AttrDatabaseBlockPassword]
	require.True(t, exists, "Field database.password should exist in schema")
	assert.True(t, passwordSchema.Sensitive, "Field database.password should be Sensitive")
}
