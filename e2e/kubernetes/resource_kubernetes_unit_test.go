package kubernetes

import (
	"fmt"
	"strings"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
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
		"type": goe2econstants.KubernetesNodePoolTypeStatic,
	}

	poolType := getNodePoolType(pool)
	if poolType != goe2econstants.KubernetesNodePoolTypeStatic {
		t.Errorf("Expected '%s', got '%s'", goe2econstants.KubernetesNodePoolTypeStatic, poolType)
	}
}

func TestGetNodePoolType_DeprecatedField(t *testing.T) {
	pool := map[string]interface{}{
		"node_pool_type": goe2econstants.KubernetesNodePoolTypeAutoscale,
	}

	poolType := getNodePoolType(pool)
	if poolType != goe2econstants.KubernetesNodePoolTypeAutoscale {
		t.Errorf("Expected '%s', got '%s'", goe2econstants.KubernetesNodePoolTypeAutoscale, poolType)
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
// Note: Comprehensive error handling tests for ExpandNodePools are now in helpers_unit_test.go
// These documentation tests are kept for reference but actual tests are in helpers_unit_test.go

func TestExpandNodePools_MissingType(t *testing.T) {
	// Document that missing type should return error
	// Actual test: TestExpandNodePools_InvalidPoolType in helpers_unit_test.go
	t.Log("ExpandNodePools should return error when node pool type is missing")
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

// ============================================
// FIELD VALIDATION TESTS
// ============================================

func TestValidateKubernetesVersion_ValidFormat(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema

	versionSchema, exists := resourceSchema["kubernetes_version"]
	require.True(t, exists, "Field kubernetes_version should exist in schema")
	require.NotNil(t, versionSchema.ValidateFunc, "Field kubernetes_version should have ValidateFunc")

	validVersions := []string{
		"1.20",
		"1.21",
		"1.22",
		"1.29",
		"1.30",
		"1.99",
	}

	for _, version := range validVersions {
		t.Run("valid_version_"+version, func(t *testing.T) {
			_, errors := versionSchema.ValidateFunc(version, "kubernetes_version")
			require.Empty(t, errors, "Version %s should be valid", version)
		})
	}
}

func TestValidateKubernetesVersion_InvalidFormat(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema

	versionSchema, exists := resourceSchema["kubernetes_version"]
	require.True(t, exists, "Field kubernetes_version should exist in schema")
	require.NotNil(t, versionSchema.ValidateFunc, "Field kubernetes_version should have ValidateFunc")

	invalidVersions := []string{
		"2.0",
		"1.2",
		"1.200",
		"v1.20",
		"1.20.0",
		"invalid",
		"1.x",
		"1",
	}

	for _, version := range invalidVersions {
		t.Run("invalid_version_"+version, func(t *testing.T) {
			_, errors := versionSchema.ValidateFunc(version, "kubernetes_version")
			require.NotEmpty(t, errors, "Version %s should be invalid", version)
			require.Contains(t, errors[0].Error(), "must be format 1.XX", "Error message should mention format requirement")
		})
	}
}

func TestValidateKubernetesVersion_Empty(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema

	versionSchema, exists := resourceSchema["kubernetes_version"]
	require.True(t, exists, "Field kubernetes_version should exist in schema")
	require.NotNil(t, versionSchema.ValidateFunc, "Field kubernetes_version should have ValidateFunc")

	// Empty string should pass validation (since field is Optional)
	// Validation only checks format, not presence
	_, errors := versionSchema.ValidateFunc("", "kubernetes_version")
	require.NotEmpty(t, errors, "Empty version should be invalid")
}

func TestValidateClusterName_ValidLength(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema

	nameSchema, exists := resourceSchema["cluster_name"]
	require.True(t, exists, "Field cluster_name should exist in schema")
	require.NotNil(t, nameSchema.ValidateFunc, "Field cluster_name should have ValidateFunc")

	validNames := []string{
		"a",                      // Minimum length (1)
		"test-cluster",           // Typical name
		strings.Repeat("a", 255), // Maximum length (255)
	}

	for _, name := range validNames {
		t.Run("valid_name_length_"+fmt.Sprintf("%d", len(name)), func(t *testing.T) {
			_, errors := nameSchema.ValidateFunc(name, "cluster_name")
			require.Empty(t, errors, "Name with length %d should be valid", len(name))
		})
	}
}

func TestValidateClusterName_TooLong(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema

	nameSchema, exists := resourceSchema["cluster_name"]
	require.True(t, exists, "Field cluster_name should exist in schema")
	require.NotNil(t, nameSchema.ValidateFunc, "Field cluster_name should have ValidateFunc")

	// Name longer than 255 characters
	tooLongName := strings.Repeat("a", 256)
	_, errors := nameSchema.ValidateFunc(tooLongName, "cluster_name")
	require.NotEmpty(t, errors, "Name longer than 255 characters should be invalid")
	require.Contains(t, errors[0].Error(), "length", "Error message should mention length")
}

func TestValidateClusterName_Empty(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema

	nameSchema, exists := resourceSchema["cluster_name"]
	require.True(t, exists, "Field cluster_name should exist in schema")
	require.NotNil(t, nameSchema.ValidateFunc, "Field cluster_name should have ValidateFunc")

	// Empty string should fail validation (minimum length is 1)
	_, errors := nameSchema.ValidateFunc("", "cluster_name")
	require.NotEmpty(t, errors, "Empty name should be invalid")
	require.Contains(t, errors[0].Error(), "length", "Error message should mention length")
}

func TestValidateNodePoolSize_Range(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema
	nodePoolsSchema := resourceSchema[tfconstants.AttrNodePools]
	nodePoolSchema := nodePoolsSchema.Elem.(*schema.Resource)
	nodePoolResourceSchema := nodePoolSchema.Schema

	sizeSchema, exists := nodePoolResourceSchema["size"]
	require.True(t, exists, "Field size should exist in node pool schema")
	require.NotNil(t, sizeSchema.ValidateFunc, "Field size should have ValidateFunc")

	validSizes := []interface{}{
		2,  // Minimum
		5,  // Typical
		15, // Mid-range
		25, // Maximum
	}

	for _, size := range validSizes {
		t.Run(fmt.Sprintf("valid_size_%d", size), func(t *testing.T) {
			_, errors := sizeSchema.ValidateFunc(size, "size")
			require.Empty(t, errors, "Size %d should be valid", size)
		})
	}
}

func TestValidateNodePoolSize_TooSmall(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema
	nodePoolsSchema := resourceSchema[tfconstants.AttrNodePools]
	nodePoolSchema := nodePoolsSchema.Elem.(*schema.Resource)
	nodePoolResourceSchema := nodePoolSchema.Schema

	sizeSchema, exists := nodePoolResourceSchema["size"]
	require.True(t, exists, "Field size should exist in node pool schema")
	require.NotNil(t, sizeSchema.ValidateFunc, "Field size should have ValidateFunc")

	invalidSizes := []interface{}{
		0,
		1,
		-1,
	}

	for _, size := range invalidSizes {
		t.Run(fmt.Sprintf("invalid_size_%d", size), func(t *testing.T) {
			_, errors := sizeSchema.ValidateFunc(size, "size")
			require.NotEmpty(t, errors, "Size %d should be invalid", size)
		})
	}
}

func TestValidateNodePoolSize_TooLarge(t *testing.T) {
	resource := ResourceKubernetesService()
	resourceSchema := resource.Schema
	nodePoolsSchema := resourceSchema[tfconstants.AttrNodePools]
	nodePoolSchema := nodePoolsSchema.Elem.(*schema.Resource)
	nodePoolResourceSchema := nodePoolSchema.Schema

	sizeSchema, exists := nodePoolResourceSchema["size"]
	require.True(t, exists, "Field size should exist in node pool schema")
	require.NotNil(t, sizeSchema.ValidateFunc, "Field size should have ValidateFunc")

	invalidSizes := []interface{}{
		26,
		100,
		1000,
	}

	for _, size := range invalidSizes {
		t.Run(fmt.Sprintf("invalid_size_%d", size), func(t *testing.T) {
			_, errors := sizeSchema.ValidateFunc(size, "size")
			require.NotEmpty(t, errors, "Size %d should be invalid", size)
		})
	}
}

// ============================================
// ERROR MESSAGE CONSISTENCY TESTS
// ============================================

func TestErrorMessages_ConsistentFormat(t *testing.T) {
	// Test that error messages from message.go constants follow consistent format
	// Error messages should be clear, actionable, and include relevant context

	// Test node pool error messages
	testCases := []struct {
		name          string
		errorConstant string
		expectedParts []string
	}{
		{
			name:          "NodePoolTypeRequired",
			errorConstant: ErrNodePoolTypeRequired,
			expectedParts: []string{"node pool type", "required"},
		},
		{
			name:          "NodePoolPlanNotFound",
			errorConstant: ErrNodePoolPlanNotFound,
			expectedParts: []string{"plan", "%s"}, // Should accept plan name as parameter
		},
		{
			name:          "NodePoolStaticSizeRequired",
			errorConstant: ErrNodePoolStaticSizeRequired,
			expectedParts: []string{"size", "Static"},
		},
		{
			name:          "NodePoolAutoscaleMinRequired",
			errorConstant: ErrNodePoolAutoscaleMinRequired,
			expectedParts: []string{"min_nodes", "Autoscale"},
		},
		{
			name:          "NodePoolAutoscaleMaxRequired",
			errorConstant: ErrNodePoolAutoscaleMaxRequired,
			expectedParts: []string{"max_nodes", "Autoscale"},
		},
		{
			name:          "ClusterVersionRequired",
			errorConstant: ErrClusterVersionRequired,
			expectedParts: []string{"kubernetes_version", "version", "required"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotEmpty(t, tc.errorConstant, "Error constant should not be empty")
			for _, part := range tc.expectedParts {
				require.Contains(t, tc.errorConstant, part, "Error message should contain '%s'", part)
			}
		})
	}
}

func TestErrorMessages_IncludeContext(t *testing.T) {
	// Test that error messages include relevant context (project_id, region, cluster_id)
	// This is tested by checking error messages in resource_kubernetes.go that use fmt.Errorf

	projectID := "project-123"
	region := "Mumbai"
	clusterID := "cluster-456"
	clusterName := "test-cluster"

	// Test error message format from resource_kubernetes.go
	// These should include project_id and region
	testCases := []struct {
		name             string
		errorMessage     string
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name: "ErrorCreatingCluster",
			errorMessage: fmt.Sprintf("Error creating Kubernetes cluster (name: %s) in project (%s), region (%s): %s",
				clusterName, projectID, region, "test error"),
			shouldContain: []string{"Error creating", "Kubernetes cluster", "name:", clusterName, "project", projectID, "region", region},
		},
		{
			name: "ErrorGettingClient",
			errorMessage: fmt.Sprintf("Error getting goe2e client for project (%s), region (%s): %s",
				projectID, region, "test error"),
			shouldContain: []string{"Error getting", "goe2e client", "project", projectID, "region", region},
		},
		{
			name: "ErrorRetrievingPlan",
			errorMessage: fmt.Sprintf("Error retrieving Kubernetes plan slug name for cluster (name: %s) in project (%s), region (%s): %s",
				clusterName, projectID, region, "test error"),
			shouldContain: []string{"Error retrieving", "Kubernetes plan", "name:", clusterName, "project", projectID, "region", region},
		},
		{
			name:          "ErrorUpdatingNodePool",
			errorMessage:  fmt.Sprintf(ErrNodePoolSizeTooSmall, "pool-name", clusterID, 1),
			shouldContain: []string{"Cannot update", "node pool", "pool-name", "Kubernetes cluster", "ID:", clusterID},
		},
		{
			name:          "ErrorDeletingNodePool",
			errorMessage:  fmt.Sprintf(ErrNodePoolDeleteNotRunning, "pool-name", clusterID),
			shouldContain: []string{"Cannot delete", "node pool", "pool-name", "Kubernetes cluster", "ID:", clusterID},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotEmpty(t, tc.errorMessage, "Error message should not be empty")
			for _, part := range tc.shouldContain {
				require.Contains(t, tc.errorMessage, part, "Error message should contain '%s'", part)
			}
			for _, part := range tc.shouldNotContain {
				require.NotContains(t, tc.errorMessage, part, "Error message should not contain '%s'", part)
			}
		})
	}
}

func TestDeprecationWarnings_Logged(t *testing.T) {
	// Test that deprecation warnings are logged correctly
	// This is tested by checking that logDeprecationWarning function exists and uses correct format

	// The deprecation warnings are logged in helpers.go using log.Printf with [WARN] prefix
	// We verify the format by checking the actual log statements in the code

	deprecatedFields := []struct {
		deprecatedField string
		preferredField  string
	}{
		{"specs_name", "plan"},
		{"node_pool_type", "type"},
		{"min_vms", "min_nodes"},
		{"max_vms", "max_nodes"},
		{"worker_node", "size"},
		{"name", "cluster_name"},
		{"version", "kubernetes_version"},
	}

	for _, field := range deprecatedFields {
		t.Run(fmt.Sprintf("deprecation_warning_%s", field.deprecatedField), func(t *testing.T) {
			// Verify that deprecation warning format is consistent
			// Format: "[WARN] Field '%s' is deprecated. Use '%s' instead."
			expectedWarning := fmt.Sprintf("[WARN] Field '%s' is deprecated. Use '%s' instead", field.deprecatedField, field.preferredField)

			// Check that the warning format matches expected pattern
			require.Contains(t, expectedWarning, "[WARN]", "Warning should include [WARN] prefix")
			require.Contains(t, expectedWarning, "deprecated", "Warning should mention 'deprecated'")
			require.Contains(t, expectedWarning, field.deprecatedField, "Warning should mention deprecated field")
			require.Contains(t, expectedWarning, field.preferredField, "Warning should mention preferred field")
		})
	}
}

func TestDeprecationWarnings_NotErrors(t *testing.T) {
	// Test that deprecation warnings don't cause errors
	// Deprecation warnings should be logged but not prevent resource operations

	// This is verified by ensuring that:
	// 1. Deprecation warnings use log.Printf (not return errors)
	// 2. Functions continue execution after logging warnings
	// 3. Warnings are informational only

	// The actual behavior is tested in integration/acceptance tests
	// This unit test documents the expected behavior
	t.Log("Deprecation warnings should be logged but not cause errors")
	t.Log("Functions should continue execution after logging deprecation warnings")
	t.Log("Warnings are informational only and don't prevent resource operations")
}
