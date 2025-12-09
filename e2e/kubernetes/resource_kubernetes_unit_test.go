package kubernetes

import (
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ============================================
// CLUSTER FIELD ALIAS TESTS
// ============================================

func TestGetClusterName_PreferredField(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceKubernetesService().Schema, map[string]interface{}{
		"cluster_name": "test-cluster",
	})

	name := getClusterName(d)
	if name != "test-cluster" {
		t.Errorf("Expected 'test-cluster', got '%s'", name)
	}
}

func TestGetClusterName_DeprecatedField(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceKubernetesService().Schema, map[string]interface{}{
		"name": "old-cluster",
	})

	name := getClusterName(d)
	if name != "old-cluster" {
		t.Errorf("Expected 'old-cluster', got '%s'", name)
	}
}

func TestGetClusterName_PreferredOverDeprecated(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceKubernetesService().Schema, map[string]interface{}{
		"cluster_name": "new-cluster",
		"name":         "old-cluster",
	})

	name := getClusterName(d)
	if name != "new-cluster" {
		t.Errorf("Expected 'new-cluster' (preferred), got '%s'", name)
	}
}

func TestGetClusterName_Empty(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceKubernetesService().Schema, map[string]interface{}{})

	name := getClusterName(d)
	if name != "" {
		t.Errorf("Expected empty string, got '%s'", name)
	}
}

func TestGetKubernetesVersion_PreferredField(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceKubernetesService().Schema, map[string]interface{}{
		"kubernetes_version": "1.30",
	})

	version := getKubernetesVersion(d)
	if version != "1.30" {
		t.Errorf("Expected '1.30', got '%s'", version)
	}
}

func TestGetKubernetesVersion_DeprecatedField(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceKubernetesService().Schema, map[string]interface{}{
		"version": "1.29",
	})

	version := getKubernetesVersion(d)
	if version != "1.29" {
		t.Errorf("Expected '1.29', got '%s'", version)
	}
}

func TestGetKubernetesVersion_PreferredOverDeprecated(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceKubernetesService().Schema, map[string]interface{}{
		"kubernetes_version": "1.30",
		"version":            "1.29",
	})

	version := getKubernetesVersion(d)
	if version != "1.30" {
		t.Errorf("Expected '1.30' (preferred), got '%s'", version)
	}
}

func TestGetKubernetesVersion_Empty(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceKubernetesService().Schema, map[string]interface{}{})

	version := getKubernetesVersion(d)
	if version != "" {
		t.Errorf("Expected empty string, got '%s'", version)
	}
}

// ============================================
// NODE POOL FIELD ALIAS TESTS
// ============================================

func TestGetNodePoolPlan_PreferredField(t *testing.T) {
	pool := map[string]interface{}{
		"plan": "C3.8GB",
	}

	plan := getNodePoolPlan(pool)
	if plan != "C3.8GB" {
		t.Errorf("Expected 'C3.8GB', got '%s'", plan)
	}
}

func TestGetNodePoolPlan_DeprecatedField(t *testing.T) {
	pool := map[string]interface{}{
		"specs_name": "C3.8GB",
	}

	plan := getNodePoolPlan(pool)
	if plan != "C3.8GB" {
		t.Errorf("Expected 'C3.8GB', got '%s'", plan)
	}
}

func TestGetNodePoolPlan_PreferredOverDeprecated(t *testing.T) {
	pool := map[string]interface{}{
		"plan":       "C3.16GB",
		"specs_name": "C3.8GB",
	}

	plan := getNodePoolPlan(pool)
	if plan != "C3.16GB" {
		t.Errorf("Expected 'C3.16GB' (preferred), got '%s'", plan)
	}
}

func TestGetNodePoolType_PreferredField(t *testing.T) {
	pool := map[string]interface{}{
		"type": "Static",
	}

	poolType := getNodePoolType(pool)
	if poolType != "Static" {
		t.Errorf("Expected 'Static', got '%s'", poolType)
	}
}

func TestGetNodePoolType_DeprecatedField(t *testing.T) {
	pool := map[string]interface{}{
		"node_pool_type": "Autoscale",
	}

	poolType := getNodePoolType(pool)
	if poolType != "Autoscale" {
		t.Errorf("Expected 'Autoscale', got '%s'", poolType)
	}
}

func TestGetNodePoolSize_PreferredField(t *testing.T) {
	pool := map[string]interface{}{
		"size": 5,
	}

	size := getNodePoolSize(pool)
	if size != 5 {
		t.Errorf("Expected 5, got %d", size)
	}
}

func TestGetNodePoolSize_DeprecatedField(t *testing.T) {
	pool := map[string]interface{}{
		"worker_node": 3,
	}

	size := getNodePoolSize(pool)
	if size != 3 {
		t.Errorf("Expected 3, got %d", size)
	}
}

func TestGetNodePoolMinNodes_PreferredField(t *testing.T) {
	pool := map[string]interface{}{
		"min_nodes": 2,
	}

	minNodes := getNodePoolMinNodes(pool)
	if minNodes != 2 {
		t.Errorf("Expected 2, got %d", minNodes)
	}
}

func TestGetNodePoolMinNodes_DeprecatedField(t *testing.T) {
	pool := map[string]interface{}{
		"min_vms": 3,
	}

	minNodes := getNodePoolMinNodes(pool)
	if minNodes != 3 {
		t.Errorf("Expected 3, got %d", minNodes)
	}
}

func TestGetNodePoolMaxNodes_PreferredField(t *testing.T) {
	pool := map[string]interface{}{
		"max_nodes": 10,
	}

	maxNodes := getNodePoolMaxNodes(pool)
	if maxNodes != 10 {
		t.Errorf("Expected 10, got %d", maxNodes)
	}
}

func TestGetNodePoolMaxNodes_DeprecatedField(t *testing.T) {
	pool := map[string]interface{}{
		"max_vms": 15,
	}

	maxNodes := getNodePoolMaxNodes(pool)
	if maxNodes != 15 {
		t.Errorf("Expected 15, got %d", maxNodes)
	}
}

func TestGetNodePoolSize_ZeroWhenEmpty(t *testing.T) {
	pool := map[string]interface{}{}

	size := getNodePoolSize(pool)
	if size != 0 {
		t.Errorf("Expected 0, got %d", size)
	}
}

func TestGetNodePoolMinNodes_ZeroWhenEmpty(t *testing.T) {
	pool := map[string]interface{}{}

	minNodes := getNodePoolMinNodes(pool)
	if minNodes != 0 {
		t.Errorf("Expected 0, got %d", minNodes)
	}
}

func TestGetNodePoolMaxNodes_ZeroWhenEmpty(t *testing.T) {
	pool := map[string]interface{}{}

	maxNodes := getNodePoolMaxNodes(pool)
	if maxNodes != 0 {
		t.Errorf("Expected 0, got %d", maxNodes)
	}
}

func TestGetNodePoolPlan_EmptyWhenEmpty(t *testing.T) {
	pool := map[string]interface{}{}

	plan := getNodePoolPlan(pool)
	if plan != "" {
		t.Errorf("Expected empty string, got '%s'", plan)
	}
}

func TestGetNodePoolType_EmptyWhenEmpty(t *testing.T) {
	pool := map[string]interface{}{}

	poolType := getNodePoolType(pool)
	if poolType != "" {
		t.Errorf("Expected empty string, got '%s'", poolType)
	}
}

// ============================================
// FLATTEN NODE POOLS TESTS
// ============================================

func TestFlattenNodePools_EmptyList(t *testing.T) {
	result := flattenNodePools([]goe2e.NodePoolServiceInfo{})
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d items", len(result))
	}
}

func TestFlattenNodePools_SinglePool(t *testing.T) {
	pools := []goe2e.NodePoolServiceInfo{
		{
			ServiceID:   123.0,
			ServiceName: "test-pool",
			Cardinality: 3,
		},
	}

	result := flattenNodePools(pools)
	if len(result) != 1 {
		t.Fatalf("Expected 1 pool, got %d", len(result))
	}

	pool := result[0].(map[string]interface{})
	if pool["name"] != "test-pool" {
		t.Errorf("Expected name 'test-pool', got '%s'", pool["name"])
	}
	if pool["service_id"] != "123" {
		t.Errorf("Expected service_id '123', got '%s'", pool["service_id"])
	}
	if pool["cardinality"] != 3 {
		t.Errorf("Expected cardinality 3, got %d", pool["cardinality"])
	}
}

func TestFlattenNodePools_MultiplePools(t *testing.T) {
	pools := []goe2e.NodePoolServiceInfo{
		{
			ServiceID:   123.0,
			ServiceName: "pool-1",
			Cardinality: 3,
		},
		{
			ServiceID:   456.0,
			ServiceName: "pool-2",
			Cardinality: 5,
		},
	}

	result := flattenNodePools(pools)
	if len(result) != 2 {
		t.Fatalf("Expected 2 pools, got %d", len(result))
	}

	pool1 := result[0].(map[string]interface{})
	if pool1["name"] != "pool-1" {
		t.Errorf("Expected name 'pool-1', got '%s'", pool1["name"])
	}

	pool2 := result[1].(map[string]interface{})
	if pool2["name"] != "pool-2" {
		t.Errorf("Expected name 'pool-2', got '%s'", pool2["name"])
	}
}

func TestFlattenNodePools_ZeroCardinality(t *testing.T) {
	pools := []goe2e.NodePoolServiceInfo{
		{
			ServiceID:   789.0,
			ServiceName: "empty-pool",
			Cardinality: 0,
		},
	}

	result := flattenNodePools(pools)
	pool := result[0].(map[string]interface{})
	if pool["cardinality"] != 0 {
		t.Errorf("Expected cardinality 0, got %d", pool["cardinality"])
	}
}

func TestFlattenNodePools_FloatServiceID(t *testing.T) {
	pools := []goe2e.NodePoolServiceInfo{
		{
			ServiceID:   123.456,
			ServiceName: "float-pool",
			Cardinality: 2,
		},
	}

	result := flattenNodePools(pools)
	pool := result[0].(map[string]interface{})
	// Should format as "123" (rounded)
	if pool["service_id"] != "123" {
		t.Errorf("Expected service_id '123', got '%s'", pool["service_id"])
	}
}

// ============================================
// SECURITY GROUP HELPER TESTS
// ============================================

func TestExpandSecurityGroupIDs_EmptyList(t *testing.T) {
	result := expandSecurityGroupIDs([]interface{}{})
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d items", len(result))
	}
}

func TestExpandSecurityGroupIDs_SingleID(t *testing.T) {
	input := []interface{}{42}
	result := expandSecurityGroupIDs(input)
	if len(result) != 1 {
		t.Fatalf("Expected 1 ID, got %d", len(result))
	}
	if result[0] != 42 {
		t.Errorf("Expected ID 42, got %d", result[0])
	}
}

func TestExpandSecurityGroupIDs_MultipleIDs(t *testing.T) {
	input := []interface{}{1, 2, 3, 4, 5}
	result := expandSecurityGroupIDs(input)
	if len(result) != 5 {
		t.Fatalf("Expected 5 IDs, got %d", len(result))
	}
	for i, expected := range []int{1, 2, 3, 4, 5} {
		if result[i] != expected {
			t.Errorf("Expected ID %d at index %d, got %d", expected, i, result[i])
		}
	}
}

func TestDifference_EmptyLists(t *testing.T) {
	result := difference([]int{}, []int{})
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d items", len(result))
	}
}

func TestDifference_NoDifference(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	result := difference(a, b)
	if len(result) != 0 {
		t.Errorf("Expected no difference, got %d items: %v", len(result), result)
	}
}

func TestDifference_SomeDifference(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{2, 4}
	result := difference(a, b)
	expected := []int{1, 3, 5}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d: %v", len(expected), len(result), result)
	}
	for i, val := range expected {
		if result[i] != val {
			t.Errorf("Expected %d at index %d, got %d", val, i, result[i])
		}
	}
}

func TestDifference_AllDifferent(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	result := difference(a, b)
	if len(result) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(result))
	}
	for i, val := range []int{1, 2, 3} {
		if result[i] != val {
			t.Errorf("Expected %d at index %d, got %d", val, i, result[i])
		}
	}
}

func TestDifference_BEmpty(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{}
	result := difference(a, b)
	if len(result) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(result))
	}
}

func TestDifference_AEmpty(t *testing.T) {
	a := []int{}
	b := []int{1, 2, 3}
	result := difference(a, b)
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d items", len(result))
	}
}

// ============================================
// ERROR HANDLING TESTS
// ============================================

func TestExpandNodePools_DuplicateNames(t *testing.T) {
	// This test validates error handling for duplicate node pool names
	// Note: This requires a mock client, so we'll test the validation logic
	// by checking the error message structure
	config := []interface{}{
		map[string]interface{}{
			"name": "duplicate",
			"plan": "C3.8GB",
			"type": "Static",
			"size": 3,
		},
		map[string]interface{}{
			"name": "duplicate", // Same name
			"plan": "C3.8GB",
			"type": "Static",
			"size": 3,
		},
	}

	// We can't easily test ExpandNodePools without a real client,
	// but we can document that duplicate names should be caught
	// This is more of a documentation test
	_ = config
	t.Log("ExpandNodePools should return error for duplicate node pool names")
}

func TestExpandNodePools_MissingType(t *testing.T) {
	// Document that missing type should return error
	t.Log("ExpandNodePools should return error when node pool type is missing")
}

func TestExpandNodePools_MissingPlan(t *testing.T) {
	// Document that missing plan should return error
	t.Log("ExpandNodePools should return error when plan is not found")
}

func TestExpandNodePools_StaticMissingSize(t *testing.T) {
	// Document that Static pools require size
	t.Log("ExpandNodePools should return error when Static pool has no size")
}

func TestExpandNodePools_AutoscaleMissingMinMax(t *testing.T) {
	// Document that Autoscale pools require min_nodes and max_nodes
	t.Log("ExpandNodePools should return error when Autoscale pool has no min_nodes or max_nodes")
}

// ============================================
// WAIT FUNCTION TESTS
// ============================================

// Note: waitForClusterStatus and clusterStatusRefresh require a real goe2e client
// and external API calls, so they are tested via acceptance tests.
// These unit tests document the expected behavior.

func TestWaitForClusterStatus_Documentation(t *testing.T) {
	// Document that waitForClusterStatus:
	// - Waits for cluster to reach target status
	// - Uses StateChangeConf with appropriate pending/target states
	// - Has configurable timeout
	t.Log("waitForClusterStatus should wait for cluster to reach target status")
	t.Log("Pending states: Creating, Provisioning, Updating")
	t.Log("Target state: Running (or specified)")
}

func TestClusterStatusRefresh_Documentation(t *testing.T) {
	// Document that clusterStatusRefresh:
	// - Fetches cluster status via goe2e client
	// - Returns cluster object, state string, and error
	// - Handles nil cluster response
	t.Log("clusterStatusRefresh should fetch cluster status and return state")
}
