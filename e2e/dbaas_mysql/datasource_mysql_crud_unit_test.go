package dbaas_mysql_test

import (
	"context"
	"errors"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mysql"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Tests for dataSourceReadMySQL
// ============================================================================

func TestDataSourceReadMySQL_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	cluster := createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning)
	cluster.MasterNode.Database.PGDetail.ID = 5
	cluster.MasterNode.Database.ID = 10
	cluster.MasterNode.Database.Database = "testdb"
	cluster.MasterNode.Database.Username = "testuser"
	cluster.MasterNode.Plan.Name = "E2E-2C-4GB"
	cluster.MasterNode.Plan.Software.Version = "8.0"
	cluster.MasterNode.Status = goe2econstants.DBaaSStatusRunning

	mockService.On("GetCluster", ctx, "123").Return(cluster, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "123", d.Id())
	assert.Equal(t, 10, d.Get(tfconstants.AttrDatabaseID))
	assert.Equal(t, "testdb", d.Get(tfconstants.AttrDatabaseName))
	assert.Equal(t, "testuser", d.Get(tfconstants.AttrDatabaseUser))
	assert.Equal(t, goe2econstants.DBaaSStatusRunning, d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "1.2.3.4", d.Get(tfconstants.AttrPublicIPAddress))
	assert.Equal(t, "10.0.0.1", d.Get(tfconstants.AttrPrivateIPAddress))
	assert.Equal(t, true, d.Get(tfconstants.AttrIsPublicIPAttached))
	assert.Equal(t, "100GB", d.Get(tfconstants.AttrDisk))
	assert.Equal(t, "E2E-2C-4GB", d.Get(tfconstants.AttrPlan))
	assert.Equal(t, "8.0", d.Get(tfconstants.AttrDatabaseVersion))
	assert.Equal(t, 5, d.Get(tfconstants.AttrParameterGroupID))
	assert.Equal(t, goe2econstants.DBaaSStatusRunning, d.Get(tfconstants.AttrPowerStatus))
	mockService.AssertExpectations(t)
}

func TestDataSourceReadMySQL_NotFound(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	mockService.On("GetCluster", ctx, "123").Return(nil, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestDataSourceReadMySQL_Error(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	mockService.On("GetCluster", ctx, "123").Return(nil, &goe2e.Response{}, errors.New("API error"))

	diags := resource.ReadContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestDataSourceReadMySQL_NoParameterGroup(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	cluster := createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning)
	cluster.MasterNode.Database.PGDetail.ID = 0 // No parameter group
	cluster.MasterNode.Database.ID = 10
	cluster.MasterNode.Database.Database = "testdb"
	cluster.MasterNode.Database.Username = "testuser"
	cluster.MasterNode.Plan.Name = "E2E-2C-4GB"
	cluster.MasterNode.Plan.Software.Version = "8.0"
	cluster.MasterNode.Status = goe2econstants.DBaaSStatusRunning

	mockService.On("GetCluster", ctx, "123").Return(cluster, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, tfconstants.DBaaSDefaultParameterGroupID, d.Get(tfconstants.AttrParameterGroupID))
	mockService.AssertExpectations(t)
}

func TestDataSourceReadMySQL_NoPublicIP(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.DataSourceMySQLDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	cluster := createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning)
	cluster.MasterNode.PublicIPAddress = "" // No public IP
	cluster.MasterNode.Database.ID = 10
	cluster.MasterNode.Database.Database = "testdb"
	cluster.MasterNode.Database.Username = "testuser"
	cluster.MasterNode.Plan.Name = "E2E-2C-4GB"
	cluster.MasterNode.Plan.Software.Version = "8.0"
	cluster.MasterNode.Status = goe2econstants.DBaaSStatusRunning

	mockService.On("GetCluster", ctx, "123").Return(cluster, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, false, d.Get(tfconstants.AttrIsPublicIPAttached))
	mockService.AssertExpectations(t)
}
