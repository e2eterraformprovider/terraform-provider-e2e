package security_group

import (
	"context"
	"fmt"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID for the security group",
				ForceNew:    true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
						"network": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "any",
							ValidateFunc: validation.StringInSlice([]string{
								"myNetwork", "manual", "any",
							}, false),
						},
						"rule_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Inbound", "Outbound",
							}, false),
						},
						"protocol_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "All",
							ValidateFunc: validation.StringInSlice([]string{
								"All", "All_TCP", "All_UDP", "ICMP", "Custom_TCP", "Custom_UDP",
							}, false),
						},
						"port_range": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "All",
						},
						"network_cidr": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "--",
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
		},
		CreateContext: resourceCreateSecurityGroup,
		ReadContext:   resourceReadSecurityGroup,
		DeleteContext: resourceDeleteSecurityGroup,
		UpdateContext: resourceUpdateSecurityGroup,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceReadSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Println("[INFO] SECURITY GROUP READ STARTS")

	name := d.Get("name").(string)
	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)

	sg, err := apiClient.GetSecurityGroup(name, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%v", sg["id"]))
	_ = d.Set("id", sg["id"])
	_ = d.Set("description", sg["description"])
	_ = d.Set("default", sg["is_default"])

	rulesRaw := sg["rules"].([]interface{})
	var ruleList []map[string]interface{}

	for _, r := range rulesRaw {
		rule := r.(map[string]interface{})

		ruleList = append(ruleList, map[string]interface{}{
			"rule_id":       rule["id"],
			"rule_type":     rule["rule_type"],
			"protocol_name": rule["protocol_name"],
			"port_range":    rule["port_range"],
			"network":       rule["network"],
			"network_cidr":  rule["network_cidr"],
			"description":   rule["description"],
			"size":          int(rule["network_size"].(float64)), // always present
		})
	}

	_ = d.Set("rules", ruleList)
	return diags
}

func resourceCreateSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Println("[INFO] SECURITY GROUP CREATE STARTS")

	rawRules := d.Get("rules").([]interface{})
	var rules []models.Rule

	for _, raw := range rawRules {
		ruleData := raw.(map[string]interface{})
		network := ruleData["network"].(string)

		var networkSizePtr *int
		if network == "myNetwork" {
			size := 512
			networkSizePtr = &size
		} else if v, ok := ruleData["size"].(int); ok && v > 0 {
			networkSizePtr = &v
		}

		rule := models.Rule{
			Rule_type:     ruleData["rule_type"].(string),
			Protocol_name: ruleData["protocol_name"].(string),
			Port_range:    ruleData["port_range"].(string),
			Network:       network,
			Network_cidr:  ruleData["network_cidr"].(string),
			Network_size:  networkSizePtr,
			Description:   ruleData["description"].(string),
		}
		rules = append(rules, rule)
	}

	payload := models.SecurityGroupCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Default:     d.Get("default").(bool),
		Rules:       rules,
	}

	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)

	if err := apiClient.CreateSecurityGroups(payload, project_id, location); err != nil {
		return diag.FromErr(err)
	}

	// Now retrieve the created SG from the API
	sg, err := apiClient.GetSecurityGroup(payload.Name, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%v", sg["id"]))
	_ = d.Set("id", sg["id"])

	rulesRaw := sg["rules"].([]interface{})
	var ruleList []map[string]interface{}

	for _, r := range rulesRaw {
		rule := r.(map[string]interface{})

		ruleList = append(ruleList, map[string]interface{}{
			"rule_id":       rule["id"],
			"rule_type":     rule["rule_type"],
			"protocol_name": rule["protocol_name"],
			"port_range":    rule["port_range"],
			"network":       rule["network"],
			"network_cidr":  rule["network_cidr"],
			"description":   rule["description"],
			"size":          int(rule["network_size"].(float64)), // guaranteed present
		})
	}

	_ = d.Set("rules", ruleList)

	return diags
}

func resourceDeleteSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Println("[INFO] SECURITY GROUP READ STARTS")

	id := d.Get("id").(string)
	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)

	err := apiClient.DeleteSecurityGroup(id, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceUpdateSecurityGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)

	if d.HasChange("rules") {
		log.Println("[INFO] SECURITY GROUP RULES UPDATE STARTS")

		rawRules := d.Get("rules").([]interface{})
		id := d.Get("id").(string)
		var rules []models.Rule

		for _, raw := range rawRules {
			ruleData := raw.(map[string]interface{})
			network := ruleData["network"].(string)

			var ruleIDPtr *int
			if v, ok := ruleData["rule_id"].(int); ok && v > 0 {
				ruleIDPtr = &v
			}

			var networkSizePtr *int
			if v, ok := ruleData["size"].(int); ok && v > 0 {
				networkSizePtr = &v
			}

			rule := models.Rule{
				Id:            ruleIDPtr,
				Rule_type:     ruleData["rule_type"].(string),
				Protocol_name: ruleData["protocol_name"].(string),
				Port_range:    ruleData["port_range"].(string),
				Network:       network,
				Network_cidr:  ruleData["network_cidr"].(string),
				Network_size:  networkSizePtr,
				Description:   ruleData["description"].(string),
			}
			rules = append(rules, rule)
		}

		payload := models.SecurityGroupUpdateRequest{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Rules:       rules,
		}

		project_id := d.Get("project_id").(string)
		location := d.Get("location").(string)

		if err := apiClient.UpdateSecurityGroups(payload, id, project_id, location); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("default") {
		project_id := d.Get("project_id").(string)
		location := d.Get("location").(string)
		id := d.Get("id")

		if err := apiClient.MakeDefaultSecurityGroup(id.(string), project_id, location); err != nil {
			return diag.FromErr(err)
		}

	}

	return resourceReadSecurityGroup(ctx, d, m)
}
