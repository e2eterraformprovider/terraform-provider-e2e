package kubernetes

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// Mock KubernetesService for Testing
// ============================================================================

// mockKubernetesService is a mock implementation of goe2e.KubernetesService
// Only implements methods needed for testing helper functions
type mockKubernetesService struct {
	mock.Mock
}

func (m *mockKubernetesService) GetWorkerPlans(ctx context.Context) ([]goe2e.KubernetesWorkerPlan, *goe2e.Response, error) {
	args := m.Called(ctx)
	var plans []goe2e.KubernetesWorkerPlan
	if args.Get(0) != nil {
		plans = args.Get(0).([]goe2e.KubernetesWorkerPlan)
	}
	var resp *goe2e.Response
	if args.Get(1) != nil {
		resp = args.Get(1).(*goe2e.Response)
	}
	return plans, resp, args.Error(2)
}

// Stub implementations for interface compliance (not used in these tests)
func (m *mockKubernetesService) GetMasterPlans(ctx context.Context) ([]goe2e.KubernetesPlan, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) Create(ctx context.Context, req *goe2e.KubernetesClusterCreateRequest) (*goe2e.KubernetesCluster, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) Get(ctx context.Context, clusterID string) (*goe2e.KubernetesCluster, *goe2e.Response, error) {
	args := m.Called(ctx, clusterID)
	var cluster *goe2e.KubernetesCluster
	if args.Get(0) != nil {
		cluster = args.Get(0).(*goe2e.KubernetesCluster)
	}
	var resp *goe2e.Response
	if args.Get(1) != nil {
		resp = args.Get(1).(*goe2e.Response)
	}
	return cluster, resp, args.Error(2)
}

func (m *mockKubernetesService) Delete(ctx context.Context, clusterID string) (*goe2e.Response, error) {
	return nil, nil
}

func (m *mockKubernetesService) GetNodePools(ctx context.Context, clusterID string) ([]goe2e.NodePoolServiceInfo, *goe2e.Response, error) {
	args := m.Called(ctx, clusterID)
	var pools []goe2e.NodePoolServiceInfo
	if args.Get(0) != nil {
		pools = args.Get(0).([]goe2e.NodePoolServiceInfo)
	}
	var resp *goe2e.Response
	if args.Get(1) != nil {
		resp = args.Get(1).(*goe2e.Response)
	}
	return pools, resp, args.Error(2)
}

func (m *mockKubernetesService) AddNodePool(ctx context.Context, clusterID string, req *goe2e.NodePoolAddRequest) (*goe2e.NodePoolAddResponse, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) UpdateNodePoolCardinality(ctx context.Context, nodePoolServiceID string, req *goe2e.NodePoolResizeRequest) (*goe2e.Response, error) {
	return nil, nil
}

func (m *mockKubernetesService) UpdateNodePoolDetails(ctx context.Context, nodePoolServiceID string, req *goe2e.NodePoolUpdateRequest) (*goe2e.NodePoolUpdateResponse, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) DeleteNodePool(ctx context.Context, nodePoolServiceID string) (*goe2e.Response, error) {
	return nil, nil
}

func (m *mockKubernetesService) CheckNodePoolStatus(ctx context.Context, clusterID string) ([]goe2e.NodePoolServiceInfo, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) ListPersistentVolumes(ctx context.Context, clusterID string) ([]goe2e.PersistentVolume, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) CreatePersistentVolume(ctx context.Context, clusterID string, req *goe2e.CreatePersistentVolumeRequest) (*goe2e.PersistentVolume, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) GetPersistentVolume(ctx context.Context, clusterID, pvID string) (*goe2e.PersistentVolume, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) DeletePersistentVolume(ctx context.Context, clusterID, pvID string) (*goe2e.Response, error) {
	return nil, nil
}

func (m *mockKubernetesService) ListAttachedSecurityGroups(ctx context.Context, clusterID string) ([]goe2e.SecurityGroupAttachment, *goe2e.Response, error) {
	return nil, nil, nil
}

func (m *mockKubernetesService) AttachSecurityGroups(ctx context.Context, clusterID string, req *goe2e.AttachSecurityGroupRequest) (*goe2e.Response, error) {
	return nil, nil
}

func (m *mockKubernetesService) DetachSecurityGroups(ctx context.Context, clusterID string, req *goe2e.DetachSecurityGroupRequest) (*goe2e.Response, error) {
	return nil, nil
}

// Helper function to create a mock goe2e.Client with KubernetesService
func createMockKubernetesClient(mockService *mockKubernetesService) *goe2e.Client {
	return &goe2e.Client{
		Kubernetes: mockService,
	}
}

// Helper function to create test worker plans
func createTestWorkerPlans() []goe2e.KubernetesWorkerPlan {
	return []goe2e.KubernetesWorkerPlan{
		{
			Plan: "C3.8GB",
			Specs: goe2e.KubernetesWorkerPlanSpecs{
				ID:      "sku-123",
				SKUName: "C3.8GB",
			},
		},
		{
			Plan: "C3.16GB",
			Specs: goe2e.KubernetesWorkerPlanSpecs{
				ID:      "sku-456",
				SKUName: "C3.16GB",
			},
		},
	}
}

// ============================================================================
// ExpandNodePools Function Tests
// ============================================================================

func TestExpandNodePools_EmptyConfig(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{}
	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_SingleStaticPool(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name": "static-pool",
			"plan": "C3.8GB",
			"type": goe2econstants.KubernetesNodePoolTypeStatic,
			"size": 3,
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "static-pool", result[0].Name)
	assert.Equal(t, "C3.8GB", result[0].SpecsName)
	assert.Equal(t, 3, result[0].WorkerNode)
	assert.Equal(t, "C3.8GB", result[0].SlugName)
	assert.Equal(t, "sku-123", result[0].SKUID)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_SingleStaticPool_DeprecatedFields(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name":           "static-pool-v2",
			"specs_name":     "C3.8GB",
			"node_pool_type": goe2econstants.KubernetesNodePoolTypeStatic,
			"worker_node":    5,
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "static-pool-v2", result[0].Name)
	assert.Equal(t, "C3.8GB", result[0].SpecsName)
	assert.Equal(t, 5, result[0].WorkerNode)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_SingleAutoscalePool(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name":            "autoscale-pool",
			"plan":            "C3.16GB",
			"type":            goe2econstants.KubernetesNodePoolTypeAutoscale,
			"node_pool_type":  goe2econstants.KubernetesNodePoolTypeAutoscale, // Also set for getElasticityDict
			"min_nodes":       2,
			"max_nodes":       10,
			"elasticity_dict": []interface{}{}, // Empty but present
			"scheduled_dict":  []interface{}{}, // Empty but present
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "autoscale-pool", result[0].Name)
	assert.Equal(t, "C3.16GB", result[0].SpecsName)
	assert.Equal(t, 0, result[0].WorkerNode) // Autoscale pools don't have WorkerNode
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_SingleAutoscalePool_DeprecatedFields(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name":            "autoscale-pool-v2",
			"specs_name":      "C3.16GB",
			"node_pool_type":  goe2econstants.KubernetesNodePoolTypeAutoscale,
			"min_vms":         3,
			"max_vms":         15,
			"elasticity_dict": []interface{}{}, // Empty but present
			"scheduled_dict":  []interface{}{}, // Empty but present
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "autoscale-pool-v2", result[0].Name)
	assert.Equal(t, "C3.16GB", result[0].SpecsName)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_MultiplePools(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name": "pool-1",
			"plan": "C3.8GB",
			"type": goe2econstants.KubernetesNodePoolTypeStatic,
			"size": 3,
		},
		map[string]interface{}{
			"name":            "pool-2",
			"plan":            "C3.16GB",
			"type":            goe2econstants.KubernetesNodePoolTypeAutoscale,
			"node_pool_type":  goe2econstants.KubernetesNodePoolTypeAutoscale, // Also set for getElasticityDict
			"min_nodes":       2,
			"max_nodes":       10,
			"elasticity_dict": []interface{}{}, // Empty but present
			"scheduled_dict":  []interface{}{}, // Empty but present
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "pool-1", result[0].Name)
	assert.Equal(t, "pool-2", result[1].Name)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_DuplicateNames(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name": "duplicate",
			"plan": "C3.8GB",
			"type": goe2econstants.KubernetesNodePoolTypeStatic,
			"size": 3,
		},
		map[string]interface{}{
			"name": "duplicate", // Same name
			"plan": "C3.8GB",
			"type": goe2econstants.KubernetesNodePoolTypeStatic,
			"size": 3,
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrNodePoolDuplicateNames)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_MissingPlan(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name": "pool",
			"plan": "NONEXISTENT",
			"type": goe2econstants.KubernetesNodePoolTypeStatic,
			"size": 3,
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no matching plan found for plan")
	assert.Contains(t, err.Error(), "NONEXISTENT")
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_StaticMissingSize(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name": "pool",
			"plan": "C3.8GB",
			"type": goe2econstants.KubernetesNodePoolTypeStatic,
			// size is missing
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrNodePoolStaticSizeRequired)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_AutoscaleMissingMinNodes(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name":      "pool",
			"plan":      "C3.16GB",
			"type":      goe2econstants.KubernetesNodePoolTypeAutoscale,
			"max_nodes": 10,
			// min_nodes is missing
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrNodePoolAutoscaleMinRequired)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_AutoscaleMissingMaxNodes(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name":      "pool",
			"plan":      "C3.16GB",
			"type":      goe2econstants.KubernetesNodePoolTypeAutoscale,
			"min_nodes": 2,
			// max_nodes is missing
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrNodePoolAutoscaleMaxRequired)
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_InvalidPoolType(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name": "pool",
			"plan": "C3.8GB",
			"type": "InvalidType",
			"size": 3,
		},
	}

	_, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	// Invalid type should be caught by validation, but if it passes through,
	// it will be treated as missing type
	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_WithElasticityDict(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name":           "autoscale-pool",
			"plan":           "C3.16GB",
			"type":           goe2econstants.KubernetesNodePoolTypeAutoscale,
			"node_pool_type": goe2econstants.KubernetesNodePoolTypeAutoscale, // Also set for getElasticityDict
			"min_nodes":      2,
			"max_nodes":      10,
			"elasticity_dict": []interface{}{
				map[string]interface{}{
					"worker": []interface{}{
						map[string]interface{}{
							"policy_paramter_type": goe2econstants.KubernetesPolicyTypeDefault,
							"parameter":            goe2econstants.KubernetesPolicyParameterCPU,
							"elasticity_policies": []interface{}{
								map[string]interface{}{
									"operator":     ">",
									"value":        80,
									"period":       5,
									"watch_period": 3,
									"cooldown":     10,
								},
							},
						},
					},
				},
			},
			"scheduled_dict": []interface{}{}, // Empty but present
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	// ElasticityDict should be populated
	assert.NotNil(t, result[0].ElasticityDict)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_WithScheduledDict(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name":            "autoscale-pool",
			"plan":            "C3.16GB",
			"type":            goe2econstants.KubernetesNodePoolTypeAutoscale,
			"node_pool_type":  goe2econstants.KubernetesNodePoolTypeAutoscale, // Also set for getScheduledDict
			"min_nodes":       2,
			"max_nodes":       10,
			"elasticity_dict": []interface{}{}, // Empty but present
			"scheduled_dict": []interface{}{
				map[string]interface{}{
					"worker": []interface{}{
						map[string]interface{}{
							"scheduled_policies": []interface{}{
								map[string]interface{}{
									"upscale_cardinality":   5,
									"upscale_recurrence":    "0 12 * * *",
									"downscale_cardinality": 2,
									"downscale_recurrence":  "0 2 * * *",
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	// ScheduledDict should be populated
	assert.NotNil(t, result[0].ScheduledDict)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_ClientError(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	mockService.On("GetWorkerPlans", ctx).Return(nil, &goe2e.Response{}, errors.New("API error"))

	config := []interface{}{
		map[string]interface{}{
			"name": "pool",
			"plan": "C3.8GB",
			"type": goe2econstants.KubernetesNodePoolTypeStatic,
			"size": 3,
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching worker plans")
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

func TestExpandNodePools_PlanNotFound(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	config := []interface{}{
		map[string]interface{}{
			"name": "pool",
			"plan": "INVALID-PLAN",
			"type": goe2econstants.KubernetesNodePoolTypeStatic,
			"size": 3,
		},
	}

	result, err := ExpandNodePools(ctx, config, client, "project-123", "Mumbai")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no matching plan found for plan")
	assert.Contains(t, err.Error(), "INVALID-PLAN")
	assert.Empty(t, result)
	mockService.AssertExpectations(t)
}

// ============================================================================
// ExpandElasticityDict Function Tests
// ============================================================================

func TestExpandElasticityDict_EmptyConfig(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{},
	}

	result, err := ExpandElasticityDict(config, 2, 10)

	assert.NoError(t, err)
	assert.Empty(t, result.Worker.ElasticityPolicies)
}

func TestExpandElasticityDict_ValidConfig(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{
			map[string]interface{}{
				"policy_paramter_type": goe2econstants.KubernetesPolicyTypeDefault,
				"parameter":            goe2econstants.KubernetesPolicyParameterCPU,
				"elasticity_policies": []interface{}{
					map[string]interface{}{
						"operator":     ">",
						"value":        80,
						"period":       5,
						"watch_period": 3,
						"cooldown":     10,
					},
					map[string]interface{}{
						"operator":     "<",
						"value":        20,
						"period":       5,
						"watch_period": 3,
						"cooldown":     10,
					},
				},
			},
		},
	}

	result, err := ExpandElasticityDict(config, 2, 10)

	assert.NoError(t, err)
	assert.Len(t, result.Worker.ElasticityPolicies, 2)
	assert.Equal(t, 2, result.Worker.MinVms)
	assert.Equal(t, 10, result.Worker.MaxVms)
	assert.Equal(t, goe2econstants.KubernetesPolicyTypeChange, result.Worker.ElasticityPolicies[0].Type)
}

func TestExpandElasticityDict_InvalidMinVMs(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{
			map[string]interface{}{
				"policy_paramter_type": goe2econstants.KubernetesPolicyTypeDefault,
				"parameter":            goe2econstants.KubernetesPolicyParameterCPU,
				"elasticity_policies": []interface{}{
					map[string]interface{}{
						"operator":     ">",
						"value":        80,
						"period":       5,
						"watch_period": 3,
						"cooldown":     10,
					},
				},
			},
		},
	}

	// MinVMs > MaxVMs should be handled by validation, but test the function behavior
	// The function itself doesn't validate min/max, it just uses them
	_, err := ExpandElasticityDict(config, 15, 10)

	// The function doesn't validate min/max relationship, it just uses the values
	assert.NoError(t, err)
}

func TestExpandElasticityDict_MultiplePolicies(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{
			map[string]interface{}{
				"policy_paramter_type": goe2econstants.KubernetesPolicyTypeDefault,
				"parameter":            goe2econstants.KubernetesPolicyParameterMemory,
				"elasticity_policies": []interface{}{
					map[string]interface{}{
						"operator":     ">",
						"value":        90,
						"period":       5,
						"watch_period": 3,
						"cooldown":     10,
					},
					map[string]interface{}{
						"operator":     "<",
						"value":        10,
						"period":       5,
						"watch_period": 3,
						"cooldown":     10,
					},
					map[string]interface{}{
						"operator":     ">=",
						"value":        85,
						"period":       5,
						"watch_period": 3,
						"cooldown":     10,
					},
				},
			},
		},
	}

	result, err := ExpandElasticityDict(config, 2, 10)

	assert.NoError(t, err)
	assert.Len(t, result.Worker.ElasticityPolicies, 3)
	// Check adjust values alternate: 1, -1, 1
	assert.Equal(t, 1, result.Worker.ElasticityPolicies[0].Adjust)
	assert.Equal(t, -1, result.Worker.ElasticityPolicies[1].Adjust)
	assert.Equal(t, 1, result.Worker.ElasticityPolicies[2].Adjust)
}

func TestExpandElasticityDict_CustomParameter(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{
			map[string]interface{}{
				"policy_paramter_type": goe2econstants.KubernetesPolicyTypeCustom,
				"parameter":            "CUSTOM_METRIC",
				"elasticity_policies": []interface{}{
					map[string]interface{}{
						"operator":     ">",
						"value":        75,
						"period":       5,
						"watch_period": 3,
						"cooldown":     10,
					},
				},
			},
		},
	}

	result, err := ExpandElasticityDict(config, 2, 10)

	assert.NoError(t, err)
	assert.Len(t, result.Worker.ElasticityPolicies, 1)
	assert.Equal(t, "CUSTOM_METRIC", result.Worker.ElasticityPolicies[0].Parameter)
}

// ============================================================================
// ExpandScheduledDict Function Tests
// ============================================================================

func TestExpandScheduledDict_EmptyConfig(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{},
	}

	result, err := ExpandScheduledDict(config, 2, 10)

	assert.NoError(t, err)
	assert.Empty(t, result.Worker.ScheduledPolicies)
}

func TestExpandScheduledDict_ValidConfig(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{
			map[string]interface{}{
				"scheduled_policies": []interface{}{
					map[string]interface{}{
						"upscale_cardinality":   5,
						"upscale_recurrence":    "0 12 * * *",
						"downscale_cardinality": 2,
						"downscale_recurrence":  "0 2 * * *",
					},
				},
			},
		},
	}

	result, err := ExpandScheduledDict(config, 2, 10)

	assert.NoError(t, err)
	assert.Len(t, result.Worker.ScheduledPolicies, 2) // upscale + downscale
	assert.Equal(t, 2, result.Worker.MinVms)
	assert.Equal(t, 10, result.Worker.MaxVms)
	assert.Equal(t, goe2econstants.KubernetesPolicyTypeCardinality, result.Worker.ScheduledPolicies[0].Type)
	assert.Equal(t, 5, result.Worker.ScheduledPolicies[0].Adjust) // upscale
	assert.Equal(t, 2, result.Worker.ScheduledPolicies[1].Adjust) // downscale
}

func TestExpandScheduledDict_InvalidMinVMs(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{
			map[string]interface{}{
				"scheduled_policies": []interface{}{
					map[string]interface{}{
						"upscale_cardinality":   15, // > max_vms
						"upscale_recurrence":    "0 12 * * *",
						"downscale_cardinality": 2,
						"downscale_recurrence":  "0 2 * * *",
					},
				},
			},
		},
	}

	result, err := ExpandScheduledDict(config, 2, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrUpscaleCardinalityRange)
	assert.Empty(t, result.Worker.ScheduledPolicies)
}

func TestExpandScheduledDict_InvalidMaxVMs(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{
			map[string]interface{}{
				"scheduled_policies": []interface{}{
					map[string]interface{}{
						"upscale_cardinality":   5,
						"upscale_recurrence":    "0 12 * * *",
						"downscale_cardinality": 1, // < min_vms
						"downscale_recurrence":  "0 2 * * *",
					},
				},
			},
		},
	}

	result, err := ExpandScheduledDict(config, 2, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrDownscaleCardinalityRange)
	assert.Empty(t, result.Worker.ScheduledPolicies)
}

func TestExpandScheduledDict_MultiplePolicies(t *testing.T) {
	config := map[string]interface{}{
		"worker": []interface{}{
			map[string]interface{}{
				"scheduled_policies": []interface{}{
					map[string]interface{}{
						"upscale_cardinality":   5,
						"upscale_recurrence":    "0 12 * * *",
						"downscale_cardinality": 2,
						"downscale_recurrence":  "0 2 * * *",
					},
					map[string]interface{}{
						"upscale_cardinality":   8,
						"upscale_recurrence":    "0 20 * * *",
						"downscale_cardinality": 3,
						"downscale_recurrence":  "0 0 * * *",
					},
				},
			},
		},
	}

	result, err := ExpandScheduledDict(config, 2, 10)

	assert.NoError(t, err)
	assert.Len(t, result.Worker.ScheduledPolicies, 4) // 2 upscale + 2 downscale
}

// ============================================================================
// ExpandNodePoolUpdate Function Tests
// ============================================================================

func TestExpandNodePoolUpdate_StaticPoolResize(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	nodePoolDetail := map[string]interface{}{
		"name":        "static-pool",
		"plan":        "C3.8GB",
		"type":        goe2econstants.KubernetesNodePoolTypeStatic,
		"cardinality": 5, // Resize to 5
	}

	result, err := ExpandNodePoolUpdate(ctx, nodePoolDetail, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Equal(t, 0, result.MinVms) // Static pools have 0 min/max
	assert.Equal(t, 0, result.MaxVms)
	assert.Equal(t, 5, result.Cardinality)
	assert.Equal(t, "sku-123", result.PlanID)
	mockService.AssertExpectations(t)
}

func TestExpandNodePoolUpdate_AutoscalePoolMinMax(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	nodePoolDetail := map[string]interface{}{
		"name":            "autoscale-pool",
		"plan":            "C3.16GB",
		"type":            goe2econstants.KubernetesNodePoolTypeAutoscale,
		"node_pool_type":  goe2econstants.KubernetesNodePoolTypeAutoscale, // Also set for updateElasticPolicies
		"min_nodes":       3,
		"max_nodes":       12,
		"cardinality":     5,
		"elasticity_dict": []interface{}{}, // Empty but present
		"scheduled_dict":  []interface{}{}, // Empty but present
	}

	result, err := ExpandNodePoolUpdate(ctx, nodePoolDetail, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Equal(t, 3, result.MinVms)
	assert.Equal(t, 12, result.MaxVms)
	assert.Equal(t, 5, result.Cardinality)
	assert.Equal(t, "sku-456", result.PlanID)
	mockService.AssertExpectations(t)
}

func TestExpandNodePoolUpdate_PlanChange(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	nodePoolDetail := map[string]interface{}{
		"name":            "autoscale-pool",
		"plan":            "C3.16GB", // Changed plan
		"type":            goe2econstants.KubernetesNodePoolTypeAutoscale,
		"node_pool_type":  goe2econstants.KubernetesNodePoolTypeAutoscale, // Also set for updateElasticPolicies
		"min_nodes":       2,
		"max_nodes":       10,
		"elasticity_dict": []interface{}{}, // Empty but present
		"scheduled_dict":  []interface{}{}, // Empty but present
	}

	result, err := ExpandNodePoolUpdate(ctx, nodePoolDetail, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.Equal(t, "sku-456", result.PlanID) // Should use new plan
	mockService.AssertExpectations(t)
}

func TestExpandNodePoolUpdate_WithElasticityPolicies(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	nodePoolDetail := map[string]interface{}{
		"name":           "autoscale-pool",
		"plan":           "C3.16GB",
		"type":           goe2econstants.KubernetesNodePoolTypeAutoscale,
		"min_nodes":      2,
		"max_nodes":      10,
		"node_pool_type": goe2econstants.KubernetesNodePoolTypeAutoscale,
		"elasticity_dict": []interface{}{
			map[string]interface{}{
				"worker": []interface{}{
					map[string]interface{}{
						"policy_paramter_type": goe2econstants.KubernetesPolicyTypeDefault,
						"parameter":            goe2econstants.KubernetesPolicyParameterCPU,
						"elasticity_policies": []interface{}{
							map[string]interface{}{
								"operator":     ">",
								"value":        80,
								"period":       5,
								"watch_period": 3,
								"cooldown":     10,
							},
						},
					},
				},
			},
		},
		"scheduled_dict": []interface{}{}, // Empty but present
	}

	result, err := ExpandNodePoolUpdate(ctx, nodePoolDetail, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.NotEmpty(t, result.ElasticityPolicy)
	mockService.AssertExpectations(t)
}

func TestExpandNodePoolUpdate_WithScheduledPolicies(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	workerPlans := createTestWorkerPlans()
	mockService.On("GetWorkerPlans", ctx).Return(workerPlans, &goe2e.Response{}, nil)

	nodePoolDetail := map[string]interface{}{
		"name":            "autoscale-pool",
		"plan":            "C3.16GB",
		"type":            goe2econstants.KubernetesNodePoolTypeAutoscale,
		"min_nodes":       2,
		"max_nodes":       10,
		"node_pool_type":  goe2econstants.KubernetesNodePoolTypeAutoscale,
		"elasticity_dict": []interface{}{}, // Empty but present
		"scheduled_dict": []interface{}{
			map[string]interface{}{
				"worker": []interface{}{
					map[string]interface{}{
						"scheduled_policies": []interface{}{
							map[string]interface{}{
								"upscale_cardinality":   5,
								"upscale_recurrence":    "0 12 * * *",
								"downscale_cardinality": 2,
								"downscale_recurrence":  "0 2 * * *",
							},
						},
					},
				},
			},
		},
	}

	result, err := ExpandNodePoolUpdate(ctx, nodePoolDetail, client, "project-123", "Mumbai")

	assert.NoError(t, err)
	assert.NotEmpty(t, result.ScheduledPolicy)
	mockService.AssertExpectations(t)
}

func TestExpandNodePoolUpdate_ClientError(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	mockService.On("GetWorkerPlans", ctx).Return(nil, &goe2e.Response{}, errors.New("API error"))

	nodePoolDetail := map[string]interface{}{
		"name": "pool",
		"plan": "C3.8GB",
		"type": goe2econstants.KubernetesNodePoolTypeStatic,
	}

	result, err := ExpandNodePoolUpdate(ctx, nodePoolDetail, client, "project-123", "Mumbai")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error fetching worker plans")
	assert.Empty(t, result.PlanID)
	mockService.AssertExpectations(t)
}

// Helper function to create test Kubernetes cluster
func createTestKubernetesCluster(id, state string) *goe2e.KubernetesCluster {
	return &goe2e.KubernetesCluster{
		ID:        id,
		ServiceID: id,
		State:     state,
		Version:   "1.30",
		VPCID:     "vpc-123",
	}
}

// ============================================================================
// Wait Function Tests
// ============================================================================

func TestWaitForClusterStatus_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	clusterID := "cluster-123"
	cluster := createTestKubernetesCluster(clusterID, goe2econstants.KubernetesClusterStatusRunning)

	// Mock returns Running status immediately
	mockService.On("Get", ctx, clusterID).Return(cluster, &goe2e.Response{}, nil)

	err := waitForClusterStatus(ctx, client, clusterID, goe2econstants.KubernetesClusterStatusRunning, 1*time.Minute)

	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestWaitForClusterStatus_Timeout(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	clusterID := "cluster-123"
	// Cluster stays in Creating state (never reaches Running)
	cluster := createTestKubernetesCluster(clusterID, goe2econstants.KubernetesClusterStatusCreating)

	// Mock returns Creating status (StateChangeConf will call multiple times until timeout)
	// Use Maybe() to allow multiple calls without exact count
	mockService.On("Get", ctx, clusterID).Return(cluster, &goe2e.Response{}, nil).Maybe()

	// Use short timeout to test timeout behavior (but long enough for at least one call)
	timeout := 200 * time.Millisecond
	err := waitForClusterStatus(ctx, client, clusterID, goe2econstants.KubernetesClusterStatusRunning, timeout)

	assert.Error(t, err)
	// Timeout error message may vary, just verify it's an error
	assert.NotNil(t, err)
}

func TestWaitForClusterStatus_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	clusterID := "cluster-123"
	cluster := createTestKubernetesCluster(clusterID, goe2econstants.KubernetesClusterStatusCreating)

	// Cancel context before starting wait
	cancel()

	// Mock may or may not be called depending on timing
	mockService.On("Get", ctx, clusterID).Return(cluster, &goe2e.Response{}, nil).Maybe()

	err := waitForClusterStatus(ctx, client, clusterID, goe2econstants.KubernetesClusterStatusRunning, 1*time.Minute)

	assert.Error(t, err)
	// Context cancellation error message varies, just check it's an error
	assert.NotNil(t, err)
}

func TestClusterStatusRefresh_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	clusterID := "cluster-123"
	cluster := createTestKubernetesCluster(clusterID, goe2econstants.KubernetesClusterStatusRunning)

	mockService.On("Get", ctx, clusterID).Return(cluster, &goe2e.Response{}, nil)

	refreshFunc := clusterStatusRefresh(ctx, client, clusterID)
	result, state, err := refreshFunc()

	assert.NoError(t, err)
	assert.Equal(t, cluster, result)
	assert.Equal(t, goe2econstants.KubernetesClusterStatusRunning, state)
	mockService.AssertExpectations(t)
}

func TestClusterStatusRefresh_NilCluster(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	clusterID := "cluster-123"

	// Mock returns nil cluster
	mockService.On("Get", ctx, clusterID).Return(nil, &goe2e.Response{}, nil)

	refreshFunc := clusterStatusRefresh(ctx, client, clusterID)
	result, state, err := refreshFunc()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cluster returned nil")
	assert.Nil(t, result)
	assert.Empty(t, state)
	mockService.AssertExpectations(t)
}

func TestClusterStatusRefresh_ClientError(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	clusterID := "cluster-123"

	// Mock returns error
	mockService.On("Get", ctx, clusterID).Return(nil, nil, errors.New("API error"))

	refreshFunc := clusterStatusRefresh(ctx, client, clusterID)
	result, state, err := refreshFunc()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
	assert.Nil(t, result)
	assert.Empty(t, state)
	mockService.AssertExpectations(t)
}

func TestClusterStatusRefresh_404Error(t *testing.T) {
	ctx := context.Background()
	mockService := &mockKubernetesService{}
	client := createMockKubernetesClient(mockService)

	clusterID := "cluster-123"

	// Mock returns 404 error (cluster deleted)
	resp := &goe2e.Response{
		Response: &http.Response{
			StatusCode: http.StatusNotFound,
		},
	}
	mockService.On("Get", ctx, clusterID).Return(nil, resp, errors.New("404 Not Found"))

	refreshFunc := clusterStatusRefresh(ctx, client, clusterID)
	result, state, err := refreshFunc()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
	assert.Nil(t, result)
	assert.Empty(t, state)
	mockService.AssertExpectations(t)
}

// ============================================================================
// Import Function Tests
// ============================================================================

func TestResourceKubernetesImport_InvalidFormat(t *testing.T) {
	ctx := context.Background()
	resource := ResourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	// Test invalid format (too many parts)
	d.SetId("project:region:cluster:extra")

	cfg := &config.Config{}
	result, err := resourceKubernetesImport(ctx, d, cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid import ID format")
	assert.Nil(t, result)
}

// ============================================================================
// Helper Function Edge Case Tests
// ============================================================================

func TestGetClusterName_NilResourceData(t *testing.T) {
	// Test with empty ResourceData (no fields set)
	resource := ResourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	name := getClusterName(d)
	assert.Empty(t, name)
}

func TestGetKubernetesVersion_InvalidType(t *testing.T) {
	// Test that function handles missing fields gracefully
	resource := ResourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	version := getKubernetesVersion(d)
	assert.Empty(t, version)
}

func TestGetNodePoolPlan_InvalidType(t *testing.T) {
	// Test with invalid type in map
	pool := map[string]interface{}{
		"plan": 123, // Invalid type (should be string)
	}

	// Function should handle type assertion gracefully
	plan := getNodePoolPlan(pool)
	// Should return empty string when type assertion fails
	assert.Empty(t, plan)
}

func TestGetNodePoolType_InvalidType(t *testing.T) {
	// Test with invalid type in map
	pool := map[string]interface{}{
		"type": 123, // Invalid type (should be string)
	}

	poolType := getNodePoolType(pool)
	// Should return empty string when type assertion fails
	assert.Empty(t, poolType)
}

func TestExpandSecurityGroupIDs_InvalidType(t *testing.T) {
	// Test with non-int values
	input := []interface{}{"not-an-int", 42, "also-not-int"}

	// This should panic with type assertion error
	// In production, schema validation should prevent this
	assert.Panics(t, func() {
		expandSecurityGroupIDs(input)
	}, "Should panic on invalid type")
}

func TestExpandSecurityGroupIDs_NegativeValues(t *testing.T) {
	// Test with negative IDs (should be validated by schema, but test function behavior)
	input := []interface{}{-1, -5, 10}

	result := expandSecurityGroupIDs(input)

	// Function doesn't validate, just converts
	assert.Len(t, result, 3)
	assert.Equal(t, -1, result[0])
	assert.Equal(t, -5, result[1])
	assert.Equal(t, 10, result[2])
	// Note: Schema validation should prevent negative values in production
}
