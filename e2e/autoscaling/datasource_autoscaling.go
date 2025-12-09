package autoscaling

import (
	"context"
	"fmt"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceScalerGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadScalerGroup,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Autoscaling-specific fields
			e2econstants.AttrID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the Scaler Group",
			},
			e2econstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the Scaler Group",
			},
			e2econstants.AttrDesired: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the desired number of nodes",
			},
			e2econstants.AttrMinNodes: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the minimum number of nodes",
			},
			e2econstants.AttrMaxNodes: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the maximum number of nodes",
			},
			e2econstants.AttrPlan: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the plan of the Scaler Group",
			},
			"vm_image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the VM image name used",
			},
			"provision_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the provision status of the Scaler Group (e.g., Running, Stopped)",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the policy type (elastic, scheduled, etc.)",
			},
			"policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":           {Type: schema.TypeString, Computed: true},
						"adjust":         {Type: schema.TypeInt, Computed: true},
						"parameter":      {Type: schema.TypeString, Computed: true},
						"operator":       {Type: schema.TypeString, Computed: true},
						"value":          {Type: schema.TypeString, Computed: true},
						"period_number":  {Type: schema.TypeString, Computed: true},
						"period_seconds": {Type: schema.TypeString, Computed: true},
						"cooldown":       {Type: schema.TypeString, Computed: true},
					},
				},
				Description: "list of elastic scaling policies (upscale and downscale)",
			},
			"scheduled_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":       {Type: schema.TypeString, Computed: true},
						"adjust":     {Type: schema.TypeString, Computed: true},
						"recurrence": {Type: schema.TypeString, Computed: true},
					},
				},
				Description: "list of scheduled scaling policies (upscale and downscale)",
			},
		},
	}
}

func dataSourceReadScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	scalerID := d.Get("id").(string)

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))
	}

	// Get scaler group using GoE2E client
	group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, scalerID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read scaler group: %w", err))
	}

	if group == nil {
		return diag.Errorf("scaler group with ID %s not found", scalerID)
	}

	d.SetId(group.ID)
	d.Set("name", group.Name)
	d.Set(e2econstants.AttrDesired, group.Desired)
	d.Set(e2econstants.AttrMinNodes, group.MinNodes)
	d.Set(e2econstants.AttrMaxNodes, group.MaxNodes)
	d.Set(e2econstants.AttrPlan, group.PlanName)
	d.Set("vm_image_name", group.VMImageName)
	d.Set("provision_status", NormalizeStatus(group.ProvisionStatus))
	d.Set("policy_type", group.PolicyType)

	// Set policy list (V2 format for datasource)
	policyList := []map[string]interface{}{
		{
			"type":           group.PolicyType,
			"adjust":         group.UpscalePolicyValue,
			"parameter":      group.PolicyMeasure,
			"operator":       group.PolicyUpscaleOperator,
			"value":          strconv.Itoa(group.UpscalePolicyValue),
			"period_number":  strconv.Itoa(group.WaitPeriod),
			"period_seconds": strconv.Itoa(group.Cooldown),
			"cooldown":       strconv.Itoa(group.Cooldown),
		},
		{
			"type":           group.PolicyType,
			"adjust":         group.DownscalePolicyValue,
			"parameter":      group.PolicyMeasure,
			"operator":       group.PolicyDownscaleOperator,
			"value":          strconv.Itoa(group.DownscalePolicyValue),
			"period_number":  strconv.Itoa(group.WaitPeriod),
			"period_seconds": strconv.Itoa(group.Cooldown),
			"cooldown":       strconv.Itoa(group.Cooldown),
		},
	}
	d.Set("policy", policyList)

	scheduledPolicyList := []map[string]interface{}{
		{
			"type":       group.ScheduledPolicyOp,
			"adjust":     strconv.Itoa(group.UpscaleAdjust),
			"recurrence": group.UpscaleRecurrence,
		},
		{
			"type":       group.ScheduledPolicyOp,
			"adjust":     strconv.Itoa(group.DownscaleAdjust),
			"recurrence": group.DownscaleRecurrence,
		},
	}
	d.Set("scheduled_policy", scheduledPolicyList)

	return diags
}
