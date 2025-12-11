package vpc

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVpcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadVpcs,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Common fields
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"vpc_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of all VPCs (you can attach these VPCs to launch resources)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the VPC network",
						},
						"pool_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "the pool size of the VPC",
						},
						tfconstants.AttrCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the creation date for the VPC",
						},
						tfconstants.AttrStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status of the VPC instance",
						},
						tfconstants.AttrName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the VPC",
						},
						"ipv4_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the IPv4 CIDR block of the VPC",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the gateway IP address of the VPC",
						},
						"is_active": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether the VPC is active",
						},
					},
				},
			},
		},
	}
}

func dataSourceReadVpcs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	log.Printf("[INFO] Inside vpcs data source ")

	vpcs, _, err := goe2eClient.Vpcs.ListVPCs(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] %v", vpcs)
	if len(vpcs) > 0 {
		d.Set("vpc_list", flattenVpcs(vpcs))
		d.SetId("vpc_list")
	} else {
		log.Printf("[ERROR] VPC list is empty in the response")
	}
	return diags
}

func flattenVpcs(vpcList []goe2e.Vpc) []interface{} {

	if len(vpcList) > 0 {
		ois := make([]interface{}, len(vpcList))

		for i, vpc := range vpcList {
			oi := make(map[string]interface{})
			oi["network_id"] = vpc.ID
			oi["pool_size"] = vpc.PoolSize
			oi["created_at"] = vpc.CreatedAt
			oi["name"] = vpc.Name
			oi["is_active"] = vpc.IsActive
			oi["gateway_ip"] = vpc.GatewayIP
			oi["ipv4_cidr"] = vpc.IPv4CIDR
			oi[tfconstants.AttrStatus] = vpc.State
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
