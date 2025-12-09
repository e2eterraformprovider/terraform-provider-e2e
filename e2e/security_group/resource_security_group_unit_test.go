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
			RuleType:     "Inbound",
			ProtocolName: "Custom_TCP",
			PortRange:    "22",
			Network:      "manual",
			NetworkCIDR:  "10.0.0.0/24",
			NetworkSize:  &networkSize,
			Description:  "SSH access",
		},
	}

	result := flattenRules(rules)

	assert.Len(t, result, 1)
	assert.Equal(t, 123, result[0]["rule_id"])
	assert.Equal(t, "Inbound", result[0]["rule_type"])
	assert.Equal(t, "Custom_TCP", result[0]["protocol_name"])
	assert.Equal(t, "22", result[0]["port_range"])
	assert.Equal(t, "manual", result[0]["network"])
	assert.Equal(t, "10.0.0.0/24", result[0]["network_cidr"])
	assert.Equal(t, 512, result[0]["size"])
	assert.Equal(t, "SSH access", result[0]["description"])
}

func TestFlattenRules_NilNetworkSize(t *testing.T) {
	rules := []goe2e.Rule{
		{
			ID:           456,
			RuleType:     "Outbound",
			ProtocolName: "All",
			PortRange:    "All",
			Network:      "any",
			NetworkCIDR:  "--",
			NetworkSize:  nil, // Nil network size
			Description:  "",
		},
	}

	result := flattenRules(rules)

	assert.Len(t, result, 1)
	assert.Equal(t, 0, result[0]["size"], "nil network size should be converted to 0")
}

func TestFlattenRules_MultipleRules(t *testing.T) {
	size256 := 256
	size512 := 512
	rules := []goe2e.Rule{
		{
			ID:           123,
			RuleType:     "Inbound",
			ProtocolName: "Custom_TCP",
			PortRange:    "22",
			Network:      "manual",
			NetworkCIDR:  "192.168.1.0/24",
			NetworkSize:  &size256,
			Description:  "SSH",
		},
		{
			ID:           456,
			RuleType:     "Inbound",
			ProtocolName: "Custom_TCP",
			PortRange:    "80",
			Network:      "myNetwork",
			NetworkCIDR:  "--",
			NetworkSize:  &size512,
			Description:  "HTTP",
		},
		{
			ID:           789,
			RuleType:     "Outbound",
			ProtocolName: "All",
			PortRange:    "All",
			Network:      "any",
			NetworkCIDR:  "--",
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
			"rule_type":     "Inbound",
			"protocol_name": "Custom_TCP",
			"port_range":    "22",
			"network":       "manual",
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "SSH access",
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, "Inbound", result[0].RuleType)
	assert.Equal(t, "Custom_TCP", result[0].ProtocolName)
	assert.Equal(t, "22", result[0].PortRange)
	assert.Equal(t, "manual", result[0].Network)
	assert.Equal(t, "10.0.0.0/24", result[0].NetworkCIDR)
	assert.NotNil(t, result[0].NetworkSize)
	assert.Equal(t, 256, *result[0].NetworkSize)
	assert.Equal(t, "SSH access", result[0].Description)
	assert.Equal(t, 0, result[0].ID, "expandRules should not set ID")
}

func TestExpandRules_MyNetworkDefaultSize(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     "Inbound",
			"protocol_name": "All",
			"port_range":    "All",
			"network":       "myNetwork",
			"network_cidr":  "--",
			"description":   "My network rule",
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, "myNetwork", result[0].Network)
	assert.NotNil(t, result[0].NetworkSize)
	assert.Equal(t, 512, *result[0].NetworkSize, "myNetwork should default to 512")
}

func TestExpandRules_AnyNetworkNoSize(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     "Outbound",
			"protocol_name": "All",
			"port_range":    "All",
			"network":       "any",
			"network_cidr":  "--",
			"description":   "Allow all outbound",
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, "any", result[0].Network)
	assert.Nil(t, result[0].NetworkSize, "any network should have nil NetworkSize")
}

func TestExpandRules_MultipleRules(t *testing.T) {
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     "Inbound",
			"protocol_name": "Custom_TCP",
			"port_range":    "22",
			"network":       "manual",
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "SSH",
		},
		map[string]interface{}{
			"rule_type":     "Inbound",
			"protocol_name": "Custom_TCP",
			"port_range":    "80",
			"network":       "myNetwork",
			"network_cidr":  "--",
			"description":   "HTTP",
		},
		map[string]interface{}{
			"rule_type":     "Outbound",
			"protocol_name": "All",
			"port_range":    "All",
			"network":       "any",
			"network_cidr":  "--",
			"description":   "All outbound",
		},
	}

	result := expandRules(rawRules)

	assert.Len(t, result, 3)
	assert.Equal(t, "Inbound", result[0].RuleType)
	assert.Equal(t, "Inbound", result[1].RuleType)
	assert.Equal(t, "Outbound", result[2].RuleType)
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
			"rule_type":     "Inbound",
			"protocol_name": "Custom_TCP",
			"port_range":    "22",
			"network":       "manual",
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "SSH access",
		},
	}

	result := expandRulesWithIDs(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, 123, result[0].ID, "expandRulesWithIDs should preserve rule ID")
	assert.Equal(t, "Inbound", result[0].RuleType)
	assert.Equal(t, "Custom_TCP", result[0].ProtocolName)
	assert.Equal(t, "22", result[0].PortRange)
}

func TestExpandRulesWithIDs_WithoutRuleIDs(t *testing.T) {
	// Test that new rules (without rule_id) get ID = 0
	rawRules := []interface{}{
		map[string]interface{}{
			"rule_type":     "Inbound",
			"protocol_name": "All",
			"port_range":    "All",
			"network":       "any",
			"network_cidr":  "--",
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
			"rule_type":     "Inbound",
			"protocol_name": "Custom_TCP",
			"port_range":    "22",
			"network":       "manual",
			"network_cidr":  "10.0.0.0/24",
			"size":          256,
			"description":   "Existing SSH rule",
		},
		map[string]interface{}{
			"rule_type":     "Inbound",
			"protocol_name": "Custom_TCP",
			"port_range":    "443",
			"network":       "any",
			"network_cidr":  "--",
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
			"rule_type":     "Inbound",
			"protocol_name": "All",
			"port_range":    "All",
			"network":       "myNetwork",
			"network_cidr":  "--",
			"description":   "My network rule",
		},
	}

	result := expandRulesWithIDs(rawRules)

	assert.Len(t, result, 1)
	assert.Equal(t, 456, result[0].ID)
	assert.Equal(t, "myNetwork", result[0].Network)
	assert.NotNil(t, result[0].NetworkSize)
	assert.Equal(t, 512, *result[0].NetworkSize, "myNetwork should default to 512")
}

func TestRoundTripConversion(t *testing.T) {
	// Test that flatten -> expand produces the same data (excluding IDs)
	size256 := 256
	originalRules := []goe2e.Rule{
		{
			ID:           123,
			RuleType:     "Inbound",
			ProtocolName: "Custom_TCP",
			PortRange:    "22",
			Network:      "manual",
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
