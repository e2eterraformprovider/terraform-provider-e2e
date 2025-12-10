package node

import (
	"context"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadNode,
		Schema: map[string]*schema.Schema{
			// Common fields
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the Node",
			},
			e2econstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the Node",
			},
			e2econstants.AttrLabel: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the label of the Node",
			},
			e2econstants.AttrPlan: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the plan of the Node",
			},
			e2econstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the Node",
			},
			e2econstants.AttrMemory: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "memory of the Node in megabytes",
			},
			e2econstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Node instance",
			},
			e2econstants.AttrDisk: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the disk information of the Node",
			},
			"price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the price details of the Node",
			},
			e2econstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the Nodes public ipv4 address",
			},
			e2econstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the Nodes private ipv4 address",
			},
			"is_bitninja_license_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the Node has BitNinja license active",
			},
			e2econstants.AttrIsLocked: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the Node has been locked",
			},
		},
	}
}
func dataSourceReadNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	log.Printf("[INFO] INSIDE NODE DATA SOURCE | read")

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	project_id, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(project_id, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	nodeId := d.Get(e2econstants.AttrNodeID).(string)
	node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return diag.Errorf("error finding Item with ID %s", nodeId)
		}
		return diag.Errorf("error finding Item with ID %s: %s", nodeId, err)
	}

	if node == nil {
		return diag.Errorf("error finding Item with ID %s: node not found", nodeId)
	}

	d.SetId(nodeId)
	log.Printf("[INFO] NODE DATA SOURCE | READ | data : %+v", node)
	d.Set(e2econstants.AttrName, node.Name)
	d.Set(e2econstants.AttrLabel, node.Label)
	d.Set(e2econstants.AttrPlan, node.Plan)
	d.Set(e2econstants.AttrCreatedAt, node.CreatedAt)
	d.Set(e2econstants.AttrMemory, node.Memory)
	d.Set(e2econstants.AttrStatus, node.Status)
	d.Set(e2econstants.AttrDisk, node.Disk)
	d.Set("price", node.Price)
	d.Set(e2econstants.AttrIsLocked, node.IsLocked)
	d.Set(e2econstants.AttrPublicIPAddress, node.PublicIPAddress)
	d.Set(e2econstants.AttrPrivateIPAddress, node.PrivateIPAddress)
	d.Set("is_bitninja_license_active", node.BitNinjaLicenseActive)
	log.Printf("[INFO] NODE DATA SOURCE | d : %+v", d)

	return diags

}
