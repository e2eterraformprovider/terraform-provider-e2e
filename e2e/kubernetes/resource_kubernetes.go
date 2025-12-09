package kubernetes

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceKubernetesService() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceKubernetesResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceKubernetesStateUpgradeV0toV1,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			e2econstants.AttrRegion: config.RegionSchema(),
			e2econstants.AttrLocation: {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Use 'region' instead. The 'location' field will be removed in v4.0.0",
				ConflictsWith: []string{e2econstants.AttrRegion},
				Description:   "the location of the cluster (deprecated, use 'region')",
			},
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// V3 PREFERRED CLUSTER FIELDS (Aliases)
			// ============================================
			"cluster_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "name of the Kubernetes cluster (preferred over 'name')",
				ConflictsWith: []string{e2econstants.AttrName},
				ValidateFunc:  validation.StringLenBetween(1, 255),
			},
			"kubernetes_version": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "the Kubernetes version (preferred over 'version')",
				ConflictsWith: []string{e2econstants.AttrVersion},
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^1\.\d{2}$`), "must be format 1.XX"),
			},

			// ============================================
			// DEPRECATED CLUSTER FIELDS (V2)
			// ============================================
			e2econstants.AttrName: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Deprecated:    "Use 'cluster_name' instead for consistency with other providers. The 'name' field will be removed in v4.0.0",
				ConflictsWith: []string{"cluster_name"},
				Description:   "name of the Kubernetes cluster (deprecated, use 'cluster_name')",
			},
			e2econstants.AttrVersion: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Deprecated:    "Use 'kubernetes_version' instead for clarity. The 'version' field will be removed in v4.0.0",
				ConflictsWith: []string{"kubernetes_version"},
				Description:   "the Kubernetes version (deprecated, use 'kubernetes_version')",
			},
			e2econstants.AttrVPCID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of the VPC for the Kubernetes cluster",
			},

			// ============================================
			// V3 NEW FEATURES
			// ============================================
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "tags to apply to the Kubernetes cluster (state-only, not sent to API)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "list of security group IDs to attach to the cluster",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"encryption_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether encryption is enabled for the cluster",
				ForceNew:    true,
			},
			"encryption_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "passphrase for cluster encryption (required if encryption_enabled is true)",
				ForceNew:    true,
			},

			// ============================================
			// REQUIRED INPUT FIELDS (Mutable via Updates)
			// ============================================
			e2econstants.AttrNodePools: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "list of node pools for the Kubernetes cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Node pool immutable fields
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "name of the node pool",
						},

						// V3 PREFERRED NODE POOL FIELDS (Aliases)
						"plan": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Description:   "the plan/SKU for worker nodes (preferred over 'specs_name')",
							ConflictsWith: []string{"specs_name"},
						},
						"type": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Description:   "the node pool type: Static or Autoscale (preferred over 'node_pool_type')",
							ConflictsWith: []string{"node_pool_type"},
							ValidateFunc:  validation.StringInSlice([]string{"Static", "Autoscale"}, false),
						},
						"size": {
							Type:          schema.TypeInt,
							Optional:      true,
							Description:   "the number of worker nodes for Static pools (preferred over 'worker_node')",
							ConflictsWith: []string{"worker_node"},
							ValidateFunc:  validation.IntBetween(2, 25),
						},
						"min_nodes": {
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       0,
							Description:   "minimum number of nodes for Autoscale pools (preferred over 'min_vms')",
							ConflictsWith: []string{e2econstants.AttrMinVMs},
							ValidateFunc:  validation.All(validation.IntAtLeast(2), validation.IntAtMost(25)),
						},
						"max_nodes": {
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       0,
							Description:   "maximum number of nodes for Autoscale pools (preferred over 'max_vms')",
							ConflictsWith: []string{e2econstants.AttrMaxVMs},
							ValidateFunc:  validation.IntAtMost(25),
						},

						// DEPRECATED NODE POOL FIELDS (V2)
						"specs_name": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Deprecated:    "Use 'plan' instead. Will be removed in v5.0.0",
							ConflictsWith: []string{"plan"},
							Description:   "the plan/SKU for worker nodes (deprecated, use 'plan')",
						},
						"node_pool_type": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							Deprecated:    "Use 'type' instead. Will be removed in v5.0.0",
							ConflictsWith: []string{"type"},
							Description:   "the node pool type (deprecated, use 'type')",
							ValidateFunc: validation.StringInSlice([]string{
								"Static",
								"Autoscale",
							}, false),
						},
						"worker_node": {
							Type:          schema.TypeInt,
							Optional:      true,
							Deprecated:    "Use 'size' instead. Will be removed in v5.0.0",
							ConflictsWith: []string{"size"},
							Description:   "number of worker nodes (deprecated, use 'size')",
							ValidateFunc:  validation.IntBetween(2, 25),
						},
						e2econstants.AttrMinVMs: {
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       0,
							Deprecated:    "Use 'min_nodes' instead. Will be removed in v5.0.0",
							ConflictsWith: []string{"min_nodes"},
							ValidateFunc:  validation.All(validation.IntAtLeast(2), validation.IntAtMost(25)),
							Description:   "the minimum number of virtual machines (Autoscale pools only, deprecated, use 'min_nodes')",
						},
						e2econstants.AttrMaxVMs: {
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       0,
							Deprecated:    "Use 'max_nodes' instead. Will be removed in v5.0.0",
							ConflictsWith: []string{"max_nodes"},
							ValidateFunc:  validation.IntAtMost(25),
							Description:   "the maximum number of virtual machines (Autoscale pools only, deprecated, use 'max_nodes')",
						},
						"node_pool_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(2),
							Description:  "the target size for resizing the node pool (must be at least 2)",
						},

						// Node pool computed fields
						"cardinality": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the current number of nodes in the pool",
						},
						"slug_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the slug name of the node pool plan",
						},
						"sku_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the SKU id of the node pool",
						},
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "id of the service for the node pool",
						},
						"elasticity_dict": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "the elasticity dictionary for the worker node pool",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"worker": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "worker settings in the elasticity dictionary",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"period_number": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "the period number",
												},
												"policy_paramter_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "the policy parameter type (Default or Custom; if Custom, you must provide the parameter field)",
													ValidateFunc: validation.StringInSlice([]string{
														"Default",
														"Custom",
													}, false),
												},
												"parameter": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "CPU",
													Description: "the parameter (e.g., CPU, Memory)",
													ValidateFunc: validation.Any(
														validation.StringInSlice([]string{"Memory", "CPU"}, false),
														validation.StringMatch(
															regexp.MustCompile(`^[A-Z0-9]([_]?[A-Z0-9])+$`),
															"Parameter Name should be at least 2 characters long with upper case characters, numbers and underscore and must be start and end with characters or numbers.",
														),
													),
												},
												"elasticity_policies": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "list of elasticity policies",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "the type (fixed value: CHANGE)",
															},
															"adjust": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "the adjust value (1 or -1)",
															},
															"operator": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "the operator for adding worker (e.g., >, >=)",
															},
															"value": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "the value for adding worker",
															},
															"period": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "the period",
															},
															"watch_period": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "the period number",
															},
															"cooldown": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "the cooldown period",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"scheduled_dict": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "the scheduled dictionary for the worker node pool",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"worker": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "worker settings in the scheduled dictionary",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scheduled_policies": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"upscale_cardinality": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "the cardinality for upscaling",
															},
															"upscale_recurrence": {
																Type:         schema.TypeString,
																Required:     true,
																Description:  "the recurrence timing for upscaling",
																ValidateFunc: validation.StringInSlice([]string{"0 12 * * *", "0 0 1 * *", "0 20 * * *", "0 9 * * 1-5", "0 9-13 * * *"}, false),
															},
															"downscale_cardinality": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "the cardinality for downscaling",
															},
															"downscale_recurrence": {
																Type:         schema.TypeString,
																Required:     true,
																Description:  "the recurrence timing for downscaling",
																ValidateFunc: validation.StringInSlice([]string{"0 2 * * *", "0 0 15 * *", "30 5 * * 1-5", "0 0 * * 6,7", "0 0 12 1 1"}, false),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the policy type for the worker node pool",
						},
						"custom_param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the custom parameter name for the worker node pool",
						},
						"custom_param_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the custom parameter value for the worker node pool",
						},
					},
				},
			},

			// ============================================
			// COMPUTED FIELDS - CLUSTER INFO
			// ============================================
			"slug_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the slug name of the Kubernetes cluster plan",
			},
			"sku_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the SKU id of the Kubernetes cluster",
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			e2econstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Kubernetes cluster instance",
			},
			e2econstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation timestamp of the Kubernetes cluster",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		CreateContext: resourceCreateKubernetesService,
		ReadContext:   resourceReadKubernetesService,
		UpdateContext: resourceUpdateKubernetesService,
		DeleteContext: resourceDeleteKubernetesService,
		Exists:        resourceExistsKubernetesService,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKubernetesImport,
		},
	}
}

func GetSlugName(ctx context.Context, d *schema.ResourceData, cfg *config.Config) (string, error) {
	log.Printf("[INFO] KUBERNETES PLAN READ STARTS")
	version := getKubernetesVersion(d)
	if version == "" {
		return "", fmt.Errorf("kubernetes_version or version is required")
	}

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return "", err
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return "", err
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return "", fmt.Errorf("error getting goe2e client for project (%s), region (%s): %s", projectIDStr, region, err)
	}

	log.Printf("--------------MAKING API CALL FOR SLUGNAME-------------")
	plans, _, err := goe2eClient.Kubernetes.GetMasterPlans(ctx)
	if err != nil {
		return "", fmt.Errorf("error retrieving Kubernetes plans for project (%s), region (%s): %s", projectIDStr, region, err.Error())
	}

	// Find plan matching the version
	for _, plan := range plans {
		if plan.K8sVersion == version {
			if err := d.Set("sku_id", plan.Specs.ID); err != nil {
				log.Printf("[WARN] Failed to set sku_id: %s", err)
			}
			return plan.Plan, nil
		}
	}

	return "", fmt.Errorf("Kubernetes plan not found for version %s in project (%s), region (%s)", version, projectIDStr, region)
}

func CreateKubernetesObject(ctx context.Context, cfg *config.Config, d *schema.ResourceData, slugName string, goe2eClient *goe2e.Client) (*goe2e.KubernetesClusterCreateRequest, diag.Diagnostics) {
	log.Printf("[INFO] KUBERNETES OBJECT CREATION STARTS")

	clusterName := getClusterName(d)
	kubernetesVersion := getKubernetesVersion(d)
	vpcID := d.Get(e2econstants.AttrVPCID).(string)
	skuID := d.Get("sku_id").(string)

	kubernetesObj := &goe2e.KubernetesClusterCreateRequest{
		Name:     clusterName,
		Version:  kubernetesVersion,
		VPCID:    vpcID,
		SKUID:    skuID,
		SlugName: slugName,
	}

	if nodePools, ok := d.GetOk(e2econstants.AttrNodePools); ok {
		nodePoolList := nodePools.([]interface{})

		// Use config helpers instead of direct Get()
		projectIDStr, err := cfg.GetProjectIDOrDefault(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		region, err := cfg.GetRegionOrDefault(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		nodePoolsDetail, err := ExpandNodePools(ctx, nodePoolList, goe2eClient, projectIDStr, region)
		if err != nil {
			return nil, diag.Errorf("error preparing node pools for Kubernetes cluster in project (%s), region (%s): %s", projectIDStr, region, err)
		}
		kubernetesObj.NodePools = nodePoolsDetail
	} else {
		kubernetesObj.NodePools = make([]goe2e.NodePool, 0)
	}

	return kubernetesObj, nil
}

func resourceCreateKubernetesService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("Error getting goe2e client for project (%s), region (%s): %s", projectIDStr, region, err)
	}

	// Log deprecation warnings
	logDeprecationWarning(d, e2econstants.AttrName, "cluster_name")
	logDeprecationWarning(d, e2econstants.AttrVersion, "kubernetes_version")
	logDeprecationWarning(d, e2econstants.AttrLocation, e2econstants.AttrRegion)

	clusterName := getClusterName(d)

	slugName, err := GetSlugName(ctx, d, cfg)
	if err != nil {
		return diag.Errorf("Error retrieving Kubernetes plan slug name for cluster (name: %s) in project (%s), region (%s): %s", clusterName, projectIDStr, region, err)
	}
	if err := d.Set("slug_name", slugName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting slug_name: %w", err))
	}

	kubernetesObject, diags := CreateKubernetesObject(ctx, cfg, d, slugName, goe2eClient)
	if diags != nil {
		return diags
	}

	log.Printf("---------KUBERNETES OBJECT CREATED---------: %+v", kubernetesObject)

	// Create cluster via goe2e
	cluster, _, err := goe2eClient.Kubernetes.Create(ctx, kubernetesObject)
	if err != nil {
		return diag.Errorf("Error creating Kubernetes cluster (name: %s) in project (%s), region (%s): %s", clusterName, projectIDStr, region, err)
	}

	if cluster == nil {
		return diag.Errorf("Cluster creation returned nil response")
	}

	// Set ID
	d.SetId(cluster.ServiceID)

	// Store tags in state (state-only, not sent to API)
	if tags, ok := d.GetOk("tags"); ok {
		if err := d.Set("tags", tags); err != nil {
			log.Printf("[WARN] Failed to set tags: %s", err)
		}
	}

	// Wait for cluster to become Running (async operation, 30 min timeout)
	err = waitForClusterStatus(ctx, goe2eClient, cluster.ServiceID, "Running", 30*time.Minute)
	if err != nil {
		return diag.Errorf("Error waiting for cluster (ID: %s) to become Running: %s", cluster.ServiceID, err)
	}

	// Attach security groups if specified
	if securityGroupIDs, ok := d.GetOk("security_group_ids"); ok {
		sgIDs := expandSecurityGroupIDs(securityGroupIDs.([]interface{}))
		attachReq := &goe2e.AttachSecurityGroupRequest{
			SecurityGroupIDs: sgIDs,
		}
		_, err = goe2eClient.Kubernetes.AttachSecurityGroups(ctx, cluster.ServiceID, attachReq)
		if err != nil {
			return diag.Errorf("Error attaching security groups to cluster (ID: %s): %s", cluster.ServiceID, err)
		}
	}

	// Read to populate all computed fields
	return resourceReadKubernetesService(ctx, d, m)
}

func resourceReadKubernetesService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("=============INSIDE KUBERNETES READ RESOURCE==========================")
	clusterID := d.Id()

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("Error getting goe2e client for project (%s), region (%s): %s", projectIDStr, region, err)
	}

	// Get cluster details
	cluster, _, err := goe2eClient.Kubernetes.Get(ctx, clusterID)
	if err != nil {
		return diag.Errorf("Error reading Kubernetes cluster (ID: %s): %s", clusterID, err)
	}

	if cluster == nil {
		log.Printf("[WARN] Kubernetes cluster (ID: %s) returned nil, removing from state", clusterID)
		d.SetId("")
		return diags
	}

	log.Printf("[INFO] KUBERNETES READ | BEFORE SETTING DATA")
	log.Printf("[INFO] SETTING--------- (1)")

	// Set computed fields
	if err := d.Set(e2econstants.AttrStatus, cluster.State); err != nil {
		log.Printf("[WARN] Failed to set status: %s", err)
	}
	if err := d.Set(e2econstants.AttrCreatedAt, cluster.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set created_at: %s", err)
	}
	if err := d.Set("slug_name", d.Get("slug_name")); err != nil {
		log.Printf("[WARN] Failed to set slug_name: %s", err)
	}
	if err := d.Set("sku_id", d.Get("sku_id")); err != nil {
		log.Printf("[WARN] Failed to set sku_id: %s", err)
	}

	// Set preferred fields (always use V3 names on read)
	if err := d.Set("cluster_name", cluster.ServiceName); err != nil {
		log.Printf("[WARN] Failed to set cluster_name: %s", err)
	}
	if err := d.Set("kubernetes_version", cluster.Version); err != nil {
		log.Printf("[WARN] Failed to set kubernetes_version: %s", err)
	}
	if err := d.Set(e2econstants.AttrRegion, region); err != nil {
		log.Printf("[WARN] Failed to set region: %s", err)
	}
	if err := d.Set(e2econstants.AttrVPCID, cluster.VPCID); err != nil {
		log.Printf("[WARN] Failed to set vpc_id: %s", err)
	}

	// Get node pools
	nodePools, _, err := goe2eClient.Kubernetes.GetNodePools(ctx, clusterID)
	if err != nil {
		return diag.Errorf("Error reading node pools for cluster (ID: %s): %s", clusterID, err)
	}

	// Flatten node pools
	if err := d.Set(e2econstants.AttrNodePools, flattenNodePools(nodePools)); err != nil {
		return diag.Errorf("Error setting node_pools: %s", err)
	}

	// Get security groups
	securityGroups, _, err := goe2eClient.Kubernetes.ListAttachedSecurityGroups(ctx, clusterID)
	if err == nil && len(securityGroups) > 0 {
		sgIDs := make([]int, len(securityGroups))
		for i, sg := range securityGroups {
			sgIDs[i] = sg.ID
		}
		if err := d.Set("security_group_ids", sgIDs); err != nil {
			log.Printf("[WARN] Failed to set security_group_ids: %s", err)
		}
	}

	return diags
}

func resourceDeleteKubernetesService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	clusterID := d.Id()

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("Error getting goe2e client for project (%s), region (%s): %s", projectIDStr, region, err)
	}

	status := d.Get("status").(string)
	invalidStates := []string{"Deleting", "Deleted", "Failed", "Error"}
	for _, s := range invalidStates {
		if status == s {
			return diag.Errorf("cannot delete Kubernetes cluster (ID: %s): cluster is in %s state in project (%s), region (%s)", clusterID, status, projectIDStr, region)
		}
	}

	_, err = goe2eClient.Kubernetes.Delete(ctx, clusterID)
	if err != nil {
		return diag.Errorf("Error deleting Kubernetes cluster (ID: %s) in project (%s), region (%s): %s", clusterID, projectIDStr, region, err)
	}

	d.SetId("")
	return diags
}

func resourceExistsKubernetesService(d *schema.ResourceData, m interface{}) (bool, error) {
	cfg := m.(*config.Config)
	ctx := context.Background()

	clusterID := d.Id()

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return false, err
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return false, err
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return false, fmt.Errorf("error getting goe2e client: %s", err)
	}

	cluster, resp, err := goe2eClient.Kubernetes.Get(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}

	return cluster != nil, nil
}

func resourceUpdateKubernetesService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	clusterID := d.Id()

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf("Error getting goe2e client for project (%s), region (%s): %s", projectIDStr, region, err)
	}

	status := d.Get("status").(string)
	invalidStates := []string{"Deleting", "Deleted", "Failed", "Error"}
	for _, s := range invalidStates {
		if status == s {
			return diag.Errorf("cannot update Kubernetes cluster (ID: %s): cluster is in %s state in project (%s), region (%s)", clusterID, status, projectIDStr, region)
		}
	}

	// Handle tags update (state-only, no API call)
	if d.HasChange("tags") {
		if tags, ok := d.GetOk("tags"); ok {
			if err := d.Set("tags", tags); err != nil {
				log.Printf("[WARN] Failed to set tags: %s", err)
			}
		}
	}

	// Handle security group updates
	if d.HasChange("security_group_ids") {
		old, new := d.GetChange("security_group_ids")
		oldSGIDs := expandSecurityGroupIDs(old.([]interface{}))
		newSGIDs := expandSecurityGroupIDs(new.([]interface{}))

		// Detach removed security groups
		toDetach := difference(oldSGIDs, newSGIDs)
		if len(toDetach) > 0 {
			detachReq := &goe2e.DetachSecurityGroupRequest{
				SecurityGroupIDs: toDetach,
			}
			_, err = goe2eClient.Kubernetes.DetachSecurityGroups(ctx, clusterID, detachReq)
			if err != nil {
				return diag.Errorf("Error detaching security groups from cluster (ID: %s): %s", clusterID, err)
			}
		}

		// Attach new security groups
		toAttach := difference(newSGIDs, oldSGIDs)
		if len(toAttach) > 0 {
			attachReq := &goe2e.AttachSecurityGroupRequest{
				SecurityGroupIDs: toAttach,
			}
			_, err = goe2eClient.Kubernetes.AttachSecurityGroups(ctx, clusterID, attachReq)
			if err != nil {
				return diag.Errorf("Error attaching security groups to cluster (ID: %s): %s", clusterID, err)
			}
		}
	}

	serviceMapping, err := GetNodePoolServiceMapping(ctx, d, m)
	if err != nil {
		return diag.Errorf("Error retrieving node pool service mapping for Kubernetes cluster (ID: %s) in project (%s), region (%s): %s", clusterID, projectIDStr, region, err)
	}
	if d.HasChange(e2econstants.AttrNodePools) {
		oldData, newData := d.GetChange(e2econstants.AttrNodePools)

		oldNodePools := oldData.([]interface{})
		newNodePools := newData.([]interface{})

		for _, oldNodePool := range oldNodePools {
			oldNodePoolMap := oldNodePool.(map[string]interface{})
			oldNPName := oldNodePoolMap["name"].(string)
			oldServiceFind := serviceMapping[oldNPName]
			if oldServiceFind == nil {
				return diag.Errorf("Cannot delete node pool '%s' from Kubernetes cluster (ID: %s): node pool does not exist in project (%s), region (%s)", oldNPName, clusterID, projectIDStr, region)
			}
			oldServiceIDFloat := oldServiceFind.(float64)
			oldServiceID := fmt.Sprintf("%.0f", oldServiceIDFloat)
			found := false
			if len(newNodePools) <= 0 {
				return diag.Errorf("Cannot delete node pool from Kubernetes cluster (ID: %s): at least one node pool must be present in a Kubernetes cluster", clusterID)
			}
			// Check if the old service_id exists in the new node pools
			for _, newNodePool := range newNodePools {
				newNodePoolMap := newNodePool.(map[string]interface{})
				newNPName := newNodePoolMap["name"].(string)
				if oldNPName == newNPName {
					found = true
					break
				}
			}
			if !found {
				nodePools, _, err := goe2eClient.Kubernetes.CheckNodePoolStatus(ctx, clusterID)
				if err != nil {
					return diag.Errorf("Error retrieving Kubernetes cluster (ID: %s) status in project (%s), region (%s): %s", clusterID, projectIDStr, region, err)
				}
				if !IsNodePoolRunning(oldServiceIDFloat, nodePools) {
					d.Set(e2econstants.AttrNodePools, oldData)
					return diag.Errorf("Cannot delete node pool '%s' from Kubernetes cluster (ID: %s): node pool must be in Running state before deletion", oldNPName, clusterID)
				}
				_, err = goe2eClient.Kubernetes.DeleteNodePool(ctx, oldServiceID)
				if err != nil {
					return diag.Errorf("Error deleting node pool '%s' from Kubernetes cluster (ID: %s) in project (%s), region (%s): %s", oldNPName, clusterID, projectIDStr, region, err)
				}
			}
		}

		for i := range newNodePools {
			newNodePoolMap := newNodePools[i].(map[string]interface{})
			newNPName := newNodePoolMap["name"].(string)
			found := false
			log.Printf("----------------------CHECKING IF THERE IS ANY ADDITION OF NODE POOLS-------------------")
			for _, oldNodePool := range oldNodePools {
				oldNodePoolMap := oldNodePool.(map[string]interface{})
				oldNPName := oldNodePoolMap["name"].(string)
				oldServiceIDFloat := serviceMapping[oldNPName].(float64)
				oldServiceID := fmt.Sprintf("%.0f", oldServiceIDFloat)
				// If exists then check if there is any change in cardinality
				if newNPName == oldNPName {
					found = true
					oldCardinality := oldNodePoolMap["cardinality"].(int)
					oldPoolType := getNodePoolType(oldNodePoolMap)
					if oldPoolType == "Static" {
						oldCardinality = getNodePoolSize(oldNodePoolMap)
					}
					node_pool_size := oldCardinality
					if newNodePoolMap["node_pool_size"].(int) != 0 {
						node_pool_size = newNodePoolMap["node_pool_size"].(int)
					}
					log.Printf("----------------PREV CARD:%+v     NEW CARD:%+v------------------", oldCardinality, node_pool_size)
					if node_pool_size < 2 {
						return diag.Errorf("Cannot update node pool '%s' in Kubernetes cluster (ID: %s): node_pool_size must be at least 2 (current value: %d)", newNPName, clusterID, node_pool_size)
					}
					if oldCardinality != node_pool_size {
						resizeReq := &goe2e.NodePoolResizeRequest{
							NodePoolSize: node_pool_size,
						}
						newNodePoolMap["cardinality"] = node_pool_size
						_, err := goe2eClient.Kubernetes.UpdateNodePoolCardinality(ctx, oldServiceID, resizeReq)
						if err != nil {
							return diag.Errorf("Error updating node pool '%s' cardinality in Kubernetes cluster (ID: %s), project (%s), region (%s): %s", newNPName, clusterID, projectIDStr, region, err)
						}
						break
					}
					new_node_pool_type := getNodePoolType(newNodePoolMap)
					// You cannot change the node pool type from Static to Autoscale and vice versa
					if oldPoolType != new_node_pool_type {
						return diag.Errorf("Cannot update node pool type for node pool '%s' in Kubernetes cluster (ID: %s): this field is immutable after node pool creation", newNPName, clusterID)
					}
					if new_node_pool_type == "Static" {
						break
					}
					nodePoolObject, err := ExpandNodePoolUpdate(ctx, newNodePoolMap, goe2eClient, projectIDStr, region)
					if err != nil {
						return diag.Errorf("Error preparing node pool '%s' update for Kubernetes cluster (ID: %s) in project (%s), region (%s): %s", newNPName, clusterID, projectIDStr, region, err)
					}
					updateReq := &goe2e.NodePoolUpdateRequest{
						MinVms:           nodePoolObject.MinVms,
						Cardinality:      nodePoolObject.Cardinality,
						MaxVms:           nodePoolObject.MaxVms,
						PlanID:           nodePoolObject.PlanID,
						ElasticityPolicy: nodePoolObject.ElasticityPolicy,
						ScheduledPolicy:  nodePoolObject.ScheduledPolicy,
						PolicyType:       nodePoolObject.PolicyType,
						CustomParamName:  nodePoolObject.CustomParamName,
						CustomParamValue: nodePoolObject.CustomParamValue,
					}
					_, _, err = goe2eClient.Kubernetes.UpdateNodePoolDetails(ctx, oldServiceID, updateReq)
					if err != nil {
						return diag.Errorf("Error updating node pool '%s' details in Kubernetes cluster (ID: %s), project (%s), region (%s): %s", newNPName, clusterID, projectIDStr, region, err)
					}
					break
				}
			}
			//If not found meaning this is a new NodePool
			if !found {
				var nodePoolList []interface{}
				nodePoolList = append(nodePoolList, newNodePools[i])
				nodePoolsDetail, err := ExpandNodePools(ctx, nodePoolList, goe2eClient, projectIDStr, region)
				if err != nil {
					return diag.Errorf("Error preparing new node pool '%s' for Kubernetes cluster (ID: %s) in project (%s), region (%s): %s", newNPName, clusterID, projectIDStr, region, err)
				}
				addReq := &goe2e.NodePoolAddRequest{
					NodePools: nodePoolsDetail,
				}
				log.Printf("----------------------ADDING A NEW NODE POOL-------------------")
				_, _, err = goe2eClient.Kubernetes.AddNodePool(ctx, clusterID, addReq)
				if err != nil {
					return diag.Errorf("Error adding node pool '%s' to Kubernetes cluster (ID: %s) in project (%s), region (%s): %s", newNPName, clusterID, projectIDStr, region, err)
				}
				continue
			}
		}
	}

	return resourceReadKubernetesService(ctx, d, m)
}

func GetNodePoolServiceMapping(ctx context.Context, d *schema.ResourceData, m interface{}) (map[string]interface{}, error) {
	cfg := m.(*config.Config)
	log.Printf("[INFO] KUBERNETES CLUSTER NODE POOLS MAPPING STARTS")
	clusterID := d.Id()

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return nil, err
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return nil, err
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return nil, fmt.Errorf("error getting goe2e client for project (%s), region (%s): %s", projectIDStr, region, err)
	}

	// Initialize the map to store service_name and service_id mappings
	serviceMapping := make(map[string]interface{})
	nodePools, _, err := goe2eClient.Kubernetes.GetNodePools(ctx, clusterID)
	if err != nil {
		return serviceMapping, fmt.Errorf("error retrieving node pool list for Kubernetes cluster (ID: %s) in project (%s), region (%s): %s", clusterID, projectIDStr, region, err.Error())
	}
	// Extract service_name and service_id from each node pool
	for _, nodePool := range nodePools {
		serviceMapping[nodePool.ServiceName] = nodePool.ServiceID
	}

	return serviceMapping, nil
}

func IsNodePoolRunning(oldServiceID float64, nodePools []goe2e.NodePoolServiceInfo) bool {
	for _, nodepool := range nodePools {
		if nodepool.ServiceID == oldServiceID {
			if nodepool.State == "Running" {
				return true
			}
		}
	}
	return false
}

// resourceKubernetesImport handles importing Kubernetes clusters
// Supports two import formats:
// 1. Simple: "cluster_id" (uses provider defaults for project_id and region)
// 2. Full: "project_id:region:cluster_id"
func resourceKubernetesImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	var projectID, region, clusterID string

	// Support two import formats:
	// 1. Simple: cluster_id (uses provider defaults for project_id and region)
	// 2. Full: project_id:region:cluster_id
	if len(parts) == 1 {
		// Simple format: just cluster ID
		clusterID = parts[0]
		cfg := m.(*config.Config)
		var err error
		region, err = cfg.GetRegionOrDefault(d)
		if err != nil {
			return nil, fmt.Errorf("error getting region: %w", err)
		}
		var projectIDStr string
		projectIDStr, err = cfg.GetProjectIDOrDefault(d)
		if err != nil {
			return nil, fmt.Errorf("error getting project ID: %w", err)
		}
		projectID = projectIDStr
	} else if len(parts) == 3 {
		// Full format: project_id:region:cluster_id
		projectID = parts[0]
		region = parts[1]
		clusterID = parts[2]
	} else {
		return nil, fmt.Errorf("invalid import ID format: expected 'cluster_id' or 'project_id:region:cluster_id', got '%s'", d.Id())
	}

	// Validate that clusterID is not empty
	if clusterID == "" {
		return nil, fmt.Errorf("cluster_id cannot be empty")
	}

	cfg := m.(*config.Config)

	// Set project_id and region in resource data
	if err := d.Set(e2econstants.AttrProjectID, projectID); err != nil {
		return nil, fmt.Errorf("error setting project_id: %w", err)
	}
	if err := d.Set(e2econstants.AttrRegion, region); err != nil {
		return nil, fmt.Errorf("error setting region: %w", err)
	}
	// Also set location for backwards compatibility
	if err := d.Set(e2econstants.AttrLocation, region); err != nil {
		return nil, fmt.Errorf("error setting location: %w", err)
	}

	// Get goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return nil, fmt.Errorf("error creating goe2e client during import: %w", err)
	}

	// Fetch cluster via goe2eClient.Kubernetes.Get(ctx, clusterID)
	cluster, _, err := goe2eClient.Kubernetes.Get(ctx, clusterID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Kubernetes cluster (ID: %s) during import: %w", clusterID, err)
	}

	if cluster == nil {
		return nil, fmt.Errorf("Kubernetes cluster (ID: %s) not found", clusterID)
	}

	// Fetch node pools via goe2eClient.Kubernetes.GetNodePools(ctx, clusterID)
	nodePools, _, err := goe2eClient.Kubernetes.GetNodePools(ctx, clusterID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving node pools for cluster (ID: %s) during import: %w", clusterID, err)
	}

	// Set cluster ID
	d.SetId(clusterID)

	// Populate V3 preferred fields (cluster_name, kubernetes_version, etc.)
	if err := d.Set("cluster_name", cluster.ServiceName); err != nil {
		return nil, fmt.Errorf("error setting cluster_name: %w", err)
	}
	if err := d.Set("kubernetes_version", cluster.Version); err != nil {
		return nil, fmt.Errorf("error setting kubernetes_version: %w", err)
	}
	if err := d.Set(e2econstants.AttrVPCID, cluster.VPCID); err != nil {
		return nil, fmt.Errorf("error setting vpc_id: %w", err)
	}

	// Set computed fields
	if err := d.Set(e2econstants.AttrStatus, cluster.State); err != nil {
		return nil, fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set(e2econstants.AttrCreatedAt, cluster.CreatedAt); err != nil {
		return nil, fmt.Errorf("error setting created_at: %w", err)
	}

	// Initialize tags as empty map (state-only, not sent to API)
	if err := d.Set("tags", make(map[string]interface{})); err != nil {
		return nil, fmt.Errorf("error setting tags: %w", err)
	}

	// Use flattenNodePools() with V3 field names
	if err := d.Set(e2econstants.AttrNodePools, flattenNodePools(nodePools)); err != nil {
		return nil, fmt.Errorf("error setting node_pools: %w", err)
	}

	// Get security groups if attached
	securityGroups, _, err := goe2eClient.Kubernetes.ListAttachedSecurityGroups(ctx, clusterID)
	if err == nil && len(securityGroups) > 0 {
		sgIDs := make([]int, len(securityGroups))
		for i, sg := range securityGroups {
			sgIDs[i] = sg.ID
		}
		if err := d.Set("security_group_ids", sgIDs); err != nil {
			log.Printf("[WARN] Failed to set security_group_ids during import: %s", err)
		}
	}

	return []*schema.ResourceData{d}, nil
}

// V0 schema and state upgrade functions are now in resource_kubernetes_state_upgrade.go
