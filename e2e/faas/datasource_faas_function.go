package faas

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceFaasFunction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFaasFunctionRead,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// FaaS-specific fields
			tfconstants.AttrFunctionID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the FaaS function",
			},
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the FaaS function",
			},
			tfconstants.AttrNamespace: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the namespace for the FaaS function",
			},
			tfconstants.AttrRuntime: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the runtime for the function",
			},
			tfconstants.AttrMemoryMB: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "memory allocation in megabytes",
			},
			tfconstants.AttrTimeoutSeconds: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "function timeout in seconds",
			},
			tfconstants.AttrMinReplicas: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the minimum number of replicas",
			},
			tfconstants.AttrMaxReplicas: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the maximum number of replicas",
			},
			"environment": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "environment variables for the function",
			},
			tfconstants.AttrEndpointURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the endpoint URL of the function",
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the FaaS function instance",
			},
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the FaaS function",
			},
			tfconstants.AttrUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the last update date for the FaaS function",
			},
		},
	}
}

func dataSourceFaasFunctionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	functionID := d.Get(tfconstants.AttrFunctionID).(string)

	log.Printf("[INFO] DATASOURCE FAAS FUNCTION READ | ID: %s", functionID)

	fn, _, err := client.FaaS.GetFunction(ctx, functionID)
	if err != nil {
		return diag.FromErr(err)
	}

	if fn == nil {
		return diag.Errorf(ErrFunctionNotFoundByIDFmt, functionID)
	}

	d.SetId(fn.ID)
	d.Set(tfconstants.AttrName, fn.Name)
	d.Set(tfconstants.AttrNamespace, fn.Namespace)
	d.Set(tfconstants.AttrRuntime, fn.Runtime)
	d.Set(tfconstants.AttrMemoryMB, fn.MemoryMB)
	d.Set(tfconstants.AttrTimeoutSeconds, fn.Timeout)
	d.Set(tfconstants.AttrMinReplicas, fn.MinReplicas)
	d.Set(tfconstants.AttrMaxReplicas, fn.MaxReplicas)
	d.Set(tfconstants.AttrEndpointURL, fn.EndpointURL)
	d.Set(tfconstants.AttrStatus, fn.Status)
	d.Set(tfconstants.AttrCreatedAt, fn.CreatedAt)
	d.Set(tfconstants.AttrUpdatedAt, fn.UpdatedAt)

	if fn.Environment != nil {
		d.Set("environment", fn.Environment)
	}

	return diags
}
