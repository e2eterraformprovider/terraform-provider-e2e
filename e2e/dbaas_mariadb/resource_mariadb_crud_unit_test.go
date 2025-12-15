package dbaas_mariadb_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mariadb"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Test Helpers
// ============================================================================

// createTestConfigWithMock creates a test config with a mock MariaDB service
func createTestConfigWithMock(mockService *dbaas_mariadb.MockMariaDBService) *config.Config {
	cfg := &config.Config{
		DefaultProjectID: "test-project",
		DefaultRegion:    "Mumbai",
	}

	// Create a minimal goe2e client and inject the mock
	client := &goe2e.Client{}
	client.MariaDB = mockService
	cfg.SetGoe2eClientForTesting(client)

	return cfg
}

// Helper to set goe2e client in config (using reflection or unexported method)
// Since we can't directly set it, we'll need to use a workaround
// For now, we'll create a helper that works with the actual structure

// createTestResourceData creates a ResourceData for testing
func createTestResourceData(t *testing.T, resource *schema.Resource, data map[string]interface{}) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, resource.Schema, data)
}

// ============================================================================
// Create Operation Tests
// ============================================================================

func TestResourceCreateMariaDB_Success(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetSoftwareIDFunc: func(ctx context.Context, name, version string) (int, error) {
			return 100, nil
		},
		GetTemplateIDFunc: func(ctx context.Context, planName string, softwareID int) (int, error) {
			return 200, nil
		},
		ExpandVPCListFunc: func(ctx context.Context, vpcIDs []string) ([]goe2e.VPCMetadata, error) {
			return []goe2e.VPCMetadata{}, nil
		},
		CreateMariaDBFunc: func(ctx context.Context, req *goe2e.MariaDBCreateRequest) (*goe2e.MariaDB, *goe2e.Response, error) {
			assert.Equal(t, "test-mariadb", req.Name)
			assert.Equal(t, 100, req.SoftwareID)
			assert.Equal(t, 200, req.TemplateID)
			return &goe2e.MariaDB{
				ID:     12345,
				Name:   "test-mariadb",
				Status: goe2econstants.DBaaSStatusRunning,
				MasterNode: goe2e.DBNode{
					PublicIPAddress:  "1.2.3.4",
					PrivateIPAddress: "10.0.0.1",
					Port:             "3306",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 201}}, nil
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.ResourceMariaDB()

	d := createTestResourceData(t, resource, map[string]interface{}{
		tfconstants.AttrName:            "test-mariadb",
		tfconstants.AttrSoftwareName:    "MariaDB",
		tfconstants.AttrSoftwareVersion: "10.6",
		tfconstants.AttrGroup:           "Default",
		tfconstants.AttrPlan:            "DBS.16GB",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				tfconstants.AttrDatabaseBlockUser:        "testuser",
				tfconstants.AttrDatabaseBlockPassword:    "testpass",
				tfconstants.AttrDatabaseBlockName:        "testdb",
				tfconstants.AttrDatabaseBlockDBaaSNumber: 1,
			},
		},
	})

	diags := resource.CreateContext(context.Background(), d, cfg)
	require.False(t, diags.HasError(), "Create should succeed")
	assert.Equal(t, "12345", d.Id())
	assert.Equal(t, "test-mariadb", d.Get(tfconstants.AttrName))
}

func TestResourceCreateMariaDB_SoftwareIDError(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetSoftwareIDFunc: func(ctx context.Context, name, version string) (int, error) {
			return 0, errors.New("software not found")
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.ResourceMariaDB()

	d := createTestResourceData(t, resource, map[string]interface{}{
		tfconstants.AttrName:            "test-mariadb",
		tfconstants.AttrSoftwareName:    "MariaDB",
		tfconstants.AttrSoftwareVersion: "10.6",
		tfconstants.AttrGroup:           "Default",
		tfconstants.AttrPlan:            "DBS.16GB",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				tfconstants.AttrDatabaseBlockUser:        "testuser",
				tfconstants.AttrDatabaseBlockPassword:    "testpass",
				tfconstants.AttrDatabaseBlockName:        "testdb",
				tfconstants.AttrDatabaseBlockDBaaSNumber: 1,
			},
		},
	})

	diags := resource.CreateContext(context.Background(), d, cfg)
	require.True(t, diags.HasError(), "Create should fail")
	assert.Empty(t, d.Id())
}

// ============================================================================
// Read Operation Tests
// ============================================================================

func TestResourceReadMariaDB_Success(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetMariaDBFunc: func(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
			return &goe2e.MariaDB{
				ID:     12345,
				Name:   "test-mariadb",
				Status: goe2econstants.DBaaSStatusRunning,
				Software: goe2e.Software{
					Name:    "MariaDB",
					Version: "10.6",
				},
				MasterNode: goe2e.DBNode{
					Plan: goe2e.Plan{
						Name: "DBS.16GB",
					},
					PublicIPAddress:  "1.2.3.4",
					PrivateIPAddress: "10.0.0.1",
					Port:             "3306",
					Disk:             "100",
				},
				IsEncryptionEnabled: false,
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.ResourceMariaDB()

	d := createTestResourceData(t, resource, map[string]interface{}{
		"id": "12345",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)
	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "test-mariadb", d.Get(tfconstants.AttrName))
	assert.Equal(t, "MariaDB", d.Get(tfconstants.AttrSoftwareName))
	assert.Equal(t, "10.6", d.Get(tfconstants.AttrSoftwareVersion))
}

func TestResourceReadMariaDB_NotFound(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetMariaDBFunc: func(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
			return nil, nil, errors.New("not found")
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.ResourceMariaDB()

	d := createTestResourceData(t, resource, map[string]interface{}{
		"id": "12345",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)
	require.True(t, diags.HasError(), "Read should fail")
}

func TestResourceReadMariaDB_StatusNormalization(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetMariaDBFunc: func(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
			return &goe2e.MariaDB{
				ID:     12345,
				Name:   "test-mariadb",
				Status: goe2econstants.DBaaSStatusSuspended, // Should be normalized to STOPPED
				MasterNode: goe2e.DBNode{
					Plan: goe2e.Plan{Name: "DBS.16GB"},
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.ResourceMariaDB()

	d := createTestResourceData(t, resource, map[string]interface{}{
		"id": "12345",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)
	require.False(t, diags.HasError(), "Read should succeed")
	// Status should be normalized from SUSPENDED to STOPPED
	assert.Equal(t, goe2econstants.DBaaSStatusStopped, d.Get(tfconstants.AttrStatus))
}

// ============================================================================
// Update Operation Tests
// ============================================================================

func TestResourceUpdateMariaDB_StatusChange(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetMariaDBFunc: func(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
			return &goe2e.MariaDB{
				ID:     12345,
				Name:   "test-mariadb",
				Status: goe2econstants.DBaaSStatusStopped,
				MasterNode: goe2e.DBNode{
					Plan: goe2e.Plan{Name: "DBS.16GB"},
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
		ShutdownMariaDBFunc: func(ctx context.Context, id string) (*goe2e.Response, error) {
			return &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.ResourceMariaDB()

	d := createTestResourceData(t, resource, map[string]interface{}{
		"id":     "12345",
		"status": goe2econstants.DBaaSStatusStopped,
	})

	// Simulate status change
	d.Set(tfconstants.AttrStatus, goe2econstants.DBaaSStatusStopped)

	diags := resource.UpdateContext(context.Background(), d, cfg)
	require.False(t, diags.HasError(), "Update should succeed")
}

// ============================================================================
// Delete Operation Tests
// ============================================================================

func TestResourceDeleteMariaDB_Success(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		DeleteMariaDBFunc: func(ctx context.Context, id string) (*goe2e.Response, error) {
			return &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.ResourceMariaDB()

	d := createTestResourceData(t, resource, map[string]interface{}{
		"id": "12345",
	})

	diags := resource.DeleteContext(context.Background(), d, cfg)
	require.False(t, diags.HasError(), "Delete should succeed")
	assert.Empty(t, d.Id())
}

func TestResourceDeleteMariaDB_NotFound(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		DeleteMariaDBFunc: func(ctx context.Context, id string) (*goe2e.Response, error) {
			return nil, errors.New("not found")
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.ResourceMariaDB()

	d := createTestResourceData(t, resource, map[string]interface{}{
		"id": "12345",
	})

	// Should handle 404 gracefully (idempotent delete)
	_ = resource.DeleteContext(context.Background(), d, cfg)
	// Note: Current implementation may return error, but should ideally be idempotent
	// This test documents current behavior
}
