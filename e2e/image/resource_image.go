package image

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceImage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			tfconstants.AttrRegion: func() *schema.Schema {
				s := config.RegionSchema()
				s.Description = "region where the Image is stored"
				return s
			}(),
			tfconstants.AttrLocation: func() *schema.Schema {
				s := config.LocationSchema()
				s.Description = "location where the Image is stored"
				return s
			}(),
			tfconstants.AttrProjectID: func() *schema.Schema {
				s := config.ProjectIDSchemaResource()
				s.Description = "id of the Project to create the Image in"
				return s
			}(),

			// ============================================
			// REQUIRED INPUT FIELDS (Immutable)
			// ============================================
			tfconstants.AttrNodeID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of the Node to create the image from",
			},
			tfconstants.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the Image. Must not contain whitespace. Can be updated in-place in V3.",
				ValidateFunc: validateName,
			},
			// ============================================
			// OPTIONAL FIELDS
			// ============================================
			tfconstants.AttrTags: {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "tags for the Image (state-only until API support is added)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// ============================================
			// COMPUTED FIELDS - IMAGE METADATA
			// ============================================
			tfconstants.AttrTemplateID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the template used to create a Node from the Image",
			},
			"image_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Image instance",
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
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "normalized state of the Image (creating, ready, error, deleted)",
			},
			"image_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "size of the Image (e.g., \"95.368 GB\")",
			},
			"cloning_ops": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "number of ongoing cloning operations",
			},
			"running_vms": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "number of VMs running from this Image",
			},
			"is_windows": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether this is a Windows image",
			},
			"vm_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of VMs created from this Image",
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the Image",
			},
		},

		CreateContext: resourceCreateImage,
		ReadContext:   resourceReadImage,
		UpdateContext: resourceUpdateImage,
		DeleteContext: resourceDeleteImage,
		Exists:        resourceExistsImage,
		CustomizeDiff: resourceImageCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImageImport,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceImageResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceImageStateUpgradeV0toV1,
				Version: 0,
			},
		},
	}
}

func validateName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceCreateImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)

	log.Printf("[INFO] IMAGE CREATE")

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Log deprecation warning if location is used
	if _, ok := d.GetOk(tfconstants.AttrLocation); ok {
		log.Printf("[WARN] Parameter 'location' is deprecated and will be removed in v4.0. Please use 'region' instead")
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))
	}

	nodeID := d.Get(tfconstants.AttrNodeID).(string)
	imageName := d.Get(tfconstants.AttrName).(string)

	// Create image via SaveImage action on node
	saveReq := &goe2e.NodeSaveImageRequest{
		ActionType: "save_images",
		Name:       imageName,
	}

	result, _, err := goe2eClient.Nodes.SaveImage(ctx, nodeID, saveReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create image: %w", err))
	}

	if result.ImageID == "" {
		return diag.Errorf("image creation succeeded but no image ID returned")
	}

	log.Printf("[INFO] IMAGE CREATION | Image ID: %s", result.ImageID)
	d.SetId(result.ImageID)

	// Store tags in state if provided
	if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
		if err := d.Set(tfconstants.AttrTags, tags); err != nil {
			return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
		}
	}

	// Poll for image to reach Ready state
	log.Printf("[INFO] IMAGE CREATION | Polling for Ready state...")
	timeout := 30 * time.Minute // Default timeout
	if err := waitForImageState(ctx, goe2eClient, result.ImageID, "ready", timeout); err != nil {
		// If timeout or error, still set the ID so user can import/manage it
		log.Printf("[WARN] Image creation initiated but polling failed: %v", err)
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Image creation initiated but not yet ready",
				Detail:   fmt.Sprintf("Image %s was created but did not reach Ready state within timeout. Error: %v", result.ImageID, err),
			},
		}
	}

	log.Printf("[INFO] IMAGE CREATION | Image is Ready")
	return resourceReadImage(ctx, d, m)
}

func resourceReadImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)

	log.Printf("[INFO] IMAGE READ | Image ID: %s", d.Id())
	imageID := d.Id()

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))
	}

	image, _, err := goe2eClient.Images.GetImage(ctx, imageID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[INFO] IMAGE READ | Image not found, removing from state")
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading image %s: %w", imageID, err))
	}

	log.Printf("[INFO] IMAGE READ | Image found: %+v", image)

	// Flatten image response to state
	flattened := flattenImageResponse(image)

	// Set all fields from flattened response
	for key, value := range flattened {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(fmt.Errorf("error setting %s: %w", key, err))
		}
	}

	// Preserve tags from state (state-only until API support)
	if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
		if err := d.Set(tfconstants.AttrTags, tags); err != nil {
			return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
		}
	}

	return nil
}

func resourceImageImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Support both formats:
	// 1. Simple: image_id (uses provider defaults)
	// 2. Full: project_id/region/image_id
	parts := config.ParseImportID(d.Id())

	var projectID, region, imageID string

	if len(parts) == 1 {
		// Simple format: just image_id
		imageID = parts[0]
		// projectID and region will be determined from provider defaults in Read
	} else if len(parts) == 3 {
		// Full format: project_id/region/image_id
		projectID = parts[0]
		region = parts[1]
		imageID = parts[2]
	} else {
		return nil, fmt.Errorf("invalid import ID format, expected 'image_id' or 'project_id/region/image_id', got: %s", d.Id())
	}

	// Set the image ID
	d.SetId(imageID)

	// Set project_id and region if provided
	if projectID != "" {
		if err := d.Set(tfconstants.AttrProjectID, projectID); err != nil {
			return nil, fmt.Errorf("error setting project_id: %w", err)
		}
	}
	if region != "" {
		if err := d.Set(tfconstants.AttrRegion, region); err != nil {
			return nil, fmt.Errorf("error setting region: %w", err)
		}
	}

	// Trigger Read to populate all fields
	diags := resourceReadImage(ctx, d, m)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading image during import: %s", diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDeleteImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("[INFO] DELETE IMAGE | Image ID: %s", d.Id())
	imageID := d.Id()

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))
	}

	// Check if image has running VMs or cloning operations (warning only)
	image, _, err := goe2eClient.Images.GetImage(ctx, imageID)
	if err == nil {
		if image.RunningVMs != "" && image.RunningVMs != "0" {
			log.Printf("[WARN] Image has %s running VMs, deletion will proceed", image.RunningVMs)
		}
		if image.CloningOps != "" && image.CloningOps != "0" {
			return diag.Errorf("cannot delete image with ongoing cloning operations (cloning_ops: %s)", image.CloningOps)
		}
	}

	// Delete the image
	result, _, err := goe2eClient.Images.DeleteImage(ctx, imageID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[INFO] Image already deleted")
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("failed to delete image: %w", err))
	}

	// Check result status if available
	if result != nil && !result.Status {
		return diag.Errorf("delete failed: %s", result.Message)
	}

	log.Printf("[INFO] DELETE IMAGE | Delete successful")

	// Optional: Poll for deletion confirmation (wait for 404)
	timeout := 5 * time.Minute
	if err := waitForImageState(ctx, goe2eClient, imageID, "deleted", timeout); err != nil {
		log.Printf("[WARN] Image deletion initiated but confirmation polling failed: %v", err)
		// Still consider it deleted since API returned success
	}

	d.SetId("")
	return diags
}

func resourceExistsImage(d *schema.ResourceData, m interface{}) (bool, error) {
	cfg := m.(*config.Config)
	ctx := context.Background()

	imageID := d.Id()
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return false, err
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return false, err
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return false, fmt.Errorf("failed to create GoE2E client: %w", err)
	}

	_, _, err = goe2eClient.Images.GetImage(ctx, imageID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// resourceUpdateImage handles updates to the image resource
// Currently only supports name updates via rename action
func resourceUpdateImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)

	log.Printf("[INFO] IMAGE UPDATE | Image ID: %s", d.Id())
	imageID := d.Id()

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))
	}

	// Handle name update
	if d.HasChange(tfconstants.AttrName) {
		newName := d.Get(tfconstants.AttrName).(string)
		log.Printf("[INFO] IMAGE UPDATE | Renaming image to: %s", newName)

		renameReq := &goe2e.RenameImageRequest{
			ActionType: "rename",
			Name:       newName,
		}

		result, _, err := goe2eClient.Images.RenameImage(ctx, imageID, renameReq)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to rename image: %w", err))
		}

		if !result.Status {
			return diag.Errorf("rename failed: %s", result.Message)
		}

		log.Printf("[INFO] IMAGE UPDATE | Rename successful")
	}

	// Handle tags update (state-only)
	if d.HasChange(tfconstants.AttrTags) {
		if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
			if err := d.Set(tfconstants.AttrTags, tags); err != nil {
				return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
			}
		}
	}

	return resourceReadImage(ctx, d, m)
}

// resourceImageCustomizeDiff handles custom diff logic
// Validates field conflicts and emits deprecation warnings
func resourceImageCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	// Emit deprecation warning if location is used
	if _, ok := d.GetOk(tfconstants.AttrLocation); ok {
		log.Printf("[WARN] Parameter 'location' is deprecated and will be removed in v4.0. Please use 'region' instead")
	}

	// Validate that region and location are not both set (handled by ConflictsWith, but double-check)
	if _, hasRegion := d.GetOk(tfconstants.AttrRegion); hasRegion {
		if _, hasLocation := d.GetOk(tfconstants.AttrLocation); hasLocation {
			return fmt.Errorf("cannot set both 'region' and 'location' parameters")
		}
	}

	return nil
}

// resourceImageResourceV0 returns the V0 schema for state migration
func resourceImageResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),
			tfconstants.AttrNodeID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrTemplateID: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"image_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_distribution": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"distro": {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceImageStateUpgradeV0toV1 upgrades state from V0 to V1
// Adds new computed fields with default values, preserves all existing fields
func resourceImageStateUpgradeV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	// Add new computed fields with default values
	if _, ok := rawState["state"]; !ok {
		// Normalize image_state to state if available
		if imageState, ok := rawState["image_state"].(string); ok {
			rawState["state"] = normalizeImageState(imageState)
		} else {
			rawState["state"] = ""
		}
	}
	if _, ok := rawState["image_size"]; !ok {
		rawState["image_size"] = ""
	}
	if _, ok := rawState["cloning_ops"]; !ok {
		rawState["cloning_ops"] = ""
	}
	if _, ok := rawState["running_vms"]; !ok {
		rawState["running_vms"] = ""
	}
	if _, ok := rawState["is_windows"]; !ok {
		rawState["is_windows"] = false
	}
	if _, ok := rawState["sku_type"]; !ok {
		rawState["sku_type"] = ""
	}
	if _, ok := rawState["vm_info"]; !ok {
		rawState["vm_info"] = []interface{}{}
	}
	if _, ok := rawState["tags"]; !ok {
		rawState["tags"] = map[string]interface{}{}
	}

	return rawState, nil
}
