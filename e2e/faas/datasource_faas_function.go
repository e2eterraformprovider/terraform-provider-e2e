package faas

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceFaasFunction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFaasFunctionRead,
		Schema: map[string]*schema.Schema{
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the FaaS function",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location/region of the function",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the FaaS function",
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The namespace for the FaaS function",
			},
			"runtime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The runtime for the function",
			},
			"memory_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Memory allocation in MB",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Function timeout in seconds",
			},
			"min_replicas": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum number of replicas",
			},
			"max_replicas": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of replicas",
			},
			"environment": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Environment variables for the function",
			},
			"endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint URL of the function",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the function",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the function was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the function was last updated",
			},
		},
	}
}

func dataSourceFaasFunctionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config.CombinedConfig)
	client := config.NewClient
	var diags diag.Diagnostics

	functionID := d.Get("function_id").(string)
	opts := &goe2e.RequestOptions{
		ProjectID: d.Get("project_id").(string),
		Location:  d.Get("location").(string),
	}

	log.Printf("[INFO] DATASOURCE FAAS FUNCTION READ | ID: %s", functionID)

	fn, _, err := client.FaaS.GetFunction(ctx, functionID, opts)
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
	d.Set("status", fn.Status)
	d.Set("created_at", fn.CreatedAt)
	d.Set("updated_at", fn.UpdatedAt)

	if fn.Environment != nil {
		d.Set("environment", fn.Environment)
	}

	return diags
}
