package container_registry_test

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/container_registry"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
)

// TestResourceContainerRegistryStateUpgradeV0toV1_Basic tests basic state upgrade adds tags field
func TestResourceContainerRegistryStateUpgradeV0toV1_Basic(t *testing.T) {
	v0State := map[string]interface{}{
		"id":            "12345",
		"project_name":  "test-registry",
		"prevent_vul":   false,
		"severity":      goe2econstants.ContainerRegistrySeverityLow,
		"status":        "active",
		"setup_status":  "active",
		"domain_name":   "test-registry.e2enetworks.net",
		"project_size":  1024.5,
		"storage_limit": 10737418240,
		"is_public":     false,
		"region":        "Mumbai",
		"project_id":    "777",
		"created_at":    "2025-12-01T10:00:00Z",
		"updated_at":    "2025-12-12T15:30:00Z",
	}

	v1State, err := container_registry.ResourceContainerRegistryStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State)

	// Verify all original fields are preserved
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "test-registry", v1State["project_name"])
	assert.Equal(t, false, v1State["prevent_vul"])
	assert.Equal(t, goe2econstants.ContainerRegistrySeverityLow, v1State["severity"])
	assert.Equal(t, "active", v1State["status"])
	assert.Equal(t, "active", v1State["setup_status"])
	assert.Equal(t, "test-registry.e2enetworks.net", v1State["domain_name"])
	assert.Equal(t, 1024.5, v1State["project_size"])
	assert.Equal(t, 10737418240, v1State["storage_limit"])
	assert.Equal(t, false, v1State["is_public"])
	assert.Equal(t, "Mumbai", v1State["region"])
	assert.Equal(t, "777", v1State["project_id"])

	// Verify tags field is added as empty map
	assert.NotNil(t, v1State[tfconstants.AttrTags], "tags field should be added")
	tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
	assert.True(t, ok, "tags should be a map[string]interface{}")
	assert.Empty(t, tags, "tags should be empty map")
}

// TestResourceContainerRegistryStateUpgradeV0toV1_WithExistingTags tests upgrade preserves existing tags
func TestResourceContainerRegistryStateUpgradeV0toV1_WithExistingTags(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "12345",
		"project_name": "test-registry",
		"prevent_vul":  true,
		"severity":     goe2econstants.ContainerRegistrySeverityHigh,
		"status":       "active",
		"region":       "Mumbai",
		"project_id":   "777",
		tfconstants.AttrTags: map[string]interface{}{
			"Environment": "production",
			"Team":        "platform",
			"Owner":       "devops",
		},
	}

	v1State, err := container_registry.ResourceContainerRegistryStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State)

	// Existing tags should be preserved
	tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Len(t, tags, 3, "should have 3 tags")
	assert.Equal(t, "production", tags["Environment"])
	assert.Equal(t, "platform", tags["Team"])
	assert.Equal(t, "devops", tags["Owner"])
}

// TestResourceContainerRegistryStateUpgradeV0toV1_PreservesAllFields tests all fields are preserved
func TestResourceContainerRegistryStateUpgradeV0toV1_PreservesAllFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":            "67890",
		"project_name":  "complex-registry",
		"prevent_vul":   true,
		"severity":      goe2econstants.ContainerRegistrySeverityCritical,
		"status":        "suspended",
		"setup_status":  "suspended",
		"domain_name":   "complex-registry.e2enetworks.net",
		"project_size":  5242880.75,
		"storage_limit": 21474836480,
		"is_public":     true,
		"region":        "Delhi",
		"location":      "Delhi",
		"project_id":    "123",
		"created_at":    "2025-01-01T00:00:00Z",
		"updated_at":    "2025-12-12T12:00:00Z",
	}

	v1State, err := container_registry.ResourceContainerRegistryStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State)

	// Verify all fields are preserved
	assert.Equal(t, "67890", v1State["id"])
	assert.Equal(t, "complex-registry", v1State["project_name"])
	assert.Equal(t, true, v1State["prevent_vul"])
	assert.Equal(t, goe2econstants.ContainerRegistrySeverityCritical, v1State["severity"])
	assert.Equal(t, "suspended", v1State["status"])
	assert.Equal(t, "suspended", v1State["setup_status"])
	assert.Equal(t, "complex-registry.e2enetworks.net", v1State["domain_name"])
	assert.Equal(t, 5242880.75, v1State["project_size"])
	assert.Equal(t, 21474836480, v1State["storage_limit"])
	assert.Equal(t, true, v1State["is_public"])
	assert.Equal(t, "Delhi", v1State["region"])
	assert.Equal(t, "Delhi", v1State["location"])
	assert.Equal(t, "123", v1State["project_id"])
	assert.Equal(t, "2025-01-01T00:00:00Z", v1State["created_at"])
	assert.Equal(t, "2025-12-12T12:00:00Z", v1State["updated_at"])

	// Verify tags field is added
	assert.NotNil(t, v1State[tfconstants.AttrTags])
}

// TestResourceContainerRegistryStateUpgradeV0toV1_EmptyState tests upgrade with minimal state
func TestResourceContainerRegistryStateUpgradeV0toV1_EmptyState(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "1",
		"project_name": "minimal",
	}

	v1State, err := container_registry.ResourceContainerRegistryStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State)

	// Verify minimal fields are preserved
	assert.Equal(t, "1", v1State["id"])
	assert.Equal(t, "minimal", v1State["project_name"])

	// Verify tags field is added
	assert.NotNil(t, v1State[tfconstants.AttrTags])
	tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty")
}

// TestResourceContainerRegistryStateUpgradeV0toV1_NilTags tests when tags field doesn't exist
func TestResourceContainerRegistryStateUpgradeV0toV1_NilTags(t *testing.T) {
	v0State := map[string]interface{}{
		"id":           "999",
		"project_name": "no-tags-registry",
		"prevent_vul":  false,
		"severity":     goe2econstants.ContainerRegistrySeverityMedium,
		"status":       "active",
		"region":       "Chennai",
		"project_id":   "456",
	}

	// Explicitly ensure tags field doesn't exist
	delete(v0State, tfconstants.AttrTags)

	v1State, err := container_registry.ResourceContainerRegistryStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State)

	// Verify all original fields preserved
	assert.Equal(t, "999", v1State["id"])
	assert.Equal(t, "no-tags-registry", v1State["project_name"])
	assert.Equal(t, false, v1State["prevent_vul"])
	assert.Equal(t, goe2econstants.ContainerRegistrySeverityMedium, v1State["severity"])

	// Verify tags field is added
	assert.NotNil(t, v1State[tfconstants.AttrTags], "tags field should be added when it doesn't exist")
	tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

// TestResourceContainerRegistryStateUpgradeV0toV1_AllSeverityLevels tests upgrade with all severity levels
func TestResourceContainerRegistryStateUpgradeV0toV1_AllSeverityLevels(t *testing.T) {
	severityLevels := []string{
		goe2econstants.ContainerRegistrySeverityLow,
		goe2econstants.ContainerRegistrySeverityMedium,
		goe2econstants.ContainerRegistrySeverityHigh,
		goe2econstants.ContainerRegistrySeverityCritical,
		goe2econstants.ContainerRegistrySeverityNone,
	}

	for _, severity := range severityLevels {
		t.Run("severity_"+severity, func(t *testing.T) {
			v0State := map[string]interface{}{
				"id":           "1",
				"project_name": "test-registry",
				"prevent_vul":  false,
				"severity":     severity,
				"status":       "active",
				"region":       "Mumbai",
				"project_id":   "123",
			}

			v1State, err := container_registry.ResourceContainerRegistryStateUpgradeV0toV1(context.Background(), v0State, nil)

			assert.NoError(t, err)
			assert.NotNil(t, v1State)
			assert.Equal(t, severity, v1State["severity"])
			assert.NotNil(t, v1State[tfconstants.AttrTags])
		})
	}
}

// TestResourceContainerRegistryStateUpgradeV0toV1_BooleanFields tests upgrade with boolean variations
func TestResourceContainerRegistryStateUpgradeV0toV1_BooleanFields(t *testing.T) {
	testCases := []struct {
		name       string
		preventVul bool
		isPublic   bool
	}{
		{
			name:       "both false",
			preventVul: false,
			isPublic:   false,
		},
		{
			name:       "prevent true, public false",
			preventVul: true,
			isPublic:   false,
		},
		{
			name:       "prevent false, public true",
			preventVul: false,
			isPublic:   true,
		},
		{
			name:       "both true",
			preventVul: true,
			isPublic:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v0State := map[string]interface{}{
				"id":           "1",
				"project_name": "test-registry",
				"prevent_vul":  tc.preventVul,
				"is_public":    tc.isPublic,
				"severity":     goe2econstants.ContainerRegistrySeverityLow,
				"status":       "active",
				"region":       "Mumbai",
				"project_id":   "123",
			}

			v1State, err := container_registry.ResourceContainerRegistryStateUpgradeV0toV1(context.Background(), v0State, nil)

			assert.NoError(t, err)
			assert.NotNil(t, v1State)
			assert.Equal(t, tc.preventVul, v1State["prevent_vul"])
			assert.Equal(t, tc.isPublic, v1State["is_public"])
			assert.NotNil(t, v1State[tfconstants.AttrTags])
		})
	}
}
