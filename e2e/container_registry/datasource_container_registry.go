package container_registry

import (
	"context"
	"fmt"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceContainerRegistry() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadContainerRegistry,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Container registry-specific fields
			tfconstants.AttrID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the Container Registry",
			},
			tfconstants.AttrProjectName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the Container Registry project",
			},
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
			tfconstants.AttrSeverity: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the severity level for vulnerability scan (low, medium, high, critical)",
			},
			tfconstants.AttrPreventVulnerabilities: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether to prevent vulnerable images from being pushed",
			},
		},
	}
}

func dataSourceReadContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	apiClient := cfg.Goe2eClient()

	id := d.Get(tfconstants.AttrID).(string)

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
		return diag.Errorf("container registry with ID %s not found", id)
	}

	d.SetId(id)
	if err := d.Set(tfconstants.AttrProjectName, registry.ProjectName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set project_name: %w", err))
	}
	if err := d.Set(tfconstants.AttrStatus, registry.State); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set status: %w", err))
	}
	if err := d.Set("setup_status", registry.State); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set setup_status: %w", err))
	}
	if err := d.Set(tfconstants.AttrSeverity, registry.Severity); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set severity: %w", err))
	}
	if err := d.Set(tfconstants.AttrPreventVulnerabilities, registry.PreventVul); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set prevent_vul: %w", err))
	}

	return nil
}
