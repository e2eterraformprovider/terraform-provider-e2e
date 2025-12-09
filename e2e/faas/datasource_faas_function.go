package faas

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceFaasFunction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFaasFunctionRead,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// FaaS-specific fields
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the FaaS function",
			},
			e2econstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the FaaS function",
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the namespace for the FaaS function",
			},
			"runtime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the runtime for the function",
			},
			"memory_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "memory allocation in megabytes",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "function timeout in seconds",
			},
			"min_replicas": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the minimum number of replicas",
			},
			"max_replicas": {
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
			"endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the endpoint URL of the function",
			},
			e2econstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the FaaS function instance",
			},
			e2econstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the FaaS function",
			},
			e2econstants.AttrUpdatedAt: {
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

	functionID := d.Get("function_id").(string)

	log.Printf("[INFO] DATASOURCE FAAS FUNCTION READ | ID: %s", functionID)

	fn, _, err := client.FaaS.GetFunction(ctx, functionID)
	if err != nil {
		return diag.FromErr(err)
	}

	if fn == nil {
		return diag.Errorf("FaaS function with ID %s not found", functionID)
	}

	d.SetId(fn.ID)
	d.Set("name", fn.Name)
	d.Set("namespace", fn.Namespace)
	d.Set("runtime", fn.Runtime)
	d.Set("memory_mb", fn.MemoryMB)
	d.Set("timeout_seconds", fn.Timeout)
	d.Set("min_replicas", fn.MinReplicas)
	d.Set("max_replicas", fn.MaxReplicas)
	d.Set("endpoint_url", fn.EndpointURL)
	d.Set(e2econstants.AttrStatus, fn.Status)
	d.Set("created_at", fn.CreatedAt)
	d.Set(e2econstants.AttrUpdatedAt, fn.UpdatedAt)

	if fn.Environment != nil {
		d.Set("environment", fn.Environment)
	}

	return diags
}
