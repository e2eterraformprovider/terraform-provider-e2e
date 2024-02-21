package sfs

import (
	// "context"

	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	//"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSfs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
         
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ValidateFunc: validateName,
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the Plan",
			},
			"vpc_id":{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "virtual private cloud id of sfs",
			},
			"disk_size":{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "size of disk to be created",
			},
			"project_id":{
				Type:        schema.TypeString,
				Required:    true,
				Description: "size of disk to be created",
			},
			"disk_iops":{
				Type:       schema.TypeInt,
				Required:   true,
				Description:  "input output per second",
			},
			"status":{
				Type:       schema.TypeString,
				Computed:   true,
				Description:  "status will be updated after creation",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location where node is to be launched",
				Default:     "Delhi",
			},
		},
		CreateContext: resourceCreateSfs,
		ReadContext:   resourceReadSfs,
		UpdateContext: resourceUpdateSfs,
		DeleteContext: resourceDeleteSfs,
		Exists:        resourceExistsSfs,
		Importer: &schema.ResourceImporter{
				State: schema.ImportStatePassthrough,
			},
}
}

func validateName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}


func resourceCreateSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.SFSClient)
	var diags diag.Diagnostics

	log.Printf("[INFO] NODE CREATE STARTS ")
	node := models.SfsCreate{
		Name:              d.Get("name").(string),
		Plan:              d.Get("plan").(string),
		Vpc_id:            d.Get("vpc_id").(int),
		Disk_size:         d.Get("disk_size").(int),
		Disk_iops:         d.Get("disk_iops").(int),
		
	}
	project_id:=d.Get("project_id").(string)
	res_Sfs, err := apiClient.NewSfs(&node,project_id)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] NODE CREATE | RESPONSE BODY | %+v", res_Sfs)
	if _, codeok := res_Sfs["code"]; !codeok {
		return diag.Errorf(res_Sfs["message"].(string))
	}

	data := res_Sfs["data"].(map[string]interface{})
	if data["is_credit_sufficient"] == false {
		return diag.Errorf(res_Sfs["message"].(string))
	}
	log.Printf("[INFO] sfs creation | before setting fields")
	sfsId := data["id"].(float64)
	sfsId = math.Round(sfsId)
	d.SetId(strconv.Itoa(int(math.Round(sfsId))))
	d.Set("status",data["status"].(string))
	return diags
}

func resourceReadSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.SFSClient)
	var diags diag.Diagnostics
	log.Printf("[info] inside node Resource read")
	nodeId := d.Id()
	project_id:=d.Get("project_id").(string)
	log.Printf("*************************====project_id type %T, value : %s \n", project_id, project_id)

	node, err := apiClient.GetSfs(nodeId,project_id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.Errorf("error finding Item with ID %s", nodeId)

		}
	}
	log.Printf("[info] node Resource read | before setting data")
	data := node["data"].(map[string]interface{})
	d.Set("name", data["name"].(string))
	d.Set("plan", data["efs_plan_name"].(string))
	d.Set("status", data["status"].(string))
	d.Set("vpc_id",data["vpc_id"].(int))
	d.Set("disk_size",data["efs_sisk_size"].(int))
	d.Set("disk_iops",data["efs_disk_iops"].(int))
	log.Printf("[info] node Resource read | after setting data")
	if d.Get("status").(string) == "Available" {
		d.Set("status", "power_on")
	}
	// if d.Get("status").(string) == "Powered off" {
	// 	d.Set("power_status", "power_off")
	// }
	return diags

}

func resourceUpdateSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// apiClient := m.(*client.SFSClient)

	

	return resourceReadSfs(ctx, d, m)

}
func resourceDeleteSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceExistsSfs(d *schema.ResourceData, m interface{}) (bool, error) {
	
	return true, nil
}

