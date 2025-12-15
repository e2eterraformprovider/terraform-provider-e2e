package autoscaling

import (
	"context"
	"errors"
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
// Mock implementations for Autoscaling datasource tests
// ============================================================================

// mockAutoscalingServiceForDatasource is a mock implementation of AutoscalingService for datasource testing
type mockAutoscalingServiceForDatasource struct {
	getScalerGroupFunc func(ctx context.Context, scalerID string) (*goe2e.ScalerGroup, *goe2e.Response, error)
}

func (m *mockAutoscalingServiceForDatasource) GetScalerGroup(ctx context.Context, scalerID string) (*goe2e.ScalerGroup, *goe2e.Response, error) {
	if m.getScalerGroupFunc != nil {
		return m.getScalerGroupFunc(ctx, scalerID)
	}
	return nil, nil, errors.New("not implemented")
}

// Unused interface methods
func (m *mockAutoscalingServiceForDatasource) CreateScalerGroup(ctx context.Context, req *goe2e.ScalerGroupCreateRequest) (*goe2e.ScalerGroup, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) UpdateScalerGroup(ctx context.Context, scalerID string, req *goe2e.ScalerGroupUpdateRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) DeleteScalerGroup(ctx context.Context, scalerID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) ListScalerGroups(ctx context.Context) ([]goe2e.ScalerGroup, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) StartScalerGroup(ctx context.Context, scalerID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) StopScalerGroup(ctx context.Context, scalerID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) UpdateScalerGroupStatus(ctx context.Context, scalerID string, status string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) UpdateDesiredNodeCount(ctx context.Context, scalerID string, count int) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) AttachVPCToScalerGroup(ctx context.Context, scalerID string, req *goe2e.VPCAttachRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) DetachVPCFromScalerGroup(ctx context.Context, scalerID string, vpcName string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) GetAttachedVPCsForScalerGroup(ctx context.Context, scalerID string) ([]goe2e.VPCPartial, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) AttachSecurityGroupToScalerGroup(ctx context.Context, scalerID string, sgID int) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) DetachSecurityGroupFromScalerGroup(ctx context.Context, scalerID string, sgID int) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) AttachPublicIPToScalerGroup(ctx context.Context, scalerID string) (*goe2e.PublicIPActionResponse, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) DetachPublicIPFromScalerGroup(ctx context.Context, scalerID string) (*goe2e.PublicIPActionResponse, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockAutoscalingServiceForDatasource) GetPublicIPStatus(ctx context.Context, scalerID string) (*goe2e.PublicIPStatus, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

// createMockConfigForAutoscalingDatasource creates a config with a mock autoscaling service
func createMockConfigForAutoscalingDatasource(t *testing.T, mockService *mockAutoscalingServiceForDatasource, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.Autoscaling = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// ============================================================================
// Test: dataSourceReadScalerGroup
// ============================================================================

func TestDataSourceReadScalerGroup_Success(t *testing.T) {
	mockService := &mockAutoscalingServiceForDatasource{
		getScalerGroupFunc: func(ctx context.Context, scalerID string) (*goe2e.ScalerGroup, *goe2e.Response, error) {
			assert.Equal(t, "scaler-123", scalerID)
			return &goe2e.ScalerGroup{
				ID:                      "scaler-123",
				Name:                    "test-scaler-group",
				Desired:                 3,
				MinNodes:                1,
				MaxNodes:                10,
				PlanName:                "c2-2c-4gb",
				VMImageName:             "ubuntu-20.04",
				ProvisionStatus:         goe2econstants.AutoscalingScalerGroupStatusRunning,
				PolicyType:              "elastic",
				PolicyMeasure:           "cpu",
				UpscalePolicyValue:      2,
				DownscalePolicyValue:    1,
				PolicyUpscaleOperator:   ">",
				PolicyDownscaleOperator: "<",
				WaitForPeriod:           5,
				Cooldown:                300,
				ScheduledPolicyOp:       "scale",
				UpscaleAdjust:           2,
				DownscaleAdjust:         1,
				UpscaleRecurrence:       "0 9 * * *",
				DownscaleRecurrence:     "0 18 * * *",
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForAutoscalingDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceScalerGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "scaler-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "scaler-123", d.Id())
	assert.Equal(t, "test-scaler-group", d.Get(tfconstants.AttrName))
	assert.Equal(t, 3, d.Get(tfconstants.AttrDesired))
	assert.Equal(t, 1, d.Get(tfconstants.AttrMinNodes))
	assert.Equal(t, 10, d.Get(tfconstants.AttrMaxNodes))
	assert.Equal(t, "c2-2c-4gb", d.Get(tfconstants.AttrPlan))
	assert.Equal(t, "ubuntu-20.04", d.Get("vm_image_name"))
	assert.Equal(t, goe2econstants.AutoscalingScalerGroupStatusRunning, d.Get("provision_status"))
	assert.Equal(t, "elastic", d.Get("policy_type"))

	// Validate policy list
	policyList := d.Get("policy").([]interface{})
	require.Len(t, policyList, 2)
	upscalePolicy := policyList[0].(map[string]interface{})
	assert.Equal(t, "elastic", upscalePolicy["type"])
	assert.Equal(t, 2, upscalePolicy["adjust"])
	assert.Equal(t, "cpu", upscalePolicy["parameter"])

	// Validate scheduled policy list
	scheduledPolicyList := d.Get("scheduled_policy").([]interface{})
	require.Len(t, scheduledPolicyList, 2)
}

func TestDataSourceReadScalerGroup_NotFound(t *testing.T) {
	mockService := &mockAutoscalingServiceForDatasource{
		getScalerGroupFunc: func(ctx context.Context, scalerID string) (*goe2e.ScalerGroup, *goe2e.Response, error) {
			return nil, nil, errors.New("scaler group not found")
		},
	}

	cfg := createMockConfigForAutoscalingDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceScalerGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "non-existent-scaler",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when scaler group not found")
	assert.Contains(t, diags[0].Summary, "failed to read scaler group")
}

func TestDataSourceReadScalerGroup_NilGroup(t *testing.T) {
	mockService := &mockAutoscalingServiceForDatasource{
		getScalerGroupFunc: func(ctx context.Context, scalerID string) (*goe2e.ScalerGroup, *goe2e.Response, error) {
			return nil, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForAutoscalingDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceScalerGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "scaler-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when scaler group is nil")
	assert.Contains(t, diags[0].Summary, "not found")
}

func TestDataSourceReadScalerGroup_StatusNormalization(t *testing.T) {
	mockService := &mockAutoscalingServiceForDatasource{
		getScalerGroupFunc: func(ctx context.Context, scalerID string) (*goe2e.ScalerGroup, *goe2e.Response, error) {
			return &goe2e.ScalerGroup{
				ID:              "scaler-456",
				Name:            "test-scaler",
				ProvisionStatus: goe2econstants.AutoscalingScalerGroupStatusStarting,
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForAutoscalingDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceScalerGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "scaler-456",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	// Status should be normalized from "Starting" to "Running"
	assert.Equal(t, goe2econstants.AutoscalingScalerGroupStatusRunning, d.Get("provision_status"))
}

func TestDataSourceReadScalerGroup_MissingRegion(t *testing.T) {
	cfg := createMockConfigForAutoscalingDatasource(t, &mockAutoscalingServiceForDatasource{}, "", "")
	resource := DataSourceScalerGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "scaler-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail without region")
}

func TestDataSourceReadScalerGroup_MissingProjectID(t *testing.T) {
	cfg := createMockConfigForAutoscalingDatasource(t, &mockAutoscalingServiceForDatasource{}, "", "us-east-1")
	resource := DataSourceScalerGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "scaler-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail without project_id")
}
