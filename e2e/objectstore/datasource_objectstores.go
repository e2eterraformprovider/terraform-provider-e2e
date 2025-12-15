package objectstore

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
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
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"bucket_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of Object Store buckets",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tfconstants.AttrID: {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "id of the Object Store bucket",
						},
						tfconstants.AttrName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Object Store bucket",
						},
						"bucket_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the size of the Object Store bucket",
						},
						tfconstants.AttrStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state of the Object Store bucket instance",
						},
						tfconstants.AttrCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "the creation date for the Object Store bucket",
						},
						tfconstants.AttrVersioningStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "whether versioning is enabled for the bucket",
						},
						tfconstants.AttrLifecycleConfigurationStatus: {
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
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	// List buckets using GoE2EClient
	buckets, _, err := goe2eClient.ObjectStorage.ListBuckets(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] BUCKETS DATA SOURCE | fetched %d buckets", len(buckets))
	// Convert []goe2e.Bucket to []interface{} for flattenBuckets
	bucketsInterface := make([]interface{}, 0, len(buckets))
	for _, b := range buckets {
		// Create a map from each bucket for compatibility with flattenBuckets
		bucketMap := map[string]interface{}{
			"id":                             b.ID,
			"name":                           b.Name,
			"bucket_size":                    b.BucketSize,
			"status":                         b.Status,
			"created_at":                     b.CreatedAt,
			"versioning_status":              b.VersioningStatus,
			"lifecycle_configuration_status": b.LifecycleConfigurationStatus,
		}
		bucketsInterface = append(bucketsInterface, bucketMap)
	}
	eos_bucket_list := flattenBuckets(bucketsInterface)
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
			// Use constants for Terraform state attribute names
			eos_bucket[tfconstants.AttrID] = bktMap["id"]
			eos_bucket[tfconstants.AttrName] = bktMap["name"]
			eos_bucket["bucket_size"] = bktMap["bucket_size"] // No constant defined for bucket_size
			eos_bucket[tfconstants.AttrCreatedAt] = bktMap["created_at"]
			eos_bucket[tfconstants.AttrStatus] = bktMap["status"]
			eos_bucket[tfconstants.AttrLifecycleConfigurationStatus] = bktMap["lifecycle_configuration_status"]
			eos_bucket[tfconstants.AttrVersioningStatus] = bktMap["versioning_status"]
			result = append(result, eos_bucket)
		}
	}
	return result
}
