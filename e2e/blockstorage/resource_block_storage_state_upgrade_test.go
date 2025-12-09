package blockstorage_test

import (
	"context"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/blockstorage"
	"github.com/stretchr/testify/assert"
)

func TestResourceBlockStorageStateUpgradeV0toV1_Basic(t *testing.T) {
	v0State := map[string]interface{}{
		"id":         "12345",
		"name":       "my-volume",
		"size":       250.0,
		"region":     "Delhi",
		"project_id": "789",
		"iops":       "3750",
		"status":     "available",
	}

	v1State, err := blockstorage.ResourceBlockStorageStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "my-volume", v1State["name"])
	assert.Equal(t, 250.0, v1State["size"])
	assert.NotNil(t, v1State["tags"], "tags field should be added")
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

func TestResourceBlockStorageStateUpgradeV0toV1_WithExistingTags(t *testing.T) {
	v0State := map[string]interface{}{
		"id":         "12345",
		"name":       "my-volume",
		"size":       250.0,
		"region":     "Delhi",
		"project_id": "789",
		"tags": map[string]interface{}{
			"Environment": "production",
		},
	}

	v1State, err := blockstorage.ResourceBlockStorageStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State["tags"])
	// Existing tags should be preserved
	tags := v1State["tags"].(map[string]interface{})
	assert.Equal(t, "production", tags["Environment"])
}

func TestResourceBlockStorageStateUpgradeV0toV1_PreservesAllFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":         "67890",
		"name":       "test-volume",
		"size":       500.0,
		"region":     "Mumbai",
		"project_id": "123",
		"iops":       "7500",
		"status":     "attached",
		"vm_id":      "999",
		"vm_name":    "test-vm",
	}

	v1State, err := blockstorage.ResourceBlockStorageStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "67890", v1State["id"])
	assert.Equal(t, "test-volume", v1State["name"])
	assert.Equal(t, 500.0, v1State["size"])
	assert.Equal(t, "Mumbai", v1State["region"])
	assert.Equal(t, "123", v1State["project_id"])
	assert.Equal(t, "7500", v1State["iops"])
	assert.Equal(t, "attached", v1State["status"])
	assert.Equal(t, "999", v1State["vm_id"])
	assert.Equal(t, "test-vm", v1State["vm_name"])
	// Tags should be added
	assert.NotNil(t, v1State["tags"])
}
