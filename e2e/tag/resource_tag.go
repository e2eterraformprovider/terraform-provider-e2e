package tag

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// Resource-specific attribute names
	attrMetadata = "metadata"
	attrLabelID  = "label_id"
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
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			tfconstants.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the tag (label_name in API)",
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			attrMetadata: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true, // Updates not supported by API
				Default:     "",
				Description: "Metadata/description for the tag",
			},
			// Computed fields
			attrLabelID: {
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
	name := d.Get(tfconstants.AttrName).(string)
	metadata := d.Get(attrMetadata).(string)

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
	d.Set(attrLabelID, tag.LabelID)
	d.Set(tfconstants.AttrProjectID, projectID)
	d.Set(tfconstants.AttrRegion, region)

	log.Printf("[INFO] Created tag: %s (ID: %s)", name, labelIDStr)

	return resourceTagRead(ctx, d, m)
}

func resourceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()

	projectID := d.Get(tfconstants.AttrProjectID).(string)
	region := d.Get(tfconstants.AttrRegion).(string)
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
	d.Set(tfconstants.AttrName, tag.LabelName)
	d.Set(attrMetadata, tag.Metadata)
	d.Set(attrLabelID, tag.LabelID)

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

	projectID := d.Get(tfconstants.AttrProjectID).(string)
	region := d.Get(tfconstants.AttrRegion).(string)
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
		d.Set(tfconstants.AttrProjectID, parts[0])
		d.Set(tfconstants.AttrRegion, parts[1])
		d.SetId(parts[2])
	} else if len(parts) == 1 {
		// Use provider defaults
		if cfg.DefaultProjectID != "" {
			d.Set(tfconstants.AttrProjectID, cfg.DefaultProjectID)
		} else {
			return nil, fmt.Errorf("project_id is required for import when not set in provider config")
		}
		if cfg.DefaultRegion != "" {
			d.Set(tfconstants.AttrRegion, cfg.DefaultRegion)
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

// isNotFoundError checks if an error indicates a tag was not found
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, goe2econstants.NotFoundSubstring) ||
		strings.Contains(errStr, goe2econstants.NotFoundCode)
}
