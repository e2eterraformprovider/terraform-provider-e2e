package security_group

import (
	"context"
	// "fmt"
	"log"
	// "math"
	// "regexp"

	// "strconv"
	//"strings"

	"github.com/devteametwoe/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/devteametwoe/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecurityGroups() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{

			"security_group_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the security group",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the security group",
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rules for the security group",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"network_size": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"security_group": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"rule_type": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"created_at": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"updated_at": {
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
									"is_active": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"deleted": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
		ReadContext: dataSourceReadSecurityGroups,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadSecurityGroups(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside images data source ")
	Response, err := apiClient.GetSecurityGroups()
	if err != nil {
		return diag.Errorf("error finding security groups")
	}
	d.Set("security_group_list", flattenSecurityGroups(&Response.Data))
	d.SetId("security_group_list")

	return diags
}

func flattenSecurityGroups(securityGroupList *[]models.SecurityGroup) []interface{} {

	if securityGroupList != nil {
		ois := make([]interface{}, len(*securityGroupList), len(*securityGroupList))

		for i, securityGroup := range *securityGroupList {
			oi := make(map[string]interface{})
			oi["name"] = securityGroup.Name
			oi["id"] = securityGroup.Id
			oi["description"] = securityGroup.Description
			oi["is_default"] = securityGroup.Is_default

			rls := make([]interface{}, len(securityGroup.Rules), len(securityGroup.Rules))
			for j, rule := range securityGroup.Rules {
				rl := make(map[string]interface{})
				rl["id"] = rule.Id
				rl["deleted"] = rule.Deleted
				rl["rule_type"] = rule.Rule_type
				rl["created_at"] = rule.Created_at
				rl["updated_at"] = rule.Updated_at
				rl["protocol_name"] = rule.Protocol_name
				rl["port_range"] = rule.Port_range
				rl["network"] = rule.Network
				rl["is_active"] = rule.Is_active
				rl["network_cidr"] = rule.Network_cidr
				rl["network_size"] = rule.Network_size
				rl["security_group"] = rule.Security_group
				rls[j] = rl
			}

			oi["rules"] = rls
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
