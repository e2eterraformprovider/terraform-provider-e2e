package objectstore

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
// Mock implementations for Object Store datasource tests
// ============================================================================

// mockObjectStorageServiceForDatasource is a mock implementation of ObjectStorageService for datasource testing
type mockObjectStorageServiceForDatasource struct {
	listBucketsFunc func(ctx context.Context) ([]goe2e.Bucket, *goe2e.Response, error)
}

func (m *mockObjectStorageServiceForDatasource) ListBuckets(ctx context.Context) ([]goe2e.Bucket, *goe2e.Response, error) {
	if m.listBucketsFunc != nil {
		return m.listBucketsFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}

// Unused interface methods
func (m *mockObjectStorageServiceForDatasource) CreateBucket(ctx context.Context, req *goe2e.BucketCreateRequest) (*goe2e.Bucket, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockObjectStorageServiceForDatasource) GetBucket(ctx context.Context, bucketName string) (*goe2e.Bucket, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockObjectStorageServiceForDatasource) DeleteBucket(ctx context.Context, bucketName string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockObjectStorageServiceForDatasource) SetBucketVersioning(ctx context.Context, bucketName string, req *goe2e.BucketVersioningRequest) (*goe2e.BucketVersioning, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

// createMockConfigForObjectStoreDatasource creates a config with a mock object storage service
func createMockConfigForObjectStoreDatasource(t *testing.T, mockService *mockObjectStorageServiceForDatasource, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.ObjectStorage = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// ============================================================================
// Test: dataSourceReadBuckets
// ============================================================================

func TestDataSourceReadBuckets_Success(t *testing.T) {
	mockService := &mockObjectStorageServiceForDatasource{
		listBucketsFunc: func(ctx context.Context) ([]goe2e.Bucket, *goe2e.Response, error) {
			return []goe2e.Bucket{
				{
					ID:                           123.0,
					Name:                         "test-bucket-1",
					BucketSize:                   "100GB",
					Status:                       "Active",
					CreatedAt:                    "2024-01-01T00:00:00Z",
					VersioningStatus:             "Enabled",
					LifecycleConfigurationStatus: "Enabled",
				},
				{
					ID:                           456.0,
					Name:                         "test-bucket-2",
					BucketSize:                   "50GB",
					Status:                       "Active",
					CreatedAt:                    "2024-01-02T00:00:00Z",
					VersioningStatus:             "Disabled",
					LifecycleConfigurationStatus: "Disabled",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForObjectStoreDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceObjectStores()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "bucket_list", d.Id())

	bucketList := d.Get("bucket_list").([]interface{})
	require.Len(t, bucketList, 2, "Should have 2 buckets")

	// Validate first bucket
	bucket1 := bucketList[0].(map[string]interface{})
	assert.Equal(t, 123.0, bucket1[tfconstants.AttrID])
	assert.Equal(t, "test-bucket-1", bucket1[tfconstants.AttrName])
	assert.Equal(t, "100GB", bucket1["bucket_size"])
	assert.Equal(t, "Active", bucket1[tfconstants.AttrStatus])
}

func TestDataSourceReadBuckets_EmptyList(t *testing.T) {
	mockService := &mockObjectStorageServiceForDatasource{
		listBucketsFunc: func(ctx context.Context) ([]goe2e.Bucket, *goe2e.Response, error) {
			return []goe2e.Bucket{}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForObjectStoreDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceObjectStores()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed with empty list")
	assert.Equal(t, "bucket_list", d.Id())

	bucketList := d.Get("bucket_list").([]interface{})
	assert.Len(t, bucketList, 0, "Should have empty list")
}

func TestDataSourceReadBuckets_APIError(t *testing.T) {
	mockService := &mockObjectStorageServiceForDatasource{
		listBucketsFunc: func(ctx context.Context) ([]goe2e.Bucket, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to list buckets")
		},
	}

	cfg := createMockConfigForObjectStoreDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceObjectStores()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
}

// ============================================================================
// Test: flattenBuckets
// ============================================================================

func TestFlattenBuckets(t *testing.T) {
	tests := []struct {
		name           string
		buckets        interface{}
		expectedLength int
		validateFunc   func(*testing.T, []interface{})
	}{
		{
			name:           "nil input - returns empty slice",
			buckets:        nil,
			expectedLength: 0,
			validateFunc: func(t *testing.T, result []interface{}) {
				assert.Len(t, result, 0)
			},
		},
		{
			name:           "empty slice - returns empty slice",
			buckets:        []interface{}{},
			expectedLength: 0,
			validateFunc: func(t *testing.T, result []interface{}) {
				assert.Len(t, result, 0)
			},
		},
		{
			name: "single bucket - all fields present",
			buckets: []interface{}{
				map[string]interface{}{
					"id":                             123.0,
					"name":                           "test-bucket",
					"bucket_size":                    "100GB",
					"status":                         "Active",
					"created_at":                     "2024-01-01T00:00:00Z",
					"versioning_status":              "Enabled",
					"lifecycle_configuration_status": "Enabled",
				},
			},
			expectedLength: 1,
			validateFunc: func(t *testing.T, result []interface{}) {
				require.Len(t, result, 1)
				bucketMap := result[0].(map[string]interface{})
				assert.Equal(t, 123.0, bucketMap[tfconstants.AttrID])
				assert.Equal(t, "test-bucket", bucketMap[tfconstants.AttrName])
				assert.Equal(t, "100GB", bucketMap["bucket_size"])
				assert.Equal(t, "Active", bucketMap[tfconstants.AttrStatus])
				assert.Equal(t, "2024-01-01T00:00:00Z", bucketMap[tfconstants.AttrCreatedAt])
				assert.Equal(t, "Enabled", bucketMap[tfconstants.AttrVersioningStatus])
				assert.Equal(t, "Enabled", bucketMap[tfconstants.AttrLifecycleConfigurationStatus])
			},
		},
		{
			name: "multiple buckets",
			buckets: []interface{}{
				map[string]interface{}{
					"id":   123.0,
					"name": "bucket-1",
				},
				map[string]interface{}{
					"id":   456.0,
					"name": "bucket-2",
				},
			},
			expectedLength: 2,
			validateFunc: func(t *testing.T, result []interface{}) {
				require.Len(t, result, 2)
				bucket1 := result[0].(map[string]interface{})
				assert.Equal(t, "bucket-1", bucket1[tfconstants.AttrName])
				bucket2 := result[1].(map[string]interface{})
				assert.Equal(t, "bucket-2", bucket2[tfconstants.AttrName])
			},
		},
		{
			name: "invalid type - returns empty slice",
			buckets: map[string]interface{}{
				"invalid": "data",
			},
			expectedLength: 0,
			validateFunc: func(t *testing.T, result []interface{}) {
				assert.Len(t, result, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenBuckets(tt.buckets)

			assert.Len(t, result, tt.expectedLength)
			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}
