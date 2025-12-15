package vpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlattenVpcs(t *testing.T) {
	tests := []struct {
		name           string
		vpcs           []goe2e.Vpc
		expectedLength int
		validateFunc   func([]interface{}) error
	}{
		{
			name:           "nil input - returns empty slice",
			vpcs:           nil,
			expectedLength: 0,
			validateFunc: func(result []interface{}) error {
				if len(result) != 0 {
					return fmt.Errorf("expected empty slice for nil input, got length %d", len(result))
				}
				return nil
			},
		},
		{
			name:           "empty slice - returns empty slice",
			vpcs:           []goe2e.Vpc{},
			expectedLength: 0,
			validateFunc: func(result []interface{}) error {
				if len(result) != 0 {
					return fmt.Errorf("expected empty slice, got length %d", len(result))
				}
				return nil
			},
		},
		{
			name: "single VPC - all fields present",
			vpcs: []goe2e.Vpc{
				{
					ID:        123.0,
					Name:      "test-vpc-1",
					State:     goe2econstants.VPCStatusActive,
					CreatedAt: "2024-01-01T00:00:00Z",
					IPv4CIDR:  "10.0.0.0/24",
					GatewayIP: "10.0.0.1",
					IsActive:  true,
					PoolSize:  254.0,
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				if len(result) != 1 {
					return fmt.Errorf("expected 1 VPC, got %d", len(result))
				}

				vpcMap := result[0].(map[string]interface{})

				if id, ok := vpcMap[tfconstants.AttrNetworkID].(float64); !ok || id != 123.0 {
					return fmt.Errorf("NetworkID = %v (type %T), want 123.0 (float64)", vpcMap[tfconstants.AttrNetworkID], vpcMap[tfconstants.AttrNetworkID])
				}

				if name, ok := vpcMap[tfconstants.AttrName].(string); !ok || name != "test-vpc-1" {
					return fmt.Errorf("Name = %v, want test-vpc-1", name)
				}

				if state, ok := vpcMap[tfconstants.AttrStatus].(string); !ok || state != goe2econstants.VPCStatusActive {
					return fmt.Errorf("Status = %v, want %s", state, goe2econstants.VPCStatusActive)
				}

				if createdAt, ok := vpcMap[tfconstants.AttrCreatedAt].(string); !ok || createdAt != "2024-01-01T00:00:00Z" {
					return fmt.Errorf("CreatedAt = %v, want 2024-01-01T00:00:00Z", createdAt)
				}

				if ipv4CIDR, ok := vpcMap[tfconstants.AttrIPv4CIDR].(string); !ok || ipv4CIDR != "10.0.0.0/24" {
					return fmt.Errorf("IPv4CIDR = %v, want 10.0.0.0/24", ipv4CIDR)
				}

				if gatewayIP, ok := vpcMap[tfconstants.AttrGatewayIP].(string); !ok || gatewayIP != "10.0.0.1" {
					return fmt.Errorf("GatewayIP = %v, want 10.0.0.1", gatewayIP)
				}

				if isActive, ok := vpcMap[tfconstants.AttrIsActive].(bool); !ok || isActive != true {
					return fmt.Errorf("IsActive = %v, want true", isActive)
				}

				if poolSize, ok := vpcMap[tfconstants.AttrPoolSize].(float64); !ok || poolSize != 254.0 {
					return fmt.Errorf("PoolSize = %v, want 254.0", poolSize)
				}

				return nil
			},
		},
		{
			name: "multiple VPCs",
			vpcs: []goe2e.Vpc{
				{
					ID:        123.0,
					Name:      "test-vpc-1",
					State:     goe2econstants.VPCStatusActive,
					CreatedAt: "2024-01-01T00:00:00Z",
					IPv4CIDR:  "10.0.0.0/24",
					GatewayIP: "10.0.0.1",
					IsActive:  true,
					PoolSize:  254.0,
				},
				{
					ID:        456.0,
					Name:      "test-vpc-2",
					State:     goe2econstants.VPCStatusActive,
					CreatedAt: "2024-01-02T00:00:00Z",
					IPv4CIDR:  "10.0.1.0/24",
					GatewayIP: "10.0.1.1",
					IsActive:  false,
					PoolSize:  128.0,
				},
			},
			expectedLength: 2,
			validateFunc: func(result []interface{}) error {
				if len(result) != 2 {
					return fmt.Errorf("expected 2 VPCs, got %d", len(result))
				}

				// Validate first VPC
				vpc1 := result[0].(map[string]interface{})
				if name, _ := vpc1[tfconstants.AttrName].(string); name != "test-vpc-1" {
					return fmt.Errorf("first VPC name = %v, want test-vpc-1", name)
				}

				// Validate second VPC
				vpc2 := result[1].(map[string]interface{})
				if name, _ := vpc2[tfconstants.AttrName].(string); name != "test-vpc-2" {
					return fmt.Errorf("second VPC name = %v, want test-vpc-2", name)
				}

				return nil
			},
		},
		{
			name: "VPC with empty fields",
			vpcs: []goe2e.Vpc{
				{
					ID:        789.0,
					Name:      "test-vpc-empty",
					State:     "",
					CreatedAt: "",
					IPv4CIDR:  "",
					GatewayIP: "",
					IsActive:  false,
					PoolSize:  0.0,
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				vpcMap := result[0].(map[string]interface{})

				if id, ok := vpcMap[tfconstants.AttrNetworkID].(float64); !ok || id != 789.0 {
					return fmt.Errorf("NetworkID = %v, want 789.0", id)
				}

				if name, ok := vpcMap[tfconstants.AttrName].(string); !ok || name != "test-vpc-empty" {
					return fmt.Errorf("Name = %v, want test-vpc-empty", name)
				}

				// Empty fields should still be set
				if state, ok := vpcMap[tfconstants.AttrStatus].(string); !ok {
					return fmt.Errorf("Status should be set even if empty")
				} else if state != "" {
					return fmt.Errorf("Status = %v, want empty string", state)
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenVpcs(tt.vpcs)

			if len(result) != tt.expectedLength {
				t.Errorf("flattenVpcs() length = %v, want %v", len(result), tt.expectedLength)
				return
			}

			if tt.validateFunc != nil {
				if err := tt.validateFunc(result); err != nil {
					t.Errorf("flattenVpcs() validation failed: %v", err)
				}
			}
		})
	}
}

// ============================================================================
// Test: dataSourceReadVpcs
// ============================================================================

func TestDataSourceReadVpcs_Success(t *testing.T) {
	mockService := &mockVpcService{
		listVPCsFunc: func(ctx context.Context) ([]goe2e.Vpc, *goe2e.Response, error) {
			return []goe2e.Vpc{
				{
					ID:        123.0,
					Name:      "test-vpc-1",
					State:     goe2econstants.VPCStatusActive,
					CreatedAt: "2024-01-01T00:00:00Z",
					IPv4CIDR:  "10.0.0.0/24",
					GatewayIP: "10.0.0.1",
					IsActive:  true,
					PoolSize:  254.0,
				},
				{
					ID:        456.0,
					Name:      "test-vpc-2",
					State:     goe2econstants.VPCStatusActive,
					CreatedAt: "2024-01-02T00:00:00Z",
					IPv4CIDR:  "10.0.1.0/24",
					GatewayIP: "10.0.1.1",
					IsActive:  false,
					PoolSize:  128.0,
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := DataSourceVpcs()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, tfconstants.AttrVPCList, d.Id())

	vpcList := d.Get(tfconstants.AttrVPCList).([]interface{})
	require.Len(t, vpcList, 2, "Should have 2 VPCs")

	// Validate first VPC
	vpc1 := vpcList[0].(map[string]interface{})
	assert.Equal(t, 123.0, vpc1[tfconstants.AttrNetworkID])
	assert.Equal(t, "test-vpc-1", vpc1[tfconstants.AttrName])
	assert.Equal(t, goe2econstants.VPCStatusActive, vpc1[tfconstants.AttrStatus])
}

func TestDataSourceReadVpcs_EmptyList(t *testing.T) {
	mockService := &mockVpcService{
		listVPCsFunc: func(ctx context.Context) ([]goe2e.Vpc, *goe2e.Response, error) {
			return []goe2e.Vpc{}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := DataSourceVpcs()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	// Empty list is valid - should not error
	require.False(t, diags.HasError(), "Read should succeed with empty list")
}

func TestDataSourceReadVpcs_APIError(t *testing.T) {
	mockService := &mockVpcService{
		listVPCsFunc: func(ctx context.Context) ([]goe2e.Vpc, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to list VPCs")
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := DataSourceVpcs()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
}

// Update mockVpcService to include ListVPCs
func (m *mockVpcService) ListVPCs(ctx context.Context) ([]goe2e.Vpc, *goe2e.Response, error) {
	if m.listVPCsFunc != nil {
		return m.listVPCsFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}
