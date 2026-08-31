package autoscaling

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceScalerGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The project ID associated with the scaler group.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location where the scaler group will be created.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The internal ID of the plan derived from plan_name.",
			},
			"sku_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SKU ID (same as plan_id) used for the scaler group.",
			},
			"slug_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The slug representation of the plan.",
			},

			"vm_image_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					trim := func(name string) string {
						if idx := strings.Index(name, "_"); idx != -1 {
							return name[:idx]
						}
						return name
					}
					trimmedOld := trim(old)
					trimmedNew := trim(new)

					log.Printf("[DEBUG] DiffSuppressFunc: old=%s → %s, new=%s → %s", old, trimmedOld, new, trimmedNew)

					return trimmedOld == trimmedNew
				},
			},
			"vm_image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"my_account_sg_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The Security Group ID to attach to the scaler group. If not provided, a default will be fetched from the API.",
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "The list of Security Group IDs currently attached to the scaler group. Used for updates.",
			},

			"is_encryption_enabled": {
				Type:     schema.TypeBool,
				Required: true,

				ForceNew:    true,
				Description: "Enable encryption for the scaler group. Defaults to false.",
			},
			"encryption_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "Passphrase for encryption (if enabled). Defaults to empty string.",
			},
			"is_public_ip_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to assign a public IP to nodes. Can only be updated when the scaler group is stopped and a VPC is attached.",
			},

			"provision_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped"}, false),
				Description:  "Set to 'Stopped' to stop the Scaler Group, or 'Running' to start it.",
			},

			"min_nodes": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 {
						errs = append(errs, fmt.Errorf("%q must be at least 1, got: %d", key, v))
					}
					return
				},
			},

			"max_nodes": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"desired": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ipv4_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"subnet_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"used_ips": {
										Type:     schema.TypeInt,
										Computed: true,
									},

									"total_ips": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"policy": {
				Type:     schema.TypeList,
				Optional: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":           {Type: schema.TypeString, Required: true},
						"adjust":         {Type: schema.TypeInt, Required: true},
						"parameter":      {Type: schema.TypeString, Required: true},
						"operator":       {Type: schema.TypeString, Required: true},
						"value":          {Type: schema.TypeString, Required: true},
						"period_number":  {Type: schema.TypeString, Required: true},
						"period_seconds": {Type: schema.TypeString, Required: true},
						"cooldown":       {Type: schema.TypeString, Required: true},
					},
				},
			},
			"scheduled_policy": {
				Type:     schema.TypeList,
				Optional: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":       {Type: schema.TypeString, Required: true},
						"adjust":     {Type: schema.TypeString, Required: true},
						"recurrence": {Type: schema.TypeString, Required: true},
					},
				},
			},
		},
		CreateContext: resourceCreateScalerGroup,
		ReadContext:   resourceReadScalerGroup,
		DeleteContext: resourceDeleteScalerGroup,
		UpdateContext: resourceUpdateScalerGroup,
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			min := diff.Get("min_nodes").(int)
			desired := diff.Get("desired").(int)
			max := diff.Get("max_nodes").(int)

			if min > desired {
				return fmt.Errorf("min_nodes (%d) cannot be greater than desired (%d)", min, desired)
			}

			if desired > max {
				return fmt.Errorf("desired (%d) cannot be greater than max_nodes (%d)", desired, max)
			}

			return nil
		},

		Importer: &schema.ResourceImporter{
			State: node.CustomImportStateFunc,
		},
	}
}

func resourceCreateScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[INFO] Starting CreateScalerGroup operation")

	apiClient := m.(*client.Client)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	imageName := d.Get("vm_image_name").(string)

	savedImage, err := apiClient.GetSavedImageByName(imageName, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch saved image details for '%s': %v", imageName, err))
	}

	if err := d.Set("vm_image_id", savedImage.ImageID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set vm_image_id: %v", err))
	}
	if err := d.Set("vm_template_id", savedImage.TemplateID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set vm_template_id: %v", err))
	}

	log.Printf("[DEBUG] Image Details → ID: %s, TemplateID: %d, Distro: %s", savedImage.ImageID, savedImage.TemplateID, savedImage.Distro)

	var sgID int
	if v, ok := d.GetOk("my_account_sg_id"); ok {
		sgID = v.(int)
		log.Printf("[INFO] Using user-provided Security Group ID: %d", sgID)
	} else {
		sgID, err = apiClient.GetDefaultSecurityGroupID(projectID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch default security group ID: %v", err))
		}
		log.Printf("[INFO] Using default Security Group ID from API: %d", sgID)
		if err := d.Set("my_account_sg_id", sgID); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set my_account_sg_id: %v", err))
		}
	}

	if err := d.Set("security_group_ids", []int{sgID}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set security_group_ids: %v", err))
	}

	req, err := expandCreateScalerGroupRequest(d, m.(*client.Client), projectID, location, sgID)
	if err != nil {
		return diag.FromErr(err)
	}

	requestJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("[DEBUG] CreateScalerGroup Request JSON:\n%s", requestJSON)

	resp, err := apiClient.CreateScalerGroup(req, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create scaler group: %v", err))
	}

	log.Printf("[INFO] ScalerGroup created with ID: %s", resp.ID)
	d.SetId(resp.ID)

	return resourceReadScalerGroup(ctx, d, m)
}

func resourceReadScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading ScalerGroup ID: %s", d.Id())

	apiClient := m.(*client.Client)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	id := d.Id()

	group, err := apiClient.GetScalerGroup(id, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read scaler group: %v", err))
	}

	log.Printf("[DEBUG] Retrieved ScalerGroup: %+v", group)

	// Suppress refresh diff on vm_image_name
	stateVMImageName := d.Get("vm_image_name").(string)
	apiVMImageName := group.VMImageName
	if !strings.HasPrefix(stateVMImageName, apiVMImageName) {
		log.Printf("[INFO] Updating vm_image_name to: %s", apiVMImageName)
		if err := d.Set("vm_image_name", apiVMImageName); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set vm_image_name: %v", err))
		}
	} else {
		log.Printf("[INFO] Keeping existing vm_image_name: %s", stateVMImageName)
	}

	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %v", err))
	}
	if err := d.Set("desired", group.Desired); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set desired: %v", err))
	}
	if err := d.Set("min_nodes", group.MinNodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set min_nodes: %v", err))
	}
	if err := d.Set("max_nodes", group.MaxNodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set max_nodes: %v", err))
	}
	if err := d.Set("plan_name", group.PlanName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set plan_name: %v", err))
	}
	if err := d.Set("plan_id", strconv.Itoa(group.PlanID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set plan_id: %v", err))
	}
	if err := d.Set("sku_id", strconv.Itoa(group.PlanID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set sku_id: %v", err))
	}
	if err := d.Set("policy_type", group.PolicyType); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set policy_type: %v", err))
	}
	if err := d.Set("vm_image_id", strconv.Itoa(group.ImageID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set vm_image_id: %v", err))
	}
	// Normalize provision_status to avoid diff when API returns transitional state
	normalizedStatus := group.ProvisionStatus
	if group.ProvisionStatus == "Stopping" {
		log.Printf("[INFO] Normalizing provision_status 'Stopping' → 'Stopped'")
		normalizedStatus = "Stopped"
	} else if group.ProvisionStatus == "Starting" {
		log.Printf("[INFO] Normalizing provision_status 'Starting' → 'Running'")
		normalizedStatus = "Running"
	}

	if err := d.Set("provision_status", normalizedStatus); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set provision_status: %v", err))
	}

	var fetchedVPCs []map[string]interface{}

	attachedVPCs, err := apiClient.GetAttachedVPCsForScalerGroup(id, projectID, location)
	if err != nil {
		log.Printf("[WARN] Failed to fetch attached VPCs from scaler group: %v", err)
	} else {
		for _, vpcPartial := range attachedVPCs {
			vpcName := vpcPartial.Name

			vpcDetail, err := apiClient.GetVpcDetailsByName(projectID, location, vpcName)
			if err != nil {
				log.Printf("[WARN] Failed to fetch VPC details for %s: %v", vpcName, err)
				continue
			}

			vpcEntry := map[string]interface{}{
				"name":       vpcDetail.Name,
				"network_id": vpcDetail.NetworkID,
				"ipv4_cidr":  vpcDetail.IPv4CIDR,
				"state":      vpcDetail.State,
			}

			var subnets []map[string]interface{}
			for _, s := range vpcDetail.Subnets {
				subnets = append(subnets, map[string]interface{}{
					"id":          s.ID,
					"subnet_name": s.SubnetName,
					"cidr":        s.CIDR,
					"used_ips":    s.UsedIPs,
					"total_ips":   s.TotalIPs,
				})
			}

			vpcEntry["subnets"] = subnets
			fetchedVPCs = append(fetchedVPCs, vpcEntry)
		}
	}

	if err := d.Set("vpc", fetchedVPCs); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set vpc details: %v", err))
	}

	if templateID, ok := d.Get("vm_template_id").(int); ok && templateID > 0 {
		_, slugName, err := apiClient.GetPlanDetailsFromPlanName(templateID, group.PlanName, projectID, location)
		if err == nil {
			d.Set("slug_name", slugName)
		} else {
			log.Printf("[WARN] Failed to recompute slug_name: %v", err)
		}
	}

	ipStatus, err := apiClient.GetPublicIPStatus(id, projectID, location)
	if err != nil {
		log.Printf("[WARN] Failed to fetch public IP status: %v", err)
	} else {
		d.Set("is_public_ip_required", ipStatus.IsPublicIPRequired)
	}

	log.Printf("[INFO] ScalerGroup ID %s state synced successfully", id)
	return nil
}

func resourceDeleteScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting ScalerGroup ID: %s", d.Id())

	apiClient := m.(*client.Client)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	id := d.Id()

	if err := apiClient.DeleteScalerGroup(id, projectID, location); err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete scaler group: %v", err))
	}

	d.SetId("")
	log.Println("[INFO] ScalerGroup deleted successfully")
	return nil
}

func expandCreateScalerGroupRequest(d *schema.ResourceData, client *client.Client, projectID, location string, sgID int) (*models.CreateScalerGroupRequest, error) {
	planName := d.Get("plan_name").(string)
	imageName := d.Get("vm_image_name").(string)

	image, err := client.GetSavedImageByName(imageName, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch saved image: %v", err)
	}

	planID, slugName, err := client.GetPlanDetailsFromPlanName(image.TemplateID, planName, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plan details: %v", err)
	}

	var elasticPolicies []models.ElasticPolicy
	if v, ok := d.GetOk("policy"); ok {
		for _, p := range v.([]interface{}) {
			pMap := p.(map[string]interface{})
			elasticPolicies = append(elasticPolicies, models.ElasticPolicy{
				Type:          pMap["type"].(string),
				Adjust:        pMap["adjust"].(int),
				Parameter:     pMap["parameter"].(string),
				Operator:      pMap["operator"].(string),
				Value:         pMap["value"].(string),
				PeriodNumber:  pMap["period_number"].(string),
				PeriodSeconds: pMap["period_seconds"].(string),
				Cooldown:      pMap["cooldown"].(string),
			})
		}
	}

	var schedPolicies []models.ScheduledPolicy
	if v, ok := d.GetOk("scheduled_policy"); ok {
		for _, s := range v.([]interface{}) {
			sMap := s.(map[string]interface{})
			schedPolicies = append(schedPolicies, models.ScheduledPolicy{
				Type:       sMap["type"].(string),
				Adjust:     sMap["adjust"].(string),
				Recurrence: sMap["recurrence"].(string),
			})
		}
	}

	var vpcDetails []models.VPCDetail
	if v, ok := d.GetOk("vpc"); ok {
		for _, vRaw := range v.([]interface{}) {
			vMap := vRaw.(map[string]interface{})
			vpcName := vMap["name"].(string)

			vpcMeta, err := client.GetVpcDetailsByName(projectID, location, vpcName)
			if err != nil {
				return nil, fmt.Errorf("failed to get VPC details for %s: %v", vpcName, err)
			}

			var selectedSubnets []models.SubnetDetail
			for _, subnet := range vpcMeta.Subnets {
				selectedSubnets = append(selectedSubnets, models.SubnetDetail{
					ID:         subnet.ID,
					SubnetName: subnet.SubnetName,
					CIDR:       subnet.CIDR,
					UsedIPs:    subnet.UsedIPs,
					TotalIPs:   subnet.TotalIPs,
				})
			}

			vpcDetails = append(vpcDetails, models.VPCDetail{
				Name:      vpcMeta.Name,
				NetworkID: vpcMeta.NetworkID,
				IPv4CIDR:  vpcMeta.IPv4CIDR,
				State:     vpcMeta.State,
				Subnets:   selectedSubnets,
			})
		}
	}

	return &models.CreateScalerGroupRequest{
		Name:                 d.Get("name").(string),
		PlanName:             planName,
		PlanID:               planID,
		SKUID:                planID,
		SlugName:             slugName,
		VMImageName:          imageName,
		VMImageID:            image.ImageID,
		VMTemplateID:         image.TemplateID,
		MyAccountSGID:        sgID,
		IsEncryptionEnabled:  d.Get("is_encryption_enabled").(bool),
		EncryptionPassphrase: d.Get("encryption_passphrase").(string),
		IsPublicIPRequired:   d.Get("is_public_ip_required").(bool),
		MinNodes:             strconv.Itoa(d.Get("min_nodes").(int)),
		MaxNodes:             strconv.Itoa(d.Get("max_nodes").(int)),
		Desired:              strconv.Itoa(d.Get("desired").(int)),
		PolicyType:           d.Get("policy_type").(string),
		Policy:               elasticPolicies,
		ScheduledPolicy:      schedPolicies,
		VPC:                  vpcDetails,
	}, nil
}

func resourceUpdateScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	id := d.Id()

	if d.HasChange("provision_status") {
		oldStatus, newStatus := d.GetChange("provision_status")
		log.Printf("[INFO] Changing provision_status from %s → %s", oldStatus, newStatus)

		intID, err := strconv.Atoi(id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid scaler group ID: %v", err))
		}

		if err := apiClient.UpdateScalerGroupStatus(intID, newStatus.(string), projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update provision_status to %s: %v", newStatus, err))
		}

		return resourceReadScalerGroup(ctx, d, m)
	}

	minNodes := d.Get("min_nodes").(int)
	maxNodes := d.Get("max_nodes").(int)
	desired := d.Get("desired").(int)

	if desired < minNodes || desired > maxNodes {
		return diag.Errorf("desired node count (%d) must be between min_nodes (%d) and max_nodes (%d)", desired, minNodes, maxNodes)
	}

	// If only desired changed, call separate API
	if d.HasChange("desired") &&
		!(d.HasChange("min_nodes") || d.HasChange("max_nodes") || d.HasChange("policy_type") || d.HasChange("policy") || d.HasChange("scheduled_policy")) {
		log.Printf("[INFO] Only desired node count changed; using separate API.")
		intID, err := strconv.Atoi(id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid scaler group ID: %v", err))
		}
		if err := apiClient.UpdateDesiredNodeCount(intID, desired, projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update desired node count: %v", err))
		}
		return resourceReadScalerGroup(ctx, d, m)
	}

	if d.HasChange("security_group_ids") {
		log.Printf("[INFO] Detected change in security_group_ids for Scaler Group %s", id)

		group, err := apiClient.GetScalerGroup(id, projectID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch scaler group status: %w", err))
		}
		if group.ProvisionStatus != "Running" {
			return diag.Errorf("Scaler group must be in 'Running' state to update security groups. Current: %s", group.ProvisionStatus)
		}

		oldRaw, newRaw := d.GetChange("security_group_ids")
		oldList := expandIntList(oldRaw.([]interface{}))
		newList := expandIntList(newRaw.([]interface{}))

		if len(newList) == 0 {
			return diag.Errorf("At least one security group must be attached to the scaler group")
		}

		oldStr := intSliceToStringSlice(oldList)
		newStr := intSliceToStringSlice(newList)

		toAttach := difference(newStr, oldStr)
		toDetach := difference(oldStr, newStr)

		for _, sgIDStr := range toAttach {
			sgID, _ := strconv.Atoi(sgIDStr)
			log.Printf("[INFO] Attaching Security Group ID %d", sgID)
			if err := apiClient.AddSecurityGroupToScalergroup(id, sgID, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach SG %d: %v", sgID, err))
			}
		}

		for _, sgIDStr := range toDetach {
			sgID, _ := strconv.Atoi(sgIDStr)
			log.Printf("[INFO] Detaching Security Group ID %d", sgID)
			if err := apiClient.DetachSecurityGroupFromScalergroup(id, sgID, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach SG %d: %v", sgID, err))
			}
		}
	}

	if d.HasChange("vpc") {

		group, err := apiClient.GetScalerGroup(d.Id(), projectID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch scaler group status for update: %w", err))
		}
		if group.ProvisionStatus != "Stopped" {
			return diag.Errorf("VPCs can only be attached or detached when the scaler group is in 'stopped' state. Current state: %q", group.ProvisionStatus)
		}

		oldRaw, newRaw := d.GetChange("vpc")
		oldList := extractVpcNames(oldRaw.([]interface{}))
		newList := extractVpcNames(newRaw.([]interface{}))

		toAttach := difference(newList, oldList)
		toDetach := difference(oldList, newList)

		for _, vpcName := range toAttach {
			vpcDetails, err := apiClient.GetVpcDetailsByName(projectID, location, vpcName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to get VPC details for name %q: %w", vpcName, err))
			}
			err = apiClient.AttachVPCToScalerGroup(d.Id(), []models.VPCDetail{*vpcDetails}, projectID, location)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach VPC %q: %w", vpcName, err))
			}
		}

		for _, vpcName := range toDetach {
			vpcDetails, err := apiClient.GetVpcDetailsByName(projectID, location, vpcName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to get VPC ID for name %q: %w", vpcName, err))
			}
			err = apiClient.DetachVPCFromScalerGroup(d.Id(), strconv.Itoa(vpcDetails.NetworkID), projectID, location)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach VPC %q: %w", vpcName, err))
			}
		}

		vpcNames := extractVpcNames(d.Get("vpc").([]interface{}))
		vpcStateList := []map[string]interface{}{}

		for _, vpcName := range vpcNames {
			vpcDetails, err := apiClient.GetVpcDetailsByName(projectID, location, vpcName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to refresh VPC details for %q: %w", vpcName, err))
			}

			subnetList := []map[string]interface{}{}
			for _, sn := range vpcDetails.Subnets {
				subnetList = append(subnetList, map[string]interface{}{
					"id":          sn.ID,
					"subnet_name": sn.SubnetName,
					"cidr":        sn.CIDR,
					"used_ips":    sn.UsedIPs,
					"total_ips":   sn.TotalIPs,
				})
			}

			vpcStateList = append(vpcStateList, map[string]interface{}{
				"name":       vpcDetails.Name,
				"network_id": vpcDetails.NetworkID,
				"ipv4_cidr":  vpcDetails.IPv4CIDR,
				"state":      vpcDetails.State,
				"subnets":    subnetList,
			})
		}

		if err := d.Set("vpc", vpcStateList); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set vpc state: %w", err))
		}
	}

	if d.HasChange("is_public_ip_required") {
		oldVal, newVal := d.GetChange("is_public_ip_required")
		log.Printf("[INFO] is_public_ip_required changed from %v to %v", oldVal, newVal)

		currentStatus := d.Get("provision_status").(string)
		if currentStatus != "Stopped" {
			return diag.Errorf("ScalerGroup must be in 'Stopped' state to attach/detach public IP")
		}

		vpcsRaw, ok := d.GetOk("vpc")
		if !ok || len(vpcsRaw.([]interface{})) == 0 {
			return diag.Errorf("At least one VPC must be attached to attach/detach public IP")
		}

		if newVal.(bool) {
			log.Printf("[INFO] Triggering Public IP ATTACH")
			_, err := apiClient.AttachPublicIP(d.Id(), projectID, location)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach public IP: %v", err))
			}
		} else {
			log.Printf("[INFO] Triggering Public IP DETACH")
			_, err := apiClient.DetachPublicIP(d.Id(), projectID, location)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach public IP: %v", err))
			}
		}
	}

	if !(d.HasChange("min_nodes") || d.HasChange("max_nodes") || d.HasChange("policy_type") || d.HasChange("policy") || d.HasChange("scheduled_policy")) {
		log.Println("[INFO] No relevant changes detected, skipping update.")
		return nil
	}

	policies := []models.ElasticPolicy{}
	for _, p := range d.Get("policy").([]interface{}) {
		pMap := p.(map[string]interface{})
		policies = append(policies, models.ElasticPolicy{
			Type:          pMap["type"].(string),
			Adjust:        pMap["adjust"].(int),
			Parameter:     pMap["parameter"].(string),
			Operator:      pMap["operator"].(string),
			Value:         pMap["value"].(string),
			PeriodNumber:  pMap["period_number"].(string),
			PeriodSeconds: pMap["period_seconds"].(string),
			Cooldown:      pMap["cooldown"].(string),
		})
	}

	schedPolicies := []models.ScheduledPolicy{}
	for _, s := range d.Get("scheduled_policy").([]interface{}) {
		sMap := s.(map[string]interface{})
		schedPolicies = append(schedPolicies, models.ScheduledPolicy{
			Type:       sMap["type"].(string),
			Adjust:     sMap["adjust"].(string),
			Recurrence: sMap["recurrence"].(string),
		})
	}

	var policyType string
	if len(policies) > 0 {
		policyType = d.Get("policy_type").(string)
	}
	req := &models.UpdateScalerGroupRequest{
		Name:            d.Get("name").(string),
		PlanID:          d.Get("plan_id").(string),
		MinNodes:        minNodes,
		MaxNodes:        maxNodes,
		PolicyType:      policyType, // empty string if no elastic policies
		Policy:          policies,
		ScheduledPolicy: schedPolicies,
	}

	log.Printf("[INFO] Updating ScalerGroup %s with new configuration...", id)
	if err := apiClient.UpdateScalerGroup(id, req, projectID, location); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update scaler group: %v", err))
	}

	return resourceReadScalerGroup(ctx, d, m)
}

func extractVpcNames(vpcs []interface{}) []string {
	var names []string
	for _, raw := range vpcs {
		if m, ok := raw.(map[string]interface{}); ok {
			if name, ok := m["name"].(string); ok {
				names = append(names, name)
			}
		}
	}
	return names
}

func difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	var diff []string
	for _, x := range a {
		if !mb[x] {
			diff = append(diff, x)
		}
	}
	return diff
}

func expandIntList(raw []interface{}) []int {
	result := make([]int, len(raw))
	for i, v := range raw {
		result[i] = v.(int)
	}
	return result
}

func intSliceToStringSlice(in []int) []string {
	result := make([]string, len(in))
	for i, v := range in {
		result[i] = strconv.Itoa(v)
	}
	return result
}
