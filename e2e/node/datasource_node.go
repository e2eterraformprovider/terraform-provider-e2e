package node

import (
	"context"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadNode,
		Schema: map[string]*schema.Schema{
			// Common fields
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the Node",
			},
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the Node",
			},
			tfconstants.AttrLabel: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the label of the Node",
			},
			tfconstants.AttrPlan: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the plan of the Node",
			},
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the Node",
			},
			tfconstants.AttrMemory: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "memory of the Node in megabytes",
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Node instance",
			},
			tfconstants.AttrDisk: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the disk information of the Node",
			},
			"price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the price details of the Node",
			},
			tfconstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the Nodes public ipv4 address",
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the Nodes private ipv4 address",
			},
			"is_bitninja_license_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the Node has BitNinja license active",
			},
			tfconstants.AttrIsLocked: {
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

	nodeId := d.Get(tfconstants.AttrNodeID).(string)
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
	d.Set(tfconstants.AttrName, node.Name)
	d.Set(tfconstants.AttrLabel, node.Label)
	d.Set(tfconstants.AttrPlan, node.Plan)
	d.Set(tfconstants.AttrCreatedAt, node.CreatedAt)
	d.Set(tfconstants.AttrMemory, node.Memory)
	d.Set(tfconstants.AttrStatus, node.Status)
	d.Set(tfconstants.AttrDisk, node.Disk)
	d.Set("price", node.Price)
	d.Set(tfconstants.AttrIsLocked, node.IsLocked)
	d.Set(tfconstants.AttrPublicIPAddress, node.PublicIPAddress)
	d.Set(tfconstants.AttrPrivateIPAddress, node.PrivateIPAddress)
	d.Set("is_bitninja_license_active", node.BitNinjaLicenseActive)
	log.Printf("[INFO] NODE DATA SOURCE | d : %+v", d)

	return diags

}
