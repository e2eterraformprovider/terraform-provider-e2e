package reserve_ip

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
// Mock implementations for Reserve IP datasource tests
// ============================================================================

// mockReserveIPServiceForDatasource is a mock implementation of ReserveIPService for datasource testing
type mockReserveIPServiceForDatasource struct {
	listReserveIPsFunc func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error)
}

func (m *mockReserveIPServiceForDatasource) ListReserveIPs(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
	if m.listReserveIPsFunc != nil {
		return m.listReserveIPsFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}

// Unused interface methods
func (m *mockReserveIPServiceForDatasource) GetReserveIP(ctx context.Context, ipAddress string) (*goe2e.ReserveIP, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForDatasource) CreateReserveIP(ctx context.Context) (*goe2e.ReserveIP, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForDatasource) DeleteReserveIP(ctx context.Context, ipAddress string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForDatasource) AttachFloatingIP(ctx context.Context, req *goe2e.FloatingIPAttachmentRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForDatasource) DetachFloatingIP(ctx context.Context, req *goe2e.FloatingIPDetachmentRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

// createMockConfigForReserveIPDatasource creates a config with a mock reserve IP service
func createMockConfigForReserveIPDatasource(t *testing.T, mockService *mockReserveIPServiceForDatasource, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.ReserveIP = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// ============================================================================
// Test: dataSourceReadReserveIps
// ============================================================================

func TestDataSourceReadReserveIps_Success(t *testing.T) {
	mockService := &mockReserveIPServiceForDatasource{
		listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
			return []goe2e.ReserveIP{
				{
					ReserveID:     "123",
					ApplianceType: "Node",
					IPAddress:     "203.0.113.1",
					ReservedType:  "IPv4",
					VMID:          12345,
					VMName:        "test-vm-1",
					Status:        "attached",
					BoughtAt:      "2024-01-01T00:00:00Z",
				},
				{
					ReserveID:     "456",
					ApplianceType: "Node",
					IPAddress:     "203.0.113.2",
					ReservedType:  "IPv4",
					VMID:          0, // Not attached
					VMName:        "",
					Status:        "available",
					BoughtAt:      "2024-01-02T00:00:00Z",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForReserveIPDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceReserveIps()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "reserve_ips_list", d.Id())

	reserveIPsList := d.Get("reserve_ips_list").([]interface{})
	require.Len(t, reserveIPsList, 2, "Should have 2 reserve IPs")

	// Validate first IP
	ip1 := reserveIPsList[0].(map[string]interface{})
	assert.Equal(t, "123", ip1[tfconstants.AttrReserveID])
	assert.Equal(t, "Node", ip1[tfconstants.AttrApplianceType])
	assert.Equal(t, "203.0.113.1", ip1[tfconstants.AttrIPAddress])
	assert.Equal(t, "attached", ip1[tfconstants.AttrStatus])
	assert.Equal(t, 12345.0, ip1[tfconstants.AttrVMID])
	assert.Equal(t, "test-vm-1", ip1[tfconstants.AttrVMName])

	// Validate second IP
	ip2 := reserveIPsList[1].(map[string]interface{})
	assert.Equal(t, "available", ip2[tfconstants.AttrStatus])
	assert.Equal(t, 0.0, ip2[tfconstants.AttrVMID])
}

func TestDataSourceReadReserveIps_EmptyList(t *testing.T) {
	mockService := &mockReserveIPServiceForDatasource{
		listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
			return []goe2e.ReserveIP{}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForReserveIPDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceReserveIps()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed with empty list")
	assert.Equal(t, "reserve_ips_list", d.Id())

	reserveIPsList := d.Get("reserve_ips_list").([]interface{})
	assert.Len(t, reserveIPsList, 0, "Should have empty list")
}

func TestDataSourceReadReserveIps_APIError(t *testing.T) {
	mockService := &mockReserveIPServiceForDatasource{
		listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to list reserve IPs")
		},
	}

	cfg := createMockConfigForReserveIPDatasource(t, mockService, "test-project", "us-east-1")
	resource := DataSourceReserveIps()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	assert.Contains(t, diags[0].Summary, "error listing reserved IPs")
}

// ============================================================================
// Test: flattenReserveIps
// ============================================================================

func TestFlattenReserveIps(t *testing.T) {
	tests := []struct {
		name           string
		reserveIPs     []goe2e.ReserveIP
		expectedLength int
		validateFunc   func(*testing.T, []interface{})
	}{
		{
			name:           "nil input - returns empty slice",
			reserveIPs:     nil,
			expectedLength: 0,
			validateFunc: func(t *testing.T, result []interface{}) {
				assert.Len(t, result, 0)
			},
		},
		{
			name:           "empty slice - returns empty slice",
			reserveIPs:     []goe2e.ReserveIP{},
			expectedLength: 0,
			validateFunc: func(t *testing.T, result []interface{}) {
				assert.Len(t, result, 0)
			},
		},
		{
			name: "single reserve IP - all fields present",
			reserveIPs: []goe2e.ReserveIP{
				{
					ReserveID:     "123",
					ApplianceType: "Node",
					IPAddress:     "203.0.113.1",
					ReservedType:  "IPv4",
					VMID:          12345,
					VMName:        "test-vm",
					Status:        "attached",
					BoughtAt:      "2024-01-01T00:00:00Z",
				},
			},
			expectedLength: 1,
			validateFunc: func(t *testing.T, result []interface{}) {
				require.Len(t, result, 1)
				ipMap := result[0].(map[string]interface{})
				assert.Equal(t, "123", ipMap[tfconstants.AttrReserveID])
				assert.Equal(t, "Node", ipMap[tfconstants.AttrApplianceType])
				assert.Equal(t, "203.0.113.1", ipMap[tfconstants.AttrIPAddress])
				assert.Equal(t, "IPv4", ipMap[tfconstants.AttrReservedType])
				assert.Equal(t, 12345, ipMap[tfconstants.AttrVMID]) // VMID from goe2e is int64
				assert.Equal(t, "test-vm", ipMap[tfconstants.AttrVMName])
				assert.Equal(t, "attached", ipMap[tfconstants.AttrStatus])
				// Note: bought_at is stored in the map but schema uses created_at
				assert.Equal(t, "2024-01-01T00:00:00Z", ipMap[tfconstants.AttrBoughtAt])
			},
		},
		{
			name: "multiple reserve IPs",
			reserveIPs: []goe2e.ReserveIP{
				{
					ReserveID: "123",
					IPAddress: "203.0.113.1",
					Status:    "attached",
					BoughtAt:  "2024-01-01T00:00:00Z",
				},
				{
					ReserveID: "456",
					IPAddress: "203.0.113.2",
					Status:    "available",
					BoughtAt:  "2024-01-02T00:00:00Z",
				},
			},
			expectedLength: 2,
			validateFunc: func(t *testing.T, result []interface{}) {
				require.Len(t, result, 2)
				ip1 := result[0].(map[string]interface{})
				assert.Equal(t, "203.0.113.1", ip1[tfconstants.AttrIPAddress])
				ip2 := result[1].(map[string]interface{})
				assert.Equal(t, "203.0.113.2", ip2[tfconstants.AttrIPAddress])
			},
		},
		{
			name: "reserve IP with empty fields",
			reserveIPs: []goe2e.ReserveIP{
				{
					ReserveID:     "789",
					ApplianceType: "",
					IPAddress:     "203.0.113.3",
					ReservedType:  "",
					VMID:          0,
					VMName:        "",
					Status:        "",
					BoughtAt:      "",
				},
			},
			expectedLength: 1,
			validateFunc: func(t *testing.T, result []interface{}) {
				require.Len(t, result, 1)
				ipMap := result[0].(map[string]interface{})
				assert.Equal(t, "789", ipMap[tfconstants.AttrReserveID])
				assert.Equal(t, "203.0.113.3", ipMap[tfconstants.AttrIPAddress])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenReserveIps(tt.reserveIPs)

			assert.Len(t, result, tt.expectedLength)
			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}
