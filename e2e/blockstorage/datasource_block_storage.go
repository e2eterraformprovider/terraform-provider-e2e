package blockstorage

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dataSourceLogReadMessage = "[INFO] INSIDE BLOCK STORAGE DATA SOURCE | read"
)

func DataSourceBlockStorage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadBlockStorage,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Block storage-specific fields
			tfconstants.AttrBlockID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the Block Storage",
			},
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the Block Storage",
			},
			tfconstants.AttrSize: {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "the size of the Block Storage in gigabytes",
			},
			tfconstants.AttrIOPS: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the IOPS of the Block Storage",
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Block Storage instance",
			},
			tfconstants.AttrVMID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the node to which the block storage is attached (null if detached)",
			},
			tfconstants.AttrVMName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the node to which the block storage is attached (null if detached)",
			},
			// "created_on": {
			// 	Type:        schema.TypeString,
			// 	Computed:    true,
			// 	Description: "Creation time of the block storage",
			// },
		},
	}
}
func dataSourceReadBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	log.Print(dataSourceLogReadMessage)

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create goe2e client with projectID and region
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	blockStorageID := d.Get(tfconstants.AttrBlockID).(string)

	blockStorage, resp, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, blockStorageID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return diag.Errorf("Block storage with ID %s not found", blockStorageID)
		}
		if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) {
			return diag.Errorf("Block storage with ID %s not found", blockStorageID)
		}
		return diag.Errorf("Error retrieving block storage (ID: %s) in project (%s), region (%s): %s", blockStorageID, projectIDStr, region, err.Error())
	}

	if blockStorage == nil {
		return diag.Errorf("Block storage with ID %s not found", blockStorageID)
	}

	d.SetId(blockStorageID)
	log.Printf("[INFO] BLOCK STORAGE DATA SOURCE | READ | data : %+v", blockStorage)

	// Set fields from the BlockStorage struct
	// Note: Size from API is in GB (int), convert to float64 for schema
	d.Set(tfconstants.AttrSize, float64(blockStorage.Size))
	d.Set(tfconstants.AttrName, blockStorage.Name)
	d.Set(tfconstants.AttrStatus, blockStorage.Status)
	d.Set(tfconstants.AttrIOPS, blockStorage.Template.TotalIOPSSec)

	// Handle VM attachment details if present
	if blockStorage.VMDetail != nil && len(blockStorage.VMDetail) > 0 {
		if vmID, ok := blockStorage.VMDetail["vm_id"]; ok {
			if vmIDFloat, ok := vmID.(float64); ok {
				d.Set(tfconstants.AttrVMID, strconv.Itoa(int(vmIDFloat)))
			} else if vmIDStr, ok := vmID.(string); ok {
				d.Set(tfconstants.AttrVMID, vmIDStr)
			}
		}
		if vmName, ok := blockStorage.VMDetail["vm_name"]; ok {
			if vmNameStr, ok := vmName.(string); ok {
				d.Set(tfconstants.AttrVMName, vmNameStr)
			}
		}
	}

	log.Printf("[INFO] BLOCK STORAGE DATA SOURCE | d : %+v", d)

	return diags
}
