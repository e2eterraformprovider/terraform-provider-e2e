package container_registry

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceContainerRegistry() *schema.Resource {
	return &schema.Resource{
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
			tfconstants.AttrProjectName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the Container Registry project",
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - SECURITY SETTINGS
			// ============================================
			tfconstants.AttrPreventVulnerabilities: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to prevent vulnerable images",
			},
			tfconstants.AttrSeverity: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "low",
				ValidateFunc: validation.StringInSlice(
					[]string{"low", "medium", "high", "critical", "none"},
					false,
				),
				Description: "vulnerability severity threshold (low, medium, high, critical, none)",
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Container Registry instance",
			},
			"setup_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use 'status' instead. This parameter will be removed in version 3.0.0",
				Description: "DEPRECATED: Use 'status' instead",
			},

			// ============================================
			// COMPUTED FIELDS - CONFIGURATION
			// ============================================
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the domain name of the Container Registry",
			},
			"is_public": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the Container Registry project is public",
			},

			// ============================================
			// COMPUTED FIELDS - STORAGE
			// ============================================
			"project_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "the size of the Container Registry project in bytes",
			},
			"storage_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the storage limit of the Container Registry in bytes",
			},

			// ============================================
			// COMPUTED FIELDS - TIMESTAMPS
			// ============================================
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the timestamp when the Container Registry was created",
			},
			tfconstants.AttrUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the timestamp when the Container Registry was last updated",
			},

			// ============================================
			// OPTIONAL FIELDS - TAGS
			// ============================================
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "map of tags to assign to the resource (state-only)",
			},
		},
		CreateContext: resourceCreateContainerRegistry,
		ReadContext:   resourceReadContainerRegistry,
		DeleteContext: resourceDeleteContainerRegistry,
		UpdateContext: resourceUpdateContainerRegistry,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	apiClient := cfg.Goe2eClient()

	projectName := d.Get(tfconstants.AttrProjectName).(string)
	preventVul := d.Get(tfconstants.AttrPreventVulnerabilities).(bool)
	severity := d.Get(tfconstants.AttrSeverity).(string)

	createReq := &goe2e.ContainerRegistryCreateRequest{
		ProjectName: projectName,
		PreventVul:  fmt.Sprintf("%t", preventVul),
		Severity:    severity,
	}

	registry, _, err := apiClient.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Container Registry: %w", err))
	}

	if registry == nil {
		return diag.FromErr(fmt.Errorf("container registry created but response was empty"))
	}

	d.SetId(fmt.Sprintf("%d", registry.ID))

	// Set all fields from API response
	if err := setContainerRegistryState(d, registry); err != nil {
		return diag.FromErr(err)
	}

	// Initialize tags if provided (state-only, not sent to API)
	if tags, ok := d.GetOk("tags"); ok {
		if err := d.Set("tags", tags); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set tags: %w", err))
		}
	}

	return nil
}

func resourceReadContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	apiClient := cfg.Goe2eClient()

	id := d.Id()

	// Parse ID to int for the API call
	registryID, err := strconv.Atoi(id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid container registry ID: %w", err))
	}

	registry, _, err := apiClient.ContainerRegistry.GetContainerRegistry(ctx, registryID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read container registry (ID: %s): %w", id, err))
	}

	if registry == nil {
		log.Printf("[INFO] Container Registry project with ID %s not found; removing from state", id)
		d.SetId("")
		return nil
	}

	// Set all fields from API response
	if err := setContainerRegistryState(d, registry); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDeleteContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	apiClient := cfg.Goe2eClient()

	projectName := d.Get(tfconstants.AttrProjectName).(string)
	crProjectID := d.Id()

	// Get the registry details to extract customer ID
	registryID, err := strconv.Atoi(crProjectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid container registry ID: %w", err))
	}

	registry, _, err := apiClient.ContainerRegistry.GetContainerRegistry(ctx, registryID)
	if err != nil {
		log.Printf("[WARN] Failed to fetch container registry details for deletion, using default customer ID: %v", err)
	}

	// Get customer ID from registry or use default
	userID := "0"
	if registry != nil {
		userID = strconv.Itoa(registry.Customer)
	}

	deleteReq := &goe2e.ContainerRegistryDeleteRequest{
		CRProjectID: crProjectID,
		ProjectName: projectName,
		UserID:      userID,
	}

	_, err = apiClient.ContainerRegistry.DeleteContainerRegistry(ctx, deleteReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Container Registry: %w", err))
	}

	d.SetId("")
	return nil
}

func resourceUpdateContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	apiClient := cfg.Goe2eClient()

	// Update security settings if they changed
	if d.HasChange(tfconstants.AttrPreventVulnerabilities) || d.HasChange(tfconstants.AttrSeverity) {
		projectName := d.Get(tfconstants.AttrProjectName).(string)
		preventVul := fmt.Sprintf("%t", d.Get(tfconstants.AttrPreventVulnerabilities).(bool))
		severity := d.Get(tfconstants.AttrSeverity).(string)

		updateReq := &goe2e.ContainerRegistryUpdateRequest{
			PreventVul: preventVul,
			Severity:   severity,
		}

		_, err := apiClient.ContainerRegistry.UpdateContainerRegistry(ctx, projectName, updateReq)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update container registry: %w", err))
		}
	}

	// Tags are state-only, no API call needed
	// Just ensure they're set in state if changed
	if d.HasChange("tags") {
		if tags, ok := d.GetOk("tags"); ok {
			if err := d.Set("tags", tags); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set tags: %w", err))
			}
		}
	}

	// Refresh the state to get any updates
	return resourceReadContainerRegistry(ctx, d, m)
}

// setContainerRegistryState sets all fields from the API response into the Terraform state
func setContainerRegistryState(d *schema.ResourceData, registry *goe2e.ContainerRegistry) error {
	if err := d.Set(tfconstants.AttrProjectName, registry.ProjectName); err != nil {
		return fmt.Errorf("failed to set project_name: %w", err)
	}
	if err := d.Set(tfconstants.AttrPreventVulnerabilities, registry.PreventVul); err != nil {
		return fmt.Errorf("failed to set prevent_vul: %w", err)
	}
	if err := d.Set(tfconstants.AttrSeverity, registry.Severity); err != nil {
		return fmt.Errorf("failed to set severity: %w", err)
	}
	if err := d.Set(tfconstants.AttrStatus, registry.State); err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}
	if err := d.Set("setup_status", registry.State); err != nil {
		return fmt.Errorf("failed to set setup_status: %w", err)
	}
	if err := d.Set("domain_name", registry.DomainName); err != nil {
		return fmt.Errorf("failed to set domain_name: %w", err)
	}
	if err := d.Set("project_size", registry.ProjectSize); err != nil {
		return fmt.Errorf("failed to set project_size: %w", err)
	}
	if err := d.Set("storage_limit", registry.StorageLimit); err != nil {
		return fmt.Errorf("failed to set storage_limit: %w", err)
	}
	if err := d.Set("is_public", registry.IsPublic); err != nil {
		return fmt.Errorf("failed to set is_public: %w", err)
	}
	if err := d.Set(tfconstants.AttrCreatedAt, registry.CreatedAt); err != nil {
		return fmt.Errorf("failed to set created_at: %w", err)
	}
	if err := d.Set(tfconstants.AttrUpdatedAt, registry.UpdatedAt); err != nil {
		return fmt.Errorf("failed to set updated_at: %w", err)
	}
	return nil
}
