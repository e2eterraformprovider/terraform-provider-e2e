package sfs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
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
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"sfs_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of all the SFS instances in the account",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tfconstants.AttrID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the id of the SFS",
						},
						tfconstants.AttrName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the name of the SFS",
						},
						tfconstants.AttrSizeGB: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the size of the SFS volume in gigabytes",
						},
						tfconstants.AttrStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the API status of the SFS instance",
						},
						tfconstants.AttrState: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the normalized state of the SFS instance",
						},
						tfconstants.AttrPrivateEndpoint: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the NFS mount endpoint for the SFS",
						},
						tfconstants.AttrPlan: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the plan of the SFS",
						},
						tfconstants.AttrIsBackupEnabled: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether backups are enabled for the SFS",
						},
						tfconstants.AttrIOPS: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the IOPS value of the SFS",
						},
						tfconstants.AttrVPCID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the id of the VPC for the SFS",
						},
						tfconstants.AttrEncryptionEnabled: {
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
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	// List all SFS instances
	sfsList, _, err := client.Sfs.ListSfss(ctx)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceDataSourceListErrorTemplate, ResourceName, projectID, region, err)
	}

	log.Printf("[INFO] SFS DATA SOURCE | Retrieved %d SFS instances", len(sfsList))

	if err := d.Set("sfs_list", flattenSfsList(sfsList)); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat("sfs_list"), err))
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
		oi[tfconstants.AttrID] = sfs.ID
		oi[tfconstants.AttrName] = sfs.Name
		oi[tfconstants.AttrStatus] = sfs.Status
		oi[tfconstants.AttrState] = normalizeSfsState(sfs.Status)
		oi[tfconstants.AttrPlan] = sfs.PlanName
		oi[tfconstants.AttrPrivateEndpoint] = sfs.PrivateIPAddress
		oi[tfconstants.AttrIsBackupEnabled] = sfs.IsBackupEnabled
		oi[tfconstants.AttrIOPS] = sfs.DiskIOPS
		oi[tfconstants.AttrVPCID] = sfs.VPCID
		oi[tfconstants.AttrEncryptionEnabled] = sfs.IsEncryptionEnabled

		// Parse disk size from string (e.g., "100GB")
		if sfs.DiskSize != "" {
			diskSizeStr := strings.TrimSpace(strings.ReplaceAll(sfs.DiskSize, "GB", ""))
			if sizeInt, err := strconv.Atoi(diskSizeStr); err == nil {
				oi[tfconstants.AttrSizeGB] = sizeInt
			}
		}

		ois[i] = oi
	}

	return ois
}
