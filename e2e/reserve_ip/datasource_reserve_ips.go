package reserve_ip

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
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
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

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
						e2econstants.AttrIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the IP address",
						},
						"reserved_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the type of IP address that is reserved",
						},
						e2econstants.AttrVMID: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the VM to which the reserved IP is attached",
						},
						e2econstants.AttrCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the date when the IP was purchased",
						},
						e2econstants.AttrVMName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the VM to which the reserved IP is attached (if any)",
						},
						e2econstants.AttrStatus: {
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
		oi[e2econstants.AttrReserveID] = ip.ReserveID
		oi[e2econstants.AttrApplianceType] = ip.ApplianceType
		oi[e2econstants.AttrBoughtAt] = ip.BoughtAt
		oi[e2econstants.AttrIPAddress] = ip.IPAddress
		oi[e2econstants.AttrReservedType] = ip.ReservedType
		oi[e2econstants.AttrStatus] = ip.Status
		oi[e2econstants.AttrVMID] = ip.VMID
		oi[e2econstants.AttrVMName] = ip.VMName
		ois[i] = oi
	}

	return ois
}
