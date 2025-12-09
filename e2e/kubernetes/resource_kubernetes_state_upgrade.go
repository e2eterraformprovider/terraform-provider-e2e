package kubernetes

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceKubernetesResourceV0 returns the schema for schema version 0
// This represents the schema before V3 changes (V2 schema)
func resourceKubernetesResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Common fields
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// Required input fields (V2 - before V3 aliases)
			e2econstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			e2econstants.AttrVersion: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			e2econstants.AttrVPCID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Node pools (V2 schema - before V3 aliases)
			e2econstants.AttrNodePools: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"specs_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"node_pool_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Static",
								"Autoscale",
							}, false),
						},
						"worker_node": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						e2econstants.AttrMinVMs: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						e2econstants.AttrMaxVMs: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"node_pool_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"cardinality": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"slug_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sku_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"elasticity_dict": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"worker": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"period_number": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"policy_paramter_type": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														"Default",
														"Custom",
													}, false),
												},
												"parameter": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "CPU",
												},
												"elasticity_policies": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"adjust": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"operator": {
																Type:     schema.TypeString,
																Required: true,
															},
															"value": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"period": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"watch_period": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"cooldown": {
																Type:     schema.TypeInt,
																Required: true,
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
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"worker": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scheduled_policies": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"upscale_cardinality": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"upscale_recurrence": {
																Type:     schema.TypeString,
																Required: true,
															},
															"downscale_cardinality": {
																Type:     schema.TypeInt,
																Required: true,
															},
															"downscale_recurrence": {
																Type:     schema.TypeString,
																Required: true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_param_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_param_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// Computed fields
			"slug_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sku_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			e2econstants.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			e2econstants.AttrCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// ResourceKubernetesStateUpgradeV0toV1 handles state migration from schema version 0 to 1
// This preserves all V2 fields and initializes new V3 fields with safe defaults
// Exported for testing purposes
func ResourceKubernetesStateUpgradeV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[INFO] Upgrading Kubernetes cluster state from V0 to V1")

	if rawState == nil {
		rawState = make(map[string]interface{})
	}

	// Initialize new V3 cluster fields with safe defaults
	if _, ok := rawState["tags"]; !ok {
		rawState["tags"] = map[string]interface{}{}
	}

	if _, ok := rawState["encryption_enabled"]; !ok {
		rawState["encryption_enabled"] = false
	}

	if _, ok := rawState["encryption_passphrase"]; !ok {
		rawState["encryption_passphrase"] = ""
	}

	if _, ok := rawState["security_group_ids"]; !ok {
		rawState["security_group_ids"] = []interface{}{}
	}

	// Initialize V3 preferred field aliases if not present
	// These will be populated from deprecated fields during read, but we initialize them here
	// to ensure the schema is consistent
	// (cluster_name and kubernetes_version will be set during read from "name" and "version"
	// to preserve backward compatibility)

	// Initialize V3 node pool aliases if not present in node pools
	// Node pool validation and field initialization is delegated to the read operation
	// to preserve backward compatibility with V2 schema

	// Preserve all V2 fields (no automatic renames, backward compatible)
	// All existing fields remain unchanged

	log.Printf("[INFO] Successfully upgraded Kubernetes cluster state to V1")
	return rawState, nil
}
