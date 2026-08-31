package autoscaling

import (
	"context"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceScalerGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadScalerGroup,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the scaler group",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID associated with the scaler group",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location of the scaler group",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the scaler group",
			},
			"desired": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Desired node count",
			},
			"min_nodes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum nodes allowed",
			},
			"max_nodes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum nodes allowed",
			},
			"plan_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Plan name for the scaler group",
			},
			"vm_image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VM image name used",
			},
			"provision_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provision status (e.g., Running, Stopped)",
			},

			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy type (elastic, scheduled, etc.)",
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
				Description: "Elastic scaling policies (upscale and downscale).",
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
				Description: "Scheduled scaling policies (upscale and downscale).",
			},
		},
	}
}

func dataSourceReadScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	scalerID := d.Get("id").(string)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	group, err := apiClient.GetScalerGroup(scalerID, projectID, location)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(group.ID))
	_ = d.Set("name", group.Name)
	_ = d.Set("desired", group.Desired)
	_ = d.Set("min_nodes", group.MinNodes)
	_ = d.Set("max_nodes", group.MaxNodes)
	_ = d.Set("plan_name", group.PlanName)
	_ = d.Set("vm_image_name", group.VMImageName)
	_ = d.Set("provision_status", group.ProvisionStatus)
	_ = d.Set("policy_type", group.PolicyType)

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
	_ = d.Set("policy", policyList)

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
	_ = d.Set("scheduled_policy", scheduledPolicyList)

	return diags
}
