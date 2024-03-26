package kubernetes

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	// "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
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
				ForceNew:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version of the Kubernetes service",
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the project. It should be unique",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Location of the block storage",
				ValidateFunc: validation.StringInSlice([]string{
					"Delhi",
					"Mumbai",
				}, false),
				Default: "Delhi",
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
						// "parameter": {
						// 	Type:        schema.TypeString,
						// 	Optional:    true,
						// 	Default:     "CPU",
						// 	Description: "Parameter (e.g., CPU, Memory)",
						// 	ValidateFunc: validation.Any(
						// 		validation.StringInSlice([]string{"Memory", "CPU"}, false),
						// 		validation.StringMatch(
						// 			regexp.MustCompile(`^[A-Z0-9]([_]?[A-Z0-9])+$`),
						// 			"Parameter Name should be at least 2 characters long with upper case characters, numbers and underscore and must be start and end with characters or numbers.",
						// 		),
						// 	),
						// },
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
																// ValidateFunc: validation.StringInSlice([]string{
																// 	"1",
																// 	"-1",
																// }, false),
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
			State: schema.ImportStatePassthrough,
		},
	}
}

func GetSlugName(ctx context.Context, d *schema.ResourceData, m interface{}) (string, error) {
	apiClient := m.(*client.Client)
	log.Printf("[INFO] KUBERNETES PLAN READ STARTS")
	version := d.Get("version").(string)
	log.Printf("--------------MAKING API CALL FOR SLUGNAME-------------")
	kubernetesPlan, err := apiClient.GetKubernetesMasterPlans(d.Get("project_id").(int), d.Get("location").(string))
	// log.Printf("--------------GOT RESPONSE FOR SLUGNAME-------------: %+v", kubernetesPlan)
	if err != nil {
		return "", fmt.Errorf("error getting Kubernetes plans: %s", err.Error())
	}
	// Extract slug_name based on the version
	data, ok := kubernetesPlan["data"].([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format: missing 'data' field or not a list")
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
				// log.Printf("--------------GOT THE SLUGNAME SUCCESSFULLY------------- : %+v", slugName)
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
	d.Set("sku_id", "1178")
	log.Printf("---------------IDHAR TOH PAHUCHA(2)-----------------")
	kubernetesObj := models.KubernetesCreate{
		Name:     d.Get("name").(string),
		Version:  d.Get("version").(string),
		VPCID:    d.Get("vpc_id").(string),
		SKUID:    d.Get("sku_id").(string),
		SlugName: slugName,
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
	// log.Printf("--------------EXPANDED THE KUBERNETES PAYLOAD SUCCESSFULLY------------ : %+v", kubernetesObj)
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
	d.Set("sku_id", "1178")
	resKubernetes, err := apiClient.NewKubernetesService(kubernetesObject, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	// log.Printf("[INFO] Kubernetes Cluster CREATE | RESPONSE BODY | %+v", resKubernetes)
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
	// clusterID, err := strconv.Atoi(clusterIDStr)
	// if err != nil {
	// 	return diag.Errorf("Failed to convert cluster ID to integer: %v", err)
	// }
	d.SetId(clusterIDStr)
	log.Printf("[INFO] Kubernetes Cluster creation | before setting fields")
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
	error := GetNodePoolServiceMapping(ctx, d, m)
	if error != nil {
		return diag.FromErr(error)
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
	if status != "Running" {
		return diag.Errorf("Kubernetes is in %s state. You can update it once it comes to the Running state.", status)
	}
	//Setting the service_id field in the node_pools list
	if d.HasChange("node_pools") {
		log.Printf("----------------------CAUGHT A CHANGE IN NODE POOLS ATLEAST-------------------")
		oldData, newData := d.GetChange("node_pools")

		oldNodePools := oldData.([]interface{})
		newNodePools := newData.([]interface{})

		for _, oldNodePool := range oldNodePools {
			oldNodePoolMap := oldNodePool.(map[string]interface{})
			oldServiceID := oldNodePoolMap["service_id"].(string)
			found := false
			if len(newNodePools) <= 0 {
				return diag.Errorf("Atleast one node pool must be present in a kubernetes cluster!")
			}
			log.Printf("----------------------SEARCHING IF ANY NODE POOL IS MISSING OR NOT-------------------")
			// Check if the old service_id exists in the new node pools
			for _, newNodePool := range newNodePools {
				newNodePoolMap := newNodePool.(map[string]interface{})
				newServiceID := newNodePoolMap["service_id"].(string)
				if oldServiceID == newServiceID {
					found = true
					break
				}
			}
			if !found {
				log.Printf("----------------------THIS NODE POOL IS MISSING: %+v-------------------", oldServiceID)
				response, err := apiClient.DeleteNodePool(oldServiceID, d.Get("project_id").(int), d.Get("location").(string))
				log.Printf("----------------RESPONSE FOR DELETE 204 NO CONTENT(Resource.go)----------------: %+v", response)
				if err != nil {
					if response == nil {
						return nil
					}
					if len(response) > 0 {
						return diag.Errorf("Error: %s", response["errors"])
					}
					return diag.FromErr(err)
				}
				return nil
			}
		}

		var nodePoolList []interface{}
		for i := range newNodePools {
			newNodePoolMap := newNodePools[i].(map[string]interface{})
			newServiceID := newNodePoolMap["service_id"].(string)
			found := false
			// Checking if the old service_id exists in the new node pools
			log.Printf("----------------------CHECKING IF THERE IS ANY ADDITION OF NODE POOLS-------------------")
			for _, oldNodePool := range oldNodePools {
				oldNodePoolMap := oldNodePool.(map[string]interface{})
				oldServiceID := oldNodePoolMap["service_id"].(string)
				// If exists then check if there is any change in cardinality
				if newServiceID == oldServiceID {
					found = true
					log.Printf("----------------------IT CAME HERE MEANING ATLEAST THIS NODE POOL IS NOT NEWLY ADDED. NOW CHECKING IF THERE IS ANY CHANGE IN ITS FIELDS -------------------")
					oldCardinality := oldNodePoolMap["cardinality"].(int)
					newCardinality := newNodePoolMap["cardinality"].(int)
					log.Printf("----------------PREV CARD:%+v     NEW CARD:%+v------------------", oldCardinality, newCardinality)
					if newCardinality < 2 {
						return diag.Errorf("Cardinality of worker nodes cannot be less than 2")
					}
					// If the cardinality has changed, call the API passing the corresponding service_id
					if oldCardinality != newCardinality {
						// nodePoolServiceID := newNodePools[i].(map[string]interface{})["service_id"].(string)
						log.Printf("----------------------THERE IS A CHANGE IN CARDINALITY-------------------")
						response, err := apiClient.UpdateNodePoolCardinality(newServiceID, d.Get("project_id").(int), d.Get("location").(string))
						if err != nil {
							if response == nil {
								return nil
							}
							if len(response) > 0 {
								return diag.Errorf("Error: %s", response["errors"])
							}
							return diag.FromErr(err)
						}
						return nil
					}
					old_node_pool_type := oldNodePoolMap["node_pool_type"].(string)
					new_node_pool_type := newNodePoolMap["node_pool_type"].(string)
					// You cannot change the node pool type from Static to Autoscale and vice versa
					if old_node_pool_type != new_node_pool_type {
						return diag.Errorf("You cannot change the node pool type")
					}
					log.Printf("----------------------GOING INTO Helper's ExpandNPUpdate Function-------------------")
					nodePoolObject, err := ExpandNPUpdate(newNodePoolMap, apiClient, d.Get("project_id").(int), d.Get("location").(string))
					if err != nil {
						return diag.FromErr(err)
					}
					log.Printf("----------------------SUCCESSFUL UPDATE OBJECT CREATION, MAKING A REQUEST NOW FOR UPDATE-------------------")
					response, err := apiClient.UpdateNodePoolDetails(&nodePoolObject, newServiceID, d.Get("project_id").(int), d.Get("location").(string))
					if err != nil {
						return diag.FromErr(err)
					}
					if _, codeOK := response["code"]; !codeOK {
						return diag.Errorf(response["message"].(string))
					}
					return nil

				}
			}
			//If not found meaning this is a new NodePool
			if !found {
				nodePoolList = append(nodePoolList, newNodePools[i])
				nodePoolsDetail, err := ExpandNodePools(nodePoolList, apiClient, d.Get("project_id").(int), d.Get("location").(string))
				if err != nil {
					return diag.FromErr(err)
				}
				kubernetesObj := models.NodePoolAdd{}
				kubernetesObj.NodePools = nodePoolsDetail
				log.Printf("----------------------ADDING A NEW NODE POOL-------------------")
				response, err := apiClient.AddNodePool(&kubernetesObj, newServiceID, d.Get("project_id").(int), d.Get("location").(string))
				if err != nil {
					return diag.FromErr(err)
				}
				if _, codeOK := response["code"]; !codeOK {
					return diag.Errorf(response["message"].(string))
				}
				return nil
			}
		}
	}

	return resourceReadKubernetesService(ctx, d, m)
}

func GetNodePoolServiceMapping(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	log.Printf("[INFO] KUBERNETES CLUSTER NODE POOLS MAPPING STARTS")
	clusterID := d.Id()
	log.Printf("--------------MAKING API CALL FOR SLUGNAME-------------")
	nodePoolList, err := apiClient.GetKubernetesNodePools(clusterID, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return fmt.Errorf("error getting list of kluster's node pools list: %s", err.Error())
	}
	if err != nil {
		return fmt.Errorf("error getting list of kluster's node pools list: %s", err.Error())
	}

	// Initialize the map to store service_name and service_id mappings
	serviceMapping := make(map[string]interface{})

	// Extract service_name and service_id from each item in the data array
	for _, nodePool := range nodePoolList["data"].([]interface{}) {
		nodePoolData := nodePool.(map[string]interface{})
		serviceName := nodePoolData["service_name"].(string)
		serviceID := nodePoolData["service_id"].(float64) // Assuming service_id is a number
		serviceMapping[serviceName] = serviceID
	}

	prevNodepool, currNodePool := d.GetChange("node_pools")
	currNodePool = currNodePool.([]interface{})
	for i, np := range prevNodepool.([]interface{}) {
		nodePool := np.(map[string]interface{})
		serviceName := nodePool["name"].(string)
		if serviceID, ok := serviceMapping[serviceName]; ok {
			d.Get("node_pools").([]interface{})[i].(map[string]interface{})["service_id"] = serviceID
		} else {
			return fmt.Errorf("service_name '%s' not found in the mapping", serviceName)
		}
	}
	return nil

}
