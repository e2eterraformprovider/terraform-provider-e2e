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
			tfconstants.AttrVPCList: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of all VPCs (you can attach these VPCs to launch resources)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tfconstants.AttrNetworkID: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the VPC network",
						},
						tfconstants.AttrPoolSize: {
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
						tfconstants.AttrIPv4CIDR: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the IPv4 CIDR block of the VPC",
						},
						tfconstants.AttrGatewayIP: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the gateway IP address of the VPC",
						},
						tfconstants.AttrIsActive: {
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
		d.Set(tfconstants.AttrVPCList, flattenVpcs(vpcs))
		d.SetId(tfconstants.AttrVPCList)
	} else {
		log.Printf("[ERROR] %s", errVPCListEmpty)
	}
	return diags
}

func flattenVpcs(vpcList []goe2e.Vpc) []interface{} {

	if len(vpcList) > 0 {
		ois := make([]interface{}, len(vpcList))

		for i, vpc := range vpcList {
			oi := make(map[string]interface{})
			oi[tfconstants.AttrNetworkID] = vpc.ID
			oi[tfconstants.AttrPoolSize] = vpc.PoolSize
			oi[tfconstants.AttrCreatedAt] = vpc.CreatedAt
			oi[tfconstants.AttrName] = vpc.Name
			oi[tfconstants.AttrIsActive] = vpc.IsActive
			oi[tfconstants.AttrGatewayIP] = vpc.GatewayIP
			oi[tfconstants.AttrIPv4CIDR] = vpc.IPv4CIDR
			oi[tfconstants.AttrStatus] = vpc.State
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
