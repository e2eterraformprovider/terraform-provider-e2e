package reserve_ip

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceReserveIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadReserveIps,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Common fields
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"reserve_ips_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of all reserved IPs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reserve_id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the reserved IP",
						},
						"appliance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the type of infrastructure to which the reserved IP is attached",
						},
						tfconstants.AttrIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the IP address",
						},
						"reserved_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the type of IP address that is reserved",
						},
						tfconstants.AttrVMID: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the VM to which the reserved IP is attached",
						},
						tfconstants.AttrCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the date when the IP was purchased",
						},
						tfconstants.AttrVMName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the VM to which the reserved IP is attached (if any)",
						},
						tfconstants.AttrStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "whether the IP is attached or available",
						},
					},
				},
			},
		},
	}
}

func dataSourceReadReserveIps(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	cfg := m.(*config.Config)
	log.Printf("[INFO] Inside reserve ips data source ")

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client: %s", err)
	}

	reserveIPs, _, err := goe2eClient.ReserveIP.ListReserveIPs(ctx)
	if err != nil {
		return diag.Errorf("error listing reserved IPs: %s", err)
	}

	d.Set("reserve_ips_list", flattenReserveIps(reserveIPs))
	d.SetId("reserve_ips_list")
	var diags diag.Diagnostics
	return diags
}
func flattenReserveIps(ips []goe2e.ReserveIP) []interface{} {
	if len(ips) == 0 {
		return make([]interface{}, 0)
	}

	ois := make([]interface{}, len(ips))
	for i, ip := range ips {
		oi := make(map[string]interface{})
		// Map goe2e model fields to Terraform schema
		// API field "bought_at" maps to schema field "created_at"
		oi[tfconstants.AttrReserveID] = ip.ReserveID
		oi[tfconstants.AttrApplianceType] = ip.ApplianceType
		oi[tfconstants.AttrBoughtAt] = ip.BoughtAt
		oi[tfconstants.AttrIPAddress] = ip.IPAddress
		oi[tfconstants.AttrReservedType] = ip.ReservedType
		oi[tfconstants.AttrStatus] = ip.Status
		oi[tfconstants.AttrVMID] = ip.VMID
		oi[tfconstants.AttrVMName] = ip.VMName
		ois[i] = oi
	}

	return ois
}
