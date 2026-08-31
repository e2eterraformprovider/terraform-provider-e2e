package security_group

import (
	"context"
	"fmt"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the security group",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)

	log.Printf("[INFO] Reading security group with name: %s", name)

	sg, err := apiClient.GetSecurityGroup(name, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%v", sg["id"]))

	_ = d.Set("project_id", sg["project_id"])
	_ = d.Set("location", sg["location"])
	_ = d.Set("description", sg["description"])
	_ = d.Set("default", sg["is_default"])

	if rulesRaw, ok := sg["rules"].([]interface{}); ok {
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
				"size":          int(rule["network_size"].(float64)),
			})
		}
		_ = d.Set("rules", ruleList)
	}

	return diags
}
