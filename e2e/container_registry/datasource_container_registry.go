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
			tfconstants.AttrSetupStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  DeprecationMessageSetupStatus,
				Description: DeprecationMessageSetupStatusAlternative,
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
		return diag.FromErr(fmt.Errorf(ErrorInvalidID, err))
	}

	registry, _, err := apiClient.ContainerRegistry.GetContainerRegistry(ctx, registryID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(ErrorReadRegistry, id, err))
	}

	if registry == nil {
		return diag.Errorf("container registry with ID %s not found", id)
	}

	d.SetId(id)
	if err := d.Set(tfconstants.AttrProjectName, registry.ProjectName); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorSetField, tfconstants.AttrProjectName, err))
	}
	if err := d.Set(tfconstants.AttrStatus, registry.State); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorSetField, tfconstants.AttrStatus, err))
	}
	if err := d.Set(tfconstants.AttrSetupStatus, registry.State); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorSetField, tfconstants.AttrSetupStatus, err))
	}
	if err := d.Set(tfconstants.AttrSeverity, registry.Severity); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorSetField, tfconstants.AttrSeverity, err))
	}
	if err := d.Set(tfconstants.AttrPreventVulnerabilities, registry.PreventVul); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorSetField, tfconstants.AttrPreventVulnerabilities, err))
	}

	return nil
}
