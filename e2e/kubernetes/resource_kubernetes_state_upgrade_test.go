package kubernetes_test

import (
	"context"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/kubernetes"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
)

func TestResourceKubernetesStateUpgradeV0toV1_Basic(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"name":    "my-cluster",
		"version": "1.30",
		"vpc_id":  "vpc-123",
		"region":  "Delhi",
		"status":  goe2econstants.KubernetesClusterStatusRunning,
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "my-cluster", v1State["name"])
	assert.Equal(t, "1.30", v1State["version"])
	assert.Equal(t, "vpc-123", v1State["vpc_id"])
	assert.Equal(t, "Delhi", v1State["region"])
	assert.Equal(t, goe2econstants.KubernetesClusterStatusRunning, v1State["status"])

	// New V3 fields initialized
	assert.NotNil(t, v1State["tags"], "tags field should be added")
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")

	assert.NotNil(t, v1State["encryption_enabled"], "encryption_enabled field should be added")
	assert.Equal(t, false, v1State["encryption_enabled"])

	assert.NotNil(t, v1State["encryption_passphrase"], "encryption_passphrase field should be added")
	assert.Equal(t, "", v1State["encryption_passphrase"])

	assert.NotNil(t, v1State["security_group_ids"], "security_group_ids field should be added")
	sgIDs, ok := v1State["security_group_ids"].([]interface{})
	assert.True(t, ok, "security_group_ids should be a list")
	assert.Empty(t, sgIDs, "security_group_ids should be empty list")
}

func TestResourceKubernetesStateUpgradeV0toV1_PreservesDeprecatedFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":       "12345",
		"name":     "old-cluster",
		"version":  "1.29",
		"location": "Delhi", // Deprecated field
		"vpc_id":   "vpc-456",
		"node_pools": []interface{}{
			map[string]interface{}{
				"name":           "pool1",
				"specs_name":     "C3.8GB",                                    // Deprecated
				"node_pool_type": goe2econstants.KubernetesNodePoolTypeStatic, // Deprecated
				"worker_node":    3,                                           // Deprecated
			},
		},
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)

	// Deprecated fields preserved (backward compatible)
	assert.Equal(t, "old-cluster", v1State["name"])
	assert.Equal(t, "1.29", v1State["version"])
	assert.Equal(t, "Delhi", v1State["location"])

	// Node pools preserved
	nodePools := v1State["node_pools"].([]interface{})
	assert.Len(t, nodePools, 1)

	pool := nodePools[0].(map[string]interface{})
	assert.Equal(t, "pool1", pool["name"])
	assert.Equal(t, "C3.8GB", pool["specs_name"])
	assert.Equal(t, goe2econstants.KubernetesNodePoolTypeStatic, pool["node_pool_type"])
	assert.Equal(t, 3, pool["worker_node"])

	// New V3 fields initialized
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["encryption_enabled"])
	assert.NotNil(t, v1State["security_group_ids"])
}

func TestResourceKubernetesStateUpgradeV0toV1_NodePools(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "67890",
		"name":    "test-cluster",
		"version": "1.30",
		"vpc_id":  "vpc-789",
		"node_pools": []interface{}{
			map[string]interface{}{
				"name":           "static-pool",
				"specs_name":     "C3.16GB",
				"node_pool_type": goe2econstants.KubernetesNodePoolTypeStatic,
				"worker_node":    5,
				"cardinality":    5,
			},
			map[string]interface{}{
				"name":           "autoscale-pool",
				"specs_name":     "C3.8GB",
				"node_pool_type": goe2econstants.KubernetesNodePoolTypeAutoscale,
				"min_vms":        2,
				"max_vms":        10,
				"cardinality":    3,
			},
		},
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)

	// Node pools preserved
	nodePools := v1State["node_pools"].([]interface{})
	assert.Len(t, nodePools, 2)

	// Static pool preserved
	staticPool := nodePools[0].(map[string]interface{})
	assert.Equal(t, "static-pool", staticPool["name"])
	assert.Equal(t, "C3.16GB", staticPool["specs_name"])
	assert.Equal(t, "Static", staticPool["node_pool_type"])
	assert.Equal(t, 5, staticPool["worker_node"])

	// Autoscale pool preserved
	autoscalePool := nodePools[1].(map[string]interface{})
	assert.Equal(t, "autoscale-pool", autoscalePool["name"])
	assert.Equal(t, "C3.8GB", autoscalePool["specs_name"])
	assert.Equal(t, goe2econstants.KubernetesNodePoolTypeAutoscale, autoscalePool["node_pool_type"])
	assert.Equal(t, 2, autoscalePool["min_vms"])
	assert.Equal(t, 10, autoscalePool["max_vms"])

	// New V3 fields initialized
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["encryption_enabled"])
	assert.NotNil(t, v1State["security_group_ids"])
}

func TestResourceKubernetesStateUpgradeV0toV1_NilState(t *testing.T) {
	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State)
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["encryption_enabled"])
	assert.NotNil(t, v1State["security_group_ids"])
}

func TestResourceKubernetesStateUpgradeV0toV1_AllFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":         "99999",
		"name":       "complete-cluster",
		"version":    "1.30",
		"vpc_id":     "vpc-999",
		"region":     "Mumbai",
		"location":   "Mumbai", // Deprecated
		"project_id": "123",
		"status":     "Running",
		"created_at": "2024-01-01T00:00:00Z",
		"slug_name":  "k8s-plan",
		"sku_id":     "sku-123",
		"node_pools": []interface{}{
			map[string]interface{}{
				"name":           "pool1",
				"specs_name":     "C3.8GB",
				"node_pool_type": goe2econstants.KubernetesNodePoolTypeStatic,
				"worker_node":    3,
				"cardinality":    3,
				"service_id":     "svc-123",
				"slug_name":      "c3-8gb",
				"sku_id":         "sku-pool-123",
			},
			map[string]interface{}{
				"name":           "pool2",
				"specs_name":     "C3.16GB",
				"node_pool_type": goe2econstants.KubernetesNodePoolTypeAutoscale,
				"min_vms":        2,
				"max_vms":        10,
				"cardinality":    5,
				"service_id":     "svc-456",
				"policy_type":    goe2econstants.KubernetesPolicyTypeDefault,
			},
		},
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)

	// All V0 fields preserved
	assert.Equal(t, "99999", v1State["id"])
	assert.Equal(t, "complete-cluster", v1State["name"])
	assert.Equal(t, "1.30", v1State["version"])
	assert.Equal(t, "vpc-999", v1State["vpc_id"])
	assert.Equal(t, "Mumbai", v1State["region"])
	assert.Equal(t, "Mumbai", v1State["location"])
	assert.Equal(t, "123", v1State["project_id"])
	assert.Equal(t, goe2econstants.KubernetesClusterStatusRunning, v1State["status"])
	assert.Equal(t, "2024-01-01T00:00:00Z", v1State["created_at"])
	assert.Equal(t, "k8s-plan", v1State["slug_name"])
	assert.Equal(t, "sku-123", v1State["sku_id"])

	// Node pools preserved with all fields
	nodePools := v1State["node_pools"].([]interface{})
	assert.Len(t, nodePools, 2)

	pool1 := nodePools[0].(map[string]interface{})
	assert.Equal(t, "pool1", pool1["name"])
	assert.Equal(t, "C3.8GB", pool1["specs_name"])
	assert.Equal(t, 3, pool1["worker_node"])
	assert.Equal(t, 3, pool1["cardinality"])
	assert.Equal(t, "svc-123", pool1["service_id"])

	pool2 := nodePools[1].(map[string]interface{})
	assert.Equal(t, "pool2", pool2["name"])
	assert.Equal(t, "C3.16GB", pool2["specs_name"])
	assert.Equal(t, 2, pool2["min_vms"])
	assert.Equal(t, 10, pool2["max_vms"])
	assert.Equal(t, goe2econstants.KubernetesPolicyTypeDefault, pool2["policy_type"])

	// New V3 fields initialized
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["encryption_enabled"])
	assert.Equal(t, false, v1State["encryption_enabled"])
	assert.NotNil(t, v1State["security_group_ids"])
}

func TestResourceKubernetesStateUpgradeV0toV1_WithExistingTags(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"name":    "my-cluster",
		"version": "1.30",
		"vpc_id":  "vpc-123",
		"tags": map[string]interface{}{
			"Environment": "production",
			"Team":        "backend",
		},
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// Existing tags should be preserved (though tags shouldn't exist in V0, this tests robustness)
	assert.NotNil(t, v1State["tags"])
	tags := v1State["tags"].(map[string]interface{})
	assert.Equal(t, "production", tags["Environment"])
	assert.Equal(t, "backend", tags["Team"])
}

func TestResourceKubernetesStateUpgradeV0toV1_WithElasticityDict(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"name":    "my-cluster",
		"version": "1.30",
		"vpc_id":  "vpc-123",
		"node_pools": []interface{}{
			map[string]interface{}{
				"name":           "autoscale-pool",
				"specs_name":     "C3.8GB",
				"node_pool_type": goe2econstants.KubernetesNodePoolTypeAutoscale,
				"min_vms":        2,
				"max_vms":        10,
				"elasticity_dict": []interface{}{
					map[string]interface{}{
						"worker": []interface{}{
							map[string]interface{}{
								"period_number":        1,
								"policy_paramter_type": goe2econstants.KubernetesPolicyTypeDefault,
								"parameter":            goe2econstants.KubernetesPolicyParameterCPU,
								"elasticity_policies": []interface{}{
									map[string]interface{}{
										"operator":     ">",
										"value":        80,
										"period":       60,
										"watch_period": 3,
										"cooldown":     300,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)

	// Elasticity dict preserved
	nodePools := v1State["node_pools"].([]interface{})
	assert.Len(t, nodePools, 1)

	pool := nodePools[0].(map[string]interface{})
	assert.NotNil(t, pool["elasticity_dict"])

	// New V3 fields initialized
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["encryption_enabled"])
	assert.NotNil(t, v1State["security_group_ids"])
}

func TestResourceKubernetesStateUpgradeV0toV1_WithScheduledDict(t *testing.T) {
	v0State := map[string]interface{}{
		"id":      "12345",
		"name":    "my-cluster",
		"version": "1.30",
		"vpc_id":  "vpc-123",
		"node_pools": []interface{}{
			map[string]interface{}{
				"name":           "autoscale-pool",
				"specs_name":     "C3.8GB",
				"node_pool_type": goe2econstants.KubernetesNodePoolTypeAutoscale,
				"min_vms":        2,
				"max_vms":        10,
				"scheduled_dict": []interface{}{
					map[string]interface{}{
						"worker": []interface{}{
							map[string]interface{}{
								"scheduled_policies": []interface{}{
									map[string]interface{}{
										"upscale_cardinality":   5,
										"upscale_recurrence":    "0 9 * * *",
										"downscale_cardinality": 2,
										"downscale_recurrence":  "0 2 * * *",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)

	// Scheduled dict preserved
	nodePools := v1State["node_pools"].([]interface{})
	assert.Len(t, nodePools, 1)

	pool := nodePools[0].(map[string]interface{})
	assert.NotNil(t, pool["scheduled_dict"])

	// New V3 fields initialized
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["encryption_enabled"])
	assert.NotNil(t, v1State["security_group_ids"])
}

func TestResourceKubernetesStateUpgradeV0toV1_WithSecurityGroups(t *testing.T) {
	v0State := map[string]interface{}{
		"id":                 "12345",
		"name":               "my-cluster",
		"version":            "1.30",
		"vpc_id":             "vpc-123",
		"security_group_ids": []interface{}{1, 2, 3},
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// Existing security_group_ids should be preserved
	assert.NotNil(t, v1State["security_group_ids"])
	sgIDs := v1State["security_group_ids"].([]interface{})
	assert.Len(t, sgIDs, 3)
	assert.Equal(t, 1, sgIDs[0])
	assert.Equal(t, 2, sgIDs[1])
	assert.Equal(t, 3, sgIDs[2])

	// Other V3 fields initialized
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["encryption_enabled"])
}

func TestResourceKubernetesStateUpgradeV0toV1_WithEncryptionFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":                    "12345",
		"name":                  "my-cluster",
		"version":               "1.30",
		"vpc_id":                "vpc-123",
		"encryption_enabled":    true,
		"encryption_passphrase": "secret-passphrase",
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// Existing encryption fields should be preserved
	assert.NotNil(t, v1State["encryption_enabled"])
	assert.Equal(t, true, v1State["encryption_enabled"])
	assert.NotNil(t, v1State["encryption_passphrase"])
	assert.Equal(t, "secret-passphrase", v1State["encryption_passphrase"])

	// Other V3 fields initialized
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["security_group_ids"])
}

func TestResourceKubernetesStateUpgradeV0toV1_WithExistingV3Fields(t *testing.T) {
	// Test that existing V3 fields are preserved (shouldn't happen in V0, but tests robustness)
	v0State := map[string]interface{}{
		"id":                 "12345",
		"name":               "my-cluster",
		"version":            "1.30",
		"vpc_id":             "vpc-123",
		"cluster_name":       "v3-cluster-name", // V3 field
		"kubernetes_version": "1.31",            // V3 field
		"tags": map[string]interface{}{
			"Environment": "test",
		},
		"security_group_ids": []interface{}{5, 6},
		"encryption_enabled": true,
	}

	v1State, err := kubernetes.ResourceKubernetesStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// All V3 fields should be preserved
	assert.Equal(t, "v3-cluster-name", v1State["cluster_name"])
	assert.Equal(t, "1.31", v1State["kubernetes_version"])
	tags := v1State["tags"].(map[string]interface{})
	assert.Equal(t, "test", tags["Environment"])
	sgIDs := v1State["security_group_ids"].([]interface{})
	assert.Len(t, sgIDs, 2)
	assert.Equal(t, true, v1State["encryption_enabled"])
}
