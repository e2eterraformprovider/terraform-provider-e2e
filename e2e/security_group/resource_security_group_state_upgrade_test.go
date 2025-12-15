package security_group_test

import (
	"context"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/security_group"
	"github.com/stretchr/testify/assert"
)

func TestResourceSecurityGroupStateUpgradeV0toV1_Basic(t *testing.T) {
	v0State := map[string]interface{}{
		"id":          "12345",
		"name":        "test-sg",
		"description": "Test security group",
		"default":     false,
		"region":      "us-east-1",
		"project_id":  "789",
		"rules": []interface{}{
			map[string]interface{}{
				"rule_id":       123,
				"rule_type":     security_group.RuleTypeInbound,
				"protocol_name": security_group.ProtocolAll,
				"port_range":    security_group.ProtocolAll,
				"network":       security_group.NetworkTypeAny,
				"network_cidr":  security_group.DefaultNetworkCIDR,
				"description":   security_group.DefaultDescription,
			},
		},
		"is_all_traffic_rule": false,
	}

	v1State, err := security_group.ResourceSecurityGroupStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.Equal(t, "12345", v1State["id"])
	assert.Equal(t, "test-sg", v1State["name"])
	assert.Equal(t, "Test security group", v1State["description"])
	assert.Equal(t, false, v1State["default"])
	assert.NotNil(t, v1State["tags"], "tags field should be added")
	tags, ok := v1State["tags"].(map[string]interface{})
	assert.True(t, ok, "tags should be a map")
	assert.Empty(t, tags, "tags should be empty map")
}

func TestResourceSecurityGroupStateUpgradeV0toV1_WithExistingTags(t *testing.T) {
	v0State := map[string]interface{}{
		"id":          "12345",
		"name":        "test-sg",
		"description": "Test security group",
		"default":     false,
		"region":      "us-east-1",
		"project_id":  "789",
		"tags": map[string]interface{}{
			"Environment": "production",
			"Team":        "platform",
		},
		"rules": []interface{}{
			map[string]interface{}{
				"rule_id":       123,
				"rule_type":     security_group.RuleTypeInbound,
				"protocol_name": security_group.ProtocolAll,
				"network":       security_group.NetworkTypeAny,
			},
		},
	}

	v1State, err := security_group.ResourceSecurityGroupStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.NotNil(t, v1State["tags"])
	// Existing tags should be preserved
	tags := v1State["tags"].(map[string]interface{})
	assert.Equal(t, "production", tags["Environment"])
	assert.Equal(t, "platform", tags["Team"])
}

func TestResourceSecurityGroupStateUpgradeV0toV1_PreservesAllFields(t *testing.T) {
	v0State := map[string]interface{}{
		"id":          "67890",
		"name":        "complex-sg",
		"description": "Complex security group",
		"default":     true,
		"region":      "us-west-2",
		"location":    "us-west-2",
		"project_id":  "123",
		"rules": []interface{}{
			map[string]interface{}{
				"rule_id":       456,
				"rule_type":     security_group.RuleTypeInbound,
				"protocol_name": security_group.ProtocolCustomTCP,
				"port_range":    "22",
				"network":       security_group.NetworkTypeManual,
				"network_cidr":  "10.0.0.0/24",
				"size":          256,
				"description":   "SSH access",
			},
			map[string]interface{}{
				"rule_id":       789,
				"rule_type":     security_group.RuleTypeOutbound,
				"protocol_name": security_group.ProtocolAll,
				"port_range":    security_group.ProtocolAll,
				"network":       security_group.NetworkTypeAny,
				"network_cidr":  security_group.DefaultNetworkCIDR,
				"description":   "Allow all outbound",
			},
		},
		"is_all_traffic_rule": false,
	}

	v1State, err := security_group.ResourceSecurityGroupStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// All V0 fields should be preserved
	assert.Equal(t, "67890", v1State["id"])
	assert.Equal(t, "complex-sg", v1State["name"])
	assert.Equal(t, "Complex security group", v1State["description"])
	assert.Equal(t, true, v1State["default"])
	assert.Equal(t, "us-west-2", v1State["region"])
	assert.Equal(t, "us-west-2", v1State["location"])
	assert.Equal(t, "123", v1State["project_id"])
	assert.Equal(t, false, v1State["is_all_traffic_rule"])

	// Rules should be preserved
	rules := v1State["rules"].([]interface{})
	assert.Equal(t, 2, len(rules))

	// Check first rule
	rule1 := rules[0].(map[string]interface{})
	assert.Equal(t, 456, rule1["rule_id"])
	assert.Equal(t, security_group.RuleTypeInbound, rule1["rule_type"])
	assert.Equal(t, security_group.ProtocolCustomTCP, rule1["protocol_name"])
	assert.Equal(t, "22", rule1["port_range"])
	assert.Equal(t, "SSH access", rule1["description"])

	// Tags should be added
	assert.NotNil(t, v1State["tags"])
}

func TestResourceSecurityGroupStateUpgradeV0toV1_DefaultFieldPreserved(t *testing.T) {
	// Test that the "default" field is preserved (not renamed to "is_default")
	// This ensures backwards compatibility
	v0State := map[string]interface{}{
		"id":          "99999",
		"name":        "default-test-sg",
		"description": "Test default field preservation",
		"default":     true, // Using deprecated "default" field
		"region":      "us-east-1",
		"project_id":  "456",
		"rules": []interface{}{
			map[string]interface{}{
				"rule_id":       111,
				"rule_type":     security_group.RuleTypeInbound,
				"protocol_name": security_group.ProtocolAll,
				"network":       security_group.NetworkTypeAny,
			},
		},
	}

	v1State, err := security_group.ResourceSecurityGroupStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	// The "default" field should still exist (not renamed)
	assert.Equal(t, true, v1State["default"])
	// "is_default" should NOT be automatically added
	assert.Nil(t, v1State["is_default"])
	// Tags should be added
	assert.NotNil(t, v1State["tags"])
}

func TestResourceSecurityGroupStateUpgradeV0toV1_EmptyRules(t *testing.T) {
	v0State := map[string]interface{}{
		"id":          "55555",
		"name":        "empty-rules-sg",
		"description": "Security group with no rules",
		"default":     false,
		"region":      "us-east-1",
		"project_id":  "789",
		"rules":       []interface{}{}, // Empty rules
	}

	v1State, err := security_group.ResourceSecurityGroupStateUpgradeV0toV1(context.Background(), v0State, nil)

	assert.NoError(t, err)
	assert.Equal(t, "55555", v1State["id"])
	assert.Equal(t, "empty-rules-sg", v1State["name"])
	assert.NotNil(t, v1State["tags"])
	assert.NotNil(t, v1State["rules"])
	rules := v1State["rules"].([]interface{})
	assert.Empty(t, rules, "rules should remain empty")
}
