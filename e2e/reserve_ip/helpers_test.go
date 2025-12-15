package reserve_ip

import (
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

// ============================================================================
// Basic smoke tests for helper functions
// Comprehensive unit tests are in resource_reserve_ip_unit_test.go
// ============================================================================

func TestGenerateReserveIPURN_Basic(t *testing.T) {
	tests := []struct {
		name      string
		region    string
		ipAddress string
		expected  string
	}{
		{
			name:      "basic URN generation",
			region:    "us-east",
			ipAddress: "164.52.220.153",
			expected:  "e2e:reserve_ip:us-east:164.52.220.153",
		},
		{
			name:      "different region",
			region:    "Mumbai",
			ipAddress: "192.168.1.1",
			expected:  "e2e:reserve_ip:Mumbai:192.168.1.1",
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

func TestParseReserveIPImportID_Basic(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		expectedProj  string
		expectedReg   string
		expectedIP    string
		expectedError bool
	}{
		{
			name:          "valid import ID",
			id:            "project-123/us-east/164.52.220.153",
			expectedProj:  "project-123",
			expectedReg:   "us-east",
			expectedIP:    "164.52.220.153",
			expectedError: false,
		},
		{
			name:          "invalid format - too few parts",
			id:            "project-123/us-east",
			expectedError: true,
		},
		{
			name:          "invalid format - too many parts",
			id:            "project-123/us-east/164.52.220.153/extra",
			expectedError: true,
		},
		{
			name:          "empty string",
			id:            "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proj, reg, ip, err := parseReserveIPImportID(tt.id)
			if tt.expectedError {
				if err == nil {
					t.Errorf("parseReserveIPImportID(%q) expected error, got nil", tt.id)
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

func TestFlattenFloatingIPAttachedNodes_Basic(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []goe2e.FloatingIPAttachedNode
		expected int
	}{
		{
			name:     "nil nodes",
			nodes:    nil,
			expected: 0,
		},
		{
			name:     "empty nodes",
			nodes:    []goe2e.FloatingIPAttachedNode{},
			expected: 0,
		},
		{
			name: "single node",
			nodes: []goe2e.FloatingIPAttachedNode{
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
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenFloatingIPAttachedNodes(tt.nodes)
			if len(result) != tt.expected {
				t.Errorf("flattenFloatingIPAttachedNodes() = %d items, want %d", len(result), tt.expected)
			}
		})
	}
}
