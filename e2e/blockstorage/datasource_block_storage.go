package blockstorage

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBlockStorage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadBlockStorage,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Block storage-specific fields
			"block_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the Block Storage",
			},
			e2econstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the Block Storage",
			},
			e2econstants.AttrSize: {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "the size of the Block Storage in gigabytes",
			},
			e2econstants.AttrIOPS: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the IOPS of the Block Storage",
			},
			e2econstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Block Storage instance",
			},
			e2econstants.AttrVMID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the node to which the block storage is attached (null if detached)",
			},
			e2econstants.AttrVMName: {
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
	log.Printf("[INFO] INSIDE BLOCK STORAGE DATA SOURCE | read")

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
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	blockStorageID := d.Get("block_id").(string)

	blockStorage, resp, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, blockStorageID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.Errorf("Block storage with ID %s not found", blockStorageID)
		}
		if strings.Contains(err.Error(), "not found") {
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
	d.Set(e2econstants.AttrSize, float64(blockStorage.Size))
	d.Set(e2econstants.AttrName, blockStorage.Name)
	d.Set(e2econstants.AttrStatus, blockStorage.Status)
	d.Set(e2econstants.AttrIOPS, blockStorage.Template.TotalIOPSSec)

	// Handle VM attachment details if present
	if blockStorage.VMDetail != nil {
		if vmID, ok := blockStorage.VMDetail["vm_id"]; ok {
			if vmIDFloat, ok := vmID.(float64); ok {
				d.Set(e2econstants.AttrVMID, strconv.Itoa(int(vmIDFloat)))
			}
		}
		if vmName, ok := blockStorage.VMDetail["vm_name"]; ok {
			if vmNameStr, ok := vmName.(string); ok {
				d.Set(e2econstants.AttrVMName, vmNameStr)
			}
		}
	} else {
		d.Set(e2econstants.AttrVMID, nil)
		d.Set(e2econstants.AttrVMName, nil)
	}

	log.Printf("[INFO] BLOCK STORAGE DATA SOURCE | d : %+v", d)

	return diags
}
