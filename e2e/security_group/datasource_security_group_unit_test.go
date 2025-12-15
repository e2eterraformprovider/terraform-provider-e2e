package security_group

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
// Mock implementations for Security Group datasource tests
// ============================================================================

// mockSecurityGroupServiceForDatasource is a mock implementation of SecurityGroupService for datasource testing
type mockSecurityGroupServiceForDatasource struct {
	getSecurityGroupListFunc func(ctx context.Context) ([]*goe2e.SecurityGroup, *goe2e.Response, error)
}

func (m *mockSecurityGroupServiceForDatasource) GetSecurityGroupList(ctx context.Context) ([]*goe2e.SecurityGroup, *goe2e.Response, error) {
	if m.getSecurityGroupListFunc != nil {
		return m.getSecurityGroupListFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}

// Unused interface methods
func (m *mockSecurityGroupServiceForDatasource) GetSecurityGroup(ctx context.Context, sgID string) (*goe2e.SecurityGroup, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockSecurityGroupServiceForDatasource) CreateSecurityGroup(ctx context.Context, req *goe2e.SecurityGroupCreateRequest) (*goe2e.SecurityGroup, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockSecurityGroupServiceForDatasource) UpdateSecurityGroup(ctx context.Context, sgID string, req *goe2e.SecurityGroupUpdateRequest) (*goe2e.SecurityGroup, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockSecurityGroupServiceForDatasource) DeleteSecurityGroup(ctx context.Context, sgID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockSecurityGroupServiceForDatasource) MakeDefaultSecurityGroup(ctx context.Context, sgID string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockSecurityGroupServiceForDatasource) AttachSecurityGroup(ctx context.Context, nodeID int, req *goe2e.SecurityGroupAttachRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockSecurityGroupServiceForDatasource) DetachSecurityGroup(ctx context.Context, nodeID int, req *goe2e.SecurityGroupAttachRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

// createMockConfigForSecurityGroupDatasource creates a config with a mock security group service
func createMockConfigForSecurityGroupDatasource(t *testing.T, mockService *mockSecurityGroupServiceForDatasource) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	mockClient := &goe2e.Client{}
	mockClient.SecurityGroups = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// ============================================================================
// Test: dataSourceSecurityGroupRead
// ============================================================================

func TestDataSourceSecurityGroupRead_Success(t *testing.T) {
	mockService := &mockSecurityGroupServiceForDatasource{
		getSecurityGroupListFunc: func(ctx context.Context) ([]*goe2e.SecurityGroup, *goe2e.Response, error) {
			return []*goe2e.SecurityGroup{
				{
					ID:          "sg-123",
					Name:        "test-sg",
					Description: "Test security group",
					IsDefault:   false,
					Rules: []goe2e.Rule{
						{
							ID:           1,
							RuleType:     "ingress",
							ProtocolName: "tcp",
							PortRange:    "80-80",
							Network:      "0.0.0.0/0",
							NetworkCIDR:  "0.0.0.0/0",
							Description:  "Allow HTTP",
						},
					},
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForSecurityGroupDatasource(t, mockService)
	resource := DataSourceSecurityGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrName: "test-sg",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "sg-123", d.Id())
	assert.Equal(t, "test-sg", d.Get(tfconstants.AttrName))
	assert.Equal(t, "Test security group", d.Get("description"))
	assert.False(t, d.Get("default").(bool))

	rules := d.Get("rules").([]interface{})
	require.Len(t, rules, 1)
	rule := rules[0].(map[string]interface{})
	assert.Equal(t, 1, rule["rule_id"])
	assert.Equal(t, "ingress", rule["rule_type"])
	assert.Equal(t, "tcp", rule["protocol_name"])
}

func TestDataSourceSecurityGroupRead_NotFound(t *testing.T) {
	mockService := &mockSecurityGroupServiceForDatasource{
		getSecurityGroupListFunc: func(ctx context.Context) ([]*goe2e.SecurityGroup, *goe2e.Response, error) {
			return []*goe2e.SecurityGroup{
				{
					ID:   "sg-123",
					Name: "other-sg",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForSecurityGroupDatasource(t, mockService)
	resource := DataSourceSecurityGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrName: "non-existent-sg",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when security group not found")
	assert.Contains(t, diags[0].Summary, "not found")
}

func TestDataSourceSecurityGroupRead_EmptyList(t *testing.T) {
	mockService := &mockSecurityGroupServiceForDatasource{
		getSecurityGroupListFunc: func(ctx context.Context) ([]*goe2e.SecurityGroup, *goe2e.Response, error) {
			return []*goe2e.SecurityGroup{}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForSecurityGroupDatasource(t, mockService)
	resource := DataSourceSecurityGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrName: "test-sg",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when security group not found")
}

func TestDataSourceSecurityGroupRead_APIError(t *testing.T) {
	mockService := &mockSecurityGroupServiceForDatasource{
		getSecurityGroupListFunc: func(ctx context.Context) ([]*goe2e.SecurityGroup, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to list security groups")
		},
	}

	cfg := createMockConfigForSecurityGroupDatasource(t, mockService)
	resource := DataSourceSecurityGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrName: "test-sg",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	assert.Contains(t, diags[0].Summary, "Error listing security groups")
}

func TestDataSourceSecurityGroupRead_DefaultSG(t *testing.T) {
	mockService := &mockSecurityGroupServiceForDatasource{
		getSecurityGroupListFunc: func(ctx context.Context) ([]*goe2e.SecurityGroup, *goe2e.Response, error) {
			return []*goe2e.SecurityGroup{
				{
					ID:          "sg-default",
					Name:        "default",
					Description: "Default security group",
					IsDefault:   true,
					Rules:       []goe2e.Rule{},
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForSecurityGroupDatasource(t, mockService)
	resource := DataSourceSecurityGroup()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrName: "default",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.True(t, d.Get("default").(bool))
}

// ============================================================================
// Test: flattenRules
// ============================================================================

func TestFlattenRules_ForDatasource(t *testing.T) {
	// Note: flattenRules is already tested in resource_security_group_unit_test.go
	// This test is just to ensure it works correctly for datasource use cases
	networkSize := 1
	rules := []goe2e.Rule{
		{
			ID:           1,
			RuleType:     "ingress",
			ProtocolName: "tcp",
			PortRange:    "80-80",
			Network:      "0.0.0.0/0",
			NetworkCIDR:  "0.0.0.0/0",
			NetworkSize:  &networkSize,
			Description:  "Allow HTTP",
		},
	}

	result := flattenRules(rules)

	require.Len(t, result, 1)
	ruleMap := result[0]
	assert.Equal(t, 1, ruleMap["rule_id"])
	assert.Equal(t, "ingress", ruleMap["rule_type"])
	assert.Equal(t, "tcp", ruleMap["protocol_name"])
	assert.Equal(t, "80-80", ruleMap["port_range"])
	assert.Equal(t, "0.0.0.0/0", ruleMap["network"])
	assert.Equal(t, "0.0.0.0/0", ruleMap["network_cidr"])
	assert.Equal(t, 1, ruleMap["size"])
	assert.Equal(t, "Allow HTTP", ruleMap["description"])
}
