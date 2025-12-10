package autoscaling_test

import (
	"context"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/autoscaling"
	"github.com/stretchr/testify/assert"
)

func TestResourceAutoscalingStateUpgradeV0toV1_Basic(t *testing.T) {
	v0State := map[string]interface{}{
		"id":                    "12345",
		"name":                  "my-scaler-group",
		"plan":                  "small",
		"vm_image_name":         "ubuntu_20.04",
		"min_nodes":             1,
		"max_nodes":             5,
		"desired":               2,
		"is_encryption_enabled": false,
		"is_public_ip_required": true,
		"provision_status":      "Running",
		"region":                "Delhi",
		"project_id":            "789",
		"plan_id":               "100",
		"vm_image_id":           "200",
		"vm_template_id":        300,
		"running":               2,
	}

	v1State, err := autoscaling.ResourceAutoscalingStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "my-scaler-group", v1State["name"])
	assert.Equal(t, 1, v1State["min_nodes"])
	assert.Equal(t, 5, v1State["max_nodes"])
	assert.Equal(t, 2, v1State["desired"])
	assert.Equal(t, false, v1State["is_encryption_enabled"])
	assert.Equal(t, true, v1State["is_public_ip_required"])
	assert.Equal(t, "Running", v1State["provision_status"])
	assert.Equal(t, 2, v1State["running"])

	// New V3 fields should be added
	assert.NotNil(t, v1State["tags"], "tags field should be added")
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")

	assert.NotNil(t, v1State["running_node_count"], "running_node_count field should be added")
	assert.Equal(t, 2, v1State["running_node_count"], "running_node_count should match running")

	assert.NotNil(t, v1State["nodes"], "nodes field should be added")
	nodes, ok := v1State["nodes"].([]interface{})
	assert.True(t, ok, "nodes should be a list")
	assert.Empty(t, nodes, "nodes should be empty list")
}

func TestResourceAutoscalingStateUpgradeV0toV1_WithExistingTags(t *testing.T) {
	v0State := map[string]interface{}{
		"id":   "12345",
		"name": "my-scaler-group",
		"plan": "small",
		"tags": map[string]interface{}{
			"Environment": "production",
			"Team":        "backend",
		},
	}

	v1State, err := autoscaling.ResourceAutoscalingStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State["tags"])
	// Existing tags should be preserved
	tags := v1State["tags"].(map[string]interface{})
	assert.Equal(t, "production", tags["Environment"])
	assert.Equal(t, "backend", tags["Team"])
}

func TestResourceAutoscalingStateUpgradeV0toV1_WithPolicies(t *testing.T) {
	v0State := map[string]interface{}{
		"id":          "12345",
		"name":        "my-scaler-group",
		"plan":        "small",
		"policy_type": "elastic",
		"policy": []interface{}{
			map[string]interface{}{
				"type":           "upscale",
				"adjust":         2,
				"parameter":      "cpu",
				"operator":       ">",
				"value":          "80",
				"period_number":  "3",
				"period_seconds": "60",
				"cooldown":       "300",
			},
		},
		"scheduled_policy": []interface{}{
			map[string]interface{}{
				"type":       "upscale",
				"adjust":     "2",
				"recurrence": "0 9 * * *",
			},
		},
	}

	v1State, err := autoscaling.ResourceAutoscalingStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "elastic", v1State["policy_type"])
	assert.NotNil(t, v1State["policy"])
	assert.NotNil(t, v1State["scheduled_policy"])

	// New V3 fields should be added
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["running_node_count"])
	assert.NotNil(t, v1State["nodes"])
}

func TestResourceAutoscalingStateUpgradeV0toV1_PreservesAllFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":                    "67890",
		"name":                  "test-scaler-group",
		"plan":                  "medium",
		"vm_image_name":         "centos_8",
		"min_nodes":             2,
		"max_nodes":             10,
		"desired":               5,
		"is_encryption_enabled": true,
		"encryption_passphrase": "secret123",
		"is_public_ip_required": false,
		"provision_status":      "Stopped",
		"region":                "Mumbai",
		"project_id":            "123",
		"plan_id":               "200",
		"sku_id":                "200",
		"slug_name":             "medium-plan",
		"vm_image_id":           "300",
		"vm_template_id":        400,
		"my_account_sg_id":      50,
		"security_group_ids":    []interface{}{50, 51},
		"running":               5,
		"policy_type":           "scheduled",
		"vpc": []interface{}{
			map[string]interface{}{
				"name":       "my-vpc",
				"network_id": 100,
			},
		},
	}

	v1State, err := autoscaling.ResourceAutoscalingStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "67890", v1State["id"])
	assert.Equal(t, "test-scaler-group", v1State["name"])
	assert.Equal(t, "medium", v1State["plan"])
	assert.Equal(t, "centos_8", v1State["vm_image_name"])
	assert.Equal(t, 2, v1State["min_nodes"])
	assert.Equal(t, 10, v1State["max_nodes"])
	assert.Equal(t, 5, v1State["desired"])
	assert.Equal(t, true, v1State["is_encryption_enabled"])
	assert.Equal(t, "secret123", v1State["encryption_passphrase"])
	assert.Equal(t, false, v1State["is_public_ip_required"])
	assert.Equal(t, "Stopped", v1State["provision_status"])
	assert.Equal(t, "Mumbai", v1State["region"])
	assert.Equal(t, "123", v1State["project_id"])
	assert.Equal(t, 5, v1State["running"])
	assert.Equal(t, "scheduled", v1State["policy_type"])
	assert.NotNil(t, v1State["vpc"])

	// New V3 fields should be added
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["running_node_count"])
	assert.Equal(t, 5, v1State["running_node_count"])
	assert.NotNil(t, v1State["nodes"])
}

func TestResourceAutoscalingStateUpgradeV0toV1_NoRunningField(t *testing.T) {
	v0State := map[string]interface{}{
		"id":   "12345",
		"name": "my-scaler-group",
		"plan": "small",
	}

	v1State, err := autoscaling.ResourceAutoscalingStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// running_node_count should default to 0 if running is not present
	assert.NotNil(t, v1State["running_node_count"])
	assert.Equal(t, 0, v1State["running_node_count"])
}
