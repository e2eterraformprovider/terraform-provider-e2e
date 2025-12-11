package kubernetes

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceKubernetesService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadKubernetes,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Kubernetes-specific fields
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of the Kubernetes service",
				ForceNew:    true,
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the Kubernetes service",
			},
			tfconstants.AttrVersion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the version of the Kubernetes service",
				ForceNew:    true,
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Kubernetes service instance",
			},
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the Kubernetes service",
			},
			"master_node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "id of the master node of the Kubernetes cluster",
				ForceNew:    true,
			},
		},
	}
}

func dataSourceReadKubernetes(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("=============INSIDE KUBERNETES READ DATA SOURCE==========================")

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("Error getting goe2e client for project (%s), region (%s): %s", projectIDStr, region, err)
	}

	kubernetesId := d.Get("service_id").(string)
	kubernetes, _, err := goe2eClient.Kubernetes.Get(ctx, kubernetesId)
	if err != nil {
		return diag.Errorf("error finding Kubernetes cluster with ID %s: %s", kubernetesId, err)
	}

	if kubernetes == nil {
		return diag.Errorf("Kubernetes cluster with ID %s not found", kubernetesId)
	}

	log.Printf("[INFO] KUBERNETES READ | BEFORE SETTING DATA")
	d.SetId(kubernetes.ServiceID)
	d.Set("name", kubernetes.ServiceName)
	d.Set(tfconstants.AttrStatus, kubernetes.State)
	d.Set(tfconstants.AttrVersion, kubernetes.Version)
	d.Set("created_at", kubernetes.CreatedAt)
	// Note: master_node_id is not available in goe2e.KubernetesCluster struct
	// This field may need to be removed or fetched from a different endpoint
	return diags
}
