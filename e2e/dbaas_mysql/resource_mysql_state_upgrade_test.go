package dbaas_mysql_test

import (
	"context"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// State Upgrade V0 → V1 Tests
// ============================================================================

func TestResourceMySQLStateUpgradeV0toV1_Basic(t *testing.T) {
	v0State := map[string]interface{}{
		"id":                 "12345",
		"version":            "8.0",
		"dbaas_name":         "my-mysql",
		"db_location":        "Delhi",
		"group":              "Default",
		"plan":               "DBS.16GB",
		"public_ip_required": true,
		"region":             "Delhi",
		"project_id":         "789",
		"database": []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "8.0", v1State["version"])
	assert.Equal(t, "my-mysql", v1State["dbaas_name"])
	assert.Equal(t, "Delhi", v1State["db_location"])
	assert.Equal(t, "Default", v1State["group"])
	assert.Equal(t, "DBS.16GB", v1State["plan"])
	assert.Equal(t, true, v1State["public_ip_required"])

	// New V1 fields should be added
	assert.NotNil(t, v1State["tags"], "tags field should be added")
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

func TestResourceMySQLStateUpgradeV0toV1_WithExistingTags(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "8.0",
		"plan":    "DBS.16GB",
		"tags": map[string]interface{}{
			"Environment": "production",
			"Team":        "backend",
		},
	}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	assert.NotNil(t, v1State["tags"])
	// Existing tags should be preserved
	tags := v1State["tags"].(map[string]interface{})
	assert.Equal(t, "production", tags["Environment"])
	assert.Equal(t, "backend", tags["Team"])
}

func TestResourceMySQLStateUpgradeV0toV1_PreservesAllFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":                    "67890",
		"version":               "8.0",
		"dbaas_name":            "test-mysql",
		"db_location":           "Mumbai",
		"group":                 "Default",
		"plan":                  "DBS.32GB",
		"public_ip_required":    true,
		"is_encryption_enabled": true,
		"encryption_passphrase": "secret123",
		"region":                "Mumbai",
		"project_id":            "123",
		"parameter_group_id":    5,
		"vpcs":                  []interface{}{1, 2},
		"size":                  100,
		"status":                "RUNNING",
		"disk":                  "100GB",
		"public_ip_address":     "1.2.3.4",
		"private_ip_address":    "10.0.0.1",
		"port":                  "3306",
		"database": []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "67890", v1State["id"])
	assert.Equal(t, "8.0", v1State["version"])
	assert.Equal(t, "test-mysql", v1State["dbaas_name"])
	assert.Equal(t, "Mumbai", v1State["db_location"])
	assert.Equal(t, "Default", v1State["group"])
	assert.Equal(t, "DBS.32GB", v1State["plan"])
	assert.Equal(t, true, v1State["public_ip_required"])
	assert.Equal(t, true, v1State["is_encryption_enabled"])
	assert.Equal(t, "secret123", v1State["encryption_passphrase"])
	assert.Equal(t, "Mumbai", v1State["region"])
	assert.Equal(t, "123", v1State["project_id"])
	assert.Equal(t, 5, v1State["parameter_group_id"])
	assert.NotNil(t, v1State["vpcs"])
	assert.Equal(t, 100, v1State["size"])
	assert.Equal(t, "RUNNING", v1State["status"])
	assert.Equal(t, "100GB", v1State["disk"])
	assert.Equal(t, "1.2.3.4", v1State["public_ip_address"])
	assert.Equal(t, "10.0.0.1", v1State["private_ip_address"])
	assert.Equal(t, "3306", v1State["port"])
	assert.NotNil(t, v1State["database"])

	// New V1 fields should be added
	assert.NotNil(t, v1State["tags"])
}

func TestResourceMySQLStateUpgradeV0toV1_WithMissingFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "8.0",
		"plan":    "DBS.16GB",
		// Missing many fields - should handle gracefully
	}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// Existing fields should be preserved
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "8.0", v1State["version"])
	assert.Equal(t, "DBS.16GB", v1State["plan"])

	// New V1 fields should be added
	assert.NotNil(t, v1State["tags"])
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

func TestResourceMySQLStateUpgradeV0toV1_EmptyState(t *testing.T) {
	v0State := map[string]interface{}{
		"id": "12345",
	}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	assert.Equal(t, "12345", v1State["id"])
	assert.NotNil(t, v1State["tags"])
}

func TestResourceMySQLStateUpgradeV0toV1_NoErrors(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "8.0",
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

	_, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)
	assert.NoError(t, err, "State upgrade should return no errors for valid V0 state")
}

func TestResourceMySQLStateUpgradeV0toV1_WithInvalidTagsType(t *testing.T) {
	v0State := map[string]interface{}{
		"id":   "12345",
		"tags": "invalid-string-type", // Wrong type - should be map
	}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	// Should handle gracefully - overwrite invalid type with empty map
	require.NoError(t, err)
	assert.NotNil(t, v1State["tags"])
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be converted to map")
	assert.Empty(t, tags, "tags should be empty map after conversion")
}

func TestResourceMySQLStateUpgradeV0toV1_WithNilValues(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": nil, // Nil value
		"plan":    "E2E-2C-4GB",
	}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	assert.Equal(t, "12345", v1State["id"])
	assert.Nil(t, v1State["version"]) // Nil values should be preserved
	assert.Equal(t, "E2E-2C-4GB", v1State["plan"])
	assert.NotNil(t, v1State["tags"])
}

func TestResourceMySQLStateUpgradeV0toV1_WithEmptyState(t *testing.T) {
	v0State := map[string]interface{}{}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	assert.NotNil(t, v1State["tags"])
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok)
	assert.Empty(t, tags)
}

func TestResourceMySQLStateUpgradeV0toV1_WithNestedStructures(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"version": "8.0",
		"database": []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
		"vpcs": []interface{}{1, 2, 3},
	}

	v1State, err := dbaas_mysql.ResourceMySQLStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// Nested structures should be preserved
	assert.NotNil(t, v1State["database"])
	assert.NotNil(t, v1State["vpcs"])
	dbList := v1State["database"].([]interface{})
	assert.Len(t, dbList, 1)
	vpcsList := v1State["vpcs"].([]interface{})
	assert.Len(t, vpcsList, 3)
	// Tags should be added
	assert.NotNil(t, v1State["tags"])
}
