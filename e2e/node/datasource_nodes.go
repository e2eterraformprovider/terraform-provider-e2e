package node

import (
	"context"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadNodes,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Common fields
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"nodes_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of all Nodes in your account",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						e2econstants.AttrID: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the Node",
						},
						e2econstants.AttrName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Node",
						},
						e2econstants.AttrIsLocked: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether the Node is locked",
						},
						e2econstants.AttrStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state of the Node instance",
						},
						e2econstants.AttrPrivateIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the Nodes private ipv4 address",
						},
						e2econstants.AttrPublicIPAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the Nodes public ipv4 address",
						},
						"rescue_mode_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the rescue mode status of the Node",
						},
					},
				},
			},
		},
	}
}

func dataSourceReadNodes(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	cfg := m.(*config.Config)
	log.Printf("[INFO] Inside nodes data source ")

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
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	nodes, _, err := goe2eClient.Nodes.ListNodes(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Found %d nodes", len(nodes))
	log.Printf("[INFO] NODES DATA SOURCE | before setting")
	d.Set("nodes_list", flattenNodes(nodes))
	d.SetId("nodes_list")

	return diags
}

func flattenNodes(nodes []goe2e.Node) []interface{} {

	if nodes != nil {
		ois := make([]interface{}, len(nodes))

		for i, node := range nodes {
			oi := make(map[string]interface{})
			// Convert ID from string to float64 for compatibility
			if id, err := strconv.ParseFloat(node.ID, 64); err == nil {
				oi[e2econstants.AttrID] = id
			} else {
				oi[e2econstants.AttrID] = 0.0
			}
			oi[e2econstants.AttrName] = node.Name
			oi[e2econstants.AttrIsLocked] = node.IsLocked
			oi[e2econstants.AttrPrivateIPAddress] = node.PrivateIPAddress
			oi[e2econstants.AttrPublicIPAddress] = node.PublicIPAddress
			// RescueModeStatus is not available in goe2e.Node, set empty string
			oi["rescue_mode_status"] = ""
			oi[e2econstants.AttrStatus] = node.Status
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
