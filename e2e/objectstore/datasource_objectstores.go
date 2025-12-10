package objectstore

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceObjectStores() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadBuckets,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Common fields
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"bucket_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of Object Store buckets",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						e2econstants.AttrID: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the Object Store bucket",
						},
						e2econstants.AttrName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Object Store bucket",
						},
						"bucket_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the size of the Object Store bucket",
						},
						e2econstants.AttrStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state of the Object Store bucket instance",
						},
						e2econstants.AttrCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "the creation date for the Object Store bucket",
						},
						e2econstants.AttrVersioningStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "whether versioning is enabled for the bucket",
						},
						e2econstants.AttrLifecycleConfigurationStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "whether lifecycle rules are configured for the bucket",
						},
					},
				},
			},
		},
	}
}

func dataSourceReadBuckets(ctx context.Context, resourceDataSource *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cfg := clientInterface.(*config.Config)
	log.Printf("[INFO] ---- Execute Get Request to fetch Buckets Data via GoE2EClient. ---- ")

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(resourceDataSource)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(resourceDataSource)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get GoE2E client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client for datasource: %s", err)
	}

	// List buckets using GoE2EClient
	buckets, _, err := goe2eClient.ObjectStorage.ListBuckets(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] BUCKETS DATA SOURCE | fetched %d buckets", len(buckets))
	eos_bucket_list := flattenBuckets(buckets)
	_ = resourceDataSource.Set("bucket_list", eos_bucket_list)
	resourceDataSource.SetId("bucket_list")
	return diags
}

// flattenBuckets converts GoE2E Bucket structs to datasource schema format
// buckets is of type []goe2e.Bucket from the GoE2EClient
func flattenBuckets(buckets interface{}) []interface{} {
	if buckets == nil {
		return make([]interface{}, 0)
	}

	// Assert to slice - buckets come from goe2e.Bucket which is marshaled as interface{}
	bucketsList, ok := buckets.([]interface{})
	if !ok {
		log.Printf("[WARN] Failed to assert buckets to []interface{}, returning empty list")
		return make([]interface{}, 0)
	}

	result := make([]interface{}, 0, len(bucketsList))

	for _, bucket := range bucketsList {
		// Each bucket should be a map[string]interface{} when unmarshaled from JSON
		if bktMap, ok := bucket.(map[string]interface{}); ok {
			log.Printf("[INFO] Processing bucket: %v", bktMap)
			eos_bucket := make(map[string]interface{})
			eos_bucket["id"] = bktMap["id"]
			eos_bucket["name"] = bktMap["name"]
			eos_bucket["bucket_size"] = bktMap["bucket_size"]
			eos_bucket["created_at"] = bktMap["created_at"]
			eos_bucket["status"] = bktMap["status"]
			eos_bucket[e2econstants.AttrLifecycleConfigurationStatus] = bktMap["lifecycle_configuration_status"]
			eos_bucket[e2econstants.AttrVersioningStatus] = bktMap["versioning_status"]
			result = append(result, eos_bucket)
		}
	}
	return result
}
