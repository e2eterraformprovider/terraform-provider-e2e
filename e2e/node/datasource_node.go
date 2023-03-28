package node

import (
	"context"
	// "fmt"
	"log"
	// "math"
	// "regexp"

	// "strconv"
	"strings"

	"github.com/devteametwoe/terraform-provider-e2e/client"
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
				Description: "The name of the Plan",
			},
			"os": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OS and its version  format : <OS>-<version>",
			},
			"backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tells you the state of your backups",
				Default:     false,
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the image you have selected",
				Default:     "CentOS-7.5-Distro",
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
				Description: "",
				Default:     false,
			},
			"enable_bitninja": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
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
				Description: "",
				Default:     false,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				Default:     "ncr",
			},
			"reserve_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				Default:     "",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				Default:     "",
			},
			"ngc_container_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "",
				Default:     nil,
			},
			"saved_image_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "",
				Default:     nil,
			},
			"security_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "",
				Default:     150,
			},
			"ssh_keys": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"price": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"power_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "power_on",
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_monitored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_bitninja_license_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lock_node": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "false",
			},
			"reboot_node": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "false",
			},
			"reinstall_node": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "false",
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
