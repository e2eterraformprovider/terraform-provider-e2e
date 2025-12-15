package security_group

import (
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/stretchr/testify/assert"
)

func TestFlattenRules_Empty(t *testing.T) {
	rules := []goe2e.Rule{}
	result := flattenRules(rules)

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestFlattenRules_SingleRule(t *testing.T) {
	networkSize := 512
	rules := []goe2e.Rule{
		{
			ID:           123,
			RuleType:     RuleTypeInbound,
			ProtocolName: ProtocolCustomTCP,
			PortRange:    "22",
			Network:      NetworkTypeManual,
			NetworkCIDR:  "10.0.0.0/24",
			NetworkSize:  &networkSize,
			Description:  "SSH access",
		},
	}

	result := flattenRules(rules)

	assert.Len(t, result, 1)
	assert.Equal(t, 123, result[0]["rule_id"])
	assert.Equal(t, RuleTypeInbound, result[0]["rule_type"])
	assert.Equal(t, ProtocolCustomTCP, result[0]["protocol_name"])
	assert.Equal(t, "22", result[0]["port_range"])
	assert.Equal(t, NetworkTypeManual, result[0]["network"])
	assert.Equal(t, "10.0.0.0/24", result[0]["network_cidr"])
	assert.Equal(t, 512, result[0]["size"])
	assert.Equal(t, "SSH access", result[0]["description"])
}

func TestFlattenRules_NilNetworkSize(t *testing.T) {
	rules := []goe2e.Rule{
		{
			ID:           456,
			RuleType:     RuleTypeOutbound,
			ProtocolName: ProtocolAll,
			PortRange:    ProtocolAll,
			Network:      NetworkTypeAny,
			NetworkCIDR:  DefaultNetworkCIDR,
			NetworkSize:  nil, // Nil network size
			Description:  DefaultDescription,
		},
	}

	result := flattenRules(rules)

	assert.Len(t, result, 1)
	assert.Equal(t, 0, result[0]["size"], "nil network size should be converted to 0")
}

func TestFlattenRules_MultipleRules(t *testing.T) {
	size256 := 256
	size512 := DefaultMyNetworkSize
	rules := []goe2e.Rule{
		{
			ID:           123,
			RuleType:     RuleTypeInbound,
			ProtocolName: ProtocolCustomTCP,
			PortRange:    "22",
			Network:      NetworkTypeManual,
			NetworkCIDR:  "192.168.1.0/24",
			NetworkSize:  &size256,
			Description:  "SSH",
		},
		{
			ID:           456,
			RuleType:     RuleTypeInbound,
			ProtocolName: ProtocolCustomTCP,
			PortRange:    "80",
			Network:      NetworkTypeMyNetwork,
			NetworkCIDR:  DefaultNetworkCIDR,
			NetworkSize:  &size512,
			Description:  "HTTP",
		},
		{
			ID:           789,
			RuleType:     RuleTypeOutbound,
			ProtocolName: ProtocolAll,
			PortRange:    ProtocolAll,
			Network:      NetworkTypeAny,
			NetworkCIDR:  DefaultNetworkCIDR,
			NetworkSize:  nil,
			Description:  "All outbound",
		},
	}

	result := flattenRules(rules)

	assert.Len(t, result, 3)
	assert.Equal(t, 123, result[0]["rule_id"])
	assert.Equal(t, 456, result[1]["rule_id"])
	assert.Equal(t, 789, result[2]["rule_id"])
}

func TestExpandRules_Empty(t *testing.T) {
	rawRules := []interface{}{}
	result := expandRules(rawRules)

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestExpandRules_SingleRule(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "22",
			"network":       NetworkTypeManual,
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "SSH access",
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, RuleTypeInbound, result[0].RuleType)
	assert.Equal(t, ProtocolCustomTCP, result[0].ProtocolName)
	assert.Equal(t, "22", result[0].PortRange)
	assert.Equal(t, NetworkTypeManual, result[0].Network)
	assert.Equal(t, "10.0.0.0/24", result[0].NetworkCIDR)
	assert.NotNil(t, result[0].NetworkSize)
	assert.Equal(t, 256, *result[0].NetworkSize)
	assert.Equal(t, "SSH access", result[0].Description)
	assert.Equal(t, 0, result[0].ID, "expandRules should not set ID")
}

func TestExpandRules_MyNetworkDefaultSize(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeMyNetwork,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   "My network rule",
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, NetworkTypeMyNetwork, result[0].Network)
	assert.NotNil(t, result[0].NetworkSize)
	assert.Equal(t, DefaultMyNetworkSize, *result[0].NetworkSize, "myNetwork should default to DefaultMyNetworkSize")
}

func TestExpandRules_AnyNetworkNoSize(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeOutbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   "Allow all outbound",
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, NetworkTypeAny, result[0].Network)
	assert.Nil(t, result[0].NetworkSize, "any network should have nil NetworkSize")
}

func TestExpandRules_MultipleRules(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "22",
			"network":       NetworkTypeManual,
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "SSH",
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "80",
			"network":       NetworkTypeMyNetwork,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   "HTTP",
		},
		map[string]interface{}{
			"rule_type":     RuleTypeOutbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   "All outbound",
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 3)
	assert.Equal(t, RuleTypeInbound, result[0].RuleType)
	assert.Equal(t, RuleTypeInbound, result[1].RuleType)
	assert.Equal(t, RuleTypeOutbound, result[2].RuleType)
}

func TestExpandRulesWithIDs_Empty(t *testing.T) {
	rawRules := []interface{}{}
	result := expandRulesWithIDs(rawRules)

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestExpandRulesWithIDs_WithRuleIDs(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_id":       123,
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "22",
			"network":       NetworkTypeManual,
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "SSH access",
		},
	}

	result := expandRulesWithIDs(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, 123, result[0].ID, "expandRulesWithIDs should preserve rule ID")
	assert.Equal(t, RuleTypeInbound, result[0].RuleType)
	assert.Equal(t, ProtocolCustomTCP, result[0].ProtocolName)
	assert.Equal(t, "22", result[0].PortRange)
}

func TestExpandRulesWithIDs_WithoutRuleIDs(t *testing.T) {
	// Test that new rules (without rule_id) get ID = 0
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   "New rule",
		},
	}

	result := expandRulesWithIDs(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, 0, result[0].ID, "New rules should have ID = 0")
}

func TestExpandRulesWithIDs_MixedRules(t *testing.T) {
	// Test mix of existing rules (with IDs) and new rules (without IDs)
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_id":       123,
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "22",
			"network":       NetworkTypeManual,
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "Existing SSH rule",
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "443",
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   "New HTTPS rule",
		},
	}

	result := expandRulesWithIDs(rawRules)

	assert.Len(t, result, 2)
	assert.Equal(t, 123, result[0].ID, "First rule should preserve ID")
	assert.Equal(t, 0, result[1].ID, "Second rule should have ID = 0 (new rule)")
}

func TestExpandRulesWithIDs_MyNetworkDefaultSize(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_id":       456,
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeMyNetwork,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   "My network rule",
		},
	}

	result := expandRulesWithIDs(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, 456, result[0].ID)
	assert.Equal(t, NetworkTypeMyNetwork, result[0].Network)
	assert.NotNil(t, result[0].NetworkSize)
	assert.Equal(t, DefaultMyNetworkSize, *result[0].NetworkSize, "myNetwork should default to DefaultMyNetworkSize")
}

func TestRoundTripConversion(t *testing.T) {
	// Test that flatten -> expand produces the same data (excluding IDs)
	size256 := 256
	originalRules := []goe2e.Rule{
		{
			ID:           123,
			RuleType:     RuleTypeInbound,
			ProtocolName: ProtocolCustomTCP,
			PortRange:    "22",
			Network:      NetworkTypeManual,
			NetworkCIDR:  "10.0.0.0/24",
			NetworkSize:  &size256,
			Description:  "SSH",
		},
	}

	// Flatten to Terraform format
	flattened := flattenRules(originalRules)
	assert.Len(t, flattened, 1)

	// Convert to []interface{} for expand
	flattenedInterface := make([]interface{}, len(flattened))
	for i, v := range flattened {
		flattenedInterface[i] = v
	}

	// Expand back (with IDs)
	expanded := expandRulesWithIDs(flattenedInterface)
	assert.Len(t, expanded, 1)

	// Verify data is preserved
	assert.Equal(t, originalRules[0].ID, expanded[0].ID)
	assert.Equal(t, originalRules[0].RuleType, expanded[0].RuleType)
	assert.Equal(t, originalRules[0].ProtocolName, expanded[0].ProtocolName)
	assert.Equal(t, originalRules[0].PortRange, expanded[0].PortRange)
	assert.Equal(t, originalRules[0].Network, expanded[0].Network)
	assert.Equal(t, originalRules[0].NetworkCIDR, expanded[0].NetworkCIDR)
	assert.Equal(t, *originalRules[0].NetworkSize, *expanded[0].NetworkSize)
	assert.Equal(t, originalRules[0].Description, expanded[0].Description)
}

// ============================================================================
// Additional Edge Cases for expandRules
// ============================================================================

// TestExpandRules_AllProtocolTypes tests expansion with all protocol types
func TestExpandRules_AllProtocolTypes(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAllTCP,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAllUDP,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolICMP,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "22",
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomUDP,
			"port_range":    "53",
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 6)
	assert.Equal(t, ProtocolAll, result[0].ProtocolName)
	assert.Equal(t, ProtocolAllTCP, result[1].ProtocolName)
	assert.Equal(t, ProtocolAllUDP, result[2].ProtocolName)
	assert.Equal(t, ProtocolICMP, result[3].ProtocolName)
	assert.Equal(t, ProtocolCustomTCP, result[4].ProtocolName)
	assert.Equal(t, ProtocolCustomUDP, result[5].ProtocolName)
}

// TestExpandRules_AllNetworkTypes tests expansion with all network types
func TestExpandRules_AllNetworkTypes(t *testing.T) {
	size1024 := 1024
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeMyNetwork,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeManual,
			"network_cidr":  "10.0.0.0/24",
			"size":          1024,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 3)
	assert.Equal(t, NetworkTypeMyNetwork, result[0].Network)
	assert.NotNil(t, result[0].NetworkSize)
	assert.Equal(t, DefaultMyNetworkSize, *result[0].NetworkSize, "myNetwork should default to DefaultMyNetworkSize")

	assert.Equal(t, NetworkTypeManual, result[1].Network)
	assert.NotNil(t, result[1].NetworkSize)
	assert.Equal(t, size1024, *result[1].NetworkSize, "manual network should use provided size")

	assert.Equal(t, NetworkTypeAny, result[2].Network)
	assert.Nil(t, result[2].NetworkSize, "any network should have nil NetworkSize")
}

// TestExpandRules_EdgeCasePortRanges tests expansion with various port range formats
func TestExpandRules_EdgeCasePortRanges(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "22", // Single port
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "80-443", // Port range
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll, // All ports
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 3)
	assert.Equal(t, "22", result[0].PortRange, "Single port should be preserved")
	assert.Equal(t, "80-443", result[1].PortRange, "Port range should be preserved")
	assert.Equal(t, ProtocolAll, result[2].PortRange, "All ports should be preserved")
}

// TestExpandRules_EmptyDescription tests expansion with empty description
func TestExpandRules_EmptyDescription(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription, // Empty description
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, DefaultDescription, result[0].Description, "Empty description should be set to DefaultDescription")
}

// ============================================================================
// Additional Edge Cases for expandRulesWithIDs
// ============================================================================

// TestExpandRulesWithIDs_ZeroRuleID tests expansion with zero rule ID
func TestExpandRulesWithIDs_ZeroRuleID(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_id":       0, // Zero rule ID
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolAll,
			"port_range":    ProtocolAll,
			"network":       NetworkTypeAny,
			"network_cidr":  DefaultNetworkCIDR,
			"description":   DefaultDescription,
		},
	}

	result := expandRulesWithIDs(rawRules)

	assert.Len(t, result, 1)
	// Zero ID should be treated as new rule (ID = 0)
	assert.Equal(t, 0, result[0].ID, "Zero rule ID should be treated as 0 (new rule)")
}

// ============================================================================
// Additional Edge Cases for flattenRules
// ============================================================================

// TestFlattenRules_AllProtocolTypes tests flattening with all protocol types
func TestFlattenRules_AllProtocolTypes(t *testing.T) {
	rules := []goe2e.Rule{
		{ID: 1, RuleType: RuleTypeInbound, ProtocolName: ProtocolAll, PortRange: ProtocolAll, Network: NetworkTypeAny, NetworkCIDR: DefaultNetworkCIDR, Description: DefaultDescription},
		{ID: 2, RuleType: RuleTypeInbound, ProtocolName: ProtocolAllTCP, PortRange: ProtocolAll, Network: NetworkTypeAny, NetworkCIDR: DefaultNetworkCIDR, Description: DefaultDescription},
		{ID: 3, RuleType: RuleTypeInbound, ProtocolName: ProtocolAllUDP, PortRange: ProtocolAll, Network: NetworkTypeAny, NetworkCIDR: DefaultNetworkCIDR, Description: DefaultDescription},
		{ID: 4, RuleType: RuleTypeInbound, ProtocolName: ProtocolICMP, PortRange: ProtocolAll, Network: NetworkTypeAny, NetworkCIDR: DefaultNetworkCIDR, Description: DefaultDescription},
		{ID: 5, RuleType: RuleTypeInbound, ProtocolName: ProtocolCustomTCP, PortRange: "22", Network: NetworkTypeAny, NetworkCIDR: DefaultNetworkCIDR, Description: DefaultDescription},
		{ID: 6, RuleType: RuleTypeInbound, ProtocolName: ProtocolCustomUDP, PortRange: "53", Network: NetworkTypeAny, NetworkCIDR: DefaultNetworkCIDR, Description: DefaultDescription},
	}

	result := flattenRules(rules)

	assert.Len(t, result, 6)
	assert.Equal(t, ProtocolAll, result[0]["protocol_name"])
	assert.Equal(t, ProtocolAllTCP, result[1]["protocol_name"])
	assert.Equal(t, ProtocolAllUDP, result[2]["protocol_name"])
	assert.Equal(t, ProtocolICMP, result[3]["protocol_name"])
	assert.Equal(t, ProtocolCustomTCP, result[4]["protocol_name"])
	assert.Equal(t, ProtocolCustomUDP, result[5]["protocol_name"])
}

// TestFlattenRules_AllNetworkTypes tests flattening with all network types
func TestFlattenRules_AllNetworkTypes(t *testing.T) {
	size512 := DefaultMyNetworkSize
	size1024 := 1024
	rules := []goe2e.Rule{
		{ID: 1, RuleType: RuleTypeInbound, ProtocolName: ProtocolAll, PortRange: ProtocolAll, Network: NetworkTypeMyNetwork, NetworkCIDR: DefaultNetworkCIDR, NetworkSize: &size512, Description: DefaultDescription},
		{ID: 2, RuleType: RuleTypeInbound, ProtocolName: ProtocolAll, PortRange: ProtocolAll, Network: NetworkTypeManual, NetworkCIDR: "10.0.0.0/24", NetworkSize: &size1024, Description: DefaultDescription},
		{ID: 3, RuleType: RuleTypeInbound, ProtocolName: ProtocolAll, PortRange: ProtocolAll, Network: NetworkTypeAny, NetworkCIDR: DefaultNetworkCIDR, NetworkSize: nil, Description: DefaultDescription},
	}

	result := flattenRules(rules)

	assert.Len(t, result, 3)
	assert.Equal(t, NetworkTypeMyNetwork, result[0]["network"])
	assert.Equal(t, DefaultMyNetworkSize, result[0]["size"], "myNetwork size should be preserved")

	assert.Equal(t, NetworkTypeManual, result[1]["network"])
	assert.Equal(t, 1024, result[1]["size"], "manual network size should be preserved")

	assert.Equal(t, NetworkTypeAny, result[2]["network"])
	assert.Equal(t, 0, result[2]["size"], "any network should have size 0")
}

// TestFlattenRules_VeryLargeNetworkSize tests flattening with very large network size
func TestFlattenRules_VeryLargeNetworkSize(t *testing.T) {
	largeSize := 2048 // Larger than default DefaultMyNetworkSize
	rules := []goe2e.Rule{
		{
			ID:           123,
			RuleType:     RuleTypeInbound,
			ProtocolName: ProtocolAll,
			PortRange:    ProtocolAll,
			Network:      NetworkTypeManual,
			NetworkCIDR:  "10.0.0.0/24",
			NetworkSize:  &largeSize,
			Description:  "Large network",
		},
	}

	result := flattenRules(rules)

	assert.Len(t, result, 1)
	assert.Equal(t, largeSize, result[0]["size"], "Very large network size should be preserved")
}

// ============================================================================
// Additional Round-Trip Conversion Tests
// ============================================================================

// TestRoundTripConversion_ExpandThenFlatten tests expand -> flatten round-trip
func TestRoundTripConversion_ExpandThenFlatten(t *testing.T) {
	// Start with Terraform schema format
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     RuleTypeInbound,
			"protocol_name": ProtocolCustomTCP,
			"port_range":    "22",
			"network":       NetworkTypeManual,
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "SSH access",
		},
	}

	// Expand to goe2e format
	expanded := expandRules(rawRules)
	assert.Len(t, expanded, 1)

	// Flatten back to Terraform format
	flattened := flattenRules(expanded)
	assert.Len(t, flattened, 1)

	// Verify data is preserved (with expected transformations)
	assert.Equal(t, RuleTypeInbound, flattened[0]["rule_type"])
	assert.Equal(t, ProtocolCustomTCP, flattened[0]["protocol_name"])
	assert.Equal(t, "22", flattened[0]["port_range"])
	assert.Equal(t, NetworkTypeManual, flattened[0]["network"])
	assert.Equal(t, "10.0.0.0/24", flattened[0]["network_cidr"])
	assert.Equal(t, 256, flattened[0]["size"])
	assert.Equal(t, "SSH access", flattened[0]["description"])
	// Note: ID will be 0 since expandRules doesn't set IDs
	assert.Equal(t, 0, flattened[0]["rule_id"])
}

// TestRoundTripConversion_FlattenThenExpand tests flatten -> expand round-trip
func TestRoundTripConversion_FlattenThenExpand(t *testing.T) {
	size256 := 256
	// Start with goe2e Rule structs
	originalRules := []goe2e.Rule{
		{
			ID:           123,
			RuleType:     RuleTypeInbound,
			ProtocolName: ProtocolCustomTCP,
			PortRange:    "22",
			Network:      NetworkTypeManual,
			NetworkCIDR:  "10.0.0.0/24",
			NetworkSize:  &size256,
			Description:  "SSH",
		},
	}

	// Flatten to Terraform schema format
	flattened := flattenRules(originalRules)
	assert.Len(t, flattened, 1)

	// Convert to []interface{} for expand
	flattenedInterface := make([]interface{}, len(flattened))
	for i, v := range flattened {
		flattenedInterface[i] = v
	}

	// Expand back to goe2e format (with IDs)
	expanded := expandRulesWithIDs(flattenedInterface)
	assert.Len(t, expanded, 1)

	// Verify data is preserved (with expected transformations)
	assert.Equal(t, originalRules[0].ID, expanded[0].ID, "ID should be preserved")
	assert.Equal(t, originalRules[0].RuleType, expanded[0].RuleType)
	assert.Equal(t, originalRules[0].ProtocolName, expanded[0].ProtocolName)
	assert.Equal(t, originalRules[0].PortRange, expanded[0].PortRange)
	assert.Equal(t, originalRules[0].Network, expanded[0].Network)
	assert.Equal(t, originalRules[0].NetworkCIDR, expanded[0].NetworkCIDR)
	assert.Equal(t, *originalRules[0].NetworkSize, *expanded[0].NetworkSize)
	assert.Equal(t, originalRules[0].Description, expanded[0].Description)
}
