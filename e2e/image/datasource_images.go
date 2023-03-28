package image

import (
	"context"
	// "fmt"
	"log"
	// "math"
	// "regexp"

	// "strconv"
	//"strings"

	"github.com/devteametwoe/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/devteametwoe/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceImages() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"image_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},
						// "vm_info":&schema.Schema{
						// 	Type:schema
						// }
						"image_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_distribution": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"distro": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"sku_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},

		ReadContext: dataSourceReadImages,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadImages(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside images data source ")
	Response, err := apiClient.GetSavedImages()
	if err != nil {
		return diag.Errorf("error finding saved images")
	}

	d.Set("image_list", flattenImages(&Response.Data))
	d.SetId("saved_image_list")
	var diags diag.Diagnostics
	return diags
}
func flattenImages(imageList *[]models.Image) []interface{} {

	if imageList != nil {

		ois := make([]interface{}, len(*imageList), len(*imageList))

		for i, image := range *imageList {

			oi := make(map[string]interface{})
			oi["template_id"] = image.Template_id
			oi["distro"] = image.Distro
			oi["image_id"] = image.Image_id
			oi["image_state"] = image.Image_state
			oi["image_type"] = image.Image_type
			oi["name"] = image.Name
			oi["sku_type"] = image.Sku_type
			oi["os_distribution"] = image.Os_distribution
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
