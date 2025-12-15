package blockstorage

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
// Mock implementations for Block Storage datasource tests
// ============================================================================

// mockBlockStorageServiceForDatasource is a mock implementation of BlockStorageService for datasource testing
type mockBlockStorageServiceForDatasource struct {
	getBlockStorageFunc func(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error)
}

func (m *mockBlockStorageServiceForDatasource) GetBlockStorage(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error) {
	if m.getBlockStorageFunc != nil {
		return m.getBlockStorageFunc(ctx, blockID)
	}
	return nil, nil, errors.New("not implemented")
}

// Unused interface methods
func (m *mockBlockStorageServiceForDatasource) CreateBlockStorage(ctx context.Context, req *goe2e.BlockStorageCreateRequest) (*goe2e.BlockStorage, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockBlockStorageServiceForDatasource) DeleteBlockStorage(ctx context.Context, blockID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockBlockStorageServiceForDatasource) UpgradeBlockStorage(ctx context.Context, blockID string, req *goe2e.BlockStorageUpgradeRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockBlockStorageServiceForDatasource) AttachBlockStorage(ctx context.Context, blockID string, req *goe2e.BlockStorageAttachRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockBlockStorageServiceForDatasource) DetachBlockStorage(ctx context.Context, blockID string, req *goe2e.BlockStorageAttachRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockBlockStorageServiceForDatasource) GetBlockStoragePlans(ctx context.Context) ([]goe2e.BlockStoragePlan, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

// createMockConfigForBlockStorageDatasource creates a config with a mock block storage service
func createMockConfigForBlockStorageDatasource(t *testing.T, mockService *mockBlockStorageServiceForDatasource, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.BlockStorage = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// ============================================================================
// Test: dataSourceReadBlockStorage
// ============================================================================

func TestDataSourceReadBlockStorage_Success(t *testing.T) {
	mockService := &mockBlockStorageServiceForDatasource{
		getBlockStorageFunc: func(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error) {
			assert.Equal(t, "block-123", blockID)
			return &goe2e.BlockStorage{
				BlockID: 123,
				Name:    "test-block-storage",
				Size:    100,
				Status:  "Active",
				Template: goe2e.ResponseTemplate{
					TotalIOPSSec: "3000",
				},
				VMDetail: map[string]interface{}{
					tfconstants.VolumeAttachmentVMDetailKeyVMID:   float64(12345),
					tfconstants.VolumeAttachmentVMDetailKeyVMName: "test-vm",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForBlockStorageDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "block-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "block-123", d.Id())
	assert.Equal(t, "test-block-storage", d.Get(tfconstants.AttrName))
	assert.Equal(t, 100.0, d.Get(tfconstants.AttrSize))
	assert.Equal(t, "Active", d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "3000", d.Get(tfconstants.AttrIOPS))
	assert.Equal(t, "12345", d.Get(tfconstants.AttrVMID))
	assert.Equal(t, "test-vm", d.Get(tfconstants.AttrVMName))
}

func TestDataSourceReadBlockStorage_Detached(t *testing.T) {
	mockService := &mockBlockStorageServiceForDatasource{
		getBlockStorageFunc: func(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error) {
			return &goe2e.BlockStorage{
				BlockID: 456,
				Name:    "test-block-detached",
				Size:    50,
				Status:  "Available",
				Template: goe2e.ResponseTemplate{
					TotalIOPSSec: "1500",
				},
				VMDetail: nil, // Not attached
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForBlockStorageDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "block-456",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	// When VMDetail is empty or nil, the fields should not be set (default to empty string in schema)
	assert.Equal(t, "", d.Get(tfconstants.AttrVMID))
	assert.Equal(t, "", d.Get(tfconstants.AttrVMName))
}

func TestDataSourceReadBlockStorage_NotFound(t *testing.T) {
	mockService := &mockBlockStorageServiceForDatasource{
		getBlockStorageFunc: func(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error) {
			return nil, &goe2e.Response{Response: &http.Response{StatusCode: http.StatusNotFound}}, fmt.Errorf("block storage with ID %s not found", blockID)
		},
	}

	cfg := createMockConfigForBlockStorageDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "non-existent-block",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when block storage not found")
	assert.Contains(t, diags[0].Summary, "not found")
}

func TestDataSourceReadBlockStorage_NotFoundError(t *testing.T) {
	mockService := &mockBlockStorageServiceForDatasource{
		getBlockStorageFunc: func(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error) {
			err := fmt.Errorf("error: %s", goe2econstants.NotFoundSubstring)
			return nil, nil, err
		},
	}

	cfg := createMockConfigForBlockStorageDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "block-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when block storage not found")
	assert.Contains(t, diags[0].Summary, "not found")
}

func TestDataSourceReadBlockStorage_NilBlockStorage(t *testing.T) {
	mockService := &mockBlockStorageServiceForDatasource{
		getBlockStorageFunc: func(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error) {
			return nil, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForBlockStorageDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "block-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when block storage is nil")
	assert.Contains(t, diags[0].Summary, "not found")
}

func TestDataSourceReadBlockStorage_APIError(t *testing.T) {
	mockService := &mockBlockStorageServiceForDatasource{
		getBlockStorageFunc: func(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to retrieve block storage")
		},
	}

	cfg := createMockConfigForBlockStorageDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "block-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	assert.Contains(t, diags[0].Summary, "Error retrieving block storage")
}

func TestDataSourceReadBlockStorage_VMDetailStringVMID(t *testing.T) {
	mockService := &mockBlockStorageServiceForDatasource{
		getBlockStorageFunc: func(ctx context.Context, blockID string) (*goe2e.BlockStorage, *goe2e.Response, error) {
			return &goe2e.BlockStorage{
				BlockID: 789,
				Name:    "test-block",
				Size:    200,
				Status:  "Active",
				Template: goe2e.ResponseTemplate{
					TotalIOPSSec: "6000",
				},
				VMDetail: map[string]interface{}{
					tfconstants.VolumeAttachmentVMDetailKeyVMID:   "67890", // String format
					tfconstants.VolumeAttachmentVMDetailKeyVMName: "test-vm-2",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForBlockStorageDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "block-789",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "67890", d.Get(tfconstants.AttrVMID))
}

func TestDataSourceReadBlockStorage_MissingRegion(t *testing.T) {
	cfg := createMockConfigForBlockStorageDatasource(t, &mockBlockStorageServiceForDatasource{}, "", "")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "block-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail without region")
	assert.Contains(t, diags[0].Summary, "region")
}

func TestDataSourceReadBlockStorage_MissingProjectID(t *testing.T) {
	cfg := createMockConfigForBlockStorageDatasource(t, &mockBlockStorageServiceForDatasource{}, "", "us-east-1")
	resource := DataSourceBlockStorage()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrBlockID: "block-123",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail without project_id")
	assert.Contains(t, diags[0].Summary, "project_id")
}
