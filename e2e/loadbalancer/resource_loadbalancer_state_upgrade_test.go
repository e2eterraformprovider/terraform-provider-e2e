package loadbalancer_test

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/loadbalancer"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
)

// TestResourceLoadBalancerStateUpgradeV0toV1_Basic tests basic state upgrade with minimal fields
func TestResourceLoadBalancerStateUpgradeV0toV1_Basic(t *testing.T) {
	tests := []struct {
		name     string
		v0State  map[string]interface{}
		validate func(t *testing.T, v1State map[string]interface{})
	}{
		{
			name: "minimal required fields only",
			v0State: map[string]interface{}{
				"id":                   "lb-12345",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				assert.Equal(t, "lb-12345", v1State["id"])
				assert.Equal(t, goe2econstants.LBPlanE2ELB2, v1State[tfconstants.AttrPlan])
				assert.Equal(t, goe2econstants.LBModeHTTP, v1State[tfconstants.AttrLbMode])

				// name field set from lb_name
				assert.Equal(t, "test-lb", v1State[tfconstants.AttrName])
				// lb_name preserved
				assert.Equal(t, "test-lb", v1State[tfconstants.AttrLbName])

				// state initialized
				assert.NotNil(t, v1State[tfconstants.AttrState])

				// tags initialized as empty map
				assert.NotNil(t, v1State[tfconstants.AttrTags])
				tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
				assert.True(t, ok, "tags should be a map")
				assert.Empty(t, tags, "tags should be empty map")
			},
		},
		{
			name: "all possible V0 fields present",
			v0State: map[string]interface{}{
				"id":                             "lb-67890",
				tfconstants.AttrPlan:             goe2econstants.LBPlanE2ELB3,
				tfconstants.AttrLbMode:           goe2econstants.LBModeHTTPS,
				tfconstants.AttrLbName:           "production-lb",
				tfconstants.AttrLocation:         "Mumbai",
				tfconstants.AttrProjectID:        "project-123",
				"lb_reserve_ip":                  "ip-456",
				"public_ip":                      "203.0.113.10",
				"private_ip":                     "10.0.1.5",
				tfconstants.AttrStatus:           goe2econstants.LBStatusRunning,
				tfconstants.AttrLbType:           goe2econstants.LBTypeExternal,
				"node_list_type":                 goe2econstants.LBNodeListTypeStatic,
				"enable_bitninja":                true,
				tfconstants.AttrRAM:              "8192",
				tfconstants.AttrDisk:             "100",
				tfconstants.AttrVCPU:             4.0,
				"host_target_ipv6":               "2001:db8::1",
				tfconstants.AttrPublicIPAddress:  "203.0.113.10",
				tfconstants.AttrPrivateIPAddress: "10.0.1.5",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// All V0 fields should be preserved
				assert.Equal(t, "lb-67890", v1State["id"])
				assert.Equal(t, goe2econstants.LBPlanE2ELB3, v1State[tfconstants.AttrPlan])
				assert.Equal(t, goe2econstants.LBModeHTTPS, v1State[tfconstants.AttrLbMode])
				assert.Equal(t, goe2econstants.LBTypeExternal, v1State[tfconstants.AttrLbType])
				assert.Equal(t, goe2econstants.LBNodeListTypeStatic, v1State["node_list_type"])
				assert.Equal(t, true, v1State["enable_bitninja"])
				assert.Equal(t, "8192", v1State[tfconstants.AttrRAM])
				assert.Equal(t, "100", v1State[tfconstants.AttrDisk])
				assert.Equal(t, 4.0, v1State[tfconstants.AttrVCPU])
				assert.Equal(t, "2001:db8::1", v1State["host_target_ipv6"])
				assert.Equal(t, "project-123", v1State[tfconstants.AttrProjectID])

				// V3 fields added
				assert.NotNil(t, v1State[tfconstants.AttrState])
				assert.NotNil(t, v1State[tfconstants.AttrTags])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1State, err := loadbalancer.ResourceLoadBalancerStateUpgradeV0toV1(context.Background(), tt.v0State, nil)
			assert.NoError(t, err)
			tt.validate(t, v1State)
		})
	}
}

// TestResourceLoadBalancerStateUpgradeV0toV1_AllFields tests upgrade with all V0 fields present
func TestResourceLoadBalancerStateUpgradeV0toV1_AllFields(t *testing.T) {
	tests := []struct {
		name     string
		v0State  map[string]interface{}
		validate func(t *testing.T, v1State map[string]interface{})
	}{
		{
			name: "all renamed fields present",
			v0State: map[string]interface{}{
				"id":                      "lb-99999",
				tfconstants.AttrPlan:      goe2econstants.LBPlanE2ELB4,
				tfconstants.AttrLbMode:    goe2econstants.LBModeBoth,
				tfconstants.AttrLbName:    "full-featured-lb",
				tfconstants.AttrLocation:  "Delhi",
				tfconstants.AttrProjectID: "project-789",
				"lb_reserve_ip":           "ip-reserved-123",
				"public_ip":               "198.51.100.42",
				"private_ip":              "172.16.0.10",
				tfconstants.AttrStatus:    goe2econstants.LBStatusRunning,
				tfconstants.AttrRAM:       "16384",
				tfconstants.AttrDisk:      "200",
				tfconstants.AttrVCPU:      8.0,
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// region set from location (location preserved)
				assert.Equal(t, "Delhi", v1State[tfconstants.AttrRegion])
				assert.Equal(t, "Delhi", v1State[tfconstants.AttrLocation])

				// name set from lb_name (lb_name preserved)
				assert.Equal(t, "full-featured-lb", v1State[tfconstants.AttrName])
				assert.Equal(t, "full-featured-lb", v1State[tfconstants.AttrLbName])

				// floating_ip_id set from lb_reserve_ip (lb_reserve_ip preserved)
				assert.Equal(t, "ip-reserved-123", v1State["floating_ip_id"])
				assert.Equal(t, "ip-reserved-123", v1State["lb_reserve_ip"])

				// public_ip_address set from public_ip (public_ip preserved)
				assert.Equal(t, "198.51.100.42", v1State[tfconstants.AttrPublicIPAddress])
				assert.Equal(t, "198.51.100.42", v1State["public_ip"])

				// private_ip_address set from private_ip (private_ip preserved)
				assert.Equal(t, "172.16.0.10", v1State[tfconstants.AttrPrivateIPAddress])
				assert.Equal(t, "172.16.0.10", v1State["private_ip"])

				// state initialized (normalized from status)
				assert.Equal(t, goe2econstants.LBStateRunning, v1State[tfconstants.AttrState])

				// tags initialized (empty map)
				tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
				assert.True(t, ok)
				assert.Empty(t, tags)

				// All field renames applied, deprecated fields preserved
				assert.NotNil(t, v1State["id"])
			},
		},
		{
			name: "nested structures preserved",
			v0State: map[string]interface{}{
				"id":                   "lb-complex-123",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB5,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTPS,
				tfconstants.AttrLbName: "complex-lb",
				tfconstants.AttrBackends: []interface{}{
					map[string]interface{}{
						"name":    "backend-1",
						"balance": goe2econstants.LBBalanceRoundRobin,
						"servers": []interface{}{
							map[string]interface{}{
								"id":   "node-1",
								"port": "8080",
							},
						},
					},
				},
				"tcp_backend": []interface{}{
					map[string]interface{}{
						"backend_name": "tcp-backend-1",
						"port":         "3306",
						"balance":      goe2econstants.LBBalanceSource,
						"servers": []interface{}{
							map[string]interface{}{
								"id":   "node-2",
								"port": "3306",
							},
						},
					},
				},
				"acl_list": []interface{}{
					map[string]interface{}{
						"acl_name":          "acl-1",
						"acl_condition":     "path_beg",
						"acl_matching_path": "/api",
					},
				},
				"acl_map": []interface{}{
					map[string]interface{}{
						"acl_name":            "acl-1",
						"acl_condition_state": true,
						"acl_backend":         "backend-1",
					},
				},
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Nested structures should be intact
				assert.NotNil(t, v1State[tfconstants.AttrBackends])
				backends := v1State[tfconstants.AttrBackends].([]interface{})
				assert.Len(t, backends, 1)

				assert.NotNil(t, v1State["tcp_backend"])
				tcpBackend := v1State["tcp_backend"].([]interface{})
				assert.Len(t, tcpBackend, 1)

				assert.NotNil(t, v1State["acl_list"])
				aclList := v1State["acl_list"].([]interface{})
				assert.Len(t, aclList, 1)

				assert.NotNil(t, v1State["acl_map"])
				aclMap := v1State["acl_map"].([]interface{})
				assert.Len(t, aclMap, 1)

				// No nested data loss, structures preserved
				assert.Equal(t, "complex-lb", v1State[tfconstants.AttrName])
				assert.NotNil(t, v1State[tfconstants.AttrTags])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1State, err := loadbalancer.ResourceLoadBalancerStateUpgradeV0toV1(context.Background(), tt.v0State, nil)
			assert.NoError(t, err)
			tt.validate(t, v1State)
		})
	}
}

// TestResourceLoadBalancerStateUpgradeV0toV1_PreservesDeprecated tests that deprecated fields are not removed
func TestResourceLoadBalancerStateUpgradeV0toV1_PreservesDeprecated(t *testing.T) {
	tests := []struct {
		name     string
		v0State  map[string]interface{}
		validate func(t *testing.T, v1State map[string]interface{})
	}{
		{
			name: "deprecated fields not removed",
			v0State: map[string]interface{}{
				"id":                     "lb-deprecated-test",
				tfconstants.AttrPlan:     goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode:   goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName:   "deprecated-lb",
				tfconstants.AttrLocation: "Mumbai",
				"lb_reserve_ip":          "ip-123",
				"public_ip":              "192.0.2.1",
				"private_ip":             "10.0.0.1",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// V1 state contains both old and new field names
				assert.Equal(t, "Mumbai", v1State[tfconstants.AttrRegion], "new field should be set")
				assert.Equal(t, "Mumbai", v1State[tfconstants.AttrLocation], "old field should be preserved")

				assert.Equal(t, "deprecated-lb", v1State[tfconstants.AttrName], "new field should be set")
				assert.Equal(t, "deprecated-lb", v1State[tfconstants.AttrLbName], "old field should be preserved")

				assert.Equal(t, "ip-123", v1State["floating_ip_id"], "new field should be set")
				assert.Equal(t, "ip-123", v1State["lb_reserve_ip"], "old field should be preserved")

				assert.Equal(t, "192.0.2.1", v1State[tfconstants.AttrPublicIPAddress], "new field should be set")
				assert.Equal(t, "192.0.2.1", v1State["public_ip"], "old field should be preserved")

				assert.Equal(t, "10.0.0.1", v1State[tfconstants.AttrPrivateIPAddress], "new field should be set")
				assert.Equal(t, "10.0.0.1", v1State["private_ip"], "old field should be preserved")

				// Backward compatibility maintained
			},
		},
		{
			name: "no forced recreation - resource ID unchanged",
			v0State: map[string]interface{}{
				"id":                   "lb-no-recreation",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB3,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTPS,
				tfconstants.AttrLbName: "stable-lb",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Resource ID unchanged, no recreation triggered
				assert.Equal(t, "lb-no-recreation", v1State["id"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1State, err := loadbalancer.ResourceLoadBalancerStateUpgradeV0toV1(context.Background(), tt.v0State, nil)
			assert.NoError(t, err)
			tt.validate(t, v1State)
		})
	}
}

// TestResourceLoadBalancerStateUpgradeV0toV1_FieldRenaming tests each field rename individually
func TestResourceLoadBalancerStateUpgradeV0toV1_FieldRenaming(t *testing.T) {
	tests := []struct {
		name     string
		v0State  map[string]interface{}
		validate func(t *testing.T, v1State map[string]interface{})
	}{
		{
			name: "location → region",
			v0State: map[string]interface{}{
				"id":                     "lb-location-test",
				tfconstants.AttrPlan:     goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode:   goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName:   "test-lb",
				tfconstants.AttrLocation: "Mumbai",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Both fields present
				assert.Equal(t, "Mumbai", v1State[tfconstants.AttrRegion])
				assert.Equal(t, "Mumbai", v1State[tfconstants.AttrLocation])
			},
		},
		{
			name: "lb_name → name",
			v0State: map[string]interface{}{
				"id":                   "lb-name-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "my-lb",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Both fields present
				assert.Equal(t, "my-lb", v1State[tfconstants.AttrName])
				assert.Equal(t, "my-lb", v1State[tfconstants.AttrLbName])
			},
		},
		{
			name: "lb_reserve_ip → floating_ip_id",
			v0State: map[string]interface{}{
				"id":                   "lb-reserve-ip-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				"lb_reserve_ip":        "ip-123",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Both fields present
				assert.Equal(t, "ip-123", v1State["floating_ip_id"])
				assert.Equal(t, "ip-123", v1State["lb_reserve_ip"])
			},
		},
		{
			name: "public_ip → public_ip_address",
			v0State: map[string]interface{}{
				"id":                   "lb-public-ip-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				"public_ip":            "1.2.3.4",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Both fields present
				assert.Equal(t, "1.2.3.4", v1State[tfconstants.AttrPublicIPAddress])
				assert.Equal(t, "1.2.3.4", v1State["public_ip"])
			},
		},
		{
			name: "private_ip → private_ip_address",
			v0State: map[string]interface{}{
				"id":                   "lb-private-ip-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				"private_ip":           "10.0.0.1",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Both fields present
				assert.Equal(t, "10.0.0.1", v1State[tfconstants.AttrPrivateIPAddress])
				assert.Equal(t, "10.0.0.1", v1State["private_ip"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1State, err := loadbalancer.ResourceLoadBalancerStateUpgradeV0toV1(context.Background(), tt.v0State, nil)
			assert.NoError(t, err)
			tt.validate(t, v1State)
		})
	}
}

// TestResourceLoadBalancerStateUpgradeV0toV1_ComputedFields tests computed field initialization
func TestResourceLoadBalancerStateUpgradeV0toV1_ComputedFields(t *testing.T) {
	tests := []struct {
		name     string
		v0State  map[string]interface{}
		validate func(t *testing.T, v1State map[string]interface{})
	}{
		{
			name: "state field initialization from status",
			v0State: map[string]interface{}{
				"id":                   "lb-state-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				tfconstants.AttrStatus: goe2econstants.LBStatusRunning,
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Computed field added - normalized from status
				assert.Equal(t, goe2econstants.LBStateRunning, v1State[tfconstants.AttrState])
			},
		},
		{
			name: "state field initialization - no status",
			v0State: map[string]interface{}{
				"id":                   "lb-no-status-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// State initialized as empty string when status not present
				assert.Equal(t, "", v1State[tfconstants.AttrState])
			},
		},
		{
			name: "tags field initialization",
			v0State: map[string]interface{}{
				"id":                   "lb-tags-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Tags field added as empty map
				tags, ok := v1State[tfconstants.AttrTags].(map[string]interface{})
				assert.True(t, ok)
				assert.Empty(t, tags)
			},
		},
		{
			name: "computed fields from API preserved",
			v0State: map[string]interface{}{
				"id":                   "lb-computed-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				tfconstants.AttrStatus: goe2econstants.LBStatusCreating,
				tfconstants.AttrRAM:    "4096",
				tfconstants.AttrDisk:   "50",
				tfconstants.AttrVCPU:   2.0,
				tfconstants.AttrState:  "custom-state", // Already set
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// No overwriting of existing computed fields
				assert.Equal(t, "custom-state", v1State[tfconstants.AttrState])
				assert.Equal(t, "4096", v1State[tfconstants.AttrRAM])
				assert.Equal(t, "50", v1State[tfconstants.AttrDisk])
				assert.Equal(t, 2.0, v1State[tfconstants.AttrVCPU])
			},
		},
		{
			name: "state normalization - Creating status",
			v0State: map[string]interface{}{
				"id":                   "lb-creating-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				tfconstants.AttrStatus: goe2econstants.LBStatusCreating,
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				assert.Equal(t, goe2econstants.LBStateCreating, v1State[tfconstants.AttrState])
			},
		},
		{
			name: "state normalization - Powered off status",
			v0State: map[string]interface{}{
				"id":                   "lb-poweredoff-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				tfconstants.AttrStatus: goe2econstants.LBStatusPoweredOff,
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				assert.Equal(t, goe2econstants.LBStateStopped, v1State[tfconstants.AttrState])
			},
		},
		{
			name: "state normalization - Upgrading status",
			v0State: map[string]interface{}{
				"id":                   "lb-upgrading-test",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				tfconstants.AttrStatus: goe2econstants.LBStatusUpgrading,
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				assert.Equal(t, goe2econstants.LBStateUpgrading, v1State[tfconstants.AttrState])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1State, err := loadbalancer.ResourceLoadBalancerStateUpgradeV0toV1(context.Background(), tt.v0State, nil)
			assert.NoError(t, err)
			tt.validate(t, v1State)
		})
	}
}

// TestResourceLoadBalancerStateUpgradeV0toV1_EdgeCases tests edge cases
func TestResourceLoadBalancerStateUpgradeV0toV1_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		v0State  map[string]interface{}
		validate func(t *testing.T, v1State map[string]interface{})
	}{
		{
			name: "empty state upgrade - minimal ID only",
			v0State: map[string]interface{}{
				"id": "lb-minimal",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Upgrade handles empty state gracefully
				assert.Equal(t, "lb-minimal", v1State["id"])
				assert.NotNil(t, v1State[tfconstants.AttrTags])
				assert.NotNil(t, v1State[tfconstants.AttrState])
			},
		},
		{
			name: "state with nil/empty optional fields",
			v0State: map[string]interface{}{
				"id":                     "lb-nil-test",
				tfconstants.AttrPlan:     goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode:   goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName:   "test-lb",
				tfconstants.AttrLocation: "",
				"lb_reserve_ip":          "",
				"public_ip":              "",
				"private_ip":             "",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// No nil pointer dereferences
				assert.Equal(t, "lb-nil-test", v1State["id"])
				// Empty string fields should not cause issues
				assert.NotPanics(t, func() {
					_ = v1State[tfconstants.AttrRegion]
					_ = v1State["floating_ip_id"]
					_ = v1State[tfconstants.AttrPublicIPAddress]
					_ = v1State[tfconstants.AttrPrivateIPAddress]
				})
			},
		},
		{
			name: "state with existing tags preserved",
			v0State: map[string]interface{}{
				"id":                   "lb-existing-tags",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode: goe2econstants.LBModeHTTP,
				tfconstants.AttrLbName: "test-lb",
				tfconstants.AttrTags: map[string]interface{}{
					"Environment": "production",
					"Team":        "platform",
				},
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Existing tags preserved
				tags := v1State[tfconstants.AttrTags].(map[string]interface{})
				assert.Equal(t, "production", tags["Environment"])
				assert.Equal(t, "platform", tags["Team"])
			},
		},
		{
			name: "state with both new and old field names",
			v0State: map[string]interface{}{
				"id":                             "lb-both-fields",
				tfconstants.AttrPlan:             goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode:           goe2econstants.LBModeHTTP,
				tfconstants.AttrName:             "new-name",
				tfconstants.AttrLbName:           "old-name",
				tfconstants.AttrRegion:           "Chennai",
				tfconstants.AttrLocation:         "Mumbai",
				"floating_ip_id":                 "ip-new",
				"lb_reserve_ip":                  "ip-old",
				tfconstants.AttrPublicIPAddress:  "203.0.113.1",
				"public_ip":                      "203.0.113.2",
				tfconstants.AttrPrivateIPAddress: "10.0.0.1",
				"private_ip":                     "10.0.0.2",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// When both present, new field values should not be overwritten
				// This tests the conditional logic in upgrade function
				assert.Equal(t, "new-name", v1State[tfconstants.AttrName])
				assert.Equal(t, "old-name", v1State[tfconstants.AttrLbName])
				assert.Equal(t, "Chennai", v1State[tfconstants.AttrRegion])
				assert.Equal(t, "Mumbai", v1State[tfconstants.AttrLocation])
				assert.Equal(t, "ip-new", v1State["floating_ip_id"])
				assert.Equal(t, "ip-old", v1State["lb_reserve_ip"])
			},
		},
		{
			name: "state with only new field names",
			v0State: map[string]interface{}{
				"id":                             "lb-new-fields-only",
				tfconstants.AttrPlan:             goe2econstants.LBPlanE2ELB2,
				tfconstants.AttrLbMode:           goe2econstants.LBModeHTTP,
				tfconstants.AttrName:             "modern-lb",
				tfconstants.AttrRegion:           "Delhi",
				"floating_ip_id":                 "ip-123",
				tfconstants.AttrPublicIPAddress:  "203.0.113.10",
				tfconstants.AttrPrivateIPAddress: "10.0.0.5",
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Old fields should not be populated when only new fields exist
				assert.Equal(t, "modern-lb", v1State[tfconstants.AttrName])
				assert.Equal(t, "Delhi", v1State[tfconstants.AttrRegion])
				assert.Equal(t, "ip-123", v1State["floating_ip_id"])
			},
		},
		{
			name: "complex backends and TCP backends preserved",
			v0State: map[string]interface{}{
				"id":                   "lb-complex-backends",
				tfconstants.AttrPlan:   goe2econstants.LBPlanE2ELB4,
				tfconstants.AttrLbMode: goe2econstants.LBModeBoth,
				tfconstants.AttrLbName: "complex-backend-lb",
				tfconstants.AttrBackends: []interface{}{
					map[string]interface{}{
						"name":        "backend-http",
						"balance":     goe2econstants.LBBalanceRoundRobin,
						"domain_name": goe2econstants.LBDefaultDomainName,
						"check_url":   goe2econstants.LBDefaultCheckURL,
						"http_check":  true,
						"servers": []interface{}{
							map[string]interface{}{
								"id":   "node-http-1",
								"port": goe2econstants.LBPortHTTP,
							},
							map[string]interface{}{
								"id":   "node-http-2",
								"port": goe2econstants.LBPortHTTP,
							},
						},
					},
					map[string]interface{}{
						"name":    "backend-https",
						"balance": goe2econstants.LBBalanceLeastConn,
						"servers": []interface{}{
							map[string]interface{}{
								"id":   "node-https-1",
								"port": goe2econstants.LBPortHTTPS,
							},
						},
					},
				},
				"tcp_backend": []interface{}{
					map[string]interface{}{
						"backend_name": "tcp-mysql",
						"port":         "3306",
						"balance":      goe2econstants.LBBalanceSource,
						"servers": []interface{}{
							map[string]interface{}{
								"id":   "node-mysql-1",
								"port": "3306",
							},
							map[string]interface{}{
								"id":   "node-mysql-2",
								"port": "3306",
							},
						},
					},
				},
			},
			validate: func(t *testing.T, v1State map[string]interface{}) {
				// Backends preserved
				backends := v1State[tfconstants.AttrBackends].([]interface{})
				assert.Len(t, backends, 2)
				backend1 := backends[0].(map[string]interface{})
				assert.Equal(t, "backend-http", backend1["name"])
				assert.Len(t, backend1["servers"].([]interface{}), 2)

				// TCP backends preserved
				tcpBackends := v1State["tcp_backend"].([]interface{})
				assert.Len(t, tcpBackends, 1)
				tcpBackend1 := tcpBackends[0].(map[string]interface{})
				assert.Equal(t, "tcp-mysql", tcpBackend1["backend_name"])
				assert.Len(t, tcpBackend1["servers"].([]interface{}), 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1State, err := loadbalancer.ResourceLoadBalancerStateUpgradeV0toV1(context.Background(), tt.v0State, nil)
			assert.NoError(t, err)
			tt.validate(t, v1State)
		})
	}
}
