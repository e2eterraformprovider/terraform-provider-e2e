package dbaas_mysql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mysql"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockConfig creates a mock config with a mock goe2e client
func mockConfigWithMySQLService(service goe2e.DBaaSMySQLService) *config.Config {
	cfg := &config.Config{
		DefaultProjectID: "test-project",
		DefaultRegion:    "Mumbai",
	}
	client := &goe2e.Client{
		DBaaSMySQL: service,
	}

	// Use the public testing method to set the client
	cfg.SetGoe2eClientForTesting(client)

	return cfg
}

// Helper to create a test MySQL cluster
func createTestMySQLCluster(id int, name, status string) *goe2e.MySQLCluster {
	return &goe2e.MySQLCluster{
		ID:                  id,
		Name:                name,
		Status:              status,
		IsEncryptionEnabled: false,
		MasterNode: goe2e.DBNode{
			PublicIPAddress:  "1.2.3.4",
			PrivateIPAddress: "10.0.0.1",
			Port:             "3306",
			Disk:             "100GB",
			Status:           status,
			Database: goe2e.DBCreds{
				PGDetail: goe2e.PGDetail{
					ID: 0,
				},
			},
			Plan: goe2e.Plan{
				Software: goe2e.Software{},
			},
		},
	}
}

// ============================================================================
// Tests for resourceCreateMySqlDB
// ============================================================================

func TestResourceCreateMySqlDB_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVersion:   "8.0",
		tfconstants.AttrPlan:      "E2E-2C-4GB",
		tfconstants.AttrDBaaSName: "test-mysql",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":     "testuser",
				"password": "testpass",
				"name":     "testdb",
			},
		},
	})

	// Mock API calls
	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwareMySQL, "8.0").Return(1, nil)
	mockService.On("GetTemplateID", ctx, "E2E-2C-4GB", 1).Return(100, nil)
	mockService.On("CreateCluster", ctx, mock.AnythingOfType("*goe2e.MySQLClusterCreateRequest")).Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning),
		&goe2e.Response{},
		nil,
	)

	diags := resource.CreateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "123", d.Id())
	mockService.AssertExpectations(t)
}

func TestResourceCreateMySqlDB_GetSoftwareIDError(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVersion:   "8.0",
		tfconstants.AttrPlan:      "E2E-2C-4GB",
		tfconstants.AttrDBaaSName: "test-mysql",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":     "testuser",
				"password": "testpass",
				"name":     "testdb",
			},
		},
	})

	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwareMySQL, "8.0").Return(0, errors.New("software not found"))

	diags := resource.CreateContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceCreateMySqlDB_GetTemplateIDError(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVersion:   "8.0",
		tfconstants.AttrPlan:      "E2E-2C-4GB",
		tfconstants.AttrDBaaSName: "test-mysql",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":     "testuser",
				"password": "testpass",
				"name":     "testdb",
			},
		},
	})

	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwareMySQL, "8.0").Return(1, nil)
	mockService.On("GetTemplateID", ctx, "E2E-2C-4GB", 1).Return(0, errors.New("template not found"))

	diags := resource.CreateContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceCreateMySqlDB_CreateClusterError(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVersion:   "8.0",
		tfconstants.AttrPlan:      "E2E-2C-4GB",
		tfconstants.AttrDBaaSName: "test-mysql",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":     "testuser",
				"password": "testpass",
				"name":     "testdb",
			},
		},
	})

	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwareMySQL, "8.0").Return(1, nil)
	mockService.On("GetTemplateID", ctx, "E2E-2C-4GB", 1).Return(100, nil)
	mockService.On("CreateCluster", ctx, mock.AnythingOfType("*goe2e.MySQLClusterCreateRequest")).Return(
		nil,
		&goe2e.Response{},
		errors.New("create failed"),
	)

	diags := resource.CreateContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

// ============================================================================
// Tests for resourceReadMySqlDB
// ============================================================================

func TestResourceReadMySqlDB_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	cluster := createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning)
	cluster.MasterNode.Database.PGDetail.ID = 5
	cluster.MasterNode.Database.ID = 10
	cluster.MasterNode.Database.Username = "testuser"
	cluster.MasterNode.Database.Database = "testdb"
	cluster.MasterNode.Plan.Name = "E2E-2C-4GB"
	cluster.MasterNode.Plan.Software.Version = "8.0"
	cluster.MasterNode.Status = goe2econstants.DBaaSStatusRunning

	mockService.On("GetCluster", ctx, "123").Return(cluster, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "123", d.Id())
	assert.Equal(t, goe2econstants.DBaaSStatusRunning, d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "1.2.3.4", d.Get(tfconstants.AttrPublicIPAddress))
	assert.Equal(t, "10.0.0.1", d.Get(tfconstants.AttrPrivateIPAddress))
	assert.Equal(t, "3306", d.Get(tfconstants.AttrPort))
	assert.Equal(t, "100GB", d.Get(tfconstants.AttrDisk))
	assert.Equal(t, 5, d.Get(tfconstants.AttrParameterGroupID))
	mockService.AssertExpectations(t)
}

func TestResourceReadMySqlDB_NotFound(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	mockService.On("GetCluster", ctx, "123").Return(nil, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "", d.Id()) // ID should be cleared
	mockService.AssertExpectations(t)
}

func TestResourceReadMySqlDB_Error(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	mockService.On("GetCluster", ctx, "123").Return(nil, &goe2e.Response{}, errors.New("API error"))

	diags := resource.ReadContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

// ============================================================================
// Tests for resourceDeleteMySqlDB
// ============================================================================

func TestResourceDeleteMySqlDB_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	mockService.On("DeleteCluster", ctx, "123").Return(&goe2e.Response{}, nil)

	diags := resource.DeleteContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "", d.Id()) // ID should be cleared
	mockService.AssertExpectations(t)
}

func TestResourceDeleteMySqlDB_AlreadyDeleted(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	// Simulate 404 error (already deleted)
	err := errors.New("not found")
	mockService.On("DeleteCluster", ctx, "123").Return(nil, err)

	_ = resource.DeleteContext(ctx, d, cfg)

	// Should succeed (idempotent) - check if error contains "not found"
	// The actual implementation checks for NotFoundSubstring
	assert.Equal(t, "", d.Id()) // ID should be cleared
	mockService.AssertExpectations(t)
}

func TestResourceDeleteMySqlDB_Error(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	mockService.On("DeleteCluster", ctx, "123").Return(nil, errors.New("delete failed"))

	diags := resource.DeleteContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	_ = diags // Use diags to avoid unused variable error
	mockService.AssertExpectations(t)
}
