package tag

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceTag returns the resource schema for E2E tags (labels)
func ResourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagCreate,
		ReadContext:   resourceTagRead,
		UpdateContext: resourceTagUpdate,
		DeleteContext: resourceTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTagImport,
		},

		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			e2econstants.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the tag (label_name in API)",
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"metadata": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true, // Updates not supported by API
				Default:     "",
				Description: "Metadata/description for the tag",
			},
			// Computed fields
			"label_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the tag (label)",
			},
		},
	}
}

func resourceTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()

	// Get project_id and region
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, err := config.GetRegionOrLocation(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Prepare request
	name := d.Get(e2econstants.AttrName).(string)
	metadata := d.Get("metadata").(string)

	createReq := &goe2e.TagCreateRequest{
		LabelName: name,
		Metadata:  metadata,
	}

	log.Printf("[DEBUG] Creating tag: name=%s, metadata=%s, project=%s, region=%s", name, metadata, projectID, region)

	// Create tag via goe2e client
	tag, _, err := client.Tags.CreateTag(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating tag %s: %s", name, err)
	}

	// Set ID and computed fields
	labelIDStr := strconv.Itoa(tag.LabelID)
	d.SetId(labelIDStr)
	d.Set("label_id", tag.LabelID)
	d.Set(e2econstants.AttrProjectID, projectID)
	d.Set(e2econstants.AttrRegion, region)

	log.Printf("[INFO] Created tag: %s (ID: %s)", name, labelIDStr)

	return resourceTagRead(ctx, d, m)
}

func resourceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()

	projectID := d.Get(e2econstants.AttrProjectID).(string)
	region := d.Get(e2econstants.AttrRegion).(string)
	tagID := d.Id()

	log.Printf("[DEBUG] Reading tag: id=%s, project=%s, region=%s", tagID, projectID, region)

	// Get tag via goe2e client
	tag, _, err := client.Tags.GetTag(ctx, tagID)
	if err != nil {
		if isNotFoundError(err) {
			log.Printf("[WARN] Tag %s not found, removing from state", tagID)
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading tag %s: %s", tagID, err)
	}

	// Update state
	d.Set(e2econstants.AttrName, tag.LabelName)
	d.Set("metadata", tag.Metadata)
	d.Set("label_id", tag.LabelID)

	return nil
}

func resourceTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Tags are ForceNew for all updatable fields since API doesn't support updates
	// This function should never be called, but we include it for safety
	return diag.Errorf("Tag updates are not supported - all changes require recreation")
}

func resourceTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()

	projectID := d.Get(e2econstants.AttrProjectID).(string)
	region := d.Get(e2econstants.AttrRegion).(string)
	tagID := d.Id()

	log.Printf("[DEBUG] Deleting tag: id=%s, project=%s, region=%s", tagID, projectID, region)

	// Delete tag via goe2e client
	_, err := client.Tags.DeleteTag(ctx, tagID)
	if err != nil {
		if isNotFoundError(err) {
			log.Printf("[WARN] Tag %s already deleted", tagID)
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting tag %s: %s", tagID, err)
	}

	log.Printf("[INFO] Deleted tag: %s", tagID)
	d.SetId("")

	return nil
}

func resourceTagImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Import format: <project_id>/<region>/<tag_id> or just <tag_id>
	parts := config.ParseImportID(d.Id())

	cfg := m.(*config.Config)

	if len(parts) == 3 {
		d.Set(e2econstants.AttrProjectID, parts[0])
		d.Set(e2econstants.AttrRegion, parts[1])
		d.SetId(parts[2])
	} else if len(parts) == 1 {
		// Use provider defaults
		if cfg.DefaultProjectID != "" {
			d.Set(e2econstants.AttrProjectID, cfg.DefaultProjectID)
		} else {
			return nil, fmt.Errorf("project_id is required for import when not set in provider config")
		}
		if cfg.DefaultRegion != "" {
			d.Set(e2econstants.AttrRegion, cfg.DefaultRegion)
		} else {
			return nil, fmt.Errorf("region is required for import when not set in provider config")
		}
		d.SetId(parts[0])
	} else {
		return nil, fmt.Errorf("invalid import ID format: expected <tag_id> or <project_id>/<region>/<tag_id>, got %s", d.Id())
	}

	// Read the tag to populate all fields
	diags := resourceTagRead(ctx, d, m)
	if diags.HasError() {
		return nil, fmt.Errorf("error importing tag: %s", diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}

// isNotFoundError checks if an error indicates a resource was not found
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "not found") || contains(errStr, "404")
}

// contains checks if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && containsSubstring(str, substr)
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
