package sfs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSfs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadSfs,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Common fields
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"sfs_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of all the SFS instances in the account",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						e2econstants.AttrID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the id of the SFS",
						},
						e2econstants.AttrName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the name of the SFS",
						},
						e2econstants.AttrSizeGB: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the size of the SFS volume in gigabytes",
						},
						e2econstants.AttrStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the API status of the SFS instance",
						},
						e2econstants.AttrState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the normalized state of the SFS instance",
						},
						"private_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the NFS mount endpoint for the SFS",
						},
						e2econstants.AttrPlan: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the plan of the SFS",
						},
						"is_backup_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether backups are enabled for the SFS",
						},
						e2econstants.AttrIOPS: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the IOPS value of the SFS",
						},
						e2econstants.AttrVPCID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the id of the VPC for the SFS",
						},
						e2econstants.AttrEncryptionEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether encryption is enabled for the SFS",
						},
					},
				},
			},
		},
	}
}

func dataSourceReadSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cfg := m.(*config.Config)

	log.Printf("[INFO] Reading SFS data source")

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

	// Create goe2e client with specific projectID and region
	client, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	// List all SFS instances
	sfsList, _, err := client.Sfs.ListSfss(ctx)
	if err != nil {
		return diag.Errorf("Error listing SFS instances in project (%s), region (%s): %s", projectID, region, err)
	}

	log.Printf("[INFO] SFS DATA SOURCE | Retrieved %d SFS instances", len(sfsList))

	if err := d.Set("sfs_list", flattenSfsList(sfsList)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sfs_list: %w", err))
	}

	d.SetId("sfs_list")

	return diags
}

// flattenSfsList converts a list of goe2e.Sfs to terraform data
func flattenSfsList(sfsList []goe2e.Sfs) []interface{} {
	if len(sfsList) == 0 {
		return make([]interface{}, 0)
	}

	ois := make([]interface{}, len(sfsList))

	for i, sfs := range sfsList {
		oi := make(map[string]interface{})
		oi[e2econstants.AttrID] = sfs.ID
		oi[e2econstants.AttrName] = sfs.Name
		oi[e2econstants.AttrStatus] = sfs.Status
		oi[e2econstants.AttrState] = normalizeSfsState(sfs.Status)
		oi[e2econstants.AttrPlan] = sfs.PlanName
		oi["private_endpoint"] = sfs.PrivateIPAddress
		oi["is_backup_enabled"] = sfs.IsBackupEnabled
		oi[e2econstants.AttrIOPS] = sfs.DiskIOPS
		oi[e2econstants.AttrVPCID] = sfs.VPCID
		oi[e2econstants.AttrEncryptionEnabled] = sfs.IsEncryptionEnabled

		// Parse disk size from string (e.g., "100GB")
		if sfs.DiskSize != "" {
			diskSizeStr := strings.TrimSpace(strings.ReplaceAll(sfs.DiskSize, "GB", ""))
			if sizeInt, err := strconv.Atoi(diskSizeStr); err == nil {
				oi[e2econstants.AttrSizeGB] = sizeInt
			}
		}

		ois[i] = oi
	}

	return ois
}
