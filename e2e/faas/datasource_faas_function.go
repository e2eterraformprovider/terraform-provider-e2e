package faas

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"

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
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	functionID := d.Get("function_id").(string)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	log.Printf("[INFO] DATASOURCE FAAS FUNCTION READ | ID: %s", functionID)

	res, err := apiClient.GetFaasFunction(functionID, projectID, location)
	if err != nil {
		return diag.FromErr(err)
	}

	if res == nil {
		return diag.Errorf("FaaS function with ID %s not found", functionID)
	}

	d.SetId(res.Data.ID)
	_ = d.Set("name", res.Data.Name)
	_ = d.Set("namespace", res.Data.Namespace)
	_ = d.Set("runtime", res.Data.Runtime)
	_ = d.Set("memory_mb", res.Data.MemoryMB)
	_ = d.Set("timeout_seconds", res.Data.Timeout)
	_ = d.Set("min_replicas", res.Data.MinReplicas)
	d.Set("max_replicas", res.Data.MaxReplicas)
	d.Set("endpoint_url", res.Data.EndpointURL)
	d.Set("status", res.Data.Status)
	d.Set("created_at", res.Data.CreatedAt)
	d.Set("updated_at", res.Data.UpdatedAt)

	if res.Data.Environment != nil {
		d.Set("environment", res.Data.Environment)
	}

	return diags
}
