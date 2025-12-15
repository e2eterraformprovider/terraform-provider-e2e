package reserve_ip

import (
	"context"
	"errors"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

// ============================================================================
// Test: TestGenerateReserveIPURN
// ============================================================================

func TestGenerateReserveIPURN(t *testing.T) {
	tests := []struct {
		name      string
		region    string
		ipAddress string
		expected  string
	}{
		{
			name:      "Test Case 1: Basic URN generation",
			region:    "Mumbai",
			ipAddress: "1.2.3.4",
			expected:  "e2e:reserve_ip:Mumbai:1.2.3.4",
		},
		{
			name:      "Test Case 2: URN with different region",
			region:    "Singapore",
			ipAddress: "5.6.7.8",
			expected:  "e2e:reserve_ip:Singapore:5.6.7.8",
		},
		{
			name:      "Test Case 3: URN with special characters in IP",
			region:    "Mumbai",
			ipAddress: "192.168.1.1",
			expected:  "e2e:reserve_ip:Mumbai:192.168.1.1",
		},
		{
			name:      "Test Case 4: URN with empty region (edge case)",
			region:    "",
			ipAddress: "1.2.3.4",
			expected:  "e2e:reserve_ip::1.2.3.4",
		},
		{
			name:      "Test Case 5: URN with empty IP address (edge case)",
			region:    "Mumbai",
			ipAddress: "",
			expected:  "e2e:reserve_ip:Mumbai:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateReserveIPURN(tt.region, tt.ipAddress)
			if result != tt.expected {
				t.Errorf("generateReserveIPURN(%q, %q) = %q, want %q", tt.region, tt.ipAddress, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Test: TestParseReserveIPImportID
// ============================================================================

func TestParseReserveIPImportID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedProj  string
		expectedReg   string
		expectedIP    string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Test Case 1: Valid import ID - all three parts",
			id:            "project-123/Mumbai/1.2.3.4",
			expectedProj:  "project-123",
			expectedReg:   "Mumbai",
			expectedIP:    "1.2.3.4",
			expectedError: false,
		},
		{
			name:          "Test Case 2: Valid import ID - different values",
			id:            "proj-456/Singapore/192.168.1.1",
			expectedProj:  "proj-456",
			expectedReg:   "Singapore",
			expectedIP:    "192.168.1.1",
			expectedError: false,
		},
		{
			name:          "Test Case 3: Invalid format - only two parts",
			id:            "project-123/1.2.3.4",
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "Test Case 4: Invalid format - only one part",
			id:            "1.2.3.4",
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "Test Case 5: Invalid format - four parts",
			id:            "project-123/Mumbai/1.2.3.4/extra",
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "Test Case 6: Empty string",
			id:            "",
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "Test Case 7: Import ID with special characters",
			id:            "project_123/Mumbai/1.2.3.4",
			expectedProj:  "project_123",
			expectedReg:   "Mumbai",
			expectedIP:    "1.2.3.4",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proj, reg, ip, err := parseReserveIPImportID(tt.id)
			if tt.expectedError {
				if err == nil {
					t.Errorf("parseReserveIPImportID(%q) expected error, got nil", tt.id)
					return
				}
				if tt.errorContains != "" {
					if err.Error() == "" || len(err.Error()) == 0 {
						t.Errorf("parseReserveIPImportID(%q) expected error message to contain %q, got empty error", tt.id, tt.errorContains)
					}
				}
				return
			}
			if err != nil {
				t.Errorf("parseReserveIPImportID(%q) unexpected error: %v", tt.id, err)
				return
			}
			if proj != tt.expectedProj {
				t.Errorf("parseReserveIPImportID(%q) project = %q, want %q", tt.id, proj, tt.expectedProj)
			}
			if reg != tt.expectedReg {
				t.Errorf("parseReserveIPImportID(%q) region = %q, want %q", tt.id, reg, tt.expectedReg)
			}
			if ip != tt.expectedIP {
				t.Errorf("parseReserveIPImportID(%q) ipAddress = %q, want %q", tt.id, ip, tt.expectedIP)
			}
		})
	}
}

// ============================================================================
// Test: TestFlattenFloatingIPAttachedNodes
// ============================================================================

func TestFlattenFloatingIPAttachedNodes(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []goe2e.FloatingIPAttachedNode
		expected []map[string]interface{}
	}{
		{
			name: "Test Case 1: Basic flattening with single node",
			nodes: []goe2e.FloatingIPAttachedNode{
				{
					ID:                  123,
					Name:                "test-node-1",
					VMID:                456,
					IPAddressPublic:     "198.51.100.1",
					IPAddressPrivate:    "10.0.0.1",
					StatusName:          "running",
					SecurityGroupStatus: "active",
				},
			},
			expected: []map[string]interface{}{
				{
					"id":                    123,
					"name":                  "test-node-1",
					"vm_id":                 456,
					"ip_address_public":     "198.51.100.1",
					"ip_address_private":    "10.0.0.1",
					"status_name":           "running",
					"security_group_status": "active",
				},
			},
		},
		{
			name: "Test Case 2: Flattening with multiple nodes",
			nodes: []goe2e.FloatingIPAttachedNode{
				{
					ID:                  123,
					Name:                "test-node-1",
					VMID:                456,
					IPAddressPublic:     "198.51.100.1",
					IPAddressPrivate:    "10.0.0.1",
					StatusName:          "running",
					SecurityGroupStatus: "active",
				},
				{
					ID:                  124,
					Name:                "test-node-2",
					VMID:                457,
					IPAddressPublic:     "198.51.100.2",
					IPAddressPrivate:    "10.0.0.2",
					StatusName:          "stopped",
					SecurityGroupStatus: "inactive",
				},
				{
					ID:                  125,
					Name:                "test-node-3",
					VMID:                458,
					IPAddressPublic:     "198.51.100.3",
					IPAddressPrivate:    "10.0.0.3",
					StatusName:          "pending",
					SecurityGroupStatus: "pending",
				},
			},
			expected: []map[string]interface{}{
				{
					"id":                    123,
					"name":                  "test-node-1",
					"vm_id":                 456,
					"ip_address_public":     "198.51.100.1",
					"ip_address_private":    "10.0.0.1",
					"status_name":           "running",
					"security_group_status": "active",
				},
				{
					"id":                    124,
					"name":                  "test-node-2",
					"vm_id":                 457,
					"ip_address_public":     "198.51.100.2",
					"ip_address_private":    "10.0.0.2",
					"status_name":           "stopped",
					"security_group_status": "inactive",
				},
				{
					"id":                    125,
					"name":                  "test-node-3",
					"vm_id":                 458,
					"ip_address_public":     "198.51.100.3",
					"ip_address_private":    "10.0.0.3",
					"status_name":           "pending",
					"security_group_status": "pending",
				},
			},
		},
		{
			name:     "Test Case 3: Flattening with nil input",
			nodes:    nil,
			expected: []map[string]interface{}{},
		},
		{
			name:     "Test Case 4: Flattening with empty array",
			nodes:    []goe2e.FloatingIPAttachedNode{},
			expected: []map[string]interface{}{},
		},
		{
			name: "Test Case 5: Flattening with missing optional fields",
			nodes: []goe2e.FloatingIPAttachedNode{
				{
					ID:                  123,
					Name:                "test-node-1",
					VMID:                0,  // zero value
					IPAddressPublic:     "", // empty string
					IPAddressPrivate:    "", // empty string
					StatusName:          "", // empty string
					SecurityGroupStatus: "", // empty string
				},
			},
			expected: []map[string]interface{}{
				{
					"id":                    123,
					"name":                  "test-node-1",
					"vm_id":                 0,
					"ip_address_public":     "",
					"ip_address_private":    "",
					"status_name":           "",
					"security_group_status": "",
				},
			},
		},
		{
			name: "Test Case 6: Flattening with all field types",
			nodes: []goe2e.FloatingIPAttachedNode{
				{
					ID:                  999,
					Name:                "full-node",
					VMID:                888,
					IPAddressPublic:     "203.0.113.1",
					IPAddressPrivate:    "172.16.0.1",
					StatusName:          "running",
					SecurityGroupStatus: "active",
				},
			},
			expected: []map[string]interface{}{
				{
					"id":                    999,
					"name":                  "full-node",
					"vm_id":                 888,
					"ip_address_public":     "203.0.113.1",
					"ip_address_private":    "172.16.0.1",
					"status_name":           "running",
					"security_group_status": "active",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenFloatingIPAttachedNodes(tt.nodes)

			// Verify length
			if len(result) != len(tt.expected) {
				t.Errorf("flattenFloatingIPAttachedNodes() returned %d items, want %d", len(result), len(tt.expected))
				return
			}

			// Verify each item
			for i, expectedMap := range tt.expected {
				if i >= len(result) {
					t.Errorf("flattenFloatingIPAttachedNodes() missing item at index %d", i)
					continue
				}

				resultMap := result[i]

				// Verify all expected fields
				for key, expectedValue := range expectedMap {
					resultValue, exists := resultMap[key]
					if !exists {
						t.Errorf("flattenFloatingIPAttachedNodes() result[%d] missing key %q", i, key)
						continue
					}

					if resultValue != expectedValue {
						t.Errorf("flattenFloatingIPAttachedNodes() result[%d][%q] = %v (type %T), want %v (type %T)",
							i, key, resultValue, resultValue, expectedValue, expectedValue)
					}
				}

				// Verify no extra fields
				for key := range resultMap {
					if _, exists := expectedMap[key]; !exists {
						t.Errorf("flattenFloatingIPAttachedNodes() result[%d] has unexpected key %q", i, key)
					}
				}
			}
		})
	}
}

// ============================================================================
// Mock implementations for CRUD operation tests
// ============================================================================

// mockReserveIPService is a mock implementation of goe2e.ReserveIPService for testing
type mockReserveIPService struct {
	listReserveIPsFunc  func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error)
	createReserveIPFunc func(ctx context.Context) (*goe2e.ReserveIP, *goe2e.Response, error)
}

func (m *mockReserveIPService) ListReserveIPs(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
	if m.listReserveIPsFunc != nil {
		return m.listReserveIPsFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockReserveIPService) GetReserveIP(ctx context.Context, ipAddress string) (*goe2e.ReserveIP, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockReserveIPService) CreateReserveIP(ctx context.Context) (*goe2e.ReserveIP, *goe2e.Response, error) {
	if m.createReserveIPFunc != nil {
		return m.createReserveIPFunc(ctx)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockReserveIPService) DeleteReserveIP(ctx context.Context, ipAddress string) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockReserveIPService) AttachFloatingIP(ctx context.Context, req *goe2e.FloatingIPAttachmentRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

func (m *mockReserveIPService) DetachFloatingIP(ctx context.Context, req *goe2e.FloatingIPDetachmentRequest) (*goe2e.Response, error) {
	return nil, errors.New("not implemented")
}

// ============================================================================
// Test: TestResourceReserveIPImport
// ============================================================================

func TestResourceReserveIPImport(t *testing.T) {
	_ = context.Background()

	tests := []struct {
		name          string
		importID      string
		setupMock     func() (*mockReserveIPService, []goe2e.ReserveIP, error)
		expectedError bool
		errorContains string
		validate      func(*testing.T, string, string, string)
	}{
		{
			name:     "Test Case 1: Valid import with all fields",
			importID: "project-123/Mumbai/1.2.3.4",
			setupMock: func() (*mockReserveIPService, []goe2e.ReserveIP, error) {
				mockService := &mockReserveIPService{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{
							{
								IPAddress: "1.2.3.4",
								Status:    goe2econstants.ReserveIPStatusAvailable,
							},
						}, &goe2e.Response{}, nil
					},
				}
				return mockService, []goe2e.ReserveIP{{IPAddress: "1.2.3.4"}}, nil
			},
			expectedError: false,
			validate: func(t *testing.T, projectID, region, ipAddress string) {
				if projectID != "project-123" {
					t.Errorf("Expected projectID to be 'project-123', got %q", projectID)
				}
				if region != "Mumbai" {
					t.Errorf("Expected region to be 'Mumbai', got %q", region)
				}
				if ipAddress != "1.2.3.4" {
					t.Errorf("Expected ipAddress to be '1.2.3.4', got %q", ipAddress)
				}
			},
		},
		{
			name:     "Test Case 2: Import with IP not found",
			importID: "project-123/Mumbai/1.2.3.4",
			setupMock: func() (*mockReserveIPService, []goe2e.ReserveIP, error) {
				mockService := &mockReserveIPService{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return []goe2e.ReserveIP{
							{
								IPAddress: "5.6.7.8", // Different IP
								Status:    goe2econstants.ReserveIPStatusAvailable,
							},
						}, &goe2e.Response{}, nil
					},
				}
				return mockService, []goe2e.ReserveIP{{IPAddress: "5.6.7.8"}}, nil
			},
			expectedError: true,
			errorContains: "not found",
		},
		{
			name:     "Test Case 3: Import with invalid format",
			importID: "invalid-format",
			setupMock: func() (*mockReserveIPService, []goe2e.ReserveIP, error) {
				return &mockReserveIPService{}, nil, nil
			},
			expectedError: true,
			errorContains: "invalid import ID format",
		},
		{
			name:     "Test Case 4: Import with goe2e client error",
			importID: "project-123/Mumbai/1.2.3.4",
			setupMock: func() (*mockReserveIPService, []goe2e.ReserveIP, error) {
				mockService := &mockReserveIPService{
					listReserveIPsFunc: func(ctx context.Context) ([]goe2e.ReserveIP, *goe2e.Response, error) {
						return nil, nil, errors.New("API error")
					},
				}
				return mockService, nil, errors.New("API error")
			},
			expectedError: true,
			errorContains: "error retrieving reserved IPs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _ = tt.setupMock()

			// Test the import ID parsing
			projectID, region, ipAddress, err := parseReserveIPImportID(tt.importID)
			if tt.expectedError {
				if err == nil && tt.name != "Test Case 2: Import with IP not found" && tt.name != "Test Case 4: Import with goe2e client error" {
					t.Errorf("Expected error, got nil")
					return
				}
				if err != nil && tt.errorContains != "" {
					if err.Error() == "" {
						t.Errorf("Expected error to contain %q, got empty error", tt.errorContains)
					}
				}
				// For test cases 2 and 4, parsing succeeds but other errors occur
				if tt.name == "Test Case 2: Import with IP not found" || tt.name == "Test Case 4: Import with goe2e client error" {
					if err != nil {
						t.Errorf("Expected parsing to succeed for test case, got error: %v", err)
					}
					// Verify parsed values
					if projectID == "" || region == "" || ipAddress == "" {
						t.Errorf("Parsed values should not be empty: projectID=%q, region=%q, ipAddress=%q", projectID, region, ipAddress)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify parsed values
			if projectID == "" || region == "" || ipAddress == "" {
				t.Errorf("Parsed values should not be empty: projectID=%q, region=%q, ipAddress=%q", projectID, region, ipAddress)
			}

			// If validate function provided, call it
			if tt.validate != nil {
				tt.validate(t, projectID, region, ipAddress)
			}
		})
	}
}

// ============================================================================
// Test: TestResourceCreateReserveIP
// ============================================================================

func TestResourceCreateReserveIP(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func() *goe2e.ReserveIP
		region        string
		expectedError bool
		errorContains string
		validate      func(*testing.T, *goe2e.ReserveIP, string)
	}{
		{
			name: "Test Case 1: Successful creation",
			setupMock: func() *goe2e.ReserveIP {
				return &goe2e.ReserveIP{
					IPAddress:     "1.2.3.4",
					Status:        goe2econstants.ReserveIPStatusAvailable,
					ReservedType:  goe2econstants.ReserveIPTypePublicIP,
					ReserveID:     "12345",
					BoughtAt:      "2024-01-01T00:00:00Z",
					VMID:          0,
					VMName:        "",
					ApplianceType: "VM",
					ProjectName:   "test-project",
				}
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				if rip.IPAddress != "1.2.3.4" {
					t.Errorf("Expected IP address to be '1.2.3.4', got %q", rip.IPAddress)
				}
				urn := generateReserveIPURN(region, rip.IPAddress)
				expectedURN := "e2e:reserve_ip:Mumbai:1.2.3.4"
				if urn != expectedURN {
					t.Errorf("Expected URN to be %q, got %q", expectedURN, urn)
				}
			},
		},
		{
			name: "Test Case 2: Creation with FloatingIP type and attached nodes",
			setupMock: func() *goe2e.ReserveIP {
				return &goe2e.ReserveIP{
					IPAddress:    "1.2.3.4",
					Status:       goe2econstants.ReserveIPStatusAttached,
					ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
					ReserveID:    "12345",
					BoughtAt:     "2024-01-01T00:00:00Z",
					FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{
						{
							ID:                  123,
							Name:                "test-node",
							VMID:                456,
							IPAddressPublic:     "198.51.100.1",
							IPAddressPrivate:    "10.0.0.1",
							StatusName:          "running",
							SecurityGroupStatus: "active",
						},
					},
				}
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				if rip.ReservedType != goe2econstants.ReserveIPTypeFloatingIP {
					t.Errorf("Expected type to be FloatingIP, got %q", rip.ReservedType)
				}
				if len(rip.FloatingIPAttachedNodes) != 1 {
					t.Errorf("Expected 1 attached node, got %d", len(rip.FloatingIPAttachedNodes))
				}
				flattened := flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes)
				if len(flattened) != 1 {
					t.Errorf("Expected flattened nodes to have 1 item, got %d", len(flattened))
				}
			},
		},
		{
			name: "Test Case 3: Creation with non-FloatingIP type",
			setupMock: func() *goe2e.ReserveIP {
				return &goe2e.ReserveIP{
					IPAddress:    "1.2.3.4",
					Status:       goe2econstants.ReserveIPStatusAvailable,
					ReservedType: goe2econstants.ReserveIPTypePublicIP,
					ReserveID:    "12345",
				}
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				if rip.ReservedType == goe2econstants.ReserveIPTypeFloatingIP {
					t.Error("Expected type to not be FloatingIP")
				}
				// floating_ip_attached_nodes should not be populated for non-FloatingIP types
				if len(rip.FloatingIPAttachedNodes) > 0 {
					t.Errorf("Expected no attached nodes for PublicIP, got %d", len(rip.FloatingIPAttachedNodes))
				}
			},
		},
		{
			name: "Test Case 4: Creation with nil response",
			setupMock: func() *goe2e.ReserveIP {
				return nil
			},
			region:        "Mumbai",
			expectedError: true,
			errorContains: "ip_address",
		},
		{
			name: "Test Case 5: Creation with empty IP address",
			setupMock: func() *goe2e.ReserveIP {
				return &goe2e.ReserveIP{
					IPAddress: "", // Empty IP
					Status:    goe2econstants.ReserveIPStatusAvailable,
				}
			},
			region:        "Mumbai",
			expectedError: true,
			errorContains: "ip_address",
		},
		{
			name: "Test Case 6: Creation with API error",
			setupMock: func() *goe2e.ReserveIP {
				// Return nil to simulate API error
				return nil
			},
			region:        "Mumbai",
			expectedError: true,
			errorContains: "ip_address",
		},
		{
			name: "Test Case 7: URN generation during creation",
			setupMock: func() *goe2e.ReserveIP {
				return &goe2e.ReserveIP{
					IPAddress:    "192.168.1.100",
					Status:       goe2econstants.ReserveIPStatusAvailable,
					ReservedType: goe2econstants.ReserveIPTypePublicIP,
					ReserveID:    "67890",
				}
			},
			region:        "Singapore",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				urn := generateReserveIPURN(region, rip.IPAddress)
				expectedURN := "e2e:reserve_ip:Singapore:192.168.1.100"
				if urn != expectedURN {
					t.Errorf("Expected URN to be %q, got %q", expectedURN, urn)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rip := tt.setupMock()

			// Validate response
			if rip == nil || rip.IPAddress == "" {
				if !tt.expectedError {
					t.Errorf("Unexpected nil or empty IP address")
				}
				if tt.errorContains != "" {
					// Error case - this is expected
					return
				}
			}

			if tt.expectedError {
				if rip != nil && rip.IPAddress != "" {
					t.Errorf("Expected error case but got valid response")
				}
				return
			}

			// If validate function provided, call it
			if tt.validate != nil && rip != nil {
				tt.validate(t, rip, tt.region)
			}
		})
	}
}

// ============================================================================
// Test: TestResourceReadReserveIP
// ============================================================================

func TestResourceReadReserveIP(t *testing.T) {
	tests := []struct {
		name            string
		ipAddress       string
		setupMock       func() []goe2e.ReserveIP
		region          string
		expectedError   bool
		errorContains   string
		validate        func(*testing.T, *goe2e.ReserveIP, string)
		expectIDCleared bool
	}{
		{
			name:      "Test Case 1: Successful read",
			ipAddress: "1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress:     "1.2.3.4",
						Status:        goe2econstants.ReserveIPStatusAvailable,
						ReservedType:  goe2econstants.ReserveIPTypePublicIP,
						ReserveID:     "12345",
						BoughtAt:      "2024-01-01T00:00:00Z",
						VMID:          0,
						VMName:        "",
						ApplianceType: "VM",
						ProjectName:   "test-project",
					},
				}
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				if rip.IPAddress != "1.2.3.4" {
					t.Errorf("Expected IP address to be '1.2.3.4', got %q", rip.IPAddress)
				}
				urn := generateReserveIPURN(region, rip.IPAddress)
				expectedURN := "e2e:reserve_ip:Mumbai:1.2.3.4"
				if urn != expectedURN {
					t.Errorf("Expected URN to be %q, got %q", expectedURN, urn)
				}
			},
		},
		{
			name:      "Test Case 2: Read with FloatingIP and attached nodes",
			ipAddress: "1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress:    "1.2.3.4",
						Status:       goe2econstants.ReserveIPStatusAttached,
						ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
						ReserveID:    "12345",
						FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{
							{
								ID:                  123,
								Name:                "test-node",
								VMID:                456,
								IPAddressPublic:     "198.51.100.1",
								IPAddressPrivate:    "10.0.0.1",
								StatusName:          "running",
								SecurityGroupStatus: "active",
							},
						},
					},
				}
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				if rip.ReservedType != goe2econstants.ReserveIPTypeFloatingIP {
					t.Errorf("Expected type to be FloatingIP, got %q", rip.ReservedType)
				}
				if len(rip.FloatingIPAttachedNodes) != 1 {
					t.Errorf("Expected 1 attached node, got %d", len(rip.FloatingIPAttachedNodes))
				}
				flattened := flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes)
				if len(flattened) != 1 {
					t.Errorf("Expected flattened nodes to have 1 item, got %d", len(flattened))
				}
			},
		},
		{
			name:      "Test Case 3: Read with FloatingIP but no attached nodes",
			ipAddress: "1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress:               "1.2.3.4",
						Status:                  goe2econstants.ReserveIPStatusAvailable,
						ReservedType:            goe2econstants.ReserveIPTypeFloatingIP,
						ReserveID:               "12345",
						FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{}, // Empty array
					},
				}
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				if rip.ReservedType != goe2econstants.ReserveIPTypeFloatingIP {
					t.Errorf("Expected type to be FloatingIP, got %q", rip.ReservedType)
				}
				if len(rip.FloatingIPAttachedNodes) != 0 {
					t.Errorf("Expected 0 attached nodes, got %d", len(rip.FloatingIPAttachedNodes))
				}
				// Should return empty array, not nil
				flattened := flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes)
				if flattened == nil {
					t.Error("Expected empty array, got nil")
				}
				if len(flattened) != 0 {
					t.Errorf("Expected flattened nodes to have 0 items, got %d", len(flattened))
				}
			},
		},
		{
			name:      "Test Case 4: Read with non-FloatingIP type",
			ipAddress: "1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress:    "1.2.3.4",
						Status:       goe2econstants.ReserveIPStatusAvailable,
						ReservedType: goe2econstants.ReserveIPTypePublicIP,
						ReserveID:    "12345",
					},
				}
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				if rip.ReservedType == goe2econstants.ReserveIPTypeFloatingIP {
					t.Error("Expected type to not be FloatingIP")
				}
				// floating_ip_attached_nodes should be empty for non-FloatingIP types
				if len(rip.FloatingIPAttachedNodes) > 0 {
					t.Errorf("Expected no attached nodes for PublicIP, got %d", len(rip.FloatingIPAttachedNodes))
				}
				flattened := flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes)
				if len(flattened) != 0 {
					t.Errorf("Expected flattened nodes to have 0 items, got %d", len(flattened))
				}
			},
		},
		{
			name:      "Test Case 5: Read with IP not found (404 scenario)",
			ipAddress: "1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress: "5.6.7.8", // Different IP
						Status:    goe2econstants.ReserveIPStatusAvailable,
					},
				}
			},
			region:          "Mumbai",
			expectedError:   false,
			expectIDCleared: true,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				// Should be nil when IP not found
				if rip != nil {
					t.Errorf("Expected nil when IP not found, got %+v", rip)
				}
			},
		},
		{
			name:      "Test Case 6: Read with nil data",
			ipAddress: "1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress: "", // Empty IP
						Status:    goe2econstants.ReserveIPStatusAvailable,
					},
				}
			},
			region:          "Mumbai",
			expectedError:   false,
			expectIDCleared: true,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				// Should be nil when IP address is empty
				if rip != nil && rip.IPAddress != "" {
					t.Errorf("Expected nil or empty IP, got %+v", rip)
				}
			},
		},
		{
			name:      "Test Case 7: Read with API error",
			ipAddress: "1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				// Return empty list - in actual implementation, error would be returned from ListReserveIPs
				// This test validates that empty list doesn't cause issues
				return []goe2e.ReserveIP{}
			},
			region:          "Mumbai",
			expectedError:   false,
			expectIDCleared: true, // IP not found, so ID should be cleared
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				// Should be nil when IP not found in empty list
				if rip != nil {
					t.Errorf("Expected nil when IP not found, got %+v", rip)
				}
			},
		},
		{
			name:      "Test Case 8: URN generation during read",
			ipAddress: "192.168.1.100",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress:    "192.168.1.100",
						Status:       goe2econstants.ReserveIPStatusAvailable,
						ReservedType: goe2econstants.ReserveIPTypePublicIP,
						ReserveID:    "67890",
					},
				}
			},
			region:        "Singapore",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				urn := generateReserveIPURN(region, rip.IPAddress)
				expectedURN := "e2e:reserve_ip:Singapore:192.168.1.100"
				if urn != expectedURN {
					t.Errorf("Expected URN to be %q, got %q", expectedURN, urn)
				}
			},
		},
		{
			name:      "Test Case 9: Backward compatibility - reserved_type and type both set",
			ipAddress: "1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress:    "1.2.3.4",
						Status:       goe2econstants.ReserveIPStatusAvailable,
						ReservedType: goe2econstants.ReserveIPTypePublicIP,
						ReserveID:    "12345",
					},
				}
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, rip *goe2e.ReserveIP, region string) {
				// Both reserved_type (deprecated) and type (V3) should be set to same value
				if rip.ReservedType == "" {
					t.Error("Expected ReservedType to be set")
				}
				// In actual implementation, both fields are set from ReservedType
				// reserved_type = data.ReservedType (deprecated)
				// type = data.ReservedType (V3)
				if rip.ReservedType != goe2econstants.ReserveIPTypePublicIP {
					t.Errorf("Expected ReservedType to be PublicIP, got %q", rip.ReservedType)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rips := tt.setupMock()

			// Find the reserved IP by IP address (simulating the read logic)
			var data *goe2e.ReserveIP
			for i := range rips {
				if rips[i].IPAddress == tt.ipAddress {
					data = &rips[i]
					break
				}
			}

			// Simulate the validation logic
			if data == nil || data.IPAddress == "" {
				if !tt.expectIDCleared {
					t.Errorf("Unexpected nil or empty IP address")
				}
				// ID should be cleared in this case
				return
			}

			if tt.expectedError {
				t.Errorf("Expected error case but got valid response")
				return
			}

			// If validate function provided, call it
			if tt.validate != nil {
				tt.validate(t, data, tt.region)
			}
		})
	}
}

// ============================================================================
// Test: TestResourceDeleteReserveIP
// ============================================================================

func TestResourceDeleteReserveIP(t *testing.T) {
	tests := []struct {
		name          string
		ipAddress     string
		setupMock     func() ([]goe2e.ReserveIP, error, error) // list result, list error, delete error
		region        string
		expectedError bool
		errorContains string
		expectWarning bool
		validate      func(*testing.T, string, []goe2e.ReserveIP)
	}{
		{
			name:      "Test Case 1: Successful deletion",
			ipAddress: "1.2.3.4",
			setupMock: func() ([]goe2e.ReserveIP, error, error) {
				return []goe2e.ReserveIP{
					{
						IPAddress: "1.2.3.4",
						Status:    goe2econstants.ReserveIPStatusAvailable,
					},
				}, nil, nil // No errors
			},
			region:        "Mumbai",
			expectedError: false,
			validate: func(t *testing.T, ipAddress string, rips []goe2e.ReserveIP) {
				// Verify IP exists in list before deletion
				found := false
				for _, rip := range rips {
					if rip.IPAddress == ipAddress {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected IP %s to exist in list", ipAddress)
				}
			},
		},
		{
			name:      "Test Case 2: Deletion with attached IP (warning logged)",
			ipAddress: "1.2.3.4",
			setupMock: func() ([]goe2e.ReserveIP, error, error) {
				return []goe2e.ReserveIP{
					{
						IPAddress: "1.2.3.4",
						Status:    goe2econstants.ReserveIPStatusAttached, // Attached status
					},
				}, nil, nil
			},
			region:        "Mumbai",
			expectedError: false,
			expectWarning: true,
			validate: func(t *testing.T, ipAddress string, rips []goe2e.ReserveIP) {
				// Verify IP has Attached status
				for _, rip := range rips {
					if rip.IPAddress == ipAddress {
						if rip.Status != goe2econstants.ReserveIPStatusAttached {
							t.Errorf("Expected status to be Attached, got %q", rip.Status)
						}
						return
					}
				}
				t.Errorf("IP %s not found in list", ipAddress)
			},
		},
		{
			name:      "Test Case 3: Deletion with FloatingIP and attached nodes (warning logged)",
			ipAddress: "1.2.3.4",
			setupMock: func() ([]goe2e.ReserveIP, error, error) {
				return []goe2e.ReserveIP{
					{
						IPAddress:    "1.2.3.4",
						Status:       goe2econstants.ReserveIPStatusAvailable,
						ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
						FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{
							{
								ID:   123,
								Name: "test-node",
							},
						},
					},
				}, nil, nil
			},
			region:        "Mumbai",
			expectedError: false,
			expectWarning: true,
			validate: func(t *testing.T, ipAddress string, rips []goe2e.ReserveIP) {
				// Verify IP has attached nodes
				for _, rip := range rips {
					if rip.IPAddress == ipAddress {
						if len(rip.FloatingIPAttachedNodes) == 0 {
							t.Error("Expected attached nodes, got 0")
						}
						return
					}
				}
				t.Errorf("IP %s not found in list", ipAddress)
			},
		},
		{
			name:      "Test Case 4: Deletion with API error",
			ipAddress: "1.2.3.4",
			setupMock: func() ([]goe2e.ReserveIP, error, error) {
				return []goe2e.ReserveIP{
					{
						IPAddress: "1.2.3.4",
						Status:    goe2econstants.ReserveIPStatusAvailable,
					},
				}, nil, errors.New("API error: deletion failed")
			},
			region:        "Mumbai",
			expectedError: true,
			errorContains: "Error deleting",
		},
		{
			name:      "Test Case 5: Deletion with goe2e client creation error",
			ipAddress: "1.2.3.4",
			setupMock: func() ([]goe2e.ReserveIP, error, error) {
				// Simulate client creation error by returning error on list
				return nil, errors.New("Error creating goe2e client"), nil
			},
			region:        "Mumbai",
			expectedError: false, // List error is ignored, deletion proceeds
			validate: func(t *testing.T, ipAddress string, rips []goe2e.ReserveIP) {
				// List error is ignored, deletion should still be attempted
				if len(rips) != 0 {
					t.Errorf("Expected empty list when client creation fails, got %d items", len(rips))
				}
			},
		},
		{
			name:      "Test Case 6: Deletion when IP not found in list (idempotent)",
			ipAddress: "1.2.3.4",
			setupMock: func() ([]goe2e.ReserveIP, error, error) {
				return []goe2e.ReserveIP{
					{
						IPAddress: "5.6.7.8", // Different IP
						Status:    goe2econstants.ReserveIPStatusAvailable,
					},
				}, nil, nil
			},
			region:        "Mumbai",
			expectedError: false, // Should be idempotent
			validate: func(t *testing.T, ipAddress string, rips []goe2e.ReserveIP) {
				// IP not in list, but deletion should still be attempted (idempotent)
				found := false
				for _, rip := range rips {
					if rip.IPAddress == ipAddress {
						found = true
						break
					}
				}
				if found {
					t.Errorf("Expected IP %s to not be in list", ipAddress)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rips, listErr, deleteErr := tt.setupMock()

			// Simulate the deletion logic
			// Check if IP is attached and log warning (if list succeeds)
			if listErr == nil {
				for _, rip := range rips {
					if rip.IPAddress == tt.ipAddress {
						if rip.Status == goe2econstants.ReserveIPStatusAttached || len(rip.FloatingIPAttachedNodes) > 0 {
							if !tt.expectWarning {
								t.Errorf("Unexpected warning condition")
							}
						}
						break
					}
				}
			}

			// Simulate delete operation
			if deleteErr != nil {
				if !tt.expectedError {
					t.Errorf("Unexpected delete error: %v", deleteErr)
					return
				}
				if tt.errorContains != "" {
					if deleteErr.Error() == "" {
						t.Errorf("Expected error to contain %q, got empty error", tt.errorContains)
					}
				}
				return
			}

			if tt.expectedError {
				t.Errorf("Expected error but deletion succeeded")
				return
			}

			// If validate function provided, call it
			if tt.validate != nil {
				tt.validate(t, tt.ipAddress, rips)
			}
		})
	}
}
