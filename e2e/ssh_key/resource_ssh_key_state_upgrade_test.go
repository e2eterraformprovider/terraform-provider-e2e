package ssh_key

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResourceSshKeyStateUpgradeV0toV1_Basic tests basic state upgrade
func TestResourceSshKeyStateUpgradeV0toV1_Basic(t *testing.T) {
	// Simulate old V0 state
	oldState := map[string]interface{}{
		"id":         "12345",
		"label":      "test-ssh-key",
		"ssh_key":    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
		"location":   "Mumbai",
		"created_at": "2025-01-01T00:00:00Z",
		"project_id": "proj-123",
	}

	// Run the upgrade function
	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify old V0 fields are preserved (including deprecated V2-style field names)
	assert.Equal(t, "12345", newState["id"])
	assert.Equal(t, "test-ssh-key", newState["label"])
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...", newState["ssh_key"])
	assert.Equal(t, "Mumbai", newState["location"])
	assert.Equal(t, "2025-01-01T00:00:00Z", newState["created_at"])
	assert.Equal(t, "proj-123", newState["project_id"])

	// Verify new V3 fields are initialized
	tags, exists := newState["tags"]
	require.True(t, exists, "tags field should exist in upgraded state")
	assert.IsType(t, map[string]interface{}{}, tags)
	tagsMap := tags.(map[string]interface{})
	assert.Empty(t, tagsMap, "tags should be initialized as empty map")

	// Verify no forced recreation - all original fields preserved
	assert.Equal(t, oldState["id"], newState["id"])
	assert.Equal(t, oldState["label"], newState["label"])
	assert.Equal(t, oldState["ssh_key"], newState["ssh_key"])
}

// TestResourceSshKeyStateUpgradeV0toV1_AllFields tests upgrade with all fields
func TestResourceSshKeyStateUpgradeV0toV1_AllFields(t *testing.T) {
	oldState := map[string]interface{}{
		"id":         "67890",
		"label":      "prod-ssh-key",
		"ssh_key":    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...",
		"location":   "Singapore",
		"created_at": "2025-01-15T10:30:00Z",
		"project_id": "proj-456",
		"region":     "Singapore", // Some V0 states might have region too
	}

	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify all original fields preserved
	assert.Equal(t, "67890", newState["id"])
	assert.Equal(t, "prod-ssh-key", newState["label"])
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...", newState["ssh_key"])
	assert.Equal(t, "Singapore", newState["location"])
	assert.Equal(t, "2025-01-15T10:30:00Z", newState["created_at"])
	assert.Equal(t, "proj-456", newState["project_id"])
	assert.Equal(t, "Singapore", newState["region"])

	// Verify new V3 fields
	tags, exists := newState["tags"]
	require.True(t, exists, "tags field should exist")
	assert.IsType(t, map[string]interface{}{}, tags)
	tagsMap := tags.(map[string]interface{})
	assert.Empty(t, tagsMap, "tags should be initialized as empty map")

	// Verify backward compatibility maintained
	assert.Equal(t, oldState["label"], newState["label"])
	assert.Equal(t, oldState["ssh_key"], newState["ssh_key"])
}

// TestResourceSshKeyStateUpgradeV0toV1_PreservesDeprecated tests that deprecated fields are preserved
func TestResourceSshKeyStateUpgradeV0toV1_PreservesDeprecated(t *testing.T) {
	oldState := map[string]interface{}{
		"id":         "11111",
		"label":      "legacy-ssh-key",
		"ssh_key":    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ...",
		"location":   "Mumbai",
		"created_at": "2025-01-10T12:00:00Z",
		"project_id": "proj-789",
	}

	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify deprecated V2-style field names (label, ssh_key, location) still work after upgrade
	assert.Equal(t, "legacy-ssh-key", newState["label"])
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ...", newState["ssh_key"])
	assert.Equal(t, "Mumbai", newState["location"])

	// Verify no data loss during upgrade
	assert.Equal(t, oldState["id"], newState["id"])
	assert.Equal(t, oldState["label"], newState["label"])
	assert.Equal(t, oldState["ssh_key"], newState["ssh_key"])
	assert.Equal(t, oldState["location"], newState["location"])
	assert.Equal(t, oldState["created_at"], newState["created_at"])
	assert.Equal(t, oldState["project_id"], newState["project_id"])

	// Verify tags initialized
	tags, exists := newState["tags"]
	require.True(t, exists)
	assert.IsType(t, map[string]interface{}{}, tags)
}

// TestResourceSshKeyStateUpgradeV0toV1_ComputedFields tests computed fields initialization
func TestResourceSshKeyStateUpgradeV0toV1_ComputedFields(t *testing.T) {
	// Test with minimal state (only required fields)
	oldState := map[string]interface{}{
		"id":      "22222",
		"label":   "minimal-ssh-key",
		"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ...",
	}

	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify tags field initialized as empty map when missing
	tags, exists := newState["tags"]
	require.True(t, exists, "tags field should be initialized even when missing in V0 state")
	assert.IsType(t, map[string]interface{}{}, tags)
	tagsMap := tags.(map[string]interface{})
	assert.Empty(t, tagsMap, "tags should be initialized as empty map")

	// Verify required fields still present
	assert.Equal(t, "22222", newState["id"])
	assert.Equal(t, "minimal-ssh-key", newState["label"])
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ...", newState["ssh_key"])
}

// TestResourceSshKeyStateUpgradeV0toV1_NilState tests upgrade with nil state
func TestResourceSshKeyStateUpgradeV0toV1_NilState(t *testing.T) {
	// The function should handle nil state gracefully (defensive programming)
	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), nil, nil)
	require.NoError(t, err)
	require.NotNil(t, newState)

	// Verify tags are initialized even for nil state
	tags, exists := newState["tags"]
	require.True(t, exists, "tags should be initialized even for nil state")
	assert.IsType(t, map[string]interface{}{}, tags)
	tagsMap := tags.(map[string]interface{})
	assert.Empty(t, tagsMap, "tags should be initialized as empty map")
}

// TestResourceSshKeyStateUpgradeV0toV1_EmptyState tests upgrade with empty state
func TestResourceSshKeyStateUpgradeV0toV1_EmptyState(t *testing.T) {
	oldState := map[string]interface{}{}

	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)
	require.NotNil(t, newState)

	// Verify new fields are initialized with defaults
	tags, exists := newState["tags"]
	require.True(t, exists, "tags should be initialized even in empty state")
	assert.IsType(t, map[string]interface{}{}, tags)
	tagsMap := tags.(map[string]interface{})
	assert.Empty(t, tagsMap, "tags should be initialized as empty map")
}

// TestResourceSshKeyStateUpgradeV0toV1_MissingOptionalFields tests upgrade with missing optional fields
func TestResourceSshKeyStateUpgradeV0toV1_MissingOptionalFields(t *testing.T) {
	// State with only id, label, ssh_key (minimal required)
	oldState := map[string]interface{}{
		"id":      "33333",
		"label":   "minimal-key",
		"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ...",
		// Missing: location, created_at, project_id
	}

	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify required fields preserved
	assert.Equal(t, "33333", newState["id"])
	assert.Equal(t, "minimal-key", newState["label"])
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQ...", newState["ssh_key"])

	// Verify tags initialized
	tags, exists := newState["tags"]
	require.True(t, exists)
	assert.IsType(t, map[string]interface{}{}, tags)
}

// TestResourceSshKeyStateUpgradeV0toV1_IntegrationWithResource tests that upgraded state works with resource
func TestResourceSshKeyStateUpgradeV0toV1_IntegrationWithResource(t *testing.T) {
	// This test verifies that the upgraded state can be used with the resource schema
	oldState := map[string]interface{}{
		"id":         "integration-test",
		"label":      "test-ssh-key",
		"ssh_key":    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
		"location":   "Mumbai",
		"created_at": "2025-01-01T00:00:00Z",
		"project_id": "proj-123",
	}

	// Run the upgrade function
	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify all required fields exist in upgraded state
	requiredFields := []string{"id", "label", "ssh_key", "tags"}
	for _, field := range requiredFields {
		_, exists := newState[field]
		assert.True(t, exists, "Field %s should exist in upgraded state", field)
	}

	// Verify the upgraded state is compatible with current resource schema
	resource := ResourceSshKey()
	require.NotNil(t, resource)
	require.NotNil(t, resource.Schema)

	// Verify that all fields in upgraded state are valid schema fields
	for field := range newState {
		// Skip fields that might be in state but not in schema (like internal Terraform fields)
		if field == "id" || field == "terraform_version" || field == "schema_version" {
			continue
		}
		// Verify field exists in schema or is a known valid field
		_, exists := resource.Schema[field]
		if !exists {
			// Some fields like location might be deprecated but still valid
			validFields := []string{"label", "ssh_key", "location", "tags", "created_at", "project_id", "region", "name", "public_key"}
			assert.Contains(t, validFields, field, "Field %s should be a valid schema field", field)
		}
	}
}

// TestResourceSshKeyStateUpgradeV0toV1_NoForcedRecreation tests that upgrade doesn't force recreation
func TestResourceSshKeyStateUpgradeV0toV1_NoForcedRecreation(t *testing.T) {
	// Create a copy of the state to preserve original for comparison
	originalState := map[string]interface{}{
		"id":         "no-recreate-test",
		"label":      "test-key",
		"ssh_key":    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
		"location":   "Mumbai",
		"created_at": "2025-01-01T00:00:00Z",
		"project_id": "proj-123",
	}

	// Copy the state since the function modifies it in place
	oldState := make(map[string]interface{})
	for k, v := range originalState {
		oldState[k] = v
	}
	originalFieldCount := len(oldState)

	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify ID is unchanged (no forced recreation)
	assert.Equal(t, originalState["id"], newState["id"], "ID should not change during upgrade")

	// Verify all original values preserved exactly
	assert.Equal(t, originalState["label"], newState["label"])
	assert.Equal(t, originalState["ssh_key"], newState["ssh_key"])
	assert.Equal(t, originalState["location"], newState["location"])
	assert.Equal(t, originalState["created_at"], newState["created_at"])
	assert.Equal(t, originalState["project_id"], newState["project_id"])

	// Verify tags field was added
	_, tagsExists := newState["tags"]
	assert.True(t, tagsExists, "tags field should be added during upgrade")

	// Verify exactly one new field was added (tags)
	assert.Equal(t, originalFieldCount+1, len(newState), "New state should have exactly one more field than old state (tags)")
}

// TestResourceSshKeyV0Schema tests that the V0 schema is correctly defined
func TestResourceSshKeyV0Schema(t *testing.T) {
	schema := resourceSshKeyResourceV0()
	require.NotNil(t, schema)

	// Verify schema structure exists
	assert.NotNil(t, schema.Schema)

	// Verify V0 has old field names (deprecated V2-style fields)
	assert.Contains(t, schema.Schema, tfconstants.AttrLabel, "V0 schema should contain 'label' field")
	assert.Contains(t, schema.Schema, tfconstants.AttrSSHKey, "V0 schema should contain 'ssh_key' field")
	assert.Contains(t, schema.Schema, tfconstants.AttrLocation, "V0 schema should contain 'location' field")
	assert.Contains(t, schema.Schema, tfconstants.AttrCreatedAt, "V0 schema should contain 'created_at' field")
	assert.Contains(t, schema.Schema, tfconstants.AttrProjectID, "V0 schema should contain 'project_id' field")

	// Verify V0 does NOT have new V3 fields (name, public_key, tags)
	// These are added in V1 schema
	assert.NotContains(t, schema.Schema, tfconstants.AttrName, "V0 schema should NOT contain 'name' field")
	assert.NotContains(t, schema.Schema, tfconstants.AttrPublicKey, "V0 schema should NOT contain 'public_key' field")
	assert.NotContains(t, schema.Schema, tfconstants.AttrTags, "V0 schema should NOT contain 'tags' field")
}

// TestStateUpgraderType tests the StateUpgrader configuration
func TestStateUpgraderType(t *testing.T) {
	resource := ResourceSshKey()
	require.NotNil(t, resource)

	// Verify StateUpgraders is configured
	assert.Greater(t, len(resource.StateUpgraders), 0, "StateUpgraders should be configured")

	// Verify upgrader points to correct upgrade function
	upgrader := resource.StateUpgraders[0]
	assert.NotNil(t, upgrader.Upgrade, "StateUpgrader should have an Upgrade function")
	assert.Equal(t, 0, upgrader.Version, "StateUpgrader should upgrade from version 0")

	// Verify SchemaVersion is 1 (current version)
	assert.Equal(t, 1, resource.SchemaVersion, "Current schema version should be 1")
}

// TestResourceSshKeyStateUpgradeV0toV1_TagsAlreadyPresent tests upgrade when tags already exist
func TestResourceSshKeyStateUpgradeV0toV1_TagsAlreadyPresent(t *testing.T) {
	// Edge case: V0 state that somehow already has tags (shouldn't happen, but test robustness)
	oldState := map[string]interface{}{
		"id":      "tags-test",
		"label":   "test-key",
		"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
		"tags": map[string]interface{}{
			"existing": "tag",
		},
	}

	newState, err := resourceSshKeyStateUpgradeV0toV1(context.Background(), oldState, nil)
	require.NoError(t, err)

	// Verify existing tags are preserved (not overwritten)
	tags, exists := newState["tags"]
	require.True(t, exists)
	assert.IsType(t, map[string]interface{}{}, tags)
	tagsMap := tags.(map[string]interface{})
	assert.Equal(t, "tag", tagsMap["existing"], "Existing tags should be preserved")
}
