package objectstore

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceObjectStore() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceObjectStoreResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceObjectStoreStateUpgradeV0toV1,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED INPUT FIELDS (Immutable)
			// ============================================
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the object store bucket",
			},

			// ============================================
			// OPTIONAL INPUT FIELDS (Mutable)
			// ============================================
			"enabling_versioning": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Deprecated:  "Use versioning_enabled instead. This field will be removed in v4.0.",
				Description: "DEPRECATED: whether to enable versioning for the bucket. Use versioning_enabled instead.",
			},

			// V3: NEW Versioning field with clearer naming
			"versioning_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to enable versioning for the bucket",
			},

			// V3: Advanced features
			"encryption_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to enable encryption for the bucket",
			},

			"lock_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to enable object lock for the bucket",
			},

			"public_access_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to enable public access to the bucket",
			},

			// V3: Tags for resource organization
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "map of tags to assign to the bucket",
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the object store bucket",
			},
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the object store bucket",
			},

			// ============================================
			// COMPUTED FIELDS - CONFIGURATION
			// ============================================
			tfconstants.AttrVersioningStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the versioning state of the bucket",
			},
			tfconstants.AttrLifecycleConfigurationStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the lifecycle configuration state of the bucket",
			},

			// V3: NEW Computed fields
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the last update date for the object store bucket",
			},

			"created_by": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the user ID who created the bucket",
			},

			"bucket_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the total size of the bucket",
			},

			"is_cdn_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether CDN is attached to the bucket",
			},

			"is_encryption_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether encryption is enabled for the bucket",
			},

			"is_lock_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether object lock is enabled for the bucket",
			},

			"is_public_access_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether public access is enabled for the bucket",
			},
		},
		CreateContext: resourceCreateBucket,
		ReadContext:   resourceReadBucket,
		UpdateContext: resourceUpdateBucket,
		DeleteContext: resourceDeleteBucket,
		CustomizeDiff: resourceObjectStoreCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceObjectStoreImport,
		},
	}
}

// resourceObjectStoreResourceV0 returns the schema for schema version 0
func resourceObjectStoreResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),
			tfconstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabling_versioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrVersioningStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrLifecycleConfigurationStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreateBucket(ctx context.Context, resourceData *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {
	cfg := clientInterface.(*config.Config)
	var diags diag.Diagnostics

	region, err := cfg.GetRegionOrDefault(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	projectIDStr, err := cfg.GetProjectIDOrDefault(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client: %s", err)
	}

	bucketName := resourceData.Get("name").(string)
	log.Printf("[INFO] BUCKET CREATE STARTS - name: %s", bucketName)

	// Log deprecation warning if enabling_versioning is used
	if enablingVersioning, ok := resourceData.GetOk("enabling_versioning"); ok && enablingVersioning.(bool) {
		log.Printf("[WARN] enabling_versioning is deprecated. Use versioning_enabled instead.")
	}

	// Determine versioning state from versioning_enabled (preferred) or enabling_versioning (deprecated)
	versioningEnabled := resourceData.Get("versioning_enabled").(bool)
	if !versioningEnabled && resourceData.Get("enabling_versioning").(bool) {
		versioningEnabled = true
	}

	createReq := &goe2e.BucketCreateRequest{
		BucketName: bucketName,
	}

	bucket, _, err := goe2eClient.ObjectStorage.CreateBucket(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating object storage bucket (name: %s) in project (%s), region (%s): %s", bucketName, projectIDStr, region, err)
	}

	log.Printf("[INFO] BUCKET CREATE | RESPONSE | %+v", bucket)
	resourceData.SetId(fmt.Sprintf("%d", bucket.ID))

	// Set all fields from API response
	_ = setObjectStoreDataFromAPI(resourceData, bucket, versioningEnabled)

	return diags
}

func resourceReadBucket(ctx context.Context, resourceData *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {

	cfg := clientInterface.(*config.Config)
	var diags diag.Diagnostics
	log.Printf("[info] inside Object Store Resource read")

	region, err := cfg.GetRegionOrDefault(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	projectIDStr, err := cfg.GetProjectIDOrDefault(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client: %s", err)
	}

	bucketName := resourceData.Get("name").(string)

	bucket, _, err := goe2eClient.ObjectStorage.GetBucket(ctx, bucketName)
	if err != nil {
		return diag.Errorf("Error retrieving object storage bucket (name: %s) in project (%s), region (%s): %s", bucketName, projectIDStr, region, err)
	}

	if bucket == nil {
		log.Printf("[INFO] Bucket not found (name: %s), removing from state", bucketName)
		resourceData.SetId("")
		return diags
	}

	log.Printf("[info] Object Store Resource read | before setting data")
	log.Printf("[INFO] Object Store Data: %+v", bucket)

	// Determine versioning state from API response
	versioningEnabled := bucket.VersioningStatus == "Enabled"

	// Set all fields from API response
	_ = setObjectStoreDataFromAPI(resourceData, bucket, versioningEnabled)

	log.Printf("[info] Object Store Resource read | after setting data")

	return diags

}

func resourceUpdateBucket(ctx context.Context, resourceData *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {

	cfg := clientInterface.(*config.Config)

	region, err := cfg.GetRegionOrDefault(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	projectIDStr, err := cfg.GetProjectIDOrDefault(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client: %s", err)
	}

	bucketName := resourceData.Get("name").(string)

	// Handle versioning updates (both old and new field names)
	if resourceData.HasChange("enabling_versioning") || resourceData.HasChange("versioning_enabled") {
		bucketstatus := resourceData.Get("status").(string)
		log.Printf("[INFO] Bucket status: %s", bucketstatus)

		// Get versioning state from new field (preferred) or old field
		versioningEnabled := resourceData.Get("versioning_enabled").(bool)
		if !versioningEnabled && resourceData.Get("enabling_versioning").(bool) {
			versioningEnabled = true
		}

		var action string
		if versioningEnabled {
			action = "Enabled"
		} else {
			action = "Suspended"
		}

		versioningReq := &goe2e.BucketVersioningRequest{
			BucketName:         bucketName,
			NewVersioningState: action,
		}

		bucketVersioning, _, err := goe2eClient.ObjectStorage.SetBucketVersioning(ctx, bucketName, versioningReq)
		if err != nil {
			return diag.Errorf("Error updating versioning (%s) for object storage bucket (name: %s) in project (%s), region (%s): %s", action, bucketName, projectIDStr, region, err)
		}
		resourceData.Set(tfconstants.AttrVersioningStatus, bucketVersioning.VersioningStatus)
		resourceData.Set("versioning_enabled", versioningEnabled)
		resourceData.Set("enabling_versioning", versioningEnabled) // Also update deprecated field for backwards compat
	}

	// Handle tags updates (state-only for now)
	if resourceData.HasChange("tags") {
		tags := resourceData.Get("tags").(map[string]interface{})
		log.Printf("[INFO] Updating tags: %v", tags)
		// Tags are state-only, no API call needed
		resourceData.Set("tags", tags)
	}

	return resourceReadBucket(ctx, resourceData, clientInterface)

}

func resourceDeleteBucket(ctx context.Context, resourceData *schema.ResourceData, clientInterface interface{}) diag.Diagnostics {
	cfg := clientInterface.(*config.Config)
	var diags diag.Diagnostics

	region, err := cfg.GetRegionOrDefault(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	projectIDStr, err := cfg.GetProjectIDOrDefault(resourceData)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client: %s", err)
	}

	bucketName := resourceData.Get("name").(string)

	// Pre-delete validation: check if bucket has lock policy
	if lockEnabled, ok := resourceData.GetOk("is_lock_enabled"); ok && lockEnabled.(bool) {
		return diag.Errorf("Cannot delete bucket with lock policy enabled (name: %s). Disable lock first.", bucketName)
	}

	log.Printf("[INFO] Deleting object storage bucket (name: %s) in project (%s), region (%s)", bucketName, projectIDStr, region)
	_, err = goe2eClient.ObjectStorage.DeleteBucket(ctx, bucketName)
	if err != nil {
		return diag.Errorf("Error deleting object storage bucket (name: %s) in project (%s), region (%s): %s", bucketName, projectIDStr, region, err)
	}

	resourceData.SetId("")
	log.Printf("[INFO] Successfully deleted object storage bucket (name: %s)", bucketName)
	return diags
}

// setObjectStoreDataFromAPI populates schema fields from API response data
func setObjectStoreDataFromAPI(resourceData *schema.ResourceData, bucket *goe2e.Bucket, versioningEnabled bool) error {
	// Status fields
	if bucket.CreatedAt != "" {
		resourceData.Set(tfconstants.AttrCreatedAt, bucket.CreatedAt)
	}
	if bucket.Status != "" {
		resourceData.Set("status", bucket.Status)
	}
	if bucket.VersioningStatus != "" {
		resourceData.Set(tfconstants.AttrVersioningStatus, bucket.VersioningStatus)
	}
	if bucket.LifecycleConfigurationStatus != "" {
		resourceData.Set(tfconstants.AttrLifecycleConfigurationStatus, bucket.LifecycleConfigurationStatus)
	}

	// Bucket size
	if bucket.BucketSize != "" {
		resourceData.Set("bucket_size", bucket.BucketSize)
	}

	// Versioning fields (both old and new for backwards compatibility)
	resourceData.Set("versioning_enabled", versioningEnabled)
	resourceData.Set("enabling_versioning", versioningEnabled)

	return nil
}

func resourceObjectStoreImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	var projectID, region, bucketName string

	// Support two import formats:
	// 1. Simple: bucket_name (uses provider defaults for project_id and region)
	// 2. Full: project_id:region:bucket_name
	if len(parts) == 1 {
		// Simple format: just bucket name
		bucketName = parts[0]
		cfg := meta.(*config.Config)
		var err error
		region, err = cfg.GetRegionOrDefault(d)
		if err != nil {
			return nil, fmt.Errorf("error getting region: %w", err)
		}
		var projectIDStr string
		projectIDStr, err = cfg.GetProjectIDOrDefault(d)
		if err != nil {
			return nil, fmt.Errorf("error getting project ID: %w", err)
		}
		projectID = projectIDStr
	} else if len(parts) == 3 {
		// Full format: project_id:region:bucket_name
		projectID = parts[0]
		region = parts[1]
		bucketName = parts[2]
	} else {
		return nil, fmt.Errorf("invalid import ID format: expected 'bucket_name' or 'project_id:region:bucket_name', got '%s'", d.Id())
	}

	if err := d.Set(tfconstants.AttrProjectID, projectID); err != nil {
		return nil, err
	}
	if err := d.Set(tfconstants.AttrRegion, region); err != nil {
		return nil, err
	}
	if err := d.Set(tfconstants.AttrName, bucketName); err != nil {
		return nil, err
	}

	// SetId will be replaced by actual ID in Read
	d.SetId(bucketName)

	return []*schema.ResourceData{d}, nil
}

// resourceObjectStoreCustomizeDiff handles field validation and deprecation warnings
func resourceObjectStoreCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Warn about deprecated enabling_versioning field
	if d.HasChange("enabling_versioning") || d.Get("enabling_versioning") != nil {
		oldVal, newVal := d.GetChange("enabling_versioning")
		if oldVal != nil || (newVal != nil && newVal != false) {
			log.Printf("[WARN] enabling_versioning is deprecated. Use versioning_enabled instead.")
		}
	}

	// Validate that both versioning fields are not both set
	enablingVersioning := d.Get("enabling_versioning").(bool)
	versioningEnabled := d.Get("versioning_enabled").(bool)

	if enablingVersioning && versioningEnabled {
		return fmt.Errorf("cannot set both 'enabling_versioning' and 'versioning_enabled'. Use 'versioning_enabled' instead")
	}

	return nil
}

// resourceObjectStoreStateUpgradeV0toV1 handles state migration from schema version 0 to 1
func resourceObjectStoreStateUpgradeV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[INFO] Upgrading Object Store state from V0 to V1")

	// Preserve all existing fields
	// If enabling_versioning is true, also set versioning_enabled to true
	if enablingVersioning, ok := rawState["enabling_versioning"]; ok && enablingVersioning.(bool) {
		rawState["versioning_enabled"] = true
	} else {
		rawState["versioning_enabled"] = false
	}

	// Initialize new V3 fields with default values
	if _, ok := rawState["encryption_enabled"]; !ok {
		rawState["encryption_enabled"] = false
	}
	if _, ok := rawState["lock_enabled"]; !ok {
		rawState["lock_enabled"] = false
	}
	if _, ok := rawState["public_access_enabled"]; !ok {
		rawState["public_access_enabled"] = false
	}
	if _, ok := rawState["tags"]; !ok {
		rawState["tags"] = map[string]interface{}{}
	}

	// Initialize computed fields with empty values if not present
	if _, ok := rawState["updated_at"]; !ok {
		rawState["updated_at"] = ""
	}
	if _, ok := rawState["created_by"]; !ok {
		rawState["created_by"] = 0
	}
	if _, ok := rawState["bucket_size"]; !ok {
		rawState["bucket_size"] = ""
	}
	if _, ok := rawState["is_cdn_attached"]; !ok {
		rawState["is_cdn_attached"] = false
	}
	if _, ok := rawState["is_encryption_enabled"]; !ok {
		rawState["is_encryption_enabled"] = false
	}
	if _, ok := rawState["is_lock_enabled"]; !ok {
		rawState["is_lock_enabled"] = false
	}
	if _, ok := rawState["is_public_access_enabled"]; !ok {
		rawState["is_public_access_enabled"] = false
	}

	log.Printf("[INFO] Object Store state upgrade V0 to V1 completed")
	return rawState, nil
}
