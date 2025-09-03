package dbaas_mysql

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMySql() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "location should specified",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the project, this should be unique",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "this field is to specifiy which version of MYSQL user wants to cerate",
			},
			"dbaas_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "have to figure out what to write",
			},
			"database": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "This is the username of the databse",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "This the passowrd of the database",
						},
						"dbaas_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "This the name of database",
						},
					},
				},
			},
			"vpcs": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "List of vpc Id which you want to attach",
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This the name of plan which user wants to select",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "This is the status of your dbaas, only to get the status from my account.",
			},
			"db_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Delhi",
				Description: "This is the location of your db",
				ForceNew:    true,
			},
			"is_encryption_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "look what to write",
			},
			"parameter_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of parameter group that need to be attached",
			},
			"public_ip_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "want to attach public ip or not",
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "additional size of disk you want to attach",
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk size",
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				Default:     "Default",
			},
		},

		CreateContext: ResourceCreateMySqlDB,
		ReadContext:   ResourceReadMySqlDB,
		UpdateContext: ResourceUpdateMySqlDB,
		DeleteContext: ResourceDeleteMySqlDB,
		Importer: &schema.ResourceImporter{
			State: node.CustomImportStateFunc,
		},
	}
}

func CreateMySqlObject(apiClient *client.Client, d *schema.ResourceData) (*models.MySqlCreate, diag.Diagnostics) {

	mySqlobject := models.MySqlCreate{
		Name:             d.Get("dbaas_name").(string),
		ParameterGroupId: d.Get("parameter_group_id").(int),
		PublicIPRequired: d.Get("public_ip_required").(bool),
		Group:            d.Get("group").(string),
	}

	dbList := d.Get("database").([]interface{})
	if len(dbList) > 0 {
		dbMap := dbList[0].(map[string]interface{})

		mySqlobject.Database = models.DBConfig{
			User:        dbMap["user"].(string),
			Password:    dbMap["password"].(string),
			Name:        dbMap["name"].(string),
			DBaaSNumber: dbMap["dbaas_number"].(int),
		}

	}

	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)
	database_version := d.Get("version").(string)

	softwareId, err := apiClient.GetSoftwareId(project_id, location, "MySQL", database_version)

	if err != nil {
		return nil, diag.FromErr(err)
	}
	mySqlobject.SoftwareID = softwareId

	templateId, err := apiClient.GetTemplateId(project_id, location, d.Get("plan").(string), strconv.Itoa(softwareId))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	mySqlobject.TemplateID = templateId

	vpcList, ok := d.GetOk("vpcs")
	if ok {
		vpcListDetail, err := ExpandVpcList(d, vpcList.(*schema.Set).List(), apiClient)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		mySqlobject.Vpcs = vpcListDetail
	} else {
		mySqlobject.Vpcs = make([]models.VPC, 0)
	}

	return &mySqlobject, nil
}

func ResourceCreateMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	apiClient := m.(*client.Client)

	mySqlObj, diags := CreateMySqlObject(apiClient, d)

	if diags != nil {
		log.Println("[ERROR] CreateMySqlObject returned error or nil object")
		return diags
	}

	response, err := apiClient.NewMySqlDb(mySqlObj, d.Get("project_id").(string), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	dataRaw, ok := response["data"]
	if !ok {
		return diag.Errorf("missing 'data' field in response")
	}

	data, ok := dataRaw.(map[string]interface{})
	if !ok {
		return diag.Errorf("invalid type for 'data' field in response")
	}

	idVal, idOK := data["id"].(float64)
	if !idOK {
		return diag.Errorf("id not found in response data")
	}

	d.SetId(strconv.Itoa(int(idVal)))
	d.Set("status", data["status"].(string))
	return diags
}

func ResourceReadMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)

	res, err := apiClient.GetMySqlDbaas(d.Id(), d.Get("project_id").(string), d.Get("location").(string))
	if err != nil {
		return diag.Errorf("finding dbaas_mysql: %v", err)
	}

	data := res.Data
	if err := d.Set("status", data.Status); err != nil {
		log.Printf("[ERROR] Failed to set status: %v", err)
	}

	d.Set("status", data.Status)
	d.Set("is_encryption_enabled", data.IsEncryptionEnabled)
	d.Set("parameter_group_id", data.MasterNode.Database.PGDetail.ID)
	d.Set("disk", data.MasterNode.Disk)

	return diags
}

func ResourceUpdateMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	mySqlDBaaSId := d.Id()
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	dbaasResp, err := apiClient.GetMySqlDbaas(mySqlDBaaSId, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch current DBaaS status: %v", err))
	}

	if dbaasResp == nil || dbaasResp.Data.Status == "" {
		return diag.Errorf("invalid response from GetMySqlDbaas while checking status")
	}

	if dbaasResp.Data.Status == "CREATING" {
		_ = ResourceReadMySqlDB(ctx, d, m)
		return diag.Errorf("Cannot update MySQL DBaaS while it is in 'CREATING' state")
	}

	if d.HasChange("status") {
		oldStatus, newStatus := d.GetChange("status")

		if oldStatus == "CREATING" {
			_ = ResourceReadMySqlDB(ctx, d, m)
			return diag.Errorf("Cannot change status while DBaaS instance is still in 'CREATING' state")
		}

		switch newStatus {
		case "start":
			_, err := apiClient.ResumeMySqlDBaaS(mySqlDBaaSId, projectID, location)
			if err != nil {
				d.Set("status", oldStatus)
				return diag.FromErr(fmt.Errorf("failed to change status of MySQL DBaaS instance to start"))
			}
		case "stop":
			_, err := apiClient.StopMySqlDBaaS(mySqlDBaaSId, projectID, location)
			if err != nil {
				d.Set("status", oldStatus)
				return diag.FromErr(fmt.Errorf("failed to change status of MySQL DBaaS instance to stop"))
			}
		case "restart":
			_, err := apiClient.RestartMySqlDBaaS(mySqlDBaaSId, projectID, location)
			if err != nil {
				d.Set("status", oldStatus)
				return diag.FromErr(fmt.Errorf("failed to change status of MySQL DBaaS instance to restart"))
			}
		}
	}

	if d.HasChange("vpcs") {
		prevRaw, newRaw := d.GetChange("vpcs")

		prevSet := prevRaw.(*schema.Set)
		newSet := newRaw.(*schema.Set)

		added := schema.NewSet(schema.HashInt, []interface{}{})
		for _, v := range newSet.List() {
			if !prevSet.Contains(v) {
				added.Add(v)
			}
		}

		removed := schema.NewSet(schema.HashInt, []interface{}{})
		for _, v := range prevSet.List() {
			if !newSet.Contains(v) {
				removed.Add(v)
			}
		}

		if added.Len() > 0 {
			vpcIDs := added.List()
			vpcDetails, err := ExpandVpcList(d, vpcIDs, apiClient)
			if err != nil {
				d.Set("vpcs", prevSet)
				return diag.FromErr(err)
			}
			attachObj := models.AttachVPCPayloadRequest{
				Action: "attach",
				VPCs:   vpcDetails,
			}
			_, err = apiClient.AttachVpcToMySql(&attachObj, mySqlDBaaSId, projectID, location)
			if err != nil {
				d.Set("vpcs", prevSet)
				return diag.FromErr(fmt.Errorf(" failed to attach VPC(s) to MySQL DBaaS instance"))
			}
		}

		if removed.Len() > 0 {
			vpcIDs := removed.List()
			vpcDetails, err := ExpandVpcList(d, vpcIDs, apiClient)
			if err != nil {
				d.Set("vpcs", prevSet)
				return diag.FromErr(err)
			}
			detachObj := models.AttachVPCPayloadRequest{
				Action: "detach",
				VPCs:   vpcDetails,
			}
			_, err = apiClient.DetachVpcFromMySql(&detachObj, mySqlDBaaSId, projectID, location)
			if err != nil {
				d.Set("vpcs", prevSet)
				return diag.FromErr(fmt.Errorf(" failed to detach VPC(s) from MySQL DBaaS instance"))
			}
		}
	}

	if d.HasChange("parameter_group_id") {
		oldVal, newVal := d.GetChange("parameter_group_id")

		oldInt, oldOK := oldVal.(int)
		newInt, newOK := newVal.(int)

		// PG added
		if oldOK && newOK && oldInt == 0 && newInt != 0 {
			_, err := apiClient.AttachPGToMySqlDBaaS(mySqlDBaaSId, strconv.Itoa(newInt), projectID, location)
			if err != nil {
				d.Set("parameter_group_id", oldVal)
				return diag.FromErr(fmt.Errorf("failed to attach PG ID %d to DBaaS: %v", newInt, err))
			}
		}

		// PG removed
		if oldOK && newOK && oldInt != 0 && newInt == 0 {
			_, err := apiClient.DetachPGFromMySqlDBaaS(mySqlDBaaSId, strconv.Itoa(oldInt), projectID, location)
			if err != nil {
				d.Set("parameter_group_id", oldVal)
				return diag.FromErr(fmt.Errorf("failed to detach PG ID %d from DBaaS: %v", oldInt, err))
			}
		}
	}

	if d.HasChange("public_ip_required") {
		oldVal, newVal := d.GetChange("public_ip_required")

		oldBool := oldVal.(bool)
		newBool := newVal.(bool)

		if !oldBool && newBool {
			_, err := apiClient.AttachPublicIPToMySql(mySqlDBaaSId, projectID, location)
			if err != nil {
				d.Set("public_ip_required", oldVal)
				return diag.FromErr(fmt.Errorf("failed to attach public IP: %v", err))
			}
		} else if oldBool && !newBool {
			_, err := apiClient.DetachPublicIPFromMySql(mySqlDBaaSId, projectID, location)
			if err != nil {
				d.Set("public_ip_required", oldVal)
				return diag.FromErr(fmt.Errorf("failed to detach public IP: %v", err))
			}
		}
	}

	if d.HasChange("plan") {
		prevPlan, currPlan := d.GetChange("plan")

		projectIDRaw, ok := d.GetOk("project_id")
		if !ok || projectIDRaw == nil {
			d.Set("plan", prevPlan)
			return diag.Errorf("project_id is required but not set")
		}
		project_id := projectIDRaw.(string)

		locationRaw, ok := d.GetOk("location")
		if !ok || locationRaw == nil {
			d.Set("plan", prevPlan)
			return diag.Errorf("location is required but not set")
		}
		location := locationRaw.(string)

		planRaw, ok := d.GetOk("plan")
		if !ok || planRaw == nil {
			d.Set("plan", prevPlan)
			return diag.Errorf("[ERROR]plan is required but not set")
		}
		plan := planRaw.(string)

		versionRaw, ok := d.GetOk("version")
		if !ok || versionRaw == nil {
			d.Set("plan", prevPlan)
			return diag.Errorf("[ERROR]database_version is required but not set")
		}
		database_version := versionRaw.(string)

		dbaas_id := d.Id()

		software_id, err := apiClient.GetSoftwareId(project_id, location, "MySQL", database_version)
		if err != nil {
			d.Set("plan", prevPlan)
			return diag.FromErr(fmt.Errorf(" error while fetching software id: %v", err))
		}

		template_id, err := apiClient.GetTemplateId(project_id, location, plan, strconv.Itoa(software_id))
		if err != nil {
			d.Set("plan", prevPlan)
			return diag.FromErr(fmt.Errorf(" error while fetching template id: %v", err))
		}

		log.Printf("[INFO] prevPlan: %s, currPlan: %s", prevPlan.(string), currPlan.(string))

		statusRaw, ok := d.GetOk("status")
		if !ok || statusRaw == nil || statusRaw.(string) != "SUSPENDED" {
			d.Set("plan", prevPlan)
			return diag.Errorf("[ERROR]Node should be stopped before any upgradation ")
		}

		_, err = apiClient.UpgradeMySQLPlan(dbaas_id, template_id, project_id, location)
		if err != nil {
			d.Set("plan", prevPlan)
			return diag.FromErr(fmt.Errorf(" error while upgrading plan: %v", err))
		}
	}

	if d.HasChange("size") {
		oldSizeRaw, newSizeRaw := d.GetChange("size")
		oldSize := oldSizeRaw.(int)
		newSize := newSizeRaw.(int)

		_, err = apiClient.ExpandMySQLDBaaSDisk(mySqlDBaaSId, newSize, projectID, location)
		if err != nil {
			d.Set("size", oldSize)
			return diag.FromErr(fmt.Errorf("failed to expand disk: %v", err))
		}

		// On success, update to oldSize + newSize
		updatedSize := oldSize + newSize
		d.Set("size", updatedSize)
	}

	return ResourceReadMySqlDB(ctx, d, m)

}

func ResourceDeleteMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	mySqlDBaaSId := d.Id()

	currentState := d.Get("status").(string)

	if currentState == "CREATING" {
		return diag.Errorf("Cannot delete DBaaS: the database is still in 'CREATING' state")
	}

	_, err := apiClient.DeleteMySqlDBaaS(mySqlDBaaSId, d.Get("project_id").(string), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf(" error while deleting dbaas instance: %s", err))
	}

	d.SetId("")
	return diags
}
