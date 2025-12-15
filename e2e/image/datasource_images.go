package image

import (
	"context"
	"fmt"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadImages,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Image-specific fields
			"image_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of all saved Images",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tfconstants.AttrTemplateID: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the template used to create a Node from the Image",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the type of the Image",
						},
						"os_distribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the OS distribution of the Image",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Image",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "id of the Image",
						},
						"distro": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the distribution type of the Image",
						},
						"sku_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the SKU type of the Image",
						},
						"image_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state of the Image instance",
						},
					},
				},
			},
		},
	}
}

func dataSourceReadImages(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	log.Printf("[INFO] Inside images data source")

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

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))
	}

	images, _, err := goe2eClient.Images.GetSavedImages(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error finding saved images: %w", err))
	}

	d.Set("image_list", flattenSavedImages(images))
	d.SetId("saved_image_list")
	return nil
}

func flattenSavedImages(images []goe2e.SavedImage) []interface{} {
	if images == nil {
		return make([]interface{}, 0)
	}

	ois := make([]interface{}, len(images))
	for i, image := range images {
		oi := make(map[string]interface{})
		oi[tfconstants.AttrTemplateID] = image.TemplateID
		oi["distro"] = image.Distro
		oi["image_id"] = image.ImageID
		oi["image_state"] = image.ImageState
		oi["image_type"] = image.ImageType
		oi["name"] = image.Name
		oi["sku_type"] = image.SKUType
		oi["os_distribution"] = image.OSDistribution
		ois[i] = oi
	}

	return ois
}
