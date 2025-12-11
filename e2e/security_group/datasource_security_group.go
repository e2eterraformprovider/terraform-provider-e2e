package security_group

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityGroupRead,
		Schema: map[string]*schema.Schema{
			// Common fields
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the Security Group",
			},
			tfconstants.AttrID: {
				Type:     schema.TypeString,
				Computed: true,
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
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	log.Printf("[INFO] Reading security group with name: %s", name)

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
		return diag.Errorf("Security group with name %s not found", name)
	}

	// Set ID
	d.SetId(sg.ID)

	// Set fields
	d.Set("description", sg.Description)
	d.Set("default", sg.IsDefault)

	// Convert rules using helper function
	ruleList := flattenRules(sg.Rules)
	d.Set("rules", ruleList)

	return diags
}
