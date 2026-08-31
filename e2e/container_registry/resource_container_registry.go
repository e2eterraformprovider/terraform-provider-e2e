package container_registry

import (
	"context"
	"fmt"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
)

func ResourceContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The project ID to associate the Container Registry with.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location/region for the Container Registry.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Container Registry project.",
			},
			"prevent_vul": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to prevent vulnerable images.",
			},
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "low",
				Description: "The severity level for vulnerabilities (low, medium, high, critical).",
			},
			"setup_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of Container Registry setup.",
			},
		},
		CreateContext: resourceCreateContainerRegistry,
		ReadContext:   resourceReadContainerRegistry,
		DeleteContext: resourceDeleteContainerRegistry,
		UpdateContext: resourceUpdateContainerRegistry,
		Importer: &schema.ResourceImporter{
			State: node.CustomImportStateFunc,
		},
	}
}

func resourceCreateContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	projectName := d.Get("project_name").(string)
	preventVul := d.Get("prevent_vul").(bool)
	severity := d.Get("severity").(string)

	req := &models.CreateContainerRegistryRequest{
		ProjectName: projectName,
		PreventVul:  fmt.Sprintf("%t", preventVul),
		Severity:    severity,
	}

	_, err := apiClient.CreateContainerRegistry(req, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Container Registry: %v", err))
	}

	projects, err := apiClient.GetContainerRegistryProjects(projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch registry list post-create: %v", err))
	}

	for _, p := range projects {
		if p.ProjectName == projectName {
			d.SetId(fmt.Sprintf("%d", p.ID))
			d.Set("setup_status", p.State)
			return nil
		}
	}

	return diag.FromErr(fmt.Errorf("container registry created but not found in list"))
}

func resourceReadContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	id := d.Id()

	projects, err := apiClient.GetContainerRegistryProjects(projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch container registry projects: %v", err))
	}

	for _, p := range projects {
		if fmt.Sprintf("%d", p.ID) == id {
			d.Set("setup_status", p.State)
			return nil
		}
	}

	log.Printf("[INFO] Container Registry project with ID %s not found; removing from state", id)
	d.SetId("")
	return nil
}

func resourceDeleteContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	projectName := d.Get("project_name").(string)
	crProjectID := d.Id()
	userID := "0" // Replace with actual user ID if needed

	err := apiClient.DeleteContainerRegistry(crProjectID, projectName, userID, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Container Registry: %v", err))
	}

	d.SetId("")
	return nil
}

func resourceUpdateContainerRegistry(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	if d.HasChange("prevent_vul") || d.HasChange("severity") {
		projectID := d.Get("project_id").(string)
		location := d.Get("location").(string)
		projectName := d.Get("project_name").(string)

		preventVul := fmt.Sprintf("%t", d.Get("prevent_vul").(bool))
		severity := d.Get("severity").(string)

		err := apiClient.UpdateContainerRegistry(projectName, preventVul, severity, projectID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update container registry: %v", err))
		}
	}

	return resourceReadContainerRegistry(ctx, d, m)
}
