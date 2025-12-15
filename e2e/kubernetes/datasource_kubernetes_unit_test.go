package kubernetes

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Mock implementations for Kubernetes datasource tests
// ============================================================================

// mockKubernetesServiceForDatasource is a mock implementation of KubernetesService for datasource testing
type mockKubernetesServiceForDatasource struct {
	getFunc func(ctx context.Context, clusterID string) (*goe2e.KubernetesCluster, *goe2e.Response, error)
}

func (m *mockKubernetesServiceForDatasource) Get(ctx context.Context, clusterID string) (*goe2e.KubernetesCluster, *goe2e.Response, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, clusterID)
	}
	return nil, nil, errors.New("not implemented")
}

// Unused interface methods
func (m *mockKubernetesServiceForDatasource) GetMasterPlans(ctx context.Context) ([]goe2e.KubernetesPlan, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) GetWorkerPlans(ctx context.Context) ([]goe2e.KubernetesWorkerPlan, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) Create(ctx context.Context, req *goe2e.KubernetesClusterCreateRequest) (*goe2e.KubernetesCluster, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) Delete(ctx context.Context, clusterID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) AddNodePool(ctx context.Context, clusterID string, req *goe2e.NodePoolAddRequest) (*goe2e.NodePoolAddResponse, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) UpdateNodePoolCardinality(ctx context.Context, clusterID string, req *goe2e.NodePoolResizeRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) UpdateNodePoolDetails(ctx context.Context, clusterID string, req *goe2e.NodePoolUpdateRequest) (*goe2e.NodePoolUpdateResponse, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) DeleteNodePool(ctx context.Context, clusterID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) GetNodePools(ctx context.Context, clusterID string) ([]goe2e.NodePoolServiceInfo, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) CheckNodePoolStatus(ctx context.Context, clusterID string) ([]goe2e.NodePoolServiceInfo, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) ListPersistentVolumes(ctx context.Context, clusterID string) ([]goe2e.PersistentVolume, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) CreatePersistentVolume(ctx context.Context, clusterID string, req *goe2e.CreatePersistentVolumeRequest) (*goe2e.PersistentVolume, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) GetPersistentVolume(ctx context.Context, clusterID string, volumeID string) (*goe2e.PersistentVolume, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) DeletePersistentVolume(ctx context.Context, clusterID string, volumeID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) ListAttachedSecurityGroups(ctx context.Context, clusterID string) ([]goe2e.SecurityGroupAttachment, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) AttachSecurityGroups(ctx context.Context, clusterID string, req *goe2e.AttachSecurityGroupRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockKubernetesServiceForDatasource) DetachSecurityGroups(ctx context.Context, clusterID string, req *goe2e.DetachSecurityGroupRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

// createMockConfigForKubernetesDatasource creates a config with a mock kubernetes service
func createMockConfigForKubernetesDatasource(t *testing.T, mockService *mockKubernetesServiceForDatasource, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.Kubernetes = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// ============================================================================
// Test: dataSourceReadKubernetes
// ============================================================================

func TestDataSourceReadKubernetes_Success(t *testing.T) {
	mockService := &mockKubernetesServiceForDatasource{
		getFunc: func(ctx context.Context, clusterID string) (*goe2e.KubernetesCluster, *goe2e.Response, error) {
			assert.Equal(t, "k8s-123", clusterID)
			return &goe2e.KubernetesCluster{
				ServiceID:   "k8s-123",
				ServiceName: "test-cluster",
				State:       "Running",
				Version:     "1.28.0",
				CreatedAt:   "2024-01-01T00:00:00Z",
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForKubernetesDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"service_id": "k8s-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "k8s-123", d.Id())
	assert.Equal(t, "test-cluster", d.Get(tfconstants.AttrName))
	assert.Equal(t, "Running", d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "1.28.0", d.Get(tfconstants.AttrVersion))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get(tfconstants.AttrCreatedAt))
}

func TestDataSourceReadKubernetes_NotFound(t *testing.T) {
	mockService := &mockKubernetesServiceForDatasource{
		getFunc: func(ctx context.Context, clusterID string) (*goe2e.KubernetesCluster, *goe2e.Response, error) {
			return nil, nil, errors.New("kubernetes cluster not found")
		},
	}

	cfg := createMockConfigForKubernetesDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"service_id": "non-existent-cluster",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when cluster not found")
	assert.Contains(t, diags[0].Detail, "error finding Kubernetes cluster")
}

func TestDataSourceReadKubernetes_NilCluster(t *testing.T) {
	mockService := &mockKubernetesServiceForDatasource{
		getFunc: func(ctx context.Context, clusterID string) (*goe2e.KubernetesCluster, *goe2e.Response, error) {
			return nil, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForKubernetesDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"service_id": "k8s-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when cluster is nil")
	assert.Contains(t, diags[0].Detail, "not found")
}

func TestDataSourceReadKubernetes_APIError(t *testing.T) {
	mockService := &mockKubernetesServiceForDatasource{
		getFunc: func(ctx context.Context, clusterID string) (*goe2e.KubernetesCluster, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to retrieve cluster")
		},
	}

	cfg := createMockConfigForKubernetesDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"service_id": "k8s-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	assert.Contains(t, diags[0].Detail, "error finding Kubernetes cluster")
}

func TestDataSourceReadKubernetes_MissingRegion(t *testing.T) {
	cfg := createMockConfigForKubernetesDatasource(t, &mockKubernetesServiceForDatasource{}, "", "")
	resource := DataSourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"service_id": "k8s-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail without region")
}

func TestDataSourceReadKubernetes_MissingProjectID(t *testing.T) {
	cfg := createMockConfigForKubernetesDatasource(t, &mockKubernetesServiceForDatasource{}, "", "us-east-1")
	resource := DataSourceKubernetesService()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"service_id": "k8s-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail without project_id")
}
