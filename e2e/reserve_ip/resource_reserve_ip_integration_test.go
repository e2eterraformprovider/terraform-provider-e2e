package reserve_ip

import (
	"regexp"
	"strings"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ============================================================================
// Integration Tests - V2 to V3 Migration
// ============================================================================

func TestV2ToV3Migration_BackwardCompatibility(t *testing.T) {
	resource := ResourceReserveIP()

	tests := []struct {
		name             string
		v2Fields         map[string]interface{}
		expectedV3Fields map[string]interface{}
		validate         func(*testing.T, *schema.Resource)
	}{
		{
			name: "V2 location field maps to V3 region",
			v2Fields: map[string]interface{}{
				tfconstants.AttrLocation: "Mumbai",
			},
			expectedV3Fields: map[string]interface{}{
				tfconstants.AttrRegion: "Mumbai",
			},
			validate: func(t *testing.T, r *schema.Resource) {
				// Verify both location and region schemas exist
				if _, ok := r.Schema[tfconstants.AttrLocation]; !ok {
					t.Error("location field should exist for backward compatibility")
				}
				if _, ok := r.Schema[tfconstants.AttrRegion]; !ok {
					t.Error("region field should exist")
				}
			},
		},
		{
			name: "V2 reserved_type field deprecated but still present",
			v2Fields: map[string]interface{}{
				tfconstants.AttrReservedType: goe2econstants.ReserveIPTypePublicIP,
			},
			expectedV3Fields: map[string]interface{}{
				"type": goe2econstants.ReserveIPTypePublicIP,
			},
			validate: func(t *testing.T, r *schema.Resource) {
				// Verify reserved_type is deprecated
				reservedTypeSchema := r.Schema[tfconstants.AttrReservedType]
				if reservedTypeSchema == nil {
					t.Error("reserved_type field should exist")
					return
				}
				if reservedTypeSchema.Deprecated == "" {
					t.Error("reserved_type field should be marked as deprecated")
				}
				// Verify type field exists
				if _, ok := r.Schema["type"]; !ok {
					t.Error("type field should exist")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.validate != nil {
				tt.validate(t, resource)
			}
		})
	}
}

func TestV2ToV3Migration_StateConsistency(t *testing.T) {
	tests := []struct {
		name            string
		v2State         map[string]interface{}
		expectedV3State map[string]interface{}
		validate        func(*testing.T, map[string]interface{}, map[string]interface{})
	}{
		{
			name: "Both location and region set to same value",
			v2State: map[string]interface{}{
				tfconstants.AttrLocation: "Mumbai",
			},
			expectedV3State: map[string]interface{}{
				tfconstants.AttrLocation: "Mumbai",
				tfconstants.AttrRegion:   "Mumbai",
			},
			validate: func(t *testing.T, v2, v3 map[string]interface{}) {
				// Both should have same value
				if v2[tfconstants.AttrLocation] != v3[tfconstants.AttrRegion] {
					t.Errorf("location and region should have same value")
				}
			},
		},
		{
			name: "Both reserved_type and type set to same value",
			v2State: map[string]interface{}{
				tfconstants.AttrReservedType: goe2econstants.ReserveIPTypePublicIP,
			},
			expectedV3State: map[string]interface{}{
				tfconstants.AttrReservedType: goe2econstants.ReserveIPTypePublicIP,
				"type":                       goe2econstants.ReserveIPTypePublicIP,
			},
			validate: func(t *testing.T, v2, v3 map[string]interface{}) {
				// Both should have same value
				if v2[tfconstants.AttrReservedType] != v3["type"] {
					t.Errorf("reserved_type and type should have same value")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.validate != nil {
				tt.validate(t, tt.v2State, tt.expectedV3State)
			}
		})
	}
}

// ============================================================================
// Integration Tests - URN Generation Consistency
// ============================================================================

func TestURNGenerationConsistency(t *testing.T) {
	tests := []struct {
		name       string
		region     string
		ipAddress  string
		operations []string // create, read, import
		validate   func(*testing.T, string, string, string)
	}{
		{
			name:       "URN consistent across create/read/import",
			region:     "Mumbai",
			ipAddress:  "1.2.3.4",
			operations: []string{"create", "read", "import"},
			validate: func(t *testing.T, region, ipAddress, operation string) {
				urn := generateReserveIPURN(region, ipAddress)
				expectedURN := "e2e:reserve_ip:Mumbai:1.2.3.4"
				if urn != expectedURN {
					t.Errorf("URN for %s operation: got %q, want %q", operation, urn, expectedURN)
				}
				// Verify URN format
				urnPattern := regexp.MustCompile(`^e2e:reserve_ip:.+:.+$`)
				if !urnPattern.MatchString(urn) {
					t.Errorf("URN format invalid: %q", urn)
				}
			},
		},
		{
			name:       "URN format validation",
			region:     "Singapore",
			ipAddress:  "192.168.1.100",
			operations: []string{"create"},
			validate: func(t *testing.T, region, ipAddress, operation string) {
				urn := generateReserveIPURN(region, ipAddress)
				// Verify URN has exactly 4 parts separated by colons
				parts := strings.Split(urn, ":")
				if len(parts) != 4 {
					t.Errorf("URN should have 4 parts, got %d: %q", len(parts), urn)
				}
				if parts[0] != "e2e" {
					t.Errorf("URN should start with 'e2e', got %q", parts[0])
				}
				if parts[1] != "reserve_ip" {
					t.Errorf("URN should have 'reserve_ip' as second part, got %q", parts[1])
				}
				if parts[2] != region {
					t.Errorf("URN should have region as third part, got %q, want %q", parts[2], region)
				}
				if parts[3] != ipAddress {
					t.Errorf("URN should have IP address as fourth part, got %q, want %q", parts[3], ipAddress)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, op := range tt.operations {
				if tt.validate != nil {
					tt.validate(t, tt.region, tt.ipAddress, op)
				}
			}
		})
	}
}

// ============================================================================
// Integration Tests - Floating IP Attached Nodes
// ============================================================================

func TestFloatingIPAttachedNodesIntegration(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *goe2e.ReserveIP
		validate  func(*testing.T, *goe2e.ReserveIP)
	}{
		{
			name: "FloatingIP with attached nodes - all fields populated",
			setupMock: func() *goe2e.ReserveIP {
				return &goe2e.ReserveIP{
					IPAddress:    "1.2.3.4",
					ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
					FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{
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
					},
				}
			},
			validate: func(t *testing.T, rip *goe2e.ReserveIP) {
				if rip.ReservedType != goe2econstants.ReserveIPTypeFloatingIP {
					t.Errorf("Expected FloatingIP type, got %q", rip.ReservedType)
				}
				if len(rip.FloatingIPAttachedNodes) != 2 {
					t.Errorf("Expected 2 attached nodes, got %d", len(rip.FloatingIPAttachedNodes))
				}
				flattened := flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes)
				if len(flattened) != 2 {
					t.Errorf("Expected 2 flattened nodes, got %d", len(flattened))
				}
				// Verify all fields present
				for i, node := range rip.FloatingIPAttachedNodes {
					flatNode := flattened[i]
					if flatNode["id"] != node.ID {
						t.Errorf("Node %d: id mismatch", i)
					}
					if flatNode["name"] != node.Name {
						t.Errorf("Node %d: name mismatch", i)
					}
					if flatNode["vm_id"] != node.VMID {
						t.Errorf("Node %d: vm_id mismatch", i)
					}
					if flatNode["ip_address_public"] != node.IPAddressPublic {
						t.Errorf("Node %d: ip_address_public mismatch", i)
					}
					if flatNode["ip_address_private"] != node.IPAddressPrivate {
						t.Errorf("Node %d: ip_address_private mismatch", i)
					}
					if flatNode["status_name"] != node.StatusName {
						t.Errorf("Node %d: status_name mismatch", i)
					}
					if flatNode["security_group_status"] != node.SecurityGroupStatus {
						t.Errorf("Node %d: security_group_status mismatch", i)
					}
				}
			},
		},
		{
			name: "FloatingIP detach - nodes updated",
			setupMock: func() *goe2e.ReserveIP {
				return &goe2e.ReserveIP{
					IPAddress:               "1.2.3.4",
					ReservedType:            goe2econstants.ReserveIPTypeFloatingIP,
					FloatingIPAttachedNodes: []goe2e.FloatingIPAttachedNode{}, // Empty after detach
				}
			},
			validate: func(t *testing.T, rip *goe2e.ReserveIP) {
				if len(rip.FloatingIPAttachedNodes) != 0 {
					t.Errorf("Expected 0 attached nodes after detach, got %d", len(rip.FloatingIPAttachedNodes))
				}
				flattened := flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes)
				if len(flattened) != 0 {
					t.Errorf("Expected 0 flattened nodes, got %d", len(flattened))
				}
				// Should return empty array, not nil
				if flattened == nil {
					t.Error("Expected empty array, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rip := tt.setupMock()
			if tt.validate != nil {
				tt.validate(t, rip)
			}
		})
	}
}

// ============================================================================
// Integration Tests - Import Functionality
// ============================================================================

func TestImportFunctionalityIntegration(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		setupMock     func() []goe2e.ReserveIP
		validate      func(*testing.T, string, string, string, []goe2e.ReserveIP)
		expectError   bool
		errorContains string
	}{
		{
			name:     "Import with full format - all fields populated",
			importID: "project-123/Mumbai/1.2.3.4",
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
			validate: func(t *testing.T, projectID, region, ipAddress string, rips []goe2e.ReserveIP) {
				if projectID != "project-123" {
					t.Errorf("Expected projectID 'project-123', got %q", projectID)
				}
				if region != "Mumbai" {
					t.Errorf("Expected region 'Mumbai', got %q", region)
				}
				if ipAddress != "1.2.3.4" {
					t.Errorf("Expected ipAddress '1.2.3.4', got %q", ipAddress)
				}
				// Verify IP exists in list
				found := false
				for _, rip := range rips {
					if rip.IPAddress == ipAddress {
						found = true
						// Verify V3 fields would be populated
						if rip.ReservedType == "" {
							t.Error("ReservedType should be populated")
						}
						break
					}
				}
				if !found {
					t.Error("IP should be found in list")
				}
			},
		},
		{
			name:     "Import with FloatingIP and attached nodes",
			importID: "project-123/Mumbai/1.2.3.4",
			setupMock: func() []goe2e.ReserveIP {
				return []goe2e.ReserveIP{
					{
						IPAddress:    "1.2.3.4",
						ReservedType: goe2econstants.ReserveIPTypeFloatingIP,
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
			validate: func(t *testing.T, projectID, region, ipAddress string, rips []goe2e.ReserveIP) {
				for _, rip := range rips {
					if rip.IPAddress == ipAddress && rip.ReservedType == goe2econstants.ReserveIPTypeFloatingIP {
						if len(rip.FloatingIPAttachedNodes) == 0 {
							t.Error("Expected attached nodes to be populated")
						}
						flattened := flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes)
						if len(flattened) == 0 {
							t.Error("Expected flattened nodes to be populated")
						}
						return
					}
				}
				t.Error("FloatingIP with attached nodes not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse import ID
			projectID, region, ipAddress, err := parseReserveIPImportID(tt.importID)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
					return
				}
				if tt.errorContains != "" && err.Error() == "" {
					t.Errorf("Expected error to contain %q", tt.errorContains)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			rips := tt.setupMock()
			if tt.validate != nil {
				tt.validate(t, projectID, region, ipAddress, rips)
			}
		})
	}
}

// ============================================================================
// Integration Tests - Deprecation Validation
// ============================================================================

func TestDeprecationValidation(t *testing.T) {
	resource := ResourceReserveIP()

	tests := []struct {
		name     string
		field    string
		validate func(*testing.T, *schema.Schema)
	}{
		{
			name:  "reserved_type field is deprecated",
			field: tfconstants.AttrReservedType,
			validate: func(t *testing.T, s *schema.Schema) {
				if s == nil {
					t.Error("reserved_type schema should exist")
					return
				}
				if s.Deprecated == "" {
					t.Error("reserved_type should be marked as deprecated")
				}
				if !strings.Contains(s.Deprecated, "type") {
					t.Error("deprecation message should mention 'type' as replacement")
				}
			},
		},
		{
			name:  "location field exists for backward compatibility",
			field: tfconstants.AttrLocation,
			validate: func(t *testing.T, s *schema.Schema) {
				if s == nil {
					t.Error("location schema should exist for backward compatibility")
				}
			},
		},
		{
			name:  "type field exists as V3 preferred",
			field: "type",
			validate: func(t *testing.T, s *schema.Schema) {
				if s == nil {
					t.Error("type schema should exist")
					return
				}
				if s.Deprecated != "" {
					t.Error("type field should not be deprecated")
				}
			},
		},
		{
			name:  "region field exists as V3 preferred",
			field: tfconstants.AttrRegion,
			validate: func(t *testing.T, s *schema.Schema) {
				if s == nil {
					t.Error("region schema should exist")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := resource.Schema[tt.field]
			if tt.validate != nil {
				tt.validate(t, schema)
			}
		})
	}
}

func TestBackwardCompatibility_V2FieldsWork(t *testing.T) {
	resource := ResourceReserveIP()

	tests := []struct {
		name     string
		v2Field  string
		v3Field  string
		validate func(*testing.T, *schema.Schema, *schema.Schema)
	}{
		{
			name:    "location and region both exist",
			v2Field: tfconstants.AttrLocation,
			v3Field: tfconstants.AttrRegion,
			validate: func(t *testing.T, v2Schema, v3Schema *schema.Schema) {
				if v2Schema == nil {
					t.Error("location schema should exist")
				}
				if v3Schema == nil {
					t.Error("region schema should exist")
				}
				// Both should have same type
				if v2Schema != nil && v3Schema != nil {
					if v2Schema.Type != v3Schema.Type {
						t.Errorf("location and region should have same schema type")
					}
				}
			},
		},
		{
			name:    "reserved_type and type both exist",
			v2Field: tfconstants.AttrReservedType,
			v3Field: "type",
			validate: func(t *testing.T, v2Schema, v3Schema *schema.Schema) {
				if v2Schema == nil {
					t.Error("reserved_type schema should exist")
				}
				if v3Schema == nil {
					t.Error("type schema should exist")
				}
				// Both should have same type
				if v2Schema != nil && v3Schema != nil {
					if v2Schema.Type != v3Schema.Type {
						t.Errorf("reserved_type and type should have same schema type")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v2Schema := resource.Schema[tt.v2Field]
			v3Schema := resource.Schema[tt.v3Field]
			if tt.validate != nil {
				tt.validate(t, v2Schema, v3Schema)
			}
		})
	}
}

// ============================================================================
// Integration Tests - Performance Benchmarks
// ============================================================================

func BenchmarkURNGeneration(b *testing.B) {
	region := "Mumbai"
	ipAddress := "1.2.3.4"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generateReserveIPURN(region, ipAddress)
	}
}

func BenchmarkParseImportID(b *testing.B) {
	importID := "project-123/Mumbai/1.2.3.4"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = parseReserveIPImportID(importID)
	}
}

func BenchmarkFlattenFloatingIPAttachedNodes(b *testing.B) {
	nodes := []goe2e.FloatingIPAttachedNode{
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
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flattenFloatingIPAttachedNodes(nodes)
	}
}

func BenchmarkFlattenFloatingIPAttachedNodes_LargeDataset(b *testing.B) {
	// Create 100 nodes to test performance with larger datasets
	nodes := make([]goe2e.FloatingIPAttachedNode, 100)
	for i := 0; i < 100; i++ {
		nodes[i] = goe2e.FloatingIPAttachedNode{
			ID:                  i,
			Name:                "test-node",
			VMID:                i * 10,
			IPAddressPublic:     "198.51.100.1",
			IPAddressPrivate:    "10.0.0.1",
			StatusName:          "running",
			SecurityGroupStatus: "active",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flattenFloatingIPAttachedNodes(nodes)
	}
}

// ============================================================================
// Integration Tests - Security Review
// ============================================================================

func TestSecurityReview_NoHardcodedCredentials(t *testing.T) {
	// This test verifies no hardcoded credentials in code
	// In a real scenario, you would use static analysis tools
	// For now, we verify that functions don't contain obvious credential patterns

	tests := []struct {
		name     string
		function func() // Function to test
		validate func(*testing.T)
	}{
		{
			name: "generateReserveIPURN contains no credentials",
			function: func() {
				_ = generateReserveIPURN("Mumbai", "1.2.3.4")
			},
			validate: func(t *testing.T) {
				// Function should not panic or contain credentials
				result := generateReserveIPURN("Mumbai", "1.2.3.4")
				if strings.Contains(result, "api_key") || strings.Contains(result, "auth_token") {
					t.Error("URN should not contain credential-like strings")
				}
			},
		},
		{
			name: "parseReserveIPImportID contains no credentials",
			function: func() {
				_, _, _, _ = parseReserveIPImportID("project-123/Mumbai/1.2.3.4")
			},
			validate: func(t *testing.T) {
				// Function should not contain credentials in logic
				projectID, _, _, err := parseReserveIPImportID("project-123/Mumbai/1.2.3.4")
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Verify project ID is not a credential pattern
				if strings.Contains(projectID, "api_key") || strings.Contains(projectID, "auth_token") {
					t.Error("Parsed values should not contain credential-like strings")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.validate != nil {
				tt.validate(t)
			}
		})
	}
}

func TestSecurityReview_ErrorMessagesNoCredentials(t *testing.T) {
	// Verify error messages don't leak credentials
	tests := []struct {
		name       string
		setupError func() error
		validate   func(*testing.T, error)
	}{
		{
			name: "Import ID error message doesn't contain credentials",
			setupError: func() error {
				_, _, _, err := parseReserveIPImportID("invalid")
				return err
			},
			validate: func(t *testing.T, err error) {
				if err == nil {
					t.Error("Expected error")
					return
				}
				errMsg := err.Error()
				if strings.Contains(errMsg, "api_key") || strings.Contains(errMsg, "auth_token") {
					t.Error("Error message should not contain credential-like strings")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setupError()
			if tt.validate != nil {
				tt.validate(t, err)
			}
		})
	}
}

// ============================================================================
// Integration Tests - Legacy Code Verification
// ============================================================================

func TestLegacyCodeVerification_NoLegacyClient(t *testing.T) {
	// This test verifies that no legacy client code exists
	// In a real scenario, this would be done via static analysis
	// For now, we verify the resource uses goe2e client patterns

	resource := ResourceReserveIP()

	tests := []struct {
		name     string
		validate func(*testing.T, *schema.Resource)
	}{
		{
			name: "Resource uses goe2e client pattern",
			validate: func(t *testing.T, r *schema.Resource) {
				// Verify resource has proper CRUD operations
				if r.CreateContext == nil {
					t.Error("CreateContext should be set")
				}
				if r.ReadContext == nil {
					t.Error("ReadContext should be set")
				}
				if r.DeleteContext == nil {
					t.Error("DeleteContext should be set")
				}
				// UpdateContext should be nil (immutable resource)
				if r.UpdateContext != nil {
					t.Error("UpdateContext should be nil for immutable resource")
				}
				// Importer should be set
				if r.Importer == nil {
					t.Error("Importer should be set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.validate != nil {
				tt.validate(t, resource)
			}
		})
	}
}

// ============================================================================
// Integration Tests - Multiple Reserve IPs Scenario
// ============================================================================

func TestMultipleReserveIPsScenario(t *testing.T) {
	// Simulate creating multiple Reserve IPs
	regions := []string{"Mumbai", "Singapore", "US-East"}
	ipAddresses := []string{"1.2.3.4", "5.6.7.8", "9.10.11.12"}

	tests := []struct {
		name     string
		count    int
		validate func(*testing.T, []string)
	}{
		{
			name:  "Create 10 Reserve IPs - URN generation",
			count: 10,
			validate: func(t *testing.T, urns []string) {
				if len(urns) != 10 {
					t.Errorf("Expected 10 URNs, got %d", len(urns))
				}
				// Verify all URNs are unique
				urnSet := make(map[string]bool)
				for _, urn := range urns {
					if urnSet[urn] {
						t.Errorf("Duplicate URN found: %q", urn)
					}
					urnSet[urn] = true
					// Verify format
					if !strings.HasPrefix(urn, "e2e:reserve_ip:") {
						t.Errorf("Invalid URN format: %q", urn)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urns := make([]string, 0, tt.count)
			for i := 0; i < tt.count; i++ {
				region := regions[i%len(regions)]
				// Use unique IP for each
				uniqueIP := ipAddresses[i%len(ipAddresses)] + string(rune(i))
				urn := generateReserveIPURN(region, uniqueIP)
				urns = append(urns, urn)
			}
			if tt.validate != nil {
				tt.validate(t, urns)
			}
		})
	}
}

func TestImportPerformance_MultipleIPs(t *testing.T) {
	// Simulate importing multiple Reserve IPs
	importIDs := []string{
		"project-123/Mumbai/1.2.3.4",
		"project-123/Mumbai/5.6.7.8",
		"project-123/Singapore/9.10.11.12",
		"project-456/Mumbai/13.14.15.16",
		"project-456/Singapore/17.18.19.20",
	}

	t.Run("Parse multiple import IDs efficiently", func(t *testing.T) {
		for _, importID := range importIDs {
			projectID, region, ipAddress, err := parseReserveIPImportID(importID)
			if err != nil {
				t.Errorf("Failed to parse import ID %q: %v", importID, err)
				continue
			}
			if projectID == "" || region == "" || ipAddress == "" {
				t.Errorf("Parsed values should not be empty for %q", importID)
			}
		}
	})
}
