package objectstore

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceObjectStore() schema.Resource {
	return schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region should specified",
			},
			"bucket_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of Buckets for the user",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Id for the My-Account bucket",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name for the My-Account bucket",
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
							ValidateFunc: validation.StringInSlice(
								[]string{
									"NEW", "CREATING", "FAILED", "WARNING", "AVAILABLE", "DELETED",
								},
								false,
							),
						},
						"bucket_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Size of the My-Account Bucket",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created Time of My-Account Bucket",
						},
						"versioning_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Is Versioning enabled?",
						},
						"lifecycle_configuration_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Is Lifecycle Rule Configured?",
						},
					},
				},
			},
		},
		ReadContext: dataSourceReadBuckets,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadBuckets(context context.Context, resourceDataSource *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := clientInterface.(*client.Client)
	log.Printf("[INFO] ---- Execute Get Request to fetch Buckets Data. ---- ")
	Response, err := apiClient.GetBuckets(resourceDataSource.Get("region").(string), resourceDataSource.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] %v", Response)
	log.Printf("[INFO] NODES DATA SOURCE | before setting")
	resourceDataSource.Set("nodes_list", flattenBuckets(&Response.Data))
	resourceDataSource.SetId("nodes_list")

	return diags
}

func flattenBuckets(buckets *[]models.ObjectStore) []interface{} {

	if buckets != nil {
		buckets_list := make([]interface{}, len(*buckets), len(*buckets))

		for i, bucket := range *buckets {
			eos_bucket := make(map[string]interface{})
			eos_bucket["id"] = bucket.ID
			eos_bucket["name"] = bucket.Name
			eos_bucket["size"] = bucket.BucketSize
			eos_bucket["created_on"] = bucket.CreatedOn
			eos_bucket["life_cycle_sonfiguration_status"] = bucket.LifecycleConfigurationStatus
			eos_bucket["versioning_status"] = bucket.VersioningStatus
			eos_bucket["status"] = bucket.Status
			buckets_list[i] = eos_bucket
		}
		return buckets_list
	}
	return make([]interface{}, 0)
}
