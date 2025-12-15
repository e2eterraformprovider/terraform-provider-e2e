package tag

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
// Mock implementations for Tag tests
// ============================================================================

// mockTagService is a mock implementation of TagService for testing
type mockTagService struct {
	createTagFunc func(ctx context.Context, req *goe2e.TagCreateRequest) (*goe2e.Tag, *goe2e.Response, error)
	getTagFunc    func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error)
	deleteTagFunc func(ctx context.Context, tagID string) (*goe2e.Response, error)
}

func (m *mockTagService) CreateTag(ctx context.Context, req *goe2e.TagCreateRequest) (*goe2e.Tag, *goe2e.Response, error) {
	if m.createTagFunc != nil {
		return m.createTagFunc(ctx, req)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockTagService) GetTag(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
	if m.getTagFunc != nil {
		return m.getTagFunc(ctx, tagID)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockTagService) DeleteTag(ctx context.Context, tagID string) (*goe2e.Response, error) {
	if m.deleteTagFunc != nil {
		return m.deleteTagFunc(ctx, tagID)
	}
	return nil, errors.New("not implemented")
}

// Unused interface methods
func (m *mockTagService) ListTags(ctx context.Context) ([]goe2e.Tag, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockTagService) UpdateTag(ctx context.Context, tagID string, req *goe2e.TagUpdateRequest) (*goe2e.Tag, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockTagService) AttachTags(ctx context.Context, resourceType, resourceID string, tagIDs []int) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockTagService) DetachTags(ctx context.Context, resourceType, resourceID string, tagIDs []int) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockTagService) GetResourceTags(ctx context.Context, resourceType, resourceID string) ([]goe2e.TagMapping, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

// createMockConfig creates a config with a mock tag service
func createMockConfig(t *testing.T, mockService *mockTagService, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.Tags = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// createTestResourceData creates a schema.ResourceData for testing
func createTestResourceData(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	resource := ResourceTag()
	d := schema.TestResourceDataRaw(t, resource.Schema, data)
	return d
}

// ============================================================================
// Test: resourceTagCreate
// ============================================================================

func TestResourceTagCreate_Success(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		setupMock     func() *mockTagService
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name: "successful create with name only",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-tag",
				tfconstants.AttrProjectID: "test-project",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockTagService {
				return &mockTagService{
					createTagFunc: func(ctx context.Context, req *goe2e.TagCreateRequest) (*goe2e.Tag, *goe2e.Response, error) {
						assert.Equal(t, "test-tag", req.LabelName)
						assert.Equal(t, "", req.Metadata)
						return &goe2e.Tag{
							LabelID:   12345,
							LabelName: "test-tag",
							Metadata:  "",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
					},
					getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
						return &goe2e.Tag{
							LabelID:   12345,
							LabelName: "test-tag",
							Metadata:  "",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "12345", d.Id())
				assert.Equal(t, "test-tag", d.Get(tfconstants.AttrName))
				assert.Equal(t, 12345, d.Get(attrLabelID))
				assert.Equal(t, "test-project", d.Get(tfconstants.AttrProjectID))
				assert.Equal(t, "us-east-1", d.Get(tfconstants.AttrRegion))
			},
		},
		{
			name: "successful create with name and metadata",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-tag-with-metadata",
				attrMetadata:              "Test metadata description",
				tfconstants.AttrProjectID: "test-project",
				tfconstants.AttrRegion:    "us-west-1",
			},
			setupMock: func() *mockTagService {
				return &mockTagService{
					createTagFunc: func(ctx context.Context, req *goe2e.TagCreateRequest) (*goe2e.Tag, *goe2e.Response, error) {
						assert.Equal(t, "test-tag-with-metadata", req.LabelName)
						assert.Equal(t, "Test metadata description", req.Metadata)
						return &goe2e.Tag{
							LabelID:   67890,
							LabelName: "test-tag-with-metadata",
							Metadata:  "Test metadata description",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
					},
					getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
						return &goe2e.Tag{
							LabelID:   67890,
							LabelName: "test-tag-with-metadata",
							Metadata:  "Test metadata description",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "67890", d.Id())
				assert.Equal(t, "test-tag-with-metadata", d.Get(tfconstants.AttrName))
				assert.Equal(t, "Test metadata description", d.Get(attrMetadata))
			},
		},
		{
			name: "successful create with empty metadata",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-tag-empty-metadata",
				attrMetadata:              "",
				tfconstants.AttrProjectID: "test-project",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockTagService {
				return &mockTagService{
					createTagFunc: func(ctx context.Context, req *goe2e.TagCreateRequest) (*goe2e.Tag, *goe2e.Response, error) {
						return &goe2e.Tag{
							LabelID:   11111,
							LabelName: "test-tag-empty-metadata",
							Metadata:  "",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
					},
					getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
						return &goe2e.Tag{
							LabelID:   11111,
							LabelName: "test-tag-empty-metadata",
							Metadata:  "",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "", d.Get(attrMetadata))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
			resource := ResourceTag()
			d := createTestResourceData(t, tt.resourceData)

			diags := resource.CreateContext(context.Background(), d, cfg)

			require.False(t, diags.HasError(), "Create should succeed")
			if tt.validateState != nil {
				tt.validateState(t, d)
			}
		})
	}
}

func TestResourceTagCreate_Errors(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		setupMock     func() *mockTagService
		expectedError bool
		errorContains string
	}{
		{
			name: "error - API create failure",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-tag",
				tfconstants.AttrProjectID: "test-project",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockTagService {
				return &mockTagService{
					createTagFunc: func(ctx context.Context, req *goe2e.TagCreateRequest) (*goe2e.Tag, *goe2e.Response, error) {
						return nil, nil, errors.New("API error: failed to create tag")
					},
				}
			},
			expectedError: true,
			errorContains: "Error creating tag",
		},
		{
			name: "error - missing project_id",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:   "test-tag",
				tfconstants.AttrRegion: "us-east-1",
			},
			setupMock: func() *mockTagService {
				return &mockTagService{}
			},
			expectedError: true,
			errorContains: "project_id",
		},
		{
			name: "error - missing region",
			resourceData: map[string]interface{}{
				tfconstants.AttrName:      "test-tag",
				tfconstants.AttrProjectID: "test-project",
			},
			setupMock: func() *mockTagService {
				return &mockTagService{}
			},
			expectedError: true,
			errorContains: "region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "", "") // No defaults
			resource := ResourceTag()
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
// Test: resourceTagRead
// ============================================================================

func TestResourceTagRead_Success(t *testing.T) {
	mockService := &mockTagService{
		getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
			assert.Equal(t, "12345", tagID)
			return &goe2e.Tag{
				LabelID:   12345,
				LabelName: "test-tag",
				Metadata:  "Test metadata",
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceTag()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-tag",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "12345", d.Id())
	assert.Equal(t, "test-tag", d.Get(tfconstants.AttrName))
	assert.Equal(t, "Test metadata", d.Get(attrMetadata))
	assert.Equal(t, 12345, d.Get(attrLabelID))
}

func TestResourceTagRead_NotFound(t *testing.T) {
	mockService := &mockTagService{
		getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
			return nil, nil, fmt.Errorf("tag with ID %s not found", tagID)
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceTag()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-tag",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("99999")

	diags := resource.ReadContext(context.Background(), d, cfg)

	// Should handle not found gracefully
	require.False(t, diags.HasError(), "Read should handle not found gracefully")
	assert.Empty(t, d.Id(), "ID should be cleared when not found")
}

func TestResourceTagRead_NotFoundError(t *testing.T) {
	mockService := &mockTagService{
		getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
			err := fmt.Errorf("error: %s", goe2econstants.NotFoundSubstring)
			return nil, nil, err
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceTag()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-tag",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("99999")

	diags := resource.ReadContext(context.Background(), d, cfg)

	// Should handle not found error gracefully
	require.False(t, diags.HasError(), "Read should handle not found error gracefully")
	assert.Empty(t, d.Id(), "ID should be cleared when not found")
}

func TestResourceTagRead_APIError(t *testing.T) {
	mockService := &mockTagService{
		getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to read tag")
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceTag()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-tag",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	assert.Contains(t, diags[0].Summary, "Error reading tag")
}

// ============================================================================
// Test: resourceTagUpdate
// ============================================================================

func TestResourceTagUpdate_NotSupported(t *testing.T) {
	cfg := createMockConfig(t, &mockTagService{}, "test-project", "us-east-1")
	resource := ResourceTag()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-tag",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.UpdateContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Update should fail (not supported)")
	assert.Contains(t, diags[0].Summary, "not supported")
}

// ============================================================================
// Test: resourceTagDelete
// ============================================================================

func TestResourceTagDelete_Success(t *testing.T) {
	mockService := &mockTagService{
		deleteTagFunc: func(ctx context.Context, tagID string) (*goe2e.Response, error) {
			assert.Equal(t, "12345", tagID)
			return &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceTag()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-tag",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Delete should succeed")
	assert.Empty(t, d.Id(), "ID should be cleared after delete")
}

func TestResourceTagDelete_NotFound(t *testing.T) {
	mockService := &mockTagService{
		deleteTagFunc: func(ctx context.Context, tagID string) (*goe2e.Response, error) {
			err := fmt.Errorf("error: %s", goe2econstants.NotFoundSubstring)
			return nil, err
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceTag()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-tag",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("99999")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	// Should handle not found gracefully (idempotent delete)
	require.False(t, diags.HasError(), "Delete should handle not found gracefully")
	assert.Empty(t, d.Id(), "ID should be cleared")
}

func TestResourceTagDelete_APIError(t *testing.T) {
	mockService := &mockTagService{
		deleteTagFunc: func(ctx context.Context, tagID string) (*goe2e.Response, error) {
			return nil, errors.New("API error: failed to delete tag")
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := ResourceTag()
	d := createTestResourceData(t, map[string]interface{}{
		tfconstants.AttrName:      "test-tag",
		tfconstants.AttrProjectID: "test-project",
		tfconstants.AttrRegion:    "us-east-1",
	})
	d.SetId("12345")

	diags := resource.DeleteContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Delete should fail on API error")
	assert.Contains(t, diags[0].Summary, "Error deleting tag")
}

// ============================================================================
// Test: resourceTagImport
// ============================================================================

func TestResourceTagImport_Success(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		setupMock     func() *mockTagService
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name:     "3-part format: project_id/region/tag_id",
			importID: "test-project/us-east-1/12345",
			setupMock: func() *mockTagService {
				return &mockTagService{
					getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
						return &goe2e.Tag{
							LabelID:   12345,
							LabelName: "imported-tag",
							Metadata:  "Imported metadata",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "12345", d.Id())
				assert.Equal(t, "test-project", d.Get(tfconstants.AttrProjectID))
				assert.Equal(t, "us-east-1", d.Get(tfconstants.AttrRegion))
				assert.Equal(t, "imported-tag", d.Get(tfconstants.AttrName))
			},
		},
		{
			name:     "1-part format: tag_id (uses provider defaults)",
			importID: "67890",
			setupMock: func() *mockTagService {
				return &mockTagService{
					getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
						return &goe2e.Tag{
							LabelID:   67890,
							LabelName: "imported-tag-2",
							Metadata:  "",
						}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
					},
				}
			},
			validateState: func(t *testing.T, d *schema.ResourceData) {
				assert.Equal(t, "67890", d.Id())
				assert.Equal(t, "test-project", d.Get(tfconstants.AttrProjectID))
				assert.Equal(t, "us-east-1", d.Get(tfconstants.AttrRegion))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
			resource := ResourceTag()
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

func TestResourceTagImport_Errors(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		setupMock     func() *mockTagService
		expectedError bool
		errorContains string
	}{
		{
			name:          "invalid format - too many parts",
			importID:      "project/region/tag/extra",
			setupMock:     func() *mockTagService { return &mockTagService{} },
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "invalid format - 2 parts",
			importID:      "project/tag",
			setupMock:     func() *mockTagService { return &mockTagService{} },
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:     "1-part format without provider defaults",
			importID: "12345",
			setupMock: func() *mockTagService {
				return &mockTagService{}
			},
			expectedError: true,
			errorContains: "project_id is required",
		},
		{
			name:     "API error during fetch",
			importID: "test-project/us-east-1/12345",
			setupMock: func() *mockTagService {
				return &mockTagService{
					getTagFunc: func(ctx context.Context, tagID string) (*goe2e.Tag, *goe2e.Response, error) {
						return nil, nil, errors.New("API error: failed to fetch tag")
					},
				}
			},
			expectedError: true,
			errorContains: "error importing tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "", "") // No defaults to test requirements
			resource := ResourceTag()
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

// ============================================================================
// Test: isNotFoundError helper
// ============================================================================

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "not found error with substring",
			err:      fmt.Errorf("error: %s", goe2econstants.NotFoundSubstring),
			expected: true,
		},
		{
			name:     "not found error with code",
			err:      fmt.Errorf("error: %s", goe2econstants.NotFoundCode),
			expected: true,
		},
		{
			name:     "generic error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "error containing not found substring",
			err:      fmt.Errorf("tag with ID %s not found", "123"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNotFoundError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
