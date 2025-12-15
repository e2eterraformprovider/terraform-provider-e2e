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
	"github.com/stretchr/testify/mock"
)

// ============================================================================
// Tests for resourceUpdateMySqlDB
// ============================================================================

func TestResourceUpdateMySqlDB_StatusChange_Stop(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	// Create resource data with new status (what user wants)
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrStatus:           goe2econstants.DBaaSStatusSuspended, // New status from config
		tfconstants.AttrPublicIPRequired: false,                               // Keep other mutable fields same
	})
	d.SetId("123")

	mockService.On("StopCluster", ctx, "123").Return(&goe2e.Response{}, nil)
	mockService.On("GetCluster", ctx, "123").Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusSuspended),
		&goe2e.Response{},
		nil,
	)

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceUpdateMySqlDB_StatusChange_Start(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	// Create resource data with new status (what user wants)
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrStatus:           goe2econstants.DBaaSStatusRunning, // Target status
		tfconstants.AttrPublicIPRequired: false,                             // Keep other mutable fields same
	})
	d.SetId("123")

	mockService.On("StartCluster", ctx, "123").Return(&goe2e.Response{}, nil)
	mockService.On("GetCluster", ctx, "123").Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning),
		&goe2e.Response{},
		nil,
	)

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceUpdateMySqlDB_StatusChange_Restart(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	// Create resource data with new status (what user wants)
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrStatus:           goe2econstants.DBaaSStatusRestarting, // Target status for restart
		tfconstants.AttrPublicIPRequired: false,                                // Keep other mutable fields same
	})
	d.SetId("123")

	mockService.On("RestartCluster", ctx, "123").Return(&goe2e.Response{}, nil)
	mockService.On("GetCluster", ctx, "123").Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning),
		&goe2e.Response{},
		nil,
	)

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceUpdateMySqlDB_PlanUpgrade_WithoutSuspendedState(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	// Test tries to upgrade plan while in RUNNING state (should fail)
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrPlan:             "E2E-2C-4GB", // New plan (upgraded)
		tfconstants.AttrStatus:           goe2econstants.DBaaSStatusRunning,
		tfconstants.AttrVersion:          "8.0",
		tfconstants.AttrPublicIPRequired: false,
	})
	d.SetId("123")

	// Mock may be called due to status change detection in unit test framework
	mockService.On("StartCluster", ctx, "123").Return(&goe2e.Response{}, nil).Maybe()

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary, "cannot upgrade plan")
}

func TestResourceUpdateMySqlDB_PlanUpgrade_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	// Test successfully upgrades plan while in SUSPENDED state
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrPlan:             "E2E-2C-4GB", // New plan (upgraded from E2E-4C-8GB)
		tfconstants.AttrStatus:           goe2econstants.DBaaSStatusSuspended,
		tfconstants.AttrVersion:          "8.0",
		tfconstants.AttrPublicIPRequired: false,
	})
	d.SetId("123")

	// Mock may be called due to status change detection in unit test framework
	mockService.On("StopCluster", ctx, "123").Return(&goe2e.Response{}, nil).Maybe()
	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwareMySQL, "8.0").Return(1, nil)
	mockService.On("GetTemplateID", ctx, "E2E-2C-4GB", 1).Return(200, nil)
	mockService.On("UpgradePlan", ctx, "123", mock.AnythingOfType("*goe2e.MySQLPlanUpgradeRequest")).Return(&goe2e.Response{}, nil)
	mockService.On("GetCluster", ctx, "123").Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusSuspended),
		&goe2e.Response{},
		nil,
	)

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceUpdateMySqlDB_DiskExpansion(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	// Initialize with new state values (what user wants)
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrSize:             10, // New size value
		tfconstants.AttrPublicIPRequired: false,
	})
	d.SetId("123")

	mockService.On("ExpandDisk", ctx, "123", mock.MatchedBy(func(req *goe2e.DiskExpansionRequest) bool {
		return req.Size == 10
	})).Return(&goe2e.Response{}, nil)
	mockService.On("GetCluster", ctx, "123").Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning),
		&goe2e.Response{},
		nil,
	)

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceUpdateMySqlDB_VPCAttach(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	// Create resource data with new VPC list
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVPCs:             []interface{}{1},
		tfconstants.AttrPublicIPRequired: false,
	})
	d.SetId("123")

	expectedVPCs := []goe2e.VPCMetadata{
		{NetworkID: "1", VPCName: "vpc1", IPv4CIDR: "10.0.0.0/24"},
	}
	mockService.On("ExpandVPCList", ctx, []string{"1"}).Return(expectedVPCs, nil)
	mockService.On("AttachVPC", ctx, "123", mock.MatchedBy(func(req *goe2e.MySQLVPCAttachRequest) bool {
		return len(req.VPCs) == 1 && req.VPCs[0].NetworkID == "1"
	})).Return(&goe2e.Response{}, nil)
	mockService.On("GetCluster", ctx, "123").Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning),
		&goe2e.Response{},
		nil,
	)

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceUpdateMySqlDB_VPCAttachError_Rollback(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVPCs:             []interface{}{1},
		tfconstants.AttrPublicIPRequired: false,
	})
	d.SetId("123")

	expectedVPCs := []goe2e.VPCMetadata{
		{NetworkID: "1", VPCName: "vpc1", IPv4CIDR: "10.0.0.0/24"},
	}
	mockService.On("ExpandVPCList", ctx, []string{"1"}).Return(expectedVPCs, nil)
	mockService.On("AttachVPC", ctx, "123", mock.Anything).Return(nil, errors.New("attach failed"))

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceUpdateMySqlDB_PublicIPAttach(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	// New desired state is true (attach public IP)
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrPublicIPRequired: true,
	})
	d.SetId("123")

	mockService.On("AttachPublicIP", ctx, "123").Return(&goe2e.Response{}, nil)
	mockService.On("GetCluster", ctx, "123").Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning),
		&goe2e.Response{},
		nil,
	)

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceUpdateMySqlDB_ParameterGroupAttach(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_mysql.MockDBaaSMySQLService{}
	cfg := mockConfigWithMySQLService(mockService)

	resource := dbaas_mysql.ResourceMySql()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrParameterGroupID: 5,
		tfconstants.AttrPublicIPRequired: false,
	})
	d.SetId("123")

	mockService.On("AttachParameterGroup", ctx, "123", "5").Return(&goe2e.Response{}, nil)
	mockService.On("GetCluster", ctx, "123").Return(
		createTestMySQLCluster(123, "test-mysql", goe2econstants.DBaaSStatusRunning),
		&goe2e.Response{},
		nil,
	)

	diags := resource.UpdateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	mockService.AssertExpectations(t)
}
