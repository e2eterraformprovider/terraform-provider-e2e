package dbaas_mysql

import (
	"context"
	"errors"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNormalizeStatus(t *testing.T) {
	tests := []struct {
		name       string
		apiStatus  string
		wantStatus string
	}{
		{
			name:       "returns status unchanged - STOPPED",
			apiStatus:  goe2econstants.DBaaSStatusStopped,
			wantStatus: goe2econstants.DBaaSStatusStopped,
		},
		{
			name:       "returns status unchanged - SUSPENDED",
			apiStatus:  goe2econstants.DBaaSStatusSuspended,
			wantStatus: goe2econstants.DBaaSStatusSuspended,
		},
		{
			name:       "returns status unchanged - RUNNING",
			apiStatus:  goe2econstants.DBaaSStatusRunning,
			wantStatus: goe2econstants.DBaaSStatusRunning,
		},
		{
			name:       "returns status unchanged - RESTARTING",
			apiStatus:  goe2econstants.DBaaSStatusRestarting,
			wantStatus: goe2econstants.DBaaSStatusRestarting,
		},
		{
			name:       "returns empty string unchanged",
			apiStatus:  "",
			wantStatus: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeStatus(tt.apiStatus)
			assert.Equal(t, tt.wantStatus, result, "normalizeStatus should be a passthrough for MySQL")
		})
	}
}

func TestExpandVPCList(t *testing.T) {
	ctx := context.Background()

	t.Run("with empty slice returns empty VPC list", func(t *testing.T) {
		mockService := &MockDBaaSMySQLService{}
		result, err := expandVPCList(ctx, &goe2e.Client{DBaaSMySQL: mockService}, []interface{}{})
		assert.NoError(t, err)
		assert.Empty(t, result)
		// Empty slice doesn't call ExpandVPCList, so no expectations to verify
	})

	t.Run("with single VPC ID", func(t *testing.T) {
		mockService := &MockDBaaSMySQLService{}
		expectedVPCs := []goe2e.VPCMetadata{
			{NetworkID: "1", VPCName: "vpc1", IPv4CIDR: "10.0.0.0/24"},
		}
		mockService.On("ExpandVPCList", ctx, []string{"1"}).Return(expectedVPCs, nil)

		result, err := expandVPCList(ctx, &goe2e.Client{DBaaSMySQL: mockService}, []interface{}{1})
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expectedVPCs, result)
		mockService.AssertExpectations(t)
	})

	t.Run("with multiple VPC IDs", func(t *testing.T) {
		mockService := &MockDBaaSMySQLService{}
		expectedVPCs := []goe2e.VPCMetadata{
			{NetworkID: "1", VPCName: "vpc1", IPv4CIDR: "10.0.0.0/24"},
			{NetworkID: "2", VPCName: "vpc2", IPv4CIDR: "10.0.1.0/24"},
			{NetworkID: "3", VPCName: "vpc3", IPv4CIDR: "10.0.2.0/24"},
		}
		mockService.On("ExpandVPCList", ctx, []string{"1", "2", "3"}).Return(expectedVPCs, nil)

		result, err := expandVPCList(ctx, &goe2e.Client{DBaaSMySQL: mockService}, []interface{}{1, 2, 3})
		require.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, expectedVPCs, result)
		mockService.AssertExpectations(t)
	})

	t.Run("handles goe2e client errors", func(t *testing.T) {
		mockService := &MockDBaaSMySQLService{}
		expectedErr := errors.New("VPC expansion failed")
		mockService.On("ExpandVPCList", ctx, []string{"1"}).Return(nil, expectedErr)

		result, err := expandVPCList(ctx, &goe2e.Client{DBaaSMySQL: mockService}, []interface{}{1})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		mockService.AssertExpectations(t)
	})
}

func TestBuildMySQLCreateRequest(t *testing.T) {
	ctx := context.Background()
	softwareID := 1
	templateID := 100

	t.Run("with minimal required fields", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			tfconstants.AttrDatabase: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user":         {Type: schema.TypeString},
						"password":     {Type: schema.TypeString},
						"name":         {Type: schema.TypeString},
						"dbaas_number": {Type: schema.TypeInt},
					},
				},
			},
			tfconstants.AttrDBaaSName:        {Type: schema.TypeString},
			tfconstants.AttrGroup:            {Type: schema.TypeString},
			tfconstants.AttrPublicIPRequired: {Type: schema.TypeBool},
		}, map[string]interface{}{
			tfconstants.AttrDatabase: []interface{}{
				map[string]interface{}{
					"user":         "testuser",
					"password":     "testpass",
					"name":         "testdb",
					"dbaas_number": 1,
				},
			},
			tfconstants.AttrDBaaSName:        "test-mysql",
			tfconstants.AttrGroup:            tfconstants.DBaaSDefaultGroupName,
			tfconstants.AttrPublicIPRequired: tfconstants.DBaaSDefaultPublicIPRequired,
		})

		mockService := &MockDBaaSMySQLService{}
		req, err := buildMySQLCreateRequest(ctx, d, &goe2e.Client{DBaaSMySQL: mockService}, softwareID, templateID)
		require.NoError(t, err)
		assert.Equal(t, "test-mysql", req.Name)
		assert.Equal(t, softwareID, req.SoftwareID)
		assert.Equal(t, templateID, req.TemplateID)
		assert.Equal(t, tfconstants.DBaaSDefaultGroupName, req.Group)
		assert.Equal(t, tfconstants.DBaaSDefaultPublicIPRequired, req.PublicIPRequired)
		assert.Equal(t, "testuser", req.Database.User)
		assert.Equal(t, "testpass", req.Database.Password)
		assert.Equal(t, "testdb", req.Database.Name)
		assert.Equal(t, 1, req.Database.DBaaSNumber)
		assert.Empty(t, req.Vpcs)
		assert.Zero(t, req.ParameterGroupId)
	})

	t.Run("with all optional fields", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			tfconstants.AttrDatabase: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user":         {Type: schema.TypeString},
						"password":     {Type: schema.TypeString},
						"name":         {Type: schema.TypeString},
						"dbaas_number": {Type: schema.TypeInt},
					},
				},
			},
			tfconstants.AttrDBaaSName:        {Type: schema.TypeString},
			tfconstants.AttrGroup:            {Type: schema.TypeString},
			tfconstants.AttrPublicIPRequired: {Type: schema.TypeBool},
			tfconstants.AttrVPCs:             {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeInt}},
			tfconstants.AttrParameterGroupID: {Type: schema.TypeInt},
		}, map[string]interface{}{
			tfconstants.AttrDatabase: []interface{}{
				map[string]interface{}{
					"user":         "testuser",
					"password":     "testpass",
					"name":         "testdb",
					"dbaas_number": 1,
				},
			},
			tfconstants.AttrDBaaSName:        "test-mysql",
			tfconstants.AttrGroup:            "CustomGroup",
			tfconstants.AttrPublicIPRequired: false,
			tfconstants.AttrVPCs:             []interface{}{1, 2}, // Use slice instead of Set for TestResourceDataRaw
			tfconstants.AttrParameterGroupID: 5,
		})

		// Manually set VPCs as a Set after creating the resource data
		vpcSet := schema.NewSet(schema.HashInt, []interface{}{1, 2})
		if err := d.Set(tfconstants.AttrVPCs, vpcSet); err != nil {
			t.Fatalf("Failed to set VPCs: %v", err)
		}

		mockService := &MockDBaaSMySQLService{}
		expectedVPCs := []goe2e.VPCMetadata{
			{NetworkID: "1", VPCName: "vpc1", IPv4CIDR: "10.0.0.0/24"},
			{NetworkID: "2", VPCName: "vpc2", IPv4CIDR: "10.0.1.0/24"},
		}
		// Use mock.MatchedBy to handle unordered VPC IDs from Set
		mockService.On("ExpandVPCList", ctx, mock.MatchedBy(func(vpcIDs []string) bool {
			if len(vpcIDs) != 2 {
				return false
			}
			return (vpcIDs[0] == "1" && vpcIDs[1] == "2") || (vpcIDs[0] == "2" && vpcIDs[1] == "1")
		})).Return(expectedVPCs, nil)

		req, err := buildMySQLCreateRequest(ctx, d, &goe2e.Client{DBaaSMySQL: mockService}, softwareID, templateID)
		require.NoError(t, err)
		assert.Equal(t, "test-mysql", req.Name)
		assert.Equal(t, "CustomGroup", req.Group)
		assert.False(t, req.PublicIPRequired)
		assert.Equal(t, expectedVPCs, req.Vpcs)
		assert.Equal(t, 5, req.ParameterGroupId)
		mockService.AssertExpectations(t)
	})

	t.Run("error handling for missing database block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			tfconstants.AttrDatabase: {Type: schema.TypeList, Required: true},
		}, map[string]interface{}{
			tfconstants.AttrDatabase: []interface{}{},
		})

		mockService := &MockDBaaSMySQLService{}
		req, err := buildMySQLCreateRequest(ctx, d, &goe2e.Client{DBaaSMySQL: mockService}, softwareID, templateID)
		assert.Error(t, err)
		assert.Nil(t, req)
		assert.Contains(t, err.Error(), tfconstants.DatabaseConfigurationRequired)
	})

	t.Run("error handling for VPC expansion errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			tfconstants.AttrDatabase: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user":         {Type: schema.TypeString},
						"password":     {Type: schema.TypeString},
						"name":         {Type: schema.TypeString},
						"dbaas_number": {Type: schema.TypeInt},
					},
				},
			},
			tfconstants.AttrDBaaSName:        {Type: schema.TypeString},
			tfconstants.AttrGroup:            {Type: schema.TypeString},
			tfconstants.AttrPublicIPRequired: {Type: schema.TypeBool},
			tfconstants.AttrVPCs:             {Type: schema.TypeSet, Elem: &schema.Schema{Type: schema.TypeInt}},
		}, map[string]interface{}{
			tfconstants.AttrDatabase: []interface{}{
				map[string]interface{}{
					"user":         "testuser",
					"password":     "testpass",
					"name":         "testdb",
					"dbaas_number": 1,
				},
			},
			tfconstants.AttrDBaaSName:        "test-mysql",
			tfconstants.AttrGroup:            tfconstants.DBaaSDefaultGroupName,
			tfconstants.AttrPublicIPRequired: true,
			tfconstants.AttrVPCs:             []interface{}{1}, // Use slice instead of Set
		})

		// Manually set VPCs as a Set after creating the resource data
		vpcSet := schema.NewSet(schema.HashInt, []interface{}{1})
		if err := d.Set(tfconstants.AttrVPCs, vpcSet); err != nil {
			t.Fatalf("Failed to set VPCs: %v", err)
		}

		mockService := &MockDBaaSMySQLService{}
		expectedErr := errors.New("VPC expansion failed")
		mockService.On("ExpandVPCList", ctx, mock.MatchedBy(func(vpcIDs []string) bool {
			return len(vpcIDs) == 1 && vpcIDs[0] == "1"
		})).Return(nil, expectedErr)

		req, err := buildMySQLCreateRequest(ctx, d, &goe2e.Client{DBaaSMySQL: mockService}, softwareID, templateID)
		assert.Error(t, err)
		assert.Nil(t, req)
		assert.Contains(t, err.Error(), "failed to expand VPC list")
		mockService.AssertExpectations(t)
	})
}

func TestCustomImportStateFunc(t *testing.T) {
	t.Run("valid import format", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			tfconstants.AttrProjectID: {Type: schema.TypeString},
		}, map[string]interface{}{})
		d.SetId("project123" + tfconstants.DBaaSImportIDSeparator + "dbaas456")

		result, err := customImportStateFunc(d, nil)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "dbaas456", result[0].Id())
		projectID, ok := result[0].Get(tfconstants.AttrProjectID).(string)
		require.True(t, ok)
		assert.Equal(t, "project123", projectID)
	})

	t.Run("invalid format - missing colon", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("project123dbaas456")

		result, err := customImportStateFunc(d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), tfconstants.DBaaSImportIDFormatDescription)
	})

	t.Run("invalid format - too many parts", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("project123" + tfconstants.DBaaSImportIDSeparator + "dbaas456" + tfconstants.DBaaSImportIDSeparator + "extra")

		result, err := customImportStateFunc(d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), tfconstants.DBaaSImportIDFormatDescription)
	})

	t.Run("invalid format - empty parts", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId(tfconstants.DBaaSImportIDSeparator)

		result, err := customImportStateFunc(d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("invalid format - empty project_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId(tfconstants.DBaaSImportIDSeparator + "dbaas456")

		result, err := customImportStateFunc(d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("invalid format - empty dbaas_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("project123" + tfconstants.DBaaSImportIDSeparator)

		result, err := customImportStateFunc(d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
