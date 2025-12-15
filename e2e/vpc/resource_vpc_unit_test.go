package vpc

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVPCImportSeparator tests the import separator constant
func TestVPCImportSeparator(t *testing.T) {
	if vpcImportSeparator != ":" {
		t.Errorf("vpcImportSeparator = %v, want ':'", vpcImportSeparator)
	}
}

// TestVPCImportFormat tests the import format constant
func TestVPCImportFormat(t *testing.T) {
	if vpcImportFormat != "project_id:vpc_id" {
		t.Errorf("vpcImportFormat = %v, want 'project_id:vpc_id'", vpcImportFormat)
	}
}

// TestVPCImportIDParsing tests the import ID parsing logic
func TestVPCImportIDParsing(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		expectedParts int
		expectedError bool
		errorContains string
	}{
		{
			name:          "valid import - single vpc_id",
			importID:      "123",
			expectedParts: 1,
			expectedError: false,
		},
		{
			name:          "valid import - project_id:vpc_id format",
			importID:      "project-123:456",
			expectedParts: 2,
			expectedError: false,
		},
		{
			name:          "invalid import - too many parts",
			importID:      "project-123:456:789",
			expectedParts: 3,
			expectedError: true,
			errorContains: errVPCImportFormat,
		},
		{
			name:          "valid import - empty string",
			importID:      "",
			expectedParts: 1, // Empty string splits to single element
			expectedError: false,
		},
		{
			name:          "valid import - vpc_id with special characters",
			importID:      "vpc-123-abc",
			expectedParts: 1,
			expectedError: false,
		},
		{
			name:          "valid import - project_id:vpc_id with special characters",
			importID:      "project-123:vpc-456-xyz",
			expectedParts: 2,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.importID, vpcImportSeparator)

			if len(parts) != tt.expectedParts {
				t.Errorf("Split(%q, %q) = %d parts, want %d", tt.importID, vpcImportSeparator, len(parts), tt.expectedParts)
			}

			// Test error condition for too many parts
			if len(parts) > 2 && !tt.expectedError {
				t.Errorf("expected error for import ID with %d parts: %q", len(parts), tt.importID)
			}

			// Test that separator constant is used
			if vpcImportSeparator != ":" {
				t.Errorf("vpcImportSeparator constant changed, test may be invalid")
			}
		})
	}
}

// ============================================================================
// Mock implementations for VPC tests
// ============================================================================

// mockVpcService is a mock implementation of VpcService for testing
type mockVpcService struct {
	createVPCFunc func(ctx context.Context, req *goe2e.VpcCreateRequest) (*goe2e.Vpc, *goe2e.Response, error)
	getVPCFunc    func(ctx context.Context, vpcID string) (*goe2e.Vpc, *goe2e.Response, error)
	listVPCsFunc  func(ctx context.Context) ([]goe2e.Vpc, *goe2e.Response, error)
	deleteVPCFunc func(ctx context.Context, vpcID string) (*goe2e.Response, error)
}

func (m *mockVpcService) CreateVPC(ctx context.Context, req *goe2e.VpcCreateRequest) (*goe2e.Vpc, *goe2e.Response, error) {
	if m.createVPCFunc != nil {
		return m.createVPCFunc(ctx, req)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockVpcService) GetVPC(ctx context.Context, vpcID string) (*goe2e.Vpc, *goe2e.Response, error) {
	if m.getVPCFunc != nil {
		return m.getVPCFunc(ctx, vpcID)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockVpcService) DeleteVPC(ctx context.Context, vpcID string) (*goe2e.Response, error) {
	if m.deleteVPCFunc != nil {
		return m.deleteVPCFunc(ctx, vpcID)
	}
	return nil, errors.New("not implemented")
}

// ListVPCs is implemented in datasource_vpcs_unit_test.go

func (m *mockVpcService) GetVPCByName(ctx context.Context, name string) (*goe2e.VpcWithSubnets, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

// createMockConfig creates a config with a mock VPC service
func createMockConfig(t *testing.T, mockService *mockVpcService, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.Vpcs = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// createTestResourceData creates a schema.ResourceData for testing
func createTestResourceData(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	resource := ResourceVpc()
	d := schema.TestResourceDataRaw(t, resource.Schema, data)
	return d
}

// ============================================================================
// Test: ResourceCreateVpc
// ============================================================================

func TestResourceCreateVpc_Success(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		setupMock     func() *mockVpcService
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name: "successful create with E2E-managed VPC",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-vpc",
				tfconstants.AttrIsE2EVPC:  true,
				tfconstants.AttrProjectID: "test-project",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockVpcService {
				return &mockVpcService{
					createVPCFunc: func(ctx context.Context, req *goe2e.VpcCreateRequest) (*goe2e.Vpc, *goe2e.Response, error) {
						assert.Equal(t, "test-vpc", req.VpcName)
						assert.Equal(t, "", req.IPv4)
						assert.True(t, req.IsE2EVpc)
						return &goe2e.Vpc{
							ID:        12345.0,
							Name:      "test-vpc",
							State:     "active",
							CreatedAt: "2024-01-01T00:00:00Z",
							IPv4CIDR:  "10.0.0.0/16",
							GatewayIP: "10.0.0.1",
							PoolSize:  65536.0,
							IsActive:  true,
						}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
					},
					getVPCFunc: func(ctx context.Context, vpcID string) (*goe2e.Vpc, *goe2e.Response, error) {
						return &goe2e.Vpc{
							ID:        12345.0,
							Name:      "test-vpc",
							State:     "active",
							CreatedAt: "2024-01-01T00:00:00Z",
							IPv4CIDR:  "10.0.0.0/16",
							GatewayIP: "10.0.0.1",
							PoolSize:  65536.0,
							IsActive:  true,
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "12345", d.Id())
				assert.Equal(t, "test-vpc", d.Get(tfconstants.AttrName))
				assert.Equal(t, 12345.0, d.Get(tfconstants.AttrNetworkID))
				assert.Equal(t, "active", d.Get(tfconstants.AttrStatus))
				assert.Equal(t, "10.0.0.0/16", d.Get(tfconstants.AttrIPv4CIDR))
				assert.Equal(t, "10.0.0.1", d.Get(tfconstants.AttrGatewayIP))
				assert.True(t, d.Get(tfconstants.AttrIsActive).(bool))
			},
		},
		{
			name: "successful create with custom IPv4 CIDR",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-vpc-custom",
				tfconstants.AttrIPv4:      "192.168.0.0/16",
				tfconstants.AttrIsE2EVPC:  false,
				tfconstants.AttrProjectID: "test-project",
				tfconstants.AttrRegion:    "us-west-1",
			},
			setupMock: func() *mockVpcService {
				return &mockVpcService{
					createVPCFunc: func(ctx context.Context, req *goe2e.VpcCreateRequest) (*goe2e.Vpc, *goe2e.Response, error) {
						assert.Equal(t, "test-vpc-custom", req.VpcName)
						assert.Equal(t, "192.168.0.0/16", req.IPv4)
						assert.False(t, req.IsE2EVpc)
						return &goe2e.Vpc{
							ID:        67890.0,
							Name:      "test-vpc-custom",
							State:     "active",
							CreatedAt: "2024-01-02T00:00:00Z",
							IPv4CIDR:  "192.168.0.0/16",
							GatewayIP: "192.168.0.1",
							PoolSize:  65536.0,
							IsActive:  true,
						}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
					},
					getVPCFunc: func(ctx context.Context, vpcID string) (*goe2e.Vpc, *goe2e.Response, error) {
						return &goe2e.Vpc{
							ID:        67890.0,
							Name:      "test-vpc-custom",
							State:     "active",
							CreatedAt: "2024-01-02T00:00:00Z",
							IPv4CIDR:  "192.168.0.0/16",
							GatewayIP: "192.168.0.1",
							PoolSize:  65536.0,
							IsActive:  true,
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "67890", d.Id())
				assert.Equal(t, "192.168.0.0/16", d.Get(tfconstants.AttrIPv4CIDR))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
			resource := ResourceVpc()
			d := createTestResourceData(t, tt.resourceData)

			diags := resource.CreateContext(context.Background(), d, cfg)

			require.False(t, diags.HasError(), "Create should succeed")
			if tt.validateState != nil {
				tt.validateState(t, d)
			}
		})
	}
}

func TestResourceCreateVpc_Errors(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		setupMock     func() *mockVpcService
		expectedError bool
		errorContains string
	}{
		{
			name: "error - API create failure",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-vpc",
				tfconstants.AttrProjectID: "test-project",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockVpcService {
				return &mockVpcService{
					createVPCFunc: func(ctx context.Context, req *goe2e.VpcCreateRequest) (*goe2e.Vpc, *goe2e.Response, error) {
						return nil, nil, errors.New("API error: failed to create VPC")
					},
				}
			},
			expectedError: true,
			errorContains: "Error creating VPC",
		},
		{
			name: "error - missing project_id",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:   "test-vpc",
				tfconstants.AttrRegion: "us-east-1",
			},
			setupMock: func() *mockVpcService {
				return &mockVpcService{}
			},
			expectedError: true,
			errorContains: "project_id",
		},
		{
			name: "error - missing region",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-vpc",
				tfconstants.AttrProjectID: "test-project",
			},
			setupMock: func() *mockVpcService {
				return &mockVpcService{}
			},
			expectedError: true,
			errorContains: "region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "", "") // No defaults
			resource := ResourceVpc()
			d := createTestResourceData(t, tt.resourceData)

			diags := resource.CreateContext(context.Background(), d, cfg)

			if tt.expectedError {
				require.True(t, diags.HasError(), "Create should fail")
				if tt.errorContains != "" {
					errorMsg := diags[0].Summary + " " + diags[0].Detail
					assert.Contains(t, errorMsg, tt.errorContains)
				}
			} else {
				require.False(t, diags.HasError(), "Create should succeed")
			}
		})
	}
}

// ============================================================================
// Test: ResourceReadVpc
// ============================================================================

func TestResourceReadVpc_Success(t *testing.T) {
	mockService := &mockVpcService{
		getVPCFunc: func(ctx context.Context, vpcID string) (*goe2e.Vpc, *goe2e.Response, error) {
			assert.Equal(t, "12345", vpcID)
			return &goe2e.Vpc{
				ID:        12345.0,
				Name:      "test-vpc",
				State:     "active",
				CreatedAt: "2024-01-01T00:00:00Z",
				IPv4CIDR:  "10.0.0.0/16",
				GatewayIP: "10.0.0.1",
				PoolSize:  65536.0,
				IsActive:  true,
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceVpc()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-vpc",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "12345", d.Id())
	assert.Equal(t, "test-vpc", d.Get(tfconstants.AttrName))
	assert.Equal(t, 12345.0, d.Get(tfconstants.AttrNetworkID))
	assert.Equal(t, "active", d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "10.0.0.0/16", d.Get(tfconstants.AttrIPv4CIDR))
	assert.Equal(t, "10.0.0.1", d.Get(tfconstants.AttrGatewayIP))
	assert.True(t, d.Get(tfconstants.AttrIsActive).(bool))
}

func TestResourceReadVpc_APIError(t *testing.T) {
	mockService := &mockVpcService{
		getVPCFunc: func(ctx context.Context, vpcID string) (*goe2e.Vpc, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to retrieve VPC")
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceVpc()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-vpc",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	assert.Contains(t, diags[0].Summary, "Error retrieving VPC")
}

// ============================================================================
// Test: ResourceDeleteVpc
// ============================================================================

func TestResourceDeleteVpc_Success(t *testing.T) {
	mockService := &mockVpcService{
		deleteVPCFunc: func(ctx context.Context, vpcID string) (*goe2e.Response, error) {
			assert.Equal(t, "12345", vpcID)
			return &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceVpc()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-vpc",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Delete should succeed")
	assert.Empty(t, d.Id(), "ID should be cleared after delete")
}

func TestResourceDeleteVpc_APIError(t *testing.T) {
	mockService := &mockVpcService{
		deleteVPCFunc: func(ctx context.Context, vpcID string) (*goe2e.Response, error) {
			return nil, errors.New("API error: failed to delete VPC")
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceVpc()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-vpc",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Delete should fail on API error")
	assert.Contains(t, diags[0].Summary, "Error deleting VPC")
}

// ============================================================================
// Test: VPC Import
// ============================================================================

func TestResourceVpcImport_Success(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name:     "single vpc_id format",
			importID: "12345",
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "12345", d.Id())
			},
		},
		{
			name:     "project_id:vpc_id format",
			importID: "test-project:67890",
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "67890", d.Id())
				assert.Equal(t, "test-project", d.Get(tfconstants.AttrProjectID))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createMockConfig(t, &mockVpcService{}, "test-project", "us-east-1")
			resource := ResourceVpc()
			d := createTestResourceData(t, map[string]interface{}{})
			d.SetId(tt.importID)

			result, err := resource.Importer.StateContext(context.Background(), d, cfg)

			require.NoError(t, err, "Import should succeed")
			require.Len(t, result, 1, "Should return one resource")
			if tt.validateState != nil {
				tt.validateState(t, result[0])
			}
		})
	}
}

func TestResourceVpcImport_InvalidFormat(t *testing.T) {
	cfg := createMockConfig(t, &mockVpcService{}, "test-project", "us-east-1")
	resource := ResourceVpc()
	d := createTestResourceData(t, map[string]interface{}{})
	d.SetId("project:region:vpc:extra")

	result, err := resource.Importer.StateContext(context.Background(), d, cfg)

	require.Error(t, err, "Import should fail")
	assert.Contains(t, err.Error(), "invalid import format")
	assert.Nil(t, result)
}
