package container_registry

import (
	"context"
	"fmt"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceContainerRegistry() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadContainerRegistry,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Container Registry ID (cr_project_id)",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID associated with the Container Registry",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location/region where the Container Registry is set up (e.g., 'Delhi')",
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the container registry project",
			},
			"setup_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Setup status of the container registry (e.g., 'CREATED')",
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Severity level for vulnerability scan ('low', 'medium', 'high', 'critical')",
			},
			"prevent_vul": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to prevent vulnerable images from being pushed",
			},
		},
	}
}

func dataSourceReadContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	id := d.Get("id").(string)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	registries, err := apiClient.GetContainerRegistryProjects(projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch container registry projects: %v", err))
	}

	for _, r := range registries {
		if strconv.Itoa(r.ID) == id {
			d.SetId(id)
			_ = d.Set("project_name", r.ProjectName)
			_ = d.Set("setup_status", r.State)
			_ = d.Set("severity", r.Severity)
			_ = d.Set("prevent_vul", r.PreventVul)
			return diags
		}
	}

	return diag.Errorf("container registry with ID %s not found", id)
}
