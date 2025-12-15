package faas

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
// Mock implementations for FaaS datasource tests
// ============================================================================

// mockFAASServiceForDatasource is a mock implementation of FaasService for datasource testing
type mockFAASServiceForDatasource struct {
	getFunctionFunc func(ctx context.Context, functionID string) (*goe2e.FaasFunction, *goe2e.Response, error)
}

func (m *mockFAASServiceForDatasource) GetFunction(ctx context.Context, functionID string) (*goe2e.FaasFunction, *goe2e.Response, error) {
	if m.getFunctionFunc != nil {
		return m.getFunctionFunc(ctx, functionID)
	}
	return nil, nil, errors.New("not implemented")
}

// Unused interface methods
func (m *mockFAASServiceForDatasource) CreateNamespace(ctx context.Context, name string) (*goe2e.FaasNamespace, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockFAASServiceForDatasource) DeleteNamespace(ctx context.Context, name string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockFAASServiceForDatasource) CreateFunction(ctx context.Context, req *goe2e.FaasFunctionCreateRequest) (*goe2e.FaasFunction, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockFAASServiceForDatasource) UpdateFunction(ctx context.Context, functionID string, req *goe2e.FaasFunctionUpdateRequest) (*goe2e.FaasFunction, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockFAASServiceForDatasource) DeleteFunction(ctx context.Context, functionID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockFAASServiceForDatasource) GetLogs(ctx context.Context, functionID string) (*goe2e.FaasLogs, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

// createMockConfigForFAASDatasource creates a config with a mock FaaS service
func createMockConfigForFAASDatasource(t *testing.T, mockService *mockFAASServiceForDatasource) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	mockClient := &goe2e.Client{}
	mockClient.FaaS = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// ============================================================================
// Test: dataSourceFaasFunctionRead
// ============================================================================

func TestDataSourceFaasFunctionRead_Success(t *testing.T) {
	mockService := &mockFAASServiceForDatasource{
		getFunctionFunc: func(ctx context.Context, functionID string) (*goe2e.FaasFunction, *goe2e.Response, error) {
			assert.Equal(t, "func-123", functionID)
			return &goe2e.FaasFunction{
				ID:          "func-123",
				Name:        "test-function",
				Namespace:   "default",
				Runtime:     "python3.9",
				MemoryMB:    256,
				Timeout:     30,
				MinReplicas: 1,
				MaxReplicas: 5,
				EndpointURL: "https://api.example.com/functions/test-function",
				Status:      "active",
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-02T00:00:00Z",
				Environment: map[string]string{
					"ENV_VAR_1": "value1",
					"ENV_VAR_2": "value2",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForFAASDatasource(t, mockService)
	resource := DataSourceFaasFunction()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrFunctionID: "func-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "func-123", d.Id())
	assert.Equal(t, "test-function", d.Get(tfconstants.AttrName))
	assert.Equal(t, "default", d.Get(tfconstants.AttrNamespace))
	assert.Equal(t, "python3.9", d.Get(tfconstants.AttrRuntime))
	assert.Equal(t, 256, d.Get(tfconstants.AttrMemoryMB))
	assert.Equal(t, 30, d.Get(tfconstants.AttrTimeoutSeconds))
	assert.Equal(t, 1, d.Get(tfconstants.AttrMinReplicas))
	assert.Equal(t, 5, d.Get(tfconstants.AttrMaxReplicas))
	assert.Equal(t, "https://api.example.com/functions/test-function", d.Get(tfconstants.AttrEndpointURL))
	assert.Equal(t, "active", d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get(tfconstants.AttrCreatedAt))
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get(tfconstants.AttrUpdatedAt))

	env := d.Get("environment").(map[string]interface{})
	assert.Equal(t, "value1", env["ENV_VAR_1"])
	assert.Equal(t, "value2", env["ENV_VAR_2"])
}

func TestDataSourceFaasFunctionRead_NoEnvironment(t *testing.T) {
	mockService := &mockFAASServiceForDatasource{
		getFunctionFunc: func(ctx context.Context, functionID string) (*goe2e.FaasFunction, *goe2e.Response, error) {
			return &goe2e.FaasFunction{
				ID:          "func-456",
				Name:        "test-function-no-env",
				Namespace:   "default",
				Runtime:     "nodejs18",
				MemoryMB:    128,
				Timeout:     60,
				MinReplicas: 0,
				MaxReplicas: 3,
				Status:      "active",
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   "2024-01-01T00:00:00Z",
				Environment: nil,
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForFAASDatasource(t, mockService)
	resource := DataSourceFaasFunction()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrFunctionID: "func-456",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	// Environment should not be set when nil
	assert.False(t, d.Get("environment").(map[string]interface{}) != nil && len(d.Get("environment").(map[string]interface{})) > 0)
}

func TestDataSourceFaasFunctionRead_NotFound(t *testing.T) {
	mockService := &mockFAASServiceForDatasource{
		getFunctionFunc: func(ctx context.Context, functionID string) (*goe2e.FaasFunction, *goe2e.Response, error) {
			return nil, nil, errors.New("function not found")
		},
	}

	cfg := createMockConfigForFAASDatasource(t, mockService)
	resource := DataSourceFaasFunction()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrFunctionID: "non-existent-func",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when function not found")
}

func TestDataSourceFaasFunctionRead_NilFunction(t *testing.T) {
	mockService := &mockFAASServiceForDatasource{
		getFunctionFunc: func(ctx context.Context, functionID string) (*goe2e.FaasFunction, *goe2e.Response, error) {
			return nil, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForFAASDatasource(t, mockService)
	resource := DataSourceFaasFunction()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrFunctionID: "func-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when function is nil")
	assert.Contains(t, diags[0].Summary, "not found")
}

func TestDataSourceFaasFunctionRead_APIError(t *testing.T) {
	mockService := &mockFAASServiceForDatasource{
		getFunctionFunc: func(ctx context.Context, functionID string) (*goe2e.FaasFunction, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to retrieve function")
		},
	}

	cfg := createMockConfigForFAASDatasource(t, mockService)
	resource := DataSourceFaasFunction()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrFunctionID: "func-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
}
