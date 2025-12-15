package dbaas_postgress_test

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_postgress"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// State Upgrade V0 → V1 Tests
// ============================================================================

func TestResourcePostgreSQLStateUpgradeV0toV1_Basic(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "14.0",
		"name":    "my-postgres",
		"plan":    "DBS.16GB",
		"group":   "Default",
		"database": []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "14.0", v1State["version"])
	assert.Equal(t, "my-postgres", v1State["name"])
	assert.Equal(t, "DBS.16GB", v1State["plan"])
	assert.Equal(t, "Default", v1State["group"])

	// New V1 fields should be added
	assert.NotNil(t, v1State["tags"], "tags field should be added")
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_WithExistingTags(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "14.0",
		"name":    "my-postgres",
		"plan":    "DBS.16GB",
		"tags": map[string]interface{}{
			"Environment": "production",
			"Team":        "backend",
		},
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	assert.NotNil(t, v1State["tags"])
	// Existing tags should be preserved
	tags := v1State["tags"].(map[string]interface{})
	assert.Equal(t, "production", tags["Environment"])
	assert.Equal(t, "backend", tags["Team"])
}

func TestResourcePostgreSQLStateUpgradeV0toV1_PreservesAllFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":         "67890",
		"version":    "14.0",
		"name":       "test-postgres",
		"plan":       "DBS.32GB",
		"group":      "Default",
		"region":     "Mumbai",
		"project_id": "123",
		"database": []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "67890", v1State["id"])
	assert.Equal(t, "14.0", v1State["version"])
	assert.Equal(t, "test-postgres", v1State["name"])
	assert.Equal(t, "DBS.32GB", v1State["plan"])
	assert.Equal(t, "Default", v1State["group"])
	assert.Equal(t, "Mumbai", v1State["region"])
	assert.Equal(t, "123", v1State["project_id"])
	assert.NotNil(t, v1State["database"])

	// New V1 fields should be added
	assert.NotNil(t, v1State["tags"])
}

func TestResourcePostgreSQLStateUpgradeV0toV1_WithMissingFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "14.0",
		"name":    "my-postgres",
		"plan":    "DBS.16GB",
		// Missing many fields - should handle gracefully
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// Existing fields should be preserved
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "14.0", v1State["version"])
	assert.Equal(t, "my-postgres", v1State["name"])
	assert.Equal(t, "DBS.16GB", v1State["plan"])

	// New V1 fields should be added
	assert.NotNil(t, v1State["tags"])
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_EmptyState(t *testing.T) {
	v0State := map[string]interface{}{
		"id": "12345",
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	assert.Equal(t, "12345", v1State["id"])
	assert.NotNil(t, v1State["tags"])
}

func TestResourcePostgreSQLStateUpgradeV0toV1_NoErrors(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "14.0",
		"name":    "my-postgres",
		"plan":    "DBS.16GB",
		"database": []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	_, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)
	assert.NoError(t, err, "State upgrade should return no errors for valid V0 state")
}

// ============================================================================
// Field Migration Tests
// ============================================================================

func TestResourcePostgreSQLStateUpgradeV0toV1_VPCListMigration(t *testing.T) {
	v0State := map[string]interface{}{
		"id":       "12345",
		"version":  "14.0",
		"name":     "my-postgres",
		"plan":     "DBS.16GB",
		"vpc_list": []interface{}{1, 2, 3},
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// vpc_list should be migrated to vpcs
	assert.NotNil(t, v1State[tfconstants.AttrVPCs], "vpcs field should be added")
	assert.Equal(t, []interface{}{1, 2, 3}, v1State[tfconstants.AttrVPCs], "vpcs should match vpc_list")
	// vpc_list should be preserved for backwards compatibility
	assert.NotNil(t, v1State[tfconstants.FieldMigrationKeyVPCList], "vpc_list should be preserved")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_VPCListNil(t *testing.T) {
	v0State := map[string]interface{}{
		"id":       "12345",
		"version":  "14.0",
		"name":     "my-postgres",
		"plan":     "DBS.16GB",
		"vpc_list": nil,
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// vpc_list should NOT be migrated to vpcs if nil (implementation only migrates non-nil values)
	_, exists := v1State[tfconstants.AttrVPCs]
	assert.False(t, exists, "vpcs field should NOT be added when vpc_list is nil")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_VPCListEmpty(t *testing.T) {
	v0State := map[string]interface{}{
		"id":       "12345",
		"version":  "14.0",
		"name":     "my-postgres",
		"plan":     "DBS.16GB",
		"vpc_list": []interface{}{},
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// vpc_list should be migrated to vpcs even if empty
	assert.NotNil(t, v1State[tfconstants.AttrVPCs], "vpcs field should be added")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_DetachPublicIPMigration_True(t *testing.T) {
	v0State := map[string]interface{}{
		"id":               "12345",
		"version":          "14.0",
		"name":             "my-postgres",
		"plan":             "DBS.16GB",
		"detach_public_ip": true,
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// detach_public_ip: true should become public_ip_required: false
	assert.Equal(t, false, v1State[tfconstants.AttrPublicIPRequired], "public_ip_required should be false when detach_public_ip is true")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_DetachPublicIPMigration_False(t *testing.T) {
	v0State := map[string]interface{}{
		"id":               "12345",
		"version":          "14.0",
		"name":             "my-postgres",
		"plan":             "DBS.16GB",
		"detach_public_ip": false,
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// detach_public_ip: false should become public_ip_required: true
	assert.Equal(t, true, v1State[tfconstants.AttrPublicIPRequired], "public_ip_required should be true when detach_public_ip is false")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_DetachPublicIPMissing(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "14.0",
		"name":    "my-postgres",
		"plan":    "DBS.16GB",
		// detach_public_ip is missing
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// public_ip_required should not be set if detach_public_ip is missing
	_, exists := v1State[tfconstants.AttrPublicIPRequired]
	assert.False(t, exists, "public_ip_required should not be set when detach_public_ip is missing")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_DetachPublicIPTypeAssertion(t *testing.T) {
	v0State := map[string]interface{}{
		"id":               "12345",
		"version":          "14.0",
		"name":             "my-postgres",
		"plan":             "DBS.16GB",
		"detach_public_ip": "not-a-bool", // Wrong type
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// Should handle type assertion gracefully - public_ip_required should not be set
	_, exists := v1State[tfconstants.AttrPublicIPRequired]
	assert.False(t, exists, "public_ip_required should not be set when detach_public_ip is not a bool")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_PowerStatusMigration_Start(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "12345",
		"version":      "14.0",
		"name":         "my-postgres",
		"plan":         "DBS.16GB",
		"power_status": tfconstants.DBaaSPowerActionStart,
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// power_status: "start" should become status: "RUNNING"
	assert.Equal(t, goe2econstants.DBaaSStatusRunning, v1State[tfconstants.AttrStatus], "status should be RUNNING when power_status is start")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_PowerStatusMigration_Stop(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "12345",
		"version":      "14.0",
		"name":         "my-postgres",
		"plan":         "DBS.16GB",
		"power_status": tfconstants.DBaaSPowerActionStop,
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// power_status: "stop" should become status: "SUSPENDED"
	assert.Equal(t, goe2econstants.DBaaSStatusSuspended, v1State[tfconstants.AttrStatus], "status should be SUSPENDED when power_status is stop")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_PowerStatusMigration_Restart(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "12345",
		"version":      "14.0",
		"name":         "my-postgres",
		"plan":         "DBS.16GB",
		"power_status": tfconstants.DBaaSPowerActionRestart,
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// power_status: "restart" should become status: "RESTARTING"
	assert.Equal(t, goe2econstants.DBaaSStatusRestarting, v1State[tfconstants.AttrStatus], "status should be RESTARTING when power_status is restart")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_PowerStatusMigration_ExistingStatus(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "12345",
		"version":      "14.0",
		"name":         "my-postgres",
		"plan":         "DBS.16GB",
		"power_status": "RUNNING", // Already a valid status
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// Existing status values should pass through unchanged
	assert.Equal(t, "RUNNING", v1State[tfconstants.AttrStatus], "status should pass through unchanged when power_status is already a valid status")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_PowerStatusMissing(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "14.0",
		"name":    "my-postgres",
		"plan":    "DBS.16GB",
		// power_status is missing
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// status should not be set if power_status is missing
	_, exists := v1State[tfconstants.AttrStatus]
	assert.False(t, exists, "status should not be set when power_status is missing")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_PowerStatusTypeAssertion(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "12345",
		"version":      "14.0",
		"name":         "my-postgres",
		"plan":         "DBS.16GB",
		"power_status": 123, // Wrong type
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// Should handle type assertion gracefully - status should not be set
	_, exists := v1State[tfconstants.AttrStatus]
	assert.False(t, exists, "status should not be set when power_status is not a string")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_PowerStatusEmpty(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "12345",
		"version":      "14.0",
		"name":         "my-postgres",
		"plan":         "DBS.16GB",
		"power_status": "",
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// Empty power_status should not set status
	_, exists := v1State[tfconstants.AttrStatus]
	assert.False(t, exists, "status should not be set when power_status is empty")
}

func TestResourcePostgreSQLStateUpgradeV0toV1_AllFieldsPresent(t *testing.T) {
	v0State := map[string]interface{}{
		"id":               "12345",
		"version":          "14.0",
		"name":             "my-postgres",
		"plan":             "DBS.16GB",
		"group":            "Default",
		"region":           "Mumbai",
		"project_id":       "123",
		"vpc_list":         []interface{}{1, 2},
		"detach_public_ip": false,
		"power_status":     "start",
		"database": []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	v1State, err := dbaas_postgress.ResourcePostgreSQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// All fields should be preserved
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "14.0", v1State["version"])
	assert.Equal(t, "my-postgres", v1State["name"])
	assert.Equal(t, "DBS.16GB", v1State["plan"])
	assert.Equal(t, "Default", v1State["group"])
	assert.Equal(t, "Mumbai", v1State["region"])
	assert.Equal(t, "123", v1State["project_id"])
	assert.NotNil(t, v1State["database"])

	// Migrations should be applied
	assert.NotNil(t, v1State[tfconstants.AttrVPCs])
	assert.Equal(t, true, v1State[tfconstants.AttrPublicIPRequired])
	assert.Equal(t, goe2econstants.DBaaSStatusRunning, v1State[tfconstants.AttrStatus])

	// New V1 fields should be added
	assert.NotNil(t, v1State["tags"])
}
