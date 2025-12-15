package container_registry

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceContainerRegistry() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceContainerRegistryV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceContainerRegistryStateUpgradeV0toV1,
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
				Default:  tfconstants.ContainerRegistryDefaultSeverity,
				ValidateFunc: validation.StringInSlice(
					[]string{
						goe2econstants.ContainerRegistrySeverityLow,
						goe2econstants.ContainerRegistrySeverityMedium,
						goe2econstants.ContainerRegistrySeverityHigh,
						goe2econstants.ContainerRegistrySeverityCritical,
						goe2econstants.ContainerRegistrySeverityNone,
					},
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
			tfconstants.AttrSetupStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  DeprecationMessageSetupStatus,
				Description: DeprecationMessageSetupStatusAlternative,
			},

			// ============================================
			// COMPUTED FIELDS - CONFIGURATION
			// ============================================
			tfconstants.AttrDomainName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the domain name of the Container Registry",
			},
			tfconstants.AttrIsPublic: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the Container Registry project is public",
			},

			// ============================================
			// COMPUTED FIELDS - STORAGE
			// ============================================
			tfconstants.AttrProjectSize: {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "the size of the Container Registry project in bytes",
			},
			tfconstants.AttrStorageLimit: {
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
			tfconstants.AttrTags: {
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
		return diag.FromErr(fmt.Errorf(ErrorCreateRegistry, err))
	}

	if registry == nil {
		return diag.Errorf(ErrorCreateResponseEmpty)
	}

	d.SetId(fmt.Sprintf("%d", registry.ID))

	// Set all fields from API response
	if err := setContainerRegistryState(d, registry); err != nil {
		return diag.FromErr(err)
	}

	// Initialize tags if provided (state-only, not sent to API)
	if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
		if err := d.Set(tfconstants.AttrTags, tags); err != nil {
			return diag.FromErr(fmt.Errorf(ErrorSetField, tfconstants.AttrTags, err))
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
		return diag.FromErr(fmt.Errorf(ErrorInvalidID, err))
	}

	registry, _, err := apiClient.ContainerRegistry.GetContainerRegistry(ctx, registryID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(ErrorReadRegistry, id, err))
	}

	if registry == nil {
		log.Printf(LogRegistryNotFound, id)
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
		return diag.FromErr(fmt.Errorf(ErrorInvalidID, err))
	}

	registry, _, err := apiClient.ContainerRegistry.GetContainerRegistry(ctx, registryID)
	if err != nil {
		log.Printf(LogDeleteWarning, err)
	}

	// Get customer ID from registry or use default
	userID := tfconstants.ContainerRegistryDefaultCustomerID
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
		return diag.FromErr(fmt.Errorf(ErrorDeleteRegistry, err))
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
			return diag.FromErr(fmt.Errorf(ErrorUpdateRegistry, err))
		}
	}

	// Tags are state-only, no API call needed
	// Just ensure they're set in state if changed
	if d.HasChange(tfconstants.AttrTags) {
		if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
			if err := d.Set(tfconstants.AttrTags, tags); err != nil {
				return diag.FromErr(fmt.Errorf(ErrorSetField, tfconstants.AttrTags, err))
			}
		}
	}

	// Refresh the state to get any updates
	return resourceReadContainerRegistry(ctx, d, m)
}

// setContainerRegistryState sets all fields from the API response into the Terraform state
func setContainerRegistryState(d *schema.ResourceData, registry *goe2e.ContainerRegistry) error {
	if err := d.Set(tfconstants.AttrProjectName, registry.ProjectName); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrProjectName, err)
	}
	if err := d.Set(tfconstants.AttrPreventVulnerabilities, registry.PreventVul); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrPreventVulnerabilities, err)
	}
	if err := d.Set(tfconstants.AttrSeverity, registry.Severity); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrSeverity, err)
	}
	if err := d.Set(tfconstants.AttrStatus, registry.State); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrStatus, err)
	}
	if err := d.Set(tfconstants.AttrSetupStatus, registry.State); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrSetupStatus, err)
	}
	if err := d.Set(tfconstants.AttrDomainName, registry.DomainName); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrDomainName, err)
	}
	if err := d.Set(tfconstants.AttrProjectSize, registry.ProjectSize); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrProjectSize, err)
	}
	if err := d.Set(tfconstants.AttrStorageLimit, registry.StorageLimit); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrStorageLimit, err)
	}
	if err := d.Set(tfconstants.AttrIsPublic, registry.IsPublic); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrIsPublic, err)
	}
	if err := d.Set(tfconstants.AttrCreatedAt, registry.CreatedAt); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrCreatedAt, err)
	}
	if err := d.Set(tfconstants.AttrUpdatedAt, registry.UpdatedAt); err != nil {
		return fmt.Errorf(ErrorSetField, tfconstants.AttrUpdatedAt, err)
	}
	return nil
}

// resourceContainerRegistryV0 returns the V0 schema (before tags field was added)
func resourceContainerRegistryV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),
			tfconstants.AttrProjectName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrPreventVulnerabilities: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			tfconstants.AttrSeverity: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  tfconstants.ContainerRegistryDefaultSeverity,
			},
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrSetupStatus: {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: DeprecationMessageSetupStatus,
			},
			tfconstants.AttrDomainName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrIsPublic: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			tfconstants.AttrProjectSize: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			tfconstants.AttrStorageLimit: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			tfconstants.AttrCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrUpdatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Note: tags field did NOT exist in V0
		},
	}
}

// ResourceContainerRegistryStateUpgradeV0toV1 upgrades the state from V0 to V1
// The main change is adding the tags field as an empty map if it doesn't exist
func ResourceContainerRegistryStateUpgradeV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	// If tags field already exists (shouldn't happen in V0, but be safe), preserve it
	if _, ok := rawState[tfconstants.AttrTags]; !ok {
		// Add empty tags map
		rawState[tfconstants.AttrTags] = make(map[string]interface{})
	}

	// All other fields are preserved as-is
	return rawState, nil
}
