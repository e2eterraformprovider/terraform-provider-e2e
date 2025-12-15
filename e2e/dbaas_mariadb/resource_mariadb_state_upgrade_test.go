package dbaas_mariadb_test

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mariadb"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// State Upgrade V0 → V1 Tests
// ============================================================================

func TestResourceMariaDBStateUpgradeV0toV1_Basic(t *testing.T) {
	v0State := map[string]interface{}{
		tfconstants.AttrID:                  "12345",
		tfconstants.AttrName:                "my-mariadb",
		tfconstants.AttrSoftwareName:        "MariaDB",
		tfconstants.AttrSoftwareVersion:     "10.6",
		tfconstants.AttrGroup:               "Default",
		tfconstants.AttrPlan:                "DBS.16GB",
		tfconstants.AttrPublicIPEnabled:     true,
		tfconstants.AttrParameterGroupID:    tfconstants.DBaaSDefaultParameterGroupID,
		tfconstants.AttrIsEncryptionEnabled: tfconstants.DBaaSDefaultIsEncryptionEnabled,
		tfconstants.AttrRegion:              "Delhi",
		tfconstants.AttrProjectID:           "789",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	v1State, err := dbaas_mariadb.ResourceMariaDBStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "12345", v1State[tfconstants.AttrID])
	assert.Equal(t, "my-mariadb", v1State[tfconstants.AttrName])
	assert.Equal(t, "MariaDB", v1State[tfconstants.AttrSoftwareName])
	assert.Equal(t, "10.6", v1State[tfconstants.AttrSoftwareVersion])
	assert.Equal(t, "Default", v1State[tfconstants.AttrGroup])
	assert.Equal(t, "DBS.16GB", v1State[tfconstants.AttrPlan])
	assert.Equal(t, true, v1State[tfconstants.AttrPublicIPEnabled])
	assert.Equal(t, tfconstants.DBaaSDefaultParameterGroupID, v1State[tfconstants.AttrParameterGroupID])
	assert.Equal(t, tfconstants.DBaaSDefaultIsEncryptionEnabled, v1State[tfconstants.AttrIsEncryptionEnabled])

	// New V1 fields should be added
	assert.NotNil(t, v1State[tfconstants.AttrTags], "tags field should be added")
	tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

func TestResourceMariaDBStateUpgradeV0toV1_WithExistingTags(t *testing.T) {
	v0State := map[string]interface{}{
		tfconstants.AttrID:   "12345",
		tfconstants.AttrName: "my-mariadb",
		tfconstants.AttrPlan: "DBS.16GB",
		tfconstants.AttrTags: map[string]interface{}{
			"Environment": "production",
			"Team":        "backend",
		},
	}

	v1State, err := dbaas_mariadb.ResourceMariaDBStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	assert.NotNil(t, v1State[tfconstants.AttrTags])
	// Existing tags should be preserved
	tags := v1State[tfconstants.AttrTags].(map[string]interface{})
	assert.Equal(t, "production", tags["Environment"])
	assert.Equal(t, "backend", tags["Team"])
}

func TestResourceMariaDBStateUpgradeV0toV1_PreservesAllFields(t *testing.T) {
	v0State := map[string]interface{}{
		tfconstants.AttrID:                   "67890",
		tfconstants.AttrName:                 "test-mariadb",
		tfconstants.AttrSoftwareName:         "MariaDB",
		tfconstants.AttrSoftwareVersion:      "10.6",
		tfconstants.AttrGroup:                "Default",
		tfconstants.AttrPlan:                 "DBS.32GB",
		tfconstants.AttrPublicIPEnabled:      true,
		tfconstants.AttrParameterGroupID:     5,
		tfconstants.AttrIsEncryptionEnabled:  true,
		tfconstants.AttrEncryptionPassphrase: "secret123",
		tfconstants.AttrRegion:               "Mumbai",
		tfconstants.AttrProjectID:            "123",
		tfconstants.AttrSoftwareID:           100,
		tfconstants.AttrTemplateID:           200,
		tfconstants.AttrPublicIPAttached:     true,
		tfconstants.AttrPublicIPAddress:      "1.2.3.4",
		tfconstants.AttrPrivateIPAddress:     "10.0.0.1",
		tfconstants.AttrPort:                 "3306",
		tfconstants.AttrTotalDiskSize:        100,
		tfconstants.AttrStatus:               goe2econstants.DBaaSStatusRunning,
		tfconstants.AttrDiskSize:             50,
		tfconstants.AttrVPCs:                 []interface{}{"vpc1", "vpc2"},
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	v1State, err := dbaas_mariadb.ResourceMariaDBStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "67890", v1State[tfconstants.AttrID])
	assert.Equal(t, "test-mariadb", v1State[tfconstants.AttrName])
	assert.Equal(t, "MariaDB", v1State[tfconstants.AttrSoftwareName])
	assert.Equal(t, "10.6", v1State[tfconstants.AttrSoftwareVersion])
	assert.Equal(t, "Default", v1State[tfconstants.AttrGroup])
	assert.Equal(t, "DBS.32GB", v1State[tfconstants.AttrPlan])
	assert.Equal(t, true, v1State[tfconstants.AttrPublicIPEnabled])
	assert.Equal(t, 5, v1State[tfconstants.AttrParameterGroupID])
	assert.Equal(t, true, v1State[tfconstants.AttrIsEncryptionEnabled])
	assert.Equal(t, "secret123", v1State[tfconstants.AttrEncryptionPassphrase])
	assert.Equal(t, "Mumbai", v1State[tfconstants.AttrRegion])
	assert.Equal(t, "123", v1State[tfconstants.AttrProjectID])
	assert.Equal(t, 100, v1State[tfconstants.AttrSoftwareID])
	assert.Equal(t, 200, v1State[tfconstants.AttrTemplateID])
	assert.Equal(t, true, v1State[tfconstants.AttrPublicIPAttached])
	assert.Equal(t, "1.2.3.4", v1State[tfconstants.AttrPublicIPAddress])
	assert.Equal(t, "10.0.0.1", v1State[tfconstants.AttrPrivateIPAddress])
	assert.Equal(t, "3306", v1State[tfconstants.AttrPort])
	assert.Equal(t, 100, v1State[tfconstants.AttrTotalDiskSize])
	assert.Equal(t, goe2econstants.DBaaSStatusRunning, v1State[tfconstants.AttrStatus])
	assert.Equal(t, 50, v1State[tfconstants.AttrDiskSize])
	assert.NotNil(t, v1State[tfconstants.AttrVPCs])
	assert.NotNil(t, v1State[tfconstants.AttrDatabase])

	// New V1 fields should be added
	assert.NotNil(t, v1State[tfconstants.AttrTags])
}

func TestResourceMariaDBStateUpgradeV0toV1_WithMissingFields(t *testing.T) {
	v0State := map[string]interface{}{
		tfconstants.AttrID:   "12345",
		tfconstants.AttrName: "my-mariadb",
		// Missing many fields - should handle gracefully
	}

	v1State, err := dbaas_mariadb.ResourceMariaDBStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	// Existing fields should be preserved
	assert.Equal(t, "12345", v1State[tfconstants.AttrID])
	assert.Equal(t, "my-mariadb", v1State[tfconstants.AttrName])

	// New V1 fields should be added
	assert.NotNil(t, v1State[tfconstants.AttrTags])
	tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

func TestResourceMariaDBStateUpgradeV0toV1_EmptyState(t *testing.T) {
	v0State := map[string]interface{}{
		tfconstants.AttrID: "12345",
	}

	v1State, err := dbaas_mariadb.ResourceMariaDBStateUpgradeV0toV1(context.Background(), v0State, nil)

	require.NoError(t, err)
	assert.Equal(t, "12345", v1State[tfconstants.AttrID])
	assert.NotNil(t, v1State[tfconstants.AttrTags])
}

func TestResourceMariaDBStateUpgradeV0toV1_NoErrors(t *testing.T) {
	v0State := map[string]interface{}{
		tfconstants.AttrID:                  "12345",
		tfconstants.AttrName:                "my-mariadb",
		tfconstants.AttrSoftwareName:        "MariaDB",
		tfconstants.AttrSoftwareVersion:     "10.6",
		tfconstants.AttrGroup:               "Default",
		tfconstants.AttrPlan:                "DBS.16GB",
		tfconstants.AttrPublicIPEnabled:     true,
		tfconstants.AttrParameterGroupID:    tfconstants.DBaaSDefaultParameterGroupID,
		tfconstants.AttrIsEncryptionEnabled: tfconstants.DBaaSDefaultIsEncryptionEnabled,
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":         "admin",
				"password":     "secret123",
				"name":         "mydb",
				"dbaas_number": 1,
			},
		},
	}

	_, err := dbaas_mariadb.ResourceMariaDBStateUpgradeV0toV1(context.Background(), v0State, nil)
	assert.NoError(t, err, "State upgrade should return no errors for valid V0 state")
}
