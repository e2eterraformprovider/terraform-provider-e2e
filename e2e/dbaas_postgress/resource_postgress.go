package dbaas_postgress

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
)

func ResourcePostgresDBaaS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "location where DBAAS is deployed",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID to which the DBaaS instance belongs",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version of database",
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of plan",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of DBaaS instance to be deployed",
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				Default:     "Default",
			},
			"public_ip_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
				Default:     true,
			},
			"power_status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start", "stop", "restart",
				}, false),
			},

			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Disk Storage for upgrade",
			},

			"database": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The database username.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The database password.",
						},
						"dbaas_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The database name.",
						},
					},
				},
			},

			"parameter_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of parameter group that need to be attached",
			},

			"vpc_list": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "List of vpc Id which you want to attach",
			},

			"detach_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			//id of dbaas
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status_title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_instances": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_exist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"connectivity_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vector_database_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_encryption_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},

		CreateContext: resourceCreatePostgress,
		ReadContext:   resourceReadPostgress,
		DeleteContext: resourceDeletePostgress,
		UpdateContext: resourceUpdatePostgress,
		Importer: &schema.ResourceImporter{
			State: node.CustomImportStateFunc,
		},
	}
}

func resourceCreatePostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] DBAAS_POSTGRESS CREATE STARTS")

	// Extract nested database config
	dbConfigList := d.Get("database").([]interface{})
	dbConfigMap := dbConfigList[0].(map[string]interface{})

	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)
	plan := d.Get("plan").(string)
	version := d.Get("version").(string)

	// Get software ID
	software_id, err := apiClient.GetSoftwareId(project_id, location, "PostgreSQL", version)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get template ID
	template_id, err := apiClient.GetTemplateId(project_id, location, plan, strconv.Itoa(software_id))
	if err != nil {
		return diag.FromErr(err)
	}

	var pgID *int
	if v, ok := d.GetOk("parameter_group_id"); ok {
		id := v.(int)
		pgID = &id
	}
	payload := models.DBCreateRequest{
		SoftwareID:       software_id,
		TemplateID:       template_id,
		Name:             d.Get("name").(string),
		Group:            d.Get("group").(string),
		PublicIPRequired: d.Get("public_ip_required").(bool),
		VPCs:             []models.VPC{},
		Database: models.DBConfig{
			User:        dbConfigMap["user"].(string),
			Password:    dbConfigMap["password"].(string),
			DBaaSNumber: dbConfigMap["dbaas_number"].(int),
			Name:        dbConfigMap["name"].(string),
		},
		PGID:                pgID, // Only set if parameter_group_id was provided
		IsEncryptionEnabled: d.Get("is_encryption_enabled").(bool),
	}
	//vpc list config
	vpcList, ok := d.GetOk("vpc_list")
	if ok {
		vpcListDetail, err := apiClient.ExpandVpcList(d, vpcList.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		payload.VPCs = vpcListDetail
	} else {
		payload.VPCs = make([]models.VPC, 0)
	}

	res, err := apiClient.CreatePostgressDB(payload, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] DBAAS POSTGRESS CREATE | RESPONSE BODY | ------ > %+v", res)

	if _, ok := res["code"]; !ok || res["is_limit_available"] == false {
		msg, _ := res["message"].(string)
		return diag.Errorf("%s", msg)
	}

	data, ok := res["data"].(map[string]interface{})
	if !ok {
		return diag.Errorf("Invalid 'data' field in response, expected object but got something else")
	}

	id := int(data["id"].(float64))
	d.SetId(strconv.Itoa(id))
	d.Set("id", id)
	d.Set("name", data["name"])
	d.Set("status", data["status"])
	d.Set("status_title", data["status_title"])
	d.Set("status_actions", data["status_actions"])
	d.Set("num_instances", int(data["num_instances"].(float64)))
	d.Set("project_name", data["project_name"])
	d.Set("snapshot_exist", data["snapshot_exist"])
	d.Set("connectivity_detail", data["connectivity_detail"])
	d.Set("vector_database_status", data["vector_database_status"])
	d.Set("is_encryption_enabled", data["is_encryption_enabled"])

	return diags
}

func resourceReadPostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] DBAAS_POSTGRESS READ STARTS")

	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)
	dbaas_id := d.Get("id").(string)

	res, err := apiClient.GetPostgressDB(dbaas_id, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] DBAAS POSTGRESS READ | RESPONSE BODY | ------ > %+v", res)

	if _, ok := res["code"]; !ok || res["is_limit_available"] == false {
		msg, _ := res["message"].(string)
		return diag.Errorf("%s", msg)
	}

	data, ok := res["data"].(map[string]interface{})
	if !ok {
		return diag.Errorf("Invalid 'data' field in response, expected object but got something else")
	}

	// Set ID and basic fields
	id := int(data["id"].(float64))
	d.SetId(strconv.Itoa(id))
	d.Set("id", id)
	d.Set("name", data["name"])
	d.Set("status", data["status"])
	d.Set("status_title", data["status_title"])
	d.Set("status_actions", data["status_actions"])
	d.Set("num_instances", int(data["num_instances"].(float64)))
	d.Set("project_name", data["project_name"])
	d.Set("snapshot_exist", data["snapshot_exist"])
	d.Set("connectivity_detail", data["connectivity_detail"])
	d.Set("vector_database_status", data["vector_database_status"])
	d.Set("is_encryption_enabled", data["is_encryption_enabled"])

	return diags

}

func resourceUpdatePostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)

	if d.HasChange("power_status") {
		status := d.Get("status").(string)
		dbaas_id := d.Get("id").(string)
		// Block operation if DBaaS is still in "Creating" state
		if status == "CREATING" {
			return diag.Errorf("Cannot perform power operations while DBaaS is in CREATING state")
		}

		powerAction := d.Get("power_status").(string)
		switch powerAction {
		case "start":
			err := apiClient.StartPostgressDB(dbaas_id, d.Get("project_id").(string), d.Get("location").(string))
			if err != nil {
				return diag.FromErr(err)
			}
		case "stop":
			err := apiClient.StopPostgressDB(dbaas_id, d.Get("project_id").(string), d.Get("location").(string))
			if err != nil {
				return diag.FromErr(err)
			}
		case "restart":
			err := apiClient.RestartPostgressDB(dbaas_id, d.Get("project_id").(string), d.Get("location").(string))
			if err != nil {
				return diag.FromErr(err)
			}
		default:
			return diag.Errorf("invalid power_status value: %s", powerAction)
		}
	}

	if d.HasChange("detach_public_ip") {
		status := d.Get("status").(string)
		dbaas_id := d.Get("id").(string)
		prev, _ := d.GetChange("detach_public_ip")

		// Block operation if DBaaS is still in "Creating" state
		if status == "CREATING" {
			d.Set("detach_public_ip", prev)
			return diag.Errorf("Cannot perform attach public ip while DBaaS is in CREATING state")

		}
		if !d.Get("detach_public_ip").(bool) {
			err := apiClient.AttachPublicIpPostgressDB(dbaas_id, d.Get("project_id").(string), d.Get("location").(string))
			if err != nil {
				d.Set("detach_public_ip", prev)
				return diag.FromErr(err)
			}
		} else {
			err := apiClient.DetachPublicIpPostgressDB(dbaas_id, d.Get("project_id").(string), d.Get("location").(string))
			if err != nil {
				d.Set("detach_public_ip", prev)
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("project_id") {
		prev, _ := d.GetChange("project_id")
		d.Set("project_id", prev)
		return diag.Errorf("Cannot change the project id once the instance is created")
	}

	if d.HasChange("vpc_list") {

		status := d.Get("status").(string)

		// Block operation if DBaaS is still in "Creating" state
		if status == "CREATING" {
			return diag.Errorf("Cannot perform vpc changes while DBaaS is in CREATING state")
		}

		prev, curr := d.GetChange("vpc_list")

		prevIDs := prev.(*schema.Set).List()
		currIDs := curr.(*schema.Set).List()

		// convert slice of interface{} to map[int]bool for fast lookup
		toSet := func(list []interface{}) map[int]bool {
			m := make(map[int]bool)
			for _, v := range list {
				m[v.(int)] = true
			}
			return m
		}

		prevMap := toSet(prevIDs)
		currMap := toSet(currIDs)

		var attachIDs, detachIDs []int

		for id := range currMap {
			if !prevMap[id] {
				attachIDs = append(attachIDs, id)
			}
		}
		for id := range prevMap {
			if !currMap[id] {
				detachIDs = append(detachIDs, id)
			}
		}

		projectID := d.Get("project_id").(string)
		location := d.Get("location").(string)
		id := d.Get("id").(string)

		if len(detachIDs) > 0 {
			vpcSet := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeInt}), convertToInterfaces(detachIDs))
			vpcListDetail, err := apiClient.ExpandVpcList(d, vpcSet.List())
			if err != nil {
				d.Set("vpc_list", prev) // rollback
				return diag.FromErr(err)
			}

			payload := models.AttachVPCPayloadRequest{
				Action: "detach",
				VPCs:   vpcListDetail,
			}

			if err := apiClient.DetachVPCPostgressDB(payload, id, projectID, location); err != nil {
				d.Set("vpc_list", prev) // rollback
				return diag.FromErr(err)
			}
		}

		if len(attachIDs) > 0 {
			vpcSet := schema.NewSet(schema.HashSchema(&schema.Schema{Type: schema.TypeInt}), convertToInterfaces(attachIDs))
			vpcListDetail, err := apiClient.ExpandVpcList(d, vpcSet.List())
			if err != nil {
				d.Set("vpc_list", prev) // rollback
				return diag.FromErr(err)
			}

			payload := models.AttachVPCPayloadRequest{
				Action: "attach",
				VPCs:   vpcListDetail,
			}

			if err := apiClient.AttachVPCPostgressDB(payload, id, projectID, location); err != nil {
				d.Set("vpc_list", prev) // rollback
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("parameter_group_id") {

		status := d.Get("status").(string)
		dbaas_id := d.Get("id")
		prev, _ := d.GetChange("parameter_group_id")

		// Block operation if DBaaS is still in "Creating" state
		if status == "CREATING" {
			d.Set("parameter_group_id", prev)
			return diag.Errorf("Cannot perform parameter group changes while DBaaS is in CREATING state")
		}

		if v, ok := d.GetOk("parameter_group_id"); ok {
			pgID := strconv.Itoa(v.(int))
			err := apiClient.UpdateParameterGroup(dbaas_id.(string), pgID, d.Get("project_id").(string), d.Get("location").(string))
			if err != nil {
				d.Set("parameter_group_id", prev)
				return diag.FromErr(err)
			}
		}

	}

	if d.HasChange("plan") {

		prevPlan, currPlan := d.GetChange("plan")
		project_id := d.Get("project_id").(string)
		location := d.Get("location").(string)
		plan := d.Get("plan").(string)
		version := d.Get("version").(string)
		dbaas_id := d.Get("id")

		// Get software ID
		software_id, err := apiClient.GetSoftwareId(project_id, location, "PostgreSQL", version)
		if err != nil {
			return diag.FromErr(err)
		}

		// Get template ID
		template_id, err := apiClient.GetTemplateId(project_id, location, plan, strconv.Itoa(software_id))
		if err != nil {
			return diag.FromErr(err)
		}

		if d.HasChange("power_status") {
			_ = waitForPoweringOffOnDBaaS(m, dbaas_id.(string), project_id, location)
		}

		log.Printf("[INFO] prevPlan %s, currPlan %s", prevPlan.(string), currPlan.(string))

		if d.Get("status").(string) != "SUSPENDED" {
			d.Set("plan", prevPlan)
			return diag.Errorf("cannot Upgrade as the node is not stopped")
		}
		_, err = apiClient.UpgradePostgressPlan(dbaas_id.(string), template_id, project_id, location)

		if err != nil {
			d.Set("plan", prevPlan)
			return diag.FromErr(err)
		}
	}

	if d.HasChange("size") {

		prevSize, currSize := d.GetChange("size")

		if d.Get("status").(string) != "SUSPENDED" {
			d.Set("size", prevSize)
			return diag.Errorf("Cannot perform power operations while DBaaS is in %s state", d.Get("status").(string))
		}

		project_id := d.Get("project_id").(string)
		location := d.Get("location").(string)
		dbaas_id := d.Get("id")

		sizeInt, ok := currSize.(int)
		if !ok {
			d.Set("size", prevSize)
			return diag.Errorf("Expected size to be int, but got %T", currSize)
		}

		err := apiClient.UpgradeDiskStorage(dbaas_id.(string), sizeInt, project_id, location)
		if err != nil {
			d.Set("size", prevSize)
			return diag.FromErr(err)
		}
	}

	return resourceReadPostgress(ctx, d, m)
}

func resourceDeletePostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)

	var diags diag.Diagnostics

	dbaas_status := d.Get("status_title").(string)
	dbaasId := d.Get("id").(string)

	if dbaas_status == constants.NODE_STATUS["CREATING"] {
		return diag.Errorf("Node in %s state", dbaas_status)
	}
	err := apiClient.DeletePostgressDB(dbaasId, d.Get("project_id").(string), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func convertToInterfaces(ids []int) []interface{} {
	out := make([]interface{}, len(ids))
	for i, id := range ids {
		out[i] = id
	}
	return out
}

func waitForPoweringOffOnDBaaS(m interface{}, dbaasID string, project_id string, location string) error {
	apiClient := m.(*client.Client)

	maxRetries := 20 // e.g., retry up to 20 times
	for i := 0; i < maxRetries; i++ {
		time.Sleep(constants.WAIT_TIMEOUT * time.Second)

		dbaasInfo, err := apiClient.GetPostgressDB(dbaasID, project_id, location)
		if err != nil {
			log.Printf("[ERROR] Error fetching DBaaS Info during wait: %s", err)
			return err
		}

		data, ok := dbaasInfo["data"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected DBaaS response format")
		}

		status, ok := data["status"].(string)
		if !ok {
			return fmt.Errorf("DBaaS status field missing in response")
		}

		log.Printf("[INFO] Waiting for DBaaS instance to power off/on, current status: %s", status)

		if status == "SUSPENDED" {
			return nil
		}
	}
	return fmt.Errorf("timeout: DBaaS did not reach SUSPENDED state in time")
}
