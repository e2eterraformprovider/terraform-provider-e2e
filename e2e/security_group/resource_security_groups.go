package security_group

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Default values for security group configuration
const (
	defaultMyNetworkSize = 512
)

func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceSecurityGroupResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceSecurityGroupStateUpgradeV0toV1,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED FIELDS - IMMUTABLE
			// ============================================
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the Security Group",
			},

			// ============================================
			// REQUIRED FIELDS - MUTABLE
			// ============================================
			"rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "firewall rules defining inbound and outbound traffic (DEPRECATED: Use e2e_security_group_rule resource instead to avoid conflicts. This inline pattern will be removed in v4.0.)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "id of the firewall rule",
						},
						"rule_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "direction of traffic (Inbound or Outbound)",
							ValidateFunc: validation.StringInSlice([]string{
								"Inbound", "Outbound",
							}, false),
						},
						"network": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "any",
							Description: "network type for the rule (myNetwork, manual, or any)",
							ValidateFunc: validation.StringInSlice([]string{
								"myNetwork", "manual", "any",
							}, false),
						},
						"protocol_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "All",
							Description: "the protocol to allow (All, All_TCP, All_UDP, ICMP, Custom_TCP, Custom_UDP)",
							ValidateFunc: validation.StringInSlice([]string{
								"All", "All_TCP", "All_UDP", "ICMP", "Custom_TCP", "Custom_UDP",
							}, false),
						},
						"port_range": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "All",
							Description: "the port range to allow (e.g., '22', '80-443', or 'All')",
						},
						"network_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "--",
							Description: "the CIDR block for the rule (format: 'vpc_<vpc_id>' for VPC or IP address for manual)",
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "the network size if myNetwork or manual network is used",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "the description of the firewall rule",
						},
					},
				},
			},

			// ============================================
			// OPTIONAL FIELDS - MANAGEMENT
			// ============================================
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the description of the Security Group",
			},
			"default": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				Description:   "whether the Security Group is the default group",
				Deprecated:    "Use `is_default` instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"is_default"},
			},
			"is_default": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				Description:   "whether the Security Group is the default group",
				ConflictsWith: []string{"default"},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "tags for the Security Group (state-only until API support is available)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// ============================================
			// COMPUTED FIELDS
			// ============================================
			tfconstants.AttrID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "id of the Security Group",
			},
			"is_all_traffic_rule": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the Security Group allows all traffic",
			},
		},
		CreateContext: resourceCreateSecurityGroup,
		ReadContext:   resourceReadSecurityGroup,
		DeleteContext: resourceDeleteSecurityGroup,
		UpdateContext: resourceUpdateSecurityGroup,
		Importer: &schema.ResourceImporter{
			StateContext: customImportSecurityGroup,
		},
	}
}

// customImportSecurityGroup handles importing security groups with the format: project_id/region/sg_name
func customImportSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format for import: expected 'project_id/region/sg_name', got: %s", d.Id())
	}

	projectID := parts[0]
	region := parts[1]
	sgName := parts[2]

	// Set the basic fields
	if err := d.Set(tfconstants.AttrProjectID, projectID); err != nil {
		return nil, fmt.Errorf("error setting project_id: %w", err)
	}
	if err := d.Set(tfconstants.AttrRegion, region); err != nil {
		return nil, fmt.Errorf("error setting region: %w", err)
	}
	if err := d.Set(tfconstants.AttrName, sgName); err != nil {
		return nil, fmt.Errorf("error setting name: %w", err)
	}

	// Now trigger Read to populate remaining fields including the actual ID
	diags := resourceReadSecurityGroup(ctx, d, m)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading security group during import: %s", diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceReadSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	log.Println("[INFO] SECURITY GROUP READ STARTS")

	name := d.Get("name").(string)

	// Get security group by name (API has no GET by ID, must list and filter)
	sgList, _, err := goe2eClient.SecurityGroups.GetSecurityGroupList(ctx)
	if err != nil {
		return diag.Errorf("Error listing security groups while searching for %s: %s", name, err)
	}

	// Find security group with matching name
	var sg *goe2e.SecurityGroup
	for _, item := range sgList {
		if item.Name == name {
			sg = item
			break
		}
	}

	if sg == nil {
		// Security group not found - it may have been deleted outside Terraform
		log.Printf("[WARN] Security group %s not found, removing from state", name)
		d.SetId("")
		return diags
	}

	// Set ID
	d.SetId(sg.ID)

	// Set basic fields
	if err := d.Set("id", sg.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting id: %w", err))
	}
	if err := d.Set("description", sg.Description); err != nil {
		return diag.FromErr(fmt.Errorf("error setting description: %w", err))
	}

	// Handle both `default` and `is_default` fields for backwards compatibility
	if _, ok := d.GetOk("is_default"); ok {
		if err := d.Set("is_default", sg.IsDefault); err != nil {
			return diag.FromErr(fmt.Errorf("error setting is_default: %w", err))
		}
	} else {
		// Use `default` field if `is_default` not set
		if err := d.Set("default", sg.IsDefault); err != nil {
			return diag.FromErr(fmt.Errorf("error setting default: %w", err))
		}
	}

	// Set computed field
	if err := d.Set("is_all_traffic_rule", false); err != nil { // TODO: Get from API when available
		return diag.FromErr(fmt.Errorf("error setting is_all_traffic_rule: %w", err))
	}

	// Convert rules from models.Rule to terraform schema format
	ruleList := flattenRules(sg.Rules)
	if err := d.Set("rules", ruleList); err != nil {
		return diag.FromErr(fmt.Errorf("error setting rules: %w", err))
	}

	// Tags are state-only, no need to set from API
	// They persist in state automatically

	return diags
}

func resourceCreateSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()

	log.Println("[INFO] SECURITY GROUP CREATE STARTS")

	// Log deprecation warning if inline rules are used
	if rawRules := d.Get("rules").([]interface{}); len(rawRules) > 0 {
		log.Printf("[WARN] Using inline rules in e2e_security_group is deprecated. " +
			"Consider using the e2e_security_group_rule resource instead to avoid conflicts. " +
			"This pattern will be removed in v4.0.")
	}

	// Expand rules from terraform schema to models.Rule
	rules := expandRules(d.Get("rules").([]interface{}))

	// Determine which field to use for default flag (backwards compatibility)
	isDefault := false
	if v, ok := d.GetOk("is_default"); ok {
		isDefault = v.(bool)
	} else if v, ok := d.GetOk("default"); ok {
		isDefault = v.(bool)
	}

	payload := &goe2e.SecurityGroupCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Default:     isDefault,
		Rules:       rules,
	}

	sgName := d.Get("name").(string)

	// Create the security group
	sg, _, err := goe2eClient.SecurityGroups.CreateSecurityGroup(ctx, payload)
	if err != nil {
		return diag.Errorf("Error creating security group (name: %s): %s", sgName, err)
	}

	// Set the ID
	d.SetId(sg.ID)

	// Read back to populate all fields including rule IDs
	return resourceReadSecurityGroup(ctx, d, m)
}

func resourceDeleteSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()

	log.Println("[INFO] SECURITY GROUP DELETE STARTS")

	id := d.Id()
	sgName := d.Get("name").(string)

	// Pre-delete validation: check if it's a default security group
	isDefault := false
	if v, ok := d.GetOk("is_default"); ok {
		isDefault = v.(bool)
	} else if v, ok := d.GetOk("default"); ok {
		isDefault = v.(bool)
	}

	if isDefault {
		return diag.Errorf("Cannot delete default security group (ID: %s, name: %s). "+
			"Unset is_default or default field before deleting.", id, sgName)
	}

	// Delete the security group
	_, err := goe2eClient.SecurityGroups.DeleteSecurityGroup(ctx, id)
	if err != nil {
		return diag.Errorf("Error deleting security group (ID: %s, name: %s): %s", id, sgName, err)
	}

	d.SetId("")
	return nil
}

func resourceUpdateSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()

	// Defensive check: name should not change (ForceNew should prevent this, but double-check)
	if d.HasChange("name") {
		return diag.Errorf("cannot update name: this field is immutable and requires resource recreation")
	}

	id := d.Id()
	sgName := d.Get("name").(string)

	// Check if any updateable fields have changed
	if d.HasChange("rules") || d.HasChange("description") {
		log.Println("[INFO] SECURITY GROUP UPDATE STARTS")

		// Log deprecation warning if inline rules are being updated
		if d.HasChange("rules") {
			log.Printf("[WARN] Updating inline rules in e2e_security_group is deprecated. " +
				"Consider using the e2e_security_group_rule resource instead. " +
				"This pattern will be removed in v4.0.")
		}

		// Expand rules from terraform schema (preserving rule IDs)
		rules := expandRulesWithIDs(d.Get("rules").([]interface{}))

		payload := &goe2e.SecurityGroupUpdateRequest{
			Name:        sgName,
			Description: d.Get("description").(string),
			Rules:       rules,
		}

		_, _, err := goe2eClient.SecurityGroups.UpdateSecurityGroup(ctx, id, payload)
		if err != nil {
			return diag.Errorf("Error updating security group (ID: %s, name: %s): %s", id, sgName, err)
		}
	}

	// Handle default/is_default field changes
	if d.HasChange("default") || d.HasChange("is_default") {
		log.Println("[INFO] SECURITY GROUP MARK DEFAULT STARTS")

		// Determine if should be set as default
		isDefault := false
		if v, ok := d.GetOk("is_default"); ok {
			isDefault = v.(bool)
		} else if v, ok := d.GetOk("default"); ok {
			isDefault = v.(bool)
		}

		if isDefault {
			_, err := goe2eClient.SecurityGroups.MakeDefaultSecurityGroup(ctx, id)
			if err != nil {
				return diag.Errorf("Error setting security group (ID: %s, name: %s) as default: %s", id, sgName, err)
			}
		} else {
			// Note: API may not support unsetting default status
			log.Printf("[WARN] Cannot unset default status for security group (ID: %s, name: %s) via API", id, sgName)
		}
	}

	// Tags are state-only, no API call needed

	return resourceReadSecurityGroup(ctx, d, m)
}

// Helper functions for rule expansion and flattening

// flattenRules converts []goe2e.Rule to terraform schema format
func flattenRules(rules []goe2e.Rule) []map[string]interface{} {
	if len(rules) == 0 {
		return []map[string]interface{}{}
	}

	ruleList := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		networkSize := 0
		if rule.NetworkSize != nil {
			networkSize = *rule.NetworkSize
		}

		ruleList[i] = map[string]interface{}{
			"rule_id":       rule.ID,
			"rule_type":     rule.RuleType,
			"protocol_name": rule.ProtocolName,
			"port_range":    rule.PortRange,
			"network":       rule.Network,
			"network_cidr":  rule.NetworkCIDR,
			"description":   rule.Description,
			"size":          networkSize,
		}
	}

	return ruleList
}

// expandRules converts terraform schema to []goe2e.Rule (for create operations, no IDs)
func expandRules(rawRules []interface{}) []goe2e.Rule {
	if len(rawRules) == 0 {
		return []goe2e.Rule{}
	}

	rules := make([]goe2e.Rule, len(rawRules))
	for i, raw := range rawRules {
		ruleData := raw.(map[string]interface{})
		network := ruleData["network"].(string)

		var networkSizePtr *int
		if network == "myNetwork" {
			size := defaultMyNetworkSize
			networkSizePtr = &size
		} else if v, ok := ruleData["size"].(int); ok && v > 0 {
			networkSizePtr = &v
		}

		rules[i] = goe2e.Rule{
			RuleType:     ruleData["rule_type"].(string),
			ProtocolName: ruleData["protocol_name"].(string),
			PortRange:    ruleData["port_range"].(string),
			Network:      network,
			NetworkCIDR:  ruleData["network_cidr"].(string),
			NetworkSize:  networkSizePtr,
			Description:  ruleData["description"].(string),
		}
	}

	return rules
}

// expandRulesWithIDs converts terraform schema to []goe2e.Rule (for update operations, preserving IDs)
func expandRulesWithIDs(rawRules []interface{}) []goe2e.Rule {
	if len(rawRules) == 0 {
		return []goe2e.Rule{}
	}

	rules := make([]goe2e.Rule, len(rawRules))
	for i, raw := range rawRules {
		ruleData := raw.(map[string]interface{})
		network := ruleData["network"].(string)

		// Extract rule ID if present (for updates)
		ruleID := 0
		if v, ok := ruleData["rule_id"].(int); ok && v > 0 {
			ruleID = v
		}

		var networkSizePtr *int
		if network == "myNetwork" {
			size := defaultMyNetworkSize
			networkSizePtr = &size
		} else if v, ok := ruleData["size"].(int); ok && v > 0 {
			networkSizePtr = &v
		}

		rules[i] = goe2e.Rule{
			ID:           ruleID,
			RuleType:     ruleData["rule_type"].(string),
			ProtocolName: ruleData["protocol_name"].(string),
			PortRange:    ruleData["port_range"].(string),
			Network:      network,
			NetworkCIDR:  ruleData["network_cidr"].(string),
			NetworkSize:  networkSizePtr,
			Description:  ruleData["description"].(string),
		}
	}

	return rules
}

// resourceSecurityGroupResourceV0 returns the V0 schema for state migration
// This represents the schema before V3 changes (without is_default and tags)
func resourceSecurityGroupResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),
			tfconstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			tfconstants.AttrID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_all_traffic_rule": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// ResourceSecurityGroupStateUpgradeV0toV1 upgrades state from v0 to v1
// Exported for testing purposes
// Changes:
// - Renames "default" to "is_default" if user is using is_default
// - Adds "tags" field with empty map
// - Preserves all other V2 fields
func ResourceSecurityGroupStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V3 tags field with default empty map
	if _, exists := rawState["tags"]; !exists {
		rawState["tags"] = make(map[string]interface{})
	}

	// No need to rename "default" to "is_default" because both fields are supported
	// for backwards compatibility. Users can continue using "default" or migrate to "is_default"
	// at their convenience.

	// Preserve all existing V2 fields (name, description, rules, default, id, is_all_traffic_rule)
	// No automatic renames or transformations

	log.Printf("[INFO] Upgraded security group state from v0 to v1: %s", rawState["id"])
	return rawState, nil
}
