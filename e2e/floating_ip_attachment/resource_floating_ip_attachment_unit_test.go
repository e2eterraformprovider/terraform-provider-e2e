package floating_ip_attachment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ============================================================================
// Mock implementations for FloatingIPAttachment tests
// ============================================================================

// mockReserveIPServiceForAttachment is a mock implementation for attachment tests
type mockReserveIPServiceForAttachment struct {
	attachFloatingIPFunc func(ctx context.Context, req *goe2e.FloatingIPAttachmentRequest) (*goe2e.Response, error)
	detachFloatingIPFunc func(ctx context.Context, req *goe2e.FloatingIPDetachmentRequest) (*goe2e.Response, error)
	listReserveIPsFunc   func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error)
	getReserveIPFunc     func(ctx context.Context, ipAddress string) (*goe2e.ReserveIP, *goe2e.Response, error)
	createReserveIPFunc  func(ctx context.Context) (*goe2e.ReserveIP, *goe2e.Response, error)
	deleteReserveIPFunc  func(ctx context.Context, ipAddress string) (*goe2e.Response, error)
}

func (m *mockReserveIPServiceForAttachment) ListReserveIPs(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
	if m.listReserveIPsFunc != nil {
		return m.listReserveIPsFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForAttachment) GetReserveIP(ctx context.Context, ipAddress string) (*goe2e.ReserveIP, *goe2e.Response, error) {
	if m.getReserveIPFunc != nil {
		return m.getReserveIPFunc(ctx, ipAddress)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForAttachment) CreateReserveIP(ctx context.Context) (*goe2e.ReserveIP, *goe2e.Response, error) {
	if m.createReserveIPFunc != nil {
		return m.createReserveIPFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForAttachment) DeleteReserveIP(ctx context.Context, ipAddress string) (*goe2e.Response, error) {
	if m.deleteReserveIPFunc != nil {
		return m.deleteReserveIPFunc(ctx, ipAddress)
	}
	return nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForAttachment) AttachFloatingIP(ctx context.Context, req *goe2e.FloatingIPAttachmentRequest) (*goe2e.Response, error) {
	if m.attachFloatingIPFunc != nil {
		return m.attachFloatingIPFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockReserveIPServiceForAttachment) DetachFloatingIP(ctx context.Context, req *goe2e.FloatingIPDetachmentRequest) (*goe2e.Response, error) {
	if m.detachFloatingIPFunc != nil {
		return m.detachFloatingIPFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

// createMockConfig creates a config with a mock client that uses the provided ReserveIP service
// Note: Since Goe2eClientForProject creates a new client each time, we set defaults
// and the resource will use those when project/region match. For full integration testing,
// use acceptance tests.
func createMockConfig(t *testing.T, mockService *mockReserveIPServiceForAttachment, defaultProjectID, defaultRegion string) *config.Config {
	// Create a config - API keys are validated but won't be used in actual API calls
	// We use a valid format to pass validation
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Set defaults so resources can use them
	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	// Create a mock client with our mock service
	mockClient := &goe2e.Client{}
	mockClient.ReserveIP = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}

// createTestResourceData creates a schema.ResourceData for testing
func createTestResourceData(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	resource := ResourceFloatingIPAttachment()
	d := schema.TestResourceDataRaw(t, resource.Schema, data)
	return d
}

// ============================================================================
// Test: resourceFloatingIPAttachmentCreate
// ============================================================================

func TestResourceFloatingIPAttachmentCreate(t *testing.T) {
	tests := []struct {
		name          string
		resourceData  map[string]interface{}
		setupMock     func() *mockReserveIPServiceForAttachment
		expectedError bool
		errorContains string
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name: "successful create - single node",
			resourceData: map[string]interface{}{
				tfconstants.AttrIPAddress: "192.168.1.100",
				tfconstants.AttrNodeIDs:   []interface{}{"node-123"},
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					attachFloatingIPFunc: func(ctx context.Context, req *goe2e.FloatingIPAttachmentRequest) (*goe2e.Response, error) {
						if req.IPAddress != "192.168.1.100" {
							return nil, fmt.Errorf("unexpected IP address: %s", req.IPAddress)
						}
						if len(req.NodeIDs) != 1 || req.NodeIDs[0] != "node-123" {
							return nil, fmt.Errorf("unexpected node IDs: %v", req.NodeIDs)
						}
						return &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
			validateState: func(t *testing.T, d *schema.ResourceData) {
				if d.Id() != "192.168.1.100" {
					t.Errorf("Expected ID to be '192.168.1.100', got %s", d.Id())
				}
				if ipAddr := d.Get(tfconstants.AttrIPAddress).(string); ipAddr != "192.168.1.100" {
					t.Errorf("Expected IP address to be '192.168.1.100', got %s", ipAddr)
				}
				nodeIDs := d.Get(tfconstants.AttrNodeIDs).([]interface{})
				if len(nodeIDs) != 1 || nodeIDs[0].(string) != "node-123" {
					t.Errorf("Expected node IDs to be ['node-123'], got %v", nodeIDs)
				}
			},
		},
		{
			name: "successful create - multiple nodes",
			resourceData: map[string]interface{}{
				tfconstants.AttrIPAddress: "192.168.1.100",
				tfconstants.AttrNodeIDs:   []interface{}{"node-123", "node-456"},
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					attachFloatingIPFunc: func(ctx context.Context, req *goe2e.FloatingIPAttachmentRequest) (*goe2e.Response, error) {
						if len(req.NodeIDs) != 2 {
							return nil, fmt.Errorf("unexpected node IDs count: %d", len(req.NodeIDs))
						}
						return &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
		},
		{
			name: "error - empty node_ids",
			resourceData: map[string]interface{}{
				tfconstants.AttrIPAddress: "192.168.1.100",
				tfconstants.AttrNodeIDs:   []interface{}{},
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{}
			},
			expectedError: true,
			errorContains: ErrNodeIDsCannotBeEmpty,
		},
		{
			name: "error - attach API failure",
			resourceData: map[string]interface{}{
				tfconstants.AttrIPAddress: "192.168.1.100",
				tfconstants.AttrNodeIDs:   []interface{}{"node-123"},
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					attachFloatingIPFunc: func(ctx context.Context, req *goe2e.FloatingIPAttachmentRequest) (*goe2e.Response, error) {
						return nil, errors.New("API error: attachment failed")
					},
				}
			},
			expectedError: true,
			errorContains: "Error attaching floating IP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "project-123", "us-east-1")

			d := createTestResourceData(t, tt.resourceData)
			ctx := context.Background()

			diags := resourceFloatingIPAttachmentCreate(ctx, d, cfg)

			if tt.expectedError {
				if !diags.HasError() {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" {
					found := false
					for _, diag := range diags {
						if strings.Contains(diag.Summary, tt.errorContains) || strings.Contains(diag.Detail, tt.errorContains) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error to contain %q, got: %v", tt.errorContains, diags)
					}
				}
			} else {
				if diags.HasError() {
					t.Errorf("Unexpected error: %v", diags)
					return
				}
				if tt.validateState != nil {
					tt.validateState(t, d)
				}
			}
		})
	}
}

// ============================================================================
// Test: resourceFloatingIPAttachmentRead
// ============================================================================

func TestResourceFloatingIPAttachmentRead(t *testing.T) {
	tests := []struct {
		name          string
		resourceID    string
		resourceData  map[string]interface{}
		setupMock     func() *mockReserveIPServiceForAttachment
		expectedError bool
		errorContains string
		validateState func(*testing.T, *schema.ResourceData)
		expectRemoved bool
	}{
		{
			name:       "successful read - with attached nodes",
			resourceID: "192.168.1.100",
			resourceData: map[string]interface{}{
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{
							{
								IPAddress:    "192.168.1.100",
								ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
								FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{
									{ID: 123},
									{ID: 456},
								},
							},
						}, &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
			validateState: func(t *testing.T, d *schema.ResourceData) {
				nodeIDs := d.Get(tfconstants.AttrNodeIDs).([]interface{})
				if len(nodeIDs) != 2 {
					t.Errorf("Expected 2 node IDs, got %d", len(nodeIDs))
				}
				// Node IDs are converted from int to string
				if nodeIDs[0].(string) != "123" && nodeIDs[0].(string) != "456" {
					t.Errorf("Unexpected node ID: %v", nodeIDs[0])
				}
			},
		},
		{
			name:       "resource removed - IP not found",
			resourceID: "192.168.1.100",
			resourceData: map[string]interface{}{
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{}, &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
			expectRemoved: true,
		},
		{
			name:       "resource removed - not FloatingIP type",
			resourceID: "192.168.1.100",
			resourceData: map[string]interface{}{
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{
							{
								IPAddress:    "192.168.1.100",
								ReservedType: goe2econstants.ReserveIPTypePublicIP,
							},
						}, &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
			expectRemoved: true,
		},
		{
			name:       "resource removed - no attached nodes",
			resourceID: "192.168.1.100",
			resourceData: map[string]interface{}{
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{
							{
								IPAddress:               "192.168.1.100",
								ReservedType:            goe2econstants.ReserveIPTypeFloatingIP,
								FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{},
							},
						}, &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
			expectRemoved: true,
		},
		{
			name:       "error - list API failure",
			resourceID: "192.168.1.100",
			resourceData: map[string]interface{}{
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return nil, nil, errors.New("API error: list failed")
					},
				}
			},
			expectedError: true,
			errorContains: "Error retrieving reserved IPs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "project-123", "us-east-1")

			d := createTestResourceData(t, tt.resourceData)
			d.SetId(tt.resourceID)

			ctx := context.Background()
			diags := resourceFloatingIPAttachmentRead(ctx, d, cfg)

			if tt.expectedError {
				if !diags.HasError() {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" {
					found := false
					for _, diag := range diags {
						if strings.Contains(diag.Summary, tt.errorContains) || strings.Contains(diag.Detail, tt.errorContains) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error to contain %q, got: %v", tt.errorContains, diags)
					}
				}
			} else {
				if diags.HasError() {
					t.Errorf("Unexpected error: %v", diags)
					return
				}
				if tt.expectRemoved {
					if d.Id() != "" {
						t.Errorf("Expected resource to be removed (ID empty), but ID is: %s", d.Id())
					}
				} else {
					if tt.validateState != nil {
						tt.validateState(t, d)
					}
				}
			}
		})
	}
}

// ============================================================================
// Test: resourceFloatingIPAttachmentUpdate
// ============================================================================

func TestResourceFloatingIPAttachmentUpdate(t *testing.T) {
	tests := []struct {
		name          string
		resourceID    string
		oldData       map[string]interface{}
		newData       map[string]interface{}
		setupMock     func() *mockReserveIPServiceForAttachment
		expectedError bool
		errorContains string
	}{
		{
			name:       "successful update - add nodes",
			resourceID: "192.168.1.100",
			oldData: map[string]interface{}{
				tfconstants.AttrNodeIDs: []interface{}{"node-123"},
			},
			newData: map[string]interface{}{
				tfconstants.AttrNodeIDs: []interface{}{"node-123", "node-456"},
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					attachFloatingIPFunc: func(ctx context.Context, req *goe2e.FloatingIPAttachmentRequest) (*goe2e.Response, error) {
						if len(req.NodeIDs) != 1 || req.NodeIDs[0] != "node-456" {
							return nil, fmt.Errorf("unexpected attach request: %v", req.NodeIDs)
						}
						return &goe2e.Response{}, nil
					},
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{
							{
								IPAddress:    "192.168.1.100",
								ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
								FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{
									{ID: 123},
									{ID: 456},
								},
							},
						}, &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
		},
		{
			name:       "successful update - remove nodes",
			resourceID: "192.168.1.100",
			oldData: map[string]interface{}{
				tfconstants.AttrNodeIDs: []interface{}{"node-123", "node-456"},
			},
			newData: map[string]interface{}{
				tfconstants.AttrNodeIDs: []interface{}{"node-123"},
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					detachFloatingIPFunc: func(ctx context.Context, req *goe2e.FloatingIPDetachmentRequest) (*goe2e.Response, error) {
						if len(req.NodeIDs) != 1 || req.NodeIDs[0] != "node-456" {
							return nil, fmt.Errorf("unexpected detach request: %v", req.NodeIDs)
						}
						return &goe2e.Response{}, nil
					},
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{
							{
								IPAddress:    "192.168.1.100",
								ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
								FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{
									{ID: 123},
								},
							},
						}, &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
		},
		{
			name:       "error - empty node_ids after update",
			resourceID: "192.168.1.100",
			oldData: map[string]interface{}{
				tfconstants.AttrNodeIDs: []interface{}{"node-123"},
			},
			newData: map[string]interface{}{
				tfconstants.AttrNodeIDs: []interface{}{},
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{}
			},
			expectedError: true,
			errorContains: ErrNodeIDsCannotBeEmptyWithContext,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "project-123", "us-east-1")

			// For Update tests, we need to simulate GetChange behavior
			// Since GetChange requires Terraform's internal state tracking,
			// we'll test the validation logic separately and note that full
			// update testing is done via acceptance tests
			// For now, test the empty node_ids validation
			if tt.name == "error - empty node_ids after update" {
				// Test validation directly
				d := createTestResourceData(t, tt.newData)
				d.SetId(tt.resourceID)
				// Manually trigger the validation by calling update with empty node_ids
				// This tests the validation logic even if GetChange doesn't work in unit tests
				ctx := context.Background()
				diags := resourceFloatingIPAttachmentUpdate(ctx, d, cfg)

				if tt.expectedError {
					if !diags.HasError() {
						t.Errorf("Expected error but got none")
						return
					}
					if tt.errorContains != "" {
						found := false
						for _, diag := range diags {
							if strings.Contains(diag.Summary, tt.errorContains) || strings.Contains(diag.Detail, tt.errorContains) {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Expected error to contain %q, got: %v", tt.errorContains, diags)
						}
					}
				}
				return
			}

			// For other update tests, we'll skip them in unit tests as they require
			// proper GetChange simulation which is better tested in acceptance tests
			t.Skip("Update tests with GetChange require acceptance test framework")
		})
	}
}

// ============================================================================
// Test: resourceFloatingIPAttachmentDelete
// ============================================================================

func TestResourceFloatingIPAttachmentDelete(t *testing.T) {
	tests := []struct {
		name          string
		resourceID    string
		resourceData  map[string]interface{}
		setupMock     func() *mockReserveIPServiceForAttachment
		expectedError bool
		errorContains string
	}{
		{
			name:       "successful delete",
			resourceID: "192.168.1.100",
			resourceData: map[string]interface{}{
				tfconstants.AttrNodeIDs:   []interface{}{"node-123", "node-456"},
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					detachFloatingIPFunc: func(ctx context.Context, req *goe2e.FloatingIPDetachmentRequest) (*goe2e.Response, error) {
						if req.IPAddress != "192.168.1.100" {
							return nil, fmt.Errorf("unexpected IP address: %s", req.IPAddress)
						}
						if len(req.NodeIDs) != 2 {
							return nil, fmt.Errorf("unexpected node IDs count: %d", len(req.NodeIDs))
						}
						return &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
		},
		{
			name:       "error - detach API failure",
			resourceID: "192.168.1.100",
			resourceData: map[string]interface{}{
				tfconstants.AttrNodeIDs:   []interface{}{"node-123"},
				tfconstants.AttrProjectID: "project-123",
				tfconstants.AttrRegion:    "us-east-1",
			},
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					detachFloatingIPFunc: func(ctx context.Context, req *goe2e.FloatingIPDetachmentRequest) (*goe2e.Response, error) {
						return nil, errors.New("API error: detach failed")
					},
				}
			},
			expectedError: true,
			errorContains: "Error detaching floating IP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "project-123", "us-east-1")

			d := createTestResourceData(t, tt.resourceData)
			d.SetId(tt.resourceID)

			ctx := context.Background()
			diags := resourceFloatingIPAttachmentDelete(ctx, d, cfg)

			if tt.expectedError {
				if !diags.HasError() {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" {
					found := false
					for _, diag := range diags {
						if strings.Contains(diag.Summary, tt.errorContains) || strings.Contains(diag.Detail, tt.errorContains) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error to contain %q, got: %v", tt.errorContains, diags)
					}
				}
			} else {
				if diags.HasError() {
					t.Errorf("Unexpected error: %v", diags)
				}
				if d.Id() != "" {
					t.Errorf("Expected ID to be cleared after delete, got: %s", d.Id())
				}
			}
		})
	}
}

// ============================================================================
// Test: resourceFloatingIPAttachmentImport
// ============================================================================

func TestResourceFloatingIPAttachmentImport(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		setupMock     func() *mockReserveIPServiceForAttachment
		expectedError bool
		errorContains string
		validateState func(*testing.T, *schema.ResourceData)
	}{
		{
			name:     "successful import - valid format",
			importID: "project-123/us-east-1/192.168.1.100",
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{
							{
								IPAddress:    "192.168.1.100",
								ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
								FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{
									{ID: 123},
								},
							},
						}, &goe2e.Response{}, nil
					},
				}
			},
			expectedError: false,
			validateState: func(t *testing.T, d *schema.ResourceData) {
				if d.Id() != "192.168.1.100" {
					t.Errorf("Expected ID to be '192.168.1.100', got %s", d.Id())
				}
				if projectID := d.Get(tfconstants.AttrProjectID).(string); projectID != "project-123" {
					t.Errorf("Expected project ID to be 'project-123', got %s", projectID)
				}
				if region := d.Get(tfconstants.AttrRegion).(string); region != "us-east-1" {
					t.Errorf("Expected region to be 'us-east-1', got %s", region)
				}
			},
		},
		{
			name:          "error - invalid format - single part",
			importID:      "192.168.1.100",
			setupMock:     func() *mockReserveIPServiceForAttachment { return &mockReserveIPServiceForAttachment{} },
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "error - invalid format - two parts",
			importID:      "project-123/192.168.1.100",
			setupMock:     func() *mockReserveIPServiceForAttachment { return &mockReserveIPServiceForAttachment{} },
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "error - invalid format - four parts",
			importID:      "project-123/us-east-1/192.168.1.100/extra",
			setupMock:     func() *mockReserveIPServiceForAttachment { return &mockReserveIPServiceForAttachment{} },
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:     "error - read failure during import",
			importID: "project-123/us-east-1/192.168.1.100",
			setupMock: func() *mockReserveIPServiceForAttachment {
				return &mockReserveIPServiceForAttachment{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return nil, nil, errors.New("API error: read failed")
					},
				}
			},
			expectedError: true,
			errorContains: "error reading floating IP attachment during import",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			cfg := createMockConfig(t, mockService, "project-123", "us-east-1")

			resource := ResourceFloatingIPAttachment()
			d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
			d.SetId(tt.importID)

			ctx := context.Background()
			result, err := resourceFloatingIPAttachmentImport(ctx, d, cfg)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorContains != "" {
					if !strings.Contains(err.Error(), tt.errorContains) {
						t.Errorf("Expected error to contain %q, got: %v", tt.errorContains, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if len(result) != 1 {
					t.Errorf("Expected 1 resource data, got %d", len(result))
					return
				}
				if tt.validateState != nil {
					tt.validateState(t, result[0])
				}
			}
		})
	}
}

// ============================================================================
// Test: Error Message Constants
// ============================================================================

func TestErrorMessageConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		contains string
	}{
		{
			name:     "ErrNodeIDsCannotBeEmpty",
			constant: ErrNodeIDsCannotBeEmpty,
			contains: "node_ids cannot be empty",
		},
		{
			name:     "ErrNodeIDsCannotBeEmptyWithContext",
			constant: ErrNodeIDsCannotBeEmptyWithContext,
			contains: "node_ids cannot be empty",
		},
		{
			name:     "ErrAttachingFloatingIP",
			constant: ErrAttachingFloatingIP,
			contains: "Error attaching floating IP",
		},
		{
			name:     "ErrDetachingFloatingIP",
			constant: ErrDetachingFloatingIP,
			contains: "Error detaching floating IP",
		},
		{
			name:     "ErrRetrievingReservedIPs",
			constant: ErrRetrievingReservedIPs,
			contains: "Error retrieving reserved IPs",
		},
		{
			name:     "ErrReadingFloatingIPAttachmentDuringImport",
			constant: ErrReadingFloatingIPAttachmentDuringImport,
			contains: "error reading floating IP attachment during import",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.constant, tt.contains) {
				t.Errorf("constant %s should contain %q, got: %q", tt.name, tt.contains, tt.constant)
			}
		})
	}
}

// ============================================================================
// Test: Import Format Constant
// ============================================================================

func TestImportFormatConstant(t *testing.T) {
	if FloatingIPAttachmentImportFormat != "project_id/region/ip_address" {
		t.Errorf("Expected FloatingIPAttachmentImportFormat to be 'project_id/region/ip_address', got: %s", FloatingIPAttachmentImportFormat)
	}
}
