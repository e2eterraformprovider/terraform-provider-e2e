package kubernetes

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceKubernetesService() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Kubernetes service",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version of the Kubernetes service",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the project. It should be unique",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location of the block storage",
			},
			"slug_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug name of the Kubernetes service",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC ID of the Kubernetes service",
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "List of security group IDs to attach to the cluster (must contain at least one ID)",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet ID of the custom VPC (applicable only if a custom VPC is used)",
			},
			"sku_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SKU ID of the Kubernetes service",
			},
			"node_pools": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of worker node pools",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the worker node pool",
						},
						"slug_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Slug name of the worker node pool",
						},
						"sku_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SKU ID of the worker node pool",
						},
						"specs_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specs name of the worker node pool",
						},
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Services ID of the worker node pool",
						},
						"node_pool_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Its value can be Autoscale or Static",
							ValidateFunc: validation.StringInSlice([]string{
								"Static",
								"Autoscale",
							}, false),
						},
						"worker_node": {
							Type:         schema.TypeInt,
							Optional:     true, //If the type is autoscale then this field is not needed. Otherwise the default value will be 3
							Description:  "Number of worker nodes in the pool",
							ValidateFunc: validation.IntBetween(2, 25),
						},
						"min_vms": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.All(validation.IntAtLeast(2), validation.IntAtMost(25)),
							Description:  "Minimum number of virtual machines",
						},
						"cardinality": {
							Type:        schema.TypeInt,
							Computed:    true, //NEW CHANGE
							Description: "Cardinality computed from min_vms during creation",
						},
						"node_pool_size": {
							Type:        schema.TypeInt,
							Optional:    true, //NEW CHANGE
							Description: "Cardinality computed from min_vms during creation",
						},
						"max_vms": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntAtMost(25),
							Description:  "Maximum number of virtual machines",
						},
						"elasticity_dict": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Elasticity dictionary for the worker node pool",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"worker": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Worker settings in the elasticity dictionary",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"period_number": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Period number",
												},
												"policy_paramter_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Its value can be Default or Custom. If it is custom then you must provide the parameter field.",
													ValidateFunc: validation.StringInSlice([]string{
														"Default",
														"Custom",
													}, false),
												},
												"parameter": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "CPU",
													Description: "Parameter (e.g., CPU, Memory)",
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
													Description: "List of elasticity policies",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "It has a fixed value, i.e, CHANGE",
															},
															"adjust": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Adjust Value. Its value can be 1 or -1",
															},
															"operator": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Operator for adding worker (e.g., >, >=)",
															},
															"value": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Value for adding worker",
															},
															"period": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Period",
															},
															"watch_period": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Period Number",
															},
															"cooldown": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Cooldown",
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
							Description: "Scheduled dictionary for the worker node pool",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"worker": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Worker settings in the scheduled dictionary",
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
																Description: "The cardinality for upscaling",
															},
															"upscale_recurrence": {
																Type:         schema.TypeString,
																Required:     true,
																Description:  "The recurrence timing for upscaling",
																ValidateFunc: validation.StringInSlice([]string{"0 12 * * *", "0 0 1 * *", "0 20 * * *", "0 9 * * 1-5", "0 9-13 * * *"}, false),
															},
															"downscale_cardinality": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The cardinality for downscaling",
															},
															"downscale_recurrence": {
																Type:         schema.TypeString,
																Required:     true,
																Description:  "The recurrence timing for downscaling",
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
							Description: "Policy type for the worker node pool",
						},
						"custom_param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom parameter name for the worker node pool",
						},
						"custom_param_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom parameter value for the worker node pool",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the status of the Kubernetes Service, only to get the status from my account.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the Kubernetes Service",
			},
		},

		CreateContext: resourceCreateKubernetesService,
		ReadContext:   resourceReadKubernetesService,
		UpdateContext: resourceUpdateKubernetesService,
		DeleteContext: resourceDeleteKubernetesService,
		Exists:        resourceExistsKubernetesService,
		Importer: &schema.ResourceImporter{
			State: KubernetesImportStateFunc,
		},
	}
}

func GetSlugName(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
	apiClient := m.(*client.Client)
	log.Printf("[INFO] KUBERNETES PLAN READ STARTS")
	version := d.Get("version").(string)
	log.Printf("--------------MAKING API CALL FOR SLUGNAME-------------")
	kubernetesPlan, err := apiClient.GetKubernetesMasterPlans(d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return "", fmt.Errorf("error getting Kubernetes plans: %s", err.Error())
	}
	// Extract slug_name based on the version
	log.Printf("[DEBUG] kubernetesPlan response: %+v", kubernetesPlan)

	// The API returns nested data: data.data is the array of plans
	dataWrapper, ok := kubernetesPlan["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format: 'data' field is not an object. Actual response: %+v", kubernetesPlan)
	}

	data, ok := dataWrapper["data"].([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format: nested 'data' field is not a list. Actual response: %+v", kubernetesPlan)
	}
	for _, plan := range data {
		planData, ok := plan.(map[string]interface{})
		if !ok {
			continue
		}
		k8sVersion, ok := planData["k8s_version"].(string)
		if !ok {
			continue
		}
		if k8sVersion == version {
			slugName, ok := planData["plan"].(string)
			if ok {
				specs := planData["specs"].(map[string]interface{})
				sku_id := specs["id"].(string)
				d.Set("sku_id", sku_id)
				return slugName, nil
			}
		}
	}
	return "", fmt.Errorf("plan not found for version %s", version)
}

func CreateKubernetesObject(m interface{}, d *schema.ResourceData, slugName string) (*models.KubernetesCreate, diag.Diagnostics) {
	apiClient, ok := m.(*client.Client)
	if !ok {
		return nil, diag.Errorf("Invalid type provided for client")
	}
	log.Printf("[INFO] KUBERNETES OBJECT CREATION STARTS")
	kubernetesObj := models.KubernetesCreate{
		Name:     d.Get("name").(string),
		Version:  d.Get("version").(string),
		VPCID:    d.Get("vpc_id").(string),
		SubnetID: d.Get("subnet_id").(string),
		SKUID:    d.Get("sku_id").(string),
		SlugName: slugName,
	}

	// Handle security_group_ids - at least one is required
	if securityGroupsList, ok := d.GetOk("security_group_ids"); ok {
		sgList := securityGroupsList.([]interface{})
		if len(sgList) == 0 {
			return nil, diag.Errorf("At least one security group ID is required")
		}
		// Set the first security group as the primary one (required by creation API)
		kubernetesObj.SecurityGroupID = sgList[0].(int)
		d.Set("security_group_ids", []int{sgList[0].(int)})
	} else {
		return nil, diag.Errorf("security_group_ids is required and must contain at least one security group ID")
	}

	if nodePools, ok := d.GetOk("node_pools"); ok {
		nodePoolList := nodePools.([]interface{})
		nodePoolsDetail, err := ExpandNodePools(nodePoolList, apiClient, d.Get("project_id").(int), d.Get("location").(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		kubernetesObj.NodePools = nodePoolsDetail
	} else {
		kubernetesObj.NodePools = make([]models.NodePool, 0)
	}
	return &kubernetesObj, nil
}

func resourceCreateKubernetesService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient, ok := m.(*client.Client)
	if !ok {
		return diag.Errorf("Invalid type provided for client")
	}
	var diags diag.Diagnostics
	slugName, err := GetSlugName(ctx, d, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("slug_name", slugName)
	kubernetesObject, diags := CreateKubernetesObject(apiClient, d, slugName)
	if diags != nil {
		return diags
	}
	log.Printf("---------KUBERNETES OBJECT CREATED---------: %+v", kubernetesObject)
	resKubernetes, err := apiClient.NewKubernetesService(kubernetesObject, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if _, codeOK := resKubernetes["code"]; !codeOK {
		return diag.Errorf(resKubernetes["message"].(string))
	}
	data, ok := resKubernetes["data"].(map[string]interface{})
	if !ok {
		return diag.Errorf("Failed to parse 'data' field in the response")
	}
	document, ok := data["DOCUMENT"].(map[string]interface{})
	if !ok {
		return diag.Errorf("Failed to parse 'DOCUMENT' field in the response")
	}
	clusterIDStr, ok := document["ID"].(string)
	if !ok {
		return diag.Errorf("Failed to parse 'ID' field in the 'DOCUMENT'")
	}
	d.SetId(clusterIDStr)
	log.Printf("[INFO] Kubernetes Cluster created successfully with ID: %s", clusterIDStr)

	return diags

}

func resourceReadKubernetesService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("=============INSIDE KUBERNETES READ RESOURCE==========================")
	kubernetesId := d.Id()
	location := d.Get("location").(string)
	kubernetes, err := apiClient.GetKubernetesServiceInfo(kubernetesId, location, d.Get("project_id").(int))
	log.Println("===========GET_KUBERNETES_RESPONSE==========", kubernetes)
	if err != nil {
		return diag.Errorf("error finding Item with ID %s", kubernetesId)
	}

	log.Printf("[INFO] KUBERNETES READ | BEFORE SETTING DATA")
	data := kubernetes["data"].([]interface{})[0].(map[string]interface{})
	log.Printf("[INFO] SETTING--------- (1)")
	serviceIDFloat, ok := data["service_id"].(float64)
	if !ok {
		return diag.Errorf("Failed to convert 'service_id' to float64")
	}
	serviceIDStr := fmt.Sprintf("%.0f", serviceIDFloat)
	d.SetId(serviceIDStr)
	d.Set("name", data["service_name"].(string))
	d.Set("status", data["state"].(string))
	d.Set("version", data["version"].(string))
	d.Set("created_at", data["created_at"].(string))

	// Refresh node_pools from the API so drift is detectable. This must stay above
	// the master VM lookup below, which returns early when no master is found.
	if err := readNodePoolsIntoState(d, apiClient, kubernetesId, location); err != nil {
		return diag.FromErr(err)
	}

	// Fetch and set security_group_ids from master node
	masterVMID := getMasterNodeVMID(data)
	if masterVMID == "" {
		log.Printf("[WARN] Could not find master VM ID in roles")
		return diags
	}

	log.Printf("[INFO] Fetching security groups for VM %s", masterVMID)
	sgResponse, err := apiClient.GetNodeSecurityGroups(masterVMID, d.Get("project_id").(int), location)
	if err != nil {
		log.Printf("[WARN] Failed to fetch security groups: %s", err.Error())
	} else if sgData, ok := sgResponse["data"].([]interface{}); ok {
		var securityGroupIDs []int
		for _, sg := range sgData {
			if sgMap, ok := sg.(map[string]interface{}); ok {
				if sgID, ok := sgMap["id"].(float64); ok {
					securityGroupIDs = append(securityGroupIDs, int(sgID))
				}
			}
		}
		log.Printf("[INFO] Setting security_group_ids: %v", securityGroupIDs)
		d.Set("security_group_ids", securityGroupIDs)
	}

	return diags
}

func resourceDeleteKubernetesService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	kubernetesID := d.Id()
	status := d.Get("status").(string)
	if status != "Running" {
		return diag.Errorf("Kubernetes is in %s state. You can delete it once it comes to the Running state.", status)
	}
	err := apiClient.DeleteKubernetesService(kubernetesID, d.Get("location").(string), d.Get("project_id").(int))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceExistsKubernetesService(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	kubernetesId := d.Id()
	location := d.Get("location").(string)
	_, err := apiClient.GetKubernetesServiceInfo(kubernetesId, location, d.Get("project_id").(int))

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func resourceUpdateKubernetesService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	status := d.Get("status").(string)
	kubernetesId := d.Id()
	if status != "Running" {
		return diag.Errorf("Kubernetes is in %s state. You can update it once it comes to the Running state.", status)
	}
	serviceMapping, err := GetNodePoolServiceMapping(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle security_group_ids changes
	if d.HasChange("security_group_ids") {
		log.Printf("[INFO] Detected change in security_group_ids for Kubernetes Cluster %s", kubernetesId)

		oldRaw, newRaw := d.GetChange("security_group_ids")
		oldSGList := oldRaw.([]interface{})
		newSGList := newRaw.([]interface{})

		// Validate that at least one security group remains
		if len(newSGList) == 0 {
			return diag.Errorf("At least one security group must be attached to the cluster. Cannot remove all security groups.")
		}

		// Convert to int slices
		oldSGs := make(map[int]bool)
		for _, sg := range oldSGList {
			oldSGs[sg.(int)] = true
		}

		newSGs := make(map[int]bool)
		for _, sg := range newSGList {
			newSGs[sg.(int)] = true
		}

		// Find security groups to attach (in new but not in old)
		var toAttach []int
		for sgID := range newSGs {
			if !oldSGs[sgID] {
				toAttach = append(toAttach, sgID)
			}
		}

		// Find security groups to detach (in old but not in new)
		var toDetach []int
		for sgID := range oldSGs {
			if !newSGs[sgID] {
				toDetach = append(toDetach, sgID)
			}
		}

		// Attach new security groups
		if len(toAttach) > 0 {
			log.Printf("[INFO] Attaching security groups %v to Kubernetes Cluster %s", toAttach, kubernetesId)
			_, err := apiClient.AttachSecurityGroupsToKubernetes(kubernetesId, toAttach, d.Get("project_id").(int), d.Get("location").(string))
			if err != nil {
				return diag.Errorf("Failed to attach security groups: %s", err.Error())
			}
		}

		// Detach removed security groups
		if len(toDetach) > 0 {
			log.Printf("[INFO] Detaching security groups %v from Kubernetes Cluster %s", toDetach, kubernetesId)
			_, err := apiClient.DetachSecurityGroupsFromKubernetes(kubernetesId, toDetach, d.Get("project_id").(int), d.Get("location").(string))
			if err != nil {
				return diag.Errorf("Failed to detach security groups: %s", err.Error())
			}
		}
	}

	if d.HasChange("node_pools") {
		oldData, newData := d.GetChange("node_pools")

		oldNodePools := oldData.([]interface{})
		newNodePools := newData.([]interface{})

		for _, oldNodePool := range oldNodePools {
			oldNodePoolMap := oldNodePool.(map[string]interface{})
			oldNPName := oldNodePoolMap["name"].(string)
			oldServiceFind := serviceMapping[oldNPName]
			if oldServiceFind == nil {
				return diag.Errorf("The Node Pool you are trying to delete does not exist!")
			}
			oldServiceID := oldServiceFind.(float64)
			found := false
			if len(newNodePools) <= 0 {
				return diag.Errorf("Atleast one node pool must be present in a kubernetes cluster!")
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
				kubernetes, err := apiClient.CheckNodePoolStatus(kubernetesId, d.Get("project_id").(int), d.Get("location").(string))
				if err != nil {
					return diag.Errorf("error finding Item with ID %s", kubernetesId)
				}
				if !IsNodePoolRunning(oldServiceID, kubernetes["data"].([]interface{})) {
					d.Set("node_pools", oldData)
					return diag.Errorf("You can delete a Node Pool once it comes to the running state")
				}
				response, err := apiClient.DeleteNodePool(oldServiceID, d.Get("project_id").(int), d.Get("location").(string))
				if err != nil {
					if response == nil {
						return nil
					}
					if len(response) > 0 {
						return diag.Errorf("Error: %s", response["Status"])
					}
					return diag.FromErr(err)
				}
				// return nil
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
				oldServiceID := serviceMapping[oldNPName].(float64)
				// If exists then check if there is any change in cardinality
				if newNPName == oldNPName {
					found = true
					oldCardinality := oldNodePoolMap["cardinality"].(int)
					if oldNodePoolMap["node_pool_type"].(string) == "Static" {
						oldCardinality = oldNodePoolMap["worker_node"].(int)
					}
					node_pool_size := oldCardinality
					if newNodePoolMap["node_pool_size"].(int) != 0 {
						node_pool_size = newNodePoolMap["node_pool_size"].(int)
					}
					log.Printf("----------------PREV CARD:%+v     NEW CARD:%+v------------------", oldCardinality, node_pool_size)
					if node_pool_size < 2 {
						return diag.Errorf("Cardinality of worker nodes cannot be less than 2")
					}
					if oldCardinality != node_pool_size {
						nodePoolResize := models.NodePoolResize{
							NodePoolSize: newNodePoolMap["node_pool_size"].(int),
						}
						newNodePoolMap["cardinality"] = newNodePoolMap["node_pool_size"].(int)
						response, err := apiClient.UpdateNodePoolCardinality(&nodePoolResize, oldServiceID, d.Get("project_id").(int), d.Get("location").(string))
						if err != nil {
							if response == nil {
								// return nil
								break
							}
							if len(response) > 0 {
								return diag.Errorf("Error: %s", response["errors"])
							}
							return diag.FromErr(err)
						}
						break
					}
					old_node_pool_type := oldNodePoolMap["node_pool_type"].(string)
					new_node_pool_type := newNodePoolMap["node_pool_type"].(string)
					// You cannot change the node pool type from Static to Autoscale and vice versa
					if old_node_pool_type != new_node_pool_type {
						return diag.Errorf("You cannot change the node pool type")
					}
					if new_node_pool_type == "Static" {
						break
					}
					nodePoolObject, err := ExpandNPUpdate(newNodePoolMap, apiClient, d.Get("project_id").(int), d.Get("location").(string))
					if err != nil {
						return diag.FromErr(err)
					}
					response, err := apiClient.UpdateNodePoolDetails(&nodePoolObject, oldServiceID, d.Get("project_id").(int), d.Get("location").(string))
					if err != nil {
						return diag.FromErr(err)
					}
					if _, codeOK := response["code"]; !codeOK {
						return diag.Errorf(response["message"].(string))
					}
					break
				}
			}
			//If not found meaning this is a new NodePool
			if !found {
				var nodePoolList []interface{}
				nodePoolList = append(nodePoolList, newNodePools[i])
				nodePoolsDetail, err := ExpandNodePools(nodePoolList, apiClient, d.Get("project_id").(int), d.Get("location").(string))
				if err != nil {
					return diag.FromErr(err)
				}
				kubernetesObj := models.NodePoolAdd{}
				kubernetesObj.NodePools = nodePoolsDetail
				log.Printf("----------------------ADDING A NEW NODE POOL-------------------")
				response, err := apiClient.AddNodePool(&kubernetesObj, kubernetesId, d.Get("project_id").(int), d.Get("location").(string))
				if err != nil {
					return diag.FromErr(err)
				}
				if _, codeOK := response["code"]; !codeOK {
					return diag.Errorf(response["message"].(string))
				}
				continue
				// return nil
			}
		}
	}

	return resourceReadKubernetesService(ctx, d, m)
}

func GetNodePoolServiceMapping(ctx context.Context, d *schema.ResourceData, m interface{}) (map[string]interface{}, error) {
	apiClient := m.(*client.Client)
	log.Printf("[INFO] KUBERNETES CLUSTER NODE POOLS MAPPING STARTS")
	clusterID := d.Id()
	// Initialize the map to store service_name and service_id mappings
	serviceMapping := make(map[string]interface{})
	nodePoolList, err := apiClient.GetKubernetesNodePools(clusterID, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return serviceMapping, fmt.Errorf("error getting list of kluster's node pools list: %s", err.Error())
	}
	if err != nil {
		return serviceMapping, fmt.Errorf("error getting list of kluster's node pools list: %s", err.Error())
	}
	// Extract service_name and service_id from each item in the data array
	for _, nodePool := range nodePoolList["data"].([]interface{}) {
		nodePoolData := nodePool.(map[string]interface{})
		serviceName := nodePoolData["service_name"].(string)
		serviceID := nodePoolData["service_id"].(float64) // Assuming service_id is a number
		serviceMapping[serviceName] = serviceID
	}

	return serviceMapping, nil
}

func IsNodePoolRunning(oldServiceID float64, nodePools []interface{}) bool {
	var status string
	for _, nodepool := range nodePools {
		npdetail := nodepool.(map[string]interface{})
		serviceID := npdetail["service_id"].(float64)
		if serviceID == oldServiceID {
			status = npdetail["state"].(string)
			if status == "Running" {
				return true
			}
		}
	}
	return false
}

// getMasterNodeVMID extracts the master node VM ID 
func getMasterNodeVMID(data map[string]interface{}) string {
	roles, ok := data["roles"].([]interface{})
	if !ok || len(roles) == 0 {
		return ""
	}

	for _, role := range roles {
		roleMap, ok := role.(map[string]interface{})
		if !ok {
			continue
		}

		if roleMap["role_name"] == "master" {
			nodeList, ok := roleMap["node_id_list"].([]interface{})
			if ok && len(nodeList) > 0 {
				if vmID, ok := nodeList[0].(float64); ok {
					return fmt.Sprintf("%.0f", vmID)
				}
			}
		}
	}
	return ""
}

// KubernetesImportStateFunc handles import for Kubernetes resources
func KubernetesImportStateFunc(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format: expected project_id/location/cluster_id")
	}

	projectIDStr := parts[0]
	location := parts[1]
	clusterID := parts[2]

	projectIDInt, err := strconv.Atoi(projectIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: must be numeric")
	}

	d.Set("project_id", projectIDInt)
	d.Set("location", location)
	d.SetId(clusterID)

	return []*schema.ResourceData{d}, nil
}

