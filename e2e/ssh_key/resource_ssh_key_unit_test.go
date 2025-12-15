package ssh_key

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
// Mock implementations for SSH Key tests
// ============================================================================

// mockSSHKeyService is a mock implementation of SSHKeyService for testing
type mockSSHKeyService struct {
	createSSHKeyFunc     func(ctx context.Context, req *goe2e.SSHKeyCreateRequest) (*goe2e.SSHKey, *goe2e.Response, error)
	getSSHKeyFunc        func(ctx context.Context, pk string) (*goe2e.SSHKey, *goe2e.Response, error)
	getSSHKeyByLabelFunc func(ctx context.Context, label string) (*goe2e.SSHKey, *goe2e.Response, error)
	listSSHKeysFunc      func(ctx context.Context) ([]goe2e.SSHKey, *goe2e.Response, error)
	deleteSSHKeyFunc     func(ctx context.Context, pk string) (*goe2e.Response, error)
}

func (m *mockSSHKeyService) CreateSSHKey(ctx context.Context, req *goe2e.SSHKeyCreateRequest) (*goe2e.SSHKey, *goe2e.Response, error) {
	if m.createSSHKeyFunc != nil {
		return m.createSSHKeyFunc(ctx, req)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockSSHKeyService) GetSSHKey(ctx context.Context, pk string) (*goe2e.SSHKey, *goe2e.Response, error) {
	if m.getSSHKeyFunc != nil {
		return m.getSSHKeyFunc(ctx, pk)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockSSHKeyService) GetSSHKeyByLabel(ctx context.Context, label string) (*goe2e.SSHKey, *goe2e.Response, error) {
	if m.getSSHKeyByLabelFunc != nil {
		return m.getSSHKeyByLabelFunc(ctx, label)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockSSHKeyService) ListSSHKeys(ctx context.Context) ([]goe2e.SSHKey, *goe2e.Response, error) {
	if m.listSSHKeysFunc != nil {
		return m.listSSHKeysFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockSSHKeyService) DeleteSSHKey(ctx context.Context, pk string) (*goe2e.Response, error) {
	if m.deleteSSHKeyFunc != nil {
		return m.deleteSSHKeyFunc(ctx, pk)
	}
	return nil, errors.New("not implemented")
}

// createMockConfig creates a config with a mock SSH key service
func createMockConfig(t *testing.T, mockService *mockSSHKeyService, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.SSHKeys = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// createTestResourceData creates a schema.ResourceData for testing
func createTestResourceData(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	resource := ResourceSshKey()
	d := schema.TestResourceDataRaw(t, resource.Schema, data)
	return d
}

// ============================================================================
// Test: resourceCreateSshKey
// ============================================================================

func TestResourceCreateSshKey_Success(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		setupMock     func() *mockSSHKeyService
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name: "successful create with V3 preferred fields",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-key",
				tfconstants.AttrPublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{
					createSSHKeyFunc: func(ctx context.Context, req *goe2e.SSHKeyCreateRequest) (*goe2e.SSHKey, *goe2e.Response, error) {
						assert.Equal(t, "test-key", req.Label)
						assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...", req.SSHKey)
						return &goe2e.SSHKey{
							PK:        12345,
							Label:     "test-key",
							SSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
							Timestamp: "2024-01-01T00:00:00Z",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "12345", d.Id())
				assert.Equal(t, "test-key", d.Get(tfconstants.AttrName))
				assert.Equal(t, "test-key", d.Get(tfconstants.AttrLabel))
				assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...", d.Get(tfconstants.AttrPublicKey))
				assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...", d.Get(tfconstants.AttrSSHKey))
				assert.Equal(t, "2024-01-01T00:00:00Z", d.Get(tfconstants.AttrCreatedAt))
			},
		},
		{
			name: "successful create with V2 deprecated fields",
			resourceData: map[string]interface{}{
				tfconstants.AttrLabel:  "test-key-v2",
				tfconstants.AttrSSHKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...",
				tfconstants.AttrRegion: "us-west-1",
			},
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{
					createSSHKeyFunc: func(ctx context.Context, req *goe2e.SSHKeyCreateRequest) (*goe2e.SSHKey, *goe2e.Response, error) {
						return &goe2e.SSHKey{
							PK:        67890,
							Label:     "test-key-v2",
							SSHKey:    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...",
							Timestamp: "2024-01-02T00:00:00Z",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "67890", d.Id())
				assert.Equal(t, "test-key-v2", d.Get(tfconstants.AttrName))
				assert.Equal(t, "test-key-v2", d.Get(tfconstants.AttrLabel))
			},
		},
		{
			name: "successful create with tags",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-key-tagged",
				tfconstants.AttrPublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
				tfconstants.AttrTags: map[string]interface{}{
					"Environment": "test",
					"ManagedBy":   "terraform",
				},
			},
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{
					createSSHKeyFunc: func(ctx context.Context, req *goe2e.SSHKeyCreateRequest) (*goe2e.SSHKey, *goe2e.Response, error) {
						return &goe2e.SSHKey{
							PK:        99999,
							Label:     "test-key-tagged",
							SSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
							Timestamp: "2024-01-03T00:00:00Z",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				tags := d.Get(tfconstants.AttrTags).(map[string]interface{})
				assert.Equal(t, "test", tags["Environment"])
				assert.Equal(t, "terraform", tags["ManagedBy"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
			resource := ResourceSshKey()
			d := createTestResourceData(t, tt.resourceData)

			diags := resource.CreateContext(context.Background(), d, cfg)

			require.False(t, diags.HasError(), "Create should succeed")
			if tt.validateState != nil {
				tt.validateState(t, d)
			}
		})
	}
}

func TestResourceCreateSshKey_Errors(t *testing.T) {
	tests := []struct {
		name           string
		resourceData   map[string]interface{}
		setupMock      func() *mockSSHKeyService
		defaultProject string
		defaultRegion  string
		expectedError  bool
		errorContains  string
	}{
		{
			name: "error - API create failure",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-key",
				tfconstants.AttrPublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
			},
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{
					createSSHKeyFunc: func(ctx context.Context, req *goe2e.SSHKeyCreateRequest) (*goe2e.SSHKey, *goe2e.Response, error) {
						return nil, nil, errors.New("API error: failed to create SSH key")
					},
				}
			},
			defaultProject: "test-project",
			defaultRegion:  "us-east-1", // Provide region so we reach API call
			expectedError:  true,
			errorContains:  "failed to create SSH key",
		},
		{
			name: "error - empty region when required",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-key",
				tfconstants.AttrPublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
			},
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{}
			},
			defaultProject: "",
			defaultRegion:  "", // No defaults to test region requirement
			expectedError:  true,
			errorContains:  "region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, tt.defaultProject, tt.defaultRegion)
			resource := ResourceSshKey()
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
// Test: resourceReadSshKey
// ============================================================================

func TestResourceReadSshKey_Success(t *testing.T) {
	tests := []struct {
		name          string
		resourceID    string
		resourceData  map[string]interface{}
		setupMock     func() *mockSSHKeyService
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name:       "successful read",
			resourceID: "12345",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-key",
				tfconstants.AttrPublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
			},
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{
					getSSHKeyFunc: func(ctx context.Context, pk string) (*goe2e.SSHKey, *goe2e.Response, error) {
						assert.Equal(t, "12345", pk)
						return &goe2e.SSHKey{
							PK:        12345,
							Label:     "test-key",
							SSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
							Timestamp: "2024-01-01T00:00:00Z",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "12345", d.Id())
				assert.Equal(t, "test-key", d.Get(tfconstants.AttrName))
				assert.Equal(t, "test-key", d.Get(tfconstants.AttrLabel))
				assert.Equal(t, "2024-01-01T00:00:00Z", d.Get(tfconstants.AttrCreatedAt))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
			resource := ResourceSshKey()
			d := createTestResourceData(t, tt.resourceData)
			d.SetId(tt.resourceID)

			diags := resource.ReadContext(context.Background(), d, cfg)

			require.False(t, diags.HasError(), "Read should succeed")
			if tt.validateState != nil {
				tt.validateState(t, d)
			}
		})
	}
}

func TestResourceReadSshKey_NotFound(t *testing.T) {
	mockService := &mockSSHKeyService{
		getSSHKeyFunc: func(ctx context.Context, pk string) (*goe2e.SSHKey, *goe2e.Response, error) {
			return nil, nil, fmt.Errorf("SSH key with ID %s not found", pk)
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceSshKey()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-key",
		tfconstants.AttrPublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
	})
	d.SetId("99999")

	diags := resource.ReadContext(context.Background(), d, cfg)

	// Should handle not found gracefully
	require.True(t, diags.HasError() || d.Id() == "", "Read should handle not found")
}

func TestResourceReadSshKey_EmptyID(t *testing.T) {
	cfg := createMockConfig(t, &mockSSHKeyService{}, "test-project", "us-east-1")
	resource := ResourceSshKey()
	d := createTestResourceData(t, map[string]interface{}{})
	d.SetId("")

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail with empty ID")
	// diag.Errorf sets Summary field, not Detail
	assert.Contains(t, diags[0].Summary, "empty")
}

// ============================================================================
// Test: resourceUpdateSshKey
// ============================================================================

func TestResourceUpdateSshKey_Tags(t *testing.T) {
	cfg := createMockConfig(t, &mockSSHKeyService{}, "test-project", "us-east-1")
	resource := ResourceSshKey()

	// Start with tags already set in the resource data
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-key",
		tfconstants.AttrPublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
		tfconstants.AttrTags: map[string]interface{}{
			"Environment": "production",
			"ManagedBy":   "terraform",
		},
	})
	d.SetId("12345")

	// Update simply preserves the tags in state (no API call needed)
	diags := resource.UpdateContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Update should succeed")
	updatedTags := d.Get(tfconstants.AttrTags).(map[string]interface{})
	assert.Equal(t, "production", updatedTags["Environment"])
	assert.Equal(t, "terraform", updatedTags["ManagedBy"])
}

// ============================================================================
// Test: resourceDeleteSshKey
// ============================================================================

func TestResourceDeleteSshKey_Success(t *testing.T) {
	mockService := &mockSSHKeyService{
		deleteSSHKeyFunc: func(ctx context.Context, pk string) (*goe2e.Response, error) {
			assert.Equal(t, "12345", pk)
			return &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceSshKey()
	d := createTestResourceData(t, map[string]interface{}{})
	d.SetId("12345")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Delete should succeed")
	assert.Empty(t, d.Id(), "ID should be cleared after delete")
}

func TestResourceDeleteSshKey_NotFound(t *testing.T) {
	mockService := &mockSSHKeyService{
		deleteSSHKeyFunc: func(ctx context.Context, pk string) (*goe2e.Response, error) {
			err := fmt.Errorf("SSH key with ID %s not found", pk)
			// Simulate not found error
			if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) || strings.Contains(err.Error(), "not found") {
				return nil, err
			}
			return nil, err
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceSshKey()
	d := createTestResourceData(t, map[string]interface{}{})
	d.SetId("99999")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	// Should handle not found gracefully (idempotent delete)
	require.False(t, diags.HasError(), "Delete should handle not found gracefully")
	assert.Empty(t, d.Id(), "ID should be cleared")
}

func TestResourceDeleteSshKey_EmptyID(t *testing.T) {
	cfg := createMockConfig(t, &mockSSHKeyService{}, "test-project", "us-east-1")
	resource := ResourceSshKey()
	d := createTestResourceData(t, map[string]interface{}{})
	d.SetId("")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Delete should fail with empty ID")
	// diag.Errorf sets Summary field, not Detail
	assert.Contains(t, diags[0].Summary, "empty")
}

func TestResourceDeleteSshKey_APIError(t *testing.T) {
	mockService := &mockSSHKeyService{
		deleteSSHKeyFunc: func(ctx context.Context, pk string) (*goe2e.Response, error) {
			return nil, errors.New("API error: failed to delete SSH key")
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceSshKey()
	d := createTestResourceData(t, map[string]interface{}{})
	d.SetId("12345")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Delete should fail on API error")
	// diag.FromErr(fmt.Errorf(...)) sets Summary field
	assert.Contains(t, diags[0].Summary, "failed to delete SSH key")
}

// ============================================================================
// Test: resourceSshKeyImport
// ============================================================================

func TestResourceSshKeyImport_Success(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		setupMock     func() *mockSSHKeyService
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name:     "2-part format: project_id:ssh_key_id",
			importID: "test-project:12345",
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{
					getSSHKeyFunc: func(ctx context.Context, pk string) (*goe2e.SSHKey, *goe2e.Response, error) {
						return &goe2e.SSHKey{
							PK:        12345,
							Label:     "imported-key",
							SSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
							Timestamp: "2024-01-01T00:00:00Z",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "12345", d.Id())
				assert.Equal(t, "test-project", d.Get(tfconstants.AttrProjectID))
				assert.Equal(t, "imported-key", d.Get(tfconstants.AttrName))
			},
		},
		{
			name:     "3-part format: project_id:region:ssh_key_id",
			importID: "test-project:us-west-1:67890",
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{
					getSSHKeyFunc: func(ctx context.Context, pk string) (*goe2e.SSHKey, *goe2e.Response, error) {
						return &goe2e.SSHKey{
							PK:        67890,
							Label:     "imported-key-2",
							SSHKey:    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...",
							Timestamp: "2024-01-02T00:00:00Z",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "67890", d.Id())
				assert.Equal(t, "test-project", d.Get(tfconstants.AttrProjectID))
				assert.Equal(t, "us-west-1", d.Get(tfconstants.AttrRegion))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
			resource := ResourceSshKey()
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

func TestResourceSshKeyImport_Errors(t *testing.T) {
	tests := []struct {
		name           string
		importID       string
		setupMock      func() *mockSSHKeyService
		defaultProject string
		defaultRegion  string
		expectedError  bool
		errorContains  string
	}{
		{
			name:           "invalid format - too many parts",
			importID:       "project:region:key:extra",
			setupMock:      func() *mockSSHKeyService { return &mockSSHKeyService{} },
			defaultProject: "",
			defaultRegion:  "",
			expectedError:  true,
			errorContains:  ImportIDInvalidFormat,
		},
		{
			name:           "invalid format - single part",
			importID:       "12345",
			setupMock:      func() *mockSSHKeyService { return &mockSSHKeyService{} },
			defaultProject: "",
			defaultRegion:  "",
			expectedError:  true,
			errorContains:  ImportIDInvalidFormat, // Single part is not valid (need at least 2)
		},
		{
			name:           "2-part format without default region",
			importID:       "test-project:12345",
			setupMock:      func() *mockSSHKeyService { return &mockSSHKeyService{} },
			defaultProject: "",
			defaultRegion:  "", // No default region for 2-part format
			expectedError:  true,
			errorContains:  ImportIDRegionRequired,
		},
		{
			name:     "API error during fetch",
			importID: "test-project:12345",
			setupMock: func() *mockSSHKeyService {
				return &mockSSHKeyService{
					getSSHKeyFunc: func(ctx context.Context, pk string) (*goe2e.SSHKey, *goe2e.Response, error) {
						return nil, nil, errors.New("API error: failed to fetch SSH key")
					},
				}
			},
			defaultProject: "test-project",
			defaultRegion:  "us-east-1", // Provide region so we reach API call
			expectedError:  true,
			errorContains:  "failed to fetch SSH key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, tt.defaultProject, tt.defaultRegion)
			resource := ResourceSshKey()
			d := createTestResourceData(t, map[string]interface{}{})
			d.SetId(tt.importID)

			result, err := resource.Importer.StateContext(context.Background(), d, cfg)

			if tt.expectedError {
				require.Error(t, err, "Import should fail")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err, "Import should succeed")
				require.NotNil(t, result)
			}
		})
	}
}
