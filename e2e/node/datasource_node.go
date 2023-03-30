package node

import (
	"context"
	// "fmt"
	"log"
	// "math"
	// "regexp"

	// "strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	// "github.com/devteametwoe/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNode() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group",
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the Plan",
			},
			"backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tells you the state of your backups",
				Default:     false,
			},

			"image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the image you have selected format :- ( os-version )",
			},
			"default_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tells us the state of default public ip",
				Default:     false,
			},
			"disable_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "can disable password as per requirement",
				Default:     false,
			},
			"enable_bitninja": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "enable bitnija as per requirement",
				Default:     false,
			},
			"is_ipv6_availed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
				Default:     false,
			},
			"is_saved_image": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "used when Creating node from a saved image",
				Default:     false,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location where node is to be launched",
				Default:     "Delhi",
			},
			"reserve_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reserve ip as per  requirement",
				Default:     "",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vpc id as per requirement",
				Default:     "Used when you need to attach a particular VPC. ",
			},
			"ngc_container_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Should be specified when launching GPU Cloud Wizard.",
				Default:     nil,
			},
			"saved_image_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "template id  is required when you save the node from saved images.Give the template id of the saved image. Required when is_saved_image field is true",
				Default:     nil,
			},
			"security_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specify the security group. Checkout security_groups datasource listing security groups",
			},
			"ssh_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the ssh keys if required. Checkout ssh_keys datasource for listing ssh keys",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the node",
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "memory of the node",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the node",
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disc info of the node",
			},
			"price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "price details of the node",
			},
			"public_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public ip address alloted to node",
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private ip address alloted to node if any",
			},
			"is_monitored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_bitninja_license_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Can check if the bitninja license is active or not",
			},
			"power_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "power_on",
				Description: "power_on to start the node and power_off to power off the node",
			},
			"lock_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Node is locked when set true .Can specify wheather to lock the node or not",
			},
			"reboot_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When set true node will be rebooted. Node should be in running state to perform rebooting.Alaways check the field. If you have an active disk-intensive process such as database, backups running, then a rebooting may lead to data corruption and data loss (best option is to reboot the machine from within Operating System). ",
			},
			"reinstall_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "for reinstalling the node. Node should be in running state to perform this action. Always check this field as it will delete all your data permenantly when set true.",
			},
			"save_image": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "For saving image of the node. The node should be in power_off state to perform this action ",
			},
			"save_image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify the name of the image to be saved. this field is required when save_image field is true. The name should be unique in the image list. Checkout images datasource to list them images",
			},
		},

		ReadContext: dataSourceReadNode,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
func dataSourceReadNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[INFO] inside node data source read")
	nodeId := d.Id()

	node, err := apiClient.GetNode(nodeId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.Errorf("error finding Item with ID %s", nodeId)

		}
	}
	data := node["data"].(map[string]interface{})

	d.Set("name", data["name"].(string))
	d.Set("label", data["label"].(string))
	d.Set("plan", data["plan"].(string))
	d.Set("backup", data["backup"].(bool))
	d.Set("is_active", data["is_active"].(bool))
	d.Set("created_at", data["created_at"].(string))
	d.Set("memory", data["memory"].(string))
	d.Set("status", data["status"].(string))
	d.Set("disk", data["disk"].(string))
	d.Set("price", data["price"].(string))
	d.Set("lock_node", data["is_locked"].(bool))
	d.Set("public_ip_address", data["public_ip_address"].(string))
	d.Set("private_ip_address", data["private_ip_address"].(string))
	d.Set("is_monitored", data["is_monitored"].(bool))
	d.Set("is_bitninja_license_active", data["is_bitninja_license_active"].(bool))
	if d.Get("status").(string) == "Running" {
		d.Set("power_status", "power_on")
	}
	if d.Get("status").(string) == "Powered off" {
		d.Set("power_status", "power_off")
	}

	return diags

}
