package autoscaling

import (
	"context"
	"fmt"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceScalerGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadScalerGroup,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Autoscaling-specific fields
			tfconstants.AttrID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the Scaler Group",
			},
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the Scaler Group",
			},
			tfconstants.AttrDesired: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the desired number of nodes",
			},
			tfconstants.AttrMinNodes: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the minimum number of nodes",
			},
			tfconstants.AttrMaxNodes: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the maximum number of nodes",
			},
			tfconstants.AttrPlan: {
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
	d.Set(tfconstants.AttrDesired, group.Desired)
	d.Set(tfconstants.AttrMinNodes, group.MinNodes)
	d.Set(tfconstants.AttrMaxNodes, group.MaxNodes)
	d.Set(tfconstants.AttrPlan, group.PlanName)
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
