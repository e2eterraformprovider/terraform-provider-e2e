package node

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Mock implementations for Node datasource tests
// ============================================================================

// mockNodeServiceForDatasource is a mock implementation of NodeService for datasource testing
type mockNodeServiceForDatasource struct {
	getNodeFunc func(ctx context.Context, nodeID string) (*goe2e.Node, *goe2e.Response, error)
}

func (m *mockNodeServiceForDatasource) GetNode(ctx context.Context, nodeID string) (*goe2e.Node, *goe2e.Response, error) {
	if m.getNodeFunc != nil {
		return m.getNodeFunc(ctx, nodeID)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) ListNodes(ctx context.Context) ([]goe2e.Node, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

// Unused interface methods for datasource tests
func (m *mockNodeServiceForDatasource) CreateNode(ctx context.Context, req *goe2e.NodeCreateRequest) (*goe2e.Node, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) UpdateNode(ctx context.Context, nodeID string, req *goe2e.NodeUpdateRequest) (*goe2e.Node, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) DeleteNode(ctx context.Context, nodeID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) PowerOn(ctx context.Context, nodeID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) PowerOff(ctx context.Context, nodeID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) Reboot(ctx context.Context, nodeID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) Reinstall(ctx context.Context, nodeID string, req *goe2e.NodeReinstallRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) SaveImage(ctx context.Context, nodeID string, req *goe2e.NodeSaveImageRequest) (*goe2e.NodeSaveImageResult, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) LockNode(ctx context.Context, nodeID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) UnlockNode(ctx context.Context, nodeID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) AttachSecurityGroup(ctx context.Context, nodeID string, req *goe2e.SecurityGroupRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) DetachSecurityGroup(ctx context.Context, nodeID string, req *goe2e.SecurityGroupRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) GetSecurityGroupList(ctx context.Context) ([]goe2e.SecurityGroupInfo, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) AttachVPC(ctx context.Context, nodeID string, req *goe2e.NodeVPCAttachRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) DetachVPC(ctx context.Context, nodeID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) GetLCMState(ctx context.Context, nodeID string) (*goe2e.NodeLCMState, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) UpgradePlan(ctx context.Context, nodeID string, req *goe2e.NodePlanUpgradeRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockNodeServiceForDatasource) UpdateSSH(ctx context.Context, nodeID string, req *goe2e.SSHUpdateRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

// createMockConfigForNodeDatasource creates a config with a mock node service
func createMockConfigForNodeDatasource(t *testing.T, mockService *mockNodeServiceForDatasource, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.Nodes = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// ============================================================================
// Test: dataSourceReadNode
// ============================================================================

func TestDataSourceReadNode_Success(t *testing.T) {
	mockService := &mockNodeServiceForDatasource{
		getNodeFunc: func(ctx context.Context, nodeID string) (*goe2e.Node, *goe2e.Response, error) {
			assert.Equal(t, "node-123", nodeID)
			return &goe2e.Node{
				ID:                    "node-123",
				Name:                  "test-node",
				Label:                 "test-label",
				Plan:                  "c2-2c-4gb",
				CreatedAt:             "2024-01-01T00:00:00Z",
				Memory:                "4096",
				Status:                goe2econstants.NodeStatusRunning,
				Disk:                  "80GB",
				Price:                 "0.10",
				PublicIPAddress:       "203.0.113.1",
				PrivateIPAddress:      "10.0.0.1",
				IsLocked:              false,
				BitNinjaLicenseActive: false,
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForNodeDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceNode()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"node_id":    "node-123",
		"region":     "us-east-1",
		"project_id": "test-project",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "node-123", d.Id())
	assert.Equal(t, "test-node", d.Get(tfconstants.AttrName))
	assert.Equal(t, "test-label", d.Get(tfconstants.AttrLabel))
	assert.Equal(t, "c2-2c-4gb", d.Get(tfconstants.AttrPlan))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get(tfconstants.AttrCreatedAt))
	assert.Equal(t, "4096", d.Get(tfconstants.AttrMemory))
	assert.Equal(t, goe2econstants.NodeStatusRunning, d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "80GB", d.Get(tfconstants.AttrDisk))
	assert.Equal(t, "0.10", d.Get("price"))
	assert.Equal(t, "203.0.113.1", d.Get(tfconstants.AttrPublicIPAddress))
	assert.Equal(t, "10.0.0.1", d.Get(tfconstants.AttrPrivateIPAddress))
	assert.False(t, d.Get(tfconstants.AttrIsLocked).(bool))
	assert.False(t, d.Get("is_bitninja_license_active").(bool))
}

func TestDataSourceReadNode_NotFound(t *testing.T) {
	mockService := &mockNodeServiceForDatasource{
		getNodeFunc: func(ctx context.Context, nodeID string) (*goe2e.Node, *goe2e.Response, error) {
			return nil, nil, fmt.Errorf("node with ID %s not found", nodeID)
		},
	}

	cfg := createMockConfigForNodeDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceNode()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"node_id":    "non-existent-node",
		"region":     "us-east-1",
		"project_id": "test-project",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when node not found")
	assert.Contains(t, diags[0].Summary, "error finding Item")
}

func TestDataSourceReadNode_NilNode(t *testing.T) {
	mockService := &mockNodeServiceForDatasource{
		getNodeFunc: func(ctx context.Context, nodeID string) (*goe2e.Node, *goe2e.Response, error) {
			return nil, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForNodeDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceNode()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"node_id":    "node-123",
		"region":     "us-east-1",
		"project_id": "test-project",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when node is nil")
	assert.Contains(t, diags[0].Summary, "node not found")
}

func TestDataSourceReadNode_APIError(t *testing.T) {
	mockService := &mockNodeServiceForDatasource{
		getNodeFunc: func(ctx context.Context, nodeID string) (*goe2e.Node, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to retrieve node")
		},
	}

	cfg := createMockConfigForNodeDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceNode()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"node_id":    "node-123",
		"region":     "us-east-1",
		"project_id": "test-project",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	assert.Contains(t, diags[0].Summary, "error finding Item")
}

func TestDataSourceReadNode_NotFoundError(t *testing.T) {
	mockService := &mockNodeServiceForDatasource{
		getNodeFunc: func(ctx context.Context, nodeID string) (*goe2e.Node, *goe2e.Response, error) {
			err := fmt.Errorf("error: %s", goe2econstants.NotFoundSubstring)
			return nil, nil, err
		},
	}

	cfg := createMockConfigForNodeDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceNode()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"node_id":    "node-123",
		"region":     "us-east-1",
		"project_id": "test-project",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when node not found")
	assert.Contains(t, diags[0].Summary, "error finding Item")
}

func TestDataSourceReadNode_MissingRegion(t *testing.T) {
	cfg := createMockConfigForNodeDatasource(t, &mockNodeServiceForDatasource{}, "", "")
	resource := DataSourceNode()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"node_id": "node-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail without region")
	assert.Contains(t, diags[0].Summary, "region")
}

func TestDataSourceReadNode_MissingProjectID(t *testing.T) {
	cfg := createMockConfigForNodeDatasource(t, &mockNodeServiceForDatasource{}, "", "us-east-1")
	resource := DataSourceNode()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"node_id": "node-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail without project_id")
	assert.Contains(t, diags[0].Summary, "project_id")
}
